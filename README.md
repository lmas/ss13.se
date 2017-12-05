[Space Station 13 Hub](https://www.ss13.se/)
================================================================================

A server hub dedicated to SS13 and possibly a better replacement for the default
server page at [Byond](https://www.byond.com/games/exadv1/spacestation13).

Status
--------------------------------------------------------------------------------

***Beta mode live at [www.ss13.se](https://www.ss13.se/)***

The code has ben run in production for some months now and seems to be stable
enough. I guess we're now back in beta?

Features For Players
--------------------------------------------------------------------------------

- A friendlier server list.

- Detailed pages for most public servers.

- Player count stats and graphs.

- Predict when it's a "good" time for you to play.

Data Source
--------------------------------------------------------------------------------

The data used for the **public** servers is scraped from the [Byond](http://www.byond.com/games/exadv1/spacestation13) page.
Relying on Byond means we're affected by their server downtime (no page, no data),
but we will automatically discover any new public servers.

Since bad hosts and owners can spoof a server's reported player count, there's
no way to guarantee that the calculated stats and graphs are 100% correct.

But then again it's just some silly numbers for a game.

Missing a server?
--------------------------------------------------------------------------------

Since we're scraping only **public** servers, we're missing any **private** that's
hidden from the Byond page (I also discarded the poller, so no more polling
**private** servers from a VIP list).

We also see a lot of **public** servers coming and going, or merely changing
names a lot. So any servers that hasn't been seen for 3 days, or more, will be
automagically removed from the list.

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

- Missing tests.

- Missign the favicon.

- Add about page (redirect to github).

- Add graph showing next week's predicted player count for all servers?

- Investigate more options for prediction.

- Investigate community approval of tracking public members (visited servers, play time, predictions)?
