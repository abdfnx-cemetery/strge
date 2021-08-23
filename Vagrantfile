Vagrant.configure("2") do |config|
  # fedora
  config.vm.define "fedora" do |g|
    g.vm.box = "fedora/28-cloud-base"
    g.vm.synced_folder ".", "/vagrant", type: "rsync",
      rsync__exclude: "bundles", rsync__args: ["-vadz", "--delete"]
    g.vm.provision "shell", inline: <<-SHELL
      sudo /vagrant/vagrant/provision.sh
    SHELL
  end

  # debian
  config.vm.define "debian" do |g|
    g.vm.box = "debian/jessie64"
    g.vm.synced_folder ".", "/vagrant", type: "rsync",
      rsync__exclude: "bundles", rsync__args: ["-vadz", "--delete"]
      g.vm.provision "shell", inline: <<-SHELL
      sudo /vagrant/vagrant/provision.sh
    SHELL
  end
end
