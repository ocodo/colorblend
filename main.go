package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// hexToRGB converts a HEX color string (e.g., "#FF00FF") to R, G, B integer components.
func hexToRGB(hexColor string) (r, g, b int, err error) {
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return 0, 0, 0, fmt.Errorf("invalid hex color format: %s", hexColor)
	}

	r, err = strconv.ParseInt(hexColor[1:3], 16, 0)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid red component: %w", err)
	}
	g, err = strconv.ParseInt(hexColor[3:5], 16, 0)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid green component: %w", err)
	}
	b, err = strconv.ParseInt(hexColor[5:7], 16, 0)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid blue component: %w", err)
	}
	return int(r), int(g), int(b), nil
}

// getGradientColor interpolates a color based on progress and returns an ANSI truecolor escape sequence.
func getGradientColor(progress float64, startHex, endHex string) (string, error) {
	startR, startG, startB, err := hexToRGB(startHex)
	if err != nil {
		return "", err
	}
	endR, endG, endB, err := hexToRGB(endHex)
	if err != nil {
		return "", err
	}

	// Interpolate each color component
	r := int(math.Round(float64(startR) + (float64(endR)-float64(startR))*progress))
	g := int(math.Round(float64(startG) + (float64(endG)-float64(startG))*progress))
	b := int(math.Round(float64(startB) + (float64(endB)-float64(startB))*progress))

	// Ensure color values are within 0-255 range
	r = int(math.Max(0, math.Min(255, float64(r))))
	g = int(math.Max(0, math.Min(255, float64(g))))
	b = int(math.Max(0, math.Min(255, float64(b))))

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b), nil
}

func main() {
	const gradStart = "#FF00FF" // Magenta
	const gradEnd = "#00FFFF"   // Cyan

	// Read all of stdin into a string
	reader := bufio.NewReader(os.Stdin)
	inputBuilder := strings.Builder{}
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		inputBuilder.WriteRune(r)
	}
	inputString := inputBuilder.String()
	runes := []rune(inputString) // Work with runes to handle multi-byte characters correctly
	totalChars := len(runes)

	if totalChars == 0 {
		fmt.Printf("\x1b[0m\n") // Reset color and print newline even if no input
		return
	}

	for i, char := range runes {
		var progress float64
		if totalChars <= 1 {
			progress = 0.0
		} else {
			progress = float64(i) / float64(totalChars-1)
		}

		colorCode, err := getGradientColor(progress, gradStart, gradEnd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting gradient color: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s%c", colorCode, char)
	}

	// Reset color at the end
	fmt.Printf("\x1b[0m\n")
}