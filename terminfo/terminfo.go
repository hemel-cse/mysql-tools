// terminfo package provides simple parser for compiled terminfo fields
// and additionally it provides trivial API for applying of capability
// of a terminal.
package terminfo

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gitlab.com/0xAX/emc/ds"
)

// Arrays of static and dynamic termifno variables
type TerminfoVariables struct {
	static  []int // Static variables A-Z
	dynamic []int // Dynamic variables a-z
}

// States of parameterized capabilities parser
const (
	initial          = iota
	percent          = iota
	condition        = iota
	pushParam        = iota
	setParam         = iota
	getParam         = iota
	charStart        = iota
	charEnd          = iota
	integerParameter = iota
	parseFlags       = iota
	skipThenExpr     = iota
	ifExprPercent    = iota
	elseExpr         = iota
	elsePercent      = iota
)

// terminfo header length
const terminfoHeaderLength = 12

// terminfo magic header
const terminfoMagic1 = 0x1a
const terminfoMagic2 = 0x01

// zero in ASCII
const zeroASCII = 48

// default path to terminfo database
const terminfoDefaultPath = "/usr/share/terminfo"

// terminfo related enviroment variables
const terminfoEnv = "TERMINFO"
const termEnv = "TERM"
const terminfoDirsEnv = "TERMINFO_DIRS"
const terminfoLibPath = "/lib/terminfo"
const terminfoUsrPath = "/usr/share/lib/terminfo"

// Terminfo base structure of terminfo library. It contains
// four fields which represent sets of terminal names
// and capabilities of a terminal
type Terminfo struct {
	Name        string            // name of a terminal
	Description string            // desription of a terminal
	Bools       map[string]bool   // map of enabled booleans of a terminal
	Numbers     map[string]uint16 // map of number values of a terminal
	Strings     map[string]string // map of string values of a terminal
}

// ParseTerminfo reads terminfo database and returns pointer
// to the *TermInfo for the currently used terminal from $TERM
// environment variable and other terminfo related enviroment
// variables.
func ParseTerminfoFromEnv() (*Terminfo, error) {
	term := os.Getenv(termEnv)
	if term == "" {
		return nil, &errorTerminfo{"$TERM is empty"}
	}
	// read $TERMINFO and use it as default path for a terminfo entry
	terminfoPath := os.Getenv(terminfoEnv)
	if terminfoPath != "" {
		return ParseTerminfoFromPath(terminfoPath + "/" + string(term[0]) + "/" + term)
	}
	// the $TERMINFO is unset, let's try to read $HOME/.terminfo
	home := os.Getenv("HOME")
	if home != "" {
		_, err := os.Stat(home + "/.terminfo")
		if err == nil {
			return ParseTerminfoFromPath(home + "/.terminfo" + "/" + string(term[0]) + "/" + term)
		}
	}
	// let's try $TERMINFO_DIRS before we go to the default path
	terminfoPath = os.Getenv(terminfoDirsEnv)
	if terminfoPath != "" {
		return ParseTerminfoFromPath(terminfoPath + "/" + string(term[0]) + "/" + term)
	}
	// We are looking for in lib/terminfo
	_, err := os.Stat(terminfoLibPath)
	if err == nil {
		return ParseTerminfoFromPath(terminfoLibPath + "/" + string(term[0]) + "/" + term)
	}
	// We are looking for in /usr/share/lib/terminfo/
	_, err = os.Stat(terminfoUsrPath)
	if err == nil {
		return ParseTerminfoFromPath(terminfoUsrPath + "/" + string(term[0]) + "/" + term)
	}
	// read a terminfo entry from the default path
	terminfoContent, err := ioutil.ReadFile(terminfoDefaultPath + "/" + string(term[0]) + "/" + term)
	if err != nil {
		return nil, err
	}
	return parseTerminfo(terminfoContent)
}

// ParseTerminfo reads terminfo database and returns pointer
// to the *TermInfo for the currently used terminal from the
// given path.
func ParseTerminfoFromPath(path string) (*Terminfo, error) {
	if path == "" {
		return nil, &errorTerminfo{"terminal path can't be \"\""}
	}
	// read terminfo for the given terminal
	terminfoContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseTerminfo(terminfoContent)
}

