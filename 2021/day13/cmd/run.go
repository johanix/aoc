/*
 *
 */
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		P, Folds = ReadInput(datafile)
		fmt.Printf("We have a paper with %d dots and %d folds\n",
			len(P.Dots), len(Folds))
		fmt.Printf("Folds: %v\n", Folds)
		P.Extents()
		for _, f := range Folds {
		    P = Folder(P, f)
		    P.Extents()
		}
		P.Print()
	},
}

var Folds = []Fold{}
var Dots = map[Dot]bool{}

var P Paper

type Paper struct {
     Dots map[Dot]bool
     Minx int
     Maxx int
     Miny int
     Maxy int
}

func (p *Paper) Extents() {
     fmt.Printf("Paper extents: [%d,%d] to [%d,%d] with a total of %d dots\n",
     		       p.Minx, p.Miny, p.Maxx, p.Maxy, len(p.Dots))
}

func (p *Paper) Print() {
     fmt.Printf("Printing paper with: ")
     p.Extents()
     var matrix = make([][]bool, p.Maxy - p.Miny + 1)
     for y := p.Miny; y <= p.Maxy ; y++ {
     	 matrix[y] = make([]bool, p.Maxx - p.Minx + 1)
     }
     for d, _ := range p.Dots {
     	 // fmt.Printf("Dot is on [%d,%d]\n", d.X, d.Y)
     	 matrix[d.Y][d.X] = true
     }
     for _, line := range matrix {
     	 for _, e := range line {
	     if e {
	     	fmt.Printf("#")
	     } else {
	       	fmt.Printf(".")
	     }
	 }
	 fmt.Println()
     }
}

type Dot struct {
     X	 int
     Y	 int
}

type Fold struct {
     Axis 	 string
     Value	 int
}

func Folder(p Paper, fold Fold) Paper {
     var newdots = map[Dot]bool{}
     maxx := p.Maxx
     maxy := p.Maxy
     var x, y int
     
     if fold.Axis == "x" {
     	for d, _ := range p.Dots {
	    x = d.X
	    if x > fold.Value {
	       x = 2 * fold.Value - x
	    }
	    newdots[Dot{ X: x, Y: d.Y}] = true
	}
	maxx = fold.Value - 1
     } else if fold.Axis == "y" {
     	for d, _ := range p.Dots {
	    y = d.Y
	    if y > fold.Value {
	       y = 2 * fold.Value - y
	    }
	    newdots[Dot{ X: d.X, Y: y }] = true
	}
	maxy = fold.Value - 1
     } else {
       log.Fatalf("Huga\n")
     }
     return Paper{ Dots: newdots, Minx: p.Minx, Maxx: maxx, Miny: p.Miny, Maxy: maxy }
}

func ReadInput(filename string) (Paper, []Fold) {
	var dots = map[Dot]bool{}
	var folds = []Fold{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: failed to open %s", filename)
	}

	var minx = 200
	var miny = 200
	var maxx, maxy int

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ",") {
		   coords := strings.Split(line, ",")
		   x, _ := strconv.Atoi(coords[0])
		   y, _ := strconv.Atoi(coords[1])
		   if x < minx {
		      minx = x
		   }
		   if x > maxx {
		      maxx = x
		   }
		   if y < miny {
		      miny = y
		   }
		   if y > maxy {
		      maxy = y
		   }
		   dots[Dot{X: x, Y: y}] = true
		   continue
		}
		if line == "" {
		   continue
		}
		if strings.Contains(line, "fold") {
		   details := strings.Split(strings.TrimLeft(line, "fold along "), "=")
		   value, _ := strconv.Atoi(details[1])
		   folds = append(folds, Fold{ Axis: details[0], Value: value})
		}
	}
	file.Close()
	p := Paper{ Dots: dots, Minx: minx, Miny: miny, Maxx: maxx, Maxy: maxy }
	
	return p, folds
}
