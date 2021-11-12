package variant

import (
	"OpenAnno/db/chromosome"
)

type CnvType string

const (
	CnvType_DUP CnvType = "DUP"
	CnvType_DEL CnvType = "DEL"
)

type Cnv struct {
	Variant
	CopyNumber int `json:"copy_number"`
}

func (c Cnv) Type() CnvType {
	if c.CopyNumber < 2 {
		return CnvType_DEL
	}
	return CnvType_DUP
}

func (c Cnv) VarType() VarType {
	return VarType_CNV
}

type Cnvs []Cnv

func (c Cnvs) GetVariant(i int) IVariant {
	return c[i]
}

func (c Cnvs) Len() int {
	return len(c)
}

func (c Cnvs) Less(i, j int) bool {
	chromOrderi, _ := chromosome.ChromList.GetByName(c[i].Chrom)
	chromOrderj, _ := chromosome.ChromList.GetByName(c[j].Chrom)
	return chromOrderi < chromOrderj || c[i].Start < c[j].Start || c[j].End < c[j].End

}

func (c Cnvs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Cnvs) FilterByChrom(chrom string) Cnvs {
	cnvs := make(Cnvs, 0)
	for _, cnv := range c {
		if cnv.Chrom == chrom {
			cnvs = append(cnvs, cnv)
		}
	}
	return cnvs
}
