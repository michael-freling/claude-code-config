# TypeScript/React Testing Patterns

## Pattern Purposes

| Pattern | Purpose |
|---------|---------|
| **Dependency Injection** | Makes components testable by allowing tests to control dependencies instead of components fetching their own |
| **Table-Driven Tests** | Reduces duplication by defining test cases as data; easier to add cases and see all scenarios |
| **Mock Context Providers** | Wraps components with controlled context values so tests don't depend on real auth, API state, etc. |
| **Test Data Builders** | Creates consistent test objects with sensible defaults; avoids repeating object construction |
| **Mocking API Calls** | Intercepts network requests to return controlled responses; tests don't depend on real APIs |

## Dependency Injection

```typescript
// Bad: Hard to test - fetches its own data internally
function ProductList() {
    const products = useProducts()  // Can't control in tests
    const { user } = useAuth()      // Can't inject test user
    return <div>{products.map(p => <Product key={p.id} product={p} />)}</div>
}

// Good: Testable - accepts data as props or injectable hooks
function ProductList({
    useFetchProducts,
    canEdit,
}: {
    useFetchProducts: (cursor: string) => { products: Product[], nextCursor: string }
    canEdit: (product: Product) => boolean
}) {
    const { products, nextCursor } = useFetchProducts(cursor)
    return <div>{products.map(p => <Product key={p.id} product={p} editable={canEdit(p)} />)}</div>
}
```

## Test Data Builders

```typescript
export function createUser(fields?: Partial<User>): User {
    return {
        id: '1',
        name: 'Test User',
        email: 'test@example.com',
        avatarUrl: '/images/default-avatar.png',
        ...fields,
    }
}

export function createProduct(fields?: Partial<Product>): Product {
    return {
        id: '1',
        name: 'Test Product',
        price: 1000,
        ...fields,
    }
}
```

## Mock Context Providers

```typescript
// Create dedicated mock providers in src/testing/providers.tsx
interface MockAuthProviderProps {
    user: User | null
    children: React.ReactNode
}

export function MockAuthProvider({ user, children }: MockAuthProviderProps) {
    return (
        <AuthContext.Provider
            value={{
                isLoading: false,
                user,
                signIn: async () => {},
                signOut: async () => {},
            }}
        >
            {children}
        </AuthContext.Provider>
    )
}
```

## Table-Driven Tests

### Jest / Vitest + React Testing Library

```typescript
import { render, screen } from '@testing-library/react'

describe('NavBar', () => {
    interface TestCase {
        name: string
        user: User | null
        expectLoginButton: boolean
        expectAvatar: boolean
    }

    const testCases: TestCase[] = [
        {
            name: 'shows login button for anonymous user',
            user: null,
            expectLoginButton: true,
            expectAvatar: false,
        },
        {
            name: 'shows avatar for authenticated user',
            user: createUser({ avatarUrl: '/avatar.png' }),
            expectLoginButton: false,
            expectAvatar: true,
        },
    ]

    test.each(testCases)('$name', (tc) => {
        render(
            <MockAuthProvider user={tc.user}>
                <NavBar />
            </MockAuthProvider>
        )

        if (tc.expectLoginButton) {
            expect(screen.getByTestId('login-button')).toBeInTheDocument()
        }
        if (tc.expectAvatar) {
            expect(screen.getByTestId('user-avatar')).toBeInTheDocument()
        }
    })
})
```

### Cypress Component Testing

```typescript
describe('NavBar', () => {
    interface TestCase {
        name: string
        user: User | null
        expectLoginButton: boolean
        expectAvatar: boolean
    }

    const testCases: TestCase[] = [
        {
            name: 'shows login button for anonymous user',
            user: null,
            expectLoginButton: true,
            expectAvatar: false,
        },
        {
            name: 'shows avatar for authenticated user',
            user: createUser({ avatarUrl: '/avatar.png' }),
            expectLoginButton: false,
            expectAvatar: true,
        },
    ]

    testCases.forEach((tc) => {
        it(tc.name, () => {
            cy.mount(
                <MockAuthProvider user={tc.user}>
                    <NavBar />
                </MockAuthProvider>
            )

            if (tc.expectLoginButton) {
                cy.get('[data-testid=login-button]').should('be.visible')
            }
            if (tc.expectAvatar) {
                cy.get('[data-testid=user-avatar]').should('be.visible')
            }
        })
    })
})
```

## Mocking API Calls

### Jest / Vitest with MSW (Mock Service Worker)

```typescript
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

const server = setupServer()

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('ContactForm', () => {
    it('submits form successfully', async () => {
        server.use(
            http.post('/api/contact', () => {
                return HttpResponse.json({ success: true })
            })
        )

        render(
            <MockAuthProvider user={createUser()}>
                <ContactForm />
            </MockAuthProvider>
        )

        await userEvent.type(screen.getByTestId('message-input'), 'Hello')
        await userEvent.click(screen.getByTestId('submit-button'))

        await waitFor(() => {
            expect(screen.getByTestId('success-message')).toBeInTheDocument()
        })
    })
})
```

### Cypress

```typescript
describe('ContactForm', () => {
    it('submits form successfully', () => {
        cy.intercept('POST', '/api/contact', {
            statusCode: 200,
            body: { success: true },
        }).as('submitForm')

        cy.mount(
            <MockAuthProvider user={createUser()}>
                <ContactForm />
            </MockAuthProvider>
        )

        cy.get('[data-testid=message-input]').type('Hello')
        cy.get('[data-testid=submit-button]').click()
        cy.wait('@submitForm')
        cy.get('[data-testid=success-message]').should('be.visible')
    })
})
```
