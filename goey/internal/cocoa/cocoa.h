#ifndef GOEY_COCOA_H
#define GOEY_COCOA_H

#include <stdint.h>

// Cannot use std bool.  The builtin type _Bool does not play well with CGO.
// Need an alternate for the binding.

typedef unsigned bool_t;

typedef struct nssize_tag {
	int32_t width;
	int32_t height;
} nssize_t;

typedef struct nscolor_tag {
	uint8_t r, g, b, a;
} nscolor_t;

/* Message dialog */
extern void messageDialog( void* window, char const* text, char const* title, char icon );
extern char const* openPanel( void* window, char const* dir, char const* base );
extern char const* savePanel( void* window, char const* dir, char const* base );
extern void dialogSendKey( unsigned key );

extern void trace( char const* func );
#ifdef NTRACE
#define TRACE() ( (void)0 )
#else
#define TRACE() trace( __func__ )
#endif

/* Window */
extern void* windowNew( char const* title, unsigned width, unsigned height );
extern void windowClose( void* handle );
extern nssize_t windowContentSize( void* handle );
extern nssize_t windowFrameSize( void* handle );
extern void* windowContentView( void* handle );
extern void windowMakeFirstResponder( void* handle, void* control );
extern void windowScreenshot( void* handle, void* data, int32_t width, int32_t height );
extern void windowSetContentSize( void* handle, int width, int height );
extern void windowSetMinSize( void* handle, int width, int height );
extern void windowSetIconImage( void* handle, void* nsimage );
extern void windowSetScrollVisible( void* handle, bool_t horz, bool_t vert );
extern void windowSetTitle( void* handle, char const* title );
extern char const* windowTitle( void* handle );

/* View */
extern void viewSetFrame( void* handle, int x, int y, int dx, int dy );
extern void viewClose( void* handle );

/* Control */
extern bool_t controlIsEnabled( void* handle );
extern void controlSetEnabled( void* handle, bool_t value );
extern nssize_t controlIntrinsicContentSize( void* handle );
extern bool_t controlMakeFirstResponder( void* handle );
extern void controlSendKey( void* handle, unsigned keyASCII );

/* Button */
extern void* buttonNew( void* superview, char const* title );
extern void* buttonNewCheck( void* window, char const* title, bool_t value );
extern bool_t buttonIsDefault( void* handle );
extern void buttonPerformClick( void* handle );
extern bool_t buttonState( void* handle );
extern void buttonSetDefault( void* handle, bool_t value );
extern void buttonSetState( void* handle, bool_t checked );
extern char const* buttonTitle( void* handle );
extern void buttonSetTitle( void* handle, char const* title );

/* Decoration */
extern void* decorationNew( void* superview, nscolor_t fill, nscolor_t stroke,
                            nssize_t radius );
extern nssize_t decorationBorderRadius( void* control );
extern nscolor_t decorationFillColor( void* control );
extern nscolor_t decorationStrokeColor( void* control );
extern void decorationSetBorderRadius( void* control, nssize_t radius );
extern void decorationSetFillColor( void* control, nscolor_t fill );
extern void decorationSetStrokeColor( void* control, nscolor_t stroke );

/* HR */
extern void* hrNew( void* superview );

/* PopUpButton */
extern void* popupbuttonNew( void* superview );
extern void popupbuttonAddItem( void* control, char const* text );
extern char const* popupbuttonItemAtIndex( void* control, int index );
extern int popupbuttonNumberOfItems( void* control );
extern void popupbuttonRemoveAllItems( void* control );
extern void popupbuttonSetValue( void* control, int index, bool_t unset );
extern int popupbuttonValue( void* control );

/* ProgressIndicator */
extern void* progressNew( void* superview, double min, double value,
                          double max );
extern double progressMax( void* handle );
extern double progressMin( void* handle );
extern double progressValue( void* handle );
extern void progressUpdate( void* handle, double min, double value,
                            double max );

/* Slider */
extern void* sliderNew( void* superview, double min, double value, double max );
extern double sliderMax( void* handle );
extern double sliderMin( void* handle );
extern double sliderValue( void* handle );
extern void sliderUpdate( void* handle, double min, double value, double max );

/* TabView */
extern void* tabviewNew( void* superview );
extern void tabviewAddItem( void* control, char const* text );
extern char const* tabviewItemAtIndex( void* control, int index );
extern int tabviewNumberOfItems( void* control );
extern void tabviewRemoveItemAtIndex( void* control, int index );
extern void tabviewSelectItem( void* control, int index );
extern void tabviewSetItemAtIndex( void* control, int index, char const* text );
extern void* tabviewContentView( void* control, int index );
extern nssize_t tabviewContentInsets( void* control );

/* Text */
extern void* textNew( void* superview, char const* text );
extern int textAlignment( void* handle );
extern int textEightyEms( void* handle );
extern int textMinHeight( void* handle, int width );
extern int textMinWidth( void* handle );
extern void textSetText( void* handle, char const* text );
extern void textSetAlignment( void* handle, int align );
extern char const* textText( void* handle );

/* TextField */
extern void* textfieldNew( void* superview, char const* text, bool_t password );
extern bool_t textfieldIsEditable( void* handle );
extern bool_t textfieldIsPassword( void* handle );
extern char const* textfieldPlaceholder( void* handle );
extern void textfieldSetEditable( void* handle, bool_t value );
extern void textfieldSetValue( void* handle, char const* text );
extern void textfieldSetPlaceholder( void* handle, char const* text );
extern char const* textfieldValue( void* handle );

/* IntField */
extern void* intfieldNew( void* superview, int64_t value, int64_t min,
                          int64_t max );
extern void intfieldClose( void* handle );
extern bool_t intfieldIsEditable( void* handle );
extern int64_t intfieldMax( void* handle );
extern int64_t intfieldMin( void* handle );
extern char const* intfieldPlaceholder( void* handle );
extern void intfieldSetEditable( void* handle, bool_t value );
extern void intfieldSetValue( void* handle, int64_t value, int64_t min,
                              int64_t max );
extern void intfieldSetPlaceholder( void* handle, char const* text );
extern int64_t intfieldValue( void* handle );
extern void intfieldSetFrame( void* handle, int x, int y, int dx, int dy );

/* TextView */
extern void* textviewNew( void* superview, char const* text );
extern void textviewSetValue( void* handle, char const* text );

/* Image */
extern void* imageNewFromRGBA( uint8_t* imageData, int width, int height,
                               int stride );
extern void* imageNewFromGray( uint8_t* imageData, int width, int height,
                               int stride );
extern void imageClose( void* handle );

/* ImageView */
extern void* imageviewNew( void* superview, void* image );
extern int imageviewImageWidth( void* control );
extern int imageviewImageHeight( void* control );
extern int imageviewImageDepth( void* control );
extern void imageviewImageData( void* control, void* data );
extern void imageviewSetImage( void* control, void* image );

#endif
