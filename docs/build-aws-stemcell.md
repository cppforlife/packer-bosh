## Building AWS BOSH stemcell with BOSH provisioner

!!! Stemcells produced using this method should NOT be used in production !!!

To quickly create BOSH stemcells for AWS:

1. Configure Packer `amazon-ebs` builder with your AWS account settings.

See Packer's [build an image](http://www.packer.io/intro/getting-started/build-image.html) page.

2. Configure BOSH provisioner with following settings:

```
"provisioners": [{
  "type": "packer-bosh",

  "assets_dir": "/your-go-dir/src/github.com/cppforlife/packer-bosh/bosh-provisioner/assets",

  "full_stemcell_compatibility": true,

  "agent_infrastructure": "aws",
  "agent_platform": "ubuntu",
  "agent_configuration": {},

  "ssh_password": "ubuntu"
}]
```

3. Remove `cloud-init` package included by default in Ubuntu AMIs
   to avoid conflicts with BOSH Agent auto configuration:

```
"provisioners": [
  ...

  { "type": "shell", "inline": "apt-get -y purge --auto-remove cloud-init" },
  { "type": "shell", "inline": "echo 'LABEL=cloudimg-rootfs /  ext4 defaults  0 0' > /etc/fstab" }
]
```

4. `packer build`

5. After a successful build, Packer will output AMI ID at the end of the output.
   Optionally make AMI public by changing its permissions
   if stemcell will be used from a different AWS account.

6. Unpack one of the officially published `light-bosh` stemcells and
   update `stemcell.MF` with new AMI reference, then repack.

7. Upload your new light stemcell to a BOSH Director and use it in your deployment.
