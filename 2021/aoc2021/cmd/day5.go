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

	"github.com/spf13/cobra"
)

var day5aCmd = &cobra.Command{
	Use:   "day5a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	     RunDay5(datafile, usediags)
	},
}

var day5bCmd = &cobra.Command{
	Use:   "day5b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
	     usediags = true
	     RunDay5(datafile, true) // force use of diagonals for part B.
	},
}

var usediags bool

func RunDay5(datafile string, usediags bool) {
		if datafile == "" {
			datafile = "data/day-5.dat"
		}

		lines := ReadDay5(datafile)
		fmt.Printf("We have %d lines\n", len(lines))
		// fmt.Printf("Lines:\n %v\n", lines)

		// compute intersections
		intersections := map[Coord]int{}
		for i, l1 := range lines {
			for j, l2 := range lines {
				if j == i {
					continue
				}
				var ip []Coord
				ip = l1.IntersectionPoints(l2)
				if len(ip) > 0 {
					for _, pt := range ip {
						intersections[pt]++
					}
				}
			}
		}

		fmt.Printf("There are %d intersection points total\n", len(intersections))
		if verbose {
			fmt.Printf("%v\n", intersections)
		}
}

func init() {
	rootCmd.AddCommand(day5aCmd, day5bCmd)
	day5aCmd.Flags().BoolVarP(&usediags, "diagonals", "D", false, "include diagonals")
}

type Line struct {
	Coords   []int
	LineType string
	Slope    int // +1 or -1
}

func NewLine(ci []int) Line {
	var lt string
	var tmp, slope int
	if ci[0] == ci[2] {
		lt = "vertical"
		if ci[3] < ci[1] {
			// fmt.Printf("NewLine: swapping Y: %d < %d\n", ci[3], ci[1])
			tmp = ci[1]
			ci[1] = ci[3]
			ci[3] = tmp
		}
	} else if ci[1] == ci[3] {
		lt = "horizontal"
		if ci[2] < ci[0] {
			// fmt.Printf("NewLine: swapping X: %d < %d\n", ci[2], ci[0])
			tmp = ci[0]
			ci[0] = ci[2]
			ci[2] = tmp
		}
	} else {
		lt = "diagonal"
		if ci[2] < ci[0] {
			// fmt.Printf("NewLine: diagonal swapping X: %d < %d\n", ci[2], ci[0])
			tmp = ci[0]
			ci[0] = ci[2]
			ci[2] = tmp
			tmp = ci[1]
			ci[1] = ci[3]
			ci[3] = tmp
		}
		if ci[3] > ci[1] {
			slope = 1
		} else {
			slope = -1
		}
	}
	l := Line{
		Coords:   ci,
		LineType: lt,
		Slope:    slope,
	}
	return l
}

type Coord struct {
	X int
	Y int
}

// if l = horizontal, then Y must match and X1 <= X <= X2
// if l = vertical, then X must match and Y1 <= Y <= Y2

func (l *Line) Intersects(pt Coord) bool {
	if l.LineType == "horizontal" {
		if pt.Y == l.Coords[1] && l.Coords[0] <= pt.X && pt.X <= l.Coords[2] {
			return true
		}
		return false
	} else if l.LineType == "vertical" {
		if pt.X == l.Coords[0] && l.Coords[1] <= pt.Y && pt.Y <= l.Coords[3] {
			return true
		}
		return false
	} else if l.LineType == "diagonal" && usediags {
		ly := l.Coords[1]
		for lx := l.Coords[0]; lx <= l.Coords[2]; lx++ {
			if pt.Y == ly && pt.X == lx {
				return true
			}
			ly = ly + l.Slope
		}
		return false
	}
	return false
}

func (l *Line) IntersectionPoints(l2 Line) []Coord {
	var ips []Coord
	l2t := l2.LineType
	if l2t == "error" || l.LineType == "error" {
		return ips
	}
	if l2t == "horizontal" {
		for x := l2.Coords[0]; x <= l2.Coords[2]; x++ {
			c := Coord{X: x, Y: l2.Coords[1]}
			if l.Intersects(c) {
				ips = append(ips, c)
			}
		}
	} else if l2t == "vertical" {
		for y := l2.Coords[1]; y <= l2.Coords[3]; y++ {
			c := Coord{X: l2.Coords[0], Y: y}
			if l.Intersects(c) {
				ips = append(ips, c)
			}
		}
	} else if l2t == "diagonal" && usediags {
		y := l2.Coords[1]
		for x := l2.Coords[0]; x <= l2.Coords[2]; x++ {
			c := Coord{X: x, Y: y}
			if l.Intersects(c) {
				ips = append(ips, c)
			}
			y = y + l2.Slope
		}
	}
	return ips
}

func (l *Line) Print() {
	fmt.Printf("%d,%d -> %d,%d (%s)\n", l.Coords[0], l.Coords[1],
		l.Coords[2], l.Coords[3], l.LineType)
}

func (l *Line) Sprint() string {
	var res string
	if l.LineType == "diagonal" {
		y := l.Coords[1]
		for x := l.Coords[0]; x <= l.Coords[2]; x++ {
			res += fmt.Sprintf(" [%d,%d]", x, y)
			y = y + l.Slope
		}
		return res
	}
	return fmt.Sprintf("%d,%d -> %d,%d", l.Coords[0], l.Coords[1],
		l.Coords[2], l.Coords[3])
}

func ReadDay5(filename string) []Line {
	var line string
	var lines = []Line{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: failed to open %s", filename)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = strings.ReplaceAll(scanner.Text(), " -> ", ",")
		coords := strings.Split(line, ",")
		var ci = make([]int, 4)
		for pos, n := range coords {
			num, err := strconv.Atoi(n)
			if err != nil {
				log.Fatalf("Error from atoi: %v\n", err)
			}
			ci[pos] = num
		}
		lines = append(lines, NewLine(ci))
	}
	file.Close()
	return lines
}
