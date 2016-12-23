package pop

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCreateEmpty(t *testing.T) {
	root, err := Generate(nil)
	defer os.RemoveAll(root)
	if err != nil {
		t.Fatalf("TestCreateNothing: unexpected error: %s", err)
	}

	if !doesDirExist(root) {
		t.Fatalf("TestCreateNothing: directory has not been created")
	}

	if countFiles(t, root) > 0 || countDirectories(t, root) > 0 {
		t.Fatalf("TestCreateNothing: directory is not empty")
	}
}

func TestCreateEmptyFromNonExistingRoot(t *testing.T) {
	root := "path/does/not/exist/"
	err := GenerateFromRoot(root, nil)
	defer os.RemoveAll("path")
	if err != nil {
		t.Fatalf("TestCreateEmptyFromNonExistingRoot: unexpected error: %s", err)
	}

	if !doesDirExist(root) {
		t.Fatalf("TestCreateEmptyFromNonExistingRoot: directory has not been created")
	}

	if countFiles(t, root) > 0 || countDirectories(t, root) > 0 {
		t.Fatalf("TestCreateEmptyFromNonExistingRoot: directory is not empty")
	}
}

func TestCreateEmptyFromExistingRoot(t *testing.T) {
	root, err := ioutil.TempDir("", "go_test_")
	defer os.RemoveAll(root)
	if err != nil {
		t.Fatalf("TestCreateEmptyFromExistingRoot: cannot create root directory: %s", err)
	}

	if err = GenerateFromRoot(root, nil); err != nil {
		t.Fatalf("TestCreateEmptyFromExistingRoot: unexpected error: %s", err)
	}

	if !doesDirExist(root) {
		t.Fatalf("TestCreateEmptyFromExistingRoot: directory has not been created")
	}

	if countFiles(t, root) > 0 || countDirectories(t, root) > 0 {
		t.Fatalf("TestCreateEmptyFromExistingRoot: directory is not empty")
	}
}

func TestCreateOneDir(t *testing.T) {
	files := Corn{
		"README.md": "# This is the title",
		"json/": Corn{
			"test1.json": bytes.NewBufferString(`{"key1":"value1","key2":"value2"}`),
			"test2.json": `{"key3":"value3","key4":"value4"}`,
		},
		"vendor/": nil,
		"src/": Corn{
			"one.cc":    "int main() {}",
			"two.cc":    "#include <iostream>",
			"empty.txt": nil,
		},
		"test/": Corn{
			".gitkeep": nil,
		},
	}

	root, err := Generate(files)
	defer os.RemoveAll(root)
	if err != nil {
		t.Fatalf("TestCreateComplex: cannot generate tree: %s", err)
	}

	for name, content := range files {
		checkTree(t, root, name, content)
	}
}

func checkTree(t *testing.T, root string, name string, content interface{}) {
	if strings.HasSuffix(name, "/") {
		checkDir(t, root, name, content)
	} else {
		checkFile(t, root, name, content)
	}
}

func checkDir(t *testing.T, root string, name string, content interface{}) {
	dirPath := path.Join(root, name)

	if !doesDirExist(dirPath) {
		t.Fatalf("checkDir: %s directory should have been generated", dirPath)
	}

	if content == nil {
		if countFiles(t, dirPath) > 0 || countDirectories(t, dirPath) > 0 {
			t.Fatalf("checkDir: %s directory should be empty", dirPath)
		}
		return
	}

	switch content := content.(type) {
	case Corn:
		for subName, subContent := range content {
			checkTree(t, dirPath, subName, subContent)
		}

	default:
		t.Fatalf("checkDir: %s directory content should be of type `pop.Corn` but got %T instead", dirPath, content)
	}
}

func contentToString(content interface{}) (string, error) {
	r, err := contentToReader(content)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func checkFile(t *testing.T, root string, name string, content interface{}) {
	filePath := path.Join(root, name)

	if !doesFileExist(filePath) {
		t.Fatalf("checkFile: %s file should have been generated", filePath)
	}

	if content == nil {
		return
	}

	text, err := contentToString(content)
	if err != nil {
		t.Fatalf("checkFile: %s file %s", filePath, err)
	}

	if text == "" {
		return
	}

	checkFileContent(t, filePath, text)
}

func checkFileContent(t *testing.T, path string, expected string) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("checkFileContent: cannot open file %s: %s", path, err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("checkFileContent: cannot read file %s: %s", path, err)
	}

	content := string(data)
	if content != expected {
		t.Fatalf("checkFileContent:\n expected: %s\n      got: %s", expected, content)
	}
}

func doesDirExist(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

func doesFileExist(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return file.Mode().IsRegular()
}

func count(t *testing.T, path string, directories bool) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatalf("countFiles: cannot read directory: %s", err)
	}

	dirCount := 0
	fileCount := 0
	for _, file := range files {
		if file.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
	}

	if directories {
		return dirCount
	} else {
		return fileCount
	}
}

func countDirectories(t *testing.T, path string) int {
	return count(t, path, true)
}

func countFiles(t *testing.T, path string) int {
	return count(t, path, false)
}
