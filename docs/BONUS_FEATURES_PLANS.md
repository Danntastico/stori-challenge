# Bonus Features Implementation Plans

> Detailed step-by-step plans for optional enhancements

---

## Table of Contents

1. [Dark Mode Toggle](#1-dark-mode-toggle)
2. [Loading Skeleton States](#2-loading-skeleton-states)
3. [Date Range Filtering](#3-date-range-filtering)
4. [Export CSV Functionality](#4-export-csv-functionality)
5. [GitHub Actions CI/CD](#5-github-actions-cicd)
6. [Terraform IaC Scripts](#6-terraform-iac-scripts)

---

## 1. Dark Mode Toggle

**⏱️ Estimated Time:** 45 minutes  
**💡 Impact:** High (visual appeal, modern UX)  
**🎯 Complexity:** Low

### Overview
Add a toggle button to switch between light and dark color schemes. Persist user preference in localStorage.

### Implementation Steps

#### Step 1: Add Dark Mode Context (10 min)

**Create:** `frontend/src/contexts/ThemeContext.tsx`

```typescript
import React, { createContext, useContext, useState, useEffect } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem('theme');
    return (saved as Theme) || 'light';
  });

  useEffect(() => {
    localStorage.setItem('theme', theme);
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) throw new Error('useTheme must be used within ThemeProvider');
  return context;
};
```

#### Step 2: Wrap App with Theme Provider (5 min)

**Modify:** `frontend/src/main.tsx`

```typescript
import { ThemeProvider } from './contexts/ThemeContext';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ThemeProvider>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
```

#### Step 3: Configure Tailwind for Dark Mode (5 min)

**Modify:** `frontend/tailwind.config.js`

```javascript
export default {
  darkMode: 'class', // Enable class-based dark mode
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [],
};
```

#### Step 4: Create Theme Toggle Component (10 min)

**Create:** `frontend/src/components/common/ThemeToggle.tsx`

```typescript
import { useTheme } from '../../contexts/ThemeContext';
import { SunIcon, MoonIcon } from '@heroicons/react/24/outline'; // or create simple SVGs

export const ThemeToggle: React.FC = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      className="p-2 rounded-lg bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      aria-label="Toggle theme"
    >
      {theme === 'light' ? (
        <MoonIcon className="w-5 h-5 text-gray-800 dark:text-gray-200" />
      ) : (
        <SunIcon className="w-5 h-5 text-gray-800 dark:text-gray-200" />
      )}
    </button>
  );
};
```

#### Step 5: Add Dark Mode Styles (15 min)

**Modify:** `frontend/src/index.css`

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  body {
    @apply bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors;
  }
}
```

**Modify components to support dark mode:**
- Add `dark:` variants to all background, text, and border colors
- Example: `bg-white dark:bg-gray-800`
- Update: Card, Button, LoadingSkeleton, ErrorAlert, etc.

**Example Updates:**
```typescript
// Card.tsx
className="bg-white dark:bg-gray-800 shadow-md dark:shadow-gray-900/50"

// Button.tsx
className="bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600"
```

#### Step 6: Add Toggle to Header (5 min)

**Modify:** `frontend/src/App.tsx`

```typescript
import { ThemeToggle } from './components/common/ThemeToggle';

// Add to top of dashboard
<div className="flex justify-between items-center mb-6">
  <h1 className="text-3xl font-bold">Financial Dashboard</h1>
  <ThemeToggle />
</div>
```

### Files to Create/Modify

**Create:**
- `frontend/src/contexts/ThemeContext.tsx`
- `frontend/src/components/common/ThemeToggle.tsx`

**Modify:**
- `frontend/src/main.tsx`
- `frontend/tailwind.config.js`
- `frontend/src/index.css`
- `frontend/src/App.tsx`
- All component files (add dark: variants)

### Testing Checklist

- [ ] Toggle switches between light and dark modes
- [ ] Preference persists after page reload
- [ ] All components are readable in both modes
- [ ] Charts render correctly in dark mode
- [ ] No color contrast issues (WCAG compliance)

---

## 2. Loading Skeleton States

**⏱️ Estimated Time:** 30 minutes  
**💡 Impact:** Medium (better perceived performance)  
**🎯 Complexity:** Low

### Overview
Replace generic loading spinners with skeleton screens that mimic the final content layout, providing better UX feedback.

### Implementation Steps

#### Step 1: Create Skeleton Components (15 min)

**Create:** `frontend/src/components/common/Skeletons.tsx`

```typescript
export const CardSkeleton: React.FC = () => (
  <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md animate-pulse">
    <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-1/2 mb-4"></div>
    <div className="h-8 bg-gray-300 dark:bg-gray-700 rounded w-3/4 mb-2"></div>
    <div className="h-3 bg-gray-300 dark:bg-gray-700 rounded w-1/3"></div>
  </div>
);

export const ChartSkeleton: React.FC = () => (
  <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md animate-pulse">
    <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-1/3 mb-4"></div>
    <div className="flex space-x-2 mb-4">
      <div className="h-8 bg-gray-300 dark:bg-gray-700 rounded w-20"></div>
      <div className="h-8 bg-gray-300 dark:bg-gray-700 rounded w-20"></div>
    </div>
    <div className="h-64 bg-gray-300 dark:bg-gray-700 rounded"></div>
  </div>
);

export const OverviewSkeleton: React.FC = () => (
  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
    <CardSkeleton />
    <CardSkeleton />
    <CardSkeleton />
  </div>
);

export const AIAdvisorSkeleton: React.FC = () => (
  <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md animate-pulse">
    <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-1/2 mb-4"></div>
    <div className="space-y-2">
      <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full"></div>
      <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-5/6"></div>
      <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-4/5"></div>
    </div>
  </div>
);
```

#### Step 2: Update Components to Use Skeletons (15 min)

**Modify:** `frontend/src/components/Dashboard/FinancialOverview.tsx`

```typescript
import { OverviewSkeleton } from '../common/Skeletons';

// Replace LoadingSkeleton with:
if (loading) return <OverviewSkeleton />;
```

**Modify:** `frontend/src/components/Charts/CategoryChart.tsx`

```typescript
import { ChartSkeleton } from '../common/Skeletons';

if (loading) return <ChartSkeleton />;
```

**Modify:** `frontend/src/components/Charts/TimelineChart.tsx`

```typescript
import { ChartSkeleton } from '../common/Skeletons';

if (loading) return <ChartSkeleton />;
```

**Modify:** `frontend/src/components/AI/AIAdvisor.tsx`

```typescript
import { AIAdvisorSkeleton } from '../common/Skeletons';

if (loading) return <AIAdvisorSkeleton />;
```

### Files to Create/Modify

**Create:**
- `frontend/src/components/common/Skeletons.tsx`

**Modify:**
- `frontend/src/components/Dashboard/FinancialOverview.tsx`
- `frontend/src/components/Charts/CategoryChart.tsx`
- `frontend/src/components/Charts/TimelineChart.tsx`
- `frontend/src/components/AI/AIAdvisor.tsx`

### Testing Checklist

- [ ] Skeletons match final content layout
- [ ] Animations are smooth (pulse effect)
- [ ] Works in both light and dark modes
- [ ] No layout shift when content loads

---

## 3. Date Range Filtering

**⏱️ Estimated Time:** 1.5 hours  
**💡 Impact:** High (user control, deeper insights)  
**🎯 Complexity:** Medium

### Overview
Allow users to filter transactions by date range. Add preset ranges (Last 30 days, Last 3 months, etc.) and custom date picker.

### Implementation Steps

#### Step 1: Backend - Add Date Filtering Support (20 min)

**Modify:** `backend/internal/handlers/summary_handler.go`

```go
func (h *SummaryHandler) GetCategorySummary(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    startDateStr := r.URL.Query().Get("startDate")
    endDateStr := r.URL.Query().Get("endDate")
    
    var startDate, endDate time.Time
    var err error
    
    if startDateStr != "" {
        startDate, err = time.Parse("2006-01-02", startDateStr)
        if err != nil {
            http.Error(w, "Invalid start date format", http.StatusBadRequest)
            return
        }
    }
    
    if endDateStr != "" {
        endDate, err = time.Parse("2006-01-02", endDateStr)
        if err != nil {
            http.Error(w, "Invalid end date format", http.StatusBadRequest)
            return
        }
    }
    
    // Get filtered summary
    summary, err := h.service.GetCategorySummary(startDate, endDate)
    // ... rest of handler
}
```

**Modify:** `backend/internal/service/analytics_service.go`

Update `GetCategorySummary` to accept optional date range parameters.

#### Step 2: Frontend - Create Date Range Picker (25 min)

**Create:** `frontend/src/components/common/DateRangePicker.tsx`

```typescript
import { useState } from 'react';

interface DateRange {
  startDate: string;
  endDate: string;
}

interface Props {
  onRangeChange: (range: DateRange) => void;
}

export const DateRangePicker: React.FC<Props> = ({ onRangeChange }) => {
  const [preset, setPreset] = useState<string>('all');
  
  const presets = [
    { id: 'all', label: 'All Time' },
    { id: '30d', label: 'Last 30 Days' },
    { id: '3m', label: 'Last 3 Months' },
    { id: '6m', label: 'Last 6 Months' },
    { id: 'ytd', label: 'Year to Date' },
    { id: 'custom', label: 'Custom' },
  ];
  
  const handlePresetChange = (presetId: string) => {
    setPreset(presetId);
    const today = new Date();
    let startDate = '';
    
    switch (presetId) {
      case '30d':
        startDate = new Date(today.setDate(today.getDate() - 30)).toISOString().split('T')[0];
        break;
      case '3m':
        startDate = new Date(today.setMonth(today.getMonth() - 3)).toISOString().split('T')[0];
        break;
      case '6m':
        startDate = new Date(today.setMonth(today.getMonth() - 6)).toISOString().split('T')[0];
        break;
      case 'ytd':
        startDate = new Date(today.getFullYear(), 0, 1).toISOString().split('T')[0];
        break;
      case 'all':
      default:
        startDate = '';
    }
    
    onRangeChange({ startDate, endDate: '' });
  };
  
  return (
    <div className="flex flex-wrap gap-2 mb-6">
      {presets.map((p) => (
        <button
          key={p.id}
          onClick={() => handlePresetChange(p.id)}
          className={`px-4 py-2 rounded-lg transition-colors ${
            preset === p.id
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
          }`}
        >
          {p.label}
        </button>
      ))}
    </div>
  );
};
```

#### Step 3: Update API Service (15 min)

**Modify:** `frontend/src/services/api.ts`

```typescript
export interface DateRangeParams {
  startDate?: string;
  endDate?: string;
}

export const getCategorySummary = async (dateRange?: DateRangeParams): Promise<CategorySummaryResponse> => {
  const params = new URLSearchParams();
  if (dateRange?.startDate) params.append('startDate', dateRange.startDate);
  if (dateRange?.endDate) params.append('endDate', dateRange.endDate);
  
  const response = await apiClient.get<CategorySummaryResponse>(
    `/summary/categories?${params.toString()}`
  );
  return response.data;
};

// Update other endpoints similarly
export const getTimelineSummary = async (dateRange?: DateRangeParams) => { /* ... */ };
```

#### Step 4: Add Date Range State to App (20 min)

**Modify:** `frontend/src/App.tsx`

```typescript
import { DateRangePicker } from './components/common/DateRangePicker';
import { useState } from 'react';

function App() {
  const [dateRange, setDateRange] = useState<{ startDate: string; endDate: string }>({
    startDate: '',
    endDate: '',
  });
  
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="container mx-auto px-4 py-8">
        <DateRangePicker onRangeChange={setDateRange} />
        
        <FinancialOverview dateRange={dateRange} />
        <CategoryChart dateRange={dateRange} />
        <TimelineChart dateRange={dateRange} />
        <AIAdvisor dateRange={dateRange} />
      </div>
    </div>
  );
}
```

#### Step 5: Update Components to Use Date Range (10 min)

**Modify all dashboard components** to accept and use `dateRange` prop:

```typescript
interface Props {
  dateRange?: { startDate: string; endDate: string };
}

export const FinancialOverview: React.FC<Props> = ({ dateRange }) => {
  const { data, loading, error } = useApiData(
    () => api.getCategorySummary(dateRange),
    [dateRange] // Re-fetch when date range changes
  );
  // ...
};
```

### Files to Create/Modify

**Backend:**
- `backend/internal/handlers/summary_handler.go`
- `backend/internal/handlers/transaction_handler.go`
- `backend/internal/service/analytics_service.go`

**Frontend:**
- `frontend/src/components/common/DateRangePicker.tsx` (create)
- `frontend/src/services/api.ts`
- `frontend/src/App.tsx`
- `frontend/src/components/Dashboard/FinancialOverview.tsx`
- `frontend/src/components/Charts/CategoryChart.tsx`
- `frontend/src/components/Charts/TimelineChart.tsx`
- `frontend/src/components/AI/AIAdvisor.tsx`

### Testing Checklist

- [ ] Preset filters work correctly
- [ ] API receives correct date parameters
- [ ] Charts update with filtered data
- [ ] "All Time" shows original data
- [ ] No errors with invalid date ranges
- [ ] URL persistence (optional: use query params)

---

## 4. Export CSV Functionality

**⏱️ Estimated Time:** 45 minutes  
**💡 Impact:** Medium (useful for users)  
**🎯 Complexity:** Low

### Overview
Allow users to export transaction data and summaries as CSV files for use in Excel or other tools.

### Implementation Steps

#### Step 1: Backend - Add CSV Export Endpoint (20 min)

**Create:** `backend/internal/handlers/export_handler.go`

```go
package handlers

import (
    "encoding/csv"
    "fmt"
    "net/http"
    "github.com/danntastico/stori-backend/internal/repository"
)

type ExportHandler struct {
    repo repository.TransactionRepository
}

func NewExportHandler(repo repository.TransactionRepository) *ExportHandler {
    return &ExportHandler{repo: repo}
}

func (h *ExportHandler) ExportTransactionsCSV(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.repo.GetAll()
    if err != nil {
        http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
        return
    }
    
    // Set headers for CSV download
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")
    
    writer := csv.NewWriter(w)
    defer writer.Flush()
    
    // Write header
    writer.Write([]string{"Date", "Amount", "Category", "Description", "Type"})
    
    // Write data
    for _, tx := range transactions {
        writer.Write([]string{
            tx.Date,
            fmt.Sprintf("%.2f", tx.Amount),
            tx.Category,
            tx.Description,
            tx.Type,
        })
    }
}
```

**Modify:** `backend/main.go`

```go
// Register export routes
exportHandler := handlers.NewExportHandler(repo)
r.Get("/api/export/transactions", exportHandler.ExportTransactionsCSV)
```

#### Step 2: Frontend - Create Export Button (15 min)

**Create:** `frontend/src/components/common/ExportButton.tsx`

```typescript
import { ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { api } from '../../services/api';

export const ExportButton: React.FC = () => {
  const handleExport = async () => {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/export/transactions`);
      const blob = await response.blob();
      
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `stori-transactions-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Failed to export transactions');
    }
  };
  
  return (
    <button
      onClick={handleExport}
      className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors"
    >
      <ArrowDownTrayIcon className="w-5 h-5" />
      Export CSV
    </button>
  );
};
```

#### Step 3: Add Export Button to UI (10 min)

**Modify:** `frontend/src/App.tsx`

```typescript
import { ExportButton } from './components/common/ExportButton';

// Add to header section
<div className="flex justify-between items-center mb-6">
  <h1 className="text-3xl font-bold">Financial Dashboard</h1>
  <div className="flex gap-4">
    <ExportButton />
    <ThemeToggle />
  </div>
</div>
```

### Files to Create/Modify

**Backend:**
- `backend/internal/handlers/export_handler.go` (create)
- `backend/main.go`

**Frontend:**
- `frontend/src/components/common/ExportButton.tsx` (create)
- `frontend/src/App.tsx`

### Testing Checklist

- [ ] CSV file downloads successfully
- [ ] File contains all transactions
- [ ] Headers are correct
- [ ] Data is properly formatted
- [ ] Works in different browsers
- [ ] Filename includes date

---

## 5. GitHub Actions CI/CD

**⏱️ Estimated Time:** 1 hour  
**💡 Impact:** High (automation, professionalism)  
**🎯 Complexity:** Medium

### Overview
Automate testing, building, and deployment using GitHub Actions. Run tests on every push, deploy on merge to main.

### Implementation Steps

#### Step 1: Backend Test Workflow (15 min)

**Create:** `.github/workflows/backend-test.yml`

```yaml
name: Backend Tests

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'backend/**'
  pull_request:
    branches: [ main ]
    paths:
      - 'backend/**'

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Install dependencies
        working-directory: ./backend
        run: go mod download
      
      - name: Run tests
        working-directory: ./backend
        run: go test ./... -v -cover
      
      - name: Run race detector
        working-directory: ./backend
        run: go test ./... -race
      
      - name: Check formatting
        working-directory: ./backend
        run: |
          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
            echo "Go files must be formatted with gofmt"
            gofmt -s -l .
            exit 1
          fi
```

#### Step 2: Frontend Test Workflow (15 min)

**Create:** `.github/workflows/frontend-test.yml`

```yaml
name: Frontend Tests

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'frontend/**'
  pull_request:
    branches: [ main ]
    paths:
      - 'frontend/**'

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: './frontend/package-lock.json'
      
      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Type check
        working-directory: ./frontend
        run: npm run type-check || npx tsc --noEmit
      
      - name: Lint
        working-directory: ./frontend
        run: npm run lint || npx eslint .
      
      - name: Build
        working-directory: ./frontend
        run: npm run build
```

#### Step 3: Deployment Workflow (30 min)

**Create:** `.github/workflows/deploy.yml`

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Build backend
        working-directory: ./backend
        run: |
          GOOS=linux GOARCH=amd64 go build -o stori-backend main.go
      
      - name: Deploy to EC2
        env:
          SSH_PRIVATE_KEY: ${{ secrets.EC2_SSH_KEY }}
          EC2_HOST: ${{ secrets.EC2_HOST }}
        run: |
          echo "$SSH_PRIVATE_KEY" > private_key.pem
          chmod 600 private_key.pem
          scp -i private_key.pem -o StrictHostKeyChecking=no \
            backend/stori-backend ec2-user@$EC2_HOST:/home/ec2-user/
          ssh -i private_key.pem -o StrictHostKeyChecking=no \
            ec2-user@$EC2_HOST 'sudo systemctl restart stori-backend'
          rm private_key.pem
  
  deploy-frontend:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'
      
      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Build frontend
        working-directory: ./frontend
        env:
          VITE_API_BASE_URL: https://stori.danntastico.dev/api
        run: npm run build
      
      - name: Deploy to EC2
        env:
          SSH_PRIVATE_KEY: ${{ secrets.EC2_SSH_KEY }}
          EC2_HOST: ${{ secrets.EC2_HOST }}
        run: |
          echo "$SSH_PRIVATE_KEY" > private_key.pem
          chmod 600 private_key.pem
          scp -i private_key.pem -o StrictHostKeyChecking=no -r \
            frontend/dist/* ec2-user@$EC2_HOST:/var/www/stori/
          rm private_key.pem
```

### GitHub Secrets to Configure

Add these in GitHub repo → Settings → Secrets:
- `EC2_SSH_KEY`: Your private SSH key content
- `EC2_HOST`: EC2 public IP (107.22.66.194)

### Files to Create

**Create:**
- `.github/workflows/backend-test.yml`
- `.github/workflows/frontend-test.yml`
- `.github/workflows/deploy.yml`

### Testing Checklist

- [ ] Push to branch triggers tests
- [ ] Tests pass/fail correctly
- [ ] Merge to main triggers deployment
- [ ] Backend deploys and restarts
- [ ] Frontend deploys successfully
- [ ] GitHub Actions badge works

---

## 6. Terraform IaC Scripts

**⏱️ Estimated Time:** 1.5 hours  
**💡 Impact:** High (reproducibility, best practice)  
**🎯 Complexity:** High

### Overview
Infrastructure as Code using Terraform to provision all AWS resources (EC2, Security Groups, etc.), making the deployment reproducible.

### Implementation Steps

#### Step 1: Setup Terraform Structure (10 min)

**Create directory structure:**
```
terraform/
├── main.tf
├── variables.tf
├── outputs.tf
├── provider.tf
└── modules/
    ├── ec2/
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    └── security-group/
        ├── main.tf
        ├── variables.tf
        └── outputs.tf
```

#### Step 2: Provider Configuration (10 min)

**Create:** `terraform/provider.tf`

```hcl
terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  backend "s3" {
    bucket = "stori-terraform-state"
    key    = "stori-challenge/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = "stori-challenge"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
```

#### Step 3: Variables Definition (15 min)

**Create:** `terraform/variables.tf`

```hcl
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "key_name" {
  description = "EC2 SSH key pair name"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
  default     = "stori.danntastico.dev"
}

variable "allowed_ssh_cidr" {
  description = "CIDR block allowed to SSH"
  type        = list(string)
  default     = ["0.0.0.0/0"] # Restrict in production!
}
```

#### Step 4: Security Group Module (20 min)

**Create:** `terraform/modules/security-group/main.tf`

```hcl
resource "aws_security_group" "main" {
  name        = "${var.name_prefix}-sg"
  description = "Security group for ${var.name_prefix}"
  
  # SSH access
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.allowed_ssh_cidr
  }
  
  # HTTP
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  # HTTPS
  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  # Egress - allow all outbound
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = {
    Name = "${var.name_prefix}-security-group"
  }
}
```

**Create:** `terraform/modules/security-group/variables.tf`

```hcl
variable "name_prefix" {
  description = "Prefix for resource names"
  type        = string
}

variable "allowed_ssh_cidr" {
  description = "CIDR blocks allowed for SSH"
  type        = list(string)
}
```

**Create:** `terraform/modules/security-group/outputs.tf`

```hcl
output "security_group_id" {
  description = "ID of the security group"
  value       = aws_security_group.main.id
}
```

#### Step 5: EC2 Instance Module (25 min)

**Create:** `terraform/modules/ec2/main.tf`

```hcl
data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]
  
  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_instance" "main" {
  ami           = data.aws_ami.amazon_linux_2023.id
  instance_type = var.instance_type
  key_name      = var.key_name
  
  vpc_security_group_ids = [var.security_group_id]
  
  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }
  
  user_data = templatefile("${path.module}/user_data.sh", {
    domain_name = var.domain_name
  })
  
  tags = {
    Name = "${var.name_prefix}-instance"
  }
}

resource "aws_eip" "main" {
  instance = aws_instance.main.id
  domain   = "vpc"
  
  tags = {
    Name = "${var.name_prefix}-eip"
  }
}
```

**Create:** `terraform/modules/ec2/user_data.sh`

```bash
#!/bin/bash
set -e

# Update system
dnf update -y

# Install nginx
dnf install -y nginx

# Install certbot
dnf install -y certbot python3-certbot-nginx

# Create app directory
mkdir -p /var/www/stori
chown -R ec2-user:ec2-user /var/www/stori

# Start nginx
systemctl enable nginx
systemctl start nginx

echo "EC2 instance initialized"
```

#### Step 6: Main Terraform Configuration (20 min)

**Create:** `terraform/main.tf`

```hcl
module "security_group" {
  source = "./modules/security-group"
  
  name_prefix      = "stori-challenge"
  allowed_ssh_cidr = var.allowed_ssh_cidr
}

module "ec2" {
  source = "./modules/ec2"
  
  name_prefix       = "stori-challenge"
  instance_type     = var.instance_type
  key_name          = var.key_name
  security_group_id = module.security_group.security_group_id
  domain_name       = var.domain_name
}
```

**Create:** `terraform/outputs.tf`

```hcl
output "instance_public_ip" {
  description = "Public IP of EC2 instance"
  value       = module.ec2.public_ip
}

output "instance_id" {
  description = "ID of EC2 instance"
  value       = module.ec2.instance_id
}

output "security_group_id" {
  description = "ID of security group"
  value       = module.security_group.security_group_id
}
```

#### Step 7: Usage Documentation (10 min)

**Create:** `terraform/README.md`

```markdown
# Terraform Infrastructure

## Prerequisites
- Terraform >= 1.5.0
- AWS CLI configured
- SSH key pair created in AWS

## Setup

1. Initialize Terraform:
```bash
cd terraform
terraform init
```

2. Create `terraform.tfvars`:
```hcl
aws_region       = "us-east-1"
environment      = "production"
instance_type    = "t2.micro"
key_name         = "stori-expenses-backend-key"
domain_name      = "stori.danntastico.dev"
allowed_ssh_cidr = ["YOUR_IP/32"]
```

3. Plan:
```bash
terraform plan
```

4. Apply:
```bash
terraform apply
```

## Outputs
- `instance_public_ip`: EC2 public IP for DNS configuration
- `instance_id`: EC2 instance ID
- `security_group_id`: Security group ID
```

### Files to Create

**Create:**
- `terraform/provider.tf`
- `terraform/variables.tf`
- `terraform/main.tf`
- `terraform/outputs.tf`
- `terraform/modules/security-group/main.tf`
- `terraform/modules/security-group/variables.tf`
- `terraform/modules/security-group/outputs.tf`
- `terraform/modules/ec2/main.tf`
- `terraform/modules/ec2/variables.tf`
- `terraform/modules/ec2/outputs.tf`
- `terraform/modules/ec2/user_data.sh`
- `terraform/README.md`

### Testing Checklist

- [ ] `terraform init` succeeds
- [ ] `terraform plan` shows expected resources
- [ ] `terraform apply` creates infrastructure
- [ ] EC2 instance is accessible
- [ ] Security groups are configured correctly
- [ ] Outputs display correct values
- [ ] `terraform destroy` cleans up resources

---

## Feature Comparison Matrix

| Feature | Time | Impact | Complexity | Demo Value | Production Value |
|---------|------|--------|------------|-----------|-----------------|
| **Dark Mode** | 45m | High | Low | ★★★★☆ | ★★★☆☆ |
| **Skeletons** | 30m | Medium | Low | ★★★☆☆ | ★★★☆☆ |
| **Date Filter** | 1.5h | High | Medium | ★★★★★ | ★★★★★ |
| **CSV Export** | 45m | Medium | Low | ★★★☆☆ | ★★★★☆ |
| **CI/CD** | 1h | High | Medium | ★★★★☆ | ★★★★★ |
| **Terraform** | 1.5h | High | High | ★★★☆☆ | ★★★★★ |

---

## Recommendations

**Best for Interview Demo (Pick 1-2):**
1. **Dark Mode** - Quick win, visually impressive, shows modern UX thinking
2. **Date Range Filtering** - Adds real functionality, shows full-stack capability

**Best for Production Readiness:**
1. **CI/CD** - Shows DevOps maturity, automation mindset
2. **Terraform** - Infrastructure as Code best practice

**Easiest Quick Wins:**
1. **Loading Skeletons** - 30 minutes, better UX
2. **CSV Export** - 45 minutes, useful feature

---

## Next Steps

**Choose your feature(s) and I'll guide you through implementation!** 🚀

Which bonus feature would you like to implement?

