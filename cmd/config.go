package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	CTYFile    = "cty.csv"
	QSLMngFile = "managers.csv"
	CFGFile    = "cfg.json"
)

var homeConfig = path.Join(".config", "qsls")

type VCfgType struct {
	DataDir             string            `mapstructure:"data-dir"`
	UserVariables       map[string]string `mapstructure:"user-var"`
	TemplateFile        string            `mapstructure:"template"`
	PriorPrefixesRegExp []string          `mapstructure:"prior-prefixes"`
}

var VCfg = VCfgType{
	filepath.Join(getHomeDir(), ".config/qsls/data"),
	map[string]string{},
	"qsls-template.tex",
	[]string{},
}

func (c *VCfgType) GetPriorPrefixesRegExp() ([]*regexp.Regexp, error) {
	var rs []*regexp.Regexp
	for _, s := range c.PriorPrefixesRegExp {
		r, err := regexp.Compile(s)
		if err != nil {
			return nil, errors.Wrapf(err, "compile regexp '%s'", s)
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (c *VCfgType) GetTemplateFileName() string {
	return c.TemplateFile
}

func (v *VCfgType) GetDataFile(name string) string {
	return filepath.Join(VCfg.DataDir, name)
}

func initViper() {
	cfgJSON, _ := json.Marshal(VCfg)
	err := viper.ReadConfig(bytes.NewBuffer(cfgJSON))
	if err != nil {
		logrus.Fatalf("Viper fatal error during reading '%s'", err)
	}
	pf := rootCmd.PersistentFlags()
	pf.String("data-dir", VCfg.DataDir, "directory with qsls databases (DXCC, QSL Managers, etc.)")
	pf.String("template", VCfg.TemplateFile, "LaTeX template (assets must be stored in same directory)")
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.Fatalf("Viper binding: '%s'", err)
	}
}

func finishConfig() {
	err := viper.Unmarshal(&VCfg)
	if err != nil {
		logrus.Fatalf("Viper unmarshal: '%s'", err)
	}
	if !filepath.IsAbs(VCfg.TemplateFile) {
		absTemplate, err := filepath.Abs(VCfg.TemplateFile)
		if err != nil {
			logrus.Warnf("can't get abs of file '%s': %s", VCfg.TemplateFile, err)
			return
		} else {
			VCfg.TemplateFile = absTemplate
		}
	}
}

func initViperFile() {
	viper.SetConfigFile(globalCfg.cfgFile)
	logrus.Debugf("Using config '%s'", viper.ConfigFileUsed())
	_, err := os.Stat(globalCfg.cfgFile)
	if err == nil {
		err := viper.ReadInConfig()
		if err != nil {
			logrus.Fatalf("Reading config '%s' returned '%s'", viper.ConfigFileUsed(), err.Error())
		}
	} else if os.IsNotExist(err) {
		logrus.Debug("No config found")
		err := createDefaultConfig()
		if err != nil {
			logrus.Fatalf("Can't create config in your home directory (%s), please specify it manually", err.Error())
		}
		err = viper.WriteConfig()
		if err != nil {
			logrus.Fatalf("Write config returned '%s'", err.Error())
		}
		logrus.Debug("Config created")
	} else {
		logrus.Fatalf("Stat config '%s' returned '%s'", viper.ConfigFileUsed(), err.Error())
	}
}

func createDefaultConfig() error {
	f, err := safeOpen(path.Join(getHomeDir(), homeConfig, "cfg.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrap(err, "create home config file")
	}
	f.Close()
	return nil
}

func getHomeDir() string {
	h, err := homedir.Dir()
	if err != nil {
		logrus.Fatalf("Access home dir: '%s'", err)
	}
	return h
}
