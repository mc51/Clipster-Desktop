#include <assert.h>  // for assert
#include <gtk/gtk.h>
#include <string.h>  // for strlen
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchange_cb( GtkComboBox *widget, gpointer data )
{
    assert( widget );
    int value = gtk_combo_box_get_active( widget );
    onChangeInt( widget, value );
}

static gboolean onfocus_combobox_cb( GtkWidget *widget, GdkEvent *event,
                                     gpointer data )
{
    assert( widget );
    assert( data );

    onFocus( data );
    return FALSE;
}

static gboolean onblur_combobox_cb( GtkWidget *widget, GdkEvent *event,
                                    gpointer data )
{
    assert( widget );
    assert( data );

    onBlur( data );
    return FALSE;
}

static void setSignals( void *widget, bool onchange, bool onfocus, bool onblur )
{
    if ( onchange ) {
        g_signal_connect( widget, "changed", G_CALLBACK( onchange_cb ), NULL );
    }
    if ( onfocus ) {
        GtkWidget *child = gtk_bin_get_child( GTK_BIN( widget ) );
        assert( child );
        g_signal_connect( child, "focus-in-event",
                          G_CALLBACK( onfocus_combobox_cb ), widget );
    }
    if ( onblur ) {
        GtkWidget *child = gtk_bin_get_child( GTK_BIN( widget ) );
        assert( child );
        g_signal_connect( child, "focus-out-event",
                          G_CALLBACK( onblur_combobox_cb ), widget );
    }
}

static void append_items( GtkComboBoxText *widget, char const *text )
{
    assert( text );

    char const *i;
    for ( i = text; *i; i += strlen( i ) + 1 ) {
        gtk_combo_box_text_append_text( widget, i );
    }
}

void *mountCombobox( void *parent, char const *items, int value, bool unset,
                     bool disabled, bool onchange, bool onfocus, bool onblur )
{
    assert( parent );
    assert( items );

    GtkWidget *widget = gtk_combo_box_text_new();
    assert( widget );
    gtk_widget_add_events( widget, GDK_FOCUS_CHANGE_MASK );
    append_items( GTK_COMBO_BOX_TEXT( widget ), items );
    if ( !unset ) {
        gtk_combo_box_set_active( GTK_COMBO_BOX( widget ), value );
    }
    gtk_widget_set_can_focus( widget, TRUE );
    gtk_widget_set_sensitive( widget, !disabled );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( widget, onchange, onfocus, onblur );

    GtkWidget *child = gtk_bin_get_child( GTK_BIN( widget ) );
    assert( child );
    gtk_widget_set_can_focus( child, true );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void comboboxUpdate( void *widget, char const *items, int value, bool unset,
                     bool disabled, bool onchange, bool onfocus, bool onblur )
{
    assert( widget );

    gtk_combo_box_text_remove_all( GTK_COMBO_BOX_TEXT( widget ) );
    append_items( GTK_COMBO_BOX_TEXT( widget ), items );
    if ( !unset ) {
        gtk_combo_box_set_active( GTK_COMBO_BOX( widget ), value );
    }
    gtk_widget_set_sensitive( widget, !disabled );
}

void *comboboxChild( void *widget )
{
    assert( widget );
    return gtk_bin_get_child( GTK_BIN( widget ) );
}

unsigned comboboxItemCount( void *widget )
{
    assert( widget );

    GtkTreeModel *model = gtk_combo_box_get_model( GTK_COMBO_BOX( widget ) );
    assert( model );

    GtkTreeIter i;
    unsigned count = 0;
    bool ok;
    for ( ok = gtk_tree_model_get_iter_first( model, &i ); ok;
          ok = gtk_tree_model_iter_next( model, &i ) ) {
        ++count;
    }

    return count;
}

char const *comboboxItem( void *widget, unsigned index )
{
    assert( widget );

    GtkTreeModel *model = gtk_combo_box_get_model( GTK_COMBO_BOX( widget ) );
    assert( model );

    GtkTreeIter i;
    bool ok;
    for ( ok = gtk_tree_model_get_iter_first( model, &i ); ok;
          ok = gtk_tree_model_iter_next( model, &i ) ) {
        if ( index == 0 ) {
            GValue value = G_VALUE_INIT;
            gtk_tree_model_get_value( model, &i, 0, &value );
            return g_value_get_string( &value );
        }
        --index;
    }

    return "";
}

int comboboxValue( void *widget )
{
    assert( widget );

    return gtk_combo_box_get_active( GTK_COMBO_BOX( widget ) );
}
