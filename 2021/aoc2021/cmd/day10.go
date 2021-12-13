/*
 *
 */
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	// "strconv"
	// "strings"

	"github.com/spf13/cobra"
)

var day10aCmd = &cobra.Command{
	Use:   "day10a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		RunDay10(datafile)
	},
}

var day10bCmd = &cobra.Command{
	Use:   "day10b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		RunDay10(datafile)
	},
}

var Slines []SyntaxLine

func RunDay10(datafile string) {
	if datafile == "" {
		datafile = "data/day-10.dat"
	}

	Slines = ReadDay10(datafile)
	var Flines = []SyntaxLine{}
	fmt.Printf("We have %d lines\n", len(Slines))
	totscore := 0
	toterrors := 0
	for _, sl := range Slines {
	    correct, score := sl.Scan()
	    if !correct {
	       toterrors++
	       totscore += score
	    } else {
	       Flines = append(Flines, sl)
	    }
	}
	fmt.Printf("There were %d errors with a total score of %d\n", toterrors, totscore)
	fmt.Printf("There are %d lines remaining, with chunkstacks:\n", len(Flines))
	cscores := []int{}
	cscore := 0
	for _, l := range Flines {
	    cscore = l.Complete()
	    cscores = append(cscores, cscore)
	    fmt.Printf("%v. Completion score: %d\n", l.ChunkStack, cscore)
	}
	sort.Ints(cscores)
	middle := cscores[len(cscores) / 2]
	fmt.Printf("Total %d completion scores. Sorted middle score: %d\n", len(cscores), middle)
}

func init() {
	rootCmd.AddCommand(day10aCmd, day10bCmd)
}

type Token struct {
     Open  bool
     Value int
}

var Tokens = map[byte]Token{
		'(': Token{ Open: true, Value: 1 },
		')': Token{ Open: false, Value: 1 },
		'[': Token{ Open: true, Value: 2 },
		']': Token{ Open: false, Value: 2 },
		'{': Token{ Open: true, Value: 3 },
		'}': Token{ Open: false, Value: 3 },
		'<': Token{ Open: true, Value: 4 },
		'>': Token{ Open: false, Value: 4 },
}

func PrependInt(x []int, y int) []int {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

var TokenScore = map[byte]int{
			')': 3,
			']': 57,
			'}': 1197,
			'>': 25137,
}

func (sl *SyntaxLine) Scan() (bool, int) {
     for _, p := range []byte(sl.Raw) {
     	 if t, ok := Tokens[p]; ok {
	    if t.Open {
	       sl.ChunkStack = PrependInt(sl.ChunkStack, t.Value)
	    } else {
	       if sl.ChunkStack[0] == t.Value {
	       	  sl.ChunkStack = sl.ChunkStack[1:]
	       } else {
	       	 // fmt.Printf("Syntax error on pos %d ('%s') with ChunkStack: %v\n", pos, string(p), sl.ChunkStack)
		 return false, TokenScore[p]
	       }
	    }
	 } else {
	   fmt.Printf("%s is not a token\n", string(p))
	 }
     }
     return true, 0
}

func (sl *SyntaxLine) Complete() (int) {
     // fmt.Printf("Line is '%s'. Chunkstack is: %v\n", sl.Raw, sl.ChunkStack)
     score := 0
     for _, c := range sl.ChunkStack {
     	 score = score * 5 + c
     }
     return score
}

type SyntaxLine struct {
	Raw        string
	ChunkStack []int
}

func ReadDay10(filename string) []SyntaxLine {
	var line string
	var slines = []SyntaxLine{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: failed to open %s", filename)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = scanner.Text()
		slines = append(slines, SyntaxLine{Raw: line, ChunkStack: []int{} })
	}
	file.Close()
	return slines
}
