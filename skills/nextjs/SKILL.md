---
name: nextjs
description: Next.js feature development following 2025 best practices for App Router, Server Components, Server Actions, Pages Router, routing, optimization, and Cypress testing. Uses pnpm package manager. Use when adding or updating Next.js applications.
allowed-tools: [Read, Write, Edit, Glob, Grep, Bash]
---

# Next.js Development Skill

A comprehensive skill for adding or updating features in Next.js projects following 2025 best practices, with emphasis on App Router, Server Components, security, and table-driven testing with Cypress.

## When to Use

Use this skill when:
- Adding new features to Next.js projects
- Updating or enhancing existing Next.js functionality
- Working with App Router or Pages Router
- Need to follow Next.js-specific patterns and conventions
- Building server or client components
- Implementing server actions or API routes

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation:

1. **Look for `.claude/design.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/design.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no design.md exists:**
   - Consider running `/document-design` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from design.md:**
- Project-specific naming conventions
- Directory structure (App Router vs Pages Router)
- Component organization patterns
- Server vs Client component guidelines
- Testing conventions
- Styling approach
- Code examples showing preferred style

### 1. Analyze Project Structure

- Check for `next.config.js`, `package.json`
- Verify Node.js version requirement (minimum 20.9 as of 2025)
- Confirm pnpm is the package manager (MUST NOT use legacy-peer-deps)
- Identify if using App Router (app/) or Pages Router (pages/)
- Review existing component patterns
- Check for styling approach (CSS modules, Tailwind, styled-components)
- Identify state management (if any)
- Review data fetching patterns
- Check for existing Data Access Layer for database operations

### 2. Search for Relevant Code

- Search for similar components or pages
- Look for related API routes or server actions
- Find existing hooks or utilities
- Identify styling patterns

## Next.js Best Practices (2025)

### Core Principles

1. **Simplicity First**: Follow DRY principles, update existing code for reusability rather than creating new functions
2. **Consistency**: Follow project-specific guidelines from `.claude/design.md`
3. **Fail-Fast**: Ensure errors are caught early and not silently suppressed
4. **Latest Versions**: Use latest stable versions of packages (Next.js 15+ with React 19+ as of 2025)
5. **Encapsulation**: Keep code private/unexported unless necessary
6. **Early Returns**: Prefer early returns and continue statements over nested conditionals

### Package Manager Requirements

- **MUST use pnpm** as the package manager (not npm or yarn)
- **MUST NOT use** `legacy-peer-deps` or `--force` flags
- Setup: `pnpm create next-app@latest` for new projects
- Commands: `pnpm i`, `pnpm dev`, `pnpm build`, `pnpm test`

### App Router (Next.js 13+)

**Server Components (Default):**
- Modern Next.js defaults to Server Components for optimal performance
- Use Server Components by default - no unnecessary JavaScript in browser
- Faster load times, smaller bundles, automatic SEO benefits
- Add 'use client' only when needed (interactivity, browser APIs, event handlers, React hooks)
- Leverage async/await in Server Components
- Direct data fetching without useEffect
- Keep client components focused and minimal

Example:
```typescript
// app/users/page.tsx - Server Component (default)
import { fetchUsers } from '@/lib/api';

export default async function UsersPage() {
  const users = await fetchUsers(); // Direct async call

  return (
    <div>
      <h1>Users</h1>
      <UserList users={users} />
    </div>
  );
}
```

**Client Components:**
- Add 'use client' directive for interactive components
- Use for browser APIs, event handlers, hooks
- Keep client components focused and minimal

Example:
```typescript
// app/users/user-list.tsx - Client Component (interactive)
'use client';

import { useState } from 'react';
import type { User } from '@/types';

interface UserListProps {
  users: User[];
}

