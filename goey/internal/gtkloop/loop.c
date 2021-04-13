#include "_cgo_export.h"
#include "thunks.h"
#include <gtk/gtk.h>

static gboolean main_context_invoke_cb(gpointer user_data) {
  // Callback into Go.
  mainContextInvokeCallback();
  // Prevent a repeat of this event.
  return FALSE;
}

void loopMainContextInvoke(void) {
  g_main_context_invoke(NULL, main_context_invoke_cb, NULL);
}

void loopIdleAdd(void) { g_idle_add(main_context_invoke_cb, NULL); }
