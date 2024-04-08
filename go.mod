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

replace CLI/MC => ./MCommands

require CLI/MC v0.0.0-00010101000000-000000000000

require CLI/s7cli v0.0.0-00010101000000-000000000000

require CLI/merlin_cli v0.0.0-00010101000000-000000000000 // indirect

require CLI/HugAI v0.0.0-00010101000000-000000000000 // indirect

require CLI/blackbox v0.0.0-00010101000000-000000000000 // indirect

require CLI/SearXNG v0.0.0-00010101000000-000000000000 // indirect

require CLI/Auth v0.0.0-00010101000000-000000000000 // indirect

require CLI/Misc v0.0.0-00010101000000-000000000000

require (
	CLI/TMDB v0.0.0-00010101000000-000000000000 // indirect
	github.com/PuerkitoBio/goquery v1.9.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/c-bata/go-prompt v0.2.6 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.7.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	github.com/rivo/tview v0.0.0-20240406141410-79d4cc321256 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
