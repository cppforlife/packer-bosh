package provisioner

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

type UserConfig struct {
	tpl      *packer.ConfigTemplate
	metadata *mapstructure.Metadata

	// Used to determine if debugging is requested
	common.PackerConfig `mapstructure:",squash"`

	// AssetsDir is a directory on a host FS with bosh-agent assets
	AssetsDir string `mapstructure:"assets_dir"`

	// ManifestPath is a path on a host FS to a deployment manifest
	ManifestPath *string `mapstructure:"manifest_path"`

	// FullStemcellCompatibility makes provisioner install additional dependencies
	FullStemcellCompatibility bool `mapstructure:"full_stemcell_compatibility"`

	// Agent configuration
	AgentInfrastructure string                 `mapstructure:"agent_infrastructure"`
	AgentPlatform       string                 `mapstructure:"agent_platform"`
	AgentConfiguration  map[string]interface{} `mapstructure:"agent_configuration"`

	// SSHPassword is used to run sudo
	SSHPassword string `mapstructure:"ssh_password"`
}

func NewUserConfig(raws ...interface{}) (UserConfig, error) {
	c := UserConfig{
		AgentInfrastructure: "warden",
		AgentPlatform:       "ubuntu",
		AgentConfiguration:  map[string]interface{}{},
	}

	metadata, err := common.DecodeConfig(&c, raws...)
	if err != nil {
		return c, err
	}

	c.metadata = metadata

	return c, nil
}

func (c *UserConfig) HasManifestPath() bool { return c.ManifestPath != nil }

func (c *UserConfig) IsDebug() bool { return c.PackerDebug }

func (c *UserConfig) SudoCmd() string {
	// Must work even if SSHPassword is not set
	return fmt.Sprintf("echo %s | sudo -S", c.SSHPassword)
}

func (c *UserConfig) Validate() error {
	var err error

	c.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return err
	}

	c.tpl.UserVars = c.PackerUserVars

	errs := common.CheckUnusedConfig(c.metadata)

	err = c.validateDirPresence(c.AssetsDir, "assets_dir")
	if err != nil {
		errs = packer.MultiErrorAppend(errs, err)
	}

	if c.ManifestPath != nil {
		err = c.validateFilePresence(*c.ManifestPath, "manifest_path")
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (c *UserConfig) validateFilePresence(path string, optionName string) error {
	if path == "" {
		return fmt.Errorf("%s must be specified", optionName)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: %s is invalid: %s", optionName, path, err)
	} else if info.IsDir() {
		return fmt.Errorf("%s: %s must be a file", optionName, path)
	}

	return nil
}

func (c *UserConfig) validateDirPresence(path string, optionName string) error {
	if path == "" {
		return fmt.Errorf("%s must be specified", optionName)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: %s is invalid: %s", optionName, path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s: %s must be a directory", optionName, path)
	}

	return nil
}
