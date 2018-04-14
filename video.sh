#!/bin/bash
DOUBLE=1 DIRECTION="0.93;0.79;6000" ./poleshift world.jpg 0 0 1920 1080 video/f.jpeg
ffmpeg -framerate 20 -pattern_type glob -i 'video/*.jpeg' -c:v libx264 -r 20 -pix_fmt yuv420p video.mp4
