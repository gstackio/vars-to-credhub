package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/colorstring"

	"golang.org/x/crypto/ssh/terminal"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Version the program version number
var Version = "(development)"

var (
	varPrefix = kingpin.Flag("prefix", "credhub path prefix for vars").Short('p').Default("/concourse/main").String()
	inputFile = kingpin.Flag("vars-file", "Pipeline vars file").Short('f').Required().File()
)

func main() {
	var app = kingpin.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	kingpin.Parse()

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)

	if bulkImport, err := Transform(*varPrefix, *inputFile); err != nil {
		log.Fatal(err)
	} else {
		encoder.Encode(bulkImport)
		outStr := b.String()
		fmt.Println(strings.Replace(outStr, "subMapValue", "value", -1))
		// taken from https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			colorstring.Println("[bold][green]Done! Double check these results, then save and run [reset]credhub import --file /path/to/file")
		}
	}
}
