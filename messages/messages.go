package messages

import (
	"strings"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Text constructs a message with a single text content entry.
func Text(role xaiapiv1.MessageRole, text string) *xaiapiv1.Message {
	return &xaiapiv1.Message{
		Role: role,
		Content: []*xaiapiv1.Content{
			{Content: &xaiapiv1.Content_Text{Text: text}},
		},
	}
}

// UserText creates a ROLE_USER message.
func UserText(text string) *xaiapiv1.Message {
	return Text(xaiapiv1.MessageRole_ROLE_USER, text)
}

// AssistantText creates a ROLE_ASSISTANT message.
func AssistantText(text string) *xaiapiv1.Message {
	return Text(xaiapiv1.MessageRole_ROLE_ASSISTANT, text)
}

// SystemText creates a ROLE_SYSTEM message.
func SystemText(text string) *xaiapiv1.Message {
	return Text(xaiapiv1.MessageRole_ROLE_SYSTEM, text)
}

// ImageURL adds an image-url content block to the supplied message.
func ImageURL(msg *xaiapiv1.Message, url string, detail xaiapiv1.ImageDetail) *xaiapiv1.Message {
	if msg == nil {
		msg = &xaiapiv1.Message{}
	}
	msg.Content = append(msg.Content, &xaiapiv1.Content{
		Content: &xaiapiv1.Content_ImageUrl{
			ImageUrl: &xaiapiv1.ImageUrlContent{
				ImageUrl: strings.TrimSpace(url),
				Detail:   detail,
			},
		},
	})
	return msg
}
