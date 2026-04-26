# React Component Patterns

## Component with Typed Props and Children

```tsx
import { type ReactNode } from "react";

// Domain types
type UserId = string & { readonly __brand: unique symbol };

interface User {
  id: UserId;
  name: string;
  email: string;
  avatarUrl: string;
}

// Props interface — named after the component, not generic
interface UserCardProps {
  user: User;
  onSelect: (userId: UserId) => void;
  children?: ReactNode;
}

function UserCard({ user, onSelect, children }: UserCardProps) {
  if (!user.name) {
    return null;
  }

  function handleClick() {
    onSelect(user.id);
  }

  return (
    <article role="article" aria-label={`User card for ${user.name}`}>
      <img src={user.avatarUrl} alt={`${user.name}'s avatar`} />
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      {children}
      <button type="button" onClick={handleClick}>
        Select
      </button>
    </article>
  );
}
```

## Custom Hook with Dependency Injection

### Cancelled flag pattern

```tsx
// Discriminated union for async state
type AsyncState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "error"; error: Error }
  | { status: "success"; data: T };

// Client interface — the hook depends on an abstraction, not a concrete implementation
interface UserClient {
  fetchUsers(query: string): Promise<User[]>;
}

function useFetchUsers(client: UserClient, query: string): AsyncState<User[]> {
  const [state, setState] = useState<AsyncState<User[]>>({ status: "idle" });

  useEffect(() => {
    if (!query) {
      setState({ status: "idle" });
      return;
    }

    let cancelled = false;
    setState({ status: "loading" });

    client
      .fetchUsers(query)
      .then((data) => {
        if (!cancelled) {
          setState({ status: "success", data });
        }
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setState({
            status: "error",
            error: error instanceof Error ? error : new Error(String(error)),
          });
        }
      });

    return () => {
      cancelled = true;
    };
  }, [client, query]);

  return state;
}
```

### AbortController pattern

Use when the underlying API supports cancellation via `AbortSignal`, so in-flight
network requests are actually cancelled rather than just ignored on completion.

```tsx
interface UserClientWithAbort {
  fetchUsers(query: string, signal: AbortSignal): Promise<User[]>;
}

function useFetchUsersWithAbort(
  client: UserClientWithAbort,
  query: string,
): AsyncState<User[]> {
  const [state, setState] = useState<AsyncState<User[]>>({ status: "idle" });

  useEffect(() => {
    if (!query) {
      setState({ status: "idle" });
      return;
    }

    const controller = new AbortController();
    setState({ status: "loading" });

    client
      .fetchUsers(query, controller.signal)
      .then((data) => {
        if (!controller.signal.aborted) {
          setState({ status: "success", data });
        }
      })
      .catch((error: unknown) => {
        if (controller.signal.aborted) {
          return; // expected cancellation, not an error
        }
        setState({
          status: "error",
          error: error instanceof Error ? error : new Error(String(error)),
        });
      });

    return () => {
      controller.abort();
    };
  }, [client, query]);

  return state;
}
```

### Usage — consumer provides the client, tests can substitute a fake

```tsx
function UserSearch({ client }: { client: UserClient }) {
  const [query, setQuery] = useState("");
  const state = useFetchUsers(client, query);

  return (
    <div>
      <input
        type="search"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        aria-label="Search users"
      />
      {state.status === "loading" && <p>Loading...</p>}
      {state.status === "error" && <p role="alert">{state.error.message}</p>}
      {state.status === "success" && (
        <ul>
          {state.data.map((user) => (
            <li key={user.id}>{user.name}</li>
          ))}
        </ul>
      )}
    </div>
  );
}
```

## Controlled Form with Zod Validation

```tsx
import { useState, type FormEvent } from "react";
import { z } from "zod";

const contactFormSchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Name is too long"),
  email: z.string().email("Invalid email address"),
  message: z
    .string()
    .min(10, "Message must be at least 10 characters")
    .max(1000, "Message is too long"),
});

type ContactFormData = z.infer<typeof contactFormSchema>;

// Per-field error map derived from the schema
type FieldErrors = Partial<Record<keyof ContactFormData, string>>;

function parseFieldErrors(error: z.ZodError<ContactFormData>): FieldErrors {
  const errors: FieldErrors = {};
  for (const issue of error.issues) {
    const field = issue.path[0] as keyof ContactFormData | undefined;
    if (field && !errors[field]) {
      errors[field] = issue.message;
    }
  }
  return errors;
}

