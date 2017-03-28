package main

import (
	"os"

	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

var name = "World"
var favNum int

func main() {
	application := console.NewApplication("eidolon/console", "0.1.0")
	application.Logo = `
                                             #
                              ###            ##
######## ### ####### #######  ###   #######  ###  ##
         ###       ##      ## ###         ## #### ##
 ####### ###  ###  ## ##   ## ###    ##   ## #######
 ###     ###  ###  ## ##   ## ###    ##   ## ### ###
 ####### ###  ######   #####  ####### #####  ###  ##
                                                   #
`

	application.AddCommand(&console.Command{
		Name:        "greet:example",
		Description: "Greet's the given user, or the world.",
		Help:        "You don't have to specify a name.",
		Configure: func(definition *console.Definition) {
			definition.AddOption(console.OptionDefinition{
				Value:  parameters.NewStringValue(&name),
				Spec:   "-n, --name=NAME",
				Desc:   "Provide a name for the greeting.",
				EnvVar: "EXAMPLE_NAME",
			})

			definition.AddArgument(console.ArgumentDefinition{
				Value: parameters.NewIntValue(&favNum),
				Spec:  "FAVOURITE_NUMBER",
				Desc:  "Provide your favourite number.",
			})
		},
		Execute: func(input *console.Input, output *console.Output) error {
			output.Printf("Hello, %s!\n", name)
			output.Printf("Your favourite number is %d.\n", favNum)
			return nil
		},
	})

	code := application.Run(os.Args[1:], os.Environ())

	os.Exit(code)
}
