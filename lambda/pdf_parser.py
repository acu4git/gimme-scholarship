import pdfplumber
import requests
from io import BytesIO
import re
from bs4 import BeautifulSoup

def fetch_latest_scholarships():
  base_url = "https://www.kit.ac.jp/campus_index/life_fee/scholarship/minkanscholarship/"
  pdf_pattern = re.compile(r"https://www\.kit\.ac\.jp/wp/wp-content/uploads/\d{4}/\d{2}/.*hpsyougakukinitiran.*\.pdf")

  html = requests.get(base_url).text
  soup = BeautifulSoup(html, "html.parser")
  
  pdf_links = [a["href"] for a in soup.find_all("a",href=True) if pdf_pattern.match(a["href"])]
  if pdf_links:
    pdf_url = pdf_links[0]
  else :
    return

  res = requests.get(pdf_url)
  with pdfplumber.open(BytesIO(res.content)) as pdf:
    # for page in pdf.pages:
    #   table = page.extract_table()
    #   headers = table[0]
    #   for row in table[1:]:
    #     row_data = dict(zip(headers,row))
    #     print(row_data)
    page = pdf.pages[0]
    table = page.extract_table()

    headers = table[0]
    for row in table[1:]:
      row_data = dict(zip(headers, row))
      print(row_data)

fetch_latest_scholarships()