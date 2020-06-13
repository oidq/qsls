package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const (
	ctyFileFlag = "cty-file"
)

var globalCfg = struct {
	cfgFile        string
	debug          bool
	verbose        bool
	saveConfig     bool
	advancedSorter bool
	version        string
}{
	filepath.Join(getHomeDir(), homeConfig, CFGFile),
	false,
	false,
	false,
	false,
	"v0.0.0-0",
}

var rootCmd = &cobra.Command{
	Use:   "qsls",
	Short: "QSLs generates qsl with LaTeX",
	Long:  `QSLs generates QSL cards with LaTeX from adif files`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatalf("Help returned error: '%s'", err.Error())
		}
	},
}

func initConfig() {
	initLogrus()
	initViperFile()
	finishConfig()
}

func Execute(version string) {
	globalCfg.version = version
	initViper()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&globalCfg.cfgFile, "config", globalCfg.cfgFile, "JSON config file")
	rootCmd.PersistentFlags().BoolVarP(&globalCfg.debug, "debug", "d", false, "debug")
	rootCmd.PersistentFlags().BoolVarP(&globalCfg.verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&globalCfg.saveConfig, "save-config", false, "save current parameters to config")
	rootCmd.PersistentFlags().BoolVarP(&globalCfg.advancedSorter, "advanced-sorter", "s", false, "use advance DXCC packet sorter")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errors.Wrap(err, "executing cobra"))
		os.Exit(1)
	}
	if viper.ConfigFileUsed() != "" && globalCfg.saveConfig {
		if err := viper.WriteConfig(); err != nil {
			logrus.Fatalf("Writing config failed: '%s'", err)
		}
	}
}

func initLogrus() {
	logrus.SetOutput(os.Stderr)
	if globalCfg.debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else if globalCfg.verbose {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
