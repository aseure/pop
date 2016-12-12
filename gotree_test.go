package gotree

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCreateEmpty(t *testing.T) {
	root, err := Generate([]File{})
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
	err := GenerateFromRoot(root, []File{})
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
	root, err := ioutil.TempDir("", "gotree_test_")
	defer os.RemoveAll(root)
	if err != nil {
		t.Fatalf("TestCreateEmptyFromExistingRoot: cannot create root directory: %s", err)
	}

	if err = GenerateFromRoot(root, []File{}); err != nil {
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
	files := []File{
		{"README.md", "# This is the title"},
		{"json/", []File{
			{"test1.json", `{"key1":"value1","key2":"value2"}`},
			{"test2.json", `{"key3":"value3","key4":"value4"}`},
		}},
		{"vendor/", nil},
		{"src/", []File{
			{"one.cc", "int main() {}"},
			{"two.cc", "#include <iostream>"},
			{"empty.txt", nil},
		}},
		{"test/", File{".gitkeep", nil}},
	}

	root, err := Generate(files)
	defer os.RemoveAll(root)
	if err != nil {
		t.Fatalf("TestCreateComplex: cannot generate tree: %s", err)
	}

	checkTree(t, root, files)
}

func checkTree(t *testing.T, root string, files []File) {
	for _, file := range files {
		if strings.HasSuffix(file.Path, "/") {
			checkDir(t, root, file)
		} else {
			checkFile(t, root, file)
		}
	}
}

func checkDir(t *testing.T, root string, file File) {
	dirPath := path.Join(root, file.Path)

	if !doesDirExist(dirPath) {
		t.Fatalf("checkDir: %s directory should have been generated", dirPath)
	}

	if file.Content == nil {
		if countFiles(t, dirPath) > 0 || countDirectories(t, dirPath) > 0 {
			t.Fatalf("checkDir: %s directory should be empty", dirPath)
		}
	} else {
		switch content := file.Content.(type) {
		case File:
			checkTree(t, dirPath, []File{content})

		case []File:
			checkTree(t, dirPath, content)

		default:
			t.Fatalf("checkDir: %s directory content should be of type `gotree.File` or `[]gotree.File` but got %T instead", dirPath, content)
		}
	}
}

func checkFile(t *testing.T, root string, file File) {
	filePath := path.Join(root, file.Path)

	if !doesFileExist(filePath) {
		t.Fatalf("checkFile: %s file should have been generated", filePath)
	}

	if file.Content == nil {
		return
	}

	content, ok := file.Content.(string)
	if !ok {
		t.Fatalf("checkFile: %s file content should be of type `string`")
	}

	if content == "" {
		return
	}

	checkFileContent(t, filePath, content)
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
