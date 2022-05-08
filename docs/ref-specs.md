# Reference specs

1. Book names can be standalone or be numbered:
    1. Standalone: Genesis, Revelation, etc.
    1. Numbered: 1 Corinthians, 2 Kings, etc.
1. References can be partial:
    1. Only book name (e.g., Genesis)
    1. Book name and chapter number (e.g., Psalm 82)
    1. Book name, chapter name and verse number (e.g., John 1:1)
1. Book numbers don't have to be expelled completely:
    1. Dots, optionally, can be included in abbreviations:
        1. Mat. should resolve to the gospel of Matthew (dot should be ignored)
        1. Mar should resolve to the gospel of Mark
    1. Ambiguous references should be unresolved:
        1. Ma should be an ambiguous reference (can resolve to both Malachi, Mark
           Matthew) and, thus, yield an error
    1. Common abbreviations should be considered when possible:
        1. Jo (common abbreviation to the gospel of John), which would be
           considered an ambiguous reference in a normal scenario (can be solved
           to both John and Job), will be solved to John, since it's a fairly
           common abbreviation used for it
1. References should support many languages:
    1. When a Portuguese translation is loaded, "Marcos" should load the correct
       book (gospel of Mark)
    1. When an English translation is loaded, "Mark" should load the correct
       book (gospel of Mark)
    1. etc.
1. References should support verse ranges and lists:
    1. Examples of verse ranges:
        1. Verses within a common chapter: Mat. 24:3-14
        1. Verses from multiple chapter: Mat. 1:1-28:20
        1. Descending ranges should be invalid: John 1:3-1
    1. Example of verse list:
        1. Adjacent verses: John 1:1, 2
        1. Verses within the same chapter: John 1:1, 2, 14
        1. Verses from multiple chapters: John 1:1; 20:28
    1. Interpolation of ranges and lists: John 1:1, 2, 14-18; 6:58; 14:1-6

## Questions

1. Is it possible to translate a reference to a universal system (e.g., book
   name to number that can be applied to any translation.)?
   1. Has to be compatible with the apocrypha.
