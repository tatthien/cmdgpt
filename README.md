# cmdgpt

`cmdgpt` is a command line tool that use ChatGPT to generate Linux commands from your natural language.

## Features

- [x] Generate commands based on your input. For example, `show git logs with author and message` => `git log --pretty=format:"%an - %s"`
- [x] Allow you to `run` or `copy` the reponse command
- [ ] Auto fix commands
- [ ] Check if the response commands are safe/risk to run

## Prerequisites

You need to have `OPENAI_API_KEY` environemnt variable configured.

```
export OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Installation

### Using go

```
go install github.com/tatthien/cmdgpt@latest
```

### Pre-built binaries

You can download the pre-built binaries for macOS, Linux, and Window from the [releases page](https://github.com/tatthien/cmdgpt/releases)

## Usage

Simply run `cmdgpt` and ask for a command. You can try to copy one of these following commands to ask:

- search for a string within an output
- create a post request with curl
- list all file types with count
- create a directory and its parent if no exist

[![asciicast](https://asciinema.org/a/3nFuZGFrsDXcRl7XBvuOfstdU.svg)](https://asciinema.org/a/3nFuZGFrsDXcRl7XBvuOfstdU)
