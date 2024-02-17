import socket

# Define a function to handle client requests
def handle_client(client_socket):
    while True:
        data = client_socket.recv(1024)
        if not data:
            break
        
        # Process the received data (e.g., perform simulation actions)
        # Simulated action: Echo the received message back to the client
        client_socket.send(data)

# Main function to start the backend server
def main():
    # Define host and port for the server
    HOST = '127.0.0.1'
    PORT = 12345

    # Create a TCP socket
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((HOST, PORT))
    server_socket.listen(5)

    print('Backend server started. Waiting for connections...')

    while True:
        client_socket, addr = server_socket.accept()
        print('Connected to', addr)

        # Handle client requests in a separate thread
        handle_client(client_socket)

    server_socket.close()

if __name__ == "__main__":
    main()
