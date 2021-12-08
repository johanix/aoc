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
	"strconv"
	"strings"

	// "github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

// day4Cmd represents the day4 command
var day8aCmd = &cobra.Command{
	Use:   "day8a",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		if datafile == "" {
			datafile = "data/day-8.dat"
		}

		file, err := os.Open(datafile)
		if err != nil {
			log.Fatalf("Error: failed to open %s", datafile)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		readings, _ := ReadInput(scanner, verbose)
		lc := CountLengths(readings)

		fmt.Printf("|x x 1 7 4 x x 8|\n")
		fmt.Printf("%v\n", lc)
		fmt.Printf("1,4,7,8 occur %d times\n", lc[2]+lc[4]+lc[3]+lc[7])

		fmt.Printf("StringToDigit: %v\n", StringToDigit)
		sum := 0
		for _, r := range readings {
			r.MatchNumber()
			r.Print()
			sum += r.Number
		}
		fmt.Printf("Part A sum: %d\n", sum)
	},
}

var day8bCmd = &cobra.Command{
	Use:   "day8b",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Runner8()
	},
}

func Runner8() {
	if datafile == "" {
		datafile = "data/day-8.dat"
	}

	file, err := os.Open(datafile)
	if err != nil {
		log.Fatalf("Error: failed to open %s", datafile)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	readings, _ := ReadInput(scanner, verbose)

	sum := 0
	outputs := []int{}
	testedmaps := 0
	for _, r := range readings {
		Permute([]rune("abcdefg"), func(s []rune) {
			testmap := GenerateMap(string(s))
			if ok, themap := ValidMap(r, testmap); ok {
				fmt.Printf("Map was valid for reading %s\n", r.Sprint())
				r.Mapping = themap
				r.MatchNumber()
				fmt.Printf("Result: %d\n", r.Number)
				outputs = append(outputs, r.Number)
				sum += r.Number
			}
			testedmaps++
		})
	}
	fmt.Printf("All outputs: %v\n", outputs)
	fmt.Printf("Total sum: %d\n", sum)

	file.Close()
}

func init() {
	rootCmd.AddCommand(day8aCmd, day8bCmd)

	// day4Cmd.PersistentFlags().String("foo", "", "A help for foo")
	// day4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Reading struct {
	Signals []string
	Output  []string
	Numstr  string
	Number  int
	Mapping map[string]string
}

/*
 *   11111
 *  2     3
 *  2     3
 *  2     3
 *   44444
 *  5     6
 *  5     6
 *  5     6
 *   77777
 */

type OkNum struct {
	Ok    bool
	Digit string
}

//      segments                digit
var ValidNumbers = map[string]OkNum{
	"123567":  OkNum{true, "0"},
	"36":      OkNum{true, "1"},
	"13457":   OkNum{true, "2"},
	"13467":   OkNum{true, "3"},
	"2346":    OkNum{true, "4"},
	"12467":   OkNum{true, "5"},
	"124567":  OkNum{true, "6"},
	"136":     OkNum{true, "7"},
	"1234567": OkNum{true, "8"},
	"123467":  OkNum{true, "9"},
}

var TestMapping = map[string]string{
	"a": "3", "b": "6", "c": "7", "d": "1", "e": "2", "f": "4", "g": "5",
}

func Permute(a []rune, f func([]rune)) {
	PermHelper(a, f, 0)
}

func PermHelper(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}
	PermHelper(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		PermHelper(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func GenerateMap(s string) map[string]string {
	if len(s) != 7 {
		log.Fatalf("GenerateMap: Error: input must be exactly 7 chars\n")
	}
	res := make(map[string]string, 7)
	for pos, c := range []byte(s) {
		res[string(c)] = strconv.Itoa(pos + 1)
	}
	return res
}

// mapping[string]string is "c" --> "2" where c is the letter in the output and 2 is the pos in the display
func ValidMap(r Reading, mapping map[string]string) (bool, map[string]string) {
	res := map[string]string{}
	for _, num := range r.Signals {
		segments := ""
		for _, d := range []byte(num) {
			segments += mapping[string(d)]
		}
		ssegments := SortString(segments)
		if v, ok := ValidNumbers[ssegments]; ok {
			// res[ssegments] = v.Digit
			res[SortString(num)] = v.Digit
			// fmt.Printf("Number '%s' was remapped to segment list '%s' (sorted: %s) which is valid (%s)\n", num, segments, ssegments, v.Digit)
		} else {
			// fmt.Printf("Number '%s' was remapped to segment list '%s' (sorted: %s) which is not valid\n", num, segments, ssegments)
			return false, res
		}
	}
	fmt.Printf("All numbers in reading: '%v' remapped to valid segment lists\n", r.Output)
	return true, res
}

func (r *Reading) MatchNumber() {
	fmt.Printf("Mapping outputs [%v] via mapping %v\n", r.Output, r.Mapping)
	num := ""
	for _, s := range r.Output {
		digit := r.Mapping[s]
		if len(digit) != 1 {
			fmt.Printf("MatchNumber: missed match: s='%s' digit='%s'\n", s, digit)
		}
		num += digit
	}
	r.Numstr = num
	r.Number, _ = strconv.Atoi(num)
}

func (r *Reading) SolveMapping() {
}

var StringToDigit = map[string]string{
	SortString("cagedb"):  "0",
	SortString("ab"):      "1",
	SortString("gcdfa"):   "2",
	SortString("fbcad"):   "3",
	SortString("eafb"):    "4",
	SortString("cdfbe"):   "5",
	SortString("cdfgeb"):  "6",
	SortString("dab"):     "7",
	SortString("acedgfb"): "8",
	SortString("cefabd"):  "9",
}

func CountLengths(d []Reading) []int {
	var counts = make([]int, 8)
	for _, r := range d {
		for _, o := range r.Output {
			counts[len(o)]++
		}
	}
	return counts
}

func (r *Reading) Print() {
	for _, s := range r.Signals {
		fmt.Printf("%s ", s)
	}
	fmt.Printf("|| ")
	for _, o := range r.Output {
		fmt.Printf("%s ", o)
	}
	fmt.Printf("|| '%s' %d\n", r.Numstr, r.Number)
}

func (r *Reading) Sprint() string {
	var res string
	//     for _, s := range r.Signals {
	//     	 fmt.Printf("%s ", s)
	//     }
	//     fmt.Printf("|| ")
	for _, o := range r.Output {
		res += fmt.Sprintf("%s ", o)
	}
	res += fmt.Sprintf("| '%s' %d\n", r.Numstr, r.Number)
	return res
}

type sortBytes []byte

func (s sortBytes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortBytes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortBytes) Len() int {
	return len(s)
}

func SortString(s string) string {
	b := []byte(s)
	sort.Sort(sortBytes(b))
	return string(b)
}

func ReadInput(scanner *bufio.Scanner, verbose bool) ([]Reading, error) {
	var line string
	var rs = []Reading{}
	var outputs, soutputs []string
	for scanner.Scan() {
		r := Reading{}
		line = scanner.Text()
		// fmt.Printf("Line: '%s'\n", line)
		foo := strings.Split(strings.TrimSpace(line), " | ")
		// fmt.Printf("foo: '%v'\n", foo)
		r.Signals = strings.Split(foo[0], " ")
		outputs = strings.Split(foo[1], " ")
		soutputs = []string{}
		for _, s := range outputs {
			soutputs = append(soutputs, SortString(s))
		}
		r.Output = soutputs
		rs = append(rs, r)
		r.Print()
	}
	return rs, nil
}
