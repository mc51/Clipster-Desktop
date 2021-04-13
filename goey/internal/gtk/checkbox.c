#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void ontoggled_cb( GtkToggleButton *button, gpointer user_data )
{
    onChangeBool( button, gtk_toggle_button_get_active( button ) );
}

static void setSignals( void *widget, bool onchange, bool onfocus, bool onblur )
{
    assert( widget );

    if ( onchange ) {
        g_signal_connect( widget, "toggled", G_CALLBACK( ontoggled_cb ), NULL );
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

void *mountCheckbox( void *parent, bool value, char const *text, bool disabled,
                     bool onchange, bool onfocus, bool onblur )
{
    assert( parent );
    assert( text );

    GtkWidget *w = gtk_check_button_new_with_label( text );
    assert( w );
    gtk_widget_add_events( w, GDK_FOCUS_CHANGE_MASK );
    gtk_toggle_button_set_active( GTK_TOGGLE_BUTTON( w ), value );
    gtk_widget_set_sensitive( w, !disabled );

    g_signal_connect( w, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( w, onchange, onfocus, onblur );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}
void checkboxUpdate( void *button, bool value, char const *text, bool disabled,
                     bool onchange, bool onfocus, bool onblur )
{
    assert( button );
    assert( text );

    gtk_button_set_label( GTK_BUTTON( button ), text );
    gtk_toggle_button_set_active( GTK_TOGGLE_BUTTON( button ), value );
    gtk_widget_set_sensitive( GTK_WIDGET( button ), !disabled );

    g_signal_handlers_disconnect_by_data( button, NULL );
    setSignals( button, onchange, onfocus, onblur );
}

void checkboxClick( void *button )
{
    assert( button );

    gtk_button_clicked( GTK_BUTTON( button ) );
}

bool checkboxValue( void *button )
{
    assert( button );
    return gtk_toggle_button_get_active( GTK_TOGGLE_BUTTON( button ) );
}

char const *checkboxText( void *button )
{
    assert( button );
    return gtk_button_get_label( GTK_BUTTON( button ) );
}
