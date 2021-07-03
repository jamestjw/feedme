# feedme

Simple tool to create an RSS feed from an instagram profile. This comes with both a CLI tool and a webapp.

## Installation
```
# Build both
make 

# Build just the webapp
make feedme-webapp

# Build just the cli tool
make feedme-cli
```

## Usage

### Webapp
``` bash
./feedme-cli -p 8080
```

Access `http://localhost:8080/instagram/user/{instagram-name}` to get the XML feed for a particular user.

### CLI tool
Suppose that we want an RSS feed of McDonald's instagram handle @mcdonalds, we just have to do the following
``` bash
./feedme-cli -t mcdonalds
```