interface ContactFormProps {
  onSubmit: (data: ContactFormData) => Promise<void>;
}

function ContactForm({ onSubmit }: ContactFormProps) {
  const [fields, setFields] = useState<ContactFormData>({
    name: "",
    email: "",
    message: "",
  });
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [submitting, setSubmitting] = useState(false);

  function handleChange(
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) {
    const { name, value } = e.target;
    setFields((prev) => ({ ...prev, [name]: value }));
    // Clear the error for the field being edited
    setFieldErrors((prev) => ({ ...prev, [name]: undefined }));
  }

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const result = contactFormSchema.safeParse(fields);
    if (!result.success) {
      setFieldErrors(parseFieldErrors(result.error));
      return;
    }

    setSubmitting(true);
    try {
      await onSubmit(result.data);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} noValidate>
      <div>
        <label htmlFor="name">Name</label>
        <input
          id="name"
          name="name"
          value={fields.name}
          onChange={handleChange}
          aria-invalid={fieldErrors.name ? true : undefined}
          aria-describedby={fieldErrors.name ? "name-error" : undefined}
        />
        {fieldErrors.name && (
          <p id="name-error" role="alert">
            {fieldErrors.name}
          </p>
        )}
      </div>

      <div>
        <label htmlFor="email">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          value={fields.email}
          onChange={handleChange}
          aria-invalid={fieldErrors.email ? true : undefined}
          aria-describedby={fieldErrors.email ? "email-error" : undefined}
        />
        {fieldErrors.email && (
          <p id="email-error" role="alert">
            {fieldErrors.email}
          </p>
        )}
      </div>

      <div>
        <label htmlFor="message">Message</label>
        <textarea
          id="message"
          name="message"
          value={fields.message}
          onChange={handleChange}
          aria-invalid={fieldErrors.message ? true : undefined}
          aria-describedby={fieldErrors.message ? "message-error" : undefined}
        />
        {fieldErrors.message && (
          <p id="message-error" role="alert">
            {fieldErrors.message}
          </p>
        )}
      </div>

      <button type="submit" disabled={submitting}>
        {submitting ? "Sending..." : "Send"}
      </button>
    </form>
  );
}
```

## Component Composition over Boolean Props

```tsx
import { type ReactNode } from "react";

// Wrong: branching render logic driven by a type prop
// function Alert({ type, children }: { type: "error" | "warning" | "info"; children: ReactNode }) {
//   const icon = type === "error" ? "X" : type === "warning" ? "!" : "i";
//   const className = `alert alert-${type}`;
//   return <div className={className}>{icon} {children}</div>;
// }

// Right: shared base with composition — each variant is its own component

interface AlertBaseProps {
  icon: ReactNode;
  className: string;
  role: "alert" | "status";
  children: ReactNode;
}

function AlertBase({ icon, className, role, children }: AlertBaseProps) {
  return (
    <div className={className} role={role}>
      <span aria-hidden="true">{icon}</span>
      <div>{children}</div>
    </div>
  );
}

function ErrorAlert({ children }: { children: ReactNode }) {
  return (
    <AlertBase icon="X" className="alert alert-error" role="alert">
      {children}
    </AlertBase>
  );
}

function WarningAlert({ children }: { children: ReactNode }) {
  return (
    <AlertBase icon="!" className="alert alert-warning" role="alert">
      {children}
    </AlertBase>
  );
}

function InfoAlert({ children }: { children: ReactNode }) {
  return (
    <AlertBase icon="i" className="alert alert-info" role="status">
      {children}
    </AlertBase>
  );
}
```

## Performance: memo, useMemo, useCallback

```tsx
import { memo, useMemo, useCallback, useState } from "react";

// Expensive child — wrap in React.memo so it only re-renders when its props change
interface ProductRowProps {
  product: Product;
  onAddToCart: (productId: string) => void;
}

const ProductRow = memo(function ProductRow({
  product,
  onAddToCart,
}: ProductRowProps) {
  return (
    <tr>
      <td>{product.name}</td>
      <td>${(product.price / 100).toFixed(2)}</td>
      <td>
        <button type="button" onClick={() => onAddToCart(product.id)}>
          Add to cart
        </button>
      </td>
    </tr>
  );
});

interface Product {
  id: string;
  name: string;
  price: number; // cents
  category: string;
}

