package renderer

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/gowtham2003/gotable/pkg/parser"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type ImageRenderer struct {
	padding    int
	cellHeight int
	fontSize   int
}

func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{
		padding:    10,
		cellHeight: 30,
		fontSize:   12,
	}
}

func (r *ImageRenderer) Render(data *parser.TableData) (string, error) {
	// Calculate dimensions
	widths := getColumnWidths(data)

	// Add padding to widths
	totalWidth := 1 // Start with 1 for left border
	for _, width := range widths {
		totalWidth += width*7 + (r.padding * 2) + 1 // Multiply by 7 for approximate character width
	}

	totalHeight := (len(data.Rows)+1)*r.cellHeight + 1

	// Create new image
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Fill background
	for x := 0; x < totalWidth; x++ {
		for y := 0; y < totalHeight; y++ {
			img.Set(x, y, color.White)
		}
	}

	// Create drawer for text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
	}

	// Draw horizontal lines
	for y := 0; y <= len(data.Rows)+1; y++ {
		yPos := y * r.cellHeight
		drawHorizontalLine(img, 0, totalWidth, yPos)
	}

	// Draw headers
	currentX := 1
	for _, header := range data.Headers {
		width := widths[header]*7 + (r.padding * 2)

		// Draw vertical line
		drawVerticalLine(img, currentX-1, 0, totalHeight)

		// Draw text
		d.Dot = fixed.Point26_6{
			X: fixed.Int26_6(currentX+r.padding) << 6,
			Y: fixed.Int26_6(r.cellHeight-r.padding) << 6,
		}
		d.DrawString(header)

		currentX += width
	}
	// Draw final vertical line
	drawVerticalLine(img, currentX, 0, totalHeight)

	// Draw data rows
	for rowIdx, row := range data.Rows {
		currentX = 1
		y := (rowIdx + 1) * r.cellHeight

		for _, header := range data.Headers {
			width := widths[header]*7 + (r.padding * 2)

			d.Dot = fixed.Point26_6{
				X: fixed.Int26_6(currentX+r.padding) << 6,
				Y: fixed.Int26_6(y+r.cellHeight-r.padding) << 6,
			}
			d.DrawString(row[header])

			currentX += width
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func drawHorizontalLine(img *image.RGBA, x1, x2, y int) {
	for x := x1; x < x2; x++ {
		img.Set(x, y, color.Black)
	}
}

func drawVerticalLine(img *image.RGBA, x, y1, y2 int) {
	for y := y1; y < y2; y++ {
		img.Set(x, y, color.Black)
	}
}
