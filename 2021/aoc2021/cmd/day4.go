/*
 *
 */
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

// day4Cmd represents the day4 command
var day4aCmd = &cobra.Command{
	Use:   "day4a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		var worstboard, bestboard Board
		var bestscore, bestdraws, worstscore, worstdraws int
		var err error

		if datafile == "" {
			datafile = "data/day-4.dat"
		}

		file, err := os.Open(datafile)
		if err != nil {
			log.Fatalf("Error: failed to open %s", datafile)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		draws := ReadDraws(scanner)
		bestboard, err = ReadBoard(scanner, verbose)
		bestdraws, bestscore = bestboard.Score(draws)
		worstboard = bestboard
		worstdraws = bestdraws
		worstscore = bestscore

		for {
		    nextboard, err := ReadBoard(scanner, verbose)
		    if err != nil {
		       	   log.Printf("Error from Readboard: %v\n", err)
			   break
		    }
		    nd, ns := nextboard.Score(draws)
		    if nd < bestdraws  {
		       bestdraws = nd
		       bestscore = ns
		       bestboard = nextboard
		    }
		    if nd > worstdraws  {
		       worstdraws = nd
		       worstscore = ns
		       worstboard = nextboard
		    }
		}

		fmt.Printf("Best board had score: %d after %d draws\n", bestscore, bestdraws)
		bestboard.Print()
		fmt.Printf("Worst board had score: %d after %d draws\n", worstscore, worstdraws)
		worstboard.Print()

		file.Close()
	},
}

func init() {
	rootCmd.AddCommand(day4aCmd)

	// day4Cmd.PersistentFlags().String("foo", "", "A help for foo")
	// day4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Board struct {
	Rows [][]int
}

func NewBoard() Board {
	b := Board{
		Rows: make([][]int, 5),
	}
	for i := 0; i < 5; i++ {
		b.Rows[i] = make([]int, 5)
	}
	return b
}

func (b *Board) Print() {
	var out []string
	for r := 0; r < 5; r++ {
		line := ""
		for c := 0; c < 5; c++ {
			line += fmt.Sprintf("%d|", b.Rows[c][r])
		}
		out = append(out, line[:len(line)-1])
	}
	fmt.Printf("%s\n", columnize.SimpleFormat(out))
}

func (b *Board) Score(draws []int) (int, int) { // numdraws, score

	mark := func(num int) (bool, int, int) {
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if b.Rows[c][r] == num {
					b.Rows[c][r] = -1
					return true, r, c
				}
			}
		}
		return false, -1, -1
	}

	check := func(r, c int) bool {
		rbingo := true
		cbingo := true
		for t := 0; t < 5; t++ {
			if b.Rows[c][t] != -1 {
				rbingo = false
			}
			if b.Rows[t][r] != -1 {
				cbingo = false
			}
		}
		return rbingo || cbingo
	}

	score := func(draw int) int {
		sum := 0
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if b.Rows[c][r] != -1 {
					sum += b.Rows[c][r]
				}
			}
		}
		return sum * draw
	}

	for count, num := range draws {
		hit, r, c := mark(num)
		if hit && check(r, c) {
		       s := score(num)
		       fmt.Printf("Board: score %d (on draw %d)\n",
				  s, count)
			return count, s
		}
	}
	fmt.Printf("Board never had bingo\n")
	b.Print()
	return -1, -1 // should not happen
}

func ReadDraws(scanner *bufio.Scanner) []int {
	var draws []int
	scanner.Scan() // read 1st line
	foo := strings.Split(scanner.Text(), ",")
	for _, n := range foo {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Printf("Error from Atoi: %v\n", err)
		}
		draws = append(draws, num)
	}
	return draws
}

func ReadBoard(scanner *bufio.Scanner, verbose bool) (Board, error) {
	var line string
	var b = NewBoard()
	var num int
	var err error
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			// fmt.Printf("Got newline as expected, all good\n")
		} else {
			fmt.Printf("No newline, huh? Line: '%s'\n", line)
		}
		for row := 0; row <= 4; row++ {
			scanner.Scan()
			line := scanner.Text()
			fmt.Printf("Line: '%s'\n", line)
			foo := strings.Split(strings.ReplaceAll(strings.TrimSpace(line), "  ", " "), " ")
			// fmt.Printf("foo: '%v'\n", foo)
			if len(foo) != 5 {
				var str []string
				for _, v := range foo {
					str = append(str, fmt.Sprintf("'%s' ", v))
				}
				fmt.Printf("Did not get 5 nums in a row when reading board. Got (len=%d): %s\n", len(foo), strings.Join(str, ","))

			}
			for col, n := range foo {
				num, err = strconv.Atoi(n)
				if err != nil {
					log.Printf("Error from Atoi: %v\n", err)
				}
				if verbose {
					fmt.Printf("Inserting '%d' in cell [%d,%d]\n",
						num, col, row)
				}
				b.Rows[col][row] = num
			}
		}
		b.Print()
		return b, nil
	}
	return b, fmt.Errorf("End of input")
}
