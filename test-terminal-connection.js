// Test WebSocket terminal connection
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8080/ws');

let containers = [];
let connected = false;

ws.on('open', function open() {
    console.log('✅ WebSocket connected');

    // Request container list (should trigger auto-connect)
    const containerListMessage = {
        type: 'container_list',
        payload: {}
    };

    ws.send(JSON.stringify(containerListMessage));
});

ws.on('message', function message(data) {
    try {
        const msg = JSON.parse(data);
        console.log(`📨 Received: ${msg.type}`, msg.payload);

        if (msg.type === 'container_list') {
            containers = msg.payload;
            console.log(`📦 Found ${containers.length} containers`);

            // Should auto-connect to first running container
            const running = containers.filter(c => c.status === 'running');
            if (running.length > 0) {
                console.log(`🔄 Should auto-connect to: ${running[0].name}`);

                // Test manual container_select message
                setTimeout(() => {
                    const selectMessage = {
                        type: 'container_select',
                        payload: {
                            container_name: running[0].name
                        }
                    };
                    console.log(`🚀 Sending container_select: ${running[0].name}`);
                    ws.send(JSON.stringify(selectMessage));
                }, 1000);
            }
        } else if (msg.type === 'status' && msg.payload.message && msg.payload.message.includes('Terminal session started')) {
            console.log('✅ Terminal session started successfully!');
            connected = true;

            // Test terminal input
            setTimeout(() => {
                const terminalInput = {
                    type: 'terminal_input',
                    payload: {
                        pane_id: 'main',
                        data: 'echo "Hello from test"\n'
                    }
                };
                console.log('⌨️  Sending test terminal input');
                ws.send(JSON.stringify(terminalInput));
            }, 1000);
        } else if (msg.type === 'terminal_output') {
            console.log(`🖥️  Terminal output: ${msg.payload.data.trim()}`);
        }

    } catch (error) {
        console.error('❌ Error parsing message:', error);
    }
});

ws.on('error', function error(err) {
    console.error('❌ WebSocket error:', err.message);
});

ws.on('close', function close() {
    console.log('🔌 WebSocket closed');
});

// Test timeout
setTimeout(() => {
    if (!connected) {
        console.log('❌ Test failed - terminal not connected after 10 seconds');
        process.exit(1);
    } else {
        console.log('✅ Test passed - terminal connection working!');
        process.exit(0);
    }
}, 10000);