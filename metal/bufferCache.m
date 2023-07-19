// go:build darwin
//  +build darwin

#import "bufferCache.h"

NSMutableArray *bufferCache = nil;

// Initialize the buffer cache.
void bufferCache_init() {
  bufferCache = [[NSMutableArray alloc] init];
  NSCAssert(bufferCache != nil, @"Failed to initialize the buffer cache");
}

// Add a buffer to the buffer cache.
int bufferCache_cache(id<MTLBuffer> buffer) {
  NSCAssert(buffer != nil, @"Missing buffer to cache");

  int bufferId = 0;

  @synchronized(bufferCache) {
    if (bufferCache == nil) {
      bufferCache_init();
    }

    [bufferCache addObject:buffer];

    // A buffer Id is its 1-based index in the array.
    bufferId = [bufferCache count];
  }

  return bufferId;
}

// Retrieve a buffer (block of memory) from the buffer cache.
id<MTLBuffer> bufferCache_retrieve(int bufferId) {
  id<MTLBuffer> buffer = nil;

  @synchronized(bufferCache) {
    if (bufferId < 1 || bufferId > [bufferCache count]) {
      NSCAssert(false, @"Invalid buffer Id %d", bufferId);
    }

    // A buffer Id is a buffer's 1-based index in the cache. We need to convert
    // it into a 0-based index to retrieve it from the cache.
    int index = bufferId - 1;

    buffer = bufferCache[index];
    NSCAssert(buffer != nil, @"Failed to find buffer in cache");
  }

  return buffer;
}
