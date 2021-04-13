#include <assert.h>
#include <gtk/gtk.h>
#include <stdbool.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

gboolean onfocus_cb( GtkWidget *widget, GdkEvent *event, gpointer user_data )
{
    assert( widget );
    onFocus( widget );
    return FALSE;
}

gboolean onblur_cb( GtkWidget *widget, GdkEvent *event, gpointer user_data )
{
    assert( widget );
    onBlur( widget );
    return FALSE;
}

void widgetClose( void *widget )
{
    assert( widget );
    gtk_widget_destroy( widget );
}

bool widgetCanFocus( void *widget )
{
    assert( widget );
    return gtk_widget_get_can_focus( widget );
}

void widgetGrabFocus( void *widget )
{
    assert( widget );
    gtk_widget_grab_focus( widget );
}

bool widgetIsFocus( void *widget )
{
    assert( widget );
    return gtk_widget_is_focus( widget );
}

static void set_key_info( GdkEventKey *evt, GdkWindow *window, guint r )
{
    assert( evt );

    evt->window = window;
    evt->time = GDK_CURRENT_TIME;
    evt->send_event = 1;

    switch ( r ) {
        case 0x1b:
            evt->keyval = GDK_KEY_Escape;
            evt->hardware_keycode = 9;
            break;
        case '\n':
            evt->keyval = GDK_KEY_Return;
            evt->hardware_keycode = 36;
            break;
        default:
            evt->keyval = r;
            break;
    }
}

void widgetSendKey( void *widget, unsigned key, bool release )
{
    GdkEvent *evt = gdk_event_new( release ? GDK_KEY_RELEASE : GDK_KEY_PRESS );
    set_key_info( (GdkEventKey *)evt, gtk_widget_get_window( widget ), key );
    gtk_widget_event( widget, evt );
}

int widgetMinHeight( void *widget )
{
    int min, natural;
    gtk_widget_get_preferred_height( widget, &min, &natural );
    return min;
}

int widgetMinHeightForWidth( void *widget, int width )
{
    int min, natural;
    gtk_widget_get_preferred_height_for_width( widget, width, &min, &natural );
    return min;
}

int widgetNaturalHeight( void *widget )
{
    int min, natural;
    gtk_widget_get_preferred_height( widget, &min, &natural );
    return natural;
}

int widgetNaturalHeightForWidth( void *widget, int width )
{
    int min, natural;
    gtk_widget_get_preferred_height_for_width( widget, width, &min, &natural );
    return natural;
}

int widgetMinWidth( void *widget )
{
    int min, natural;
    gtk_widget_get_preferred_width( widget, &min, &natural );
    return min;
}

int widgetNaturalWidth( void *widget )
{
    int min, natural;
    gtk_widget_get_preferred_width( widget, &min, &natural );
    return natural;
}

int widgetNaturalWidthForHeight( void *widget, int height )
{
    int min, natural;
    gtk_widget_get_preferred_width_for_height( widget, height, &min, &natural );
    return natural;
}

void widgetNaturalSize( void *widget, int *width, int *height )
{
    int min;
    gtk_widget_get_preferred_width( widget, &min, width );
    gtk_widget_get_preferred_height( widget, &min, height );
}

extern void widgetSetBounds( void *widget, int x, int y, int width, int height )
{
    assert( widget );

    GtkWidget *parent = gtk_widget_get_parent( widget );
    GtkLayout *layout = GTK_LAYOUT( parent );
    gtk_layout_move( layout, widget, x, y );
    GtkAllocation alloc = {x, y, width, height};
    gtk_widget_size_allocate( widget, &alloc );
}

bool widgetSensitive( void *widget )
{
    assert( widget );
    return gtk_widget_get_sensitive( GTK_WIDGET( widget ) );
}

bool widgetCanDefault( void *widget )
{
    assert( widget );
    return gtk_widget_get_can_default( GTK_WIDGET( widget ) );
}

void widgetSetSizeRequest( void *widget, int width, int height )
{
    assert( widget );
    gtk_widget_set_size_request( GTK_WIDGET( widget ), width, height );
}
