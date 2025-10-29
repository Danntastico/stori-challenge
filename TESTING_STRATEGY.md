# Testing Strategy

> Comprehensive testing approach for Stori Financial Tracker

---

## Current State (MVP)

### What Was Tested

**Manual Testing:**
- ✅ All API endpoints functional (health, transactions, summaries, advice)
- ✅ Mobile responsiveness (iPhone, Android browsers)
- ✅ CORS working cross-origin (frontend → backend)
- ✅ SSL certificate valid on all devices
- ✅ AI integration with real OpenAI API
- ✅ Charts render correctly with various data sizes
- ✅ Error states display properly
- ✅ Loading states work as expected

**Automated Tests:**
- ✅ Domain model validation (`internal/domain/transaction_test.go`)
- ✅ Date parsing and period calculations
- ✅ Transaction type validation

### Manual Testing Checklist Used

```
Desktop Testing:
- [ ] Visit https://stori.danntastico.dev/
- [ ] Financial overview displays correct totals
- [ ] Category pie chart renders and is interactive
- [ ] Timeline chart shows all 10 months
- [ ] Click "Get AI Advice" → receives response
- [ ] Hover over chart elements → tooltips work
- [ ] Resize window → components adapt responsively

Mobile Testing:
- [ ] Open URL on phone
- [ ] Charts resize properly
- [ ] Navigation works with touch
- [ ] AI advice button accessible
- [ ] No horizontal scrolling
- [ ] SSL certificate accepted (no warnings)

API Testing:
- [ ] GET /api/health → 200 OK
- [ ] GET /api/transactions → returns 114 transactions
- [ ] GET /api/summary/categories → valid percentages
- [ ] GET /api/summary/timeline → 10 data points
- [ ] POST /api/advice → structured response
```

---

## Production Testing Strategy

### Testing Pyramid

```
        /\
       /E2E\         Small number, high value
      /------\
     /  API   \      More tests, medium speed
    /----------\
   /    Unit    \    Many tests, fast execution
  /--------------\
```

---

## 1. Unit Tests

**Goal:** Test individual functions/components in isolation

### Backend Unit Tests

**Test Coverage Targets:**

| Component | Test Focus | Priority |
|-----------|-----------|----------|
| `domain/transaction.go` | ✅ Already tested | High |
| `service/analytics_service.go` | Aggregation logic | High |
| `service/ai_service.go` | Prompt building, response parsing | Medium |
| `repository/json_repository.go` | Date filtering, data loading | High |
| `handlers/*` | Input validation, error responses | Medium |

**Example Test Cases:**

```go
// internal/service/analytics_service_test.go
func TestCalculateCategorySummary(t *testing.T) {
    tests := []struct {
        name     string
        input    []domain.Transaction
        expected domain.CategorySummary
    }{
        {
            name: "single income transaction",
            input: []domain.Transaction{
                {Type: "income", Amount: 1000, Category: "salary"},
            },
            expected: domain.CategorySummary{
                Summary: domain.FinancialSummary{
                    TotalIncome: 1000,
                    TotalExpenses: 0,
                    NetSavings: 1000,
                    SavingsRate: 100.0,
                },
            },
        },
        {
            name: "mixed transactions",
            input: []domain.Transaction{
                {Type: "income", Amount: 1000, Category: "salary"},
                {Type: "expense", Amount: 300, Category: "rent"},
            },
            expected: domain.CategorySummary{
                Summary: domain.FinancialSummary{
                    TotalIncome: 1000,
                    TotalExpenses: 300,
                    NetSavings: 700,
                    SavingsRate: 70.0,
                },
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

func TestGetTimeline(t *testing.T) {
    // Test monthly aggregation
    // Test date range filtering
    // Test empty data handling
}
```

**Edge Cases to Test:**
- Empty transaction list
- Single transaction
- All income / all expenses
- Invalid dates
- Transactions on same day
- Missing category field
- Zero amounts

---

### Frontend Unit Tests

**Component Testing with React Testing Library:**

```typescript
// components/Dashboard/FinancialOverview.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { FinancialOverview } from './FinancialOverview';
import * as api from '../../services/api';

jest.mock('../../services/api');

describe('FinancialOverview', () => {
  it('displays loading state initially', () => {
    render(<FinancialOverview />);
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('renders financial metrics when data loads', async () => {
    const mockData = {
      summary: {
        total_income: 56000,
        total_expenses: 47000,
        net_savings: 9000,
        savings_rate: 16.1
      }
    };
    
    (api.getCategorySummary as jest.Mock).mockResolvedValue(mockData);
    
    render(<FinancialOverview />);
    
    await waitFor(() => {
      expect(screen.getByText('$56,000.00')).toBeInTheDocument();
      expect(screen.getByText('16.1%')).toBeInTheDocument();
    });
  });

  it('displays error message when API fails', async () => {
    (api.getCategorySummary as jest.Mock).mockRejectedValue(
      new Error('Network error')
    );
    
    render(<FinancialOverview />);
    
    await waitFor(() => {
      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });
  });
});
```

