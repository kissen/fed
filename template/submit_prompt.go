package template

import "math/rand"

var submitPrompts = []string{
	"What are you thinking about?",
	"What would you like to tell the world?",
	"What are you doing right now?",
	"What are you planning for tomorrow?",
}

func SubmitPrompt() string {
	totalChoices := int32(len(submitPrompts))
	idx := rand.Int31n(totalChoices)
	return submitPrompts[idx]
}
