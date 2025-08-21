package commands

import "fmt"

// DrawCommand draws a pixel at a specific location
type DrawCommand struct {
	BaseCommand
	X, Y     int
	NewPixel Pixel
	OldPixel Pixel
	Layer    int
	executed bool
}

// NewDrawCommand creates a new draw command
func NewDrawCommand(x, y int, pixel Pixel, layer int) *DrawCommand {
	return &DrawCommand{
		BaseCommand: NewBaseCommand(fmt.Sprintf("Draw '%c' at (%d, %d)", pixel.Char, x, y)),
		X:           x,
		Y:           y,
		NewPixel:    pixel,
		Layer:       layer,
		executed:    false,
	}
}

// Execute draws the pixel
func (dc *DrawCommand) Execute(ctx *CommandContext) error {
	// Store the old pixel for undo
	if !dc.executed {
		oldPixel, err := ctx.Canvas.GetPixel(dc.X, dc.Y)
		if err != nil {
			// If pixel doesn't exist, create an empty one
			dc.OldPixel = Pixel{Char: ' '}
		} else {
			dc.OldPixel = oldPixel
		}
		dc.executed = true
	}

	// Set the new pixel
	return ctx.Canvas.SetPixel(dc.X, dc.Y, dc.NewPixel)
}

// Undo restores the old pixel
func (dc *DrawCommand) Undo(ctx *CommandContext) error {
	if !dc.executed {
		return fmt.Errorf("cannot undo unexecuted command")
	}

	return ctx.Canvas.SetPixel(dc.X, dc.Y, dc.OldPixel)
}

// EraseCommand erases a pixel at a specific location
type EraseCommand struct {
	BaseCommand
	X, Y     int
	OldPixel Pixel
	Layer    int
	executed bool
}

// NewEraseCommand creates a new erase command
func NewEraseCommand(x, y int, layer int) *EraseCommand {
	return &EraseCommand{
		BaseCommand: NewBaseCommand(fmt.Sprintf("Erase at (%d, %d)", x, y)),
		X:           x,
		Y:           y,
		Layer:       layer,
		executed:    false,
	}
}

// Execute erases the pixel
func (ec *EraseCommand) Execute(ctx *CommandContext) error {
	// Store the old pixel for undo
	if !ec.executed {
		oldPixel, err := ctx.Canvas.GetPixel(ec.X, ec.Y)
		if err != nil {
			return err
		}
		ec.OldPixel = oldPixel
		ec.executed = true
	}

	// Erase by setting to empty pixel
	emptyPixel := Pixel{Char: ' '}
	return ctx.Canvas.SetPixel(ec.X, ec.Y, emptyPixel)
}

// Undo restores the erased pixel
func (ec *EraseCommand) Undo(ctx *CommandContext) error {
	if !ec.executed {
		return fmt.Errorf("cannot undo unexecuted command")
	}

	return ctx.Canvas.SetPixel(ec.X, ec.Y, ec.OldPixel)
}

// MoveCursorCommand moves the cursor to a new position
type MoveCursorCommand struct {
	BaseCommand
	NewX, NewY int
	OldX, OldY int
	executed   bool
}

// NewMoveCursorCommand creates a new move cursor command
func NewMoveCursorCommand(x, y int) *MoveCursorCommand {
	return &MoveCursorCommand{
		BaseCommand: NewBaseCommand(fmt.Sprintf("Move cursor to (%d, %d)", x, y)),
		NewX:        x,
		NewY:        y,
		executed:    false,
	}
}

// Execute moves the cursor
func (mc *MoveCursorCommand) Execute(ctx *CommandContext) error {
	if !mc.executed {
		mc.OldX = ctx.Cursor.X
		mc.OldY = ctx.Cursor.Y
		mc.executed = true
	}

	ctx.Cursor.X = mc.NewX
	ctx.Cursor.Y = mc.NewY
	return nil
}

// Undo restores the cursor position
func (mc *MoveCursorCommand) Undo(ctx *CommandContext) error {
	if !mc.executed {
		return fmt.Errorf("cannot undo unexecuted command")
	}

	ctx.Cursor.X = mc.OldX
	ctx.Cursor.Y = mc.OldY
	return nil
}
