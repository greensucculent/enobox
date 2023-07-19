// go:build darwin
//  +build darwin

#ifndef HEADER_BUFFERCACHE
#define HEADER_BUFFERCACHE

#import <Metal/Metal.h>

int bufferCache_cache(id<MTLBuffer> buffer);
id<MTLBuffer> bufferCache_retrieve(int bufferId);

#endif