
locals {
  # Automatically load region-level variables
  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl"))
  aws_region  = local.region_vars.locals.aws_region
}

# Generate an AWS provider block
generate "provider" {
  path      = "${get_terragrunt_dir()}/aws_providers_override.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
    terraform {
      required_providers {
        aws = {
          source  = "hashicorp/aws"
          version = "3.27.0"
        }
      }
    }
  EOF
}
