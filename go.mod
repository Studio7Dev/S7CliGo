module CLI

go 1.22.1

replace CLI/s7cli => ./s7cli

replace CLI/merlin_cli => ./Utils/Merlin

replace CLI/HugAI => ./Utils/HuggingFace

replace CLI/blackbox => ./Utils/BlackBox

replace CLI/SearXNG => ./Utils/Searx

replace CLI/TMDB => ./Utils/Tmdb

replace CLI/Auth => ./Auth

replace CLI/Misc => ./Misc

// replace CLI/MC => ./MCommands

// require CLI/MC v0.0.0-00010101000000-000000000000

require CLI/s7cli v0.0.0-00010101000000-000000000000

require CLI/merlin_cli v0.0.0-00010101000000-000000000000 // indirect

require CLI/HugAI v0.0.0-00010101000000-000000000000

require CLI/blackbox v0.0.0-00010101000000-000000000000

require CLI/SearXNG v0.0.0-00010101000000-000000000000

require CLI/Auth v0.0.0-00010101000000-000000000000 // indirect

require (
	CLI/Misc v0.0.0-00010101000000-000000000000
	CLI/TMDB v0.0.0-00010101000000-000000000000
)

require (
	github.com/PuerkitoBio/goquery v1.9.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/c-bata/go-prompt v0.2.6 // indirect
	github.com/catppuccin/go v0.2.0 // indirect
	github.com/charmbracelet/bubbles v0.18.0 // indirect
	github.com/charmbracelet/bubbletea v0.25.0 // indirect
	github.com/charmbracelet/huh v0.3.0 // indirect
	github.com/charmbracelet/lipgloss v0.9.1 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.7.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	github.com/rivo/tview v0.0.0-20240406141410-79d4cc321256 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sahilm/fuzzy v0.1.1-0.20230530133925-c48e322e2a8f // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
