import re

from bs4 import BeautifulSoup as bs


def parse_chapters_wuxiabox(req):
    """Parses the chapter content from the response."""
    soup = bs(req.text, "html.parser")

    chapter_article = soup.find("article", id="chapter-article")
    if not chapter_article:
        return None  # No valid chapter found

    title = chapter_article.find("h2").text.strip()

    body_paragraphs = chapter_article.find("div", class_="chapter-content").find_all("p", recursive=False)

    # Extract text from all <p> elements and join them with line breaks
    body = "\n\n".join(p.get_text(strip=True) for p in body_paragraphs) if body_paragraphs else ""

    return {"title": title, "body": body}


def parse_chapters_novtales(req):
    """Parses the chapter content from the response."""
    soup = bs(req.text, "html.parser")
    title_tag = soup.find("meta", property="og:title") or soup.find("meta", name="twitter:title")
    description_tag = soup.find("meta", property="og:description") or soup.find("meta", name="twitter:description")

    if title_tag and description_tag:
        title = title_tag.get('content')
        body = description_tag.get('content')

        if title and body:
            return {"title": title, "body": body}

    return None


def parse_latest_chap(req):
    soup = bs(req, "html.parser")
    latest_chapter = 0
    latest_chapter_div = soup.find("div", class_="bubble-element Group baTaHaNh bubble-r-container flex row")
    if latest_chapter_div.find("h2").text == "Latest Chapter:":
        latest_chapter_text = latest_chapter_div.find("div").text.strip()

        # Extract only the number using regex
        match = re.search(r"Chapter (\d+)", latest_chapter_text)

        if match:
            chapter_number = match.group(1)  # Get the number
            latest_chapter = chapter_number

    return latest_chapter


def parse_series(req, latest_chapter_param):
    soup = bs(req, "html.parser")
    article = soup.find("div", id="novel")
    novel_info = article.find("div", class_="books")

    title = novel_info.find(class_="title").text.strip()

    image = novel_info.find("img").get("data-src")

    info_meta = article.find(class_="info info-meta")

    authors = []
    genres = []
    tags = []
    status = ""
    year = "N/A"
    for li in info_meta.find_all("li"):
        if li.find("h3").text == "Author:":
            for author in li.find_all("a"):
                authors.append(
                    {
                        "name": author.text.strip()
                    }
                )

        elif li.find("h3").text == "Genre:":
            for g in li.findAll("a"):
                genres.append(
                    {
                        "name": g.text.strip(),
                    }
                )

        elif li.find("h3").text == "Tag:":
            for t in li.find("div", class_="tag-container").findAll("a"):
                tags.append(
                    {
                        "name": t.text.strip(),
                    }
                )
        elif li.find("h3").text == "Status:":
            status = li.find("a").text.strip()

        elif li.find("h3").text == "Year of publishing:":
            year = li.find("a").text.strip()

    rating_div = article.find("div", class_="rate-info")
    rating = rating_div.find("input").get("value").strip()

    latest_chapter_div = article.find("div", class_="l-chapter")
    latest_chapter_text = latest_chapter_div.find("a").text.strip()

    # Extract only the number using regex
    match = re.search(r"Chapter (\d+)", latest_chapter_text)

    latest_chapter = 0
    if match:
        chapter_number = match.group(1)  # Get the number
        latest_chapter = chapter_number

    if int(latest_chapter) < int(latest_chapter_param):
        latest_chapter = latest_chapter_param

    description_tab = article.find("div", class_="tab-content")
    description = description_tab.find("div", class_="desc-text").text.strip()

    result = {
        "title": title,
        "image": image,
        "genre": genres,
        "tags": tags,
        "rating": rating,
        "language": {"name": "English"},
        "authors": authors,
        "year": year,
        "status": status,
        "release_freq": "N/A",
        "description": description,
        "latest_chapter": latest_chapter,
    }
    return result
