package main 

import (
	"fmt"
	"os"
	"lzw"
)

func usage() {
	progName := os.Args[0]
	fmt.Printf("%s:\n", progName)
	fmt.Printf("\t[-e | -d]: Specifies to encode or decode\n")
	fmt.Printf("\tinput_file_name:  the input file (file to becompressed for encode, or compressed file for decode\n")
	fmt.Printf("\toutput_file_name: the output file (name of the created compressed file for encode, or name of newly decompressed fle for decode\n")
}

func main() {

	numArgs := len(os.Args)
	if numArgs < 4 {
		usage()
		os.Exit(1)
	}

	method := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
                   
    if(method == "-e") {
   		lzw.Encode(inputFile, outputFile)
    } else {
    	lzw.Decode(inputFile, outputFile)
    }
}
