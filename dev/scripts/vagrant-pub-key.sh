#!/bin/bash

set -e

mkdir ~/.ssh

wget -qO- https://raw.github.com/mitchellh/vagrant/master/keys/vagrant.pub >> ~/.ssh/authorized_keys
