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
	"sort"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		P, recipes = ReadInput(datafile)
		fmt.Printf("We have a polymer '%s' with %d chars, %d pairs (%v) and %d recipes\n",
			P.Seq, len(P.Seq), len(P.Pairs), P.Pairs, len(recipes))
		fmt.Printf("Part 1:\n")
		for i := 1 ; i <= 10 ; i++ {
		    P = P.Apply(recipes)
		    fmt.Printf("Sequence after %d steps (len %d) \n", i, len(P.Seq))
		}
		fmt.Printf("Counts: %v\n", P.Counts())
		fmt.Printf("Part 2:\n")
		P2, recipes = ReadInput(datafile)
		for i := 1 ; i <= 40 ; i++ {
		    P2 = P2.ApplyNG(recipes)
		    fmt.Printf("Counts after %d steps (len %d): %v\n", i, len(P2.Pairs), P2.Pairs)
		}
		c := P2.CountsNG()
		fmt.Printf("Counts: %v\n", c)
		res := c[len(c)-1] - c[0]
		fmt.Printf("Result (largest - smallest): %d\n", res)		
	},
}

var P, P2 Polymer
var recipes map[string]string

type Polymer struct {
     Seq     string
     Pairs   map[string]int	// frequency of each pair of genes
     Freq    map[string]int	// frequency of individual genes
     First   string
     Last    string
}

func (p *Polymer) Counts() []int {
     counts := []int{}
     p.Freq = map[string]int{}
     for _, v := range p.Seq {
     	 p.Freq[string(v)]++
     }
     for _, c := range p.Freq {
     	 counts = append(counts, c)
     }
     sort.Ints(counts)
     return counts
}

func (p *Polymer) CountsNG() []int {
     counts := []int{}
     p.Freq = map[string]int{}
     fmt.Printf("Counting pairs: ")
     for k, v := range p.Pairs {
     	 fmt.Printf(" [%s: %d]", k, v)
     	 k1 := string(k[0])
     	 k2 := string(k[1])
     	 p.Freq[k1] += v
     	 p.Freq[k2] += v
     }
     p.Freq[p.First]++
     p.Freq[p.Last]++
     for k, v := range p.Freq {
     	 p.Freq[k] = v / 2
     }
     fmt.Println()
     for _, c := range p.Freq {
     	 counts = append(counts, c)
     }
     fmt.Printf("CountsNG: first: '%s', last: '%s', p.Freq: %v\n",
     			   p.First, p.Last, p.Freq)
     fmt.Printf("CountsNG: counts: %v\n", counts)
     sort.Ints(counts)
     return counts
}

func (p *Polymer) Apply(recipes map[string]string) Polymer {
     seq := ""
     var c1, c2 string
     for i := 0 ; i < len(p.Seq)-1 ; i++ {
     	 c1 = string(p.Seq[i])
     	 c2 = string(p.Seq[i+1])
	 seq += c1
	 if v, exists := recipes[c1+c2]; exists {
	    seq += v
	 }
     }
     seq += string(p.Seq[len(p.Seq)-1])
     return Polymer{ Seq: seq }
}

func (p *Polymer) ApplyNG(recipes map[string]string) Polymer {
     pairs := map[string]int{}
     fmt.Printf("Applying %d recipes:\n", len(recipes))
     for r, v := range recipes {
     	 fmt.Printf("Recipe: %s-->%s: ", r, v)
     	 if ec, ok := p.Pairs[r]; ok {
     	    fmt.Printf(" [%s: %d]", r, ec)
	    newk1 := string(r[0])+string(v)
	    newk2 := string(v)+string(r[1])
	    pairs[newk1] += ec
    	    pairs[newk2] += ec
	    fmt.Printf("-->> [%s:%d] + [%s:%d]",
	    		     newk1, pairs[newk1], newk2, pairs[newk2])
	 }
	 fmt.Println()
     }
     fmt.Println()
     return Polymer{ Pairs: pairs, First: p.First, Last: p.Last }
}

func ReadInput(filename string) (Polymer, map[string]string) {
        r := map[string]string{}
	pairs := map[string]int{}
	pair := ""
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: failed to open %s", filename)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	line := scanner.Text()
	firstc := string(line[0])
	lastc  := string(line[len(line)-1])
	fmt.Printf("ReadInput: first: '%s' last: '%s'\n", firstc, lastc)
	for i := 0; i <= len(line) - 2; i++ {
	    pair = string(line[i])+string(line[i+1])
	    pairs[pair]++
	}
	p := Polymer{
	     	     Seq: line,
		     Pairs: pairs,
		     First: firstc,
		     Last: lastc,
	     }

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, " -> ") {
		   tmp := strings.Split(line, " -> ")
		   r[tmp[0]] = tmp[1]
		}
	}
	file.Close()
	
	fmt.Printf("ReadInput: Polymer: %v\n", p)
	return p, r
}
