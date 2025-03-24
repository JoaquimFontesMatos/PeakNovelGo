import re

from bs4 import BeautifulSoup as bs

from .novel import NovelMetadata


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
    try:
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
    except Exception:
        return 0


def parse_series_novelbin(req):
    soup = bs(req, "html.parser")
    article = soup.find("div", id="novel")
    novel_info = article.find("div", class_="books")

    title = novel_info.find(class_="title").text.strip()

    image_tag = soup.find("meta", property="og:image")
    image=""
    if  image_tag:
        image = image_tag.get('content').strip()

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
        latest_chapter = int(chapter_number)

    description_tab = article.find("div", class_="tab-content")
    description = description_tab.find("div", class_="desc-text").text.strip()

    novel = NovelMetadata(
        title=title,
        image=image,
        genre=genres,
        tags=tags,
        rating=rating,
        language={"name": "English"},
        authors=authors,
        year=year,
        status=status,
        release_freq="N/A",
        description=description,
        latest_chapter=latest_chapter
    )

    return novel


def parse_series_lightnovelworld(req):
    title = ""
    image = ""
    authors = []
    genres = []
    tags = []
    status = ""
    year = "N/A"
    latest_chapter = 0

    soup = bs(req, "html.parser")

    article = soup.find("article", id="novel")

    main_head = article.find("div", class_="main-head")

    title_tag = soup.find("meta", property="og:title") or soup.find("meta", property="twitter:title")

    image_tag = soup.find("meta", property="og:image") or soup.find("meta", property="twitter:image")

    if title_tag and image_tag:
        title = title_tag.get('content').strip().replace(" | Light Novel World", "")
        image = image_tag.get('content').strip()

    rating = main_head.find("div", class_="rating").find("div", class_="rating-star").find("strong").text.strip()

    for author in article.find("div", class_="author").find_all("a"):
        authors.append(
            {
                "name": author.text.strip()
            }
        )

    for g in article.find("div", class_="categories").findAll("a"):
        genres.append(
            {
                "name": g.text.strip(),
            }
        )

    header_stats = article.find("div", class_="header-stats")
    for span in header_stats.find_all("span"):
        if span.find("small").text.strip() == "Chapters":
            latest_chapter = int(span.find("strong").text.strip())
        elif span.find("small").text.strip() == "Status":
            status = span.find("strong").text.strip()

    info = article.find("div", class_="novel-body container").find("section", id="info")

    description_paragraphs = info.find("div", class_="summary").find("div", class_="content").find_all("p",
                                                                                                       recursive=False)

    description = "\n\n".join(p.get_text(strip=True) for p in description_paragraphs) if description_paragraphs else ""

    for t in info.find("div", class_="tags").find("ul").findAll("a"):
        tags.append(
            {
                "name": t.text.strip(),
            }
        )

    novel = NovelMetadata(
        title=title,
        image=image,
        genre=genres,
        tags=tags,
        rating=rating,
        language={"name": "English"},
        authors=authors,
        year=year,
        status=status,
        release_freq="N/A",
        description=description,
        latest_chapter=latest_chapter
    )

    return novel


def parse_series_description(req):
    soup = bs(req, "html.parser")

    description_tag = soup.find("meta", property="og:description") or soup.find("meta", property="twitter:description")

    if description_tag and description_tag:
        description = description_tag.get('content')

        if description:
            return description

    return ""
