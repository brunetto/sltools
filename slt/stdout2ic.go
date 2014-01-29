package slt

import (
	"path/filepath"
	"strings"
)


// Create the output file name that will be the new IC for the restart
func Stdout2Ic (inFileName, fileN string) (outFileName string) {
	var (
		extension string
		baseName string
		file string
		dir string
	)
	
	dir = filepath.Dir(inFileName)
	file = filepath.Base(inFileName)
	extension = filepath.Ext(inFileName)
	baseName = strings.TrimSuffix(file, extension)
	baseName = strings.TrimPrefix(baseName, "new_")
	outFileName = filepath.Join(dir, baseName) + "-IC-" + fileN + extension //FIXME detectare nOfFiles
	return outFileName
}