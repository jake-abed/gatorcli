# gator
A Command Line Blog Aggregator written in Go.

## Requisites:

`postgres`
`go`

## Setup

You will need a config file like this one:

```json
{
    "db_url": "yr/postgres",
    "current_user_name": "YOU!"
}
```

Installation

Create a PostgreSQL Database first, then:

`go install https://github.com/jake-abed/gatorcli`

Go nuts!

## Usage

`gator register "name" //Registers a user and logs them in`
`gator login "name" //Logs in a user"`
`gator addfeed "name" "url" // Adds a new feed for the currently logged in user`
`gator follow "url" // Follows a feed if not already followed`
`gator unfollow "url" // Unfollows a feed`
`gator agg "timestring" // Aggregates all followed feeds at a time string rate like 5s, 1m`
`gator browse [optional: limit] // browses the posts from feeds you follow`
