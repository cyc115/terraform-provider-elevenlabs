# Changelog

## v0.2.0 (unreleased)

- chore(release): add GoReleaser config with full HashiCorp platform matrix and GPG signing
- ci: add GitHub Actions release workflow (triggered on v* tags) and test workflow (PR + main)
- chore(release): add `terraform-registry-manifest.json` (protocol_versions 6.0)
- docs: generate provider documentation via `tfplugindocs` (initial `docs/` tree)
- docs: add `docs/PUBLISH_CHECKLIST.md` runbook for first Terraform Registry publish
- chore(release): add `scripts/release/gen-gpg-key.sh` for dedicated signing key generation
- docs(readme): point installation at `registry.terraform.io/providers/cyc115/elevenlabs`

## v0.1.1

- fix(examples): parameterize `agent_name` and `webhook_post_call_id` in clone-restaurant-agent and replicate-agent
- fix(examples): standardize auth to variable+sensitive pattern; add `sk_` prefix validation
- fix(examples): add per-example `.gitignore` blocking `terraform.tfvars`; add `tfvars.example` templates
- docs(examples): add `README.md` to all three examples and a top-level `examples/README.md` index

## v0.1.0

- feat: initial provider release with `elevenlabs_convai_agent`, `elevenlabs_workspace_webhook`, `elevenlabs_convai_phone_number` resources
- feat: full CRUD + import support for all three resources
- feat: acceptance tests and diff-agents.py utility script
