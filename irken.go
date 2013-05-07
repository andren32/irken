package main

import (
	"fmt"
	"irken/app"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	a := app.NewIrkenApp("../config.cfg")
	go func() {
		a.GUI().StartMain()
	}()

	fmt.Println("Got here!")
	select {}
}
