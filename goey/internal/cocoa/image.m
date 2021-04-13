#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

void* imageNewFromRGBA( uint8_t* imageData, int width, int height,
                        int stride ) {
	assert( imageData );
	assert( width > 0 && height > 0 );

	NSBitmapImageRep* imagerep =
	    [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:NULL
	                                            pixelsWide:width
	                                            pixelsHigh:height
	                                         bitsPerSample:8
	                                       samplesPerPixel:4
	                                              hasAlpha:YES
	                                              isPlanar:NO
	                                        colorSpaceName:NSDeviceRGBColorSpace
	                                           bytesPerRow:4 * width
	                                          bitsPerPixel:32];
	assert( imagerep );

	// Copy over the image data.
	if ( [imagerep bytesPerRow] == stride ) {
		// Check for overflow
		assert( width * height > width && width * height > height );
		assert( ( width * height * 4 ) > width * height );
		// Copy the data
		assert( [imagerep bitmapData] );
		memcpy( [imagerep bitmapData], imageData, width * height * 4 );
	} else {
		assert( false ); // not implemented
	}

	// Create the image
	NSImage* image = [[NSImage alloc] initWithSize:NSMakeSize( width, height )];
	assert( image );
	[image addRepresentation:imagerep];
	[imagerep release];
	return image;
}

void* imageNewFromGray( uint8_t* imageData, int width, int height,
                        int stride ) {
	assert( imageData );
	assert( width > 0 && height > 0 );

	NSBitmapImageRep* imagerep = [[NSBitmapImageRep alloc]
	    initWithBitmapDataPlanes:NULL
	                  pixelsWide:width
	                  pixelsHigh:height
	               bitsPerSample:8
	             samplesPerPixel:1
	                    hasAlpha:NO
	                    isPlanar:NO
	              colorSpaceName:NSDeviceWhiteColorSpace
	                 bytesPerRow:width
	                bitsPerPixel:8];
	assert( imagerep );

	// Copy over the image data.
	if ( [imagerep bytesPerRow] == stride ) {
		// Check for overflow
		assert( width * height > width && width * height > height );
		// Copy the data
		assert( [imagerep bitmapData] );
		memcpy( [imagerep bitmapData], imageData, width * height );
	} else {
		assert( false ); // not implemented
	}

	NSImage* image = [[NSImage alloc] initWithSize:NSMakeSize( width, height )];
	assert( image );
	[image addRepresentation:imagerep];
	[imagerep release];
	return image;
}

void imageClose( void* image ) {
	assert( image && [(id)image isKindOfClass:[NSImage class]] );

	[(NSImage*)image release];
}

void* imageviewNew( void* superview, void* image ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( image && [(id)image isKindOfClass:[NSImage class]] );

	// Create the control
	NSImageView* control = [[NSImageView alloc] init];
	[control setImage:(NSImage*)image];
	[control setImageScaling:NSImageScaleAxesIndependently];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

int imageviewImageWidth( void* control ) {
	assert( control && [(id)control isKindOfClass:[NSImageView class]] );

	NSArray* reps = [[(NSImageView*)control image] representations];
	NSImageRep* rep = [reps objectAtIndex:0];
	return [rep pixelsWide];
}

int imageviewImageHeight( void* control ) {
	assert( control && [(id)control isKindOfClass:[NSImageView class]] );

	NSArray* reps = [[(NSImageView*)control image] representations];
	NSImageRep* rep = [reps objectAtIndex:0];
	return [rep pixelsHigh];
}

int imageviewImageDepth( void* control ) {
	assert( control && [(id)control isKindOfClass:[NSImageView class]] );

	NSArray* reps = [[(NSImageView*)control image] representations];
	assert( reps );
	NSImageRep* rep = [reps objectAtIndex:0];
	assert( rep );

	if ( [rep colorSpaceName] == NSDeviceWhiteColorSpace ) {
		int bits = [rep bitsPerSample];
		if ( [rep hasAlpha] ) {
			bits *= 2;
		}
		return bits;
	}

	assert( [rep colorSpaceName] == NSDeviceRGBColorSpace );
	int bits = [rep bitsPerSample];
	bits *= [rep hasAlpha] ? 4 : 3;
	return bits;
}

void imageviewImageData( void* control, void* data ) {
	assert( control && [(id)control isKindOfClass:[NSImageView class]] );

	NSArray* reps = [[(NSImageView*)control image] representations];
	assert( reps );
	NSBitmapImageRep* rep = [reps objectAtIndex:0];
	assert( rep );

	NSInteger const bytesPerRow = [rep bytesPerRow];
	int i;
	for ( i = 0; i < [rep pixelsHigh]; i++ ) {
		unsigned char* dst = (unsigned char*)( data ) + i * bytesPerRow;
		unsigned char* src = [rep bitmapData] + i * bytesPerRow;
		memcpy( dst, src, bytesPerRow );
	}
}

void imageviewSetImage( void* control, void* image ) {
	assert( control && [(id)control isKindOfClass:[NSImageView class]] );
	assert( image );

	[(NSImageView*)control setImage:(NSImage*)image];
}
