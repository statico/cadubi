package io

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cadubi/cadubi/pkg/canvas"
	"github.com/cadubi/cadubi/pkg/commands"
)

func TestSaveAndLoadFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.cadubi")

	// Create a test canvas with some content
	c := canvas.NewCanvas(80, 24)
	
	// Add some test pixels
	pixel1 := commands.Pixel{
		Char:       '*',
		Foreground: "red",
		Background: "",
		Bold:       true,
	}
	pixel2 := commands.Pixel{
		Char:       '#',
		Foreground: "blue",
		Background: "",
		Bold:       false,
	}
	
	// Set pixels at different positions
	err := c.SetPixel(5, 3, pixel1)
	if err != nil {
		t.Fatalf("Failed to set pixel 1: %v", err)
	}
	
	err = c.SetPixel(10, 7, pixel2)
	if err != nil {
		t.Fatalf("Failed to set pixel 2: %v", err)
	}

	// Test saving
	t.Run("Save file", func(t *testing.T) {
		err := SaveToFile(testFile, c)
		if err != nil {
			t.Errorf("SaveToFile failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("File was not created: %s", testFile)
		}
	})

	// Test loading
	t.Run("Load file", func(t *testing.T) {
		loadedCanvas, layerInfos, err := LoadFromFile(testFile)
		if err != nil {
			t.Errorf("LoadFromFile failed: %v", err)
			return
		}

		// Verify we got back the same content
		if len(layerInfos) == 0 {
			t.Errorf("No layers loaded")
			return
		}

		// Check that the pixels are where we expect them
		loadedPixel1, err := loadedCanvas.GetPixel(5, 3)
		if err != nil {
			t.Errorf("Failed to get pixel 1: %v", err)
			return
		}

		if loadedPixel1.Char != '*' {
			t.Errorf("Expected char '*', got '%c'", loadedPixel1.Char)
		}
		if loadedPixel1.Foreground != "red" {
			t.Errorf("Expected foreground 'red', got '%s'", loadedPixel1.Foreground)
		}
		if !loadedPixel1.Bold {
			t.Errorf("Expected bold=true, got %v", loadedPixel1.Bold)
		}

		loadedPixel2, err := loadedCanvas.GetPixel(10, 7)
		if err != nil {
			t.Errorf("Failed to get pixel 2: %v", err)
			return
		}

		if loadedPixel2.Char != '#' {
			t.Errorf("Expected char '#', got '%c'", loadedPixel2.Char)
		}
		if loadedPixel2.Foreground != "blue" {
			t.Errorf("Expected foreground 'blue', got '%s'", loadedPixel2.Foreground)
		}
		if loadedPixel2.Bold {
			t.Errorf("Expected bold=false, got %v", loadedPixel2.Bold)
		}
	})

	// Test loading non-existent file
	t.Run("Load non-existent file", func(t *testing.T) {
		_, _, err := LoadFromFile("does-not-exist.cadubi")
		if err == nil {
			t.Errorf("Expected error when loading non-existent file, got nil")
		}
	})
}

func TestSaveEmptyCanvas(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.cadubi")

	// Create empty canvas
	c := canvas.NewCanvas(40, 20)

	// Save empty canvas
	err := SaveToFile(testFile, c)
	if err != nil {
		t.Errorf("Failed to save empty canvas: %v", err)
	}

	// Load it back
	loadedCanvas, layerInfos, err := LoadFromFile(testFile)
	if err != nil {
		t.Errorf("Failed to load empty canvas: %v", err)
		return
	}

	// Should have at least one layer
	if len(layerInfos) == 0 {
		t.Errorf("Empty canvas should have at least one layer")
	}

	// Verify dimensions (for empty canvas, GetDimensions returns defaults)
	width, height := loadedCanvas.GetDimensions()
	// Empty canvas returns default dimensions since no pixels are set
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive dimensions, got %dx%d", width, height)
	}
}

func TestCanvasModifiedTracking(t *testing.T) {
	c := canvas.NewCanvas(10, 10)

	// New canvas should not be modified
	if c.IsModified() {
		t.Errorf("New canvas should not be modified")
	}

	// Set a pixel - should mark as modified
	pixel := commands.Pixel{Char: 'X', Foreground: "white"}
	err := c.SetPixel(1, 1, pixel)
	if err != nil {
		t.Fatalf("Failed to set pixel: %v", err)
	}

	if !c.IsModified() {
		t.Errorf("Canvas should be modified after setting pixel")
	}

	// Mark as not modified
	c.SetModified(false)
	if c.IsModified() {
		t.Errorf("Canvas should not be modified after SetModified(false)")
	}

	// Test filename tracking
	if c.HasFilename() {
		t.Errorf("New canvas should not have filename")
	}

	c.SetFilename("test.cadubi")
	if !c.HasFilename() {
		t.Errorf("Canvas should have filename after SetFilename")
	}

	if c.GetFilename() != "test.cadubi" {
		t.Errorf("Expected filename 'test.cadubi', got '%s'", c.GetFilename())
	}
}