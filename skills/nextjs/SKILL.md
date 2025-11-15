---
name: nextjs
description: Next.js feature development following best practices for App Router, Server Components, Server Actions, Pages Router, routing, and optimization. Use when adding or updating Next.js applications.
---

# Next.js Development Skill

A comprehensive skill for adding or updating features in Next.js projects following best practices.

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
- Identify if using App Router (app/) or Pages Router (pages/)
- Review existing component patterns
- Check for styling approach (CSS modules, Tailwind, styled-components)
- Identify state management (if any)
- Review data fetching patterns

### 2. Search for Relevant Code

- Search for similar components or pages
- Look for related API routes or server actions
- Find existing hooks or utilities
- Identify styling patterns

## Next.js Best Practices

### App Router (Next.js 13+)

**Server Components (Default):**
- Use Server Components by default
- Add 'use client' only when needed (interactivity, browser APIs)
- Leverage async/await in Server Components
- Direct data fetching without useEffect

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

**Server Actions:**
- Define server actions in separate files or at top of server components
- Use proper form validation
- Return structured data with success/error states
- Use revalidatePath/revalidateTag for cache updates

Example:
```typescript
// app/actions/user.ts
'use server';

import { revalidatePath } from 'next/cache';
import { z } from 'zod';

const createUserSchema = z.object({
  name: z.string().min(1).max(100),
  email: z.string().email(),
});

export async function createUser(formData: FormData) {
  const parsed = createUserSchema.safeParse({
    name: formData.get('name'),
    email: formData.get('email'),
  });

  if (!parsed.success) {
    return { success: false, errors: parsed.error.flatten() };
  }

  try {
    await db.user.create({ data: parsed.data });
    revalidatePath('/users');
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to create user' };
  }
}
```

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

### API Routes

**App Router API Routes:**
```typescript
// app/api/users/[id]/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { z } from 'zod';

const updateUserSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  email: z.string().email().optional(),
});

export async function PATCH(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const body = await request.json();
    const validated = updateUserSchema.parse(body);

    const user = await db.user.update({
      where: { id: params.id },
      data: validated,
    });

    return NextResponse.json({ user, message: 'User updated' });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return NextResponse.json(
        { error: 'Validation failed', details: error.errors },
        { status: 400 }
      );
    }
    return NextResponse.json(
      { error: 'Failed to update user' },
      { status: 500 }
    );
  }
}
```

### Testing

**CRITICAL: Always use table-driven tests (test.each) as the primary testing approach.**

Table-driven tests provide the same benefits in Next.js projects:
- Reduce code duplication
- Make it easy to add new test cases
- Improve test readability
- Ensure consistent test structure
- Make test coverage gaps obvious

**Use table-driven tests for:**
- Component rendering with different props
- Utility functions with multiple test cases
- Server Actions with various inputs
- API route handlers
- Edge cases and error scenarios

**Table-Driven Test Examples for Next.js:**

