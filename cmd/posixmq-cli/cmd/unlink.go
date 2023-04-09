package cmd

import (
	"log"

	"github.com/narslan/posixmq"

	"github.com/spf13/cobra"
)

// unlinkCmd represents the unlink command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "delete queue",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		qname := args[0]

		err := posixmq.Unlink(qname)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
