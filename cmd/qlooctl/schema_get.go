package main

import (
	"github.com/spf13/cobra"
	"github.com/solo-io/qloo/pkg/api/types/v1"
	"github.com/pkg/errors"
	"io/ioutil"
)

var schemaGetCmd = &cobra.Command{
	Use:   "get NAME",
	Short: "return a schema by its name",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.Errorf("requires exactly 1 argument")
		}
		msg, err := getSchema(args[0])
		if err != nil {
			return err
		}
		return printAsYaml(msg)
	},
}

func init() {
	schemaCmd.AddCommand(schemaGetCmd)
}

func getSchema(name string) (*v1.Schema, error) {
	cli, err := makeClient()
	if err != nil {
		return nil, err
	}
	return cli.V1().Schemas().Get(name)
}

var createSchemaOpts struct {
	FromFile       string
	UseResolverMap string
}{}

var schemaCreateCmd = &cobra.Command{
	Use:   "create NAME --from-file <path/to/your/graphql/schema>",
	Short: "upload a schema to QLoo from a local GraphQL Schema file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.Errorf("requires exactly 1 argument")
		}
		if err := createSchema(args[0], createSchemaOpts.FromFile, createSchemaOpts.UseResolverMap); err != nil {
			return err
		}

	},
}

func init() {
	schemaCreateCmd.PersistentFlags().StringVarP(&createSchemaOpts.FromFile, "from-file", "f", "", "path to a "+
		"graphql schema file from which to create the QLoo schema object")
	schemaCreateCmd.PersistentFlags().StringVarP(&createSchemaOpts.UseResolverMap, "resolvermap", "r", "", "The name of a "+
		"ResolverMap to connect to this Schema. If none is specified, an empty ResolverMap will be generated for you, which "+
		"you can then configure with qlooctl")
	schemaCmd.AddCommand(schemaCreateCmd)
}

func createSchema(name, filename, resolvermap string) error {
	if name == "" {
		return errors.Errorf("schema name must be set")
	}
	if filename == "" {
		return errors.Errorf("filename must be set")
	}
	cli, err := makeClient()
	if err != nil {
		return err
	}
	inlineSchemaBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	schema := &v1.Schema{
		Name:         name,
		InlineSchema: string(inlineSchemaBytes),
		ResolverMap:  resolvermap,
	}
	_, err = cli.V1().Schemas().Create(schema)
	return err
}