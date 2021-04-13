#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchange_cb( GtkRange *range, gpointer user_data )
{
    assert( range );
    onChangeFloat64( range, gtk_range_get_value( range ) );
}

static void setSignals( void *widget, bool onchange, bool onfocus, bool onblur )
{
    assert( widget );

    if ( onchange ) {
        g_signal_connect( widget, "value-changed", G_CALLBACK( onchange_cb ),
                          NULL );
    }
    if ( onfocus ) {
        g_signal_connect( widget, "focus-in-event", G_CALLBACK( onfocus_cb ),
                          NULL );
    }
    if ( onblur ) {
        g_signal_connect( widget, "focus-out-event", G_CALLBACK( onblur_cb ),
                          NULL );
    }
}

extern void *mountSlider( void *parent, double value, bool disabled, double min,
                          double max, bool onchange, bool onfocus, bool onblur )
{
    assert( parent );
    assert( min <= max );

    GtkWidget *w = gtk_scale_new_with_range( GTK_ORIENTATION_HORIZONTAL, min,
                                             max, ( max - min ) / 10 );
    assert( w );
    gtk_widget_add_events( w, GDK_FOCUS_CHANGE_MASK );
    gtk_range_set_value( GTK_RANGE( w ), value );
    gtk_scale_set_draw_value( GTK_SCALE( w ), FALSE );
    gtk_widget_set_sensitive( w, !disabled );

    g_signal_connect( w, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( w, onchange, onfocus, onblur );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}

extern void sliderUpdate( void *widget, double value, bool disabled, double min,
                          double max, bool onchange, bool onfocus, bool onblur )
{
    assert( widget );
    gtk_range_set_range( GTK_RANGE( widget ), min, max );
    gtk_range_set_value( GTK_RANGE( widget ), value );
    gtk_widget_set_sensitive( GTK_WIDGET( widget ), !disabled );

    g_signal_handlers_disconnect_by_data( widget, NULL );
    setSignals( widget, onchange, onfocus, onblur );
}

extern double sliderValue( void *widget )
{
    assert( widget );
    return gtk_range_get_value( GTK_RANGE( widget ) );
}
