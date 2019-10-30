package cmd

import (
	"log"
	"os"
	"presenter-pad/internal/pkg/mapper"

	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start the gamepad to keystroke mapper",
		Run: func(cmd *cobra.Command, args []string) {
			err := mapper.Run(device)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	device string
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&device, "device", "d", "", "gamepad device name")
}
