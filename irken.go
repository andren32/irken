package main

import (
	"irken/frontend/app"
	"irken/frontend/handlers"
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
