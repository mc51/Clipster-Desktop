#ifndef GOEY_CALLBACK_H
#define GOEY_CALLBACK_H

#include <gtk/gtk.h>

void ondestroy_cb( GtkWidget *, gpointer );
gboolean onfocus_cb( GtkWidget *, GdkEvent *, gpointer );
gboolean onblur_cb( GtkWidget *, GdkEvent *, gpointer );

#endif