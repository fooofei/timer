#coding=utf-8

'''
the butter lib is from  http://blitz.works/butter/file/tip

DEP :  you can install a higher version of cffi by install rpm file which download from
https://cbs.centos.org/koji/buildinfo?buildID=20864

WARNING: the butter lib has <linux/memfd.h>, which is avaliabale on linux kernel 3.17,
if you fail python setup.py install. you can comment line memfd in setup.py , and
setup again.

'''

import os
import sys

from time import sleep
from time import time
from datetime import datetime as DateTime
from select import epoll
from select import EPOLLIN
from butter import timerfd
from random import randint

def log(start, s):
    t = time()-start
    t = int(t)
    print('{} s {}'.format(t, s))

def example():

    t = timerfd.Timer(timerfd.CLOCK_MONOTONIC,
                      timerfd.TFD_NONBLOCK | timerfd.TFD_CLOEXEC) #

    t.after(3,0).every(3,0)
    t.update()

    efd = epoll()
    efd.register(t.fileno(), EPOLLIN)
    start = time()

    while True:
        log(start, 'into poll')
        evs = efd.poll(-1, 10)
        log(start, 'over poll')
        for fileno, ev in evs:
            assert fileno == t.fileno()
            if ev & EPOLLIN:
                r = t.read_event()
                log(start, 'read {} from {}'.format(r, fileno))
                wk = randint(0, 3)
                log(start, 'goto work {}'.format(wk))
                sleep(wk)
                log(start, 'work done')
            print('---------------------------------------------')




if __name__ == '__main__':
    example()
