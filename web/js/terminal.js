class DashTerminal {
    constructor() {
        this.terminal = null;
        this.websocket = null;
        this.connected = false;
        this.statusIndicator = document.getElementById('status-indicator');
        this.containerSelect = document.getElementById('container-select');

        this.init();
    }

    init() {
        this.initTerminal();
        this.connectWebSocket();
        this.containerSelect.addEventListener('change', (event) => {
            this.connectToContainer(event.target.value);
        });
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
                this.sendTerminalInput(data);
            }
        });

        this.terminal.onResize(size => {
            if (this.connected && this.websocket) {
                this.sendTerminalResize(size.cols, size.rows);
            }
        });
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        this.websocket = new WebSocket(wsUrl);

        this.websocket.onopen = () => {
            this.connected = true;
            this.updateStatus('connected', 'Connected');
            this.requestContainerList();
        };

        this.websocket.onmessage = (event) => {
            this.handleWebSocketMessage(event.data);
        };

        this.websocket.onclose = () => {
            this.connected = false;
            this.updateStatus('disconnected', 'Disconnected');
            setTimeout(() => this.connectWebSocket(), 3000);
        };

        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    handleWebSocketMessage(data) {
        try {
            const message = JSON.parse(data);

            switch (message.type) {
                case 'status':
                    this.handleStatusMessage(message.payload);
                    break;
                case 'container_list':
                    this.handleContainerList(message.payload);
                    break;
                case 'terminal_output':
                    this.terminal.write(message.payload.data);
                    break;
                case 'error':
                    this.handleError(message.payload);
                    break;
                default:
                    console.log('Unknown message type:', message.type);
            }
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    }

    handleStatusMessage(payload) {
        console.log('Status:', payload.message);
    }

    handleContainerList(containers) {
        this.containerSelect.innerHTML = '';
        const placeholder = document.createElement('option');
        placeholder.disabled = true;
        placeholder.selected = true;
        placeholder.textContent = 'Select a container...';
        this.containerSelect.appendChild(placeholder);

        containers.forEach(container => {
            const option = document.createElement('option');
            option.value = container.name;
            option.textContent = `${container.name} (${container.status})`;
            option.disabled = container.status !== 'running';
            this.containerSelect.appendChild(option);
        });
    }

    handleError(error) {
        this.terminal.write(`\x1b[31mError ${error.code}: ${error.message}\x1b[0m\r\n`);
    }

    sendTerminalInput(data) {
        const message = {
            type: 'terminal_input',
            payload: { data: data }
        };
        this.websocket.send(JSON.stringify(message));
    }

    sendTerminalResize(cols, rows) {
        const message = {
            type: 'terminal_resize',
            payload: { width: cols, height: rows }
        };
        this.websocket.send(JSON.stringify(message));
    }

    requestContainerList() {
        const message = { type: 'container_list', payload: {} };
        this.websocket.send(JSON.stringify(message));
    }

    connectToContainer(containerName) {
        const message = {
            type: 'container_select',
            payload: { container_name: containerName }
        };
        this.websocket.send(JSON.stringify(message));
    }

    updateStatus(status, message) {
        this.statusIndicator.className = `status-indicator ${status}`;
        this.statusIndicator.textContent = message;
    }
}

window.addEventListener('load', () => {
    if (typeof Terminal === 'undefined' || typeof FitAddon === 'undefined') {
        console.error('xterm.js or FitAddon not loaded');
        return;
    }
    window.dashTerminal = new DashTerminal();
});
