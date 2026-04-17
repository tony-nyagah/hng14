# Stage 0 Frontend Task: Build a Testable Todo Item Card

| Field    | Details                              |
|----------|--------------------------------------|
| Deadline | 4/17/2026 1:59am                     |
| Task     | Build a Testable Todo Item Card      |

---

## Overview

Build a clean, high-fidelity Task Card (or a single-card page) that serves as the core UI element for a productivity app. It needs to feel more "alive" and interactive than a standard card, capturing the polished energy of an early-stage startup.

---

## Evaluation Criteria / Acceptance Criteria

Build the component using standard semantic elements, mapping the provided `testIds` to their respective roles for automated testing.

### Required `data-testid` Attributes

| Element            | `data-testid`                          |
|--------------------|----------------------------------------|
| Card container     | `test-todo-card`                       |
| Task title         | `test-todo-title`                      |
| Task description   | `test-todo-description`                |
| Priority badge     | `test-todo-priority`                   |
| Due date           | `test-todo-due-date`                   |
| Time remaining     | `test-todo-time-remaining`             |
| Status indicator   | `test-todo-status`                     |
| Checkbox           | `test-todo-complete-toggle`            |
| Categories list    | `test-todo-tags`                       |
| Each tag           | `test-todo-tag-{tag-name}` (optional)  |
| "Edit" button/icon | `test-todo-edit-button`                |
| "Delete" button/icon | `test-todo-delete-button`            |

### Notes

1. All due dates and hints must be formatted nicely (e.g. `"Due Feb 18, 2026"`, `"Due in 3 days"`, `"Overdue by 2 hours"` etc.) and must be accurate to the current time. Update the data reasonably (about once every 30–60 seconds).
2. The checkbox must be a real `<input type="checkbox">` or a properly labeled button with `role="checkbox"`.

---

## HTML & Semantics Recommendations

| Element             | Recommended Tag                                      |
|---------------------|------------------------------------------------------|
| Card root           | `<article>` or `<section role="region">`             |
| Title               | `<h2>` or `<h3>`                                     |
| Description         | `<p>`                                                |
| Priority & Status   | `<span>` or `<strong>` (add `aria-label` if visual-only) |
| Due date & time remaining | `<time>` element (with `datetime` attribute if possible) |
| Checkbox            | Real `<input type="checkbox">` + `<label>`           |
| Tags                | `<ul role="list">` of `<li>`, or `<div role="list">` with chips |
| Buttons             | `<button>` — not `<div>` (add `aria-label` if icon-only) |

---

## Submission Format

1. Live URL (Vercel / Netlify / GitHub Pages)
2. GitHub repo link with README including:
   - How to run locally
   - Decisions made
   - Any trade-offs

**Optional (Bonus Points):**
- Basic tests (React Testing Library / Vitest / Cypress)

**Submission Link:** https://docs.google.com/forms/d/e/1FAIpQLSd57saqCAZ34jRlAqD7y3gkopz7YZQ46-etzGmj1fZcVjEWwA/viewform
