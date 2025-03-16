import json
import sys

import requests

from . import parsers
from .request import PlaywrightScraper


class Client:
    def __init__(self):
        self.req = PlaywrightScraper()

    def series_info(self, series_id):
        """Gets information about a series."""
        try:
            latest_chapter_url=f"https://novtales.com/novel/{series_id}"
            url = f"https://novelbin.com/b/{series_id}"
            content = self.req.scrape(url)
            content2 = self.req.scrape(latest_chapter_url)
            if content is None:
                return {"status": 503, "error": "Failed to fetch the page"}

            # Check for 404 status (if applicable)
            if "Page not found" in content:  # Replace with a condition specific to your website
                return {"status": 404, "error": "Series not found"}

            # Parse the series information
            latest_chapter=parsers.parse_latest_chap(content2)

            return parsers.parse_series(content, latest_chapter)
        except Exception as e:
            return {"status": 503, "error": f"Unexpected error: {e}"}

    def chapters(self, series_id, i=1):
        """Gets the chapters of a series.

        Parameters
        ----------
        series_id : :class:`int`
            The id of the series. (/series/{ID})

        Returns
        -------
        :class:`dict`
            A dictionary containing chapter information, or an error indicator.
        """
        urls = {
            "wuxiabox": f"https://www.wuxiabox.com/novel/{series_id}_{i}.html",
            "novtales": f"https://novtales.com/chapter/{series_id}-{i}"
        }

        error = {}

        for urlKey in urls.keys():
            url = urls[urlKey]
            req = requests.get(url)  # Make the request

            tries = 5

            if not isinstance(req, requests.Response):
                error = {"status": 400, "chapter_no": i, "error": "Invalid response object"}
                continue

            while tries > 0 and req.status_code != 200:
                if req.status_code == 404:
                    error = {"status": 404, "chapter_no": i, "error": "Chapter not found"}
                    break
                else:
                    tries -= 1
                    req = requests.get(url)  # Retry request

            if req.status_code != 200:
                error = {
                    "status": 500,
                    "chapter_no": i,
                    "error": "Failed to retrieve chapter",
                }
                continue

            # Parse the chapter content
            if urlKey == 'wuxiabox':
                chapter_data = parsers.parse_chapters_wuxiabox(req)
            elif urlKey == "novtales":
                chapter_data = parsers.parse_chapters_novtales(req)

            if not chapter_data or not chapter_data["body"].strip():
                error = {"status": 204, "chapter_no": i, "error": "Empty chapter"}
                continue

            return {
                "status": 200,
                "chapter_no": i,
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
                chaptNo = sys.argv[3]
                result = client.chapters(series_id, chaptNo)
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
