# config
config.go provides a general purpose command line argument and config file parser. All configuration values are strings with bools being 
"yes" or "no"

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
const dfltConf string = "/etc/myapp/myapp.conf"
const dfltPort string = "12345"                  // Default UDP port

func main() {
	// Make a new configuration map
	var tconfig config.Config

	// Add some options to read from the command line and/or config file
	tconfig.AddOption("verbose", "v", false, "Output log messages to the console", "no")
	tconfig.AddOption("show", "show", false, "List the current config and exit", "no")
	tconfig.AddOption("port", "p", true, "UDP port on which to listen", dfltPort)

	// Read the command line and config file
	options := tconfig.Read(dfltConf)

	// Print program usage with an optional title
	tconfig.PrintUsage("test - usage:")
	
	// Read some of the options specified by the user
	if options["verbose"] == "yes" {
	    fmt.Println("Verbose mode is set")
	}
	
	fmt.Printf("Specified port is %s\n", options["port"])
}
```

Exported resources are:
* Config type
* AddOption(name string, argument-flag string, needs-option bool, help-text string, default-val string) method
* Read(config_filename string) method
* PrintUsage(title string) method
