package commands

import (
	"encoding/json"
	"fmt"
	"testing"
)

// JSONCommand represents a command in JSON format for testing
type JSONCommand struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// JSONTestScenario represents a complete test scenario
type JSONTestScenario struct {
	Name      string        `json:"name"`
	Commands  []JSONCommand `json:"commands"`
	Assertions []Assertion  `json:"assertions"`
}

// Assertion represents a test assertion
type Assertion struct {
	Type     string                 `json:"type"`     // "pixel", "cursor", "history_count"
	Expected map[string]interface{} `json:"expected"`
}

// ExecuteJSONCommands executes a sequence of commands from JSON
func ExecuteJSONCommands(ctx *CommandContext, history *CommandHistory, jsonCommands []JSONCommand) error {
	for _, jsonCmd := range jsonCommands {
		cmd, err := CreateCommandFromJSON(jsonCmd)
		if err != nil {
			return err
		}
		
		if err := history.Execute(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

// CreateCommandFromJSON creates a Command from JSON representation
func CreateCommandFromJSON(jsonCmd JSONCommand) (Command, error) {
	switch jsonCmd.Type {
	case "draw":
		x := int(jsonCmd.Params["x"].(float64))
		y := int(jsonCmd.Params["y"].(float64))
		char := rune(jsonCmd.Params["char"].(string)[0])
		
		pixel := Pixel{Char: char}
		if fg, exists := jsonCmd.Params["foreground"]; exists {
			pixel.Foreground = fg.(string)
		}
		if bg, exists := jsonCmd.Params["background"]; exists {
			pixel.Background = bg.(string)
		}
		if bold, exists := jsonCmd.Params["bold"]; exists {
			pixel.Bold = bold.(bool)
		}
		
		layer := 0
		if l, exists := jsonCmd.Params["layer"]; exists {
			layer = int(l.(float64))
		}
		
		return NewDrawCommand(x, y, pixel, layer), nil
		
	case "erase":
		x := int(jsonCmd.Params["x"].(float64))
		y := int(jsonCmd.Params["y"].(float64))
		layer := 0
		if l, exists := jsonCmd.Params["layer"]; exists {
			layer = int(l.(float64))
		}
		return NewEraseCommand(x, y, layer), nil
		
	case "move_cursor":
		x := int(jsonCmd.Params["x"].(float64))
		y := int(jsonCmd.Params["y"].(float64))
		return NewMoveCursorCommand(x, y), nil
		
	default:
		return nil, fmt.Errorf("unknown command type: %s", jsonCmd.Type)
	}
}

// ValidateAssertions validates test assertions
func ValidateAssertions(t *testing.T, ctx *CommandContext, history *CommandHistory, assertions []Assertion) {
	for _, assertion := range assertions {
		switch assertion.Type {
		case "pixel":
			x := int(assertion.Expected["x"].(float64))
			y := int(assertion.Expected["y"].(float64))
			expectedChar := rune(assertion.Expected["char"].(string)[0])
			
			pixel, err := ctx.Canvas.GetPixel(x, y)
			if err != nil {
				t.Errorf("Failed to get pixel at (%d, %d): %v", x, y, err)
				continue
			}
			
			if pixel.Char != expectedChar {
				t.Errorf("Expected pixel at (%d, %d) to be '%c', got '%c'", x, y, expectedChar, pixel.Char)
			}
			
		case "cursor":
			expectedX := int(assertion.Expected["x"].(float64))
			expectedY := int(assertion.Expected["y"].(float64))
			
			if ctx.Cursor.X != expectedX || ctx.Cursor.Y != expectedY {
				t.Errorf("Expected cursor at (%d, %d), got (%d, %d)", 
					expectedX, expectedY, ctx.Cursor.X, ctx.Cursor.Y)
			}
			
		case "history_count":
			expectedCount := int(assertion.Expected["count"].(float64))
			total, _ := history.GetStats()
			
			if total != expectedCount {
				t.Errorf("Expected %d commands in history, got %d", expectedCount, total)
			}
		}
	}
}

func TestJSONDrivenCommands(t *testing.T) {
	// Test scenario as JSON
	scenarioJSON := `{
		"name": "Basic drawing and cursor movement",
		"commands": [
			{
				"type": "draw",
				"params": {
					"x": 10,
					"y": 5,
					"char": "*",
					"foreground": "red",
					"bold": true
				}
			},
			{
				"type": "move_cursor",
				"params": {
					"x": 15,
					"y": 10
				}
			},
			{
				"type": "draw",
				"params": {
					"x": 15,
					"y": 10,
					"char": "#",
					"foreground": "blue"
				}
			},
			{
				"type": "erase",
				"params": {
					"x": 10,
					"y": 5
				}
			}
		],
		"assertions": [
			{
				"type": "pixel",
				"expected": {
					"x": 10,
					"y": 5,
					"char": " "
				}
			},
			{
				"type": "pixel",
				"expected": {
					"x": 15,
					"y": 10,
					"char": "#"
				}
			},
			{
				"type": "cursor",
				"expected": {
					"x": 15,
					"y": 10
				}
			},
			{
				"type": "history_count",
				"expected": {
					"count": 4
				}
			}
		]
	}`
	
	// Parse JSON scenario
	var scenario JSONTestScenario
	if err := json.Unmarshal([]byte(scenarioJSON), &scenario); err != nil {
		t.Fatalf("Failed to parse test scenario: %v", err)
	}
	
	// Set up test context
	canvas := NewMockCanvas(80, 24)
	ctx := &CommandContext{
		Canvas: canvas,
		Layers: NewMockLayerManager(),
		Cursor: CursorState{X: 0, Y: 0},
	}
	history := NewCommandHistory(100)
	
	// Execute commands
	if err := ExecuteJSONCommands(ctx, history, scenario.Commands); err != nil {
		t.Fatalf("Failed to execute JSON commands: %v", err)
	}
	
	// Validate assertions
	ValidateAssertions(t, ctx, history, scenario.Assertions)
	
	t.Logf("Successfully executed JSON test scenario: %s", scenario.Name)
}

func TestJSONCommandUndo(t *testing.T) {
	// Test undo functionality with JSON commands
	scenarioJSON := `{
		"name": "Draw and undo",
		"commands": [
			{
				"type": "draw",
				"params": {
					"x": 5,
					"y": 5,
					"char": "A"
				}
			},
			{
				"type": "draw",
				"params": {
					"x": 6,
					"y": 5,
					"char": "B"
				}
			}
		],
		"assertions": [
			{
				"type": "pixel",
				"expected": {
					"x": 5,
					"y": 5,
					"char": "A"
				}
			},
			{
				"type": "pixel",
				"expected": {
					"x": 6,
					"y": 5,
					"char": "B"
				}
			}
		]
	}`
	
	var scenario JSONTestScenario
	if err := json.Unmarshal([]byte(scenarioJSON), &scenario); err != nil {
		t.Fatalf("Failed to parse test scenario: %v", err)
	}
	
	// Set up test context
	canvas := NewMockCanvas(80, 24)
	ctx := &CommandContext{
		Canvas: canvas,
		Layers: NewMockLayerManager(),
		Cursor: CursorState{X: 0, Y: 0},
	}
	history := NewCommandHistory(100)
	
	// Execute commands
	if err := ExecuteJSONCommands(ctx, history, scenario.Commands); err != nil {
		t.Fatalf("Failed to execute JSON commands: %v", err)
	}
	
	// Validate initial state
	ValidateAssertions(t, ctx, history, scenario.Assertions)
	
	// Undo last command
	if err := history.Undo(ctx); err != nil {
		t.Fatalf("Failed to undo: %v", err)
	}
	
	// Verify B was undone but A remains
	pixelA, _ := ctx.Canvas.GetPixel(5, 5)
	pixelB, _ := ctx.Canvas.GetPixel(6, 5)
	
	if pixelA.Char != 'A' {
		t.Error("Expected pixel A to remain after undo")
	}
	if pixelB.Char != ' ' {
		t.Error("Expected pixel B to be undone")
	}
	
	t.Log("JSON-driven undo test passed")
}