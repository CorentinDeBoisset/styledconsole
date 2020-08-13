# StyledConsole

Helper functions to pretty-print messages in a terminal.

This project is heavily inspired on the awesome php helper [SymfonyStyle from SensioLabs](https://symfony.com/doc/current/console/style.html).

# Installation

Simply run:

    go get github.com/corentindeboisset/styledconsole

# Contrubuting

If you want to open an MR, be sure to run the tests with:

    find . -type f -regex ".*\.go$" | xargs goimports -l
    golint ./...
    go vet ./...
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
styledconsole.Title("My Title")
```


---------------------


`Section(title string)`

It displays the given string as the title of some command section.
This is only needed in complex commands which want to better separate their contents:

```golang
styledconsole.Section("My section")
```

### Content methods

`Text(content string)`

It displays the given string or array of strings as regular text. This is useful to render help messages and instructions for the user running the command:

```golang
styledconsole.Text("Some awesome multi-line\ntext")
```


---------------------


`Listing(items []string)`

It displays an unordered list of elements passed as an array:

```golang
styledconsole.Listing(["item A", "item B", "item δ"])
```


---------------------


`Table(headers []string, rows [][]string)`

It displays the given array of headers and rows as a compact table:

```golang
styledconsole.Table(["header1", "header2"], [["cell 1-1", "cell 1-2"], ["cell 2-1", "cell 2-2"]])
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

### Admonition Methods

`Note(content string)`

It displays the given string or array of strings as a highlighted admonition.
Use this helper sparingly to avoid cluttering command's output:

```golang
styledconsole.Note("Note should be taken about this particular point")
```


---------------------


`Caution(content string)`

Similar to the note() helper, but the contents are more prominently highlighted.
The resulting contents resemble an error message, so you should avoid using this helper unless strictly necessary:

```golang
styledconsole.Caution("Wow, be careful about this or that")
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

`Ask(myquestion string, validator func(string) bool) string`

`AskWithDefault(myquestion string, defaultAnswer string, validator func(string) bool) string`

It asks the user to provide some value, with or without a default value.
The callback is used to validate the input:

```golang
// Ask some question with no pre-defined answer
styledconsole.Ask("What did you forget to do?", func (res string) bool {
    return res != "Nothing I swear"
})

// Ask a question with a default, so we can only hit <ENTER> and move on
styledconsole.AskWithDefault("Where is the sea?", "All around, this is an island", func (res srting) bool {
    return res != "I don't know"
})
```


---------------------


`AskHidden(question string, validator func(string) bool) string`

It's very similar to the `ask(...)` method but the user's input will be hidden.
Use it when asking for sensitive information:

```golang
res := styledconsole.AskHidden("What do you want to hide from your neighbors?", func (res string) bool {
    return len(res) > 0
})
```


---------------------


`Confirm(question string) bool`

`ConfirmWithDefault(question string, defaultAnswer bool) bool`

It asks a Yes/No question to the user and it only returns true or false:

```golang
// Ask a non-determined yes/no question
styledconsole.Confirm("Is this real life?")

// Ask a yes/no question with a default, so we can only hit <ENTER> and move on
styledconsole.ConfirmWithDefault("Is false as true to true while true can be false to false?", false)

```


---------------------


`Choice(question string, choices []string) int`

`ChoiceWithDefault(question string, choices []string, defaultAnswer int) int`

It asks a question whose answer is constrained to the given list of valid answers.
The value returned is the index of the selected answer:

```golang
// Ask between some choices
styledconsole.Choice("To be or not to be?", ["To be", "Not to be"])
// Ask between some choices with a default, so we can hit <ENTER> and move on (...and go?)
styledconsole.ChoiceWithDefault("Should I stay or should I go", ["Stay", "Go"], 1)
```

### Result methods

`Success(content string)`

It displays the given string or array of strings highlighted as a successful message (with a green background and the [OK] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Success("The atomic bomb model was a dud, you're safe.")
```


---------------------


`Warning(content string)`

It displays the given string or array of strings highlighted as a warning message (with a red background and the [WARNING] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Warning("Your backpack just started ticking.")
```


---------------------


`Error(content string)`

It displays the given string or array of strings highlighted as an error message (with a red background and the [ERROR] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Error("Someone lit a cigarette in the firework storing area.")
```
