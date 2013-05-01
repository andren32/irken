package client

import (
	"irken/test"
	"testing"
)

func TestLexValid(t *testing.T) {
	message := ":prefix COMMAND param1 param2 :param 3 :-) yeah!?"
	l, err := lexMsg(message)
	if err != nil {
		t.Error(err)
	}
	prefix := l.Src
	command := l.Cmd
	params := l.Args

	test.Check(t, prefix, "prefix")
	test.Check(t, command, "COMMAND")
	test.Check(t, params[0], "param1")
	test.Check(t, params[1], "param2")
	test.Check(t, params[2], "param 3 :-) yeah!?")
}

func TestLexInValid(t *testing.T) {
	message := ":prefix"
	message2 := ":prefix "
	_, err := lexMsg(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
	_, err = lexMsg(message2)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}

}

func TestLexNoParams(t *testing.T) {
	message := "COMMAND"
	_, err := lexMsg(message)
	if err == nil {
		t.Errorf("No prefix - should not parse")
	}
}

func TestJoin(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com JOIN #chan"
	l, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	msg := l.output
	cont := l.Context
	expMsg := "_mrx has joined #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestMode(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com MODE #chan +i -l"
	msg, cont, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx changed mode +i -l for #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestQuit(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com QUIT :Later suckerz"
	l, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	msg := l.output
	cont := l.Context
	expMsg := "_mrx has quit (Later suckerz)"
	expCont := ""
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPart(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PART #chan"
	l, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	msg := l.output
	cont := l.Context
	expMsg := "_mrx has left #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPrivMsg(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PRIVMSG #chan :Octotastic! I like pie btw :)"
	l, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	msg := l.output
	cont := l.Context
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx: Octotastic! I like pie btw :)"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestResolveValidPrefix(t *testing.T) {
	input := "_mrx!blabla@haxxor.com"
	nick, ident, host, src, err := resolvePrefix(input)
	expNick := "_mrx"
	expIdent := "blabla"
	expHost := "haxxor.com"
	expSrc := "_mrx!blabla@haxxor.com"
	if err != nil {
		t.Errorf("Should parse")
	}

	test.Check(t, nick, expNick)
	test.Check(t, ident, expIdent)
	test.Check(t, host, expHost)
	test.Check(t, src, expSrc)
}

func TestResolveInvalidPrefix(t *testing.T) {
	input := ""
	_, _, _, _, err := resolvePrefix(input)
	if err == nil {
		t.Errorf("Should not parse")
	}
}

func TestResolveServer(t *testing.T) {
	input := "blabla.haxxor.com"
	nick, ident, host, src, err := resolvePrefix(input)
	expNick := "blabla.haxxor.com"
	expIdent := ""
	expHost := ""
	expSrc := "blabla.haxxor.com"
	if err != nil {
		t.Errorf("Should parse")
	}

	test.Check(t, nick, expNick)
	test.Check(t, ident, expIdent)
	test.Check(t, host, expHost)
	test.Check(t, src, expSrc)
}
