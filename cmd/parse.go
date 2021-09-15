package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

var parseCommand = &cobra.Command{
	Use:   "parse <path>",
	Short: "Parse a fixture into JSON format",
	Long: `Parse a fixture into JSON format.

The 'parse' command takes the path to a fixture and returns the JSON format that 
would need to be used when writing the Rego rules.

For example, to parse a Terraform file run the following command:
$ snyk-iac-rules parse test.tf --format hcl2
The '--format' flag can be left out when parsing Terraform files, as we default to hcl2.

And to parse a YAML file run the following command:
$ snyk-iac-rules parse test.yaml --format yaml

The output of this command can be used when writing tests. Run the following command to find out how:
$ snyk-iac-rules test --help
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Expected a path to be provided via the command")
		}
		if len(args) > 1 {
			return errors.New("Too many paths provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.RunParse(args, parseParams)
	},
}

func newParseCommandParams() *internal.ParseCommandParams {
	return &internal.ParseCommandParams{
		Format: util.NewEnumFlag(internal.HCL2, []string{internal.HCL2, internal.YAML}),
	}
}

var parseParams = newParseCommandParams()

func init() {
	parseCommand.Flags().VarP(&parseParams.Format, "format", "f", "choose the format for the parser")
	RootCommand.AddCommand(parseCommand)
}
