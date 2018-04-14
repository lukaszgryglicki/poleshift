#!/bin/sh
for alpha in 0 15 30 45 60 75 90 105 120 135 150 165 180 195 210 225 240 255 270 285 300 315 330 345
do
  for beta in 0 15 30 45 60 75 90 105 120 135 150 165 180 195 210 225 240 255 270 285 300 315 330 345
  do
    echo "alpha: ${alpha}, beta: ${beta}..."
    ./poleshift world.jpg "${alpha}" "${beta}" 5400 3200 "./images/rotated_${alpha}_${beta}.jpg"
    ls -l "./images/rotated_${alpha}_${beta}.jpg"
  done
done
