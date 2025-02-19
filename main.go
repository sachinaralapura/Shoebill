package main

import (
	"os"

	"github.com/sachinaralapura/shoebill/evaluator"
	filereader "github.com/sachinaralapura/shoebill/fileReader"
	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/object"
	"github.com/sachinaralapura/shoebill/parser"
	"github.com/sachinaralapura/shoebill/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
	// create channels
	fileToLexChan := make(chan []byte)

	fileReader := filereader.New(fileToLexChan)
	fileReader.SetFileNameFromArgs()
	go fileReader.ReadChunk()

	lexer := lexer.New(fileToLexChan)
	lexer.LoadBuffer()

	parser := parser.New(lexer)
	program := parser.ParseProgram()

	env := object.NewEnvirnoment()
	evaluator.Eval(program, env)
}
