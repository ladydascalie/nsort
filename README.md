# nsort

This is a utility for sorting files by kind or by extension.

## Usage

```man
nsort [options] [path]
```

## Options

```man
-t [dir]
    Set the target directory.

-by-kind
    Sort by kind.

-upd [ext]:[kind]
    Update a mapping.
    Example: nsort -upd mp3:Audio

-del [ext]:[kind]
    Delete a mapping.
    Example: nsort -del mp3:Audio

-map [ext]:[kind]
    Add a mapping.
    Example: nsort -map mp4:Video
```
