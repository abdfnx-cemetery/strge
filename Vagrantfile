Vagrant.configure("2") do |config|
  config.vm.define "debian" do |g|
    g.vm.box = "debian/jessie64"
    g.vm.synced_folder ".", "/vagrant", type: "rsync",
      rsync__exclude: "bundles", rsync__args: ["-vadz", "--delete"]
      g.vm.provision "shell", inline: <<-SHELL
      sudo /vagrant/vagrant/provision.sh
    SHELL
  end
end
