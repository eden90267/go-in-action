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
// main.go
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

```go
import (
  "log"
  "os"

  _ "github.com/eden90267/go-in-action/cn02.quick-start/sample/matchers"
  "github.com/eden90267/go-in-action/cn02.quick-start/sample/search"
)
```

import 就是導入一段代碼，讓用戶可以訪問其中的標飾符，如類型、函數、常量和接口。

所有處於同一個文件夾裡的代碼文件，必須使用同一個包名。按照慣例，包和文件夾同名。就像之前說的，一個包定義一組編譯後的代碼，每段代碼都描述包的一部分。

```go
  _ "github.com/eden90267/go-in-action/cn02.quick-start/sample/matchers"
```

這個技術是為了讓 Go 語言對包作初始化操作，但是並不使用包裡的標飾符。為了讓程序的可讀性更強，Go 編譯器不允許聲明導入某個包卻不使用。下劃線讓編譯器接受這類導入，並且調用對應包內的所有代碼裡定義的 init 函數。對這個程序來說，這樣做的目的是調用 matchers 包中的 rss.go 代碼文件裡的 init 函數，註冊 RSS 匹配器，以便後用。後面會展示具體的工作方式。

```go
// main.go
// init 在 main 之前調用
func init() {
  // 將日誌輸出到標準輸出
  log.SetOutput(os.Stdout);
}
```

程序中的每個代碼文件裡的 init 函數都會在 main 函數執行前調用。這個 init 函數將標準庫裡日誌類的輸出，從默認的標準錯誤 (stderr)，設置為標準輸出 (stdout) 設備。在第七章，我們會近一步討論 log 包和標準庫裡其他重要包。

```go
func main() {
  // 使用特定的項做搜索
  search.Run("president")
}
```

這個函數包含程序核心業務邏輯。一旦 Run 函數退出，程序就會終止。

## search 包

整個程序都圍繞匹配器來運作。這個程序裡的匹配器，是指包含特定信息、用於處理某類數據源的實例。在這示例有兩個匹配器。框架本身實現了一個無法獲取任何信息的默認匹配器，而在 matchers 包裡實現了 RSS 匹配器。RSS 知道如何獲取、讀入並查找 RSS 數據源。隨後會擴展這個程序，加入能讀取 JSON 文檔或 CSV 文件的匹配器。

```go
package search

import (
  "log"
  "sync"
)

// 註冊用於搜索的匹配器的映射
var matchers = make(map[string]Matcher)
```

可看到，每個代碼文件都以 package 關鍵字開頭，隨後跟著包的名字。文件夾 search 下的每個代碼文件都使用 search 作為包名。

與第三包不同，從標準庫導入代碼，只需要給出要導入的包名。編譯器查找包的時候，總是會到 GOROOT 和 GOPATH 環境變數引用的位置去查找。

```shell
GOROOT=/Users/me/go
GOPATH=/Users/me/spaces/go/projects
```

log 包提供打印日誌信息到標準輸出 (stdout)、標準錯誤 (stderr) 或者自定義設備的功能。sync 包提供同步 goroutine 的功能。這個示例程序需要用到同步功能。

```go
// 註冊用於搜索的匹配器的映射
var matchers = make(map[string]Matcher)
```

這個變數沒有定義在任何函數作用域內，所以會被當成包級變數。這個變數使用 var 關鍵字聲明，而且聲明為 Matcher 類型的映射 (map)，這個映射以 string 類型值作為鍵，Matcher 類型值作為映射後的值。Matcher 類型在代碼文件 matcher.go 中聲明，後面再講這個類型的用途。這個變數聲明還有一個地方要強調一下：變數名 matchers 是以小寫字母開頭的。

在 Go 語言裡，標飾符要馬從包裡公開，要馬不從包裡公開。當代碼導入了一個包時，程序可以直接訪問這個包中任意一個公開的標飾符。這些標飾符以大寫字母開頭。以小寫字母開頭的標飾符是不公開的，不能被其他包中的代碼直接訪問。但是，其他包可以間接訪問不公開的標飾符。例如，一個函數可以返回一個未公開類型的值，那麼這個函數的任何調用者，哪怕調用者不是在這個包裡聲明的，都可以訪問這個值。

這個變數聲明還使用賦值運算符和特殊的內置函數 make 初始化了變數：

```go
make(map[string]Matcher)
```

map 是 Go 語言裡的一個引用類型，需要使用 make 來構造。如果不先構造 map 並將構造後的值復值給變數，會在試圖使用這個 map 變數時收到出錯信息。這是因為 map 變數默認的零值是 nil。第四章會進一步了解關於映射的細節。

在 Go 語言中，所有變數都被初始化為零值。

- 數值類型，零值是 0
- 字符串類型，零值是空字符串
- 布林類型，零值是false
- 指針類型，零值是 nil
- 引用類型，引用的底層數據結構會被初始化為對應的零值。但是被聲明為零值的引用類型的變數，會返回 nil 作為其值。

現在來看看 search.Run 函數的內容：

```go
// search/search.go

```