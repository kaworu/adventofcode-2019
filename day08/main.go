package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	// Black is a dark Pixel
	Black = iota
	// White is a bright Pixel
	White
	// Trans is a transparent Pixel
	Trans
)

// Pixel are the basic elements of a Layer.
type Pixel uint8

// Layer of the Space Image Format.
type Layer struct {
	width, height int
	pixels        []Pixel
}

// NewLayer return an image Layer of the provided dimensions with the given
// pixels data. An error is returned when the dimensions doesn't match the
// count of pixels elements.
func NewLayer(width int, height int, pixels []Pixel) (*Layer, error) {
	if width < 0 || height < 0 || len(pixels) != width*height {
		return nil, fmt.Errorf("wrong buffer dimensions")
	}
	return &Layer{width, height, pixels}, nil
}

// Flatten take a stack of layer and combine them from the top (first) to the
// bottom (last). It return the combined Layer.
func Flatten(layers []*Layer) *Layer {
	if len(layers) == 0 {
		return nil
	}
	width, height := layers[0].width, layers[0].height
	pixels := make([]Pixel, width*height)
	flat := &Layer{width, height, pixels}
	// Black being the zero-value for Pixel, we start with an all-black Layer
	// here. Since we never check for Pixel values in flat it is not an issue.
	for i := range flat.pixels {
	layersLoop:
		for _, l := range layers {
			switch p := l.pixels[i]; p {
			case Black, White:
				flat.pixels[i] = p
				break layersLoop
			}
		}
	}
	return flat
}

// MinLayerBy returns the first Layer l that evaluate to the minimum value
// under f across all layers.  When more than one layer evaluate to the same
// minimum value, the first one (in all's order) is returned. It is guaranteed
// that f is called once and only once for every Layer in l.
func MinLayerBy(all []*Layer, f func(*Layer) int) *Layer {
	var min *Layer
	var fmin int
	for _, l := range all {
		fl := f(l)
		if min == nil || fl < fmin {
			min, fmin = l, fl
		}
	}
	return min
}

// Count returns the number of pixels matching ref in the Layer.
func (l Layer) Count(ref Pixel) int {
	n := 0
	for _, p := range l.pixels {
		if p == ref {
			n++
		}
	}
	return n
}

// String implement fmt Stringer interface for Layer.
func (l Layer) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			switch p := l.pixels[y*l.width+x]; p {
			// We use UTF-8 'BLACK' for white and conversely assuming output to
			// a something-on-black terminal because that is what I use.
			case Black:
				buf.WriteString("⬜") // 'WHITE LARGE SQUARE' (U+2B1C)
			case White:
				buf.WriteString("⬛") // 'BLACK LARGE SQUARE' (U+2B1B)
			case Trans:
				buf.WriteString("  ")
			default:
				buf.WriteString("??")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// main parse the puzzle provided on stdin and then compute the number of 1
// digits multiplied by the number of 2 digits from the layer having the fewest
// 0 digits.
func main() {
	width, height := 25, 6
	layers, err := Parse(width, height, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %s\n", err)
		os.Exit(1)
	}
	l := MinLayerBy(layers, func(l *Layer) int {
		return l.Count(0)
	})
	one, two := l.Count(1), l.Count(2)
	fmt.Printf("The number of 1 digits multiplied by the number of 2 digits is %v * %v = %v.\n", one, two, one*two)
	flat := Flatten(layers)
	fmt.Printf("The message after decoding the image is:\n\n%v\n", flat)
}

// Parse the Space Image Format into its pixels layers.
// It returns the Layer stack and any read or parsing error encountered.
func Parse(width int, height int, r io.Reader) ([]*Layer, error) {
	// Scan all the pixels
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)
	var pixels []Pixel
	for scanner.Scan() {
		b := scanner.Text()
		if b == "\n" {
			continue // ignore newlines
		}
		i, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		pixels = append(pixels, Pixel(i))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// Create as many layers as needed, slicing the pixels.
	lpc := width * height    // layer pixel count
	tlc := len(pixels) / lpc // total layer count
	layers := make([]*Layer, tlc)
	for i := 0; i < tlc; i++ {
		chunk := pixels[i*lpc : (i+1)*lpc]
		l, err := NewLayer(width, height, chunk)
		if err != nil {
			return nil, err
		}
		layers[i] = l
	}
	return layers, nil
}
