// go:build darwin
//  +build darwin

#ifndef HEADER_BUFFER
#define HEADER_BUFFER

#import <Metal/Metal.h>

// Functions that must be called once for every buffer used as an argument to a
// metal function
int buffer_newBuffer(int numBytes);
void *buffer_retrieveBuffer(int bufferId);

#endif