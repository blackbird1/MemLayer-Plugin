# Gemini CLI Setup

This repository uses the **standard Gemini CLI plugin flow**: build a local binary, install the plugin manifest, and let the binary read MCP configuration at runtime.

## Standard Install (Gemini CLI)

1. **Build the plugin binary** from the repo root:
   ```bash
   go build -o memlayer .
   ```
2. **Create MCP configuration** at `~/.mcp.json`:
   ```json
   {
     "server_url": "YOUR_MCP_SERVER_URL",
     "api_key": "YOUR_API_KEY"
   }
   ```
3. **Install the plugin** with Gemini CLI:
   ```bash
   gemini plugin install ./plugin.yaml
   ```
4. **Verify the plugin can reach MCP** (example):
   ```bash
   ./memlayer prociq_get_memory_stats
   ```

## What Gemini Loads (and What It Does Not)

- **Gemini CLI loads `plugin.yaml`** and launches the binary specified under `runtime.binary` (`./memlayer`).
- **Gemini CLI does not load MCP config itself** for this repo. The `memlayer` binary reads `~/.mcp.json` at runtime.

## Installing to a Global Path (Optional)

If you want Gemini to call a globally installed binary, update `plugin.yaml` to use an absolute path (or a wrapper script on your `PATH`) and then re-run:

```bash
gemini plugin install ./plugin.yaml
```

## Troubleshooting

### `run_shell_command` errors

If you see errors about `run_shell_command` not being found, the install command is being attempted inside an LLM chat/tool environment rather than your system shell. The Gemini CLI does **not** expose a `run_shell_command` tool. Run `gemini plugin install ./plugin.yaml` directly in your terminal (or a CI step).

### `spawn ~/.mcp.json EACCES`

If Gemini CLI reports `Error during discovery for MCP server 'prociq': spawn ~/.mcp.json EACCES`, it is trying to execute the JSON file as a command. The `~/.mcp.json` file is **only** read by the `memlayer` binary and should not be configured as an MCP server command in Gemini. Ensure your Gemini MCP server configuration does **not** point to `~/.mcp.json` as an executable, and confirm the file is readable (for example: `chmod 600 ~/.mcp.json` and owned by your user).
