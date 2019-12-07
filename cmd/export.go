package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jamesdobson/gogitmoji/tmpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a commit template",
	Long:  `Export a commit template.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("The export command expects a single argument: the name of the commit template to export.\n")
		}

		export(args[0])
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func export(templateName string) {
	templates := viper.GetStringMap("templates")
	tmpl.LoadTemplates(templates)

	t, ok := tmpl.TemplateLookup[templateName]

	if !ok {
		log.Fatalf("\nUnknown commit template: \"%s\"\n\n", templateName)
	}

	var result map[string]interface{}

	err := mapstructure.Decode(t, &result)

	if err != nil {
		log.Fatalf("\nError converting template: %v\n\n", err)
	}

	var wrapper = map[string]interface{}{
		"templates": map[string]interface{}{
			templateName: result,
		},
	}

	out, err := yaml.Marshal(wrapper)

	if err != nil {
		log.Fatalf("\nUnable to output template as YAML: %v\n\n", err)
	}

	_, err = os.Stdout.Write(out)

	if err != nil {
		log.Fatalf("\nUnable to write output: %v\n\n", err)
	}

	fmt.Println()
}
