package main

import (
	"fmt"

	filereader "github.com/sachinaralapura/shoebill/fileReader"
	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/parser"
)

func main() {
	// repl.Start(os.Stdin, os.Stdout)
	// create channels
	fileToLexChan := make(chan []byte)

	fileReader := filereader.New(fileToLexChan)
	fileReader.SetFileNameFromArgs()
	go fileReader.ReadChunk()

	lexer := lexer.New(fileToLexChan)
	lexer.LoadBuffer()

	parser := parser.New(lexer)
	program := parser.ParseProgram()
	fmt.Println(lexer)

	fmt.Println(program)
}
