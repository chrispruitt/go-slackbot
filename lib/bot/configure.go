package bot

type Config struct {
	SlackBotToken    string
	SlackAppToken    string
	BotName          string
	Debug            bool
	ShellMode        bool
	ShellModeChannel string
}
