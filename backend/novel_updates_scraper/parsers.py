from bs4 import BeautifulSoup as bs
import re


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


def parseChapters(req):
    soup = bs(req.text, "html.parser")

    chapter_article = soup.find("article", id="chapter-article")
    section = chapter_article.find("section", class_="page-in content-wrap")
    title = section.find("h1").find("span", class_="chapter-title").text.strip()

    body = section.find("div", id="chapter-container")

    body = body.decode_contents().strip()

    result = {"title": title, "body": body, "err": False}

    return result


def parseSeries(req):
    soup = bs(req.text, "html.parser")

    header = soup.find("header", class_="novel-header")
    novel_body = soup.find("div", class_="novel-body container")

    image = header.find("figure", class_="cover").find("img").get("data-src")

    novel_info = header.find("div", class_="novel-info")

    title = novel_info.find("h1").text.strip()

    author_div = novel_info.find("div", class_="author")
    authors = []
    for author in author_div.find_all("a"):
        link = author.get("href")
        authorName = author.find("span").text
        authors.append({"name": authorName, "link": link})

    rating_div = novel_info.find("div", class_="rating")
    rating = rating_div.find("div", class_="rating-star").find("strong").text.strip()

    categories_div = novel_info.find("div", class_="categories")
    genre = []
    for g in categories_div.find("ul").find_all("li"):
        category = g.find("a")
        genre.append(
            {
                "name": category.text,
                "link": category.get("href"),
                "description": category.get("title"),
            }
        )

    stats = novel_info.find("div", class_="header-stats")
    latest_chapter = stats.find("span").find("strong").text.strip()
    second_span = stats.find("span").find_next_sibling("span")
    views = second_span.find("strong").text.strip()

    third_span = stats.find("span").find_next_sibling("span").find_next_sibling("span")
    bookmarks = third_span.find("strong").text.strip()

    fourth_span = (
        stats.find("span")
        .find_next_sibling("span")
        .find_next_sibling("span")
        .find_next_sibling("span")
    )
    status = fourth_span.find("strong").text.strip()

    info = novel_body.find("section", id="info")
    summary_div = info.find("div", class_="summary")
    description = (
        summary_div.find("div", class_="content expand-wrapper").find("p").text.strip()
    )

    tags_div = info.find("div", class_="tags")
    tag = []
    for t in tags_div.find("div", class_="expand-wrapper").find("ul").find_all("li"):
        tag_item = t.find("a")
        tag.append(
            {
                "name": tag_item.text,
                "link": tag_item.get("href"),
                "description": tag_item.get("title"),
            }
        )

    result = {
        "title": title,
        "image": image,
        "genre": genre,
        "tags": tag,
        "rating": rating,
        "language": {"name": "English"},
        "authors": authors,
        "year": "N/A",
        "status": status,
        "release_freq": "N/A",
        "description": description,
    }
    return result
