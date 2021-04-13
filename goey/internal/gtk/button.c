#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onclicked_cb( GtkButton *button, gpointer user_data )
{
    assert( button );
    onClick( button );
}

static void setSignals( GtkButton *widget, bool onclick, bool onfocus,
                        bool onblur )
{
    if ( onclick ) {
        g_signal_connect( widget, "clicked", G_CALLBACK( onclicked_cb ), NULL );
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

void *mountButton( void *parent, char const *text, bool disabled, bool def,
                   bool onclick, bool onfocus, bool onblur )
{
    assert( parent );
    assert( text );

    GtkWidget *widget = gtk_button_new_with_label( text );
    assert( widget );
    gtk_widget_add_events( widget, GDK_FOCUS_CHANGE_MASK );
    gtk_widget_set_sensitive( widget, !disabled );
    gtk_widget_set_can_default( widget, def );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( GTK_BUTTON( widget ), onclick, onfocus, onblur );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void buttonUpdate( void *widget, char const *text, bool disabled, bool def,
                   bool onclick, bool onfocus, bool onblur )
{
    assert( widget );
    assert( text );

    gtk_button_set_label( GTK_BUTTON( widget ), text );
    gtk_widget_set_sensitive( GTK_WIDGET( widget ), !disabled );
    gtk_widget_set_can_default( GTK_WIDGET( widget ), def );

    g_signal_handlers_disconnect_by_data( widget, NULL );
    setSignals( GTK_BUTTON( widget ), onclick, onfocus, onblur );
}

void buttonClick( void *button )
{
    assert( button );

    gtk_button_clicked( GTK_BUTTON( button ) );
}

char const *buttonText( void *button )
{
    assert( button );
    return gtk_button_get_label( GTK_BUTTON( button ) );
}