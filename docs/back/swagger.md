# Endpoints

Base URL: `http://localhost:8080/api`

Every endpoint reads the `Accept-Language` header (`en` by default, or `fr`) to choose the language of the `message` and `error` fields.

## Home

### GET /home

Returns global statistics about the library.

- Access: public
- 200 response: `{ book_count, journal_count, available_copy_count, user_count }`

## Authentication

### POST /login

Authenticates a user or a librarian.

- Access: public
- Body: `{ login, password }`
- 200 response: `{ id, last_name, first_name, role, message }` — `role`: `user` | `librarian`
- 400 response: missing fields
- 401 response: invalid credentials

## Items

### GET /items

Searches items with optional filters.

- Access: public
- Query parameters (all optional): `title`, `author`, `isbn`, `barcode`, `issue_number`, `available` (boolean)
- 200 response: array of items `{ id, title, deposit, type, isbn?, author?, issue_number?, publication_date?, available_copies, location? }` — `type`: `book` | `journal`, `location`: `{ bay_number, shelf_number, level }`
- 204 response: no result
- 400 response: invalid parameter
- 500 response: server error

### GET /authors

Returns every author.

- Access: librarian
- 200 response: array `{ id, last_name, first_name }`
- 500 response: server error

### GET /locations

Returns the available locations.

- Access: librarian
- 200 response: array `{ id, bay_number, shelf_number, level, category_name }`
- 500 response: server error

### POST /books

Creates a new book. The author can be an existing one (`author_id`) or a new one (`author_last_name` + `author_first_name`).

- Access: librarian
- Body: `{ title, isbn, deposit, location_id, author_id?, author_last_name?, author_first_name? }`
- 201 response: `{ id, message }`
- 400 response: invalid data
- 500 response: server error

### POST /journals

Creates a new journal issue.

- Access: librarian
- Body: `{ title, issue_number, publication_date, deposit, location_id }`
- 201 response: `{ id, message }`
- 400 response: invalid data
- 500 response: server error

## Copies

### GET /items/:id/copies

Returns the copies of an item that are not on loan.

- Access: public
- Path parameter: `id` (integer, item ID)
- 200 response: array `{ id, barcode }`
- 204 response: no copy available
- 400 response: invalid ID
- 500 response: server error

### GET /items/:id/copies/all

Returns every copy of an item (available and on loan).

- Access: librarian
- Path parameter: `id` (integer)
- 200 response: array `{ id, barcode, is_borrowed, loan_due_date? }`
- 204 response: no copy
- 400 response: invalid ID
- 500 response: server error

### POST /items/:id/copies

Creates a copy of a given item.

- Access: librarian
- Path parameter: `id` (integer)
- Body: `{ barcode, loan_period_days }`
- 201 response: `{ id }`
- 400 response: invalid data

## Users

### GET /departments

Returns the school departments.

- Access: librarian
- 200 response: array `{ id, name }`
- 500 response: server error

### POST /users

Creates a student, a teacher or an individual.

- Access: librarian
- Body: `{ last_name, first_name, phone_number, birth_date, email, status, study_year?, department_id?, street_number, street_name, postal_code, city }`
- `status`: `student` | `teacher` | `individual`
- 201 response: `{ login, password, message }`
- 400 response: invalid data or wrong format
- 409 response: `{ message, existing_login }` — an account already exists with this login
- 500 response: server error

### GET /users/search

Searches users by criteria.

- Access: librarian
- Query parameters (all optional): `last_name`, `first_name`, `postal_code`, `phone_number`
- 200 response: array `{ id, last_name, first_name, phone_number, role }` — `role`: `student` | `teacher` | `user`
- 204 response: no result
- 400 response: invalid parameter
- 500 response: server error

### GET /users/:id/deposit

Returns a user's deposit balance.

- Access: librarian
- Path parameter: `id` (integer)
- 200 response: `{ deposit_balance, total_deposit }`
- 400 response: invalid ID
- 404 response: user not found

### PUT /users/:id/deposit

Updates a user's total deposit and recomputes their balance.

- Access: librarian
- Path parameter: `id` (integer)
- Body: `{ total_deposit }`
- 200 response: `{ deposit_balance, total_deposit }`
- 400 response: invalid data

## Loans

### GET /loans/check

Checks whether a loan is possible and returns a preview (impact on the deposit).

- Access: librarian
- Query parameters: `user_id` (integer), `barcode` (string)
- 200 response: `{ title, barcode, deposit, current_balance, new_balance }`
- 400 response: missing or invalid parameter
- 422 response: loan not possible (insufficient deposit, copy already on loan, etc.)

### POST /loans

Records a loan.

- Access: librarian
- Body: `{ user_id, barcode }`
- 201 response: `{ message }`
- 400 response: invalid data
- 422 response: loan not possible

### DELETE /loans/:id

Returns a copy and refunds the matching deposit.

- Access: librarian
- Path parameter: `id` (integer, copy ID)
- 200 response: `{ message }`
- 400 response: invalid ID
- 422 response: copy not on loan or not found

### GET /loans

Returns a user's active loans.

- Access: logged-in user
- Query parameter: `user_id` (integer)
- 200 response: array `{ id, barcode, title, start_date, due_date, overdue }`
- 204 response: no active loan
- 400 response: invalid parameter
- 500 response: server error

### GET /loans/overdue

Returns every loan whose due date has passed.

- Access: librarian
- 200 response: array `{ id, barcode, title, due_date, last_name, first_name, phone_number }`
- 204 response: no overdue loan
- 500 response: server error
