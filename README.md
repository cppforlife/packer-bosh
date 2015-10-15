## BOSH provisioner for Packer

BOSH provisioner allows to provision VM by specifying regular BOSH deployment manifest.

You might also be interested in [vagrant-bosh](https://github.com/cppforlife/vagrant-bosh).


### Usage

1. Run `go get github.com/cppforlife/packer-bosh` to download `packer-bosh`. 
  Make sure to run `git submodule update --init` from inside `packer-bosh` directory.

2. Configure Packer's `$HOME/.packerconfig` to know about `packer-bosh` executable:

```
{
  "provisioners": {
    "packer-bosh": "/your-bin-dir/packer-bosh"
  }
}
```

See [Packer core configuration](http://www.packer.io/docs/other/core-configuration.html) 
for more details on how to configure custom provisioners.

2. Add new provisioning step to your `template.json`. For example:

```
{
  ...

  "provisioners": [{
    "type": "packer-bosh",
    "manifest_path": "example-bosh-lite-manifest.yml",
    "assets_dir": "/your-go-dir/src/github.com/cppforlife/packer-bosh/bosh-provisioner/assets",
    "ssh_password": "ubuntu"
  }]
}
```

See [dev/template-vbox-*.json](dev/template-vbox-bosh-lite.json) for example using VirtualBox.

3. Create a deployment manifest and specify it via `manifest_path` option.
   See [dev/example-bosh-manifest.yml](dev/example-bosh-manifest.yml) 
   for an example deployment manifest used to deploy BOSH Director with BOSH Warden CPI.

4. Run `packer build` to build a VM


### Provisioner options

- `manifest_path` (String, default: `nil`)
  should contain path to a full BOSH deployment manifest

- `assets_dir` (String, default: `nil`)
  should contain path to directory with `bosh-provisioner` assets

- `ssh_password` (String, default: `nil`)
  should contain password of SSH user to run sudo; 
  can be empty if password-less sudo is configured

- `full_stemcell_compatibility` (Boolean, default: `false`)
  forces provisioner to install all (not just minimum) dependencies usually found on a stemcell

- `agent_infrastructure` (String, default: `warden`)
  configures BOSH Agent infrastructure (e.g. `aws`, `openstack`)

- `agent_platform` (String, default: `ubuntu`)
  configures BOSH Agent platform (e.g. `ubuntu`, `centos`)

- `agent_configuration` (Hash, default: '{ ... }')


### Using provisioner to build BOSH stemcells

See [building AWS Stemcell](docs/build-aws-stemcell.md).
