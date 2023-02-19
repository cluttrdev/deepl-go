# DeepL Go Library

This is an **unofficial** client library and command line interface for the [DeepL API][api-docs].

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

Import the package and construct a `Translator`.

```go
import (
    "fmt"
    "log"
    
    deepl "github.com/cluttrdev/deepl-go/pkg/api"
)

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
```

## Command Line Interface

### Installation

Using the Go tools:

```shell
go install github.com/cluttrdev/deepl-go/cmd/deepl
```

### Usage

Currently, in order to use the command line interface the authentication key must be set as an environment
variable.

```shell
$ export DEEPL_AUTH_KEY="f63c02c5-f056..."
$ deepl translate "Hello, world!" --target-lang FR
Bonjour, le monde !
```

To get an overview of the available commands, run `deepl help`.

#### Autocompletion

You can use the `deepl completion` command to generate autocompletion scripts for the cli.


[api-docs]: https://www.deepl.com/docs-api?utm_source=github

[create-account]: https://www.deepl.com/pro?utm_source=github#developer
