
PIXC_ROOT=../..
CC=$(PIXC_ROOT)/bin/cc-arm-raspian
AR=$(PIXC_ROOT)/install/binutils/bin/ar

all: arm-cc-eg

libhello.so: hello.o
	$(CC) -shared -o $@ $<


libhello.a: hello.o
	$(AR) rcs $@ $<

.o: .c
	$(CC) -c $<

arm-cc-eg: main.o libhello.so
	$(CC) -o $@ main.o -lhello -L.

clean:
	$(RM) *.o *.a *.so