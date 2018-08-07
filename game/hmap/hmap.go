// Command go-heightmap provides a basic Square Diamond heightmap / terrain generator.
package hmap

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	optSize := flag.Int("size", 1024, "Image size, must be power of two.")
	optSamples := flag.Int("samples", 128, "Number of samples. Lower is more rough terrain. Must be power of two.")
	optBlur := flag.Int("blur", 2, "Blur/Smooth size. Somewhere between 1 and 5.")
	optScale := flag.Float64("scale", 1.0, "Scale. Most lileky 1.0")
	optOut := flag.String("o", "", "Output png filename.")
	flag.Parse()

	if *optOut == "" {
		flag.Usage()
		os.Exit(1)
	}

	h := NewHeightmap(*optSize)
	h.generate(*optSamples, *optScale)
	h.normalize()
	h.blur(*optBlur)
	h.png(*optOut)
}

// Heightmap generates a heightmap based on the Square Diamond algorithm.
type Heightmap struct {
	random        *rand.Rand
	points        []float64
	width, height int
}

type generator interface {
	Generate()
}

// NewHeightmap initializes a new Heightmap using the specified size.
func NewHeightmap(size int) *Heightmap {
	h := &Heightmap{}
	h.random = rand.New(rand.NewSource(time.Now().Unix()))
	h.points = make([]float64, size*size)
	h.width = size
	h.height = size

	// init/randomize
	for x := 0; x < h.width; x++ {
		for y := 0; y < h.height; y++ {
			h.set(x, y, h.frand())
		}
	}
	return h
}

func (h *Heightmap) png(fname string) {
	rect := image.Rect(0, 0, h.width, h.height)
	img := image.NewGray16(rect)

	for x := 0; x < h.width; x++ {
		for y := 0; y < h.height; y++ {
			val := h.get(x, y)
			col := color.Gray16{uint16(val * 0xffff)}
			img.Set(x, y, col)
		}
	}

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func (h *Heightmap) normalize() {
	var min = 1.0
	var max = 0.0

	for i := 0; i < h.width*h.height; i++ {
		if h.points[i] < min {
			min = h.points[i]
		}
		if h.points[i] > max {
			max = h.points[i]
		}
	}
	rat := max - min
	for i := 0; i < h.width*h.height; i++ {
		h.points[i] = (h.points[i] - min) / rat
	}
}

func (h *Heightmap) blur(size int) {
	for x := 0; x < h.width; x++ {
		for y := 0; y < h.height; y++ {
			count := 0
			total := 0.0

			for x0 := x - size; x0 <= x+size; x0++ {
				for y0 := y - size; y0 <= y+size; y0++ {
					total += h.get(x0, y0)
					count++
				}
			}
			if count > 0 {
				h.set(x, y, total/float64(count))
			}
		}
	}
}

func (h *Heightmap) frand() float64 {
	return (h.random.Float64() * 2.0) - 1.0
}

func (h *Heightmap) get(x, y int) float64 {
	return h.points[(x&(h.width-1))+((y&(h.height-1))*h.width)]
}

func (h *Heightmap) set(x, y int, val float64) {
	h.points[(x&(h.width-1))+((y&(h.height-1))*h.width)] = val
}

func (h *Heightmap) generate(samples int, scale float64) {
	for samples > 0 {
		h.squarediamond(samples, scale)
		samples /= 2
		scale /= 2.0
	}
}

func (h *Heightmap) squarediamond(step int, scale float64) {
	half := step / 2
	for y := half; y < h.height+half; y += step {
		for x := half; x < h.width+half; x += step {
			h.square(x, y, step, h.frand()*scale)
		}
	}
	for y := 0; y < h.height; y += step {
		for x := 0; x < h.width; x += step {
			h.diamond(x+half, y, step, h.frand()*scale)
			h.diamond(x, y+half, step, h.frand()*scale)
		}
	}
}

func (h *Heightmap) square(x, y, size int, val float64) {
	half := size / 2
	a := h.get(x-half, y-half)
	b := h.get(x+half, y-half)
	c := h.get(x-half, y+half)
	d := h.get(x+half, y+half)
	h.set(x, y, ((a+b+c+d)/4.0)+val)
}

func (h *Heightmap) diamond(x, y, size int, val float64) {
	half := size / 2
	a := h.get(x-half, y)
	b := h.get(x+half, y)
	c := h.get(x, y-half)
	d := h.get(x, y+half)
	h.set(x, y, ((a+b+c+d)/4.0)+val)
}
