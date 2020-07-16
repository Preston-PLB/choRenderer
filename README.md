# choRenderer
This is the module that will be used to render out images for chord charts in propresenter7

## What does this do?

This cho renderer go module takes a Chord Pro file and cuts it up into seperate PNG files for each section denoted by a {Comment: "*"} tag.

## How do I do the thing?

Pretty simple

```
go get github.com/Preston-PLB/choRendere
```

then in your code construct a Song struct

```golang
song := choRenderer.Song{Name: "Overcome", PathToFile: "testing/overcome-A.cho", Resolution: canvas.Rect{H: 1920, W:1080}}
```

then render the song

```
song.RenderSong()
```

after the song has been rendered a directory in the directory of the .cho file will be created with all of the images inside

## But Why?

Because I needed a way to elegantly display chord charts on a stage display using pro presenter 7

## Is there going to be an easy to use client for this?

I hope so... I hope so...
