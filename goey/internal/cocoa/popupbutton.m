#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GPopUpButton : NSPopUpButton
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)onchange;
@end

@implementation GPopUpButton

- (void)onchange {
	NSInteger s = [self indexOfSelectedItem];
	popupbuttonOnChange( self, s );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		popupbuttonOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		popupbuttonOnBlur( self );
	}
	return rc;
}

@end

void* popupbuttonNew( void* superview ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the button
	GPopUpButton* control = [[GPopUpButton alloc] init];
	[control setTarget:control];
	[control setAction:@selector( onchange )];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void popupbuttonAddItem( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );
	assert( text );

	NSString* nstext = [[NSString alloc] initWithUTF8String:text];
	[(GPopUpButton*)handle addItemWithTitle:nstext];
	[nstext release];
}

char const* popupbuttonItemAtIndex( void* handle, int index ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );

	NSString* text = [(GPopUpButton*)handle itemTitleAtIndex:index];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}

int popupbuttonNumberOfItems( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );

	return [(GPopUpButton*)handle numberOfItems];
}

void popupbuttonRemoveAllItems( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );

	[(GPopUpButton*)handle removeAllItems];
}

void popupbuttonSetValue( void* handle, int index, bool_t unset ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );

	if ( unset ) {
		[(GPopUpButton*)handle selectItem:nil];
	} else {
		[(GPopUpButton*)handle selectItemAtIndex:index];
	}
}

int popupbuttonValue( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GPopUpButton class]] );

	return [(GPopUpButton*)handle indexOfSelectedItem];
}
