package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {

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
        case sig := <-sigs: // シグナルを受け取ったらここに入る
            fmt.Println("Got signal!", sig)
            cancel() // cancelを呼び出して全ての処理を終了させる
        }
    }()

    if err := doTask(ctx); err != nil {
        fmt.Printf("failed to doTask: %v", err)
        cancel() // 何らかのエラーが発生した場合、他の処理も全てcancelさせる
        return
    }
    fmt.Println("done successfully.")
}

func doTask(ctx context.Context) error {
    defer fmt.Println("done doTask")
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("received done")
            return ctx.Err()
        default:
        }
        // // エラー時の挙動が見たい場合はここのコメントアウトを外す
        // if i == 3 {
        //  return fmt.Errorf("error happened")
        // }

        // do something
        fmt.Println("sleep 1. count:", i)
        time.Sleep(1 * time.Second)
    }
    return nil
}