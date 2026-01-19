package main

import (
	"flag"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet

	// Flag to exclude types by name (comma-separated list)
	excludeTypes := flags.String("exclude", "", "comma-separated list of message names to exclude from generation")

	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(gen *protogen.Plugin) error {
		// Parse excluded types into a set
		excluded := make(map[string]bool)
		if *excludeTypes != "" {
			for _, name := range strings.Split(*excludeTypes, ",") {
				excluded[strings.TrimSpace(name)] = true
			}
		}

		config := &GeneratorConfig{
			ExcludedTypes: excluded,
		}

		// Track which packages have had ProtoValue generated
		generatedPackages := make(map[protogen.GoImportPath]bool)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := generateFile(gen, f, config, generatedPackages); err != nil {
				return err
			}
		}
		return nil
	})
}
