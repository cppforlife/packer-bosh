package provisioner

import (
	"fmt"
	"os"

	"github.com/mitchellh/packer/packer"
)

type Provisioner struct {
	userConfig   UserConfig
	remoteConfig RemoteConfig
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	var err error

	p.userConfig, err = NewUserConfig(raws...)
	if err != nil {
		return fmt.Errorf("Building user config: %s", err)
	}

	err = p.userConfig.Validate()
	if err != nil {
		return fmt.Errorf("Validating user config: %s", err)
	}

	vmProvisionerConfig := VMProvisionerConfig{
		FullStemcellCompatibility: p.userConfig.FullStemcellCompatibility,

		AgentInfrastructure: p.userConfig.AgentInfrastructure,
		AgentPlatform:       p.userConfig.AgentPlatform,
		AgentConfiguration:  p.userConfig.AgentConfiguration,
	}

	localManifest := NewLocalManifest(p.userConfig.ManifestPath)

	assets := NewAssets(p.userConfig.AssetsDir)

	p.remoteConfig = NewRemoteConfig(
		"/opt/bosh-provisioner",
		vmProvisionerConfig,
		localManifest,
		p.userConfig.RemoteManifestPath,
		assets,
	)

	return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	cmds := NewSimpleCmds(p.userConfig.SudoCmd(), ui, comm)

	err := p.remoteConfig.Upload(cmds)
	if err != nil {
		return fmt.Errorf("Uploading config: %s", err)
	}

	err = cmds.ChmodX(p.remoteConfig.ExePath())
	if err != nil {
		return fmt.Errorf("Setting provisioner perms: %s", err)
	}

	err = cmds.RunPriv(p.buildCmd())
	if err != nil {
		return fmt.Errorf("Executing provisioner: %s", err)
	}

	return nil
}

func (p *Provisioner) Cancel() {
	os.Exit(0)
}

func (p *Provisioner) buildCmd() string {
	var stdErrRedirect string

	logPath := p.remoteConfig.ExeLogPath()

	if p.userConfig.IsDebug() {
		stdErrRedirect = fmt.Sprintf("2> >(tee %s >&2)", logPath)
	} else {
		stdErrRedirect = fmt.Sprintf("2>%s", logPath)
	}

	return fmt.Sprintf(
		"%s %s -configPath=%s %s",
		p.userConfig.SudoCmd(),
		p.remoteConfig.ExePath(),
		p.remoteConfig.ConfigPath(),
		stdErrRedirect,
	)
}
