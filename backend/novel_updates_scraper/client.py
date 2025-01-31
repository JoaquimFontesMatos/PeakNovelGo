import sys
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
        id : :class:`int`
            The id of the series. (/series/{ID})

        Returns
        -------
        :class:`dict`
            A dictionary containing information about the series.
            Contains all information and links for the series.
        """
        req = self.req.get(f"https://www.novelupdates.com/series/{series_id}")
        return parsers.parseSeries(req)

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
        series_id = sys.argv[1]
        result = client.series_info(series_id)
        print(json.dumps(result))
    else:
        print(json.dumps({"error": "No series ID provided"}))
        sys.exit(1)
