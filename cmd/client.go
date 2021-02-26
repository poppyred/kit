package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kit/generator"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:     "client",
	Short:   "Generate simple client lib",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide a name for the service")
			return
		}

		pbImportPath := viper.GetString("g_c_pb_import_path")
		if viper.GetString("g_c_transport") == "grpc" {
			if pbImportPath == "" {
				logrus.Error("You must provide pb import path by --pb_import_path or -i, because transport is grpc")
				return
			}
		}
		g := generator.NewGenerateClient(
			args[0],
			viper.GetString("g_c_transport"),
			pbImportPath,
		)
		if err := g.Generate(); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("transport", "t", "http", "The transport you want your client to be initiated")
	clientCmd.Flags().StringP("pb_import_path", "i", "", "Specify path to import pb")
	viper.BindPFlag("g_c_transport", clientCmd.Flags().Lookup("transport"))
	viper.BindPFlag("g_c_pb_import_path", clientCmd.Flags().Lookup("pb_import_path"))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
