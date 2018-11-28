package toolbox

import (
	"fmt"
	"strings"
)

//Matcher represents a matcher, that matches input from offset position, it returns number of characters matched.
type Matcher interface {
	//Match matches input starting from offset, it return number of characters matched
	Match(input string, offset int) (matched int)
}

//Token a matchable input
type Token struct {
	Token   int
	Matched string
}

//Tokenizer represents a token scanner.
type Tokenizer struct {
	matchers       map[int]Matcher
	Input          string
	Index          int
	InvalidToken   int
	EndOfFileToken int
}

//Nexts matches the first of the candidates
func (t *Tokenizer) Nexts(candidates ...int) *Token {
	for _, candidate := range candidates {
		result := t.Next(candidate)
		if result.Token != t.InvalidToken {
			return result

		}
	}
	return &Token{t.InvalidToken, ""}
}

//Next tries to match a candidate, it returns token if imatching is successful.
func (t *Tokenizer) Next(candidate int) *Token {
	offset := t.Index
	if !(offset < len(t.Input)) {
		return &Token{t.EndOfFileToken, ""}
	}

	if candidate == t.EndOfFileToken {
		return &Token{t.InvalidToken, ""}
	}
	if matcher, ok := t.matchers[candidate]; ok {
		matchedSize := matcher.Match(t.Input, offset)
		if matchedSize > 0 {
			t.Index = t.Index + matchedSize
			return &Token{candidate, t.Input[offset : offset+matchedSize]}
		}

	} else {
		panic(fmt.Sprintf("failed to lookup matcher for %v", candidate))
	}
	return &Token{t.InvalidToken, ""}
}

//NewTokenizer creates a new NewTokenizer, it takes input, invalidToken, endOfFileToeken, and matchers.
func NewTokenizer(input string, invalidToken int, endOfFileToken int, matcher map[int]Matcher) *Tokenizer {
	return &Tokenizer{
		matchers:       matcher,
		Input:          input,
		Index:          0,
		InvalidToken:   invalidToken,
		EndOfFileToken: endOfFileToken,
	}
}

//CharactersMatcher represents a matcher, that matches any of Chars.
type CharactersMatcher struct {
	Chars string //characters to be matched
}

//Match matches any characters defined in Chars in the input, returns 1 if character has been matched
func (m CharactersMatcher) Match(input string, offset int) (matched int) {
	var result = 0
outer:
	for i := 0; i < len(input)-offset; i++ {
		aChar := input[offset+i : offset+i+1]
		for j := 0; j < len(m.Chars); j++ {
			if aChar == m.Chars[j:j+1] {
				result++
				continue outer
			}
		}
		break
	}
	return result
}

func isLetter(aChar string) bool {
	return (aChar >= "a" && aChar <= "z") || (aChar >= "A" && aChar <= "Z")
}

func isDigit(aChar string) bool {
	return (aChar >= "0" && aChar <= "9")
}

//EOFMatcher represents end of input matcher
type EOFMatcher struct {
}

//Match returns 1 if end of input has been reached otherwise 0
func (m EOFMatcher) Match(input string, offset int) (matched int) {
	if offset+1 == len(input) {
		return 1
	}
	return 0
}

//IntMatcher represents a matcher that finds any int in the input
type IntMatcher struct{}

//Match matches a literal in the input, it returns number of character matched.
func (m IntMatcher) Match(input string, offset int) (matched int) {
	if !isDigit(input[offset : offset+1]) {
		return 0
	}
	var i = 1
	for ; i < len(input)-offset; i++ {
		aChar := input[offset+i : offset+i+1]
		if !isDigit(aChar) {
			break
		}
	}
	return i
}

//NewIntMatcher returns a new integer matcher
func NewIntMatcher() Matcher {
	return &IntMatcher{}
}

//LiteralMatcher represents a matcher that finds any literals in the input
type LiteralMatcher struct{}

