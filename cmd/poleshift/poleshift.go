package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"
)

// input: lon, lat in degrees
// alpha and beta rotation angles in degrees
// alpha rotation is clockwise by using Y axis which goes through (90, 0) and (-90, 0)
// beta rotation is clockwise by using Z axis which goes through unrotated North and South poles
func poleShift(lon, lat, a, b float64) (float64, float64) {
	// (lon, lat): convert to radians
	lonR := lon * math.Pi / 180.0
	latR := lat * math.Pi / 180.0
	// Convert to radians
	aR := a * math.Pi / 180.0
	bR := b * math.Pi / 180.0
	// Convert to x,y,z coord system:
	// Looking from up to the North Pole, (lon 0, lat 0) is to the right (1, 0, 0)
	// So (lon 90, lat 0) is (0, 1, 0)
	// (lon ?, lat 90) is (0, 0, 1) North Pole
	x := math.Cos(lonR) * math.Cos(latR)
	y := math.Sin(lonR) * math.Cos(latR)
	z := math.Sin(latR)
	// Rotate
	// First rotation is clockwise by using Y axis which goes through (90, 0) and (-90, 0)
	// Second rotation is clockwise by using Z axis which goes through unrotated North and South poles
	nx := math.Cos(aR)*math.Cos(bR)*x + math.Sin(bR)*y + math.Sin(aR)*math.Cos(bR)*z
	ny := -math.Cos(aR)*math.Sin(bR)*x + math.Cos(bR)*y - math.Sin(aR)*math.Sin(bR)*z
	nz := -math.Sin(aR)*x + math.Cos(aR)*z
	// fmt.Printf("Intermediate: (%f,%f,%f) -> (%f,%f,%f)\n", x, y, z, nx, ny, nz)
	// Convert to new (lon, lat) with rotated poles
	nlonR := math.Atan2(ny, nx)
	nlatR := math.Asin(nz)
	// Convert to degrees
	nlon := nlonR / (math.Pi / 180.0)
	nlat := nlatR / (math.Pi / 180.0)
	return nlon, nlat
}

func poleShiftArgs() error {
	// Get longitude
	lon, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		return err
	}
	if lon < -180.0 || lon > 180.0 {
		return fmt.Errorf("longitude must be from [-180.0, 180.0] range")
	}
	// Get latitude
	lat, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		return err
	}
	if lat < -90.0 || lat > 90.0 {
		return fmt.Errorf("latitude must be from [-90.0, 90.0] range")
	}
	// Get alpha
	a, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		return err
	}
	if a < 0.0 || a > 360.0 {
		return fmt.Errorf("alpha must be from [0.0, 360.0] range")
	}
	// Get beta
	b, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		return err
	}
	if b < 0.0 || b > 360.0 {
		return fmt.Errorf("beta must be from [0.0, 360.0] range")
	}
	nlon, nlat := poleShift(lon, lat, a, b)
	fmt.Printf("(lon %f,lat %f) rotated by (%f, %f) ==> (lon %f, lat %f)\n", lon, lat, a, b, nlon, nlat)
	return nil
}

