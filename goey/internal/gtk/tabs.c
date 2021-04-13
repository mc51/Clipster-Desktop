#include <assert.h>  // for assert
#include <gtk/gtk.h>
#include <string.h>  // for strlen
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchange_cb( GtkNotebook *notebook, GtkWidget *page, guint page_num,
                         gpointer user_data )
{
    onChangeTab( notebook, page_num );
}

static void append_item( GtkNotebook *widget, char const *text )
{
    // Every tab needs some contents.  We will use a layout so that we can
    // custom layout of the controls.
    GtkWidget *contents = gtk_layout_new( NULL, NULL );
    assert( contents );

    // Create a label for the tabs.
    GtkWidget *label = gtk_label_new( text );
    assert( label );

    // Append the new page to the notebook.
    gtk_notebook_append_page( widget, contents, label );
    gtk_widget_show( contents );
}

static void append_items( GtkNotebook *widget, char const *text )
{
    assert( text );

    char const *i;
    for ( i = text; *i; i += strlen( i ) + 1 ) {
        append_item( widget, i );
    }
}

static void remove_items( GtkNotebook *widget, int currentPage, int length )
{
    assert( currentPage >= 0 );
    assert( length >= 0 );

    int i;
    for ( i = currentPage; i < length; ++i ) {
        gtk_notebook_remove_page( GTK_NOTEBOOK( widget ), currentPage );
    }
}

void *mountTabs( void *parent, int value, char *tabs, bool onchange )
{
    assert( parent );
    assert( tabs );

    GtkWidget *widget = gtk_notebook_new();
    assert( widget );
    gtk_widget_add_events( widget, GDK_FOCUS_CHANGE_MASK );
    append_items( GTK_NOTEBOOK( widget ), tabs );
    gtk_notebook_set_current_page( GTK_NOTEBOOK( widget ), value );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    g_signal_connect( widget, "switch-page", G_CALLBACK( onchange_cb ), NULL );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void tabsUpdate( void *widget, int value, char *tabs, bool onchange )
{
    assert( widget );
    assert( tabs );

    gint const len = gtk_notebook_get_n_pages( GTK_NOTEBOOK( widget ) );
    int currentPage = 0;
    char const *i;
    for ( i = tabs; *i; i += strlen( i ) + 1 ) {
        if ( currentPage < len ) {
            GtkWidget *page = gtk_notebook_get_nth_page( GTK_NOTEBOOK( widget ),
                                                         currentPage );
            assert( page );
            GtkWidget *label =
                gtk_notebook_get_tab_label( GTK_NOTEBOOK( widget ), page );
            assert( label );
            gtk_label_set_text( GTK_LABEL( label ), i );
        } else {
            append_item( GTK_NOTEBOOK( widget ), i );
        }
        ++currentPage;
    }
    remove_items( GTK_NOTEBOOK( widget ), currentPage, len );

    gtk_notebook_set_current_page( GTK_NOTEBOOK( widget ), value );
}

void *tabsGetTabParent( void *widget, int value )
{
    assert( widget );
    assert( value >= 0 );

    GtkWidget *parent =
        gtk_notebook_get_nth_page( GTK_NOTEBOOK( widget ), value );
    assert( parent );
    return GTK_LAYOUT( parent );
}

int tabsItemCount( void *widget )
{
    assert( widget );

    return gtk_notebook_get_n_pages( GTK_NOTEBOOK( widget ) );
}

char const *tabsItemCaption( void *widget, int index )
{
    assert( widget );
    assert( index >= 0 );

    GtkWidget *page =
        gtk_notebook_get_nth_page( GTK_NOTEBOOK( widget ), index );
    assert( page );
    GtkWidget *label =
        gtk_notebook_get_tab_label( GTK_NOTEBOOK( widget ), page );
    assert( label );
    return gtk_label_get_text( GTK_LABEL( label ) );
}
