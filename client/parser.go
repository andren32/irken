package client

import (
	"strings"
)

// 
func ParseInData(data string) string {
	if strings.HasPrefix(data, ":") {
		senderEnd := strings.Index(data, " ")
		sender := data[1 : senderEnd-1]
		return sender
	}
	return ""
}

// LexIRC scans a IRC message and outputs its tokens.
func LexIRC(data string) (prefix string, command string, params []string,
	err error) {

	var prefixEnd int
	if strings.HasPrefix(data, ":") {
		prefixEnd := strings.Index(data, " ")
		prefix = data[1:prefixEnd]
	}

	var trailing string
	trailingStart := strings.Index(data, " :")
	if trailingStart >= 0 {
		trailing = data[trailingStart+2:]
	} else {
		trailingStart = len(data)
	}

	tmp := data[prefixEnd+1 : trailingStart-prefixEnd]
	cmdAndParams := strings.Fields(tmp)

	_ = cmdAndParams
	_ = trailing
	return

}

func ParseOutData(data string) string {
	return ""
}
