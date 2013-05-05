package main

import (
	"bufio"
	"fmt"
	"irken/client"
	"irken/irc"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := irc.NewConn("irc.freenode.net")

	go func() {
		for {
			s, _ := conn.Read()
			s = strings.TrimSpace(strings.Replace(s, "\n", "", -1))
			line, err := client.ParseServerMsg(s)
			if err != nil {
				fmt.Println(s)
			} else {
				fmt.Println(line.Output())
			}
		}
	}()
	go func() {
		for {
			in := bufio.NewReader(os.Stdin)
			input, err := in.ReadString('\n')
			if err != nil {
				continue
			}
			conn.Write(input)
		}
	}()

	time.Sleep(time.Second * 3)
	err = conn.Write("USER testurnstf irken irken:Hejsan karlsson")
	err = conn.Write("NICK testurnstf")
	if err != nil {
		fmt.Println(err)
	}

	select {}

	//fmt.Println(client.ParseInData(":nyuszika7h!nyuszika7h@pdpc/supporter/active/nyuszika7h PRIVMSG #freenode :ppooooo :) sfdsldkfökakfa :(:(:(:(:(:(:(:())))))))"))
}