export function UserList({ users }: UserListProps) {
  const [filter, setFilter] = useState('');

  const filtered = users.filter(u =>
    u.name.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div>
      <input
        type="text"
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
        placeholder="Filter users..."
      />
      {filtered.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
}
```

**Server Actions (2025 Best Practices):**

Server Actions are the idiomatic way to perform writes/mutations in Next.js App Router. Starting from Next.js 15/React 19, Server Actions are a type of Server Function.

**Key Guidelines:**
- Use for side-effects like saving changes, mutations, or logging out
- Define with `'use server'` directive at function or module level
- Always create separate `actions.ts` files (better organization and security)
- Use POST method only (automatically enforced for Server Actions)
- For reading data from Client Components, prefer Route Handlers (GET requests)
- Implement proper form validation with Zod or TypeScript interfaces
- Return structured data with success/error states
- Use revalidatePath/revalidateTag for cache updates
- Verify Origin header matches Host header (automatic CSRF protection)

**Security Considerations:**
- Server Actions use POST method, preventing most CSRF vulnerabilities
- Origin/Host header comparison automatically protects against CSRF
- Never expose database packages or environment variables outside Data Access Layer
- Create isolated Data Access Layer for database operations

**File Organization:**
```typescript
// Good: Separate actions file with module-level directive
// app/actions/user.ts
'use server';

import { revalidatePath } from 'next/cache';
import { z } from 'zod';

const createUserSchema = z.object({
  name: z.string().min(1).max(100),
  email: z.string().email(),
});

export async function createUser(formData: FormData) {
  // Early return on validation failure
  const parsed = createUserSchema.safeParse({
    name: formData.get('name'),
    email: formData.get('email'),
  });

  if (!parsed.success) {
    return { success: false, errors: parsed.error.flatten() };
  }

  try {
    // Use Data Access Layer instead of direct db access
    await db.user.create({ data: parsed.data });
    revalidatePath('/users');
    return { success: true };
  } catch (error) {
    // Fail-fast with explicit error
    return { success: false, error: 'Failed to create user' };
  }
}
```

**Using Server Actions in Client Components:**
```typescript
// Bad: Cannot define server action directly in client component
'use client';
export default function Form() {
  async function createUser() {
    'use server'; // This will not work!
    // ...
  }
}

// Good: Import server action from separate file
'use client';
import { createUser } from '@/app/actions/user';

export default function Form() {
  return (
    <form action={createUser}>
      {/* form fields */}
    </form>
  );
}
```

**Server Actions vs Route Handlers:**
- Use **Server Actions** for mutations (create, update, delete) - POST only
- Use **Route Handlers** for fetching data from Client Components - supports GET requests
- Route Handlers are suitable for reading data as GET requests can be cached

**Metadata and SEO:**
```typescript
// app/users/[id]/page.tsx
import { Metadata } from 'next';

interface UserPageProps {
  params: { id: string };
}

export async function generateMetadata({ params }: UserPageProps): Promise<Metadata> {
  const user = await getUserById(params.id);

  return {
    title: user ? `${user.name} - Profile` : 'User Not Found',
    description: user?.bio,
  };
}

export default async function UserPage({ params }: UserPageProps) {
  const user = await getUserById(params.id);

  if (!user) {
    notFound();
  }

  return <UserProfile user={user} />;
}
```

### Pages Router (Legacy)

**Data Fetching:**
- Use getServerSideProps for dynamic data
- Use getStaticProps + ISR for static content
- Implement getStaticPaths for dynamic routes

Example:
```typescript
// pages/users/[id].tsx
import type { GetServerSideProps } from 'next';
import type { User } from '@/types';

interface UserPageProps {
  user: User;
}

export default function UserPage({ user }: UserPageProps) {
  return (
    <div>
      <h1>{user.name}</h1>
      <p>{user.email}</p>
    </div>
  );
}

export const getServerSideProps: GetServerSideProps<UserPageProps> = async (context) => {
  const { id } = context.params!;

  try {
    const user = await fetchUser(id as string);

    if (!user) {
      return { notFound: true };
    }

    return { props: { user } };
  } catch (error) {
    return { notFound: true };
  }
};
```

### Component Patterns

**Prop Typing:**
```typescript
// Good: Well-typed reusable component
import { memo } from 'react';

interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
}

export const Button = memo<ButtonProps>(({
  variant = 'primary',
  size = 'md',
  disabled = false,
  onClick,
  children
}) => {
  return (
    <button
      className={`btn btn-${variant} btn-${size}`}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
});

Button.displayName = 'Button';
```

### Optimization

**Image Optimization:**
```typescript
import Image from 'next/image';

// Good: Optimized image
<Image
  src="/profile.jpg"
  alt="Profile"
  width={200}
  height={200}
  priority // for LCP images
/>
```

**Client-Side Navigation:**
```typescript
import Link from 'next/link';

// Good: Client-side navigation
<Link href="/dashboard">Dashboard</Link>
```

**Code Splitting:**
```typescript
import dynamic from 'next/dynamic';

// Good: Dynamic import with loading state
const HeavyChart = dynamic(() => import('@/components/HeavyChart'), {
  loading: () => <div>Loading chart...</div>,
  ssr: false, // disable SSR if not needed
});
```

### API Routes (Route Handlers)

**When to Use Route Handlers:**
- Fetching data from Client Components (prefer GET requests)
- External API integrations
- Webhooks
- Custom API endpoints for third-party services

**App Router Route Handlers (2025 Best Practices):**
```typescript
// app/api/users/[id]/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { z } from 'zod';

const updateUserSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  email: z.string().email().optional(),
});

