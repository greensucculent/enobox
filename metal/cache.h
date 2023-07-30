// go:build darwin
//  +build darwin

#ifndef HEADER_CACHE
#define HEADER_CACHE

int cache_cache(void *item);
void *cache_retrieve(int cacheId);

#endif