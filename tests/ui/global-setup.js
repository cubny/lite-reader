/**
 * Global setup for UI tests.
 * Runs the setup wizard to create the initial admin user and enable signups
 * so that existing test flows (which create users via the signup page) continue to work.
 */

import { request } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
const ADMIN_EMAIL = 'admin@test.example';
const ADMIN_PASSWORD = 'AdminPassword123!';

export default async function globalSetup() {
  const ctx = await request.newContext({ baseURL: BASE_URL });

  // Check if setup is needed
  const statusResp = await ctx.get('/setup/status');
  const status = await statusResp.json();

  if (status.needs_setup) {
    // Run setup wizard: create admin and allow signups for test users
    const setupResp = await ctx.post('/setup', {
      data: {
        email: ADMIN_EMAIL,
        password: ADMIN_PASSWORD,
        confirm_password: ADMIN_PASSWORD,
        allow_signup: true,
      },
    });

    if (!setupResp.ok()) {
      throw new Error(`Setup failed: ${setupResp.status()} ${await setupResp.text()}`);
    }
  } else if (!status.allow_signup) {
    // Setup already done but signup disabled — enable it for tests
    const loginResp = await ctx.post('/login', {
      data: { email: ADMIN_EMAIL, password: ADMIN_PASSWORD },
    });
    const loginData = await loginResp.json();

    await ctx.put('/settings', {
      data: { allow_signup: true },
      headers: { Authorization: `Bearer ${loginData.access_token}` },
    });
  }

  await ctx.dispose();
}
