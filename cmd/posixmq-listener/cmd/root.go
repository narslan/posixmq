package cmd

import (
	"log"
	"os"

	"github.com/narslan/posixmq"
	"github.com/narslan/posixmq/poll"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "posixmq listener a queue for sent events",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		qname := args[0]
		// Initialize netpoll instance. We will use it to be noticed about incoming
		// events from listener of user connections.
		poller, err := poll.New(nil)
		if err != nil {
			log.Fatal(err)
		}

		var (
			exit = make(chan struct{})
		)

		cfg := &posixmq.Config{
			QueueSize:   10,
			MessageSize: 4096,
			Name:        qname,
		}

		mq, err := posixmq.Open(cfg)
		if err != nil {
			log.Fatal(mq, err)
		}

		// Create poll event descriptor for mq.
		// We want to handle only read events of it.
		desc := poll.Must(poll.HandleRead(mq))

		// Subscribe to events of mq.
		poller.Start(desc, func(ev poll.Event) {
			if ev&poll.EventHup != 0 {
				// When a Hup received, stop polling  .
				poller.Stop(desc)
				mq.Close()
				return
			}
			resp, err := mq.Receive()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("epoll receives %q", string(resp))
		})

		<-exit
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
