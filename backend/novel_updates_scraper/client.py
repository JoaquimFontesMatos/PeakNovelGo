import sys

import requests
from .request import Request
from . import parsers
import json


class Client:
    def __init__(self):
        self.req = Request()

    def get_latest_feed(self):
        """Gets the latest updates from NovelUpdates.

        Parameters
        ----------
        None

        Returns
        -------
        :class:`list`
            A dictionary containing the latest novel updates from NovelUpdates.
            Contains all information and links for each update.
        """
        req = self.req.get("https://www.novelupdates.com/")
        return parsers.parseFeed(req)

    def search_series(self, name):
        """Searches for a series and gets back the top 25 results (first page).

        Parameters
        ----------
        name : :class:`str`
            The name of the series to search for.

        Returns
        -------
        :class:`list`
            A dictionary containing the top 25 results for the search.
            Contains all information and links for each result.
        """
        req = self.req.get(f"https://www.novelupdates.com/?s={name}")
        return parsers.parseSearch(req)

    def series_info(self, series_id):
        """Gets information about a series.

        Parameters
        ----------
        series_id : :class:`str`
            The ID of the series. (/series/{ID})

        Returns
        -------
        :class:`dict`
            A dictionary containing information about the series.
            If an error occurs, the dictionary will contain an "error" key.
        """
        try:
            # Make the HTTP request
            req = self.req.get(
                f"https://www.lightnovelworld.co/novel/{series_id}"
            )

            if req.status_code == 404:
                return {"status": 404, "error": "Series not found"}

            req.raise_for_status()  # Raise an exception for HTTP errors (4xx, 5xx)

            # Parse the response
            return parsers.parseSeries(req)
        except requests.exceptions.ConnectionError:
            return {"status": 503, "error": "Network connection down"}
        except requests.exceptions.Timeout:
            return {"status": 503, "error": "Request timed out"}
        except requests.exceptions.HTTPError as e:
            return {"status": 503, "error": f"Source website down: {e}"}
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
        url = f"https://www.lightnovelworld.co/novel/{series_id}/chapter-{i}"
        req = self.req.get(url)  # Make the request

        tries = 5

        while tries > 0 and req.status_code != 200:
            if req.status_code == 404:
                return {"status": 404, "chapter_no": i, "error": "Chapter not found"}
            else:
                tries -= 1
                req = self.req.get(url)  # Retry request

        if req.status_code != 200:
            return {
                "status": 500,
                "chapter_no": i,
                "error": "Failed to retrieve chapter",
            }

        # Parse the chapter content
        chapter_data = parsers.parseChapters(req)

        if not chapter_data["body"].strip():
            return {"status": 204, "chapter_no": i, "error": "Empty chapter"}

        return {
            "status": 200,
            "chapter_no": i,
            "title": chapter_data["title"],
            "url": url,
            "body": chapter_data["body"],
        }

    def series_groups(self, series_id):
        """Gets the groups that are translating a series.

        Parameters
        ----------
        id : :class:`int`
            The id of the series. (/series/{ID})

        Returns
        -------
        :class:`list`
            A dictionary containing information about the groups.
            Contains all information and links for each group.
        """
        req = self.req.get(f"https://www.novelupdates.com/series/{series_id}")
        extras = parsers.parseSeries(req, extras=True)
        data = {
            "action": "nd_getgroupnovel",
            "mygrr": extras["grr_groups"],
            "mypostid": extras["postid"],
        }
        req2 = self.req.post(
            "https://www.novelupdates.com/wp-admin/admin-ajax.php", data=data
        )
        return req2.text


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
