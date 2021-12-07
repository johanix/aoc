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

	//"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var day6aCmd = &cobra.Command{
	Use:   "day6a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Runner()
	},
}

var day6bCmd = &cobra.Command{
	Use:   "day6b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Runner()
	},
}

func Runner() {
	if datafile == "" {
		datafile = "data/day-6.dat"
	}

	fish := ReadFish(datafile)
	sfish := SortFish(fish)
	days := ReadDays()
	PrintFishNG("Initial state", 0, sfish)
	for day := 1; day <= days; day++ {
		sfish = RunDayNG(sfish)
		if days < 30 {
			PrintFishNG(fmt.Sprintf("After %d days", day), day, sfish)
		}
	}
	PrintFishNG(fmt.Sprintf("After %d days", days), days, sfish)
}

func SortFish(fish []int) []int {
	sfish := make([]int, 9)

	for _, f := range fish {
		sfish[f]++
	}
	return sfish
}

func RunDay(fish []int) []int {
	for pos, f := range fish {
		switch f {
		case 0:
			fish[pos] = 6
			fish = append(fish, 8)

		case 1, 2, 3, 4, 5, 6, 7, 8:
			fish[pos]--
		}
	}
	return fish
}

func RunDayNG(sortedfish []int) []int {
	sfish := make([]int, 9)

	for phase, count := range sortedfish {
		switch phase {
		case 0:
			sfish[6] = count
			sfish[8] = count

		case 1, 2, 3, 4, 5, 6, 8:
			sfish[phase-1] = count
		case 7:
			sfish[6] += count
		}
	}
	return sfish
}

func ReadDays() int {
	foo := 18
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%d]: ", "Days to iterate", foo)
	text, _ := reader.ReadString('\n')
	if text == "\n" {
		fmt.Printf("[empty response, using default]\n")
		return foo
	} else {
		days, err := strconv.Atoi(strings.TrimSuffix(text, "\n"))
		if err != nil {
			log.Fatalf("Error from atoi: %v\n", err)
		}
		return days
	}
}

func PrintFish(str string, days int, fish []int) {
	l := len(fish)
	if l == 0 {
		return
	}
	fmt.Printf("OG: %s [%d total] %d", str, l, fish[0])
	for _, f := range fish[1:] {
		fmt.Printf(",%d", f)
	}
	fmt.Printf(" [total after %d days is %d]\n", days, l)
}

func PrintFishNG(str string, days int, fish []int) {
	sum := func(fish []int) int {
		total := 0
		for _, count := range fish {
			total += count
		}
		return total
	}

	l := sum(fish)
	if l == 0 {
		return
	}
	fmt.Printf("NG: %s [%d total]", str, l)
	for phase, count := range fish {
		fmt.Printf("[phase %d:%d],", phase, count)
	}
	fmt.Printf(" [total after %d days is %d]\n", days, l)
}

func init() {
	rootCmd.AddCommand(day6aCmd)

	// day4Cmd.PersistentFlags().String("foo", "", "A help for foo")
	// day4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ReadFish(datafile string) []int {
	var fish []int

	file, err := os.Open(datafile)
	if err != nil {
		log.Fatalf("Error: failed to open %s", datafile)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan() // read 1st line
	foo := strings.Split(scanner.Text(), ",")
	for _, n := range foo {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Printf("Error from Atoi: %v\n", err)
		}
		fish = append(fish, num)
	}
	file.Close()
	return fish
}
