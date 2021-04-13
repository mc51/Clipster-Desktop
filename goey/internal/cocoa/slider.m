#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GSlider : NSSlider
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)onchange;
@end

@implementation GSlider

- (void)onchange {
	double s = [self doubleValue];
	sliderOnChange( self, s );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		sliderOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		sliderOnBlur( self );
	}
	return rc;
}

@end

void* sliderNew( void* superview, double min, double value, double max ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the slider
	NSSlider* control = [[GSlider alloc] init];
	//[control setSliderType:NSSliderTypeLinear];
	[control setTarget:control];
	[control setAction:@selector( onchange )];
	sliderUpdate( control, min, value, max );

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

double sliderMax( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSSlider class]] );
	return [(NSSlider*)handle maxValue];
}

double sliderMin( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSSlider class]] );
	return [(NSSlider*)handle minValue];
}

double sliderValue( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSSlider class]] );
	return [(NSSlider*)handle doubleValue];
}

void sliderUpdate( void* handle, double min, double value, double max ) {
	assert( handle && [(id)handle isKindOfClass:[NSSlider class]] );

	[(NSSlider*)handle setMinValue:min];
	[(NSSlider*)handle setMaxValue:max];
	[(NSSlider*)handle setDoubleValue:value];
}
