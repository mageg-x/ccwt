# CCWT - Claude Code Web Terminal

<p align="center">
  <img src="./snapshot/home.png" width="800" alt="CCWT Desktop" />
</p>

> Self-hosted multi-user Web workspace for Claude Code CLI | User Isolation · Project Isolation · Cross-Platform

<p align="center">
  English | <a href="./README.md">中文</a>
</p>

<p align="center">
  <a href="https://github.com/ccwt/ccwt/releases">
    <img src="https://img.shields.io/badge/version-1.0.0-blue.svg" alt="version">
  </a>
  <a href="https://github.com/ccwt/ccwt/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="license">
  </a>
  <img src="https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-blue.svg" alt="platform">
</p>

---

## What Problems Does It Solve?

- **Shared Development Server**: When multiple developers share a server and all use Claude Code, their authentication credentials and configurations get mixed up and overwrite each other.
- **Multiple Projects**: Context learned by Claude in Project A disappears when you switch to Project B.
- **Mobile Office**: Unable to use Claude Code on phones or tablets.

**CCWT** was created to solve all these problems — a clean, focused, multi-user, project-isolated Claude Code cloud workspace.

---

## Key Features

### 🔐 Multi-User Isolation · Project Isolation

Each CCWT user has their own completely independent Claude configuration space:

```
~/.ccwt/users/
├── alice/
│   ├── .claude/      # Alice's own OAuth credentials, settings, history
│   └── workspace/    # Alice's project folder
└── bob/
    ├── .claude/      # Bob's independent config, completely isolated from Alice
    └── workspace/    # Bob's project folder
```

When User A completes OAuth authentication, User B still needs to authenticate separately when logging into CCWT — they cannot reuse User A's Token.

### 🖥️ 100% Native Terminal Experience

- Support for all `/slash` commands
- Support for MCP protocol
- Support for interactive input
- Terminal scrollback buffer (5MB)
- Session recovery after page refresh

### 📁 Project-Level Context Isolation

Different projects for the same user are stored in separate folders. Claude Code's context and session history are automatically isolated by project. Environment variables and working directory sync automatically when switching projects.

### 🔧 Built-in SOCKS5 Proxy

Solves the IP drift problem when first logging into Claude on remote servers. Administrators can enable the SOCKS5 proxy with one click, and users can complete OAuth authentication after configuring their local proxy.

### 🎤 Offline Voice Input

Uses `whisper.cpp` for offline speech recognition. User audio data is processed only in memory and never uploaded to third parties, protecting privacy.

<p align="center">
  <img src="./snapshot/voice.png" width="400" alt="Voice Input" />
</p>

### 📱 Cross-Platform Adaptive

- **Desktop**: Left sidebar file tree + right side terminal, clear view division
- **Mobile**: Hamburger menu slides out sidebar, terminal full-screen display
- Virtual function key bar (Ctrl, Tab, Esc, ↑, ↓, etc.)

---

## Quick Start

### One-Click Deployment

CCWT is packaged as a **single binary file** with zero dependencies:

```bash
# Download the binary for your platform
curl -L https://github.com/ccwt/ccwt/releases/latest/download/ccwt-linux-amd64 -o ccwt
chmod +x ccwt

# Start directly (default port 3000)
./ccwt
```

### Enable Registration (Optional)

```bash
# Start with invite code
./ccwt -code=your-secret-code

# Disable registration (default)
./ccwt
```

### First-Time Usage

1. Visit `http://your-server-ip:3000`
2. Register your first user (automatically becomes admin)
3. Go to Settings → Enable SOCKS5 Proxy
4. Configure SOCKS5 proxy `your-server-ip:1080` in your local browser
5. Enter `claude` in the terminal to complete OAuth authentication
6. Disable the proxy and start using

---

## Tech Stack

| Layer | Tech Stack |
|-------|------------|
| Backend | Go + Gin |
| Frontend | Vue 3 + Vite |
| Terminal | xterm.js + go-pty |
| Database | SQLite |
| Auth | JWT |
| Isolation | Bubblewrap (bwrap) |

---

## Roadmap

- [ ] Team collaboration and project sharing
- [ ] MCP server GUI management
- [ ] More terminal themes

---

## Contributing

Issues and Pull Requests are welcome!

If you're interested in CCWT, feel free to give it a Star.

**GitHub**: https://github.com/ccwt/ccwt

---

*CCWT - Enabling developers to use Claude Code in the purest, most isolated way, from anywhere, on any device.*
