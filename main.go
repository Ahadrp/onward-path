package main

import (
	"onward-path/internal/mm"
	"onward-path/internal/xui"
)

func main() {
	mm.Load()
	mm.Run()

	xui.Login("root", "123")

	select {}
}
