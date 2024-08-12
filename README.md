# chirpy
A guided project where I created a minimum viable project version of X. We named it Chirpy because Lane is pretty funny and I like chirpin'. This was created as part of my boot.dev training. This repo contains a compiled binary that can be run on any server. The server will support many endpoints that are documented below. At long last, you can get a bit chirpy!

# Installation
Installation is super simple, just download the repo. It contains the compiled server binary.

## Usage
To launch the server, navigate to the repo on the command line and run ./chirpy
Once this is running, the base url is http://localhost:8080

## Endpoints
The server has the following endpoints:
1. http://localhost:8080/ - This will take you to a sample home page
2. http://localhost:8080/admin/metrics - This will show the number of hits over the lifetime of the server (Only accepts GET)
3. http://localhost:8080/api/reset - This will reset the hit counter for analytics
4. http://localhost:8080/api/users -
