Here is some info about the example programs:

BRAINFRICK.NMS:
A barebones brainf*ck interpreter with a few limitations:
- It works best with an input file in binary read mode (commandline argument -i). 
- It does not support reading or writing past slot 29999 or before slot 0, and doing so will result in unintended behaviour. 
- Cells can be of negative value, and can go past 255. 
- The , command always returns 0.

ECHO.NMS:
A simple program that reads input, then outputs that plus a space, until the value read is -1.