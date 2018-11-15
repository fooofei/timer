


// gcc timer-clock_nanosleep.c -o timer-nano

#include <stdio.h>
#include <signal.h>
#include <time.h>
#include <unistd.h>
#include <stdbool.h>
#include <errno.h>
#include <stdint.h>




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


// 定时任务
static void timer(){
    struct timespec res;
    time_t start;
    time_t now;
    time(&start);
    clock_gettime(CLOCK_REALTIME, &res);
    res.tv_sec += 3;
    int wk=0; // 限制： 超时任务不能 >3 ，比如用了4s，那么sleep立马返回，
    // 下一次sleep是以此为基线的下一个3s，会导致整体往后延迟
    for(;!g_force_quit;){

        // 含义是要 sleep 到哪个绝对时间
        for(;0 != clock_nanosleep(CLOCK_REALTIME, TIMER_ABSTIME, &res, &res) 
            && errno == EINTR && !g_force_quit;);
        clock_gettime(CLOCK_REALTIME, &res);
        res.tv_sec += 3; // 相对于这个时间的 3s 即便后面拖延了几秒 但是也能跳跃过去
        time(&now);
        printf("hit at %d ", now-start);
        time_t t1;
        time(&t1);
        wk = rand()%5;
        wk *= 1000*1000;
        usleep(wk);
        time_t t2;
        time(&t2);
        printf("work take %d(s)\n", t2-t1);
    }
}


int main(){

    timer();

    return 0;
}
