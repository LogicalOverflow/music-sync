grammar Shell;

line: command (SPACE parameter)*;

command: commandString;

parameter: commandString;

commandString: rawString;

rawString: character*;

character: NORMAL_CHARACTER | ESCAPE_CHARACTER SPACE | ESCAPE_CHARACTER ESCAPE_CHARACTER | ESCAPE_CHARACTER;


ESCAPE_CHARACTER: '\\';
SPACE: ' ';
NORMAL_CHARACTER: .;
