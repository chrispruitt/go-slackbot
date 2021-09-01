package bot

import "github.com/slack-go/slack/slackevents"

type ScriptContext struct {
	Arguments  map[string]string
	SlackEvent *slackevents.MessageEvent
}
