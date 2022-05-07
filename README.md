# bib

Terminal Bible reader.

There's not much yet, just simple navigation and side by side visualization of
all available translations.

## Development

### Debugging

To debug the application, invoke the debugger console in one shell:

```shell
dlv debug --headless --listen :4747 main.go
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
