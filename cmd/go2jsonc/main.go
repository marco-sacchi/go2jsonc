// go2jsonc standalone executable and go generator application package.
package main

import (
	"flag"
	"github.com/marco-sacchi/go2jsonc"
	"log"
	"os"
)

const version = "0.2.0"

func main() {
	flag.Usage = usage
	typeName := flag.String("type", "", "struct type name for which generate JSONC")
	output := flag.String("out", "", "output JSONC filepath")

	flag.Parse()

	if *typeName == "" {
		println("Flag -type is mandatory.\n")
		flag.Usage()
		os.Exit(1)
	}

	dirs := flag.Args()

	dir := "."
	switch len(dirs) {
	case 0:
		println("No directory specified, using current working dir.")

	case 1:
		dir = dirs[0]

	default:
		println("Only one directory can be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	code, err := go2jsonc.Generate(dir, *typeName)
	if err != nil {
		log.Fatal(err)
	}

	if *output != "" {
		err = os.WriteFile(*output, []byte(code), 0666)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		print(code)
	}

}

func usage() {
	println("go2jsonc v" + version + " Copyright 2022 Marco Sacchi")
	println()
	println("Usage:")
	println("  go2jsonc -type <type-name> [-out jsonc-filename] [package-dir]")
	println()
	println("When -out flag is omitted the code is written to stdout.")
	println("When package-dir is omitted, the current working directory will be used.")

	flag.PrintDefaults()
}
