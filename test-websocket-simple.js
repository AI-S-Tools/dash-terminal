// Simple WebSocket test without external dependencies
const WebSocket = require('ws');

console.log('🧪 Testing WebSocket terminal connection...');

const ws = new WebSocket('ws://localhost:8080/ws');

ws.on('open', function() {
    console.log('✅ WebSocket connected');

    // Request container list
    ws.send(JSON.stringify({
        type: 'container_list',
        payload: {}
    }));
});

ws.on('message', function(data) {
    const msg = JSON.parse(data);
    console.log(`📨 ${msg.type}:`, msg.payload);

    if (msg.type === 'container_list') {
        const running = msg.payload.filter(c => c.status === 'running');
        if (running.length > 0) {
            console.log(`🚀 Auto-connecting to: ${running[0].name}`);

            // Connect to first running container
            ws.send(JSON.stringify({
                type: 'container_select',
                payload: { container_name: running[0].name }
            }));
        }
    } else if (msg.type === 'status') {
        if (msg.payload.message && msg.payload.message.includes('Terminal session started')) {
            console.log('✅ TERMINAL SESSION STARTED!');

            // Test command
            ws.send(JSON.stringify({
                type: 'terminal_input',
                payload: {
                    pane_id: 'main',
                    data: 'echo "Hello from real terminal"\\n'
                }
            }));
        }
    } else if (msg.type === 'terminal_output') {
        console.log(`🖥️  Terminal output: ${msg.payload.data.trim()}`);
        if (msg.payload.data.includes('Hello from real terminal')) {
            console.log('✅ REAL TERMINAL WORKING!');
            process.exit(0);
        }
    }
});

ws.on('error', function(err) {
    console.error('❌ WebSocket error:', err.message);
    process.exit(1);
});

ws.on('close', function() {
    console.log('🔌 WebSocket closed');
});

// Timeout
setTimeout(() => {
    console.log('⏰ Test timeout');
    process.exit(1);
}, 10000);