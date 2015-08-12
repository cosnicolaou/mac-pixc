package main

// #include <stdlib.h>
//
// #include "clib/hello.h"
//
//	static char** newCArray(int size) {
//      return calloc(sizeof(char*), size);
//  }
//
//  static void setCArray(char **a, int n, char *s) {
//      a[n] = s;
//  }
//
//  static void freeCArray(char **a, int size) {
//      for (int i = 0; i < size; i++) {
//          free(a[i]);
//      }
//      free(a);
//}
import "C"

func Msg(m []string) {
	l := C.int(len(m))
	argv := C.newCArray(l)
	for i, v := range m {
		C.setCArray(argv, C.int(i), C.CString(v))
	}
	defer C.freeCArray(argv, l)
	C.msg(l, argv)
}
