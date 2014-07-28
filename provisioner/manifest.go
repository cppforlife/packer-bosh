package provisioner

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manifest represent a deployment manifest that can be uploaded to a remote location
type Manifest struct {
	localPath *string
}

func NewManifest(localPath *string) Manifest {
	return Manifest{localPath: localPath}
}

func (m Manifest) IsPresent() bool {
	return m.localPath != nil
}

func (m Manifest) Upload(remotePath *string, cmds SimpleCmds) error {
	if remotePath != nil && m.IsPresent() {
		return m.uploadFile(cmds, *remotePath, *m.localPath)
	}

	return nil
}

func (m Manifest) uploadFile(cmds SimpleCmds, dstPath, srcPath string) error {
	err := cmds.MkdirP(filepath.Dir(dstPath))
	if err != nil {
		return err
	}

	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("Opening %s: %s", srcPath, err)
	}

	defer f.Close()

	err = cmds.Upload(dstPath, f)
	if err != nil {
		return fmt.Errorf("Uploading %s: %s", dstPath, err)
	}

	return nil
}
