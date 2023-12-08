# DeepL Go Library

This is an **unofficial** client library and command line interface for the
[DeepL API][api-docs].

## Getting an authentication ley

To use the DeepL Go Library, you'll need an API authentication key. To get a
key, [please create an account here][create-account]. With a DeepL API Free
account you can translate up to 500,000 characters/month for free.

## Library

### Installation

Using the Go tools, from inside your project:

```shell
go get github.com/cluttrdev/deepl-go
```

### Usage

Import the package and create a `Translator`.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cluttrdev/deepl-go/deepl"
)

func main() {
    authKey := "f63c02c5-f056-..."  // Replace with your key

    translator, err := deepl.NewTranslator(authKey)
    if err != nil {
        log.Fatal(err)
    }

    translations, err := translator.TranslateText([]string{"Hello, world!"}, "FR")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(translations[0].Text)  // "Bonjour, le monde !"
}
```

## Command Line Interface

### Installation
 
To install the `deepl` cli you can download a [prebuilt binary][releases] that
matches your system and place it in a directory that's part of your system's
search path, e.g.

```shell
# system information
OS=linux
ARCH=amd64

# install dir (must exist)
BIN_DIR=$HOME/.local/bin

# download latest release
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/cluttrdev/deepl-go/releases/latest | jq -r '.tag_name')
curl \
    -L https://github.com/cluttrdev/deepl-go/releases/download/${RELEASE_TAG}/deepl_${RELEASE_TAG}_${OS}_${ARCH}.tar.gz \
    -o deepl.tar.gz

# create install dir (if necessary)
BIN_DIR=$HOME/.local/bin  # assuming this is part of your system's $PATH
mkdir -p ${BIN_DIR}

# install
tar -C ${BIN_DIR} -zxf deepl.tar.gz deepl
```

Alternatively, you can use the Go tools:

```shell
go install github.com/cluttrdev/deepl-go/cmd/deepl@latest
```

### Usage

To use the cli the authentication key must either be set as an environment
variable or passed via the `--auth-key` option.

```shell
$ export DEEPL_AUTH_KEY="f63c02c5-f056..."
$ deepl translate --to FR "Hello, World!"
Bonjour, le monde !
```

To get an overview of the available commands, run `deepl --help`.

<!-- Links -->
[releases]: https://github.com/cluttrdev/deepl-go/releases/latest
[api-docs]: https://www.deepl.com/docs-api
[create-account]: https://www.deepl.com/pro#developer
