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
                if (strstr(buffer, "\r\n\r\n")) {
                        break;
                }
        }
        // printf("%.*s", (int)bytes_read, buffer);
        char *http_response = "HTTP/1.1 200 OK\r\n"
                              "Content-Type: text/plain\r\n"
                              "Content-Length: 13\r\n"
                              "Connection: close\r\n"
                              "\r\n"
                              "Hello world!\n";
        int wn = write(clietn_fd, http_response, strlen(http_response));
        if (wn < 0) {
                perror("write");
                return -1;
        }
        return 0;
}

int main() {
        // handling singlas
        signal(SIGINT, handlesignal);
        signal(SIGTERM, handlesignal);
        // geting a socket for our server
        int server_fd = socket(AF_INET, SOCK_STREAM, 0);
        if (server_fd == -1) {
                perror("Socket Failed");
                exit(EXIT_FAILURE);
        };
        int opt = 1;
        //  set socket behviour
        if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) < 0) {
                perror("setsockopt");
                exit(EXIT_FAILURE);
        }
        // netwrok lyaer connectin with the socket
        struct sockaddr_in addr = {0};
        addr.sin_family = AF_INET;         // ipv4
        addr.sin_port = htons(PORT);       // correct the endeinness
        addr.sin_addr.s_addr = INADDR_ANY; // 0.0.0.0 (bad everything in you pulbic network has access to it )
        if (bind(server_fd, (struct sockaddr *)&addr, sizeof(addr)) == -1) {
                perror("bind");
                exit(EXIT_FAILURE);
        }
        // make both ack and sync array for the server socket
        // listen wiht 10 backlog
        if (listen(server_fd, 10) < 0) {
                perror("listen");
                exit(EXIT_FAILURE);
        }
        printf("we are waiting for a connection \n");
        // start reading the for the ack array
        while (keep_running) {
                struct sockaddr_in client_add = {0};
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
