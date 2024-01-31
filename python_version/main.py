import requests
from bs4 import BeautifulSoup


url = "https://www.espn.com/nfl/boxscore/_/gameId/401547639"
response = requests.get(url)
html_content = response.text

soup = BeautifulSoup(html_content, "html.parser")

data = soup.find("div", class_="Boxscore__Category")
print(soup)
