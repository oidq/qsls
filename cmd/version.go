package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "QSLs version",
	Long:  `Print QSLs version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running QSLs version %s\n", globalCfg.version)
	},
}
