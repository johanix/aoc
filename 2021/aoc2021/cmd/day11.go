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

var day11aCmd = &cobra.Command{
	Use:   "day11a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		RunDay11(datafile)
	},
}

var day11bCmd = &cobra.Command{
	Use:   "day11b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		RunDay11(datafile)
	},
}

func RunDay11(datafile string) {
	if datafile == "" {
		datafile = "data/day-9.dat"
	}

	TheArea = ReadDay11(datafile)
	fmt.Printf("We have %d lines. minx: %d, miny: %d maxx: %d, maxy: %d\n", len(TheArea.Hlines), TheArea.Minx, TheArea.Miny, TheArea.Maxx, TheArea.Maxy)

	maxr := ReadInt("NUmber of rounds", 2)
	TheArea.Print()
	totflash := 0
	r := 0
	for r = 1; r <= maxr; r++ {
		TheArea.Increment()
		f := TheArea.Flash()
		if f == 100 {
		   fmt.Printf("********* All octi flashed during round %d\n", r)
		}
		totflash += f
		fmt.Printf("Round %d: %d flashes\n", r, f)
	}
	fmt.Printf("Totflash %d after %d rounds\n", totflash, r-1)
}

func (pt *Coord) OctoVal() int {
	foo, _ := strconv.Atoi(string(Hlines[pt.Y][pt.X]))
	return foo
}

func (a *Area) Energy(x, y int) int {
	// foo, _ := strconv.Atoi(string(a.Hlines[y][x]))
	foo := a.Hlines[y][x]
	return foo
}

func (a *Area) Print() {
	fmt.Printf("-----------\n")
	for y := a.Miny; y <= a.Maxy; y++ {
		for x := a.Minx; x <= a.Maxx; x++ {
			fmt.Printf("%d", a.Energy(x, y))
		}
		fmt.Println()
	}
}

func (a *Area) Increment() {
	for y := a.Miny; y <= a.Maxy; y++ {
		for x := a.Minx; x <= a.Maxx; x++ {
			// v := a.Energy(x,y)
			// fmt.Printf("Old val: %d new val: %v\n", v,byte(v + 1))
			a.Hlines[y][x] = a.Energy(x, y) + 1
		}
	}
}

func (a *Area) Flash() int {
	totflashes := 0
	a.HasFlashed = map[Coord]bool{} // clear all flash notations
	for {
		a.NewFlash = false
		for y := a.Miny; y <= a.Maxy; y++ {
			for x := a.Minx; x <= a.Maxx; x++ {
				if a.Energy(x, y) > 9 {
					pt := Coord{X: x, Y: y}
					if !a.HasFlashed[pt] {
						a.HasFlashed[pt] = true
						pt.IncNeighbors()
						a.NewFlash = true
						totflashes++
					}
				}
			}
		}
		if !a.NewFlash {
		   break
		}
	}
	checkflashes :=0
	for y := a.Miny; y <= a.Maxy; y++ {
		for x := a.Minx; x <= a.Maxx; x++ {
			if a.Energy(x, y) > 9 {
				a.Hlines[y][x] = 0
				checkflashes++
			}
		}
	}
	//fmt.Printf("Flash: totflashes: %d checkflashes: %d\n",
	//		   totflashes, checkflashes)
	return totflashes
}

func (pt *Coord) OctoNeighbors() []Coord {
	x := pt.X
	y := pt.Y
	res := []Coord{}
	tmp := []Coord{
		Coord{X: x - 1, Y: y - 1},
		Coord{X: x - 1, Y: y},
		Coord{X: x - 1, Y: y + 1},
		Coord{X: x, Y: y - 1},
		Coord{X: x, Y: y + 1},
		Coord{X: x + 1, Y: y - 1},
		Coord{X: x + 1, Y: y},
		Coord{X: x + 1, Y: y + 1},
	}
	for _, p := range tmp {
		if p.X < TheArea.Minx || p.X > TheArea.Maxx {
			continue
		}
		if p.Y < TheArea.Miny || p.Y > TheArea.Maxy {
			continue
		}
		res = append(res, p)
	}
	return res
}

func (pt *Coord) IncNeighbors() {
	neighbors := pt.OctoNeighbors()
	for _, n := range neighbors {
		TheArea.Hlines[n.Y][n.X]++
	}
}

func init() {
	rootCmd.AddCommand(day11aCmd, day11bCmd)
}

type IntLine []int

var TheArea = Area{}

type Area struct {
	Hlines     []IntLine
	TotFlashes int
	Maxx       int
	Maxy       int
	Minx       int
	Miny       int
	HasFlashed map[Coord]bool
	NewFlash   bool
}

func ReadDay11(filename string) Area {
	var line string
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
	return Area{
		Hlines:     hlines,
		Maxx:       len(hlines[0]) - 1,
		Maxy:       len(hlines) - 1,
		HasFlashed: map[Coord]bool{},
	}
}
