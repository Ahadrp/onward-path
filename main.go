package main

import (
	"onward-path/internal/mm"
)

func main() {
	mm.Load()
	mm.Run()

	select {}
}
