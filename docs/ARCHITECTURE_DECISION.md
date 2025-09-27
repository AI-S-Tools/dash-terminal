# Architecture Decision: Zellij vs Tmux for Terminal Multiplexing

## Decision: Use Zellij Instead of Tmux

### Why Not Tmux

Tmux was designed for SSH sessions, not web terminals. This creates friction:

1. **Control Mode Parsing** - Tmux control mode output is difficult to parse reliably
2. **Keybinding Conflicts** - Prefix keys (Ctrl+B) interfere with browser shortcuts
3. **State Synchronization** - Keeping tmux state and browser UI in sync is complex
4. **Limited API** - No structured data format, just escape sequences
5. **UI Constraints** - Hard to implement modern web UI features on top of tmux

### Why Zellij

Zellij is designed with modern use cases in mind:

1. **Better API** - Structured communication protocol
2. **Plugin System** - WebAssembly plugins can bridge to browser
3. **Modern Architecture** - Built for programmatic control
4. **Floating Panes** - Better suited for web UI paradigms
5. **No Keybinding Theft** - Can be run in headless mode

### Implementation Strategy

#### Option 1: Go + Zellij Subprocess (Recommended for MVP)
```go
// Keep existing Go server, use Zellij as multiplexer
cmd := exec.Command("zellij", "attach", "--create", sessionName)
```

**Pros:**
- Reuse existing PTY code
- Fast to implement
- Can migrate later if needed

**Cons:**
- Extra process overhead
- Some API limitations

#### Option 2: Rust Rewrite with Native Zellij
```rust
// Direct Zellij integration in Rust
use zellij_utils::data::*;
```

**Pros:**
- Single binary per container
- Maximum performance
- Full Zellij API access

**Cons:**
- Complete rewrite needed
- Longer development time

## Recommendation

Start with **Option 1** (Go + Zellij subprocess) for fast iteration, with clear migration path to Option 2 if performance demands it.

The key insight: Tmux steals too much control from a modern web terminal. Zellij provides multiplexing without the baggage.