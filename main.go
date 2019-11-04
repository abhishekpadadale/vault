package main // import "github.com/abhishekpadadale/vault"

import (
	"os"

	"github.com/abhishekpadadale/vault/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
