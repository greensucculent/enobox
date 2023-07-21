// go:build darwin
//  +build darwin

#include "pipelineCache.h"

NSMutableArray *pipelineCache = nil;

// Initialize the pipeline cache.
void pipelineCache_init() {
  pipelineCache = [[NSMutableArray alloc] init];
  NSCAssert(pipelineCache != nil, @"Failed to initialize the pipeline cache");
}

// Add a pipeline to the pipeline cache.
int pipelineCache_cache(_pipeline *pipeline) {
  NSCAssert(pipeline != nil, @"Missing pipeline");

  int pipelineId = 0;

  @synchronized(pipelineCache) {
    if (pipelineCache == nil) {
      pipelineCache_init();
    }

    // We cannot store the struct into the cache directly. Instead, we need to
    // encode it into an NSValue and store that.
    NSValue *value = [NSValue valueWithBytes:pipeline
                                    objCType:@encode(_pipeline)];
    [pipelineCache addObject:value];

    // A pipeline Id is its 1-based index in the array.
    pipelineId = [pipelineCache count];
  }

  return pipelineId;
}

// Retrieve a pipeline from the pipeline cache.
_pipeline *pipelineCache_retrieve(int pipelineId) {
  _pipeline *pipeline = nil;

  pipeline = pipeline_alloc();

  @synchronized(pipelineCache) {
    NSCAssert(pipelineId >= 1, @"Invalid pipeline Id %d", pipelineId);
    NSCAssert(pipelineId <= [pipelineCache count], @"Invalid pipeline Id %d",
              pipelineId);

    // A pipeline Id is a pipeline's 1-based index in the cache. We need to
    // convert it into a 0-based index to retrieve it from the cache.
    int index = pipelineId - 1;

    // Retrieve and decode the encoded struct.
    NSValue *value = pipelineCache[index];
    [value getValue:pipeline];
  }

  return pipeline;
}