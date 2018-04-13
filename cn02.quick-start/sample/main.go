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