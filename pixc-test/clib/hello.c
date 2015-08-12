#include <stdio.h>

#include "hello.h"

void msg(int ac, char **av) {
	printf("messages:\n");
	for (int i = 0; i < ac; i++) {
		printf("%d: %s\n",i,av[i]);
	}
}
