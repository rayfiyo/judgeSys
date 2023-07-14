package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	//"strings"
)

func main() {
	var fileName string
	fileName = (os.Args[1])
	// .c ついてたら消す、len() は 1-based index に注意
	if fileName[len(fileName)-2:] == ".c" {
		fileName = fileName[:len(fileName)-2]
	}
	// .Run() は、エラー(エラーコード)のみ出力
	exec.Command("bash", "-c", "echo \""+"1"+"\" > "+fileName+".txt").Run() // 依存無しでWA(= 1)と書き込み

	cmdRun("fileName", "/usr/bin/gcc", fileName+".c", "-lm", "-o", fileName+".out")
	for _, testCase := range testCase(fileName) {
		cmdRun(fileName, "bash", "-c", "echo \""+testCase+"\" | ./"+fileName+".out >> "+fileName+".txt")
		// exec.Command("bash", "-c", "echo \""+testCase+"\" | ./"+fileName+".out >> "+fileName+".txt").Run() // 変形無しで書き込み
	}
	cmdRun(fileName, "diff", "-q", "ans.txt", fileName+".txt")          // 回答との差分比較
	cmdRun(fileName, "bash", "-c", "sed -i '1s/1/0/' "+fileName+".txt") // ACしたことを書き込み(= 先頭の1を0にする)
	cmdRun(fileName, "rm", fileName+".c", fileName+".out")              // 副産物ファイル削除
	fmt.Printf("Judge system done, and AC.")
	fmt.Println("<br>")
}

func cmdRun(fileName string, cmd ...string) {
	var fullCmd *exec.Cmd
	if len(cmd) == 1 {
		fullCmd = exec.Command(cmd[0])
	} else {
		fullCmd = exec.Command(cmd[0], cmd[1:]...)
	}
	fullCmd.CombinedOutput() /*
		output, err := fullCmd.CombinedOutput()
		if err != nil {
			printCmd := fullCmd.String()                           // 型変換
			printOutput := strings.TrimRight(string(output), "\n") // 改行削除
			addMessage := "fullCmd: " + printCmd + "<br>\n" + "output: " + printOutput + "<br>\n" + "err: " + err.Error() + "<br>\n"
			errProcess(err, fileName, "err1: シェルコマンド実行エラー", addMessage)
			log.Fatal("Judge system done, but WA.")
		}
		fmt.Println(string(output))
		//*/
}

func testCase(fileName string) []string {
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

	/*sample := make([]string, size)
	for i := 0; i < size; { //入力できる間
		sc.Scan()
		text := sc.Text()
		if text != "" {
			sample[i] += text + " "
		} else if len(sample[i]) > 0 { // 1ケースが終わったとき
			sample[i] = sample[i][:len(sample[i])-1] // ケツのスペース消す
			i++
		}
	}*/

	if err := sc.Err(); err != nil {
		fmt.Printf("err4: テストケース読み込みエラー")
		fmt.Println("<br>")
		errProcess(err, fileName, "err4: テストケース読み込みエラー", "")
		log.Fatal("エラー処理終了")
	}

	return sample
}

func errProcess(errMessage error, fileName, errCode, addMessage string) {
	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		cmdRun(fileName, "rm", fileName+".c", fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err5: エラーを書き込むファイル展開に失敗<br>")
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintln(errMessage))
	_, err = file.WriteString(fmt.Sprintln(addMessage))
	if err != nil {
		cmdRun(fileName, "rm", fileName+".c", fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err6: エラーメッセージの書き込み失敗<br>")
	}
	_, err = file.WriteString(fmt.Sprintln(errCode))
	if err != nil {
		cmdRun(fileName, "rm", fileName+".c", fileName+".out") // 副産物ファイル削除
		fmt.Println(err)
		log.Fatal("err7: エラーコードの書き込み失敗<br>")
	}

	cmdRun(fileName, "rm", fileName+".c", fileName+".out") // 副産物ファイル削除
}
