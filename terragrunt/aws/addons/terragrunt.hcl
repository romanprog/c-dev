terraform {
  source = "github.com/shalb/cluster.dev?ref=v0.1.10/terraform/aws/addons"

}

inputs = {
  cluster_cloud_domain        = "terragrunt-test.cluster.dev"
  region                      = "eu-central-1"
  cluster_name                = "terragrunt-test"
  config_path                 = "~/.kube/config"
  eks                         = "true"
}

include {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../minikube"]
}
