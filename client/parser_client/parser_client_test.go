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
	expCmd := "CCOMMAND"
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

func TestLexClientMessageWithArgs(t *testing.T) {
	input := "/command arg1 arg2"
	l, err := lexClientMsg(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	arg1 := l.Args()[0]
	arg2 := l.Args()[1]
	cmd := l.Cmd()
	expArg1 := "arg1"
	expArg2 := "arg2"
	expCmd := "CCOMMAND"
	test.Check(t, arg1, expArg1)
	test.Check(t, arg2, expArg2)
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
	l, o, err := Parse(input, nick, context)
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
	l, o, err := Parse(input, nick, context)
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
	l, o, err := Parse(input, nick, context)
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
	input := "/part"
	nick := "user"
	context := "#chan"
	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "#chan"
	expOut := "PART #chan"
	expPr := "user has left #chan"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestClientQuit(t *testing.T) {
	input := "/quit Going to sleep"
	nick := "user"
	context := "#chan"
	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "#chan"
	expOut := "QUIT :Going to sleep"
	expPr := "user has quit (Going to sleep)"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestClientMsg(t *testing.T) {
	input := "/msg otheruser Hi there"
	nick := "user"
	context := "#chan"
	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := "otheruser"
	expOut := "PRIVMSG otheruser :Hi there"
	expPr := "user: Hi there"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestDisconnect(t *testing.T) {
	input := "/disconnect Bye bye"
	nick := "user"
	context := "#chan"

	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()

	expPr := "user disconnected"
	expOut := "QUIT :Bye bye"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
}

func TestServer(t *testing.T) {
	input := "/connect irc.freenode.net"
	nick := "user"
	context := ""

	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()
	cont := l.Context()

	expCont := ""
	expPr := "user connected to irc.freenode.net"
	expOut := ""
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
	test.Check(t, cont, expCont)
}

func TestPing(t *testing.T) {
	input := "/ping target"
	nick := "user"
	context := ""
	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()

	expPr := "user pinged target"
	expOut := "PING target"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)

	input = "/ping"
	l, o, err = Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr = l.OutputMsg()

	expPr = "/ping: You must supply a ping target!"
	expOut = ""
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
}

func TestSilentPing(t *testing.T) {
	input := "/silentping target"
	nick := "user"
	context := ""

	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()

	expPr := ""
	expOut := "PING target"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
}

func TestSilentPong(t *testing.T) {
	input := "/silentpong message"
	nick := "user"
	context := ""

	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()

	expPr := ""
	expOut := "PONG :message"
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)
}

func TestHelp(t *testing.T) {
	input := "/help arg"
	nick := "user"
	context := ""
	l, o, err := Parse(input, nick, context)
	if err != nil {
		test.UnExpErr(t, err)
	}
	pr := l.OutputMsg()

	expPr := ""
	expOut := ""
	test.Check(t, o, expOut)
	test.Check(t, pr, expPr)

}
