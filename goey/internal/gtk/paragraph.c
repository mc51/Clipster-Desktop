#include <assert.h>
#include <gtk/gtk.h>
#include "callback.h"
#include "thunks.h"

static GtkJustification tojustify( char align )
{
    switch ( align ) {
        case 0:
            return GTK_JUSTIFY_LEFT;
        case 1:
            return GTK_JUSTIFY_CENTER;
        case 2:
            return GTK_JUSTIFY_RIGHT;
        default:
            return GTK_JUSTIFY_FILL;
    }
}

static GtkAlign toalign( char align )
{
    switch ( align ) {
        case 0:
            return GTK_ALIGN_START;
        case 1:
            return GTK_ALIGN_CENTER;
        case 2:
            return GTK_ALIGN_END;
        default:
            return GTK_ALIGN_START;
    }
}

void *mountParagraph( void *parent, char const *text, char align )
{
    GtkWidget *w = gtk_label_new( text );
    assert( w );
    gtk_label_set_single_line_mode( GTK_LABEL( w ), false );
    gtk_label_set_justify( GTK_LABEL( w ), tojustify( align ) );
    gtk_widget_set_halign( w, toalign( align ) );
    gtk_label_set_line_wrap( GTK_LABEL( w ), true );

    g_signal_connect( w, "destroy", G_CALLBACK( ondestroy_cb ), NULL );

    gtk_container_add( GTK_CONTAINER( parent ), w );
    gtk_widget_show( w );

    return w;
}

void paragraphUpdate( void *widget, char const *text, char align )
{
    gtk_label_set_text( GTK_LABEL( widget ), text );
    gtk_label_set_justify( GTK_LABEL( widget ), tojustify( align ) );
    gtk_widget_set_halign( GTK_WIDGET( widget ), toalign( align ) );
}

char const *paragraphText( void *widget )
{
    assert( widget );
    return gtk_label_get_text( GTK_LABEL( widget ) );
}

char paragraphAlign( void *widget )
{
    assert( widget );
    GtkJustification j = gtk_label_get_justify( GTK_LABEL( widget ) );

    /* See paragraph.go in main repository for ordering */
    switch ( j ) {
        case GTK_JUSTIFY_LEFT:
            return 0;
        case GTK_JUSTIFY_CENTER:
            return 1;
        case GTK_JUSTIFY_RIGHT:
            return 2;
        default:
            assert( j == GTK_JUSTIFY_FILL );
            return 3;
    }
}

void paragraphSetText( void *widget, char const *text )
{
    assert( widget );
    assert( text );
    gtk_label_set_text( GTK_LABEL( widget ), text );
}
