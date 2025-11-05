# vibe-cli

**vibe-cli** is the command-line component of **Vibe** — a local, AI-assisted coding environment that pairs with Gemini CLI and Vibe Mobile app.  
It provides lightweight version tracking, automatic change detection, and snapshot management for your development projects — all locally, with no cloud dependency.

---

## 🧠 Overview

`vibe-cli` acts as your local project “memory.”  
It records code changes as Git-style patches called **vibes**, allowing you to:

- Save, list, and restore checkpoints of your work.  
- Automatically snapshot edits made by Gemini CLI.  
- Run entirely offline, with zero configuration.

---

## ⚙️ Installation

```bash
git clone https://github.com/SacredNightmare99/vibe-cli.git
cd vibe-cli/cmd
go build -o ../vibe .

