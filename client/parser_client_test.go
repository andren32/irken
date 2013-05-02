package client

import (
	"irken/test"
	"testing"
)

func TestLexClientCommand(t *testing.T) {
	input := "/command argument"
	l, err := lexClientMsg(input)
	if err != nil {
		t.Errorf("Should parse")
	}
	arg := l.args[0]
	cmd := l.cmd
	expArg := "argument"
	expCmd := "COMMAND"
	test.Check(t, arg, expArg)
	test.Check(t, cmd, expCmd)
}

func TestLexClientMessage(t *testing.T) {
	input := "Hallo from another world!"
	l, err := lexClientMsg(input)
	if err != nil {
		t.Errorf("Should parse")
	}
	arg := l.args[0]
	cmd := l.cmd
	expArg := "Hallo from another world!"
	expCmd := "PRIVMSG"
	test.Check(t, arg, expArg)
	test.Check(t, cmd, expCmd)
}

func TestLexClientMessageEscapeChar(t *testing.T) {
	input := "//Hallo from another world!"
	l, err := lexClientMsg(input)
	if err != nil {
		t.Errorf("Should parse")
	}
	arg := l.args[0]
	cmd := l.cmd
	expArg := "/Hallo from another world!"
	expCmd := "PRIVMSG"
	test.Check(t, arg, expArg)
	test.Check(t, cmd, expCmd)
}
