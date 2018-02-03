// Generated from D:/dev/go-dev/src/github.com/LogicalOverflow/music-sync/ssh/parser\Shell.g4 by ANTLR 4.7.

package shell_parser // Shell
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 5, 43, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 3,
	2, 3, 2, 3, 2, 7, 2, 18, 10, 2, 12, 2, 14, 2, 21, 11, 2, 3, 3, 3, 3, 3,
	4, 3, 4, 3, 5, 3, 5, 3, 6, 7, 6, 30, 10, 6, 12, 6, 14, 6, 33, 11, 6, 3,
	7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 41, 10, 7, 3, 7, 2, 2, 8, 2, 4,
	6, 8, 10, 12, 2, 2, 2, 41, 2, 14, 3, 2, 2, 2, 4, 22, 3, 2, 2, 2, 6, 24,
	3, 2, 2, 2, 8, 26, 3, 2, 2, 2, 10, 31, 3, 2, 2, 2, 12, 40, 3, 2, 2, 2,
	14, 19, 5, 4, 3, 2, 15, 16, 7, 4, 2, 2, 16, 18, 5, 6, 4, 2, 17, 15, 3,
	2, 2, 2, 18, 21, 3, 2, 2, 2, 19, 17, 3, 2, 2, 2, 19, 20, 3, 2, 2, 2, 20,
	3, 3, 2, 2, 2, 21, 19, 3, 2, 2, 2, 22, 23, 5, 8, 5, 2, 23, 5, 3, 2, 2,
	2, 24, 25, 5, 8, 5, 2, 25, 7, 3, 2, 2, 2, 26, 27, 5, 10, 6, 2, 27, 9, 3,
	2, 2, 2, 28, 30, 5, 12, 7, 2, 29, 28, 3, 2, 2, 2, 30, 33, 3, 2, 2, 2, 31,
	29, 3, 2, 2, 2, 31, 32, 3, 2, 2, 2, 32, 11, 3, 2, 2, 2, 33, 31, 3, 2, 2,
	2, 34, 41, 7, 5, 2, 2, 35, 36, 7, 3, 2, 2, 36, 41, 7, 4, 2, 2, 37, 38,
	7, 3, 2, 2, 38, 41, 7, 3, 2, 2, 39, 41, 7, 3, 2, 2, 40, 34, 3, 2, 2, 2,
	40, 35, 3, 2, 2, 2, 40, 37, 3, 2, 2, 2, 40, 39, 3, 2, 2, 2, 41, 13, 3,
	2, 2, 2, 5, 19, 31, 40,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'\\'", "' '",
}
var symbolicNames = []string{
	"", "ESCAPE_CHARACTER", "SPACE", "NORMAL_CHARACTER",
}

var ruleNames = []string{
	"line", "command", "parameter", "commandString", "rawString", "character",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type ShellParser struct {
	*antlr.BaseParser
}

func NewShellParser(input antlr.TokenStream) *ShellParser {
	this := new(ShellParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Shell.g4"

	return this
}

// ShellParser tokens.
const (
	ShellParserEOF              = antlr.TokenEOF
	ShellParserESCAPE_CHARACTER = 1
	ShellParserSPACE            = 2
	ShellParserNORMAL_CHARACTER = 3
)

// ShellParser rules.
const (
	ShellParserRULE_line          = 0
	ShellParserRULE_command       = 1
	ShellParserRULE_parameter     = 2
	ShellParserRULE_commandString = 3
	ShellParserRULE_rawString     = 4
	ShellParserRULE_character     = 5
)

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) Command() ICommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandContext)
}

func (s *LineContext) AllSPACE() []antlr.TerminalNode {
	return s.GetTokens(ShellParserSPACE)
}

func (s *LineContext) SPACE(i int) antlr.TerminalNode {
	return s.GetToken(ShellParserSPACE, i)
}

func (s *LineContext) AllParameter() []IParameterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IParameterContext)(nil)).Elem())
	var tst = make([]IParameterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IParameterContext)
		}
	}

	return tst
}

func (s *LineContext) Parameter(i int) IParameterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParameterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IParameterContext)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *ShellParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ShellParserRULE_line)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(12)
		p.Command()
	}
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ShellParserSPACE {
		{
			p.SetState(13)
			p.Match(ShellParserSPACE)
		}
		{
			p.SetState(14)
			p.Parameter()
		}

		p.SetState(19)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICommandContext is an interface to support dynamic dispatch.
type ICommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandContext differentiates from other interfaces.
	IsCommandContext()
}

type CommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandContext() *CommandContext {
	var p = new(CommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_command
	return p
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) CommandString() ICommandStringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandStringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandStringContext)
}

func (s *CommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterCommand(s)
	}
}

func (s *CommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitCommand(s)
	}
}

func (p *ShellParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, ShellParserRULE_command)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		p.CommandString()
	}

	return localctx
}

