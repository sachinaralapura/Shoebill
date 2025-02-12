package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/parser"
)

const PROMPT = ">>"

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

		fmt.Println(program)
	}
}
