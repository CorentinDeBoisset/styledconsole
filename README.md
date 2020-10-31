# StyledConsole

Helper functions to pretty-print messages in a terminal.

This project is heavily inspired on the awesome php helper [SymfonyStyle from SensioLabs](https://symfony.com/doc/current/console/style.html).

# Installation

Simply run:

    go get github.com/corentindeboisset/styledconsole

# Contrubuting

If you want to open an MR, be sure to run the tests with:

    golangci-lint run
    go test ./...

If you want to run all these tests automatically before every commit, add the custom git-hooks with:

    git config core.hooksPath .githooks

# Usage

## Styling tags

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


## Methods

### Titling methods

`Title(mytitle string)`

It displays the given string as the command title.
This method is meant to be used only once in a given command, but nothing prevents you to use it repeatedly:

```golang
styledconsole.Title("Some Title")
```


---------------------


`Section(title string)`

It displays the given string as the title of some command section.
This is only needed in complex commands which want to better separate their contents:

```golang
styledconsole.Section("Some section")
```

### Content methods

`Text(content string)`

It displays the given string as regular text. This is useful to render help messages and instructions for the user running the command:

```golang
styledconsole.Text("Some \ntext\nthat can even handle <fg=red>multi-line\nstyling</>")
```


---------------------


`Listing(items []string)`

It displays an unordered list of elements passed as an array:

```golang
styledconsole.Listing([]string{"item A", "item B", "item δ"})
```


---------------------


`Table(headers []string, rows [][]string)`

It displays the given array of headers and rows as a compact table:

```golang
styledconsole.Table([]string{"header 1", "header 2"}, [][]string{[]string{"cell 1-1", "cell 1-2"}, []string{"cell 2-1", "cell 2-2"}})
```

`NewLine()`

`NewLines(newLineCount int)`

It displays one or more blank lines in the command output.
Although it may seem useful, most of the times you won't need it at all.
The reason is that every helper already adds their own blank lines, so you don't have to care about the vertical spacing:

```golang
// Prints a new line
styledconsole.NewLine()

// Prints 10 new lines
styledconsole.NewLines(10)
```

### Progress Bar Methods

`ProgressStart(totalSteps int)`

It displays a progress bar with a number of steps equal to the argument passed to the method:

```golang
styledconsole.ProgressStart(10)
```


---------------------


`ProgressAdvance(stepCount int) int`

It makes the progress bar advance the given number of steps.
Returns the number of steps left:

```golang
styledconsole.ProgressAdvance(3)
```


---------------------


`ProgressFinish()`

It finishes the progress bar:

```golang
styledconsole.ProgressFinish()
```

### User Input Methods

Note that in this section, when there is a validator function as an argument, it is possible to give `nil` to accept any given answer.

`Ask(myquestion string, validator func(string) bool) string`

`AskWithDefault(myquestion string, defaultAnswer string, validator func(string) bool) string`

It asks the user to provide some value, with or without a default value.
The callback is used to validate the input:

```golang
// Ask some question with no pre-defined answer
styledconsole.Ask("Do you have anything (more than 10char) to say?", func (res string) bool {
    return len(res) > 10
})

// Ask a question with a default, so we can only hit <ENTER> and move on
styledconsole.AskWithDefault("What is your favourite colour?", "Blue. No, yel-- auuuuuuuugh! ", nil)
```

---------------------


`AskHidden(question string, validator func(string) bool) string`

It's very similar to the `ask(...)` method but the user's input will be hidden.
Use it when asking for sensitive information:

```golang
res := styledconsole.AskHidden("What is your password?", func (res string) bool {
    return len(res) > 0
})
```


---------------------


`Confirm(question string) bool`

`ConfirmWithDefault(question string, defaultAnswer bool) bool`

It asks a Yes/No question to the user and it only returns true or false:

```golang
// Ask a non-determined yes/no question
styledconsole.Confirm("You are about to wipe your computer, are you sure?")

// Ask a yes/no question with a default, so we can only hit <ENTER> and move on
styledconsole.ConfirmWithDefault("Should we continue?", true)
```


---------------------


`Choice(question string, choices []string) int`

`ChoiceWithDefault(question string, choices []string, defaultAnswer int) int`

It asks a question whose answer is constrained to the given list of valid answers.
The value returned is the index of the selected answer:

```golang
// Ask between some choices
styledconsole.Choice("To be or not to be?", []string{"To be", "Not to be"})
// Ask between some choices with a default, so we can hit <ENTER> and move on
styledconsole.ChoiceWithDefault("Should I stay or should I go", []string{"Stay", "Go"}, 1)
```

### Result methods

`Success(content string)`

It displays the given string or array of strings highlighted as a successful message (with a green background and the [OK] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Success("Everything is fine.")
```


---------------------


`Warning(content string)`

It displays the given string or array of strings highlighted as a warning message (with a red background and the [WARNING] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Warning("You should worry about this.")
```


---------------------


`Error(content string)`

It displays the given string or array of strings highlighted as an error message (with a red background and the [ERROR] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Error("There was an error.")
```
