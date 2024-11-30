#!/bin/bash
watchman-make -p '*.go' '*/**/*.go' -r 'go build -o build/pacman.exe . && ./build/pacman.exe'

