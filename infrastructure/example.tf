# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}
variable "public_key_path" {}
variable "region" {
  default = "ams2"
}
variable "cluster_size" {
  default = "1"
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_ssh_key" "default" {
  name       = "Terraform Example Dino"
  public_key = "${file("${var.public_key_path}")}"
}

# Create a web server
resource "digitalocean_droplet" "k8s-node" {
  image    = "ubuntu-16-04-x64"
  region   = "${var.region}"
  size     = "1gb"
  ssh_keys = ["${digitalocean_ssh_key.default.id}"]
  name     = "k8s-node-${count.index + 1}"
  count    = "${var.cluster_size}"

  provisioner "local-exec" {
    command = "echo ${digitalocean_droplet.k8s-node.ipv4_address} >> private_ips.txt"
  }
}

