# ProcIQ Memory System (Gemini Extension)

Use this extension to integrate Gemini with the ProcIQ memory system for persistent learning across sessions.

## When to Use

- Before starting any non-trivial task, retrieve relevant context.
- After completing or failing a task, log an episode.
- When debugging or encountering repeated issues, search past episodes or patterns.
- Periodically trigger consolidation to extract patterns and promote skills.

## Core Tools

| Tool | When to Use |
|------|-------------|
| `prociq_retrieve_context` | Before starting any non-trivial task |
| `prociq_log_episode` | After completing or failing a task |
| `prociq_search_episodes` | When looking for specific past experiences |
| `prociq_get_memory_stats` | To check system health |
| `prociq_trigger_consolidation` | To process accumulated episodes into patterns |

## Logging Episodes

Always log episodes after completing significant work:

```python
prociq_log_episode(
    task_goal="What you were trying to accomplish",
    approach_taken="How you approached it and what worked/failed",
    outcome="success" | "partial" | "failure",
    error_message="Details if failed",
    tools_used=["Read", "Edit", "Bash"],
    file_patterns=["*.py", "src/**/*.ts"],
    component_types=["api", "database", "ui"],
    importance_hint=0.8  # Higher for important discoveries
)
```

### Importance Hints

| Situation | Hint Value |
|-----------|------------|
| Normal success | 0.2-0.3 |
| First-time task type | 0.5-0.6 |
| Learned something new | 0.7-0.8 |
| Critical discovery/failure | 0.9-1.0 |

## Retrieving Context

Before starting complex tasks, retrieve relevant past experiences:

```python
prociq_retrieve_context(
    task_description="Specific description of what you're about to do",
    error_state="Current error message if debugging",
    tools=["Read", "Edit"]  # Tools you're likely to use
)
```

The system returns:
- **Episodes**: Past experiences with similar tasks
- **Patterns**: Derived strategies from multiple experiences
- **Suggested skills**: Related skills to consider

## Reflection and Patterns

For failures, analyze and store learnings:
1. Identify root cause
2. Propose a strategy for next time
3. Store as a reflection/pattern candidate

For manual knowledge, use `memory.teach`:
- Injects knowledge directly
- Higher importance than regular episodes
- Fast-tracks to pattern status

## Consolidation

Run `prociq_trigger_consolidation` periodically to:
- Cluster similar episodes
- Extract patterns
- Decay old, low-value episodes
- Promote patterns to skills

Auto-consolidation triggers every 10 episodes and at session end.
