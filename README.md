# Crow

A file watcher script: given an dir, it will

1. Wait for a file to be created, written, and chmodded (usually means the file is done being written)
2.  Return the file in question to Stdout

I wrote it so that I can automate some of the repeated tasks in testing.

I wrote it and used in OSX, but the go library *should* be os-indenpendent.  

## Installation

```bash
% go get github.com/ericychoi/crow
```

## Usage

```bash
% crow /path/to/incoming
```

The way I use it (for email processing testing for instance)

```bash
% some_script_that_sends_file && FILE=$(crow /path/to/incoming | sed s/\.idx//) && mv "$FILE" "$FILE.eml" && open "$FILE.eml"
```

This is a one-liner that will process a file, create an MIME file in /path/to/incoming, and open the file with default email viewer.
