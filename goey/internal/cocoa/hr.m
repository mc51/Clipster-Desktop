#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GHR : NSView
- (void)drawRect:(NSRect)dirtyRect;
- (BOOL)isOpaque;
@end

@implementation GHR

- (void)drawRect:(NSRect)dirtyRect {
	NSSize size = [self frame].size;
	[[NSColor blackColor] set];
	[NSBezierPath
	    strokeLineFromPoint:NSMakePoint( 0, size.height / 2 )
	                toPoint:NSMakePoint( size.width, size.height / 2 )];
}

- (BOOL)isOpaque {
	return NO;
}

@end

void* hrNew( void* superview ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the horizontal rule.
	GHR* control = [[GHR alloc] init];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}
