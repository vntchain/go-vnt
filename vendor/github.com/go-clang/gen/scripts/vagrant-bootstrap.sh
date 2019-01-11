#!/bin/bash

set -exuo pipefail

export CODENAME=$(lsb_release --codename --short)

# Update and upgrade
apt-get update
apt-get -V upgrade -y
apt-get -V autoremove -y

# Install needed packages
apt-get -V install -y git make

# Change the owner of the go directory. This is needed because the missing parent folders up to our synced folder are created with the user "root" instead of "vagrant".
chown -R vagrant:vagrant /home/vagrant/go
