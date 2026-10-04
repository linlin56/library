/**
 * Mocha root hooks — e2e test environment setup.
 *
 * Before each test session:
 *  1. Applies insert.sql to put the database back into a known state.
 *  2. Starts the Go backend if it is not already running.
 *
 * Environment variables (all optional, default values = dev):
 *   API_URL, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
 */

import { ChildProcess, spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

import { Client } from 'pg';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';

// Paths resolved from test/e2e/
const BACK_DIR = path.resolve(__dirname, '..', '..', '..', 'back');
const SQL_FILE = path.resolve(__dirname, '..', '..', '..', 'db', 'scripts', 'insert.sql');

// Database settings (same values as back/db/db.env)
const DB_HOST = process.env['DB_HOST'] ?? '127.0.0.1';
const DB_PORT = process.env['DB_PORT'] ?? '5435';
const DB_USER = process.env['DB_USER'] ?? 'postgres';
const DB_PASSWORD = process.env['DB_PASSWORD'] ?? 'admin_library';
const DB_NAME = process.env['DB_NAME'] ?? 'library';

let backendProcess: ChildProcess | null = null;

// ─── Helpers ──────────────────────────────────────────────────────────────────

async function seedDatabase(): Promise<void> {
  console.log('\n  → Initializing the test database (insert.sql)...');
  const sql = fs.readFileSync(SQL_FILE, 'utf-8');
  const client = new Client({
    host: DB_HOST,
    port: parseInt(DB_PORT),
    user: DB_USER,
    password: DB_PASSWORD,
    database: DB_NAME,
  });
  await client.connect();
  await client.query(sql);
  await client.end();
  console.log('  ✔ Database initialized with the test data');
}

async function isBackendUp(): Promise<boolean> {
  try {
    await fetch(`${API_URL}/api/items`, { signal: AbortSignal.timeout(1500) });
    return true;
  } catch {
    return false;
  }
}

async function waitForBackend(timeoutMs = 20000): Promise<void> {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    if (await isBackendUp()) return;
    await new Promise((r) => setTimeout(r, 500));
  }
  throw new Error(`Backend (${API_URL}) unreachable after ${timeoutMs / 1000}s`);
}

// ─── Hooks Mocha ──────────────────────────────────────────────────────────────

export const mochaHooks = {
  async beforeAll() {
    // If the backend is already running, do not run insert.sql again:
    // DROP + CREATE TABLE would break the backend's active connections (→ 500 on queries).
    // Each test file calls TestContext.before(), which runs a TRUNCATE (no DROP).
    if (await isBackendUp()) {
      console.log('  ✔ Backend already running — reusing it (no re-seed)\n');
      return;
    }

    // Backend not available: create the schema + initial data, then start it.
    await seedDatabase();

    console.log('  → Starting the backend (go run main.go)...');
    backendProcess = spawn('go', ['run', 'main.go'], {
      cwd: BACK_DIR,
      env: {
        ...process.env,
        DB_HOST,
        DB_PORT,
        DB_USER,
        DB_PASSWORD,
        DB_NAME,
      },
      stdio: 'pipe',
    });

    backendProcess.on('error', (err) => {
      throw new Error(`Unable to start the backend: ${err.message}`);
    });

    await waitForBackend();
    console.log('  ✔ Backend ready\n');
  },

  async afterAll() {
    if (backendProcess) {
      backendProcess.kill('SIGTERM');
      backendProcess = null;
    }
  },
};
