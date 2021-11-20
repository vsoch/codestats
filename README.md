# Org-Stats

This is a small project that will allow for easily calculating stats for all
repos across an organization.

🚧️ *under development* 🚧️


## Usage

To build the library:

```bash
$ make
```

This will compile an executable, `org-stats` that you can interact with.

## Commands

### Stats

If you want to get stats for an org:

```bash
$ ./org-stats stats buildsi

build-notes
build-si-modeling
Smeagle
build-abi-tests
build-abi-containers
build-sandbox
...
```

Currently we print the repository names, and then print the result object. 
This will eventually generate json output that we can save and pipe into a web interface (to be developed).


### TODO

 - save output to json
 - allow for custom config to specify metrics/stats desired
 - allow to provide GitHub token for use
 - need way to also customize rendering of stats - should go handle the UI generation?
 - should there be a common format for a stat, beyond hard coding? E.g., most seem like checking if something exists - this could be YAML


## Previous Art

This was inspired by [bloodorange/oci-stats](https://github.com/bloodorangeio/oci-stats/blob/main/gen-html-for-repo.sh)!
