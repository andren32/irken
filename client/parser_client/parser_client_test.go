package parser_client

import (
	"irken/test"
	"testing"
)

func TestLexClientCommand(t *testing.T) {
	input := "/command argument"
	l, err := lexClientMsg(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	arg := l.Args()[0]
	cmd := l.Cmd()
	expArg := "argument"
	expCmd := "COMMAND"
	test.Check(t, arg, expArg)
	test.Check(t, cmd, expCmd)
}

func TestLexClientMessage(t *testing.T) {
	input := "Hallo from another world!"
	l, err := lexClientMsg(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	arg := l.Args()[0]
	cmd := l.Cmd()
	expArg := "Hallo from another world!"
	expCmd := "CHAN"
	test.Check(t, arg, expArg)
	test.Check(t, cmd, expCmd)
}

func TestLexClientMessageEscapeChar(t *testing.T) {
	input := "\\/Hallo from another world!"
	l, err := lexClientMsg(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	arg := l.Args()[0]
	cmd := l.Cmd()
	expArg := "/Hallo from another world!"
	expCmd := "CHAN"
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

func TestClientChanMsg(t *testing.T) {
	input := "testing testing 123"
	nick := "user"
	context := "#chan"
	l, o, err := parseClientMsg(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "#chan"
	expOut := "PRIVMSG #chan :testing testing 123"
	expPr := "user: testing testing 123"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestClientJoinChan(t *testing.T) {
	input := "/join #chan"
	nick := "user"
	context := ""
	l, o, err := parseClientMsg(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "#chan"
	expOut := "JOIN #chan"
	expPr := "user joined #chan"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestClientMe(t *testing.T) {
	input := "/me is testing IRC"
	nick := "user"
	context := "#chan"
	l, o, err := parseClientMsg(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "#chan"
	expOut := "PRIVMSG #chan :user is testing IRC"
	expPr := "user is testing IRC"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
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
