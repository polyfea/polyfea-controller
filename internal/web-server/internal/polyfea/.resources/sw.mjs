try {
  self["workbox:core:7.0.0"] && _();
} catch {
}
const Le = (s, ...e) => {
  let t = s;
  return e.length > 0 && (t += ` :: ${JSON.stringify(e)}`), t;
}, qe = Le;
class f extends Error {
  /**
   *
   * @param {string} errorCode The error code that
   * identifies this particular error.
   * @param {Object=} details Any relevant arguments
   * that will help developers identify issues should
   * be added as a key on the context object.
   */
  constructor(e, t) {
    const n = qe(e, t);
    super(n), this.name = e, this.details = t;
  }
}
const g = {
  googleAnalytics: "googleAnalytics",
  precache: "precache-v2",
  prefix: "workbox",
  runtime: "runtime",
  suffix: typeof registration < "u" ? registration.scope : ""
}, K = (s) => [g.prefix, s, g.suffix].filter((e) => e && e.length > 0).join("-"), Oe = (s) => {
  for (const e of Object.keys(g))
    s(e);
}, O = {
  updateDetails: (s) => {
    Oe((e) => {
      typeof s[e] == "string" && (g[e] = s[e]);
    });
  },
  getGoogleAnalyticsName: (s) => s || K(g.googleAnalytics),
  getPrecacheName: (s) => s || K(g.precache),
  getPrefix: () => g.prefix,
  getRuntimeName: (s) => s || K(g.runtime),
  getSuffix: () => g.suffix
};
function oe(s, e) {
  const t = e();
  return s.waitUntil(t), t;
}
try {
  self["workbox:precaching:7.0.0"] && _();
} catch {
}
const Pe = "__WB_REVISION__";
function Ie(s) {
  if (!s)
    throw new f("add-to-cache-list-unexpected-type", { entry: s });
  if (typeof s == "string") {
    const i = new URL(s, location.href);
    return {
      cacheKey: i.href,
      url: i.href
    };
  }
  const { revision: e, url: t } = s;
  if (!t)
    throw new f("add-to-cache-list-unexpected-type", { entry: s });
  if (!e) {
    const i = new URL(t, location.href);
    return {
      cacheKey: i.href,
      url: i.href
    };
  }
  const n = new URL(t, location.href), r = new URL(t, location.href);
  return n.searchParams.set(Pe, e), {
    cacheKey: n.href,
    url: r.href
  };
}
class Ne {
  constructor() {
    this.updatedURLs = [], this.notUpdatedURLs = [], this.handlerWillStart = async ({ request: e, state: t }) => {
      t && (t.originalRequest = e);
    }, this.cachedResponseWillBeUsed = async ({ event: e, state: t, cachedResponse: n }) => {
      if (e.type === "install" && t && t.originalRequest && t.originalRequest instanceof Request) {
        const r = t.originalRequest.url;
        n ? this.notUpdatedURLs.push(r) : this.updatedURLs.push(r);
      }
      return n;
    };
  }
}
class Ue {
  constructor({ precacheController: e }) {
    this.cacheKeyWillBeUsed = async ({ request: t, params: n }) => {
      const r = (n == null ? void 0 : n.cacheKey) || this._precacheController.getCacheKeyForURL(t.url);
      return r ? new Request(r, { headers: t.headers }) : t;
    }, this._precacheController = e;
  }
}
let x;
function Ae() {
  if (x === void 0) {
    const s = new Response("");
    if ("body" in s)
      try {
        new Response(s.body), x = !0;
      } catch {
        x = !1;
      }
    x = !1;
  }
  return x;
}
async function je(s, e) {
  let t = null;
  if (s.url && (t = new URL(s.url).origin), t !== self.location.origin)
    throw new f("cross-origin-copy-response", { origin: t });
  const n = s.clone(), i = {
    headers: new Headers(n.headers),
    status: n.status,
    statusText: n.statusText
  }, a = Ae() ? n.body : await n.blob();
  return new Response(a, i);
}
const Me = (s) => new URL(String(s), location.href).href.replace(new RegExp(`^${location.origin}`), "");
function ce(s, e) {
  const t = new URL(s);
  for (const n of e)
    t.searchParams.delete(n);
  return t.href;
}
async function Fe(s, e, t, n) {
  const r = ce(e.url, t);
  if (e.url === r)
    return s.match(e, n);
  const i = Object.assign(Object.assign({}, n), { ignoreSearch: !0 }), a = await s.keys(e, i);
  for (const o of a) {
    const l = ce(o.url, t);
    if (r === l)
      return s.match(o, n);
  }
}
class Be {
  /**
   * Creates a promise and exposes its resolve and reject functions as methods.
   */
  constructor() {
    this.promise = new Promise((e, t) => {
      this.resolve = e, this.reject = t;
    });
  }
}
const _e = /* @__PURE__ */ new Set();
async function Ke() {
  for (const s of _e)
    await s();
}
function Re(s) {
  return new Promise((e) => setTimeout(e, s));
}
try {
  self["workbox:strategies:7.0.0"] && _();
} catch {
}
function N(s) {
  return typeof s == "string" ? new Request(s) : s;
}
class We {
  /**
   * Creates a new instance associated with the passed strategy and event
   * that's handling the request.
   *
   * The constructor also initializes the state that will be passed to each of
   * the plugins handling this request.
   *
   * @param {workbox-strategies.Strategy} strategy
   * @param {Object} options
   * @param {Request|string} options.request A request to run this strategy for.
   * @param {ExtendableEvent} options.event The event associated with the
   *     request.
   * @param {URL} [options.url]
   * @param {*} [options.params] The return value from the
   *     {@link workbox-routing~matchCallback} (if applicable).
   */
  constructor(e, t) {
    this._cacheKeys = {}, Object.assign(this, t), this.event = t.event, this._strategy = e, this._handlerDeferred = new Be(), this._extendLifetimePromises = [], this._plugins = [...e.plugins], this._pluginStateMap = /* @__PURE__ */ new Map();
    for (const n of this._plugins)
      this._pluginStateMap.set(n, {});
    this.event.waitUntil(this._handlerDeferred.promise);
  }
  /**
   * Fetches a given request (and invokes any applicable plugin callback
   * methods) using the `fetchOptions` (for non-navigation requests) and
   * `plugins` defined on the `Strategy` object.
   *
   * The following plugin lifecycle methods are invoked when using this method:
   * - `requestWillFetch()`
   * - `fetchDidSucceed()`
   * - `fetchDidFail()`
   *
   * @param {Request|string} input The URL or request to fetch.
   * @return {Promise<Response>}
   */
  async fetch(e) {
    const { event: t } = this;
    let n = N(e);
    if (n.mode === "navigate" && t instanceof FetchEvent && t.preloadResponse) {
      const a = await t.preloadResponse;
      if (a)
        return a;
    }
    const r = this.hasCallback("fetchDidFail") ? n.clone() : null;
    try {
      for (const a of this.iterateCallbacks("requestWillFetch"))
        n = await a({ request: n.clone(), event: t });
    } catch (a) {
      if (a instanceof Error)
        throw new f("plugin-error-request-will-fetch", {
          thrownErrorMessage: a.message
        });
    }
    const i = n.clone();
    try {
      let a;
      a = await fetch(n, n.mode === "navigate" ? void 0 : this._strategy.fetchOptions);
      for (const o of this.iterateCallbacks("fetchDidSucceed"))
        a = await o({
          event: t,
          request: i,
          response: a
        });
      return a;
    } catch (a) {
      throw r && await this.runCallbacks("fetchDidFail", {
        error: a,
        event: t,
        originalRequest: r.clone(),
        request: i.clone()
      }), a;
    }
  }
  /**
   * Calls `this.fetch()` and (in the background) runs `this.cachePut()` on
   * the response generated by `this.fetch()`.
   *
   * The call to `this.cachePut()` automatically invokes `this.waitUntil()`,
   * so you do not have to manually call `waitUntil()` on the event.
   *
   * @param {Request|string} input The request or URL to fetch and cache.
   * @return {Promise<Response>}
   */
  async fetchAndCachePut(e) {
    const t = await this.fetch(e), n = t.clone();
    return this.waitUntil(this.cachePut(e, n)), t;
  }
  /**
   * Matches a request from the cache (and invokes any applicable plugin
   * callback methods) using the `cacheName`, `matchOptions`, and `plugins`
   * defined on the strategy object.
   *
   * The following plugin lifecycle methods are invoked when using this method:
   * - cacheKeyWillByUsed()
   * - cachedResponseWillByUsed()
   *
   * @param {Request|string} key The Request or URL to use as the cache key.
   * @return {Promise<Response|undefined>} A matching response, if found.
   */
  async cacheMatch(e) {
    const t = N(e);
    let n;
    const { cacheName: r, matchOptions: i } = this._strategy, a = await this.getCacheKey(t, "read"), o = Object.assign(Object.assign({}, i), { cacheName: r });
    n = await caches.match(a, o);
    for (const l of this.iterateCallbacks("cachedResponseWillBeUsed"))
      n = await l({
        cacheName: r,
        matchOptions: i,
        cachedResponse: n,
        request: a,
        event: this.event
      }) || void 0;
    return n;
  }
  /**
   * Puts a request/response pair in the cache (and invokes any applicable
   * plugin callback methods) using the `cacheName` and `plugins` defined on
   * the strategy object.
   *
   * The following plugin lifecycle methods are invoked when using this method:
   * - cacheKeyWillByUsed()
   * - cacheWillUpdate()
   * - cacheDidUpdate()
   *
   * @param {Request|string} key The request or URL to use as the cache key.
   * @param {Response} response The response to cache.
   * @return {Promise<boolean>} `false` if a cacheWillUpdate caused the response
   * not be cached, and `true` otherwise.
   */
  async cachePut(e, t) {
    const n = N(e);
    await Re(0);
    const r = await this.getCacheKey(n, "write");
    if (!t)
      throw new f("cache-put-with-no-response", {
        url: Me(r.url)
      });
    const i = await this._ensureResponseSafeToCache(t);
    if (!i)
      return !1;
    const { cacheName: a, matchOptions: o } = this._strategy, l = await self.caches.open(a), c = this.hasCallback("cacheDidUpdate"), u = c ? await Fe(
      // TODO(philipwalton): the `__WB_REVISION__` param is a precaching
      // feature. Consider into ways to only add this behavior if using
      // precaching.
      l,
      r.clone(),
      ["__WB_REVISION__"],
      o
    ) : null;
    try {
      await l.put(r, c ? i.clone() : i);
    } catch (h) {
      if (h instanceof Error)
        throw h.name === "QuotaExceededError" && await Ke(), h;
    }
    for (const h of this.iterateCallbacks("cacheDidUpdate"))
      await h({
        cacheName: a,
        oldResponse: u,
        newResponse: i.clone(),
        request: r,
        event: this.event
      });
    return !0;
  }
  /**
   * Checks the list of plugins for the `cacheKeyWillBeUsed` callback, and
   * executes any of those callbacks found in sequence. The final `Request`
   * object returned by the last plugin is treated as the cache key for cache
   * reads and/or writes. If no `cacheKeyWillBeUsed` plugin callbacks have
   * been registered, the passed request is returned unmodified
   *
   * @param {Request} request
   * @param {string} mode
   * @return {Promise<Request>}
   */
  async getCacheKey(e, t) {
    const n = `${e.url} | ${t}`;
    if (!this._cacheKeys[n]) {
      let r = e;
      for (const i of this.iterateCallbacks("cacheKeyWillBeUsed"))
        r = N(await i({
          mode: t,
          request: r,
          event: this.event,
          // params has a type any can't change right now.
          params: this.params
          // eslint-disable-line
        }));
      this._cacheKeys[n] = r;
    }
    return this._cacheKeys[n];
  }
  /**
   * Returns true if the strategy has at least one plugin with the given
   * callback.
   *
   * @param {string} name The name of the callback to check for.
   * @return {boolean}
   */
  hasCallback(e) {
    for (const t of this._strategy.plugins)
      if (e in t)
        return !0;
    return !1;
  }
  /**
   * Runs all plugin callbacks matching the given name, in order, passing the
   * given param object (merged ith the current plugin state) as the only
   * argument.
   *
   * Note: since this method runs all plugins, it's not suitable for cases
   * where the return value of a callback needs to be applied prior to calling
   * the next callback. See
   * {@link workbox-strategies.StrategyHandler#iterateCallbacks}
   * below for how to handle that case.
   *
   * @param {string} name The name of the callback to run within each plugin.
   * @param {Object} param The object to pass as the first (and only) param
   *     when executing each callback. This object will be merged with the
   *     current plugin state prior to callback execution.
   */
  async runCallbacks(e, t) {
    for (const n of this.iterateCallbacks(e))
      await n(t);
  }
  /**
   * Accepts a callback and returns an iterable of matching plugin callbacks,
   * where each callback is wrapped with the current handler state (i.e. when
   * you call each callback, whatever object parameter you pass it will
   * be merged with the plugin's current state).
   *
   * @param {string} name The name fo the callback to run
   * @return {Array<Function>}
   */
  *iterateCallbacks(e) {
    for (const t of this._strategy.plugins)
      if (typeof t[e] == "function") {
        const n = this._pluginStateMap.get(t);
        yield (i) => {
          const a = Object.assign(Object.assign({}, i), { state: n });
          return t[e](a);
        };
      }
  }
  /**
   * Adds a promise to the
   * [extend lifetime promises]{@link https://w3c.github.io/ServiceWorker/#extendableevent-extend-lifetime-promises}
   * of the event event associated with the request being handled (usually a
   * `FetchEvent`).
   *
   * Note: you can await
   * {@link workbox-strategies.StrategyHandler~doneWaiting}
   * to know when all added promises have settled.
   *
   * @param {Promise} promise A promise to add to the extend lifetime promises
   *     of the event that triggered the request.
   */
  waitUntil(e) {
    return this._extendLifetimePromises.push(e), e;
  }
  /**
   * Returns a promise that resolves once all promises passed to
   * {@link workbox-strategies.StrategyHandler~waitUntil}
   * have settled.
   *
   * Note: any work done after `doneWaiting()` settles should be manually
   * passed to an event's `waitUntil()` method (not this handler's
   * `waitUntil()` method), otherwise the service worker thread my be killed
   * prior to your work completing.
   */
  async doneWaiting() {
    let e;
    for (; e = this._extendLifetimePromises.shift(); )
      await e;
  }
  /**
   * Stops running the strategy and immediately resolves any pending
   * `waitUntil()` promises.
   */
  destroy() {
    this._handlerDeferred.resolve(null);
  }
  /**
   * This method will call cacheWillUpdate on the available plugins (or use
   * status === 200) to determine if the Response is safe and valid to cache.
   *
   * @param {Request} options.request
   * @param {Response} options.response
   * @return {Promise<Response|undefined>}
   *
   * @private
   */
  async _ensureResponseSafeToCache(e) {
    let t = e, n = !1;
    for (const r of this.iterateCallbacks("cacheWillUpdate"))
      if (t = await r({
        request: this.request,
        response: t,
        event: this.event
      }) || void 0, n = !0, !t)
        break;
    return n || t && t.status !== 200 && (t = void 0), t;
  }
}
class k {
  /**
   * Creates a new instance of the strategy and sets all documented option
   * properties as public instance properties.
   *
   * Note: if a custom strategy class extends the base Strategy class and does
   * not need more than these properties, it does not need to define its own
   * constructor.
   *
   * @param {Object} [options]
   * @param {string} [options.cacheName] Cache name to store and retrieve
   * requests. Defaults to the cache names provided by
   * {@link workbox-core.cacheNames}.
   * @param {Array<Object>} [options.plugins] [Plugins]{@link https://developers.google.com/web/tools/workbox/guides/using-plugins}
   * to use in conjunction with this caching strategy.
   * @param {Object} [options.fetchOptions] Values passed along to the
   * [`init`](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters)
   * of [non-navigation](https://github.com/GoogleChrome/workbox/issues/1796)
   * `fetch()` requests made by this strategy.
   * @param {Object} [options.matchOptions] The
   * [`CacheQueryOptions`]{@link https://w3c.github.io/ServiceWorker/#dictdef-cachequeryoptions}
   * for any `cache.match()` or `cache.put()` calls made by this strategy.
   */
  constructor(e = {}) {
    this.cacheName = O.getRuntimeName(e.cacheName), this.plugins = e.plugins || [], this.fetchOptions = e.fetchOptions, this.matchOptions = e.matchOptions;
  }
  /**
   * Perform a request strategy and returns a `Promise` that will resolve with
   * a `Response`, invoking all relevant plugin callbacks.
   *
   * When a strategy instance is registered with a Workbox
   * {@link workbox-routing.Route}, this method is automatically
   * called when the route matches.
   *
   * Alternatively, this method can be used in a standalone `FetchEvent`
   * listener by passing it to `event.respondWith()`.
   *
   * @param {FetchEvent|Object} options A `FetchEvent` or an object with the
   *     properties listed below.
   * @param {Request|string} options.request A request to run this strategy for.
   * @param {ExtendableEvent} options.event The event associated with the
   *     request.
   * @param {URL} [options.url]
   * @param {*} [options.params]
   */
  handle(e) {
    const [t] = this.handleAll(e);
    return t;
  }
  /**
   * Similar to {@link workbox-strategies.Strategy~handle}, but
   * instead of just returning a `Promise` that resolves to a `Response` it
   * it will return an tuple of `[response, done]` promises, where the former
   * (`response`) is equivalent to what `handle()` returns, and the latter is a
   * Promise that will resolve once any promises that were added to
   * `event.waitUntil()` as part of performing the strategy have completed.
   *
   * You can await the `done` promise to ensure any extra work performed by
   * the strategy (usually caching responses) completes successfully.
   *
   * @param {FetchEvent|Object} options A `FetchEvent` or an object with the
   *     properties listed below.
   * @param {Request|string} options.request A request to run this strategy for.
   * @param {ExtendableEvent} options.event The event associated with the
   *     request.
   * @param {URL} [options.url]
   * @param {*} [options.params]
   * @return {Array<Promise>} A tuple of [response, done]
   *     promises that can be used to determine when the response resolves as
   *     well as when the handler has completed all its work.
   */
  handleAll(e) {
    e instanceof FetchEvent && (e = {
      event: e,
      request: e.request
    });
    const t = e.event, n = typeof e.request == "string" ? new Request(e.request) : e.request, r = "params" in e ? e.params : void 0, i = new We(this, { event: t, request: n, params: r }), a = this._getResponse(i, n, t), o = this._awaitComplete(a, i, n, t);
    return [a, o];
  }
  async _getResponse(e, t, n) {
    await e.runCallbacks("handlerWillStart", { event: n, request: t });
    let r;
    try {
      if (r = await this._handle(t, e), !r || r.type === "error")
        throw new f("no-response", { url: t.url });
    } catch (i) {
      if (i instanceof Error) {
        for (const a of e.iterateCallbacks("handlerDidError"))
          if (r = await a({ error: i, event: n, request: t }), r)
            break;
      }
      if (!r)
        throw i;
    }
    for (const i of e.iterateCallbacks("handlerWillRespond"))
      r = await i({ event: n, request: t, response: r });
    return r;
  }
  async _awaitComplete(e, t, n, r) {
    let i, a;
    try {
      i = await e;
    } catch {
    }
    try {
      await t.runCallbacks("handlerDidRespond", {
        event: r,
        request: n,
        response: i
      }), await t.doneWaiting();
    } catch (o) {
      o instanceof Error && (a = o);
    }
    if (await t.runCallbacks("handlerDidComplete", {
      event: r,
      request: n,
      response: i,
      error: a
    }), t.destroy(), a)
      throw a;
  }
}
class E extends k {
  /**
   *
   * @param {Object} [options]
   * @param {string} [options.cacheName] Cache name to store and retrieve
   * requests. Defaults to the cache names provided by
   * {@link workbox-core.cacheNames}.
   * @param {Array<Object>} [options.plugins] {@link https://developers.google.com/web/tools/workbox/guides/using-plugins|Plugins}
   * to use in conjunction with this caching strategy.
   * @param {Object} [options.fetchOptions] Values passed along to the
   * {@link https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters|init}
   * of all fetch() requests made by this strategy.
   * @param {Object} [options.matchOptions] The
   * {@link https://w3c.github.io/ServiceWorker/#dictdef-cachequeryoptions|CacheQueryOptions}
   * for any `cache.match()` or `cache.put()` calls made by this strategy.
   * @param {boolean} [options.fallbackToNetwork=true] Whether to attempt to
   * get the response from the network if there's a precache miss.
   */
  constructor(e = {}) {
    e.cacheName = O.getPrecacheName(e.cacheName), super(e), this._fallbackToNetwork = e.fallbackToNetwork !== !1, this.plugins.push(E.copyRedirectedCacheableResponsesPlugin);
  }
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    const n = await t.cacheMatch(e);
    return n || (t.event && t.event.type === "install" ? await this._handleInstall(e, t) : await this._handleFetch(e, t));
  }
  async _handleFetch(e, t) {
    let n;
    const r = t.params || {};
    if (this._fallbackToNetwork) {
      const i = r.integrity, a = e.integrity, o = !a || a === i;
      n = await t.fetch(new Request(e, {
        integrity: e.mode !== "no-cors" ? a || i : void 0
      })), i && o && e.mode !== "no-cors" && (this._useDefaultCacheabilityPluginIfNeeded(), await t.cachePut(e, n.clone()));
    } else
      throw new f("missing-precache-entry", {
        cacheName: this.cacheName,
        url: e.url
      });
    return n;
  }
  async _handleInstall(e, t) {
    this._useDefaultCacheabilityPluginIfNeeded();
    const n = await t.fetch(e);
    if (!await t.cachePut(e, n.clone()))
      throw new f("bad-precaching-response", {
        url: e.url,
        status: n.status
      });
    return n;
  }
  /**
   * This method is complex, as there a number of things to account for:
   *
   * The `plugins` array can be set at construction, and/or it might be added to
   * to at any time before the strategy is used.
   *
   * At the time the strategy is used (i.e. during an `install` event), there
   * needs to be at least one plugin that implements `cacheWillUpdate` in the
   * array, other than `copyRedirectedCacheableResponsesPlugin`.
   *
   * - If this method is called and there are no suitable `cacheWillUpdate`
   * plugins, we need to add `defaultPrecacheCacheabilityPlugin`.
   *
   * - If this method is called and there is exactly one `cacheWillUpdate`, then
   * we don't have to do anything (this might be a previously added
   * `defaultPrecacheCacheabilityPlugin`, or it might be a custom plugin).
   *
   * - If this method is called and there is more than one `cacheWillUpdate`,
   * then we need to check if one is `defaultPrecacheCacheabilityPlugin`. If so,
   * we need to remove it. (This situation is unlikely, but it could happen if
   * the strategy is used multiple times, the first without a `cacheWillUpdate`,
   * and then later on after manually adding a custom `cacheWillUpdate`.)
   *
   * See https://github.com/GoogleChrome/workbox/issues/2737 for more context.
   *
   * @private
   */
  _useDefaultCacheabilityPluginIfNeeded() {
    let e = null, t = 0;
    for (const [n, r] of this.plugins.entries())
      r !== E.copyRedirectedCacheableResponsesPlugin && (r === E.defaultPrecacheCacheabilityPlugin && (e = n), r.cacheWillUpdate && t++);
    t === 0 ? this.plugins.push(E.defaultPrecacheCacheabilityPlugin) : t > 1 && e !== null && this.plugins.splice(e, 1);
  }
}
E.defaultPrecacheCacheabilityPlugin = {
  async cacheWillUpdate({ response: s }) {
    return !s || s.status >= 400 ? null : s;
  }
};
E.copyRedirectedCacheableResponsesPlugin = {
  async cacheWillUpdate({ response: s }) {
    return s.redirected ? await je(s) : s;
  }
};
class ze {
  /**
   * Create a new PrecacheController.
   *
   * @param {Object} [options]
   * @param {string} [options.cacheName] The cache to use for precaching.
   * @param {string} [options.plugins] Plugins to use when precaching as well
   * as responding to fetch events for precached assets.
   * @param {boolean} [options.fallbackToNetwork=true] Whether to attempt to
   * get the response from the network if there's a precache miss.
   */
  constructor({ cacheName: e, plugins: t = [], fallbackToNetwork: n = !0 } = {}) {
    this._urlsToCacheKeys = /* @__PURE__ */ new Map(), this._urlsToCacheModes = /* @__PURE__ */ new Map(), this._cacheKeysToIntegrities = /* @__PURE__ */ new Map(), this._strategy = new E({
      cacheName: O.getPrecacheName(e),
      plugins: [
        ...t,
        new Ue({ precacheController: this })
      ],
      fallbackToNetwork: n
    }), this.install = this.install.bind(this), this.activate = this.activate.bind(this);
  }
  /**
   * @type {workbox-precaching.PrecacheStrategy} The strategy created by this controller and
   * used to cache assets and respond to fetch events.
   */
  get strategy() {
    return this._strategy;
  }
  /**
   * Adds items to the precache list, removing any duplicates and
   * stores the files in the
   * {@link workbox-core.cacheNames|"precache cache"} when the service
   * worker installs.
   *
   * This method can be called multiple times.
   *
   * @param {Array<Object|string>} [entries=[]] Array of entries to precache.
   */
  precache(e) {
    this.addToCacheList(e), this._installAndActiveListenersAdded || (self.addEventListener("install", this.install), self.addEventListener("activate", this.activate), this._installAndActiveListenersAdded = !0);
  }
  /**
   * This method will add items to the precache list, removing duplicates
   * and ensuring the information is valid.
   *
   * @param {Array<workbox-precaching.PrecacheController.PrecacheEntry|string>} entries
   *     Array of entries to precache.
   */
  addToCacheList(e) {
    const t = [];
    for (const n of e) {
      typeof n == "string" ? t.push(n) : n && n.revision === void 0 && t.push(n.url);
      const { cacheKey: r, url: i } = Ie(n), a = typeof n != "string" && n.revision ? "reload" : "default";
      if (this._urlsToCacheKeys.has(i) && this._urlsToCacheKeys.get(i) !== r)
        throw new f("add-to-cache-list-conflicting-entries", {
          firstEntry: this._urlsToCacheKeys.get(i),
          secondEntry: r
        });
      if (typeof n != "string" && n.integrity) {
        if (this._cacheKeysToIntegrities.has(r) && this._cacheKeysToIntegrities.get(r) !== n.integrity)
          throw new f("add-to-cache-list-conflicting-integrities", {
            url: i
          });
        this._cacheKeysToIntegrities.set(r, n.integrity);
      }
      if (this._urlsToCacheKeys.set(i, r), this._urlsToCacheModes.set(i, a), t.length > 0) {
        const o = `Workbox is precaching URLs without revision info: ${t.join(", ")}
This is generally NOT safe. Learn more at https://bit.ly/wb-precache`;
        console.warn(o);
      }
    }
  }
  /**
   * Precaches new and updated assets. Call this method from the service worker
   * install event.
   *
   * Note: this method calls `event.waitUntil()` for you, so you do not need
   * to call it yourself in your event handlers.
   *
   * @param {ExtendableEvent} event
   * @return {Promise<workbox-precaching.InstallResult>}
   */
  install(e) {
    return oe(e, async () => {
      const t = new Ne();
      this.strategy.plugins.push(t);
      for (const [i, a] of this._urlsToCacheKeys) {
        const o = this._cacheKeysToIntegrities.get(a), l = this._urlsToCacheModes.get(i), c = new Request(i, {
          integrity: o,
          cache: l,
          credentials: "same-origin"
        });
        await Promise.all(this.strategy.handleAll({
          params: { cacheKey: a },
          request: c,
          event: e
        }));
      }
      const { updatedURLs: n, notUpdatedURLs: r } = t;
      return { updatedURLs: n, notUpdatedURLs: r };
    });
  }
  /**
   * Deletes assets that are no longer present in the current precache manifest.
   * Call this method from the service worker activate event.
   *
   * Note: this method calls `event.waitUntil()` for you, so you do not need
   * to call it yourself in your event handlers.
   *
   * @param {ExtendableEvent} event
   * @return {Promise<workbox-precaching.CleanupResult>}
   */
  activate(e) {
    return oe(e, async () => {
      const t = await self.caches.open(this.strategy.cacheName), n = await t.keys(), r = new Set(this._urlsToCacheKeys.values()), i = [];
      for (const a of n)
        r.has(a.url) || (await t.delete(a), i.push(a.url));
      return { deletedURLs: i };
    });
  }
  /**
   * Returns a mapping of a precached URL to the corresponding cache key, taking
   * into account the revision information for the URL.
   *
   * @return {Map<string, string>} A URL to cache key mapping.
   */
  getURLsToCacheKeys() {
    return this._urlsToCacheKeys;
  }
  /**
   * Returns a list of all the URLs that have been precached by the current
   * service worker.
   *
   * @return {Array<string>} The precached URLs.
   */
  getCachedURLs() {
    return [...this._urlsToCacheKeys.keys()];
  }
  /**
   * Returns the cache key used for storing a given URL. If that URL is
   * unversioned, like `/index.html', then the cache key will be the original
   * URL with a search parameter appended to it.
   *
   * @param {string} url A URL whose cache key you want to look up.
   * @return {string} The versioned URL that corresponds to a cache key
   * for the original URL, or undefined if that URL isn't precached.
   */
  getCacheKeyForURL(e) {
    const t = new URL(e, location.href);
    return this._urlsToCacheKeys.get(t.href);
  }
  /**
   * @param {string} url A cache key whose SRI you want to look up.
   * @return {string} The subresource integrity associated with the cache key,
   * or undefined if it's not set.
   */
  getIntegrityForCacheKey(e) {
    return this._cacheKeysToIntegrities.get(e);
  }
  /**
   * This acts as a drop-in replacement for
   * [`cache.match()`](https://developer.mozilla.org/en-US/docs/Web/API/Cache/match)
   * with the following differences:
   *
   * - It knows what the name of the precache is, and only checks in that cache.
   * - It allows you to pass in an "original" URL without versioning parameters,
   * and it will automatically look up the correct cache key for the currently
   * active revision of that URL.
   *
   * E.g., `matchPrecache('index.html')` will find the correct precached
   * response for the currently active service worker, even if the actual cache
   * key is `'/index.html?__WB_REVISION__=1234abcd'`.
   *
   * @param {string|Request} request The key (without revisioning parameters)
   * to look up in the precache.
   * @return {Promise<Response|undefined>}
   */
  async matchPrecache(e) {
    const t = e instanceof Request ? e.url : e, n = this.getCacheKeyForURL(t);
    if (n)
      return (await self.caches.open(this.strategy.cacheName)).match(n);
  }
  /**
   * Returns a function that looks up `url` in the precache (taking into
   * account revision information), and returns the corresponding `Response`.
   *
   * @param {string} url The precached URL which will be used to lookup the
   * `Response`.
   * @return {workbox-routing~handlerCallback}
   */
  createHandlerBoundToURL(e) {
    const t = this.getCacheKeyForURL(e);
    if (!t)
      throw new f("non-precached-url", { url: e });
    return (n) => (n.request = new Request(e), n.params = Object.assign({ cacheKey: t }, n.params), this.strategy.handle(n));
  }
}
try {
  self["workbox:routing:7.0.0"] && _();
} catch {
}
const Ee = "GET", j = (s) => s && typeof s == "object" ? s : { handle: s };
class He {
  /**
   * Constructor for Route class.
   *
   * @param {workbox-routing~matchCallback} match
   * A callback function that determines whether the route matches a given
   * `fetch` event by returning a non-falsy value.
   * @param {workbox-routing~handlerCallback} handler A callback
   * function that returns a Promise resolving to a Response.
   * @param {string} [method='GET'] The HTTP method to match the Route
   * against.
   */
  constructor(e, t, n = Ee) {
    this.handler = j(t), this.match = e, this.method = n;
  }
  /**
   *
   * @param {workbox-routing-handlerCallback} handler A callback
   * function that returns a Promise resolving to a Response
   */
  setCatchHandler(e) {
    this.catchHandler = j(e);
  }
}
class Ve extends He {
  /**
   * If the regular expression contains
   * [capture groups]{@link https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp#grouping-back-references},
   * the captured values will be passed to the
   * {@link workbox-routing~handlerCallback} `params`
   * argument.
   *
   * @param {RegExp} regExp The regular expression to match against URLs.
   * @param {workbox-routing~handlerCallback} handler A callback
   * function that returns a Promise resulting in a Response.
   * @param {string} [method='GET'] The HTTP method to match the Route
   * against.
   */
  constructor(e, t, n) {
    const r = ({ url: i }) => {
      const a = e.exec(i.href);
      if (a && !(i.origin !== location.origin && a.index !== 0))
        return a.slice(1);
    };
    super(r, t, n);
  }
}
class $e {
  /**
   * Initializes a new Router.
   */
  constructor() {
    this._routes = /* @__PURE__ */ new Map(), this._defaultHandlerMap = /* @__PURE__ */ new Map();
  }
  /**
   * @return {Map<string, Array<workbox-routing.Route>>} routes A `Map` of HTTP
   * method name ('GET', etc.) to an array of all the corresponding `Route`
   * instances that are registered.
   */
  get routes() {
    return this._routes;
  }
  /**
   * Adds a fetch event listener to respond to events when a route matches
   * the event's request.
   */
  addFetchListener() {
    self.addEventListener("fetch", (e) => {
      const { request: t } = e, n = this.handleRequest({ request: t, event: e });
      n && e.respondWith(n);
    });
  }
  /**
   * Adds a message event listener for URLs to cache from the window.
   * This is useful to cache resources loaded on the page prior to when the
   * service worker started controlling it.
   *
   * The format of the message data sent from the window should be as follows.
   * Where the `urlsToCache` array may consist of URL strings or an array of
   * URL string + `requestInit` object (the same as you'd pass to `fetch()`).
   *
   * ```
   * {
   *   type: 'CACHE_URLS',
   *   payload: {
   *     urlsToCache: [
   *       './script1.js',
   *       './script2.js',
   *       ['./script3.js', {mode: 'no-cors'}],
   *     ],
   *   },
   * }
   * ```
   */
  addCacheListener() {
    self.addEventListener("message", (e) => {
      if (e.data && e.data.type === "CACHE_URLS") {
        const { payload: t } = e.data, n = Promise.all(t.urlsToCache.map((r) => {
          typeof r == "string" && (r = [r]);
          const i = new Request(...r);
          return this.handleRequest({ request: i, event: e });
        }));
        e.waitUntil(n), e.ports && e.ports[0] && n.then(() => e.ports[0].postMessage(!0));
      }
    });
  }
  /**
   * Apply the routing rules to a FetchEvent object to get a Response from an
   * appropriate Route's handler.
   *
   * @param {Object} options
   * @param {Request} options.request The request to handle.
   * @param {ExtendableEvent} options.event The event that triggered the
   *     request.
   * @return {Promise<Response>|undefined} A promise is returned if a
   *     registered route can handle the request. If there is no matching
   *     route and there's no `defaultHandler`, `undefined` is returned.
   */
  handleRequest({ request: e, event: t }) {
    const n = new URL(e.url, location.href);
    if (!n.protocol.startsWith("http"))
      return;
    const r = n.origin === location.origin, { params: i, route: a } = this.findMatchingRoute({
      event: t,
      request: e,
      sameOrigin: r,
      url: n
    });
    let o = a && a.handler;
    const l = e.method;
    if (!o && this._defaultHandlerMap.has(l) && (o = this._defaultHandlerMap.get(l)), !o)
      return;
    let c;
    try {
      c = o.handle({ url: n, request: e, event: t, params: i });
    } catch (h) {
      c = Promise.reject(h);
    }
    const u = a && a.catchHandler;
    return c instanceof Promise && (this._catchHandler || u) && (c = c.catch(async (h) => {
      if (u)
        try {
          return await u.handle({ url: n, request: e, event: t, params: i });
        } catch (m) {
          m instanceof Error && (h = m);
        }
      if (this._catchHandler)
        return this._catchHandler.handle({ url: n, request: e, event: t });
      throw h;
    })), c;
  }
  /**
   * Checks a request and URL (and optionally an event) against the list of
   * registered routes, and if there's a match, returns the corresponding
   * route along with any params generated by the match.
   *
   * @param {Object} options
   * @param {URL} options.url
   * @param {boolean} options.sameOrigin The result of comparing `url.origin`
   *     against the current origin.
   * @param {Request} options.request The request to match.
   * @param {Event} options.event The corresponding event.
   * @return {Object} An object with `route` and `params` properties.
   *     They are populated if a matching route was found or `undefined`
   *     otherwise.
   */
  findMatchingRoute({ url: e, sameOrigin: t, request: n, event: r }) {
    const i = this._routes.get(n.method) || [];
    for (const a of i) {
      let o;
      const l = a.match({ url: e, sameOrigin: t, request: n, event: r });
      if (l)
        return o = l, (Array.isArray(o) && o.length === 0 || l.constructor === Object && // eslint-disable-line
        Object.keys(l).length === 0 || typeof l == "boolean") && (o = void 0), { route: a, params: o };
    }
    return {};
  }
  /**
   * Define a default `handler` that's called when no routes explicitly
   * match the incoming request.
   *
   * Each HTTP method ('GET', 'POST', etc.) gets its own default handler.
   *
   * Without a default handler, unmatched requests will go against the
   * network as if there were no service worker present.
   *
   * @param {workbox-routing~handlerCallback} handler A callback
   * function that returns a Promise resulting in a Response.
   * @param {string} [method='GET'] The HTTP method to associate with this
   * default handler. Each method has its own default.
   */
  setDefaultHandler(e, t = Ee) {
    this._defaultHandlerMap.set(t, j(e));
  }
  /**
   * If a Route throws an error while handling a request, this `handler`
   * will be called and given a chance to provide a response.
   *
   * @param {workbox-routing~handlerCallback} handler A callback
   * function that returns a Promise resulting in a Response.
   */
  setCatchHandler(e) {
    this._catchHandler = j(e);
  }
  /**
   * Registers a route with the router.
   *
   * @param {workbox-routing.Route} route The route to register.
   */
  registerRoute(e) {
    this._routes.has(e.method) || this._routes.set(e.method, []), this._routes.get(e.method).push(e);
  }
  /**
   * Unregisters a route with the router.
   *
   * @param {workbox-routing.Route} route The route to unregister.
   */
  unregisterRoute(e) {
    if (!this._routes.has(e.method))
      throw new f("unregister-route-but-not-found-with-method", {
        method: e.method
      });
    const t = this._routes.get(e.method).indexOf(e);
    if (t > -1)
      this._routes.get(e.method).splice(t, 1);
    else
      throw new f("unregister-route-route-not-registered");
  }
}
function Qe(s) {
  _e.add(s);
}
function ve(s) {
  s.then(() => {
  });
}
function Ge(s) {
  O.updateDetails(s);
}
class le extends k {
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    let n = await t.cacheMatch(e), r;
    if (!n)
      try {
        n = await t.fetchAndCachePut(e);
      } catch (i) {
        i instanceof Error && (r = i);
      }
    if (!n)
      throw new f("no-response", { url: e.url, error: r });
    return n;
  }
}
class Je extends k {
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    const n = await t.cacheMatch(e);
    if (!n)
      throw new f("no-response", { url: e.url });
    return n;
  }
}
const Ce = {
  /**
   * Returns a valid response (to allow caching) if the status is 200 (OK) or
   * 0 (opaque).
   *
   * @param {Object} options
   * @param {Response} options.response
   * @return {Response|null}
   *
   * @private
   */
  cacheWillUpdate: async ({ response: s }) => s.status === 200 || s.status === 0 ? s : null
};
class Xe extends k {
  /**
   * @param {Object} [options]
   * @param {string} [options.cacheName] Cache name to store and retrieve
   * requests. Defaults to cache names provided by
   * {@link workbox-core.cacheNames}.
   * @param {Array<Object>} [options.plugins] [Plugins]{@link https://developers.google.com/web/tools/workbox/guides/using-plugins}
   * to use in conjunction with this caching strategy.
   * @param {Object} [options.fetchOptions] Values passed along to the
   * [`init`](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters)
   * of [non-navigation](https://github.com/GoogleChrome/workbox/issues/1796)
   * `fetch()` requests made by this strategy.
   * @param {Object} [options.matchOptions] [`CacheQueryOptions`](https://w3c.github.io/ServiceWorker/#dictdef-cachequeryoptions)
   * @param {number} [options.networkTimeoutSeconds] If set, any network requests
   * that fail to respond within the timeout will fallback to the cache.
   *
   * This option can be used to combat
   * "[lie-fi]{@link https://developers.google.com/web/fundamentals/performance/poor-connectivity/#lie-fi}"
   * scenarios.
   */
  constructor(e = {}) {
    super(e), this.plugins.some((t) => "cacheWillUpdate" in t) || this.plugins.unshift(Ce), this._networkTimeoutSeconds = e.networkTimeoutSeconds || 0;
  }
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    const n = [], r = [];
    let i;
    if (this._networkTimeoutSeconds) {
      const { id: l, promise: c } = this._getTimeoutPromise({ request: e, logs: n, handler: t });
      i = l, r.push(c);
    }
    const a = this._getNetworkPromise({
      timeoutId: i,
      request: e,
      logs: n,
      handler: t
    });
    r.push(a);
    const o = await t.waitUntil((async () => await t.waitUntil(Promise.race(r)) || // If Promise.race() resolved with null, it might be due to a network
    // timeout + a cache miss. If that were to happen, we'd rather wait until
    // the networkPromise resolves instead of returning null.
    // Note that it's fine to await an already-resolved promise, so we don't
    // have to check to see if it's still "in flight".
    await a)());
    if (!o)
      throw new f("no-response", { url: e.url });
    return o;
  }
  /**
   * @param {Object} options
   * @param {Request} options.request
   * @param {Array} options.logs A reference to the logs array
   * @param {Event} options.event
   * @return {Promise<Response>}
   *
   * @private
   */
  _getTimeoutPromise({ request: e, logs: t, handler: n }) {
    let r;
    return {
      promise: new Promise((a) => {
        r = setTimeout(async () => {
          a(await n.cacheMatch(e));
        }, this._networkTimeoutSeconds * 1e3);
      }),
      id: r
    };
  }
  /**
   * @param {Object} options
   * @param {number|undefined} options.timeoutId
   * @param {Request} options.request
   * @param {Array} options.logs A reference to the logs Array.
   * @param {Event} options.event
   * @return {Promise<Response>}
   *
   * @private
   */
  async _getNetworkPromise({ timeoutId: e, request: t, logs: n, handler: r }) {
    let i, a;
    try {
      a = await r.fetchAndCachePut(t);
    } catch (o) {
      o instanceof Error && (i = o);
    }
    return e && clearTimeout(e), (i || !a) && (a = await r.cacheMatch(t)), a;
  }
}
class Ye extends k {
  /**
   * @param {Object} [options]
   * @param {Array<Object>} [options.plugins] [Plugins]{@link https://developers.google.com/web/tools/workbox/guides/using-plugins}
   * to use in conjunction with this caching strategy.
   * @param {Object} [options.fetchOptions] Values passed along to the
   * [`init`](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters)
   * of [non-navigation](https://github.com/GoogleChrome/workbox/issues/1796)
   * `fetch()` requests made by this strategy.
   * @param {number} [options.networkTimeoutSeconds] If set, any network requests
   * that fail to respond within the timeout will result in a network error.
   */
  constructor(e = {}) {
    super(e), this._networkTimeoutSeconds = e.networkTimeoutSeconds || 0;
  }
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    let n, r;
    try {
      const i = [
        t.fetch(e)
      ];
      if (this._networkTimeoutSeconds) {
        const a = Re(this._networkTimeoutSeconds * 1e3);
        i.push(a);
      }
      if (r = await Promise.race(i), !r)
        throw new Error(`Timed out the network response after ${this._networkTimeoutSeconds} seconds.`);
    } catch (i) {
      i instanceof Error && (n = i);
    }
    if (!r)
      throw new f("no-response", { url: e.url, error: n });
    return r;
  }
}
class Ze extends k {
  /**
   * @param {Object} [options]
   * @param {string} [options.cacheName] Cache name to store and retrieve
   * requests. Defaults to cache names provided by
   * {@link workbox-core.cacheNames}.
   * @param {Array<Object>} [options.plugins] [Plugins]{@link https://developers.google.com/web/tools/workbox/guides/using-plugins}
   * to use in conjunction with this caching strategy.
   * @param {Object} [options.fetchOptions] Values passed along to the
   * [`init`](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters)
   * of [non-navigation](https://github.com/GoogleChrome/workbox/issues/1796)
   * `fetch()` requests made by this strategy.
   * @param {Object} [options.matchOptions] [`CacheQueryOptions`](https://w3c.github.io/ServiceWorker/#dictdef-cachequeryoptions)
   */
  constructor(e = {}) {
    super(e), this.plugins.some((t) => "cacheWillUpdate" in t) || this.plugins.unshift(Ce);
  }
  /**
   * @private
   * @param {Request|string} request A request to run this strategy for.
   * @param {workbox-strategies.StrategyHandler} handler The event that
   *     triggered the request.
   * @return {Promise<Response>}
   */
  async _handle(e, t) {
    const n = t.fetchAndCachePut(e).catch(() => {
    });
    t.waitUntil(n);
    let r = await t.cacheMatch(e), i;
    if (!r)
      try {
        r = await n;
      } catch (a) {
        a instanceof Error && (i = a);
      }
    if (!r)
      throw new f("no-response", { url: e.url, error: i });
    return r;
  }
}
const et = (s, e) => e.some((t) => s instanceof t);
let ue, he;
function tt() {
  return ue || (ue = [
    IDBDatabase,
    IDBObjectStore,
    IDBIndex,
    IDBCursor,
    IDBTransaction
  ]);
}
function st() {
  return he || (he = [
    IDBCursor.prototype.advance,
    IDBCursor.prototype.continue,
    IDBCursor.prototype.continuePrimaryKey
  ]);
}
const ke = /* @__PURE__ */ new WeakMap(), Q = /* @__PURE__ */ new WeakMap(), xe = /* @__PURE__ */ new WeakMap(), W = /* @__PURE__ */ new WeakMap(), Z = /* @__PURE__ */ new WeakMap();
function nt(s) {
  const e = new Promise((t, n) => {
    const r = () => {
      s.removeEventListener("success", i), s.removeEventListener("error", a);
    }, i = () => {
      t(w(s.result)), r();
    }, a = () => {
      n(s.error), r();
    };
    s.addEventListener("success", i), s.addEventListener("error", a);
  });
  return e.then((t) => {
    t instanceof IDBCursor && ke.set(t, s);
  }).catch(() => {
  }), Z.set(e, s), e;
}
function rt(s) {
  if (Q.has(s))
    return;
  const e = new Promise((t, n) => {
    const r = () => {
      s.removeEventListener("complete", i), s.removeEventListener("error", a), s.removeEventListener("abort", a);
    }, i = () => {
      t(), r();
    }, a = () => {
      n(s.error || new DOMException("AbortError", "AbortError")), r();
    };
    s.addEventListener("complete", i), s.addEventListener("error", a), s.addEventListener("abort", a);
  });
  Q.set(s, e);
}
let G = {
  get(s, e, t) {
    if (s instanceof IDBTransaction) {
      if (e === "done")
        return Q.get(s);
      if (e === "objectStoreNames")
        return s.objectStoreNames || xe.get(s);
      if (e === "store")
        return t.objectStoreNames[1] ? void 0 : t.objectStore(t.objectStoreNames[0]);
    }
    return w(s[e]);
  },
  set(s, e, t) {
    return s[e] = t, !0;
  },
  has(s, e) {
    return s instanceof IDBTransaction && (e === "done" || e === "store") ? !0 : e in s;
  }
};
function it(s) {
  G = s(G);
}
function at(s) {
  return s === IDBDatabase.prototype.transaction && !("objectStoreNames" in IDBTransaction.prototype) ? function(e, ...t) {
    const n = s.call(z(this), e, ...t);
    return xe.set(n, e.sort ? e.sort() : [e]), w(n);
  } : st().includes(s) ? function(...e) {
    return s.apply(z(this), e), w(ke.get(this));
  } : function(...e) {
    return w(s.apply(z(this), e));
  };
}
function ot(s) {
  return typeof s == "function" ? at(s) : (s instanceof IDBTransaction && rt(s), et(s, tt()) ? new Proxy(s, G) : s);
}
function w(s) {
  if (s instanceof IDBRequest)
    return nt(s);
  if (W.has(s))
    return W.get(s);
  const e = ot(s);
  return e !== s && (W.set(s, e), Z.set(e, s)), e;
}
const z = (s) => Z.get(s);
function De(s, e, { blocked: t, upgrade: n, blocking: r, terminated: i } = {}) {
  const a = indexedDB.open(s, e), o = w(a);
  return n && a.addEventListener("upgradeneeded", (l) => {
    n(w(a.result), l.oldVersion, l.newVersion, w(a.transaction), l);
  }), t && a.addEventListener("blocked", (l) => t(
    // Casting due to https://github.com/microsoft/TypeScript-DOM-lib-generator/pull/1405
    l.oldVersion,
    l.newVersion,
    l
  )), o.then((l) => {
    i && l.addEventListener("close", () => i()), r && l.addEventListener("versionchange", (c) => r(c.oldVersion, c.newVersion, c));
  }).catch(() => {
  }), o;
}
function ct(s, { blocked: e } = {}) {
  const t = indexedDB.deleteDatabase(s);
  return e && t.addEventListener("blocked", (n) => e(
    // Casting due to https://github.com/microsoft/TypeScript-DOM-lib-generator/pull/1405
    n.oldVersion,
    n
  )), w(t).then(() => {
  });
}
const lt = ["get", "getKey", "getAll", "getAllKeys", "count"], ut = ["put", "add", "delete", "clear"], H = /* @__PURE__ */ new Map();
function de(s, e) {
  if (!(s instanceof IDBDatabase && !(e in s) && typeof e == "string"))
    return;
  if (H.get(e))
    return H.get(e);
  const t = e.replace(/FromIndex$/, ""), n = e !== t, r = ut.includes(t);
  if (
    // Bail if the target doesn't exist on the target. Eg, getAll isn't in Edge.
    !(t in (n ? IDBIndex : IDBObjectStore).prototype) || !(r || lt.includes(t))
  )
    return;
  const i = async function(a, ...o) {
    const l = this.transaction(a, r ? "readwrite" : "readonly");
    let c = l.store;
    return n && (c = c.index(o.shift())), (await Promise.all([
      c[t](...o),
      r && l.done
    ]))[0];
  };
  return H.set(e, i), i;
}
it((s) => ({
  ...s,
  get: (e, t, n) => de(e, t) || s.get(e, t, n),
  has: (e, t) => !!de(e, t) || s.has(e, t)
}));
try {
  self["workbox:expiration:7.0.0"] && _();
} catch {
}
const ht = "workbox-expiration", D = "cache-entries", fe = (s) => {
  const e = new URL(s, location.href);
  return e.hash = "", e.href;
};
class dt {
  /**
   *
   * @param {string} cacheName
   *
   * @private
   */
  constructor(e) {
    this._db = null, this._cacheName = e;
  }
  /**
   * Performs an upgrade of indexedDB.
   *
   * @param {IDBPDatabase<CacheDbSchema>} db
   *
   * @private
   */
  _upgradeDb(e) {
    const t = e.createObjectStore(D, { keyPath: "id" });
    t.createIndex("cacheName", "cacheName", { unique: !1 }), t.createIndex("timestamp", "timestamp", { unique: !1 });
  }
  /**
   * Performs an upgrade of indexedDB and deletes deprecated DBs.
   *
   * @param {IDBPDatabase<CacheDbSchema>} db
   *
   * @private
   */
  _upgradeDbAndDeleteOldDbs(e) {
    this._upgradeDb(e), this._cacheName && ct(this._cacheName);
  }
  /**
   * @param {string} url
   * @param {number} timestamp
   *
   * @private
   */
  async setTimestamp(e, t) {
    e = fe(e);
    const n = {
      url: e,
      timestamp: t,
      cacheName: this._cacheName,
      // Creating an ID from the URL and cache name won't be necessary once
      // Edge switches to Chromium and all browsers we support work with
      // array keyPaths.
      id: this._getId(e)
    }, i = (await this.getDb()).transaction(D, "readwrite", {
      durability: "relaxed"
    });
    await i.store.put(n), await i.done;
  }
  /**
   * Returns the timestamp stored for a given URL.
   *
   * @param {string} url
   * @return {number | undefined}
   *
   * @private
   */
  async getTimestamp(e) {
    const n = await (await this.getDb()).get(D, this._getId(e));
    return n == null ? void 0 : n.timestamp;
  }
  /**
   * Iterates through all the entries in the object store (from newest to
   * oldest) and removes entries once either `maxCount` is reached or the
   * entry's timestamp is less than `minTimestamp`.
   *
   * @param {number} minTimestamp
   * @param {number} maxCount
   * @return {Array<string>}
   *
   * @private
   */
  async expireEntries(e, t) {
    const n = await this.getDb();
    let r = await n.transaction(D).store.index("timestamp").openCursor(null, "prev");
    const i = [];
    let a = 0;
    for (; r; ) {
      const l = r.value;
      l.cacheName === this._cacheName && (e && l.timestamp < e || t && a >= t ? i.push(r.value) : a++), r = await r.continue();
    }
    const o = [];
    for (const l of i)
      await n.delete(D, l.id), o.push(l.url);
    return o;
  }
  /**
   * Takes a URL and returns an ID that will be unique in the object store.
   *
   * @param {string} url
   * @return {string}
   *
   * @private
   */
  _getId(e) {
    return this._cacheName + "|" + fe(e);
  }
  /**
   * Returns an open connection to the database.
   *
   * @private
   */
  async getDb() {
    return this._db || (this._db = await De(ht, 1, {
      upgrade: this._upgradeDbAndDeleteOldDbs.bind(this)
    })), this._db;
  }
}
class ft {
  /**
   * To construct a new CacheExpiration instance you must provide at least
   * one of the `config` properties.
   *
   * @param {string} cacheName Name of the cache to apply restrictions to.
   * @param {Object} config
   * @param {number} [config.maxEntries] The maximum number of entries to cache.
   * Entries used the least will be removed as the maximum is reached.
   * @param {number} [config.maxAgeSeconds] The maximum age of an entry before
   * it's treated as stale and removed.
   * @param {Object} [config.matchOptions] The [`CacheQueryOptions`](https://developer.mozilla.org/en-US/docs/Web/API/Cache/delete#Parameters)
   * that will be used when calling `delete()` on the cache.
   */
  constructor(e, t = {}) {
    this._isRunning = !1, this._rerunRequested = !1, this._maxEntries = t.maxEntries, this._maxAgeSeconds = t.maxAgeSeconds, this._matchOptions = t.matchOptions, this._cacheName = e, this._timestampModel = new dt(e);
  }
  /**
   * Expires entries for the given cache and given criteria.
   */
  async expireEntries() {
    if (this._isRunning) {
      this._rerunRequested = !0;
      return;
    }
    this._isRunning = !0;
    const e = this._maxAgeSeconds ? Date.now() - this._maxAgeSeconds * 1e3 : 0, t = await this._timestampModel.expireEntries(e, this._maxEntries), n = await self.caches.open(this._cacheName);
    for (const r of t)
      await n.delete(r, this._matchOptions);
    this._isRunning = !1, this._rerunRequested && (this._rerunRequested = !1, ve(this.expireEntries()));
  }
  /**
   * Update the timestamp for the given URL. This ensures the when
   * removing entries based on maximum entries, most recently used
   * is accurate or when expiring, the timestamp is up-to-date.
   *
   * @param {string} url
   */
  async updateTimestamp(e) {
    await this._timestampModel.setTimestamp(e, Date.now());
  }
  /**
   * Can be used to check if a URL has expired or not before it's used.
   *
   * This requires a look up from IndexedDB, so can be slow.
   *
   * Note: This method will not remove the cached entry, call
   * `expireEntries()` to remove indexedDB and Cache entries.
   *
   * @param {string} url
   * @return {boolean}
   */
  async isURLExpired(e) {
    if (this._maxAgeSeconds) {
      const t = await this._timestampModel.getTimestamp(e), n = Date.now() - this._maxAgeSeconds * 1e3;
      return t !== void 0 ? t < n : !0;
    } else
      return !1;
  }
  /**
   * Removes the IndexedDB object store used to keep track of cache expiration
   * metadata.
   */
  async delete() {
    this._rerunRequested = !1, await this._timestampModel.expireEntries(1 / 0);
  }
}
class mt {
  /**
   * @param {ExpirationPluginOptions} config
   * @param {number} [config.maxEntries] The maximum number of entries to cache.
   * Entries used the least will be removed as the maximum is reached.
   * @param {number} [config.maxAgeSeconds] The maximum age of an entry before
   * it's treated as stale and removed.
   * @param {Object} [config.matchOptions] The [`CacheQueryOptions`](https://developer.mozilla.org/en-US/docs/Web/API/Cache/delete#Parameters)
   * that will be used when calling `delete()` on the cache.
   * @param {boolean} [config.purgeOnQuotaError] Whether to opt this cache in to
   * automatic deletion if the available storage quota has been exceeded.
   */
  constructor(e = {}) {
    this.cachedResponseWillBeUsed = async ({ event: t, request: n, cacheName: r, cachedResponse: i }) => {
      if (!i)
        return null;
      const a = this._isResponseDateFresh(i), o = this._getCacheExpiration(r);
      ve(o.expireEntries());
      const l = o.updateTimestamp(n.url);
      if (t)
        try {
          t.waitUntil(l);
        } catch {
        }
      return a ? i : null;
    }, this.cacheDidUpdate = async ({ cacheName: t, request: n }) => {
      const r = this._getCacheExpiration(t);
      await r.updateTimestamp(n.url), await r.expireEntries();
    }, this._config = e, this._maxAgeSeconds = e.maxAgeSeconds, this._cacheExpirations = /* @__PURE__ */ new Map(), e.purgeOnQuotaError && Qe(() => this.deleteCacheAndMetadata());
  }
  /**
   * A simple helper method to return a CacheExpiration instance for a given
   * cache name.
   *
   * @param {string} cacheName
   * @return {CacheExpiration}
   *
   * @private
   */
  _getCacheExpiration(e) {
    if (e === O.getRuntimeName())
      throw new f("expire-custom-caches-only");
    let t = this._cacheExpirations.get(e);
    return t || (t = new ft(e, this._config), this._cacheExpirations.set(e, t)), t;
  }
  /**
   * @param {Response} cachedResponse
   * @return {boolean}
   *
   * @private
   */
  _isResponseDateFresh(e) {
    if (!this._maxAgeSeconds)
      return !0;
    const t = this._getDateHeaderTimestamp(e);
    if (t === null)
      return !0;
    const n = Date.now();
    return t >= n - this._maxAgeSeconds * 1e3;
  }
  /**
   * This method will extract the data header and parse it into a useful
   * value.
   *
   * @param {Response} cachedResponse
   * @return {number|null}
   *
   * @private
   */
  _getDateHeaderTimestamp(e) {
    if (!e.headers.has("date"))
      return null;
    const t = e.headers.get("date"), r = new Date(t).getTime();
    return isNaN(r) ? null : r;
  }
  /**
   * This is a helper method that performs two operations:
   *
   * - Deletes *all* the underlying Cache instances associated with this plugin
   * instance, by calling caches.delete() on your behalf.
   * - Deletes the metadata from IndexedDB used to keep track of expiration
   * details for each Cache instance.
   *
   * When using cache expiration, calling this method is preferable to calling
   * `caches.delete()` directly, since this will ensure that the IndexedDB
   * metadata is also cleanly removed and open IndexedDB instances are deleted.
   *
   * Note that if you're *not* using cache expiration for a given cache, calling
   * `caches.delete()` and passing in the cache's name should be sufficient.
   * There is no Workbox-specific method needed for cleanup in that case.
   */
  async deleteCacheAndMetadata() {
    for (const [e, t] of this._cacheExpirations)
      await self.caches.delete(e), await t.delete();
    this._cacheExpirations = /* @__PURE__ */ new Map();
  }
}
try {
  self["workbox:cacheable-response:7.0.0"] && _();
} catch {
}
class pt {
  /**
   * To construct a new CacheableResponse instance you must provide at least
   * one of the `config` properties.
   *
   * If both `statuses` and `headers` are specified, then both conditions must
   * be met for the `Response` to be considered cacheable.
   *
   * @param {Object} config
   * @param {Array<number>} [config.statuses] One or more status codes that a
   * `Response` can have and be considered cacheable.
   * @param {Object<string,string>} [config.headers] A mapping of header names
   * and expected values that a `Response` can have and be considered cacheable.
   * If multiple headers are provided, only one needs to be present.
   */
  constructor(e = {}) {
    this._statuses = e.statuses, this._headers = e.headers;
  }
  /**
   * Checks a response to see whether it's cacheable or not, based on this
   * object's configuration.
   *
   * @param {Response} response The response whose cacheability is being
   * checked.
   * @return {boolean} `true` if the `Response` is cacheable, and `false`
   * otherwise.
   */
  isResponseCacheable(e) {
    let t = !0;
    return this._statuses && (t = this._statuses.includes(e.status)), this._headers && t && (t = Object.keys(this._headers).some((n) => e.headers.get(n) === this._headers[n])), t;
  }
}
class yt {
  /**
   * To construct a new CacheableResponsePlugin instance you must provide at
   * least one of the `config` properties.
   *
   * If both `statuses` and `headers` are specified, then both conditions must
   * be met for the `Response` to be considered cacheable.
   *
   * @param {Object} config
   * @param {Array<number>} [config.statuses] One or more status codes that a
   * `Response` can have and be considered cacheable.
   * @param {Object<string,string>} [config.headers] A mapping of header names
   * and expected values that a `Response` can have and be considered cacheable.
   * If multiple headers are provided, only one needs to be present.
   */
  constructor(e) {
    this.cacheWillUpdate = async ({ response: t }) => this._cacheableResponse.isResponseCacheable(t) ? t : null, this._cacheableResponse = new pt(e);
  }
}
try {
  self["workbox:background-sync:7.0.0"] && _();
} catch {
}
const me = 3, gt = "workbox-background-sync", y = "requests", T = "queueName";
class wt {
  constructor() {
    this._db = null;
  }
  /**
   * Add QueueStoreEntry to underlying db.
   *
   * @param {UnidentifiedQueueStoreEntry} entry
   */
  async addEntry(e) {
    const n = (await this.getDb()).transaction(y, "readwrite", {
      durability: "relaxed"
    });
    await n.store.add(e), await n.done;
  }
  /**
   * Returns the first entry id in the ObjectStore.
   *
   * @return {number | undefined}
   */
  async getFirstEntryId() {
    const t = await (await this.getDb()).transaction(y).store.openCursor();
    return t == null ? void 0 : t.value.id;
  }
  /**
   * Get all the entries filtered by index
   *
   * @param queueName
   * @return {Promise<QueueStoreEntry[]>}
   */
  async getAllEntriesByQueueName(e) {
    const n = await (await this.getDb()).getAllFromIndex(y, T, IDBKeyRange.only(e));
    return n || new Array();
  }
  /**
   * Returns the number of entries filtered by index
   *
   * @param queueName
   * @return {Promise<number>}
   */
  async getEntryCountByQueueName(e) {
    return (await this.getDb()).countFromIndex(y, T, IDBKeyRange.only(e));
  }
  /**
   * Deletes a single entry by id.
   *
   * @param {number} id the id of the entry to be deleted
   */
  async deleteEntry(e) {
    await (await this.getDb()).delete(y, e);
  }
  /**
   *
   * @param queueName
   * @returns {Promise<QueueStoreEntry | undefined>}
   */
  async getFirstEntryByQueueName(e) {
    return await this.getEndEntryFromIndex(IDBKeyRange.only(e), "next");
  }
  /**
   *
   * @param queueName
   * @returns {Promise<QueueStoreEntry | undefined>}
   */
  async getLastEntryByQueueName(e) {
    return await this.getEndEntryFromIndex(IDBKeyRange.only(e), "prev");
  }
  /**
   * Returns either the first or the last entries, depending on direction.
   * Filtered by index.
   *
   * @param {IDBCursorDirection} direction
   * @param {IDBKeyRange} query
   * @return {Promise<QueueStoreEntry | undefined>}
   * @private
   */
  async getEndEntryFromIndex(e, t) {
    const r = await (await this.getDb()).transaction(y).store.index(T).openCursor(e, t);
    return r == null ? void 0 : r.value;
  }
  /**
   * Returns an open connection to the database.
   *
   * @private
   */
  async getDb() {
    return this._db || (this._db = await De(gt, me, {
      upgrade: this._upgradeDb
    })), this._db;
  }
  /**
   * Upgrades QueueDB
   *
   * @param {IDBPDatabase<QueueDBSchema>} db
   * @param {number} oldVersion
   * @private
   */
  _upgradeDb(e, t) {
    t > 0 && t < me && e.objectStoreNames.contains(y) && e.deleteObjectStore(y), e.createObjectStore(y, {
      autoIncrement: !0,
      keyPath: "id"
    }).createIndex(T, T, { unique: !1 });
  }
}
class bt {
  /**
   * Associates this instance with a Queue instance, so entries added can be
   * identified by their queue name.
   *
   * @param {string} queueName
   */
  constructor(e) {
    this._queueName = e, this._queueDb = new wt();
  }
  /**
   * Append an entry last in the queue.
   *
   * @param {Object} entry
   * @param {Object} entry.requestData
   * @param {number} [entry.timestamp]
   * @param {Object} [entry.metadata]
   */
  async pushEntry(e) {
    delete e.id, e.queueName = this._queueName, await this._queueDb.addEntry(e);
  }
  /**
   * Prepend an entry first in the queue.
   *
   * @param {Object} entry
   * @param {Object} entry.requestData
   * @param {number} [entry.timestamp]
   * @param {Object} [entry.metadata]
   */
  async unshiftEntry(e) {
    const t = await this._queueDb.getFirstEntryId();
    t ? e.id = t - 1 : delete e.id, e.queueName = this._queueName, await this._queueDb.addEntry(e);
  }
  /**
   * Removes and returns the last entry in the queue matching the `queueName`.
   *
   * @return {Promise<QueueStoreEntry|undefined>}
   */
  async popEntry() {
    return this._removeEntry(await this._queueDb.getLastEntryByQueueName(this._queueName));
  }
  /**
   * Removes and returns the first entry in the queue matching the `queueName`.
   *
   * @return {Promise<QueueStoreEntry|undefined>}
   */
  async shiftEntry() {
    return this._removeEntry(await this._queueDb.getFirstEntryByQueueName(this._queueName));
  }
  /**
   * Returns all entries in the store matching the `queueName`.
   *
   * @param {Object} options See {@link workbox-background-sync.Queue~getAll}
   * @return {Promise<Array<Object>>}
   */
  async getAll() {
    return await this._queueDb.getAllEntriesByQueueName(this._queueName);
  }
  /**
   * Returns the number of entries in the store matching the `queueName`.
   *
   * @param {Object} options See {@link workbox-background-sync.Queue~size}
   * @return {Promise<number>}
   */
  async size() {
    return await this._queueDb.getEntryCountByQueueName(this._queueName);
  }
  /**
   * Deletes the entry for the given ID.
   *
   * WARNING: this method does not ensure the deleted entry belongs to this
   * queue (i.e. matches the `queueName`). But this limitation is acceptable
   * as this class is not publicly exposed. An additional check would make
   * this method slower than it needs to be.
   *
   * @param {number} id
   */
  async deleteEntry(e) {
    await this._queueDb.deleteEntry(e);
  }
  /**
   * Removes and returns the first or last entry in the queue (based on the
   * `direction` argument) matching the `queueName`.
   *
   * @return {Promise<QueueStoreEntry|undefined>}
   * @private
   */
  async _removeEntry(e) {
    return e && await this.deleteEntry(e.id), e;
  }
}
const _t = [
  "method",
  "referrer",
  "referrerPolicy",
  "mode",
  "credentials",
  "cache",
  "redirect",
  "integrity",
  "keepalive"
];
class S {
  /**
   * Converts a Request object to a plain object that can be structured
   * cloned or JSON-stringified.
   *
   * @param {Request} request
   * @return {Promise<StorableRequest>}
   */
  static async fromRequest(e) {
    const t = {
      url: e.url,
      headers: {}
    };
    e.method !== "GET" && (t.body = await e.clone().arrayBuffer());
    for (const [n, r] of e.headers.entries())
      t.headers[n] = r;
    for (const n of _t)
      e[n] !== void 0 && (t[n] = e[n]);
    return new S(t);
  }
  /**
   * Accepts an object of request data that can be used to construct a
   * `Request` but can also be stored in IndexedDB.
   *
   * @param {Object} requestData An object of request data that includes the
   *     `url` plus any relevant properties of
   *     [requestInit]{@link https://fetch.spec.whatwg.org/#requestinit}.
   */
  constructor(e) {
    e.mode === "navigate" && (e.mode = "same-origin"), this._requestData = e;
  }
  /**
   * Returns a deep clone of the instances `_requestData` object.
   *
   * @return {Object}
   */
  toObject() {
    const e = Object.assign({}, this._requestData);
    return e.headers = Object.assign({}, this._requestData.headers), e.body && (e.body = e.body.slice(0)), e;
  }
  /**
   * Converts this instance to a Request.
   *
   * @return {Request}
   */
  toRequest() {
    return new Request(this._requestData.url, this._requestData);
  }
  /**
   * Creates and returns a deep clone of the instance.
   *
   * @return {StorableRequest}
   */
  clone() {
    return new S(this.toObject());
  }
}
const pe = "workbox-background-sync", Rt = 60 * 24 * 7, V = /* @__PURE__ */ new Set(), ye = (s) => {
  const e = {
    request: new S(s.requestData).toRequest(),
    timestamp: s.timestamp
  };
  return s.metadata && (e.metadata = s.metadata), e;
};
class Et {
  /**
   * Creates an instance of Queue with the given options
   *
   * @param {string} name The unique name for this queue. This name must be
   *     unique as it's used to register sync events and store requests
   *     in IndexedDB specific to this instance. An error will be thrown if
   *     a duplicate name is detected.
   * @param {Object} [options]
   * @param {Function} [options.onSync] A function that gets invoked whenever
   *     the 'sync' event fires. The function is invoked with an object
   *     containing the `queue` property (referencing this instance), and you
   *     can use the callback to customize the replay behavior of the queue.
   *     When not set the `replayRequests()` method is called.
   *     Note: if the replay fails after a sync event, make sure you throw an
   *     error, so the browser knows to retry the sync event later.
   * @param {number} [options.maxRetentionTime=7 days] The amount of time (in
   *     minutes) a request may be retried. After this amount of time has
   *     passed, the request will be deleted from the queue.
   * @param {boolean} [options.forceSyncFallback=false] If `true`, instead
   *     of attempting to use background sync events, always attempt to replay
   *     queued request at service worker startup. Most folks will not need
   *     this, unless you explicitly target a runtime like Electron that
   *     exposes the interfaces for background sync, but does not have a working
   *     implementation.
   */
  constructor(e, { forceSyncFallback: t, onSync: n, maxRetentionTime: r } = {}) {
    if (this._syncInProgress = !1, this._requestsAddedDuringSync = !1, V.has(e))
      throw new f("duplicate-queue-name", { name: e });
    V.add(e), this._name = e, this._onSync = n || this.replayRequests, this._maxRetentionTime = r || Rt, this._forceSyncFallback = !!t, this._queueStore = new bt(this._name), this._addSyncListener();
  }
  /**
   * @return {string}
   */
  get name() {
    return this._name;
  }
  /**
   * Stores the passed request in IndexedDB (with its timestamp and any
   * metadata) at the end of the queue.
   *
   * @param {QueueEntry} entry
   * @param {Request} entry.request The request to store in the queue.
   * @param {Object} [entry.metadata] Any metadata you want associated with the
   *     stored request. When requests are replayed you'll have access to this
   *     metadata object in case you need to modify the request beforehand.
   * @param {number} [entry.timestamp] The timestamp (Epoch time in
   *     milliseconds) when the request was first added to the queue. This is
   *     used along with `maxRetentionTime` to remove outdated requests. In
   *     general you don't need to set this value, as it's automatically set
   *     for you (defaulting to `Date.now()`), but you can update it if you
   *     don't want particular requests to expire.
   */
  async pushRequest(e) {
    await this._addRequest(e, "push");
  }
  /**
   * Stores the passed request in IndexedDB (with its timestamp and any
   * metadata) at the beginning of the queue.
   *
   * @param {QueueEntry} entry
   * @param {Request} entry.request The request to store in the queue.
   * @param {Object} [entry.metadata] Any metadata you want associated with the
   *     stored request. When requests are replayed you'll have access to this
   *     metadata object in case you need to modify the request beforehand.
   * @param {number} [entry.timestamp] The timestamp (Epoch time in
   *     milliseconds) when the request was first added to the queue. This is
   *     used along with `maxRetentionTime` to remove outdated requests. In
   *     general you don't need to set this value, as it's automatically set
   *     for you (defaulting to `Date.now()`), but you can update it if you
   *     don't want particular requests to expire.
   */
  async unshiftRequest(e) {
    await this._addRequest(e, "unshift");
  }
  /**
   * Removes and returns the last request in the queue (along with its
   * timestamp and any metadata). The returned object takes the form:
   * `{request, timestamp, metadata}`.
   *
   * @return {Promise<QueueEntry | undefined>}
   */
  async popRequest() {
    return this._removeRequest("pop");
  }
  /**
   * Removes and returns the first request in the queue (along with its
   * timestamp and any metadata). The returned object takes the form:
   * `{request, timestamp, metadata}`.
   *
   * @return {Promise<QueueEntry | undefined>}
   */
  async shiftRequest() {
    return this._removeRequest("shift");
  }
  /**
   * Returns all the entries that have not expired (per `maxRetentionTime`).
   * Any expired entries are removed from the queue.
   *
   * @return {Promise<Array<QueueEntry>>}
   */
  async getAll() {
    const e = await this._queueStore.getAll(), t = Date.now(), n = [];
    for (const r of e) {
      const i = this._maxRetentionTime * 60 * 1e3;
      t - r.timestamp > i ? await this._queueStore.deleteEntry(r.id) : n.push(ye(r));
    }
    return n;
  }
  /**
   * Returns the number of entries present in the queue.
   * Note that expired entries (per `maxRetentionTime`) are also included in this count.
   *
   * @return {Promise<number>}
   */
  async size() {
    return await this._queueStore.size();
  }
  /**
   * Adds the entry to the QueueStore and registers for a sync event.
   *
   * @param {Object} entry
   * @param {Request} entry.request
   * @param {Object} [entry.metadata]
   * @param {number} [entry.timestamp=Date.now()]
   * @param {string} operation ('push' or 'unshift')
   * @private
   */
  async _addRequest({ request: e, metadata: t, timestamp: n = Date.now() }, r) {
    const a = {
      requestData: (await S.fromRequest(e.clone())).toObject(),
      timestamp: n
    };
    switch (t && (a.metadata = t), r) {
      case "push":
        await this._queueStore.pushEntry(a);
        break;
      case "unshift":
        await this._queueStore.unshiftEntry(a);
        break;
    }
    this._syncInProgress ? this._requestsAddedDuringSync = !0 : await this.registerSync();
  }
  /**
   * Removes and returns the first or last (depending on `operation`) entry
   * from the QueueStore that's not older than the `maxRetentionTime`.
   *
   * @param {string} operation ('pop' or 'shift')
   * @return {Object|undefined}
   * @private
   */
  async _removeRequest(e) {
    const t = Date.now();
    let n;
    switch (e) {
      case "pop":
        n = await this._queueStore.popEntry();
        break;
      case "shift":
        n = await this._queueStore.shiftEntry();
        break;
    }
    if (n) {
      const r = this._maxRetentionTime * 60 * 1e3;
      return t - n.timestamp > r ? this._removeRequest(e) : ye(n);
    } else
      return;
  }
  /**
   * Loops through each request in the queue and attempts to re-fetch it.
   * If any request fails to re-fetch, it's put back in the same position in
   * the queue (which registers a retry for the next sync event).
   */
  async replayRequests() {
    let e;
    for (; e = await this.shiftRequest(); )
      try {
        await fetch(e.request.clone());
      } catch {
        throw await this.unshiftRequest(e), new f("queue-replay-failed", { name: this._name });
      }
  }
  /**
   * Registers a sync event with a tag unique to this instance.
   */
  async registerSync() {
    if ("sync" in self.registration && !this._forceSyncFallback)
      try {
        await self.registration.sync.register(`${pe}:${this._name}`);
      } catch {
      }
  }
  /**
   * In sync-supporting browsers, this adds a listener for the sync event.
   * In non-sync-supporting browsers, or if _forceSyncFallback is true, this
   * will retry the queue on service worker startup.
   *
   * @private
   */
  _addSyncListener() {
    "sync" in self.registration && !this._forceSyncFallback ? self.addEventListener("sync", (e) => {
      if (e.tag === `${pe}:${this._name}`) {
        const t = async () => {
          this._syncInProgress = !0;
          let n;
          try {
            await this._onSync({ queue: this });
          } catch (r) {
            if (r instanceof Error)
              throw n = r, n;
          } finally {
            this._requestsAddedDuringSync && !(n && !e.lastChance) && await this.registerSync(), this._syncInProgress = !1, this._requestsAddedDuringSync = !1;
          }
        };
        e.waitUntil(t());
      }
    }) : this._onSync({ queue: this });
  }
  /**
   * Returns the set of queue names. This is primarily used to reset the list
   * of queue names in tests.
   *
   * @return {Set<string>}
   *
   * @private
   */
  static get _queueNames() {
    return V;
  }
}
class vt {
  /**
   * @param {string} name See the {@link workbox-background-sync.Queue}
   *     documentation for parameter details.
   * @param {Object} [options] See the
   *     {@link workbox-background-sync.Queue} documentation for
   *     parameter details.
   */
  constructor(e, t) {
    this.fetchDidFail = async ({ request: n }) => {
      await this._queue.pushRequest({ request: n });
    }, this._queue = new Et(e, t);
  }
}
class ee extends Ve {
  /**
   * Creates a new instance of the PolyfeaRoute class.
   * @param route - The Polyfea route options.
   */
  constructor(e) {
    const t = new RegExp(e.pattern || ".*");
    let n = e.prefix || "";
    if (n) {
      let o = decodeURIComponent(
        new URL(globalThis.location.href).searchParams.get("base-path") || ""
      );
      o || (o = new URL(globalThis.location.href).pathname.split("/").slice(0, -1).join("/")), n = new URL(n, `http://host${o}/`).pathname;
    }
    let r;
    const i = [];
    switch (e.strategy !== "network-only" && i.push(new yt({ statuses: e.statuses || [0, 200, 201, 202, 204] })), e.maxAgeSeconds && i.push(new mt({
      maxAgeSeconds: e.maxAgeSeconds
    })), e.syncRetentionMinutes && i.push(new vt("polyfea", {
      maxRetentionTime: e.syncRetentionMinutes
    })), e.strategy) {
      case "cache-first":
        r = new le({ plugins: i });
        break;
      case "cache-only":
        r = new Je({ plugins: i });
        break;
      case "network-first":
        r = new Xe({ plugins: i });
        break;
      case "network-only":
        r = new Ye({ plugins: i });
        break;
      case "stale-while-revalidate":
        r = new Ze({ plugins: i });
        break;
      default:
        r = new le({ plugins: i });
        break;
    }
    super(t, r, e.method || "GET");
    const a = this.match;
    this.match = (o) => e.destination && o.request.destination !== e.destination || e.prefix && !o.url.pathname.startsWith(e.prefix) ? !1 : a(o);
  }
  /**
   * Creates a new instance of the PolyfeaRoute class from the given route options.
   * @param route - The Polyfea route options.
   * @returns A new instance of the PolyfeaRoute class.
   */
  static from(e) {
    return new ee(e);
  }
}
var F = { exports: {} };
function Ct(s) {
  try {
    return JSON.stringify(s);
  } catch {
    return '"[Circular]"';
  }
}
var kt = xt;
function xt(s, e, t) {
  var n = t && t.stringify || Ct, r = 1;
  if (typeof s == "object" && s !== null) {
    var i = e.length + r;
    if (i === 1)
      return s;
    var a = new Array(i);
    a[0] = n(s);
    for (var o = 1; o < i; o++)
      a[o] = n(e[o]);
    return a.join(" ");
  }
  if (typeof s != "string")
    return s;
  var l = e.length;
  if (l === 0)
    return s;
  for (var c = "", u = 1 - r, h = -1, m = s && s.length || 0, d = 0; d < m; ) {
    if (s.charCodeAt(d) === 37 && d + 1 < m) {
      switch (h = h > -1 ? h : 0, s.charCodeAt(d + 1)) {
        case 100:
        case 102:
          if (u >= l || e[u] == null)
            break;
          h < d && (c += s.slice(h, d)), c += Number(e[u]), h = d + 2, d++;
          break;
        case 105:
          if (u >= l || e[u] == null)
            break;
          h < d && (c += s.slice(h, d)), c += Math.floor(Number(e[u])), h = d + 2, d++;
          break;
        case 79:
        case 111:
        case 106:
          if (u >= l || e[u] === void 0)
            break;
          h < d && (c += s.slice(h, d));
          var P = typeof e[u];
          if (P === "string") {
            c += "'" + e[u] + "'", h = d + 2, d++;
            break;
          }
          if (P === "function") {
            c += e[u].name || "<anonymous>", h = d + 2, d++;
            break;
          }
          c += n(e[u]), h = d + 2, d++;
          break;
        case 115:
          if (u >= l)
            break;
          h < d && (c += s.slice(h, d)), c += String(e[u]), h = d + 2, d++;
          break;
        case 37:
          h < d && (c += s.slice(h, d)), c += "%", h = d + 2, d++, u--;
          break;
      }
      ++u;
    }
    ++d;
  }
  return h === -1 ? s : (h < m && (c += s.slice(h)), c);
}
const ge = kt;
F.exports = b;
const L = Bt().console || {}, Dt = {
  mapHttpRequest: U,
  mapHttpResponse: U,
  wrapRequestSerializer: $,
  wrapResponseSerializer: $,
  wrapErrorSerializer: $,
  req: U,
  res: U,
  err: be,
  errWithCause: be
};
function M(s, e) {
  return s === "silent" ? 1 / 0 : e.levels.values[s];
}
const te = Symbol("pino.logFuncs"), J = Symbol("pino.hierarchy"), Tt = {
  error: "log",
  fatal: "error",
  warn: "error",
  info: "log",
  debug: "log",
  trace: "log"
};
function we(s, e) {
  const t = {
    logger: e,
    parent: s[J]
  };
  e[J] = t;
}
function St(s, e, t) {
  const n = {};
  e.forEach((r) => {
    n[r] = t[r] ? t[r] : L[r] || L[Tt[r] || "log"] || q;
  }), s[te] = n;
}
function Lt(s, e) {
  return Array.isArray(s) ? s.filter(function(n) {
    return n !== "!stdSerializers.err";
  }) : s === !0 ? Object.keys(e) : !1;
}
function b(s) {
  s = s || {}, s.browser = s.browser || {};
  const e = s.browser.transmit;
  if (e && typeof e.send != "function")
    throw Error("pino: transmit option must have a send function");
  const t = s.browser.write || L;
  s.browser.write && (s.browser.asObject = !0);
  const n = s.serializers || {}, r = Lt(s.browser.serialize, n);
  let i = s.browser.serialize;
  Array.isArray(s.browser.serialize) && s.browser.serialize.indexOf("!stdSerializers.err") > -1 && (i = !1);
  const a = Object.keys(s.customLevels || {}), o = ["error", "fatal", "warn", "info", "debug", "trace"].concat(a);
  typeof t == "function" && o.forEach(function(p) {
    t[p] = t;
  }), (s.enabled === !1 || s.browser.disabled) && (s.level = "silent");
  const l = s.level || "info", c = Object.create(t);
  c.log || (c.log = q), St(c, o, t), we({}, c), Object.defineProperty(c, "levelVal", {
    get: h
  }), Object.defineProperty(c, "level", {
    get: m,
    set: d
  });
  const u = {
    transmit: e,
    serialize: r,
    asObject: s.browser.asObject,
    formatters: s.browser.formatters,
    levels: o,
    timestamp: jt(s)
  };
  c.levels = qt(s), c.level = l, c.setMaxListeners = c.getMaxListeners = c.emit = c.addListener = c.on = c.prependListener = c.once = c.prependOnceListener = c.removeListener = c.removeAllListeners = c.listeners = c.listenerCount = c.eventNames = c.write = c.flush = q, c.serializers = n, c._serialize = r, c._stdErrSerialize = i, c.child = P, e && (c._logEvent = X());
  function h() {
    return M(this.level, this);
  }
  function m() {
    return this._level;
  }
  function d(p) {
    if (p !== "silent" && !this.levels.values[p])
      throw Error("unknown level " + p);
    this._level = p, v(this, u, c, "error"), v(this, u, c, "fatal"), v(this, u, c, "warn"), v(this, u, c, "info"), v(this, u, c, "debug"), v(this, u, c, "trace"), a.forEach((C) => {
      v(this, u, c, C);
    });
  }
  function P(p, C) {
    if (!p)
      throw new Error("missing bindings for child Pino");
    C = C || {}, r && p.serializers && (C.serializers = p.serializers);
    const ne = C.serializers;
    if (r && ne) {
      var I = Object.assign({}, n, ne), re = s.browser.serialize === !0 ? Object.keys(I) : r;
      delete p.serializers, se([p], re, I, this._stdErrSerialize);
    }
    function ie(ae) {
      this._childLevel = (ae._childLevel | 0) + 1, this.bindings = p, I && (this.serializers = I, this._serialize = re), e && (this._logEvent = X(
        [].concat(ae._logEvent.bindings, p)
      ));
    }
    ie.prototype = this;
    const B = new ie(this);
    return we(this, B), B.level = this.level, B;
  }
  return c;
}
function qt(s) {
  const e = s.customLevels || {}, t = Object.assign({}, b.levels.values, e), n = Object.assign({}, b.levels.labels, Ot(e));
  return {
    values: t,
    labels: n
  };
}
function Ot(s) {
  const e = {};
  return Object.keys(s).forEach(function(t) {
    e[s[t]] = t;
  }), e;
}
b.levels = {
  values: {
    fatal: 60,
    error: 50,
    warn: 40,
    info: 30,
    debug: 20,
    trace: 10
  },
  labels: {
    10: "trace",
    20: "debug",
    30: "info",
    40: "warn",
    50: "error",
    60: "fatal"
  }
};
b.stdSerializers = Dt;
b.stdTimeFunctions = Object.assign({}, { nullTime: Te, epochTime: Se, unixTime: Mt, isoTime: Ft });
function Pt(s) {
  const e = [];
  s.bindings && e.push(s.bindings);
  let t = s[J];
  for (; t.parent; )
    t = t.parent, t.logger.bindings && e.push(t.logger.bindings);
  return e.reverse();
}
function v(s, e, t, n) {
  if (Object.defineProperty(s, n, {
    value: M(s.level, t) > M(n, t) ? q : t[te][n],
    writable: !0,
    enumerable: !0,
    configurable: !0
  }), !e.transmit && s[n] === q)
    return;
  s[n] = Nt(s, e, t, n);
  const r = Pt(s);
  r.length !== 0 && (s[n] = It(r, s[n]));
}
function It(s, e) {
  return function() {
    return e.apply(this, [...s, ...arguments]);
  };
}
function Nt(s, e, t, n) {
  return /* @__PURE__ */ function(r) {
    return function() {
      const a = e.timestamp(), o = new Array(arguments.length), l = Object.getPrototypeOf && Object.getPrototypeOf(this) === L ? L : this;
      for (var c = 0; c < o.length; c++)
        o[c] = arguments[c];
      if (e.serialize && !e.transmit && se(o, this._serialize, this.serializers, this._stdErrSerialize), e.asObject || e.formatters ? r.call(l, Ut(this, n, o, a, e.formatters)) : r.apply(l, o), e.transmit) {
        const u = e.transmit.level || s._level, h = t.levels.values[u], m = t.levels.values[n];
        if (m < h)
          return;
        At(this, {
          ts: a,
          methodLevel: n,
          methodValue: m,
          transmitLevel: u,
          transmitValue: t.levels.values[e.transmit.level || s._level],
          send: e.transmit.send,
          val: M(s._level, t)
        }, o);
      }
    };
  }(s[te][n]);
}
function Ut(s, e, t, n, r = {}) {
  const {
    level: i,
    log: a = (m) => m
  } = r, o = t.slice();
  let l = o[0];
  const c = {};
  if (n && (c.time = n), i) {
    const m = i(e, s.levels.values[e]);
    Object.assign(c, m);
  } else
    c.level = s.levels.values[e];
  let u = (s._childLevel | 0) + 1;
  if (u < 1 && (u = 1), l !== null && typeof l == "object") {
    for (; u-- && typeof o[0] == "object"; )
      Object.assign(c, o.shift());
    l = o.length ? ge(o.shift(), o) : void 0;
  } else
    typeof l == "string" && (l = ge(o.shift(), o));
  return l !== void 0 && (c.msg = l), a(c);
}
function se(s, e, t, n) {
  for (const r in s)
    if (n && s[r] instanceof Error)
      s[r] = b.stdSerializers.err(s[r]);
    else if (typeof s[r] == "object" && !Array.isArray(s[r]) && e)
      for (const i in s[r])
        e.indexOf(i) > -1 && i in t && (s[r][i] = t[i](s[r][i]));
}
function At(s, e, t) {
  const n = e.send, r = e.ts, i = e.methodLevel, a = e.methodValue, o = e.val, l = s._logEvent.bindings;
  se(
    t,
    s._serialize || Object.keys(s.serializers),
    s.serializers,
    s._stdErrSerialize === void 0 ? !0 : s._stdErrSerialize
  ), s._logEvent.ts = r, s._logEvent.messages = t.filter(function(c) {
    return l.indexOf(c) === -1;
  }), s._logEvent.level.label = i, s._logEvent.level.value = a, n(i, s._logEvent, o), s._logEvent = X(l);
}
function X(s) {
  return {
    ts: 0,
    messages: [],
    bindings: s || [],
    level: { label: "", value: 0 }
  };
}
function be(s) {
  const e = {
    type: s.constructor.name,
    msg: s.message,
    stack: s.stack
  };
  for (const t in s)
    e[t] === void 0 && (e[t] = s[t]);
  return e;
}
function jt(s) {
  return typeof s.timestamp == "function" ? s.timestamp : s.timestamp === !1 ? Te : Se;
}
function U() {
  return {};
}
function $(s) {
  return s;
}
function q() {
}
function Te() {
  return !1;
}
function Se() {
  return Date.now();
}
function Mt() {
  return Math.round(Date.now() / 1e3);
}
function Ft() {
  return new Date(Date.now()).toISOString();
}
function Bt() {
  function s(e) {
    return typeof e < "u" && e;
  }
  try {
    return typeof globalThis < "u" || Object.defineProperty(Object.prototype, "globalThis", {
      get: function() {
        return delete Object.prototype.globalThis, this.globalThis = this;
      },
      configurable: !0
    }), globalThis;
  } catch {
    return s(self) || s(window) || s(this) || {};
  }
}
F.exports.default = b;
var Kt = F.exports.pino = b, A = F.exports;
let Y = self.__POLYFEA_SW_LOGS_LEVEL === void 0 ? self.__POLYFEA_LOGS_LEVEL : self.__POLYFEA_SW_LOGS_LEVEL;
Y === void 0 && (Y = A.levels.values.info);
const R = Kt({
  level: A.levels.labels[Y],
  timestamp: A.stdTimeFunctions.isoTime,
  browser: {
    asObject: !0,
    write: (s) => {
      var l;
      const e = {
        trace: "#95a5a6",
        debug: "#7f8c8d",
        log: "#2ecc71",
        info: "#3498db",
        warn: "#f39c12",
        error: "#c0392b",
        fatal: "#c0392b"
      }, t = A.levels.labels[s.level], n = [
        `background: ${e[t] || "#000"}`,
        "border-radius: 0.5em",
        "color: white",
        "font-weight: bold",
        "padding: 2px 0.5em"
      ];
      let r = "polyfea", i = new Error();
      if (!i.stack)
        try {
          throw i;
        } catch {
        }
      let a = (l = i.stack) == null ? void 0 : l.toString().split(/\r\n|\n/);
      s.component && (r += "/" + s.component), s = Object.assign(s, { module: "polyfea", level: t, src: (a == null ? void 0 : a[1]) || void 0 });
      const o = ["%c" + r, n.join(";")];
      console[t](...o, s);
    }
  }
}).child({ component: "sw" });
class Wt {
  /**
   * Creates an instance of PolyfeaServiceWorker.
   * @param scope - The service worker global scope.
   */
  constructor(e = self) {
    this.scope = e, this.router = new $e(), Ge({
      prefix: "polyfea",
      suffix: "v1",
      precache: "install-time",
      runtime: "run-time"
    }), this.precacheController = new ze();
    const t = new URL(globalThis.location.href).searchParams.get("reconcile-interval");
    this.reconcilationInterval = (parseInt(t || "") || 60 * 30) * 1e3;
  }
  /**
   * Starts the service worker by adding event listeners and setting up route reconciliation.
   */
  async start() {
    this.scope.addEventListener("install", (e) => this.install(e)), this.scope.addEventListener("activate", (e) => this.activate(e)), this.scope.addEventListener("fetch", (e) => this.precache(e)), this.scope.addEventListener("fetch", (e) => this.runtime(e)), this.scope.addEventListener("fetch", (e) => this.fallback(e)), setInterval(() => {
      this.reconcileRoutes();
    }, this.reconcilationInterval);
  }
  /**
   * @private
   * Reconciles the routes by fetching the caching configuration and updating the precache and router.
   */
  async reconcileRoutes(e = !1) {
    var i, a, o;
    const t = await this.getLastReconciliationTime();
    let n = 0;
    if (t && (n = Date.now() + 1e3 - parseInt(t)), !e && n && n < this.reconcilationInterval) {
      R.debug("Skipping reconciliation - data are fresh ");
      return;
    }
    const r = decodeURIComponent(
      new URL(globalThis.location.href).searchParams.get("caching-config") || "./polyfea-caching.json"
    );
    try {
      const l = await fetch(r);
      if (l.status < 300) {
        const c = await l.json();
        this.precacheController.addToCacheList((c.precache || []).filter((u) => {
          const h = typeof u == "string" ? u : u.url;
          return !this.precacheController.getCacheKeyForURL(h);
        })), this.router.routes.clear(), (i = c.routes) == null || i.map(ee.from).forEach((u) => this.router.registerRoute(u)), R.info(`Service worker reconciled: precached ${((a = c.precache) == null ? void 0 : a.length) || 0} files and added ${((o = c.routes) == null ? void 0 : o.length) || 0} routes`);
      }
      await this.setLastReconciliationTime(Date.now().toString());
    } catch (l) {
      R.warn({ err: l }, "Failed to reconcile routes");
    }
  }
  async getLastReconciliationTime() {
    return new Promise((e, t) => {
      const n = indexedDB.open("polyfeaDB", 1);
      n.onerror = () => t(n.error), n.onupgradeneeded = () => {
        n.result.createObjectStore("reconciliationTime", { keyPath: "id" });
      }, n.onsuccess = () => {
        const o = n.result.transaction("reconciliationTime", "readonly").objectStore("reconciliationTime").get("lastReconciliationTime");
        o.onerror = () => e(null), o.onsuccess = () => {
          var l;
          return e((l = o.result) == null ? void 0 : l.value);
        };
      };
    });
  }
  async setLastReconciliationTime(e) {
    return new Promise((t, n) => {
      const r = indexedDB.open("polyfeaDB", 1);
      r.onerror = () => n(r.error), r.onupgradeneeded = () => {
        r.result.createObjectStore("reconciliationTime");
      }, r.onsuccess = () => {
        const l = r.result.transaction("reconciliationTime", "readwrite").objectStore("reconciliationTime").put({ id: "lastReconciliationTime", value: e });
        l.onerror = () => n(l.error), l.onsuccess = () => t();
      };
    });
  }
  /**
   * @private
   * Installs the service worker by reconciling routes and installing the precache.
   * @param event - The install event.
   */
  install(e) {
    e.waitUntil((async () => {
      R.debug("Installing"), await this.reconcileRoutes(!0), await this.precacheController.install(e);
    })());
  }
  /**
   * @private
   * Activates the service worker by activating the precache.
   * @param event - The activate event.
   */
  activate(e) {
    e.waitUntil((async () => {
      R.debug("Activating"), this.precacheController.activate(e);
    })());
  }
  /**
   * @private
   * Handles the fetch event by responding from the precache if the URL is in the precache.
   * @param event - The fetch event.
   */
  precache(e) {
    const t = R.child({ request: e.request });
    t.trace({ request: e.request }, `trying to fetch from the precache ${e.request.url}`);
    const { request: n } = e, r = this.precacheController.getCacheKeyForURL(n.url);
    if (r) {
      t.debug(`Responded from precache: ${e.request.url}`), e.respondWith(caches.match(r));
      return;
    }
  }
  /**
   * @private
   * Handles the fetch event by responding from the router if a matching route is found.
   * @param event - The fetch event.
   */
  runtime(e) {
    const t = R.child({ request: e.request });
    t.trace({ request: e.request }, `trying to fetch from the router ${e.request.url}`);
    const { request: n } = e, r = this.router.handleRequest({
      event: e,
      request: n
    });
    r && (t.debug(`Responded from router: ${e.request.url}`), e.respondWith(r));
  }
  /**
   * @private
   * Handles the fetch event when a route is not found.
   * @param event - The fetch event.
   */
  fallback(e) {
    R.debug(`Route not found, ignoring: ${e.request.url}`);
  }
}
new Wt().start();
//# sourceMappingURL=sw.mjs.map
