# StyledConsole

Helper functions to pretty-print messages in a terminal.

This project is heavily inspired on the awesome php helper [SymfonyStyle from SensioLabs](https://symfony.com/doc/current/console/style.html).
(Sources can be found [here](https://github.com/symfony/console/blob/master/Style/SymfonyStyle.php))

# Installation

Simply run:

    go get github.com/coreoas/styledconsole

# Contrubuting

If you want to open an MR, be sure to run the tests with:

    find . -type f -regex ".*\.go$" | xargs goimports -l
    golint ./...
    go vet ./...
    go test ./...

If you want to run all these tests automatically before every commit, add the custom git-hooks with:

    git config core.hooksPath .githooks

# Methods

## Titling methods

`Title(mytitle string)`

It displays the given string as the command title.
This method is meant to be used only once in a given command, but nothing prevents you to use it repeatedly:

```golang
styledconsole.Title("My Title")
```


---------------------


`Section(mysection string)`

It displays the given string as the title of some command section.
This is only needed in complex commands which want to better separate their contents:

```golang
styledconsole.Section("My section")
```

## Content methods

`Text(sometext string)`

It displays the given string or array of strings as regular text. This is useful to render help messages and instructions for the user running the command:

```golang
styledconsole.Text("Some awesome multi-line\ntext")
```


---------------------


`Listing(myitems []string)`

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

## Admonition Methods

`Note(sometext string)`

It displays the given string or array of strings as a highlighted admonition.
Use this helper sparingly to avoid cluttering command's output:

```golang
styledconsole.Note("Note should be taken about this particular point")
```


---------------------


`Caution(sometext string)`

Similar to the note() helper, but the contents are more prominently highlighted.
The resulting contents resemble an error message, so you should avoid using this helper unless strictly necessary:

```golang
styledconsole.Caution("Wow, be careful about this or that")
```

## Progress Bar Methods

`ProgressStart(nbSteps int)`

It displays a progress bar with a number of steps equal to the argument passed to the method:

```golang
styledconsole.ProgressStart(10)
```


---------------------


`ProgressAdvance(nbSteps int) int`

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

## User Input Methods

`Ask(myquestion string, func(string) bool) string`

`AskWithDefault(myquestion string, defaultAnswer string, func(string) bool) string`

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


`AskHidden(myquestion string, func(string) bool) string`

It's very similar to the `ask(...)` method but the user's input will be hidden.
Use it when asking for sensitive information:

```golang
res := styledconsole.AskHidden("What do you want to hide from your neighbors?", func (res string) bool {
    return len(res) > 0
})
```


---------------------


`Confirm(myquestion string) bool`

`ConfirmWithDefault(myquestion string, defaultAnswer bool) bool`

It asks a Yes/No question to the user and it only returns true or false:

```golang
// Ask a non-determined yes/no question
styledconsole.Confirm("Is this real life?")

// Ask a yes/no question with a default, so we can only hit <ENTER> and move on
styledconsole.ConfirmWithDefault("Is false as true to true while true can be false to false?", false)

```


---------------------


`Choice(myquestion, choices []string)`

`ChoiceWithDefault(myquestion, choices []string, defaultAnswer string)`

It asks a question whose answer is constrained to the given list of valid answers:

```golang
// Ask between some choices
styledconsole.Choice("To be or not to be?", ["To be", "Not to be"])
// Ask between some choices with a default, so we can hit <ENTER> and move on (...and go?)
styledconsole.ChoiceWithDefault("Should I stay or should I go", ["Stay", "Umpf", "Go"], "Umpf")
```

## Result methods

`Success(sometext string)`

It displays the given string or array of strings highlighted as a successful message (with a green background and the [OK] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Success("The atomic bomb model was a dud, you're safe.")
```


---------------------


`Warning(sometext string)`

It displays the given string or array of strings highlighted as a warning message (with a red background and the [WARNING] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Warning("Your backpack just started ticking.")
```


---------------------


`Error(sometext string)`

It displays the given string or array of strings highlighted as an error message (with a red background and the [ERROR] label).
It's meant to be used once to display the final result of executing the given command, but you can use it repeatedly during the execution of the command:

```golang
styledconsole.Error("Someone lit a cigarette in the firework storing area.")
```
