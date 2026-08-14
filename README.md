# 🚀 letsgo

The simplest Go live-realoader - run any command when a file changes

---

# Installing

The recommended way of providing yourself the latest version of letsgo is by using:

```sh
go install github.com/ThibaultLonguepee/letsgo@latest
```

Make sure you `GOPATH` is present in your shell `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

It is recommended you add the line above at the bottom of your shell configuration. (`~/.bashrc`, `~/.zshrc`; for example)

# Usage

As of this v0. The current use of `letsgo` is strangely simmilar as to the usage of `entr`.
```sh
<filelist> | letsgo command [args...]
```

Example usage:
```sh
find . -name "*.go" | letsgo go run .
```
This will listen for any change in any `.go` files present in this directory (and children recursively).<br>
Once a change has been detected, it will run `go run .`.

> [!NOTE]
> You can technically manually enter each file by hand when running `letsgo command [args...]`<br>
> List files one by one, line by line, and press `CTRL+D` (SIGTERM) to end your input.

---

# Motivation

I personally do not recommend using `letsgo` for web servers, scripts, etc; as other OSS projects do it better I think.<br>
While working on a termbox app written in Go, I was faced with having to re-run the application a lot.<br>
I tried many alternatives to live-reloading like `Air`, `wgo`, `entr`... But none worked the way I was trying to make them work.<br>
I ended up quickly putting together the v0 for this, inspired heavily by `entr`.<br>

# Roadmap

Since `letsgo` satisfied my personnal needs, I'll take thins a step further:
- [x] Live-reloading compatible with TTY-raw-ers & stdin hostaging (termbox, tcell, ...)
- [ ] Configuration files so you don't have to type the command everytime
- [ ] ???
