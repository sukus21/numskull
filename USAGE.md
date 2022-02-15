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
 Shows the help page for the parameter `-v`. The program stops once help has been printed.

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
 Input file will be read as binary by default. Look under [File IO](#reading--writing-data) for more information.

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
 If the `-i` argument isn't present, this argument does nothing.

### `-o`, `--output <path>`
 Print output to a given file.
 <br>
 When outputting, you can choose to also write that output to a file.
 <br>
 The file is treated as a byte array.

 If this argument isn't present, the output of the program is displayed in the console.
 If you still want console output AND saving to a file, use the `-c` argument.
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
 Use this argument to reenable it, while also writing the output to a file, using `-o`.
 <br>
 If the `-o` argument isn't present, this argument does nothing.

## Reading / writing data
 All output is by default treated as a console, which characters can be written to. When writing via the `!` operator, multiple characters are written, and when outputting via the `#` operator, only one character is written.

 By default, this data only goes to the console the interpreter runs in. If a file is passed in using the [`-o`](#-o---output-path) argument, 