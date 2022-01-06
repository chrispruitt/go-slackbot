package bot

import (
	"regexp"
	"strings"
)

type Matcher string

func (m Matcher) getArguments(text string) (args map[string]string) {
	var re = regexp.MustCompile(m.toPerlSyntax())
	match := re.FindStringSubmatch(text)

	args = make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			args[name] = match[i]
		}
	}
	return args
}

// Matcher syntax to perl syntax
func (m Matcher) toPerlSyntax() (result string) {
	result = strings.ReplaceAll(string(m), "<", "(?P<")

	// Match any characters up until the next whitespace
	result = strings.ReplaceAll(result, ">", ">[^\\s]*)")

	// if ">+" is used (i.e. "<myarg>+" ) then spaces are allowed in the arg - mimicking the regex "one or more" syntax
	result = strings.ReplaceAll(result, ">[^\\s]*)+", ">.*)+")

	return result
}

// Matcher syntax to regex
func (m Matcher) toRegex() (result string) {
	re := regexp.MustCompile("<([^>]*)>")
	return re.ReplaceAllString(string(m), "(.*)")
}
