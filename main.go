package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jamestjw/feedme/instagram"
)

func main() {
	var targetID = flag.String("t", "", "ID of account to generate feeds for (required)")
	var outputFilename = flag.String("o", "", "Output filename")
	var allowPrivate = flag.Bool("p", false, "Allow returning feed of private accounts (only minimum amount of data will be returned)")
	flag.Parse()

	if *targetID == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := instagram.FetchFeed(targetID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !*allowPrivate && data.Data.User.IsPrivate {
		fmt.Fprintf(os.Stderr, "%s is a private account, run with -p to bypass this\n", *targetID)
		os.Exit(1)
	}

	outString, err := instagram.GenerateXML(data)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *outputFilename != "" {
		err := ioutil.WriteFile(*outputFilename, []byte(outString), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Println(outString)
	}
}
