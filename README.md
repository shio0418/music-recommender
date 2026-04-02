# music-recommender

YouTube URL から楽曲情報を取得し、簡易スコア付きで一覧表示する MVP です。

## 現在できること

- 楽曲の追加（手入力: `title`, `artist`, `tag`）
- YouTube URL 追加（`youtubeUrl`）
  - `title`（動画タイトル）
  - `artist`（チャンネル名を代用）
  - `thumbnailUrl`
  - `embedUrl`
  - `viewCount` / `likeCount` / `commentCount`
  - `score`（再生数に対する反応率を重視）
- 楽曲一覧表示
- おすすめ表示（現状はタグ `夜` のみを返す簡易実装）

## 構成

- `backend/` Go + Echo API
- `frontend/` React + Vite

## セットアップ

### 1. APIキー設定

`backend/.env` を作成して設定:

```env
YOUTUBE_API_KEY=your_api_key
```

> `backend/.env` は `.gitignore` 済みです。

### 2. バックエンド起動

```bash
cd backend
go run main.go
```

### 3. フロントエンド起動

```bash
cd frontend
npm install
npm run dev
```

## API（現状）

- `GET /songs` 楽曲一覧
- `POST /songs` 楽曲追加
  - 手入力例:
    ```json
    { "title": "夜に駆ける", "artist": "YOASOBI", "tag": "夜" }
    ```
  - YouTube 例:
    ```json
    { "youtubeUrl": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "tag": "夜" }
    ```
- `GET /recommend` おすすめ一覧（現状は `tag == "夜"`）

## スコア式（現状）

```text
engagementRate = (likes + comments*3) / (views + 300)
score = engagementRate*1000 + 0.15*log(1+views)
```

「再生数が少ないわりに、いいね・コメントが多い」動画を上に出しやすい設計です。

## 今後やりたいこと

- おすすめロジックの改善
  - 固定タグではなく、入力曲や選択タグベースで推薦
  - 複数要素（タグ、score近傍、チャンネル傾向）でランキング
- サンプル収集機能の追加
  - `POST /seed` のような一括収集 API（検索クエリから初期データを作る）
- データ永続化
  - 現在メモリ保存のため、再起動で消える
  - SQLite / Postgres などへ移行
- UI改善
  - ローディング/エラー表示
  - カードUI、並び替え、フィルタ
- 品質改善
  - テスト追加
  - ルーティング・ロジック分割（`main.go` の責務分離）