//Match matches a literal in the input, it returns number of character matched.
func (m LiteralMatcher) Match(input string, offset int) (matched int) {
	if !isLetter(input[offset : offset+1]) {
		return 0
	}
	var i = 1
	for ; i < len(input)-offset; i++ {
		aChar := input[offset+i : offset+i+1]
		if !((isLetter(aChar)) || isDigit(aChar) || aChar == "_" || aChar == ".") {
			break
		}
	}
	return i
}

//LiteralMatcher represents a matcher that finds any literals in the input
type IdMatcher struct{}

//Match matches a literal in the input, it returns number of character matched.
func (m IdMatcher) Match(input string, offset int) (matched int) {
	if !isLetter(input[offset:offset+1]) && !isDigit(input[offset:offset+1]) {
		return 0
	}
	var i = 1
	for ; i < len(input)-offset; i++ {
		aChar := input[offset+i : offset+i+1]
		if !((isLetter(aChar)) || isDigit(aChar) || aChar == "_" || aChar == ".") {
			break
		}
	}
	return i
}

//SequenceMatcher represents a matcher that finds any sequence until find provided terminators
type SequenceMatcher struct {
	Terminators   []string
	CaseSensitive bool
}

func (m *SequenceMatcher) hasTerminator(candidate string) bool {
	var candidateLength = len(candidate)
	for _, terminator := range m.Terminators {
		terminatorLength := len(terminator)
		if len(terminator) > candidateLength {
			continue
		}

		if !m.CaseSensitive {
			if strings.ToLower(terminator) == strings.ToLower(string(candidate[:terminatorLength])) {
				return true
			}
		}

		if terminator == string(candidate[:terminatorLength]) {
			return true
		}
	}
	return false
}

//Match matches a literal in the input, it returns number of character matched.
func (m *SequenceMatcher) Match(input string, offset int) (matched int) {
	var i = 0
	for ; i < len(input)-offset; i++ {
		if m.hasTerminator(string(input[offset+i:])) {
			break
		}
	}
	return i
}

//NewSequenceMatcher creates a new matcher that finds any sequence until find provided terminators
func NewSequenceMatcher(terminators ...string) Matcher {
	return &SequenceMatcher{
		Terminators: terminators,
	}
}

//remainingSequenceMatcher represents a matcher that matches all reamining input
type remainingSequenceMatcher struct{}

//Match matches a literal in the input, it returns number of character matched.
func (m *remainingSequenceMatcher) Match(input string, offset int) (matched int) {
	return len(input) - offset
}

//Creates a matcher that matches all remaining input
func NewRemainingSequenceMatcher() Matcher {
	return &remainingSequenceMatcher{}
}

//CustomIdMatcher represents a matcher that finds any literals with additional custom set of characters in the input
type customIdMatcher struct {
	Allowed map[string]bool
}

func (m *customIdMatcher) isValid(aChar string) bool {
	if isLetter(aChar) {
		return true
	}
	if isDigit(aChar) {
		return true
	}
	return m.Allowed[aChar]
}

//Match matches a literal in the input, it returns number of character matched.
func (m *customIdMatcher) Match(input string, offset int) (matched int) {

	if !m.isValid(input[offset : offset+1]) {
		return 0
	}
	var i = 1
	for ; i < len(input)-offset; i++ {
		aChar := input[offset+i : offset+i+1]
		if !m.isValid(aChar) {
			break
		}
	}
	return i
}

//NewCustomIdMatcher creates new custom matcher
func NewCustomIdMatcher(allowedChars ...string) Matcher {
	var result = &customIdMatcher{
		Allowed: make(map[string]bool),
	}

	if len(allowedChars) == 1 && len(allowedChars[0]) > 0 {
		for _, allowed := range allowedChars[0] {
			result.Allowed[string(allowed)] = true
		}
	}
	for _, allowed := range allowedChars {
		result.Allowed[allowed] = true
	}
	return result
}

//LiteralMatcher represents a matcher that finds any literals in the input
type BodyMatcher struct {
	Begin string
	End   string
}

