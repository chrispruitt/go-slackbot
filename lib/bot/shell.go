package bot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	c "github.com/chrispruitt/go-slackbot/lib/config"
	"github.com/common-nighthawk/go-figure"
	"github.com/mattes/go-asciibot"
	"github.com/slack-go/slack/slackevents"
)

func shell() {
	banner()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("slackbot> ")
		cmdString, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if strings.TrimSpace(cmdString) != "" {
			err = runCommand(cmdString)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func banner() {
	myFigure := figure.NewColorFigure("Slack bot", "", "green", true)
	myFigure.Print()
	fmt.Println(asciibot.Random())
	fmt.Println("Welcome to the slackbot shell! Type 'help' for help, 'exit' to exit.")
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	switch strings.TrimSpace(arrCommandStr[0]) {
	case "exit":
		os.Exit(0)
	default:
		// Mock Slack MessageEvent
		event := &slackevents.MessageEvent{
			Text:    fmt.Sprintf("%s %s", c.BotName, commandStr),
			Channel: os.Getenv("SHELL_MODE_CHANNEL"),
		}
		HandleMessageEvent(event)
		return nil
	}
	cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
