package main

import (
	"github.com/barancanatbas/v2b/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
