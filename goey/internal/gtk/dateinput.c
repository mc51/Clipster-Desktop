#include <assert.h>
#include <gtk/gtk.h>
#include "_cgo_export.h"
#include "callback.h"
#include "thunks.h"

static void onchange_cb( GtkCalendar *widget, gpointer user_data )
{
    assert( widget );

    guint year, month, day;
    gtk_calendar_get_date( widget, &year, &month, &day );
    assert( month <= 11 );
    assert( day >= 1 && day <= 31 );
    onChangeTime( widget, year, month + 1, day );
}

static void setSignals( GtkCalendar *widget, bool onclick, bool onfocus,
                        bool onblur )
{
    if ( onclick ) {
        g_signal_connect( widget, "day-selected", G_CALLBACK( onchange_cb ),
                          NULL );
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

void *mountDateInput( void *parent, int year, unsigned month, unsigned day,
                      bool disabled, bool onchange, bool onfocus, bool onblur )
{
    assert( parent );
    assert( month >= 1 && month <= 12 );
    assert( day >= 1 && day <= 31 );

    GtkWidget *widget = gtk_calendar_new();
    assert( widget );
    gtk_widget_add_events( widget, GDK_FOCUS_CHANGE_MASK );
    gtk_calendar_select_month( GTK_CALENDAR( widget ), month - 1, year );
    gtk_calendar_select_day( GTK_CALENDAR( widget ), day );
    gtk_widget_set_sensitive( widget, !disabled );

    g_signal_connect( widget, "destroy", G_CALLBACK( ondestroy_cb ), NULL );
    setSignals( GTK_CALENDAR( widget ), onchange, onfocus, onblur );

    gtk_container_add( GTK_CONTAINER( parent ), widget );
    gtk_widget_show( widget );

    return widget;
}

void dateInputUpdate( void *widget, int year, unsigned month, unsigned day,
                      bool disabled, bool onchange, bool onfocus, bool onblur )
{
    assert( widget );

    gtk_calendar_select_month( GTK_CALENDAR( widget ), month - 1, year );
    gtk_calendar_select_day( GTK_CALENDAR( widget ), day );
    gtk_widget_set_sensitive( widget, !disabled );

    g_signal_handlers_disconnect_by_data( widget, NULL );
    setSignals( GTK_CALENDAR( widget ), onchange, onfocus, onblur );
}

int dateInputYear( void *widget )
{
    assert( widget );

    guint year, month, day;
    gtk_calendar_get_date( GTK_CALENDAR( widget ), &year, &month, &day );
    return year;
}

unsigned dateInputMonth( void *widget )
{
    assert( widget );

    guint year, month, day;
    gtk_calendar_get_date( GTK_CALENDAR( widget ), &year, &month, &day );
    return month + 1;
}

unsigned dateInputDay( void *widget )
{
    assert( widget );

    guint year, month, day;
    gtk_calendar_get_date( GTK_CALENDAR( widget ), &year, &month, &day );
    return day;
}
