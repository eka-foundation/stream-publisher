# stream-publisher
Creates video stream and serves the files on network.

It also publishes the service information.

## Quick start

Prequisite is to have a functional Go toolchain installed. See https://golang.org/dl/ for more info.

- git clone
- go build
- ./stream-publisher

## FFMPEG : Building from source on ubuntu 16.04LTS and later
Reference: https://trac.ffmpeg.org/wiki/CompilationGuide/Ubuntu
- Requirements
  * Install `libx264-dev` `libasound2-dev` `build-essential`
- Download FFMPEG (4.1.3) source
  * `wget https://ffmpeg.org/releases/ffmpeg-4.1.3.tar.bz2`
- Unzip source file to a directory
  * `tar -xjf ffmpeg-4.1.3.tar.bz2`
- Configure FFMPEG
 * `cd <src_dir>`
 * `./configure --enable-indev=alsa --enable-outdev=alsa --enable-gpl --enable-libx264`
- Compile and install
 * `make -j3`
 * `sudo make install`

## Run ffmpeg to produce HLS stream from webcam
- `ffmpeg -f v4l2 -i /dev/video0 -vcodec h264 -f alsa -i pulse -acodec aac -strict experimental -pix_fmt yuv420p -hls_playlist_type event -hls_flags append_list <dd_mm_yyyy_unique-title>.m3u8`
- Reference: https://ffmpeg.org/ffmpeg-formats.html (See `hls` section)
