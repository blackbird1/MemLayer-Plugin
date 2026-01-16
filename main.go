package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath" // Re-add this import
	"strconv"
	"strings"
)

type Headers struct {
	Authorization string `json:"Authorization"`
}

type MemlayerConfig struct {
	Type      string  `json:"type"`
	ServerURL string  `json:"url"`
	Headers   Headers `json:"headers"`
}

type McpServers struct {
	Memlayer MemlayerConfig `json:"memlayer"`
}

type McpConfig struct {
	McpServers McpServers `json:"mcpServers"`
}

var mcpConfig McpConfig

func loadMcpConfig() error {
	configPath := ".mcp.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".mcp.json")
	}

	fmt.Printf("Attempting to load MCP config from: %s\n", configPath)
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .mcp.json at %s: %v\n", configPath, err)
		return fmt.Errorf("failed to read .mcp.json: %w", err)
	}
	fmt.Printf("Successfully read .mcp.json. Content:\n%s\n", string(fileContent))

	err = json.Unmarshal(fileContent, &mcpConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing .mcp.json: %v\n", err)
		return fmt.Errorf("failed to parse .mcp.json: %w", err)
	}
	fmt.Printf("Successfully parsed .mcp.json. ServerURL: %s, ApiKey: %s\n", mcpConfig.McpServers.Memlayer.ServerURL, mcpConfig.McpServers.Memlayer.Headers.Authorization)
	return nil
}

func callProciqApi(endpoint string, jsonPayload []byte) ([]byte, error) {
	reqBody := bytes.NewBuffer(jsonPayload)

	baseURL := mcpConfig.McpServers.Memlayer.ServerURL
	var fullURL string
	if strings.Contains(baseURL, "?") {
		parts := strings.SplitN(baseURL, "?", 2)
		fullURL = parts[0] + endpoint + "?" + parts[1]
	} else {
		fullURL = baseURL + endpoint
	}

	req, err := http.NewRequest("POST", fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", mcpConfig.McpServers.Memlayer.Headers.Authorization)

	tr := &http.Transport{}
	if os.Getenv("PROCIQ_SKIP_TLS_VERIFY") == "true" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for %s: %w", endpoint, err)
	}

	return respBody, nil
}

