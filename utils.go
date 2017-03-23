package frame

// #include <glib.h>
import "C"

func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func goBool(b C.gboolean) bool {
	if b != 0 {
		return true
	}
	return false
}
