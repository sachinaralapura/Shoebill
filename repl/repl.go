package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/parser"
)

const PROMPT = ">>"
const p = ``

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		outchannel := make(chan []byte)
		l := lexer.New(outchannel)
		go func() {
			outchannel <- []byte(line)
			close(outchannel)
		}()
		l.LoadBuffer()
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}

}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
