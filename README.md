**Description**

This is a simple hubot like slack bot. It will listen to slack MessageEvents prefixed with a given BOT_NAME and parse the message to run your custom scripts.

**Interactive shell for testing message events**

[![asciicast](https://asciinema.org/a/433605.svg)](https://asciinema.org/a/433605)

**Setup**

1. Create a slack app
1. Set app to socket mode
1. Get/Generate your App Token and Bot Token
1. Enable Events and Subscribe to `message.channels` events
1. Create a golang project with a main.go file and a scripts/ directory with the below example files

**Environment Variables**

Set Environment Vars
```bash
export SLACK_APP_TOKEN=xapp-blahblah
export SLACK_BOT_TOKEN=xoxb-blahblah

# The Script Matcher will match against messages prefixed this value
export BOT_NAME=bender2

# If set, when using shell mode, bot.PostMessage will post to this slack channel instead of the console
# export SHELL_MODE_CHANNEL=SLACKCHANNELID

# Enable debug mode to log every event
export DEBUG=true
```

**Shell Mode**
```bash
SHELL_MODE=true go run main.go
```

main.go example

```go
package main

import (
	"os"

	"github.com/chrispruitt/go-slackbot/bot"
	_ "<yourModuleName>/scripts"
)

func main() {
	bot.Start()
}
```

scripts/exampleScript.go

```go
package scripts

import (
	"fmt"
	"regexp"

	"github.com/chrispruitt/go-slackbot/bot"

	"github.com/slack-go/slack/slackevents"
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

	// Script with some custom parameter syntax
	bot.RegisterScript(bot.Script{
		Name:        "ship",
		Matcher:     `ship <app> to <env>`,
		Description: "Usage: 'ship app1@v1.0.0 to dev' or 'ship app1@v1.0.0 app2@v1.0.0 to dev",
		Function: func(context *bot.ScriptContext) {
			// TODO Validation
			env := context.Arguments["env"]
			apps := strings.Split(context.Arguments["app"], " ")

			for _, a := range apps {
				app := strings.Split(a, "@")
				bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("Shipping App: %s Version: %s to %s", app[0], app[1], env))
			}
		},
	})
}
```

**In Slack**

Add your slack bot to a channel.

Then, execute a script via slack by typing the given BOT_NAME followed by a command that will match a script matcher.

`${BOT_NAME} help` is a built in script that will list all your commands using the Description and Matcher fields.

**Roadmap**

- Provide terraform module for quick setup in fargate
- Update readme with a "how to" to set up slack bot or publish one
- Add native script authorization via roles
- Give go-slackbot a brain via dynamodb or s3 json file
- Catch fatal script errors and prevent exit
- Dockerize
