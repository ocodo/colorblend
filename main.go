package main

import (
    "bufio"
    "fmt"
    "io"
    "math"
    "os"
    "strings"
    
    "github.com/spf13/cobra"
    "github.com/lucasb-eyer/go-colorful"
)

var osExit = os.Exit

func blendHCLWithDirection(c1, c2 colorful.Color, t float64, hueDirection string) colorful.Color {
    h1, c1Chroma, l1 := c1.Hcl()
    h2, c2Chroma, l2 := c2.Hcl()

    h1 = math.Mod(h1+360, 360)
    h2 = math.Mod(h2+360, 360)

    var deltaH float64
    switch hueDirection {
    case "shortest", "short", "sh":
        deltaH = h2 - h1
        if deltaH > 180 {
            deltaH -= 360
        } else if deltaH < -180 {
            deltaH += 360
        }
    case "clockwise", "cw":
        deltaH = math.Mod(h2-h1+360, 360)
    case "counter-clockwise", "counterclockwise", "ccw":
        deltaH = -math.Mod(h1-h2+360, 360)
    default:
        deltaH = h2 - h1
        if deltaH > 180 {
            deltaH -= 360
        } else if deltaH < -180 {
            deltaH += 360
        }
    }

    h := math.Mod(h1+t*deltaH+360, 360)
    c := c1Chroma + t*(c2Chroma-c1Chroma)
    l := l1 + t*(l2-l1)

    return colorful.Hcl(h, c, l)
}

func getGradientColor(progress float64, startHex, endHex, hueDirection string) (string, error) {
    startColor, err := colorful.Hex(startHex)
    if err != nil {
        return "", fmt.Errorf("invalid start hex color: %s (%w)", startHex, err)
    }
    endColor, err := colorful.Hex(endHex)
    if err != nil {
        return "", fmt.Errorf("invalid end hex color: %s (%w)", endHex, err)
    }

    // Interpolate in HCL with directional hue
    interpolated := blendHCLWithDirection(startColor, endColor, progress, hueDirection)
    r, g, b := interpolated.Clamped().RGB255()
    return fmt.Sprintf("38;2;%d;%d;%dm", r, g, b), nil
}

var (
    startColor        string
    endColor          string
    gradientDirection string
    hueDirection      string
    steps             int
    invert            bool
)

var rootCmd = &cobra.Command{
    Use:   "colorblend",
    Short: "Applies a color gradient to text",
    Run: func(cmd *cobra.Command, args []string) {
        // Validate colors
        if _, err := colorful.Hex(startColor); err != nil {
            fmt.Fprintf(os.Stderr, "Error: Invalid format for --start-color: %s. Must be a 7-character hex string (e.g., #RRGGBB). Details: %v\n\n", startColor, err)
            cmd.Usage()
            os.Exit(1)
        }
        if _, err := colorful.Hex(endColor); err != nil {
            fmt.Fprintf(os.Stderr, "Error: Invalid format for --end-color: %s. Must be a 7-character hex string (e.g., #RRGGBB). Details: %v\n\n", endColor, err)
            cmd.Usage()
            os.Exit(1)
        }

        if gradientDirection != "horizontal" && gradientDirection != "vertical" {
            fmt.Fprintf(os.Stderr, "Error: Invalid value for --gradient-direction: %s. Must be 'horizontal' or 'vertical'.\n\n", gradientDirection)
            cmd.Usage()
            os.Exit(1)
        }

        if steps < 0 {
            fmt.Fprintf(os.Stderr, "Error: --steps cannot be negative.\n\n")
            cmd.Usage()
            os.Exit(1)
        }

        if len(args) > 0 {
            fmt.Fprintf(os.Stderr, "Error: Unexpected arguments: %s\n\n", strings.Join(args, " "))
            cmd.Usage()
            os.Exit(1)
        }

        // Read input lines
        reader := bufio.NewReader(os.Stdin)
        var lines [][]rune
        for {
            lineBytes, _, err := reader.ReadLine()
            if err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
                os.Exit(1)
            }
            lines = append(lines, []rune(string(lineBytes)))
        }

        if len(lines) == 0 {
            fmt.Printf("\x1b[0m\n")
            return
        }

        var totalGradientUnits int
        if gradientDirection == "horizontal" {
            for _, line := range lines {
                totalGradientUnits += len(line)
            }
            if totalGradientUnits == 0 {
                for range lines {
                    fmt.Printf("\x1b[0m\n")
                }
                return
            }
        } else {
            totalGradientUnits = len(lines)
        }

        charCountHorizontal := 0

        for lineIndex, line := range lines {
            if gradientDirection == "vertical" && len(line) == 0 && totalGradientUnits > 1 {
                progress := 0.0
                if totalGradientUnits > 1 {
                    progress = float64(lineIndex) / float64(totalGradientUnits-1)
                }
                if invert {
                    progress = 1.0 - progress
                }
                if steps > 0 {
                    progress = math.Round(progress*float64(steps)) / float64(steps)
                }

                colorPart, err := getGradientColor(progress, startColor, endColor, hueDirection)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Error getting gradient color for empty line: %v\n", err)
                    os.Exit(1)
                }

                fmt.Printf("\x1b[%s\n", colorPart)
                continue
            }

            for _, char := range line {
                progress := 0.0
                if gradientDirection == "horizontal" {
                    if totalGradientUnits > 1 {
                        progress = float64(charCountHorizontal) / float64(totalGradientUnits-1)
                    }
                    charCountHorizontal++
                } else {
                    if totalGradientUnits > 1 {
                        progress = float64(lineIndex) / float64(totalGradientUnits-1)
                    }
                }

                if invert {
                    progress = 1.0 - progress
                }
                if steps > 0 {
                    progress = math.Round(progress*float64(steps)) / float64(steps)
                }

                colorPart, err := getGradientColor(progress, startColor, endColor, hueDirection)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Error getting gradient color: %v\n", err)
                    os.Exit(1)
                }

                fmt.Printf("\x1b[%s%c", colorPart, char)
            }
            if !(gradientDirection == "vertical" && len(line) == 0 && totalGradientUnits > 1) {
                fmt.Printf("\n")
            }
        }

        fmt.Printf("\x1b[0m\n")
    },
}

