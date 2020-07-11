terraform {
  source = "github.com/shalb/cluster.dev?ref=v0.1.10/terraform/aws/backend"
}

inputs = {
  region            = "eu-central-1"
  s3_backend_bucket = "terragrunt-test"
}
