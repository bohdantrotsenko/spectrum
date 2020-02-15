package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

// NN is the number of frequncies in the array.
const NN = 456

var width, height, hrepeat int
var framesPerSeconds int
var framesPrintingThreshold int
var p6header []byte
var printFrequencies = flag.Bool("freq", false, "prints assumed frequencies")
var printSinglePicture = flag.Bool("pic", false, "print a single picture")

func init() {
	flag.IntVar(&width, "width", 1920, "width of the picture")
	flag.IntVar(&height, "maxheight", 1080, "maximal height of the picture")
	flag.IntVar(&framesPerSeconds, "fps", 30, "frames per second in the resulting video")
}

func writeP6(buf []byte, off int, rightBlue bool, dst io.Writer) {
	dst.Write(p6header)
	stride := make([]byte, 0, width*3)
	for y := NN - 1; y >= 0; y-- {

		for x := 0; x < width; x++ {
			el := buf[(x+off)%width+y*width]
			if rightBlue && x*2 > width {
				stride = append(stride, 0, 0, el)
			} else {
				stride = append(stride, el, el, el)
			}
		}
		for i := 0; i < hrepeat; i++ {
			dst.Write(stride)
		}
		stride = stride[:0]
	}
}

func run(src io.Reader, dst io.Writer) error {
	slice := make([]byte, NN)

	buf := make([]byte, width*height)
	xOffset, zeroSlices := 0, 0
	if !*printSinglePicture {
		xOffset = width / 2
		zeroSlices = width / 2
	}

	printingStarted := false
	framesCounter := 0

	for {
		chunk, err := readFull(src, slice)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read input: %s", err)
		}

		if chunk < NN {
			if zeroSlices == 0 {
				break
			}
			zeroSlices--
			for i := 0; i < 456; i++ {
				slice[i] = 0
			}
		}

		for y, el := range slice {
			buf[xOffset+y*width] = el
		}

		xOffset++

		if xOffset == width && *printSinglePicture {
			writeP6(buf, 0, false, dst)
			return nil
		}

		if !printingStarted {
			if xOffset == width {
				writeP6(buf, 0, true, dst)
				printingStarted = true
			}
		} else { // printingStarted = true
			framesCounter++
		}

		xOffset = xOffset % width

		if framesCounter >= framesPrintingThreshold {
			writeP6(buf, xOffset, true, dst)
			framesCounter -= framesPrintingThreshold
		}
	}

	if *printSinglePicture {
		writeP6(buf, 0, false, dst)
	}

	return nil
}

func main() {
	flag.Parse()

	if *printFrequencies {
		// 24.499714748859326 to 17485.357993994654
		for i := -200; i < 256; i++ {
			f := 440.0 * math.Pow(2, float64(i)/48.0)
			fmt.Printf("%.6f, ", f)
		}
		fmt.Println()
		return
	}

	height -= height % NN
	hrepeat = height / NN
	p6header = []byte(fmt.Sprintf("P6\n%d %d\n255\n", width, height))

	framesPrintingThreshold = 96000 / 40 / framesPerSeconds

	dst := bufio.NewWriter(os.Stdout)
	defer dst.Flush()
	if err := run(bufio.NewReader(os.Stdin), dst); err != nil {
		log.Fatal(err)
	}
}

func readFull(src io.Reader, target []byte) (int, error) {
	off := 0
	for {
		n, err := src.Read(target[off:])
		off += n
		if err != nil {
			return off, err
		}
		if off == len(target) {
			return off, nil
		}
	}
}