export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    // Use Data Access Layer for database operations
    const user = await getUserById(params.id);

    // Early return for not found
    if (!user) {
      return NextResponse.json(
        { error: 'User not found' },
        { status: 404 }
      );
    }

    return NextResponse.json({ user });
  } catch (error) {
    // Fail-fast with explicit error
    return NextResponse.json(
      { error: 'Failed to fetch user' },
      { status: 500 }
    );
  }
}

export async function PATCH(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const body = await request.json();

    // Early return on validation failure (fail-fast)
    const result = updateUserSchema.safeParse(body);
    if (!result.success) {
      return NextResponse.json(
        { error: 'Validation failed', details: result.error.errors },
        { status: 400 }
      );
    }

    // Use Data Access Layer
    const user = await db.user.update({
      where: { id: params.id },
      data: result.data,
    });

    return NextResponse.json({ user, message: 'User updated' });
  } catch (error) {
    // Fail-fast with explicit error
    return NextResponse.json(
      { error: 'Failed to update user' },
      { status: 500 }
    );
  }
}
```

**Type Safety for API Responses (2025):**
```typescript
// Good: Use TypeScript interfaces for strict API response validation
interface UserResponse {
  user: {
    id: string;
    name: string;
    email: string;
  };
}

// Or use Zod for runtime validation
const userResponseSchema = z.object({
  user: z.object({
    id: z.string(),
    name: z.string(),
    email: z.string().email(),
  }),
});
```

### Testing

**CRITICAL: Always use table-driven tests as the primary testing approach.**

Table-driven tests provide the same benefits in Next.js projects:
- Reduce code duplication
- Make it easy to add new test cases
- Improve test readability
- Ensure consistent test structure
- Make test coverage gaps obvious
- Avoid redundant, meaningless test cases with the same purpose

**IMPORTANT: Separate happy path and error path test cases**
- Create separate test case arrays for success and error scenarios
- Use `successTestCases` for happy path scenarios
- Use `errorTestCases` for error scenarios
- This separation improves test readability and maintenance
- Makes test intent clearer
- Split happy and error test sets if they are complicated

**Testing Principles:**
1. Do not add redundant test cases with the same purpose
2. Prefer injecting values over changing global states
3. Define test inputs as test case fields, not function arguments
4. Keep tests independent - each test should set up its own environment
5. Do not rely on other tests

**Use table-driven tests for:**
- Component rendering with different props
- Utility functions with multiple test cases
- Server Actions with various inputs
- API route handlers
- Edge cases and error scenarios
- E2E tests with Cypress

### Cypress Testing (2025 Best Practices)

**Setup and Configuration:**
```bash
# Install Cypress
pnpm add -D cypress start-server-and-test

# Add to package.json scripts
{
  "scripts": {
    "cypress:open": "cypress open",
    "cypress:run": "cypress run",
    "test:e2e": "start-server-and-test start http://localhost:3000 cypress:run"
  }
}
```

**cypress.config.js:**
```javascript
import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
  component: {
    devServer: {
      framework: 'next',
      bundler: 'webpack',
    },
  },
});
```

**Key Cypress Best Practices:**
1. **E2E vs Component Testing**: Use E2E for full user flows, Component Testing for individual components
2. **API Mocking**: Use `cy.intercept()` for mocking API calls
3. **Fixtures**: Store mock data in `/cypress/fixtures` directory
4. **Test Isolation**: Each test should be independent with its own setup
5. **data-testid Attributes**: Use for reliable element selection (not CSS classes)
6. **CI/CD Integration**: Run Cypress tests in your CI/CD pipeline

**Component Testing Notes:**
- Component tests do not launch a Next.js server
- `<Image />` and `getServerSideProps` will not work without workarounds
- See Cypress Next.js docs for solutions

**Table-Driven Test Pattern:**

```typescript
// Good: Separate happy path and error path test cases
describe('Feature', () => {
  // Success scenarios
  const successTestCases = [
    { name: 'valid case 1', input: '...', expected: '...' },
    { name: 'valid case 2', input: '...', expected: '...' },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      // test implementation
    });
  });

  // Error scenarios
  const errorTestCases = [
    { name: 'invalid case 1', input: '...', expectedError: '...' },
    { name: 'invalid case 2', input: '...', expectedError: '...' },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      // test implementation
    });
  });
});
```

**Unit and Integration Test Examples (Jest/Vitest):**

```typescript
// Good: Testing component rendering with different props
import { render, screen } from '@testing-library/react';
import { Button } from './Button';

