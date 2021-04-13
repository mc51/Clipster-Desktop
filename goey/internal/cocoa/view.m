#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void viewClose( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSView class]] );

	[(NSView*)handle removeFromSuperview /* WithoutNeedingDisplay? */];
	[(NSView*)handle release];
}

void viewSetFrame( void* handle, int x, int y, int dx, int dy ) {
	assert( handle && [(id)handle isKindOfClass:[NSView class]] );
	assert( dx >= 0 && dy >= 0 );

	NSRect frame = [[(NSView*)handle superview] frame];
	frame = NSMakeRect( x, frame.size.height - y - dy, dx, dy );
	[(NSView*)handle setFrame:frame];
	[(NSView*)handle display];
}
