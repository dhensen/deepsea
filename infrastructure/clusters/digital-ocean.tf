module "digital-ocean-deepsea" {
#  source = "git::https://github.com/poseidon/typhoon//digital-ocean/container-linux/kubernetes?ref=a441f5c6e03e46557e3eafa6f244f649bb00095f"
  source = "/home/dino/work/typhoon/digital-ocean/container-linux/kubernetes"

  region   = "ams2"
  dns_zone = "deepsea.hensenitsolutions.nl"

  cluster_name     = "deepsea"
  image            = "coreos-stable"
  controller_count = 1
  controller_type  = "2gb"
  worker_count     = 2
  worker_type      = "1gb"
  ssh_fingerprints = ["9b:c4:0f:b5:4e:86:42:34:81:03:71:0e:4f:d3:1e:40"]

  # output assets dir
  asset_dir = "/home/dino/.secrets/clusters/deepsea"
}
