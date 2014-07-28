package provisioner

// Assets represents assets directory that can be uploaded to a remote location
type Assets struct {
	localDir string
}

func NewAssets(localDir string) Assets {
	return Assets{localDir: localDir}
}

func (a Assets) Upload(remoteDir string, cmds SimpleCmds) error {
	err := cmds.MkdirP(remoteDir)
	if err != nil {
		return err
	}

	// Make sure there is a trailing "/"
	// so that the directory isn't created on the other side.
	if a.localDir[len(a.localDir)-1] != '/' {
		a.localDir = a.localDir + "/"
	}

	return cmds.UploadDir(remoteDir, a.localDir, nil)
}
