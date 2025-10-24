# Stori Financial Tracker - Frontend

React + TypeScript web application for visualizing financial data with AI-powered insights.

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 📋 Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Build tool & dev server
- **TailwindCSS** - Utility-first styling
- **Recharts** - Chart library for visualizations
- **Axios** - HTTP client for API calls

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── Dashboard/      # Financial overview components
│   │   ├── Charts/         # Chart components (category, timeline)
│   │   ├── AI/            # AI advisor components
│   │   └── common/        # Reusable UI components
│   ├── services/
│   │   └── api.ts         # API client & service methods with types
│   ├── hooks/             # Custom React hooks
│   ├── utils/
│   │   └── formatters.ts  # Currency, date, number formatting
│   ├── App.tsx            # Main app component
│   ├── main.tsx           # Entry point
│   └── vite-env.d.ts      # Vite environment types
├── public/                # Static assets
└── index.html            # HTML template
```

## 🎯 Features

- 📊 Financial overview dashboard
- 📈 Category spending breakdown (pie/donut chart)
- 📉 Income vs expense timeline (line/area chart)
- 🤖 AI-powered financial advice
- 📱 Mobile-responsive design
- 🎨 Modern UI with TailwindCSS

## ⚙️ Configuration

### Environment Variables

Create a `.env` file (copy from `env.example`):

```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api
```

### API Endpoints

The frontend connects to the backend API:

- `GET /api/health` - Health check
- `GET /api/transactions` - Get transactions
- `GET /api/summary/categories` - Category breakdown
- `GET /api/summary/timeline` - Timeline data
- `POST /api/advice` - AI financial advice

## 🔧 Development

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Backend API running on port 8080

### Available Scripts

```bash
npm run dev        # Start dev server (http://localhost:5173)
npm run build      # Build for production
npm run preview    # Preview production build
npm run lint       # Run ESLint
```

### Development Server

The dev server runs on `http://localhost:5173` with:
- Hot module replacement (HMR)
- Fast refresh
- Proxy to backend API (if configured)

## 🎨 Styling

TailwindCSS utility classes with custom configuration:

```javascript
// Custom color palette
primary: {
  50: '#f0f9ff',
  // ... more shades
  900: '#0c4a6e',
}
```

### Custom Components

Predefined component classes in `index.css`:

- `.card` - Card container
- `.btn-primary` - Primary button
- `.btn-secondary` - Secondary button

## 📊 Chart Components

Powered by Recharts with responsive design:

- **CategoryChart** - Pie/Donut chart for spending breakdown
- **TimelineChart** - Line/Area chart for income vs expenses

## 🏗️ Architecture

### Component Hierarchy

```
App
├── Header (nav + API status)
├── Dashboard
│   ├── FinancialOverview (summary cards)
│   ├── CategoryChart (spending breakdown)
│   └── TimelineChart (monthly trends)
├── AIAdvisor (financial advice)
└── Footer
```

### Data Flow

1. Components call API service methods
2. API service uses axios to fetch data
3. Data formatted using utility functions
4. Components render with formatted data

## 🐳 Docker

```bash
# Build image
docker build -t stori-frontend .

# Run container
docker run -p 5173:5173 stori-frontend
```

## 🚀 Deployment

### Build for Production

```bash
npm run build
```

Output in `dist/` directory:
- Optimized bundle
- Minified CSS
- Static assets

### Deploy to AWS S3

```bash
# Build
npm run build

# Sync to S3
aws s3 sync dist/ s3://your-bucket-name

# Invalidate CloudFront (if using)
aws cloudfront create-invalidation \
  --distribution-id YOUR_DIST_ID \
  --paths "/*"
```

## 📱 Mobile Responsive

Tailwind breakpoints:
- `sm:` - 640px and up
- `md:` - 768px and up
- `lg:` - 1024px and up
- `xl:` - 1280px and up

## 🧪 Testing

```bash
# Run tests (when implemented)
npm test

# Run with coverage
npm test -- --coverage
```

## 🔗 Related

- Backend API: `../backend/`
- Documentation: `../docs/`
- Deployment: See Day 3 deployment guide

## 📚 Learn More

- [React Documentation](https://react.dev)
- [Vite Documentation](https://vitejs.dev)
- [TailwindCSS Documentation](https://tailwindcss.com)
- [Recharts Documentation](https://recharts.org)
