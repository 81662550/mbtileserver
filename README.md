mbtileserver
[![Build Status](https://travis-ci.org/consbio/mbtileserver.svg?branch=master)](https://travis-ci.org/consbio/mbtileserver)
[![Coverage Status](https://coveralls.io/repos/github/consbio/mbtileserver/badge.svg?branch=master)](https://coveralls.io/github/consbio/mbtileserver?branch=master)
[![GoDoc](https://godoc.org/github.com/consbio/mbtileserver?status.svg)](http://godoc.org/github.com/consbio/mbtileserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/consbio/mbtileserver)](https://goreportcard.com/report/github.com/consbio/mbtileserver)

============

A simple Go-based server for map tiles stored in [mbtiles](https://github.com/mapbox/mbtiles-spec) 
format.

Requires Go 1.8.

It currently provides support for `png`, `jpg`, and `pbf` (vector tile)
tilesets according to version 1.0 of the mbtiles specification.  Tiles
are served following the XYZ tile scheme, based on the Web Mercator
coordinate reference system. UTF8 Grids are also supported.

In addition to tile-level access, it provides:
* TileJSON 2.1.0 endpoint for each tileset, with full metadata
from the mbtiles file.
* a preview map for exploring each tileset.
* a minimal ArcGIS tile map service API (work in progress)

We have been able to host a bunch of tilesets on an 
[AWS t2.nano](https://aws.amazon.com/about-aws/whats-new/2015/12/introducing-t2-nano-the-smallest-lowest-cost-amazon-ec2-instance/)
virtual machine without any issues.


## Goals
* Provide a web tile API for map tiles stored in mbtiles format
* Be fast
* Run on small resource cloud hosted machines (limited memory & CPU)
* Be easy to install and operate


## Installation
You can install this project with 

```sh
go get github.com/consbio/mbtileserver
```

This will create and install an executable called `mbtileserver`.


This uses `Govendor` tool, so dependencies ship with the repo for easier builds.

Dependencies:
* [github.com/labstack/echo](https://github.com/labstack/echo)
* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
* [github.com/spf13/cobra](https://github.com/spf13/cobra)
* [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus)
* [golang.org/x/crypto/acme/autocert](https://golang.org/x/crypto/acme/autocert)

Development Dependencies (only needed when modifying the code):
* [github.com/shurcooL/vfsgen](https://github.com/shurcooL/vfsgen)

On Windows, it is necessary to install `gcc` in order to compile `mattn/go-sqlite3`.  
MinGW or [TDM-GCC](https://sourceforge.net/projects/tdm-gcc/) should work fine.


If you experience very slow builds each time, it may be that you need to first run
```
go build -a .
```

to make subsequent builds much faster.


## Usage
From within the repository root ($GOPATH/bin needs to be in in your $PATH):
```
$  mbtileserver --help
Serve tiles from mbtiles files.

Usage:
  mbtileserver [flags]

Flags:
  -c, --cert string     X.509 TLS certificate filename.  If present, will be used to enable SSL on the server.
  -d, --dir string      Directory containing mbtiles files. (default "./tilesets")
      --domain string   Domain name of this server
      --dsn string      Sentry DSN
  -h, --help            help for mbtileserver
  -k, --key string      TLS private key
      --path string     URL root path of this server (if behind a proxy)
  -p, --port int        Server port. (default 8000)
  -t, --tls				Auto TLS using Let's Encrypt
  -r, --redirect		Redirect HTTP to HTTPS
  -v, --verbose         Verbose logging
```

So hosting tiles is as easy as putting your mbtiles files in the `tilesets`
directory and starting the server.  Woo hoo!

You can have multiple directories in your `tilesets` directory; these will be converted into appropriate URLs:

`mytiles/foo/bar/baz.mbtiles` will be available at `/services/foo/bar/baz`.

When you want to remove, modify, or add new tilesets, simply restart the server process.

If a valid Sentry DSN is provided, warnings, errors, fatal errors, and panics will be reported to Sentry.

If `redirect` option is provided, the server also listens on port 80 and redirects to port 443.

## Specifications
* expects mbtiles files to follow version 1.0 of the [mbtiles specification](https://github.com/mapbox/mbtiles-spec).  Version 1.1 is preferred.
* implements [TileJSON 2.1.0](https://github.com/mapbox/tilejson-spec)


## Creating Tiles
You can create mbtiles files using a variety of tools.  We have created
tiles for use with mbtileserver using:
* [TileMill](https://www.mapbox.com/tilemill/)  (image tiles)
* [tippecanoe](https://github.com/mapbox/tippecanoe)   (vector tiles)
* [tpkutils](https://github.com/consbio/tpkutils)  (image tiles from ArcGIS tile packages)


## Examples 

TileJSON API for each tileset:
`http://localhost/services/states_outline`

returns something like this;
```
{
  "bounds": [
    -179.23108,
    -14.601813,
    179.85968,
    71.441055
  ],
  "center": [
    0.314297,
    28.419622,
    1
  ],
  "credits": "US Census Bureau",
  "description": "States",
  "format": "png",
  "id": "states_outline",
  "legend": "[{\"elements\": [{\"label\": \"\", \"imageData\": \"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAUCAYAAACNiR0NAAAAAXNSR0IB2cksfwAAAAlwSFlzAAAOxAAADsQBlSsOGwAAAGFJREFUOI3tlDEOgEAIBClI5kF+w0fxwXvQdjZywcZEtDI31YaQgWrdPsYzAPFGJCmmEAhJGzCash0wSVE/HHnlKcDMfrPXYgmXcAl/JswK6lCrz89BdGVm1+qrH0bbWDgA3WwmgzD8ueEAAAAASUVORK5CYII=\"}], \"name\": \"tl_2015_us_state\"}]",
  "map": "http://localhost/services/states_outline/map",
  "maxzoom": 4,
  "minzoom": 0,
  "name": "states_outline",
  "scheme": "xyz",
  "tags": "states",
  "tilejson": "2.1.0",
  "tiles": [
    "http://localhost/services/states_outline/tiles/{z}/{x}/{y}.png"
  ],
  "type": "overlay",
  "version": "1.0.0"
}
```

It provides most elements of the `metadata` table in the mbtiles file.


XYZ tile endpoint for individual tiles:
`http://localhost/services/states_outline/tiles/{z}/{x}/{y}.png`


If UTF-8 Grid data are present in the mbtiles file, they will be served up over the
grid endpoint:
`http://localhost/services/states_outline/tiles/{z}/{x}/{y}.json`

Grids are assumed to be gzip or zlib compressed in the mbtiles file.  These grids
are automatically spliced with any grid key/value data if such exists in the mbtiles
file.



The map endpoint:
`http://localhost/services/states_outline/map`

provides an interactive Leaflet map for image tiles, including a few
helpful plugins like a legend (if compatible legend elements found in
TileJSON) and a transparency slider.  Vector tiles are previewed using
Mapbox GL.


## ArcGIS API
This project currently provides a minimal ArcGIS tiled map service API for tiles stored in an mbtiles file.
This should be sufficient for use with online platforms such as [Data Basin](https://databasin.org).  Because the ArcGIS API relies on a number of properties that are not commonly available within an mbtiles file, so certain aspects are stubbed out with minimal information.

This API is not intended for use with more full-featured ArcGIS applications such as ArcGIS Desktop.


## Live Examples
These are hosted on a free dyno by Heroku (thanks Heroku!), so there might be a small delay when you first access these.

* [List of services](http://frozen-island-41032.herokuapp.com/services)
* [TileJSON](http://frozen-island-41032.herokuapp.com/services/geography-class-png) for a PNG based tileset generated using TileMill.
* [Map Preview ](http://frozen-island-41032.herokuapp.com/services/geography-class-png/map) for a map preview of the above.
* [ArcGIS Map Service](http://frozen-island-41032.herokuapp.com/arcgis/rest/services/geography-class-png/MapServer)

## Roadmap
See the issues tagged to the [0.5 version](https://github.com/consbio/mbtileserver/milestone/1)
for our near term features and improvements.

In short, we are planning to:
* add tests and benchmarks
* get things production ready


## Development
Development of the templates and static assets likely requires using
`node` and `npm`.  Install these tools in the normal way.

From the `handlers/templates/static` folder, run
```
$npm install
```

to pull in the static dependencies.  These are referenced in the
`package.json` file.

Then to build the minified version, run:
```
$gulp build
```

Modifying the `.go` files always requires re-running `go build .`.

In case you have modified the templates and static assets, you need to run `go
generate ./handlers` to ensure that your modifications
are embedded into the executable. For this to work, you must have
[github.com/shurcooL/vfsgen)[https://github.com/shurcooL/vfsgen) installed.
This will rewrite the `assets_vfsdata.go` which you must commit along with your
modification. Also you should run `go build` after `go generate`.

During the development cycle you may use `go build -tags dev .` to build the
binary, in which case it will always take the assets from the relative file
path `handlers/templates/` directly and you can omit the `go generate` step.  (note: this is currently not working properly)
But do not forget to perform it in the end.





