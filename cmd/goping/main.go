package main

import (
	"os"

	"github.com/ale-batt/goping"
)

func main() {
	if os.Args[1] != "" {
		goping.Ping(os.Args[1])
	}
}
