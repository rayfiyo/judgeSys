package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const ( // エラーコード
	argErr int = iota + 1
	cmdRunArg
	cmdRunShell
	testCaseOpen
	testCaseConvAtoi
	errProcessManyArg
	errProcessOpen
	errProcessMessage
)

func main() {
	if len(os.Args) < 2 {
		errProcess(nil, argErr, "実行時の引数不足")
	}
	fileName := os.Args[1]                  // os.Args は 値を保持しているので、二回呼び出しではない
	if fileName[len(fileName)-2:] == ".c" { // .c ついてたら消す
		fileName = fileName[:len(fileName)-2]
	}

	exec.Command("bash", "-c", "echo \""+"1"+"\" > "+fileName+".txt").Run() // 依存無し,エラー詳細無しでWA(= 1)と書き込み

	cmdRun(fileName, "/usr/bin/gcc", fileName+".c", "-lm", "-o", fileName+".out")
	for _, oneCase := range testCase(fileName) {
		cmdRun(fileName, "bash", "-c", "echo "+oneCase+" | ./"+fileName+".out >> "+fileName+".txt")
		// exec.Command("bash", "-c", "echo \""+oneCase+"\" | ./"+fileName+".out >> "+fileName+".txt").Run() // 変形無しで書き込み
	}
	cmdRun(fileName, "diff", "-q", "ans.txt", fileName+".txt")          // 回答との差分比較
	cmdRun(fileName, "bash", "-c", "sed -i '1s/1/0/' "+fileName+".txt") // ACしたことを書き込み(= 先頭の1を0にする)
	cmdRun(fileName, "rm", fileName+".c", fileName+".out")              // 副産物ファイル削除
	fmt.Println("End of Program: Judge system done, and AC.", " <br>")
}

func cmdRun(fileName string, cmd ...string) {
	var fullCmd *exec.Cmd
	cRunning := false
	if len(cmd) == 1 {
		fullCmd = exec.Command(cmd[0])
	} else if len(cmd) > 1 {
		fullCmd = exec.Command(cmd[0], cmd[1:]...)
		if strings.Contains(cmd[2], "echo ") && strings.Contains(cmd[2], ".out >> ") {
			cRunning = true
		}
	} else { // len(cmd) < 1
		log.Panicf("Error Code: %d\nError Message: cmdRunの引数不足", cmdRunArg)
	}
	output, err := fullCmd.CombinedOutput()

	if err != nil {
		errText := err.Error()
		if cRunning {
			printCmd := fullCmd.String()
			addMessage := "ignore:   " + errText + " <br>\n" +
				"full cmd: " + printCmd + " <br>\n"
			fmt.Println(addMessage)
		} else {
			printCmd := fullCmd.String()
			printOutput := strings.TrimRight(string(output), "\n") // 望まない改行削除
			addMessage := "シェルコマンド実行エラー <br>\n" +
				"full cmd: " + printCmd + " <br>\n" +
				"output:   " + printOutput
			errProcess(err, cmdRunShell, addMessage, fileName)
			defer fmt.Println("End of Program: Judge system done, but WA.", " <br>")
		}
	}
}

func testCase(fileName string) []string {
	file, err := os.Open("sample.txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		errProcess(err, testCaseOpen, "テストケースファイルの展開に失敗", fileName)
	}

	sc := bufio.NewScanner(file)
	sc.Scan()
	size, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println(err)
		errProcess(err, testCaseConvAtoi, "入力文字をint型に変換失敗", fileName)
	}

	var sample []string
	var currentCase string
	for i := 0; i < size; { //入力できる間
		sc.Scan()
		text := sc.Text()
		if text == "" {
			if len(currentCase) > 0 {
				sample = append(sample, currentCase[:len(currentCase)-1]) // ケツのスペース消す
				currentCase = ""
				i++
			}
		} else {
			currentCase += text + " "
		}
	}

	if err := sc.Err(); err != nil {
		errProcess(err, errProcessOpen, "テストケース読み込みエラー", fileName)
	}

	return sample
}

func errProcess(defaultErr error, errCode int, nameOrMessage ...string) {

	var errMessage, fileName string
	if len(nameOrMessage) == 1 {
		errMessage = nameOrMessage[0]
	} else if len(nameOrMessage) == 2 {
		errMessage = nameOrMessage[0]
		fileName = nameOrMessage[1]
	} else {
		log.Panicf("Error Code: %d\nError Message: エラー処理の引数が多い\n", errProcessManyArg)
	}
	errText := "<br>\n" +
		"err:           " + defaultErr.Error() + " <br>\n" +
		"Error Code:    " + fmt.Sprint(errCode) + " <br>\n" +
		"Error Message: " + errMessage + " <br>\n"

	exec.Command("rm", fileName+".c", fileName+".out").Run() // 副産物ファイル削除

	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		log.Panicf("Error Code: %d\nError Message: エラーを書き込むファイルの展開に失敗\n", errProcessOpen)
	}

	_, err = file.WriteString(fmt.Sprintf(errText))
	if err != nil {
		fmt.Println(err)
		log.Panicf("Error Code: %d\nError Message: エラーメッセージの書き込みに失敗\n", errProcessMessage)
	}

	log.Panicf(errText)
}
