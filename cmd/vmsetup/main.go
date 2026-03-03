package main

import (
	"fmt"
	"os"

	"github.com/minor-industries/vmsetup"
)

func main() {
	if err := vmsetup.Run(&vmsetup.Config{
		SshKeys:        nil,
		Username:       "ubuntu",
		CloudImageURL:  "https://cloud-images.ubuntu.com/noble/20260108/noble-server-cloudimg-amd64.img",
		CloudImageHash: "00786c0936a7dd91a6b07941ca60bb56652975e0e72f9dacf73c887ada420966",
	}); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
