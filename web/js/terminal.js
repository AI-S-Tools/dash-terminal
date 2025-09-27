class DashTerminal {
    constructor() {
        this.terminal = null;
        this.websocket = null;
        this.connected = false;
        this.sessionId = null;
        this.statusIndicator = document.getElementById('status-indicator');

        this.init();
    }

    init() {
        this.initTerminal();
        this.createSimpleSession();
    }

    initTerminal() {
        const TerminalClass = window.Terminal || Terminal;
        this.terminal = new TerminalClass({
            theme: {
                background: '#1a1a1a',
                foreground: '#ffffff',
                cursor: '#ffffff',
                selection: '#ffffff33',
            },
            fontSize: 14,
            fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
            cursorBlink: true,
        });

        const fitAddon = new FitAddon.FitAddon();
        this.terminal.loadAddon(fitAddon);
        this.terminal.open(document.getElementById('terminal-container'));
        fitAddon.fit();

        window.addEventListener('resize', () => fitAddon.fit());

        this.terminal.onData(data => {
            if (this.connected && this.websocket) {
                this.websocket.send(data);
            }
        });

        this.terminal.onResize(size => {
            if (this.connected && this.websocket) {
                this.sendTerminalResize(size.cols, size.rows);
            }
        });
    }

    async createSimpleSession() {
        this.updateStatus('connecting', 'Creating session...');

        try {
            const response = await fetch('/api/sessions', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name: 'terminal' })
            });

            if (response.ok) {
                const session = await response.json();
                this.sessionId = session.id;
                this.connectToSession(session.id);
            } else {
                throw new Error(`Failed to create session: ${response.status}`);
            }
        } catch (error) {
            console.error('Failed to create session:', error);
            this.updateStatus('error', 'Failed to create session');
            this.terminal.write('\\x1b[31mError: Failed to create terminal session\\x1b[0m\\r\\n');
        }
    }

    connectToSession(sessionId) {
        this.updateStatus('connecting', 'Connecting to terminal...');

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?session=${sessionId}`;

        this.websocket = new WebSocket(wsUrl);

        this.websocket.onopen = () => {
            this.connected = true;
            this.updateStatus('connected', 'Terminal ready');
            this.terminal.focus();
        };

        this.websocket.onmessage = (event) => {
            if (event.data instanceof Blob) {
                // Binary PTY data
                const reader = new FileReader();
                reader.onload = () => {
                    this.terminal.write(new Uint8Array(reader.result));
                };
                reader.readAsArrayBuffer(event.data);
            } else {
                // JSON messages (resize, control) or text data
                try {
                    const message = JSON.parse(event.data);
                    this.handleControlMessage(message);
                } catch (e) {
                    // Plain text terminal data
                    this.terminal.write(event.data);
                }
            }
        };

        this.websocket.onclose = () => {
            this.connected = false;
            this.updateStatus('disconnected', 'Connection lost');
            this.terminal.write('\\r\\n\\x1b[31mConnection lost. Reconnecting...\\x1b[0m\\r\\n');

            // Auto-reconnect after 3 seconds
            setTimeout(() => {
                this.createSimpleSession();
            }, 3000);
        };

        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.updateStatus('error', 'Connection error');
        };
    }

    handleControlMessage(message) {
        switch (message.type) {
            case 'error':
                console.error('Session error:', message.message);
                this.terminal.write(`\\x1b[31mError: ${message.message}\\x1b[0m\\r\\n`);
                break;
            default:
                console.log('Unknown control message:', message.type);
        }
    }

    sendTerminalResize(cols, rows) {
        if (this.websocket && this.connected) {
            const message = {
                type: 'resize',
                width: cols,
                height: rows
            };
            this.websocket.send(JSON.stringify(message));
        }
    }

    updateStatus(status, message) {
        if (this.statusIndicator) {
            this.statusIndicator.className = `status ${status}`;
            this.statusIndicator.textContent = message;
        }
    }
}

// Initialize when page loads
document.addEventListener('DOMContentLoaded', () => {
    new DashTerminal();
});