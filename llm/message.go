package llm

import (
	"fmt"
	"strings"
)

const (
	RoleAssistant = "assistant"
	RoleUser      = "user"
	RoleSystem    = "system"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func UserMsg(content string) Message {
	return Message{Role: RoleUser, Content: content}
}

type MessageList []Message

func (ml MessageList) Latest() Message {
	if len(ml) == 0 {
		return Message{}
	}
	return ml[len(ml)-1]
}

func (ml MessageList) User() []Message {
	msgs := []Message{}
	for i := range ml {
		if ml[i].Role == RoleUser {
			msgs = append(msgs, ml[i])
		}
	}
	return msgs
}

func (ml MessageList) System() []Message {
	msgs := []Message{}
	for i := range ml {
		if ml[i].Role == RoleSystem {
			msgs = append(msgs, ml[i])
		}
	}
	return msgs
}

func (ml MessageList) Assistant() []Message {
	msgs := []Message{}
	for i := range ml {
		if ml[i].Role == RoleAssistant {
			msgs = append(msgs, ml[i])
		}
	}
	return msgs
}

func (ml MessageList) Formated() string {
	builder := strings.Builder{}
	for i := range ml {
		builder.WriteString(fmt.Sprintf("%s: %s\n", ml[i].Role, ml[i].Content))
	}
	return builder.String()
}
