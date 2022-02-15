# Using the Numskull interpreter
 `Numskull interpreter v0.1.0`

 This document will assume the numskull interpreter can be accessed through a command line by simply typing `numskull`.

 Also note that a majority of this file is copy-pasted from the program itself. use `numskull -h` to get the information you get here.

## Usage
 To execute a Numskull program, open a command line interface and enter the interpreter path. Give the interpreter the different flags you need, and at last give it a path to the Numskull program you wrote.
 <br>
 `numskull [-i file] [-t] [-o file] [-c] <program-file>`

## Available arguments
 There are a couple arguments written into the interpreter:
 ```
 -h, --help <argument>   Prints usage for given argument
 -v, --version           Prints program version number
 -i, --input <path>      File to read input from
 -t, --type              Tells program to read input file as text
 -o, --output <path>     File to print output to
 -c, --console           Force program output to console
 ```

### `-h`, `--help <argument>`
 Prints usage for given argument.

 *Example:* `numskull --help v`
 <br>
 Shows the help page for the parameter [`-v`](#-v---version). The program stops once help has been printed.

### `-v`, `--version`
 Prints program and language version number.
 <br>
 That's it, really.

### `-i`, `--input <path>`
 Read input from a file.
 <br>
 When reading after the end of file, the result is always `-1`.
 <br>
 If this argument isn't present, input is given through the console.
 <br>
 Input file will be read as binary by default. Pass in [`-t`](#-t---type) to prevent this. Look under [Reading / writing data](#reading--writing-data) for more information.

 *Example:* `numskull -i numbers.bin program.nms`
 <br>
 Opens the program `program.nms`, and reads from `numbers.bin` when reading input.

### `-t`, `--type`
 Tells program to read input file as text, rather than as binary.

 Normally when reading input from a file, the file will be read as binary.
 That is, one byte per input, each consisting of a number from 0 to 255.
 Some may not want this behaviour, so passing in `--type` will make input read as text.
 
 Entries are read as numbers, seperated by whitespace (tabs, spaces, or newlines).
 An incorrectly formatted entry will cause an error when trying to read it.
 <br>
 If the [`-i`](#-i---input-path) argument isn't present, this argument does nothing.

### `-o`, `--output <path>`
 Print output to a given file.
 <br>
 When outputting, you can choose to also write that output to a file.
 <br>
 The file is treated as a byte array.

 If this argument isn't present, the output of the program is displayed in the console.
 If you still want console output AND saving to a file, use the [`-c`](#-c---console) argument.
 <br>
 If the program stops due to an error, the file is still saved.

 *Example:* `numskull --output numbers.bin program.nms`
 <br>
 Opens `program.nms`, and saves output to `numbers.bin`, once the program stops running.

### `-c`, `--console`
 Force program output to the console.
 <br>
 When outputting to a file, console output is turned off by default.
 <br>
 Use this argument to reenable it, while also writing the output to a file, using [`-o`](#-o---output-path).
 <br>
 If the [`-o`](#-o---output-path) argument isn't present, this argument does nothing.

## Reading / writing data
 All output is by default treated as a console, which characters can be written to. When writing via the `!` operator, multiple characters are written, and when outputting via the `#` operator, only one character is written.

 By default, this data only goes to the console the interpreter runs in. If a file is passed in using the [`-o`](#-o---output-path) argument, output only goes to the file and not to the console. If you want both, pass in both [`-o`](#-o---output-path) and [`-c`](#-c---console).

 When requesting input using the `"` operator, the interpreter will by default pause execution until the user writes something in the console. To stop this behaviour, pass in [`-i`](#-i---input-path) to make it instantly read input from a file. By default, the input read will be binary, but this behaviour can be changed by passing in [`-t`](#-t---type). This makes the program read numbers as text instead.

 When the end of a file has been reached, the `"` operator will always return `-1`.
 Input files are processed before the program runs, so if a text based input file contains a syntax error, the program won't run.
 
 ---
 
 This project is maintained on [Github](http://github.com/sukus21/numskull).
