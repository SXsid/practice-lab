#include <netinet/in.h> //http protocla and config
#include <signal.h>
#include <stdio.h>      //for input out
#include <stdlib.h>     // funcion to handle meroy and core funtin
#include <string.h>     //all operation on char* DS
#include <sys/socket.h> //all the server / netowrk warpper
#include <unistd.h>     //sys call

#define PORT 8080
#define MAX_REQUEST_SIZE 8192

volatile sig_atomic_t keep_running = 1;
void handlesignal() {
        keep_running = 0;
}
int handleConnection(int clietn_fd) {
        char buffer[MAX_REQUEST_SIZE];
        size_t bytes_read = 0;
        while (1) {
                int n = read(clietn_fd, buffer + bytes_read, sizeof(buffer) - bytes_read);
                if (n < 0) {
                        if (!keep_running)
                                break;
                        perror("read");
                        return -1;
                }
                if (n == 0) {
                        break;
                }
                bytes_read += n;
                buffer[bytes_read] = '\0';
                if ((int)bytes_read > MAX_REQUEST_SIZE) {
                        printf("request is too large");
                        return -1;
                }
                if (strstr(buffer, "\r\n\r\n")) {
                        break;
                }
        }
        printf("%.*s", (int)bytes_read, buffer);
        char *http_response = "HTTP/1.1 200 OK\r\n"
                              "Content-Type: text/plain\r\n"
                              "Content-Length: 12\r\n"
                              "Connection: close\r\n"
                              "\r\n"
                              "Hello world!";
        int wn = write(clietn_fd, http_response, strlen(http_response));
        if (wn < 0) {
                perror("write");
                return -1;
        }
        return 0;
}

int main() {
        signal(SIGINT, handlesignal);
        signal(SIGTERM, handlesignal);
        int server_fd = socket(AF_INET, SOCK_STREAM, 0);
        if (server_fd == 0) {
                perror("Socket Failed");
                exit(EXIT_FAILURE);
        };
        int opt = 1;
        if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) < 0) {
                perror("setsockopt");
                exit(EXIT_FAILURE);
        }
        struct sockaddr_in addr;
        addr.sin_family = AF_INET;
        addr.sin_port = htons(PORT);
        addr.sin_addr.s_addr = INADDR_ANY;
        if (bind(server_fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
                perror("bind");
                exit(EXIT_FAILURE);
        }
        if (listen(server_fd, 10) < 0) {
                perror("listen");
                exit(EXIT_FAILURE);
        }
        printf("we are waiting for a connection \n");
        while (keep_running) {
                struct sockaddr_in client_add;
                socklen_t client_addr_len = sizeof(client_add);
                int clientFd = accept(server_fd, (struct sockaddr *)&client_add, &client_addr_len);
                if (clientFd < 0) {
                        if (!keep_running)
                                break;
                        perror("accept");
                        continue;
                }
                printf("a client connected wiht %d\n", client_add.sin_port);
                if (handleConnection(clientFd) < 0) {
                        perror("Error handling client");
                        close(clientFd);
                        printf("Conecction closed%d\n", client_add.sin_port);
                        continue;
                }
                close(clientFd);
                printf("Conecction closed%d\n", client_add.sin_port);
        }
        close(server_fd);
}
