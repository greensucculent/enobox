// go:build darwin
//  +build darwin

#include <stdlib.h>

int initKernel(const char *metalCode, const char *funcName);
void *addBuffer(int kernelId, int numBytes);
void run(int kernelId);