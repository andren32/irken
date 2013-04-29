package client

import (
	"irken/test"
	"testing"
)

func TestLexValid(t *testing.T) {
	message := ":prefix COMMAND param1 param2 :param 3 :-) yeah!?"
	prefix, command, params, err := LexIRC(message)
	if err != nil {
		t.Error(err)
	}
	test.Check(t, prefix, "prefix")
	test.Check(t, command, "COMMAND")
	test.Check(t, params[0], "param1")
	test.Check(t, params[1], "param2")
	test.Check(t, params[2], "param 3 :-) yeah!?")
}

func TestLexInValid(t *testing.T) {
	message := ":prefix"
	message2 := ":prefix "
	_, _, _, err := LexIRC(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
	_, _, _, err := LexIRC(message2)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}

}

func TestLexNoParams(t *testing.T) {
	message := "COMMAND"
	prefix, command, params, err := LexIRC(message)
	if err != nil {
		t.Error(err)
	}
	test.Check(t, prefix, "")
	test.Check(t, command, "COMMAND")
	if len(params) != 0 {
		t.Errorf("Reported fake parameters")
	}
}
