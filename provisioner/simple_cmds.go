package provisioner

import (
	"fmt"
	"io"

	"github.com/mitchellh/packer/packer"
)

type SimpleCmds struct {
	sudoCmd string
	ui      packer.Ui
	comm    packer.Communicator
}

func NewSimpleCmds(sudoCmd string, ui packer.Ui, comm packer.Communicator) SimpleCmds {
	return SimpleCmds{sudoCmd: sudoCmd, ui: ui, comm: comm}
}

func (c SimpleCmds) Upload(dstPath string, input io.Reader) error {
	c.ui.Message(fmt.Sprintf("Uploading %s", dstPath))

	return c.comm.Upload(dstPath, input, nil)
}

func (c SimpleCmds) UploadDir(dstDir string, srcDir string, excl []string) error {
	c.ui.Message(fmt.Sprintf("Uploading %s", dstDir))

	return c.comm.UploadDir(dstDir, srcDir, excl)
}

func (c SimpleCmds) RunPriv(cmd string) error {
	return c.runCmdPriv(cmd)
}

// MkdirPNonPriv creates directory chowned to user of this SSH connection
func (c SimpleCmds) MkdirPNonPriv(dir string) error {
	err := c.runCmdPriv(fmt.Sprintf("mkdir -p %s", dir))
	if err != nil {
		return fmt.Errorf("Creating dir: %s", err)
	}

	err = c.runCmdPriv(fmt.Sprintf("chown `whoami`:`whoami` %s", dir))
	if err != nil {
		return fmt.Errorf("Chowning dir: %s", err)
	}

	return nil
}

func (c SimpleCmds) MkdirP(dir string) error {
	return c.runCmd(fmt.Sprintf("mkdir -p %s", dir))
}

func (c SimpleCmds) ChmodX(path string) error {
	return c.runCmd(fmt.Sprintf("chmod +x %s", path))
}

func (c SimpleCmds) runCmdPriv(cmd string) error {
	return c.runCmd(fmt.Sprintf("%s %s", c.sudoCmd, cmd))
}

func (c SimpleCmds) runCmd(cmd string) error {
	remoteCmd := &packer.RemoteCmd{Command: cmd}

	err := remoteCmd.StartWithUi(c.comm, c.ui)
	if err != nil {
		return fmt.Errorf("Starting command: %s", err)
	}

	if remoteCmd.ExitStatus != 0 {
		return fmt.Errorf("Non-zero exit status: %d", remoteCmd.ExitStatus)
	}

	return nil
}
