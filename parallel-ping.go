// TODO
// - 同時実行数の制御
// - 結果の表示
// - ネイティブにping実行

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

func main() {
	// 引数処理
	args := os.Args
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

	// ping実行
	var wg sync.WaitGroup
	cpus := runtime.NumCPU()
	c := make(chan int, cpus)
	for _, ip := range iplist {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			c <- 1
			worker(ip)
			<-c
		}(ip)
	}
	wg.Wait()
}

func usage() {
	fmt.Println("[Usage] parallel-ping IP_LIST_FILE")
}

func worker(ip string) {
	err := exec.Command("ping", "-c", "1", "-W", "1", ip).Run()

	if err != nil {
		fmt.Println(ip, " : NG")
	} else {
		fmt.Println(ip, " : OK")
	}
}
