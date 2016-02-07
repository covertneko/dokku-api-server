# -*- mode: ruby -*-
# vi: set ft=ruby :

BOX_NAME = ENV['BOX_NAME'] || 'trusty'
BOX_URI = ENV['BOX_URI'] || 'https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box'
DOKKU_DOMAIN = ENV['DOKKU_DOMAIN'] || 'dokku.local'
DOKKU_IP = ENV['DOKKU_IP'] || '10.0.0.2'
# User's public key to administer Dokku
DOKKU_LOCAL_KEY_FILE = ENV['DOKKU_LOCAL_KEY_FILE'] || "#{ENV['HOME']}/.ssh/id_rsa.pub"

Vagrant.configure(2) do |config|
  config.ssh.forward_agent = true
  config.vm.box = BOX_NAME
  config.vm.box_url = BOX_URI
  config.vm.network :forwarded_port, guest: 80, host: 8080
  config.vm.hostname = "#{DOKKU_DOMAIN}"
  config.vm.network :private_network, ip: DOKKU_IP

  config.vm.provider 'virtualbox' do |vb|
    vb.customize ['modifyvm', :id, '--natdnshostresolver1', 'on']
  end

  # Install dokku
  config.vm.provision 'shell' do |s|
    # Copy local public key for dokku authorized keys
    pubkey = File.readlines("#{DOKKU_LOCAL_KEY_FILE}").first.strip

    s.inline = <<-SHELL
      set -eo pipefail
      export DEBIAN_FRONTEND=noninteractive

      echo "Installing prerequisites..."

      # install prerequisites
      #sudo apt-get update -qq > /dev/null
      #sudo apt-get install -qq -y apt-transport-https

      # install docker
      #wget -nv -O - https://get.docker.com/ | sh

      # add dokku source
      #wget -nv -O - https://packagecloud.io/gpg.key | apt-key add -
      #echo "deb https://packagecloud.io/dokku/dokku/ubuntu/ trusty main" | sudo tee /etc/apt/sources.list.d/dokku.list
      #sudo apt-get update -qq > /dev/null

      echo "Configuring dokku..."
      # configure dokku installation options:
        # - disable web config
        # - enable vhost-based app deployment
        # - set hostname
        # - skip key file check (will add manually after installation)
      debconf-set-selections <<< "
        dokku dokku/web_config boolean false
        dokku dokku/vhost_enable boolean true
        dokku dokku/hostname string #{DOKKU_DOMAIN}
        dokku dokku/skip_key_file boolean true
        dokku dokku/key_file string /root/.ssh/id_rsa.pub"

      echo "Installing dokku..."
      # install dokku
      sudo apt-get -y install dokku
      sudo dokku plugin:install-dependencies --core

      # Add public key for dokku user
      echo "#{pubkey}" | sudo sshcommand acl-add dokku #{`whoami`}
    SHELL
  end
end
