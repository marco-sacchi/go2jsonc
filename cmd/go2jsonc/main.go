// go2jsonc standalone executable and go generator application package.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/marco-sacchi/go2jsonc"
)

const version = "0.3.2"

func main() {
	flag.Usage = usage
	typeName := flag.String("type", "", "struct type name for which generate JSONC; mandatory")
	docTypeMode := flag.String("doc-types", "",
		"pipe-separated bits representing struct fields types for which do not\n"+
			"render the type in JSONC comments; when omitted all types will be\nrendered for all fields")
	output := flag.String("out", "", "output JSONC filepath; when omitted the code is written to stdout")

	flag.Parse()

	if *typeName == "" {
		println("Flag -type is mandatory.\n")
		flag.Usage()
		os.Exit(1)
	}

	docMode := go2jsonc.AllFields

	if *docTypeMode != "" {
		bits := strings.Split(*docTypeMode, "|")
		for _, bit := range bits {
			if bit == "NotFields" {
				docMode = go2jsonc.NotFields
			} else {
				switch bit {
				case "NotStructFields":
					docMode |= go2jsonc.NotStructFields

				case "NotArrayFields":
					docMode |= go2jsonc.NotArrayFields

				case "NotMapFields":
					docMode |= go2jsonc.NotMapFields

				default:
					fmt.Printf("Invalid bit name %s for -doc-types flag.\n\n", bit)
					flag.Usage()
					os.Exit(1)
				}
			}
		}
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

	code, err := go2jsonc.Generate(dir, *typeName, docMode)
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
	println("go2jsonc v" + version + " Copyright 2022 Marco Sacchi\n")

	println("Usage:")
	println("  go2jsonc -type <type-name> [-doc-types bits] [-out jsonc-filename] [package-dir]\n")

	flag.PrintDefaults()

	println("\npackage-dir: directory that contains the go file where specified type is")
	println("defined; when omitted, current working directory will be used\n")

	println("Allowed constants for -doc-types flag:")
	println("  NotFields        Does not display type in all fields;")
	println("  NotStructFields  Does not display type in fields of type struct;")
	println("  NotArrayFields   Does not display type in fields of type array or slice;")
	println("  NotMapFields     Does not display type in fields of type map.")
}
