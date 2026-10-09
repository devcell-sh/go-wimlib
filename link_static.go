//go:build cgo && wimlib_static

package wimlib

// Static linking: the resulting binary has no runtime dependency on a system
// libwim. The per-platform libwim.a is NOT committed — build it once locally
// with `task static-lib` (see Taskfile.yml), which compiles a pinned wimlib
// release with --enable-static and drops lib/<goos>_<goarch>/libwim.a.
//
// libwim is LGPLv3: binaries that statically link it must ship the license
// text (lib/THIRD_PARTY_LICENSES/) and allow relinking — trivially satisfied
// while the consuming application's source is published.

/*
#cgo CFLAGS: -I${SRCDIR}/lib/include
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/darwin_arm64/libwim.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/darwin_amd64/libwim.a
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/linux_arm64/libwim.a -lpthread
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/linux_amd64/libwim.a -lpthread
*/
import "C"
