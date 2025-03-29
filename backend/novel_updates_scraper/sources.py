from enum import Enum, auto


class NovelSource(Enum):
    """
    Represents the source website for novels.

    Each member represents a different website from which novels can be sourced.
    """
    NOVELBIN = auto()  # Represents NovelBin website
    LIGHTNOVELWORLD = auto()  # Represents LightNovelWorld website
    WUXIABOX = auto()  # Represents WuxiaBox website
    NOVTALES = auto()  # Represents NovTales website


class ChapterSource(Enum):
    """
    Represents the source website for individual chapters.

    Similar to NovelSource, but specifically for chapters, allowing for different chapter sources even if the novel originates from a different site.
    """
    WUXIABOX = auto()  # Represents WuxiaBox website (for chapters)
    NOVTALES = auto()  # Represents NovTales website (for chapters)
