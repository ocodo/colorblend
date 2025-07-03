```sh
Usage:
  colorblend [flags]

Flags:
  -c, --color-direction string      Direction for color hue interpolation (shortest, clockwise, counter-clockwise) (default "shortest")
  -e, --end-color string            Ending HEX color (e.g., #00FFFF for cyan) (default "#00FFFF")
  -g, --gradient-direction string   Direction of the gradient (horizontal, vertical) (default "horizontal")
  -h, --help                        Show help message
  -i, --invert                      Invert the gradient direction
  -s, --start-color string          Starting HEX color (e.g., #FF00FF for magenta) (default "#FF00FF")
  -t, --steps int                   Number of discrete color steps (0 for smooth gradient)
  -v, --version                     Show version information

Examples:
  echo "Hello, World!" | colorblend
  echo "Colorful!" | colorblend --start-color #FF0000 --end-color #00FF00
  cat my_file.txt | colorblend --start-color #FFFF00 --end-color #0000FF
  echo "Stepped!" | colorblend --steps 5 --start-color #FF0000 --end-color #0000FF
  echo "Vertical!" | colorblend --gradient-direction vertical --start-color #FF0000 --end-color #0000FF
  echo "Inverted!" | colorblend --invert
  echo "Hue Direction" | colorblend --hue-direction clockwise
```
