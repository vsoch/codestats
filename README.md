# Code Stats

This is a small project that will allow for easily calculating stats for all
repos across an organization, or a single repository of interest. Other metrics
could be added for other entities beyond repository health too!

🚧️ *under development* 🚧️


## Usage

To build the library:

```bash
$ make
```

This will compile an executable, `codestats` that you can interact with.

## Commands

### Repository Stats

```bash
$ ./codestats repo buildsi/build-abi-containers
```

You can also save to file:

```bash
$ ./codestats repo buildsi/build-abi-containers --outfile examples/repo.json
```

You can also pretty print json:

```bash
$ go run main.go repo buildsi/build-abi-containers --pretty
build-abi-containers

{
        "Stats": [
            {
                "Name": "Has-Codeowners",
                "Pass": false
            },
            {
                "Name": "Has-Maintainers",
                "Pass": false
            },
            {
                "Name": "Has-GitHub-Actions",
                "Pass": true
            },
            {
                "Name": "Has-CircleCI",
                "Pass": false
            },
            {
                "Name": "Has-Travis",
                "Pass": false
            },
            {
                "Name": "Has-PullApprove",
                "Pass": false
            },
            {
                "Name": "Has-Glide",
                "Pass": false
            }
        ],
        "Name": "buildsi/build-abi-containers",
        "Branch": "main",
        "Url": "https://github.com/buildsi/build-abi-containers",
        "Stars": 2,
        "Forks": 1,
        "Issues": 3,
        "Language": "Python",
        "Archived": false,
        "CreatedAt": "2021-05-23T20:02:51Z",
        "UpdatedAt": "2021-10-16T12:37:19Z"
    }
}
```


### Organization Stats

If you want to get stats for an org:

```bash
$ ./codestats org buildsi

build-notes
build-si-modeling
Smeagle
build-abi-tests
build-abi-containers
build-sandbox
...
```

You can also add an optional pattern to only parse a subset of repos, or add `--pretty` to pretty print the json.

```bash
$ go run main.go stats buildsi --pattern build-abi-containers --pretty
build-abi-containers
build-abi-containers-results
{
        {
            "Stats": [
                {
                    "Name": "Has-Codeowners",
                    "Pass": false
                },
                {
                    "Name": "Has-Maintainers",
                    "Pass": false
                },
                {
                    "Name": "Has-GitHub-Actions",
                    "Pass": true
                },
                {
                    "Name": "Has-CircleCI",
                    "Pass": false
                },
                {
                    "Name": "Has-Travis",
                    "Pass": false
                },
                {
                    "Name": "Has-PullApprove",
                    "Pass": false
                },
                {
                    "Name": "Has-Glide",
                    "Pass": false
                }
            ],
            "Name": "buildsi/build-abi-containers",
            "Branch": "main",
            "Url": "https://github.com/buildsi/build-abi-containers",
            "Stars": 2,
            "Forks": 1,
            "Issues": 3,
            "Language": "Python",
            "Archived": false,
            "CreatedAt": "2021-05-23T20:02:51Z",
            "UpdatedAt": "2021-10-16T12:37:19Z"
        },
        {
            "Stats": [
                {
                    "Name": "Has-Codeowners",
                    "Pass": false
                },
                {
                    "Name": "Has-Maintainers",
                    "Pass": false
                },
                {
                    "Name": "Has-GitHub-Actions",
                    "Pass": true
                },
                {
                    "Name": "Has-CircleCI",
                    "Pass": false
                },
                {
                    "Name": "Has-Travis",
                    "Pass": false
                },
                {
                    "Name": "Has-PullApprove",
                    "Pass": false
                },
                {
                    "Name": "Has-Glide",
                    "Pass": false
                }
            ],
            "Name": "buildsi/build-abi-containers-results",
            "Branch": "main",
            "Url": "https://github.com/buildsi/build-abi-containers-results",
            "Stars": 1,
            "Forks": 0,
            "Issues": 0,
            "Language": "Python",
            "Archived": false,
            "CreatedAt": "2021-06-08T23:44:24Z",
            "UpdatedAt": "2021-08-29T14:25:50Z"
        }
}
```

or again, save to output file:

```bash
$ go run main.go org buildsi --outfile examples/org.json
```

The output file is the same format, so you can use the same visualization code in [docs](docs) for a web interface.

### Customize Stats

For either the org or repo command, you can provide an optional yaml file
to define a list of stats you want:

```bash
$ go run main.go repo buildsi/build-abi-containers --config examples/all-stats.yaml 
```

Or save to output file, perhaps for a web interface!

```bash
$ go run main.go org buildsi --outfile examples/data.json 
```


### Web Interface

When you generate stats, you can put them in the same directory as the provided [docs](docs)
to generate a nice interface, which will be organized by type and repository.
You can see an example at [https://vsoch.github.io/codestats](https://vsoch.github.io/codestats)

### TODO

 - add better organization / filtering to results UI
 - should there be a common format for a stat, beyond hard coding? E.g., most seem like checking if something exists - this could be YAML


## Previous Art

This was inspired by [bloodorange/oci-stats](https://github.com/bloodorangeio/oci-stats/blob/main/gen-html-for-repo.sh)!
