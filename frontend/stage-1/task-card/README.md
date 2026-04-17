# Task Card

A neubrutalist todo item card built with Vue 3, Tailwind CSS v4, and DaisyUI v5. Built for HNG14 Stage 0 (Frontend).

## How to run locally

**Prerequisites:** [Bun](https://bun.sh) installed.

```sh
bun install
bun dev
```

Open [http://localhost:5173](http://localhost:5173).

To preview the production build:

```sh
bun run build
bun run preview
```

## Decisions made

- **Vue 3 + `<script setup>`** — composition API keeps reactive timer logic and computed display values co-located and easy to follow.
- **Tailwind CSS v4 via `@tailwindcss/vite`** — no config file required; single `@import "tailwindcss"` in `style.css`.
- **DaisyUI v5** — loaded as a Tailwind plugin for its reset and theme tokens. Pinned `data-theme="light"` on `<html>` to prevent the default dark/cyan theme from overriding colors.
- **Neubrutalist style** — hard offset box-shadows, thick solid black borders, zero border-radius, and high-contrast accent colors (`#ff3b3b`, `#ffe500`, `#b2ff59`) implemented as scoped CSS classes to avoid fighting Tailwind's utility layer.
- **Live time display** — a `setInterval` running every 30 seconds updates a `now` ref; `timeRemaining` and `timeColor` are computed values derived from it. The interval is cleaned up in `onUnmounted`.
- **All `data-testid` attributes** on the correct semantic elements (`<article>`, `<h2>`, `<p>`, `<time>`, `<input type="checkbox">`, `<ul>`, `<button>`) for automated testing compatibility.

## Trade-offs

- **Single component** — all markup, logic, and scoped styles live in `App.vue`. Fine for a single card; would need splitting into a proper `TodoCard.vue` component if the project grew.
- **Hardcoded data** — task content and due date are static. A real app would accept props or fetch from an API.
- **30-second interval** — the spec calls for updates "about once every 30–60 seconds". 30s was chosen for responsiveness; it has negligible performance impact but means the tab holds an active timer while open.
- **`vite preview` in Docker** — the container runs `vite preview` instead of a dedicated static file server. This is fine for a demo/staging deployment but `vite preview` is not intended for production at scale.