describe('Button', () => {
  const successTestCases = [
    { name: 'renders primary variant', variant: 'primary', expectedClass: 'btn-primary' },
    { name: 'renders secondary variant', variant: 'secondary', expectedClass: 'btn-secondary' },
    { name: 'renders danger variant', variant: 'danger', expectedClass: 'btn-danger' },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      render(<Button variant={testCase.variant}>Click me</Button>);
      const button = screen.getByRole('button');
      expect(button).toHaveClass(testCase.expectedClass);
    });
  });
});

// Good: Testing server actions with separated success/error cases
import { createUser } from './actions';

describe('createUser', () => {
  const successTestCases = [
    {
      name: 'creates user with valid data',
      formData: new FormData([['name', 'John'], ['email', 'john@example.com']]),
      expected: { success: true },
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      const result = await createUser(testCase.formData);
      expect(result).toMatchObject(testCase.expected);
    });
  });

  const errorTestCases = [
    {
      name: 'rejects missing name',
      formData: new FormData([['email', 'john@example.com']]),
      expectedError: { success: false, error: expect.any(Object) },
    },
    {
      name: 'rejects invalid email',
      formData: new FormData([['name', 'John'], ['email', 'invalid']]),
      expectedError: { success: false, error: expect.any(Object) },
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      const result = await createUser(testCase.formData);
      expect(result).toMatchObject(testCase.expectedError);
    });
  });
});

// Good: Testing utility functions
import { formatPrice } from '@/lib/utils';

describe('formatPrice', () => {
  const successTestCases = [
    { name: 'formats USD cents correctly', amount: 1000, currency: 'USD', expected: '$10.00' },
    { name: 'formats USD with decimals', amount: 2550, currency: 'USD', expected: '$25.50' },
    { name: 'formats zero amount', amount: 0, currency: 'USD', expected: '$0.00' },
    { name: 'formats EUR cents correctly', amount: 1000, currency: 'EUR', expected: '€10.00' },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      expect(formatPrice(testCase.amount, testCase.currency)).toBe(testCase.expected);
    });
  });
});

// Good: Testing API routes with separated success/error cases
import { GET, POST } from '@/app/api/users/route';
import { NextRequest } from 'next/server';

describe('/api/users', () => {
  describe('POST', () => {
    const successTestCases = [
      {
        name: 'creates user successfully',
        body: { name: 'John Doe', email: 'john@example.com' },
        expectedStatus: 200,
        expectedBody: { success: true },
      },
    ];

    successTestCases.forEach((testCase) => {
      it(testCase.name, async () => {
        const request = new NextRequest('http://localhost/api/users', {
          method: 'POST',
          body: JSON.stringify(testCase.body),
        });

        const response = await POST(request);
        expect(response.status).toBe(testCase.expectedStatus);

        const data = await response.json();
        expect(data).toMatchObject(testCase.expectedBody);
      });
    });

    const errorTestCases = [
      {
        name: 'rejects invalid email',
        body: { name: 'John Doe', email: 'invalid' },
        expectedStatus: 400,
        expectedBody: { error: expect.any(String) },
      },
      {
        name: 'rejects missing name',
        body: { email: 'john@example.com' },
        expectedStatus: 400,
        expectedBody: { error: expect.any(String) },
      },
    ];

    errorTestCases.forEach((testCase) => {
      it(testCase.name, async () => {
        const request = new NextRequest('http://localhost/api/users', {
          method: 'POST',
          body: JSON.stringify(testCase.body),
        });

        const response = await POST(request);
        expect(response.status).toBe(testCase.expectedStatus);

        const data = await response.json();
        expect(data).toMatchObject(testCase.expectedBody);
      });
    });
  });
});

// Good: Testing React hooks
import { renderHook, act } from '@testing-library/react';
import { useCounter } from '@/hooks/useCounter';

describe('useCounter', () => {
  const successTestCases = [
    { name: 'increments from 0 by 1', initialValue: 0, incrementBy: 1, expected: 1 },
    { name: 'increments from 5 by 2', initialValue: 5, incrementBy: 2, expected: 7 },
    { name: 'decrements from 10 by 3', initialValue: 10, incrementBy: -3, expected: 7 },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      const { result } = renderHook(() => useCounter(testCase.initialValue));

      act(() => {
        result.current.increment(testCase.incrementBy);
      });

      expect(result.current.count).toBe(testCase.expected);
    });
  });
});

// Good: Testing with React Testing Library user interactions
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SearchForm } from './SearchForm';

