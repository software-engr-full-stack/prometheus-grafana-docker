### Notes

1. `mkdir data && chmod 1777 data`

2. `docker compose up --build`

No data will be displayed on first run. Make sure to make a request to the end points so you'll have some data in `/metrics`. Don't forget to directly `curl` `/metrics` to see if Prometheus is actually scraping.
