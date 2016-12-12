package gotree

// File represents both directories and regular files on disk. `Path` is the
// location of the file from the root directory. If it ends with a `/` (slash),
// the file is considered as a directory. Otherwise, it is considered as a
// regular file.
//
// If it is a regular file, an empty string or nil as `Content` will generate
// an empty file, whereas a non-empty string will be used as the file content.
//
// If it is a directory, a nil `Content` will only generate an empty directory.
// However, if the content is a `gotree.File` or a `[]gotree.File`, the files
// inside of it will be generated.
type File struct {
	Path    string
	Content interface{}
}
