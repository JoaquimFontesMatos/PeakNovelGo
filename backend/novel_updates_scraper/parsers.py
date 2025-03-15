import re

from bs4 import BeautifulSoup as bs


def parseFeed(req):
    soup = bs(req.text, "html.parser")
    feed = []
    for table in soup.find_all("table", id="myTable", class_="tablesorter"):
        for entry in table.find("tbody").find_all("tr"):
            title_elem = entry.select("td")[0].find("a")
            title = title_elem.get("title") if title_elem else None
            nuLink = title_elem.get("href") if title_elem else None

            release_elem = entry.select("td")[1].find(
                "span"
            )  # Using <span> tag for release
            release = release_elem.get("title") if release_elem else None
            releaseLink = (
                None  # Release link is not available in this case, so set it to None
            )

            group_elem = entry.select("td")[2].find("a")
            groupName = group_elem.get("title") if group_elem else None
            groupLink = group_elem.get("href") if group_elem else None

            feed.append(
                {
                    "title": title,
                    "nuLink": nuLink,
                    "group": {"name": groupName, "link": groupLink},
                    "release": {"name": release, "link": releaseLink},
                }
            )

    return feed


def parseSearch(req):
    soup = bs(req.text, "html.parser")
    results = []
    for result in soup.find_all("div", class_="search_main_box_nu"):
        body = result.find("div", class_="search_body_nu")
        imageBody = result.find("div", class_="search_img_nu")

        title = body.find("div", class_="search_title").find("a").text
        link = body.find("div", class_="search_title").find("a").get("href")

        image = imageBody.find("img").get("src")
        if image.endswith("noimagemid.jpg"):
            image = None
        imageBody.find("div", class_="search_ratings").find("span").decompose()
        searchRating = re.sub(
            r"[()]", "", imageBody.find("div", class_="search_ratings").text.strip()
        )

        ogDescription = body.find(text=True, recursive=False).strip()
        moreDescription = body.find("span", class_="testhide")
        for p in moreDescription.find_all("p", style="margin-top:-5px;"):
            p.decompose()
        moreDescription.find("span", class_="morelink list").decompose()
        description = ogDescription + moreDescription.text

        stats = body.find("div", class_="search_stats").find_all(
            "span", class_="ss_desk"
        )
        releases = stats[0].text.strip()
        updateFreq = stats[1].text.strip()
        nuReaders = stats[2].text.strip()
        nuReviews = stats[3].text.strip()
        lastUpdated = stats[4].text.strip()

        genres = []
        for genre in body.find("div", class_="search_genre").find_all("a"):
            genreName = genre.text
            genreLink = genre.get("href")
            genres.append({"name": genreName, "link": genreLink})

        results.append(
            {
                "title": title,
                "link": link,
                "image": image,
                "search_rating": searchRating,
                "description": description[:-1],
                "releases": releases,
                "update_freq": updateFreq,
                "nu_readers": nuReaders,
                "nu_reviews": nuReviews,
                "last_updated": lastUpdated,
                "genres": genres,
            }
        )
    return results


def parse_chapters(req):
    """Parses the chapter content from the response."""
    soup = bs(req.text, "html.parser")

    chapter_article = soup.find("div", class_="chapter container")
    if not chapter_article:
        return None  # No valid chapter found

    title = chapter_article.find("h2").find("a").text.strip()

    body_paragraphs = chapter_article.find("div", id="chr-content").find_all("p", recursive=False)

    # Extract text from all <p> elements and join them with line breaks
    body = "\n\n".join(p.get_text(strip=True) for p in body_paragraphs) if body_paragraphs else ""

    return {"title": title, "body": body}


def parse_series(req):
    soup = bs(req.text, "html.parser")
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