describe('SearchForm', () => {
  const successTestCases = [
    {
      name: 'submits valid search query',
      input: 'test query',
      expectedValue: 'test query',
    },
    {
      name: 'trims whitespace from query',
      input: '  test query  ',
      expectedValue: 'test query',
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      const onSubmit = jest.fn();
      const user = userEvent.setup();

      render(<SearchForm onSubmit={onSubmit} />);

      const searchInput = screen.getByRole('textbox');
      const submitButton = screen.getByRole('button');

      await user.type(searchInput, testCase.input);
      await user.click(submitButton);

      expect(onSubmit).toHaveBeenCalledWith(testCase.expectedValue);
    });
  });

  const errorTestCases = [
    {
      name: 'does not submit empty query',
      input: '',
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      const onSubmit = jest.fn();
      const user = userEvent.setup();

      render(<SearchForm onSubmit={onSubmit} />);

      const searchInput = screen.getByRole('textbox');
      const submitButton = screen.getByRole('button');

      await user.type(searchInput, testCase.input);
      await user.click(submitButton);

      expect(onSubmit).not.toHaveBeenCalled();
    });
  });
});

// Good: Testing metadata generation
import { generateMetadata } from '@/app/users/[id]/page';

describe('generateMetadata', () => {
  const successTestCases = [
    {
      name: 'generates metadata for existing user',
      params: { id: '123' },
      mockUser: { id: '123', name: 'John Doe', bio: 'Developer' },
      expected: {
        title: 'John Doe - Profile',
        description: 'Developer',
      },
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      // Mock the getUserById function
      jest.mock('@/lib/api', () => ({
        getUserById: jest.fn().mockResolvedValue(testCase.mockUser),
      }));

      const metadata = await generateMetadata({ params: testCase.params });
      expect(metadata).toMatchObject(testCase.expected);
    });
  });

  const errorTestCases = [
    {
      name: 'generates metadata for non-existent user',
      params: { id: '999' },
      mockUser: null,
      expected: {
        title: 'User Not Found',
        description: undefined,
      },
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, async () => {
      // Mock the getUserById function
      jest.mock('@/lib/api', () => ({
        getUserById: jest.fn().mockResolvedValue(testCase.mockUser),
      }));

      const metadata = await generateMetadata({ params: testCase.params });
      expect(metadata).toMatchObject(testCase.expected);
    });
  });
});
```

**E2E Testing with Cypress:**

Cypress tests should also use the table-driven pattern with separated success and error cases.

**Using data-testid Attributes:**
```typescript
// Good: Use data-testid in your components for reliable testing
export function LoginForm() {
  return (
    <form>
      <input
        data-testid="email-input"
        type="email"
        name="email"
      />
      <input
        data-testid="password-input"
        type="password"
        name="password"
      />
      <button data-testid="submit-button">Login</button>
    </form>
  );
}
```

**Cypress E2E Test Examples:**

```typescript
// cypress/e2e/login.cy.ts
describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('/login');
  });

  // Success scenarios
  const successTestCases = [
    {
      name: 'successfully logs in with valid credentials',
      email: 'user@example.com',
      password: 'password123',
      expectedUrl: '/dashboard',
      expectedMessage: 'Welcome back!',
    },
    {
      name: 'remembers user when "Remember me" is checked',
      email: 'user@example.com',
      password: 'password123',
      rememberMe: true,
      expectedUrl: '/dashboard',
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      // Use data-testid for reliable element selection
      cy.get('[data-testid="email-input"]').type(testCase.email);
      cy.get('[data-testid="password-input"]').type(testCase.password);

      if (testCase.rememberMe) {
        cy.get('[data-testid="remember-me-checkbox"]').check();
      }

      cy.get('[data-testid="submit-button"]').click();

      cy.url().should('include', testCase.expectedUrl);

      if (testCase.expectedMessage) {
        cy.contains(testCase.expectedMessage).should('be.visible');
      }
    });
  });

  // Error scenarios
  const errorTestCases = [
    {
      name: 'shows error with invalid email',
      email: 'invalid-email',
      password: 'password123',
      expectedError: 'Please enter a valid email address',
    },
    {
      name: 'shows error with empty password',
      email: 'user@example.com',
      password: '',
      expectedError: 'Password is required',
    },
    {
      name: 'shows error with wrong credentials',
      email: 'user@example.com',
      password: 'wrongpassword',
      expectedError: 'Invalid email or password',
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      if (testCase.email) {
        cy.get('[data-testid="email-input"]').type(testCase.email);
      }

      if (testCase.password) {
        cy.get('[data-testid="password-input"]').type(testCase.password);
      }

      cy.get('[data-testid="submit-button"]').click();

      cy.contains(testCase.expectedError).should('be.visible');
      cy.url().should('include', '/login'); // Should stay on login page
    });
  });
});

