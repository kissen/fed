package ap

import (
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"gitlab.cs.fau.de/kissen/fed/db"
)

func convertPostToNote(post *db.FedPost) vocab.ActivityStreamsNote {
	note := streams.NewActivityStreamsNote()

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString("Note")
	note.SetActivityStreamsName(name)

	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString(post.Content)
	note.SetActivityStreamsContent(content)

	return note
}

func convertPostsToNotes(posts []*db.FedPost) vocab.ActivityStreamsOrderedItemsProperty {
	oi := streams.NewActivityStreamsOrderedItemsProperty()

	for _, post := range posts {
		note := convertPostToNote(post)
		oi.AppendActivityStreamsNote(note)
	}

	return oi
}
