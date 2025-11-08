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

Projects are tracked centrally (in `~/.vibe/tracked.json`), but all patches and vibe logs are stored locally within your project's `.vibes` directory (just like `.git`).

---

## Commands

| Command | Description |
| --- | --- |
| `vibe init [id]` | Initializes `.vibes` session in current dir. Registers it centrally with an optional `id` (defaults to folder name). |
| `vibe save "<message>"` | Saves your current working diff as a **vibe patch**, storing it under `.vibes/patches/`. |
| `vibe list` | Displays all saved vibes (newest first) for the *current* project. |
| `vibe diff <id>` | Prints the patch contents for a specific vibe ID from the *current* project. |
| `vibe checkout <id>` | Reverts your workspace to the state of a saved vibe from the *current* project. |
| `vibe reset` | Removes all tracked vibes and clears the `.vibes/` directory for the *current* project. |
| `vibe watch [id]` | Starts a daemon. Watches all tracked projects, or only the project specified by `[id]`. |
