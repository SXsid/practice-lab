#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
        if (argc != 3) {
                printf("Usages:\n\t./client <server_addr> <port>");
                return 0;
        }
        int ServerPort = atoi(argv[2]);
        int fd = socket(AF_INET, SOCK_STREAM, 0);
        if (fd == -1) {
                perror("socket");
                return -1;
        }
        struct sockaddr_in addr = {0};
        addr.sin_family = AF_INET;
        addr.sin_addr.s_addr = inet_addr(argv[1]);
        addr.sin_port = htons(ServerPort);
        // connect to server adder to our client socket fd
        int res = connect(fd, (struct sockaddr *)&addr, sizeof(addr));
        if (res == -1) {
                perror("connect");
                close(fd);
                return 0;
        }
        close(fd);

        return 0;
}
