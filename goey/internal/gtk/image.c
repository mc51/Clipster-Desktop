#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

void *mountImage( void *parent, unsigned char const *data, int width,
                  int height, int rowStride )
{
    assert( parent );
    assert( data );

    GdkPixbuf *pixbuf =
        gdk_pixbuf_new_from_data( data, GDK_COLORSPACE_RGB, TRUE, 8, width,
                                  height, rowStride, NULL, NULL );
    assert( pixbuf );

    GtkWidget *widget = gtk_image_new_from_pixbuf( pixbuf );
    assert( widget );
    g_object_unref( pixbuf );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void *imageUpdate( void *widget, unsigned char const *data, int width,
                   int height, int rowStride )
{
    assert( widget );
    assert( data );

    GdkPixbuf *pixbuf =
        gdk_pixbuf_new_from_data( data, GDK_COLORSPACE_RGB, TRUE, 8, width,
                                  height, rowStride, NULL, NULL );
    assert( pixbuf );

    gtk_image_set_from_pixbuf( GTK_IMAGE( widget ), pixbuf );
    g_object_unref( pixbuf );
}

unsigned imageColorSpace( void *widget )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    return gdk_pixbuf_get_colorspace( pixbuf );
}

bool imageHasAlpha( void *widget )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    return gdk_pixbuf_get_has_alpha( pixbuf );
}

unsigned imageImageWidth( void *widget )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    return gdk_pixbuf_get_width( pixbuf );
}

unsigned imageImageHeight( void *widget )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    return gdk_pixbuf_get_height( pixbuf );
}

unsigned imageImageStride( void *widget )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    return gdk_pixbuf_get_rowstride( pixbuf );
}

char const *imageImageData( void *widget, size_t *length )
{
    assert( widget );

    GdkPixbuf *pixbuf = gtk_image_get_pixbuf( GTK_IMAGE( widget ) );
    assert( pixbuf );

    guint length2 = 0;
    guchar *data = gdk_pixbuf_get_pixels_with_length( pixbuf, &length2 );
    *length = length2;
    return data;
}
