// go:build darwin
//  +build darwin

#ifndef HEADER_BUFFERCACHE
#define HEADER_BUFFERCACHE

#import "buffer.h"
#import <Metal/Metal.h>

int bufferCache_cache(_buffer *buffer);
_buffer *bufferCache_retrieve(int bufferId);

#endif