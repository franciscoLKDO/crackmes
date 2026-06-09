#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <time.h>

#define SOCKET_PATH "/tmp/keygen.sock"

time_t (*real_time)(time_t *) = NULL;

time_t time(time_t *t) {
    int sock;
    struct sockaddr_un addr;
    time_t now;
    char *date_str;

    // 1. Create Socket
    sock = socket(AF_UNIX, SOCK_STREAM, 0);
    if (sock == -1) {
        perror("socket");
        return now;
    }

    // 2. Connect to Server
    memset(&addr, 0, sizeof(struct sockaddr_un));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, SOCKET_PATH, sizeof(addr.sun_path) - 1);

    if (connect(sock, (struct sockaddr*)&addr, sizeof(struct sockaddr_un)) == -1) {
        perror("connect");
        return now;
    }

    // 3. Get Date from real time lib
    real_time = dlsym(RTLD_NEXT, "time");
    now = real_time(t);
    // 4. Send Data
    send(sock, &now, sizeof(time_t), 0);

    close(sock);
    return now;
}