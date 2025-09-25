# GEMINI Project Analysis: dash-terminal Phase 2 (p2-FAIL)

## Project Overview

This directory contains the project management and recovery plan for **Phase 2 of the dash-terminal project**. The ultimate goal is to create a fully functional, web-based PWA that allows users to manage and interact with `tmux` sessions running inside different **LXC containers** through a real-time terminal interface.

## Current Status & Recovery Plan

The initial implementation of Phase 2 failed due to a critical integration gap between the frontend and the backend. The components were built in isolation and did not communicate, rendering the application non-functional.

A new, concrete recovery plan has been established and is detailed in the `tasks/` directory. This plan is designed to be implemented without mocks and verified through a rigorous, interactive testing regime.

The strategy is as follows:

1.  **Implement a Concrete Frontend (`t2-3`):** The entire frontend will be rewritten in **TypeScript** (`/src`). The UI will guide the user to first select an LXC container and then view/select the tmux sessions within it. It will manage multiple `xterm.js` instances to provide true session isolation.

2.  **Implement an Integrated Backend (`t2-4`):** The backend logic will be updated to execute all `tmux` and shell commands *inside* the user-selected LXC container via a pseudo-terminal (PTY). This ensures a real, non-mocked terminal experience.

3.  **Implement a Context-Aware Router (`t2-1`):** The WebSocket handler will be the intelligent link between the frontend and backend, managing the state (which container/session each user is connected to) for all concurrent clients.

4.  **Perform Interactive Verification (`t2-5`):** Once implemented, the entire system will be validated through a strict, interactive testing process led by the AI and verified by the user across multiple platforms (Mac, iPhone) to guarantee all requirements have been met.

## Key Planning Documents

*   **`FASE-2-MANGLER.md`**: The original analysis document detailing the initial implementation failure.
*   **`phase.yaml`**: Defines the high-level goal for Phase 2.
*   **`tasks/t2-1-websocket-handler.yaml`**: Defines the work for the **context-aware WebSocket router**.
*   **`tasks/t2-3-terminal-wrapper.yaml`**: Defines the work for the **concrete TypeScript frontend**.
*   **`tasks/t2-4-native-tmux-manager.yaml`**: Defines the work for the **LXC-integrated backend**.
*   **`tasks/t2-5-integration-testing.yaml`**: Defines the **final, interactive verification regime** to be executed upon completion of the other tasks.

## Next Steps

The next step is to begin the implementation work as detailed in tasks `t2-1`, `t2-3`, and `t2-4`. Once those are complete, the final verification in `t2-5` can be initiated.