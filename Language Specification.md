# Numskull language specification
 `numskull version 1.0`

## History
 The idea for Numskull was concieved the 11th of February 2022.
 <br>
 The idea was to make a programming language, where numbers could be equal to other numbers. The first version of the interpreter `(v0.1.0)` was written in GO over that same weekend `(feb. 11-13 2022)`. Language features, language syntax and interpreter functionality was all decided while writing the interpreter.

## General information
 In this language there are no variables, and letters are forbidden. Only numbers and symbols are allowed. The language's gimmick is the ability to assign a number to be equal to another number. This means that all numbers are variables, containing themselves by default. This means all integer numbers, negative numbers, decimal numbers and decimal negative numbers can store values and are valid values to be stored. All values (and keys) are represented by a 64-bit float. The language only supports real numbers, no NaN's, infinities or others.
 <br>
 There is currently no way to write comments in programs.



## Instructions
 Numskull programs are made up of individual instructions. Each line of code contains one instructions, and required newlines between them. An instruction is structured like this: 
 <br>
 `<lefthand> <operation> [righthand] [bracket]`

 Both lefthand and righthand are numbers, but righthand isn't required for all operations, and brackets are only a necessity for the comparison operator. Below are examples of instructions.
 ```
 Assign a value:             10 = 60
 Decrement a value:          -5++
 Addition plus assignment:   7.56 += 7
 ```

## Operations
 There are only a handful of valid operations in Numskull. These are:

 - ### `--`: Decrement
    Decrements value of lefthand. Usually doesn't work well with decimal numbers due to floating point weirdness.
    <br>
    *Example:* `5++` would add `1` to the value stored in `5`.

 - ### `++`: Increment
    Increments value of lefthand. Usually doesn't work well with decimal numbers due to floating point weirdness.
    <br>*Example:* `6--` would subtract `1` from the value stored in `6`.

 - ### `=`: Assign
    Assigns the value of lefthand to the value of the righthand.
    <br>
    *Example:* `44.2 = -7` sets `44.2` to be equal to the value of `-7`.

 - ### `+=`: Add
    Adds the value of the righthand to the value of the lefthand and stores the result in the lefthand.
    <br>
    *Example:* `5 += 2` adds the values of `5` and `2` together, and stores the result in `5`.


 - ### `-=`: Subtract
    Subtracts value of the righthand from the value of the lefthand and stores the result in the lefthand.

 - ### `*=`: Multiply
    Same as `+=`, but multiplies.

 - ### `/=`: Divide
    Same as `+=`, but divides.

 - ### `!`: Print number
    Outputs the number stored in the lefthand as a string.
    <br>
    *Example:* `17!` will output the string "`17`".

 - ### `#`: Print character
    Outputs the number stored in the lefthand as a unicode character.
    <br>
    *Example:* `32#` will output a space character. Look up an ascii/unicode table for the characters you wish to print.

 - ### `"`: Read input
    Reads a number from the input and stored it in the lefthand. The number read can be either from a file, either as binary or text, or can be input via the commandline.
    <br>
    *Example:* `-60"` reads a value and stores it in `-60`.

 - ### `?x`: Comparison
    Requires a bracket at the end of the instruction. Tests the value of the lefthand compared to the value of the righthand. If the condition is not met, the program jumps to the next closing bracket of matching type at the same level and continues executing from there. Further documented in the next chapter.

## Conditions
 Numskull supports branching using the comparison operator, as shown above. The different kinds of comparisons possible are the following:

 - ### `?=`: Equals
    Check if values are equal. Is true if the value of the lefthand is equal to the value of the righthand.

 - ### `?!`: Not equals
    Check if values are different. Same as `?=`, but is true if the values AREN'T the same.

 - ### `?>`: Greater than
    Is true if the value of the lefthand is greater than the value of the righthand.

 - ### `?>=`: Greater than or equal to
    Is true if the value of the lefthand is greater than or equal to the value of the righthand.
  
 - ### `?<`: Lesser than
    Is true if the value of the lefthand is lesser than the value of the righthand.

 - ### `?<=`: Lesser than or equal to
    Is true if the value of the lefthand is lesser than or equal to the value of the righthand.
 
 A comparison operator requires an opening bracket after the righthand operator. The bracket should be on the same line as the instruction, but it is not required. If a condition isn't met, the program skips ahead to the next closing bracket, and continues execution from there. In that way, the comparison operator acts like an equivelant to the `if`-statement in C. Closing brackets should always be on their own line, separate from any instructions.
 
 Below are two example programs and their outputs, to hopefully demonstrate how the comparison instruction works.
 ```c
  __________________     __________________
 | Program one:     |   | Program two:     |
 |__________________|   |__________________|
 |                  |   |                  |
 | 10 ?! 0 {        |   | 10 ?< 5 {        |
 |     10 = 60      |   |     10 = 40      |
 |     10!          |   |     10!          |
 |     10!          |   |     10!          |
 |     10!          |   |     10!          |
 | }                |   | }                |
 | 20!              |   | 20!              |
 |__________________|   |__________________|
 | Output:          |   | Output:          |
 | 606060           |   | 20               |
 |__________________|   |__________________|
 ```
 Loops can be constructed using the comparison operator instead, by simple using square brackets `[]` instead of curly brackets `{}`. When a closing square bracket is encountered ( `]` ), the program skips back up to the matching opening brackets condition statement, and continues from there. If the condition is still true, the loop is run again. Otherwise the program skips to the closing bracket and continues from there. Using curly brackets to open and curly brackets to close (or the other way around) is not permitted.

 Below is an example program and its output, to demonstrate how a loop works.
 ```c
  __________________
 | Example program: |
 |__________________|
 |                  |
 | 1 = 10           |
 | 1 ?> 5 [         |
 |     1!           |
 |     32#          |
 |     1--          |
 | ]                |
 |__________________|
 | Output:          |
 | 10 9 8 7 6       |
 |__________________|
 ```


## Lefthand chaining
 The lefthand operator supports chaining using the `+` and `-` signs, but the way this works might not be obvious. 
 
 The leftmost value is the base, and is read as an immediate value, aka the number itself, not the value within. Every link in the chain thereafter is not immediate, and the value added or subtracted will instead be the value of said number. Whitespace between each link and number is optional. Lefthand chaining is available for all operation types. The righthand cannot be chained, and must always be just one number.
 <br>
 `<base> [+/- offset] [+/- offset]...`

 It is easier to illustrate using examples:
 - `5+8+6` sets the lefthand to be equal to `5` plus the value of `8` plus the value of `6`.

 - `5.5 - -7` sets the lefthand to be equal to `5.5` minus the value of `-7`.


 Below is an example program to better explain lefthand chaining.
 ```c
  __________________
 | Example program: |
 |__________________|
 |                  |
 | 1 = 10           |
 | 6+1!             |
 | 32#              |
 | 6+1+7!           |
 |__________________|
 | Output:          |
 | 16 23            |
 |__________________|
 ```

## More information
 The interpreter is commandline based, and will likely remain so.
 For more information on how it works and what options are available, run the interpreter from the commandline with `--help`.

 The original Numskull interpreter is open source, and can be found on [Github](http://github.com/sukus21/numskull).