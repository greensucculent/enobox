// go:build darwin
//  +build darwin

#import "buffer.h"

// Allocate a block of memory accessible to both the CPU and GPU that is large
// enough to hold the number of bytes specified, cache it, and return a buffer
// Id that can be used to retrieve the memory from the cache.
int buffer_newBuffer(int numBytes) {
  // Create a new buffer. At this point, we're only allocating a run of bytes
  // and don't need to declare any type or meaning to them. We're using
  // storageModeShared so both the CPU and GPU can access the buffers.
  // TODO: look into MTLStorageModeManaged
  id<MTLBuffer> buffer =
      [device newBufferWithLength:numBytes
                          options:MTLResourceStorageModeShared];
  NSCAssert(buffer != nil, @"Failed to create buffer");

  // Add the buffer to the buffer cache and return its unique Id.
  return bufferCache_cache(buffer);
}

// Retrieve a buffer from the cache.
void *buffer_retrieveBuffer(int bufferId) {
  // Retrieve the buffer from the cache.
  id<MTLBuffer> buffer = bufferCache_retrieve(bufferId);
  NSCAssert(buffer != nil, @"Failed to retrieve buffer from cache");

  // Return a reference to the block of memory. This keeps the buffer
  // wrapping type opaque to the consumer.
  return buffer.contents;
}