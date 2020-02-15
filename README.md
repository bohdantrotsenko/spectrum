# spectrum

A sample code to play with my spectrogram algorithm via web service.

Please refer to https://www.trotsenko.com.ua for the description of my invention (spectrogram algorithm).

## Prerequisites

- [ffmpeg](https://www.ffmpeg.org/)
- [go](https://golang.org/dl/) (tested with 1.13)

## Step 0. Build spectrogram -> pictures converter

    go get github.com/bohdantrotsenko/spectrum

## Step 1. Prepare data

    ffmpeg -i your_media_file -t 0:2:00 -ar 96000 -ac 1 -sample_fmt s16 temp.wav

where `-t 0:2:00` limits duration to two minutes,
`-ar 96000` sets audio sampling rate to 96khz,
`-ac 1` makes it mono 
and `-sample_fmt s16` sets format to 16bit PCM.
You can also specify e.g. `-ss 0:0:15` to set the starting point in `your_media_file`
and use other options.

## Step 2. Transform

    curl --data-binary @temp.wav https://fra.trotsenko.com.ua/v1 > temp.bin

**Note** that the service logs IP addresses and hashsums of the payload.

`temp.bin` will have a sequence of arrays.

Each array contains 456 bytes, each corresponding to frequencies
(24.499715, 24.856071, ..., 17485.357994)  
_(note the geometrical progression)_.  
Each array corresponds to 40 samples of input as is normalized.

## Step 3. See it on video

    cat temp.bin | spectrum | ffmpeg -y -i temp.wav -r 30 -i - -c:v libx264 -pix_fmt yuv420p -crf 17 -vf fps=30 -profile:v high -level 4.2 -c:a aac temp.mp4

You can then open temp.mp4 for viewing.