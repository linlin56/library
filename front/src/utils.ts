/**
 * Pure utility functions (no Angular dependency) used by the services.
 * They can be tested with Mocha directly under Node.js.
 */

export interface SearchFilters {
  title?: string;
  author?: string;
  isbn?: string;
  barcode?: string;
  issueNumber?: string;
  available?: boolean;
}

export interface LoggedInUser {
  id: number;
  last_name: string;
  first_name: string;
  role: string;
  message: string;
}

/**
 * Converts a SearchFilters object into a dictionary of URL parameters.
 * Empty or whitespace-only fields are left out.
 */
export function buildSearchParams(filters: SearchFilters): Record<string, string> {
  const params: Record<string, string> = {};
  if (filters.title?.trim()) params['title'] = filters.title.trim();
  if (filters.author?.trim()) params['author'] = filters.author.trim();
  if (filters.isbn?.trim()) params['isbn'] = filters.isbn.trim();
  if (filters.barcode?.trim()) params['barcode'] = filters.barcode.trim();
  if (filters.issueNumber?.trim()) params['issue_number'] = filters.issueNumber.trim();
  if (filters.available) params['available'] = 'true';
  return params;
}

/**
 * Deserializes raw JSON into a LoggedInUser.
 * Returns null if the JSON is invalid.
 */
export function parseUser(raw: string): LoggedInUser | null {
  try {
    return JSON.parse(raw) as LoggedInUser;
  } catch {
    return null;
  }
}

/**
 * Tells whether the logged-in user has the librarian role.
 */
export function isLibrarian(user: LoggedInUser | null): boolean {
  return user?.role === 'librarian';
}