function ProductTable({ products }: { products: Product[] }) {
  const [filter, setFilter] = useState("");
  const [cart, setCart] = useState<string[]>([]);

  // useMemo — recompute filtered list only when products or filter change,
  // not on every render caused by unrelated state (e.g., cart updates)
  const filtered = useMemo(() => {
    const lower = filter.toLowerCase();
    return products.filter((p) => p.name.toLowerCase().includes(lower));
  }, [products, filter]);

  // useCallback — stable function reference so ProductRow (wrapped in memo) does
  // not re-render when unrelated parent state changes
  const handleAddToCart = useCallback((productId: string) => {
    setCart((prev) => [...prev, productId]);
  }, []);

  return (
    <div>
      <input
        type="search"
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
        placeholder="Filter products"
        aria-label="Filter products"
      />
      <p>{cart.length} items in cart</p>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Price</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((product) => (
            <ProductRow
              key={product.id}
              product={product}
              onAddToCart={handleAddToCart}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

## Testing with @testing-library/react

```tsx
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";

// --- Testing UserCard ---

describe("UserCard", () => {
  function createUser(fields?: Partial<User>): User {
    return {
      id: "user-1" as UserId,
      name: "Alice",
      email: "alice@example.com",
      avatarUrl: "/alice.png",
      ...fields,
    };
  }

  it("renders user details and calls onSelect when clicked", async () => {
    const user = createUser();
    const onSelect = vi.fn();

    render(<UserCard user={user} onSelect={onSelect} />);

    expect(screen.getByRole("heading", { name: "Alice" })).toBeInTheDocument();
    expect(screen.getByText("alice@example.com")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "Select" }));
    expect(onSelect).toHaveBeenCalledWith(user.id);
  });

  it("returns null when user name is empty", () => {
    const { container } = render(
      <UserCard user={createUser({ name: "" })} onSelect={vi.fn()} />,
    );

    expect(container.firstChild).toBeNull();
  });
});

// --- Testing ContactForm ---

describe("ContactForm", () => {
  it("shows validation errors for empty fields on submit", async () => {
    const onSubmit = vi.fn();

    render(<ContactForm onSubmit={onSubmit} />);

    await userEvent.click(screen.getByRole("button", { name: "Send" }));

    expect(screen.getByText("Name is required")).toBeInTheDocument();
    expect(screen.getByText("Invalid email address")).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("submits valid form data", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<ContactForm onSubmit={onSubmit} />);

    await userEvent.type(screen.getByLabelText("Name"), "Alice");
    await userEvent.type(screen.getByLabelText("Email"), "alice@example.com");
    await userEvent.type(
      screen.getByLabelText("Message"),
      "Hello, this is a test message.",
    );
    await userEvent.click(screen.getByRole("button", { name: "Send" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Alice",
        email: "alice@example.com",
        message: "Hello, this is a test message.",
      });
    });
  });

  it("clears field error when user edits the field", async () => {
    render(<ContactForm onSubmit={vi.fn()} />);

    // Trigger validation
    await userEvent.click(screen.getByRole("button", { name: "Send" }));
    expect(screen.getByText("Name is required")).toBeInTheDocument();

    // Start typing in the name field
    await userEvent.type(screen.getByLabelText("Name"), "A");
    expect(screen.queryByText("Name is required")).not.toBeInTheDocument();
  });
});

// --- Testing useFetchUsers with a fake client ---

describe("UserSearch", () => {
  function createFakeClient(
    response: User[] | Error,
  ): UserClient {
    return {
      fetchUsers: response instanceof Error
        ? vi.fn().mockRejectedValue(response)
        : vi.fn().mockResolvedValue(response),
    };
  }

  it("displays fetched users", async () => {
    const users: User[] = [
      {
        id: "1" as UserId,
        name: "Alice",
        email: "a@b.com",
        avatarUrl: "/a.png",
      },
    ];
    const client = createFakeClient(users);

    render(<UserSearch client={client} />);

    await userEvent.type(screen.getByLabelText("Search users"), "Alice");

    await waitFor(() => {
      expect(screen.getByText("Alice")).toBeInTheDocument();
    });
  });

  it("displays error message on failure", async () => {
    const client = createFakeClient(new Error("Network failure"));

    render(<UserSearch client={client} />);

    await userEvent.type(screen.getByLabelText("Search users"), "Alice");

    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent("Network failure");
    });
  });
});
```
