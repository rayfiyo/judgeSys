package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	var fileName string
	// fmt.Printf("file name:")
	fileName = (os.Args[1])                                                            // fmt.Scanln(&fileName)
	exec.Command("/usr/bin/bash", "-c", "echo \""+"0"+"\" > "+fileName+".txt").CombinedOutput() // WA(= 0)と書き込み

	cmdRun("\"gcc ./"+fileName+".c -o ./"+fileName+".out\"", fileName)
	for _, testCase := range testCase(fileName) {
		cmdRun("echo \""+testCase+"\" | ./"+fileName+".out >> "+fileName+".txt", fileName)
	}
	cmdRun("diff ./ans.txt ./"+fileName+".txt", fileName)    // 回答との差分比較
	cmdRun("sed -i 's/0/1/1' "+fileName+".txt", fileName)     // ACしたことを書き込み(= 先頭の0を1にする)
	cmdRun("rm ./"+fileName+".c "+fileName+".out", fileName) // 副産物削除
}

func cmdRun(cmd, fileName string) {
	fullCmd := exec.Command("/usr/bin/bash", "-c", cmd)
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("err1シェルコマンド実行エラー")
		fmt.Println("<br>")
		fmt.Printf("fullCmd: %s", fullCmd)
		fmt.Println("<br>")
		fmt.Printf("err: %s", err)
		fmt.Println("<br>")
		errProcess(err, fileName, "err1: シェルコマンド実行エラー")
		log.Fatal("エラー処理終了")
	}
	fmt.Println(string(output))
}

func testCase(fileName string) []string {
	file, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)
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
		fmt.Printf("err2: テストケース読み込みエラー")
		fmt.Println("<br>")
		errProcess(err, fileName, "err2: テストケース読み込みエラー")
		log.Fatal("エラー処理終了")
	}

	return sample
}

func errProcess(errMessage error, fileName, errCode string) {
	file, err := os.Open(fileName + ".txt")
	if err != nil {
		log.Fatal(err)
		log.Fatal("err3: エラーを書き込むファイル展開に失敗")
	}

	_, err = file.WriteString(fmt.Sprintln(errMessage))
	if err != nil {
		log.Fatal("err4: エラーメッセージの書き込み失敗")
	}

	_, err = file.WriteString(fmt.Sprintln(errCode))
	if err != nil {
		log.Fatal("err5: エラーコードの書き込み失敗")
	}
}
