---
description: Manually inject knowledge into the memory system
argument-hint: "<lesson or best practice to remember>"
---

# Manually Inject Knowledge

Add knowledge directly to the memory system without an execution episode.

## Instructions

Inject knowledge provided by the user into the memory system.

Parse $ARGUMENTS to understand what is being taught. The user might provide:
- A lesson learned from experience
- A best practice or pattern
- A warning about something to avoid
- A strategy for handling a specific situation

1. Create an episode using `prociq_log_episode` with:
   - task_goal: "Manual teaching: [summary of what was taught]"
   - approach_taken: The full knowledge/lesson from $ARGUMENTS
   - outcome: "success" (teaching is always successful)
   - importance_hint: 0.8 (manually taught knowledge is important)
   - tools_used: Extract any tool names mentioned
   - component_types: Extract any component/technology types mentioned

2. If the teaching includes clear trigger conditions and strategy, also store it as a reflection using `prociq_store_reflection` to fast-track it toward becoming a pattern

3. Confirm what was stored and suggest running `prociq_trigger_consolidation` if enough teachings have accumulated

Example usage:
- `memory.teach When using Detox with animations, always add explicit waitFor timeouts of at least 5000ms`
- `memory.teach The auth module requires running migrations before tests`
