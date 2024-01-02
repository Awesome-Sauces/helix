import socket
import sys

def create_response(status_code, content):
    response = f"HTTP/1.1 {status_code}\r\n"
    response += "Content-Type: text/html\r\n"
    response += f"Content-Length: {len(content)}\r\n"
    response += "\r\n"
    response += content
    return response

def main():
    host = '127.0.0.1'
    port = 8080

    try:
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_socket.bind((host, port))
        server_socket.listen(1)
        print(f"Server listening on http://{host}:{port}")
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

    while True:
        client_socket, client_address = server_socket.accept()
        request_data = client_socket.recv(1024).decode('utf-8')
        
        if not request_data:
            continue

        # Basic routing
        if request_data.startswith("GET / HTTP/1.1"):
            response_content = "<h1>Hello, World!</h1>"
            response = create_response("200 OK", response_content)
        else:
            response_content = "<h1>Not Found</h1>"
            response = create_response("404 Not Found", response_content)
        
        client_socket.send(response.encode('utf-8'))
        client_socket.close()

if __name__ == "__main__":
    main()
