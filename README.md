# vibe-cli

**vibe-cli** is the command-line component of **Vibe** — a local, AI-assisted coding environment that pairs with Gemini CLI and Vibe Mobile app.  
It provides lightweight version tracking, automatic change detection, and snapshot management for your development projects — all locally, with no cloud dependency.

---

## Overview

`vibe-cli` acts as your local project “memory.”  
It records code changes as Git-style patches called **vibes**, allowing you to:

- Save, list, and restore checkpoints of your work.  
- Automatically snapshot edits made by Gemini CLI.  
- Run entirely offline, with zero configuration.

---

## Commands

| Command                 | Description                                                                                              |
| ----------------------- | -------------------------------------------------------------------------------------------------------- |
| `vibe init`             | Initializes a new `.vibes` session in the current directory and automatically marks it for tracking.     |
| `vibe save "<message>"` | Saves your current working diff as a **vibe patch**, storing it under `.vibes/patches/`.                 |
| `vibe list`             | Displays all saved vibes (newest first).                                                                 |
| `vibe diff <id>`        | Prints the patch contents for a specific vibe ID.                                                        |
| `vibe checkout <id>`    | Reverts your workspace to the state of a saved vibe.                                                     |
| `vibe reset`            | Removes all tracked vibes and clears the `.vibes/` directory.                                            |
| `vibe watch`            | Starts a daemon that watches marked directories for file changes and automatically triggers `vibe save`. |
