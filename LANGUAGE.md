# Numskull language specification
 `numskull version 1.1`

## History
 The idea for Numskull was concieved the 11th of February 2022.
 
 The idea was to make a programming language, where numbers weren't constant, but mutable. The first version of the interpreter was written over that same weekend `(feb. 11-13 2022)`. Language features and syntax was decided as I was implementing them, and were not planned out from the start.

## General information
 In this language there are no variables, and letters are forbidden. Only numbers and symbols are allowed. The language's gimmick is the ability to assign a number to be equal to another number. This means that all numbers are variables, containing themselves by default. This means all integer numbers, negative numbers, decimal numbers and decimal negative numbers can store values and are valid values to be stored. 
 
 All values (and keys) are represented by a 64-bit float. The language only supports real numbers as input, but NaN and infinity values can be created during runtime (doing this is strongly discouraged, unless you know what you're doing).



## Instructions
 Numskull programs are made up of individual instructions. Each line of code contains one instructions, and required newlines between them. An instruction is structured like this: 
 
 `<lefthand> <operation> [righthand] [bracket]`

 Both lefthand and righthand are numbers, but righthand isn't required for all operations. Brackets are only a necessity for the comparison operations. Below are a few examples of some instructions:
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
    <br>*Example:* `44.2 = -7` sets `44.2` to be equal to the value of `-7`.

 - ### `+=`: Add
    Adds the value of the righthand to the value of the lefthand and stores the result in the lefthand.
    <br>*Example:* `5 += 2` adds the values of `5` and `2` together, and stores the result in `5`.


 - ### `-=`: Subtract
    Subtracts value of the righthand from the value of the lefthand and stores the result in the lefthand.

 - ### `*=`: Multiply
    Same as `+=`, but multiplies.

 - ### `/=`: Divide
    Same as `+=`, but divides.

 - ### `!`: Print number
    Outputs the number stored in the lefthand as a string.
    <br>*Example:* `17!` will output the string "`17`".

 - ### `#`: Print character
    Outputs the number stored in the lefthand as a unicode character.
    <br>*Example:* `32#` will output a space character. Look up an ascii/unicode table for the characters you wish to print.

 - ### `"`: Read input
    Reads a number from the input and stored it in the lefthand. The number read can be either from a file, either as binary or text, or can be input via the commandline.
    <br>*Example:* `-60"` reads a value and stores it in `-60`.

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
 
 A comparison operator requires an opening bracket after the righthand operator. The bracket should be on the same line as the condition instruction. If a condition isn't met, the program skips ahead to the next closing bracket of the same type at the same depth, and continues execution from there. In that way, the comparison operator acts like an equivelant to the `if`-statement in C. Closing brackets should always be on their own line, separate from any instructions.
 
 Below are two example programs and their outputs, to hopefully demonstrate how the comparison instruction works:
 ```c
 //Example program 1
 10 ?= 0 {    //Is 10 equal to 0?
     10 = 60  //Set 10 to 60
     10!      //Print value of 10
     10!
     10!
 }            //End of if-statement
 20!          //Print value of 20
 ```
 Output:
 ```
 60606020
 ```
 ```c
 //Example program 2
 10 ?< 5 {    //Is 10 below 5?
     10 = 40  //Set 10 to 40
     10!      //Print value of 10
     10!
     10!
 }            //End of if-statement
 20!          //Print value of 20
 ```
 Output:
 ```
 20
 ```
 Loops can be constructed using the comparison operator instead, by simple using square brackets `[]` instead of curly brackets `{}`. When a closing square bracket is encountered ( `]` ), the program skips back up to the matching opening brackets condition statement, and continues from there. If the condition is still true, the loop is run again. Otherwise the program skips to the closing bracket and continues from there. Brackets only close each other, meaning a `]` is required to close a `[`. The same applies to `}` and `{`.

 Below is an example program and its output, to demonstrate how a loop works:
 ```c
 1 = 10     //Set 1 to 10
 1 ?> 5 [   //Is 1 greater than 5?
     1!     //Print contents of 1
     32#    //Print a space
     1--    //Decrement 1
 ]
 ```
 Output:
 ```
 10 9 8 7 6
 ```


## Lefthand chaining
 The lefthand operator supports chaining using the `+` and `-` signs, but the way this works might not be obvious. 
 
 The leftmost value is the base, and is read as an immediate value, and not the value contained in it. Every link in the chain thereafter is not immediate, and the value added or subtracted will instead be the value of said number. Therefore it can be useful to keep 0 free, and use that as a base, if you don't care for a base offset.
 
 `<base> [+/- offset] [+/- offset]...`

 It is easier to illustrate using examples:
 - `5+8+6` sets the lefthand to be equal to 5 plus the value of 8 plus the value of 6.

 - `5.5 - -7` sets the lefthand to be equal to 5.5 minus the value of -7. (Important: When chaining by subtraction, whitespace is always reqired between the chain operator and the number itself.)


 Below is an example program to better explain lefthand chaining:
 ```c
 1 = 10  //Set 1 to 10
 6+1!    //Print value at (6+10) = 16 (1 contains 10)
 32#     //Print space
 6+1+7!  //Print number at (6+10+7) = 23
 ```
 Output:
 ```
 16 23
 ```

 The righthand cannot be chained, and must always be just one number.

## Other language features
 - As of version 1.1, Numskull supports code comments. A comment can be started using `//`, and anything that comes after it on the same line will be ignored. Comments can also be multilined. These are started with `/*` and terminated with `*/`. Anything between the two will be ignored at runtime.

 - The interpreter does not care if brackets are opened and closed out of pairs. This means you could use brackets in the following order: `{ [ } ]`. If the first `{` is taken, it jumps to the `}`, inside the body of the `[ ]`, and execution continues as normal from there, and when the `]` is encountered, the condition at the `[` is evaluated as normal. The same technique could be used to exit the body of the `[ ]` early.


## More information
 The interpreter is commandline based, and will likely remain so.
 For more information on how it works and what options are available, run the interpreter from the commandline with `--help`, or check the [usage document](http://github.com/sukus21/numskull).
 
 ---
 
 This project is maintained on [Github](http://github.com/sukus21/numskull).
