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

var day1aCmd = &cobra.Command{
	Use: "day1a",
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
		fmt.Printf("Increased: %d decreased: %d lines: %d\n",
			increased, decreased, len(data))
	},
}

var day1bCmd = &cobra.Command{
	Use: "day1b",
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
		fmt.Printf("Increased: %d decreased: %d lines: %d\n",
			increased, decreased, len(data))
	},
}

var day2aCmd = &cobra.Command{
	Use: "day2a",
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
	Use: "day2b",
	Run: func(cmd *cobra.Command, args []string) {

		data := ReadDay2("day-2a.dat")

		var depth, dist, aim int
		for _, str := range data {
			parts := strings.Split(str, " ")
			num, _ := strconv.Atoi(parts[1])
			switch str[:2] {
			case "up":
				aim -= num
			case "do":
				aim += num
			case "fo":
				dist += num
				depth += aim * num
			}
		}
		fmt.Printf("Distance: %d depth: %d product: %d\n",
			dist, depth, dist*depth)
	},
}

var day3aCmd = &cobra.Command{
	Use: "day3a",
	Run: func(cmd *cobra.Command, args []string) {

		data := ReadDay2("day-3.dat")

		var pcounters = make(map[int]int, 1000)
		var ncounters = make(map[int]int, 1000)
		for _, str := range data {
			for pos, c := range strings.Split(str, "") {
				if c == "1" {
					pcounters[pos]++
				} else {
					ncounters[pos]++
				}
			}
		}
		var gammastr []string
		gamma := 0

		for i := 0; i < 12; i++ {
			fmt.Printf("Char %d: p=%d n=%d. Gamma+: %d\n", i, pcounters[i],
				ncounters[i], 1<<(11-i))
			if pcounters[i] > ncounters[i] {
				gammastr = append(gammastr, "1")
				gamma = 1<<(11-i) + gamma
			} else {
				gammastr = append(gammastr, "0")
			}
		}
		eps := 1<<12 - 1 - gamma
		fmt.Printf("Gamma str: '%s' num: %d eps: %d result: %d\n",
			gammastr, gamma, eps, gamma*eps)
	},
}

var day3bCmd = &cobra.Command{
	Use: "day3b",
	Run: func(cmd *cobra.Command, args []string) {

		if datafile == "" {
			datafile = "day-3.dat"
		}
		oxydata := ReadDay2(datafile)
		co2data := ReadDay2(datafile)

		keylen := len(oxydata[0])
		fmt.Printf("Key length: %d\n", keylen)

		fmt.Printf("*** Generating oxygen data\n")
		oxygen := Remover2(oxydata, "1", "more", 0, keylen)
		fmt.Printf("*** Generating CO2 data\n")
		co2scrub := Remover2(co2data, "0", "less", 0, keylen)

		oxy := BinToInt(oxygen)
		co2 := BinToInt(co2scrub)

		fmt.Printf("Oxygen: %s CO2 scrub: %s. Product: %d\n",
			oxygen, co2scrub, oxy*co2)
	},
}

func BinToInt(str string) int {
	mexp := len(str) - 1 // should be 11
	sum := 0
	for i, v := range str {
		if v == '1' {
			sum = 1<<(mexp-i) + sum
		}
	}
	return sum
}

// map version, would have been great if using original counters rather than
// recomputing for each step
func Remover(dmap map[string]bool, keeper rune, keylen int, pc, nc map[int]int) string {
	var remaining = len(dmap)

	for n := 0; n < keylen; n++ {
		// keep := "0"
//		if pc[n] > nc[n] {
//			keep = "1"
//		} else if pc[n] < nc[n] {
//			keep = "0"
//		} else if pc[n] == nc[n] {
//			keep = string(keeper)
//		}
		for k, _ := range dmap {
			if pc[n] > nc[n] { // keep data with pos(n) = 1
				if k[n] == '1' {
					continue
				} else {
					delete(dmap, k)
					remaining--
				}
			} else if pc[n] < nc[n] {
				if k[n] == '0' {
					continue
				} else {
					delete(dmap, k)
					remaining--
				}
			} else {
				if string(k[n]) == string(keeper) {

					continue
				} else {
					delete(dmap, k)
					remaining--
				}
			}
			if remaining == 1 {
				var thekey string
				for k, _ := range dmap {
					thekey = k
				}
				fmt.Printf("Only one item left (%s). Returning.\n", thekey)
				return thekey
			}
		}
	}
	return "" // bad if this happens
}

func Remover2(data []string, keeper, want string, pos, keylen int) string {

	counter := func(curpos int, want string) (map[int]int, map[int]int) {
		onecount := map[int]int{}
		zerocount := map[int]int{}
		for _, str := range data {
			for pos, c := range strings.Split(str, "") {
				if c == "1" {
					onecount[pos]++
				} else {
					zerocount[pos]++
				}
			}
		}
		if want == "more" {
		   return onecount, zerocount
		}
		return zerocount, onecount
	}

	cleaner := func(d []string) []string {
		totlen := len(d)
		dirty := true
		for dirty {
			dirty = false
			for k := 0 ; k < totlen ; k++ {
			    v := d[k]
				if string(v[0]) == "-" {
					d[k] = d[totlen-1]
					totlen--
					dirty = true
				}
			}
		}
		return d[:totlen]
	}

	marker := func(d []string, i int) {
		d[i] = "----------------------"
	}

	var remaining = len(data)
	for n := 0; n < keylen; n++ {
		ones, zeros := counter(n, want)
		for line, s := range data {
			if ones[n] > zeros[n] {
				if s[n] == '1' {
					continue
				} else {
					marker(data, line)
					remaining--
				}
			} else if ones[n] < zeros[n] {
				if s[n] == '0' {
					continue
				} else {
					marker(data, line)
					remaining--
				}
			} else {
				if string(s[n]) == keeper {
					continue
				} else {
					marker(data, line)
					remaining--
				}
			}
			if remaining == 1 {
				data = cleaner(data)
				return data[0]
			}
		}
		data = cleaner(data)
	}
	return "" // bad if this happens
}

var datafile string
var verbose bool

func init() {
	rootCmd.AddCommand(day1aCmd, day1bCmd, day2aCmd, day2bCmd, day3aCmd, day3bCmd)
	rootCmd.PersistentFlags().StringVarP(&datafile, "data", "d", "", "data input")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}

func PrintKeys(dmap map[string]bool) {
	for k, _ := range dmap {
		fmt.Printf("key: %s\n", k)
	}
}

func PrintSlice(data []string) {
	for p, s := range data {
		fmt.Printf("data[%d]: %s\n", p, s)
	}
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
