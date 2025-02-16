package repl

import (
	"bufio"
	"fmt"
	"strings"

	"io"

	"github.com/sachinaralapura/shoebill/evaluator"
	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/object"
	"github.com/sachinaralapura/shoebill/parser"
)

const PROMPT = ">>"
const p = ``

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvirnoment()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if strings.ToLower(line) == "quit" || strings.ToLower(line) == "exit" {
			break
		}
		p := createParser(line)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
			continue
		}
		// -----------------------------
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func createParser(line string) *parser.Parser {
	outchannel := make(chan []byte)
	l := lexer.New(outchannel)
	go func() {
		outchannel <- []byte(line)
		close(outchannel)
	}()
	l.LoadBuffer()
	p := parser.New(l)
	return p
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
