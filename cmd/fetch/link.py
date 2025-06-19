import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin, urlparse
import re
import json


def scrape_scholarship_info(url: str) -> dict:
    """
    指定されたURLから奨学会名と関連リンクの情報をスクレイピングし、辞書として返す。

    Args:
        url (str): スクレイピング対象のURL。

    Returns:
        dict: 奨学会名をキー、リンク情報のリストをバリューとする辞書。
              リンク情報は {'text': 'リンクテキスト', 'url': 'リンクURL'} の辞書形式。
    """

    # URLからベースURLを抽出（相対パスを絶対パスに変換するため）
    parsed_url = urlparse(url)
    base_url = f"{parsed_url.scheme}://{parsed_url.netloc}"

    # 最終的なデータを格納する辞書
    scholarship_data = {}

    try:
        # Webページを取得
        response = requests.get(url)
        response.raise_for_status()
        response.encoding = response.apparent_encoding

        # BeautifulSoupオブジェクトを作成
        soup = BeautifulSoup(response.text, "html.parser")

        # ページ内のすべての<table>を検索
        tables = soup.find_all("table")

        if not tables:
            print("ページ内にテーブルが見つかりませんでした。")
            return scholarship_data

        for table in tables:
            # table内のすべての行(tr)をループ処理
            # a.find_parent('tr') は、リンクがヘッダー行にある場合を除外するのに役立つ
            for row in table.find_all("tr"):
                # 行からセル(td)をすべて取得
                cells = row.find_all("td")

                # セルが2つ以上ある行のみを処理対象とする（ヘッダー行や空の行を除外）
                if len(cells) >= 2:
                    # 1番目のセルから奨学会名を取得
                    scholarship_name = cells[0].get_text(strip=True)

                    if not scholarship_name:
                        continue

                    # --- ここからが修正箇所 ---
                    # 正規表現を使い、奨学会名末尾の更新日情報（例: " 6/16up"）を除去する
                    # パターン：(半角スペースが0個以上) + (数字1-2桁) + / + (数字1-2桁) + "up" + (文字列の末尾)
                    cleaned_name = re.sub(
                        r"\s*\d{1,2}/\d{1,2}up$", "", scholarship_name
                    ).strip()

                    # 2番目のセルからリンク情報を取得
                    links_cell = cells[1]
                    links = []

                    for a_tag in links_cell.find_all("a"):
                        link_text = a_tag.get_text(strip=True)
                        raw_url = a_tag.get("href")

                        # 相対URLを絶対URLに変換
                        absolute_url = urljoin(base_url, raw_url)

                        links.append({"text": link_text, "url": absolute_url})

                    # 辞書に格納
                    if links:  # リンクが一つでも見つかった場合のみ追加
                        scholarship_data[cleaned_name] = links

    except requests.exceptions.RequestException as e:
        print(f"URLへのアクセス中にエラーが発生しました: {e}")
    except Exception as e:
        print(f"処理中にエラーが発生しました: {e}")

    return scholarship_data


# --- 実行 ---
if __name__ == "__main__":
    target_url = (
        "https://www.kit.ac.jp/campus_index/life_fee/scholarship/minkanscholarship/"
    )
    result_dict = scrape_scholarship_info(target_url)

    # 結果をJSON形式で分かりやすく表示
    print(json.dumps(result_dict, indent=2, ensure_ascii=False))
