/*
 *
 */
package cmd

import (
	"bufio"
	"fmt"
	"log"
	// "math"
	"os"
	"strconv"
	"strings"

	//"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var day7aCmd = &cobra.Command{
	Use:   "day7a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Runner7()
	},
}

var day7bCmd = &cobra.Command{
	Use:   "day7b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Runner7()
	},
}

func Runner7() {
	if datafile == "" {
		datafile = "data/day-7.dat"
	}

	crabs, count, minc, maxc := ReadCrabs(datafile)
	fmt.Printf("We have %d crabs with min pos %d and max pos %d\n", count, minc, maxc)
	FuelCosts = make([]int, maxc+1)
	StepsToFuel(maxc)

	var minpos, minpos2 int
	minfuel := FuelNeed(minc, crabs)
	minfuel2 := FuelNeed2(minc, crabs)
	for p := minc ; p <= maxc ; p++ {
	    f1 := FuelNeed(p, crabs)
	    if f1 < minfuel {
	       minfuel = f1
	       minpos = p
	    }
	    f2 := FuelNeed2(p, crabs)
	    if f2 < minfuel2 {
	       minfuel2 = f2
	       minpos2 = p
	    }
	}
	fmt.Printf("minfuel: method 1: %d (pos %d) method 2: %d (pos %d)\n", minfuel, minpos, minfuel2, minpos2)
}

func FuelNeed(pos int, crabs []int) int {
     fuel := 0
     for _, c := range crabs {
     	 if pos > c {
	    fuel += pos - c
	 } else {
	    fuel += c - pos
	 }
     }
     return fuel
}

var FuelCosts []int

func StepsToFuel(n int)  {
     for p:= 1 ; p <= n ; p++ {
     	 FuelCosts[p] = p + FuelCosts[p-1]
     }
}

func FuelNeed2(pos int, crabs []int) int {
     fuel := 0
     steps := 0
     for _, c := range crabs {
     	 if pos > c {
	    steps = pos - c
	 } else {
	    steps = c - pos
	 }
	 fuel += FuelCosts[steps]
     }
     return fuel
}

func init() {
	rootCmd.AddCommand(day7aCmd, day7bCmd)
}

func ReadCrabs(datafile string) ([]int, int, int, int) {
	var crabs []int

	file, err := os.Open(datafile)
	if err != nil {
		log.Fatalf("Error: failed to open %s", datafile)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan() // read 1st line
	maxc := 0
	minc := 0
	foo := strings.Split(scanner.Text(), ",")
	sum := 0
	var n string
	for _, n = range foo {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Printf("Error from Atoi: %v\n", err)
		}
		crabs = append(crabs, num)
		sum += num
		if num > maxc {
		   maxc = num
		}
		if num < minc {
		   minc = num
		}
	}
	file.Close()
	l := len(crabs)
	return crabs, l, minc, maxc
}
