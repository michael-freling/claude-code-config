# React Component Patterns

## Functional Component with Props Interface

```typescript
interface UserCardProps {
  user: User;
  onEdit: (userId: string) => void;
}

function UserCard({ user, onEdit }: UserCardProps) {
  return (
    <div className="rounded-lg border p-4">
      <h3 className="text-lg font-semibold">{user.name}</h3>
      <p className="text-sm text-gray-600">{user.email}</p>
      <button
        type="button"
        onClick={() => onEdit(user.id)}
        className="mt-2 rounded bg-blue-600 px-3 py-1 text-white"
      >
        Edit
      </button>
    </div>
  );
}
```

## Component with Loading and Error States

```typescript
function UserProfile({ userId }: { userId: string }) {
  const { data: user, error, isLoading } = useUser(userId);

  if (isLoading) return <UserProfileSkeleton />;
  if (error) return <ErrorMessage error={error} />;
  if (!user) return <NotFound resource="user" />;

  return (
    <div>
      <h1>{user.name}</h1>
      <p>{user.email}</p>
      <OrderList userId={user.id} />
    </div>
  );
}
```

## Custom Hook for Data Fetching

```typescript
interface UseUserReturn {
  data: User | undefined;
  error: Error | undefined;
  isLoading: boolean;
}

function useUser(userId: string): UseUserReturn {
  const [data, setData] = useState<User>();
  const [error, setError] = useState<Error>();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const controller = new AbortController();

    async function fetchUser() {
      setIsLoading(true);
      try {
        const user = await userClient.getById(userId, {
          signal: controller.signal,
        });
        setData(user);
      } catch (err) {
        if (!controller.signal.aborted) {
          setError(err instanceof Error ? err : new Error(String(err)));
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    fetchUser();
    return () => controller.abort();
  }, [userId]);

  return { data, error, isLoading };
}
```

## Composition over Prop Drilling

```typescript
function OrderPage({ orderId }: { orderId: string }) {
  const { data: order } = useOrder(orderId);
  if (!order) return null;

  return (
    <OrderLayout>
      <OrderHeader order={order} />
      <OrderItems items={order.items} />
      <OrderSummary total={order.total} status={order.status} />
    </OrderLayout>
  );
}

function OrderSummary({ total, status }: { total: number; status: OrderStatus }) {
  return (
    <div className="mt-4 border-t pt-4">
      <p className="text-lg font-bold">${total.toFixed(2)}</p>
      <StatusBadge status={status} />
    </div>
  );
}
```

## Form with Controlled State

```typescript
interface CreateUserFormProps {
  onSubmit: (data: CreateUserInput) => Promise<void>;
}

function CreateUserForm({ onSubmit }: CreateUserFormProps) {
  const [formData, setFormData] = useState<CreateUserInput>({
    name: "",
    email: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string>();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setIsSubmitting(true);
    setError(undefined);

    try {
      await onSubmit(formData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <label htmlFor="name">Name</label>
      <input
        id="name"
        value={formData.name}
        onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
        required
      />
      <label htmlFor="email">Email</label>
      <input
        id="email"
        type="email"
        value={formData.email}
        onChange={(e) => setFormData((prev) => ({ ...prev, email: e.target.value }))}
        required
      />
      {error && <p role="alert" className="text-red-600">{error}</p>}
      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Creating..." : "Create User"}
      </button>
    </form>
  );
}
```

## Testing with Testing Library

```typescript
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

describe("UserCard", () => {
  const user: User = {
    id: "user-1",
    name: "Alice",
    email: "alice@example.com",
  };

  it("renders user info", () => {
    render(<UserCard user={user} onEdit={vi.fn()} />);
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("alice@example.com")).toBeInTheDocument();
  });

  it("calls onEdit with user id when clicked", async () => {
    const onEdit = vi.fn();
    render(<UserCard user={user} onEdit={onEdit} />);

    await userEvent.click(screen.getByRole("button", { name: /edit/i }));
    expect(onEdit).toHaveBeenCalledWith("user-1");
  });
});
```

## Storybook Story Co-location

```typescript
// UserCard.stories.tsx — co-located next to UserCard.tsx
import type { Meta, StoryObj } from "@storybook/react";
import { UserCard } from "./UserCard";

const meta: Meta<typeof UserCard> = {
  component: UserCard,
  tags: ["autodocs"],
};
export default meta;

type Story = StoryObj<typeof UserCard>;

export const Default: Story = {
  args: {
    user: { id: "user-1", name: "Alice", email: "alice@example.com" },
    onEdit: () => {},
  },
};

export const LongName: Story = {
  args: {
    user: { id: "user-2", name: "Alexander Bartholomew Charleston", email: "abc@example.com" },
    onEdit: () => {},
  },
};
```
