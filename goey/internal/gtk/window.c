#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "thunks.h"

void ondestroy_cb( GtkWidget *widget, gpointer user_data )
{
    onDestroy( widget );
}

static void ondeleteevent_cb( GtkWidget *widget, gpointer user_data )
{
    onDeleteEvent( widget );
}

static void onsizeallocate_cb( GtkWidget *widget, GdkRectangle *rectangle,
                               gpointer user_data )
{
    onSizeAllocate( widget, rectangle->width, rectangle->height );
}

void *mountWindow( char const *text )
{
    GtkWidget *window = gtk_window_new( GTK_WINDOW_TOPLEVEL );
    gtk_window_set_title( GTK_WINDOW( window ), text );
    gtk_container_set_border_width( GTK_CONTAINER( window ), 0 );

    GtkWidget *scroll = gtk_scrolled_window_new( NULL, NULL );
    gtk_scrolled_window_set_policy( GTK_SCROLLED_WINDOW( scroll ),
                                    GTK_POLICY_NEVER, GTK_POLICY_NEVER );
    gtk_container_add( GTK_CONTAINER( window ), scroll );

    GtkWidget *layout = gtk_layout_new( NULL, NULL );
    gtk_container_add( GTK_CONTAINER( scroll ), layout );

    g_signal_connect( window, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    g_signal_connect( window, "delete-event", G_CALLBACK( ondeleteevent_cb ),
                      NULL );
    g_signal_connect( window, "size-allocate", G_CALLBACK( onsizeallocate_cb ),
                      NULL );

    return window;
}

windowsize_t windowSize( void *window )
{
    windowsize_t size;
    gtk_window_get_size( GTK_WINDOW( window ), &size.width, &size.height );
    return size;
}

void *windowScrolledWindow( void *window )
{
    GtkWidget *ss = gtk_bin_get_child( GTK_BIN( window ) );
    assert( ss );
    return ss;
}

void *windowLayout( void *window )
{
    GtkWidget *ss = gtk_bin_get_child( GTK_BIN( window ) );
    assert( ss );
    GtkWidget *layout = gtk_bin_get_child( GTK_BIN( ss ) );
    assert( layout );
    return layout;
}

void windowSetLayoutSize( void *window, unsigned width, unsigned height )
{
    GtkWidget *ss = gtk_bin_get_child( GTK_BIN( window ) );
    assert( ss );
    GtkWidget *layout = gtk_bin_get_child( GTK_BIN( ss ) );
    assert( layout );
    gtk_layout_set_size( GTK_LAYOUT( layout ), width, height );
}

char const *windowTitle( void *window )
{
    return gtk_window_get_title( GTK_WINDOW( window ) );
}

void windowSetTitle( void *window, char const *text )
{
    gtk_window_set_title( GTK_WINDOW( window ), text );
}

unsigned windowVScrollbarWidth( void *window )
{
    GtkWidget *oldChild = gtk_bin_get_child( GTK_BIN( window ) );
    assert( oldChild );
    gtk_container_remove( GTK_CONTAINER( window ), oldChild );

    GtkWidget *sb = gtk_scrollbar_new( GTK_ORIENTATION_VERTICAL, NULL );
    assert( sb );

    gtk_container_add( GTK_CONTAINER( window ), sb );
    gtk_widget_show( sb );
    int min, nominal;
    gtk_widget_get_preferred_width( sb, &min, &nominal );
    gtk_widget_destroy( sb );
    gtk_container_add( GTK_CONTAINER( window ), oldChild );
    return nominal;
}

unsigned windowHScrollbarHeight( void *window )
{
    GtkWidget *oldChild = gtk_bin_get_child( GTK_BIN( window ) );
    assert( oldChild );
    gtk_container_remove( GTK_CONTAINER( window ), oldChild );

    GtkWidget *sb = gtk_scrollbar_new( GTK_ORIENTATION_HORIZONTAL, NULL );
    assert( sb );

    gtk_container_add( GTK_CONTAINER( window ), sb );
    gtk_widget_show( sb );
    int min, nominal;
    gtk_widget_get_preferred_height( sb, &min, &nominal );
    gtk_widget_destroy( sb );
    gtk_container_add( GTK_CONTAINER( window ), oldChild );
    return nominal;
}

void windowShowScrollbars( void *window, bool horz, bool vert )
{
    GtkWidget *ss = gtk_bin_get_child( GTK_BIN( window ) );
    assert( ss );
    gtk_scrolled_window_set_policy(
        GTK_SCROLLED_WINDOW( ss ), horz ? GTK_POLICY_ALWAYS : GTK_POLICY_NEVER,
        vert ? GTK_POLICY_ALWAYS : GTK_POLICY_NEVER );
}

void windowShow( void *window )
{
    gtk_widget_show_all( GTK_WIDGET( window ) );
}

void windowSetDefaultSize( void *window, int width, int height )
{
    assert( window );
    gtk_window_set_default_size( GTK_WINDOW( window ), width, height );
}

void windowSetIcon( void *window, unsigned char const *data, int width,
                    int height, int rowStride )
{
    assert( window );

    if ( !data ) {
        gtk_window_set_icon( GTK_WINDOW( window ), NULL );
        return;
    }

    GdkPixbuf *pixbuf =
        gdk_pixbuf_new_from_data( data, GDK_COLORSPACE_RGB, TRUE, 8, width,
                                  height, rowStride, NULL, NULL );
    assert( pixbuf );

    gtk_window_set_icon( GTK_WINDOW( window ), pixbuf );
    g_object_unref( pixbuf );
}

void *windowScreenshot( void *window, void **data, size_t *dataLen,
                        bool *haveAlpha, int *width, int *height,
                        unsigned *stride )
{
    assert( window );
    assert( width );
    assert( height );

    GdkScreen *screen = gtk_window_get_screen( GTK_WINDOW( window ) );
    assert( screen );

    GdkWindow *rw = gdk_screen_get_root_window( screen );
    assert( rw );

    GdkWindow *ww = gtk_widget_get_window( GTK_WIDGET( window ) );
    assert( ww );
    gint x, y, w, h;
    gdk_window_get_origin( ww, &x, &y );
    gdk_window_get_geometry( ww, NULL, NULL, &w, &h );

    // The offsets to the dimensions below are to capture the title bar and
    // the borders for the window.  This is tuned to XFCE, and will likely need
    // to be adjusted with any other DE.
    GdkPixbuf *pix =
        gdk_pixbuf_get_from_window( rw, x - 1, y - 25, w + 2, h + 26 );
    assert( pix );

    guint dataLen2;
    *data = gdk_pixbuf_get_pixels_with_length( pix, &dataLen2 );
    *dataLen = dataLen2;
    *haveAlpha = gdk_pixbuf_get_has_alpha( pix );
    *width = gdk_pixbuf_get_width( pix );
    *height = gdk_pixbuf_get_height( pix );
    *stride = gdk_pixbuf_get_rowstride( pix );

    return pix;
}
