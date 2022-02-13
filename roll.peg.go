package main

// Code generated by peg -inline roll.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulee
	rulee1
	rulee2
	rulee3
	ruleeDice
	rulee4
	rulevalue
	rulenumber
	ruleidentifier
	rulesub
	ruleadd
	ruleminus
	rulemultiply
	ruledivide
	rulemodulus
	ruleexponentiation
	ruleopen
	ruleclose
	rulesp
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	rulePegText
)

var rul3s = [...]string{
	"Unknown",
	"e",
	"e1",
	"e2",
	"e3",
	"eDice",
	"e4",
	"value",
	"number",
	"identifier",
	"sub",
	"add",
	"minus",
	"multiply",
	"divide",
	"modulus",
	"exponentiation",
	"open",
	"close",
	"sp",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"PegText",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type DiceRollParser struct {
	RollExpression

	Buffer string
	buffer []rune
	rules  [33]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *DiceRollParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *DiceRollParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *DiceRollParser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *DiceRollParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *DiceRollParser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *DiceRollParser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func (p *DiceRollParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddOperator(TypeHalt)
		case ruleAction1:
			p.AddOperator(TypeAdd)
		case ruleAction2:
			p.AddOperator(TypeSubtract)
		case ruleAction3:
			p.AddOperator(TypeMultiply)
		case ruleAction4:
			p.AddOperator(TypeDivide)
		case ruleAction5:
			p.AddOperator(TypeModulus)
		case ruleAction6:
			p.AddOperator(TypeExponentiation)
		case ruleAction7:
			p.AddOperator(TypeDice)
		case ruleAction8:
			p.AddValue("1")
			p.AddOperator(TypeSwap)
			p.AddOperator(TypeDice)
		case ruleAction9:
			p.AddOperator(TypeNegation)
		case ruleAction10:
			p.AddValue(string(text))
		case ruleAction11:
			p.AddLoadVarname(string(text))

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func Pretty(pretty bool) func(*DiceRollParser) error {
	return func(p *DiceRollParser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*DiceRollParser) error {
	return func(p *DiceRollParser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *DiceRollParser) Init(options ...func(*DiceRollParser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 e <- <(sp e1 !. Action0)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulesp]() {
					goto l0
				}
				if !_rules[rulee1]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				{
					add(ruleAction0, position)
				}
				add(rulee, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 e1 <- <(e2 ((add e2 Action1) / (minus e2 Action2))*)> */
		func() bool {
			position4, tokenIndex4 := position, tokenIndex
			{
				position5 := position
				if !_rules[rulee2]() {
					goto l4
				}
			l6:
				{
					position7, tokenIndex7 := position, tokenIndex
					{
						position8, tokenIndex8 := position, tokenIndex
						{
							position10 := position
							if buffer[position] != rune('+') {
								goto l9
							}
							position++
							if !_rules[rulesp]() {
								goto l9
							}
							add(ruleadd, position10)
						}
						if !_rules[rulee2]() {
							goto l9
						}
						{
							add(ruleAction1, position)
						}
						goto l8
					l9:
						position, tokenIndex = position8, tokenIndex8
						if !_rules[ruleminus]() {
							goto l7
						}
						if !_rules[rulee2]() {
							goto l7
						}
						{
							add(ruleAction2, position)
						}
					}
				l8:
					goto l6
				l7:
					position, tokenIndex = position7, tokenIndex7
				}
				add(rulee1, position5)
			}
			return true
		l4:
			position, tokenIndex = position4, tokenIndex4
			return false
		},
		/* 2 e2 <- <(e3 ((multiply e3 Action3) / (divide e3 Action4) / (modulus e3 Action5))*)> */
		func() bool {
			position13, tokenIndex13 := position, tokenIndex
			{
				position14 := position
				if !_rules[rulee3]() {
					goto l13
				}
			l15:
				{
					position16, tokenIndex16 := position, tokenIndex
					{
						position17, tokenIndex17 := position, tokenIndex
						{
							position19 := position
							if buffer[position] != rune('*') {
								goto l18
							}
							position++
							if !_rules[rulesp]() {
								goto l18
							}
							add(rulemultiply, position19)
						}
						if !_rules[rulee3]() {
							goto l18
						}
						{
							add(ruleAction3, position)
						}
						goto l17
					l18:
						position, tokenIndex = position17, tokenIndex17
						{
							position22 := position
							if buffer[position] != rune('/') {
								goto l21
							}
							position++
							if !_rules[rulesp]() {
								goto l21
							}
							add(ruledivide, position22)
						}
						if !_rules[rulee3]() {
							goto l21
						}
						{
							add(ruleAction4, position)
						}
						goto l17
					l21:
						position, tokenIndex = position17, tokenIndex17
						{
							position24 := position
							if buffer[position] != rune('%') {
								goto l16
							}
							position++
							if !_rules[rulesp]() {
								goto l16
							}
							add(rulemodulus, position24)
						}
						if !_rules[rulee3]() {
							goto l16
						}
						{
							add(ruleAction5, position)
						}
					}
				l17:
					goto l15
				l16:
					position, tokenIndex = position16, tokenIndex16
				}
				add(rulee2, position14)
			}
			return true
		l13:
			position, tokenIndex = position13, tokenIndex13
			return false
		},
		/* 3 e3 <- <(eDice (exponentiation eDice Action6)*)> */
		func() bool {
			position26, tokenIndex26 := position, tokenIndex
			{
				position27 := position
				if !_rules[ruleeDice]() {
					goto l26
				}
			l28:
				{
					position29, tokenIndex29 := position, tokenIndex
					{
						position30 := position
						{
							position31, tokenIndex31 := position, tokenIndex
							if buffer[position] != rune('^') {
								goto l32
							}
							position++
							if !_rules[rulesp]() {
								goto l32
							}
							goto l31
						l32:
							position, tokenIndex = position31, tokenIndex31
							if buffer[position] != rune('*') {
								goto l29
							}
							position++
							if buffer[position] != rune('*') {
								goto l29
							}
							position++
							if !_rules[rulesp]() {
								goto l29
							}
						}
					l31:
						add(ruleexponentiation, position30)
					}
					if !_rules[ruleeDice]() {
						goto l29
					}
					{
						add(ruleAction6, position)
					}
					goto l28
				l29:
					position, tokenIndex = position29, tokenIndex29
				}
				add(rulee3, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 4 eDice <- <((e4 ('d' e4 ('k' e4)? Action7)*) / ('d' e4 ('k' e4)? Action8)*)> */
		func() bool {
			{
				position35 := position
				{
					position36, tokenIndex36 := position, tokenIndex
					if !_rules[rulee4]() {
						goto l37
					}
				l38:
					{
						position39, tokenIndex39 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l39
						}
						position++
						if !_rules[rulee4]() {
							goto l39
						}
						{
							position40, tokenIndex40 := position, tokenIndex
							if buffer[position] != rune('k') {
								goto l40
							}
							position++
							if !_rules[rulee4]() {
								goto l40
							}
							goto l41
						l40:
							position, tokenIndex = position40, tokenIndex40
						}
					l41:
						{
							add(ruleAction7, position)
						}
						goto l38
					l39:
						position, tokenIndex = position39, tokenIndex39
					}
					goto l36
				l37:
					position, tokenIndex = position36, tokenIndex36
				l43:
					{
						position44, tokenIndex44 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l44
						}
						position++
						if !_rules[rulee4]() {
							goto l44
						}
						{
							position45, tokenIndex45 := position, tokenIndex
							if buffer[position] != rune('k') {
								goto l45
							}
							position++
							if !_rules[rulee4]() {
								goto l45
							}
							goto l46
						l45:
							position, tokenIndex = position45, tokenIndex45
						}
					l46:
						{
							add(ruleAction8, position)
						}
						goto l43
					l44:
						position, tokenIndex = position44, tokenIndex44
					}
				}
			l36:
				add(ruleeDice, position35)
			}
			return true
		},
		/* 5 e4 <- <((minus value Action9) / value)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				{
					position50, tokenIndex50 := position, tokenIndex
					if !_rules[ruleminus]() {
						goto l51
					}
					if !_rules[rulevalue]() {
						goto l51
					}
					{
						add(ruleAction9, position)
					}
					goto l50
				l51:
					position, tokenIndex = position50, tokenIndex50
					if !_rules[rulevalue]() {
						goto l48
					}
				}
			l50:
				add(rulee4, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 6 value <- <((number Action10) / (identifier Action11) / sub)> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				{
					position55, tokenIndex55 := position, tokenIndex
					{
						position57 := position
						{
							position58 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l56
							}
							position++
						l59:
							{
								position60, tokenIndex60 := position, tokenIndex
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l60
								}
								position++
								goto l59
							l60:
								position, tokenIndex = position60, tokenIndex60
							}
							add(rulePegText, position58)
						}
						if !_rules[rulesp]() {
							goto l56
						}
						add(rulenumber, position57)
					}
					{
						add(ruleAction10, position)
					}
					goto l55
				l56:
					position, tokenIndex = position55, tokenIndex55
					{
						position63 := position
						{
							position64, tokenIndex64 := position, tokenIndex
							if buffer[position] != rune('d') {
								goto l64
							}
							position++
							goto l62
						l64:
							position, tokenIndex = position64, tokenIndex64
						}
						{
							position65 := position
							{
								position68, tokenIndex68 := position, tokenIndex
								{
									position69, tokenIndex69 := position, tokenIndex
									if buffer[position] != rune('!') {
										goto l70
									}
									position++
									goto l69
								l70:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('"') {
										goto l71
									}
									position++
									goto l69
								l71:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('#') {
										goto l72
									}
									position++
									goto l69
								l72:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('$') {
										goto l73
									}
									position++
									goto l69
								l73:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('%') {
										goto l74
									}
									position++
									goto l69
								l74:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('&') {
										goto l75
									}
									position++
									goto l69
								l75:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('\'') {
										goto l76
									}
									position++
									goto l69
								l76:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('(') {
										goto l77
									}
									position++
									goto l69
								l77:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune(')') {
										goto l78
									}
									position++
									goto l69
								l78:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('*') {
										goto l79
									}
									position++
									goto l69
								l79:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('+') {
										goto l80
									}
									position++
									goto l69
								l80:
									position, tokenIndex = position69, tokenIndex69
									if c := buffer[position]; c < rune(',') || c > rune('.') {
										goto l81
									}
									position++
									goto l69
								l81:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('/') {
										goto l82
									}
									position++
									goto l69
								l82:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune(':') {
										goto l83
									}
									position++
									goto l69
								l83:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune(';') {
										goto l84
									}
									position++
									goto l69
								l84:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('<') {
										goto l85
									}
									position++
									goto l69
								l85:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('=') {
										goto l86
									}
									position++
									goto l69
								l86:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('>') {
										goto l87
									}
									position++
									goto l69
								l87:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('?') {
										goto l88
									}
									position++
									goto l69
								l88:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('@') {
										goto l89
									}
									position++
									goto l69
								l89:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('[') {
										goto l90
									}
									position++
									goto l69
								l90:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('\\') {
										goto l91
									}
									position++
									goto l69
								l91:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune(']') {
										goto l92
									}
									position++
									goto l69
								l92:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('^') {
										goto l93
									}
									position++
									goto l69
								l93:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('_') {
										goto l94
									}
									position++
									goto l69
								l94:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('`') {
										goto l95
									}
									position++
									goto l69
								l95:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('{') {
										goto l96
									}
									position++
									goto l69
								l96:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('|') {
										goto l97
									}
									position++
									goto l69
								l97:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('}') {
										goto l98
									}
									position++
									goto l69
								l98:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('~') {
										goto l99
									}
									position++
									goto l69
								l99:
									position, tokenIndex = position69, tokenIndex69
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l100
									}
									position++
									goto l69
								l100:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune(' ') {
										goto l101
									}
									position++
									goto l69
								l101:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('\t') {
										goto l102
									}
									position++
									goto l69
								l102:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('\n') {
										goto l103
									}
									position++
									goto l69
								l103:
									position, tokenIndex = position69, tokenIndex69
									if buffer[position] != rune('\r') {
										goto l68
									}
									position++
								}
							l69:
								goto l62
							l68:
								position, tokenIndex = position68, tokenIndex68
							}
							if !matchDot() {
								goto l62
							}
						l66:
							{
								position67, tokenIndex67 := position, tokenIndex
								{
									position104, tokenIndex104 := position, tokenIndex
									{
										position105, tokenIndex105 := position, tokenIndex
										if buffer[position] != rune('!') {
											goto l106
										}
										position++
										goto l105
									l106:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('"') {
											goto l107
										}
										position++
										goto l105
									l107:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('#') {
											goto l108
										}
										position++
										goto l105
									l108:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('$') {
											goto l109
										}
										position++
										goto l105
									l109:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('%') {
											goto l110
										}
										position++
										goto l105
									l110:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('&') {
											goto l111
										}
										position++
										goto l105
									l111:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('\'') {
											goto l112
										}
										position++
										goto l105
									l112:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('(') {
											goto l113
										}
										position++
										goto l105
									l113:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune(')') {
											goto l114
										}
										position++
										goto l105
									l114:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('*') {
											goto l115
										}
										position++
										goto l105
									l115:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('+') {
											goto l116
										}
										position++
										goto l105
									l116:
										position, tokenIndex = position105, tokenIndex105
										if c := buffer[position]; c < rune(',') || c > rune('.') {
											goto l117
										}
										position++
										goto l105
									l117:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('/') {
											goto l118
										}
										position++
										goto l105
									l118:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune(':') {
											goto l119
										}
										position++
										goto l105
									l119:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune(';') {
											goto l120
										}
										position++
										goto l105
									l120:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('<') {
											goto l121
										}
										position++
										goto l105
									l121:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('=') {
											goto l122
										}
										position++
										goto l105
									l122:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('>') {
											goto l123
										}
										position++
										goto l105
									l123:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('?') {
											goto l124
										}
										position++
										goto l105
									l124:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('@') {
											goto l125
										}
										position++
										goto l105
									l125:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('[') {
											goto l126
										}
										position++
										goto l105
									l126:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('\\') {
											goto l127
										}
										position++
										goto l105
									l127:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune(']') {
											goto l128
										}
										position++
										goto l105
									l128:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('^') {
											goto l129
										}
										position++
										goto l105
									l129:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('_') {
											goto l130
										}
										position++
										goto l105
									l130:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('`') {
											goto l131
										}
										position++
										goto l105
									l131:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('{') {
											goto l132
										}
										position++
										goto l105
									l132:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('|') {
											goto l133
										}
										position++
										goto l105
									l133:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('}') {
											goto l134
										}
										position++
										goto l105
									l134:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('~') {
											goto l135
										}
										position++
										goto l105
									l135:
										position, tokenIndex = position105, tokenIndex105
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l136
										}
										position++
										goto l105
									l136:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune(' ') {
											goto l137
										}
										position++
										goto l105
									l137:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('\t') {
											goto l138
										}
										position++
										goto l105
									l138:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('\n') {
											goto l139
										}
										position++
										goto l105
									l139:
										position, tokenIndex = position105, tokenIndex105
										if buffer[position] != rune('\r') {
											goto l104
										}
										position++
									}
								l105:
									goto l67
								l104:
									position, tokenIndex = position104, tokenIndex104
								}
								if !matchDot() {
									goto l67
								}
								goto l66
							l67:
								position, tokenIndex = position67, tokenIndex67
							}
							add(rulePegText, position65)
						}
						if !_rules[rulesp]() {
							goto l62
						}
						add(ruleidentifier, position63)
					}
					{
						add(ruleAction11, position)
					}
					goto l55
				l62:
					position, tokenIndex = position55, tokenIndex55
					{
						position141 := position
						{
							position142 := position
							if buffer[position] != rune('(') {
								goto l53
							}
							position++
							if !_rules[rulesp]() {
								goto l53
							}
							add(ruleopen, position142)
						}
						if !_rules[rulee1]() {
							goto l53
						}
						{
							position143 := position
							if buffer[position] != rune(')') {
								goto l53
							}
							position++
							if !_rules[rulesp]() {
								goto l53
							}
							add(ruleclose, position143)
						}
						add(rulesub, position141)
					}
				}
			l55:
				add(rulevalue, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 7 number <- <(<[0-9]+> sp)> */
		nil,
		/* 8 identifier <- <(!'d' <(!('!' / '"' / '#' / '$' / '%' / '&' / '\'' / '(' / ')' / '*' / '+' / [,-.] / '/' / ':' / ';' / '<' / '=' / '>' / '?' / '@' / '[' / '\\' / ']' / '^' / '_' / '`' / '{' / '|' / '}' / '~' / [0-9] / ' ' / '\t' / '\n' / '\r') .)+> sp)> */
		nil,
		/* 9 sub <- <(open e1 close)> */
		nil,
		/* 10 add <- <('+' sp)> */
		nil,
		/* 11 minus <- <('-' sp)> */
		func() bool {
			position148, tokenIndex148 := position, tokenIndex
			{
				position149 := position
				if buffer[position] != rune('-') {
					goto l148
				}
				position++
				if !_rules[rulesp]() {
					goto l148
				}
				add(ruleminus, position149)
			}
			return true
		l148:
			position, tokenIndex = position148, tokenIndex148
			return false
		},
		/* 12 multiply <- <('*' sp)> */
		nil,
		/* 13 divide <- <('/' sp)> */
		nil,
		/* 14 modulus <- <('%' sp)> */
		nil,
		/* 15 exponentiation <- <(('^' sp) / ('*' '*' sp))> */
		nil,
		/* 16 open <- <('(' sp)> */
		nil,
		/* 17 close <- <(')' sp)> */
		nil,
		/* 18 sp <- <(' ' / '\t')*> */
		func() bool {
			{
				position157 := position
			l158:
				{
					position159, tokenIndex159 := position, tokenIndex
					{
						position160, tokenIndex160 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l161
						}
						position++
						goto l160
					l161:
						position, tokenIndex = position160, tokenIndex160
						if buffer[position] != rune('\t') {
							goto l159
						}
						position++
					}
				l160:
					goto l158
				l159:
					position, tokenIndex = position159, tokenIndex159
				}
				add(rulesp, position157)
			}
			return true
		},
		/* 20 Action0 <- <{ p.AddOperator(TypeHalt) }> */
		nil,
		/* 21 Action1 <- <{ p.AddOperator(TypeAdd) }> */
		nil,
		/* 22 Action2 <- <{ p.AddOperator(TypeSubtract) }> */
		nil,
		/* 23 Action3 <- <{ p.AddOperator(TypeMultiply) }> */
		nil,
		/* 24 Action4 <- <{ p.AddOperator(TypeDivide) }> */
		nil,
		/* 25 Action5 <- <{ p.AddOperator(TypeModulus) }> */
		nil,
		/* 26 Action6 <- <{ p.AddOperator(TypeExponentiation) }> */
		nil,
		/* 27 Action7 <- <{ p.AddOperator(TypeDice) }> */
		nil,
		/* 28 Action8 <- <{ p.AddValue("1"); p.AddOperator(TypeSwap); p.AddOperator(TypeDice) }> */
		nil,
		/* 29 Action9 <- <{ p.AddOperator(TypeNegation) }> */
		nil,
		/* 30 Action10 <- <{ p.AddValue(string(text)) }> */
		nil,
		/* 31 Action11 <- <{ p.AddLoadVarname(string(text)) }> */
		nil,
		nil,
	}
	p.rules = _rules
	return nil
}