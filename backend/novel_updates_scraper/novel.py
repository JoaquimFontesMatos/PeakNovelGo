from typing import TypedDict, List, Dict, Union, NotRequired


class NovelMetadata(TypedDict):
    """
    Represents metadata for a novel.  Fields are optional (NotRequired).

    Attributes:
        title (str, optional): The title of the novel.
        image (str, optional): URL or path to an image representing the novel.
        genre (List[Dict[str, str]], optional): A list of dictionaries, each representing a genre.  Each dictionary should have at least one key-value pair (e.g., {"name": "Fantasy"}).
        tags (List[Dict[str, str]], optional): A list of dictionaries, each representing a tag. Each dictionary should have at least one key-value pair (e.g., {"name": "Romance"}).
        rating (Union[str, float], optional): The rating of the novel (can be a string or a float).
        language (Dict[str, str], optional): A dictionary representing the language of the novel (e.g., {"name": "English").
        authors (List[Dict[str, str]], optional): A list of dictionaries, each representing an author. Each dictionary should have at least one key-value pair (e.g., {"name": "John Doe"}).
        year (str, optional): The year the novel was published or first released.
        status (str, optional): The status of the novel (e.g., "Completed", "Ongoing").
        release_freq (str, optional): The release frequency of the novel (e.g., "Weekly", "Monthly").
        description (str, optional): A description of the novel.
        latest_chapter (int, optional): The number of the latest chapter.


    Example:
        >>> metadata = {
        ...     "title": "The Lord of the Rings",
        ...     "image": "https://exampleimage.com",
        ...     "genre": [{"name": "Fantasy"}],
        ...     "tags": [{"name": "Action"}],
        ...     "authors": [{"name": "J.R.R. Tolkien"}],
        ...     "year": "1954",
        ...     "rating": 4.8,
        ...     "release_freq": "N/A",
        ...     "description": "A story about funny little rings",
        ...     "language": {"name": "English"},
        ...     "status": "Completed",
        ...     "latest_chapter": 0
        ... }
        >>>
    """
    title: NotRequired[str]
    image: NotRequired[str]
    genre: NotRequired[List[Dict[str, str]]]
    tags: NotRequired[List[Dict[str, str]]]
    rating: NotRequired[Union[str, float]]
    language: NotRequired[Dict[str, str]]
    authors: NotRequired[List[Dict[str, str]]]
    year: NotRequired[str]
    status: NotRequired[str]
    release_freq: NotRequired[str]
    description: NotRequired[str]
    latest_chapter: NotRequired[int]
