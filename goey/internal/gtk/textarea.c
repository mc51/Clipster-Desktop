#include <assert.h>  // for assert
#include <gtk/gtk.h>
#include <stdlib.h>  // for free
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchanged_cb( GtkTextBuffer *widget, gpointer user_data )
{
    assert( widget );
    GtkTextIter si, ei;
    gtk_text_buffer_get_start_iter( widget, &si );
    gtk_text_buffer_get_end_iter( widget, &ei );
    char *text = gtk_text_buffer_get_text( widget, &si, &ei, TRUE );

    onChangeString( user_data, text );

    free( text );
}

static gboolean onfocus_textarea_cb( GtkWidget *widget, GdkEvent *event,
                                     gpointer user_data )
{
    assert( user_data );
    onFocus( user_data );
    return FALSE;
}

static gboolean onblur_textarea_cb( GtkWidget *widget, GdkEvent *event,
                                    gpointer user_data )
{
    assert( user_data );
    onBlur( user_data );
    return FALSE;
}

static void setSignals( GtkScrolledWindow *widget, GtkWidget *textview,
                        bool onchange, bool onfocus, bool onblur )
{
    // Get the text view
    assert( widget );
    assert( textview );

    if ( onchange ) {
        GtkTextBuffer *buffer =
            gtk_text_view_get_buffer( GTK_TEXT_VIEW( textview ) );
        g_signal_connect( buffer, "changed", G_CALLBACK( onchanged_cb ),
                          widget );
    }
    if ( onfocus ) {
        g_signal_connect( textview, "focus-in-event",
                          G_CALLBACK( onfocus_textarea_cb ), widget );
    }
    if ( onblur ) {
        g_signal_connect( textview, "focus-out-event",
                          G_CALLBACK( onblur_textarea_cb ), widget );
    }
}

extern void *mountTextarea( void *parent, char const *text, bool disabled,
                            bool readonly, bool onchange, bool onfocus,
                            bool onblur )
{
    GtkTextBuffer *buffer = gtk_text_buffer_new( NULL );
    assert( buffer );
    gtk_text_buffer_set_text( buffer, text, -1 );

    GtkWidget *w = gtk_text_view_new_with_buffer( buffer );
    assert( w );
    g_object_unref( buffer );

    gtk_text_view_set_left_margin( GTK_TEXT_VIEW( w ), 3 );
    gtk_text_view_set_right_margin( GTK_TEXT_VIEW( w ), 3 );
    #if GTK_CHECK_VERSION( 3, 18, 0 )
    gtk_text_view_set_top_margin( GTK_TEXT_VIEW( w ), 3 );
    gtk_text_view_set_bottom_margin( GTK_TEXT_VIEW( w ), 3 );
    #endif
    gtk_text_view_set_wrap_mode( GTK_TEXT_VIEW( w ), GTK_WRAP_WORD );
    gtk_widget_set_sensitive( w, !disabled );
    gtk_text_view_set_editable( GTK_TEXT_VIEW( w ), !readonly );

    GtkWidget *sw = gtk_scrolled_window_new( NULL, NULL );
    assert( sw );
    gtk_container_add( GTK_CONTAINER( sw ), w );
    gtk_scrolled_window_set_policy( GTK_SCROLLED_WINDOW( sw ), GTK_POLICY_NEVER,
                                    GTK_POLICY_AUTOMATIC );
    gtk_scrolled_window_set_shadow_type( GTK_SCROLLED_WINDOW( sw ),
                                         GTK_SHADOW_IN );
    gtk_widget_set_vexpand( w, TRUE );

    setSignals( GTK_SCROLLED_WINDOW( sw ), w, onchange, onfocus, onblur );

    gtk_container_add( GTK_CONTAINER( parent ), sw );
    gtk_widget_show( sw );

    return sw;
}

extern void textareaUpdate( void *sw, char const *text, bool disabled,
                            bool readonly, bool onchange, bool onfocus,
                            bool onblur )
{
    assert( sw );
    assert( text );

    GtkWidget *widget = gtk_bin_get_child( GTK_BIN( sw ) );
    assert( widget );

    GtkTextBuffer *buffer = gtk_text_view_get_buffer( GTK_TEXT_VIEW( widget ) );
    assert( buffer );
    gtk_text_buffer_set_text( buffer, text, -1 );
    gtk_widget_set_sensitive( GTK_WIDGET( widget ), !disabled );
    gtk_text_view_set_editable( GTK_TEXT_VIEW( widget ), !readonly );

    g_signal_handlers_disconnect_by_data( widget, widget );
    setSignals( GTK_SCROLLED_WINDOW( sw ), widget, onchange, onfocus, onblur );
}

char const *textareaText( void *widget )
{
    assert( widget );
    GtkWidget *textview = gtk_bin_get_child( GTK_BIN( widget ) );
    assert( textview );

    GtkTextBuffer *buffer =
        gtk_text_view_get_buffer( GTK_TEXT_VIEW( textview ) );
    assert( buffer );
    GtkTextIter si, ei;
    gtk_text_buffer_get_start_iter( buffer, &si );
    gtk_text_buffer_get_end_iter( buffer, &ei );
    return gtk_text_buffer_get_text( buffer, &si, &ei, TRUE );
}

char const *textareaPlaceholder( void *widget )
{
    return "";
}

bool textareaReadOnly( void *widget )
{
    assert( widget );
    GtkWidget *textview = gtk_bin_get_child( GTK_BIN( widget ) );
    assert( textview );
    return !gtk_text_view_get_editable( GTK_TEXT_VIEW( textview ) );
}

void *textareaTextview( void *widget )
{
    assert( widget );
    GtkWidget *tv = gtk_bin_get_child( GTK_BIN( widget ) );
    assert( tv );
    return tv;
}
