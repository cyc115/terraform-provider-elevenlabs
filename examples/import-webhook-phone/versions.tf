terraform {
  required_providers {
    elevenlabs = {
      # dev_overrides: install via `make install` in repo root
      # sets ~/.terraform.d/plugins/registry.terraform.io/cyc115/elevenlabs/
      source  = "registry.terraform.io/cyc115/elevenlabs"
      version = "~> 0.1.0"
    }
  }
}