```typescript
// Good: Testing component rendering with different props
import { render, screen } from '@testing-library/react';
import { Button } from './Button';

describe('Button', () => {
  test.each([
    { variant: 'primary', expectedClass: 'btn-primary' },
    { variant: 'secondary', expectedClass: 'btn-secondary' },
    { variant: 'danger', expectedClass: 'btn-danger' },
  ])('renders with $variant variant', ({ variant, expectedClass }) => {
    render(<Button variant={variant}>Click me</Button>);
    const button = screen.getByRole('button');
    expect(button).toHaveClass(expectedClass);
  });
});

// Good: Testing server actions
import { createUser } from './actions';

describe('createUser', () => {
  test.each([
    {
      name: 'valid user data',
      formData: new FormData([['name', 'John'], ['email', 'john@example.com']]),
      expected: { success: true },
    },
    {
      name: 'missing name',
      formData: new FormData([['email', 'john@example.com']]),
      expected: { success: false, error: expect.any(Object) },
    },
    {
      name: 'invalid email',
      formData: new FormData([['name', 'John'], ['email', 'invalid']]),
      expected: { success: false, error: expect.any(Object) },
    },
  ])('$name', async ({ formData, expected }) => {
    const result = await createUser(formData);
    expect(result).toMatchObject(expected);
  });
});

// Good: Testing utility functions
import { formatPrice } from '@/lib/utils';

describe('formatPrice', () => {
  test.each([
    { amount: 1000, currency: 'USD', expected: '$10.00' },
    { amount: 2550, currency: 'USD', expected: '$25.50' },
    { amount: 0, currency: 'USD', expected: '$0.00' },
    { amount: 1000, currency: 'EUR', expected: '€10.00' },
  ])('formats $amount cents as $expected', ({ amount, currency, expected }) => {
    expect(formatPrice(amount, currency)).toBe(expected);
  });
});

// Good: Testing API routes
import { GET, POST } from '@/app/api/users/route';
import { NextRequest } from 'next/server';

describe('/api/users', () => {
  describe('POST', () => {
    test.each([
      {
        name: 'creates user successfully',
        body: { name: 'John Doe', email: 'john@example.com' },
        expectedStatus: 200,
        expectedBody: { success: true },
      },
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
    ])('$name', async ({ body, expectedStatus, expectedBody }) => {
      const request = new NextRequest('http://localhost/api/users', {
        method: 'POST',
        body: JSON.stringify(body),
      });

      const response = await POST(request);
      expect(response.status).toBe(expectedStatus);

      const data = await response.json();
      expect(data).toMatchObject(expectedBody);
    });
  });
});

// Good: Testing React hooks
import { renderHook, act } from '@testing-library/react';
import { useCounter } from '@/hooks/useCounter';

describe('useCounter', () => {
  test.each([
    { initialValue: 0, incrementBy: 1, expected: 1 },
    { initialValue: 5, incrementBy: 2, expected: 7 },
    { initialValue: 10, incrementBy: -3, expected: 7 },
  ])('increments from $initialValue by $incrementBy', ({ initialValue, incrementBy, expected }) => {
    const { result } = renderHook(() => useCounter(initialValue));

    act(() => {
      result.current.increment(incrementBy);
    });

    expect(result.current.count).toBe(expected);
  });
});

// Good: Testing with React Testing Library user interactions
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SearchForm } from './SearchForm';

describe('SearchForm', () => {
  test.each([
    {
      name: 'submits valid search query',
      input: 'test query',
      shouldSubmit: true,
    },
    {
      name: 'does not submit empty query',
      input: '',
      shouldSubmit: false,
    },
    {
      name: 'trims whitespace from query',
      input: '  test query  ',
      shouldSubmit: true,
      expectedValue: 'test query',
    },
  ])('$name', async ({ input, shouldSubmit, expectedValue }) => {
    const onSubmit = jest.fn();
    const user = userEvent.setup();

    render(<SearchForm onSubmit={onSubmit} />);

    const searchInput = screen.getByRole('textbox');
    const submitButton = screen.getByRole('button');

    await user.type(searchInput, input);
    await user.click(submitButton);

    if (shouldSubmit) {
      expect(onSubmit).toHaveBeenCalledWith(expectedValue || input);
    } else {
      expect(onSubmit).not.toHaveBeenCalled();
    }
  });
});

// Good: Testing metadata generation
import { generateMetadata } from '@/app/users/[id]/page';

describe('generateMetadata', () => {
  test.each([
    {
      name: 'generates metadata for existing user',
      params: { id: '123' },
      mockUser: { id: '123', name: 'John Doe', bio: 'Developer' },
      expected: {
        title: 'John Doe - Profile',
        description: 'Developer',
      },
    },
    {
      name: 'generates metadata for non-existent user',
      params: { id: '999' },
      mockUser: null,
      expected: {
        title: 'User Not Found',
        description: undefined,
      },
    },
  ])('$name', async ({ params, mockUser, expected }) => {
    // Mock the getUserById function
    jest.mock('@/lib/api', () => ({
      getUserById: jest.fn().mockResolvedValue(mockUser),
    }));

    const metadata = await generateMetadata({ params });
    expect(metadata).toMatchObject(expected);
  });
});
```

