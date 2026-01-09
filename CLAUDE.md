# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLIで遊べるオセロ（リバーシ）ゲーム。

## Technology Stack

- 言語: Go
- テスト: `go test`
- Lint: `gofmt`, `go vet`

## Project Structure

- `src/` - オセロアプリのソースコード（Goモジュール）
- `tests/` - テストコード（`*_test.go`）
- `doc/` - 仕様書・ドキュメント

## Commands

```bash
# テスト実行
go test ./...

# フォーマットチェック
gofmt -l .

# フォーマット適用
gofmt -w .

# 静的解析
go vet ./...
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

### Commit Message Format

コミットメッセージの先頭にTDDのステータスを `[]` で囲んで記載する。

```
[RED] テスト内容の説明
[GREEN] 実装内容の説明
[REFACTOR] リファクタリング内容の説明
```

例:
- `[RED] 盤面初期化のテストを追加`
- `[GREEN] 盤面初期化を実装`
- `[REFACTOR] 盤面クラスのメソッドを整理`
