//go:build cgo

package wimlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailable_CGO(t *testing.T) {
	assert.True(t, Available())
}

func TestCreateWIM_AllCompressions(t *testing.T) {
	for _, tc := range []struct {
		name string
		comp Compression
	}{
		{"None", None},
		{"LZX", LZX},
		{"LZMS", LZMS},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, err := CreateWIM(tc.comp)
			require.NoError(t, err)
			defer w.Close()

			count, err := w.ImageCount()
			require.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	}
}

func TestCreateWIM_InvalidCompression(t *testing.T) {
	_, err := CreateWIM(Compression(99))
	assert.Error(t, err)
}

func TestAddEmptyImage(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	idx, err := w.AddEmptyImage("test-image")
	require.NoError(t, err)
	assert.Equal(t, 1, idx)

	count, err := w.ImageCount()
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAddEmptyImage_Multiple(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	idx1, err := w.AddEmptyImage("first")
	require.NoError(t, err)
	assert.Equal(t, 1, idx1)

	idx2, err := w.AddEmptyImage("second")
	require.NoError(t, err)
	assert.Equal(t, 2, idx2)

	count, err := w.ImageCount()
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestAddEmptyImage_EmptyName(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	idx, err := w.AddEmptyImage("")
	require.NoError(t, err)
	assert.Equal(t, 1, idx)
}

func TestWriteAndReopen(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("img1")
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	count, err := w2.ImageCount()
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUpdateImageAdd_AndExtract(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcFile := filepath.Join(tmp, "hello.txt")
	extractDir := filepath.Join(tmp, "extracted")

	require.NoError(t, os.WriteFile(srcFile, []byte("hello wimlib"), 0644))

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("test")
	require.NoError(t, err)

	err = w.UpdateImageAdd(1, srcFile, `\hello.txt`)
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	require.NoError(t, os.MkdirAll(extractDir, 0755))
	err = w2.ExtractImage(1, extractDir, nil)
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(extractDir, "hello.txt"))
	require.NoError(t, err)
	assert.Equal(t, "hello wimlib", string(got))
}

func TestUpdateImageAddTree(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcDir := filepath.Join(tmp, "src")
	extractDir := filepath.Join(tmp, "extracted")

	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "sub"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("aaa"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "sub", "b.txt"), []byte("bbb"), 0644))

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("tree-test")
	require.NoError(t, err)

	err = w.UpdateImageAddTree(1, srcDir, `\data`)
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	require.NoError(t, os.MkdirAll(extractDir, 0755))
	err = w2.ExtractImage(1, extractDir, nil)
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(extractDir, "data", "a.txt"))
	require.NoError(t, err)
	assert.Equal(t, "aaa", string(got))

	got, err = os.ReadFile(filepath.Join(extractDir, "data", "sub", "b.txt"))
	require.NoError(t, err)
	assert.Equal(t, "bbb", string(got))
}

func TestUpdateImageDelete(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcFile := filepath.Join(tmp, "deleteme.txt")
	extractDir := filepath.Join(tmp, "extracted")

	require.NoError(t, os.WriteFile(srcFile, []byte("gone"), 0644))

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("del-test")
	require.NoError(t, err)

	err = w.UpdateImageAdd(1, srcFile, `\deleteme.txt`)
	require.NoError(t, err)

	err = w.UpdateImageDelete(1, `\deleteme.txt`)
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	require.NoError(t, os.MkdirAll(extractDir, 0755))
	err = w2.ExtractImage(1, extractDir, nil)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(extractDir, "deleteme.txt"))
	assert.True(t, os.IsNotExist(err))
}

