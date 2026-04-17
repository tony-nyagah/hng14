# Stage 1B – Profile Card

A small, accessible, responsive Profile Card built with plain HTML, CSS, and JavaScript as part of the HNG14 Frontend track.

## Running locally

No build step required. Serve the directory over HTTP (opening `index.html` directly with `file://` can block CSS/JS loading in some browsers):

```bash
# Python (usually pre-installed)
python3 -m http.server 5174

# or Node / npx
npx serve .

# or Bun
bunx serve .
```

Then open **http://localhost:5174** in your browser.

## Features

| Feature | Detail |
|---|---|
| Live epoch time | Updates every 1 s via `setInterval`, shown in milliseconds |
| Avatar | External image URL; graceful alt text for screen readers |
| Social links | GitHub, Twitter/X, LinkedIn — open in new tab with `rel="noopener noreferrer"` |
| Hobbies & Dislikes | Two visually distinct lists (green / red accent) |
| Responsive layout | Avatar + text side-by-side on tablet+; stacked on mobile (≤ 480 px) |
| Accessibility | `aria-live="polite"` on the time element; visible focus rings; WCAG AA contrast |

## `data-testid` reference

| Element | `data-testid` |
|---|---|
| Card root | `test-profile-card` |
| Name | `test-user-name` |
| Bio | `test-user-bio` |
| Epoch time | `test-user-time` |
| Avatar | `test-user-avatar` |
| Social links container | `test-user-social-links` |
| GitHub link | `test-user-social-github` |
| Twitter/X link | `test-user-social-twitter` |
| LinkedIn link | `test-user-social-linkedin` |
| Hobbies list | `test-user-hobbies` |
| Dislikes list | `test-user-dislikes` |

## Design

Neobrutalist style — hard black borders, flat box shadows, bold uppercase labels, and yellow/green/red accent colours consistent with the rest of the HNG14 frontend work.
