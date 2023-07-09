package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	var fileName string
	fmt.Printf("ac.c を実行します:")
	fileName = "ac"

	cmdRun("gcc " + fileName + ".c -o " + fileName + ".out")
	cmdRun("echo \"" + "1" + "\" > " + "ans.txt") // ans.txt の初期化(比較時はWAなので1)
	for _, testCase := range testCase() {
		cmdRun("echo \"" + string(testCase) + "\" | ./" + fileName + ".out >> " + "ans.txt")
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
	sc.Scan()
	size, err := strconv.Atoi(sc.Text())
	if err != nil {
		panic(err)
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
		log.Fatal("err2: テストケース読み込みエラー")
	}

	return sample
}
