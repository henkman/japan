package main

import (
	"unsafe"

	ui "github.com/libui-ng/golang-ui"
)

/*
	#cgo pkg-config: gtk+-3.0
	#include <gtk/gtk.h>
*/
import "C"

func SetIconFromFile(win *ui.Window, iconfile string) {
	cif := C.CString(iconfile)
	C.gtk_window_set_icon_from_file((*C.GtkWindow)(unsafe.Pointer(win.Handle())), cif, nil)
	C.free(unsafe.Pointer(cif))
}
