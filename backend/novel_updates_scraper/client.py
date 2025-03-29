import json
import sys

import requests

from . import parsers
from .novel import NovelMetadata
from .request import PlaywrightScraper
from .sources import NovelSource, ChapterSource


class Client:
    def __init__(self):
        self.req = PlaywrightScraper()

    def series_info(self, novel_id):
        """Gets the metadata of a series.

       Parameters
       ----------
       novel_id : :class:`string`
            The id of the series. (/series/{novel_id})

       Returns
       -------
       :class:`dict`
           A dictionary containing novel information, or an error indicator.
       """
        urls = {
            NovelSource.NOVELBIN: f"https://novelbin.com/b/{novel_id}",
            NovelSource.LIGHTNOVELWORLD: f"https://www.lightnovelworld.co/novel/{novel_id}",
            NovelSource.NOVTALES: f"https://novtales.com/novel/{novel_id}"
        }

        try:
            error = {}

            default_content = self.req.scrape(urls[NovelSource.NOVTALES])

            # Parse the series information
            latest_chapter = parsers.parse_latest_chap(default_content)
            description = parsers.parse_series_description(default_content)

            for urlKey in urls.keys():
                url = urls[urlKey]

                content = self.req.scrape(url)

                if content is None:
                    return {"status": 503, "error": "Failed to fetch the page"}

                if "Page not found" in content:
                    return {"status": 404, "error": "Series not found"}

                chapter_data: NovelMetadata = {}

                if urlKey == NovelSource.NOVELBIN:
                    chapter_data = parsers.parse_series_novelbin(content)
                elif urlKey == NovelSource.LIGHTNOVELWORLD:
                    chapter_data = parsers.parse_series_lightnovelworld(content)

                if not description == "Read the best web novels on Novtales, featuring Xianxia, Korean, cultivation, and fantasy novels. Enjoy high-quality translations and explore legendary stories today!":
                    chapter_data["description"] = description

                if chapter_data["latest_chapter"] < int(latest_chapter):
                    chapter_data["latest_chapter"] = latest_chapter

                return chapter_data
        except Exception as e:
            return {"status": 503, "error": f"Unexpected error: {e}"}

    def chapters(self, novel_id="", chapter_no="1"):
        """Gets the chapters of a series.

        Parameters
        ----------
        novel_id : :class:`string`
            The id of the series. (/series/{novel_id})
        chapter_no : :class:`string`
            The chapter number. (/series/{novel_id}/chapter.{chapter_no})

        Returns
        -------
        :class:`dict`
            A dictionary containing chapter information, or an error indicator.
        """
        urls = {
            ChapterSource.WUXIABOX: f"https://www.wuxiabox.com/novel/{novel_id}_{chapter_no}.html",
            ChapterSource.NOVTALES: f"https://novtales.com/chapter/{novel_id}-{chapter_no}"
        }

        error = {}

        for urlKey in urls.keys():
            url = urls[urlKey]
            req = requests.get(url)
            tries = 5

            if not isinstance(req, requests.Response):
                error = {"status": 400, "chapter_no": chapter_no, "error": "Invalid response object"}
                continue

            while tries > 0 and req.status_code != 200:
                if req.status_code == 404:
                    error = {"status": 404, "chapter_no": chapter_no, "error": "Chapter not found"}
                    break
                else:
                    tries -= 1
                    req = requests.get(url)  # Retry request

            if req.status_code != 200:
                error = {
                    "status": 500,
                    "chapter_no": chapter_no,
                    "error": "Failed to retrieve chapter",
                }
                continue

            chapter_data = {}

            # Parse the chapter content
            if urlKey == ChapterSource.WUXIABOX:
                chapter_data = parsers.parse_chapters_wuxiabox(req)
            elif urlKey == ChapterSource.NOVTALES:
                chapter_data = parsers.parse_chapters_novtales(req)

            if not chapter_data or not chapter_data["body"].strip():
                error = {"status": 204, "chapter_no": chapter_no, "error": "Empty chapter"}
                continue

            return {
                "status": 200,
                "chapter_no": chapter_no,
                "title": chapter_data["title"],
                "url": url,
                "body": chapter_data["body"],
            }

        return error


if __name__ == "__main__":
    client = Client()

    # Read the argument passed from Go
    if len(sys.argv) > 1:
        action = sys.argv[1]

        if action == "import-novel":
            if len(sys.argv) > 2:
                series_id = sys.argv[2]
                result = client.series_info(series_id)
                print(json.dumps(result))
            else:
                print(json.dumps({"status": 400, "error": "No series ID provided"}))
                sys.exit(1)
        elif action == "import-chapter":
            if len(sys.argv) > 3:
                series_id = sys.argv[2]
                chapterNo = sys.argv[3]
                result = client.chapters(series_id, chapterNo)
                print(json.dumps(result))
            else:
                print(json.dumps({"status": 400, "error": "Invalid arguments for import-chapter"}))
                sys.exit(1)
        else:
            print(json.dumps({"status": 400, "error": "Invalid action"}))
            sys.exit(1)
    else:
        print(json.dumps({"status": 400, "error": "No action provided"}))
        sys.exit(1)
