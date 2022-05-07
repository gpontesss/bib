# bib

Terminal Bible resources for study.

The project is intended to be a versatile studying tool supporting multiple
translations, commentaries, cross references, and whatever else useful to
extracting more meaning from the Biblical text.

So far, just simple navigation and side by side visualization of all available
translations are implemented.

## Goals

+ **Support of multiple languages:** Deep study usually involves reading
    passages in their original language (Greek for the NT; Hebrew and Aramaic
    for the OT). There are also translations that became part of the tradition
    of the faith (Septuagint, abbreviated as LXX, a translation of the OT to
    Greek, which is the source of quotes from the NT; Vulgate, a translation of
    the Scriptures to Latin by Jerome, which became the main translation used
    throughout the middle ages in the West; etc.) which require the support of
    ancient languages. Also, extensibility is a core principle for the project.
    That being, allowing the import of custom translations in other languages is
    a generic functionality wanted;
+ **Navigation to Bible commentaries**: There are various commentaries from
    different perspectives available (for a glimpse, checkout [Bible Hub's
    page], which more than 100 commentaries). Gathering them in one place is a
    wanted feature to avoid wasting time navigating through websites to compare
    different commentaries. Also, it's desired for it to be extensible,
    allowing for custom commentaries to be imported;
+ **Linking of cross references**: Cross references link passages related to a
    verse that may be the place quoted, deal with the same subject, provide the
    context needed to understand the current verse, etc. Multiple commentators
    also have concordances, which provides this kind of linking. Easy access to
    these is a wanted a feature;
+ **Vim-like bindings**: I mostly live in the terminal. Most of the features I
    listed already exist in GUI software (Take a look at [Logos], for example).
    Since I'm most comfortable with Vim bindings, present in a lot of terminal
    applications, one of the main goals is to bring this idiom to the project.
    Nonetheless, with a well-built core, other front-ends can be created to
    satisfy different users.

## Development

### Debugging

To debug the application, invoke the debugger console in one shell:

```shell
dlv debug --headless --listen :4747 main.go --version=/path/to/version.tsv
```

And execute the application in another:

```shell
dlv connect :4747
```

### Troubleshooting

MacOS seems to have a problem related with pkg-config .pc files related with
`ncurses`. If you have related with it, there's a [tutorial from Michael Cook
that should fix it in a hacky
way](https://mrcook.uk/how-to-install-go-ncurses-on-mac-osx).

[Bible Hub's page]:https://biblehub.com/commentaries/
[Logos]: https://www.logos.com/