// cypress/e2e/user-profile.cy.ts
describe('User Profile Page', () => {
  beforeEach(() => {
    cy.login('user@example.com', 'password123'); // Custom command
  });

  const successTestCases = [
    {
      name: 'updates profile name successfully',
      field: 'name',
      value: 'John Smith',
      expectedMessage: 'Profile updated successfully',
    },
    {
      name: 'updates profile bio successfully',
      field: 'bio',
      value: 'Software developer',
      expectedMessage: 'Profile updated successfully',
    },
    {
      name: 'updates avatar successfully',
      field: 'avatar',
      file: 'avatar.jpg',
      expectedMessage: 'Avatar updated successfully',
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      cy.visit('/profile/edit');

      if (testCase.file) {
        cy.get(`[data-testid="${testCase.field}-input"]`).selectFile(
          `cypress/fixtures/${testCase.file}`
        );
      } else {
        cy.get(`[data-testid="${testCase.field}-input"]`).clear().type(testCase.value);
      }

      cy.get('[data-testid="save-button"]').click();

      cy.contains(testCase.expectedMessage).should('be.visible');
    });
  });

  const errorTestCases = [
    {
      name: 'rejects name that is too long',
      field: 'name',
      value: 'a'.repeat(101),
      expectedError: 'Name must be less than 100 characters',
    },
    {
      name: 'rejects empty name',
      field: 'name',
      value: '',
      expectedError: 'Name is required',
    },
    {
      name: 'rejects invalid file type for avatar',
      field: 'avatar',
      file: 'document.pdf',
      expectedError: 'Please upload an image file',
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      cy.visit('/profile/edit');

      if (testCase.file) {
        cy.get(`[data-testid="${testCase.field}-input"]`).selectFile(
          `cypress/fixtures/${testCase.file}`,
          { force: true }
        );
      } else {
        cy.get(`[data-testid="${testCase.field}-input"]`).clear();
        if (testCase.value) {
          cy.get(`[data-testid="${testCase.field}-input"]`).type(testCase.value);
        }
      }

      cy.get('[data-testid="save-button"]').click();

      cy.contains(testCase.expectedError).should('be.visible');
    });
  });
});

// cypress/e2e/api-integration.cy.ts
describe('API Integration', () => {
  const successTestCases = [
    {
      name: 'fetches user data successfully',
      endpoint: '/api/users/123',
      method: 'GET',
      expectedStatus: 200,
      expectedBody: { id: '123', name: 'John Doe' },
    },
    {
      name: 'creates new post successfully',
      endpoint: '/api/posts',
      method: 'POST',
      body: { title: 'Test Post', content: 'Test content' },
      expectedStatus: 201,
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      cy.request({
        method: testCase.method,
        url: testCase.endpoint,
        body: testCase.body,
      }).then((response) => {
        expect(response.status).to.eq(testCase.expectedStatus);

        if (testCase.expectedBody) {
          expect(response.body).to.include(testCase.expectedBody);
        }
      });
    });
  });

  const errorTestCases = [
    {
      name: 'returns 404 for non-existent user',
      endpoint: '/api/users/999',
      method: 'GET',
      expectedStatus: 404,
    },
    {
      name: 'returns 400 for invalid post data',
      endpoint: '/api/posts',
      method: 'POST',
      body: { title: '' }, // Invalid: empty title
      expectedStatus: 400,
    },
    {
      name: 'returns 401 for unauthorized request',
      endpoint: '/api/private-data',
      method: 'GET',
      expectedStatus: 401,
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      cy.request({
        method: testCase.method,
        url: testCase.endpoint,
        body: testCase.body,
        failOnStatusCode: false, // Don't fail on 4xx/5xx
      }).then((response) => {
        expect(response.status).to.eq(testCase.expectedStatus);
      });
    });
  });
});

