package client

import (
	"errors"
	"strings"
)

// LexIRC scans a IRC message and outputs its tokens.
func LexIRC(data string) (prefix, command string, params []string, err error) {

	// grab prefix if present
	prefixEnd := -1
	if strings.HasPrefix(data, ":") {
		prefixEnd = strings.Index(data, " ")
		if prefixEnd == -1 {
			err = errors.New("Message with only a prefix")
			return
		}
		prefix = data[1:prefixEnd]
	}

	// grab trailing param if present
	var trailing string
	trailingStart := strings.Index(data, " :")
	if trailingStart >= 0 {
		trailing = data[trailingStart+2:]
	} else {
		trailingStart = len(data)
	}

	tmp := data[prefixEnd+1 : trailingStart]
	cmdAndParams := strings.Fields(tmp)
	if len(cmdAndParams) < 1 {
		err = errors.New("Cannot lex command")
		return
	}

	command = cmdAndParams[0]
	params = cmdAndParams[1:]
	if trailing != "" {
		params = append(params, trailing)
	}

	return

}
