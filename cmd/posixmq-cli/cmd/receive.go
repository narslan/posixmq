package cmd

import (
	"log"

	"github.com/narslan/posixmq"
	"github.com/spf13/cobra"
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "fetch a message from queue",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		qname := args[0]

		cfg := &posixmq.Config{
			QueueSize:   10,
			MessageSize: 4096,
			Name:        qname,
		}
		mq, err := posixmq.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := mq.Receive()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("received message %q", string(resp))

	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
