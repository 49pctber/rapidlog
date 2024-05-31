# Rapid Log

A simple command-line tool to create a digital [rapid log](https://bulletjournal.com/blogs/faq/what-is-rapid-logging-understand-rapid-logging-bullets-and-signifiers).

## About

A rapid log is an efficient way to record notes, to-dos, events, feelings, and questions.
Each entry is prepended with one of five symbols according to the following table:

| Symbol | Description |
|--------|-------------|
| - | notes, thoughts, etc. |
| . | to-do list items |
| o | events |
| = | feelings, mindfullness |
| ? | questions |

## Installation

You can install `rapidlog` using the command

```
go install github.com/49pctber/rapidlog@latest
```

## Usage

Open a terminal and type `rapidlog` and press enter.
Type the character corresponding to the entry type you would like to create, and type your entry.
Note that you can input multiple entries at once.
When you're done, simply type `exit` or `quit`.

For example:

```
rapidlog
- wrote a book today
= ecstatic that I'm done
o broke my leg at the gym
. get cast on leg
? how long will my leg take to heal?
quit
```

### Reading Your Entries

You can list all of your entries using `rapidlog list`.
This will list your entries and group them by date.
For more information, you can append the `-v` flag to enable verbose printing.
This will give a timestamp and the entry's ID (e.g. so you can delete entries later if necessary).

If you only want to list certain items (e.g. to-do items prepened with `.`, or questions prepended with `?`), use the command `rapidlog list -e .`.

Alternatively, you can call `rapidlog summary` which will render all of your entries as HTML.
Your default browser will be opened to display your rapid log.

### Editing an Entry

You can edit your entries in your default text editor by runing `rapidlog edit <id>`.
(A typical entry ID looks like `2hCcN9LzWH0whbkF8vzSSdKCVfA`.)
Edit your entry like you would any other text file, save, and close your editor.
The saved version of your entry will replace your old entry.

### Deleting an Entry

You can delete an entry using the comand `rapidlog delete <id>`, where ID can be obtained by either

1. reading the ID from using `rapidlog list -v`
2. clicking on the entry from `summary.html` to copy it to your keyboard

### Importing Another Database File

If you have multiple database files, you can merge them using `rapidlog import <path>` to incorporate your other database into your active one.

### `--help`

If you want more information about how to use this program, run `rapidlog --help` in your terminal.