// mapWithPolesShifted: args: input_name alpha beta output_width output_height output_name
func mapWithPolesShifted(args []string) error {
	// If DOUBLE env is set, it will create 2x wider image:
	// 50%{actal_image}50%} by wrapping
	dbl := os.Getenv("DOUBLE") != ""

	// Get alpha
	a, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return err
	}
	if a < 0.0 || a > 360.0 {
		return fmt.Errorf("alpha must be from [0.0, 360.0] range")
	}

	// Get beta
	b, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}
	if b < 0.0 || b > 360.0 {
		return fmt.Errorf("beta must be from [0.0, 360.0] range")
	}

	// Multiple steps in a given direction?
	dir := os.Getenv("DIRECTION")
	var (
		steps  int
		prefix bool
		alphas []float64
		betas  []float64
	)
	if dir == "" {
		alphas = append(alphas, a)
		betas = append(betas, b)
		steps = 1
		prefix = false
	} else {
		ary := strings.Split(dir, ";")
		if len(ary) != 3 {
			return fmt.Errorf("required 3 args 'alphaInc;betaInc;steps', go '%s'", dir)
		}
		alphaInc, err := strconv.ParseFloat(ary[0], 64)
		if err != nil {
			return err
		}
		if alphaInc < -180.0 || alphaInc > 180.0 {
			return fmt.Errorf("alpha increment must be from [-180, 180], got %f", alphaInc)
		}
		betaInc, err := strconv.ParseFloat(ary[1], 64)
		if err != nil {
			return err
		}
		if betaInc < -180.0 || betaInc > 180.0 {
			return fmt.Errorf("beta increment must be from [-180, 180], got %f", betaInc)
		}
		steps, err = strconv.Atoi(ary[2])
		if err != nil {
			return err
		}
		if steps < 2 {
			return fmt.Errorf("there should be at least 2 steps in incremental mode, got %d", steps)
		}
		currAlpha := a
		currBeta := b
		for f := 0; f < steps; f++ {
			alphas = append(alphas, currAlpha)
			betas = append(betas, currBeta)
			currAlpha += alphaInc
			currBeta += betaInc
			if currAlpha >= 360.0 {
				currAlpha -= 360.0
			}
			if currAlpha < 0.0 {
				currAlpha += 360.0
			}
			if currBeta >= 360.0 {
				currBeta -= 360.0
			}
			if currBeta < 0.0 {
				currBeta += 360.0
			}
		}
		prefix = true
	}

	// Read input_file (use "world.jpg" 43200 x 21600 pixels for best results)
	reader, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()
	m, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	bounds := m.Bounds()
	x := bounds.Max.X
	y := bounds.Max.Y
	// Width
	width, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	if width <= 0 || width > x {
		return fmt.Errorf("width must be from (0, %d] range", x)
	}
	// Height
	height, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	if height <= 0 || height > y {
		return fmt.Errorf("height must be from (0, %d] range", y)
	}

	// Output image(s)
	var target []*image.RGBA
	if dbl {
		width2 := width * 2
		widthOff1 := width / 2
		widthOff2 := width + widthOff1
		for f := 0; f < steps; f++ {
			target = append(target, image.NewRGBA(image.Rect(0, 0, width2, height)))
		}
		for i := 0; i < width; i++ {
			if i%50 == 0 {
				fmt.Printf("Generating: %f%%\n", float64(i)/float64(width)*100.0)
			}
			iOff1 := i + widthOff1
			iOff2 := i + widthOff2
			if iOff2 > width2 {
				iOff2 -= width2
			}
			lon := -180.0 + float64(i)*360.0/float64(width)
			for j := 0; j < height; j++ {
				lat := -90.0 + float64(j)*180.0/float64(height)
				for f := 0; f < steps; f++ {
					nlon, nlat := poleShift(lon, lat, alphas[f], betas[f])
					ix := int((nlon + 180.0) / 360.0 * float64(x))
					iy := int((nlat + 90.0) / 180.0 * float64(y))
					target[f].Set(iOff1, j, m.At(ix, iy))
					target[f].Set(iOff2, j, m.At(ix, iy))
				}
			}
		}
	} else {
		for f := 0; f < steps; f++ {
			target = append(target, image.NewRGBA(image.Rect(0, 0, width, height)))
		}
		for i := 0; i < width; i++ {
			if i%50 == 0 {
				fmt.Printf("Generating: %f%%\n", float64(i)/float64(width)*100.0)
			}
			lon := -180.0 + float64(i)*360.0/float64(width)
			for j := 0; j < height; j++ {
				lat := -90.0 + float64(j)*180.0/float64(height)
				for f := 0; f < steps; f++ {
					nlon, nlat := poleShift(lon, lat, alphas[f], betas[f])
					ix := int((nlon + 180.0) / 360.0 * float64(x))
					iy := int((nlat + 90.0) / 180.0 * float64(y))
					target[f].Set(i, j, m.At(ix, iy))
				}
			}
		}
	}

	// Save output JPEGs, PNGs or GIFs
	for f := 0; f < steps; f++ {
		var fn string
		if prefix {
			ary := strings.Split(args[5], "/")
			lAry := len(ary)
			last := ary[lAry-1]
			ary[lAry-1] = fmt.Sprintf("%09d_%s", f+1, last)
			fn = strings.Join(ary, "/")
		} else {
			fn = args[5]
		}
		fi, err := os.Create(fn)
		if err != nil {
			return err
		}
		lfn := strings.ToLower(args[5])
		var ierr error
		if strings.Contains(lfn, ".png") {
			ierr = png.Encode(fi, target[f])
		} else if strings.Contains(lfn, ".jpg") || strings.Contains(lfn, ".jpeg") {
			ierr = jpeg.Encode(fi, target[f], nil)
		} else if strings.Contains(lfn, ".gif") {
			ierr = gif.Encode(fi, target[f], nil)
		}
		if ierr != nil {
			return ierr
		}
		err = fi.Close()
		if err != nil {
			return err
		}
		fmt.Printf("Saved: %s\n", fn)
	}
	return nil
}

func main() {
	if len(os.Args) == 5 {
		fmt.Printf("Looking from up to the North Pole\n")
		fmt.Printf("(lon 0, lat 0) is to the right: (1, 0, 0)\n")
		fmt.Printf("(lon 90, lat 0) is to the front: (0, 1, 0)\n")
		fmt.Printf("(lon +/-180, lat 0) is to the left: (-1, 0, 0)\n")
		fmt.Printf("(lon -90, lat 0) is to the back: (0, -1, 0)\n")
		fmt.Printf("North Pole (lon ?, lat 90) is up: (0, 0, 1)\n")
		fmt.Printf("South Pole (lon ?, lat -90) is down: (0, 0, -1)\n")
		fmt.Printf("First rotation is clockwise by using Y axis which goes through (90, 0) and (-90, 0)\n")
		fmt.Printf("Second rotation is clockwise by using Z axis which goes through unrotated North and South poles\n")
		err := poleShiftArgs()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else if len(os.Args) == 7 {
		err := mapWithPolesShifted(os.Args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Provide either 4 args: lon lat alpha beta or 6 args: input_file alpha beta output_width output_height output_name\n")
	}
}
