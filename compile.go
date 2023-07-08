package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var fileName string
	fileName = (os.Args[1])
	if fileName[len(fileName)-2:] == ".c" { // len は 1-based index
		fileName = fileName[:len(fileName)-2] // len は 1-based index
	}
	exec.Command("bash", "-c", "echo \""+"1"+"\" > "+fileName+".txt").CombinedOutput() // 依存無しでWA(= 1)と書き込み

	cmdRun(fileName, "bash", "-c", "gcc ./"+fileName+".c -o ./"+fileName+".out")
	for _, testCase := range testCase(fileName) {
		cmdRun(fileName, "bash", "-c", "echo \""+testCase+"\" | ./"+fileName+".out >> "+fileName+".txt")
	}
	cmdRun(fileName, "diff", "-q", "ans.txt", fileName+".txt")             // 回答との差分比較
	cmdRun(fileName, "bash", "-c", "sed -i 's/1/0/1' "+fileName+".txt")    // ACしたことを書き込み(= 先頭の1を0にする)
	cmdRun(fileName, "bash", "-c", "rm ./"+fileName+".c "+fileName+".out") // 副産物ファイル削除
	fmt.Printf("Judge system done.")
	fmt.Println("<br>")
}

func cmdRun(fileName string, cmd ...string) {
	var fullCmd *exec.Cmd
	if len(cmd) == 1 {
		fullCmd = exec.Command(cmd[0])
	} else {
		fullCmd = exec.Command(cmd[0], cmd[1:]...)
	}
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		printCmd := fullCmd.String()                           // 型変換
		printOutput := strings.TrimRight(string(output), "\n") // 改行削除
		addMessage := "fullCmd: " + printCmd + "<br>\n" + "output: " + printOutput + "<br>\n" + "err: " + err.Error() + "<br>\n"
		errProcess(err, fileName, "err1: シェルコマンド実行エラー", addMessage)
		log.Fatal("Judge system done, but not AC.")
	}
	fmt.Println(string(output))
}

func testCase(fileName string) []string {
	file, err := os.Open("sample.txt")
	if err != nil {
		fmt.Println(err)
		log.Fatal("err2: エラーを書き込むファイル展開に失敗")
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	var sample []string

	// 二行で1空白区切りスライス
	for sc.Scan() {
		line1 := sc.Text()
		if !sc.Scan() {
			break
		}
		line2 := sc.Text()

		sample = append(sample, line1+" "+line2)
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("err3: テストケース読み込みエラー")
		fmt.Println("<br>")
		errProcess(err, fileName, "err3: テストケース読み込みエラー", "")
		log.Fatal("エラー処理終了")
	}

	return sample
}

func errProcess(errMessage error, fileName, errCode, addMessage string) {
	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		cmdRun(fileName, "bash", "-c", "rm ./"+fileName+".c "+fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err4: エラーを書き込むファイル展開に失敗<br>")
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintln(errMessage))
	_, err = file.WriteString(fmt.Sprintln(addMessage))
	if err != nil {
		cmdRun(fileName, "bash", "-c", "rm ./"+fileName+".c "+fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err5: エラーメッセージの書き込み失敗<br>")
	}
	_, err = file.WriteString(fmt.Sprintln(errCode))
	if err != nil {
		cmdRun(fileName, "bash", "-c", "rm ./"+fileName+".c "+fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err6: エラーコードの書き込み失敗<br>")
	}
}
