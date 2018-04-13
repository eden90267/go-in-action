# Chap 02. 快速開始一個 Go 程序

> 學習寫一個複雜的 Go 程序
> 聲明類型、變數、函數和方法
> 啟動並同步操作 goroutine
> 使用接口寫通用的代碼
> 處理程序邏輯和錯誤

為了能更高效使用語言進行編碼，Go 語言有自己的哲學和編程習慣。

Go 語言的設計者們從編程效率出發設計了這門語言，但又不會丟掉訪問底層程序結構的能力。設計者們透過最少的關鍵字、內置的方法和語法，最終平衡了這兩方面。

Go 語言也提供了完善的標準庫。標準庫提供了建構實際的基於 Web 和基於網絡的程序所需的所有核心庫。

現在透過一個完整的 Go 語言程序，看 Go 如何實現這些功能：

- 從不同數據源拉取數據，將數據內容與一組搜索項做對比，然後將匹配的內容顯示在終端窗口。

這個程序會讀取文本文件，進行網絡調用，解碼 XML 和 JSON 成為結構化類型數據，並且利用 Go 語言的併發機制保證這些操作的速度。

## 程序架構

![](https://i.imgur.com/UyGF6Cg.png)

這個程序分成多個不同步驟，在多個不同的 goroutine 里運行。我們會根據流程展示代碼，從主 goroutine 開始，一直到執行搜索的 goroutine 和跟蹤結果的 goroutine，最後回到主 goroutine。

以下是整個項目的結構：

```
- data
  data.json  -- 包含一組數據源
- matchers
  rss.go     -- 搜索 rss 源的匹配器
- search
  default.go -- 搜索數據用的默認匹配器
  feed.go    -- 用於讀取 json 數據文件
  match.go   -- 用於支持不同匹配器的接口
  search.go  -- 執行搜索的主控制邏輯
main.go      -- 程序的入口
```

- 文件夾 data 的 JSON 文檔，其內容是程序要拉取和處理的數據源
- 文件夾 matchers 中包含程序裡用於支持搜索不同數據源的代碼。目前程序只完成了支持處理 RSS 類型的數據源的匹配器
- 文件夾 search 中包含使用不同匹配器進行搜索的業務邏輯。
- main.go 是整個程序入口

## main 包

```go
package sample

import (
  "log"
  "os"

  _ "github.com/eden90267/go-in-action/cn02.quick-start/sample/matchers"
  "github.com/eden90267/go-in-action/cn02.quick-start/sample/search"
)

// init 在 main 之前調用
func init() {
  // 將日誌輸出到標準輸出
  log.SetOutput(os.Stdout)
}

func main() {
  // 使用特定的項做搜索
  search.Run("president")
}
```

Go 程序都有兩個明顯特徵：

1. main 函數。構建程序在構建可執行文件時，需要找到這個已經聲明的 main 函數，把它作為程序的入口
2. 包名 main

    如果 main 函數不在 main 包裡，構建工具就不會生成可執行的文件

Go 語言的每一個代碼文件都屬於一個包，main.go 也不例外。包這個特性對於 Go 語言來說很重要，會在第三章接觸到更多細節。現在先簡單瞭解以下內容：一個包定義一組編譯過的代碼，包的名字類似命名空間，可以用來間接訪問包內聲明的標飾符。這個特性可以把不同包中定義的同名標飾符區別開。