//Match matches a literal in the input, it returns number of character matched.
func (m *BodyMatcher) Match(input string, offset int) (matched int) {
	beginLen := len(m.Begin)
	endLen := len(m.End)
	uniEclosed := m.Begin == m.End

	if offset+beginLen >= len(input) {
		return 0
	}
	if input[offset:offset+beginLen] != m.Begin {
		return 0
	}

	var depth = 1
	var i = 1
	for ; i < len(input)-offset; i++ {

		canCheckEnd := offset+i+endLen <= len(input)
		if !canCheckEnd {
			return 0
		}
		if !uniEclosed {
			canCheckBegin := offset+i+beginLen <= len(input)
			if canCheckBegin {
				if string(input[offset+i:offset+i+beginLen]) == m.Begin {
					depth++
				}
			}
		}
		if string(input[offset+i:offset+i+endLen]) == m.End {
			depth--
		}
		if depth == 0 {
			i += endLen
			break
		}
	}
	return i
}

//KeywordMatcher represents a keyword matcher
type KeywordMatcher struct {
	Keyword       string
	CaseSensitive bool
}

//Match matches keyword in the input,  it returns number of character matched.
func (m KeywordMatcher) Match(input string, offset int) (matched int) {
	if !(offset+len(m.Keyword)-1 < len(input)) {
		return 0
	}

	if m.CaseSensitive {
		if input[offset:offset+len(m.Keyword)] == m.Keyword {
			return len(m.Keyword)
		}
	} else {
		if strings.ToLower(input[offset:offset+len(m.Keyword)]) == strings.ToLower(m.Keyword) {
			return len(m.Keyword)
		}
	}
	return 0
}

//KeywordsMatcher represents a matcher that finds any of specified keywords in the input
type KeywordsMatcher struct {
	Keywords      []string
	CaseSensitive bool
}

//Match matches any specified keyword,  it returns number of character matched.
func (m KeywordsMatcher) Match(input string, offset int) (matched int) {
	for _, keyword := range m.Keywords {
		if len(input)-offset < len(keyword) {
			continue
		}
		if m.CaseSensitive {
			if input[offset:offset+len(keyword)] == keyword {
				return len(keyword)
			}
		} else {
			if strings.ToLower(input[offset:offset+len(keyword)]) == strings.ToLower(keyword) {
				return len(keyword)
			}
		}
	}
	return 0
}

//NewKeywordsMatcher returns a matcher for supplied keywords
func NewKeywordsMatcher(caseSensitive bool, keywords ...string) Matcher {
	return &KeywordsMatcher{CaseSensitive: caseSensitive, Keywords: keywords}
}

//IllegalTokenError represents illegal token error
type IllegalTokenError struct {
	Illegal  *Token
	Message  string
	Expected []int
	Position int
}

func (e *IllegalTokenError) Error() string {
	return fmt.Sprintf("%v; illegal token at %v [%v], expected %v, but had: %v", e.Message, e.Position, e.Illegal.Matched, e.Expected, e.Illegal.Token)
}

//NewIllegalTokenError create a new illegal token error
func NewIllegalTokenError(message string, expected []int, position int, found *Token) error {
	return &IllegalTokenError{
		Message:  message,
		Illegal:  found,
		Expected: expected,
		Position: position,
	}
}

//ExpectTokenOptionallyFollowedBy returns second matched token or error if first and second group was not matched
func ExpectTokenOptionallyFollowedBy(tokenizer *Tokenizer, first int, errorMessage string, second ...int) (*Token, error) {
	_, _ = ExpectToken(tokenizer, "", first)
	return ExpectToken(tokenizer, errorMessage, second...)
}

//ExpectToken returns the matched token or error
func ExpectToken(tokenizer *Tokenizer, errorMessage string, candidates ...int) (*Token, error) {
	token := tokenizer.Nexts(candidates...)
	hasMatch := HasSliceAnyElements(candidates, token.Token)
	if !hasMatch {
		return nil, NewIllegalTokenError(errorMessage, candidates, tokenizer.Index, token)
	}
	return token, nil
}
