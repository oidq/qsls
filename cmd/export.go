package cmd

import (
	"github.com/oIdq/qsls/converter"
	"github.com/oIdq/qsls/sorter"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rwestlund/gotex"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

var exportCfg = struct {
	OutputFile  string
	OutputLaTeX bool
}{
	"QSLs.pdf",
	false,
}

func init() {
	setupQSLCommand(exportCmd)
	setupExportCommand(exportCmd)
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export QSOs",
	Long:  `Export QSL cards in PDF`,
	Run: func(cmd *cobra.Command, args []string) {
		q, e := loadAdif()
		var qsls []*converter.QSLCard
		prior, err := VCfg.GetPriorPrefixesRegExp()
		if err != nil {
			logrus.Fatalf("Get prior prefixes: %s", err)
		}
		s := sorter.NewSorter(prior)
		if globalCfg.advancedSorter {
			packets := s.SortQSLsByDXCC(e, q)
			qsls = packets.Convert()
		} else {
			qsls = s.SortQSLsByAlphabet(q)
		}
		temp, err := templateQSLs(qsls)
		if err != nil {
			logrus.Fatalf("Templating QSLs: %s", err)
		}
		if exportCfg.OutputLaTeX {
			fmt.Print(temp)
			return
		}
		pdf, err := processLatex(temp)
		if err != nil {
			logrus.Fatalf("Processing LaTeX: %s", err)
		}
		pdfFile, err := safeOpen(exportCfg.OutputFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			logrus.Fatalf("Opening output file: %s", err)
		}
		_, err = pdfFile.Write(pdf)
		if err != nil {
			logrus.Fatalf("Writing output file: %s", err)
		}
		pdfFile.Close()
	},
}

type templateForm struct {
	QSLs []*converter.QSLCard
	User map[string]string
}

func templateQSLs(qsls []*converter.QSLCard) (string, error) {
	tempFileName := VCfg.GetTemplateFileName()
	qslsTemplate, err := template.New(filepath.Base(tempFileName)).ParseFiles(tempFileName)
	if err != nil {
		return "", errors.Wrap(err, "read template file")
	}
	buf := bytes.NewBufferString("")
	form := templateForm{qsls, VCfg.UserVariables}
	err = qslsTemplate.Execute(buf, form)
	if err != nil {
		return "", errors.Wrap(err, "execute template")
	}
	return buf.String(), nil
}

func processLatex(latex string) ([]byte, error) {
	pdf, err := gotex.Render(latex, gotex.Options{
		Texinputs: filepath.Dir(VCfg.GetTemplateFileName()),
	})
	return pdf, errors.Wrap(err, "render pdf")
}

func setupExportCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&exportCfg.OutputFile, "output-file", "o", exportCfg.OutputFile, "output PDF file name")
	cmd.Flags().BoolVar(&exportCfg.OutputLaTeX, "output-latex", exportCfg.OutputLaTeX, "output LaTeX file to stdout")
}