func init() {
    rootCmd.Flags().StringVarP(&startColor, "start-color", "s", "#FF00FF", "Starting HEX color (e.g., #FF00FF for magenta)")
    rootCmd.Flags().StringVarP(&endColor, "end-color", "e", "#00FFFF", "Ending HEX color (e.g., #00FFFF for cyan)")
    rootCmd.Flags().StringVarP(&gradientDirection, "gradient-direction", "g", "horizontal", "Direction of the gradient (horizontal, vertical)")
    rootCmd.Flags().StringVarP(&hueDirection, "hue-direction", "u", "shortest", "Direction for hue interpolation (shortest, clockwise, counter-clockwise)")
    rootCmd.Flags().IntVarP(&steps, "steps", "t", 0, "Number of discrete color steps (0 for smooth gradient)")
    rootCmd.Flags().BoolVarP(&invert, "invert", "i", false, "Invert the gradient direction")
    rootCmd.Flags().BoolP("help", "h", false, "Show help message")
    rootCmd.Flags().BoolP("version", "v", false, "Show version information")

    rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
        cmd.SetOut(os.Stderr)
        cmd.Usage()

        fmt.Fprintln(os.Stderr, "\nExamples:")
        fmt.Fprintln(os.Stderr, "  echo \"Hello, World!\" | colorblend")
        fmt.Fprintln(os.Stderr, "  echo \"Colorful!\" | colorblend --start-color #FF0000 --end-color #00FF00")
        fmt.Fprintln(os.Stderr, "  cat my_file.txt | colorblend --start-color #FFFF00 --end-color #0000FF")
        fmt.Fprintln(os.Stderr, "  echo \"Stepped!\" | colorblend --steps 5 --start-color #FF0000 --end-color #0000FF")
        fmt.Fprintln(os.Stderr, "  echo \"Vertical!\" | colorblend --gradient-direction vertical --start-color #FF0000 --end-color #0000FF")
        fmt.Fprintln(os.Stderr, "  echo \"Inverted!\" | colorblend --invert")
        fmt.Fprintln(os.Stderr, "  echo \"Hue Direction\" | colorblend --hue-direction clockwise")
    })
    
    rootCmd.SetVersionTemplate("colorblend v1.0.0\n")

    rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
        if v, _ := cmd.Flags().GetBool("version"); v {
            fmt.Println("colorblend v1.0.0")
            os.Exit(0)
        }
        if h, _ := cmd.Flags().GetBool("help"); h {
            cmd.Usage()
            os.Exit(0)
        }
    }
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }
}
