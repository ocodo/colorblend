package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

var osExit = os.Exit

// getGradientColor interpolates a color based on progress and returns the ANSI truecolor color part (e.g., "38;2;R;G;Bm").
// It now uses go-colorful for robust color space handling and interpolation.
func getGradientColor(progress float64, startHex, endHex string, colorspace, hueDirection string) (string, error) {
	// Parse start and end hex colors using go-colorful
	startColor, err := colorful.Hex(startHex)
	if err != nil {
		return "", fmt.Errorf("invalid start hex color: %s (%w)", startHex, err)
	}
	endColor, err := colorful.Hex(endHex)
	if err != nil {
		return "", fmt.Errorf("invalid end hex color: %s (%w)", endHex, err)
	}

	var interpolatedColor colorful.Color

	// Determine the blending method based on colorspace and hueDirection
	switch colorspace {
	case "rgb":
		interpolatedColor = startColor.BlendRgb(endColor, progress)
	case "hcl":
		// Your installed go-colorful v1.2.0 does not define HuePath or accept it in BlendHcl.
		// BlendHcl will use its internal default hue path (likely shortest).
		interpolatedColor = startColor.BlendHcl(endColor, progress)
	case "lab":
		interpolatedColor = startColor.BlendLab(endColor, progress)
	default:
		return "", fmt.Errorf("unsupported colorspace: %s", colorspace)
	}

	r, g, b := interpolatedColor.Clamped().RGB255()

	return fmt.Sprintf("38;2;%d;%d;%dm", r, g, b), nil
}

