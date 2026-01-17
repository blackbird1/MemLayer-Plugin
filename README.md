# MemLayer Plugin

A self-learning memory system plugin for [Claude Code](https://docs.anthropic.com/en/docs/claude-code) that enables persistent learning across task executions.
Must be used with [MemLayer](https://prociq.ai).

## Overview

MemLayer provides Claude Code with episodic memory capabilities, allowing it to:

- **Log task executions** as episodes with outcomes, errors, and context
- **Extract patterns** from past experiences, especially failures
- **Promote proven strategies** into reusable skills
- **Retrieve relevant context** before starting new tasks
- **Learn from mistakes** to avoid repeating them

## Installation

1. Add the marketplace to Claude Code:
   ```
   /plugin marketplace add shafty023/MemLayer-Plugin
   ```

2. Install the memory plugin:
   ```
   /plugin install memory@ProcIQ
   ```

3. Configure the prociq MCP server with your API key (see [prociq.ai](https://prociq.ai) for setup)

For more details on plugin installation, see the [official documentation](https://code.claude.com/docs/en/plugin-marketplaces).

## Gemini Extension

The Gemini extension mirrors the Claude plugin commands and usage guidance. The Gemini CLI expects a `gemini-extension.json` file at the repo root, so you can run:

```
gemini extensions install ./
```

If you need the raw manifest used by the extension, it is located at:

```
plugins/memory/.gemini-extension/extension.json
```

Once loaded, use the `prociq.audit`, `prociq.teach`, and `prociq.forget` commands to manage ProcIQ memory from Gemini.

The extension relies on the ProcIQ MCP server being configured in your Gemini environment (with your API key). The MCP server is what actually makes the HTTP requests to the ProcIQ backend when Gemini invokes `prociq_*` tools.

### MCP Configuration (Gemini)

The MCP server configuration is included with the Gemini extension and lives at:

```
plugins/memory/.gemini-extension/mcp.json
```

The bundled configuration expects the ProcIQ API key (bearer token) to be provided via the `PROCIQ_TOKEN` environment variable. A typical entry includes the endpoint and environment variable reference:

```json
{
  "name": "prociq",
  "endpoint": "http://prociq-alb-2037713618.us-east-1.elb.amazonaws.com/mcp",
  "apiKey": "${PROCIQ_TOKEN}",
  "transport": "sse"
}
```

Replace the values with the configuration details from [prociq.ai](https://prociq.ai), then restart Gemini so it loads the new MCP server configuration.

Note: installing the extension does not automatically register MCP servers with Gemini. You must copy the `mcp.json` entry into your Gemini MCP configuration file/location and restart Gemini after updating it.

#### Troubleshooting MCP Discovery

If you see errors like `Error during discovery for MCP server 'prociq': fetch failed`, check:

1. `PROCIQ_TOKEN` is set in your environment and is a valid bearer token.
2. Your network can reach `http://prociq-alb-2037713618.us-east-1.elb.amazonaws.com/mcp` (no proxy or firewall blocking HTTP).
3. Your Gemini MCP configuration file includes the `prociq` entry and Gemini was restarted after changes.

If you see `Not Acceptable: Client must accept text/event-stream`, ensure the MCP configuration uses `transport: "sse"` so Gemini negotiates Server-Sent Events.

## Project Structure

```
MemLayer-Plugin/
├── .claude-plugin/
│   └── marketplace.json      # Marketplace registration
└── plugins/
    └── memory/
        ├── .claude-plugin/
        │   └── plugin.json   # Plugin manifest
        ├── .gemini-extension/
        │   ├── extension.json  # Gemini extension manifest
        │   ├── commands/       # Gemini command prompts
        │   └── instructions/   # Gemini usage guide
        ├── commands/         # CLI commands
        │   ├── audit.md      # /memory:audit - inspect memory state
        │   ├── teach.md      # /memory:teach - inject knowledge manually
        │   └── forget.md     # /memory:forget - remove episodes
        ├── hooks/            # Integration hooks
        │   ├── hooks.json
        │   └── scripts/
        │       ├── session-start.sh
        │       └── user-prompt.sh
        └── skills/
            └── memory-usage/
                └── SKILL.md  # Memory system usage guide
```

## Commands

### `/memory:audit [episodes|patterns|skills]`
Inspect the current state of the memory system. Shows statistics, recent episodes, high-confidence patterns, and skill inventory.

### `/memory:teach <lesson>`
Manually inject knowledge into the memory system without task execution.

```
/memory:teach When using Detox with animations, always add explicit waitFor timeouts of at least 5000ms
```

### `/memory:forget <episode-id|query>`
Remove specific episodes from memory. Can search by query or delete by ID.

## Core Concepts

### Episodes
Records of task execution containing:
- Task goal and approach taken
- Outcome (success/partial/failure)
- Error details if applicable
- Tools used and file patterns involved
- Importance score (0.0–1.0)

### Patterns
Derived learnings extracted from multiple episodes:
- Root cause analysis
- Recommended strategy
- Trigger conditions (errors, keywords, tools)

### Skills
Mature, high-confidence patterns promoted to reusable knowledge that gets surfaced when relevant tasks arise.

### Notes
Persistent freeform knowledge entries that never decay (unlike episodes):
- Recording important discoveries that shouldn't fade
- Documenting project-specific knowledge
- Manual teaching via `/memory:teach` command
- Reference material that should always be findable

### Consolidation
Automatic processing that:
- Clusters similar episodes
- Extracts patterns from failures
- Decays old/low-value episodes
- Promotes patterns to skills

## Memory Tools (MCP)

### Episode Tools
| Tool | Purpose |
|------|---------|
| `prociq_log_episode` | Record task execution (async, non-blocking) |
| `prociq_retrieve_context` | Get relevant past experiences before a task |
| `prociq_search_episodes` | Search with filters (outcome, error_type, project) |
| `prociq_get_episode` | Retrieve full episode by ID |
| `prociq_search_episodes_full` | Semantic search returning full episodes |
| `prociq_forget_episodes` | Delete episodes permanently |
| `prociq_archive_episode` | Soft-delete episodes (reversible) |

### Note Tools
| Tool | Purpose |
|------|---------|
| `prociq_log_note` | Store persistent freeform knowledge |
| `prociq_update_note` | Modify existing note |
| `prociq_get_note` | Retrieve note by ID |
| `prociq_search_notes` | Search notes by content or tags |
| `prociq_delete_note` | Remove note from storage |

### Pattern & Skill Tools
| Tool | Purpose |
|------|---------|
| `prociq_search_patterns` | Search patterns with filters |
| `prociq_list_skills` | List all available skills |
| `prociq_get_skill_content` | Retrieve skill markdown by ID |

### System Tools
| Tool | Purpose |
|------|---------|
| `prociq_get_memory_stats` | View memory health and statistics |
| `prociq_trigger_consolidation` | Manually run memory maintenance |

## Best Practices

### When to Log

**Do log:**
- Failures (always, before retrying)
- Non-obvious solutions requiring investigation
- First-time task types
- Recurring problem categories (config, debugging, integration)

**Don't log:**
- Trivial fixes (typos, missing imports)
- Routine CRUD operations
- Pure research/exploration tasks

### Importance Scoring

| Scenario | Score |
|----------|-------|
| Normal success | 0.2–0.3 |
| First-time task type | 0.5–0.6 |
| Learned something new | 0.7–0.8 |
| Critical discovery/failure | 0.9–1.0 |

### Critical Rule: Log Failures First

Always log a failure **before** retrying. This captures the exact error context that would otherwise be lost after a successful retry.

## How Hooks Work

1. **SessionStart** — Reminds Claude about available memory tools
2. **UserPromptSubmit** — Injects memory workflow into TodoWrite (check memory first, log outcome last)
3. **Stop** — Reminds Claude to log failures and suggests reflection

## License

MIT License — see [LICENSE](LICENSE) for details.

## Author

Daniel Ochoa ([@shafty023](https://github.com/shafty023))

## Links

- [prociq.ai](https://prociq.ai) — Memory system backend
- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
