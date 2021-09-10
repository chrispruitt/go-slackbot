package bot

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
)

var botStartDateTime time.Time

func init() {

	botStartDateTime = time.Now().UTC()

	RegisterScript(Script{
		Name:        "Status",
		Matcher:     "status",
		Description: "Display status.",
		Function:    statusScript,
	})
}

func statusScript(context *ScriptContext) {
	duration := time.Now().UTC().Sub(botStartDateTime)
	PostMessage(context.SlackEvent.Channel, fmt.Sprintf(":wave: Hi, I've been running for %s", durafmt.ParseShort(duration)))
}
