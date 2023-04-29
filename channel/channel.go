package channel

import (
	"fmt"
	"os"
	"time"

	"github.com/eduncan911/podcast"
)

const DefaultInvidiousInstance = "https://inv.bp.projectsegfau.lt"

type Channel struct {
	InvidiousInstance string
	Author            string
	Podcast           podcast.Podcast
}

func New(author string) Channel {
	ret := Channel{
		Author:            author,
		InvidiousInstance: DefaultInvidiousInstance,
		Podcast:           podcast.Podcast{},
	}
	return ret
}

func (c Channel) StoreOutput(outputFile string) error {
	var err error
	if c.Podcast.Title == "" {
		c.Podcast, err = c.createPodcast()
		if err != nil {
			return err
		}
	}
	outputWriter := os.Stdout
	if outputFile != "" {
		outputWriter, err = os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not write output to file '%v': %v", outputFile, err)
		}
		defer outputWriter.Close()
	}
	err = c.Podcast.Encode(outputWriter)
	if err != nil {
		return fmt.Errorf("could not write the feed to the desired output: %v", err)
	}
	return nil
}

func (c Channel) createPodcast() (podcast.Podcast, error) {
	if c.Author == "" {
		return podcast.Podcast{}, fmt.Errorf("ERROR: Please specify author")
	}

	channel, err := GetChannel(c.InvidiousInstance, c.Author)

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

	return p, nil
}