func parseTerminfo(terminfoContent []byte) (*Terminfo, error) {
	var namesSectionSize uint16 = 0
	var booleanSectionSize uint16 = 0
	var shortIntsSectionSize uint16 = 0
	var offsetsNumber uint16 = 0

	terminfo := &Terminfo{}
	terminfo.Bools = make(map[string]bool)
	terminfo.Numbers = make(map[string]uint16)
	terminfo.Strings = make(map[string]string)

	// check magic header of a terminfo file
	if terminfoContent[0] != terminfoMagic1 && terminfoContent[1] != terminfoMagic2 {
		return nil, &errorTerminfo{"wrong terminfo file"}
	}

	// parse header of a terminfo
	namesSectionSize = binary.LittleEndian.Uint16(terminfoContent[2:4])
	booleanSectionSize = binary.LittleEndian.Uint16(terminfoContent[4:6])
	shortIntsSectionSize = binary.LittleEndian.Uint16(terminfoContent[6:8])
	offsetsNumber = binary.LittleEndian.Uint16(terminfoContent[8:10])

	// parse terminal name and description
	terminalName := strings.Split(string(terminfoContent[terminfoHeaderLength:terminfoHeaderLength+namesSectionSize]), "|")
	terminfo.Name = terminalName[0]
	terminfo.Description = terminalName[1]

	// parse booleans
	offset := int(terminfoHeaderLength + namesSectionSize)
	for i := 0; i < int(booleanSectionSize)-1; i++ {
		val := terminfoContent[offset+i : int(offset+i+1)][0]
		if val == 0 {
			terminfo.Bools[GetTerminfoBooleanCodes()[i]] = false
		} else {
			terminfo.Bools[GetTerminfoBooleanCodes()[i]] = true
		}
	}

	// pad offset to short integer boundary
	offset += int(booleanSectionSize) + (offset+int(booleanSectionSize))%2

	// parse numbers
	numbericIndex := 0
	for i := 0; i < int(shortIntsSectionSize*2); i += 2 {
		val := binary.LittleEndian.Uint16(terminfoContent[offset+i : int(offset+i+2)])
		if val != 0 && val != 0x3737 && val != 0xffff {
			terminfo.Numbers[GetTerminfoNumericCodes()[numbericIndex]] = val
		}
		numbericIndex++
	}
	offset += int(2 * shortIntsSectionSize)
	// parse strings
	stringTableStart := offset + int(offsetsNumber*2)
	termCapabilityStr := ""
	stringIndex := 0
	for i := 0; i < int(offsetsNumber*2); i += 2 {
		stringTableOffset := binary.LittleEndian.Uint16(terminfoContent[offset+i : int(offset+i+2)])
		if stringTableOffset != 0xffff {
			stringTable := terminfoContent[stringTableStart+int(stringTableOffset):]
			for _, t := range stringTable {
				if t != 0 {
					termCapabilityStr += string(t)
					continue
				}
				break
			}
			terminfo.Strings[GetTerminfoStringCodes()[stringIndex]] = termCapabilityStr
			termCapabilityStr = ""
		}
		stringIndex++
	}
	return terminfo, nil
}

