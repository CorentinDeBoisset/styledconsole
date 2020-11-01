# StyledConsole

💅⚡️🖥

Make your GUI tools user-friendly with methods to pretty-print text, lists, tables, user prompts, progress bars...

This project is heavily inspired on the awesome php helper [SymfonyStyle from SensioLabs](https://symfony.com/doc/current/console/style.html).


## ❯ Install

Simply run:

    go get github.com/corentindeboisset/styledconsole

## ❯ Usage

<img src="https://github.com/corentindeboisset/styledconsole/raw/master/assets/demo.svg" alt="example of styledconsole usage in a terminal" />

See the full usage documentation [here](https://pkg.go.dev/github.com/corentindeboisset/styledconsole).

## ❯ About styling tags

Any text can be augmented with style, by enclosing the text with tags like this:

    <fg=blue>Blue foreground text</fg=blue>


It is possible to imbricate the tags like this:

    <fg=blue>Text with blue foreground, <bg=red>text with blue foreground and red background</bg=red>, text with blue foreground again</fg=blue>

If you want to, you can use the `</>` shorthand to close the last opened tag. The previous text thus becomes:

    <fg=blue>Text with blue foreground, <bg=red>text with blue foreground and red background</>, text with blue foreground again</>

One can set multiple properties in a single tag, by separating them with `;`:

    <fg=blue;bg=red>text with blue fg and red bg</>

The available properties are the following:

* `fg=color`: sets the foreground color (value must be in the set below).
* `bg=color`: sets the background color (the value must be in the set below).
* `href=http://link/to/resource`: adds a hypertext link to the given location.
* `options=opt1,opt2,opt3`: Adds additionnal text decorations. Available options are `bold`, `underscore`, `blink`, `reverse` and `conceal`.

The available colors are the standard ANSI set:

* black
* red
* green
* yellow
* blue
* magenta
* cyan
* white
* default (use the terminal's default)

## ❯ Contrubuting

If you want to open an MR, be sure to run the tests with:

    golangci-lint run
    go test ./...

If you want to run all these tests automatically before every commit, add the custom git-hooks with:

    git config core.hooksPath .githooks
