terraform {
  source = "github.com/shalb/cluster.dev?ref=v0.1.10/terraform/aws/minikube"

  after_hook "after_hook" {
    commands     = ["apply", "plan"]
    execute      = ["./pull_kubeconf.sh", "terragrunt-test"]
    run_on_error = false
  }
}

inputs = {
  hosted_zone       = "terragrunt-test.cluster.dev"
  region            = "eu-central-1"
  cluster_name      = "terragrunt-test"
  aws_instance_type = "m5.large"
}

include {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../vpc"]
}