// cypress/e2e/api-mocking.cy.ts - Using cy.intercept() for API mocking
describe('API Mocking with cy.intercept()', () => {
  beforeEach(() => {
    // Load mock data from fixtures
    cy.fixture('users.json').as('usersData');
  });

  const successTestCases = [
    {
      name: 'displays user data from mocked API',
      route: '/api/users/*',
      fixture: 'users.json',
      expectedText: 'John Doe',
    },
    {
      name: 'handles loading state correctly',
      route: '/api/posts/*',
      delay: 1000,
      fixture: 'posts.json',
    },
  ];

  successTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      // Mock API response
      cy.intercept('GET', testCase.route, {
        fixture: testCase.fixture,
        delay: testCase.delay || 0,
      }).as('apiCall');

      cy.visit('/users');

      if (testCase.delay) {
        // Verify loading state appears
        cy.get('[data-testid="loading-spinner"]').should('be.visible');
      }

      cy.wait('@apiCall');

      if (testCase.expectedText) {
        cy.contains(testCase.expectedText).should('be.visible');
      }
    });
  });

  const errorTestCases = [
    {
      name: 'displays error message on API failure',
      route: '/api/users/*',
      statusCode: 500,
      expectedError: 'Failed to load users',
    },
    {
      name: 'handles network timeout gracefully',
      route: '/api/posts/*',
      forceNetworkError: true,
      expectedError: 'Network error',
    },
  ];

  errorTestCases.forEach((testCase) => {
    it(testCase.name, () => {
      // Mock API error
      cy.intercept('GET', testCase.route, {
        statusCode: testCase.statusCode || 500,
        forceNetworkError: testCase.forceNetworkError,
      }).as('apiError');

      cy.visit('/users');

      cy.wait('@apiError');

      cy.contains(testCase.expectedError).should('be.visible');
    });
  });
});
```

**When NOT to use table-driven tests:**
- Single test case with complex mocking or setup
- Visual snapshot tests
- Tests where the table structure becomes too complex
- Tests with very different setup/teardown requirements

**Best Practices:**
1. **Always separate success and error test cases** into different arrays
2. Use descriptive test case names with the `name` field
3. Keep test cases focused and related
4. Use objects for test cases (easier to understand and maintain)
5. Mock external dependencies consistently across test cases
6. In Cypress tests, use `data-testid` attributes for reliable element selection
7. For API tests, use `failOnStatusCode: false` when testing error responses

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Verify pnpm is the package manager (check package.json)
- Identify App Router vs Pages Router
- Extract component patterns
- Note styling approach
- Check for Data Access Layer pattern

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into components
- Identify server vs client components (default to server)
- Plan for simplicity and reusability (DRY principle)
- Consider updating existing code rather than creating new

**Step 3: Read Existing Code**
- Review similar pages/components
- Understand data fetching patterns
- Note routing conventions
- Check for existing Server Actions
- Identify test patterns

**Step 4: Implement Changes**
- Follow patterns from design.md
- Use Server Components by default (Next.js 15+ with React 19+)
- Add 'use client' only when needed (interactivity, hooks, browser APIs)
- Create separate actions.ts files for Server Actions
- Use Data Access Layer for database operations
- Implement proper TypeScript types
- Add data-testid attributes for testing
- Prefer early returns over nested conditionals

**Step 5: Write Tests**
- Use table-driven tests with separate success/error cases
- Add data-testid attributes to components
- Write Cypress E2E tests for user flows
- Use cy.intercept() for API mocking
- Store mock data in fixtures
- Ensure test independence

**Step 6: Ensure Quality and Fix All Issues**

**CRITICAL: All quality checks must pass before considering the task complete.**

1. **Run type checking**
   ```bash
   pnpm type-check
   # or
   pnpm tsc --noEmit
   ```
   - If type errors exist, fix them immediately
   - Do not proceed until type checking passes

2. **Run Next.js linter**
   ```bash
   pnpm lint
   ```
   - If lint errors exist, try auto-fix first: `pnpm lint --fix`
   - Fix any remaining issues manually
   - Do not proceed until linter passes

3. **Build the project**
   ```bash
   pnpm build
   ```
   - If build fails, fix the errors immediately
   - Next.js build catches runtime errors and type issues
   - Do not proceed until build succeeds

4. **Run unit/integration tests (if available)**
   ```bash
   pnpm test
   ```
   - If tests fail, fix them immediately
   - Add new tests if implementing new features or components
   - Update existing tests if modifying behavior
   - Do not proceed until all tests pass

5. **Run Cypress E2E tests (if available)**
   ```bash
   pnpm cypress:run
   # or
   pnpm test:e2e
   ```
   - If tests fail, fix them immediately
   - Add new E2E tests for new user flows
   - Update existing tests if user flows changed
   - Do not proceed until all tests pass

6. **Manual browser testing**
   - Start dev server: `pnpm dev`
   - Test the feature in browser
   - Check for console errors
   - Verify responsive behavior
   - Test different user interactions

7. **Update tests for your changes**
   - If you added a new component, add corresponding tests
   - If you modified behavior, update existing tests
   - Add data-testid attributes for Cypress testing
   - Ensure test coverage is maintained or improved

**Iterative Fix Process:**
- If any step fails, fix the issues and re-run ALL previous steps
- Continue iterating until ALL checks pass:
  - ✅ No type errors
  - ✅ No lint issues
  - ✅ Build succeeds (very important for Next.js)
  - ✅ All unit/integration tests pass
  - ✅ All Cypress E2E tests pass
  - ✅ No console errors in browser
  - ✅ Feature works as expected

## Commands Reference

**IMPORTANT: Use pnpm instead of npm**

```bash
# Setup new project
pnpm create next-app@latest

# Install dependencies
pnpm i

