package main

import (
	"log"
	"fmt"
	"flag"
	"go/build"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need one package name to process")
	}

	target := args[0]
	pkgs, err := findDeps(target)
	if err != nil {
		log.Fatal(err)
	}
	for name, _ := range pkgs {
		if name == target {
			continue
		}
		fmt.Printf("%s\n", name)
	}
}

func findDeps(p string) (map[string]*build.Package, error) {
	ret := make(map[string]*build.Package)
	err := addDeps(ret, p)
	return ret, err
}

func addDeps(pkgs map[string]*build.Package, p string) error {
	pkg, err := build.Import(p, "./", 0)
	if err != nil {
		return fmt.Errorf("failed to import %s: %v", p, err)
	}

	if _, found := pkgs[pkg.ImportPath]; found {
		return nil
	}

	pkgs[pkg.ImportPath] = pkg
	
	for _, imp := range pkg.Imports {
		err = addDeps(pkgs, imp)
		if err != nil {
			return err
		}
	}

	return nil
}
