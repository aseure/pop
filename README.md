# pop.Corn
[![Build Status](https://travis-ci.org/aseure/pop.svg?branch=master)](https://travis-ci.org/aseure/pop)

Go standard library makes file manipulation very easy. However, it's a little
cumbersome to generate a whole file architecture easily. This `pop` Go package
is here to help! Simply describe your tree with a few *pop.Corn*, call
`pop.Generate` or `pop.GenerateFromRoot` on it and all the intermediate
directories and dummy files will be ready in no time! Populating a directory
for integration testing has never been so fast and fun.

![pop](https://github.com/aseure/pop/raw/master/pop.gif)

## What's included?

Only one structure and two functions are exported.

The `pop.Corn` structure is actually a `map[string]interface{}` holding all
your files and directories, because a 14-character difference is a lot.

The `pop.Generate` generates the given files and directories in a new
directory. Its path is returned as long as an error which may be non-nil if
anything went wrong during the generation.

The `pop.GenerateFromRoot` does the same as `pop.Generate` but will produce the
resulting files and directories within the providided directory path. It also
returns a non-nil error if anything went wrong during the generation. If the
root directory already exists, it will be deleted first.

## Example

In the following example, we are generating a tree of files under the
automatically generated `root` directory. It contains:
 - a `README.md` file
 - a `json/` directory with two JSON files in it
 - an empty `vendor/` directory
 - a `src/` directory with two C++ files in it and an empty file
 - a `test/` directory with an empty file in it

```go
files := pop.Corn{
    "README.md": "# This is the title",
    "json/": pop.Corn{
        "test1.json": bytes.NewBufferString(`{"key1":"value1","key2":"value2"}`),
        "test2.json": `{"key3":"value3","key4":"value4"}`,
    },
    "vendor/": nil,
    "src/": pop.Corn{
        "one.cc":    "int main() {}",
        "two.cc":    "#include <iostream>",
        "empty.txt": nil,
    },
    "test/": pop.Corn{
        ".gitkeep": nil,
    },
}

root, err := pop.Generate(files)
```

All directory names must end with a slash and can either be `nil` or contain a
new `pop.Corn` instance.

Files are represented by non-slash-terminating string names and their content
as a `string` or an  `io.Reader`. If the content is `nil` or an empty string,
the generated file will be empty.
