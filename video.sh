#!/bin/bash
DOUBLE=1 DIRECTION="0.93;0.79;3000" ./poleshift world.jpg 0 0 1080 540 video/f.jpeg
ffmpeg -framerate 25 -pattern_type glob -i 'video/*.jpeg' -c:v libx264 -r 25 -pix_fmt yuv420p video.mp4
