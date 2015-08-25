package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var defaultDirsToKeep = []string{
	"/lib",
	"/usr/lib",
	"/usr/include",
}

type dirList []string

func (dirs *dirList) Set(d string) error {
	*dirs = append(*dirs, strings.Split(d, ",")...)
	return nil
}

func (dirs dirList) String() string {
	return strings.Join(dirs, "\n")
}

var sysroot string
var dirsToKeep []string
var preview bool

func usage() {
	fmt.Printf(
		`clean-sysroot is a utility to clean out unneeded directories from a sysroot
directory tree. The --keep argument is used to specify a comma-separated list of
directories to be kept. --keep may be repeated. There is a default list of
directories to be kept: `)
	fmt.Printf("%s.\n", strings.Join(defaultDirsToKeep, ", "))
	fmt.Println(
		`By default clean-sysroot operates in 'preview' mode and shows
the action it would take without actually removing files; use --preview=false
to delete files.`)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func main() {

	var dirs dirList
	flag.StringVar(&sysroot, "sysroot", "", "root of sysroot directory tree")
	flag.BoolVar(&preview, "preview", true, "set to true to actually delete files")
	flag.Var(&dirs, "keep", "directories to keep")
	flag.Parse()
	if len(sysroot) == 0 {
		fmt.Fprintf(os.Stderr, "please supply a sysroot arg.\n")
		os.Exit(1)
	}
	if len(dirs) == 0 {
		dirs.Set(strings.Join(defaultDirsToKeep, ","))
	}

	for _, d := range dirs {
		dirsToKeep = append(dirsToKeep, filepath.Join(sysroot, d))
	}
	clean(sysroot)
}

func keep(pathname string) bool {
	if !filepath.HasPrefix(pathname, sysroot) {
		// never delete anything outside of our sysroot.
		return false
	}
	for _, k := range dirsToKeep {
		if filepath.HasPrefix(pathname, k) {
			return true
		}
		// Don't delete directories that are prefixes of paths
		// to keep.
		if filepath.HasPrefix(k, pathname) {
			return true
		}
	}
	return false
}

func clean(root string) {
	di, err := os.Open(root)
	if err != nil {
		fmt.Printf("failed to open %q: %s\n", root, err)
		return
	}
	files, err := di.Readdir(-1)
	if err != nil {
		fmt.Printf("failed to readdir for %q: %s\n", di.Name(), err)
		return
	}
	di.Close()
	for _, file := range files {
		pathname := filepath.Join(root, file.Name())
		keeping := keep(pathname)
		if preview {
			verb := map[bool]string{true: "Keeping", false: "Deleting"}[keeping]
			if file.IsDir() {
				fmt.Printf("%s: %q/...\n", verb, pathname)
			} else {
				fmt.Printf("%s: %q\n", verb, pathname)
			}
		} else {
			if !keeping {
				os.RemoveAll(pathname)
				continue
			}
		}
		if file.IsDir() {
			clean(pathname)
		}
	}
}
