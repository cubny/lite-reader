/**
 * Database utilities for UI tests
 */

import { exec } from 'child_process';
import { promisify } from 'util';
import * as fs from 'fs';
import * as path from 'path';

const execAsync = promisify(exec);

// Test database path
const TEST_DB_PATH = process.env.TEST_DB_PATH || 'data/test-dev.db';

/**
 * Clean up the test database
 */
export async function cleanDatabase() {
  const dbPath = path.join(process.cwd(), TEST_DB_PATH);
  
  if (fs.existsSync(dbPath)) {
    fs.unlinkSync(dbPath);
  }
  
  // Ensure the data directory exists
  const dataDir = path.dirname(dbPath);
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }
}

/**
 * Generate a unique username for testing
 */
export function generateUsername(prefix = 'testuser') {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 1000);
  return `${prefix}_${timestamp}_${random}`;
}

/**
 * Generate a unique email for testing
 */
export function generateEmail(username) {
  const user = username || generateUsername();
  return `${user}@test.example`;
}

/**
 * Generate a test password
 */
export function generatePassword() {
  return 'TestPassword123!';
}

/**
 * Wait for a condition to be true
 */
export async function waitFor(
  condition,
  timeout = 5000,
  interval = 100
) {
  const startTime = Date.now();
  
  while (Date.now() - startTime < timeout) {
    if (await condition()) {
      return;
    }
    await new Promise(resolve => setTimeout(resolve, interval));
  }
  
  throw new Error(`Condition not met within ${timeout}ms`);
}

/**
 * Mock feed URLs for testing
 */
export const MOCK_FEEDS = {
  techNews: 'http://localhost:3002/feeds/tech-news.xml',
  scienceBlog: 'http://localhost:3002/feeds/science-blog.xml',
  empty: 'http://localhost:3002/feeds/empty.xml',
};
