package search

import (
  "log"
  "sync"
)

var matchers = make(map[string]Matcher)

// Run 執行搜索邏輯
func Run(searchTerm string) {
  // 獲取需要搜索的數據源列表
  feeds, err := RetrieveFeeds()
  if err != nil {
    log.Fatal(err)
  }

  // 創建一個無緩衝的通道，接收匹配後的結果
  results := make(chan *Result)

  // 構造一個 waitGroup，以便處理所有數據源
  var waitGroup sync.WaitGroup

  // 設置需要等待處理
  // 每個數據源的 goroutine 的數量
  waitGroup.Add(len(feeds))

  // 為每個數據源啟動一個 goroutine 來查找結果
}