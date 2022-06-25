// Copyright © 2021 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

func MainTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package main

import "{{ .PkgName }}/pkg/{{ .AppName }}"

func main() {
	{{ .AppName }}.Execute()
}
`)
}

func RootTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package {{ .AppName }}

import (
{{- if .Viper }}
	"fmt"{{ end }}
	"os"

	//+cobra:commandsImport
	"github.com/spf13/cobra"
{{- if .Viper }}
	"github.com/spf13/viper"{{ end }}
)

{{ if .Viper -}}
var cfgFile string
{{- end }}

// rootCmd represents the base command when called without any subcommands
func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "{{ .AppName }}",
		Short: "A brief description of your application",
		Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	//+cobra:subcommands

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
{{ if .Viper }}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ else }}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ end }}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := BuildRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

{{- if .Viper }}
func init() {
	cobra.OnInitialize(initConfig)
}
{{ end }}

{{ if .Viper -}}
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".{{ .AppName }}" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{{ .AppName }}")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
{{- end }}
`)
}

func AddCommandTemplate() []byte {
	return []byte(`/*
{{ .Project.Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// {{ title .CmdName }}Cmd represents the {{ .CmdName }} command
func Build{{ title .CmdName }}Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{{ .CmdName }}",
		Short: "A brief description of your command",
		Long: ` + "`" + `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("{{ .CmdName }} called")
		},
	}
	return cmd
}
`)
}