# Development
pnpm dev

# Type checking
pnpm type-check
pnpm tsc --noEmit

# Linting
pnpm lint
pnpm lint --fix

# Build (CRITICAL - catches many issues)
pnpm build

# Production
pnpm start

# Unit/Integration Tests
pnpm test
pnpm test:watch

# E2E Tests with Cypress
pnpm cypress:open     # Interactive mode
pnpm cypress:run      # Headless mode
pnpm test:e2e         # Run with server (if configured with start-server-and-test)
```

## Common Pitfalls to Avoid

### Client Component Overuse
```typescript
// Bad: Unnecessary client component
'use client';
export default function Page() {
  const [data, setData] = useState(null);
  useEffect(() => {
    fetch('/api/data').then(r => r.json()).then(setData);
  }, []);
  // ...
}

// Good: Server component with direct data fetching
export default async function Page() {
  const data = await fetchData();
  return <ClientComponent data={data} />;
}
```

### Missing Loading States
```typescript
// Bad: No loading UI
export default async function Page() {
  const data = await fetchData();
  return <div>{data}</div>;
}

// Good: With loading.tsx
// app/loading.tsx
export default function Loading() {
  return <div>Loading...</div>;
}

// app/page.tsx
export default async function Page() {
  const data = await fetchData();
  return <div>{data}</div>;
}
```

### Improper Image Usage
```typescript
// Bad: Regular img tag
<img src="/photo.jpg" alt="Photo" />

// Good: Next.js Image component
<Image src="/photo.jpg" alt="Photo" width={500} height={300} />
```

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Identify App Router vs Pages Router
- [ ] Extract project-specific conventions
- [ ] Review existing component patterns
- [ ] Check styling approach

### During Implementation
- [ ] Use Server Components by default
- [ ] Add 'use client' only when needed
- [ ] Implement proper TypeScript types
- [ ] Use Next.js optimization features (Image, Link)
- [ ] Follow existing naming conventions
- [ ] Add loading and error states
- [ ] Write table-driven tests alongside code (separate success/error cases)

### After Implementation - MUST ALL PASS
- [ ] Run type checking (`pnpm type-check`) - **FIX ALL TYPE ERRORS**
- [ ] Run linters (`pnpm lint`) - **FIX ALL LINT ISSUES**
- [ ] Build successfully (`pnpm build`) - **MUST SUCCEED** (critical for Next.js)
- [ ] Run unit tests if available (`pnpm test`) - **ALL TESTS MUST PASS**
- [ ] Run Cypress E2E tests if available (`pnpm cypress:run` or `pnpm test:e2e`) - **ALL TESTS MUST PASS**
- [ ] Add/update table-driven tests (separate success/error cases) for new or modified components
- [ ] Add data-testid attributes for Cypress testing
- [ ] Test in browser - **NO CONSOLE ERRORS**
- [ ] Verify responsive design
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation

## Key Principles (2025)

1. **Package Manager**: Always use pnpm (MUST NOT use legacy-peer-deps)
2. **Project Guidelines First**: Always read and follow `.claude/design.md`
3. **Simplicity & DRY**: Update existing code for reusability, avoid creating new when possible
4. **Server First**: Use Server Components by default (Next.js 15+ with React 19+)
5. **Security**: Use Server Actions for mutations, Data Access Layer for database
6. **Type Safety**: Use TypeScript with Zod for runtime validation
7. **User Experience**: Implement loading and error states
8. **Optimization**: Leverage Next.js built-in optimizations
9. **Performance**: Code split heavy components
10. **Fail-Fast**: Early returns, explicit error handling
11. **Testing**: Table-driven tests with Cypress E2E, separate success/error cases

## Version History

### Version 2.0 (2025-01-16)
- Updated for 2025 best practices (Next.js 15+, React 19+)
- Added mandatory pnpm package manager requirement
- Enhanced Server Actions guidelines with security best practices
- Added Data Access Layer pattern
- Expanded Cypress testing section with setup, configuration, and examples
- Added API mocking with cy.intercept() examples
- Emphasized data-testid attributes for testing
- Incorporated coding guidelines: simplicity, DRY, fail-fast, early returns
- Added testing principles: avoid redundant tests, injection over global state
- Enhanced Type Safety section for API responses
- Clarified Server Actions vs Route Handlers usage
- Added Core Principles section (simplicity, consistency, latest versions, encapsulation)
- Updated all commands to use pnpm instead of npm
- Added Node.js version requirement (20.9+)

### Version 1.0 (Previous)
- Initial Next.js skill with App Router and Pages Router support
- Basic Server Components and Client Components guidance
- Server Actions examples
- Table-driven testing pattern
- Jest/Vitest examples
