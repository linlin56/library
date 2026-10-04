# Frontend architecture

## Folder structure

```
front/src/app/
├── home/            Home page
├── login/           Login form
├── search/          Item search + loan modal
├── my-loans/        Active loans of the logged-in user
├── create-user/     User creation form
├── create-item/     Book or journal creation form
├── user-loans/      Loan management for librarians
├── overdue-loans/   List of overdue loans
├── guards/          Route guards
├── i18n/            Translations, language service, pipes, interceptor
├── services/        API calls
└── shared/          Shared components (navbar)
```

## Routes

| Route | Component | Access |
|-------|-----------|--------|
| `/` | HomeComponent | Public |
| `/login` | LoginComponent | Logged out (guestGuard) |
| `/search` | SearchComponent | Public |
| `/my-loans` | MyLoansComponent | Logged in (authGuard) |
| `/create-user` | CreateUserComponent | Librarian (librarianGuard) |
| `/create-item` | CreateItemComponent | Librarian (librarianGuard) |
| `/user-loans` | UserLoansComponent | Librarian (librarianGuard) |
| `/overdue-loans` | OverdueLoansComponent | Librarian (librarianGuard) |
| `**` | — | Redirects to `/` |

## Guards

- `authGuard`: requires a logged-in user, redirects to `/login` otherwise
- `librarianGuard`: requires the `librarian` role, redirects to `/` otherwise
- `guestGuard`: blocks `/login` when already logged in, redirects to `/`

## Session management

`AuthService` is the central point of the session:

- Persists the logged-in user in `sessionStorage`
- Exposes a `BehaviorSubject<LoggedInUser | null>` (`loggedInUser$`)
- Methods: `login()`, `logout()`, `isLoggedIn()`, `getLoggedInUser()`

## Internationalization

- `LanguageService` holds the current language in a signal and remembers it in `localStorage`
- `TranslatePipe` (`| t`) translates a key, `AmountPipe` (`| amount`) formats euro amounts for the current language
- `languageInterceptor` adds the `Accept-Language` header to every API call

## API services

| Service | Called endpoints |
|---------|------------------|
| `HomeService` | `GET /home` |
| `AuthService` | `POST /login` |
| `SearchService` | `GET /items` |
| `LoanService` | `GET /loans/check`, `POST /loans`, `DELETE /loans/:id`, `GET /loans`, `GET /loans/overdue`, `GET /items/:id/copies`, `GET /items/:id/copies/all`, `POST /items/:id/copies` |
| `UserService` | `GET /users/search`, `GET /users/:id/deposit`, `PUT /users/:id/deposit` |
| `CreateUserService` | `GET /departments`, `POST /users` |
| `ItemService` | `GET /authors`, `GET /locations`, `POST /books`, `POST /journals` |

## TypeScript models

### Session

```typescript
LoggedInUser  { id, last_name, first_name, role, message }
UserSummary   { id, last_name, first_name, phone_number, role }
```

### Items

```typescript
Item        { id, title, deposit, type, isbn?, author?, issue_number?, available_copies, location? }
Author      { id, last_name, first_name }
Location    { id, bay_number, shelf_number, level, category_name }
Department  { id, name }
```

### Copies

```typescript
AvailableCopy { id, barcode }
CopyDetail    { id, barcode, is_borrowed, loan_due_date? }
```

### Loans

```typescript
Loan         { id, barcode, title, start_date, due_date, overdue }
LoanPreview  { title, barcode, deposit, current_balance, new_balance }
OverdueLoan  { id, barcode, title, due_date, last_name, first_name, phone_number }
```

### Deposit and creation

```typescript
DepositInfo         { deposit_balance, total_deposit }
CreateUserResponse  { login, password, message }
CreateItemResponse  { id, message }
SearchFilters       { title?, author?, isbn?, barcode?, issueNumber?, available? }
```

## Pages — functional roles

**Home**: shows the 4 global statistics (books, journals, available copies, users).

**Login**: login/password form. Redirects automatically when already logged in.

**Search**: filters items by title, author, ISBN, barcode, journal issue number or availability. Logged-in users can borrow from a modal; librarians can also manage copies in a modal.

**My loans**: lists the logged-in user's active loans with due dates, an overdue indicator and a "mark as returned" button.

**Create user**: reactive form with conditional fields depending on the status (year of study for students, department for teachers). Shows the generated login and password after creation.

**Create item**: two-step creation — step 1: book or journal details (with author autocompletion); step 2: adding copies with a barcode.

**User loans**: find a user, view their active loans and update their total deposit.

**Overdue loans**: list of loans past their due date, with the borrower's contact details.
