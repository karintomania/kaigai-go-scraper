# Structure
This folder contains endpoints for these app.

1. scraper.go: endpoint for scraper
  - store_links.go: get links for the page 
  - scrape_html.go: download HTML and scrape info
  - translate.go: translation via LLM
  - generate_article.go: generate the hugo template for a hugo article
  - publish.go: run shell script to publish the articles

2. tweet.go: endpoint for tweet automation
