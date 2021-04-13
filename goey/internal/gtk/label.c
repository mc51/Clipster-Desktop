#include <assert.h>
#include <gtk/gtk.h>
#include "thunks.h"

void *mountLabel( void *parent, char const *text )
{
    assert( parent );
    assert( text );

    GtkWidget *w = gtk_label_new( text );
    assert( w );

    gtk_label_set_single_line_mode( GTK_LABEL( w ), FALSE );
    gtk_label_set_justify( GTK_LABEL( w ), GTK_JUSTIFY_LEFT );
    gtk_widget_set_halign( w, GTK_ALIGN_START );
    gtk_label_set_line_wrap( GTK_LABEL( w ), FALSE );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}

void labelUpdate( void *label, char const *text )
{
    assert( label );
    assert( text );

    gtk_label_set_text( GTK_LABEL( label ), text );
}

char const *labelText( void *label )
{
    assert( label );

    return gtk_label_get_text( GTK_LABEL( label ) );
}
