package anno

import (
	"fmt"
	"log"
	"open-anno/pkg/io"
	"strings"

	"github.com/spf13/cobra"
)

func NewAnnoCmd(varType string) *cobra.Command {
	varType = strings.ToLower(varType)
	if varType != "snv" && varType != "cnv" {
		log.Fatalln("only 'snv' or 'cnv'")
	}
	cmd := &cobra.Command{
		Use:   varType,
		Short: fmt.Sprintf("Annotate for %s", strings.ToUpper(varType)),
		Run: func(cmd *cobra.Command, args []string) {
			var param AnnoParam
			param.Input, _ = cmd.Flags().GetString("avinput")
			param.DBpath, _ = cmd.Flags().GetString("dbpath")
			param.Builder, _ = cmd.Flags().GetString("builder")
			param.AAshort, _ = cmd.Flags().GetBool("aashort")
			param.Exon, _ = cmd.Flags().GetBool("exon")
			param.Overlap, _ = cmd.Flags().GetFloat64("overlap")
			dbtypes, _ := cmd.Flags().GetString("dbtypes")
			dbnames, _ := cmd.Flags().GetString("dbnames")
			if len(args) <= 3 {
				err := cmd.Help()
				if err != nil {
					log.Panic(err)
				}
			} else {
				dbNames := strings.Split(dbnames, ",")
				dbTypes := strings.Split(dbtypes, ",")
				errChan := make(chan error, len(dbNames))
				for i := 0; i < len(dbNames); i++ {
					param.DBname = strings.TrimSpace(dbNames[i])
					param.DBType = strings.TrimSpace(dbTypes[i])
					go param.RunAnno(varType, errChan)
				}
				for i := 0; i < len(dbNames); i++ {
					err := <-errChan
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}
	cmd.Flags().StringP("avinput", "i", "", "Annotated Variants Input File")
	cmd.Flags().StringP("outprefix", "o", "", "Output Prefix")
	cmd.Flags().StringP("dbpath", "d", "", "Database Directory")
	cmd.Flags().StringP("dbnames", "n", "", "Database Names")
	cmd.Flags().StringP("dbtypes", "t", "", "Database Types")
	cmd.Flags().StringP("builder", "b", "hg19", "Database Builder")
	if varType == "snv" {
		cmd.Flags().BoolP("aashort", "s", false, "Database Builder")
		cmd.Flags().BoolP("exon", "e", false, "Output ExonOrder Instead of TypeOrder")
	}
	if varType == "cnv" {
		cmd.Flags().Float64P("overlap", "p", 0.75, "CNV Overlap Threshold")
	}
	return cmd
}

func NewMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge Annotation",
		Run: func(cmd *cobra.Command, args []string) {
			avinput, _ := cmd.Flags().GetString("avinput")
			genebased, _ := cmd.Flags().GetString("genebased")
			otherbaseds, _ := cmd.Flags().GetStringArray("otherbaseds")
			outfile, _ := cmd.Flags().GetString("outfile")
			if avinput == "" || genebased == "" || outfile == "" {
				err := cmd.Help()
				if err != nil {
					log.Panic(err)
				}
			} else {
				io.MergeAnno(outfile, avinput, genebased, otherbaseds...)
			}
		},
	}
	cmd.Flags().StringP("avinput", "i", "", "Annotated Variants Input File")
	cmd.Flags().StringP("genebased", "g", "", "GeneBased Annotation File")
	cmd.Flags().StringArrayP("otherbaseds", "d", []string{}, "FilterBased or RegionBased Annotation Files")
	cmd.Flags().StringP("outfile", "o", "", "Output Merged Annotation File")
	return cmd
}
