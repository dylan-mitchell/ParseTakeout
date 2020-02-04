# ParseTakeout

Recently I did a [Google Takeout](https://takeout.google.com/settings/takeout) which lets you download all of the data that Google has collected from you. I thought it would be interesting to analyze/ play around with this data. The unfortunate part is the data doesn't have a standard data type. There are `.json`, `.csv`, `.txt`, and `.html` files. Of these, I think the `.html` are the trickiest to parse through as there was no standard between them. (I.E the search.html has a different structure than the youtube-watch-history.html file)

I added `My-Activity-Developers.html` as an example file because it contains only the data collected on me from `developers.google.com` and contains a lower amount of personal information compared to data collected from other Google services (Chrome, Search, Android, ...)
