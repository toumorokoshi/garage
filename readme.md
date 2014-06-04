# Garage

Garage is a shell command management tool. It includes:

* package management of shell commands
* multi-level namespacing of commands
* smart completion helper to find the commands that do what you want

## Smart completion

If you type:

    garage :find

This will bring up an ncurses-like interface with a blank text prompt. Garage looks through all existing commands,
and looks for strings (on a single line) that match the format:

    gfind: <message> ; <command_arguments>

The text typed by the user will match against <message>, and if
executed will execute the commands as dictated in <command_arguments>.

You can choose not to add a '; <command_arguments>'. This results in
the command not being runnable, but will still appear in the
completion.

For example, let's say we want to create a gfind entry for grep. A couple examples would be:

    # no command_arguments, will complete but will not be executable
    # gfind: search for specific lines of text for files in a folder

    # with command_arguments, will complete and be executable
    # gfind: search for the text "foo" recursively in the current directory; -r 'foo' *
