# OfflineRSS

OfflineRSS is a tool written in Go for syncing [RSS](https://en.wikipedia.org/wiki/RSS) between a hosted [feed aggregator](https://en.wikipedia.org/wiki/News_aggregator) and local RSS readers.

It generates databases for supported RSS readers that can be used offline. Any changes done to these databases will be synced back to the aggregator next time the program is run whilst also updating the local databases with everything new.

## Why?

* Allows you to use any supported reader with any supported feed aggregator
* Allows you to use your favorite reader offline. Great if you don't always have an internet connection.
* Makes all user interaction feel instant because they are synced later.

## Supported servers (feed aggregators)

* [Miniflux](https://miniflux.app) Apache-2.0
* [NewsBlur](https://newsblur.com) MIT

## Supported clients (readers)

* [FeedReader](https://jangernert.github.io/FeedReader) GPL-3.0
* [Newsboat](https://newsboat.org) MIT
* [QuiteRSS](https://quiterss.org) GPL-3.0

## Get started

Install OfflineRSS (it will install to `~/go/bin/offlinerss` unless changed)

```
go install github.com/limero/offlinerss@latest
```

You can now run OfflineRSS and it will prompt you for what server and clients to use and save this information in `~/.config/offlinerss/config.json`. It will then create local databases for your chosen client(s).

By default, the paths for the generated databases will be in `~/.local/share/offlinerss`. You can either change these or symlink them to the correct locations, see instructions below.

**WARNING! This will remove any existing databases at these locations!**

### FeedReader

```
ln -sf ~/.local/share/offlinerss/feedreader/feedreader-7.db ~/.local/share/feedreader/data/feedreader-7.db
```

### Newsboat

```
ln -sf ~/.local/share/offlinerss/newsboat/cache.db ~/.local/share/newsboat/cache.db
ln -sf ~/.local/share/offlinerss/newsboat/urls ~/.config/newsboat/urls
```

### QuiteRSS

```
ln -sf ~/.local/share/offlinerss/quiterss/feeds.db ~/.local/share/QuiteRss/QuiteRss/feeds.db
```

## Limitations

* Only one server can be used at a time. However, all clients can be enabled simultaneously.
* Clients can't be synced with each other, without syncing with the server first.

## License

Licensed under the MIT license. (See LICENSE)
