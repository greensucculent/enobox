// go:build darwin
//  +build darwin

#include "kernel.h"

_kernel **kernels;
int numKernels;

// Set up a new kernel.
_kernel *newKernel() {
  _kernel *kernel;

  kernel = malloc(sizeof(_kernel));
  kernel->buffers = [[NSMutableArray alloc] init];

  return kernel;
}

// Save a kernel into the global space.
int registerKernel(_kernel *kernel) {
  int kernelId;

  kernelId = numKernels;

  numKernels++;
  kernels = realloc(kernels, sizeof(_kernel *) * (numKernels));

  kernels[kernelId] = kernel;

  return kernelId;
}

// Load a kernel from the global space.
_kernel *retrieveKernel(int kernelId) {
  if (kernelId < 0 || kernelId >= numKernels) {
    return nil;
  }

  return kernels[kernelId];
}
