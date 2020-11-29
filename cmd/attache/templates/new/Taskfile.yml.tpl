# https://taskfile.dev

version: "3"

tasks:
  test:
    cmds:
      - go test -v ./...
  
  start:
    cmds:
      - go run .