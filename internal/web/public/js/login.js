import { render } from 'preact';

import { html } from './util/html.js';
import { LoginPage } from './pages/LoginPage.js';

const root = document.getElementById('app');
render(html`<${LoginPage} />`, root);
