#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <ctype.h> // for tolower

bool_t controlIsEnabled( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	return [(NSControl*)handle isEnabled];
}

void controlSetEnabled( void* handle, bool_t value ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	[(NSControl*)handle setEnabled:value];
}

nssize_t controlIntrinsicContentSize( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	// Note that accessing the cell is deprecated, but GNUstep does not have
	// the newer methods needed to gather this information.
	NSCell* cell = [(NSControl*)handle cell];
	NSSize size = [cell cellSize];

	// Return the values
	nssize_t ret = {size.width, size.height};
	return ret;
}

bool_t controlMakeFirstResponder( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );
	assert( ![(NSControl*)handle refusesFirstResponder] );

	NSWindow* window = [(NSControl*)handle window];
	assert( window );
	return [window makeFirstResponder:(NSControl*)handle];
}

void controlSendKey( void* handle, unsigned keyASCII ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	// Return is used for the enter key on MacOS, not new-line as on other
	// platforms.
	if ( keyASCII== '\n' ) {
		keyASCII = '\r';
	}

	char const bytes = keyASCII;
	char const bytes2 = tolower( bytes );

	NSString* characters =
	    [[NSString alloc] initWithBytes:&bytes
	                             length:1
	                           encoding:NSASCIIStringEncoding];
	assert( characters );
	NSString* characters2 =
	    [[NSString alloc] initWithBytes:&bytes2
	                             length:1
	                           encoding:NSASCIIStringEncoding];
	assert( characters2 );

	NSTimeInterval timestamp = [[NSProcessInfo processInfo] systemUptime];

	NSEvent* evt = [NSEvent keyEventWithType:NSKeyDown
	                                location:NSZeroPoint
	                           modifierFlags:0
	                               timestamp:0
	                            windowNumber:0
	                                 context:nil
	                              characters:characters
	             charactersIgnoringModifiers:characters2
	                               isARepeat:NO
	                                 keyCode:0];
	assert( evt );
	[characters release];
	[characters2 release];

	if ( [(NSControl*)handle currentEditor] ) {
		[[(NSControl*)handle currentEditor] keyDown:evt];
	} else {
		[(NSControl*)handle keyDown:evt];
	}
	//[evt release];
}