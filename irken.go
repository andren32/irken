package main

import (
	"irken/app"
	"irken/handlers"
	"log"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile)
}

func main() {
	a := app.NewIrkenApp("config.cfg")
	handlers.Init(a)
	a.GUI().StartMain()
}
