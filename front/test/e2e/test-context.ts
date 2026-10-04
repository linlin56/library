/**
 * TestContext — resets the database between test files.
 *
 * Usage in each .e2e.ts file:
 *
 *   const ctx = new TestContext();
 *
 *   describe('Suite ...', () => {
 *     before(async () => {
 *       await ctx.before();               // empties every table, resets the sequences
 *       await ctx.exec(`INSERT ...`);     // inserts the data specific to this file
 *     });
 *     after(async () => await ctx.after());
 *     // ... nested describes ...
 *   });
 */

import { Client } from 'pg';

export class TestContext {
  private readonly client: Client;

  constructor() {
    this.client = new Client({
      host: process.env['DB_HOST'] ?? '127.0.0.1',
      port: parseInt(process.env['DB_PORT'] ?? '5435'),
      user: process.env['DB_USER'] ?? 'postgres',
      password: process.env['DB_PASSWORD'] ?? 'admin_library',
      database: process.env['DB_NAME'] ?? 'library',
    });
  }

  /**
   * Opens the connection, applies the missing schema migrations,
   * then empties every table with RESTART IDENTITY CASCADE.
   * Uses TRUNCATE (not DROP) so that the backend connections stay valid.
   */
  async before(): Promise<void> {
    await this.client.connect();

    // Columns required by the backend but missing from the initial schema
    await this.client.query(
      `ALTER TABLE items ADD COLUMN IF NOT EXISTS location_id INT REFERENCES locations(id)`,
    );
    await this.client.query(`ALTER TABLE journals ADD COLUMN IF NOT EXISTS publication_date DATE`);

    await this.client.query(`
      TRUNCATE
        teachers, students, ONLY users,
        librarians,
        copies,
        books, journals, ONLY items,
        authors, addresses,
        locations, categories,
        departments
      RESTART IDENTITY CASCADE
    `);
  }

  /** Closes the connection. */
  async after(): Promise<void> {
    await this.client.end();
  }

  /** Runs one or more INSERT statements (separated by ;) in the test context. */
  async exec(sql: string): Promise<void> {
    const statements = sql
      .split(';')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);
    for (const stmt of statements) {
      await this.client.query(stmt);
    }
  }
}
