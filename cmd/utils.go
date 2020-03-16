package cmd

import (
	"bitbucket.org/olik636/qsls/converter"
	"bitbucket.org/olik636/qsls/dxcc"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoDataFiles error = errors.Errorf("data files not found")

type WriteCounter struct {
	Total         uint64
	ContentLength int64
	Filename      string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 50))
	if wc.ContentLength > 0 {
		fmt.Printf("\rDownloading %s ... %s / %s - %d%% complete",
			wc.Filename,
			humanize.Bytes(wc.Total),
			humanize.Bytes(uint64(wc.ContentLength)),
			wc.Total*100/uint64(wc.ContentLength))
	} else {
		fmt.Printf("\rDownloading %s ... %s / ? complete",
			wc.Filename,
			humanize.Bytes(wc.Total))
	}
}

func safeOpen(filename string, flags int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(filename, flags, perm)
	if err == nil {
		return f, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	logrus.Infof("Creating directory '%s' for '%s'", filepath.Dir(filename), filename)
	err = os.MkdirAll(filepath.Dir(filename), 0744)
	if err != nil {
		return nil, err
	}
	f, err = os.OpenFile(filename, flags, perm) // file didnt exist so open again
	return f, errors.Wrap(err, "open file")
}

var adifCfg = struct {
	AdifFile string
}{
	"log.adif",
}

func loadAdif() ([]*converter.QSLCard, *dxcc.Entities) {
	e, err := initEntityDB()
	if err == ErrNoDataFiles {
		logrus.Fatalf("Data files in %s not found, try running `qsls update` first", VCfg.DataDir)
	} else if err != nil {
		logrus.Fatalf("Initializing DXCC entities: '%s'", err)
	}
	f, err := os.Open(adifCfg.AdifFile)
	if err != nil {
		logrus.Fatalf("Opening ADIF input file: '%s'", err)
	}
	q, err := converter.ParserWrapper(f, e)
	if err != nil {
		logrus.Fatalf("Parsing ADIF: '%s'", err)
	}
	return q, e
}

func setupQSLCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&adifCfg.AdifFile, "adif", "a", adifCfg.AdifFile, "input adif file")
}

func initEntityDB() (*dxcc.Entities, error) {
	f, err := os.Open(VCfg.GetDataFile(CTYFile))
	if os.IsNotExist(err) {
		return nil, ErrNoDataFiles
	} else if err != nil {
		return nil, errors.Wrap(err, "open CTY file")
	}
	defer f.Close()
	e, err := dxcc.NewEntityDB(f)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &e, nil
}
