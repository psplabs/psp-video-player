# Perspectives Video Player

An MPEG-DASH video streaming solution for Perspectives Study Program, which supports both hosted and offline / downloadable videos.

### Background / Introduction

We're currently using Vimeo to host and stream the Perspectives Instructor Video Sessions, by embedding them in an iframe within a modal. This is very insecure because anyone can inspect the frame source and locate the video's download URL.

This repo is a proof of concept demonstration for the following requirements:

 1) Secure, self-hosted, adaptive bitrate streaming
 2) Support for all major, modern browsers and devices
 3) Downloadable, offline version that expires after X days or after class is finished (for classes without dependable internet)

I decided to write this demo in Google's GO programming language, primarily because it is a cross-platform, fully-compiled language and thanks to [packr](https://github.com/gobuffalo/packr) we're able to compile the web app, videos, and all static assets down into a single binary executable that works out of the box! (Seriously, how cool is that!?!)

## Getting Started

The two dependencies are GO and Bento4's SDK for preparing your MPEG DASH videos for streaming.

### Prerequisites

First, make sure you have GO installed (or [download it here](https://golang.org/dl/))):

```
brew install go
```

You'll also need Bento4 to encode / transmux your videos into DASH compatible mpd/m4s. There's a homebrew version, but it doesn't come with the `mp4dash` python script, so just download the SDK from [https://www.bento4.com/downloads/](https://www.bento4.com/downloads/) and either add it to your PATH or move it over to `/usr/local/Cellar/bento4`.


### Installing

Clone this repo down locally:

```
git clone git@github.com:psplabs/psp-video-player.git
```

Next, install the GO dependencies:

```
go get -u github.com/gobuffalo/packr/...
go get github.com/djherbis/times
```


## Usage

### Prepare your videos for Streaming

[Read this guide](https://www.bento4.com/developers/dash/) for instructions on how to fragment your video and generate DASH presentations. After you generate the `mpd` file, rename it "stream.mpd" and move it and the `audio` and `video` directories over to the `media` directory in the root of this project. The build script will copy them over to `/standalone/stream` for streaming.

You can run the app by calling the `./build` script and passing it the directory name of your video under `/media` (in this example, we'll pretend we have a video called `demo` under `/media/demo`):

```
./build demo
```

### /hosted

This folder contains an example of a self-hosted player with the "Download" UI.

### /standalone

This directory gets compiled into a single binary and downloaded to the user's computer. They can copy it and move it to another computer, but the binary will self destruct after 14 days.

## Authors

* **Daniel Bodnar** - *Initial work* - [DanielBodnar](https://github.com/DanielBodnar)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* The awesome dashjs team!
* Bento4 for the awesome mp4dash tools
* Google for GoLang
