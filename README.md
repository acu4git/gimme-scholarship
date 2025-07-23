# このリポジトリについて

京都工芸繊維大学の学生を対象とした奨学金一覧へのアクセスを便利にした Web アプリ「[KIT クレクレ奨学金](https://www.kit-gimme-scholarship.com/)」のバックエンドのリポジトリです．旧版の CLI アプリ([pdfxtractor](https://github.com/acu4git/pdfxtractor))を Web 版として再実装しました．

**フロントエンド**: https://github.com/acu4git/gimme-scholarship-front/<br>
**インフラ**: https://github.com/acu4git/gimme-scholarship-terraform/

# 開発環境

- Ubuntu 22.04.5 (WSL2)
- Go 1.24
- Python 3.12.0
- MySQL 8.0
- Docker version 28.2.2
- tbls 1.84.0
- sql-migrate v1.7.1

# クラウド

AWS を採用しました．リソースが豊富であり，今回実装するアプリケーションに必要なリソースを組み合わせることが容易であったためです．

## アーキテクチャ図

![gs-architecture-ver2](クレクレ奨学金_aws02.png)

# ER 図

![Entity Relationship Diagram](doc/schema/schema.svg)

# 機能

## api

このリポジトリのメインとなる部分です．ユーザー登録や奨学金検索，その他 API の責務を果たします．Web API として用意することでネイティブアプリを実装する際にも使えるようにしています．

### 主な機能

- ユーザー登録/削除
- 奨学金一覧取得
- お気に入り登録

## migrate

文字通り DB へのマイグレーションを行います．migration ディレクトリ配下にあるファイルを読み込ませることで，テーブルやカラムの作成・修正します．コンテナ化して ECS で稼働させることで任意のタイミングでの実行を可能にしています．

## fetch

奨学金一覧情報を取得して DB の更新を行います．毎日 AM10:00(JST)にバッチタスクとして起動します．

> [!Important]
> 奨学金一覧の更新フローは
>
> 1. 奨学金 PDF をスクレイピングにより取得
> 2. pdfplumber によりテーブル情報を抽出し，無効な情報を除外しながら奨学金情報をリストに追加
> 3. リストの内容と DB の内容で差分更新を行う
>
> の順番で行われています．

### AWS Lambda で良くない？

-> 駄目です．RDS と外部の Web サイトに接続しようとすると Internet Gateway を置く必要があり，これが中々お金がかかるので今回のようにしました．

## task

更新情報のメール通知などのバッチタスクです．

# メール通知
Eメール通知はAWS SES v2を利用しており，Goの標準ライブラリのembedとtext/templateを用いてメール文の埋め込み/生成を行っています．

# その他

- マイグレーションには sql-migrate を採用．他の golang-migrate 等に比べてバージョンが安定しており，ライブラリや CLI での利用のしやすさが魅力的であったため．
- DB テーブルの関係図を tbls によって描画し，ドキュメントとして整理．
- Go での実装はレイヤードアーキテクチャを採用．なんちゃって DDD になっている．
