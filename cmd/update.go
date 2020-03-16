package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

const (
	ctyURL = "http://www.country-files.com/cty/cty.csv"
)

type UpdateCfg struct {
}

func init() {
	var cfg = UpdateCfg{
	}
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "update qsls databases",
		Long:  `Update all qsls databases`,
		Run: func(cmd *cobra.Command, args []string) {
			Update(cfg)
		},
	}
	rootCmd.AddCommand(updateCmd)
}

func Update(c UpdateCfg) {
	logrus.Debugf("Writing CTY to %s", VCfg.GetDataFile(CTYFile))
	f := wrapSafeOpen(VCfg.GetDataFile(CTYFile))
	defer f.Close()
	downloadToFile(ctyURL, f)
}

func wrapSafeOpen(file string) *os.File {
	f, err := safeOpen(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Fatalf("Opening '%s': %s", file, err)
	}
	return f
}

func downloadToFile(url string, outputFile *os.File) {
	r, err := http.Get(url)
	if err != nil {
		logrus.Fatalf("Downloading CTY.csv file: '%s'", err)
	} else if r.StatusCode != 200 {
		r.Body.Close()
		logrus.Fatalf("Downloading CTY.csv returned status %s", r.Status)
	}
	defer r.Body.Close()
	wc := &WriteCounter{ContentLength: r.ContentLength}
	_, err = io.Copy(outputFile, io.TeeReader(r.Body, wc))
	fmt.Println() // newline after updated content
	if err != nil {
		logrus.Fatalf("Writing downloaded output: %s", err)
	}
	return
}
