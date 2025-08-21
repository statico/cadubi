package commands

import (
	"fmt"
	"math"
)

// LineCommand draws a line between two points using Bresenham's algorithm
type LineCommand struct {
	BaseCommand
	X1, Y1    int
	X2, Y2    int
	Pixel     Pixel
	Layer     int
	OldPixels []PixelPosition // Store pixels that were overwritten
	executed  bool
}

// PixelPosition represents a pixel and its position
type PixelPosition struct {
	X, Y  int
	Pixel Pixel
}

// NewLineCommand creates a new line command
func NewLineCommand(x1, y1, x2, y2 int, pixel Pixel, layer int) *LineCommand {
	return &LineCommand{
		BaseCommand: NewBaseCommand(fmt.Sprintf("Draw line from (%d,%d) to (%d,%d)", x1, y1, x2, y2)),
		X1:          x1,
		Y1:          y1,
		X2:          x2,
		Y2:          y2,
		Pixel:       pixel,
		Layer:       layer,
		OldPixels:   make([]PixelPosition, 0),
		executed:    false,
	}
}

// Execute draws the line
func (lc *LineCommand) Execute(ctx *CommandContext) error {
	if !lc.executed {
		// Store old pixels for undo before drawing
		points := lc.getLinePoints()
		lc.OldPixels = make([]PixelPosition, len(points))

		for i, point := range points {
			oldPixel, err := ctx.Canvas.GetPixel(point.X, point.Y)
			if err != nil {
				oldPixel = Pixel{Char: ' '}
			}
			lc.OldPixels[i] = PixelPosition{
				X:     point.X,
				Y:     point.Y,
				Pixel: oldPixel,
			}
		}

		lc.executed = true
	}

	// Draw the line
	points := lc.getLinePoints()
	for _, point := range points {
		if err := ctx.Canvas.SetPixel(point.X, point.Y, lc.Pixel); err != nil {
			return err
		}
	}

	return nil
}

// Undo restores the old pixels
func (lc *LineCommand) Undo(ctx *CommandContext) error {
	if !lc.executed {
		return fmt.Errorf("cannot undo unexecuted command")
	}

	// Restore old pixels
	for _, oldPixel := range lc.OldPixels {
		if err := ctx.Canvas.SetPixel(oldPixel.X, oldPixel.Y, oldPixel.Pixel); err != nil {
			return err
		}
	}

	return nil
}

// getLinePoints calculates all points on the line using Bresenham's algorithm
func (lc *LineCommand) getLinePoints() []struct{ X, Y int } {
	points := make([]struct{ X, Y int }, 0)

	x1, y1, x2, y2 := lc.X1, lc.Y1, lc.X2, lc.Y2

	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))

	sx := 1
	if x1 > x2 {
		sx = -1
	}

	sy := 1
	if y1 > y2 {
		sy = -1
	}

	err := dx - dy
	x, y := x1, y1

	for {
		points = append(points, struct{ X, Y int }{x, y})

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x += sx
		}

		if e2 < dx {
			err += dx
			y += sy
		}
	}

	return points
}

// RectangleCommand draws a rectangle
type RectangleCommand struct {
	BaseCommand
	X1, Y1    int
	X2, Y2    int
	Pixel     Pixel
	Layer     int
	Filled    bool
	OldPixels []PixelPosition
	executed  bool
}

// NewRectangleCommand creates a new rectangle command
func NewRectangleCommand(x1, y1, x2, y2 int, pixel Pixel, layer int, filled bool) *RectangleCommand {
	fillType := "outline"
	if filled {
		fillType = "filled"
	}

	return &RectangleCommand{
		BaseCommand: NewBaseCommand(fmt.Sprintf("Draw %s rectangle from (%d,%d) to (%d,%d)", fillType, x1, y1, x2, y2)),
		X1:          x1,
		Y1:          y1,
		X2:          x2,
		Y2:          y2,
		Pixel:       pixel,
		Layer:       layer,
		Filled:      filled,
		OldPixels:   make([]PixelPosition, 0),
		executed:    false,
	}
}

// Execute draws the rectangle
func (rc *RectangleCommand) Execute(ctx *CommandContext) error {
	if !rc.executed {
		// Store old pixels for undo
		points := rc.getRectanglePoints()
		rc.OldPixels = make([]PixelPosition, len(points))

		for i, point := range points {
			oldPixel, err := ctx.Canvas.GetPixel(point.X, point.Y)
			if err != nil {
				oldPixel = Pixel{Char: ' '}
			}
			rc.OldPixels[i] = PixelPosition{
				X:     point.X,
				Y:     point.Y,
				Pixel: oldPixel,
			}
		}

		rc.executed = true
	}

	// Draw the rectangle
	points := rc.getRectanglePoints()
	for _, point := range points {
		if err := ctx.Canvas.SetPixel(point.X, point.Y, rc.Pixel); err != nil {
			return err
		}
	}

	return nil
}

// Undo restores the old pixels
func (rc *RectangleCommand) Undo(ctx *CommandContext) error {
	if !rc.executed {
		return fmt.Errorf("cannot undo unexecuted command")
	}

	// Restore old pixels
	for _, oldPixel := range rc.OldPixels {
		if err := ctx.Canvas.SetPixel(oldPixel.X, oldPixel.Y, oldPixel.Pixel); err != nil {
			return err
		}
	}

	return nil
}

// getRectanglePoints calculates all points for the rectangle
func (rc *RectangleCommand) getRectanglePoints() []struct{ X, Y int } {
	points := make([]struct{ X, Y int }, 0)

	// Normalize coordinates
	x1, y1 := rc.X1, rc.Y1
	x2, y2 := rc.X2, rc.Y2

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	if rc.Filled {
		// Filled rectangle
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				points = append(points, struct{ X, Y int }{x, y})
			}
		}
	} else {
		// Rectangle outline
		// Top and bottom lines
		for x := x1; x <= x2; x++ {
			points = append(points, struct{ X, Y int }{x, y1})
			if y2 != y1 {
				points = append(points, struct{ X, Y int }{x, y2})
			}
		}

		// Left and right lines (excluding corners to avoid duplication)
		for y := y1 + 1; y < y2; y++ {
			points = append(points, struct{ X, Y int }{x1, y})
			if x2 != x1 {
				points = append(points, struct{ X, Y int }{x2, y})
			}
		}
	}

	return points
}
