package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var fileName string
	fileName = "ac"

	fmt.Println("ls コマンド:")
	rawOutput, _ := exec.Command("ls").Output()
	fmt.Println(strings.Join(strings.Fields(string(bytes.TrimRight(rawOutput, "\n"))), "  ")) // 整形

	fmt.Println("ディレクトリの構成例(makeAns抜き): ")
	fmt.Println("ac.c  index.php  sample.txt  send.php")

	fmt.Printf("何か入力後、実行: ")
	var input string
	fmt.Scanln(&input)

	fmt.Printf("実行コマンド: ")
	cmdRun("fileName", "/usr/bin/gcc", fileName+".c", "-lm", "-o", fileName+".out")
	cmdRun(fileName, "bash", "-c", "echo \""+"1"+"\" > "+"ans.txt") // ans.txt の初期化(比較時はWAなので1)
	for _, testCase := range testCase() {
		cmdRun(fileName, "bash", "-c", "echo \""+testCase+"\" | ./"+fileName+".out >> "+"ans.txt")
	}
	cmdRun(fileName, "rm", fileName+".out") // 実行ファイル削除
	fmt.Println("ans.txt の作成完了")
	fmt.Println("ディレクトリの構成例(makeAnsあり)")
	fmt.Println("ac.c  ans.txt  index.php  makeAns  sample.txt  send.php")
}

func cmdRun(fileName string, cmd ...string) {
	var fullCmd *exec.Cmd
	if len(cmd) == 1 {
		fullCmd = exec.Command(cmd[0])
	} else {
		fullCmd = exec.Command(cmd[0], cmd[1:]...)
	}
	fmt.Println(fullCmd)
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		printCmd := fullCmd.String()                           // 型変換
		printOutput := strings.TrimRight(string(output), "\n") // 改行削除
		addMessage := "fullCmd: " + printCmd + "<br>\n" + "output: " + printOutput + "<br>\n" + "err: " + err.Error() + "<br>\n"
		fmt.Println(err, fileName, "err1: シェルコマンド実行エラー", addMessage)
		log.Fatal("Judge system done, but WA.")
	}
	fmt.Println(string(output))
}

func testCase() []string {
	file, err := os.Open("sample.txt")
	if err != nil {
		fmt.Println(err)
		log.Fatal("err2: エラーを書き込むファイル展開に失敗")
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	sc.Scan()
	size, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println(err)
		log.Fatal("err3: 入力文字をint型に変換失敗")
	}

	sample := make([]string, size)

	for i := 0; i < size; { //入力できる間
		sc.Scan()
		text := sc.Text()
		if text != "" {
			sample[i] += text + " "
		} else if len(sample[i]) > 0 { // 1ケースが終わったとき
			sample[i] = sample[i][:len(sample[i])-1] // ケツのスペース消す
			i++
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatal("err4: テストケース読み込みエラー")
	}

	return sample
}
