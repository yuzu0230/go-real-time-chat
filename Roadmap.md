---
title: Go Real-time Chat Project Roadmap

---

# Go Real-time Chat 專案學習路線圖

這是一個透過實作即時對話 App 來學習 Go 語言、單元測試、CI/CD 與分散式系統的 project 計畫。

## 核心學習目標

1.  **Go 語言 proficiency**: 從零基礎到能開發網路應用。
2.  **Unit Testing**: 學習撰寫可維護、可信賴的程式碼。
3.  **CI/CD**: 自動化測試與建置流程。
4.  **分散式系統設計**: 理解如何建構高可用、可擴展的系統。
5.  **後端架構**: 專注於後端邏輯，前端為輔。

---

## Phase 0: Go 語言基礎與環境設定

**目標**: 掌握 Go 的基本語法、專案結構，並寫出你的第一個 Go 程式。

### 學習重點 (對應知識點 #1)

* **環境安裝**: 安裝 Go，理解 `GOPATH` 與 Go Modules 的差異。
* **基本語法**: 變數 (`var`, `:=`)、資料型別 (structs, maps, slices)、流程控制 (`if`, `for`, `switch`)、函式。
* **Go 的獨特之處**:
    * `package` 概念
    * `public`/`private` (首字母大寫/小寫)
    * `error` interface 的錯誤處理機制
* **併發基礎**: 初步了解 `goroutine` (輕量級執行緒) 與 `channel` (安全的溝通管道)。

### 第一個小任務

