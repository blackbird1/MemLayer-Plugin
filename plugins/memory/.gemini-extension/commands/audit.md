---
description: Inspect the current state of the memory system
author: ProcIQ
---

# Inspect Memory State

Use the ProcIQ MCP tools to retrieve memory statistics, recent episodes, and high-confidence patterns.

## Instructions

1. Call `prociq_get_memory_stats` to retrieve memory health and recent counts.
2. If the user asks for episodes, call `prociq_search_episodes` with a reasonable limit.
3. If the user asks for patterns, call `prociq_search_patterns` and return the most relevant items.
4. If the user asks for skills, call `prociq_list_skills` and summarize the inventory.

Summarize the results clearly and highlight any urgent issues (high failure rate, large backlog, etc.).