**When NOT to use table-driven tests:**
- Single test case with complex mocking
- Visual snapshot tests
- E2E tests with Playwright/Cypress
- Tests where the table structure becomes too complex

**Best Practices:**
1. Use descriptive test case names (use the `name` field or `$variable` syntax)
2. Keep test cases focused and related
3. Use objects for complex test cases, arrays for simple ones
4. Mock external dependencies consistently across test cases
5. Test both successful and error scenarios
6. Consider using React Testing Library for component tests

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Identify App Router vs Pages Router
- Extract component patterns
- Note styling approach

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into components
- Identify server vs client components

**Step 3: Read Existing Code**
- Review similar pages/components
- Understand data fetching patterns
- Note routing conventions

**Step 4: Implement Changes**
- Follow patterns from design.md
- Use Server Components by default
- Add 'use client' only when needed
- Implement proper types

**Step 5: Ensure Quality and Fix All Issues**

**CRITICAL: All quality checks must pass before considering the task complete.**

1. **Run type checking**
   ```bash
   npm run type-check
   # or
   npx tsc --noEmit
   ```
   - If type errors exist, fix them immediately
   - Do not proceed until type checking passes

2. **Run Next.js linter**
   ```bash
   npm run lint
   ```
   - If lint errors exist, try auto-fix first: `npm run lint -- --fix`
   - Fix any remaining issues manually
   - Do not proceed until linter passes

3. **Build the project**
   ```bash
   npm run build
   ```
   - If build fails, fix the errors immediately
   - Next.js build catches runtime errors and type issues
   - Do not proceed until build succeeds

4. **Run tests (if available)**
   ```bash
   npm test
   ```
   - If tests fail, fix them immediately
   - Add new tests if implementing new features or components
   - Update existing tests if modifying behavior
   - Do not proceed until all tests pass

5. **Manual browser testing**
   - Start dev server: `npm run dev`
   - Test the feature in browser
   - Check for console errors
   - Verify responsive behavior
   - Test different user interactions

6. **Update tests for your changes**
   - If you added a new component, add corresponding tests
   - If you modified behavior, update existing tests
   - Ensure test coverage is maintained or improved

**Iterative Fix Process:**
- If any step fails, fix the issues and re-run ALL previous steps
- Continue iterating until ALL checks pass:
  - ✅ No type errors
  - ✅ No lint issues
  - ✅ Build succeeds (very important for Next.js)
  - ✅ All tests pass
  - ✅ No console errors in browser
  - ✅ Feature works as expected

## Commands Reference

```bash
# Development
npm run dev

# Type checking
npm run type-check
npx tsc --noEmit

# Linting
npm run lint
npm run lint -- --fix

# Build (CRITICAL - catches many issues)
npm run build

# Production
npm start

# Tests
npm test
npm run test:watch
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
- [ ] Write table-driven tests alongside code

### After Implementation - MUST ALL PASS
- [ ] Run type checking (`npm run type-check`) - **FIX ALL TYPE ERRORS**
- [ ] Run linters (`npm run lint`) - **FIX ALL LINT ISSUES**
- [ ] Build successfully (`npm run build`) - **MUST SUCCEED** (critical for Next.js)
- [ ] Run tests if available (`npm test`) - **ALL TESTS MUST PASS**
- [ ] Add/update table-driven tests for new or modified components
- [ ] Test in browser - **NO CONSOLE ERRORS**
- [ ] Verify responsive design
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Server First**: Use Server Components by default
3. **Optimization**: Leverage Next.js built-in optimizations
4. **Type Safety**: Use TypeScript for all components
5. **User Experience**: Implement loading and error states
6. **Performance**: Code split heavy components
7. **Testing**: Use table-driven tests to ensure components work correctly
