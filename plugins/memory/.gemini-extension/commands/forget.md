---
description: Remove episodes from the memory system
argument-hint: "<episode-id|query>"
---

# Forget Episodes

Remove specific episodes from the memory system. This can be used to delete sensitive data or incorrect records.

## Instructions

1. Parse $ARGUMENTS. If it looks like an episode ID, call `prociq_forget_episodes` with that ID.
2. If it looks like a query, call `prociq_search_episodes` first to list candidate episodes, then confirm which ones to delete.
3. After deletion, confirm the result and remind the user that deletion is permanent.

Example usage:
- `prociq.forget episode_1234`
- `prociq.forget failed auth0 migration`
