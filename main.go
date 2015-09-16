package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var Version string

var formats = map[string]func(string, []FormatVariable) string{
	"export":  formatExport,
	"envfile": formatEnvfile,
}

func main() {
	app := cli.NewApp()
	app.Name = "envar"
	app.Usage = "Manage environment variable at one place"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:   "print",
			Usage:  "Print environment variables",
			Action: printCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "envar.yml",
					Usage: "environment variables definition file",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "export",
					Usage: fmt.Sprintf("output format of variables [%s]", strings.Join(getFormatsList(), ", ")),
				},
			},
		},
	}

	app.Run(os.Args)
}

func printCmd(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Environment name must be specified\n")
		os.Exit(1)
	}
	environmentName := c.Args()[0]
	format, ok := formats[c.String("output")]
	if !ok {
		fmt.Fprintf(os.Stderr, "Output format must be one of [%s]\n", strings.Join(getFormatsList(), ", "))
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(c.String("file"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read YAML file: %s\n", err)
		os.Exit(1)
	}
	config, errs := parse(string(contents))
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(1)
	}
	if (environmentExists(config.Environments, environmentName) == false) {
		fmt.Fprintf(os.Stderr, "No such a environment: %s\n", environmentName)
		os.Exit(1)
	}

	variables := FormatVariables{}
	for name, values := range config.Variables {
		variables = append(variables, FormatVariable{name, values[environmentName]})
	}
	sort.Sort(variables)

	fmt.Fprintf(os.Stdout, format(environmentName, variables))
}

func getFormatsList() []string {
	list := []string{}
	for k, _ := range formats {
		list = append(list, k)
	}
	return list
}

type FormatVariable struct {
	Name  string
	Value interface{}
}

type FormatVariables []FormatVariable

func (p FormatVariables) Len() int {
	return len(p)
}

func (p FormatVariables) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p FormatVariables) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
