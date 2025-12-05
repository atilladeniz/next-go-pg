---
name: shadcn-components
description: Add and use shadcn/ui components. Use when building UI, adding new components, or styling pages with shadcn.
allowed-tools: Read, Edit, Write, Bash, Glob
---

# shadcn/ui Components (FSD)

This project uses shadcn/ui with the neutral theme and all components pre-installed.

## FSD Pfad

All components are in `frontend/src/shared/ui/`:

### Layout

- `card` - Card container with header, content, footer
- `separator` - Visual divider
- `aspect-ratio` - Maintain aspect ratios
- `scroll-area` - Custom scrollbars
- `resizable` - Resizable panels

### Forms

- `button` - Buttons with variants
- `input` - Text input
- `textarea` - Multi-line input
- `label` - Form labels
- `checkbox` - Checkboxes
- `radio-group` - Radio buttons
- `select` - Dropdown select
- `switch` - Toggle switch
- `slider` - Range slider
- `form` - Form with validation
- `input-otp` - OTP input

### Feedback

- `alert` - Alert messages
- `alert-dialog` - Confirmation dialogs
- `dialog` - Modal dialogs
- `drawer` - Slide-out drawer
- `sheet` - Side sheet
- `sonner` - Toast notifications
- `progress` - Progress bar
- `skeleton` - Loading skeleton
- `spinner` - Loading spinner

### Navigation

- `tabs` - Tab navigation
- `navigation-menu` - Main navigation
- `menubar` - Menu bar
- `breadcrumb` - Breadcrumb nav
- `pagination` - Page navigation
- `sidebar` - App sidebar

### Data Display

- `table` - Data tables
- `avatar` - User avatars
- `badge` - Status badges
- `calendar` - Date picker calendar
- `chart` - Charts (recharts)
- `carousel` - Image carousel

### Overlay

- `dropdown-menu` - Dropdown menus
- `context-menu` - Right-click menu
- `popover` - Popovers
- `tooltip` - Tooltips
- `hover-card` - Hover cards
- `command` - Command palette

### Other

- `accordion` - Collapsible sections
- `collapsible` - Toggle content
- `toggle` - Toggle button
- `toggle-group` - Button group

## Usage Examples

### Button Variants

```tsx
import { Button } from "@shared/ui/button"

<Button>Default</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="destructive">Destructive</Button>
<Button variant="link">Link</Button>
```

### Card

```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@shared/ui/card"

<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
    <CardDescription>Description</CardDescription>
  </CardHeader>
  <CardContent>
    Content here
  </CardContent>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>
```

### Form with Input

```tsx
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"

<div className="space-y-2">
  <Label htmlFor="email">Email</Label>
  <Input id="email" type="email" placeholder="name@example.com" />
</div>
```

### Dialog

```tsx
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@shared/ui/dialog"

<Dialog>
  <DialogTrigger asChild>
    <Button>Open Dialog</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Dialog Title</DialogTitle>
      <DialogDescription>Dialog description here.</DialogDescription>
    </DialogHeader>
    <p>Dialog content</p>
  </DialogContent>
</Dialog>
```

### Toast (Sonner)

```tsx
import { toast } from "sonner"

toast.success("Success message")
toast.error("Error message")
toast.info("Info message")
```

## Dark Mode

Components automatically support dark mode via ThemeProvider.

Use `ModeToggle` from Header Widget:

```tsx
import { Header } from "@widgets/header"

// ModeToggle is integrated in the Header
<Header user={session.user} />
```

## Adding New Components

If a component is missing (unlikely, all are installed):

```bash
bunx shadcn@latest add [component-name]
```

## Styling

Components use Tailwind CSS classes. Customize with `className`:

```tsx
<Button className="w-full mt-4">Full Width Button</Button>
```

Use `cn()` utility for conditional classes:

```tsx
import { cn } from "@shared/lib"

<div className={cn("p-4", isActive && "bg-primary")}>
```
