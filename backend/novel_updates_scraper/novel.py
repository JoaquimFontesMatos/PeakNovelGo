from typing import TypedDict, List, Dict, Union, NotRequired

class NovelMetadata(TypedDict):
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
