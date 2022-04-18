package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/sync/errgroup"
	//"downloader"
)

func main() {
	// 実行時にURLを取得
	var url string
	var divideNum int
	var outputPath string
	flag.StringVar(&url, "h", "http://localhost:8080", "requst url")
	flag.IntVar(&divideNum, "d", 5, "number of download divisions")
	flag.StringVar(&outputPath, "o", "output/sample.jpg", "output path (only jpg)")
	flag.Parse()

	/**
	 * そのサイトが分割ダウンロードに対応しているかの確認
	 * 対応している場合は、コンテンツのLengthを取ってくる
	 */
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	// クライアントの設定
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// Range Accessの存在チェック
	if _, ok := res.Header["Accept-Ranges"]; !ok {
		fmt.Println(ok)
	}

	contentLength, err := strconv.Atoi(res.Header["Content-Length"][0])
	if err != nil {
		fmt.Println(err)
	}

	/**
	 * byte数を選択し、分割ダウンロード
	 * これをgoroutineにて行う
	 */

	// 1アクセスごとにByte数を分割する
	byteRangeArray := make([]string, divideNum)
	chank := contentLength / divideNum
	for i := 0; i < divideNum; i++ {
		begin := chank * i
		if i == divideNum-1 {
			byteRangeArray[i] = fmt.Sprintf("%d-%d", begin, contentLength)
		} else {
			end := chank*(i+1) - 1
			byteRangeArray[i] = fmt.Sprintf("%d-%d", begin, end)
		}
	}

	/**
	 * ファイルダウンロード処理
	 */
	eg, ctx := errgroup.WithContext(context.Background())

	byteDataArrary := make([][]byte, divideNum)
	for i := 0; i < divideNum; i++ {
		i := i
		eg.Go(func() error {
			fmt.Printf("[info] %d番目のダウンロード開始\n", i+1)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Set("Range", "bytes="+byteRangeArray[i])

			client := &http.Client{}
			res, err := client.Do(req)

			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			if err != nil {
				fmt.Println(err)
			}

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			}

			byteDataArrary[i] = data
			fmt.Printf("[info] %d番目のダウンロード完了\n", i+1)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}

	/**
	 * byte列をファイルに書き出す
	 */

	buf := &bytes.Buffer{}

	for i := 0; i < divideNum; i++ {
		err := binary.Write(buf, binary.BigEndian, byteDataArrary[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(file)
	}
	defer file.Close()

	fmt.Println("[info] ファイルへ書き込み中です。。。")
	file.Write(buf.Bytes())

	fmt.Println("[info] ファイルへの書き込みが完了しました！")
}
