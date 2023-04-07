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

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "fetch a message from queue",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		qname := args[0]
		ctx := context.Background()

		cfg := &posixmq.Config{
			QueueSize:   10,
			MessageSize: 4096,
			Name:        qname,
		}
		mq, err := posixmq.Open(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := mq.Receive(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("received message %q", string(resp))

	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
