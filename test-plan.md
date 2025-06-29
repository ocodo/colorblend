# 🧪 Test Plan for `colorblend`

The program applies a color gradient to text input via standard input. It supports various options like color space, gradient direction, formatting (bold/italic/underline), and more.

This document outlines a comprehensive test plan covering functionality, edge cases, error handling, and integration scenarios.

---

## 🎯 1. **Test Objectives**

- Ensure the program correctly applies color gradients to text input.
- Validate that all command-line flags work as expected.
- Confirm correct behavior with various input types (empty, single line, multi-line).
- Handle invalid inputs gracefully with helpful messages.
- Ensure compatibility with common terminal emulators supporting ANSI truecolor.

---

## 📋 2. **Test Categories**

| Category | Description |
|---------|-------------|
| ✅ Functional Testing | Validate each flag and combination of flags works as intended. |
| ⚠️ Edge Case Testing | Test boundary conditions and special inputs. |
| ❌ Error Handling | Verify proper handling of invalid colors, unsupported flags, etc. |
| 🔀 Integration Testing | Ensure interaction between flags (e.g., invert + steps) works. |
| 💥 Regression Testing | Ensure future updates do not break existing behavior. |

---

## 🧩 3. **Test Scenarios**

### ✅ Functional Tests

| Test ID | Description | Input Example | Expected Output |
|--------|-------------|---------------|-----------------|
| F01 | Basic horizontal gradient | `echo "Hello World" | colorblend` | Gradient from magenta (`#FF00FF`) to cyan (`#00FFFF`) applied per character |
| F02 | Vertical gradient on multiple lines | `echo -e "Line1\nLine2" | colorblend --gradient-direction vertical` | First line = start color, second = end color |
| F03 | Bold + Italic + Underline formatting | `echo "Styled Text" | colorblend --bold --italic --underline` | Styled text in gradient with bold, italic, underline |
| F04 | Invert gradient | `echo "Inverted" | colorblend --invert` | Start color appears at end of string |
| F05 | Stepped gradient | `echo "Stepped!" | colorblend --steps 3 --start-color #FF0000 --end-color #0000FF` | 3 distinct red → blue steps |
| F06 | HSL color space | `echo "HSL" | colorblend --colorspace hsl` | Colors interpolated using HSL logic |
| F07 | LAB color space | `echo "LAB" | colorblend --colorspace lab` | Colors interpolated using perceptually uniform space |
| F08 | Hue direction clockwise | `echo "Clockwise Hue" | colorblend --colorspace hsl --hue-direction clockwise` | Hue interpolation follows clockwise path |
| F09 | Version flag | `colorblend --version` | Outputs `colorblend v1.0.0` |
| F10 | Help flag | `colorblend --help` | Displays usage, options, examples |

---

### ⚠️ Edge Case Tests

| Test ID | Description | Input Example | Expected Output |
|--------|-------------|---------------|-----------------|
| E01 | Empty input | `echo "" | colorblend` | Only reset code printed (`\x1b[0m\n`) |
| E02 | Multi-line empty input | `echo -e "\n\n\n" | colorblend --gradient-direction vertical` | Each line gets appropriate gradient step |
| E03 | Single character input | `echo "A" | colorblend` | Single character with start color |
| E04 | Very long line | `yes 'A' | head -c 1000 | colorblend` | All characters colored smoothly |
| E05 | Zero steps | `echo "Smooth" | colorblend --steps 0` | Smooth gradient (same as no steps) |
| E06 | Max steps | `echo "MaxSteps" | colorblend --steps 1000` | Step count capped logically (no crash) |
| E07 | Negative steps | `echo "Negative" | colorblend --steps -1` | Error message shown, help displayed |

---

### ❌ Error Handling Tests

| Test ID | Description | Input Example | Expected Output |
|--------|-------------|---------------|-----------------|
| ER01 | Invalid start color | `echo "Bad Color" | colorblend --start-color #GGGGGG` | Error: Invalid format message |
| ER02 | Invalid end color | `echo "Bad Color" | colorblend --end-color #XYZXYZ` | Same as above |
| ER03 | Invalid gradient direction | `echo "Invalid Dir" | colorblend --gradient-direction diagonal` | Error message |
| ER04 | Invalid colorspace | `echo "Bad CS" | colorblend --colorspace hsv` | Error message |
| ER05 | Invalid hue direction | `echo "Bad HD" | colorblend --colorspace hsl --hue-direction around-the-world` | Error message |
| ER06 | Unexpected arguments | `echo "Oops" | colorblend unused_arg` | Error: unexpected argument |

---

### 🔀 Integration Tests

| Test ID | Description | Input Example | Expected Output |
|--------|-------------|---------------|-----------------|
| I01 | Invert + Steps | `echo "InvertStep" | colorblend --invert --steps 5` | Stepped gradient, inverted order |
| I02 | HSL + Clockwise Hue | `echo "HuePath" | colorblend --colorspace hsl --hue-direction clockwise` | Correct hue interpolation |
| I03 | Bold + Italic + Underline + Color | `echo "BoldItalicUnderline" | colorblend --bold --italic --underline` | Text styled with all attributes |
| I04 | Vertical + Bold + Invert | `echo -e "Line1\nLine2" | colorblend --gradient-direction vertical --bold --invert` | Line2 gets start color due to inversion |
| I05 | Multiple Flags Together | `echo "MultiFlag" | colorblend --start-color #00FF00 --end-color #0000FF --colorspace lab --bold --italic --underline --invert --steps 5` | All flags respected together |

---

## 🧹 4. **Cleanup & Reset**

After each test:

- Ensure terminal output resets properly using `\x1b[0m`.
- If ANSI escape codes are misapplied or left hanging, it should be flagged as a bug.

---

## 📦 5. **Optional Future Enhancements (Not Required for Current Scope)**

| Enhancement | Description |
|-------------|-------------|
| Output to file | Add `--output FILE` support |
| Background colors | Support background color gradients |
| Terminal auto-detection | Detect if stdout is a TTY and disable colors if not |
| ANSI escape stripping | Strip ANSI before applying gradient if needed |

---

## ✅ 6. **Pass/Fail Criteria**

- **Pass:** Program behaves as expected per test case; outputs correct ANSI-colored text or error message.
- **Fail:** Incorrect color application, panic, crash, or incorrect flag handling.
- **Skip:** If external dependencies (like terminal capabilities) interfere with test.

---

## 📄 7. **Notes for Implementation**

- Use `go test` or shell scripts to automate execution of these tests.
- For visual inspection, use terminals that support truecolor (e.g., iTerm2, Windows Terminal, Alacritty).
- Consider writing unit tests for `getGradientColor()` to verify accurate color blending across different spaces.



