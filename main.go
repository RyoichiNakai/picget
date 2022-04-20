package main

import (
  "os"
  "os/signal"
  "flag"
  "syscall"
  "context"
  "fmt"
	"time"

	"picget/downloader"
)

func main() {
  // 初期化
	var url string
	var split int
	var path string
	flag.StringVar(&url, "h", "http://localhost:8080", "requst url")
	flag.IntVar(&split, "d", 5, "number of download divisions")
	flag.StringVar(&path, "f", "output/sample.jpg", "output path (only jpg)")
	flag.Parse()

  // クライアントの初期化
  dc := &downloader.Client{
    Url: url,
    Split: split,
    Path: path,
		ByteData: make([][]byte, split),
  }

	// キャンセル処理の設定
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		// シグナルの受付を終了する
		signal.Stop(sigs)
		cancel()
	}()

	go func() {
		select {
		// シグナルを受け取ったらここに入る
		case sig := <-sigs:
			fmt.Println("\nGot signal!", sig)
			// cancelを呼び出して全ての処理を終了させる
			cancel()
		}
	}()

	// ダウンロードの開始
  if err := Run(ctx, dc); err != nil {
		cancel()
		fmt.Println(err)
	}
}

/*
ダウンロードの実行
*/
func Run(ctx context.Context, dp downloader.PicGet) error {
	now := time.Now() 
	
  if err := dp.ValidClient(); err != nil {
    return err
  }

  if err := dp.Download(ctx); err != nil {
    return err
  }

	if err := dp.FileMerge(ctx); err != nil {
		return err
	}

	defer fmt.Printf("[info] 経過: %vms\n", time.Since(now).Milliseconds()) 

  return nil
}