// ApplyCapability function applies the given terminfo capability.
// Additionally it privdes parser for parameterized capabilities.
// TODO
//   * make normal errors and names
//
func (terminfo *Terminfo) ApplyCapability(capability string, args ...interface{}) (string, error) {
	result := ""
	format := ""
	strNumber := ""
	state := initial
	stack := stack.New()
	variables := &TerminfoVariables{static: make([]int, 26), dynamic: make([]int, 26)}
	ifExprNestedLevel := 0

	if seq, ok := terminfo.Strings[capability]; ok {
		for _, cur := range seq {
			switch state {
			case initial:
				if cur == '%' {
					state = percent
				} else {
					result += string(cur)
				}
				continue
			case percent:
				switch cur {
				case '%':
					result += "%"
					state = initial
					break
				case 'x', 'X', 'o', 's', 'd':
					result += fmt.Sprintf("%"+string(cur), stack.Pop())
					state = initial
					break
				case 'c':
					value := stack.Pop()
					if value == 0 {
						result += string(0x80)
					} else {
						result += string(value.(uint8))
					}
					state = initial
					break
				case 'p':
					state = pushParam
					break
				case 'P':
					state = setParam
					break
				case 'l':
					stack.Push(len(fmt.Sprintf("%s", stack.Pop())))
					state = initial
					break
				case 'g':
					state = getParam
					break
				case '<', 'O', 'A', '=', '>':
					set := 0
					v1 := stack.Pop().(int)
					v2 := stack.Pop().(int)
					switch cur {
					case 'O':
						if v1 > 0 || v2 > 0 {
							stack.Push(1)
							set = 1
							state = initial
						}
					case 'A':
						if v1 > 0 && v2 > 0 {
							stack.Push(1)
							state = initial
							set = 1
						}
					case '=':
						if v1 == v2 {
							stack.Push(1)
							state = initial
							set = 1
						}
					case '<':
						if v1 < v2 {
							stack.Push(1)
							state = initial
							set = 1
						}
					case '>':
						if v1 > v2 {
							stack.Push(1)
							state = initial
							set = 1
						}
					}
					if set == 1 {
						break
					}
					stack.Push(0)
					state = initial
					break
				case '\'':
					state = charStart
				case '+':
				case '-':
				case '/':
				case '*':
				case '^':
				case '&':
				case '|':
				case 'm':
					v1 := stack.Pop().(int)
					v2 := stack.Pop().(int)
					switch cur {
					case '+':
						stack.Push(v1 + v2)
						break
					case '-':
						stack.Push(v1 - v2)
						break
					case '/':
						stack.Push(v1 / v2)
						break
					case '*':
						stack.Push(v1 * v2)
						break
					case '^':
						stack.Push(v1 ^ v2)
						break
					case '&':
						stack.Push(v1 & v2)
						break
					case '|':
						stack.Push(v1 | v2)
						break
					case 'm':
						stack.Push(v1 % v2)
						break
					}
					break
				case '!':
					value := stack.Pop()
					if value.(int) > 0 {
						stack.Push(0)
					} else {
						stack.Push(1)
					}
					break
				case '~':
					value := stack.Pop()
					stack.Push(^value.(int))
					break
				case '{':
					state = integerParameter
					break
				case 'i':
					args[0] = args[0].(int) + 1
					args[1] = args[1].(int) + 1
					state = initial
					break
				case ':':
					state = parseFlags
					break
				case ';':
				case '?':
					state = initial
					break
				case 't':
					value := stack.Pop()
					if value == 0 {
						state = skipThenExpr
						break
					}
					state = initial
					break
				case 'e':
					state = elseExpr
					break
				default:
					state = initial
					break
				}
				continue
			case parseFlags:
				switch cur {
				case '+':
				case '-':
				case '#':
				case ' ':
				case '.':
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					format += string(cur)
				case '%':
					result += fmt.Sprintf("%"+format, stack.Pop())
					state = initial
					format = ""
				}
				break
			case charStart:
				stack.Push(string(cur))
				state = charEnd
				break
			case charEnd:
				if cur != '\'' {
					return "", &errorTerminfo{"a char must be closed with '"}
				}
				state = initial
				break
			case integerParameter:
				if cur == '}' {
					i, err := strconv.Atoi(strNumber)
					if err != nil {
						return "", &errorTerminfo{"{nn} must be integer number"}
					}
					stack.Push(i)
					strNumber = ""
					state = initial
				} else {
					i := strconv.Itoa(int(64 - cur - 9))
					strNumber += i
				}
				break
			case getParam:
				if cur >= 'a' && cur <= 'z' {
					cur = cur - 97
					stack.Push(variables.dynamic[cur])
					state = initial
					continue
				}

				if cur >= 'A' && cur <= 'Z' {
					cur = cur - 66
					stack.Push(variables.static[cur])
					state = initial
					continue
				}
				return "", &errorTerminfo{"'g' must be followed with [a-zA-Z]"}
			case setParam:
				value := stack.Pop()
				if cur >= 'a' && cur <= 'z' {
					cur = cur - 97
					variables.dynamic[cur] = value.(int)
					state = initial
					continue
				}

				if cur >= 'A' && cur <= 'Z' {
					cur = cur - 66
					variables.static[cur] = value.(int)
					state = initial
					continue
				}
				return "", &errorTerminfo{"'P' must be followed with [a-zA-Z]"}
			case pushParam:
				stack.Push(args[cur-49])
				state = initial
				break
			case skipThenExpr:
				if cur == '%' {
					state = ifExprPercent
				}
				break
			case ifExprPercent:
				if cur == ';' {
					if ifExprNestedLevel == 0 {
						state = initial
					} else {
						ifExprNestedLevel -= 1
						state = skipThenExpr
					}
				} else if cur == 'e' && ifExprNestedLevel == 0 {
					state = initial
				} else if cur == '?' {
					ifExprNestedLevel += 1
					state = skipThenExpr
				} else {
					state = skipThenExpr
				}
				break
			case elseExpr:
				if cur == '%' {
					state = elsePercent
				}
				break
			case elsePercent:
				if cur == ';' {
					if ifExprNestedLevel == 0 {
						state = initial
					} else {
						ifExprNestedLevel -= 1
						state = elseExpr
					}
				} else if cur == '?' {
					ifExprNestedLevel += 1
					state = elseExpr
				} else {
					state = elseExpr
				}
				break
			}
		}
		return result, nil
	}

	return "", &errorTerminfo{"Your terminal does not support such capability"}
}
