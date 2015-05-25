[Space Station 13 Hub](http://www.ss13.se/)
================================================================================

A server hub dedicated to SS13 and possibly a better replacement for the default
server page at [Byond](http://www.byond.com/games/exadv1/spacestation13).

Features
--------------------------------------------------------------------------------

- A friendlier server list, which is also sortable.

- Dedicated pages for each public server, with more detailed info.

- Pages for private servers too, upon request.

- Player count stats and graphs.

Data Source
--------------------------------------------------------------------------------
The data used for most of the **public** servers is scraped from the [Byond](http://www.byond.com/games/exadv1/spacestation13) page.
Relying on Byond means we're affected by their server downtime (no page, no data),
but we will automatically discover any new public servers.

All **private** servers, and some public ones, are being polled directly for data.
That means a list of private servers have to be maintained manually, but we're
not affected by any downtimes (as long as the game servers themselves are up
and running). This is usually a more reliable method, but it's more expensive to
run since we have to send multiple network requests (instead of one to Byond).

Both methods can be affected by spoofing attacks, done by bad servers, and so
there's no way to guarantee that the calculated stats and graphs are 100% correct.

But then again it's just some silly numbers for a bunch of games.

Source Code
--------------------------------------------------------------------------------
This is a open source project running on [Python](https://www.python.org/) and [Django](https://www.djangoproject.com/) (for the dynamic web
stuff), a little [Go](https://golang.org/) (for some behind-the-scene scripting) and some [gnuplot](http://www.gnuplot.info/) (for some
ugly but interesting graphs).

The source code itself is hosted at [github.com](https://github.com/lmas/ss13_se) and is **open for any and all
contributions**.

License
--------------------------------------------------------------------------------
MIT License, see the LICENSE file.

Credits
--------------------------------------------------------------------------------
- A. Svensson - Author and host.

- [stuicey](https://www.reddit.com/user/stuicey) - Thanks for original idea and [first work](https://www.reddit.com/r/SS13/comments/2p6znr/hub_population_data/).

- [headswe](https://www.reddit.com/user/headswe) - Thanks for [reverse engineered code](http://www.reddit.com/r/SS13/comments/31b5im/a_bunch_of_graphs_for_all_servers/cq11nld) for polling SS13 servers.

- [Hugo14453](https://github.com/Hugo14453) - Thanks for new corgi favicon.

Todo
--------------------------------------------------------------------------------

**Bugs**

- Better colors for the warning/offline server notices.

**New features**

- Show note about data source for each server (scraped/polled).

- Rewrite the whole update script to go.

- Show fancy web stats from the web server logs.

- stats page:
    - time to run update
    - log of recently added/removed servers
    - server graphs
    - player graphs + average

**Suggestions**

- Frontpage of some sort?

- Page to send in requests to add new private servers?

- Player growth rate for each server (+/- compared to avg.)?

- Use some kind of banner/logo?
