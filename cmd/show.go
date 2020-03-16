package cmd

import (
	"bitbucket.org/olik636/qsls/sorter"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	setupQSLCommand(showCmd)
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show sorting of QSOs",
	Long:  `Show sorting of QSOs with some additional info`,
	Run: func(cmd *cobra.Command, args []string) {
		q, e := loadAdif()
		prior, err := VCfg.GetPriorPrefixesRegExp()
		if err != nil {
			logrus.Fatalf("Get prior prefixes: %s", err)
		}
		s := sorter.NewSorter(prior)
		if globalCfg.advancedSorter {
			packets := s.SortQSLsByDXCC(e, q)
			for _, v := range packets {
				ent, _ := e.LookupEntityCode(int64(v.DXCC))
				fmt.Printf("%s - %s\n", v.Prefix, ent.Entity)
				for _, q := range v.QSLs {
					if q.QSLVia != "" {
						fmt.Printf("\t%s - %d - via %s\n", q.Callsign, len(q.QSOs), q.QSLVia)
					} else {
						fmt.Printf("\t%s - %d\n", q.Callsign, len(q.QSOs))
					}
				}
			}
		} else {
			qsos := s.SortQSLsByAlphabet(q)
			for _, qso := range qsos {
				if qso.QSLVia != "" {
					fmt.Printf("%s - %d - via %s\n", qso.Callsign, len(qso.QSOs), qso.QSLVia)
				} else {
					fmt.Printf("%s - %d\n", qso.Callsign, len(qso.QSOs))
				}
			}
		}
	},
}
