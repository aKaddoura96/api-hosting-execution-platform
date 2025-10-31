# Frontend - API Platform

Next.js-based marketplace and dashboard for the API hosting platform.

## Features

- Public API marketplace
- Developer dashboard (API management, analytics, earnings)
- Consumer dashboard (usage, billing)
- API documentation viewer
- Payment integration

## Getting Started

```bash
# Install dependencies
npm install

# Run development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the app.

## Project Structure

```
app/
├── (auth)/           # Authentication pages
├── (marketplace)/    # Public API catalog
├── dashboard/        # Developer dashboard
│   ├── apis/        # API management
│   ├── analytics/   # Usage metrics
│   └── earnings/    # Revenue tracking
└── consumer/         # Consumer dashboard
```

## TODO

- [ ] Implement authentication UI
- [ ] Create API marketplace browse/search
- [ ] Build developer dashboard
- [ ] Add consumer usage dashboard
- [ ] Integrate payment UI (Stripe/PayTabs)
- [ ] Add API documentation viewer
