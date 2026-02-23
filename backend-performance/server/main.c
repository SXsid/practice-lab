#include <netinet/in.h> //http protocla and config
#include <stdio.h>      //for input out
#include <stdlib.h>     // funcion to handle meroy and core funtin
#include <string.h>     //all operation on char* DS
#include <sys/socket.h> //all the server / netowrk warpper
#include <unistd.h>     //sys call
// stdio.h        → printing logs
// stdlib.h       → memory management
// string.h       → handle buffers/messages
// unistd.h       → OS syscalls
// sys/socket.h   → create sockets
// netinet/in.h   → define IP + ports

#define PORT 8080
int main(int argc, char *argv[]) {
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
  listen(server_fd, 10);
}
