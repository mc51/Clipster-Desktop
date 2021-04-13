#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the text view
	NSText* control = [[NSText alloc] init];
	[control setDrawsBackground:NO];
	textSetText( control, text );
	[control setEditable:NO];

	// Add the control as the view for the window
	[(NSView*)superview addSubview:control];

	return control;
}

int textAlignment( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );
	switch ( [(NSText*)handle alignment] ) {
	default:
	case NSTextAlignmentLeft:
		return 0;

	case NSTextAlignmentCenter:
		return 1;

	case NSTextAlignmentRight:
		return 2;

	case NSTextAlignmentJustified:
		return 3;
	}
}

static NSSize measureTextSize( NSFont* font, NSString* text, NSSize size ) {
	assert( font && text );

	// Create objects.
	NSTextContainer* container =
	    [[NSTextContainer alloc] initWithContainerSize:size];
	NSLayoutManager* manager = [[NSLayoutManager alloc] init];
	NSTextStorage* storage = [[NSTextStorage alloc] initWithString:text];

	// Configure the objects
	[container setLineFragmentPadding:0.0];
	[manager addTextContainer:container];
	[storage addLayoutManager:manager];
	[storage addAttribute:NSFontAttributeName
	                value:font
	                range:NSMakeRange( 0, [storage length] )];

	// Force layout
	[manager glyphRangeForTextContainer:container];
	// Find the size
	size = [manager usedRectForTextContainer:container].size;

	// Release objects
	[manager release];
	[container release];
	[storage release];

	return size;
}

int textEightyEms( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	static char const eightyEms[] = "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm"
	                                "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm";

	// Documentation is unclear, but assuming that the returned string is
	// in the autopool.
	NSString* ems = [NSString stringWithCString:eightyEms length:80];
	// Measure the width of the string.
	NSSize size = measureTextSize( [(NSText*)handle font], ems,
	                               NSMakeSize( FLT_MAX, FLT_MAX ) );

	return size.width;
}

int textMinHeight( void* handle, int width ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	NSSize size =
	    measureTextSize( [(NSText*)handle font], [(NSText*)handle string],
	                     NSMakeSize( width, FLT_MAX ) );
	return size.height;
}

int textMinWidth( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	NSSize size =
	    measureTextSize( [(NSText*)handle font], [(NSText*)handle string],
	                     NSMakeSize( FLT_MAX, FLT_MAX ) );
	return size.width;
}

void textSetText( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );
	assert( text );

	NSString* nsText = [[NSString alloc] initWithUTF8String:text];
	[(NSText*)handle setText:nsText];
	[nsText release];
}

void textSetAlignment( void* handle, int align ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	switch ( align ) {
	default:
	case 0:
		[(NSText*)handle setAlignment:NSTextAlignmentLeft];
		break;
	case 1:
		[(NSText*)handle setAlignment:NSTextAlignmentCenter];
		break;
	case 2:
		[(NSText*)handle setAlignment:NSTextAlignmentRight];
		break;
	case 3:
		[(NSText*)handle setAlignment:NSTextAlignmentJustified];
		break;
	}
}

char const* textText( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	NSString* text = [(NSText*)handle text];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}
