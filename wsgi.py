import server

if __name__ == "__main__":
    server.socketio.run(server.app, host='0.0.0.0', port=5000)