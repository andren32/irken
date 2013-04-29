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
	check(t, "prefix", prefix)
	check(t, "COMMAND", command)
	check(t, "param1", params[0])
	check(t, "param2", params[1])
	check(t, "param 3 :-) yeah!?", params[2])
}

func TestLexInValid(t *testing.T) {
	message := ":prefix"
	_, _, _, err := LexIRC(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
}

func check(t *testing.T, exp, res interface{}) {
	if mess, diff := test.Diff(exp, res); diff {
		t.Errorf("%s", mess)
	}

}
