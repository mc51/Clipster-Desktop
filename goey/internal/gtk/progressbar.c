#include <assert.h>
#include <gtk/gtk.h>
#include "callback.h"
#include "thunks.h"

extern void *mountProgressbar( void *parent, double value )
{
    assert( parent );

    GtkWidget *w = gtk_progress_bar_new();
    assert( w );
    gtk_widget_add_events( w, GDK_FOCUS_CHANGE_MASK );
    gtk_progress_bar_set_fraction( GTK_PROGRESS_BAR( w ), value );
    g_signal_connect( w, "destroy", G_CALLBACK( ondestroy_cb ), NULL );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}

extern void progressbarUpdate( void *widget, double value )
{
    assert( widget );
    gtk_progress_bar_set_fraction( GTK_PROGRESS_BAR( widget ), value );
}

extern double progressbarValue( void *widget )
{
    assert( widget );
    return gtk_progress_bar_get_fraction( GTK_PROGRESS_BAR( widget ) );
}
