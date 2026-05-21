"""TDD tests for scripts/diff-agents.py.

Test 1 (RED → GREEN): agents differing only in volatile keys (agent_id, name) → exit 0
Test 2 (RED → GREEN): agents differing in prompt.prompt → exit 1
Test 3: agents differing in created_at_unix_secs (volatile) → exit 0
Test 4: agents differing in a non-volatile nested field → exit 1
"""

import json
import subprocess
import sys
from pathlib import Path

import pytest

SCRIPT = Path(__file__).parent.parent / "scripts" / "diff-agents.py"

BASE_AGENT = {
    "agent_id": "agent_aaa",
    "name": "test-agent",
    "created_at_unix_secs": 1700000000,
    "last_call_time_unix_secs": 1700001000,
    "access_info": {"is_deleted": False},
    "conversation_config": {
        "agent": {
            "prompt": {
                "prompt": "You are a helpful assistant.",
                "llm": "claude-3-5-sonnet",
                "temperature": 0,
                "max_tokens": -1,
                "tools": [],
                "knowledge_base": [],
            },
            "first_message": "Hello!",
            "language": "en",
        },
        "asr": {"quality": "high", "provider": "elevenlabs", "user_input_audio_format": "ulaw_8000"},
        "tts": {
            "model_id": "eleven_turbo_v2_5",
            "voice_id": "voice_abc",
            "agent_output_audio_format": "ulaw_8000",
            "optimize_streaming_latency": 3,
            "stability": 0.5,
            "speed": 1.0,
        },
        "conversation": {"max_duration_seconds": 300, "client_events": []},
        "turn": {"turn_timeout": 7, "silence_end_call_timeout": -1, "mode": "turn"},
    },
    "platform_settings": {
        "widget": None,
        "elevenlabs_conversational_ai": None,
        "data_collection": {},
    },
}


def _make_agent(**overrides):
    """Deep-copy BASE_AGENT and apply top-level key overrides."""
    import copy
    a = copy.deepcopy(BASE_AGENT)
    for k, v in overrides.items():
        a[k] = v
    return a


def _run_diff(agent_a: dict, agent_b: dict) -> int:
    """Write agents to tmp files and invoke diff-agents.py --file mode."""
    import tempfile, os

    with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as fa:
        json.dump(agent_a, fa)
        fa_path = fa.name
    with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as fb:
        json.dump(agent_b, fb)
        fb_path = fb.name
    try:
        result = subprocess.run(
            ["uv", "run", str(SCRIPT), "--file", fa_path, fb_path],
            capture_output=True,
            text=True,
        )
        return result.returncode
    finally:
        os.unlink(fa_path)
        os.unlink(fb_path)


class TestDiffAgents:
    def test_volatile_only_diff_exits_zero(self):
        """Agents differ only in agent_id + name → volatile → exit 0."""
        a = _make_agent(agent_id="agent_aaa", name="original-name")
        b = _make_agent(agent_id="agent_bbb", name="clone-name")
        assert _run_diff(a, b) == 0

    def test_prompt_diff_exits_one(self):
        """Agents differ in prompt.prompt → non-volatile → exit 1."""
        import copy
        a = _make_agent()
        b = copy.deepcopy(a)
        b["conversation_config"]["agent"]["prompt"]["prompt"] = "You are a different assistant."
        assert _run_diff(a, b) == 1

    def test_created_at_volatile_exits_zero(self):
        """created_at_unix_secs is volatile → exit 0."""
        import copy
        a = _make_agent()
        b = copy.deepcopy(a)
        b["created_at_unix_secs"] = 9999999999
        assert _run_diff(a, b) == 0

    def test_tts_voice_diff_exits_one(self):
        """voice_id change is non-volatile → exit 1."""
        import copy
        a = _make_agent()
        b = copy.deepcopy(a)
        b["conversation_config"]["tts"]["voice_id"] = "voice_xyz"
        assert _run_diff(a, b) == 1

    def test_identical_agents_exits_zero(self):
        """Identical agents → exit 0."""
        a = _make_agent()
        assert _run_diff(a, a) == 0
