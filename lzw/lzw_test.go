package lzw

import (
	"testing"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func createFile(fName, contents string) {
	f, err1 := os.Create(fName)
	_, err2 := f.WriteString(contents)
	if err1 != nil || err2 != nil {
		panic("Unable to create test file!")
	} 
	f.Close();
}

func getFileContents(fName string) (text string) {
	contents, err := ioutil.ReadFile(fName)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file %s", fName))
	}
	text = string(contents)
	return
}

func singleTest(t *testing.T, testString, fPath string) {

	fInputName := fmt.Sprintf("%s/tempTest.txt", fPath)
	fEncodedName := fmt.Sprintf("%s/tempTest.enc", fPath)
	fDecodedName := fmt.Sprintf("%s/tempTest.dec", fPath)

	createFile(fInputName, testString)

	// run encode and decode
	Encode(fInputName, fEncodedName)
	Decode(fEncodedName, fDecodedName)

	// check that decoded and original are the same 
	decodedText := getFileContents(fDecodedName)
	
	if testString != decodedText {
		t.Errorf("expected file contents: \"%s\", but got: \"%s\"", testString, decodedText)
	}

	// clean up
	os.Remove(fInputName)
	os.Remove(fEncodedName)
	os.Remove(fDecodedName)
}

func TestEncodeDecode(t *testing.T) {
	curDir := filepath.Dir(".")
	testDir, err := ioutil.TempDir(curDir, "tempDir")
	if err != nil {
		panic("Unable to create the temporary test directory!")
	}

	testFileContents := []string{ 
							  "first simple text file", 
							  "second test?!  with some non-a1phaNum3r1c Chars!",
							  "Hi,\n\tThis is the third test, OK?\n\n",
							}


	for i := range testFileContents {
		singleTest(t, testFileContents[i], testDir)
	}

	os.Remove(testDir)
}
