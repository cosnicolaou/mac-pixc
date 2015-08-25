package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type absoluteLink struct {
	source, target string
}

var (
	preview bool
	sysroot string
)

func usage() {
	fmt.Println(
		`rewrite-links will rewrite soft links with a sysroot directory tree
that are absolute to be relative to the root of that tree`)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func main() {
	flag.BoolVar(&preview, "preview", true, "set to true to actually rewrite links")
	flag.StringVar(&sysroot, "sysroot", "", "sysroot directory")
	flag.Parse()
	if len(sysroot) == 0 {
		fmt.Fprintf(os.Stderr, "please supply a sysroot arg\n")
		os.Exit(1)
	}
	ch := make(chan *absoluteLink, 100)
	errch := make(chan error, 1)
	go func() {
		findAbsoluteLinks(ch, errch, sysroot)
		close(ch)
	}()
	for {
		select {
		case err := <-errch:
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		case al := <-ch:
			if al == nil {
				return
			}
			rewrite(al.source, al.target, sysroot)
		}
	}
}

func rewrite(pathname, target, sysroot string) {
	newTarget := filepath.Join(sysroot, target)
	if preview {
		fmt.Printf("rewrite: %s -> %q to %q\n", pathname, target, newTarget)
		return
	}
	os.Remove(pathname)
	if err := os.Symlink(newTarget, pathname); err != nil {
		fmt.Fprintf(os.Stderr, "%q -> %q: %s\n", pathname, newTarget, err)
	}
}

func findAbsoluteLinks(ch chan *absoluteLink, errch chan error, root string) {
	di, err := os.Open(root)
	if err != nil {
		errch <- err
		return
	}
	files, err := di.Readdir(-1)
	if err != nil {
		errch <- err
		return
	}
	di.Close()
	for _, file := range files {
		pathname := filepath.Join(root, file.Name())
		if file.IsDir() {
			findAbsoluteLinks(ch, errch, pathname)
			continue
		}
		if (file.Mode() & os.ModeSymlink) != 0 {
			target, err := os.Readlink(pathname)
			if err != nil {
				errch <- err
			}
			if filepath.IsAbs(target) && !filepath.HasPrefix(target, sysroot) {
				ch <- &absoluteLink{pathname, target}
			}
		}
	}
}
