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
	input := "\\/Hallo from another world!"
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

func TestLexClientInvalidMessage(t *testing.T) {
	input := "/*^Hallo from another world!"
	_, err := lexClientMsg(input)
	if err == nil {
		t.Errorf("Should not parse")
	}
}

func TestClientJoinChan(t *testing.T) {
	input := "/join #chan"
	nick := "user"
	context := ""
	l, o, err := parseClientMsg(input, nick, context)
	if err != nil {
		t.Errorf("Should parse")
	}
	pr := l.output

	expOut := "JOIN #chan"
	expPr := "user joined #chan"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
}

func TestClientMe(t *testing.T) {
	// TODO
}

func TestClientHelp(t *testing.T) {
	// TODO
}

func TestClientPart(t *testing.T) {
	// TODO
}

func TestClientQuit(t *testing.T) {
	// TODO
}

func TestClientMsg(t *testing.T) {
	// TODO
}
