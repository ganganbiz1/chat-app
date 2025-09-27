# Chat App

このプロジェクトは、フルスタックのチャットアプリケーションです。バックエンドAPIはGoとEchoフレームワークで構築され、MySQL（永続化）とRedis（キャッシュ・セッション管理）のハイブリッド構成でレイヤードアーキテクチャを採用しています。

## 📁 プロジェクト構造

```
chat-app/
├── backend/                    # Go バックエンドAPI
│   ├── cmd/server/            # アプリケーションエントリーポイント
│   ├── internal/              # プライベートアプリケーションコード
│   │   ├── handlers/          # HTTPハンドラー
│   │   ├── services/          # ビジネスロジック
│   │   ├── repositories/      # データアクセス層
│   │   ├── models/           # ドメインモデル
│   │   └── config/           # 設定
│   ├── pkg/                  # 共有ライブラリ
│   ├── migrations/           # データベースマイグレーション
│   ├── Dockerfile
│   └── go.mod
├── frontend/                  # フロントエンド（今後実装予定）
├── docker-compose.yml         # 全体のDocker構成
└── README.md
```

## 🏗️ アーキテクチャ

- **Presentation Layer**: HTTPハンドラー、ミドルウェア
- **Business Logic Layer**: ビジネスロジック、サービス
- **Data Access Layer**: リポジトリパターン（MySQL + Redis）
- **Domain Models**: エンティティ、エラー定義

## 🚀 技術スタック

- **Language**: Go 1.22
- **Framework**: Echo v4
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **ORM**: GORM
- **Containerization**: Docker + Docker Compose

## 📋 API エンドポイント

### ヘルスチェック
- `GET /health` - サーバーステータス確認

### ルーム管理
- `POST /api/rooms` - チャットルーム作成
- `GET /api/rooms` - ルーム一覧取得（ページネーション対応）
- `GET /api/rooms/:id` - 特定ルーム情報取得
- `PUT /api/rooms/:id` - ルーム情報更新
- `DELETE /api/rooms/:id` - ルーム削除

### メッセージ管理
- `POST /api/rooms/:id/messages` - メッセージ送信
- `GET /api/rooms/:id/messages` - メッセージ履歴取得（ページネーション対応）
- `GET /api/rooms/:id/messages/recent` - 最新メッセージ取得（キャッシュ優先）
- `GET /api/messages/:id` - 特定メッセージ取得
- `DELETE /api/messages/:id` - メッセージ削除
- `GET /api/users/:userId/messages` - ユーザーメッセージ履歴取得

## 🚦 起動方法

### 前提条件
- Docker & Docker Compose

### 起動手順

#### 🚀 開発環境（ホットリロード有効）
```bash
# プロジェクトディレクトリに移動
cd chat-app

# 開発環境起動（Air によるホットリロード）
docker compose up --build -d

# 動作確認
curl http://localhost:8080/health
```

#### 💻 ローカル開発（推奨：Docker使用）
```bash
# 最も簡単な開発方法
docker compose up --build -d

# すべてがコンテナで動作し、環境構築不要
```

### サービス構成

- **API サーバー**: `http://localhost:8080`
- **MySQL**: `localhost:3306`
- **Redis**: `localhost:6379`

## 📝 リクエスト例

### ルーム作成
```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "一般チャット",
    "description": "みんなで話そう！"
  }'
```

### メッセージ送信
```bash
curl -X POST http://localhost:8080/api/rooms/{room_id}/messages \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "username": "太郎",
    "content": "こんにちは！"
  }'
```

### メッセージ取得
```bash
curl "http://localhost:8080/api/rooms/{room_id}/messages?page=1&limit=20"
```

## 🗄️ データベース設計

### Rooms テーブル
- `id` (CHAR(36), PRIMARY KEY)
- `name` (VARCHAR(255), NOT NULL)
- `description` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Messages テーブル
- `id` (CHAR(36), PRIMARY KEY)
- `room_id` (CHAR(36), FOREIGN KEY)
- `user_id` (VARCHAR(255), NOT NULL)
- `username` (VARCHAR(255), NOT NULL)
- `content` (TEXT, NOT NULL)
- `created_at` (TIMESTAMP)

## 🔧 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `DB_HOST` | localhost | MySQL ホスト |
| `DB_PORT` | 3306 | MySQL ポート |
| `DB_USER` | chatuser | MySQL ユーザー |
| `DB_PASSWORD` | chatpass | MySQL パスワード |
| `DB_NAME` | chatapp | MySQL データベース名 |
| `REDIS_URL` | localhost:6379 | Redis URL |
| `REDIS_PASSWORD` | (空) | Redis パスワード |
| `REDIS_DB` | 0 | Redis DB番号 |
| `PORT` | 8080 | API サーバーポート |

## 🛠️ 開発

### ローカル開発のオプション

#### Option 1: Docker + Air（推奨）
```bash
# ホットリロード付きで開発環境起動
docker compose up --build

# ファイル変更時に自動でリロードされます
```

#### Option 2: ローカル実行 + Air
```bash
# データベースのみDocker起動
docker compose up mysql redis

# Airをインストール（初回のみ）
go install github.com/cosmtrek/air@v1.49.0

# Airでホットリロード実行
cd backend
air
```

#### Option 3: 通常のローカル実行
```bash
# データベースのみDocker起動
docker compose up mysql redis

# 通常実行
cd backend
go run cmd/server/main.go
```

### テスト実行
```bash
go test ./...
```

## 🏭 本番環境デプロイ

本番環境では個別のDockerコンテナ運用を想定しています：

### Backend用Dockerコマンド
```bash
# 本番用イメージビルド
docker build --target prod -t chat-backend:latest ./backend

# コンテナ実行（環境変数で接続先指定）
docker run -d \
  --name chat-backend \
  -p 8080:8080 \
  -e DB_HOST=your-mysql-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e REDIS_URL=your-redis-host:6379 \
  chat-backend:latest
```

### データベース・キャッシュ
- **MySQL**: 外部管理（RDS等）
- **Redis**: 外部管理（ElastiCache等）

## 📚 今後の拡張予定

- [ ] WebSocket対応（リアルタイム通信）
- [ ] ユーザー認証・認可
- [ ] ファイルアップロード機能
- [ ] プッシュ通知
- [ ] 全文検索機能
- [ ] API レート制限
- [ ] Kubernetes対応
- [ ] CI/CD パイプライン