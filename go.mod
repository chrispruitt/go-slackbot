module github.com/chrispruitt/go-slackbot

go 1.16

replace github.com/chrispruitt/go-slackbot/lib/bot => ./bot

require (
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/mattes/go-asciibot v0.0.0-20190603170252-3fa6d766c482
	github.com/pkg/errors v0.9.1 // indirect
	github.com/slack-go/slack v0.9.4
	github.com/stretchr/testify v1.6.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
