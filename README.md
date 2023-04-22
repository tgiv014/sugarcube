# sugarcube

<img align="right" width="159px" src="https://raw.githubusercontent.com/tgiv014/sugarcube/main/web/static/sugarcube.png">

`sugarcube` is a *simple* and *sweet* dashboard for viewing blood glucose data in the browser.

# Features
- Realtime glucose readings through the Dexcom Share API
- Smart glucose reading fetches (It knows when to expect your next reading and starts asking for a reading then)
- Dashboard updates instantly whena new glucose reading becomes available (Server-Sent Events!)

<img src="https://raw.githubusercontent.com/tgiv014/sugarcube/main/screenshot.png">

# Setup

`⚒️ TODO`

# Development

`sugarcube` uses a [Go](https://go.dev/) backend and [SvelteKit](https://kit.svelte.dev/) frontend.

The production build embeds the frontend in it so all you need is the binary. The dev build starts up the server without the embedded frontend, starts [vite](https://vitejs.dev/) and proxies unhandled requests to it. This enables hot reload for frontend development with one command and no cross origin issues.

1. [Install `node` and `npm`](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
1. Run ```make dev```
