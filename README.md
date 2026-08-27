# go-wimlib

Go bindings for [wimlib](https://wimlib.net), the library behind `wimlib-imagex`. It covers what you need to work with Windows WIM images from Go: open and create archives, extract full images or single paths, add and delete files in an image, export images between archives (ESD to WIM included), and set image names and boot flags.

Two extras go beyond the usual imaging calls. `Decompressor` gives you raw LZX/LZMS block decompression, and the `dcs` package uses it to parse CBS delta-compressed stores, the format Windows packages hide their payloads in. `ListChildren` enumerates a directory inside an image without extracting anything.

The cgo bindings sit behind the `wimlib` build tag, so plain `go build` compiles a stub that returns errors instead, and packages that import this one still build on machines without libwim. Build with `-tags wimlib` (and libwim installed) to get the real thing. This library is part of the Windows tooling behind [devcell](https://github.com/DimmKirr/devcell); pin a version, there is no v1 yet.

```sh
go get github.com/devcell-sh/go-wimlib
```
