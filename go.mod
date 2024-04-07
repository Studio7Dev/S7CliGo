module CLI

go 1.22.1

replace CLI/s7cli => ./s7cli

replace CLI/merlin_cli => ./Utils/Merlin

replace CLI/HugAI => ./Utils/HuggingFace

replace CLI/blackbox => ./Utils/BlackBox

replace CLI/SearXNG => ./Utils/Searx

replace CLI/TMDB => ./Utils/Tmdb

replace CLI/Auth => ./Auth

require CLI/s7cli v0.0.0-00010101000000-000000000000

require CLI/merlin_cli v0.0.0-00010101000000-000000000000

require CLI/HugAI v0.0.0-00010101000000-000000000000

require CLI/blackbox v0.0.0-00010101000000-000000000000

require CLI/SearXNG v0.0.0-00010101000000-000000000000

require CLI/Auth v0.0.0-00010101000000-000000000000

require CLI/TMDB v0.0.0-00010101000000-000000000000

require (
	github.com/PuerkitoBio/goquery v1.9.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/c-bata/go-prompt v0.2.6 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/jpillora/overseer v1.1.6 // indirect
	github.com/jpillora/s3 v1.1.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)
