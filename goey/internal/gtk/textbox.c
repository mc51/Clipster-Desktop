#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchanged_cb( GtkEditable *widget, gpointer user_data )
{
    char const *text = gtk_entry_get_text( GTK_ENTRY( widget ) );
    // Go does not understand constness.
    onChangeString( widget, (char *)text );
}

static void onactivate_cb( GtkEntry *entry, gpointer user_data )
{
    char const *text = gtk_entry_get_text( entry );
    // Go does not understand constness.
    onEnterKey( entry, (char *)text );
}

static void setSignals( void *widget, bool onchange, bool onfocus, bool onblur,
                        bool onenterkey )
{
    if ( onchange ) {
        g_signal_connect( widget, "changed", G_CALLBACK( onchanged_cb ), NULL );
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

void *mountTextbox( void *parent, char const *text, char const *placeholder,
                    bool disabled, bool password, bool readonly, bool onchange,
                    bool onfocus, bool onblur, bool onenterkey )
{
    assert( parent );
    assert( text );

    GtkWidget *w = gtk_entry_new();
    assert( w );
    gtk_entry_set_text( GTK_ENTRY( w ), text );
    gtk_entry_set_placeholder_text( GTK_ENTRY( w ), placeholder );
    gtk_widget_add_events( w, GDK_FOCUS_CHANGE_MASK );
    gtk_widget_set_sensitive( w, !disabled );
    gtk_entry_set_visibility( GTK_ENTRY( w ), !password );
    if ( password ) {
        gtk_entry_set_input_purpose( GTK_ENTRY( w ),
                                     GTK_INPUT_PURPOSE_PASSWORD );
    }
    gtk_editable_set_editable( GTK_EDITABLE( w ), !readonly );

    setSignals( w, onchange, onfocus, onblur, onenterkey );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}

void textboxUpdate( void *widget, char const *text, char const *placeholder,
                    bool disabled, bool password, bool readonly, bool onchange,
                    bool onfocus, bool onblur, bool onenterkey )
{
    assert( widget );
    assert( text );

    gtk_entry_set_text( GTK_ENTRY( widget ), text );
    gtk_entry_set_placeholder_text( GTK_ENTRY( widget ), placeholder );
    gtk_widget_set_sensitive( GTK_WIDGET( widget ), !disabled );
    gtk_entry_set_visibility( GTK_ENTRY( widget ), !password );
    gtk_entry_set_input_purpose(
        GTK_ENTRY( widget ),
        password ? GTK_INPUT_PURPOSE_PASSWORD : GTK_INPUT_PURPOSE_FREE_FORM );
    gtk_editable_set_editable( GTK_EDITABLE( widget ), !readonly );

    g_signal_handlers_disconnect_by_data( widget, NULL );
    setSignals( widget, onchange, onfocus, onblur, onenterkey );
}

char const *textboxText( void *widget )
{
    assert( widget );
    return gtk_entry_get_text( GTK_ENTRY( widget ) );
}

char const *textboxPlaceholder( void *widget )
{
    assert( widget );
    return gtk_entry_get_placeholder_text( GTK_ENTRY( widget ) );
}

bool textboxPassword( void *widget )
{
    assert( widget );
    return !gtk_entry_get_visibility( GTK_ENTRY( widget ) );
}

bool textboxReadOnly( void *widget )
{
    assert( widget );
    return !gtk_editable_get_editable( GTK_EDITABLE( widget ) );
}
