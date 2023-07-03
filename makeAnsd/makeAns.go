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
	fmt.Printf("ac.c を実行します:")
	fileName = "ac"

	cmdRun("gcc " + fileName + ".c -o " + fileName + ".out")
	cmdRun("echo \"" + "1" + "\" > " + "ans.txt") // ans.txt の初期化(比較時はWAなので1)
	for _, testCase := range testCase() {
		cmdRun("echo \"" + testCase + "\" | ./" + fileName + ".out >> " + "ans.txt")
	}
	cmdRun("rm ./" + fileName + ".out") // 実行ファイル削除
	fmt.Println("作成が完了しました。")
}

func cmdRun(cmd string) {
	fullCmd := exec.Command("bash", "-c", cmd)
	fmt.Println(fullCmd)
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		log.Fatal("err1: コマンド実行エラー")
	}
	fmt.Println(string(output))
}

func testCase() []string {
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
		log.Fatal("err2: テストケース読み込みエラー")
	}

	return sample
}
