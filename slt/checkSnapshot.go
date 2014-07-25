package slt

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
)

func CheckSnapshot(inFileName string) {
	var (
		err     error
		inFile  *os.File
		nReader *bufio.Reader
		fZip    *gzip.Reader
		ext     string
	)
	// 	log.Println("Checking ", inFileName)
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	ext = filepath.Ext(inFileName)

	switch ext {
	case ".txt":
		{
			nReader = bufio.NewReader(inFile)
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
		}
	default:
		{
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extention ", ext)
		}
	}

	for {
		if _, err = ReadOutSnapshot(nReader); err != nil {
			break
		}
	}
}
