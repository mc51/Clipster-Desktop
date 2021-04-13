#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* progressNew( void* superview, double min, double value, double max ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the button
	NSProgressIndicator* control = [[NSProgressIndicator alloc] init];
	progressUpdate( control, min, value, max );

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

double progressMax( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSProgressIndicator class]] );
	return [(NSProgressIndicator*)handle maxValue];
}

double progressMin( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSProgressIndicator class]] );
	return [(NSProgressIndicator*)handle minValue];
}

double progressValue( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSProgressIndicator class]] );
	return [(NSProgressIndicator*)handle doubleValue];
}

void progressUpdate( void* handle, double min, double value, double max ) {
	assert( handle && [(id)handle isKindOfClass:[NSProgressIndicator class]] );

	[(NSProgressIndicator*)handle setMinValue:min];
	[(NSProgressIndicator*)handle setMaxValue:max];
	[(NSProgressIndicator*)handle setDoubleValue:value];
}