**Test Coverage Targets:**

| Component | Test Focus |
|-----------|-----------|
| `FinancialOverview` | Data rendering, error states |
| `CategoryChart` | Data transformation, Recharts integration |
| `TimelineChart` | Date formatting, multiple series |
| `AIAdvisor` | Loading states, response parsing |
| `useApiData` hook | Fetch logic, error handling, retries |

---

## 2. Integration Tests

**Goal:** Test interactions between components

### Backend Integration Tests

```go
// internal/handlers/summary_handler_test.go
func TestCategorySummaryEndpoint(t *testing.T) {
    // Setup
    repo := repository.NewJSONRepository(testData)
    service := service.NewAnalyticsService(repo)
    handler := handlers.NewSummaryHandler(service)
    
    router := chi.NewRouter()
    router.Get("/api/summary/categories", handler.GetCategorySummary)
    
    // Test
    req := httptest.NewRequest("GET", "/api/summary/categories", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response domain.CategorySummary
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Greater(t, response.Summary.TotalIncome, 0.0)
}

func TestAdviceEndpoint(t *testing.T) {
    // Test full flow: Handler → Analytics → AI → Response
}

func TestCORSMiddleware(t *testing.T) {
    // Test CORS headers with allowed/disallowed origins
}
```

**Integration Test Scenarios:**
- Complete request → response flow
- Middleware chain execution
- Error propagation through layers
- JSON serialization/deserialization
- Query parameter parsing

---

### Frontend Integration Tests

```typescript
// integration/api-integration.test.ts
describe('API Integration', () => {
  beforeAll(() => {
    // Start mock API server
  });

  it('fetches and displays complete dashboard', async () => {
    render(<App />);
    
    // Should fetch all endpoints
    await waitFor(() => {
      expect(screen.getByText(/Total Income/i)).toBeInTheDocument();
      expect(screen.getByRole('img', { name: /pie chart/i })).toBeInTheDocument();
      expect(screen.getByRole('img', { name: /timeline/i })).toBeInTheDocument();
    });
  });

  it('handles network failures gracefully', async () => {
    // Mock network failure
    server.use(
      rest.get('/api/summary/categories', (req, res, ctx) => {
        return res.networkError('Failed to connect')
      })
    );
    
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText(/connection failed/i)).toBeInTheDocument();
    });
  });
});
```

---

## 3. End-to-End Tests

**Goal:** Test complete user flows in real browser

### E2E Test Framework

**Tool:** Playwright (cross-browser, fast, reliable)

```typescript
// e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Financial Dashboard', () => {
  test('user can view complete financial overview', async ({ page }) => {
    await page.goto('https://stori.danntastico.dev/');
    
    // Wait for data to load
    await expect(page.getByText('Total Income')).toBeVisible();
    
    // Check all sections render
    await expect(page.getByText('Spending by Category')).toBeVisible();
    await expect(page.getByText('Income vs Expenses Timeline')).toBeVisible();
    await expect(page.getByRole('button', { name: /Get AI Advice/i })).toBeVisible();
    
    // Take screenshot for visual regression
    await page.screenshot({ path: 'dashboard.png', fullPage: true });
  });

  test('user can request AI advice', async ({ page }) => {
    await page.goto('https://stori.danntastico.dev/');
    
    // Click AI advice button
    await page.click('text=Get AI Advice');
    
    // Wait for loading state
    await expect(page.getByText(/generating/i)).toBeVisible();
    
    // Wait for advice to load (max 10 seconds)
    await expect(page.getByText(/INSIGHTS/i)).toBeVisible({ timeout: 10000 });
    await expect(page.getByText(/RECOMMENDATIONS/i)).toBeVisible();
  });

  test('dashboard is responsive on mobile', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('https://stori.danntastico.dev/');
    
    // Charts should resize
    const chart = page.locator('[data-testid="category-chart"]');
    await expect(chart).toBeVisible();
    
    // No horizontal scrolling
    const scrollWidth = await page.evaluate(() => document.body.scrollWidth);
    const clientWidth = await page.evaluate(() => document.body.clientWidth);
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth + 1); // +1 for rounding
  });
});
```

**E2E Test Scenarios:**
- Initial page load and data display
- Chart interactions (hover, click)
- AI advice generation
- Error recovery
- Mobile responsiveness
- Cross-browser compatibility (Chrome, Firefox, Safari)
- SSL certificate validation

---

## 4. Performance Tests

**Goal:** Ensure system handles load gracefully

### Load Testing with k6

