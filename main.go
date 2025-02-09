package main

import (
	"fmt"

	filereader "github.com/sachinaralapura/shoebill/fileReader"
	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/parser"
)

func main() {
	// create channels
	fileToLexChan := make(chan []byte)

	fileReader := filereader.New(fileToLexChan)
	fileReader.SetFileNameFromArgs()
	go fileReader.ReadChunk()

	lexer := lexer.New(fileToLexChan)
	lexer.LoadBuffer()

	// fmt.Println(lexer)

	parser := parser.New(lexer)
	program := parser.ParseProgram()

	fmt.Println(program)
}
