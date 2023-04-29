package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/misssparrow/v2cast/channel"
)

func storeOutput(p podcast.Podcast, outputFile string) error {

	outputWriter := os.Stdout
	var err error
	if outputFile != "" {
		outputWriter, err = os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not write output to file '%v': %v", outputFile, err)
		}
		defer outputWriter.Close()
	}
	err = p.Encode(outputWriter)
	if err != nil {
		return fmt.Errorf("could not write the feed to the desired output: %v", err)
	}
	return nil
}

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

	channel, err := channel.GetChannel(invidiousInstance, author)

	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		os.Exit(1)
	}

	var p podcast.Podcast
	now := time.Now()
	p = podcast.New(channel.Author, channel.Url, channel.Description, &now, &now)
	//The URL returned by AuthorThumbnails starts with // instead of https. Probably a bug in govidious
	p.AddImage(channel.ImageUrl)

	for _, video := range channel.Videos {
		podcastItem := podcast.Item{
			Title:       video.Title,
			Link:        video.YtUrl,
			Description: video.Description,
		}
		podcastItem.AddImage(video.ThumbnailUrl)
		podcastItem.AddEnclosure(video.StreamUrl, podcast.M4A, video.StreamLengthBytes)
		podcastItem.AddDuration(int64(video.StreamLengthSeconds))
		pubDate := time.Unix(video.PublicationDate, 0)
		podcastItem.AddPubDate(&pubDate)
		// We need to manually set the GUID to the YT-URL of the video. Otherwise, it would change on every creation (as the media link changes)
		podcastItem.GUID = podcastItem.Link
		_, err = p.AddItem(podcastItem)
		if err != nil {
			fmt.Printf("Error adding item: %v\n", err)
		}
	}

	// Processing the feed finished. Write the output

	err = storeOutput(p, outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write feed to output: %v\n", err)
		os.Exit(1)
	}
}
