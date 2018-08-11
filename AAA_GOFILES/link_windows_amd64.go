// 13 december 2015

package ui

// #cgo LDFLAGS: ${SRCDIR}/libui_windows_amd64.a
// /* note the order; also note the lack of uuid */
// #cgo LDFLAGS: -luser32 -lkernel32 -lusp10 -lgdi32 -lcomctl32 -luxtheme -lmsimg32 -lcomdlg32 -ld2d1 -ldwrite -lole32 -loleaut32 -loleacc -static -static-libgcc -static-libstdc++
import "C"
