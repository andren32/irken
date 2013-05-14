package sock

import "testing"

func TestNewConn(t *testing.T) {
	_, err := NewConn("irc.freenode.net", "le_tester1234")

	if err != nil {
		t.Error(err)
	}
}
