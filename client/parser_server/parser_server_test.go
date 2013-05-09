package parser_server

import (
	"irken/test"
	"testing"
)

func TestLexValidServerMsg(t *testing.T) {
	message := ":prefix COMMAND param1 param2 :param 3 :-) yeah!?"
	l, err := lexServerMsg(message)
	if err != nil {
		test.UnExpErr(t, err)
	}
	prefix := l.Src()
	command := l.Cmd()
	params := l.Args()

	test.Check(t, prefix, "prefix")
	test.Check(t, command, "COMMAND")
	test.Check(t, params[0], "param1")
	test.Check(t, params[1], "param2")
	test.Check(t, params[2], "param 3 :-) yeah!?")
}

func TestLexInValidServerMsg(t *testing.T) {
	message := ":prefix"
	message2 := ":prefix "
	_, err := lexServerMsg(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
	_, err = lexServerMsg(message2)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}

}

func TestLexNoParams(t *testing.T) {
	message := "COMMAND"
	_, err := lexServerMsg(message)
	if err != nil {
		test.UnExpErr(t, err)
	}
}

func TestJoin(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com JOIN #chan"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "_mrx has joined #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestMode(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com MODE #chan +i -l"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "_mrx changed mode +i -l for #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestNotice(t *testing.T) {
	input := ":blabla.haxxor.com NOTICE * :Welcome"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "blabla.haxxor.com: Welcome"
	expCont := "*"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestNoticeNoPrefix(t *testing.T) {
	input := "NOTICE * :Welcome"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "<Server>: Welcome"
	expCont := "*"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestQuit(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com QUIT :Later suckerz"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "_mrx has quit (Later suckerz)"
	expCont := ""
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPart(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PART #chan"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "_mrx has left #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPrivMsg(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PRIVMSG #chan :Octotastic! I like pie btw :)"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	if err != nil {
		test.UnExpErr(t, err)
	}
	expMsg := "_mrx: Octotastic! I like pie btw :)"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)

	input = ":_mrx!blabla@haxxor.com PRIVMSG user :Octotastic! I like pie btw :)"
	l, err = Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg = l.OutputMsg()
	cont = l.Context()
	cmd := l.Cmd()
	if err != nil {
		test.UnExpErr(t, err)
	}
	expCmd := "P2PMSG"
	expMsg = "_mrx: Octotastic! I like pie btw :)"
	expCont = "_mrx"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
	test.Check(t, cmd, expCmd)
}

func TestNick(t *testing.T) {
	input := ":WiZ!jto@tolsun.oulu.fi NICK Kilroy"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	expMsg := "WiZ changed nick to Kilroy"
	test.Check(t, msg, expMsg)
}

func TestTopic(t *testing.T) {
	input := ":blabla.haxxor.com 332 axelri #chan :Welcome to chan!"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "Topic for #chan is \"Welcome to chan!\""
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestTopicSetBy(t *testing.T) {
	input := ":blabla.haxxor.com 333 user #chan marienz 1365217959"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "Topic set by marienz on Sat, 06 Apr 2013 03:12:39 UTC"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestChannelURL(t *testing.T) {
	input := ":services. 328 user #chan :http://chan.org/"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "URL for #chan: http://chan.org/"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestChannelCreation(t *testing.T) {
	input := ":blabla.haxxor.com 329 user #chan 981760584"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expMsg := "Channel created on Fri, 09 Feb 2001 23:16:24 UTC"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestNicks(t *testing.T) {
	input := ":blabla.haxxor.com 353 user = #chan :user1 user2"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expCont := "#chan"
	expMsg := "user1 user2"
	test.Check(t, cont, expCont)
	test.Check(t, msg, expMsg)
}

func TestPing(t *testing.T) {
	input := "PING :arg"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expCont := ""
	expMsg := "Pinged: arg"
	test.Check(t, cont, expCont)
	test.Check(t, msg, expMsg)
}

func TestNoSuchTarget(t *testing.T) {
	input := ":blabla.haxxor.com 401 user somenick :No such nick/channel"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expCont := "somenick"
	expMsg := "somenick - No such nick/channel"
	test.Check(t, cont, expCont)
	test.Check(t, msg, expMsg)
}

func TestAction(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PRIVMSG #chan :\001ACTION is tired\001"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cmd := l.Cmd()
	cont := l.Context()
	expCmd := "PRIVMSG"
	expCont := "#chan"
	expMsg := "*_mrx* is tired"
	test.Check(t, cont, expCont)
	test.Check(t, msg, expMsg)
	test.Check(t, cmd, expCmd)
}

func TestForwarding(t *testing.T) {
	input := ":blabla.haxxor.com 470 user #oldchan #newchan :Forwarding to new channel"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	cont := l.Context()
	expCont := "#newchan"
	expMsg := "#oldchan --> #newchan: Forwarding to new channel"
	test.Check(t, cont, expCont)
	test.Check(t, msg, expMsg)
}

func TestNumeric(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com 008 user :Something in the beginning"
	l, err := Parse(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	msg := l.OutputMsg()
	expMsg := "Something in the beginning"
	test.Check(t, msg, expMsg)
}

func TestResolveValidPrefix(t *testing.T) {
	input := "_mrx!blabla@haxxor.com"
	nick, ident, host, src, err := resolvePrefix(input)
	expNick := "_mrx"
	expIdent := "blabla"
	expHost := "haxxor.com"
	expSrc := "_mrx!blabla@haxxor.com"
	if err != nil {
		test.UnExpErr(t, err)
	}

	test.Check(t, nick, expNick)
	test.Check(t, ident, expIdent)
	test.Check(t, host, expHost)
	test.Check(t, src, expSrc)
}

func TestResolveEmptyPrefix(t *testing.T) {
	input := ""
	nick, ident, host, src, err := resolvePrefix(input)
	if err != nil {
		test.UnExpErr(t, err)
	}
	expNick := "<Server>"
	expIdent := ""
	expHost := ""
	expSrc := ""

	test.Check(t, nick, expNick)
	test.Check(t, ident, expIdent)
	test.Check(t, host, expHost)
	test.Check(t, src, expSrc)
}

func TestResolveServer(t *testing.T) {
	input := "blabla.haxxor.com"
	nick, ident, host, src, err := resolvePrefix(input)
	expNick := "blabla.haxxor.com"
	expIdent := ""
	expHost := ""
	expSrc := "blabla.haxxor.com"
	if err != nil {
		test.UnExpErr(t, err)
	}

	test.Check(t, nick, expNick)
	test.Check(t, ident, expIdent)
	test.Check(t, host, expHost)
	test.Check(t, src, expSrc)
}
