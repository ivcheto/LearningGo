package lzw 

import (
	"io"
	"fmt"
	"os"
	"strconv"
)

func getEncodeDictionary() (dict map[string]int) {
	dict = make(map[string]int)

	for i := 0; i < 128; i++ {
		letter := fmt.Sprintf("%c", i)
		dict [letter] = (i)
	}

	return 
}

func getDecodeDictionary() (dict map[int]string) {
	dict = make(map[int]string)

	for i := 0; i < 128; i++ {
		letter := fmt.Sprintf("%c", i)
		dict[i] = letter
	}

	return 
}

/*
 openInput accepts a string representing a file to be read, and returns
 a file pointer to that file.  This function panics on error (including if 
 it gets a bad file descriptor)

 TODO - this function doesn't quite work yet - when passed a filename, it 
 throws a bad file descriptor, but when the same code is pasted into the calling
 method, it works like a charm o.O Still need to look into this
 */
func openInput(inputFileName string) (in *os.File) {
	//
	// open input file
	in, err := os.OpenFile(inputFileName, os.O_RDONLY, 644)
	if err != nil {
		panic(fmt.Sprintf("Tried to open input file but got error: %v", err))
	}

	// defering the annonymous function, both defining it here, and invoking it - pretty cool :)
	defer func(){
		if 	err := in.Close(); err != nil {
			panic(fmt.Sprintf("Tried closing input file %s but got error: %v", inputFileName, err))
		}
	} ()

	return 
}

/*
 openOutput accepts a string representing a file to be written, and returns
 a file pointer to that file.  This function panics on error (including if 
 it gets a bad file descriptor)

 TODO - this function doesn't quite work yet - when passed a filename, it 
 throws a bad file descriptor, but when the same code is pasted into the calling
 method, it works like a charm o.O Still need to look into this
 */
func openOutput(outputFileName string) (out *os.File) {
	//
	// open output file
	out, err := os.Create(outputFileName)
	if err != nil {
		panic(fmt.Sprintf("Tried to open output file %s, but got error: %v", outputFileName, err))
	}

	defer func(){
		if err := out.Close(); err != nil {
			panic(fmt.Sprintf("Tried to close output file %s, but got error: %v", outputFileName, err))
		}
	} ()

	return
}



/*
 Encode accepts the name of the file to be compresed, as well as the name 
 the resulting file, and runs the LZW compresion algorithm to compress the 
 input file 

 inputFileName - file to be compressed
 outputFileName - name of resulting, compressed file
*/
func Encode(inputFileName, outputFileName string) {

	//
	// open input file 
	inFile, err := os.OpenFile(inputFileName, os.O_RDONLY, 677)
	if err != nil {
		panic(fmt.Sprintf("Tried to open input file but got error: %v", err))
	}

	// defering the annonymous function, both defining it here, and invoking it - pretty cool :)
	defer func(){
		if 	err := inFile.Close(); err != nil {
			panic(fmt.Sprintf("Tried closing input file %s but got error: %v", inputFileName, inFile))
		}
	} ()

	//
	// open output file
	outFile, err := os.Create(outputFileName)
	if err != nil {
		panic(fmt.Sprintf("Tried to open output file %s, but got error: %v", outputFileName, err))
	}

	defer func(){
		if err := outFile.Close(); err != nil {
			panic(fmt.Sprintf("Tried to close output file %s, but got error: %v", outputFileName, err))
		}
	} ()

	runEncoding(inFile, outFile)
}

/*
 Decode accepts the name of a compresed file, as well as the name 
 the newly uncompresed file, and runs the reverse of the LZW compresion 
 algorithm to uncompress the input file 

 compressedFileName - file to be uncompressed
 resultingFileName - name of resulting file
*/
func Decode(compressedFileName, resultFileName string) {
	// open the files
	compressedFile, err := os.OpenFile(compressedFileName, os.O_RDONLY, 644)
	if err != nil {
		panic(fmt.Sprintf("Tried to open input file but got error: %v", err))
	}

	// defering the annonymous function, both defining it here, and invoking it - pretty cool :)
	defer func(){
		if 	err := compressedFile.Close(); err != nil {
			panic(fmt.Sprintf("Tried closing input file %s but got error: %v", compressedFileName, err))
		}
	} ()

	resultFile, err := os.Create(resultFileName)
	if err != nil {
		panic(fmt.Sprintf("Tried to open output file %s, but got error: %v", resultFileName, err))
	}

	defer func(){
		if err := resultFile.Close(); err != nil {
			panic(fmt.Sprintf("Tried to close output file %s, but got error: %v", resultFileName, err))
		}
	} ()

	runDecoding(compressedFile, resultFile)
}

