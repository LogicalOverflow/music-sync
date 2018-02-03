// Generated from D:/dev/go-dev/src/github.com/LogicalOverflow/music-sync/ssh/parser\Shell.g4 by ANTLR 4.7.

package shell_parser // Shell
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ShellListener is a complete listener for a parse tree produced by ShellParser.
type ShellListener interface {
	antlr.ParseTreeListener

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// EnterCommand is called when entering the command production.
	EnterCommand(c *CommandContext)

	// EnterParameter is called when entering the parameter production.
	EnterParameter(c *ParameterContext)

	// EnterCommandString is called when entering the commandString production.
	EnterCommandString(c *CommandStringContext)

	// EnterRawString is called when entering the rawString production.
	EnterRawString(c *RawStringContext)

	// EnterCharacter is called when entering the character production.
	EnterCharacter(c *CharacterContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)

	// ExitCommand is called when exiting the command production.
	ExitCommand(c *CommandContext)

	// ExitParameter is called when exiting the parameter production.
	ExitParameter(c *ParameterContext)

	// ExitCommandString is called when exiting the commandString production.
	ExitCommandString(c *CommandStringContext)

	// ExitRawString is called when exiting the rawString production.
	ExitRawString(c *RawStringContext)

	// ExitCharacter is called when exiting the character production.
	ExitCharacter(c *CharacterContext)
}
