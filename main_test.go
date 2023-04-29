package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eduncan911/podcast"
)

func Test_storeOutput(t *testing.T) {
	f, err := os.CreateTemp("", "v2casttest")
	if err != nil {
		log.Fatal(err)
	}
	tmpFile := f.Name()
	f.Close()
	defer os.Remove(f.Name())
	f, err = os.CreateTemp("", "v2casttestro")
	if err != nil {
		log.Fatal(err)
	}
	tmpFileRo := f.Name()
	f.Close()
	defer os.Remove(f.Name())
	err = os.Chmod(tmpFileRo, 0444)
	if err != nil {
		log.Fatal(err)
	}

	podcastDate := time.Date(1970, 1, 1, 22, 33, 00, 00, time.UTC)
	testPodcast := podcast.New("Footitle", "https://foo.acme", "Foo description", &podcastDate, &podcastDate)

	type args struct {
		p          podcast.Podcast
		outputFile string
		wantErr    bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Write to file",
			args{
				testPodcast,
				tmpFile,
				false,
			},
		},
		{"Overwrite to existing file",
			args{
				testPodcast,
				tmpFile,
				false,
			},
		},
		{"Writing to readonly file",
			args{testPodcast, tmpFileRo, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storeOutput(tt.args.p, tt.args.outputFile)
			if err != nil {
				if !tt.args.wantErr {
					t.Errorf("Error while storing the file: %v", err)
					return
				}
			} else {
				if tt.args.wantErr {
					t.Errorf("No error occured, even though one was expeced")
				}
			}

			storedFile, err := os.OpenFile(tmpFile, os.O_RDONLY, 0)
			if err != nil {
				t.Errorf("File could not be opened: %v", err)
				return
			}
			defer storedFile.Close()
			contentBin, err := ioutil.ReadAll(storedFile)
			if err != nil {
				t.Errorf("File could not be read: %v", err)
				return
			}
			feedWritten := string(contentBin)
			if len(feedWritten) == 0 {
				t.Errorf("Feed written is empty")
				return
			}
			if !strings.HasPrefix(feedWritten, "<?xml") || !strings.HasSuffix(feedWritten, "</rss>") {
				t.Errorf("Feed written to file does not have the expected format")
				return
			}
		})

	}
}
