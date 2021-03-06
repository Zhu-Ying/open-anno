package db

import (
	"fmt"
	"open-anno/pkg"
	"open-anno/pkg/io"
	"strings"
)

func AnnoRegionBased(infile, dbfile, outfile string, overlap float64) error {
	// 读取变异
	variantMap, err := io.ReadVariantMap(infile)
	if err != nil {
		return err
	}
	// 读取BED DB文件
	var regMap map[string]io.BEDs
	refReader, err := io.NewIoReader(infile)
	if err != nil {
		return err
	}
	defer refReader.Close()
	var scanner io.DBRegScanner
	scanner = io.NewDBRegScanner(refReader)
	for scanner.Scan() {
		row, err := scanner.Row()
		if err != nil {
			return err
		}
		if _, ok := regMap[row.Chrom]; ok {
			regMap[row.Chrom] = append(regMap[row.Chrom], row)
		} else {
			regMap[row.Chrom] = io.BEDs{row}
		}
	}
	// 定义输出
	writer, err := io.NewIoWriter(outfile)
	if err != nil {
		return err
	}
	defer writer.Close()
	// 开始注释
	// 输出表头
	headers := strings.Split(scanner.Header, "\t")
	fmt.Fprintf(writer, "Chr\tStart\tEnd\tRef\tAlt\t%s\n", headers[3])
	for chrom, variants := range variantMap {
		if regs, ok := regMap[chrom]; ok {
			for _, variant := range variants {
				var annos []string
				for _, reg := range regs {
					if variant.End >= reg.Start && variant.Start <= reg.End {
						vlen := variant.End - variant.Start + 1
						olen := pkg.Min(variant.End, reg.End) - pkg.Max(variant.Start, reg.Start) + 1
						if float64(olen)/float64(vlen) >= overlap {
							annos = append(annos, reg.Name)
						}
					}
				}
				if len(annos) > 0 {
					fmt.Fprintf(writer, "%s\t%d\t%d\t%s\t%s\t%s\n",
						variant.Chrom, variant.Start, variant.End,
						variant.Ref, variant.Alt, strings.Join(annos, ","))
				}
			}
		}
	}
	return err
}
