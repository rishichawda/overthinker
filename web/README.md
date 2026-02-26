# Website Directory

This directory contains the source files for the overthink project website, which is deployed to the `gh-pages` branch.

## Files

- `index.html` - Main website page with project information, features, and installation methods
- `styles.css` - Stylesheet with the theatrical aesthetic matching the tool's tone
- `.nojekyll` - File that tells GitHub Pages to skip Jekyll processing

## Deployment

The website is automatically deployed to GitHub Pages via the `.github/workflows/deploy-website.yml` workflow whenever changes are pushed to the `web/` directory on the `main` branch.

**Important:** The workflow preserves existing files in the `gh-pages` branch (like `overthink-archive-keyring.gpg` used for APT package distribution), so the website deployment won't overwrite package distribution infrastructure.

## Local Development

To preview the website locally, simply open `web/index.html` in a web browser, or use a local web server:

```bash
# Using Python 3
python -m http.server --directory web

# Using Node.js
cd web && npx http-server
```

Then visit `http://localhost:8000` (or the port shown in your terminal).

## Styling

The website uses a theatrical color scheme matching the tool's dramatic personality:

- **Primary Color:** Purple (`#6b21a8`, `#9333ea`) for headers and accents
- **Backgrounds:** Light gradients with white cards for clean, readable content
- **Typography:** System fonts for fast loading, semantic sizing hierarchy

The design is responsive and works well on mobile, tablet, and desktop devices.