func main() {
	if err := loadMcpConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading MCP configuration: %v\n", err)
		// We can still proceed if the config is not loaded, but some tools might not work.
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		fmt.Println("memory_reminder")
		fmt.Println("memory_workflow")
		fmt.Println("get_memory_documentation")
		// Add prociq tools to list output
		fmt.Println("prociq_retrieve_context")
		fmt.Println("prociq_log_episode")
		fmt.Println("prociq_search_episodes")
		fmt.Println("prociq_get_memory_stats")
		fmt.Println("prociq_trigger_consolidation")
	case "memory_reminder":
		fmt.Println(`<prociq-memory-reminder>
Memory tools available via prociq_mcp:
- prociq_retrieve_context: Get relevant past experiences before starting a task
- prociq_log_episode: Record what happened after completing/failing a task
- prociq_search_episodes: Search past experiences by query
- prociq_get_memory_stats: Check memory system status

Best practice: Call prociq_retrieve_context at the start of non-trivial tasks.
</prociq-memory-reminder>`)
	case "memory_workflow":
		fmt.Println(`<prociq-memory>
## REQUIRED: Create TodoWrite with Memory Tasks

Your FIRST action must be calling TodoWrite with these items at the TOP of your list:

1. "Check prociq memory for relevant context" (status: in_progress)
2. [Your actual task items here...]
3. "Log episode outcome" (status: pending) - ALWAYS include this as the LAST item

## Memory Task Details

**Check memory** - Call ` + "`prociq_retrieve_context`" + ` with:
- Task description (what you're doing)
- Error message (if debugging)
- This surfaces past solutions and prevents repeating mistakes

## When to Log SUCCESS Episodes

Log with ` + "`prociq_log_episode`" + ` outcome='success' ONLY when:
- Solved a recurring problem type (config, debugging, integration)
- Required investigation to find non-obvious cause
- Approach wasn't obvious (chose between multiple solutions)
- Established a pattern for future work
- First time doing X in this codebase

DO NOT log success when:
- One-time task (migration, rename) - won't recur
- Trivial fix (typo, missing import) - no learning
- Just following explicit instructions - no decision-making
- Pure research/reading - no action taken
- Routine CRUD operations - template work

TEST: "If I encounter a similar situation, would knowing what I did help?" No = skip.

## CRITICAL: On ANY Error or Failure

**STOP before retrying.** When a command fails or produces an error:

1. **STOP** - Do not immediately retry
2. **LOG THE FAILURE** - Call ` + "`prociq_log_episode`" + ` with outcome='failure', error_message, and error_type
3. **THEN RETRY** - Only after logging, try the corrected approach

This applies to:
- Command errors (exit code != 0)
- Module not found / import errors
- Test failures
- Build failures
- Any error that causes you to change approach

**Why this matters**: Logging failures BEFORE retrying captures the exact error context. If you retry first, you lose the original error details.

## When to Log FAILURE Episodes

ALWAYS log failures - they prevent repeating mistakes:
- User says "still broken", "doesn't work", "same error"
- Your code/command produces errors
- Tests/build fail after your changes
- You pivot approaches (log failed attempt BEFORE trying new one)
- You abandon/give up on an approach

Include error_message and error_type for failures.
</prociq-memory>`)
	case "get_memory_documentation":
		fmt.Println(`---
	description: This skill should be used when working with the prociq memory system, logging episodes, retrieving context, or managing patterns and skills.
	Triggers include debugging errors, starting complex tasks, or when the user asks about past experiences.
	---
	
	# ProcIQ Memory System Usage
	
	Guide for using the self-learning memory system effectively.
	
	## Core Tools
	
	| Tool | When to Use |
	|------|-------------|
	| ` + "`prociq_retrieve_context`" + ` | Before starting any non-trivial task |
	| ` + "`prociq_log_episode`" + ` | After completing or failing a task |
	| ` + "`prociq_search_episodes`" + ` | When looking for specific past experiences |
	| ` + "`prociq_get_memory_stats`" + ` | To check system health |
	| ` + "`prociq_trigger_consolidation`" + ` | To process accumulated episodes into patterns |
	
	## Logging Episodes
	
	Always log episodes after completing significant work:
	
	` + "```python" + `
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
	` + "```" + `
	
	### Importance Hints
	
	| Situation | Hint Value |
	|-----------|------------|
	| Normal success | 0.2-0.3 |
	| First-time task type | 0.5-0.6 |
	| Learned something new | 0.7-0.8 |
	| Critical discovery/failure | 0.9-1.0 |
	
	## Retrieving Context
	
	Before starting complex tasks, retrieve relevant past experiences:
	
	` + "```python" + `
	prociq_retrieve_context(
	    task_description="Specific description of what you're about to do",
	    error_state="Current error message if debugging",
	    tools=["Read", "Edit"]  # Tools you're likely to use
	)
	` + "```" + `
	
	The system returns:
	- **Episodes**: Past experiences with similar tasks
	- **Patterns**: Derived strategies from multiple experiences
	- **Suggested skills**: Related skills to consider
	
	## Reflection and Patterns
	
	For failures, use ` + "`/reflect`" + ` to analyze and store learnings:
	1. Identifies root cause
	2. Proposes strategy for next time
	3. Stores as pattern candidate
	
	For manual knowledge, use ` + "`/teach`" + `:
	- Injects knowledge directly
	- Higher importance than regular episodes
	- Fast-tracks to pattern status
	
	## Consolidation
	
	Run ` + "`/consolidate`" + ` periodically to:
	- Cluster similar episodes
	- Extract patterns
	- Decay old, low-value episodes
	- Promote patterns to skills
	
	Auto-consolidation triggers every 10 episodes and at session end.`)
	case "prociq_retrieve_context":
		type RetrieveContextArgs struct {
			TaskDescription string `json:"task_description"`
			ErrorState      string `json:"error_state"`
			Tools           string `json:"tools"`
		}
		args := RetrieveContextArgs{
			TaskDescription: os.Args[2],
			ErrorState:      os.Args[3],
			Tools:           os.Args[4],
		}
		jsonPayload, err := json.Marshal(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON for prociq_retrieve_context: %v\n", err)
			os.Exit(1)
		}

		respBody, err := callProciqApi("/retrieve_context", jsonPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling prociq_retrieve_context API: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(respBody))
	case "prociq_log_episode":
		type LogEpisodeArgs struct {
			TaskGoal       string  `json:"task_goal"`
			ApproachTaken  string  `json:"approach_taken"`
			Outcome        string  `json:"outcome"`
			ErrorMessage   string  `json:"error_message"`
			ToolsUsed      string  `json:"tools_used"`
			FilePatterns   string  `json:"file_patterns"`
			ComponentTypes string  `json:"component_types"`
			ImportanceHint float64 `json:"importance_hint"`
		}

		importanceHint, err := strconv.ParseFloat(os.Args[9], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing importance_hint for prociq_log_episode: %v\n", err)
			os.Exit(1)
		}

		args := LogEpisodeArgs{
			TaskGoal:       os.Args[2],
			ApproachTaken:  os.Args[3],
			Outcome:        os.Args[4],
			ErrorMessage:   os.Args[5],
			ToolsUsed:      os.Args[6],
			FilePatterns:   os.Args[7],
			ComponentTypes: os.Args[8],
			ImportanceHint: importanceHint,
		}
		jsonPayload, err := json.Marshal(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON for prociq_log_episode: %v\n", err)
			os.Exit(1)
		}
		respBody, err := callProciqApi("/log_episode", jsonPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling prociq_log_episode API: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(respBody))
	case "prociq_search_episodes":
		type SearchEpisodesArgs struct {
			Query string `json:"query"`
		}
		args := SearchEpisodesArgs{
			Query: os.Args[2],
		}
		jsonPayload, err := json.Marshal(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON for prociq_search_episodes: %v\n", err)
			os.Exit(1)
		}

		respBody, err := callProciqApi("/search_episodes", jsonPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for prociq_search_episodes: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(respBody))
	case "prociq_get_memory_stats":
		reqBody := bytes.NewBuffer([]byte("{}")) // Empty JSON body for no parameters

		respBody, err := callProciqApi("/get_memory_stats", reqBody.Bytes())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for prociq_get_memory_stats: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(respBody))
	case "prociq_trigger_consolidation":
		reqBody := bytes.NewBuffer([]byte("{}")) // Empty JSON body for no parameters

		respBody, err := callProciqApi("/trigger_consolidation", reqBody.Bytes())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for prociq_trigger_consolidation: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(respBody))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
