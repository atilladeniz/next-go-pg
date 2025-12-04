---
allowed-tools: Bash(goca:*), Bash(cd:*), Bash(make:*), Read, Edit, Write
argument-hint: <FeatureName> --fields "field1:type,field2:type"
description: Generate full-stack feature with Goca CLI
---

# Feature Generator

Generate a complete full-stack feature using Goca CLI.

## Arguments

$ARGUMENTS

## Steps

1. **Generate Backend Feature** with Goca:
   ```bash
   cd backend
   goca feature $ARGUMENTS
   ```

2. **Register Entity** in `backend/internal/domain/registry.go`:
   - Add the new entity to `AllEntities()`: `&<FeatureName>{}`

3. **Generate API Client**:
   ```bash
   cd ..
   make api
   ```

4. **Show next steps** for Frontend integration:
   - Create page in `frontend/src/app/(protected)/<feature>/page.tsx`
   - Use HydrationBoundary pattern with prefetchQuery
   - Create client component with useGet<Feature> hook

## Field Types

- `string` - Text
- `int` - Integer
- `float64` - Decimal
- `bool` - Boolean
- `time.Time` - Timestamp

## Example

```
/feature Product --fields "name:string,price:float64,stock:int"
```
