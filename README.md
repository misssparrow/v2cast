# v2cast

## About

This is a CLI tool that converts a YouTube channel into a RSS-feed usable by any Podcast application.

## Usage

Use the `-a` flag to specify the author of the channel you want to convert and use the `-o` flag to specify an output file.

```bash
v2cast -a "The Linux Foundation" -o feed.xml
```

The invidious instance, which must provide API access, to use can be specified by using the `-i` switch.
If this is omitted, a hardcoded default instance is used.

### Caveats

The tool does not provide any web server capabilities. In order to use it, you need to copy the XML file created into a webserver's directory and point your podcatcher to its URL.

The download links for the media files expire after aproximately 3 hours. To make the feed continuously usable, you have to recreate it peridocally via a cronjob:

```cron
Min Hou Dom Mon Dow   command
*   */2 *   *   *     v2cast -a "The Linux Foundation" -o /var/www/html/feed.xml
```

Downloading the media files is extremely slow.

## Build

Simply run `go build` in the source directory.
