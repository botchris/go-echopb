package shared

import (
	"regexp"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
)

// loremRegex is the regex to match @lorem(<number>) patterns.
// For example: @lorem(5) will be replaced with 5 lorem ipsum words randomly generated.
var loremRegex = regexp.MustCompile(`@lorem\((\d+)\)`)

// MessageGenerator generates messages based on a template.
// It supports generating lorem ipsum text by using the `@lorem(<number>)` pattern.
// Otherwise, it returns the template as is.
type MessageGenerator struct {
	template   string
	loremWords int
}

// NewMessageGenerator creates a new MessageGenerator based on the provided template.
func NewMessageGenerator(template string) (*MessageGenerator, error) {
	g := &MessageGenerator{
		template: template,
	}

	matches := loremRegex.FindStringSubmatch(template)
	if len(matches) == 2 {
		lw, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		g.loremWords = lw
	}

	return g, nil
}

// Get generates and returns the message based on the template.
func (mg *MessageGenerator) Get() string {
	if mg.loremWords > 0 {
		return gofakeit.LoremIpsumSentence(mg.loremWords)
	}

	return mg.template
}
