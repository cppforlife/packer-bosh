#!/bin/bash

set -e

apt-get -y update
apt-get -y install linux-headers-$(uname -r) build-essential dkms puppet-common
apt-get -y clean

mount -o loop /tmp/VBoxGuestAdditions.iso /media/cdrom
/media/cdrom/VBoxLinuxAdditions.run || true
umount /media/cdrom
