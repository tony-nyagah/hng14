# Invoice Management Application

A full-stack, responsive Invoice Management Application built with React and Vanilla CSS. This project was developed to strictly adhere to Figma design specifications as part of the HNG14 Stage 2 frontend task.

## Features

* **Full CRUD Functionality:** Create, Read, Update, and Delete invoices.
* **Persistent State Management:** Invoice data and user preferences are saved locally using `localStorage`.
* **Theme Toggling:** Switch seamlessly between Light and Dark mode.
* **Form Validation:** Robust form handling and strict validation using `react-hook-form` and `zod`.
* **Dynamic Calculations:** Invoice totals update in real-time as item quantities or prices are adjusted.
* **Save as Draft:** Bypass strict validation to save partial invoices as drafts.
* **Filtering:** Filter invoices by their current status (Draft, Pending, Paid) via a custom dropdown.
* **Responsive Layout:** Optimized for Mobile, Tablet, and Desktop screens, featuring a dynamic slide-out form drawer.

## Setup Instructions

Ensure you have [Node.js](https://nodejs.org/) and [Bun](https://bun.sh/) installed on your machine.

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd inventory-app
   ```
2. **Install dependencies:**
   ```bash
   bun install
   ```
3. **Start the development server:**
   ```bash
   bun run dev
   ```
4. **Build for production:**
   ```bash
   bun run build
   ```

### Generating Fake Data (Optional)
During development, you may want to populate the application with mock data instead of starting from an empty state. A seeding script is included which utilizes `@faker-js/faker` to generate 10 schema-compliant, random invoices and saves them to `src/data.ts`.

To run the seeding script:
```bash
bun run scripts/seed.ts
```
**Note:** After running the script, you must clear your browser's `localStorage` and refresh the page to allow the app to initialize with the newly generated data.

## Architecture Explanation

The application follows a modular, client-side React architecture:
* **UI Layer:** Built with functional React components, utilizing standard React Router DOM for routing between the List view and Detail view.
* **State Management:** The Context API (`InvoiceContext`, `ThemeContext`) is used as a lightweight global store to manage invoices and the active theme, avoiding prop-drilling without adding the overhead of Redux or Zustand.
* **Data Layer:** `localStorage` acts as our persistent database, synchronized via Context. Forms use `react-hook-form` connected to a strictly typed `Zod` validation schema (`src/schema.ts`) to ensure data integrity before writing to state.
* **Styling Strategy:** Vanilla CSS is strictly used. A robust design system of CSS Variables (`index.css`) handles complex Light/Dark mode transitions across all components, intentionally avoiding reliance on utility frameworks like Tailwind or heavy component libraries.

## Trade-offs

* **Vanilla CSS vs. Tailwind/UI Libraries:** Choosing Vanilla CSS provided exact control over the complex, custom Figma designs (like the custom dropdowns and animations) and avoided external dependencies, but increased the sheer volume of CSS written and manually managed.
* **Context API vs. External State Management (Zustand/Redux):** Context API is perfect for the scale of this app. However, as the app grows, Context could lead to unnecessary re-renders across the tree. A trade-off was made favoring simplicity and zero-dependencies over atomic state optimization.
* **LocalStorage vs. Backend Database:** Using `localStorage` meets the requirements for a frontend-focused task and allows for immediate persistence, but it means data is local to the device and cannot be shared across sessions or users.

## Accessibility Notes

* **Semantic HTML:** The application uses semantic elements like `<header>`, `<main>`, `<section>`, and `<nav>` to ensure screen readers can accurately interpret the layout and landmarks.
* **Visual Contrast:** The color palettes exactly match the provided Figma specifications, which were designed with high-contrast ratios in mind for both Light and Dark modes.
* **Keyboard Navigation:** Forms are fully navigable via the `Tab` key, and inputs include visible focus states (`outline`, `border-color` transitions) to aid users relying on keyboards.
* **Responsive Touch Targets:** Buttons and interactive elements on mobile views are appropriately sized to prevent accidental mis-taps.

## Improvements Beyond Requirements

* **Automated Data Seeding:** Implemented a backend-style Node/Bun script utilizing `faker-js` to automatically generate realistic, schema-compliant dummy data instantly.
* **Reactive Form Totals:** Upgraded the New/Edit Invoice forms to dynamically calculate and render item totals and grand totals instantly upon typing, creating a much richer UX.
* **Polished Interactions:** Added a custom `click-outside` event listener to gracefully close the custom filtering dropdown, and implemented a smooth, hardware-accelerated slide-out animation for the side drawer across all viewports.
* **Refined Drawer Overlap Logic:** Calculated negative offset transformations on Desktop views to ensure the hidden side-drawer sits cleanly off-screen without interfering with the fixed sidebar.
