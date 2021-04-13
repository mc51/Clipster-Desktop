#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GIntField : NSTextField <NSTextFieldDelegate> {
	NSStepper* stepper;
}
@property( nonatomic, retain ) NSStepper* stepper;
- (void)dealloc;
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)controlTextDidChange:(NSNotification*)obj;
@end

@implementation GIntField

@synthesize stepper;

- (void)dealloc {
	[stepper release];
	[super dealloc];
}

- (void)controlTextDidChange:(NSNotification*)obj {
	NSInteger value = [self integerValue];
	double minValue = [[self stepper] minValue];
	double maxValue = [[self stepper] maxValue];
	if ( value > maxValue ) {
		value = maxValue;
		[self setIntegerValue:value];
	} else if ( value < minValue ) {
		value = minValue;
		[self setIntegerValue:value];
	}
	[[self stepper] setIntegerValue:value];
	intfieldOnChange( self, value );
}

- (void)onclick {
	NSInteger value = [[self stepper] integerValue];
	[self setIntegerValue:value];
	intfieldOnChange( self, value );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		intfieldOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		intfieldOnBlur( self );
	}
	return rc;
}

@end

void* intfieldNew( void* superview, int64_t value, int64_t min, int64_t max ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( min <= value && value <= max );

	// Create the button
	GIntField* control = [[GIntField alloc] init];
	intfieldSetValue( control, value, min, max );
	[control setEditable:YES];
	//[control setUsesSingleLineMode:YES];
	[control setDelegate:control];
	assert( sizeof( NSInteger ) >= sizeof( int64_t ) );
	[control setIntegerValue:value];

	NSStepper* stepper = [[NSStepper alloc] init];
	[stepper setMaxValue:max];
	[stepper setMinValue:min];
	[stepper setIntegerValue:value];
	[stepper setTarget:control];
	[stepper setAction:@selector( onclick )];
	[control setStepper:stepper];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];
	[(NSView*)superview addSubview:stepper];
	[stepper release];

	// Return handle to the control
	return control;
}

void intfieldClose( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	[(GIntField*)handle setStepper:nil];
	viewClose( handle );
}

bool_t intfieldIsEditable( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	return [(GIntField*)handle isEditable];
}

// Because Cocoa uses double (or float64 in Go) to store the range, it cannot
// keep full precision for values near the minimum or maximum of the int64
// range.  This function is just for Props, and is used to adjust the int64
// to get a match.
static int64_t toInt64( double scale ) {
	int64_t a = scale;
	if ( (double)a == scale ) {
		return a;
	}
	if ( (double)( a - 1 ) == scale ) {
		return a - 1;
	}
	printf( "mismatch...%ld %f\n", a, scale );
	return a;
}

int64_t intfieldMax( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	NSStepper* stepper = [(GIntField*)handle stepper];
	return toInt64( [stepper maxValue] );
}

int64_t intfieldMin( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	NSStepper* stepper = [(GIntField*)handle stepper];
	return toInt64( [stepper minValue] );
}

char const* intfieldPlaceholder( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	NSString* text = [[(GIntField*)handle cell] placeholderString];
	return [text cStringUsingEncoding:NSUTF8StringEncoding];
}

void intfieldSetEditable( void* handle, bool_t value ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	[(GIntField*)handle setEditable:value];
}

void intfieldSetValue( void* handle, int64_t value, int64_t min, int64_t max ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );
	assert( min <= value && value <= max );

	NSStepper* stepper = [(GIntField*)handle stepper];
	[stepper setMaxValue:max];
	[stepper setMinValue:min];

	if ( value != [(GIntField*)handle integerValue] ) {
		[(GIntField*)handle setIntegerValue:value];
		[[(GIntField*)handle stepper] setIntegerValue:value];
	}
}

void intfieldSetPlaceholder( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );
	assert( text );

	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[[(GIntField*)handle cell] setPlaceholderString:title];
	[title release];
}

int64_t intfieldValue( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );

	return [(GIntField*)handle integerValue];
}

void intfieldSetFrame( void* handle, int x, int y, int dx, int dy ) {
	assert( handle && [(id)handle isKindOfClass:[GIntField class]] );
	assert( dx >= 0 && dy >= 0 );

	NSRect frame = [[(NSView*)handle superview] frame];
	frame = NSMakeRect( x, frame.size.height - y - dy, dx - 16, dy );
	[(GIntField*)handle setFrame:frame];
	[(GIntField*)handle display];
	frame.origin.x = x + dx - 16;
	frame.size.width = 16;
	NSStepper* tmp = [(GIntField*)handle stepper];
	[tmp setFrame:frame];
	[tmp display];
}
