remote_state {
  backend = "s3"
  config = {
    bucket         = "terragrunt-test"
    key            = "states/terraform-${path_relative_to_include()}.state"
    region         = "eu-central-1"
    encrypt        = false
    dynamodb_table = "terragrunt-test-state"
  }
}

