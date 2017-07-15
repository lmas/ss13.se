[Space Station 13 Hub](https://www.ss13.se/)
================================================================================

A server hub dedicated to SS13 and possibly a better replacement for the default
server page at [Byond](https://www.byond.com/games/exadv1/spacestation13).

Status
--------------------------------------------------------------------------------

***In the cloner***

Please note that the code is still in a **experimental stage** at this time and
there is still a lot of work to be done before it's ready for production use.

Features For Players
--------------------------------------------------------------------------------

- A friendlier server list.

- Player count stats and graphs.

- Dedicated pages for each public server, with more detailed info.

- Pages for private servers too, upon request.

Features for server Owners
--------------------------------------------------------------------------------

TODO

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

Add new private server
--------------------------------------------------------------------------------

If you would like to add your server to the private list, and accept being polled
about 4 times an hour, please open a [new ticket](https://github.com/lmas/ss13_se/issues/new)
on the issue tracker on github.

Please provide the following info for your server and write it in your new ticket:

    Title - The public title of your server.
    Game URL - The publicly open address to the game server.
    Site URL - The address to your server's web page, if you have one.

See the file `TODO` for examples.

License
--------------------------------------------------------------------------------

MIT License, see the LICENSE file.

Credits
--------------------------------------------------------------------------------

- [stuicey](https://www.reddit.com/user/stuicey) - Thanks for original idea and [first work](https://www.reddit.com/r/SS13/comments/2p6znr/hub_population_data/).

- [headswe](https://www.reddit.com/user/headswe) - Thanks for [reverse engineered code](http://www.reddit.com/r/SS13/comments/31b5im/a_bunch_of_graphs_for_all_servers/cq11nld) for polling SS13 servers.

- [Hugo14453](https://github.com/Hugo14453) - Thanks for new corgi favicon.

Todo
--------------------------------------------------------------------------------

- Missing tests (a must have for the scraper and poller).

- Missign the favicon.

