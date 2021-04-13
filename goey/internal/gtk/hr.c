#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

extern void *mountHR( void *parent )
{
    assert( parent );

    GtkWidget *w = gtk_separator_new( GTK_ORIENTATION_HORIZONTAL );
    assert( w );
    g_signal_connect( w, "destroy", G_CALLBACK( ondestroy_cb ), NULL );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}