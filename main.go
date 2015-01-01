package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"sort"
)

type pkgList struct {
	pmap  map[string]*build.Package
	plist []string
}

func newPkgList() *pkgList {
	ret := new(pkgList)
	ret.pmap = make(map[string]*build.Package)
	ret.plist = make([]string, 0, 100)
	return ret
}

func (self *pkgList) have(p string) bool {
	return self.pmap[p] != nil
}

func (self *pkgList) add(p string, pkg *build.Package) bool {
	ret := self.have(p)
	if !ret {
		self.pmap[p] = pkg
		self.plist = append(self.plist, p)
	}
	return ret
}

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
	sort.Strings(pkgs.plist)
	for _, name := range pkgs.plist {
		if name == target {
			continue
		}
		fmt.Printf("%s\n", name)
	}
}

func findDeps(p string) (*pkgList, error) {
	ret := newPkgList()
	err := addDeps(ret, p)
	return ret, err
}

func addDeps(pkgs *pkgList, p string) error {
	if p == "C" {
		return nil
	}

	pkg, err := build.Import(p, "./", 0)
	if err != nil {
		return fmt.Errorf("failed to import %s: %v", p, err)
	}

	if pkgs.have(p) {
		return nil
	}

	pkgs.add(p, pkg)
	for _, imp := range pkg.Imports {
		err = addDeps(pkgs, imp)
		if err != nil {
			return err
		}
	}

	return nil
}
