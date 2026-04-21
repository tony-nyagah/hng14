# Stage 2 Frontend Task: Build an Invoice Management App

| Field    | Details                              |
|----------|--------------------------------------|
| Deadline | 4/23/2026 2:00am                     |
| Task     | Build an Invoice Management App      |

---

## Overview

Build a responsive, full-stack Invoice Management Application based on the provided Figma design.

**[FIGMA DESIGN LINK](https://www.figma.com/design/e3MtRefbZw41Ts897CQF4N/invoice-app?node-id=0-1&m=dev&t=pJoJoOU92dYwiC5p-1)**


---

## Core Objective

Build a fully functional invoice app that allows users to:

- Create invoices
- Read (view) invoices
- Update invoices
- Delete invoices
- Save drafts
- Mark invoices as paid
- Filter by invoice status
- Toggle light/dark mode
- Experience full responsiveness
- See hover states on interactive elements
- Persist state using: LocalStorage, IndexedDB, or a backend (Node/Express, Next.js API, etc.)

---

## Core Features

### 1️⃣ CRUD

| Action | Steps |
|--------|-------|
| Create | Open invoice form → fill required fields → save |
| Read   | View invoice list → click to view full details |
| Update | Edit existing invoice → persist updated values |
| Delete | Delete invoice → show confirmation modal first |

### 2️⃣ Form Validation

When creating or editing invoices:

- Required fields must be validated
- Invalid fields should show an error message, have a visual error state, and prevent submission
- Example validations:
  - Client name required
  - Valid email format
  - At least one invoice item
  - Quantity and price must be positive numbers

### 3️⃣ Draft & Payment Flow

Invoices have one of three statuses: **Draft**, **Pending**, **Paid**

- Users can save an invoice as Draft
- Draft invoices can be edited later
- Pending invoices can be marked as Paid
- Paid invoices **cannot** be reverted to Draft
- Status must be clearly reflected in the list view, detail view, and status badge color/style

### 4️⃣ Filter by Status

- Filter by: All, Draft, Pending, Paid
- Filter control should be intuitive (dropdown or checkbox)
- Filtered list updates immediately
- Empty state displays when no invoices match the filter

### 5️⃣ Light & Dark Mode Toggle

- Theme applies globally across all components
- Preference persisted in LocalStorage
- Good color contrast in both modes (WCAG AA)

### 6️⃣ Responsive Design

| Breakpoint | Width   |
|------------|---------|
| Mobile     | 320px+  |
| Tablet     | 768px+  |
| Desktop    | 1024px+ |

- Invoice list adapts to screen size
- Forms are usable on mobile
- No horizontal overflow
- Proper spacing and visual hierarchy

### 7️⃣ Hover & Interactive States

All interactive elements must have visible hover states: buttons, links, invoice list items, status filters, and form inputs.

---

## Recommended Architecture

Use **React only**. Suggested component structure:

- `InvoiceListPage`
- `InvoiceDetailPage` / Drawer
- `InvoiceForm` component
- `StatusBadge` component
- `Filter` component
- `ThemeProvider` / Context

---

## Accessibility Expectations

- Proper semantic HTML
- Form fields with `<label>`
- Buttons must be `<button>`
- Modal must: trap focus, close via ESC key, and be keyboard navigable
- Good color contrast (WCAG AA)

---

## Evaluation Criteria / Acceptance Criteria

- CRUD functionality works end-to-end
- Form validation prevents invalid submissions
- Status logic (Draft → Pending → Paid) behaves correctly
- Filtering works accurately with empty state
- Theme toggle persists across reload
- Fully responsive layout (mobile → desktop)
- Clean component structure
- No console errors
- Good accessibility practices

---

## Submission Format

1. Live URL (Vercel / Netlify / etc.)
2. GitHub repository
3. README including:
   - Setup instructions
   - Architecture explanation
   - Trade-offs
   - Accessibility notes
   - Any improvements beyond requirements

**Submission Link:** https://docs.google.com/forms/d/e/1FAIpQLScwBeDPp672SP-hUQ3PTN2NFAUbvaWmHdeDbie560iJMDQC_w/viewform

**Points:** 100 | **Pass Mark:** 70
