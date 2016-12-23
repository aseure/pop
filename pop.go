package pop

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Corn is a `map[string]interface{}` used to describe the files and
// directories to generate.
type Corn map[string]interface{}

// Generate populates a new tree of files on disk populated with the given
// `Corn` (map). The root directory path is returned as a string. If an error
// occured during the generation, a non-nil error is returned.
func Generate(files Corn) (root string, err error) {
	root, err = ioutil.TempDir(root, "pop")
	if err != nil {
		err = fmt.Errorf("pop: cannot generate root directory: %s", err)
		return
	}

	err = GenerateFromRoot(root, files)
	return
}

// GenerateFromRoot populates the given `root` directory with the given `Corn`
// (map). If an error occured during the generation, a non-nil error is
// returned.
func GenerateFromRoot(root string, files Corn) (err error) {
	if root == "" {
		return fmt.Errorf("pop: root directory cannot be nil", err)
	}

	if err := os.RemoveAll(root); err != nil {
		return fmt.Errorf("pop: cannot delete pre-existing root directory %s: %s", root, err)
	}

	if err = createDir(root); err != nil {
		return err
	}

	for name, content := range files {
		if err = generate(root, name, content); err != nil {
			return err
		}
	}

	return nil
}

func createDir(path string) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return fmt.Errorf("pop: cannot create directory %s: %s", path, err)
	}

	return nil
}

func generate(root string, name string, content interface{}) error {
	if strings.HasSuffix(name, "/") {
		return generateDir(root, name, content)
	} else {
		return generateFile(root, name, content)
	}
}

func generateDir(root string, name string, content interface{}) error {
	var err error
	dirPath := path.Join(root, name)

	// Create the directory
	if err = createDir(dirPath); err != nil {
		return err
	}

	// Return without error if the content is nil
	if content == nil {
		return nil
	}

	// Generate the directory content
	switch content := content.(type) {
	case Corn:
		for subName, subContent := range content {
			if err = generate(dirPath, subName, subContent); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("pop: directory content is typed %T instead of Corn", content)
	}
}

func generateFile(root string, name string, content interface{}) error {
	filePath := path.Join(root, name)

	// Open the file, which must not already exist
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("pop: cannot create file %s: %s", filePath, err)
	}
	defer f.Close()

	// Generate the content only if it is non-nil or a non-empty string
	r, err := contentToReader(content)
	if err != nil {
		return fmt.Errorf("pop: file content of %s is not valid: %s", filePath, err)
	}
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("pop: cannot write file %s: %s", filePath, err)
	}
	return nil
}

func contentToReader(content interface{}) (io.Reader, error) {
	if content == nil {
		return bytes.NewReader([]byte{}), nil
	}
	switch v := content.(type) {
	case string:
		return bytes.NewBufferString(v), nil
	case io.Reader:
		return v, nil
	}
	return nil, fmt.Errorf("unsuported type")
}
