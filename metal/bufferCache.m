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
int bufferCache_cache(_buffer *buffer) {
  NSCAssert(buffer != nil, @"Missing buffer to cache");

  int bufferId = 0;

  @synchronized(bufferCache) {
    if (bufferCache == nil) {
      bufferCache_init();
    }

    // We cannot store the struct into the cache directly. Instead, we need to
    // encode it into an NSValue and store that.
    NSValue *value = [NSValue valueWithBytes:buffer objCType:@encode(_buffer)];
    [bufferCache addObject:value];

    // A buffer Id is its 1-based index in the array.
    bufferId = [bufferCache count];
  }

  return bufferId;
}

// Retrieve a buffer (block of memory) from the buffer cache.
_buffer *bufferCache_retrieve(int bufferId) {
  _buffer *buffer = nil;

  buffer = malloc(sizeof(_buffer));
  NSCAssert(buffer != nil, @"Failed to initialize new buffer");

  @synchronized(bufferCache) {
    NSCAssert(bufferId >= 1, @"Invalid buffer Id %d", bufferId);
    NSCAssert(bufferId <= [bufferCache count], @"Invalid buffer Id %d",
              bufferId);

    // A buffer Id is a buffer's 1-based index in the cache. We need to convert
    // it into a 0-based index to retrieve it from the cache.
    int index = bufferId - 1;

    // Retrieve and decode the encoded struct.
    NSValue *value = bufferCache[index];
    [value getValue:buffer];
  }

  return buffer;
}