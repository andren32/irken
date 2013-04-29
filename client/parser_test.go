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
	check(t, prefix, "prefix")
	check(t, command, "COMMAND")
	check(t, params[0], "param1")
	check(t, params[1], "param2")
	check(t, params[2], "param 3 :-) yeah!?")
}

func TestLexInValid(t *testing.T) {
	message := ":prefix"
	_, _, _, err := LexIRC(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
}

func check(t *testing.T, res, exp interface{}) {
	if mess, diff := test.Diff(res, exp); diff {
		t.Errorf("%s", mess)
	}

}
