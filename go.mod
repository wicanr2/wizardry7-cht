module github.com/wicanr2/wizardry7-cht

go 1.24.0

toolchain go1.24.13

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/wicanr2/dosgolem v0.0.0
)

replace github.com/wicanr2/dosgolem => ./worktrees/dosgolem

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.4.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
