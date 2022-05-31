package seq

import (
	"errors"
	"fmt"
	"strings"
)

var GENOME_HG19 = map[string]int{
	"1":  249250621,
	"2":  243199373,
	"3":  198022430,
	"4":  191154276,
	"5":  180915260,
	"6":  171115067,
	"7":  159138663,
	"8":  146364022,
	"9":  141213431,
	"10": 135534747,
	"11": 135006516,
	"12": 133851895,
	"13": 115169878,
	"14": 107349540,
	"15": 102531392,
	"16": 90354753,
	"17": 81195210,
	"18": 78077248,
	"19": 59128983,
	"20": 63025520,
	"21": 48129895,
	"22": 51304566,
	"X":  155270560,
	"Y":  59373566,
}

var GENOME_HG38 = map[string]int{
	"1":  248956422,
	"2":  242193529,
	"3":  198295559,
	"4":  190214555,
	"5":  181538259,
	"6":  170805979,
	"7":  159345973,
	"8":  145138636,
	"9":  138394717,
	"10": 133797422,
	"11": 135086622,
	"12": 133275309,
	"13": 114364328,
	"14": 107043718,
	"15": 101991189,
	"16": 90338345,
	"17": 83257441,
	"18": 80373285,
	"19": 58617616,
	"20": 64444167,
	"21": 46709983,
	"22": 50818468,
	"X":  156040895,
	"Y":  57227415,
}
var GENOME_MT = map[string]int{
	"MT": 16569,
}

var GENOME map[string]int

func SetGenome(builder string) error {
	switch strings.ToLower(builder) {
	case "hg19", "grch37":
		GENOME = GENOME_HG19
	case "hg38", "grch38":
		GENOME = GENOME_HG38
	case "m", "mt", "mito":
		GENOME = GENOME_MT
	default:
		return errors.New(fmt.Sprintf("error builder: %s, the choice is (hg19, hg38, grch38, grch37, m, mt, mito)", builder))
	}
	return nil
}
