# sv

Frontend project, drop in replacement for templ and htmx

- UI now can support RESTful API with JSON for document content
- Markdown now move to Svelte Preprocessor, which can support more component

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.

### Update dependancy

There might be some warning/error when deal with Svelte build, we should update regularly to prevent problem.

```sh
npm install @sveltejs/kit@latest
```
