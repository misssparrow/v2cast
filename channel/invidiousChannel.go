package channel

import (
	"fmt"
	"net/http"
	"strconv"

	"git.sr.ht/~greenfoo/govidious"
)

const DefaultInvidiousInstance = "https://inv.bp.projectsegfau.lt"

type VideoContent struct {
	Title               string
	YtUrl               string
	Description         string
	ThumbnailUrl        string
	StreamUrl           string
	StreamLengthBytes   int64
	StreamLengthSeconds int64
	PublicationDate     int64
}

type ChannelContent struct {
	Author      string
	Url         string
	Description string
	ImageUrl    string
	Videos      []VideoContent
}

func GetChannel(invidiousInstance, author string) (ChannelContent, error) {

	g := govidious.New(invidiousInstance, nil)
	ret := ChannelContent{}

	//channels, err := g.ChannelsSearch(author)
	channels, err := g.Search(author, 1, "", "", "", "", "", "")
	if err != nil {
		return ChannelContent{}, fmt.Errorf("could not find the author: %v", err)
	}

	authorId := ""
	for _, channel := range channels.Channels {
		if channel.Author == author {
			authorId = channel.AuthorId
			ret.Author = author
			ret.Url = channel.AuthorUrl
			ret.Description = channel.Description
			//The URL returned by AuthorThumbnails starts with // instead of https. Probably a bug in govidious
			authorImage := "https:" + channel.AuthorThumbnails[len(channel.AuthorThumbnails)-1].Url
			ret.ImageUrl = authorImage
			break
		}

	}

	if authorId == "" {
		return ChannelContent{}, fmt.Errorf("channel not found")
	}

	channel, err := g.ChannelsVideos(authorId, "")
	if err != nil {
		return ChannelContent{}, fmt.Errorf("could not find videos in channel: %v", err)
	}

	ret.Videos = make([]VideoContent, len(channel.Videos))
	for i, video := range channel.Videos {

		videoInfos, err := g.Videos(video.VideoId, "")
		if err != nil {
			// We continue if there is an error with a single video
			continue
		}
		for _, f := range videoInfos.AdaptiveFormats {
			if f.Container == "m4a" {
				streamResp, err := http.Head(f.Url)
				if err != nil {
					fmt.Printf("Stream not accessible: %v\n", err)
					continue
				}
				streamLength, err := strconv.ParseInt(streamResp.Header.Get("Content-Length"), 10, 64)
				if err != nil {
					fmt.Printf("Can't get length of stream: %v\n", err)
					continue
				}

				videoContent := VideoContent{
					Title:               video.Title,
					YtUrl:               "https://www.youtube.com/watch?v=" + video.VideoId,
					Description:         videoInfos.Description,
					ThumbnailUrl:        videoInfos.VideoThumbnails[len(videoInfos.VideoThumbnails)-1].Url,
					StreamUrl:           f.Url,
					StreamLengthBytes:   streamLength,
					StreamLengthSeconds: int64(video.LengthSeconds),
					PublicationDate:     videoInfos.Published,
				}
				ret.Videos[i] = videoContent
				break
			}
		}
	}

	return ret, nil
}
