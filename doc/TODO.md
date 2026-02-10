# TODO

## Phase 0: インフラ構築
- [x] docker-compose.yml 作成
- [x] MySQL初期化SQL作成
- [x] Goバックエンド初期化（go.mod, Dockerfile, config, CORS, main.go）
- [x] フロントエンド初期化（Vite + React + TypeScript, Dockerfile, nginx.conf）
- [x] docker compose up で3コンテナ起動確認

## Phase 1: Goバックエンド - モデル & リポジトリ（TDD）
- [x] model/model.go - 型定義
- [x] repository/repository.go - インターフェース + MySQL実装
- [x] testutil/testutil.go - テストDBヘルパー
- [x] repository/repository_test.go - インテグレーションテスト

## Phase 2: Goバックエンド - HTTPハンドラー（TDD）
- [x] handler/handler.go - 3エンドポイント実装
- [x] handler/handler_test.go - モックRepositoryテスト（100% coverage）
- [x] main.go - ルーティング・サーバー起動

## Phase 3: フロントエンド - ゲームロジック（TDD）
- [x] types/game.ts - 型定義
- [x] logic/board.ts + board.test.ts - 盤面ロジック（17 tests pass）
- [x] logic/game.ts + game.test.ts - ゲーム状態管理（10 tests pass）

## Phase 4: フロントエンド - UIコンポーネント
- [x] api/client.ts - APIクライアント
- [x] components/Cell.tsx + Cell.css
- [x] components/Board.tsx + Board.css
- [x] components/GameInfo.tsx
- [x] components/GameControls.tsx
- [x] hooks/useGame.ts
- [x] App.tsx + App.css

## Phase 5: 統合テスト & 仕上げ
- [x] docker compose up でフルフロー検証
- [x] Goテストカバレッジ80%以上確認（handler: 100%, config: 100%, middleware: 100%）
- [x] gofmt, go vet 確認
- [x] CLAUDE.md 更新
