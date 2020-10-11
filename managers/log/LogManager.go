package log

import (
	"ProjetoUnivesp2020/managers/config"
	"ProjetoUnivesp2020/utils"
	"bufio"
	"fmt"
	"os"
	"time"
)

var LogManager = createLog()

type logLevel struct {
	Name    string
	LogFile *os.File
}

func (l logLevel) AppendLog(log string) {
	if _, e := l.LogFile.WriteString(formatLog(log)); e != nil {
		utils.CheckPanic(&e)
	}
}

func (l logLevel) ClearLog() {
	e := l.LogFile.Truncate(0)

	utils.CheckPanic(&e)

	_, e = l.LogFile.Seek(0, 0)

	utils.CheckPanic(&e)

	e = l.LogFile.Sync()

	utils.CheckPanic(&e)
}

func (l logLevel) ReadLog() []string {
	var lines []string

	scanner := bufio.NewScanner(l.LogFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

type logManagerType []logLevel

func (l logManagerType) GetLogLevel(nome string) *logLevel {
	for i := 0; i < len(l); i++ {
		if l[i].Name == nome {
			return &l[i]
		}
	}

	return nil
}

func (l logManagerType) GetAllLogsLevels() []string {
	levels := make([]string, len(l))

	for i := 0; i < len(l); i++ {
		levels[i] = l[i].Name
	}

	return levels
}

func createLog() *logManagerType {
	if _, err := os.Stat(config.Configs.LogPath); os.IsNotExist(err) {
		err = os.MkdirAll(config.Configs.LogPath, 0777)

		utils.CheckPanic(&err)
	}

	path := config.Configs.LogPath + "/"

	return &logManagerType{
		logLevel{
			Name:    "warning",
			LogFile: openLog(path + "warning.log"),
		},
		logLevel{
			Name:    "infoAdmin",
			LogFile: openLog(path + "infoAdmin.log"),
		},
		logLevel{
			Name:    "docker",
			LogFile: openLog(path + "docker.log"),
		},
		logLevel{
			Name:    "infoUser",
			LogFile: openLog(path + "info.log"),
		},
	}
}

func openLog(path string) *os.File {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		_, e = os.Create(path)

		utils.CheckPanic(&e)
	}
	f, e := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

	utils.CheckPanic(&e)

	return f
}

func formatLog(log string) string {
	return fmt.Sprintf("%d,%s\n", time.Now().Unix(), log)
}
