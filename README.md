## PicGet

### 概要

画像のURLからダウンロードを並列に行う処理の開発

### 対応拡張子

- png
- jpeg

### コマンド

| オプション | 説明                         | デフォルト            | 
| ---------- | ---------------------------- | --------------------- | 
| -h         | requst url                   | http://localhost:8080 | 
| -d         | number of download divisions | 5                     | 
| -f         | output path (only jpeg/png)  | output/sample.jpg     | 

### 実行結果

#### 正常処理

```bash
vscode ➜ /workspaces/picget (main ✗) $ go run main.go -h https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/100MB.png
[Info] ダウンロード開始...
[Info] ダウンロード終了
[info] ファイルへ書き込み中です。。。
[info] ファイルへの書き込みが完了しました！
[info] 経過: 13340ms
```

#### キャンセル処理

```bash
vscode ➜ /workspaces/picget (main ✗) $ go run main.go -h https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/100MB.png
[Info] ダウンロード開始...
^C
Got signal! interrupt
context canceled
```
