package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"net/http"
	"strconv"

	// "golang.org/x/sync/errgourp"

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
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// Range Accessの存在チェック
	if val, ok := res.Header["Accept-Ranges"]; ok {
		fmt.Println(val)
	}

	contentLength, err := strconv.Atoi(res.Header["Content-Length"][0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(contentLength)
	
	/**
	 * byte数を選択し、分割ダウンロード
	 * これをgoroutineにて行う
	 */

	// 1アクセスごとにByte数を分割する
	byteRangeArray := make([]string, divideNum)
	accessRange := contentLength / divideNum
	for i := 0; i < divideNum; i++ {
		begin := accessRange * i
		if i == divideNum - 1 {
			byteRangeArray[i] = strconv.Itoa(begin) + "-" + strconv.Itoa(contentLength)
		} else {
			end := accessRange * (i + 1) - 1
			byteRangeArray[i] = strconv.Itoa(begin) + "-" + strconv.Itoa(end)
		}
	}

	// ファイルダウンロード
	DownloadFile(outputPath, url)

}

func DownloadFile(filepath string, url string) error {

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

		req.Header.Set("Range", "0-499")
		// クライアントの設定
		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		fmt.Println(res.Header)

    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
		defer file.Close()

    _, err = io.Copy(file, res.Body)
    return err
}

