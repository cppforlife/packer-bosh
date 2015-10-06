package provisioner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
)

type RemoteConfig struct {
	vmProvisionerConfig VMProvisionerConfig

	// Deployment manifest and its path
	localManifest LocalManifest
	manifestPath  *string

	// bosh-provisioner assets (e.g. bosh-agent, monit)
	assets     Assets
	configPath string
	exePath    string
	exeLogPath string

	// Usually /opt/bosh-provisioner
	baseDir   string
	assetsDir string
	reposDir  string
	tmpDir    string

	localBlobstoreDir string
}

type VMProvisionerConfig struct {
	FullStemcellCompatibility bool

	AgentInfrastructure string
	AgentPlatform       string
	AgentConfiguration  map[string]interface{}
}

func NewRemoteConfig(
	baseDir string,
	vmProvisionerConfig VMProvisionerConfig,
	localManifest LocalManifest,
	remoteManifestPath string,
	assets Assets,
) RemoteConfig {
	assetsDir := filepath.Join(baseDir, "assets")

	c := RemoteConfig{
		vmProvisionerConfig: vmProvisionerConfig,

		localManifest: localManifest,

		assets:     assets,
		configPath: filepath.Join(baseDir, "config.json"),
		exePath:    filepath.Join(assetsDir, "bosh-provisioner"),
		exeLogPath: filepath.Join(baseDir, "provisioner.log"),

		baseDir:   baseDir,
		assetsDir: assetsDir,
		reposDir:  filepath.Join(baseDir, "repos"),
		tmpDir:    filepath.Join(baseDir, "tmp"),

		localBlobstoreDir: filepath.Join(baseDir, "blobstore"),
	}

	if localManifest.IsPresent() {
		manifestPath := filepath.Join(baseDir, "manifest.yml")
		c.manifestPath = &manifestPath
	}

	if remoteManifestPath != "" {
		c.manifestPath = &remoteManifestPath
	}

	return c
}

func (c RemoteConfig) ConfigPath() string { return c.configPath }
func (c RemoteConfig) ExePath() string    { return c.exePath }
func (c RemoteConfig) ExeLogPath() string { return c.exeLogPath }

func (c RemoteConfig) Upload(cmds SimpleCmds) error {
	// Create base directory for non-privileged user so that upload can succeed
	err := cmds.MkdirPNonPriv(c.baseDir)
	if err != nil {
		return fmt.Errorf("Creating base dir: %s", err)
	}

	err = c.localManifest.Upload(c.manifestPath, cmds)
	if err != nil {
		return fmt.Errorf("Uploading manifest: %s", err)
	}

	err = c.assets.Upload(c.assetsDir, cmds)
	if err != nil {
		return fmt.Errorf("Uploading assets: %s", err)
	}

	dstPath := c.configPath

	config, err := c.build()
	if err != nil {
		return fmt.Errorf("Building config: %s", err)
	}

	err = cmds.Upload(dstPath, config)
	if err != nil {
		return fmt.Errorf("Uploading config %s: %s", dstPath, err)
	}

	return nil
}

func (c RemoteConfig) build() (io.Reader, error) {
	type h map[string]interface{}

	config := h{
		"assets_dir": c.assetsDir,
		"repos_dir":  c.reposDir,
		"tmp_dir":    c.tmpDir,

		"event_log": h{
			"device_type": "text",
		},

		"blobstore": h{
			"provider": "local",
			"options": h{
				"blobstore_path": c.localBlobstoreDir,
			},
		},

		"vm_provisioner": h{
			"full_stemcell_compatibility": c.vmProvisionerConfig.FullStemcellCompatibility,

			"agent_provisioner": h{
				"infrastructure": c.vmProvisionerConfig.AgentInfrastructure,
				"platform":       c.vmProvisionerConfig.AgentPlatform,
				"configuration":  c.vmProvisionerConfig.AgentConfiguration,

				"mbus": "https://user:password@127.0.0.1:4321/agent",
			},
		},

		"deployment_provisioner": h{
			"manifest_path": c.manifestPath,
		},
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("Marshalling provisioner config: %s", err)
	}

	return bytes.NewBuffer(configBytes), nil
}
