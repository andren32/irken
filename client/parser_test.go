package client

import "testing"

const message = ":prefix COMMAND param1 param2 :param 3 :-) yeah!?"

func TestLexValidMessage(t *testing.T) {
	prefix, command, params, err := LexIRC(message)
	if err != nil {
		t.Error(err)
	}
	check(t, "prefix", prefix)
	check(t, "COMMAND", command)
	check(t, "param1", params[0])
	check(t, "param2", params[1])
	check(t, "param 3 :-) yeah!?", params[2])
}

func check(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
