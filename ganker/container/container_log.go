package container

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const LogFileName = "container.log"

func ShowContainerLog(containerId string) {
	logFilePath := path.Join(ContainerRootPath+containerId, LogFileName)
	file, err := os.Open(logFilePath)
	if err != nil {
		logrus.Errorf("Open Log File %s Error %v", logFilePath, err)
		return
	}

	defer file.Close()

	content, err1 := io.ReadAll(file)
	if err1 != nil {
		logrus.Errorf("Read Log File %s Error %v", logFilePath, err)
		return
	}
	if _, err := fmt.Fprint(os.Stdout, string(content)); err != nil {
		logrus.Errorf("Print Log File %s Error %v", logFilePath, err)
		return
	}
}

func recordContainerLog(containerId string) *os.File {

	logFilePath := path.Join(ContainerRootPath+containerId, LogFileName)
	if exist, _ := checkFileOrDirExist(logFilePath); !exist {
		if err := os.MkdirAll(path.Dir(logFilePath), 0666); err != nil {
			logrus.Errorf("Create Log File %s Error %v", logFilePath, err)
			return nil
		}

		file, err := os.Create(logFilePath)
		if err != nil {
			logrus.Errorf("Create Log File %s Error %v", logFilePath, err)
			return nil
		}
		return file
	} else {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			logrus.Errorf("Open Log File %s Error %v", logFilePath, err)
			return nil
		}
		return file
	}
}
