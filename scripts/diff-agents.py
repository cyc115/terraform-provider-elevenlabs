#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.12"
# dependencies = ["httpx>=0.27", "deepdiff>=7.0"]
# ///
"""Compare two ElevenLabs ConvAI agents, ignoring volatile fields.

Usage (live API):   diff-agents.py <agent_id_a> <agent_id_b>
Usage (file mode):  diff-agents.py --file <path_a.json> <path_b.json>

Exit 0 — no meaningful diff
Exit 1 — agents differ in non-volatile fields
Exit 2 — error (network, auth, bad args)
"""

import json
import os
import sys
from pathlib import Path

import httpx
from deepdiff import DeepDiff

# Top-level keys stripped before comparison
_VOLATILE_TOP = {
    "agent_id",
    "name",
    "created_at_unix_secs",
    "last_call_time_unix_secs",
    "access_info",
    # Per-agent identity / versioning
    "version_id",
    "branch_id",
    "main_branch_id",
    # API-managed, not configurable via provider
    "coaching_settings",
    "procedure_settings",
    "metadata",
    "whatsapp_accounts",
    "phone_numbers",
    "workflow",
    "tags",
    "trust_context",
    "smb_metadata",
    "ban",
    "procedures",
    "procedures_enabled",
    "procedure_settings",
    "summary_language",
    "alerting",
    "safety",
}

# DeepDiff exclude_paths for non-configurable nested fields
_EXCLUDE_PATHS = {
    # Widget settings (not in provider schema)
    "root['platform_settings']['widget']",
    # Attached tests (workspace-specific)
    "root['platform_settings']['testing']",
    # Built-in tools (API manages, not configurable)
    "root['conversation_config']['agent']['prompt']['built_in_tools']",
    # tools array mirrors built_in_tools (redundant representation)
    "root['conversation_config']['agent']['prompt']['tools']",
    # TTS suggested audio tags (read-only API field)
    "root['conversation_config']['tts']['suggested_audio_tags']",
    # Overrides / workspace_overrides schema (non-configurable)
    "root['platform_settings']['overrides']",
    # Language presets (not in provider schema)
    "root['conversation_config']['language_presets']",
}


def _strip_volatile(obj: dict) -> dict:
    return {k: v for k, v in obj.items() if k not in _VOLATILE_TOP}


def _fetch(agent_id: str) -> dict:
    api_key = os.environ.get("ELEVENLABS_API_KEY")
    if not api_key:
        print("ERROR: ELEVENLABS_API_KEY not set", file=sys.stderr)
        sys.exit(2)
    url = f"https://api.elevenlabs.io/v1/convai/agents/{agent_id}"
    resp = httpx.get(url, headers={"xi-api-key": api_key}, timeout=30)
    if resp.status_code != 200:
        print(f"ERROR: GET {url} → {resp.status_code}: {resp.text}", file=sys.stderr)
        sys.exit(2)
    return resp.json()


def _load_file(path: str) -> dict:
    try:
        return json.loads(Path(path).read_text())
    except Exception as e:
        print(f"ERROR: reading {path}: {e}", file=sys.stderr)
        sys.exit(2)


def main() -> None:
    args = sys.argv[1:]

    if len(args) >= 3 and args[0] == "--file":
        a = _load_file(args[1])
        b = _load_file(args[2])
    elif len(args) == 2 and not args[0].startswith("--"):
        a = _fetch(args[0])
        b = _fetch(args[1])
    else:
        print(__doc__, file=sys.stderr)
        sys.exit(2)

    a_stripped = _strip_volatile(a)
    b_stripped = _strip_volatile(b)

    diff = DeepDiff(a_stripped, b_stripped, ignore_order=True, exclude_paths=_EXCLUDE_PATHS)

    if diff:
        print("DIFF:", diff.to_json(indent=2))
        sys.exit(1)
    else:
        print("No diff.")
        sys.exit(0)


if __name__ == "__main__":
    main()
