# Space Engineers: Link Steam & Discord

A [Torch][torch] server plugin and associated backend for
[Space Engineers][space-engineers], offering players the ability to link their
Steam & Discord together, and allowing server administrators further
integrations with the data.

Both the plugin and backend are open-source so it can be self-hosted, and
server administrators remain control of the data.

Source code: on [github][github].

[torch]: https://torchapi.net/
[space-engineers]: https://www.spaceengineersgame.com/
[github]: https://github.com/StalkR/Space-Engineers-Link-Steam-Discord

## Plugin

See the `plugin/` subdirectory.

Releases: on the [Torch plugin page][plugin].

Users can use the `!link` command to link their Steam & Discord. If they
haven't linked yet, it will open their browser to the linking backend.

Server admininistrators can configure the backend URL in the plugin settings.

[plugin]: https://torchapi.net/plugins/item/64827390-c847-41fb-adc3-2719b1d1b536

## Backend

See the `backend/` subdirectory.

A small web app written in Go, implementing both Steam OpenID 2.0 and Discord
OAuth2 sign-in. It saves the Steam ID & Discord ID mapping to an SQLite3
database.

Configuration:
- see `config.example.json`
- create a Discord OAuth application, write down the Client ID & Secret
- configure `secret` with random bytes, it's used for state authentication
- configure some `api_keys` if you need apps to use the `/lookup` API
- run with `go run .` or build the binary and deploy elsewhere
- expose the app behind an HTTPS reverse proxy, or tweak it to your needs

API:
- `/steam/check/<steamID>`: returns 200 if linked, 404 if not; used by the
  plugin to check if a given steam ID is linked or not
- `/lookup/steam/<steamID>`: (API key required) returns the Discord ID if
  found, otherwise 404
- `/lookup/discord/<discordID>`: (API key required) returns the Steam ID if
  found, otherwise 404

## Origins

It is based on the `!sedb link` feature from [@Bishbash777][bish] in
[SEDiscordBridge][sedb] ([github][sedb-github]), which also offers a fully
featured discord bridge, but it uses a fixed backend that is not open-source,
anyone can query, and does not authenticate steam ID so fake mappings can be
inserted.

[bish]: https://github.com/Bishbash777
[sedb]: https://torchapi.net/plugins/item/3cd3ba7f-c47c-4efe-8cf1-bd3f618f5b9c
[sedb-github]: https://github.com/Bishbash777/SEDB-RELOADED

## Bugs, comments, questions

Create a [new issue][new-issue].

[new-issue]: https://github.com/StalkR/Space-Engineers-Link-Steam-Discord/issues/new
