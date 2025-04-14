# StandBy Web

[![CircleCI](https://dl.circleci.com/status-badge/img/circleci/5Q5Wtj99uUSh3smo3BU9Fu/AkAG8p4SHPuRmyxAWrgwwN/tree/main.svg?style=svg&circle-token=CCIPRJ_2fmHaM97dUucEZvqYHdRzV_1f3362ba02bd572b98bf3248bf92f32b297c010e)](https://dl.circleci.com/status-badge/redirect/circleci/5Q5Wtj99uUSh3smo3BU9Fu/AkAG8p4SHPuRmyxAWrgwwN/tree/main)

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

## Nginx Config

```nginx
    location ^~ /standby {
        root   /opt/homebrew/var/www/;
        index  index.html index.htm;
        try_files  $uri $uri/ /standby/index.html;
        
        location ~* \.(otf|woff2)$ {
            expires 1y;
            add_header Cache-Control "public, max-age=31536000, immutable";
        }
    }
```