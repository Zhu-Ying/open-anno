package fitler

import (
	"OpenAnno/anno"
	"OpenAnno/variant"
)

type FilterBasedAnno map[string]string

func (f FilterBasedAnno) AnnoType() anno.AnnoType {
	return anno.AnnoType_FILTER
}

func ReadFilterBasedFields(headers []string, fields []string) (variant.Variant, FilterBasedAnno) {
	_anno := make(FilterBasedAnno)
	_variant := variant.NewVariant(fields[0], fields[1], fields[2], fields[3], fields[4])
	for i := 5; i < len(fields); i++ {
		_anno[headers[i]] = fields[i]
	}
	return _variant, _anno
}
