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
	cmdRunArgErr
	shellErr
	errProcessOpen
	errProcessMessage
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Error Code: %d\nError Message: 実行時の引数不足", argErr)
	}
	fileName := os.Args[1]                  // os.Args は 値を保持しているので、二回呼び出しではない
	if fileName[len(fileName)-2:] == ".c" { // .c ついてたら消す
		fileName = fileName[:len(fileName)-2]
	}

	exec.Command("bash", "-c", "echo \""+"1"+"\" > "+fileName+".txt").Run() // 依存無し,エラー詳細無しでWA(= 1)と書き込み

	cmdRun(fileName, "/usr/bin/gcc", fileName+".c", "-lm", "-o", fileName+".out")
	for _, testCase := range testCase(fileName) {
		cmdRun(fileName, "bash", "-c", "echo "+testCase+" | ./"+fileName+".out >> "+fileName+".txt")
		// exec.Command("bash", "-c", "echo \""+testCase+"\" | ./"+fileName+".out >> "+fileName+".txt").Run() // 変形無しで書き込み
	}
	cmdRun(fileName, "diff", "-q", "ans.txt", fileName+".txt")          // 回答との差分比較
	cmdRun(fileName, "bash", "-c", "sed -i '1s/1/0/' "+fileName+".txt") // ACしたことを書き込み(= 先頭の1を0にする)
	cmdRun(fileName, "rm", fileName+".c", fileName+".out")              // 副産物ファイル削除
	fmt.Printf("Judge system done, and AC.\n<br>\n")
}

func cmdRun(fileName string, cmd ...string) {
	var fullCmd *exec.Cmd
	if len(cmd) < 1 {
		log.Fatalf("Error Code: %d\nError Message: cmdRunの引数不足", cmdRunArgErr)
	} else if len(cmd) == 1 {
		fullCmd = exec.Command(cmd[0])
	} else {
		fullCmd = exec.Command(cmd[0], cmd[1:]...)
	}
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ExitCode() == 3 { // exit status 3 のとき、エラーで強制終了しない
			printCmd := fullCmd.String()                           // 型変換
			printOutput := strings.TrimRight(string(output), "\n") // 改行削除
			addMessage := "fullCmd: " + printCmd + "<br>\n" + "err: " + err.Error() + "<br>\n" + "output: " + printOutput + "<br>\n"
			fmt.Println(addMessage)
			fmt.Printf("output:")
		} else {
			printCmd := fullCmd.String()                           // 型変換
			printOutput := strings.TrimRight(string(output), "\n") // 改行削除
			addMessage := "fullCmd: " + printCmd + "<br>\n" + "err: " + err.Error() + "<br>\n" + "output: " + printOutput + "<br>\n"
			errProcess(err, fileName, "err1: シェルコマンド実行エラー", addMessage)
			defer fmt.Println("Error Message: Judge system done, but WA.")
			log.Fatalf("Error Code: %d\nError Message: シェルがエラーで終了", shellErr)
		}
	}
	fmt.Println(output)
}

func testCase(fileName string) []string {
	file, err := os.Open("sample.txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error Code: %d\nError Message: テストケースファイルの展開に失敗", testCaseOpen)
	}

	sc := bufio.NewScanner(file)
	sc.Scan()
	size, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error Code: %d\nError Message: 入力文字をint型に変換失敗", testCaseConvAtoi)
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
		errProcess(fileName, err, errProcessOpen, "Error Message: テストケース読み込みエラー\n")
		log.Fatal("エラー処理終了")
	}

	return sample
}

func errProcess(fileName string, defaultErr error, errCode int, errMessage string) {
	defer exec.Command("rm", fileName+".c", fileName+".out").Run() // 副産物ファイル削除
	errText := "Error " + fmt.Sprint(errCode) + ": " + errMessage

	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error Code: %d\nError Message: エラーを書き込むファイルの展開に失敗\n", errProcessOpen)
	}

	_, err = file.WriteString(fmt.Sprintln(errCode, errMessage, addMessage))
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error Code: %d\nError Message: エラーメッセージの書き込みに失敗\n", errProcessMessage)
	}
	
	log.Fatalln(errCode, errMessage, addMessage)
}
