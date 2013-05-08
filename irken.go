package main

import (
	"irken/app"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	a := app.NewIrkenApp("config.cfg")
	a.GUI().StartMain()
}
