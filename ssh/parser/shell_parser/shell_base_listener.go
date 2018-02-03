// Generated from D:/dev/go-dev/src/github.com/LogicalOverflow/music-sync/ssh/parser\Shell.g4 by ANTLR 4.7.

package shell_parser // Shell
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseShellListener is a complete listener for a parse tree produced by ShellParser.
type BaseShellListener struct{}

var _ ShellListener = &BaseShellListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseShellListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseShellListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseShellListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseShellListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterLine is called when production line is entered.
func (s *BaseShellListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *BaseShellListener) ExitLine(ctx *LineContext) {}

// EnterCommand is called when production command is entered.
func (s *BaseShellListener) EnterCommand(ctx *CommandContext) {}

// ExitCommand is called when production command is exited.
func (s *BaseShellListener) ExitCommand(ctx *CommandContext) {}

// EnterParameter is called when production parameter is entered.
func (s *BaseShellListener) EnterParameter(ctx *ParameterContext) {}

// ExitParameter is called when production parameter is exited.
func (s *BaseShellListener) ExitParameter(ctx *ParameterContext) {}

// EnterCommandString is called when production commandString is entered.
func (s *BaseShellListener) EnterCommandString(ctx *CommandStringContext) {}

// ExitCommandString is called when production commandString is exited.
func (s *BaseShellListener) ExitCommandString(ctx *CommandStringContext) {}

// EnterRawString is called when production rawString is entered.
func (s *BaseShellListener) EnterRawString(ctx *RawStringContext) {}

// ExitRawString is called when production rawString is exited.
func (s *BaseShellListener) ExitRawString(ctx *RawStringContext) {}

// EnterCharacter is called when production character is entered.
func (s *BaseShellListener) EnterCharacter(ctx *CharacterContext) {}

// ExitCharacter is called when production character is exited.
func (s *BaseShellListener) ExitCharacter(ctx *CharacterContext) {}
