package client

import "irken/client/msg"
import "strings"

// IRCChannel is a channel on the irc
type IRCChannel struct {
	Ch    chan *msg.Line
	Nicks map[string]string // Key is the nick and the value is the permission: @, %, + or nothing.
}

// AddNicks adds a string of nicks to the nick collection. Indata is for example: nick1 nick2 @nick3
func (ircch *IRCChannel) AddNicks(nicks string) {
	nicks = strings.TrimSpace(nicks)
	nickArray := strings.Split(nicks, " ")
	for _, v := range nickArray {
		if string(v[0]) == "@" || string(v[0]) == "%" || string(v[0]) == "+" { // REGEXP? Other solution?
			ircch.Nicks[string(v[1:])] = string(v[0])
		} else {
			ircch.Nicks[v] = ""
		}
	}
}

// Removes the old nick and creates the new nick with the same permission
func (ircch *IRCChannel) ChangeNick(prevNick, newNick string) {
	perm := ircch.Nicks[prevNick]
	delete(ircch.Nicks, prevNick)
	ircch.Nicks[newNick] = perm
}

// RemoveNick removes the nick from the collection
func (ircch *IRCChannel) RemoveNick(nick string) {
	delete(ircch.Nicks, nick)
}

// AddNick changes the permission of the nick.
// Can also be used to change the nick
func (ircch *IRCChannel) AddNick(nick, perm string) {
	ircch.Nicks[nick] = perm
}