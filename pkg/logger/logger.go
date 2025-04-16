package logger

import (
	"auth/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const Console = false

var mutex sync.Mutex
var fullPathToFile string

var blue = color.New(color.FgBlue).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJSON - конвертирует объект в JSON строку
func toJSON(object any) string {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		return ""
	}
	n := len(jsonByte)             //Find the length of the byte array
	result := string(jsonByte[:n]) //convert to string

	return jsonPrettyPrint(result)
}

// getFuncName - получаем имя функции из которой вызван логгер
func getFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		parseDetails := strings.Split(details.Name(), "/")
		pkgFuncName := ""
		if len(parseDetails) == 2 {
			pkgFuncParse := strings.Split(parseDetails[1], ".")
			if len(pkgFuncParse) == 2 {
				pkgFuncName = fmt.Sprintf("[P]%s=>[F]%s", pkgFuncParse[0], pkgFuncParse[1])
			} else if len(pkgFuncParse) == 3 {
				pkgFuncName = fmt.Sprintf("[P]%s=>[S]%s=>[F]%s", pkgFuncParse[0], pkgFuncParse[1], pkgFuncParse[1])
			}
		} else if len(parseDetails) == 1 {
			pkgFuncParse := strings.Split(parseDetails[0], ".")
			if len(pkgFuncParse) == 2 {
				pkgFuncName = fmt.Sprintf("[P]%s=>[F]%s", pkgFuncParse[0], pkgFuncParse[1])
			} else if len(pkgFuncParse) == 3 {
				if pkgFuncParse[0] == pkgFuncParse[1] {
					pkgFuncName = fmt.Sprintf("[P]%s=>[F]%s=>[AF]: %s", pkgFuncParse[0], pkgFuncParse[1],
						pkgFuncParse[2])
				}
			}
		} else if len(parseDetails) > 2 {
			pkgFuncParse := strings.Split(parseDetails[len(parseDetails)-1], ".")
			if len(pkgFuncParse) == 2 {
				pkgFuncName = fmt.Sprintf("[P]%s=>[F]%s", pkgFuncParse[0], pkgFuncParse[1])
			} else if len(pkgFuncParse) == 3 {
				pkgFuncName = fmt.Sprintf("[P]%s=>[S]%s=>[F]%s", pkgFuncParse[0], pkgFuncParse[1],
					pkgFuncParse[2])
			}
		}

		return pkgFuncName
	}
	return "undefined"
}

func writeToLog(str string, a ...any) error {
	mutex.Lock()
	// Открываем файл с логом
	file, err := os.OpenFile(fullPathToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		mutex.Unlock()
		file.Close()
	}()

	// Дозаписываем данные в файл
	data := []byte(fmt.Sprintf(str+"\n", a...))
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf(str+"\n", a...)

	return nil
}

// New - инициализация логгера, создание папки и файла с логом
func New(cfg *config.Config) error {
	//Файл с логом
	fileName := cfg.Paths.LogName

	//Папка куда будет складывать лог
	folderName := filepath.Join(cfg.Paths.Files, cfg.Paths.LogDir)

	//Полный путь к файлу с логом
	fullPathToFile = filepath.Join(folderName, fileName)

	// Создаем папку с правами доступа для текущего пользователя
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	// Создаем файл с логом
	file, err := os.OpenFile(fullPathToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	mutex = sync.Mutex{}

	return nil
}

// Info - лог с пометкой info
func Info(format string, a ...any) {
	funcName := fmt.Sprintf("%s ->", getFuncName())
	go func() {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Ошибка при загрузке локации: %v\n", err)
			return
		}
		date := time.Now().In(loc)
		dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.04"))
		str := fmt.Sprintf("%s %s [%s]: %s", dateFormat, funcName, blue("INFO"), format)
		if Console {
			fmt.Printf(str, a...)
			fmt.Println()
		} else {
			err := writeToLog(str, a...)
			if err != nil {
				fmt.Printf("Ошибка при записи в лог: %v\n", err)
				return
			}
		}
	}()
}

// Error - лог с пометкой error
func Error(format string, a ...any) {
	funcName := fmt.Sprintf("%s ->", getFuncName())
	go func() {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Ошибка при загрузке локации: %v\n", err)
			return
		}
		date := time.Now().In(loc)
		dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.04"))
		str := fmt.Sprintf("%s %s [%s]: %s", dateFormat, funcName, red("ERROR"), format)
		if Console {
			fmt.Printf(str, a...)
			fmt.Println()
		} else {
			err := writeToLog(str, a...)
			if err != nil {
				fmt.Printf("Ошибка при записи в лог: %v\n", err)
				return
			}
		}
	}()
}

// Warn - лог с пометкой warn
func Warn(format string, a ...any) {
	funcName := fmt.Sprintf("%s ->", getFuncName())
	go func() {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Ошибка при загрузке локации: %v\n", err)
			return
		}
		date := time.Now().In(loc)
		dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.04"))
		str := fmt.Sprintf("%s %s [%s]: %s", dateFormat, funcName, yellow("WARN"), format)
		if Console {
			fmt.Printf(str, a...)
			fmt.Println()
		} else {
			err := writeToLog(str, a...)
			if err != nil {
				fmt.Printf("Ошибка при записи в лог: %v\n", err)
				return
			}
		}
	}()
}
