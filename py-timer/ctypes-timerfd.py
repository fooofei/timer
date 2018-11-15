#coding=utf-8

__all__ = [
	"TFD_CLOEXEC",
	"TFD_NONBLOCK",
	"TFD_TIMER_ABSTIME",

	"CLOCK_REALTIME",
	"CLOCK_MONOTONIC",

	"timespec",
	"itimerspec",
	"Timer",
]


from ctypes import Structure
from ctypes import c_long
from ctypes import CDLL
from ctypes.util import find_library as c_find_library
from ctypes import get_errno
from os import strerror
from ctypes import pointer
from os import read as _read
from os import close as _close

TFD_CLOEXEC         = 0o02000000
TFD_NONBLOCK        = 0o00004000

TFD_TIMER_ABSTIME   = 0x00000001

CLOCK_REALTIME  = 0
CLOCK_MONOTONIC = 1


libc = CDLL(c_find_library("c"), use_errno=True)


def errcheck(result, func, argu):
    if result <0:
        errno = get_errno()
        raise OSError(errno, strerror(errno))
    return result

class timespec(Structure):
    _fields_ = [
        ('tv_sec', c_long),
        ('tv_nsec', c_long),
    ]
    def _str(self):
        s = 'tv_sec={0}'.format(self.tv_sec)
        ns = 'tv_nsec={0}'.format(self.tv_nsec)
        return '{0} {1}'.format(s, ns)
    def __str__(self):
        return self._str()
    def __repr__(self):
        return 'timerfd.timespec({0})'.format(self._str())

class itimerspec(Structure):
    _fields_=[
        ('it_interval', timespec),
        ('it_value', timespec),
    ]

    def _str(self):
        v1 = 'it_interval={0}'.format(self.it_interval)
        v2 = 'it_value={0}'.format(self.it_value)
        return '{0} {1}'.format(v1, v2)

    def __str__(self):
        return self._str()
    def __repr__(self):
        return 'timerfd.itimerspec({0})'.format(self._str())



class Timer(object):
        def __init__(self, timerspecval=None, clock_id = CLOCK_MONOTONIC,
                     flags=TFD_CLOEXEC):
            if timerspecval:
                self._timerspec = timerspecval
            else:
                self._timerspec = itimerspec()
            self._fd = libc.timerfd_create(clock_id, flags)

        def every(self, seconds= None, nano_seconds=None):
            if seconds is not None:
                self._timerspec.it_interval.tv_sec = seconds
            if nano_seconds is not None:
                self._timerspec.it_interval.tv_nsec = nano_seconds
            return self

        def after(self, seconds= None, nano_seconds=None):
            if seconds is not None:
                self._timerspec.it_value.tv_sec = seconds
            if nano_seconds is not None:
                self._timerspec.it_value.tv_nsec = nano_seconds
            return self

        def offset(self, seconds=None, nano_seconds=None):
            return self.after(seconds, nano_seconds)

        def repeat(self, seconds=None, nano_seconds=None):
            return self.every(seconds, nano_seconds)

        def fileno(self):
            return self._fd

        def get_current(self):
            cur = itimerspec()
            rc = libc.timerfd_gettime(self.fileno(), pointer(cur))
            if rc != 0:
                errno = get_errno()
                raise OSError(errno, strerror(errno))
            return cur

        def update(self, absolute=False):
            flags = TFD_TIMER_ABSTIME if absolute else 0
            old_timer = itimerspec()
            rc = libc.timerfd_settime(self.fileno(), flags,
                                      pointer(self._timerspec), pointer(old_timer))
            if rc != 0:
                errno = get_errno()
                raise OSError(errno, strerror(errno))
            return old_timer

        def close(self):
            _close(self.fileno())
            self._fd = -1
        def readev(self):
            data = _read(self.fileno(), 8)
            return data

        def __del__(self):
            self.close()

        def __repr__(self):
            v = '<{} (fd ={} after=({} s, {} ns) every=({} s, {} ns))>'.format(
                self.__class__.__name__,
                self.fileno(),
                self._timerspec.it_value.tv_sec,self._timerspec.it_value.tv_nsec,
                self._timerspec.it_interval.tv_sec,self._timerspec.it_interval.tv_nsec
            )
            return v

        def __str__(self):
            return self.__repr__()

def log(start, s):
    from time import time as Time
    t = Time()-start
    t = int(t)
    print('{} s {}'.format(t, s))

def example():
    from time import sleep
    from select import epoll
    from select import EPOLLIN
    from random import randint
    from time import time as Time



    sched = Timer(clock_id=CLOCK_MONOTONIC,
                  flags=TFD_NONBLOCK | TFD_CLOEXEC)

    sched.after(3,0).every(3,0)
    sched.update(absolute=False)
    efd = epoll()
    efd.register(sched.fileno(), EPOLLIN)
    start= Time()

    while True:
        log(start, 'into poll')
        evs = efd.poll(-1,10)
        log(start, 'over poll')
        for fileno, ev in evs:
            assert fileno == sched.fileno()
            if ev & EPOLLIN:
                r = sched.readev()
                log(start, 'read {} from {}'.format(r, fileno))
                wk = randint(0,3)
                log(start, 'goto work {}'.format(wk))
                sleep(wk)
                log(start, 'work done')
            print('---------------------------------------------')




if __name__ == '__main__':
    example()
