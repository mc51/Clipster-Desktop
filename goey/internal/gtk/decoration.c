#include <assert.h>  // for assert
#include <gtk/gtk.h>
#include <math.h>
#include <stdint.h>  // for uint32_t
#include <stdlib.h>  // for malloc, free
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void ondestroy_decoration_cb( GtkWidget *widget, gpointer user_data )
{
    free( user_data );
    ondestroy_cb( widget, NULL );
}

struct decoration_info_tag {
    uint32_t fill;
    uint32_t stroke;
    int radius;
};
typedef struct decoration_info_tag decoration_info_t;

gboolean ondraw_cb( GtkWidget *widget, cairo_t *cr, gpointer data )
{
    assert( widget );
    assert( cr );

    decoration_info_t const *const info =
        g_object_get_data( G_OBJECT( widget ), "drawinfo" );

    GtkAllocation a;
    gtk_widget_get_allocation( widget, &a );

    if ( info->radius > 0 ) {
        int const width = a.width;
        int const height = a.height;
        int radius = info->radius;

        if ( 2 * radius > width ) {
            radius = width / 2;
        }
        if ( 2 * radius > height ) {
            radius = height / 2;
        }
        cairo_move_to( cr, 0, radius );
        cairo_arc( cr, radius, radius, radius, M_PI, 3 * M_PI_2 );
        cairo_line_to( cr, width - radius, 0 );
        cairo_arc( cr, width - radius, radius, radius, 3 * M_PI_2, 2 * M_PI );
        cairo_line_to( cr, width, height - radius );
        cairo_arc( cr, width - radius, height - radius, radius, 0, M_PI_2 );
        cairo_line_to( cr, radius, height );
        cairo_arc( cr, radius, height - radius, radius, M_PI_2, M_PI );
        cairo_close_path( cr );
    } else {
        cairo_rectangle( cr, 0, 0, a.width, a.height );
    }
    int const fr = info->fill & 0xff;
    int const fg = ( info->fill >> 8 ) & 0xff;
    int const fb = ( info->fill >> 16 ) & 0xff;
    int const fa = ( info->fill >> 24 ) & 0xff;
    int const sr = info->stroke & 0xff;
    int const sg = ( info->stroke >> 8 ) & 0xff;
    int const sb = ( info->stroke >> 16 ) & 0xff;
    int const sa = ( info->stroke >> 24 ) & 0xff;

    if ( fa > 0 && sa > 0 ) {
        cairo_set_source_rgb( cr, ( (double)fr ) / 0xff, ( (double)fg ) / 0xff,
                              ( (double)fb ) / 0xff );
        cairo_fill_preserve( cr );
        cairo_set_source_rgb( cr, ( (double)sr ) / 0xff, ( (double)sg ) / 0xff,
                              ( (double)sb ) / 0xff );
        cairo_stroke( cr );
    } else if ( fa > 0 ) {
        cairo_set_source_rgb( cr, ( (double)fr ) / 0xff, ( (double)fg ) / 0xff,
                              ( (double)fb ) / 0xff );
        cairo_fill( cr );
    } else if ( sa > 0 ) {
        cairo_set_source_rgb( cr, ( (double)sr ) / 0xff, ( (double)sg ) / 0xff,
                              ( (double)sb ) / 0xff );
        cairo_stroke( cr );
    }
    return false;
}

void *mountDecoration( void *parent, unsigned fill, unsigned stroke,
                       int radius )
{
    assert( parent );
    assert( radius >= 0 );

    GtkWidget *widget = gtk_drawing_area_new();
    assert( widget );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_decoration_cb ),
                      NULL );
    g_signal_connect( widget, "draw", G_CALLBACK( ondraw_cb ), NULL );

    decoration_info_t *info = malloc( sizeof( decoration_info_t ) );
    assert( info );
    info->fill = fill;
    info->stroke = stroke;
    info->radius = radius;
    g_object_set_data( G_OBJECT( widget ), "drawinfo", info );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void decorationUpdate( void *widget, unsigned fill, unsigned stroke,
                       int radius )
{
    assert( widget );
    assert( radius >= 0 );

    decoration_info_t *info = g_object_get_data( widget, "drawinfo" );
    assert( info );
    info->fill = fill;
    info->stroke = stroke;
    info->radius = radius;
}

unsigned decorationFill( void *widget )
{
    assert( widget );
    decoration_info_t const *info = g_object_get_data( widget, "drawinfo" );
    assert( info );
    return info->fill;
}

unsigned decorationStroke( void *widget )
{
    assert( widget );
    decoration_info_t const *info = g_object_get_data( widget, "drawinfo" );
    assert( info );
    return info->stroke;
}

int decorationRadius( void *widget )
{
    assert( widget );
    decoration_info_t const *info = g_object_get_data( widget, "drawinfo" );
    assert( info );
    return info->radius;
}

void decorationSetRadius( void *widget, int radius )
{
    assert( widget );
    assert( radius >= 0 );
    decoration_info_t *info = g_object_get_data( widget, "drawinfo" );
    assert( info );
    info->radius = radius;
}