/* 
  Encode helper function.  Reads in a fixed buffer from 
  the inStream, runs the LZW algorithm and writes the 
  resulting codes to the outStream
 
  TODO: currenlty this function only processes the first
  1024 bytes - make this support variable length files  
*/
func runEncoding(inStream, outStream *os.File) {

	dictionary := getEncodeDictionary()
	nextCode := len(dictionary)

	// 1024 should be a constant, even after the buffer
	// is made as a sliding window
	buf := make([]byte, 1024)

	bytesRead := readFromFile(buf, inStream)

	// prev and cur will keep track of the substring we match in the dictionary 
	// until we're ready to print a code to the outfile (and consequently add the new 
	// substr to the dictionary)
	prev := 0
	cur := 1
	prevSubstr := ""
	reachedEnd := false;
	for {
		// ¿how inefficient is it to call Srprintf for every substr?  str is 
		// immutable anyway, so we'd can't just append to the end, we have to create
		// a new string - should we just use []byte instead?  and if so, will the equals 
		// operator handle that?
		if cur == bytesRead {
			reachedEnd = true
		}
		
		substr := fmt.Sprintf("%s", buf[prev:cur])
		if dictionary[substr] == 0 {
			// we've found the new substr to add to dictionary; 
			// add it to the dictionary and append the prev code to the output
			dictionary[substr] = nextCode
			nextCode ++
			// this should have error checking around it, but that'll be added later, let's first
			// confirm the algorithm
			io.WriteString(outStream, fmt.Sprintf("%03d", dictionary[prevSubstr]))
			prev = cur - 1
		
		} else {
			
			if reachedEnd {
				// dump last stuff to buffer 
				io.WriteString(outStream, fmt.Sprintf("%03d",dictionary[substr]))
				break
			} else {
				cur ++
			}
		}
		prevSubstr = substr
		
	}
}

/* 
  Decode helper function.  Reads in a fixed buffer from 
  the inStream, runs the reverse of the LZW algorithm to 
  decode the file and writes the resulting uncompressed \
  text to the outStream
 
  TODO: currenlty this function only processes the first
  1024 bytes - make this support variable length files  
*/
func runDecoding(inStream, outStream *os.File) {

	dictionary := getDecodeDictionary()
	nextCodeVal := len(dictionary)

	buf := make([]byte, 1024)

	bytesRead, err := inStream.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	curCode := 0
	nextCode := 0

	for i:= 0; i < bytesRead; i+= 3 {	

		curCode = getNextCode(buf, i)
		curString := dictionary[curCode]

		if (i + 3) < bytesRead {	
			nextCode = getNextCode(buf, i + 3)
		} else {
			nextCode = 0
		}

		nextString := ""
		if dictionary[nextCode] != "" {
			nextString = dictionary[nextCode]
		} else {
			nextString = curString
		}  

		newString := fmt.Sprintf("%s%s", curString, nextString[:1])

		dictionary[nextCodeVal] = newString
		nextCodeVal ++
		
		io.WriteString(outStream, dictionary[curCode])
	}
}

/*
 * Reads the compressed (encoded) file and returns the next code as an integer
 */
func getNextCode(buf []byte, pos int) (code int) {
	if pos <  0 || pos >= (len(buf) - 1) {
		// we're trying to read past the end of the buffer,
		// so throw an error 
		//
		// TODO: perhaps it'd be nice to pass in  a custom error message to 
		// print in this case, since this function on its own doesn't have 
		// enough context what buffer its reading from)
		panic("Index is either past the end of buffer or negative!  Bailing")
	}
	codeStr := fmt.Sprintf("%c%c%c", buf[pos], buf[pos + 1], buf[pos + 2])
	code, err := strconv.Atoi(codeStr)

	if err != nil {
		panic(fmt.Sprintf("Trying to convert %s to an int, but got error %v", code, err))
	}

	return 
}

func readFromFile(buf []byte, stream *os.File) (bytesRead int) {
	bytesRead, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
	return
}


