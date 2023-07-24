# connect-go-embed-vue

Vue.js の SPA を go:embed で埋め込んで connect-go で実装した API を呼び出すサンプル
https://blog.arithmeticoverflow.com/posts/20230801-go-connect-embed-vue/

## Requirements

+ Go 1.20
+ NodeJS 18
+ make command

### Confirmed

+ Ubuntu 22.04.2 LTS
+ go version go1.20.6 linux/amd64
+ node v18.16.0

## Getting Started

1. Clone this repository.
2. Install dependencies. 
   + `make prepare`
3. Run the dev servers.
   + `make dev`
   + `cd frontend; npm run dev;` (on another terminal)

## Build

Run `make build`
