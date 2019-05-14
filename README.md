## acsnip

This code is awful, and you shouldn't use it. But if I need to do this again in the future, I'll miss it.

## How to Use

1. `go build`
1. Edit the asciinema recording file, delete any content you wish to skip over.
1. Find the timestamp of a step that is a long time after a previous step (the time period you want to "snip")
1. Run acsnip with the filename of the cast as one input, and the timestamp as another, for example `./acsnip ~/asciicast.cast 6.016997`
1. acsnip will output a file in the same location as the original file, with the extension `.snip.cast`, for instance `~/asciicast.snip.cast`
1. Snip out all uncomfortable breaks/undesireable output.
1. :parrot:
