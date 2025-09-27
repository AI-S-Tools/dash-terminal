# Dash Terminal

A web-based terminal for LXC containers.

## 🎯 Vision

Create a premium web-based terminal for LXC containers that is simple and fast.

**Target:** iOS/Android PWA that feels like a native terminal app with professional UX.

**Core Value:** Zero-config terminal access to LXC containers.

## 🏗️ Architecture

- **Backend:** Go + WebSocket + PTY
- **Frontend:** Vanilla JS + xterm.js
- **Container:** LXC exec
- **Mobile:** PWA + Service Worker

## 🚀 Development

This project uses **Test-Driven Development** and **anti-scope creep** discipline:

### Setup
```bash
git clone https://github.com/AI-S-Tools/dash-terminal.git
cd dash-terminal
```

### DPPM Project Management
- Project managed with [DPPM](https://github.com/AI-S-Tools/dppm)
- See `docs/` for complete architecture and task breakdown
- TDD required for all components
- Strict MVP scope enforcement

### Performance Targets
- Sub-50ms terminal response time
- 60fps animations
- <2MB total PWA size

## 📖 Documentation

- **[Project Canvas](docs/lxc-pwa.md)** - Architecture and phases
- **[AI Task Control](docs/tasklist.md)** - Anti-scope creep system with TDD
- **[DPPM Integration](docs.project.yaml)** - Project management

## 🎯 MVP Scope

**Include:**
- ✅ Basic terminal with xterm.js
- ✅ LXC container selection
- ✅ PWA installation
- ✅ Mobile-responsive design

**Explicitly Forbidden:**
- ❌ Authentication systems
- ❌ File upload/download
- ❌ Multi-server support
- ❌ Theming system

## 🤖 AI Development

This project is designed for AI-driven development with:
- Comprehensive task breakdown
- TDD requirements
- Multi-AI handoff protocols
- Scope compliance checking

See `docs/tasklist.md` for AI collaboration guidelines.