1.  **完成官方教學**: [A Tour of Go](https://go.dev/tour/welcome/1)
2.  **建立簡單伺服器**: 撰寫一個 HTTP Server，當訪問 `http://localhost:8080/ping` 時，回傳 "pong"。

```golang
package main

import (
	"fmt"
	"net/http"
	"log"
)
func main() {
    // 當有請求訪問 /ping 路徑時，執行這個函式
    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        // 將 "pong" 字串寫入到 HTTP Response 中，這樣瀏覽器才會收到
        fmt.Fprintf(w, "pong")
    })
    // 啟動伺服器並監聽 8080 port
	// 如果啟動失敗，log.Fatal 會印出錯誤並結束程式
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Phase 1: 建立一個 WebSocket Echo Server (廣播聊天室)

**目標**: 理解 WebSocket 運作原理，並用 Go 實作一個廣播伺服器。這是專案的 "Hello, World!"。

### 學習重點 (對應知識點 #1, #2)

* **Go 網路程式設計**: 使用 `net/http` 處理 HTTP 請求並升級為 WebSocket 連線。
* **外部庫管理**: 使用 Go Modules 引入並管理 `gorilla/websocket` 套件。
* **Goroutine 應用**: 為每一個 WebSocket 連線建立一個獨立的 `goroutine` 進行處理，避免阻塞。
* **Channel 初體驗**: 建立一個中央 `channel`，將所有連線收到的訊息匯集於此，再由一個廣播 `goroutine` 統一發送給所有連線。
* **單元測試 (Unit Test)**: 使用 `testing` package 和 `go test` 指令，為獨立的純函式撰寫第一個測試案例。

### 待思考問題

> 當一個新的使用者連線進來，或是一個使用者斷線時，你要如何管理這些連線的清單？用什麼資料結構來儲存它們 (`slice`? `map`?) 會比較有效率？為什麼？

---

## Phase 2: 實現一對一與群組聊天

**目標**: 重構架構，加入「房間」和「使用者」的概念，支援更複雜的訊息路由。

### 學習重點 (對應知識點 #1, #4)

* **架構設計**: 引入 "Hub" (或稱 "Broker") 的設計模式。一個 Hub 代表一個聊天室，負責管理其內部的所有使用者連線與訊息轉發。
* **資料結構**: 大量使用 `map` 來管理房間與使用者，例如 `map[roomID]*Hub` 和 `map[userID]*Client`。
* **併發安全**: 當多個 `goroutine` 同時存取共享資源 (如 `map`) 時，學習使用 `sync.Mutex` 或透過 `channel` 來避免 race condition。
* **分散式系統思維**: 開始思考單機架構的瓶頸。如果記憶體不足或伺服器當機該如何處理？為後續的架構演進埋下伏筆。

### 待思考問題

> 訊息的格式應該是什麼？我們應該用 JSON 嗎？一個訊息的 JSON 應該包含哪些欄位（例如：`fromUser`, `toUser`, `roomID`, `content`, `timestamp`）才能同時滿足一對一和群組聊天的需求？

---

## Phase 3: 訊息持久化與使用者系統

**目標**: 整合資料庫，儲存使用者資訊、房間資訊和歷史訊息，讓資料不再因伺服器重啟而遺失。

### 學習重點 (對應知識點 #1, #4)

* **資料庫操作**:
    * 學習 Go 標準庫 `database/sql` 的用法。
    * 從輕量的 `SQLite` 開始，再逐步過渡到 `PostgreSQL` 或 `MySQL`。
* **系統分層**: 將程式碼解耦為不同層次：
    * `handlers`: 處理網路請求。
    * `services`: 處理商業邏輯。
    * `repositories`: 處理資料庫操作。
* **ORM (Object-Relational Mapping)**: 了解 `GORM` 或 `sqlx` 等工具如何簡化資料庫操作。

### 待思考問題

> 當使用者上線時，他需要看到離線期間的未讀訊息。這個邏輯應該在哪裡實現？是使用者連上 WebSocket 後主動跟後端要，還是後端主動推送？哪種方式對系統的擴展性比較好？

---

## Phase 4: CI/CD 與專案自動化

**目標**: 建立自動化的流程，當程式碼推送到 Git Repo 時，自動執行測試與建置。

### 學習重點 (對應知識點 #2, #3)

* **CI/CD 概念**: 理解持續整合 (CI) 與持續部署 (CD) 的核心思想。
* **工具學習**: 使用 **GitHub Actions** 建立你的第一個 workflow。
* **基本流程**: 在 `.github/workflows/main.yml` 中定義以下步驟：
    1.  `checkout`: 拉取最新程式碼。
    2.  `setup-go`: 設定 Go 環境。
    3.  `test`: 執行 `go test ./...` 跑所有單元測試。
    4.  `build`: 執行 `go build -v ./...` 確認程式可被編譯。

### 待思考問題

> 除了跑測試和建置，CI/CD 還能做什麼？（提示：程式碼品質掃描 `lint`、產生 Docker Image、自動部署到伺服器）。

---

## Phase 5: 邁向分散式架構

**目標**: 將單機的有狀態 (Stateful) 架構，演化為可以水平擴展 (horizontal scaling) 的分散式架構。

### 學習重點 (對應知識點 #4)

* **識別瓶頸**: 理解為什麼單機**有狀態**的服務難以擴展。當 User A 和 User B 連在不同伺服器上時，訊息無法直接傳遞。
* **解決方案：訊息佇列 (Message Queue)**:
    * 引入 `Redis Pub/Sub`, `NATS`, 或 `Kafka` 作為中間層。
    * **新流程**: Server A 收到訊息後，不直接找 User B，而是將訊息**發布**到 MQ。所有伺服器都**訂閱**它們關心的主題。持有 User B 連線的 Server B 從 MQ 收到訊息後，再發送給 User B。
* **無狀態服務 (Stateless)**: 透過 MQ，Go 應用伺服器本身變得「幾乎」無狀態，只負責維持連線與轉發，易於擴展和容錯。
* **連線狀態管理**: 使用 `Redis` 等高速快取來記錄哪個使用者目前在哪一台伺服器上。

### 最終思考問題

> 在這種分散式架構下，如何保證訊息「只被傳送一次」 (exactly-once delivery)？如果 Server B 從訊息佇列拿到訊息後，還沒發給 User B 就掛了，這個訊息會不會遺失？