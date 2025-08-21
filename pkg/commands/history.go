package commands

import (
	"fmt"
	"sync"
)

// CommandHistory manages the undo/redo stack
type CommandHistory struct {
	commands    []Command
	currentPos  int // Position in the command stack
	maxCommands int // Maximum number of commands to keep
	mutex       sync.RWMutex
}

// NewCommandHistory creates a new command history
func NewCommandHistory(maxCommands int) *CommandHistory {
	return &CommandHistory{
		commands:    make([]Command, 0, maxCommands),
		currentPos:  -1,
		maxCommands: maxCommands,
	}
}

// Execute executes a command and adds it to the history
func (ch *CommandHistory) Execute(ctx *CommandContext, cmd Command) error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	// Execute the command
	if err := cmd.Execute(ctx); err != nil {
		return err
	}

	// Clear any commands after current position (for redo)
	if ch.currentPos < len(ch.commands)-1 {
		ch.commands = ch.commands[:ch.currentPos+1]
	}

	// Add the new command
	ch.commands = append(ch.commands, cmd)
	ch.currentPos++

	// Limit the number of commands
	if len(ch.commands) > ch.maxCommands {
		// Remove the oldest command
		copy(ch.commands, ch.commands[1:])
		ch.commands = ch.commands[:len(ch.commands)-1]
		ch.currentPos--
	}

	return nil
}

// Undo undoes the last command
func (ch *CommandHistory) Undo(ctx *CommandContext) error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	if ch.currentPos < 0 {
		return fmt.Errorf("no commands to undo")
	}

	cmd := ch.commands[ch.currentPos]
	if err := cmd.Undo(ctx); err != nil {
		return fmt.Errorf("failed to undo command: %w", err)
	}

	ch.currentPos--
	return nil
}

// Redo redoes the next command
func (ch *CommandHistory) Redo(ctx *CommandContext) error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	if ch.currentPos >= len(ch.commands)-1 {
		return fmt.Errorf("no commands to redo")
	}

	ch.currentPos++
	cmd := ch.commands[ch.currentPos]
	if err := cmd.Execute(ctx); err != nil {
		ch.currentPos-- // Rollback position on error
		return fmt.Errorf("failed to redo command: %w", err)
	}

	return nil
}

// CanUndo returns true if there are commands to undo
func (ch *CommandHistory) CanUndo() bool {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	return ch.currentPos >= 0
}

// CanRedo returns true if there are commands to redo
func (ch *CommandHistory) CanRedo() bool {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	return ch.currentPos < len(ch.commands)-1
}

// GetHistory returns a copy of the command history for display
func (ch *CommandHistory) GetHistory() []string {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	history := make([]string, len(ch.commands))
	for i, cmd := range ch.commands {
		marker := "  "
		if i <= ch.currentPos {
			marker = "● "
		}
		history[i] = marker + cmd.Description()
	}

	return history
}

// Clear clears the command history
func (ch *CommandHistory) Clear() {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	ch.commands = ch.commands[:0]
	ch.currentPos = -1
}

// GetStats returns statistics about the command history
func (ch *CommandHistory) GetStats() (total, position int) {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	return len(ch.commands), ch.currentPos + 1
}
