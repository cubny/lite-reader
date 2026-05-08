import { render } from 'preact';

import { html } from './util/html.js';
import { AppPage } from './pages/AppPage.js';

const root = document.getElementById('app');
render(html`<${AppPage} />`, root);
