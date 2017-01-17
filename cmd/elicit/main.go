package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	pkg   = flag.String("package", "", "The package which will contain the generated tests. The default is [directory]_test.")
	steps = flag.String("steps", ".", "The directory to search for step implementations.")
)

func usage() {
	name := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
	fmt.Fprintf(os.Stderr, "  %s [-package pkg] [-steps stepsdir] [directory]\n\n", name)
	fmt.Fprintf(os.Stderr, "  The default directory is \".\"\n")
	fmt.Fprintf(os.Stderr, "  The default stepsdir is \".\"\n")
	fmt.Fprintf(os.Stderr, "  The default pkg is [directory]_test\n")
	flag.PrintDefaults()
}

func main() {
	dir := parseArgs()
	registerSteps(*steps)
	transformSpecs(dir)
}

func parseArgs() string {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	} else if len(args) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if !isDirectory(*steps) {
		log.Fatalf("%q is not a directory", *steps)
		os.Exit(2)
	}

	dir := args[0]
	if !isDirectory(dir) {
		log.Fatalf("%q is not a directory", dir)
		os.Exit(3)
	}

	if len(*pkg) == 0 {
		if absDir, err := filepath.Abs(dir); err != nil {
			log.Fatalf("determining absolute path for output package: %s", err)
		} else {
			*pkg = filepath.Join(filepath.Dir(absDir), filepath.Base(absDir)+"_test")
		}
	}

	if err := os.MkdirAll(*pkg, 0755); err != nil {
		log.Fatalf("creating output directory %q: %s", *pkg, err)
	}

	return dir
}

func isDirectory(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("checking directory %q: %s", dir, err)
	}
	return info.IsDir()
}

func registerSteps(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {

			stepPkgImport := getPathPackageImport(path)
			stepImpls := parseStepPackage(path)

			g := stepImplGenerator{
				stepPkg:   stepPkgImport,
				stepImpls: stepImpls,
				outPkg:    *pkg,
			}
			g.generate()
		}
		return nil
	})
}

func getPathPackageImport(path string) string {
	if stepsAbs, err := filepath.Abs(path); err != nil {
		log.Fatalf("determing absolute path to %q: %s", path, err)
	} else {
		for _, goPath := range filepath.SplitList(os.Getenv("GOPATH")) {
			prefix := filepath.Join(goPath, "src") + string(filepath.Separator)
			path = strings.TrimPrefix(stepsAbs, prefix)
		}
	}

	return path
}

func transformSpecs(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".spec") {
			spec := parseSpecFile(path)
			g := specGenerator{
				specsRoot: directory,
				spec:      spec,
				outPkg:    *pkg,
			}
			g.generate()
		}
		return nil
	})
}
