package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Start() {
	if getenv("SHELL_MODE", false).(bool) {
		Shell()
	} else {
		DEBUG := getenv("DEBUG", false).(bool)

		webApi := slack.New(
			os.Getenv("SLACK_BOT_TOKEN"),
			slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
			slack.OptionDebug(DEBUG),
			slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		)
		socketMode := socketmode.New(
			webApi,
			socketmode.OptionDebug(false),
			socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
		)
		authTest, authTestErr := webApi.AuthTest()
		if authTestErr != nil {
			fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
			os.Exit(1)
		}
		selfUserId := authTest.UserID

		go func() {
			for envelope := range socketMode.Events {
				switch envelope.Type {
				case socketmode.EventTypeEventsAPI:
					// Events API:

					// Acknowledge the eventPayload first
					socketMode.Ack(*envelope.Request)

					eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
					switch eventPayload.Type {
					case slackevents.CallbackEvent:
						switch event := eventPayload.InnerEvent.Data.(type) {
						case *slackevents.MessageEvent:
							if event.User != selfUserId {
								HandleMessageEvent(event)
							}
						case *slackevents.AppMentionEvent:
							socketMode.Debugf("Skipped: %v", event)
						default:
							socketMode.Debugf("Skipped: %v", event)
						}
					default:
						socketMode.Debugf("unsupported Events API eventPayload received")
					}
				case socketmode.EventTypeInteractive:
					// Shortcuts:

					payload, _ := envelope.Data.(slack.InteractionCallback)
					switch payload.Type {
					case slack.InteractionTypeShortcut:
						if payload.CallbackID == "socket-mode-shortcut" {
							socketMode.Ack(*envelope.Request)
							modalView := slack.ModalViewRequest{
								Type:       "modal",
								CallbackID: "modal-id",
								Title: slack.NewTextBlockObject(
									"plain_text",
									"New Task",
									false,
									false,
								),
								Submit: slack.NewTextBlockObject(
									"plain_text",
									"Submit",
									false,
									false,
								),
								Close: slack.NewTextBlockObject(
									"plain_text",
									"Cancel",
									false,
									false,
								),
								Blocks: slack.Blocks{
									BlockSet: []slack.Block{
										slack.NewInputBlock(
											"input-task",
											slack.NewTextBlockObject(
												"plain_text",
												"Task Description",
												false,
												false,
											),
											// multiline is not yet supported
											slack.NewPlainTextInputBlockElement(
												slack.NewTextBlockObject(
													"plain_text",
													"Describe the task in detail with its timeline",
													false,
													false,
												),
												"input",
											),
										),
									},
								},
							}
							resp, err := webApi.OpenView(payload.TriggerID, modalView)
							if err != nil {
								log.Printf("Failed to opemn a modal: %v", err)
							}
							socketMode.Debugf("views.open response: %v", resp)
						}
					case slack.InteractionTypeViewSubmission:
						// View Submission:
						if payload.CallbackID == "modal-id" {
							socketMode.Debugf("Submitted Data: %v", payload.View.State.Values)
							socketMode.Ack(*envelope.Request)
						}
					default:
						// Others
						socketMode.Debugf("Skipped: %v", payload)
					}

				default:
					socketMode.Debugf("Skipped: %v", envelope.Type)
				}
			}
		}()

		socketMode.Run()
	}
}
