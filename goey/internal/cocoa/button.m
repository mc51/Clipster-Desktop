#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GButton : NSButton
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)onclick;
- (void)onchange;
@end

@implementation GButton

- (void)onclick {
	buttonOnClick( self );
}

- (void)onchange {
	NSInteger s = [self state];
	buttonOnChange( self, s == NSOnState );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		buttonOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		buttonOnBlur( self );
	}
	return rc;
}

@end

void* buttonNew( void* superview, char const* title ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( title );

	// Create the button
	GButton* control = [[GButton alloc] init];
	buttonSetTitle( control, title );
	[control setTarget:control];
	[control setAction:@selector( onclick )];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void* buttonNewCheck( void* superview, char const* title, bool_t value ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( title );

	// Create the button
	GButton* control = [[GButton alloc] init];
	buttonSetTitle( control, title );
	[control setTarget:control];
	[control setAction:@selector( onchange )];
	[control setButtonType:NSSwitchButton];
	[control setState:value];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void buttonPerformClick( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	[[(GButton*)handle cell] performClick:nil];
}

bool_t buttonIsDefault( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	return [[(GButton*)handle keyEquivalent] length] != 0;
}

bool_t buttonState( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	NSInteger s = [(GButton*)handle state];
	return s == NSOnState;
}

void buttonSetDefault( void* handle, bool_t isdefault ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	NSString* ke = 0;
	if ( isdefault ) {
		ke = [[[NSString alloc] initWithBytesNoCopy:"\r"
		                                     length:1
		                                   encoding:NSUTF8StringEncoding
		                               freeWhenDone:NO] autorelease];
	}
	return [(GButton*)handle setKeyEquivalent:ke];
}

void buttonSetState( void* handle, bool_t state ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	NSInteger s = state ? NSOnState : NSOffState;
	[(GButton*)handle setState:s];
}

char const* buttonTitle( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );

	NSString* title = [(GButton*)handle title];
	return [title cStringUsingEncoding:NSUTF8StringEncoding];
}

void buttonSetTitle( void* handle, char const* title ) {
	assert( handle && [(id)handle isKindOfClass:[GButton class]] );
	assert( title );

	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];
	[(GButton*)handle setTitle:nsTitle];
	[nsTitle release];
}
