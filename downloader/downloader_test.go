package downloader

import (
	"context"
	"testing"
	"strings"
	"net/http"
)

type clientTestCase struct {
	downloader Client
	expected	 string
}

// ValidClientの単体テスト
func TestValidClient(t *testing.T) {
	clientCases := []clientTestCase{
		{
			// 異常系：Pathの拡張子のエラー
			downloader: Client{
				Url: 	 "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
				Split: 5,
				Path:  "output/sample.pdf",
			},
			expected: "invalid exetention",
		},
		{
			// 異常系：URLの拡張子のエラー
			downloader: Client{
				Url: 	 "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.pdf",
				Split: 5,
				Path:  "output/sample.jpg",
			},
			expected: "invalid exetention",
		},
		{
			// 異常系：分割数エラー
			downloader: Client{
				Url: 	 "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
				Split: -1,
				Path:  "output/sample.jpg",
			},
			expected: "the number of download divisions is 1 or more",
		},
		{
			// 正常系
			downloader: Client{
				Url: 	 "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
				Split: 5,
				Path:  "output/sample.jpg",
			},
			expected: "",
		},
	}

	for _, c := range clientCases {
		actual := c.downloader.ValidClient(); 
		if actual != nil && actual.Error() != c.expected {
			t.Errorf("ValidClient() == %q, expect %q", actual, c.expected)
		}
	}
}

// doRequestの単体テスト
func TestDoRequest(t *testing.T) {
	// キャンセル処理にて必要
	ctx, _ := context.WithCancel(context.Background())
	
	clientCases := []clientTestCase{
		// 正常系
		{
			downloader: Client{
				Url: "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
			},
		},
		// 異常系
		{
			downloader: Client{
				Url: "http://localhost:8080",
			},
		},
	}

	for _, c := range clientCases {
		req, _ := http.NewRequest("GET", c.downloader.Url, nil)
		_, err := c.downloader.doRequest(ctx, req) 
		if err != nil && !strings.Contains(err.Error(), "connect: connection refused"){
			t.Errorf("%q", err)
		}
	}

}

// canRangeAccessの単体テスト 
func TestHasRangeAccess(t *testing.T) {
	// キャンセル処理にて使用
	ctx, _ := context.WithCancel(context.Background())
	
	clientCases := []clientTestCase{
		// 正常系：Accept Rangeが存在する
		{
			downloader: Client{
				Url: "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
			},
		},
		// 正常系：Accept-Rangeが存在しない
		{
			downloader: Client{
				Url: "https://farm8.staticflickr.com/7151/6760135001_14c59a1490_o.jpg",
			},
		},
		// 正常系：Content-Lengthが存在しない
		{
			downloader: Client{
				Url: "https://www.ricoh-imaging.co.jp/japan/dc/past/rdc/7/img/rdc7_sample01b.jpg",
			},
		},
		// 異常系：リクエストできない
		{
			downloader: Client{
				Url: "http://localhost:8080",
			},
		},
	}

	for _, c := range clientCases {
		_, _, err := c.downloader.canRangeAccess(ctx)
		if err != nil && !strings.Contains(err.Error(), "connect: connection refused"){
			t.Errorf("%q", err)
		}
	}
}

func TestDownload(t *testing.T) {
	// キャンセル処理にて使用
	ctx, _ := context.WithCancel(context.Background())

	clientCases := []clientTestCase{
		// 正常系：Accept Rangeが存在する
		// {
		// 	downloader: Client{
		// 		Url: "https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/10MB.png",
		// 		Split: 5,
		// 		ByteData: make([][]byte, 5),
		// 	},
		// },
		// // 正常系：Accept-Rangeが存在しない
		// {
		// 	downloader: Client{
		// 		Url: "https://farm8.staticflickr.com/7151/6760135001_14c59a1490_o.jpg",
		// 		Split: 5,
		// 		ByteData: make([][]byte, 5),
		// 	},
		// },
		// // 正常系：Content-Lengthが存在しない
		// {
		// 	downloader: Client{
		// 		Url: "https://www.ricoh-imaging.co.jp/japan/dc/past/rdc/7/img/rdc7_sample01b.jpg",
		// 		Split: 5,
		// 		ByteData: make([][]byte, 5),
		// 	},
		// },
		// 異常系：リクエストできない
		{
			downloader: Client{
				Url: "http://localhost:8080",
				Split: 5,
				ByteData: make([][]byte, 5),
			},
		},
	}

	for _, c := range clientCases {
		err := c.downloader.Download(ctx)
		if err != nil && !strings.Contains(err.Error(), "connect: connection refused"){
			t.Errorf("%q", err)
		}
	}
}