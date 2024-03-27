package file

import (
	"os"
	"strings"
)

func CreateFolderForFile(filePath string) error {
	filePathSplit := strings.Split(filePath, "/")
	filePathSplit = filePathSplit[0 : len(filePathSplit)-1]
	pathWithoutFileJoined := strings.Join(filePathSplit, "/")
	return os.MkdirAll(pathWithoutFileJoined, os.ModePerm)
}
