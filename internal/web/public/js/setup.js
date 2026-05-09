import { render } from 'preact';

import { html } from './util/html.js';
import { SetupPage } from './pages/SetupPage.js';

const root = document.getElementById('app');
render(html`<${SetupPage} />`, root);
