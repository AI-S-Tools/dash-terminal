// Quick test to verify xterm.js fix
const http = require('http');

console.log('🧪 Testing terminal fix...');

// Test 1: HTTP server is responding
http.get('http://localhost:8080/', (res) => {
    console.log('✅ HTTP server responding:', res.statusCode);

    let data = '';
    res.on('data', (chunk) => {
        data += chunk;
    });

    res.on('end', () => {
        // Test 2: Check if xterm.js script is in HTML
        const hasXtermJs = data.includes('cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.js');
        console.log(hasXtermJs ? '✅ xterm.js CDN found in HTML' : '❌ xterm.js CDN missing');

        // Test 3: Check if terminal.js is included
        const hasTerminalJs = data.includes('js/terminal.js');
        console.log(hasTerminalJs ? '✅ terminal.js included' : '❌ terminal.js missing');

        // Test 4: Check if terminal element exists
        const hasTerminalDiv = data.includes('id="terminal"');
        console.log(hasTerminalDiv ? '✅ Terminal div found' : '❌ Terminal div missing');

        console.log('\n🌐 Open http://localhost:8080 in browser and check console for:');
        console.log('   - typeof Terminal (should not be "undefined")');
        console.log('   - WebSocket connection');
        console.log('   - Terminal output appearing in #terminal div');
    });
}).on('error', (err) => {
    console.error('❌ HTTP test failed:', err.message);
});