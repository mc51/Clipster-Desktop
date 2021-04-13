#include <assert.h>
#include <gtk/gtk.h>
#include "thunks.h"

unsigned dialogRun( void* dialog )
{
    return gtk_dialog_run( GTK_DIALOG( dialog ) );
}

void dialogAddFilter( void* dialog, char const* name, char const* pattern )
{
    GtkFileFilter* ff = gtk_file_filter_new();
    assert( ff );

    gtk_file_filter_set_name( ff, name );
    gtk_file_filter_add_pattern( ff, pattern );
    gtk_file_chooser_add_filter( dialog, ff );
    g_object_unref( ff );
}

unsigned dialogResponseAccept( void )
{
    return GTK_RESPONSE_ACCEPT;
}

unsigned dialogResponseCancel( void )
{
    return GTK_RESPONSE_CANCEL;
}

void* mountMessageDialog( void* window, char const* title, unsigned icon,
                          char const* text )
{
    GtkWidget* dialog =
        gtk_message_dialog_new( GTK_WINDOW( window ), GTK_DIALOG_MODAL, icon,
                                GTK_BUTTONS_OK, "%s", text );
    assert( dialog );
    gtk_window_set_title( GTK_WINDOW( dialog ), text );
    return dialog;
}

unsigned messageDialogWithError( void )
{
    return GTK_MESSAGE_ERROR;
}

unsigned messageDialogWithWarn( void )
{
    return GTK_MESSAGE_WARNING;
}

unsigned messageDialogWithInfo( void )
{
    return GTK_MESSAGE_INFO;
}

void* mountOpenDialog( void* window, char const* title, char const* filename )
{
    GtkWidget* dialog = gtk_file_chooser_dialog_new(
        title, GTK_WINDOW( window ), GTK_FILE_CHOOSER_ACTION_OPEN, "_Open",
        GTK_RESPONSE_ACCEPT, "_Cancel", GTK_RESPONSE_CANCEL, NULL );
    assert( dialog );
    gtk_file_chooser_set_filename( GTK_FILE_CHOOSER( dialog ), filename );
    return dialog;
}

void* mountSaveDialog( void* window, char const* title, char const* filename )
{
    GtkWidget* dialog = gtk_file_chooser_dialog_new(
        title, GTK_WINDOW( window ), GTK_FILE_CHOOSER_ACTION_SAVE, "_Save",
        GTK_RESPONSE_ACCEPT, "_Cancel", GTK_RESPONSE_CANCEL, NULL );
    assert( dialog );
    gtk_file_chooser_set_filename( GTK_FILE_CHOOSER( dialog ), filename );
    return dialog;
}

char const* dialogGetFilename( void* dialog )
{
    return gtk_file_chooser_get_filename( GTK_FILE_CHOOSER( dialog ) );
}
