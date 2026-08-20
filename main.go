package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

const (
	pageChars = 3000  // characters of text per page
	maxPad    = 16384 // max filler bytes in the Info dictionary
)

func randomText(rng *rand.Rand, n int) string {
	var sb strings.Builder
	sb.Grow(n + 16)
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for sb.Len() < n {
		wl := rng.Intn(8) + 3
		for i := 0; i < wl; i++ {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
		sb.WriteByte(' ')
	}
	return sb.String()
}

// build creates a document with `pages` text pages and `pad` bytes of
// filler in the Info dictionary, so the final size can be hit exactly.
func build(seed int64, pages int, pad int) *gofpdf.Fpdf {
	rng := rand.New(rand.NewSource(seed))
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetCompression(false)
	pdf.SetFont("Arial", "", 12)
	if pad > 0 {
		pdf.SetKeywords(strings.Repeat("x", pad), false)
	}
	for i := 0; i < pages; i++ {
		pdf.AddPage()
		pdf.MultiCell(0, 5, randomText(rng, pageChars), "", "L", false)
	}
	return pdf
}

type countWriter struct{ n int64 }

func (c *countWriter) Write(p []byte) (int, error) { c.n += int64(len(p)); return len(p), nil }

func size(seed int64, pages, pad int) int64 {
	var c countWriter
	if err := build(seed, pages, pad).Output(&c); err != nil {
		log.Fatalf("render failed: %v", err)
	}
	return c.n
}

func main() {
	sizeMB := flag.Int("size", 5, "target size in MB")
	outFile := flag.String("out", "out.pdf", "output path")
	seed := flag.Int64("seed", 0, "random seed (0 = time-based)")
	flag.Parse()
	if *sizeMB <= 0 {
		log.Fatalf("invalid -size: %d", *sizeMB)
	}
	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	target := int64(*sizeMB) * 1024 * 1024
	start := time.Now()
	log.Printf("target %d bytes, seed %d", target, *seed)

	// Probe per-page cost, then estimate the page count.
	const probe = 20
	base := size(*seed, 0, 0)
	perPage := float64(size(*seed, probe, 0)-base) / probe
	pages := int(float64(target-base) / perPage)
	if pages < 1 {
		pages = 1
	}

	// Converge: adjust page count proportionally while over target, then pad the remainder.
	pad := 0
	for i := 0; i < 10; i++ {
		got := size(*seed, pages, pad)
		diff := target - got
		log.Printf("pass %d: %d pages, %d pad -> %d bytes (diff %d)", i+1, pages, pad, got, diff)
		if diff == 0 {
			break
		}
		if diff < 0 && int64(pad) < -diff {
			over := float64(-diff - int64(pad))
			pages -= int(over/perPage) + 1
			pad = 0
			continue
		}
		pad += int(diff)
		if pad > maxPad { // keep the filler string well under the 32767-byte PDF string limit
			pages += pad / int(perPage)
			pad = 0
		}
	}

	f, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("create failed: %v", err)
	}
	w := bufio.NewWriterSize(f, 1<<20)
	if err := build(*seed, pages, pad).Output(w); err != nil {
		log.Fatalf("render failed: %v", err)
	}
	if err := w.Flush(); err != nil {
		log.Fatalf("write failed: %v", err)
	}
	if err := f.Close(); err != nil {
		log.Fatalf("close failed: %v", err)
	}
	st, _ := os.Stat(*outFile)
	fmt.Printf("wrote %s: %d bytes (%d pages, %d pad) in %.1fs\n", *outFile, st.Size(), pages, pad, time.Since(start).Seconds())
	if st.Size() != target {
		log.Printf("WARNING: size %d differs from target %d", st.Size(), target)
	}
}
