[Space Station 13 Hub](http://www.ss13.se/)
================================================================================

A server hub dedicated to SS13 and possibly a better replacement for the default
server page at [Byond](http://www.byond.com/games/exadv1/spacestation13).

The source code has recently been rewritten from Python to Go, for better
performance and stability, as well as cleaner structure and ease of use for
the server host.

Please note that the code is still in a **experimental stage** at this time and
there is still a lot of work to be done before it's ready for production use.

Features For Players
--------------------------------------------------------------------------------

- A friendlier server list, which is also sortable.

- Player count stats and graphs.

- Dedicated pages for each public server, with more detailed info.

- Pages for private servers too, upon request.

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

- Prevent locking the whole db when updating.

- Update static files to newer versions.

- Fix and clean up the tooltips in the server details template.

- Use the same format for the verbose timestamps.

**New features**

- Tests (a must have for the scraper and poller).

- Show note about data source for each server (scraped/polled).

- Index page:
    - TODO

- About page:
    - Move all notices about Byond sources to this page.
    - A way to contact me (reddit, github etc.).
    - Info on how to request adding a new server to be polled.

- Stats page:
    - Server graphs
    - Player graphs + average
    - Number of online/warn/offline servers.
    - Total/average number of online players.
    - Log of recently added/removed servers.
    - The number of data points since start.
    - Time since last update.
    - Time to run update

**Suggestions**

- Live player count API and script to embed in external sites (https://github.com/lmas/ss13_se/issues/2)?

- Player growth rate for each server (+/- compared to avg.)?

