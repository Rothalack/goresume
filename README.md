**Single Point Rendering**
npm run dev

this runs the following
air
npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch
npx webpack --mode development --watch --no-cache

**smol todos**

**BIG TODOs**
**logging/error logging**
Something to compile errors and logs into somewhere searchable. Google has a free option I think. Used Relix before.

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
