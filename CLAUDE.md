# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Webブラウザで遊べるオセロ（リバーシ）ゲーム。React + TypeScript フロントエンド、Go APIバックエンド、MySQL データベースを Docker Compose で構成。

## Technology Stack

- フロントエンド: React + TypeScript (Vite)
- バックエンド: Go
- データベース: MySQL 8.0
- インフラ: Docker Compose
- テスト: Vitest (フロントエンド), `go test` (バックエンド)
- Lint: `gofmt`, `go vet`

## Project Structure

- `frontend/` - React + TypeScript フロントエンド
  - `src/logic/` - ゲームロジック（盤面操作、合法手判定、勝敗判定）
  - `src/components/` - UIコンポーネント
  - `src/hooks/` - カスタムフック
  - `src/api/` - APIクライアント
- `backend/` - Go APIサーバー
  - `handler/` - HTTPハンドラー
  - `model/` - リクエスト/レスポンス型
  - `repository/` - データベースアクセス
  - `config/` - 設定
  - `middleware/` - CORSなど
- `mysql/init/` - DB初期化SQL
- `doc/` - 仕様書・ドキュメント

## Commands

```bash
# 全コンテナ起動
docker compose up --build

# フロントエンドテスト実行
cd frontend && npm test

# バックエンドテスト実行（Docker経由、ローカルにGoなし）
cd backend && docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine go test ./...

# バックエンドフォーマットチェック
cd backend && docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine gofmt -l .

# バックエンド静的解析
cd backend && docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine go vet ./...
```

## Development Methodology

Kent BeckのTDD（テスト駆動開発）に従って開発を進める。

### Task Management

タスクは `doc/TODO.md` で管理する。

- 実装開始前にタスク分解を行い、TODO.mdにチェックリストとして記載する
- 実装中に新しいタスクが発生した場合は、TODO.mdに追記する
- 完了したタスクはチェックを付ける

### TDDサイクル

1. **RED**: 失敗するテストを書く → コミット
2. **GREEN**: テストが通る最小限の実装を書く → コミット
3. **REFACTOR**: コードを整理する → コミット

### テストパスの条件

以下がすべて成功することをテストパスとする：

1. `go test ./...` が成功する
2. `gofmt -l .` で差分が出ない（フォーマットが正しい）
3. `go vet ./...` でエラーが出ない
4. `cd frontend && npm test` が成功する

### Commit Message Format

コミットメッセージの先頭にTDDのステータスを `[]` で囲んで記載する。
 <noreply@anthropic.com>などのメッセージを入れないようにしてください。

```
[RED] テスト内容の説明
[GREEN] 実装内容の説明
[REFACTOR] リファクタリング内容の説明
```

例:
- `[RED] 盤面初期化のテストを追加`
- `[GREEN] 盤面初期化を実装`
- `[REFACTOR] 盤面クラスのメソッドを整理`
