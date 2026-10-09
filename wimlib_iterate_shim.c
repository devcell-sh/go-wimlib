//go:build cgo

#include <wimlib.h>
#include <stdint.h>

extern int gowimlibDirEntryGo(char *name, uintptr_t handle);

static int gowimlib_dir_cb(const struct wimlib_dir_entry *dentry, void *ctx) {
	return gowimlibDirEntryGo((char *)dentry->filename, (uintptr_t)ctx);
}

int gowimlibIterateChildren(WIMStruct *w, int image, const wimlib_tchar *path, uintptr_t handle) {
	return wimlib_iterate_dir_tree(w, image, path,
		WIMLIB_ITERATE_DIR_TREE_FLAG_CHILDREN,
		gowimlib_dir_cb, (void *)handle);
}
