package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("invalid argument")
		os.Exit(1)
	}

	destPath := args[0]
	funcMap := map[string]interface{}{
		"package_path": func() string {
			fmt.Println(destPath)
			return destPath
		},
		"package_name": func() string {
			fmt.Println(path.Base(destPath))
			return path.Base(destPath)
		},
	}

	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ";") {
		sourcePath := filepath.Join(gopath, "src/github.com/qor/admin/bindatafs/templates")
		err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err == nil {
				var relativePath = strings.TrimPrefix(path, sourcePath)

				if info.IsDir() {
					err = os.MkdirAll(filepath.Join(destPath, relativePath), os.ModePerm)
				} else if info.Mode().IsRegular() {
					if source, err := ioutil.ReadFile(path); err == nil {
						var tmpl *template.Template
						if tmpl, err = template.New("").Funcs(funcMap).Parse(string(source)); err == nil {
							var result = bytes.NewBufferString("")
							if err = tmpl.Execute(result, ""); err != nil {
								return err
							}
							source = result.Bytes()
						} else {
							return err
						}
						if err = ioutil.WriteFile(filepath.Join(destPath, strings.TrimSuffix(relativePath, ".template")), source, os.ModePerm); err != nil {
							fmt.Println(err)
						}
					}
				}
			}
			return err
		})

		if err != nil {
			fmt.Println("failed to copy files:", err)
		}
	}
}
