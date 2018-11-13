# config
congig.go is a general purpose command line argument and config file parser.

1. Create a new config
2. Add config item specifications to that config
3. Call the Read() method to read the command line args, config file items and combine them (with optional default values)
4. Use as required by looking up the config by name in the map

E.g.:
```go
package main

import ( 
	"fmt"
	"config"
)


// Configuration defaults
const dfltConf string = "/etc/tnsrids/tnsrids.conf"
const dfltPort string = "12345"                  // Default UDP port on whic alert messages are received

func main() {
	var tconfig config.Config

	tconfig.AddOption("verbose", "v", false, "Output log messages to the console", "no")
	tconfig.AddOption("show", "show", false, "List the current block rules and exit", "no")
	tconfig.AddOption("port", "p", true, "UDP port on which to listen for alert messages", dfltPort)

	options := tconfig.Read(dfltConf)

	fmt.Printf("Configuration: %v\n", options)

	tconfig.PrintUsage("test - usage:")
	
	if tconfig['verbose'] == "yes" {
	    fmt.Println("Verbose mode is set")
	}
}
```

Exported resources are:
* Config type
* AddOption(name string, argument-flag string, needs-option bool, help-text string, default-val string) method
* Read(config_filename string) method
* PrintUsage(title string) method
