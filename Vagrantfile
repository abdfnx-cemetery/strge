Vagrant.configure("2") do |config|
  config.vm.define "fedora" do |g|
    c.vm.box = "fedora/28-cloud-base"
    c.vm.synced_folder ".", "/scripts/vagrant", type: "rsync",
      rsync__exclude: "bundles", rsync__args: ["-vadz", "--delete"]
    c.vm.provision "shell", inline: <<-SHELL
      sudo /vagrant/vagrant/provision.sh
    SHELL
  end
  config.vm.define "debian" do |g|
  end
end
