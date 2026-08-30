//go:build wimlib && !wimlib_static

package wimlib

// Dynamic linking: libwim comes from the system (brew install wimlib,
// nix profile install nixpkgs#wimlib, apt install libwim-dev). pkg-config
// supplies both the include path and -lwim.

/*
#cgo pkg-config: wimlib
*/
import "C"
