// TODO
**Dockerize**
    go app
    mysql
    redis? not currently needed

**Componentize Frontend**
    react, vue? - only for components, not going full SPA
    vue
    npm install vue@latest webpack webpack-cli webpack-dev-server vue-loader@next --save-dev

**Single Point Rendering**
In order for the dev server to run while coding, I have to have air running, and npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch
And possibly soon one for running vue builds for javascript
How should I combine these into one thing that runs the dev server.
Prod is only relying on