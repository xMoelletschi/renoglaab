package main

import (
	"os"

	"github.com/xMoelletschi/renoglaab/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
