terraform {
  source = "github.com/shalb/cluster.dev?ref=v0.1.10/terraform/aws/route53"
}

inputs = {
  region            = "eu-central-1"
  cluster_name      = "terragrunt-test"
  cluster_domain    = "cluster.dev"
  zone_delegation   = "false"
 
}

include {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../backend"]
}