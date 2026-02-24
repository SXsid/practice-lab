
#include <netinet/in.h> //http protocla and config
#include <stdio.h>      //for input out
#include <stdlib.h>     // funcion to handle meroy and core funtin
#include <string.h>     //all operation on char* DS
#include <sys/socket.h> //all the server / netowrk warpper
#include <unistd.h>     //sys call


#define PORT 8080
#define MAX_REQUEST_SIZE 8192


int handleConnection(int clietn_fd) {
        char buffer[MAX_REQUEST_SIZE];
        size_t bytes_read = 0;
        while (1) {
                int n = read(clietn_fd, buffer + bytes_read, sizeof(buffer) - bytes_read);
                if (n < 0) {
                        perror("read");
                        return -1;
                }
                if (n == 0) {
                        // client closed
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
        return 0;
}


int main() {
        // create a socket
        int server_fd = socket(AF_INET, SOCK_STREAM, 0);
        if (server_fd == 0) {
                perror("Socket Failed");
                exit(EXIT_FAILURE);
        };

        // bind a socket
        struct sockaddr_in addr;
        addr.sin_family = AF_INET;
        addr.sin_port = htons(PORT);
        addr.sin_addr.s_addr = INADDR_ANY;
        if (bind(server_fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
                perror("bind");
                exit(EXIT_FAILURE);
        }
        // now we have synce que and accept que(client who has sended ack)
        if (listen(server_fd, 10) < 0) {
                perror("listen");
                exit(EXIT_FAILURE);
        }
        printf("we are waiting for a connection \n");
        while (1) {
                struct sockaddr_in client_add;
                socklen_t client_addr_len = sizeof(client_add);
                int clientFd = accept(server_fd, (struct sockaddr *)&client_add, &client_addr_len);
                if (clientFd < 0) {
                        perror("accept");
                        continue;
                }
                printf("a client connected wiht %d\n", client_add.sin_port);
                if (handleConnection(clientFd) < 0) {
                        perror("Error handling client");
                        close(clientFd);
                        continue;
                }
                close(clientFd);
        }
        close(server_fd);
}
