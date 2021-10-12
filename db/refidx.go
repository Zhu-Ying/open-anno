package db

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type RefIndex struct {
	Chrom       string
	Start       int
	End         int
	Transcripts []string
}

func (r RefIndex) Range() (int, int) {
	order, _ := ChromArray.GetByName(r.Chrom)
	start := order*1e9 + r.Start
	end := order*1e9 + r.End
	return start, end
}

func (r *RefIndex) SetTranscript(refgenes Refgenes) {
	for _, refgene := range refgenes {
		starti, endi := r.Range()
		startj, endj := refgene.Range()
		if starti <= endj && endi >= startj {
			r.Transcripts = append(r.Transcripts, refgene.SN())
		}
	}
}

type RefIndexes []RefIndex

func (r RefIndexes) SetTranscript(refgenes Refgenes, refIndexesChannel chan RefIndexes) {
	for i := 0; i < len(r); i++ {
		r[i].SetTranscript(refgenes)
	}
	refIndexesChannel <- r
}

func (r RefIndexes) Len() int {
	return len(r)
}

func (r RefIndexes) Less(i, j int) bool {
	starti, endi := r[i].Range()
	startj, endj := r[j].Range()
	if starti == startj {
		return endi < endj
	} else {
		return starti < startj
	}
}

func (r RefIndexes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func InitRefIndexes() RefIndexes {
	refIndexes := make(RefIndexes, 0)
	for _, chrom := range ChromArray {
		for i := 0; i < chrom.Length; i += RefIndexStepLen {
			end := i + RefIndexStepLen
			if end > chrom.Length {
				end = chrom.Length
			}
			refIndex := RefIndex{Chrom: chrom.Name, Start: i + 1, End: end}
			refIndexes = append(refIndexes, refIndex)
		}
	}
	sort.Sort(refIndexes)
	return refIndexes
}

func NewRefIndexes(refIndexFile string) RefIndexes {
	refIndexes := make(RefIndexes, 0)
	if fp, err := os.Open(refIndexFile); err == nil {
		defer func(fp *os.File) {
			err := fp.Close()
			if err != nil {
				log.Panic(err.Error())
			}
		}(fp)
		reader := bufio.NewReader(fp)
		for {
			if line, err := reader.ReadString('\n'); err == nil {
				line = strings.TrimSpace(line)
				if len(line) == 0 || line[0] == '#' {
					continue
				}
				fields := strings.Split(line, "\t")
				start, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Panic(err.Error())
				}
				end, err := strconv.Atoi(fields[2])
				if err != nil {
					log.Panic(err.Error())
				}
				refIndexes = append(refIndexes, RefIndex{
					Chrom: fields[0], Start: start, End: end, Transcripts: strings.Split(fields[3], ","),
				})
			} else {
				if err == io.EOF {
					break
				} else {
					log.Panic(err.Error())
				}
			}
		}
	} else {
		log.Panic(err.Error())
	}
	return refIndexes
}

func CreateRefIndex(refgenes Refgenes, refgeneIndexFile string) {
	refIndexes := InitRefIndexes()
	fo, err := os.Create(refgeneIndexFile)
	if err == nil {
		defer func(fo *os.File) {
			err := fo.Close()
			if err != nil {
				log.Panic(err.Error())
			}
		}(fo)
		refgeneMap := make(map[string]Refgenes)
		for _, refgene := range refgenes {
			if rs, ok := refgeneMap[refgene.Chrom]; ok {
				refgeneMap[refgene.Chrom] = append(rs, refgene)
			} else {
				refgeneMap[refgene.Chrom] = Refgenes{refgene}
			}
		}
		refIndexMap := make(map[string]RefIndexes)
		for _, refIndex := range refIndexes {
			if rs, ok := refIndexMap[refIndex.Chrom]; ok {
				refIndexMap[refIndex.Chrom] = append(rs, refIndex)
			} else {
				refIndexMap[refIndex.Chrom] = RefIndexes{refIndex}
			}
		}
		refIndexesChannel := make(chan RefIndexes)
		for _, chrom := range ChromArray {
			go refIndexMap[chrom.Name].SetTranscript(refgeneMap[chrom.Name], refIndexesChannel)
		}
		for _, _ = range ChromArray {
			if refIndexes_callback, ok := <-refIndexesChannel; ok {
				for _, refIndex := range refIndexes_callback {
					if len(refIndex.Transcripts) > 0 {
						if _, err := fo.WriteString(fmt.Sprintf(
							"%s\t%d\t%d\t%s\n",
							refIndex.Chrom, refIndex.Start, refIndex.End,
							strings.Join(refIndex.Transcripts, ","),
						)); err != nil {
							log.Panic(err.Error())
						}
					}
				}

			} else {
				break
			}
		}
		close(refIndexesChannel)
	}
}
