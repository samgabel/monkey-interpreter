package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/samgabel/monkey-interpreter/lexer"
	"github.com/samgabel/monkey-interpreter/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// first print our the prompt
		fmt.Fprintf(out, PROMPT)
		if !scanner.Scan() {
			return
		}

		// grab the text that was input into the REPL and feed it into our lexer
		line := scanner.Text()
		l := lexer.NewLexer(line)

		// (Lexer).NextToken() serves to increment the token in the input string,
		// this happends until the TokenType of a token is EOF.
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok) // %+v shows the values and fields of the token
		}
	}
}
