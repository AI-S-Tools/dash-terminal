```markdown
# Project: dash-terminal

## Project Overview

`dash-terminal` is a web-based Progressive Web App (PWA) designed to provide a premium terminal experience for managing and interacting with LXC containers. It aims to replicate the elegant simplicity and session management capabilities of tools like SecureShellfish, with a strong focus on a native-app feel for mobile devices.

The core functionality revolves around providing zero-configuration terminal access to LXC containers and leveraging native `tmux` integration for persistent and structured sessions.

**Key Technologies:**

*   **Backend:** Go
*   **Frontend:** Vanilla JavaScript, `xterm.js`, and Web Components
*   **Communication:** WebSockets
*   **Containerization:** LXC
*   **Session Management:** `tmux`

## Architecture

The application follows a client-server architecture:

*   A **Go backend** exposes a WebSocket endpoint. This backend is responsible for managing LXC containers (listing, executing commands within them) and handling `tmux` sessions.
*   A **JavaScript frontend** provides the user interface. It uses `xterm.js` to render the terminal and communicates with the backend via WebSockets to send commands and receive output.
*   The UI is structured to mirror `tmux`'s hierarchy:
    1.  **Session Tabs:** For switching between different `tmux` sessions.
    2.  **Window Tabs:** For switching between windows within a session.
    3.  **Panes:** For displaying multiple terminal panes within a window.

## Building and Running

**Backend (Go):**

To run the backend server, execute the following command from the project root:

```bash
# TODO: Verify the exact command once cmd/server/main.go is available.
# It is likely one of the following:
go run ./cmd/server
# or
go build -o dash-terminal-server ./cmd/server && ./dash-terminal-server
```

**Frontend:**

The frontend is composed of static HTML, CSS, and JavaScript files served directly. No build step seems to be required. The application can be accessed by opening `web/index.html` in a browser, assuming the backend server is running and serving the `web` directory.

## Development Conventions

*   **Project Structure:** The project is organized into `cmd` (application entry points), `internal` (private application logic), and `web` (frontend assets).
*   **Backend:** The backend is written in Go. It uses the `gorilla/websocket` library for WebSocket communication and `creack/pty` for pseudo-terminal management.
*   **Frontend:** The frontend is built with vanilla JavaScript and `xterm.js`. It is designed as a PWA, with a service worker (`sw.js`) for offline capabilities.
*   **Communication Protocol:** The client and server communicate using a JSON-based WebSocket protocol defined in `internal/websocket/protocol.go`. Messages are structured with a `type` and a `payload`.
*   **Project Management:** The project uses a custom project management system (`DPPM`), with task breakdowns and architectural documents located in the `docs/` and `dbdocs/` directories.
```