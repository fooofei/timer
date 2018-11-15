

// gcc timer-timerfd.c -o timer-timerfd

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include <sys/timerfd.h>
#include <stdbool.h>
#include <signal.h>

typedef unsigned long long uint64;

static bool g_force_quit = false;

static void signal_handler(int signum) {
    fprintf(stdout, "Received signal %d", signum);
    if (signum == SIGINT || signum == SIGTERM || signum == SIGALRM) {
        switch (signum)
        {
        case SIGINT: fprintf(stdout, " SIGINT");g_force_quit=true; break;
        case SIGTERM: fprintf(stdout, " SIGTERM"); g_force_quit=true; break;
        case SIGALRM: fprintf(stdout, " SIGALRM"); break;
        default:
            break;
        }
    }
    fprintf(stdout, "\n");
    fflush(stdout);
}

static void main_setup_signals(void) {
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);
    signal(SIGALRM, signal_handler);
}


static void timer(){

    int tfd = 0;
    struct itimerspec utmr;
    struct itimerspec otmr;
    memset(&utmr, 0, sizeof(utmr));
    memset(&otmr, 0, sizeof(otmr));
    int rc;
    time_t start=0;
    int wk=0;

    // 只有 timerfd 无法工作
    tfd = timerfd_create(CLOCK_MONOTONIC, 0); // TFD_NONBLOCK
    if(tfd <=0){
        fprintf(stderr, "fail timerfd_create errno=%d %s", errno, strerror(errno));
    }
    utmr.it_value.tv_sec = 3; // first timeout time
    utmr.it_interval.tv_sec = 3;
    rc = timerfd_settime(tfd, 0,
        &utmr, &otmr);
    if(rc !=0){
        fprintf(stderr, "fail timerfd_settime errno=%d %s", errno, strerror(errno));
    }
    //
    // 同样的，与 clock_nanosleep 一样，如果某次任务超过了心跳周期，
    // 在这个时间点之后的任务整体时间会延迟
    time(&start);
    for(;!g_force_quit;){
        uint64 v=0;
        printf("%llu into read\n", time(NULL)-start); fflush(stdout);
        rc = read(tfd, &v, sizeof v);
        printf("%llu over read rc=%d v=%llu\n", time(NULL)-start,rc, v); fflush(stdout);

        wk = rand()%6;
        printf("%llu begin work, will take %llu sec\n", time(NULL)-start, wk); fflush(stdout);
        usleep(wk * 1000 * 1000);
        printf("%llu done work \n", time(NULL)-start);
        printf("--------------------------------------------\n");
    }
    printf("%llu end of timer\n", time(NULL)-start);
    close(tfd);
}


int main(){



    timer();
    return 0;
}
