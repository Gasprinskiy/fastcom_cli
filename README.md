# CLI application to internet speed test through fast.com api

> **Note:** The application is not yet available through any package manager (pacman/AUR, apt, brew, etc.).
> For now, the only way to use it is to compile it yourself.

## How to compile

### Requirements

- [Go](https://go.dev/dl/) (the minimum version is listed in `go.mod`)
- Git

Check that Go is installed:

```bash
go version
```

### Build

```bash
git clone https://github.com/<your-username>/<repo-name>.git
cd <repo-name>
go build -o fastcom_cli
```

After the build, the `fastcom_cli` binary will appear in the project directory.

Optionally, move it somewhere in your `PATH` to run it from anywhere:

```bash
sudo install -Dm755 fastcom_cli /usr/local/bin/fastcom_cli
```

## Run

```bash
./fastcom_cli
```

or, if you installed it into your `PATH`:

```bash
fastcom_cli
```

The application shows the current (`CUR`) and average (`AVG`) download speed, updating in place.