func main() {
	// Define command-line flags
	showHelp := flag.Bool("help", false, "Show this help message")
	showVersion := flag.Bool("version", false, "Show version information")

	startColor := flag.String("start-color", "#FF00FF", "Starting HEX color (e.g., #FF00FF for magenta)")
	endColor := flag.String("end-color", "#00FFFF", "Ending HEX color (e.g., #00FFFF for cyan)")

	gradientDirection := flag.String("gradient-direction", "horizontal", "Direction of the gradient (horizontal, vertical).")
	colorspace := flag.String("colorspace", "rgb", "Color space for interpolation (rgb, hcl, lab).")
	hueDirection := flag.String("hue-direction", "shortest", "Direction for hue interpolation in HCL (shortest, clockwise, counter-clockwise). Only applies if colorspace is HCL or LAB.")
	steps := flag.Int("steps", 0, "Number of discrete color steps (0 for smooth gradient).")
	invert := flag.Bool("invert", false, "Invert the gradient direction (e.g., end color at start).")

	// Set a custom usage function for --help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Applies a color gradient to text read from standard input.")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults() // Prints default usage for defined flags
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  echo \"Hello, World!\" | colorblend")
		fmt.Fprintln(os.Stderr, "  echo \"Colorful!\" | colorblend --start-color #FF0000 --end-color #00FF00")
		fmt.Fprintln(os.Stderr, "  cat my_file.txt | colorblend -start-color #FFFF00 -end-color #0000FF")
		fmt.Fprintln(os.Stderr, "  echo \"Stepped!\" | colorblend --steps 5 --start-color #FF0000 --end-color #0000FF")
		fmt.Fprintln(os.Stderr, "  echo \"Vertical!\" | colorblend --gradient-direction vertical --start-color #FF0000 --end-color #0000FF")
		fmt.Fprintln(os.Stderr, "  echo \"Inverted!\" | colorblend --invert")
		fmt.Fprintln(os.Stderr, "  echo \"HCL Gradient!\" | colorblend -s #FF0000 -e #0000FF --colorspace hcl --hue-direction clockwise")
	}

	// Parse command-line arguments
	flag.Parse()

	// Handle --version flag
	if *showVersion {
		fmt.Println("colorblend v1.0.0")
		osExit(0)
	}

	// Handle --help flag explicitly
	if *showHelp {
		flag.Usage()
		osExit(0)
	}

	// Validate color format using colorful.Hex
	_, err := colorful.Hex(*startColor)
	if err != nil { // Corrected: Check for error != nil
		fmt.Fprintf(os.Stderr, "Error: Invalid format for --start-color: %s. Must be a 7-character hex string (e.g., #RRGGBB). Details: %v\n\n", *startColor, err)
		flag.Usage()
		osExit(1)
	}
	_, err = colorful.Hex(*endColor)
	if err != nil { // Corrected: Check for error != nil
		fmt.Fprintf(os.Stderr, "Error: Invalid format for --end-color: %s. Must be a 7-character hex string (e.g., #RRGGBB). Details: %v\n\n", *endColor, err)
		flag.Usage()
		osExit(1)
	}

	// Validate other flag values
	if *gradientDirection != "horizontal" && *gradientDirection != "vertical" {
		fmt.Fprintf(os.Stderr, "Error: Invalid value for --gradient-direction: %s. Must be 'horizontal' or 'vertical'.\n\n", *gradientDirection)
		flag.Usage()
		osExit(1)
	}
	if *colorspace != "rgb" && *colorspace != "hcl" && *colorspace != "lab" {
		fmt.Fprintf(os.Stderr, "Error: Invalid value for --colorspace: %s. Must be 'rgb', 'hcl', or 'lab'.\n\n", *colorspace)
		flag.Usage()
		osExit(1)
	}
	// Hue direction is only relevant for HCL/LAB.
	// Since HuePath is not defined in your installed module, the --hue-direction flag
	// will have no effect on HCL/LAB blending, as BlendHcl/BlendLab will use their default.
	if (*colorspace == "hcl" || *colorspace == "lab") && (*hueDirection != "shortest" && *hueDirection != "clockwise" && *hueDirection != "counter-clockwise") {
		fmt.Fprintf(os.Stderr, "Error: Invalid value for --hue-direction: %s. Must be 'shortest', 'clockwise', or 'counter-clockwise' when using HCL/LAB colorspace.\n\n", *hueDirection)
		flag.Usage()
		osExit(1)
	}
	if *steps < 0 {
		fmt.Fprintf(os.Stderr, "Error: --steps cannot be negative.\n\n")
		flag.Usage()
		osExit(1)
	}

	// If there are any remaining (non-flag) arguments, it might indicate incorrect usage.
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Error: Unexpected arguments: %s\n\n", strings.Join(flag.Args(), " "))
		flag.Usage()
		osExit(1)
	}

	// Read all of stdin into lines (necessary for vertical gradient)
	reader := bufio.NewReader(os.Stdin)
	var lines [][]rune
	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			osExit(1)
		}
		lineRunes := []rune(string(lineBytes)) // Convert line bytes to runes
		lines = append(lines, lineRunes)
	}

	if len(lines) == 0 {
		fmt.Printf("\x1b[0m\n") // Reset color and print newline even if no input
		return
	}

	// Determine total characters for horizontal, or total lines for vertical
	var totalGradientUnits int
	if *gradientDirection == "horizontal" {
		totalGradientUnits = 0
		for _, line := range lines {
			totalGradientUnits += len(line)
		}
		if totalGradientUnits == 0 { // Edge case for empty lines (only newlines in input)
			// Print newlines with reset color for each empty line
			for range lines {
				fmt.Printf("\x1b[0m\n")
			}
			return
		}
	} else { // vertical
		totalGradientUnits = len(lines)
	}

	// Character counter for horizontal gradient progress
	charCountHorizontal := 0

	// Process lines based on gradient direction
	for lineIndex, line := range lines {
		if *gradientDirection == "vertical" && len(line) == 0 && totalGradientUnits > 1 {
            // Handle empty lines specifically for vertical gradients to ensure they contribute to progress
            // but still print a newline with reset.
            var progress float64
            if totalGradientUnits <= 1 {
                progress = 0.0
            } else {
                progress = float64(lineIndex) / float64(totalGradientUnits-1)
            }
            if *invert {
                progress = 1.0 - progress
            }
            if *steps > 0 {
                progress = math.Round(progress*float64(*steps)) / float64(*steps)
            }
            // Get color for the "empty" line based on its vertical position
            colorPart, err := getGradientColor(progress, *startColor, *endColor, *colorspace, *hueDirection)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error getting gradient color for empty line: %v\n", err)
                osExit(1)
            }
            // Print the empty line with its calculated color and format
            fmt.Printf("\x1b[%s%s\n", colorPart)
            continue // Move to next line
        }

		for _, char := range line {		
			var progress float64

			if *gradientDirection == "horizontal" {
				if totalGradientUnits <= 1 {
					progress = 0.0
				} else {
					progress = float64(charCountHorizontal) / float64(totalGradientUnits-1)
				}
				charCountHorizontal++
			} else { // vertical
				if totalGradientUnits <= 1 {
					progress = 0.0
				} else {
					progress = float64(lineIndex) / float64(totalGradientUnits-1)
				}
			}

			// Apply invert if flag is set
			if *invert {
				progress = 1.0 - progress
			}

			// Apply steps if flag is set (quantize progress)
			if *steps > 0 {
				progress = math.Round(progress*float64(*steps)) / float64(*steps)
			}

			// Get the ANSI color code for the current character/line using go-colorful
			colorPart, err := getGradientColor(progress, *startColor, *endColor, *colorspace, *hueDirection)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting gradient color: %v\n", err)
				osExit(1)
			}

			fmt.Printf("\x1b[%s%s%c", colorPart, char)
		}
		// Print newline at end of line (original line breaks), but only if not an empty line already handled
		if !(*gradientDirection == "vertical" && len(line) == 0 && totalGradientUnits > 1) {
            fmt.Printf("\n")
        }
	}

	// Reset all formatting at the end
	fmt.Printf("\x1b[0m\n")
}
