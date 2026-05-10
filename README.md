[![Build Status](https://travis-ci.org/omakoto/a2h.svg?branch=master)](https://travis-ci.org/omakoto/a2h)
# a2h — ANSI to HTML converter

`a2h` converts terminal output containing ANSI escape sequences into a self-contained HTML file, preserving colors and text attributes.

## Features

- Supports standard text attributes: bold, faint, italic, underline, blink, reverse video, conceal, strikethrough
- Supports standard 8/16 ANSI colors, xterm 256-color palette, and 24-bit (truecolor) RGB
- Visualizes control characters (e.g. BS → `^H`)
- Outputs a complete, self-contained HTML page with embedded CSS

## Installation

```
go install github.com/omakoto/a2h@latest
```

## Usage

```
a2h [flags] [file ...]
```

Read from files or stdin, write HTML to stdout:

```sh
# Convert a file
a2h output.txt > output.html

# Pipe command output
some-command 2>&1 | a2h > output.html
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-title` | `A2H` | HTML page title |
| `-bg-color` | `#000000` | Background color |
| `-text-color` | `#ffffff` | Default text color |
| `-font-size` | `9pt` | Font size |
| `-gamma` | `1.0` | Gamma correction for RGB color conversion |
| `-auto-flush` | false | Flush output after each line |
| `-no-convert-controls` | false | Don't visualize control characters |

## See also

[The Rust version](https://github.com/omakoto/a2h-rs)
