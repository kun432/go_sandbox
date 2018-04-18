package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	// 引数処理
	args := os.Args
	fmt.Println(args)
	if len(args[1:]) != 1 {
		fmt.Println("Error! irregal arguments.")
		usage()
		os.Exit(1)
	}

	// リストファイルオープン
	r, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(2)
	}
	defer r.Close()

	var iplist []string

	// リスト読み込み
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line[:1] == "#" {
			continue
		}
		iplist = append(iplist, line)
	}

	// カウンタ初期化
	//cnt_total := len(iplist)
	//cnt_success := 0
	//cnt_fail := 0

	// ping実行
	var wg sync.WaitGroup
	for _, ip := range iplist {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			worker(ip)
		}(ip)
	}
	wg.Wait()

	// 結果出力
	//fmt.Printf("%d件中、%d成功しました\n", cnt_total, cnt_success)
}

func usage() {
	fmt.Println("[Usage] expingo IP_LIST_FILE")
}

func worker(ip string) {
	err := exec.Command("ping", "-c", "1", "-W", "2", ip).Run()

	if err != nil {
		fmt.Println(ip, " : NG")
	} else {
		fmt.Println(ip, " : OK")
	}
}
