package bot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	c "github.com/chrispruitt/go-slackbot/lib/config"
	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	SlackClient     *slack.Client
	config          *Config
	scripts         map[string]Script
	periodicScripts map[string]PeriodicScript
	cronJobs        *cron.Cron
)

type ScriptFunction func(*ScriptContext)

type Script struct {
	Name        string
	Matcher     Matcher
	Description string
	Function    ScriptFunction
}

func init() {
	SlackClient = slack.New(
		c.SlackBotToken,
		slack.OptionAppLevelToken(c.SlackAppToken),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
}

func RegisterScript(script Script) {
	if scripts == nil {
		scripts = make(map[string]Script)
	}
	scripts[script.Name] = script
}

func PostMessage(channelID string, message string) (string, string, error) {
	if c.ShellMode && c.ShellModeChannel == "" {
		fmt.Println(message)
		return "", "", nil
	} else {
		return SlackClient.PostMessage(channelID, slack.MsgOptionText(message, false))
	}
}

func HandleMessageEvent(event *slackevents.MessageEvent) {

	if strings.HasPrefix(event.Text, fmt.Sprintf("%s ", c.BotName)) {

		// Strip out BotName
		re := regexp.MustCompile(fmt.Sprintf(`^%s *`, c.BotName))
		event.Text = re.ReplaceAllString(event.Text, "")

		for _, script := range scripts {
			if match(script.Matcher.toRegex(), event.Text) {

				ScriptContext := &ScriptContext{
					SlackEvent: event,
				}

				ScriptContext.Arguments = script.Matcher.getArguments(event.Text)

				// Execute scripts asynchronously so bot can keep listening
				go script.Function(ScriptContext)

				return
			}
		}

		PostMessage(event.Channel, "Sorry, I don't know that command.")
	}
}

func match(matcher string, content string) bool {
	re := regexp.MustCompile(matcher)
	return re.MatchString(content)
}
