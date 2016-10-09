package file

import (
	"os"
	"strconv"
)

const (
	HOST        = "localhost"
	PORT        = "5556"
	TYPE        = "TCP"
	BUFFER_SIZE = 1024
)

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func getFileStats(file *os.File) (string, string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return "", "", err
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)

	return fileSize, fileName, nil
}
