import pdfplumber
import requests
from io import BytesIO
import re
from bs4 import BeautifulSoup
from datetime import date

def is_white_background(cell_bbox, rects):
  for rect in rects:
    if 'non_stroking_color' in rect:
      color = rect['non_stroking_color']
      if color != (1,):  # RGB白でない
        # バウンディングボックスの重なり判定
        rx0, ry0, rx1, ry1 = rect['x0'], rect['top'], rect['x1'], rect['bottom']
        cx0, cy0, cx1, cy1 = cell_bbox
        if not (cx1 < rx0 or cx0 > rx1 or cy1 < ry0 or cy0 > ry1):
          return False
  return True

def reiwa_to_date(date_str):
  match = re.match(r"令和(\d+)年(\d+)月(\d+)日",date_str)
  if not match:
    raise ValueError("reiwa format error")
  
  reiwa_year, month, day = map(int, match.groups())
  year = 2018 + reiwa_year
  return date(year,month,day)

# 奨学金のレコード(辞書型)のリストを返す
def fetch_latest_scholarships():
  base_url = "https://www.kit.ac.jp/campus_index/life_fee/scholarship/minkanscholarship/"
  pdf_pattern = re.compile(r"https://www\.kit\.ac\.jp/wp/wp-content/uploads/\d{4}/\d{2}/.*hpsyougakukinitiran.*\.pdf")

  html = requests.get(base_url).text
  soup = BeautifulSoup(html, "html.parser")
  
  pdf_links = [a["href"] for a in soup.find_all("a",href=True) if pdf_pattern.match(a["href"])]
  if not pdf_links:
    return []

  pdf_url = pdf_links[0]
  res = requests.get(pdf_url)
  result = []

  # 無効な奨学金情報を見つけたらフラグを立てる
  flag = False 
  with pdfplumber.open(BytesIO(res.content)) as pdf:
    for page in pdf.pages:
      # print(f"page {page}")
      table = page.extract_table()
      if not table:
        continue

      rects = page.rects
      headers = table[0]
      for row in table[1:]:
        row_data = dict(zip(headers,row))
        first_cell_text = row[0]

        # セルの位置取得
        match_word = next((w for w in page.extract_words() if w['text'].strip() == first_cell_text.strip()), None)
        if not match_word:
          continue

        bbox = (match_word['x0'], match_word['top'], match_word['x1'], match_word['bottom'])
        # 背景が灰色の項目が現れると以降の項目は不要
        if not is_white_background(bbox, rects):
          flag = True
          break

        # 整形
        if row_data["掲示日"] != "":
          row_data["掲示日"] = reiwa_to_date(row_data["掲示日"])
        row_data["対象(学部・院)"] = row_data["対象(学部・院)"].replace("\n","")
        if row_data["申請期限等"] != "" and row_data["申請期限等"].startswith("令和"):
          row_data["申請期日"] = reiwa_to_date(row_data["申請期限等"])
        else:
          row_data["申請期日"] = None

        # print(row_data)
        result.append(row_data)
        
      if flag :
        break

  return result

fetch_latest_scholarships()