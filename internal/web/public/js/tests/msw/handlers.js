import { http, HttpResponse } from 'msw';

export const state = {
  users: new Map(), // email -> password
  feeds: [],
  items: new Map(), // feedID -> []item
  nextFeedId: 1,
  nextItemId: 1,
};

export function reset() {
  state.users.clear();
  state.feeds = [];
  state.items.clear();
  state.nextFeedId = 1;
  state.nextItemId = 1;
}

function makeItem(feedId, overrides = {}) {
  const id = state.nextItemId++;
  return {
    id,
    title: `Item ${id}`,
    desc: '<p>body</p>',
    link: 'http://example.com/a',
    is_new: true,
    starred: false,
    timestamp: new Date().toISOString(),
    dir: 'ltr',
    feed_id: feedId,
    ...overrides,
  };
}

function unreadAll() {
  return [...state.items.values()].flat().filter((i) => i.is_new);
}
function starredAll() {
  return [...state.items.values()].flat().filter((i) => i.starred);
}

export const handlers = [
  http.post('/signup', async ({ request }) => {
    const body = await request.json();
    if (!body.email || !body.password) {
      return HttpResponse.json({ message: 'invalid request' }, { status: 400 });
    }
    if (state.users.has(body.email)) {
      return HttpResponse.json({ message: 'email already in use' }, { status: 400 });
    }
    state.users.set(body.email, body.password);
    return new HttpResponse(null, { status: 201 });
  }),

  http.post('/login', async ({ request }) => {
    const body = await request.json();
    if (state.users.get(body.email) !== body.password) {
      return HttpResponse.json({ message: 'invalid email or password' }, { status: 400 });
    }
    return HttpResponse.json({
      access_token: 'tok-' + body.email,
      refresh_token: 'r-' + body.email,
      expires_in: 3600,
    });
  }),

  http.get('/feeds', () => HttpResponse.json(state.feeds)),

  http.post('/feeds', async ({ request }) => {
    const body = await request.json();
    const id = state.nextFeedId++;
    const feed = {
      id,
      title: `Feed ${id}`,
      desc: '',
      link: body.url,
      url: body.url,
      lang: 'en',
      updated_at: new Date().toISOString(),
      unread_count: 0,
    };
    state.feeds.push(feed);
    state.items.set(feed.id, [makeItem(feed.id), makeItem(feed.id)]);
    feed.unread_count = state.items.get(feed.id).filter((i) => i.is_new).length;
    return HttpResponse.json(feed);
  }),

  http.delete('/feeds/:id', ({ params }) => {
    const id = Number(params.id);
    state.feeds = state.feeds.filter((f) => f.id !== id);
    state.items.delete(id);
    return new HttpResponse(null, { status: 204 });
  }),

  http.put('/feeds/:id/fetch', ({ params }) => {
    const id = Number(params.id);
    const list = state.items.get(id) || [];
    return HttpResponse.json(list);
  }),

  http.post('/feeds/:id/read', ({ params }) => {
    const id = Number(params.id);
    (state.items.get(id) || []).forEach((i) => { i.is_new = false; });
    const f = state.feeds.find((f) => f.id === id);
    if (f) f.unread_count = 0;
    return HttpResponse.json({});
  }),

  http.post('/feeds/:id/unread', ({ params }) => {
    const id = Number(params.id);
    (state.items.get(id) || []).forEach((i) => { i.is_new = true; });
    const f = state.feeds.find((f) => f.id === id);
    if (f) f.unread_count = state.items.get(id).length;
    return HttpResponse.json({});
  }),

  http.get('/feeds/:id/items', ({ params }) => {
    const id = Number(params.id);
    return HttpResponse.json(state.items.get(id) || []);
  }),

  http.put('/items/:id', async ({ params, request }) => {
    const id = Number(params.id);
    const body = await request.json();
    for (const list of state.items.values()) {
      const it = list.find((i) => i.id === id);
      if (it) {
        it.is_new = !!body.is_new;
        it.starred = !!body.starred;
        break;
      }
    }
    return HttpResponse.json({});
  }),

  http.get('/items/unread', () => HttpResponse.json(unreadAll())),
  http.get('/items/starred', () => HttpResponse.json(starredAll())),
  http.get('/items/unread/count', () => HttpResponse.json({ count: unreadAll().length })),
  http.get('/items/starred/count', () => HttpResponse.json({ count: starredAll().length })),
];
