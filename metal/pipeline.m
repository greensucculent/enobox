// go:build darwin
//  +build darwin

#include "pipeline.h"

// Initialize a new pipeline.
_pipeline *pipeline_newPipeline() {
  _pipeline *pipeline = nil;

  pipeline = malloc(sizeof(_pipeline));
  NSCAssert(pipeline != nil, @"Failed to initialize new pipeline");

  return pipeline;
}