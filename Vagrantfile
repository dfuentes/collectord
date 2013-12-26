# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

BOX_NAME = "ubuntu"
BOX_URI = "http://files.vagrantup.com/precise64.box"

$script = <<SCRIPT
# grab and install go 1.2
/usr/bin/wget -O /usr/local/src/go1.2.linux-amd64.tar.gz https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
/bin/tar -C /usr/local -xzf /usr/local/src/go1.2.linux-amd64.tar.gz

# install git (for go get)
/usr/bin/apt-get install -y git

# make workspace
/bin/mkdir -p /home/vagrant/go/src/github.com/dfuentes/
/bin/mkdir -p /home/vagrant/go/pkg
/bin/mkdir -p /home/vagrant/go/bin

# setup path export
/bin/echo 'export PATH=/home/vagrant/go/bin:/usr/local/go/bin:$PATH' >> /home/vagrant/.profile
/bin/echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.profile

source /home/vagrant/.profile

# grab dependencies
cd /home/vagrant/go/src/github.com/dfuentes/collectord
go get -d

# give ownership of workspace to vagrant user
/bin/chown -R vagrant /home/vagrant/go

SCRIPT

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = BOX_NAME
  config.vm.box_url = BOX_URI

  config.vm.provision :shell, :inline => $script

  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/dfuentes/collectord", owner: "vagrant", group: "vagrant"
end
