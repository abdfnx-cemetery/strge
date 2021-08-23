Vagrant.configure("2") do |config|
  config.vm.define "fedora" do |g|
    c.vm.box = "fedora/28-cloud-base"
    c.vm.synced_folder ".", "/scripts/vagrant", type: "rsync",
      rsync__exclude: "bundles", rsync__args: ["-vadz", "--delete"]
  end
end
