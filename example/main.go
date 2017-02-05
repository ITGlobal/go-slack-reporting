package main

import (
	slrep "github.com/itglobal/go-slack-reporting"
	"github.com/mkideal/cli"
	"log"
	"os"
	"time"
)

type argT struct {
	cli.Helper
	AccessToken string `cli:"token"    usage:"Slack access token" dft:"$SLACK_TOKEN"`
	Channel     string `cli:"channel"  usage:"Slack channel"      dft:"$SLACK_CHANNEL"`
}

const (
	maxSteps = 10
	delay    = time.Second * 2
)

var operations = []string{
	"Checking out sources...",
	"Cleaning temporary directories...",
	"Compiling project...",
	"*Warning* There are some compilation warnings",
	":exclamation: Test _my-sample-test_ is taking too long",
	":arrow_forward: Test run results:\n> Test 1/3 : OK\n> Test 2/3 : OK\n> Test 3/3 : OK",
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		l := log.New(os.Stderr, "slack ", log.Ltime)
		log.SetPrefix("      ")
		log.SetFlags(log.Ltime)

		config := slrep.NewConfig(argv.AccessToken)
		config.SetChannel(argv.Channel)
		config.SetLogger(l)
		config.SetUsername("Test bot")
		config.SetIcon(slrep.NewIconFromEmoji(":information_source:"))

		log.Printf("Connecting to Slack\n")
		reporter, err := config.CreateReporter()
		if err != nil {
			return err
		}

		log.Printf("Posting initial status message\n")
		msg, err := reporter.BeginMessage("`[          ]` Preparing some long-term operation...")
		if err != nil {
			return err
		}

		for i := 0; i < maxSteps; i++ {
			time.Sleep(delay)
			text := "`["

			j := 0
			for ; j < i; j++ {
				text += "#"
			}
			for ; j < maxSteps; j++ {
				text += " "
			}

			text += "]` "
			text += operations[i%len(operations)]

			log.Printf("Updating status message (%d/%d)\n", i+1, maxSteps)
			err = msg.Update(text)
			if err != nil {
				return err
			}
		}

		time.Sleep(delay)
		log.Printf("Deleting status message\n")
		err = msg.Delete()
		if err != nil {
			return err
		}

		return nil
	})
}
