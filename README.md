# Deezshots

I wanted a simple alternative to gyazo that will let me take screenshots and upload it to my own server, heavily inspired by [myazo](https://github.com/migueldemoura/myazo) (I rewrote this in golang because I didn't want to waste an hour trying to install python)

## Building & Usage

### Client

Make sure you have [golang](https://golang.org/) installed and run the following commands:

```bash
go build -o deezshots-client main.go
```

You should be able to use the binary by running:

```
./deezshots-client
```

### Server

Very similar to the client, make sure you have [golang](https://golang.org/) installed and run the following commands:

```bash
go build -o deezshots-server main.go
```

```
./deezshots-server
```

Replace keys and stuff in the code too lazy to make a config file.