package cmd

import (
	"context"
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

		ctx := context.Background()
		mq, err := posixmq.Open(ctx, cfg)
		if err != nil {
			log.Fatal(mq, err)
		}

		data := []byte(msg)

		err = mq.Send(ctx, data, 2)
		if err != nil {
			log.Fatal(mq, err)
		}

		mq.Close(ctx)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
