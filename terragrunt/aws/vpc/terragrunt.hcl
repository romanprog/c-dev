terraform {
  source = "github.com/shalb/cluster.dev?ref=v0.1.10/terraform/aws/vpc"
}

inputs = {
  vpc_id            = "default"
  availability_zones = ["eu-central-1b", "eu-central-1c"]
  vpc_cidr          = "10.0.0.0/18"
  region            = "eu-central-1"
  cluster_name      = "terragrunt-test"
}

include {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../route53"]
}