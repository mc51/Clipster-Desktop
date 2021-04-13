#include "cocoa.h"
#import <Foundation/NSThread.h>
#include <stdio.h>

void trace( char const* func ) {
	printf( "%s\t%p\n", func, [NSThread currentThread] );
	fflush( stdout );
}