func TestExtractPaths(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcDir := filepath.Join(tmp, "src")
	extractDir := filepath.Join(tmp, "extracted")

	require.NoError(t, os.MkdirAll(srcDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "keep.txt"), []byte("keep"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "skip.txt"), []byte("skip"), 0644))

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("paths-test")
	require.NoError(t, err)

	err = w.UpdateImageAddTree(1, srcDir, `\`)
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	require.NoError(t, os.MkdirAll(extractDir, 0755))
	err = w2.ExtractPaths(1, extractDir, []string{`\keep.txt`})
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(extractDir, "keep.txt"))
	require.NoError(t, err)
	assert.Equal(t, "keep", string(got))

	_, err = os.Stat(filepath.Join(extractDir, "skip.txt"))
	assert.True(t, os.IsNotExist(err))
}

func TestSetImageProperty(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	_, err = w.AddEmptyImage("prop-test")
	require.NoError(t, err)

	err = w.SetImageProperty(1, "FLAGS", "9")
	require.NoError(t, err)
}

func TestSetImageName(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")

	w, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = w.AddEmptyImage("original")
	require.NoError(t, err)

	err = w.SetImageName(1, "renamed", "a test description")
	require.NoError(t, err)

	err = w.Write(wimPath)
	require.NoError(t, err)
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	desc, err := w2.ImageDescription(1)
	require.NoError(t, err)
	assert.Equal(t, "a test description", desc)
}

func TestSetBootImage(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	_, err = w.AddEmptyImage("boot")
	require.NoError(t, err)

	err = w.SetBootImage(1)
	require.NoError(t, err)
}

func TestImageDescription_Empty(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	defer w.Close()

	_, err = w.AddEmptyImage("no-desc")
	require.NoError(t, err)

	desc, err := w.ImageDescription(1)
	require.NoError(t, err)
	assert.Empty(t, desc)
}

func TestExportImage(t *testing.T) {
	tmp := t.TempDir()
	srcFile := filepath.Join(tmp, "export.txt")
	wimPath := filepath.Join(tmp, "dest.wim")

	require.NoError(t, os.WriteFile(srcFile, []byte("exported"), 0644))

	src, err := CreateWIM(LZX)
	require.NoError(t, err)

	_, err = src.AddEmptyImage("source")
	require.NoError(t, err)
	require.NoError(t, src.UpdateImageAdd(1, srcFile, `\export.txt`))

	dst, err := CreateWIM(LZX)
	require.NoError(t, err)

	err = src.ExportImage(1, dst, LZX)
	require.NoError(t, err)
	src.Close()

	count, err := dst.ImageCount()
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	err = dst.Write(wimPath)
	require.NoError(t, err)
	dst.Close()

	w, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w.Close()

	extractDir := filepath.Join(tmp, "extracted")
	require.NoError(t, os.MkdirAll(extractDir, 0755))
	require.NoError(t, w.ExtractImage(1, extractDir, nil))

	got, err := os.ReadFile(filepath.Join(extractDir, "export.txt"))
	require.NoError(t, err)
	assert.Equal(t, "exported", string(got))
}

func TestOverwrite(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcFile := filepath.Join(tmp, "added.txt")

	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	_, err = w.AddEmptyImage("overwrite-test")
	require.NoError(t, err)
	require.NoError(t, w.Write(wimPath))
	w.Close()

	require.NoError(t, os.WriteFile(srcFile, []byte("new content"), 0644))

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	require.NoError(t, w2.UpdateImageAdd(1, srcFile, `\added.txt`))
	err = w2.Overwrite()
	require.NoError(t, err)
	w2.Close()

	w3, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w3.Close()

	extractDir := filepath.Join(tmp, "extracted")
	require.NoError(t, os.MkdirAll(extractDir, 0755))
	require.NoError(t, w3.ExtractImage(1, extractDir, nil))

	got, err := os.ReadFile(filepath.Join(extractDir, "added.txt"))
	require.NoError(t, err)
	assert.Equal(t, "new content", string(got))
}

func TestListChildren(t *testing.T) {
	tmp := t.TempDir()
	wimPath := filepath.Join(tmp, "test.wim")
	srcDir := filepath.Join(tmp, "src")

	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "alpha"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "beta"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "alpha", "x.txt"), []byte("x"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "beta", "y.txt"), []byte("y"), 0644))

	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	_, err = w.AddEmptyImage("list-test")
	require.NoError(t, err)
	require.NoError(t, w.UpdateImageAddTree(1, srcDir, `\`))
	require.NoError(t, w.Write(wimPath))
	w.Close()

	w2, err := OpenWIM(wimPath)
	require.NoError(t, err)
	defer w2.Close()

	children, err := w2.ListChildren(1, `\`)
	require.NoError(t, err)
	assert.Contains(t, children, "alpha")
	assert.Contains(t, children, "beta")
}

func TestClose_DoubleClose(t *testing.T) {
	w, err := CreateWIM(LZX)
	require.NoError(t, err)
	w.Close()
	w.Close() // must not panic
}

func TestOpenWIM_NonexistentFile_CGO(t *testing.T) {
	_, err := OpenWIM("/nonexistent/path.wim")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wimlib_open_wim")
}
