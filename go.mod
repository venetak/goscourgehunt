module scourgehunt

go 1.24.2

require (
	github.com/hajimehoshi/ebiten v1.12.12
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	github.com/rendering v0.0.0-00010101000000-000000000000
	github/actor v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200707082815-5321531c36a2 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/mobile v0.0.0-20210208171126-f462b3930c8f // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace github/actor => ./actor

replace github/input => ./input

replace github.com/rendering => ./rendering

replace utils => ./utils
