#ifndef GOEY_COCOA_LOOP_H
#define GOEY_COCOA_LOOP_H

// Cannot use std bool.  The builtin type _Bool does not play well with CGO.
// Need an alternate for the binding.
typedef unsigned bool_t;

extern void init( void );
extern void run( void );
extern void performOnMainThread( void );
extern void stop( void );
extern bool_t isMainThread( void );

extern void trace( char const* func );
#ifdef NTRACE
#define TRACE() ( (void)0 )
#else
#define TRACE() trace( __func__ )
#endif

#endif
