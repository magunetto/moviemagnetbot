# Movie Magnet Bot

🤖 A telegram bot for movies

[![Build Status](https://travis-ci.org/magunetto/moviemagnetbot.svg)](https://travis-ci.org/magunetto/moviemagnetbot)
[![Coverage Status](https://coveralls.io/repos/github/magunetto/moviemagnetbot/badge.svg?branch=master)](https://coveralls.io/github/magunetto/moviemagnetbot?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/magunetto/moviemagnetbot)](https://goreportcard.com/report/github.com/magunetto/moviemagnetbot)

## What can I do with this bot

- Search information about movies and TVs
- Search magnet links of movies and TVs and download them
- Download any links (Magnet/eD2k/FTP) automatically
- View and clear download history

*[Complete feature list and future plans](https://github.com/magunetto/moviemagnetbot/wiki/Features)*

## How to use it

- Prepare dependencies: PostgreSQL database, [Telegram bot](https://core.telegram.org/bots), and [TMDB API](https://www.themoviedb.org/documentation/api)
- Build, configure [environment variables](https://github.com/magunetto/moviemagnetbot/search?q=Getenv), and run the bot server
- Talk to your bot user on Telegram (*[Video example 1](https://t.me/moviemagnetfm/6), [Video example 2](https://t.me/moviemagnetfm/7)*)
- Subscribe personal RSS feed in download tools

## Ways to get involved

- Join our [channel](https://t.me/moviemagnetfm) and [user group](https://t.me/moviemagnetusers) on Telegram
- [Open an issue](https://github.com/magunetto/moviemagnetbot/issues/new/choose) when you have an idea or found a bug

## How to contribute

1. Have Go installed
1. Fork it and start hacking

    ```bash
    cd moviemagnetbot
    cd cmd/moviemagnetbot
    go build && PORT="8000" ./moviemagnetbot
    ```

1. Open a pull request when you improved something

## Alternatives

- [Netflix](https://www.netflix.com/)
- [showRSS](https://showrss.info/)
- [Popcorn Time](https://popcorn-time.to/)
