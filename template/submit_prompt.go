package template

import "math/rand"

var submitPrompts = []string{
	"What are you thinking about?",
	"What would you like to tell the world?",
	"What are you doing right now?",
	"What are you planning for tomorrow?",
}

// Return a random string from submitPrompts that prompts
// the user to enter something fun.
//
// It's used at the placeholder string in the submit text
// field.
func submitPrompt() string {
	totalChoices := int32(len(submitPrompts))
	idx := rand.Int31n(totalChoices)
	return submitPrompts[idx]
}
