# User Stories - Stori Challenge

## MVP Requirements (Must Implement - Days 1-3)

### Epic 1: Financial Overview Dashboard
**Priority: P0**
- As a user, I want to see my total income, expenses, and net savings at a glance
- As a user, I want to see the time period covered by my financial data
- As a user, I want to quickly understand my overall financial health

**Acceptance Criteria:**
- Display total income, total expenses, net savings
- Show date range (start and end dates)
- Clear visual hierarchy with key metrics prominent

---

### Epic 2: Spending Analysis by Category
**Priority: P0** (Required by challenge)
- As a user, I want to see a visual chart of my spending breakdown by category
- As a user, I want to see how much I spent in each expense category with dollar amounts
- As a user, I want to understand what percentage each category represents of my total spending

**Acceptance Criteria:**
- Visual chart (pie, donut, or bar chart) showing expense categories
- Display both amounts and percentages
- Categories: rent, groceries, utilities, dining, transportation, entertainment, shopping, healthcare
- Clear labeling and color coding

---

### Epic 3: Income & Expense Timeline
**Priority: P0** (Required by challenge)
- As a user, I want to see a timeline chart of my income and expenses over time
- As a user, I want to identify trends and patterns in my cash flow
- As a user, I want to compare my monthly income vs expenses visually

**Acceptance Criteria:**
- Timeline chart showing income and expenses over time (monthly aggregation)
- Clear differentiation between income and expense lines/bars
- X-axis: time periods, Y-axis: amounts
- Easy to spot months with high spending or low savings

---

### Epic 4: AI-Powered Financial Advice
**Priority: P0** (Required by challenge)
- As a user, I want to receive personalized financial advice based on my actual spending data
- As a user, I want actionable recommendations to improve my savings
- As a user, I want the advice to be easy to understand and contextual
- As a user, I want to request new advice when needed

**Acceptance Criteria:**
- AI advisor analyzes actual transaction data
- Provides specific insights (e.g., "You spent $800 on dining in October")
- Gives actionable recommendations
- Plain language, no jargon
- Button/trigger to request advice

---

### Epic 5: Mobile-First Responsive Design
**Priority: P0** (Required by challenge)
- As a user, I want to access the app from my mobile phone
- As a user, I want charts and data to be readable on small screens
- As a user, I want touch-friendly interactions

**Acceptance Criteria:**
- Responsive design that works on mobile (320px+)
- Charts resize appropriately
- Touch targets are at least 44px
- No horizontal scrolling required

---

## Phase 2: Enhanced Features (Document Only - Future Improvements)

### Epic 6: Transaction Management (CRUD Operations)
**Priority: P1** (Nice to have - Time permitting Day 4)
- As a user, I want to add new transactions manually
- As a user, I want to edit existing transactions to correct mistakes
- As a user, I want to delete transactions I added by mistake
- As a user, I want to see my changes reflected immediately in charts

**Technical Notes:**
- Requires database migration from JSON to PostgreSQL/DynamoDB
- Add POST, PUT, DELETE endpoints
- Form validation (date, amount, category)
- Optimistic UI updates

---

### Epic 7: Advanced Filtering & Views
**Priority: P2**
- As a user, I want to filter transactions by custom date ranges
- As a user, I want to filter by specific categories
- As a user, I want to export my data as CSV
- As a user, I want to see detailed transaction lists with search

**Technical Notes:**
- Query parameters for API filtering
- Frontend date picker components
- CSV generation endpoint
- Pagination for large datasets

---

## Phase 3: Advanced Features (Long-term Vision - Document in Architecture)

### Epic 8: Multi-Account Management
**Priority: P3**
- As a user, I want to track multiple accounts (Cash, Checking, Savings, Credit Card)
- As a user, I want to see balances per account
- As a user, I want to categorize transactions by account
- As a user, I want to transfer money between accounts

**Technical Notes:**
- Add Account entity and relationships
- User authentication system (AWS Cognito)
- Multi-tenant data isolation
- Account-level transaction tracking

---

### Epic 9: Budget Tracking & Alerts
**Priority: P3**
- As a user, I want to set monthly budgets per expense category
- As a user, I want to track progress toward my budgets
- As a user, I want to receive alerts when approaching budget limits
- As a user, I want to see budget vs actual spending comparisons

**Technical Notes:**
- Budget entity and rules engine
- Email/SMS notification service (SNS)
- Real-time budget calculation
- Visual progress indicators

---

### Epic 10: Intelligent Insights & Predictions
**Priority: P4**
- As a user, I want automatic detection of recurring transactions
- As a user, I want predictions of future expenses
- As a user, I want anomaly detection for unusual spending
- As a user, I want personalized savings recommendations

**Technical Notes:**
- ML model for pattern recognition
- Time-series forecasting
- Anomaly detection algorithms
- Enhanced AI prompts with historical context

---

### Epic 11: Financial Goals & Planning
**Priority: P4**
- As a user, I want to set savings goals (e.g., emergency fund, vacation)
- As a user, I want to track progress toward my goals
- As a user, I want suggestions on how to reach goals faster
- As a user, I want to visualize my path to financial milestones

**Technical Notes:**
- Goals entity with target amounts and dates
- Progress calculation engine
- Goal-oriented AI advice
- Gamification elements (progress bars, achievements)

---

## Implementation Priority Summary

```
MVP (Build):
✅ Epic 1: Financial Overview Dashboard
✅ Epic 2: Spending Analysis by Category  
✅ Epic 3: Income & Expense Timeline
✅ Epic 4: AI-Powered Financial Advice
✅ Epic 5: Mobile-First Responsive Design

Phase 2 (Document + Optional Implementation):
📝 Epic 6: Transaction Management (CRUD)
📝 Epic 7: Advanced Filtering & Views

Phase 3+ (Document as Future Vision):
📝 Epic 8: Multi-Account Management
📝 Epic 9: Budget Tracking & Alerts
📝 Epic 10: Intelligent Insights & Predictions
📝 Epic 11: Financial Goals & Planning
```

---

## Alignment with Challenge Requirements

| Challenge Requirement | Epic Coverage |
|----------------------|---------------|
| "See a summary of spending by expense category" | Epic 1, Epic 2 |
| "See a timeline of income and expenses" | Epic 3 |
| "Get advice about how to manage spending to save more money" | Epic 4 |
| "Mobile-friendly web application" | Epic 5 |
| "REST APIs" | All MVP epics (backend) |
| "AWS cloud hosting" | Infrastructure (Day 3) |

---

## Notes for System Design Document

**What to emphasize in writeup:**
- Started with MVP scope aligned with challenge requirements
- Designed architecture to accommodate Phases 2-3 (repository pattern, clean separation)
- Prioritized working prototype over feature richness
- Focused on core value: insights and visualization
- Future roadmap shows product thinking and scalability awareness