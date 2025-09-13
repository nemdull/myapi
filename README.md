# MyAPI: Go-based Article & Comment REST API

Go + Gorilla Mux を採用した、記事の投稿・取得・いいね・コメント機能を提供するシンプルな REST API です。  
データベースは MySQL、サーバーは port 8080 で起動します。

---

## プロジェクト概要
Go + Gorilla Mux による記事投稿・一覧・詳細表示・いいね機能・コメント機能を実装した小規模 API。サンプルクライアントからのリクエストを受け、サービス層がビジネスロジックを担い、リポジトリ層が DB との I/O を抽象化します。

---

## 技術スタック

- 言語/フレームワーク
  - Go (モジュール: github.com/nemdull/myapi)
  - Gorilla Mux
- データベース/ORM系
  - MySQL (github.com/go-sql-driver/mysql)
- その他
  - Docker Compose 対応（DB/サービスの切り出し想定）
  - テスト: Go の標準 testing パッケージ
  - CI/CD: GitHub Actions workflows（.github/workflows 配置想定）
- ファイル/構成の出典
  - go.mod: module github.com/nemdull/myapi
  - main.go, handlers/handlers.go, models/models.go, services/article_service.go などの実装群

---

## 主な機能一覧

- GET /hello
  - Hello のテスト用エンドポイント
- POST /article
  - 記事を新規作成
- GET /article/list?page=N
  - 指定ページの記事一覧を取得（ページネーション、デフォルトは 1 ページ目）
- GET /article/{id}
  - 指定記事の詳細表示（コメント一覧を含む）
- POST /article/nice
  - 指定記事の「いいね」数を増加
- POST /comment
  - 記事へコメントを投稿

---

## 設計・実装の工夫

- 層構造
  - ハンドラ層（Handlers）で HTTP 入力を受け取り、サービス層へ委譲
  - サービス層（Services）でビジネスロジックを実装
  - リポジトリ層（Repositories）で DB との I/O を抽象化
- データモデル
  - models.Article と models.Comment のデータ構造を使用
  - Article は CommentList を持ち、記事に紐づくコメントを含む
- エラーハンドリングの方針
  - JSON デコード失敗時は 400 Bad Request
  - 内部処理エラー時は 500 Internal Server Error
- API の挙動
  - ArticleDetail で記事情報とコメントを統合して返却
  - Nice 更新後は NiceNum をインクリメントした新しい記事情報を返却
- クエリパラメータ処理
  - /article/list の page パラメータを受け取り、未指定時は 1 をデフォルト

---

## セットアップ & 動作確認方法

- 要件
  - Go 環境（推奨: Go 1.24.x 系以上）
  - Docker/Docker Compose（DB 環境を Docker で起動する場合）
  - Git（ブランチ戦略として feature/readme を使用可能）

- ローカル実行
  1. 依存関係の取得
     - go mod download
  2. アプリの起動
     - go run main.go
     - ブラウザ/curl で http://localhost:8080 へアクセス
  3. テストの実行
     - go test ./...

- Docker Compose を利用する場合
  - DB 環境とアプリを一括で起動
    - docker-compose up -d
  - DB の準備/初期化は testdata/setupDB.sql などのスクリプトを参照
  - アプリはデフォルトで port 8080 をリスニング

- 主なエンドポイントの curl 例
  - Hello
    - curl http://localhost:8080/hello
  - 記事作成
    - curl -X POST -H "Content-Type: application/json" -d '{"article_id":1,"title":"サンプル","contents":"本文","user_name":"tester"}' http://localhost:8080/article
  - 記事一覧
    - curl http://localhost:8080/article/list?page=1
  - 記事詳細
    - curl http://localhost:8080/article/1
  - いいね
    - curl -X POST -H "Content-Type: application/json" -d '{"article_id":1,"title":"サンプル","contents":"本文","user_name":"tester","nice":0,"created_at":"2025-01-01T00:00:00Z"}' http://localhost:8080/article/nice
  - コメント投稿
    - curl -X POST -H "Content-Type: application/json" -d '{"article_id":1,"message":"コメントです"}' http://localhost:8080/comment

---

## 改善ポイント / TODO

- テスト
  - 現状のエンドポイントに対するユニットテスト/統合テストを追加
  - 例: PostArticleHandler のデコードエラーパスのテスト、GetArticleList のページ境界テスト、ArticleDetail の不正 ID のテスト
- エラーハンドリング
  - PostArticleHandler で JSON デコード失敗時に早期リターンを追加して、本来の処理を止めるべき
  - 具体的には decode エラー時に return を追加
- 設計の安定性
  - article.CommentList のマージ処理の意図を明示化（空配列初期化・append の正確性）
  - DB 接続の再利用/接続プールの活用検討
- CI/CD
  - CI 側の lint 判定を有効化（現在の設定状態を再確認・有効化を検討）
  - DB 初期化スクリプトの自動実行や tests の自動化の強化
- ドキュメントの補完
  - コード構造の概要、API 仕様、データモデルの説明を README に詳述
  - testdata ディレクトリの利用方法やサンプルデータの説明を追記

---

## 本リポジトリで強調したいポイント

- 層化アーキテクチャ（ハンドラ -> サービス -> リポジトリ）
- Go での REST API 実装の実践例
- JSON 入出力のスキーマ設計と、データモデルの整合性
- ページネーション・データ結合の基本パターン
- エラーハンドリングの標準的な対応

---

## 変更履歴の要約

- README.md の作成: リポジトリの概要・技術スタック・主な機能・設計方針・セットアップ方法・TODO のカテゴリ分けを含む
- 今後の改善案を TODO に分類してドキュメント化
