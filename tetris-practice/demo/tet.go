package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

// Grid dimensions
const gridRows = 20
const gridCols = 12

// Tetris shapes (relative coordinates)
var shapes = [][][][]int{
	// I shape
	{{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, {{0, 0}, {1, 0}, {2, 0}, {3, 0}}},
	// J shape
	{{{0, 0}, {1, 0}, {1, 1}, {1, 2}}, {{0, 1}, {1, 1}, {2, 1}, {2, 0}}, {{1, 0}, {1, 1}, {1, 2}, {2, 2}}, {{0, 2}, {1, 2}, {2, 2}, {2, 3}}}, // Corrected J shape rotations
	// L shape
	{{{0, 2}, {1, 0}, {1, 1}, {1, 2}}, {{0, 1}, {1, 1}, {2, 1}, {2, 2}}, {{1, 0}, {1, 1}, {1, 2}, {2, 0}}, {{0, 0}, {1, 0}, {2, 0}, {2, 1}}}, // Corrected L shape rotations
	// O shape
	{{{0, 0}, {0, 1}, {1, 0}, {1, 1}}}, // O shape has only one rotation state
	// S shape
	{{{0, 1}, {0, 2}, {1, 0}, {1, 1}}, {{0, 0}, {1, 0}, {1, 1}, {2, 1}}},
	// T shape
	{{{0, 1}, {1, 0}, {1, 1}, {1, 2}}, {{0, 1}, {1, 1}, {1, 2}, {2, 1}}, {{1, 0}, {1, 1}, {1, 2}, {2, 1}}, {{0, 1}, {1, 0}, {1, 1}, {2, 1}}}, // Corrected T shape rotations
	// Z shape
	{{{0, 0}, {0, 1}, {1, 1}, {1, 2}}, {{0, 2}, {1, 1}, {1, 2}, {2, 1}}},
}

var shapeColors = []string{
	"\033[36m", // Cyan (I)
	"\033[34m", // Blue (J)
	"\033[33m", // Yellow (L)
	"\033[32m", // Green (O)
	"\033[31m", // Red (S)
	"\033[35m", // Magenta (T)
	"\033[37m", // White (Z)
}

const resetColor = "\033[0m"

type Piece struct {
	Shape     [][]int
	Color     string
	Row       int // Top-left corner of the piece's bounding box
	Col       int // Top-left corner of the piece's bounding box
	Rotation  int
	ShapeType int // Index in the shapes array
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize keyboard
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Persistent grid for landed cells
	grid := make([][]string, gridRows)
	for i := range grid {
		grid[i] = make([]string, gridCols)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	// Game state
	var currentPiece *Piece
	score := 0
	gameOver := false

	// Clear screen initially (ANSI: \033[2J)
	fmt.Print("\033[2J")

	// Channel for keyboard input
	keyChan := make(chan rune)
	go func() {
		for {
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEsc || char == 'q' || char == 'Q' {
				keyChan <- 'q' // Signal quit
				return
			}
			if key == keyboard.KeyArrowLeft {
				keyChan <- 'a'
			} else if key == keyboard.KeyArrowRight {
				keyChan <- 'd'
			} else if key == keyboard.KeyArrowDown {
				keyChan <- 's'
			} else if key == keyboard.KeyArrowUp {
				keyChan <- 'w' // Use 'w' for rotation for now
			} else {
				keyChan <- char
			}
		}
	}()

	// Timer for automatic falling
	fallTicker := time.NewTicker(500 * time.Millisecond) // Adjust speed here
	defer fallTicker.Stop()

	// Game loop
	for !gameOver {
		// Spawn new piece if none exists
		if currentPiece == nil {
			currentPiece = spawnNewPiece()
			if checkCollision(grid, currentPiece, currentPiece.Row, currentPiece.Col, currentPiece.Rotation) {
				gameOver = true
				break
			}
		}

		// Redraw grid
		fmt.Print("\033[1;1H") // Move cursor to top-left
		drawGrid(grid, currentPiece)
		fmt.Printf("\nScore: %d\n", score) // Display score

		select {
		case char := <-keyChan:
			if char == 'q' {
				gameOver = true
				break
			}
			handleInput(grid, currentPiece, char, &gameOver, &score)
		case <-fallTicker.C:
			handleFall(grid, currentPiece, &currentPiece, &gameOver, &score)
		}
	}

	// Game Over message
	fmt.Print("\033[22;1H\n") // Move cursor below grid
	fmt.Println("Game Over!")
	fmt.Printf("Final Score: %d\n", score)
}

// spawnNewPiece creates a new random Tetris piece
func spawnNewPiece() *Piece {
	shapeType := rand.Intn(len(shapes))
	shape := shapes[shapeType][0] // Start with the first rotation state
	color := shapeColors[shapeType]
	// Calculate initial column to center the piece
	maxCol := 0
	for _, coord := range shape {
		if coord[1] > maxCol {
			maxCol = coord[1]
		}
	}
	initialCol := (gridCols - (maxCol + 1)) / 2

	return &Piece{
		Shape:     shape,
		Color:     color,
		Row:       0, // Start at the top
		Col:       initialCol,
		Rotation:  0,
		ShapeType: shapeType,
	}
}

// handleInput processes keyboard input
func handleInput(grid [][]string, piece *Piece, char rune, gameOver *bool, score *int) {
	newRow := piece.Row
	newCol := piece.Col
	newRotation := piece.Rotation

	moved := false

	switch char {
	case 'a': // Move left
		newCol--
		moved = true
	case 'd': // Move right
		newCol++
		moved = true
	case 's': // Move down (instant drop)
		// Find the lowest possible row the piece can occupy
		dropRow := piece.Row
		for {
			if checkCollision(grid, piece, dropRow+1, piece.Col, piece.Rotation) {
				break
			}
			dropRow++
		}
		newRow = dropRow
		moved = true
		// Land the piece immediately after dropping
		landPiece(grid, piece, newRow, newCol, score)
		*score += (newRow - piece.Row) // Add points for dropping
		// Set piece to nil to spawn a new one
		piece = nil // This won't update the piece in main's scope directly
		// We need a way to signal main to spawn a new piece
		// For now, we'll let the fallTicker handle the next piece spawn
		return // Exit after handling drop
	case 'w': // Rotate
		newRotation = (piece.Rotation + 1) % len(shapes[piece.ShapeType])
		// Get the shape for the new rotation state
		rotatedShape := shapes[piece.ShapeType][newRotation]
		// Check collision for the rotated shape at the current position
		if !checkCollision(grid, piece, piece.Row, piece.Col, newRotation) {
			piece.Shape = rotatedShape
			piece.Rotation = newRotation
		} else {
			// Simple wall kick attempt (try moving one step left/right)
			if !checkCollision(grid, piece, piece.Row, piece.Col-1, newRotation) {
				piece.Shape = rotatedShape
				piece.Rotation = newRotation
				piece.Col--
			} else if !checkCollision(grid, piece, piece.Row, piece.Col+1, newRotation) {
				piece.Shape = rotatedShape
				piece.Rotation = newRotation
				piece.Col++
			}
			// More advanced wall kick logic could be implemented here
		}
		// No need to update row/col if only rotating
		return
	}

	// If moved, check for collision before updating position
	if moved && !checkCollision(grid, piece, newRow, newCol, piece.Rotation) {
		piece.Row = newRow
		piece.Col = newCol
	} else if moved && char == 's' { // If moved down and collided, land the piece
		// This case is handled by the instant drop logic above, but keeping it for clarity
		// if 's' was for single step down
	}
}

// handleFall moves the piece down automatically or lands it
func handleFall(grid [][]string, piece *Piece, currentPiece **Piece, gameOver *bool, score *int) {
	if checkCollision(grid, piece, piece.Row+1, piece.Col, piece.Rotation) {
		// Collision below, land the piece
		landPiece(grid, piece, piece.Row, piece.Col, score)
		// Spawn new piece
		*currentPiece = spawnNewPiece()
		// Check for game over immediately after spawning
		if checkCollision(grid, *currentPiece, (*currentPiece).Row, (*currentPiece).Col, (*currentPiece).Rotation) {
			*gameOver = true
		}
	} else {
		// No collision, move down
		piece.Row++
	}
}

// checkCollision checks if the piece at the given position and rotation collides with the grid or boundaries
func checkCollision(grid [][]string, piece *Piece, testRow, testCol, testRotation int) bool {
	shapeToCheck := shapes[piece.ShapeType][testRotation]

	for _, coord := range shapeToCheck {
		gridRow := testRow + coord[0]
		gridCol := testCol + coord[1]

		// Check boundaries
		if gridRow < 0 || gridRow >= gridRows || gridCol < 0 || gridCol >= gridCols {
			return true // Collision with boundary
		}

		// Check collision with landed pieces
		if grid[gridRow][gridCol] != "." {
			return true // Collision with landed piece
		}
	}
	return false // No collision
}

// landPiece adds the current piece to the grid and checks for cleared lines
func landPiece(grid [][]string, piece *Piece, finalRow, finalCol int, score *int) {
	for _, coord := range piece.Shape {
		gridRow := finalRow + coord[0]
		gridCol := finalCol + coord[1]
		// Ensure coordinates are within bounds before placing
		if gridRow >= 0 && gridRow < gridRows && gridCol >= 0 && gridCol < gridCols {
			grid[gridRow][gridCol] = "x" // Mark as landed
		}
	}

	// Check for full lines (bottom-up)
	linesCleared := 0
	for r := gridRows - 1; r >= 0; r-- {
		isFull := true
		for c := 0; c < gridCols; c++ {
			if grid[r][c] == "." {
				isFull = false
				break
			}
		}
		if isFull {
			linesCleared++
			// Clear the row
			for c := 0; c < gridCols; c++ {
				grid[r][c] = "."
			}
			// Shift rows above down
			for shiftR := r - 1; shiftR >= 0; shiftR-- {
				for c := 0; c < gridCols; c++ {
					grid[shiftR+1][c] = grid[shiftR][c]
					grid[shiftR][c] = "." // Clear the original position
				}
			}
			r++ // Re-check this row after shifting
		}
	}

	// Update score based on lines cleared
	if linesCleared > 0 {
		// Simple scoring: 100 per line, bonus for multiple lines
		*score += linesCleared * 100
		if linesCleared == 2 {
			*score += 100 // Bonus for double
		} else if linesCleared == 3 {
			*score += 300 // Bonus for triple
		} else if linesCleared == 4 {
			*score += 800 // Bonus for Tetris!
		}
	}
}

// drawGrid draws the grid with the current shape
func drawGrid(grid [][]string, piece *Piece) {
	for i := 0; i < gridRows; i++ {
		for j := 0; j < gridCols; j++ {
			isPieceCell := false
			// Check if current grid position is part of the moving shape
			if piece != nil {
				for _, coord := range piece.Shape {
					shapeRow := piece.Row + coord[0]
					shapeCol := piece.Col + coord[1]
					if i == shapeRow && j == shapeCol {
						fmt.Print(piece.Color + "x" + resetColor)
						isPieceCell = true
						break
					}
				}
			}

			if !isPieceCell {
				// Draw from persistent grid
				if grid[i][j] == "." {
					fmt.Print(".")
				} else {
					// Draw landed pieces with a default color or different color
					fmt.Print("\033[90m" + "x" + resetColor) // Grey for landed pieces
				}
			}
		}
		fmt.Println() // New line after each row
	}
}
