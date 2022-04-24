package downloader

import (
  "os"
	"bytes"
	"encoding/binary"
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"
	"strconv"
  "fmt"
	"io/ioutil"

	"golang.org/x/sync/errgroup"
	"golang.org/x/net/context/ctxhttp"
)

type Client struct {
	Url 	   string
	Split	   int
	Path     string
	ByteData [][]byte
}

type PicGet interface {
	ValidClient() error
	Download (ctx context.Context) error
	FileMerge (ctx context.Context) error
}

/*
作成したDownloaderのクライアントを検証
@param なし
@return error
*/
func (c *Client) ValidClient() error {
	// I/Oの拡張子チェック
	var exts = [6]string{".jpg", ".jpeg", ".png"}
	urlExt := false
	pathExt := false
	
	for _, v := range exts {
		if v == filepath.Ext(c.Url) {
			urlExt = true
		}

		if v == filepath.Ext(c.Path) {
			pathExt = true
		}
	}
	
	if !urlExt || !pathExt{
		return errors.New("invalid exetention")
	}

	// 分割数チェック
	// 分割数が
	if c.Split < 1 {
		return errors.New("the number of download divisions is 1 or more")
	}

	return nil
}

/*
HTTPリクエストを行う
*/
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := &http.Client {
		Timeout: 60 * time.Second,
	}

	res, err := ctxhttp.Do(ctx, client, req)
	if err != nil {
		return nil, err
	}

  return res, nil
}

/*
HTTPリクエストを行い、RangeAccessが可能かどうかを判断する
同時にContentLengthのデータを取得する
*/
func (c *Client) canRangeAccess(ctx context.Context) (bool, int, error) {
  req, _ := http.NewRequest("HEAD", c.Url, nil)

	res, err := c.doRequest(ctx, req)
  if err != nil {
    return false, 0, err
  }
  defer res.Body.Close()

  // Range Accessできるかのチェック
  _, hasRangeAccess := res.Header["Accept-Ranges"];

	// Content Lengthがない場合は分割ダウンロードをさせない
	if _, hasContentLength := res.Header["Content-Length"]; !hasContentLength {
		return false, 0, nil
	}

	// Content Lengthが存在する場合は、
	contentLength, err := strconv.Atoi(res.Header["Content-Length"][0])
  if err != nil {
    return false, 0, err
  }

	fmt.Println(hasRangeAccess, contentLength)

  return hasRangeAccess, contentLength, nil
}

/*
Downloadを行う
Clientのsplitが2以上であれば、分割ダウンロードを行う
*/
func (c *Client) Download(ctx context.Context) error {
	eg, _ := errgroup.WithContext(ctx)
  byteRange := make([]string, c.Split)

  hasRangeAccess, contentLength, err := c.canRangeAccess(ctx)
  if err != nil {
    return err
  }

  if !hasRangeAccess {
    c.Split = 1
		fmt.Println("[Info] RangeAccessができないため、分割数を1に設定します")
  } else {
    byteRange = c.getByteRange(contentLength)
  }

	fmt.Println("[Info] ダウンロード開始...")
  for i := 0; i < c.Split; i++ {
    i := i
    eg.Go(func() error {
			req, _ := http.NewRequest("GET", c.Url, nil)

			if hasRangeAccess {
				req.Header.Set("Range", "bytes="+byteRange[i])
			}

			res, err := c.doRequest(ctx, req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			c.ByteData[i] = data

			return nil
    })
  }

  if err := eg.Wait(); err != nil {
    return err
  }

	fmt.Println("[Info] ダウンロード終了")

  return nil
}

/*
分割ダウンロードを行う際のbyteの繁栄を決定
*/
func (c *Client) getByteRange(contentLength int) []string {
  chunk := contentLength / c.Split
	byteRange := make([]string, c.Split)

	for i := 0; i < c.Split; i++ {
		begin := chunk * i
		if i == c.Split-1 {
			byteRange[i] = fmt.Sprintf("%d-%d", begin, contentLength)
		} else {
			end := chunk*(i+1) - 1
			byteRange[i] = fmt.Sprintf("%d-%d", begin, end)
		}
	}
  return byteRange
}

/*
取得したbyte列をマージして、ファイルを作成。その後、コピーする
*/
func (c *Client) FileMerge(ctx context.Context) error {
	buf := &bytes.Buffer{}

	for i := 0; i < c.Split; i++ {
		err := binary.Write(buf, binary.BigEndian, c.ByteData[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	file, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("[info] ファイルへ書き込み中です。。。")
	if _, err := file.Write(buf.Bytes()); err != nil {
		return err
	}

	fmt.Println("[info] ファイルへの書き込みが完了しました！")

  return nil
}
