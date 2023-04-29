package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/misssparrow/v2cast/channel"
)

func main() {
	var author string
	var ignoreTlsCertificateVerification bool
	outputFile := ""
	var invidiousInstance string

	flag.StringVar(&author, "a", "", "Author of the Youtube channel to parse")
	flag.StringVar(&outputFile, "o", "", "Outputfile to store the feed into. File will be overwritten. If this is not present, the feed will be printed to stdout.")
	flag.StringVar(&invidiousInstance, "i", channel.DefaultInvidiousInstance, "URL of the Invidious instance to use. Need to offer API access. A list is available at https://api.invidious.io/")
	flag.BoolVar(&ignoreTlsCertificateVerification, "t", false, "Ignore TLS certificate errors (use with caution)")

	flag.Parse()

	if author == "" {
		fmt.Printf("ERROR: Please specify author\n")
		flag.Usage()
		os.Exit(-1)
	}

	if ignoreTlsCertificateVerification {
		// Configure the HTTP library to skip TLS certificate verification
		// This might be necessary on old systems not having access to current certificates
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	channel := channel.New(author)
	channel.InvidiousInstance = invidiousInstance

	err := channel.StoreOutput(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write feed to output: %v\n", err)
		os.Exit(1)
	}
}
