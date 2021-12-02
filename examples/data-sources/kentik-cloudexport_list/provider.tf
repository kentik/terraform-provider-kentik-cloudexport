# Note: configuration for locally-built provider
terraform {
  required_providers {
    kentik-cloudexport = {
      version = ">= 0.3.0"
      source  = "kentik/automation/kentik-cloudexport"
    }
  }
}

provider "kentik-cloudexport" {}
