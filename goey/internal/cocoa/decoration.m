#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GDecoration : NSView {
	NSColor* fillColor;
	NSColor* strokeColor;
	NSSize borderRadius;
}
@property( nonatomic ) NSSize borderRadius;
@property( nonatomic, retain ) NSColor* fillColor;
@property( nonatomic, retain ) NSColor* strokeColor;
- (void)dealloc;
- (void)drawRect:(NSRect)dirtyRect;
- (BOOL)isOpaque;
@end

@implementation GDecoration

@synthesize borderRadius;
@synthesize fillColor;
@synthesize strokeColor;

- (void)dealloc {
	[fillColor release];
	[strokeColor release];
	[super dealloc];
}

- (void)drawRect:(NSRect)dirtyRect {
	NSRect frame = [self frame];
	frame.origin.x = 0;
	frame.origin.y = 0;
	if ( self.borderRadius.width > 0 ) {
		NSBezierPath* path =
		    [NSBezierPath bezierPathWithRoundedRect:frame
		                                    xRadius:self.borderRadius.width
		                                    yRadius:self.borderRadius.height];
		if ( self.fillColor ) {
			[self.fillColor set];
			[path fill];
		}
		if ( self.strokeColor ) {
			[self.strokeColor set];
			[path stroke];
		}
		//[path release];
	} else {
		if ( self.fillColor ) {
			[self.fillColor set];
			[NSBezierPath fillRect:frame];
		}
		if ( self.strokeColor ) {
			[self.strokeColor set];
			[NSBezierPath strokeRect:frame];
		}
	}
}

- (BOOL)isOpaque {
	return self.borderRadius.width <= 0;
}

@end

void* decorationNew( void* superview, nscolor_t fill, nscolor_t stroke,
                     nssize_t radius ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the button
	GDecoration* control = [[GDecoration alloc] init];
	decorationSetFillColor( control, fill );
	decorationSetStrokeColor( control, stroke );
	decorationSetBorderRadius( control, radius );

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

nssize_t decorationBorderRadius( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );

	NSSize size = [(GDecoration*)handle borderRadius];
	nssize_t rc = {size.width, size.height};
	return rc;
}

nscolor_t decorationFillColor( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );

	NSColor* clr = [(GDecoration*)handle fillColor];
	if ( !clr ) {
		nscolor_t clr = {0, 0, 0, 0};
		return clr;
	}

	CGFloat r, g, b, a;
	[clr getRed:&r green:&g blue:&b alpha:&a];
	nscolor_t ret = {r * 255, g * 255, b * 255, a * 255};
	return ret;
}

nscolor_t decorationStrokeColor( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );

	NSColor* clr = [(GDecoration*)handle strokeColor];
	if ( !clr ) {
		nscolor_t clr = {0, 0, 0, 0};
		return clr;
	}

	CGFloat r, g, b, a;
	[clr getRed:&r green:&g blue:&b alpha:&a];
	nscolor_t ret = {r * 255, g * 255, b * 255, a * 255};
	return ret;
}

static NSColor* createColor( nscolor_t clr ) {
	if ( clr.a == 0 ) {
		return NULL;
	}

	return [NSColor colorWithDeviceRed:clr.r / 255.0
	                             green:clr.g / 255.0
	                              blue:clr.b / 255.0
	                             alpha:clr.a / 255.0];
}

void decorationSetFillColor( void* handle, nscolor_t fill ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );
	[(GDecoration*)handle setFillColor:createColor( fill )];
}

void decorationSetStrokeColor( void* handle, nscolor_t stroke ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );
	[(GDecoration*)handle setStrokeColor:createColor( stroke )];
}

void decorationSetBorderRadius( void* handle, nssize_t r ) {
	assert( handle && [(id)handle isKindOfClass:[GDecoration class]] );
	NSSize radius = NSMakeSize( r.width, r.height );
	[(GDecoration*)handle setBorderRadius:radius];
}
