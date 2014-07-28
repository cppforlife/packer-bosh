package main

import (
	"github.com/mitchellh/packer/packer/plugin"

	pbprov "github.com/cppforlife/packer-bosh/provisioner"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterProvisioner(new(pbprov.Provisioner))

	server.Serve()
}
