#! /bin/bash

# BE CAREFUL!
# "noninteractive", means that whatever distro commands
# you run will say yes to defaults.
export DEBIAN_FRONTEND=noninteractive
apt-get update

# install anything you need in the OS to get your app working
apt-get install -y --no-install-recommends python-pip
pip install requests

# clean up
apt-get clean
rm -rf /var/lib/apt

# We need to symlink everything in /etc/sv to /etc/service
mkdir -p /etc/service
ln -s /etc/sv/* /etc/service/
