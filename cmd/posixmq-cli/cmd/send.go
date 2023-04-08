package cmd

import (
	"log"

	"github.com/narslan/posixmq"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send a message to a queue",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		qname := args[0]
		msg := args[1]
		cfg := &posixmq.Config{
			QueueSize:   10,
			MessageSize: 4096,
			Name:        qname,
		}

		mq, err := posixmq.Open(cfg)
		if err != nil {
			log.Fatal(mq, err)
		}

		data := []byte(msg)

		err = mq.Send(data, 2)
		if err != nil {
			log.Fatal(mq, err)
		}

		mq.Close()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
