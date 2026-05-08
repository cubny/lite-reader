import { render } from 'preact';

import { html } from './util/html.js';
import { SignupPage } from './pages/SignupPage.js';

const root = document.getElementById('app');
render(html`<${SignupPage} />`, root);
