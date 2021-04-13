#include <assert.h>  // for assert
#include <gtk/gtk.h>
#include <stdlib.h>  // for atol
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void oncchange_cb( GtkSpinButton *widget, gpointer user_data )
{
    assert( widget );
    onChangeInt64( widget, gtk_spin_button_get_value( widget ) );
}

static void onactivate_cb( GtkEntry *entry, gpointer user_data )
{
    char const *text = gtk_entry_get_text( entry );
    assert( text );
    onEnterKeyInt64( entry, atol( text ) );
}

static void setSignals( GtkSpinButton *widget, bool onchange, bool onfocus,
                        bool onblur, bool onenterkey )
{
    if ( onchange ) {
        g_signal_connect( widget, "value-changed", G_CALLBACK( oncchange_cb ),
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
    if ( onenterkey ) {
        g_signal_connect( widget, "activate", G_CALLBACK( onactivate_cb ),
                          NULL );
    }
}

void *mountIntInput( void *parent, long value, char const *placeholder,
                     bool disabled, long min, long max, bool onchange,
                     bool onfocus, bool onblur, bool onenterkey )
{
    assert( parent );

    GtkWidget *widget = gtk_spin_button_new_with_range( min, max, 1 );
    assert( widget );
    gtk_spin_button_set_value( GTK_SPIN_BUTTON( widget ), value );
    gtk_spin_button_set_increments( GTK_SPIN_BUTTON( widget ), 1, 10 );
    gtk_entry_set_placeholder_text( GTK_ENTRY( widget ), placeholder );
    gtk_widget_set_sensitive( widget, !disabled );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( GTK_SPIN_BUTTON( widget ), onchange, onfocus, onblur,
                onenterkey );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void intinputUpdate( void *widget, long value, char const *placeholder,
                     bool disabled, long min, long max, bool onchange,
                     bool onfocus, bool onblur, bool onenterkey )
{
    gtk_spin_button_set_value( GTK_SPIN_BUTTON( widget ), value );
    gtk_spin_button_set_increments( GTK_SPIN_BUTTON( widget ), 1, 10 );
    gtk_entry_set_placeholder_text( GTK_ENTRY( widget ), placeholder );
    gtk_widget_set_sensitive( GTK_WIDGET( widget ), !disabled );

    g_signal_handlers_disconnect_by_data( widget, NULL );
    setSignals( GTK_SPIN_BUTTON( widget ), onchange, onfocus, onblur,
                onenterkey );
}

long intinputValue( void *widget )
{
    assert( widget );
    double value = gtk_spin_button_get_value( GTK_SPIN_BUTTON( widget ) );
    return value;
}

double intinputMin( void *widget )
{
    assert( widget );
    double min, max;
    gtk_spin_button_get_range( GTK_SPIN_BUTTON( widget ), &min, &max );
    return min;
}

double intinputMax( void *widget )
{
    assert( widget );
    double min, max;
    gtk_spin_button_get_range( GTK_SPIN_BUTTON( widget ), &min, &max );
    return max;
}
