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

## Technology Stack

* **Core:** React, TypeScript, Vite
* **Styling:** Vanilla CSS (Custom properties for theming, responsive media queries)
* **Routing:** React Router DOM
* **Forms & Validation:** React Hook Form, Zod
* **Icons:** Lucide React
* **Utilities:** date-fns (date formatting), uuid (unique IDs)
* **Development/Testing:** Faker.js (mock data generation)

## Getting Started

### Prerequisites

Ensure you have [Node.js](https://nodejs.org/) and [Bun](https://bun.sh/) installed.

### Installation

1. Clone the repository and navigate to the project directory:
   ```bash
   cd inventory-app
   ```
2. Install the dependencies:
   ```bash
   bun install
   ```
3. Start the development server:
   ```bash
   bun run dev
   ```

### Generating Fake Data

During development, you may want to populate the application with mock data instead of starting from an empty state.

A seeding script is included which utilizes `@faker-js/faker` to generate 10 schema-compliant, random invoices and saves them to `src/data.ts`.

To run the seeding script:

```bash
bun run scripts/seed.ts
```

**Note:** After running the script, you must clear your browser's `localStorage` and refresh the page to allow the app to initialize with the newly generated data.

## Project Structure

* `/src/components` - Reusable UI components (Buttons, Badges, Modals, Forms, Cards).
* `/src/context` - React Context providers (`ThemeContext`, `InvoiceContext`) for global state.
* `/src/pages` - Main application views (`InvoiceListPage`, `InvoiceDetailPage`).
* `/src/utils` - Helper functions and utility scripts.
* `/scripts` - Node scripts for development tasks (e.g., data seeding).
