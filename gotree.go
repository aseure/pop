package gotree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Generate populates a new tree of files on disk populated with the given
// `File` slice. The root directory path is returned as a string. If an error
// occured during the generation, a non-nil error is returned.
func Generate(files []File) (root string, err error) {
	root, err = ioutil.TempDir(root, "gotree")
	if err != nil {
		err = fmt.Errorf("gotree: cannot generate root directory: %s", err)
		return
	}

	err = GenerateFromRoot(root, files)
	return
}

// GenerateFromRoot populates the given `root` directory with the given `File`
// slice. If an error occured during the generation, a non-nil error is
// returned.
func GenerateFromRoot(root string, files []File) (err error) {
	if root == "" {
		return fmt.Errorf("gotree: root directory cannot be nil", err)
	}

	if err = createDir(root); err != nil {
		return err
	}

	for _, file := range files {
		if err = generate(root, file); err != nil {
			return err
		}
	}

	return nil
}

func createDir(path string) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return fmt.Errorf("gotree: cannot create directory %s: %s", path, err)
	}

	return nil
}

func generate(root string, file File) error {
	if strings.HasSuffix(file.Path, "/") {
		return generateDir(root, file)
	} else {
		return generateFile(root, file)
	}
}

func generateDir(root string, file File) error {
	var err error
	dirPath := path.Join(root, file.Path)

	// Create the directory
	if err = createDir(dirPath); err != nil {
		return err
	}

	// Return without error if the content is nil
	if file.Content == nil {
		return nil
	}

	// Generate the directory content
	switch content := file.Content.(type) {
	case File:
		return generate(dirPath, content)

	case []File:
		for _, subFile := range content {
			if err = generate(dirPath, subFile); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("gotree: directory content is not a valid file")
	}
}

func generateFile(root string, file File) error {
	filePath := path.Join(root, file.Path)

	// Open the file, which must not already exist
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("gotree: cannot create file %s: %s", filePath, err)
	}
	defer f.Close()

	// Generate the content only if `Content` is non-nil or a non-empty string
	if file.Content != nil {
		content, ok := file.Content.(string)
		if !ok {
			return fmt.Errorf("gotree: file content of %s should be a string", filePath)
		}

		if _, err = f.WriteString(content); err != nil {
			return fmt.Errorf("gotree: cannot write file %s: %s", filePath, err)
		}
	}

	return nil
}
