package client

import (
	"irken/test"
	"testing"
)

func TestLexValid(t *testing.T) {
	message := ":prefix COMMAND param1 param2 :param 3 :-) yeah!?"
	prefix, command, params, err := lexMsg(message)
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
	_, _, _, err := lexMsg(message)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}
	_, _, _, err = lexMsg(message2)
	if err == nil {
		t.Errorf("Illegal message is not error reported")
	}

}

func TestLexNoParams(t *testing.T) {
	message := "COMMAND"
	prefix, command, params, err := lexMsg(message)
	if err != nil {
		t.Error(err)
	}
	test.Check(t, prefix, "")
	test.Check(t, command, "COMMAND")
	if len(params) != 0 {
		t.Errorf("Reported fake parameters")
	}
}

func TestJoin(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com JOIN #chan"
	msg, cont, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx has joined #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestMode(t *testing.T) {
	// TODO
}

func TestQuit(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com QUIT :Later suckerz"
	msg, cont, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx has quit (Later suckerz)"
	expCont := ""
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPart(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PART #chan"
	msg, cont, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx has left #chan"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestPrivMsg(t *testing.T) {
	input := ":_mrx!blabla@haxxor.com PRIVMSG #chan :Octotastic!"
	msg, cont, err := ParseServerMsg(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	expMsg := "_mrx: Octotastic!"
	expCont := "#chan"
	test.Check(t, msg, expMsg)
	test.Check(t, cont, expCont)
}

func TestResolveNick(t *testing.T) {
	input := "_mrx!blabla@haxxor.com"
	nick, err := resolveNick(input)
	if err != nil {
		t.Errorf("Should parse!")
	}
	exp := "_mrx"
	test.Check(t, nick, exp)
}

func TestResolveInvalidNick(t *testing.T) {
	input := "_mrxblabla@haxxor.com"
	_, err := resolveNick(input)
	if err == nil {
		t.Errorf("Should not parse")
	}
}
