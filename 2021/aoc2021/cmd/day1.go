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

// day1Cmd represents the day1 command
var day1aCmd = &cobra.Command{
	Use:   "day1a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	     	data := ReadDay1("day-1a.dat")
		
		var previous, increased, decreased int
		for line, n := range data {
			if line == 0 {
				previous = n
				continue
			}
			if n > previous {
				increased++
			} else if n < previous {
				decreased++
			}
			previous = n
		}
		fmt.Printf("Increased: %d decreased: %d lines: %d\n", increased, decreased, len(data))
	},
}

var day1bCmd = &cobra.Command{
	Use:   "day1b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	     	data := ReadDay1("day-1a.dat")
		
		var previous, increased, decreased int
		for line, n := range data[:len(data)-2] {
		    	this := n + data[line+1] + data[line+2]
			if line == 0 {
				previous = this
				continue
			}
			if this > previous {
				increased++
			} else if this < previous {
				decreased++
			}
			previous = this
		}
		fmt.Printf("Increased: %d decreased: %d lines: %d\n", increased, decreased, len(data))
	},
}

var day2aCmd = &cobra.Command{
	Use:   "day2a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	     	data := ReadDay2("day-2a.dat")
		
		var depth, dist int
		for _, str := range data {
		        parts := strings.Split(str, " ")
			num, _ := strconv.Atoi(parts[1])
		    	switch foo := str[:2]; foo {
			case "up":
			     depth -= num
			case "do":
			     depth += num
			case "fo":
			     dist += num
			}
		}
		fmt.Printf("Distance: %d depth: %d product: %d\n",
				      dist, depth, dist*depth)
	},
}

var day2bCmd = &cobra.Command{
	Use:   "day2b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	     	data := ReadDay2("day-2a.dat")
		
		var depth, dist, aim int
		for _, str := range data {
		        parts := strings.Split(str, " ")
			num, _ := strconv.Atoi(parts[1])
		    	switch foo := str[:2]; foo {
			case "up":
			     aim -= num
			case "do":
			     aim += num
			case "fo":
			     dist += num
			     depth += aim*num
			}
		}
		fmt.Printf("Distance: %d depth: %d product: %d\n",
				      dist, depth, dist*depth)
	},
}

func init() {
	rootCmd.AddCommand(day1aCmd, day1bCmd, day2aCmd, day2bCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// day1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// day1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ReadDay1(filename string) []int {
		var data []int
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error: failed to open %s", filename)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var num int
		for scanner.Scan() {
			num, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatalf("Error from atoi: %v\n", err)
			}
			data = append(data, num)
		}

		file.Close()
		return data
}

func ReadDay2(filename string) []string {
		var data []string
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error: failed to open %s", filename)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		// var num int
		for scanner.Scan() {
			// num, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatalf("Error from atoi: %v\n", err)
			}
			data = append(data, scanner.Text())
		}

		file.Close()
		return data
}
