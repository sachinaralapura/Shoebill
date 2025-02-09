package filereader

import (
	"fmt"
	"log"
	"os"

	"github.com/sachinaralapura/shoebill/constants"
)

type FileReader struct {
	FileName string
	sendChan chan<- []byte
}

// sets the FileName from the command line argument
func (filereader *FileReader) GetFileName() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}
	return "", fmt.Errorf("Command Line argument not found")
}

func (fr *FileReader) SetFileNameFromArgs() {
	filename, err := fr.GetFileName()
	if err != nil {
		log.Fatal(err)
	}
	fr.FileName = filename
}

func (filereader *FileReader) ReadChunk() {
	file, err := os.Open(filereader.FileName)
	if err != nil {
		log.Fatalf("Error opening file %s : file not found", filereader.FileName)
	}
	log.Println("reading file : ", filereader.FileName)
	defer file.Close()

	buffer := make([]byte, constants.ChunkSize) //  buffer
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			chunk := make([]byte, n) // Create a slice of exact size
			copy(chunk, buffer[:n])  // Copy only read bytes
			filereader.sendChan <- chunk
		}

		if err != nil {
			if err.Error() == "EOF" {
				log.Println("End of file:", err)
			}

			log.Println("closing file to lex channel")
			close(filereader.sendChan)
			break
		}
	}
}

func New(in chan<- []byte) *FileReader {
	fr := &FileReader{sendChan: in}
	// fr.SetFileNameFromArgs()
	return fr
}
