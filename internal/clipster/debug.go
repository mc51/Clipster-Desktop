// Debug logging. Anything that may contain private data, like clip content, is only logged here
package clipster

import (
	"log"
	"os"
	"strings"
)

// DEBUG_ENV_VAR enables debug messages when set to a non-empty value other than 0 or false
const DEBUG_ENV_VAR = "CLIPSTER_DEBUG"

// debugEnabled reports whether debug messages are written to the log
var debugEnabled = envFlag(os.Getenv(DEBUG_ENV_VAR))

func envFlag(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v != "" && v != "0" && v != "false"
}

// debugf logs a message prefixed with "Debug:", but only if debugging is enabled
func debugf(format string, args ...any) {
	if debugEnabled {
		log.Printf("Debug: "+format, args...)
	}
}
