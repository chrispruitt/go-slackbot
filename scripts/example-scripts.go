package scripts

import (
	"fmt"

	"github.com/chrispruitt/go-slackbot/lib/bot"
)

func init() {
	// Simple script
	bot.RegisterScript(bot.Script{
		Name:        "lulz",
		Matcher:     "lulz",
		Description: "lulz",
		Function: func(context *bot.ScriptContext) {
			bot.PostMessage(context.SlackEvent.Channel, "lol")
		},
	})

	// Script with parameters
	bot.RegisterScript(bot.Script{
		Name:        "echo",
		Matcher:     "echo <message>",
		Description: "Echo a message",
		Function: func(context *bot.ScriptContext) {
			message := context.Arguments["message"]
			bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("You said, \"%s\"", message))
		},
	})
}
