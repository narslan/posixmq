/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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
		cfg := &posixmq.Config{Name: qname}
		mq := posixmq.MessageQueue{Config: cfg}

		err := mq.Unlink(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unlinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unlinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
