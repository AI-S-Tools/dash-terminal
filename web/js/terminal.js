// Terminal.js - T2.3 xterm.js integration with WebSocket
class DashTerminal {
    constructor() {
        this.terminal = null;
        this.websocket = null;
        this.connected = false;
        this.statusIndicator = document.querySelector('.status-indicator');

        this.init();
    }

    init() {
        this.initTerminal();
        this.connectWebSocket();
    }

    initTerminal() {
        // Create xterm.js terminal instance
        this.terminal = new Terminal({
            theme: {
                background: '#1a1a1a',
                foreground: '#ffffff',
                cursor: '#ffffff',
                selection: '#ffffff33',
            },
            fontSize: 14,
            fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
            cursorBlink: true,
            rows: 24,
            cols: 80
        });

        // Open terminal in DOM
        this.terminal.open(document.getElementById('terminal'));

        // Show welcome message
        this.terminal.write('🚀 Dash Terminal PWA\r\n');
        this.terminal.write('Connecting to WebSocket...\r\n\r\n');

        // Handle user input - send to WebSocket when connected
        this.terminal.onData(data => {
            if (this.connected && this.websocket) {
                this.sendTerminalInput(data);
            }
        });
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        this.websocket = new WebSocket(wsUrl);

        this.websocket.onopen = () => {
            console.log('WebSocket connected');
            this.connected = true;
            this.updateStatus('connected', 'Connected');
            this.terminal.write('✅ WebSocket connected\r\n');

            // Request container list on connection
            this.requestContainerList();
        };

        this.websocket.onmessage = (event) => {
            this.handleWebSocketMessage(event.data);
        };

        this.websocket.onclose = () => {
            console.log('WebSocket disconnected');
            this.connected = false;
            this.updateStatus('disconnected', 'Disconnected');
            this.terminal.write('\r\n❌ WebSocket disconnected\r\n');

            // Attempt reconnection after 3 seconds
            setTimeout(() => {
                this.connectWebSocket();
            }, 3000);
        };

        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.terminal.write('❌ WebSocket error\r\n');
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

                case 'container_info':
                    this.handleContainerInfo(message.payload);
                    break;

                case 'terminal_output':
                    this.handleTerminalOutput(message.payload);
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
        this.terminal.write(`📟 ${payload.message}\r\n`);
    }

    handleContainerList(containers) {
        this.terminal.write('📦 Available containers:\r\n');
        containers.forEach(container => {
            const status = container.status === 'running' ? '🟢' : '🔴';
            this.terminal.write(`  ${status} ${container.name} (${container.status})\r\n`);
        });
        this.terminal.write('\r\n💡 Type container name to connect (T2.4 will implement)\r\n');
        this.terminal.write('$ ');
    }

    handleContainerInfo(container) {
        this.terminal.write(`📊 Container: ${container.name}\r\n`);
        this.terminal.write(`   Status: ${container.status}\r\n`);
        this.terminal.write(`   Type: ${container.type}\r\n\r\n`);
    }

    handleTerminalOutput(payload) {
        // T2.4 will implement actual terminal output
        this.terminal.write(payload.data);
    }

    handleError(error) {
        this.terminal.write(`❌ Error ${error.code}: ${error.message}\r\n`);
    }

    // Send terminal input to WebSocket
    sendTerminalInput(data) {
        const message = {
            type: 'terminal_input',
            payload: {
                pane_id: 'main', // T2.4 will implement proper pane management
                data: data
            }
        };

        this.websocket.send(JSON.stringify(message));
    }

    // Request list of containers
    requestContainerList() {
        const message = {
            type: 'container_list',
            payload: {}
        };

        this.websocket.send(JSON.stringify(message));
    }

    // Request container info
    requestContainerInfo(containerName) {
        const message = {
            type: 'container_info',
            payload: {
                container_name: containerName
            }
        };

        this.websocket.send(JSON.stringify(message));
    }

    // Update status indicator
    updateStatus(status, message) {
        this.statusIndicator.className = `status-indicator ${status}`;
        this.statusIndicator.textContent = message;
    }

    // Resize terminal
    resize() {
        if (this.terminal) {
            this.terminal.fit();
        }
    }
}

// Initialize terminal when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const dashTerminal = new DashTerminal();

    // Handle window resize
    window.addEventListener('resize', () => {
        dashTerminal.resize();
    });

    // Make terminal available globally for debugging
    window.dashTerminal = dashTerminal;
});