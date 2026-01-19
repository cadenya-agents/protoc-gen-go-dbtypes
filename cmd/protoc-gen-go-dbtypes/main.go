package main

import (
	"flag"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet

	// Flag to exclude types by name (comma-separated list)
	excludeTypes := flags.String("exclude", "", "comma-separated list of message names to exclude from generation")
	// Flag to only generate for a specific package
	onlyPackage := flags.String("package", "", "only generate for this proto package (e.g., 'example.v1')")

	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(gen *protogen.Plugin) error {
		// Declare support for proto3 optional fields
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		// Parse excluded types into a set
		excluded := make(map[string]bool)
		if *excludeTypes != "" {
			for _, name := range strings.Split(*excludeTypes, ",") {
				excluded[strings.TrimSpace(name)] = true
			}
		}

		config := &GeneratorConfig{
			ExcludedTypes: excluded,
			OnlyPackage:   strings.TrimSpace(*onlyPackage),
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
