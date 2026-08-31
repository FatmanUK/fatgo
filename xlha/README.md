# xlha

This is a file-extraction port of lhasa to Go by way of Gemini (and debugged by Claude). Original by Simon Howard. See https://github.com/fragglet/lhasa . See the LICENSE file in this directory.

Gemini Pro did a good job (mostly) of analysing the lhasa repo and creating a single-file extraction library in Go. I mean it works, I'm not saying it's necessarily the most efficient way to do this. I haven't bothered, because LHA archives tend to be very small.

There was one bug Gemini Pro just couldn't track down. I had to ask Claude (Sonnet 5) to fix it instead.
