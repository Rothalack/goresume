**Prod release**
pull master
npm run prod
go build
systemctl restart goresume

**commands**
npm run dev
concurrently "air" "npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --minify --watch --no-cache" "npx webpack --mode development --watch --no-cache"

npm run prod
npx webpack --mode production --no-cache && npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --minify

go run sync_base_data
sync warcraftlogs data tree

**smol todos**
Something on the homepage. Anything?
Link to github repo
Dark/Light mode?
Get all raid images. Write a scrapper? I don't know how I would do that because the images I've gotten so far have come from official blizzard press release material and it's consistent. Maybe get images from wowhead, which could be more orderly

**BIG TODOs**
**Mobile reactive layout**
Probably entirely tailwind classes

**logging/error logging**
Something to compile errors and logs into somewhere searchable. Prometheus+Grafana?

**Release System**
Simple and just have releases set up on github?
Go full, have release commands synced with github actions?
Migrations system? Just hard coded for documentation the migrations in /migrations for now.
Docker required to go full release hooked up with github actions, with migrations and all?

**Dockerize**
go app
mysql
redis? not currently needed
redis for storing access_token i think

compartmentalize into releasble framework with example pages
login/account system
backups?
maintanence mode
admin area with database tool? can be behind cloudflare auth etc
