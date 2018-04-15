# poleshift
Allows to choose earth map image (or any other) and regenerate it with North/South poles shifted

# Config

Looking from up to the North Pole.

(lon 0, lat 0) is to the right: (1, 0, 0).

(lon 90, lat 0) is to the front: (0, 1, 0).

(lon +/-180, lat 0) is to the left: (-1, 0, 0).

(lon -90, lat 0) is to the back: (0, -1, 0).

North Pole (lon ?, lat 90) is up: (0, 0, 1).

South Pole (lon ?, lat -90) is down: (0, 0, -1).

First rotation is clockwise by using Y axis which goes through (90, 0) and (-90, 0).

Second rotation is clockwise by using Z axis which goes through unrotated North and South poles.

# Usage

Provide either:
- 4 args: lon lat alpha beta
- 6 args: input_file alpha beta output_width output_height output_name

Use `DOUBLE=1` env variable to make x resultion double, to include wrapping 50% from the left and 50% from the right.

Use `DIRECTION="alpha_inc;beta_inc;n_frames"` to generate `n_fromes` output files (frame num prefix prepended) changing alpha and beta by `alpha_inc`, `beta_inc` each step.

# Generate images

- Use `images.sh` script to generate example rotation images for all 15 degrees rotations for alpha and beta angles

# Generate video

- Use `video.sh` script that will create final 3000 frames which means 120s (2 minutes) of 25 FPS movie with rotating poles, file is is not available directly (it is 264Mb forbiden by GitHub).
- It is available on YouTube [here](https://www.youtube.com/watch?v=MwESyNhfXYg). make sure you select HD quality.
- Another YT vide is [here](https://www.youtube.com/watch?v=0F8KWVFGkNo&feature=youtu.be).

