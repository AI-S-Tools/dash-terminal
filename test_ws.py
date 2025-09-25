#!/usr/bin/env python3

import websocket
import json
import time

def on_message(ws, message):
    print(f"Received: {message}")

    # Parse message
    try:
        msg = json.loads(message)
        print(f"Type: {msg['type']}")
        if msg['type'] == 'status' and msg['payload']['connected']:
            print("Connected! Sending container_list request...")
            # Send container list request
            request = {
                "type": "container_list",
                "payload": {}
            }
            ws.send(json.dumps(request))
        elif msg['type'] == 'container_list':
            containers = msg['payload']
            print(f"Found {len(containers)} containers:")
            for c in containers:
                print(f"  - {c['name']} ({c['status']}, {c['type']})")

            # Test container info for first container
            if containers:
                container_name = containers[0]['name']
                print(f"Requesting info for {container_name}...")
                info_request = {
                    "type": "container_info",
                    "payload": {
                        "container_name": container_name
                    }
                }
                ws.send(json.dumps(info_request))
        elif msg['type'] == 'container_info':
            container = msg['payload']
            print(f"Container info: {container}")

            # Test container select
            print(f"Testing container select for {container['name']}...")
            select_request = {
                "type": "container_select",
                "payload": {
                    "container_name": container['name']
                }
            }
            ws.send(json.dumps(select_request))
        elif msg['type'] == 'error':
            error = msg['payload']
            print(f"ERROR {error['code']}: {error['message']}")

    except json.JSONDecodeError:
        print(f"Invalid JSON: {message}")

def on_error(ws, error):
    print(f"WebSocket error: {error}")

def on_close(ws, close_status_code, close_msg):
    print("WebSocket connection closed")

def on_open(ws):
    print("WebSocket connection opened")

if __name__ == "__main__":
    websocket.enableTrace(False)
    ws = websocket.WebSocketApp("ws://localhost:8080/ws",
                              on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    # Run for 5 seconds
    def close_connection():
        time.sleep(5)
        ws.close()

    import threading
    threading.Thread(target=close_connection, daemon=True).start()

    ws.run_forever()