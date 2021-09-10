**Description**

This is a simple hubot like slack bot. It will listen to slack MessageEvents prefixed with a given BOT_NAME and parse the message to run your custom scripts.

**Interactive shell for testing message events**

[![asciicast](https://asciinema.org/a/433605.svg)](https://asciinema.org/a/433605)

**Setup**

1. Create a slack app
1. Set app to socket mode
1. Get/Generate your App Token and Bot Token
1. Enable Events and Subscribe to `message.channels` events

```go
package main

import (
	"fmt"
	"strings"

	"github.com/chrispruitt/go-slackbot/lib/bot"
	c "github.com/chrispruitt/go-slackbot/lib/config"
)

func main() {
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

	// Script with some custom arguments syntax
	bot.RegisterScript(bot.Script{
		Name:        "ship",
		Matcher:     `ship <app> to <env>`,
		Description: "Usage: 'ship app1@v1.0.0 to dev' or 'ship app1@v1.0.0 app2@v1.0.0 to dev",
		Function: func(context *bot.ScriptContext) {

			env := context.Arguments["env"]
			apps := strings.Split(context.Arguments["app"], " ")

			for _, a := range apps {
				app := strings.Split(a, "@")
				bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("Shipping App: %s Version: %s to %s", app[0], app[1], env))
			}
		},
	})

	botConfig := &bot.Config{
		SlackAppToken:    c.SlackAppToken,
		SlackBotToken:    c.SlackBotToken,
		BotName:          c.BotName,

		 // Set to true to enable shell mode
		ShellMode:        c.ShellMode,
		
		// If set, when using shell mode, scripts will bot.PostMessage will post message in given slack channel
		ShellModeChannel: c.ShellModeChannel,
	}

	bot.Start(botConfig)
}
```


**In Slack**

Add your slack bot to a channel.

Then, execute a script via slack by typing the given BOT_NAME followed by a command that will match a script matcher.

`${BOT_NAME} help <filter>` is a built in script that will list all your commands using the Description and Matcher fields.

**Build in commands**

`${BOT_NAME} help ?<filter>` - filter arg is optional - list all available commands
`${BOT_NAME} status` - check status of bot

**Roadmap**

- Provide terraform module for quick setup in fargate
- Update readme with a "how to" to set up slack bot or publish one
- Add native script authorization via roles
- Give go-slackbot a brain via dynamodb or s3 json file
- Catch fatal script errors and prevent exit
- Dockerize
