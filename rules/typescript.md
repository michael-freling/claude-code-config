---
paths: **/*.ts, **/*.tsx
---
# TypeScript Coding Guidelines

- Implement UI components based on a mobile-first approach
- Import assets from files instead of inlining SVG, base64, or XML data:
  ```tsx
  // Bad - inline SVG
  const Icon = () => <svg viewBox="0 0 24 24"><path d="M12 2C6.48..."/></svg>

  // Good - import from file
  import Icon from './icon.svg'
  ```
## React Hooks

- Only use `useMemo` when there is a real performance need (expensive computation or preventing unnecessary child re-renders). Add a brief explanation of why it's needed.
- Minimize `useEffect` hooks per component; consolidate related effects that depend on the same data source instead of creating separate effects for each field.
