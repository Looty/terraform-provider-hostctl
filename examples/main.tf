terraform {
  required_providers {
    hostctl = {
      version = "0.0.1"
      source  = "hashicorp.com/Looty/hostctl"
    }
  }
}

provider "hostctl" {}

resource "hostctl_profile" "example" {
  profile = "homelab"

  domain = [
    "my-awesome.project.loc",
    "my-awesome-ui.project.loc",
  ]
}
