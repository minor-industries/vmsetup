package main

import (
	"fmt"
	"os"

	"github.com/minor-industries/vmsetup"
)

func main() {
	if err := vmsetup.Run(nil, "ubuntu"); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
