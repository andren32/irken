package main

import (
	"bufio"
	"fmt"
	"irken/client"
	"os"
)

func main() {
	cfg := client.NewConfig("config.cfg")
	cfg.Load()
	cfgValues := cfg.GetCfgValues()

	nick := cfgValues["NICK"]
	realName := cfgValues["REAL_NAME"]

	cs, err := client.NewConnectSession("irc.freenode.net", nick, realName)
	if err != nil {
		fmt.Println(err)
	}

	cs.IrcChannels["#freenode"] = &client.IRCChannel{Ch: make(chan *client.Line)}
	cs.ReadToChannels()

	go func() {
		for {
			line := <-cs.IrcChannels[""].Ch
			fmt.Println(line.Raw())
			fmt.Print("SERVER WINDOW: ")
			fmt.Println(line.Output())
		}
	}()

	go func() {
		for {
			line := <-cs.IrcChannels["#freenode"].Ch
			fmt.Println(line.Raw())
			fmt.Print("FREENODE: ")
			fmt.Println(line.Output())
		}
	}()

	go func() {
		for {
			in := bufio.NewReader(os.Stdin)
			input, err := in.ReadString('\n')
			if err != nil {
				continue
			}
			cs.Conn.Write(input)
		}
	}()

	select {}
}
