package provisioner

import (
	"fmt"
	"os"
	"path/filepath"
)

// LocalManifest represent a deployment manifest that can be uploaded to a remote location
type LocalManifest struct {
	localPath *string
}

func NewLocalManifest(localPath *string) LocalManifest {
	return LocalManifest{localPath: localPath}
}

func (m LocalManifest) IsPresent() bool {
	return m.localPath != nil
}

func (m LocalManifest) Upload(remotePath *string, cmds SimpleCmds) error {
	if remotePath != nil && m.IsPresent() {
		return m.uploadFile(cmds, *remotePath, *m.localPath)
	}

	return nil
}

func (m LocalManifest) uploadFile(cmds SimpleCmds, dstPath, srcPath string) error {
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
