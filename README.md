# devfix
# fuck me this is written by ai I aint writing a novel for this

A comprehensive, all-in-one developer system toolkit.

`devfix` helps you clean up, diagnose, and fix your local developer environment instantly. Instead of remembering 20 different CLI tools and scripts, `devfix` packages everything into a beautiful interactive terminal menu.

## Features

- **System Diagnostics:** 
  - `doctor`: Holistic system diagnosis (checks ports, temps, git, PATH) and gives your system a Health Score.
  - `ports` & `kill`: List listening TCP ports with PIDs and process names. Easily kill blocked processes by port number.
  - `sys`: Audit your PATH variable for broken links and duplicates.
- **Cleanup & Fixes:**
  - `clean`: Recursively finds old unused dev folders (`node_modules`, `target`, `.venv`) and cleans them up.
  - `nuke`: Quick fix for crashed projects (instantly removes locks and dependency folders).
  - `distress`: Deep clean tool caches (npm, pip, docker) and view extreme factory reset instructions.
- **Project & Git:**
  - `env`: Audit local `.env` files against `.env.example` and automatically fix missing keys.
  - `drift`: Detect tool version mismatches across your project.
  - `repos`: Recursively scan for Git repos that have uncommitted or unpushed changes.
  - `guard`: Security check that scans your git diff for accidentally exposed API keys or passwords before you commit.
- **Misc Utilities:**
  - `serve`: Fast local HTTP server in your current directory.
  - `flush`: Flush local DNS cache and check your hosts file.
  - `history`: Export a history of your actions to a PowerShell script.

## Installation

### For Windows Users
Open an Administrator PowerShell and run:
```powershell
irm https://raw.githubusercontent.com/Domcho214/devfix/main/install.ps1 | iex
```

### Manual Installation (All Platforms)
Check out the `Releases` tab on GitHub to download the standalone binary for Windows, macOS, or Linux. Just drop it into your system's PATH.

## Usage

Simply open your terminal anywhere and type:
```bash
devfix
```
This will open the interactive UI. You can also run specific commands directly, e.g.:
```bash
devfix nuke
```

---
*For support or questions, contact me on Discord: @domco00*
