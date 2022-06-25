package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/tpl"
)

var templateFuncs = template.FuncMap{
	"title": strings.Title,
}

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0755); err != nil {
			return err
		}
	}

	// create cmd/<AppName>/
	cmdAppDir := filepath.Join(p.AbsolutePath, "cmd", p.AppName)
	if _, err := os.Stat(cmdAppDir); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(cmdAppDir, 0755))
	}

	// create cmd/<AppName>/main.go
	mainFile, err := os.Create(filepath.Join(p.AbsolutePath, "cmd", p.AppName, "main.go"))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Funcs(templateFuncs).Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create pkg/<AppName>/
	pkgAppDir := filepath.Join(p.AbsolutePath, "pkg", p.AppName)
	if _, err := os.Stat(pkgAppDir); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(pkgAppDir, 0755))
	}

	// create pkg/<AppName>/root.go
	if _, err := os.Stat(filepath.Join(pkgAppDir, "root.go")); os.IsNotExist(err) {
		rootFile, err := os.Create(filepath.Join(pkgAppDir, "root.go"))
		if err != nil {
			return err
		}
		defer rootFile.Close()

		rootTemplate := template.Must(template.New("root").Funcs(templateFuncs).Parse(string(tpl.RootTemplate())))
		err = rootTemplate.Execute(rootFile, p)
		if err != nil {
			return err
		}
	}

	// create license
	return p.createLicenseFile()
}

func (p *Project) createLicenseFile() error {
	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Funcs(templateFuncs).Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

func (c *Command) Create() error {
	// create pkg/<AppName>/commands/ if it doesn't exist
	commandsDir := filepath.Join(c.AbsolutePath, "pkg", c.AppName, "commands")
	rootFile := filepath.Join(c.AbsolutePath, "pkg", c.AppName, "root.go")
	rootFileData, err := os.ReadFile(rootFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(commandsDir, 0755))
	}

	// If needed, update root.go by replacing the comment
	// //+cobra:commandsImport
	// with "{{ .PkgName }}/pkg/{{ .AppName }}/commands"
	rootFileData = bytes.Replace(rootFileData, []byte(`//+cobra:commandsImport`), []byte(fmt.Sprintf(`"%s/pkg/%s/commands"`, c.PkgName, c.AppName)), 1)

	cmdFile, err := os.Create(filepath.Join(commandsDir, c.CmdName+".go"))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Funcs(templateFuncs).Parse(string(tpl.AddCommandTemplate())))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}

	// and add a new entry above '//+cobra:subcommands'
	rootFileData = bytes.Replace(rootFileData, []byte("//+cobra:subcommands"), []byte(fmt.Sprintf("rootCmd.AddCommand(commands.Build%sCmd())\n\t//+cobra:subcommands", strings.Title(c.CmdName))), 1)

	// write back to root.go
	if err := os.WriteFile(rootFile, rootFileData, 0644); err != nil {
		return err
	}
	return nil
}
