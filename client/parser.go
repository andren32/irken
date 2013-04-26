package client

import (
	"strings"
)

func ParseInData(data string) string {
	if strings.HasPrefix(data, ":") {
		senderEnd := strings.Index(data, " ")
		sender := data[1 : senderEnd-1]
		return sender
	}
	return ""
}

func ParseOutData(data string) string {
	return ""
}
