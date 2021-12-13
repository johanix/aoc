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
	// "strings"

	"github.com/spf13/cobra"
)

var day9aCmd = &cobra.Command{
	Use:   "day9a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		RunDay9(datafile)
	},
}

var day9bCmd = &cobra.Command{
	Use:   "day9b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		RunDay9(datafile)
	},
}

// type Coord struct {
//     X	   int
//     Y	   int
// }

// var Hlines []HLine
var Hlines []IntLine
var Basins = map[Coord]Basin{}

func RunDay9(datafile string) {
	if datafile == "" {
		datafile = "data/day-9.dat"
	}

	Hlines = ReadDay9(datafile)
	fmt.Printf("We have %d lines\n", len(Hlines))

	lows := map[Coord]int{}
	// basins := map[Coord]Basin{}
	maxy := len(Hlines) - 1
	maxx := 0
	h := 0
	var pt Coord
	for y, l := range Hlines {
		maxx = len(l) - 1
		for x, val := range l {
			lowpoint := true
			pt = Coord{X: x, Y: y}
			for _, n := range pt.Neighbors(maxx, maxy) {
				if Hlines[n.Y][n.X] <= val {
					lowpoint = false
					break
				}
			}
			if lowpoint {
				h, _ = strconv.Atoi(string(val))
				lows[Coord{X: x, Y: y}] = h
				Basins[Coord{X: x, Y: y}] = Basin{}
			}

		}
	}

	risksum := 0
	for _, v := range lows {
		risksum += v + 1
	}
	fmt.Printf("There are %d low points total, risk sum is %d\n", len(lows), risksum)
	fmt.Printf("Low point locations: %v\n", lows)
	for y, l := range Hlines {
		maxx = len(l) - 1
		for x, _ := range l {
			pt := Coord{X: x, Y: y}
			err := pt.FindBasin(maxx, maxy)
			if err != nil {
				fmt.Printf("Error from findbasin\n")
			}
		}
	}
	//		startx := ReadInt("FindBasin start x", 0)
	//		starty := ReadInt("FindBasin start y", 0)
	//		pt = Coord{ X: startx, Y: starty}
	//		fmt.Printf("Will look for lowest point in basin starting at %s\n", pt.Sprintf())
	//		pt.FindBasin(maxx, maxy)
	fmt.Printf("Part 2: Basins: there are %d basins with centers and sizes:\n", len(Basins))
	for pt, b := range Basins {
		fmt.Printf("%s size: %d\n", pt.Sprintf(), b.Size)
	}
}

var LowPointMap = make(map[Coord]LowPoint, 10000)

type LowPoint struct {
	LP    Coord
	Known bool
}

type Basin struct {
	X       int
	Y       int
	Members []Coord
	Size    int
}

func (pt *Coord) Height() int {
	// foo, _ := strconv.Atoi(string(Hlines[pt.Y][pt.X]))
	foo := Hlines[pt.Y][pt.X]
	return foo
}

func (pt *Coord) Sprintf() string {
	return fmt.Sprintf("[%d,%d (%d)]", pt.X, pt.Y, pt.Height())
}

func (pt *Coord) FindBasin(maxx, maxy int) error {
	myheight := pt.Height()
	if myheight == 9 {
		fmt.Printf("Point %s has height == 9 and is not part of any basin\n", pt.Sprintf())
		return nil
	}
	lowestheight := myheight
	lowp := *pt
	lowprospects := pt.Neighbors(maxx, maxy)
	seen := make(map[Coord]bool, 10000)
	seen[*pt] = true
	var pheight int
	fmt.Printf("%s Low prospects: %v\n", pt.Sprintf(), lowprospects)
	for {
		if len(lowprospects) == 0 {
			break
		}
		p := lowprospects[0]
		lowprospects = lowprospects[1:]
		fmt.Printf("Analysing prospect %s\n", p.Sprintf())
		if p.Height() > myheight {
			continue
		}
		// rest := lowprospects[1:]
		if LowPointMap[p].Known {
			lowp = LowPointMap[p].LP
			break
		}
		if p.Height() < lowestheight {
			lowp = p
			lowestheight = p.Height()
		}
		prospneighbors := p.Neighbors(maxx, maxy)
		fmt.Printf("%s Prospect %s neighbors: %v\n", pt.Sprintf(), p.Sprintf(), prospneighbors)
		for _, n := range prospneighbors {
			if seen[n] {
				continue // don't add a neighbor that is already analysed
			}
			if n.Height() <= pheight {
				lowprospects = append(lowprospects, n) // only add relevant neighbors
			}
		}
	}
	fmt.Printf("%s: Low point in basin is %s with height %d\n", pt.Sprintf(), lowp.Sprintf(), lowestheight)
	b := Basins[lowp]
	b.Size += 1
	b.Members = append(b.Members, *pt)
	Basins[lowp] = b
	lp := LowPoint{Known: true, LP: lowp}
	LowPointMap[*pt] = lp
	return nil
}

// func (pt *Coord) LowestNeighbor(

func (pt *Coord) PrintNeighbors(maxx, maxy int) {
	x := pt.X
	y := pt.Y
	var h int
	printline := func(y int) {
		if x > 0 {
			h, _ = strconv.Atoi(string(Hlines[y][x-1]))
			fmt.Printf("%d", h)
		} else {
			fmt.Printf("|")
		}
		h, _ = strconv.Atoi(string(Hlines[y][x]))
		fmt.Printf("%d", h)
		if x < maxx {
			h, _ = strconv.Atoi(string(Hlines[y][x+1]))
			fmt.Printf("%d", h)
		} else {
			fmt.Printf("|")
		}
		fmt.Println()
	}
	fmt.Printf("Low point: [%d,%d]:\n", x, y)
	if y > 0 {
		printline(y - 1)
	} else {
		fmt.Printf("---\n")
	}
	printline(y)
	if y < maxy {
		printline(y + 1)
	} else {
		fmt.Printf("---\n")
	}
	fmt.Println()
}

func (pt *Coord) Neighbors(maxx, maxy int) []Coord {
	res := []Coord{}
	if pt.X > 0 {
		res = append(res, Coord{X: pt.X - 1, Y: pt.Y})
	}
	if pt.X < maxx {
		res = append(res, Coord{X: pt.X + 1, Y: pt.Y})
	}
	if pt.Y > 0 {
		res = append(res, Coord{X: pt.X, Y: pt.Y - 1})
	}
	if pt.Y < maxy {
		res = append(res, Coord{X: pt.X, Y: pt.Y + 1})
	}
	return res
}

func init() {
	rootCmd.AddCommand(day9aCmd, day9bCmd)
}

type HLine []byte

func ReadDay9(filename string) []IntLine {
	var line string
	// var hlines = []HLine{}
	var hlines = []IntLine{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: failed to open %s", filename)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = scanner.Text()
		iline := []int{}
		tmp := []byte(line)
		for _, n := range tmp {
		    v, _ := strconv.Atoi(string(n))
		    iline = append(iline, v)
		}
		hlines = append(hlines, iline)
	}
	file.Close()
	return hlines
}
