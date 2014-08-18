#!/bin/bash

set -e

tmp_path=/tmp/kernel_debs
mkdir $tmp_path
cd $tmp_path

debs=(
  linux-headers-3.13.0-32_3.13.0-32.56_all.deb
  linux-headers-3.13.0-32-generic_3.13.0-32.56_amd64.deb
  linux-image-3.13.0-32-generic_3.13.0-32.56_amd64.deb
  linux-image-extra-3.13.0-32-generic_3.13.0-32.56_amd64.deb
)

dload="https://raw.githubusercontent.com/cloudfoundry/bosh/master/stemcell_builder/stages/system_kernel/assets/trusty/"

for deb in ${debs[@]}; do
  wget "${dload}/${deb}"
  dpkg -i $deb
done