// IParameterContext is an interface to support dynamic dispatch.
type IParameterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParameterContext differentiates from other interfaces.
	IsParameterContext()
}

type ParameterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParameterContext() *ParameterContext {
	var p = new(ParameterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_parameter
	return p
}

func (*ParameterContext) IsParameterContext() {}

func NewParameterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParameterContext {
	var p = new(ParameterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_parameter

	return p
}

func (s *ParameterContext) GetParser() antlr.Parser { return s.parser }

func (s *ParameterContext) CommandString() ICommandStringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandStringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandStringContext)
}

func (s *ParameterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParameterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParameterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterParameter(s)
	}
}

func (s *ParameterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitParameter(s)
	}
}

func (p *ShellParser) Parameter() (localctx IParameterContext) {
	localctx = NewParameterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, ShellParserRULE_parameter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(22)
		p.CommandString()
	}

	return localctx
}

// ICommandStringContext is an interface to support dynamic dispatch.
type ICommandStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandStringContext differentiates from other interfaces.
	IsCommandStringContext()
}

type CommandStringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandStringContext() *CommandStringContext {
	var p = new(CommandStringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_commandString
	return p
}

func (*CommandStringContext) IsCommandStringContext() {}

func NewCommandStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandStringContext {
	var p = new(CommandStringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_commandString

	return p
}

func (s *CommandStringContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandStringContext) RawString() IRawStringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRawStringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRawStringContext)
}

func (s *CommandStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterCommandString(s)
	}
}

func (s *CommandStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitCommandString(s)
	}
}

func (p *ShellParser) CommandString() (localctx ICommandStringContext) {
	localctx = NewCommandStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ShellParserRULE_commandString)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(24)
		p.RawString()
	}

	return localctx
}

// IRawStringContext is an interface to support dynamic dispatch.
type IRawStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRawStringContext differentiates from other interfaces.
	IsRawStringContext()
}

type RawStringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRawStringContext() *RawStringContext {
	var p = new(RawStringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_rawString
	return p
}

func (*RawStringContext) IsRawStringContext() {}

func NewRawStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawStringContext {
	var p = new(RawStringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_rawString

	return p
}

func (s *RawStringContext) GetParser() antlr.Parser { return s.parser }

func (s *RawStringContext) AllCharacter() []ICharacterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICharacterContext)(nil)).Elem())
	var tst = make([]ICharacterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICharacterContext)
		}
	}

	return tst
}

func (s *RawStringContext) Character(i int) ICharacterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICharacterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICharacterContext)
}

func (s *RawStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RawStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RawStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterRawString(s)
	}
}

func (s *RawStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitRawString(s)
	}
}

func (p *ShellParser) RawString() (localctx IRawStringContext) {
	localctx = NewRawStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ShellParserRULE_rawString)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ShellParserESCAPE_CHARACTER || _la == ShellParserNORMAL_CHARACTER {
		{
			p.SetState(26)
			p.Character()
		}

		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICharacterContext is an interface to support dynamic dispatch.
type ICharacterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCharacterContext differentiates from other interfaces.
	IsCharacterContext()
}

type CharacterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCharacterContext() *CharacterContext {
	var p = new(CharacterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ShellParserRULE_character
	return p
}

func (*CharacterContext) IsCharacterContext() {}

func NewCharacterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CharacterContext {
	var p = new(CharacterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ShellParserRULE_character

	return p
}

func (s *CharacterContext) GetParser() antlr.Parser { return s.parser }

func (s *CharacterContext) NORMAL_CHARACTER() antlr.TerminalNode {
	return s.GetToken(ShellParserNORMAL_CHARACTER, 0)
}

func (s *CharacterContext) AllESCAPE_CHARACTER() []antlr.TerminalNode {
	return s.GetTokens(ShellParserESCAPE_CHARACTER)
}

func (s *CharacterContext) ESCAPE_CHARACTER(i int) antlr.TerminalNode {
	return s.GetToken(ShellParserESCAPE_CHARACTER, i)
}

func (s *CharacterContext) SPACE() antlr.TerminalNode {
	return s.GetToken(ShellParserSPACE, 0)
}

func (s *CharacterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CharacterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CharacterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.EnterCharacter(s)
	}
}

func (s *CharacterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ShellListener); ok {
		listenerT.ExitCharacter(s)
	}
}

func (p *ShellParser) Character() (localctx ICharacterContext) {
	localctx = NewCharacterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, ShellParserRULE_character)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(32)
			p.Match(ShellParserNORMAL_CHARACTER)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(33)
			p.Match(ShellParserESCAPE_CHARACTER)
		}
		{
			p.SetState(34)
			p.Match(ShellParserSPACE)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(35)
			p.Match(ShellParserESCAPE_CHARACTER)
		}
		{
			p.SetState(36)
			p.Match(ShellParserESCAPE_CHARACTER)
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(37)
			p.Match(ShellParserESCAPE_CHARACTER)
		}

	}

	return localctx
}
