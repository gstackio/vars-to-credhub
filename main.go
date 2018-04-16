package main

import (
	"log"
	"os"

	"github.com/mitchellh/colorstring"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	varPrefix = kingpin.Flag("prefix", "credhub path prefix for vars").Short('p').Default("/concourse/main").String()
	inputFile = kingpin.Flag("vars-file", "Pipeline vars file").Short('f').Required().File()
)

func main() {
	kingpin.Parse()
	encoder := yaml.NewEncoder(os.Stdout)

	if bulkImport, err := Transform(*varPrefix, *inputFile); err != nil {
		log.Fatal(err)
	} else {
		encoder.Encode(bulkImport)
		// taken from https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			colorstring.Println("[bold][green]Done! Save this file and run [reset]credhub bulk-import --file /path/to/file")
		}
	}
}
