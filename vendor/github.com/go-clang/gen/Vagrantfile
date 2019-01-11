# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"

  config.vm.provision :shell do |sh|
    sh.path = "scripts/vagrant-bootstrap.sh"
  end

  config.vm.provision :shell do |sh|
    sh.path = "scripts/vagrant-environment.sh"
    sh.privileged = false
  end

  config.vm.synced_folder "..", "/home/vagrant/go/src/github.com/go-clang"
end
