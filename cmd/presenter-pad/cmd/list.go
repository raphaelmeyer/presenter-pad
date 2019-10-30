package cmd

import (
	"log"
	"os"
	"presenter-pad/internal/pkg/mapper"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all input devices",
	Run: func(cmd *cobra.Command, args []string) {
		err := mapper.ListDevices()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
