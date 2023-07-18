// go:build darwin
//  +build darwin

#ifndef HEADER_PIPELINECACHE
#define HEADER_PIPELINECACHE

#include "pipeline.h"

int pipelineCache_cache(_pipeline *pipeline);
_pipeline *pipelineCache_retrieve(int pipelineId);

#endif