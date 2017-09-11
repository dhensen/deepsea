output "first_k8s_node_address" {
  value = "${digitalocean_droplet.k8s-node.0.ipv4_address}"
}

output "all_addresses" {
  value = ["${digitalocean_droplet.k8s-node.*.ipv4_address}"]
}