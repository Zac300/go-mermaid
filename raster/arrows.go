package raster

import (
	"image"
	"image/color"
	"math"
	"regexp"
	"strconv"
)

// oksvg does not render SVG markers, so arrowheads on marker-end edges are
// missing. We synthesize them: find marker-end <path>/<line> elements, take
// the last segment, and fill a triangle at the tip.

var (
	markerElemRe = regexp.MustCompile(`(?s)<(path|line)\b([^>]*marker-end="url\([^"]*"[^>]*)>`)
	numRe        = regexp.MustCompile(`-?[0-9]+(?:\.[0-9]+)?`)
)

// drawArrows draws arrowheads for marker-end edges onto img.
func drawArrows(img *image.RGBA, svg string, scale float64) {
	gx, gy := 0.0, 0.0
	if loc := groupRe.FindStringSubmatch(svg); loc != nil {
		gx = numAttr(loc[1], "0")
		gy = numAttr(loc[2], "0")
	}
	for _, m := range markerElemRe.FindAllStringSubmatch(svg, -1) {
		attrs := parseAttrs(m[2])
		col := colorOf(attrs["stroke"])
		var tx, ty, px, py float64
		var ok bool
		if m[1] == "line" {
			tx, ty = numAttr(attrs["x2"], "0"), numAttr(attrs["y2"], "0")
			px, py = numAttr(attrs["x1"], "0"), numAttr(attrs["y1"], "0")
			ok = true
		} else {
			tx, ty, px, py, ok = lastSegment(attrs["d"])
		}
		if !ok {
			continue
		}
		arrowhead(img, (tx+gx)*scale, (ty+gy)*scale, (px+gx)*scale, (py+gy)*scale, scale, col)
	}
}

// lastSegment returns the final point and the one before it from a path's d.
func lastSegment(d string) (tx, ty, px, py float64, ok bool) {
	nums := numRe.FindAllString(d, -1)
	if len(nums) < 4 {
		return 0, 0, 0, 0, false
	}
	tx, _ = strconv.ParseFloat(nums[len(nums)-2], 64)
	ty, _ = strconv.ParseFloat(nums[len(nums)-1], 64)
	px, _ = strconv.ParseFloat(nums[len(nums)-4], 64)
	py, _ = strconv.ParseFloat(nums[len(nums)-3], 64)
	return tx, ty, px, py, true
}

// arrowhead fills a triangle at (tx,ty) pointing away from (px,py).
func arrowhead(img *image.RGBA, tx, ty, px, py, scale float64, col color.Color) {
	dx, dy := tx-px, ty-py
	d := math.Hypot(dx, dy)
	if d == 0 {
		return
	}
	dx, dy = dx/d, dy/d
	l, half := 9*scale, 4.5*scale
	bx, by := tx-dx*l, ty-dy*l
	perpx, perpy := -dy, dx
	fillTriangle(img,
		tx, ty,
		bx+perpx*half, by+perpy*half,
		bx-perpx*half, by-perpy*half,
		col)
}

// fillTriangle rasterizes a solid triangle with a simple scanline fill.
func fillTriangle(img *image.RGBA, x0, y0, x1, y1, x2, y2 float64, col color.Color) {
	minX := int(math.Floor(min3(x0, x1, x2)))
	maxX := int(math.Ceil(max3(x0, x1, x2)))
	minY := int(math.Floor(min3(y0, y1, y2)))
	maxY := int(math.Ceil(max3(y0, y1, y2)))
	area := edge(x0, y0, x1, y1, x2, y2)
	if area == 0 {
		return
	}
	b := img.Bounds()
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			w0 := edge(x1, y1, x2, y2, fx, fy)
			w1 := edge(x2, y2, x0, y0, fx, fy)
			w2 := edge(x0, y0, x1, y1, fx, fy)
			if (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0) {
				if image.Pt(x, y).In(b) {
					img.Set(x, y, col)
				}
			}
		}
	}
}

func edge(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

func min3(a, b, c float64) float64 { return math.Min(a, math.Min(b, c)) }
func max3(a, b, c float64) float64 { return math.Max(a, math.Max(b, c)) }