```javascript
// load-tests/api-load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Ramp up
    { duration: '3m', target: 50 },   // Stay at 50 users
    { duration: '1m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% requests under 500ms
    http_req_failed: ['rate<0.01'],   // Error rate under 1%
  },
};

export default function () {
  const baseUrl = 'https://stori.danntastico.dev/api';
  
  // Test health endpoint
  let res = http.get(`${baseUrl}/health`);
  check(res, { 'health check ok': (r) => r.status === 200 });
  
  // Test category summary
  res = http.get(`${baseUrl}/summary/categories`);
  check(res, { 'categories ok': (r) => r.status === 200 });
  
  sleep(1);
}
```

**Performance Metrics to Track:**
- Response time (p50, p95, p99)
- Throughput (requests/second)
- Error rate
- Resource usage (CPU, memory)
- Database query time (when added)
- OpenAI API latency

---

## 5. Security Tests

**Security Testing Checklist:**

```
Authentication (when implemented):
- [ ] Weak password rejection
- [ ] SQL injection attempts (when DB added)
- [ ] XSS attack prevention
- [ ] CSRF token validation

API Security:
- [ ] CORS policy enforcement
- [ ] Rate limiting (prevent abuse)
- [ ] Input validation (malformed JSON, invalid dates)
- [ ] API key exposure (not in responses)

Infrastructure:
- [ ] SSL certificate validity
- [ ] HTTPS enforcement (HTTP → HTTPS redirect)
- [ ] Security headers (CSP, X-Frame-Options)
- [ ] Dependency vulnerabilities (npm audit, go mod)
```

---

## Test Automation & CI/CD

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Run tests
        run: |
          cd backend
          go test ./... -v -cover
          go test ./... -race

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Run tests
        run: cd frontend && npm test
      - name: Type check
        run: cd frontend && npm run type-check

  e2e-tests:
    runs-on: ubuntu-latest
    needs: [backend-tests, frontend-tests]
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - name: Install Playwright
        run: npx playwright install --with-deps
      - name: Run E2E tests
        run: npx playwright test
      - name: Upload screenshots
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-screenshots
          path: playwright-results/
```

---

## Test Coverage Goals

| Component | Current | Target |
|-----------|---------|--------|
| **Backend Domain** | 80% | 90% |
| **Backend Services** | 0% | 85% |
| **Backend Handlers** | 0% | 75% |
| **Frontend Components** | 0% | 80% |
| **Frontend Hooks** | 0% | 90% |
| **E2E Critical Paths** | Manual | 100% |

---

## Testing Best Practices

### General Principles

1. **Test Behavior, Not Implementation**
   - ❌ Bad: `expect(component.state.isLoading).toBe(true)`
   - ✅ Good: `expect(screen.getByText('Loading...')).toBeVisible()`

2. **Arrange-Act-Assert Pattern**
   ```typescript
   test('calculates savings rate', () => {
     // Arrange
     const income = 1000;
     const expenses = 300;
     
     // Act
     const rate = calculateSavingsRate(income, expenses);
     
     // Assert
     expect(rate).toBe(70.0);
   });
   ```

3. **Independent Tests**
   - Each test should run in isolation
   - No shared state between tests
   - Use `beforeEach` to reset state

4. **Descriptive Test Names**
   - ❌ Bad: `test('handler works')`
   - ✅ Good: `test('returns 400 when date format is invalid')`

5. **Test Edge Cases**
   - Empty inputs
   - Boundary values (0, negative, very large)
   - Invalid data types
   - Network failures

---

## Rationale for Current Approach

**Why manual testing for MVP?**

✅ **Speed:** Faster to deliver working product  
✅ **Flexibility:** Easier to iterate on design  
✅ **Learning:** Understanding system behavior before writing tests  
✅ **Pragmatic:** Interview demo doesn't need 100% coverage  

**Why comprehensive strategy document?**

✅ **Shows Knowledge:** Demonstrates testing expertise  
✅ **Production Ready:** Clear path for real implementation  
✅ **Risk Awareness:** Acknowledges current gaps  
✅ **Professionalism:** Plans for quality, not just features  

---

## Implementation Timeline

**Phase 1 (Week 1):** Backend unit tests
- Domain and service layers
- 80%+ coverage

**Phase 2 (Week 2):** Frontend component tests
- All major components
- Custom hooks

**Phase 3 (Week 3):** Integration tests
- API endpoints
- Full data flow

**Phase 4 (Week 4):** E2E tests
- Critical user paths
- Cross-browser testing

**Phase 5 (Ongoing):** CI/CD integration
- Automated test runs
- Coverage reporting
- Performance monitoring

---

## Conclusion

This testing strategy balances **pragmatic MVP delivery** with **production quality standards**. The current manual testing validates core functionality, while the documented strategy provides a clear roadmap for comprehensive automated testing.

**Key Takeaway:** Testing isn't an afterthought—it's a planned evolution from MVP to production-ready system.

---

**Last Updated:** October 2025  
**Current Coverage:** Domain models (automated), Full system (manual)  
**Target Coverage:** 85%+ automated within 4 weeks of production launch

