/* Copyright (c) 2018-2019 Rubicon Communications, LLC (Netgate)
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// config.go contaains all of the funtions required to process the comamnd line arguments, config file values and defaults
// This was started as an exercise to learn Go flags, methods, structures and maps, but has turned out to be useful here
// This file can be moved to its own package, or incorporated in another project as here.
package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// A Config is a list of configuration Items that specify the option details
type Config struct {
	//	filename string
	Items []ConfigItem
}

type ConfigItem struct {
	Name   string // The name of this config item (used  as a map key)
	Arg    string // Command line argument that sets it
	Hasval bool   // Does this command line flag have an associated value string
	Descr  string // Description of the item used in constructing usage/help
	Dflt   string // Default value for this item
}

// Add a new config item specification to the configuration parser
func (cfg *Config) AddOption(name string, arg string, hasval bool, descr string, dflt string) {
	cfg.Items = append(cfg.Items, ConfigItem{name, arg, hasval, descr, dflt})
}

// Print a table of options and help strings
func (cfg Config) PrintUsage(title string) {
	option := ""

	fmt.Println(title)
	for idx := 0; idx < len(cfg.Items); idx++ {
		if len(cfg.Items[idx].Arg) == 0 {
			continue
		}

		if cfg.Items[idx].Hasval {
			option = fmt.Sprintf("  -%s <%s>", cfg.Items[idx].Arg, cfg.Items[idx].Name)
		} else {
			option = fmt.Sprintf("  -%s", cfg.Items[idx].Arg)
		}

		fmt.Printf("   %-20s : %s\n", option, cfg.Items[idx].Descr)
	}
}

// Read the command line arguments
// Read the config file values
// Combine the two plus the defaults
func (cfg *Config) Read(cfgname string) (map[string]string, error) {
	cfgpath := ""

	// These two options are added by default so the program knows where to find the config file
	// and can provide help
	cfg.AddOption("help", "h", false, "Output usage information to the console", "no")
	if len(cfgname) > 0 {
		cfg.AddOption("cfgpath", "c", true, "Path to configuration file", cfgname)
	}

	argmap := cfg.readArgs()

	if len(argmap["cfgpath"]) > 0 {
		cfgpath = argmap["cfgpath"]
	} else {
		cfgpath = cfgname
	}

	confmap := make(map[string]string)
	var err error
	if len(cfgname) > 0 {
		confmap, err = readConfigFile(cfgpath)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}
	}

	return cfg.mergeItems(argmap, confmap), nil
}

// Read the command line arguments by creating a flag entry for each option, then parsing the flags
func (cfg Config) readArgs() map[string]string {
	args := make(map[string]*string)
	boolargs := make(map[string]*bool)
	combo := make(map[string]string)

	// Options expecting sting arguments, and boolean options (which do not) are added differently
	for idx := 0; idx < len(cfg.Items); idx++ {
		if cfg.Items[idx].Hasval {
			args[cfg.Items[idx].Name] = flag.String(cfg.Items[idx].Arg, "", cfg.Items[idx].Descr)
		} else {
			boolargs[cfg.Items[idx].Name] = flag.Bool(cfg.Items[idx].Arg, false, cfg.Items[idx].Descr)
		}
	}

	flag.Parse()

	// Now that there is a map of pointers to command line options, translate that to a map of strings
	for k, v := range boolargs {
		if *v {
			combo[k] = "yes"
		} else {
			combo[k] = ""
		}
	}

	for k, v := range args {
		combo[k] = *v
	}

	return combo
}

// If a command line argument is provided, use it, otherwise use the config file value or the default
func merge(arg string, conf string, dflt string) string {
	if len(arg) == 0 {
		if len(conf) != 0 {
			return conf
		} else {
			return dflt
		}
	}

	return arg
}

// Iterate over the list of options, merging the command line, config file and defaults
func (cfg Config) mergeItems(args map[string]string, conf map[string]string) map[string]string {
	mergedmap := make(map[string]string)

	for _, ci := range cfg.Items {
		mergedmap[ci.Name] = merge(args[ci.Name], conf[ci.Name], ci.Dflt)
	}

	return mergedmap
}

// Read a config file and return its contents in a map
// There are many Go config file packages available, but most are more complicated than needed here
func readConfigFile(filename string) (map[string]string, error) {
	cfg := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return cfg, errors.New("Unable to open configuration file. Using default values")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Ignore comment lines
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

		s := strings.SplitN(scanner.Text(), "=", 2)
		// Ignore mal-formed lines
		if len(s) != 2 {
			continue
		}

		// Trim white space from front and back, delete any quotes and make the key lower case
		cfg[strings.ToLower(strings.TrimSpace(s[0]))] = strings.Replace(strings.TrimSpace(s[1]), "\"", "", -1)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
