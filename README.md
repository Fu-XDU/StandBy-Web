<h1 align="center">
  <br>
  <a href="https://github.com/Fu-XDU/StandBy-Web"><img src="https://github.com/Fu-XDU/StandBy-Web/blob/main/src/assets/logo.png?raw=true" alt="standby-web" width="200"></a>
  <br>
Standby Web
  <br>
</h1>

<h4 align="center">A minimal Web Clock</a>.</h4>

<p align="center">
<a href="https://dl.circleci.com/status-badge/redirect/circleci/5Q5Wtj99uUSh3smo3BU9Fu/AkAG8p4SHPuRmyxAWrgwwN/tree/main"><img src="https://dl.circleci.com/status-badge/img/circleci/5Q5Wtj99uUSh3smo3BU9Fu/AkAG8p4SHPuRmyxAWrgwwN/tree/main.svg?style=svg&circle-token=CCIPRJ_2fmHaM97dUucEZvqYHdRzV_1f3362ba02bd572b98bf3248bf92f32b297c010e" alt="CircleCI"></a>
<a href="/CONTRIBUTING.md" title="Go to contributions doc"><img src="https://img.shields.io/badge/contributions-welcome-blue" alt="contributions - welcome"></a>
<a href="#license"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
</p>

## Screenshot

![Float](Screenshot/Float.png)

## See Also

[Fu-XDU/StandBy: A minimal Clock application for iPhone / iPad and Mac.](https://github.com/Fu-XDU/StandBy)

## Nginx (generic subpath)

The production build uses a relative asset base (`base: './'` in Vite), so you can host the app under any URL prefix. Replace **`YOUR_PREFIX`** in the snippet (e.g. `/standby`, `/clock`) and keep the on-disk layout aligned with the last `try_files` argument.

With **`root`**, put static files under **`<root>/YOUR_PREFIX/`** (example: `root` is `/var/www/html` and prefix is `/clock` → entry file `/var/www/html/clock/index.html`):

```nginx
location ^~ /YOUR_PREFIX {
    root /path/to/parent;
    index index.html index.htm;
    try_files $uri $uri/ /YOUR_PREFIX/index.html;

    location ~* \.(otf|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, max-age=31536000, immutable";
    }
}
```

If you **`proxy_pass`** a prefix to an upstream static host, relative assets still work as long as the browser’s URL matches that prefix; configure SPA fallback to `index.html` on the upstream or gateway as needed.

## Remote config API

The clock page registers a random `deviceId` (stored in browser `localStorage` as `standby_device_id`). Use that id to read or patch **only your own** config. The page syncs from the server on load and every minute.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/remote/register` | Register `deviceId` (done automatically by the page) |
| POST | `/v1/remote/sync` | Full config sync with `updatedAt` |
| GET | `/v1/remote/config?deviceId=&pageId=float` | Read server config |
| PATCH | `/v1/remote/config` | Update **only the fields you send** |

`pageId` for the Float clock is `float`.

### PATCH examples (partial update)

Only brightness:

```bash
curl -X PATCH 'http://127.0.0.1:1423/v1/remote/config' \
  -H 'Content-Type: application/json' \
  -d '{"deviceId":"YOUR_DEVICE_ID","pageId":"float","brightness":0.6}'
```

Color preset index (0–13):

```bash
curl -X PATCH 'http://127.0.0.1:1423/v1/remote/config' \
  -H 'Content-Type: application/json' \
  -d '{"deviceId":"YOUR_DEVICE_ID","pageId":"float","selectedColorIndex":5}'
```

Night mode schedule (ISO date-time strings, two elements):

```bash
curl -X PATCH 'http://127.0.0.1:1423/v1/remote/config' \
  -H 'Content-Type: application/json' \
  -d '{
    "deviceId":"YOUR_DEVICE_ID",
    "pageId":"float",
    "autoNightMode":true,
    "nightModeRange":["2025-01-01T22:00:00.000Z","2025-01-01T07:00:00.000Z"]
  }'
```

Other patchable fields: `autoInvisible`, `invisibleRange`, `invisibleDayEnable`, `invisibleDay` (7 booleans, Sun–Sat).

The response includes the merged full config and a new `updatedAt`; the browser applies it on the next sync.