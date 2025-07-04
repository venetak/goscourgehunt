module scourgehunt

go 1.24.2

replace github.com/actor => ./actor

replace github.com/input => ./input

replace github.com/rendering => ./rendering

replace github.com/utils => ./utils

replace github.com/game => ./game

require (
	github.com/game v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.8.8
)

require (
	github.com/actor v0.0.0-00010101000000-000000000000 // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/gameplay v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/player v0.0.0-00010101000000-000000000000 // indirect
	github.com/rendering v0.0.0-00010101000000-000000000000 // indirect
	github.com/utils v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace github.com/gameplay => ./gameplay

replace github.com/player => ./player
