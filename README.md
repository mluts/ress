# ress
A Fast RSS Aggregator

## Motivation
I'm working on this rss aggregator as an alternative for tiny tiny rss, 
it should have more consistent, minimalistic api and should be faster
when running on my raspberry pi.

## API

```
GET 		/feeds
POST 		/feeds
GET 		/feeds/[0-9]+
DELETE 	/feeds/[0-9]+
GET 		/feeds/[0-9]+/items
POST 		/feeds/[0-9]+/items/read
DEETE 	/feeds/[0-9]+/items/read
```
