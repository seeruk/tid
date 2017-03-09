package cli

import "github.com/eidolon/console"

// CreateApplication builds the console application instance. Providing it with some basic
// information like the name and version.
func CreateApplication() *console.Application {
	application := console.NewApplication("tid", "0.2.0-alpha.1")
	application.Logo = `
######## ### #######
   ###   ###       ##
   ###   ###  ###  ##
   ###   ###  ###  ##
   ###   ###  ######
`

	return application
}
