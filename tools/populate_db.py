import requests
import sys

if len(sys.argv) != 3:
    print("E1: populate_db expects exactly two (2) inputs, received %d." % (len(sys.argv)-1))
    exit(1)

api = "https://api.github.com"
owner = sys.argv[0]
repos = sys.argv[1]

request_url = "https://api.github.com/repos/%s/%s/contributors" % (owner, repos)
r = requests.get(request_url)

if r.status_code != 200:
    print("E2: repository not found.")
    exit(2)

headers = r.headers
