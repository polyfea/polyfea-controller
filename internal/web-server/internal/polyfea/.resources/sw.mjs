//#region \0rolldown/runtime.js
var e = (e, t) => () => (t || (e((t = { exports: {} }).exports, t), e = null), t.exports);
//#endregion
//#region node_modules/workbox-core/_version.js
try {
	self["workbox:core:7.3.0"] && _();
} catch {}
var t = (e, ...t) => {
	let n = e;
	return t.length > 0 && (n += ` :: ${JSON.stringify(t)}`), n;
}, n = class extends Error {
	constructor(e, n) {
		let r = t(e, n);
		super(r), this.name = e, this.details = n;
	}
}, r = {
	googleAnalytics: "googleAnalytics",
	precache: "precache-v2",
	prefix: "workbox",
	runtime: "runtime",
	suffix: typeof registration < "u" ? registration.scope : ""
}, i = (e) => [
	r.prefix,
	e,
	r.suffix
].filter((e) => e && e.length > 0).join("-"), a = (e) => {
	for (let t of Object.keys(r)) e(t);
}, o = {
	updateDetails: (e) => {
		a((t) => {
			typeof e[t] == "string" && (r[t] = e[t]);
		});
	},
	getGoogleAnalyticsName: (e) => e || i(r.googleAnalytics),
	getPrecacheName: (e) => e || i(r.precache),
	getPrefix: () => r.prefix,
	getRuntimeName: (e) => e || i(r.runtime),
	getSuffix: () => r.suffix
};
//#endregion
//#region node_modules/workbox-core/_private/waitUntil.js
function s(e, t) {
	let n = t();
	return e.waitUntil(n), n;
//#endregion
//#region node_modules/workbox-core/_private/waitUntil.js
function s(e, t) {
	let n = t();
	return e.waitUntil(n), n;
}
//#endregion
//#region node_modules/workbox-precaching/_version.js
//#endregion
//#region node_modules/workbox-precaching/_version.js
try {
	self["workbox:precaching:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-precaching/utils/createCacheKey.js
var c = "__WB_REVISION__";
function l(e) {
	if (!e) throw new n("add-to-cache-list-unexpected-type", { entry: e });
	if (typeof e == "string") {
		let t = new URL(e, location.href);
		return {
			cacheKey: t.href,
			url: t.href
		};
	}
	let { revision: t, url: r } = e;
	if (!r) throw new n("add-to-cache-list-unexpected-type", { entry: e });
	if (!t) {
		let e = new URL(r, location.href);
		return {
			cacheKey: e.href,
			url: e.href
		};
	}
	let i = new URL(r, location.href), a = new URL(r, location.href);
	return i.searchParams.set(c, t), {
		cacheKey: i.href,
		url: a.href
	};
	self["workbox:precaching:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-precaching/utils/createCacheKey.js
var c = "__WB_REVISION__";
function l(e) {
	if (!e) throw new n("add-to-cache-list-unexpected-type", { entry: e });
	if (typeof e == "string") {
		let t = new URL(e, location.href);
		return {
			cacheKey: t.href,
			url: t.href
		};
	}
	let { revision: t, url: r } = e;
	if (!r) throw new n("add-to-cache-list-unexpected-type", { entry: e });
	if (!t) {
		let e = new URL(r, location.href);
		return {
			cacheKey: e.href,
			url: e.href
		};
	}
	let i = new URL(r, location.href), a = new URL(r, location.href);
	return i.searchParams.set(c, t), {
		cacheKey: i.href,
		url: a.href
	};
}
//#endregion
//#region node_modules/workbox-precaching/utils/PrecacheInstallReportPlugin.js
var u = class {
	constructor() {
		this.updatedURLs = [], this.notUpdatedURLs = [], this.handlerWillStart = async ({ request: e, state: t }) => {
			t && (t.originalRequest = e);
		}, this.cachedResponseWillBeUsed = async ({ event: e, state: t, cachedResponse: n }) => {
			if (e.type === "install" && t && t.originalRequest && t.originalRequest instanceof Request) {
				let e = t.originalRequest.url;
				n ? this.notUpdatedURLs.push(e) : this.updatedURLs.push(e);
			}
			return n;
		};
	}
}, d = class {
	constructor({ precacheController: e }) {
		this.cacheKeyWillBeUsed = async ({ request: e, params: t }) => {
			let n = t?.cacheKey || this._precacheController.getCacheKeyForURL(e.url);
			return n ? new Request(n, { headers: e.headers }) : e;
		}, this._precacheController = e;
	}
}, f;
function p() {
	if (f === void 0) {
		let e = new Response("");
		if ("body" in e) try {
			new Response(e.body), f = !0;
		} catch {
			f = !1;
		}
		f = !1;
	}
	return f;
//#endregion
//#region node_modules/workbox-precaching/utils/PrecacheInstallReportPlugin.js
var u = class {
	constructor() {
		this.updatedURLs = [], this.notUpdatedURLs = [], this.handlerWillStart = async ({ request: e, state: t }) => {
			t && (t.originalRequest = e);
		}, this.cachedResponseWillBeUsed = async ({ event: e, state: t, cachedResponse: n }) => {
			if (e.type === "install" && t && t.originalRequest && t.originalRequest instanceof Request) {
				let e = t.originalRequest.url;
				n ? this.notUpdatedURLs.push(e) : this.updatedURLs.push(e);
			}
			return n;
		};
	}
}, d = class {
	constructor({ precacheController: e }) {
		this.cacheKeyWillBeUsed = async ({ request: e, params: t }) => {
			let n = t?.cacheKey || this._precacheController.getCacheKeyForURL(e.url);
			return n ? new Request(n, { headers: e.headers }) : e;
		}, this._precacheController = e;
	}
}, f;
function p() {
	if (f === void 0) {
		let e = new Response("");
		if ("body" in e) try {
			new Response(e.body), f = !0;
		} catch {
			f = !1;
		}
		f = !1;
	}
	return f;
}
//#endregion
//#region node_modules/workbox-core/copyResponse.js
async function m(e, t) {
	let r = null;
	if (e.url && (r = new URL(e.url).origin), r !== self.location.origin) throw new n("cross-origin-copy-response", { origin: r });
	let i = e.clone(), a = {
		headers: new Headers(i.headers),
		status: i.status,
		statusText: i.statusText
	}, o = t ? t(a) : a, s = p() ? i.body : await i.blob();
	return new Response(s, o);
//#endregion
//#region node_modules/workbox-core/copyResponse.js
async function m(e, t) {
	let r = null;
	if (e.url && (r = new URL(e.url).origin), r !== self.location.origin) throw new n("cross-origin-copy-response", { origin: r });
	let i = e.clone(), a = {
		headers: new Headers(i.headers),
		status: i.status,
		statusText: i.statusText
	}, o = t ? t(a) : a, s = p() ? i.body : await i.blob();
	return new Response(s, o);
}
//#endregion
//#region node_modules/workbox-core/_private/getFriendlyURL.js
var h = (e) => new URL(String(e), location.href).href.replace(RegExp(`^${location.origin}`), "");
//#endregion
//#region node_modules/workbox-core/_private/cacheMatchIgnoreParams.js
function g(e, t) {
	let n = new URL(e);
	for (let e of t) n.searchParams.delete(e);
	return n.href;
//#endregion
//#region node_modules/workbox-core/_private/getFriendlyURL.js
var h = (e) => new URL(String(e), location.href).href.replace(RegExp(`^${location.origin}`), "");
//#endregion
//#region node_modules/workbox-core/_private/cacheMatchIgnoreParams.js
function g(e, t) {
	let n = new URL(e);
	for (let e of t) n.searchParams.delete(e);
	return n.href;
}
async function v(e, t, n, r) {
	let i = g(t.url, n);
	if (t.url === i) return e.match(t, r);
	let a = Object.assign(Object.assign({}, r), { ignoreSearch: !0 }), o = await e.keys(t, a);
	for (let t of o) if (i === g(t.url, n)) return e.match(t, r);
}
//#endregion
//#region node_modules/workbox-core/_private/Deferred.js
var y = class {
	constructor() {
		this.promise = new Promise((e, t) => {
			this.resolve = e, this.reject = t;
		});
	}
}, b = /* @__PURE__ */ new Set();
//#endregion
//#region node_modules/workbox-core/_private/executeQuotaErrorCallbacks.js
async function x() {
	for (let e of b) await e();
async function v(e, t, n, r) {
	let i = g(t.url, n);
	if (t.url === i) return e.match(t, r);
	let a = Object.assign(Object.assign({}, r), { ignoreSearch: !0 }), o = await e.keys(t, a);
	for (let t of o) if (i === g(t.url, n)) return e.match(t, r);
}
//#endregion
//#region node_modules/workbox-core/_private/Deferred.js
var y = class {
	constructor() {
		this.promise = new Promise((e, t) => {
			this.resolve = e, this.reject = t;
		});
	}
}, b = /* @__PURE__ */ new Set();
//#endregion
//#region node_modules/workbox-core/_private/executeQuotaErrorCallbacks.js
async function x() {
	for (let e of b) await e();
}
//#endregion
//#region node_modules/workbox-core/_private/timeout.js
function S(e) {
	return new Promise((t) => setTimeout(t, e));
//#endregion
//#region node_modules/workbox-core/_private/timeout.js
function S(e) {
	return new Promise((t) => setTimeout(t, e));
}
//#endregion
//#region node_modules/workbox-strategies/_version.js
//#endregion
//#region node_modules/workbox-strategies/_version.js
try {
	self["workbox:strategies:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-strategies/StrategyHandler.js
function C(e) {
	return typeof e == "string" ? new Request(e) : e;
}
var w = class {
	constructor(e, t) {
		this._cacheKeys = {}, Object.assign(this, t), this.event = t.event, this._strategy = e, this._handlerDeferred = new y(), this._extendLifetimePromises = [], this._plugins = [...e.plugins], this._pluginStateMap = /* @__PURE__ */ new Map();
		for (let e of this._plugins) this._pluginStateMap.set(e, {});
		this.event.waitUntil(this._handlerDeferred.promise);
	}
	async fetch(e) {
		let { event: t } = this, r = C(e);
		if (r.mode === "navigate" && t instanceof FetchEvent && t.preloadResponse) {
			let e = await t.preloadResponse;
			if (e) return e;
		}
		let i = this.hasCallback("fetchDidFail") ? r.clone() : null;
		try {
			for (let e of this.iterateCallbacks("requestWillFetch")) r = await e({
				request: r.clone(),
				event: t
			});
		} catch (e) {
			if (e instanceof Error) throw new n("plugin-error-request-will-fetch", { thrownErrorMessage: e.message });
		}
		let a = r.clone();
		try {
			let e;
			e = await fetch(r, r.mode === "navigate" ? void 0 : this._strategy.fetchOptions);
			for (let n of this.iterateCallbacks("fetchDidSucceed")) e = await n({
				event: t,
				request: a,
				response: e
			});
			return e;
		} catch (e) {
			throw i && await this.runCallbacks("fetchDidFail", {
				error: e,
				event: t,
				originalRequest: i.clone(),
				request: a.clone()
			}), e;
		}
	}
	async fetchAndCachePut(e) {
		let t = await this.fetch(e), n = t.clone();
		return this.waitUntil(this.cachePut(e, n)), t;
	}
	async cacheMatch(e) {
		let t = C(e), n, { cacheName: r, matchOptions: i } = this._strategy, a = await this.getCacheKey(t, "read"), o = Object.assign(Object.assign({}, i), { cacheName: r });
		n = await caches.match(a, o);
		for (let e of this.iterateCallbacks("cachedResponseWillBeUsed")) n = await e({
			cacheName: r,
			matchOptions: i,
			cachedResponse: n,
			request: a,
			event: this.event
		}) || void 0;
		return n;
	}
	async cachePut(e, t) {
		let r = C(e);
		await S(0);
		let i = await this.getCacheKey(r, "write");
		if (!t) throw new n("cache-put-with-no-response", { url: h(i.url) });
		let a = await this._ensureResponseSafeToCache(t);
		if (!a) return !1;
		let { cacheName: o, matchOptions: s } = this._strategy, c = await self.caches.open(o), l = this.hasCallback("cacheDidUpdate"), u = l ? await v(c, i.clone(), ["__WB_REVISION__"], s) : null;
		try {
			await c.put(i, l ? a.clone() : a);
		} catch (e) {
			if (e instanceof Error) throw e.name === "QuotaExceededError" && await x(), e;
		}
		for (let e of this.iterateCallbacks("cacheDidUpdate")) await e({
			cacheName: o,
			oldResponse: u,
			newResponse: a.clone(),
			request: i,
			event: this.event
		});
		return !0;
	}
	async getCacheKey(e, t) {
		let n = `${e.url} | ${t}`;
		if (!this._cacheKeys[n]) {
			let r = e;
			for (let e of this.iterateCallbacks("cacheKeyWillBeUsed")) r = C(await e({
				mode: t,
				request: r,
				event: this.event,
				params: this.params
			}));
			this._cacheKeys[n] = r;
		}
		return this._cacheKeys[n];
	}
	hasCallback(e) {
		for (let t of this._strategy.plugins) if (e in t) return !0;
		return !1;
	}
	async runCallbacks(e, t) {
		for (let n of this.iterateCallbacks(e)) await n(t);
	}
	*iterateCallbacks(e) {
		for (let t of this._strategy.plugins) if (typeof t[e] == "function") {
			let n = this._pluginStateMap.get(t);
			yield (r) => {
				let i = Object.assign(Object.assign({}, r), { state: n });
				return t[e](i);
			};
		}
	}
	waitUntil(e) {
		return this._extendLifetimePromises.push(e), e;
	}
	async doneWaiting() {
		for (; this._extendLifetimePromises.length;) {
			let e = this._extendLifetimePromises.splice(0), t = (await Promise.allSettled(e)).find((e) => e.status === "rejected");
			if (t) throw t.reason;
		}
	}
	destroy() {
		this._handlerDeferred.resolve(null);
	}
	async _ensureResponseSafeToCache(e) {
		let t = e, n = !1;
		for (let e of this.iterateCallbacks("cacheWillUpdate")) if (t = await e({
			request: this.request,
			response: t,
			event: this.event
		}) || void 0, n = !0, !t) break;
		return n || t && t.status !== 200 && (t = void 0), t;
	}
}, T = class {
	constructor(e = {}) {
		this.cacheName = o.getRuntimeName(e.cacheName), this.plugins = e.plugins || [], this.fetchOptions = e.fetchOptions, this.matchOptions = e.matchOptions;
	}
	handle(e) {
		let [t] = this.handleAll(e);
		return t;
	}
	handleAll(e) {
		e instanceof FetchEvent && (e = {
			event: e,
			request: e.request
		});
		let t = e.event, n = typeof e.request == "string" ? new Request(e.request) : e.request, r = "params" in e ? e.params : void 0, i = new w(this, {
			event: t,
			request: n,
			params: r
		}), a = this._getResponse(i, n, t);
		return [a, this._awaitComplete(a, i, n, t)];
	}
	async _getResponse(e, t, r) {
		await e.runCallbacks("handlerWillStart", {
			event: r,
			request: t
		});
		let i;
		try {
			if (i = await this._handle(t, e), !i || i.type === "error") throw new n("no-response", { url: t.url });
		} catch (n) {
			if (n instanceof Error) {
				for (let a of e.iterateCallbacks("handlerDidError")) if (i = await a({
					error: n,
					event: r,
					request: t
				}), i) break;
			}
			if (!i) throw n;
		}
		for (let n of e.iterateCallbacks("handlerWillRespond")) i = await n({
			event: r,
			request: t,
			response: i
		});
		return i;
	}
	async _awaitComplete(e, t, n, r) {
		let i, a;
		try {
			i = await e;
		} catch {}
		try {
			await t.runCallbacks("handlerDidRespond", {
				event: r,
				request: n,
				response: i
			}), await t.doneWaiting();
		} catch (e) {
			e instanceof Error && (a = e);
		}
		if (await t.runCallbacks("handlerDidComplete", {
			event: r,
			request: n,
			response: i,
			error: a
		}), t.destroy(), a) throw a;
	}
}, E = class e extends T {
	constructor(t = {}) {
		t.cacheName = o.getPrecacheName(t.cacheName), super(t), this._fallbackToNetwork = t.fallbackToNetwork !== !1, this.plugins.push(e.copyRedirectedCacheableResponsesPlugin);
	}
	async _handle(e, t) {
		return await t.cacheMatch(e) || (t.event && t.event.type === "install" ? await this._handleInstall(e, t) : await this._handleFetch(e, t));
	}
	async _handleFetch(e, t) {
		let r, i = t.params || {};
		if (this._fallbackToNetwork) {
			let n = i.integrity, a = e.integrity, o = !a || a === n;
			r = await t.fetch(new Request(e, { integrity: e.mode === "no-cors" ? void 0 : a || n })), n && o && e.mode !== "no-cors" && (this._useDefaultCacheabilityPluginIfNeeded(), await t.cachePut(e, r.clone()));
		} else throw new n("missing-precache-entry", {
			cacheName: this.cacheName,
			url: e.url
		});
		return r;
	}
	async _handleInstall(e, t) {
		this._useDefaultCacheabilityPluginIfNeeded();
		let r = await t.fetch(e);
		if (!await t.cachePut(e, r.clone())) throw new n("bad-precaching-response", {
			url: e.url,
			status: r.status
		});
		return r;
	}
	_useDefaultCacheabilityPluginIfNeeded() {
		let t = null, n = 0;
		for (let [r, i] of this.plugins.entries()) i !== e.copyRedirectedCacheableResponsesPlugin && (i === e.defaultPrecacheCacheabilityPlugin && (t = r), i.cacheWillUpdate && n++);
		n === 0 ? this.plugins.push(e.defaultPrecacheCacheabilityPlugin) : n > 1 && t !== null && this.plugins.splice(t, 1);
	}
};
E.defaultPrecacheCacheabilityPlugin = { async cacheWillUpdate({ response: e }) {
	return !e || e.status >= 400 ? null : e;
} }, E.copyRedirectedCacheableResponsesPlugin = { async cacheWillUpdate({ response: e }) {
	return e.redirected ? await m(e) : e;
} };
//#endregion
//#region node_modules/workbox-precaching/PrecacheController.js
var D = class {
	constructor({ cacheName: e, plugins: t = [], fallbackToNetwork: n = !0 } = {}) {
		this._urlsToCacheKeys = /* @__PURE__ */ new Map(), this._urlsToCacheModes = /* @__PURE__ */ new Map(), this._cacheKeysToIntegrities = /* @__PURE__ */ new Map(), this._strategy = new E({
			cacheName: o.getPrecacheName(e),
			plugins: [...t, new d({ precacheController: this })],
			fallbackToNetwork: n
		}), this.install = this.install.bind(this), this.activate = this.activate.bind(this);
	}
	get strategy() {
		return this._strategy;
	}
	precache(e) {
		this.addToCacheList(e), this._installAndActiveListenersAdded ||= (self.addEventListener("install", this.install), self.addEventListener("activate", this.activate), !0);
	}
	addToCacheList(e) {
		let t = [];
		for (let r of e) {
			typeof r == "string" ? t.push(r) : r && r.revision === void 0 && t.push(r.url);
			let { cacheKey: e, url: i } = l(r), a = typeof r != "string" && r.revision ? "reload" : "default";
			if (this._urlsToCacheKeys.has(i) && this._urlsToCacheKeys.get(i) !== e) throw new n("add-to-cache-list-conflicting-entries", {
				firstEntry: this._urlsToCacheKeys.get(i),
				secondEntry: e
			});
			if (typeof r != "string" && r.integrity) {
				if (this._cacheKeysToIntegrities.has(e) && this._cacheKeysToIntegrities.get(e) !== r.integrity) throw new n("add-to-cache-list-conflicting-integrities", { url: i });
				this._cacheKeysToIntegrities.set(e, r.integrity);
			}
			if (this._urlsToCacheKeys.set(i, e), this._urlsToCacheModes.set(i, a), t.length > 0) {
				let e = `Workbox is precaching URLs without revision info: ${t.join(", ")}\nThis is generally NOT safe. Learn more at https://bit.ly/wb-precache`;
				console.warn(e);
			}
		}
	}
	install(e) {
		return s(e, async () => {
			let t = new u();
			this.strategy.plugins.push(t);
			for (let [t, n] of this._urlsToCacheKeys) {
				let r = this._cacheKeysToIntegrities.get(n), i = this._urlsToCacheModes.get(t), a = new Request(t, {
					integrity: r,
					cache: i,
					credentials: "same-origin"
				});
				await Promise.all(this.strategy.handleAll({
					params: { cacheKey: n },
					request: a,
					event: e
				}));
			}
			let { updatedURLs: n, notUpdatedURLs: r } = t;
			return {
				updatedURLs: n,
				notUpdatedURLs: r
			};
		});
	}
	activate(e) {
		return s(e, async () => {
			let e = await self.caches.open(this.strategy.cacheName), t = await e.keys(), n = new Set(this._urlsToCacheKeys.values()), r = [];
			for (let i of t) n.has(i.url) || (await e.delete(i), r.push(i.url));
			return { deletedURLs: r };
		});
	}
	getURLsToCacheKeys() {
		return this._urlsToCacheKeys;
	}
	getCachedURLs() {
		return [...this._urlsToCacheKeys.keys()];
	}
	getCacheKeyForURL(e) {
		let t = new URL(e, location.href);
		return this._urlsToCacheKeys.get(t.href);
	}
	getIntegrityForCacheKey(e) {
		return this._cacheKeysToIntegrities.get(e);
	}
	async matchPrecache(e) {
		let t = e instanceof Request ? e.url : e, n = this.getCacheKeyForURL(t);
		if (n) return (await self.caches.open(this.strategy.cacheName)).match(n);
	}
	createHandlerBoundToURL(e) {
		let t = this.getCacheKeyForURL(e);
		if (!t) throw new n("non-precached-url", { url: e });
		return (n) => (n.request = new Request(e), n.params = Object.assign({ cacheKey: t }, n.params), this.strategy.handle(n));
	}
};
//#endregion
//#region node_modules/workbox-routing/_version.js
	self["workbox:strategies:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-strategies/StrategyHandler.js
function C(e) {
	return typeof e == "string" ? new Request(e) : e;
}
var w = class {
	constructor(e, t) {
		this._cacheKeys = {}, Object.assign(this, t), this.event = t.event, this._strategy = e, this._handlerDeferred = new y(), this._extendLifetimePromises = [], this._plugins = [...e.plugins], this._pluginStateMap = /* @__PURE__ */ new Map();
		for (let e of this._plugins) this._pluginStateMap.set(e, {});
		this.event.waitUntil(this._handlerDeferred.promise);
	}
	async fetch(e) {
		let { event: t } = this, r = C(e);
		if (r.mode === "navigate" && t instanceof FetchEvent && t.preloadResponse) {
			let e = await t.preloadResponse;
			if (e) return e;
		}
		let i = this.hasCallback("fetchDidFail") ? r.clone() : null;
		try {
			for (let e of this.iterateCallbacks("requestWillFetch")) r = await e({
				request: r.clone(),
				event: t
			});
		} catch (e) {
			if (e instanceof Error) throw new n("plugin-error-request-will-fetch", { thrownErrorMessage: e.message });
		}
		let a = r.clone();
		try {
			let e;
			e = await fetch(r, r.mode === "navigate" ? void 0 : this._strategy.fetchOptions);
			for (let n of this.iterateCallbacks("fetchDidSucceed")) e = await n({
				event: t,
				request: a,
				response: e
			});
			return e;
		} catch (e) {
			throw i && await this.runCallbacks("fetchDidFail", {
				error: e,
				event: t,
				originalRequest: i.clone(),
				request: a.clone()
			}), e;
		}
	}
	async fetchAndCachePut(e) {
		let t = await this.fetch(e), n = t.clone();
		return this.waitUntil(this.cachePut(e, n)), t;
	}
	async cacheMatch(e) {
		let t = C(e), n, { cacheName: r, matchOptions: i } = this._strategy, a = await this.getCacheKey(t, "read"), o = Object.assign(Object.assign({}, i), { cacheName: r });
		n = await caches.match(a, o);
		for (let e of this.iterateCallbacks("cachedResponseWillBeUsed")) n = await e({
			cacheName: r,
			matchOptions: i,
			cachedResponse: n,
			request: a,
			event: this.event
		}) || void 0;
		return n;
	}
	async cachePut(e, t) {
		let r = C(e);
		await S(0);
		let i = await this.getCacheKey(r, "write");
		if (!t) throw new n("cache-put-with-no-response", { url: h(i.url) });
		let a = await this._ensureResponseSafeToCache(t);
		if (!a) return !1;
		let { cacheName: o, matchOptions: s } = this._strategy, c = await self.caches.open(o), l = this.hasCallback("cacheDidUpdate"), u = l ? await v(c, i.clone(), ["__WB_REVISION__"], s) : null;
		try {
			await c.put(i, l ? a.clone() : a);
		} catch (e) {
			if (e instanceof Error) throw e.name === "QuotaExceededError" && await x(), e;
		}
		for (let e of this.iterateCallbacks("cacheDidUpdate")) await e({
			cacheName: o,
			oldResponse: u,
			newResponse: a.clone(),
			request: i,
			event: this.event
		});
		return !0;
	}
	async getCacheKey(e, t) {
		let n = `${e.url} | ${t}`;
		if (!this._cacheKeys[n]) {
			let r = e;
			for (let e of this.iterateCallbacks("cacheKeyWillBeUsed")) r = C(await e({
				mode: t,
				request: r,
				event: this.event,
				params: this.params
			}));
			this._cacheKeys[n] = r;
		}
		return this._cacheKeys[n];
	}
	hasCallback(e) {
		for (let t of this._strategy.plugins) if (e in t) return !0;
		return !1;
	}
	async runCallbacks(e, t) {
		for (let n of this.iterateCallbacks(e)) await n(t);
	}
	*iterateCallbacks(e) {
		for (let t of this._strategy.plugins) if (typeof t[e] == "function") {
			let n = this._pluginStateMap.get(t);
			yield (r) => {
				let i = Object.assign(Object.assign({}, r), { state: n });
				return t[e](i);
			};
		}
	}
	waitUntil(e) {
		return this._extendLifetimePromises.push(e), e;
	}
	async doneWaiting() {
		for (; this._extendLifetimePromises.length;) {
			let e = this._extendLifetimePromises.splice(0), t = (await Promise.allSettled(e)).find((e) => e.status === "rejected");
			if (t) throw t.reason;
		}
	}
	destroy() {
		this._handlerDeferred.resolve(null);
	}
	async _ensureResponseSafeToCache(e) {
		let t = e, n = !1;
		for (let e of this.iterateCallbacks("cacheWillUpdate")) if (t = await e({
			request: this.request,
			response: t,
			event: this.event
		}) || void 0, n = !0, !t) break;
		return n || t && t.status !== 200 && (t = void 0), t;
	}
}, T = class {
	constructor(e = {}) {
		this.cacheName = o.getRuntimeName(e.cacheName), this.plugins = e.plugins || [], this.fetchOptions = e.fetchOptions, this.matchOptions = e.matchOptions;
	}
	handle(e) {
		let [t] = this.handleAll(e);
		return t;
	}
	handleAll(e) {
		e instanceof FetchEvent && (e = {
			event: e,
			request: e.request
		});
		let t = e.event, n = typeof e.request == "string" ? new Request(e.request) : e.request, r = "params" in e ? e.params : void 0, i = new w(this, {
			event: t,
			request: n,
			params: r
		}), a = this._getResponse(i, n, t);
		return [a, this._awaitComplete(a, i, n, t)];
	}
	async _getResponse(e, t, r) {
		await e.runCallbacks("handlerWillStart", {
			event: r,
			request: t
		});
		let i;
		try {
			if (i = await this._handle(t, e), !i || i.type === "error") throw new n("no-response", { url: t.url });
		} catch (n) {
			if (n instanceof Error) {
				for (let a of e.iterateCallbacks("handlerDidError")) if (i = await a({
					error: n,
					event: r,
					request: t
				}), i) break;
			}
			if (!i) throw n;
		}
		for (let n of e.iterateCallbacks("handlerWillRespond")) i = await n({
			event: r,
			request: t,
			response: i
		});
		return i;
	}
	async _awaitComplete(e, t, n, r) {
		let i, a;
		try {
			i = await e;
		} catch {}
		try {
			await t.runCallbacks("handlerDidRespond", {
				event: r,
				request: n,
				response: i
			}), await t.doneWaiting();
		} catch (e) {
			e instanceof Error && (a = e);
		}
		if (await t.runCallbacks("handlerDidComplete", {
			event: r,
			request: n,
			response: i,
			error: a
		}), t.destroy(), a) throw a;
	}
}, E = class e extends T {
	constructor(t = {}) {
		t.cacheName = o.getPrecacheName(t.cacheName), super(t), this._fallbackToNetwork = t.fallbackToNetwork !== !1, this.plugins.push(e.copyRedirectedCacheableResponsesPlugin);
	}
	async _handle(e, t) {
		return await t.cacheMatch(e) || (t.event && t.event.type === "install" ? await this._handleInstall(e, t) : await this._handleFetch(e, t));
	}
	async _handleFetch(e, t) {
		let r, i = t.params || {};
		if (this._fallbackToNetwork) {
			let n = i.integrity, a = e.integrity, o = !a || a === n;
			r = await t.fetch(new Request(e, { integrity: e.mode === "no-cors" ? void 0 : a || n })), n && o && e.mode !== "no-cors" && (this._useDefaultCacheabilityPluginIfNeeded(), await t.cachePut(e, r.clone()));
		} else throw new n("missing-precache-entry", {
			cacheName: this.cacheName,
			url: e.url
		});
		return r;
	}
	async _handleInstall(e, t) {
		this._useDefaultCacheabilityPluginIfNeeded();
		let r = await t.fetch(e);
		if (!await t.cachePut(e, r.clone())) throw new n("bad-precaching-response", {
			url: e.url,
			status: r.status
		});
		return r;
	}
	_useDefaultCacheabilityPluginIfNeeded() {
		let t = null, n = 0;
		for (let [r, i] of this.plugins.entries()) i !== e.copyRedirectedCacheableResponsesPlugin && (i === e.defaultPrecacheCacheabilityPlugin && (t = r), i.cacheWillUpdate && n++);
		n === 0 ? this.plugins.push(e.defaultPrecacheCacheabilityPlugin) : n > 1 && t !== null && this.plugins.splice(t, 1);
	}
};
E.defaultPrecacheCacheabilityPlugin = { async cacheWillUpdate({ response: e }) {
	return !e || e.status >= 400 ? null : e;
} }, E.copyRedirectedCacheableResponsesPlugin = { async cacheWillUpdate({ response: e }) {
	return e.redirected ? await m(e) : e;
} };
//#endregion
//#region node_modules/workbox-precaching/PrecacheController.js
var D = class {
	constructor({ cacheName: e, plugins: t = [], fallbackToNetwork: n = !0 } = {}) {
		this._urlsToCacheKeys = /* @__PURE__ */ new Map(), this._urlsToCacheModes = /* @__PURE__ */ new Map(), this._cacheKeysToIntegrities = /* @__PURE__ */ new Map(), this._strategy = new E({
			cacheName: o.getPrecacheName(e),
			plugins: [...t, new d({ precacheController: this })],
			fallbackToNetwork: n
		}), this.install = this.install.bind(this), this.activate = this.activate.bind(this);
	}
	get strategy() {
		return this._strategy;
	}
	precache(e) {
		this.addToCacheList(e), this._installAndActiveListenersAdded ||= (self.addEventListener("install", this.install), self.addEventListener("activate", this.activate), !0);
	}
	addToCacheList(e) {
		let t = [];
		for (let r of e) {
			typeof r == "string" ? t.push(r) : r && r.revision === void 0 && t.push(r.url);
			let { cacheKey: e, url: i } = l(r), a = typeof r != "string" && r.revision ? "reload" : "default";
			if (this._urlsToCacheKeys.has(i) && this._urlsToCacheKeys.get(i) !== e) throw new n("add-to-cache-list-conflicting-entries", {
				firstEntry: this._urlsToCacheKeys.get(i),
				secondEntry: e
			});
			if (typeof r != "string" && r.integrity) {
				if (this._cacheKeysToIntegrities.has(e) && this._cacheKeysToIntegrities.get(e) !== r.integrity) throw new n("add-to-cache-list-conflicting-integrities", { url: i });
				this._cacheKeysToIntegrities.set(e, r.integrity);
			}
			if (this._urlsToCacheKeys.set(i, e), this._urlsToCacheModes.set(i, a), t.length > 0) {
				let e = `Workbox is precaching URLs without revision info: ${t.join(", ")}\nThis is generally NOT safe. Learn more at https://bit.ly/wb-precache`;
				console.warn(e);
			}
		}
	}
	install(e) {
		return s(e, async () => {
			let t = new u();
			this.strategy.plugins.push(t);
			for (let [t, n] of this._urlsToCacheKeys) {
				let r = this._cacheKeysToIntegrities.get(n), i = this._urlsToCacheModes.get(t), a = new Request(t, {
					integrity: r,
					cache: i,
					credentials: "same-origin"
				});
				await Promise.all(this.strategy.handleAll({
					params: { cacheKey: n },
					request: a,
					event: e
				}));
			}
			let { updatedURLs: n, notUpdatedURLs: r } = t;
			return {
				updatedURLs: n,
				notUpdatedURLs: r
			};
		});
	}
	activate(e) {
		return s(e, async () => {
			let e = await self.caches.open(this.strategy.cacheName), t = await e.keys(), n = new Set(this._urlsToCacheKeys.values()), r = [];
			for (let i of t) n.has(i.url) || (await e.delete(i), r.push(i.url));
			return { deletedURLs: r };
		});
	}
	getURLsToCacheKeys() {
		return this._urlsToCacheKeys;
	}
	getCachedURLs() {
		return [...this._urlsToCacheKeys.keys()];
	}
	getCacheKeyForURL(e) {
		let t = new URL(e, location.href);
		return this._urlsToCacheKeys.get(t.href);
	}
	getIntegrityForCacheKey(e) {
		return this._cacheKeysToIntegrities.get(e);
	}
	async matchPrecache(e) {
		let t = e instanceof Request ? e.url : e, n = this.getCacheKeyForURL(t);
		if (n) return (await self.caches.open(this.strategy.cacheName)).match(n);
	}
	createHandlerBoundToURL(e) {
		let t = this.getCacheKeyForURL(e);
		if (!t) throw new n("non-precached-url", { url: e });
		return (n) => (n.request = new Request(e), n.params = Object.assign({ cacheKey: t }, n.params), this.strategy.handle(n));
	}
};
//#endregion
//#region node_modules/workbox-routing/_version.js
try {
	self["workbox:routing:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-routing/utils/normalizeHandler.js
var O = (e) => e && typeof e == "object" ? e : { handle: e }, k = class {
	constructor(e, t, n = "GET") {
		this.handler = O(t), this.match = e, this.method = n;
	}
	setCatchHandler(e) {
		this.catchHandler = O(e);
	}
}, A = class extends k {
	constructor(e, t, n) {
		super(({ url: t }) => {
			let n = e.exec(t.href);
			if (n && !(t.origin !== location.origin && n.index !== 0)) return n.slice(1);
		}, t, n);
	}
}, j = class {
	constructor() {
		this._routes = /* @__PURE__ */ new Map(), this._defaultHandlerMap = /* @__PURE__ */ new Map();
	}
	get routes() {
		return this._routes;
	}
	addFetchListener() {
		self.addEventListener("fetch", ((e) => {
			let { request: t } = e, n = this.handleRequest({
				request: t,
				event: e
			});
			n && e.respondWith(n);
		}));
	}
	addCacheListener() {
		self.addEventListener("message", ((e) => {
			if (e.data && e.data.type === "CACHE_URLS") {
				let { payload: t } = e.data, n = Promise.all(t.urlsToCache.map((t) => {
					typeof t == "string" && (t = [t]);
					let n = new Request(...t);
					return this.handleRequest({
						request: n,
						event: e
					});
				}));
				e.waitUntil(n), e.ports && e.ports[0] && n.then(() => e.ports[0].postMessage(!0));
			}
		}));
	}
	handleRequest({ request: e, event: t }) {
		let n = new URL(e.url, location.href);
		if (!n.protocol.startsWith("http")) return;
		let r = n.origin === location.origin, { params: i, route: a } = this.findMatchingRoute({
			event: t,
			request: e,
			sameOrigin: r,
			url: n
		}), o = a && a.handler, s = e.method;
		if (!o && this._defaultHandlerMap.has(s) && (o = this._defaultHandlerMap.get(s)), !o) return;
		let c;
		try {
			c = o.handle({
				url: n,
				request: e,
				event: t,
				params: i
			});
		} catch (e) {
			c = Promise.reject(e);
		}
		let l = a && a.catchHandler;
		return c instanceof Promise && (this._catchHandler || l) && (c = c.catch(async (r) => {
			if (l) try {
				return await l.handle({
					url: n,
					request: e,
					event: t,
					params: i
				});
			} catch (e) {
				e instanceof Error && (r = e);
			}
			if (this._catchHandler) return this._catchHandler.handle({
				url: n,
				request: e,
				event: t
			});
			throw r;
		})), c;
	}
	findMatchingRoute({ url: e, sameOrigin: t, request: n, event: r }) {
		let i = this._routes.get(n.method) || [];
		for (let a of i) {
			let i, o = a.match({
				url: e,
				sameOrigin: t,
				request: n,
				event: r
			});
			if (o) return i = o, (Array.isArray(i) && i.length === 0 || o.constructor === Object && Object.keys(o).length === 0 || typeof o == "boolean") && (i = void 0), {
				route: a,
				params: i
			};
		}
		return {};
	}
	setDefaultHandler(e, t = "GET") {
		this._defaultHandlerMap.set(t, O(e));
	}
	setCatchHandler(e) {
		this._catchHandler = O(e);
	}
	registerRoute(e) {
		this._routes.has(e.method) || this._routes.set(e.method, []), this._routes.get(e.method).push(e);
	}
	unregisterRoute(e) {
		if (!this._routes.has(e.method)) throw new n("unregister-route-but-not-found-with-method", { method: e.method });
		let t = this._routes.get(e.method).indexOf(e);
		if (t > -1) this._routes.get(e.method).splice(t, 1);
		else throw new n("unregister-route-route-not-registered");
	}
};
//#endregion
//#region node_modules/workbox-core/registerQuotaErrorCallback.js
function M(e) {
	b.add(e);
}
//#endregion
//#region node_modules/workbox-core/_private/dontWaitFor.js
function N(e) {
	e.then(() => {});
}
//#endregion
//#region node_modules/workbox-core/_private/resultingClientExists.js
var P = 2e3;
async function ee(e) {
	if (!e) return;
	let t = await self.clients.matchAll({ type: "window" }), n = new Set(t.map((e) => e.id)), r, i = performance.now();
	for (; performance.now() - i < P && (t = await self.clients.matchAll({ type: "window" }), r = t.find((t) => e ? t.id === e : !n.has(t.id)), !r);) await S(100);
	return r;
}
//#endregion
//#region node_modules/workbox-core/setCacheNameDetails.js
function te(e) {
	o.updateDetails(e);
}
//#endregion
//#region node_modules/workbox-strategies/CacheFirst.js
var ne = class extends T {
	async _handle(e, t) {
		let r = await t.cacheMatch(e), i;
		if (!r) try {
			r = await t.fetchAndCachePut(e);
		} catch (e) {
			e instanceof Error && (i = e);
		}
		if (!r) throw new n("no-response", {
			url: e.url,
			error: i
		});
		return r;
	}
}, re = class extends T {
	async _handle(e, t) {
		let r = await t.cacheMatch(e);
		if (!r) throw new n("no-response", { url: e.url });
		return r;
	}
}, ie = { cacheWillUpdate: async ({ response: e }) => e.status === 200 || e.status === 0 ? e : null }, ae = class extends T {
	constructor(e = {}) {
		super(e), this.plugins.some((e) => "cacheWillUpdate" in e) || this.plugins.unshift(ie), this._networkTimeoutSeconds = e.networkTimeoutSeconds || 0;
	}
	async _handle(e, t) {
		let r = [], i = [], a;
		if (this._networkTimeoutSeconds) {
			let { id: n, promise: o } = this._getTimeoutPromise({
				request: e,
				logs: r,
				handler: t
			});
			a = n, i.push(o);
		}
		let o = this._getNetworkPromise({
			timeoutId: a,
			request: e,
			logs: r,
			handler: t
		});
		i.push(o);
		let s = await t.waitUntil((async () => await t.waitUntil(Promise.race(i)) || await o)());
		if (!s) throw new n("no-response", { url: e.url });
		return s;
	}
	_getTimeoutPromise({ request: e, logs: t, handler: n }) {
		let r;
		return {
			promise: new Promise((t) => {
				r = setTimeout(async () => {
					t(await n.cacheMatch(e));
				}, this._networkTimeoutSeconds * 1e3);
			}),
			id: r
		};
	}
	async _getNetworkPromise({ timeoutId: e, request: t, logs: n, handler: r }) {
		let i, a;
		try {
			a = await r.fetchAndCachePut(t);
		} catch (e) {
			e instanceof Error && (i = e);
		}
		return e && clearTimeout(e), (i || !a) && (a = await r.cacheMatch(t)), a;
	}
}, oe = class extends T {
	constructor(e = {}) {
		super(e), this._networkTimeoutSeconds = e.networkTimeoutSeconds || 0;
	}
	async _handle(e, t) {
		let r, i;
		try {
			let n = [t.fetch(e)];
			if (this._networkTimeoutSeconds) {
				let e = S(this._networkTimeoutSeconds * 1e3);
				n.push(e);
			}
			if (i = await Promise.race(n), !i) throw Error(`Timed out the network response after ${this._networkTimeoutSeconds} seconds.`);
		} catch (e) {
			e instanceof Error && (r = e);
		}
		if (!i) throw new n("no-response", {
			url: e.url,
			error: r
		});
		return i;
	}
}, se = class extends T {
	constructor(e = {}) {
		super(e), this.plugins.some((e) => "cacheWillUpdate" in e) || this.plugins.unshift(ie);
	}
	async _handle(e, t) {
		let r = t.fetchAndCachePut(e).catch(() => {});
		t.waitUntil(r);
		let i = await t.cacheMatch(e), a;
		if (!i) try {
			i = await r;
		} catch (e) {
			e instanceof Error && (a = e);
		}
		if (!i) throw new n("no-response", {
			url: e.url,
			error: a
		});
		return i;
	}
}, ce = (e, t) => t.some((t) => e instanceof t), le, ue;
function de() {
	return le ||= [
		IDBDatabase,
		IDBObjectStore,
		IDBIndex,
		IDBCursor,
		IDBTransaction
	];
}
function fe() {
	return ue ||= [
		IDBCursor.prototype.advance,
		IDBCursor.prototype.continue,
		IDBCursor.prototype.continuePrimaryKey
	];
}
var F = /* @__PURE__ */ new WeakMap(), I = /* @__PURE__ */ new WeakMap(), L = /* @__PURE__ */ new WeakMap(), R = /* @__PURE__ */ new WeakMap(), z = /* @__PURE__ */ new WeakMap();
function pe(e) {
	let t = new Promise((t, n) => {
		let r = () => {
			e.removeEventListener("success", i), e.removeEventListener("error", a);
		}, i = () => {
			t(V(e.result)), r();
		}, a = () => {
			n(e.error), r();
		};
		e.addEventListener("success", i), e.addEventListener("error", a);
	});
	return t.then((t) => {
		t instanceof IDBCursor && F.set(t, e);
	}).catch(() => {}), z.set(t, e), t;
}
function me(e) {
	if (I.has(e)) return;
	let t = new Promise((t, n) => {
		let r = () => {
			e.removeEventListener("complete", i), e.removeEventListener("error", a), e.removeEventListener("abort", a);
		}, i = () => {
			t(), r();
		}, a = () => {
			n(e.error || new DOMException("AbortError", "AbortError")), r();
		};
		e.addEventListener("complete", i), e.addEventListener("error", a), e.addEventListener("abort", a);
	});
	I.set(e, t);
}
var B = {
	get(e, t, n) {
		if (e instanceof IDBTransaction) {
			if (t === "done") return I.get(e);
			if (t === "objectStoreNames") return e.objectStoreNames || L.get(e);
			if (t === "store") return n.objectStoreNames[1] ? void 0 : n.objectStore(n.objectStoreNames[0]);
		}
		return V(e[t]);
	},
	set(e, t, n) {
		return e[t] = n, !0;
	},
	has(e, t) {
		return e instanceof IDBTransaction && (t === "done" || t === "store") ? !0 : t in e;
	}
};
function he(e) {
	B = e(B);
}
function ge(e) {
	return e === IDBDatabase.prototype.transaction && !("objectStoreNames" in IDBTransaction.prototype) ? function(t, ...n) {
		let r = e.call(H(this), t, ...n);
		return L.set(r, t.sort ? t.sort() : [t]), V(r);
	} : fe().includes(e) ? function(...t) {
		return e.apply(H(this), t), V(F.get(this));
	} : function(...t) {
		return V(e.apply(H(this), t));
	};
}
function _e(e) {
	return typeof e == "function" ? ge(e) : (e instanceof IDBTransaction && me(e), ce(e, de()) ? new Proxy(e, B) : e);
}
function V(e) {
	if (e instanceof IDBRequest) return pe(e);
	if (R.has(e)) return R.get(e);
	let t = _e(e);
	return t !== e && (R.set(e, t), z.set(t, e)), t;
}
var H = (e) => z.get(e);
//#endregion
//#region node_modules/idb/build/index.js
function U(e, t, { blocked: n, upgrade: r, blocking: i, terminated: a } = {}) {
	let o = indexedDB.open(e, t), s = V(o);
	return r && o.addEventListener("upgradeneeded", (e) => {
		r(V(o.result), e.oldVersion, e.newVersion, V(o.transaction), e);
	}), n && o.addEventListener("blocked", (e) => n(e.oldVersion, e.newVersion, e)), s.then((e) => {
		a && e.addEventListener("close", () => a()), i && e.addEventListener("versionchange", (e) => i(e.oldVersion, e.newVersion, e));
	}).catch(() => {}), s;
}
function ve(e, { blocked: t } = {}) {
	let n = indexedDB.deleteDatabase(e);
	return t && n.addEventListener("blocked", (e) => t(e.oldVersion, e)), V(n).then(() => void 0);
}
var ye = [
	"get",
	"getKey",
	"getAll",
	"getAllKeys",
	"count"
], be = [
	"put",
	"add",
	"delete",
	"clear"
], W = /* @__PURE__ */ new Map();
function G(e, t) {
	if (!(e instanceof IDBDatabase && !(t in e) && typeof t == "string")) return;
	if (W.get(t)) return W.get(t);
	let n = t.replace(/FromIndex$/, ""), r = t !== n, i = be.includes(n);
	if (!(n in (r ? IDBIndex : IDBObjectStore).prototype) || !(i || ye.includes(n))) return;
	let a = async function(e, ...t) {
		let a = this.transaction(e, i ? "readwrite" : "readonly"), o = a.store;
		return r && (o = o.index(t.shift())), (await Promise.all([o[n](...t), i && a.done]))[0];
	};
	return W.set(t, a), a;
}
he((e) => ({
	...e,
	get: (t, n, r) => G(t, n) || e.get(t, n, r),
	has: (t, n) => !!G(t, n) || e.has(t, n)
}));
//#endregion
//#region node_modules/workbox-expiration/_version.js
//#endregion
//#region node_modules/workbox-expiration/_version.js
try {
	self["workbox:expiration:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-expiration/models/CacheTimestampsModel.js
var xe = "workbox-expiration", K = "cache-entries", q = (e) => {
	let t = new URL(e, location.href);
	return t.hash = "", t.href;
}, Se = class {
	constructor(e) {
		this._db = null, this._cacheName = e;
	}
	_upgradeDb(e) {
		let t = e.createObjectStore(K, { keyPath: "id" });
		t.createIndex("cacheName", "cacheName", { unique: !1 }), t.createIndex("timestamp", "timestamp", { unique: !1 });
	}
	_upgradeDbAndDeleteOldDbs(e) {
		this._upgradeDb(e), this._cacheName && ve(this._cacheName);
	}
	async setTimestamp(e, t) {
		e = q(e);
		let n = {
			url: e,
			timestamp: t,
			cacheName: this._cacheName,
			id: this._getId(e)
		}, r = (await this.getDb()).transaction(K, "readwrite", { durability: "relaxed" });
		await r.store.put(n), await r.done;
	}
	async getTimestamp(e) {
		return (await (await this.getDb()).get(K, this._getId(e)))?.timestamp;
	}
	async expireEntries(e, t) {
		let n = await this.getDb(), r = await n.transaction(K).store.index("timestamp").openCursor(null, "prev"), i = [], a = 0;
		for (; r;) {
			let n = r.value;
			n.cacheName === this._cacheName && (e && n.timestamp < e || t && a >= t ? i.push(r.value) : a++), r = await r.continue();
		}
		let o = [];
		for (let e of i) await n.delete(K, e.id), o.push(e.url);
		return o;
	}
	_getId(e) {
		return this._cacheName + "|" + q(e);
	}
	async getDb() {
		return this._db ||= await U(xe, 1, { upgrade: this._upgradeDbAndDeleteOldDbs.bind(this) }), this._db;
	}
}, Ce = class {
	constructor(e, t = {}) {
		this._isRunning = !1, this._rerunRequested = !1, this._maxEntries = t.maxEntries, this._maxAgeSeconds = t.maxAgeSeconds, this._matchOptions = t.matchOptions, this._cacheName = e, this._timestampModel = new Se(e);
	}
	async expireEntries() {
		if (this._isRunning) {
			this._rerunRequested = !0;
			return;
		}
		this._isRunning = !0;
		let e = this._maxAgeSeconds ? Date.now() - this._maxAgeSeconds * 1e3 : 0, t = await this._timestampModel.expireEntries(e, this._maxEntries), n = await self.caches.open(this._cacheName);
		for (let e of t) await n.delete(e, this._matchOptions);
		this._isRunning = !1, this._rerunRequested && (this._rerunRequested = !1, N(this.expireEntries()));
	}
	async updateTimestamp(e) {
		await this._timestampModel.setTimestamp(e, Date.now());
	}
	async isURLExpired(e) {
		if (this._maxAgeSeconds) {
			let t = await this._timestampModel.getTimestamp(e), n = Date.now() - this._maxAgeSeconds * 1e3;
			return t === void 0 ? !0 : t < n;
		} else return !1;
	}
	async delete() {
		this._rerunRequested = !1, await this._timestampModel.expireEntries(Infinity);
	}
}, we = class {
	constructor(e = {}) {
		this.cachedResponseWillBeUsed = async ({ event: e, request: t, cacheName: n, cachedResponse: r }) => {
			if (!r) return null;
			let i = this._isResponseDateFresh(r), a = this._getCacheExpiration(n);
			N(a.expireEntries());
			let o = a.updateTimestamp(t.url);
			if (e) try {
				e.waitUntil(o);
			} catch {}
			return i ? r : null;
		}, this.cacheDidUpdate = async ({ cacheName: e, request: t }) => {
			let n = this._getCacheExpiration(e);
			await n.updateTimestamp(t.url), await n.expireEntries();
		}, this._config = e, this._maxAgeSeconds = e.maxAgeSeconds, this._cacheExpirations = /* @__PURE__ */ new Map(), e.purgeOnQuotaError && M(() => this.deleteCacheAndMetadata());
	}
	_getCacheExpiration(e) {
		if (e === o.getRuntimeName()) throw new n("expire-custom-caches-only");
		let t = this._cacheExpirations.get(e);
		return t || (t = new Ce(e, this._config), this._cacheExpirations.set(e, t)), t;
	}
	_isResponseDateFresh(e) {
		if (!this._maxAgeSeconds) return !0;
		let t = this._getDateHeaderTimestamp(e);
		return t === null ? !0 : t >= Date.now() - this._maxAgeSeconds * 1e3;
	}
	_getDateHeaderTimestamp(e) {
		if (!e.headers.has("date")) return null;
		let t = e.headers.get("date"), n = new Date(t).getTime();
		return isNaN(n) ? null : n;
	}
	async deleteCacheAndMetadata() {
		for (let [e, t] of this._cacheExpirations) await self.caches.delete(e), await t.delete();
		this._cacheExpirations = /* @__PURE__ */ new Map();
	}
};
//#endregion
//#region node_modules/workbox-cacheable-response/_version.js
try {
	self["workbox:cacheable-response:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-cacheable-response/CacheableResponse.js
var Te = class {
	constructor(e = {}) {
		this._statuses = e.statuses, this._headers = e.headers;
	}
	isResponseCacheable(e) {
		let t = !0;
		return this._statuses && (t = this._statuses.includes(e.status)), this._headers && t && (t = Object.keys(this._headers).some((t) => e.headers.get(t) === this._headers[t])), t;
	}
}, Ee = class {
	constructor(e) {
		this.cacheWillUpdate = async ({ response: e }) => this._cacheableResponse.isResponseCacheable(e) ? e : null, this._cacheableResponse = new Te(e);
	}
};
//#endregion
//#region node_modules/workbox-background-sync/_version.js
try {
	self["workbox:background-sync:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-background-sync/lib/QueueDb.js
var De = 3, Oe = "workbox-background-sync", J = "requests", Y = "queueName", ke = class {
	constructor() {
		this._db = null;
	}
	async addEntry(e) {
		let t = (await this.getDb()).transaction(J, "readwrite", { durability: "relaxed" });
		await t.store.add(e), await t.done;
	}
	async getFirstEntryId() {
		return (await (await this.getDb()).transaction(J).store.openCursor())?.value.id;
	}
	async getAllEntriesByQueueName(e) {
		return await (await this.getDb()).getAllFromIndex(J, Y, IDBKeyRange.only(e)) || [];
	}
	async getEntryCountByQueueName(e) {
		return (await this.getDb()).countFromIndex(J, Y, IDBKeyRange.only(e));
	}
	async deleteEntry(e) {
		await (await this.getDb()).delete(J, e);
	}
	async getFirstEntryByQueueName(e) {
		return await this.getEndEntryFromIndex(IDBKeyRange.only(e), "next");
	}
	async getLastEntryByQueueName(e) {
		return await this.getEndEntryFromIndex(IDBKeyRange.only(e), "prev");
	}
	async getEndEntryFromIndex(e, t) {
		return (await (await this.getDb()).transaction(J).store.index(Y).openCursor(e, t))?.value;
	}
	async getDb() {
		return this._db ||= await U(Oe, De, { upgrade: this._upgradeDb }), this._db;
	}
	_upgradeDb(e, t) {
		t > 0 && t < De && e.objectStoreNames.contains(J) && e.deleteObjectStore(J), e.createObjectStore(J, {
			autoIncrement: !0,
			keyPath: "id"
		}).createIndex(Y, Y, { unique: !1 });
	}
}, Ae = class {
	constructor(e) {
		this._queueName = e, this._queueDb = new ke();
	}
	async pushEntry(e) {
		delete e.id, e.queueName = this._queueName, await this._queueDb.addEntry(e);
	}
	async unshiftEntry(e) {
		let t = await this._queueDb.getFirstEntryId();
		t ? e.id = t - 1 : delete e.id, e.queueName = this._queueName, await this._queueDb.addEntry(e);
	}
	async popEntry() {
		return this._removeEntry(await this._queueDb.getLastEntryByQueueName(this._queueName));
	}
	async shiftEntry() {
		return this._removeEntry(await this._queueDb.getFirstEntryByQueueName(this._queueName));
	}
	async getAll() {
		return await this._queueDb.getAllEntriesByQueueName(this._queueName);
	}
	async size() {
		return await this._queueDb.getEntryCountByQueueName(this._queueName);
	}
	async deleteEntry(e) {
		await this._queueDb.deleteEntry(e);
	}
	async _removeEntry(e) {
		return e && await this.deleteEntry(e.id), e;
	}
}, je = [
	"method",
	"referrer",
	"referrerPolicy",
	"mode",
	"credentials",
	"cache",
	"redirect",
	"integrity",
	"keepalive"
], Me = class e {
	static async fromRequest(t) {
		let n = {
			url: t.url,
			headers: {}
		};
		t.method !== "GET" && (n.body = await t.clone().arrayBuffer());
		for (let [e, r] of t.headers.entries()) n.headers[e] = r;
		for (let e of je) t[e] !== void 0 && (n[e] = t[e]);
		return new e(n);
	}
	constructor(e) {
		e.mode === "navigate" && (e.mode = "same-origin"), this._requestData = e;
	}
	toObject() {
		let e = Object.assign({}, this._requestData);
		return e.headers = Object.assign({}, this._requestData.headers), e.body &&= e.body.slice(0), e;
	}
	toRequest() {
		return new Request(this._requestData.url, this._requestData);
	}
	clone() {
		return new e(this.toObject());
	}
}, Ne = "workbox-background-sync", Pe = 1440 * 7, X = /* @__PURE__ */ new Set(), Fe = (e) => {
	let t = {
		request: new Me(e.requestData).toRequest(),
		timestamp: e.timestamp
	};
	return e.metadata && (t.metadata = e.metadata), t;
}, Ie = class {
	constructor(e, { forceSyncFallback: t, onSync: r, maxRetentionTime: i } = {}) {
		if (this._syncInProgress = !1, this._requestsAddedDuringSync = !1, X.has(e)) throw new n("duplicate-queue-name", { name: e });
		X.add(e), this._name = e, this._onSync = r || this.replayRequests, this._maxRetentionTime = i || Pe, this._forceSyncFallback = !!t, this._queueStore = new Ae(this._name), this._addSyncListener();
	}
	get name() {
		return this._name;
	}
	async pushRequest(e) {
		await this._addRequest(e, "push");
	}
	async unshiftRequest(e) {
		await this._addRequest(e, "unshift");
	}
	async popRequest() {
		return this._removeRequest("pop");
	}
	async shiftRequest() {
		return this._removeRequest("shift");
	}
	async getAll() {
		let e = await this._queueStore.getAll(), t = Date.now(), n = [];
		for (let r of e) {
			let e = this._maxRetentionTime * 60 * 1e3;
			t - r.timestamp > e ? await this._queueStore.deleteEntry(r.id) : n.push(Fe(r));
		}
		return n;
	}
	async size() {
		return await this._queueStore.size();
	}
	async _addRequest({ request: e, metadata: t, timestamp: n = Date.now() }, r) {
		let i = {
			requestData: (await Me.fromRequest(e.clone())).toObject(),
			timestamp: n
		};
		switch (t && (i.metadata = t), r) {
			case "push":
				await this._queueStore.pushEntry(i);
				break;
			case "unshift":
				await this._queueStore.unshiftEntry(i);
				break;
		}
		this._syncInProgress ? this._requestsAddedDuringSync = !0 : await this.registerSync();
	}
	async _removeRequest(e) {
		let t = Date.now(), n;
		switch (e) {
			case "pop":
				n = await this._queueStore.popEntry();
				break;
			case "shift":
				n = await this._queueStore.shiftEntry();
				break;
		}
		if (n) {
			let r = this._maxRetentionTime * 60 * 1e3;
			return t - n.timestamp > r ? this._removeRequest(e) : Fe(n);
		} else return;
	}
	async replayRequests() {
		let e;
		for (; e = await this.shiftRequest();) try {
			await fetch(e.request.clone());
		} catch {
			throw await this.unshiftRequest(e), new n("queue-replay-failed", { name: this._name });
		}
	}
	async registerSync() {
		if ("sync" in self.registration && !this._forceSyncFallback) try {
			await self.registration.sync.register(`${Ne}:${this._name}`);
		} catch {}
	}
	_addSyncListener() {
		"sync" in self.registration && !this._forceSyncFallback ? self.addEventListener("sync", (e) => {
			e.tag === `${Ne}:${this._name}` && e.waitUntil((async () => {
				this._syncInProgress = !0;
				let t;
				try {
					await this._onSync({ queue: this });
				} catch (e) {
					if (e instanceof Error) throw t = e, t;
				} finally {
					this._requestsAddedDuringSync && !(t && !e.lastChance) && await this.registerSync(), this._syncInProgress = !1, this._requestsAddedDuringSync = !1;
				}
			})());
		}) : this._onSync({ queue: this });
	}
	static get _queueNames() {
		return X;
	}
}, Le = class {
	constructor(e, t) {
		this.fetchDidFail = async ({ request: e }) => {
			await this._queue.pushRequest({ request: e });
		}, this._queue = new Ie(e, t);
	}
};
//#endregion
//#region node_modules/workbox-broadcast-update/_version.js
try {
	self["workbox:broadcast-update:7.3.0"] && _();
} catch {}
//#endregion
//#region node_modules/workbox-broadcast-update/responsesAreSame.js
var Re = (e, t, n) => n.some((n) => e.headers.has(n) && t.headers.has(n)) ? n.every((n) => {
	let r = e.headers.has(n) === t.headers.has(n), i = e.headers.get(n) === t.headers.get(n);
	return r && i;
}) : !0, ze = "CACHE_UPDATED", Be = "workbox-broadcast-update", Ve = [
	"content-length",
	"etag",
	"last-modified"
], He = /^((?!chrome|android).)*safari/i.test(navigator.userAgent);
function Ue(e) {
	return {
		cacheName: e.cacheName,
		updatedURL: e.request.url
	};
}
var We = class {
	constructor({ generatePayload: e, headersToCheck: t, notifyAllClients: n } = {}) {
		this._headersToCheck = t || Ve, this._generatePayload = e || Ue, this._notifyAllClients = n ?? !0;
	}
	async notifyIfUpdated(e) {
		if (e.oldResponse && !Re(e.oldResponse, e.newResponse, this._headersToCheck)) {
			let t = {
				type: ze,
				meta: Be,
				payload: this._generatePayload(e)
			};
			if (e.request.mode === "navigate") {
				let t;
				e.event instanceof FetchEvent && (t = e.event.resultingClientId), (!await ee(t) || He) && await S(3500);
			}
			if (this._notifyAllClients) {
				let e = await self.clients.matchAll({ type: "window" });
				for (let n of e) n.postMessage(t);
			} else e.event instanceof FetchEvent && (await self.clients.get(e.event.clientId))?.postMessage(t);
		}
	}
}, Ge = class {
	constructor(e) {
		this.cacheDidUpdate = async (e) => {
			N(this._broadcastUpdate.notifyIfUpdated(e));
		}, this._broadcastUpdate = new We(e);
	}
}, Ke = class e extends A {
	constructor(e) {
		let t = new RegExp(e.pattern || ".*"), n = e.prefix || "";
		if (n) {
			let e = decodeURIComponent(new URL(globalThis.location.href).searchParams.get("base-path") || "");
			e ||= new URL(self.registration.scope).pathname, n = new URL(n, `http://host${e.replace(/\/$/, "")}/`).pathname;
		}
		let r, i = [];
		e.strategy !== "network-only" && i.push(new Ee({ statuses: e.statuses || [
			0,
			200,
			201,
			202,
			204
		] })), e.maxAgeSeconds && i.push(new we({ maxAgeSeconds: e.maxAgeSeconds })), e.syncRetentionMinutes && i.push(new Le("polyfea", { maxRetentionTime: e.syncRetentionMinutes })), e.strategy === "stale-while-revalidate" && i.push(new Ge());
		let a = e.cacheName || "polyfea-run-time-v1";
		switch (e.strategy) {
			case "cache-first":
				r = new ne({
					cacheName: a,
					plugins: i
				});
				break;
			case "cache-only":
				r = new re({
					cacheName: a,
					plugins: i
				});
				break;
			case "network-first":
				r = new ae({
					cacheName: a,
					plugins: i
				});
				break;
			case "network-only":
				r = new oe({ plugins: i });
				break;
			case "stale-while-revalidate":
				r = new se({
					cacheName: a,
					plugins: i
				});
				break;
			default:
				r = new ne({
					cacheName: a,
					plugins: i
				});
				break;
		}
		super(t, r, e.method || "GET");
		let o = this.match;
		this.match = (t) => e.destination && t.request.destination !== e.destination || n && !t.url.pathname.startsWith(n) ? !1 : o(t);
	}
	static from(t) {
		return new e(t);
	}
}, qe = /* @__PURE__ */ e(((e, t) => {
	function n(e) {
		try {
			return JSON.stringify(e);
		} catch {
			return "\"[Circular]\"";
		}
	}
	t.exports = r;
	function r(e, t, r) {
		var i = r && r.stringify || n, a = 1;
		if (typeof e == "object" && e) {
			var o = t.length + a;
			if (o === 1) return e;
			var s = Array(o);
			s[0] = i(e);
			for (var c = 1; c < o; c++) s[c] = i(t[c]);
			return s.join(" ");
		}
		if (typeof e != "string") return e;
		var l = t.length;
		if (l === 0) return e;
		for (var u = "", d = 1 - a, f = -1, p = e && e.length || 0, m = 0; m < p;) {
			if (e.charCodeAt(m) === 37 && m + 1 < p) {
				switch (f = f > -1 ? f : 0, e.charCodeAt(m + 1)) {
					case 100:
					case 102:
						if (d >= l || t[d] == null) break;
						f < m && (u += e.slice(f, m)), u += Number(t[d]), f = m + 2, m++;
						break;
					case 105:
						if (d >= l || t[d] == null) break;
						f < m && (u += e.slice(f, m)), u += Math.floor(Number(t[d])), f = m + 2, m++;
						break;
					case 79:
					case 111:
					case 106:
						if (d >= l || t[d] === void 0) break;
						f < m && (u += e.slice(f, m));
						var h = typeof t[d];
						if (h === "string") {
							u += "'" + t[d] + "'", f = m + 2, m++;
							break;
						}
						if (h === "function") {
							u += t[d].name || "<anonymous>", f = m + 2, m++;
							break;
						}
						u += i(t[d]), f = m + 2, m++;
						break;
					case 115:
						if (d >= l) break;
						f < m && (u += e.slice(f, m)), u += String(t[d]), f = m + 2, m++;
						break;
					case 37:
						f < m && (u += e.slice(f, m)), u += "%", f = m + 2, m++, d--;
						break;
				}
				++d;
			}
			++m;
		}
		return f === -1 ? e : (f < p && (u += e.slice(f)), u);
	}
})), Z = (/* @__PURE__ */ e(((e, t) => {
	var n = qe();
	t.exports = f;
	var r = N().console || {}, i = {
		mapHttpRequest: E,
		mapHttpResponse: E,
		wrapRequestSerializer: D,
		wrapResponseSerializer: D,
		wrapErrorSerializer: D,
		req: E,
		res: E,
		err: w,
		errWithCause: w
	};
	function a(e, t) {
		return e === "silent" ? Infinity : t.levels.values[e];
	}
	var o = Symbol("pino.logFuncs"), s = Symbol("pino.hierarchy"), c = {
		error: "log",
		fatal: "error",
		warn: "error",
		info: "log",
		debug: "log",
		trace: "log"
	};
	function l(e, t) {
		t[s] = {
			logger: t,
			parent: e[s]
		};
	}
	function u(e, t, n) {
		let i = {};
		t.forEach((e) => {
			i[e] = n[e] ? n[e] : r[e] || r[c[e] || "log"] || O;
		}), e[o] = i;
	}
	function d(e, t) {
		return Array.isArray(e) ? e.filter(function(e) {
			return e !== "!stdSerializers.err";
		}) : e === !0 ? Object.keys(t) : !1;
	}
	function f(e) {
		e ||= {}, e.browser = e.browser || {};
		let t = e.browser.transmit;
		if (t && typeof t.send != "function") throw Error("pino: transmit option must have a send function");
		let n = e.browser.write || r;
		e.browser.write && (e.browser.asObject = !0);
		let i = e.serializers || {}, o = d(e.browser.serialize, i), s = e.browser.serialize;
		Array.isArray(e.browser.serialize) && e.browser.serialize.indexOf("!stdSerializers.err") > -1 && (s = !1);
		let c = Object.keys(e.customLevels || {}), f = [
			"error",
			"fatal",
			"warn",
			"info",
			"debug",
			"trace"
		].concat(c);
		typeof n == "function" && f.forEach(function(e) {
			n[e] = n;
		}), (e.enabled === !1 || e.browser.disabled) && (e.level = "silent");
		let m = e.level || "info", h = Object.create(n);
		h.log ||= O, u(h, f, n), l({}, h), Object.defineProperty(h, "levelVal", { get: y }), Object.defineProperty(h, "level", {
			get: b,
			set: S
		});
		let v = {
			transmit: t,
			serialize: o,
			asObject: e.browser.asObject,
			asObjectBindingsOnly: e.browser.asObjectBindingsOnly,
			formatters: e.browser.formatters,
			reportCaller: e.browser.reportCaller,
			levels: f,
			timestamp: T(e),
			messageKey: e.messageKey || "msg",
			onChild: e.onChild || O
		};
		h.levels = p(e), h.level = m, h.isLevelEnabled = function(e) {
			return this.levels.values[e] ? this.levels.values[e] >= this.levels.values[this.level] : !1;
		}, h.setMaxListeners = h.getMaxListeners = h.emit = h.addListener = h.on = h.prependListener = h.once = h.prependOnceListener = h.removeListener = h.removeAllListeners = h.listeners = h.listenerCount = h.eventNames = h.write = h.flush = O, h.serializers = i, h._serialize = o, h._stdErrSerialize = s, h.child = function(...e) {
			return w.call(this, v, ...e);
		}, t && (h._logEvent = C());
		function y() {
			return a(this.level, this);
		}
		function b() {
			return this._level;
		}
		function S(e) {
			if (e !== "silent" && !this.levels.values[e]) throw Error("unknown level " + e);
			this._level = e, g(this, v, h, "error"), g(this, v, h, "fatal"), g(this, v, h, "warn"), g(this, v, h, "info"), g(this, v, h, "debug"), g(this, v, h, "trace"), c.forEach((e) => {
				g(this, v, h, e);
			});
		}
		function w(n, r, a) {
			if (!r) throw Error("missing bindings for child Pino");
			a ||= {}, o && r.serializers && (a.serializers = r.serializers);
			let s = a.serializers;
			if (o && s) {
				var c = Object.assign({}, i, s), u = e.browser.serialize === !0 ? Object.keys(c) : o;
				delete r.serializers, x([r], u, c, this._stdErrSerialize);
			}
			function d(e) {
				this._childLevel = (e._childLevel | 0) + 1, this.bindings = r, c && (this.serializers = c, this._serialize = u), t && (this._logEvent = C([].concat(e._logEvent.bindings, r)));
			}
			d.prototype = this;
			let f = new d(this);
			return l(this, f), f.child = function(...e) {
				return w.call(this, n, ...e);
			}, f.level = a.level || this.level, n.onChild(f), f;
		}
		return h;
	}
	function p(e) {
		let t = e.customLevels || {};
		return {
			values: Object.assign({}, f.levels.values, t),
			labels: Object.assign({}, f.levels.labels, m(t))
		};
	}
	function m(e) {
		let t = {};
		return Object.keys(e).forEach(function(n) {
			t[e[n]] = n;
		}), t;
	}
	f.levels = {
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
	}, f.stdSerializers = i, f.stdTimeFunctions = Object.assign({}, {
		nullTime: k,
		epochTime: A,
		unixTime: j,
		isoTime: M
	});
	function h(e) {
		let t = [];
		e.bindings && t.push(e.bindings);
		let n = e[s];
		for (; n.parent;) n = n.parent, n.logger.bindings && t.push(n.logger.bindings);
		return t.reverse();
	}
	function g(e, t, n, r) {
		if (Object.defineProperty(e, r, {
			value: a(e.level, n) > a(r, n) ? O : n[o][r],
			writable: !0,
			enumerable: !0,
			configurable: !0
		}), e[r] === O) {
			if (!t.transmit) return;
			let i = a(t.transmit.level || e.level, n);
			if (a(r, n) < i) return;
		}
		e[r] = y(e, t, n, r);
		let i = h(e);
		i.length !== 0 && (e[r] = v(i, e[r]));
	}
	function v(e, t) {
		return function() {
			return t.apply(this, [...e, ...arguments]);
		};
	}
	function y(e, t, n, i) {
		return (function(o) {
			return function() {
				let s = t.timestamp(), c = Array(arguments.length), l = Object.getPrototypeOf && Object.getPrototypeOf(this) === r ? r : this;
				for (var u = 0; u < c.length; u++) c[u] = arguments[u];
				var d = !1;
				if (t.serialize && (x(c, this._serialize, this.serializers, this._stdErrSerialize), d = !0), t.asObject || t.formatters) {
					let e = b(this, i, c, s, t);
					if (t.reportCaller && e && e.length > 0 && e[0] && typeof e[0] == "object") try {
						let t = P();
						t && (e[0].caller = t);
					} catch {}
					o.call(l, ...e);
				} else {
					if (t.reportCaller) try {
						let e = P();
						e && c.push(e);
					} catch {}
					o.apply(l, c);
				}
				if (t.transmit) {
					let r = t.transmit.level || e._level, o = a(r, n), l = a(i, n);
					if (l < o) return;
					S(this, {
						ts: s,
						methodLevel: i,
						methodValue: l,
						transmitLevel: r,
						transmitValue: n.levels.values[t.transmit.level || e._level],
						send: t.transmit.send,
						val: a(e._level, n)
					}, c, d);
				}
			};
		})(e[o][i]);
	}
	function b(e, t, r, i, a) {
		let { level: o, log: s = (e) => e } = a.formatters || {}, c = r.slice(), l = c[0], u = {}, d = (e._childLevel | 0) + 1;
		if (d < 1 && (d = 1), i && (u.time = i), o) {
			let n = o(t, e.levels.values[t]);
			Object.assign(u, n);
		} else u.level = e.levels.values[t];
		if (a.asObjectBindingsOnly) {
			if (typeof l == "object" && l) for (; d-- && typeof c[0] == "object";) Object.assign(u, c.shift());
			return [s(u), ...c];
		} else {
			if (typeof l == "object" && l) {
				for (; d-- && typeof c[0] == "object";) Object.assign(u, c.shift());
				l = c.length ? n(c.shift(), c) : void 0;
			} else typeof l == "string" && (l = n(c.shift(), c));
			return l !== void 0 && (u[a.messageKey] = l), [s(u)];
		}
	}
	function x(e, t, n, r) {
		for (let i in e) if (r && e[i] instanceof Error) e[i] = f.stdSerializers.err(e[i]);
		else if (typeof e[i] == "object" && !Array.isArray(e[i]) && t) for (let r in e[i]) t.indexOf(r) > -1 && r in n && (e[i][r] = n[r](e[i][r]));
	}
	function S(e, t, n, r = !1) {
		let i = t.send, a = t.ts, o = t.methodLevel, s = t.methodValue, c = t.val, l = e._logEvent.bindings;
		r || x(n, e._serialize || Object.keys(e.serializers), e.serializers, e._stdErrSerialize === void 0 ? !0 : e._stdErrSerialize), e._logEvent.ts = a, e._logEvent.messages = n.filter(function(e) {
			return l.indexOf(e) === -1;
		}), e._logEvent.level.label = o, e._logEvent.level.value = s, i(o, e._logEvent, c), e._logEvent = C(l);
	}
	function C(e) {
		return {
			ts: 0,
			messages: [],
			bindings: e || [],
			level: {
				label: "",
				value: 0
			}
		};
	}
	function w(e) {
		let t = {
			type: e.constructor.name,
			msg: e.message,
			stack: e.stack
		};
		for (let n in e) t[n] === void 0 && (t[n] = e[n]);
		return t;
	}
	function T(e) {
		return typeof e.timestamp == "function" ? e.timestamp : e.timestamp === !1 ? k : A;
	}
	function E() {
		return {};
	}
	function D(e) {
		return e;
	}
	function O() {}
	function k() {
		return !1;
	}
	function A() {
		return Date.now();
	}
	function j() {
		return Math.round(Date.now() / 1e3);
	}
	function M() {
		return new Date(Date.now()).toISOString();
	}
	/* istanbul ignore next */
	function N() {
		function e(e) {
			return e !== void 0 && e;
		}
		try {
			return typeof globalThis < "u" || Object.defineProperty(Object.prototype, "globalThis", {
				get: function() {
					return delete Object.prototype.globalThis, this.globalThis = this;
				},
				configurable: !0
			}), globalThis;
		} catch {
			return e(self) || e(window) || e(this) || {};
		}
	}
	t.exports.default = f, t.exports.pino = f;
	/* istanbul ignore next */
	function P() {
		let e = (/* @__PURE__ */ Error()).stack;
		if (!e) return null;
		let t = e.split("\n");
		for (let e = 1; e < t.length; e++) {
			let n = t[e].trim();
			if (/(^at\s+)?(createWrap|LOG|set\s*\(|asObject|Object\.apply|Function\.apply)/.test(n) || n.indexOf("browser.js") !== -1 || n.indexOf("node:internal") !== -1 || n.indexOf("node_modules") !== -1) continue;
			let r = n.match(/\((.*?):(\d+):(\d+)\)/);
			if (r ||= n.match(/at\s+(.*?):(\d+):(\d+)/), r) {
				let e = r[1], t = r[2], n = r[3];
				return e + ":" + t + ":" + n;
			}
		}
		return null;
	}
})))(), Q = self.__POLYFEA_SW_LOGS_LEVEL === void 0 ? self.__POLYFEA_LOGS_LEVEL : self.__POLYFEA_SW_LOGS_LEVEL;
Q === void 0 && (Q = Z.levels.values.info);
var $ = (0, Z.pino)({
	level: Z.levels.labels[Q],
	timestamp: Z.stdTimeFunctions.isoTime,
	browser: {
		asObject: !0,
		write: (e) => {
			let t = {
				trace: "#95a5a6",
				debug: "#7f8c8d",
				log: "#2ecc71",
				info: "#3498db",
				warn: "#f39c12",
				error: "#c0392b",
				fatal: "#c0392b"
			}, n = Z.levels.labels[e.level], r = [
				`background: ${t[n] || "#000"}`,
				"border-radius: 0.5em",
				"color: white",
				"font-weight: bold",
				"padding: 2px 0.5em"
			], i = "polyfea", a = /* @__PURE__ */ Error();
			/* v8 ignore next 3 */
			if (!a.stack) try {
				throw a;
			} catch {}
			let o = a.stack?.toString().split(/\r\n|\n/);
			e.component && (i += "/" + e.component), e = Object.assign(e, {
				module: "polyfea",
				level: n,
				src: o?.[1] || void 0
			});
			let s = ["%c" + i, r.join(";")];
			console[n](...s, e);
		}
	}
}).child({ component: "sw" });
//#endregion
//#region src/sw.ts
new class {
	constructor(e = self) {
		this.scope = e, this.router = new j(), this.interceptors = [], this.routesRestored = !1, te({
			prefix: "polyfea",
			suffix: "v1",
			precache: "install-time",
			runtime: "run-time"
		}), this.precacheController = new D();
		let t = new URL(globalThis.location.href).searchParams.get("reconcile-interval");
		this.reconcilationInterval = (parseInt(t || "") || 1800) * 1e3;
	}
	async start() {
		this.scope.addEventListener("install", (e) => this.install(e)), this.scope.addEventListener("activate", (e) => this.activate(e)), this.scope.addEventListener("fetch", (e) => this.handleFetch(e));
	}
	async reconcileRoutes(e = !1) {
		if (!this.routesRestored) {
			this.routesRestored = !0;
			let e = await this.getStoredConfig();
			e && await this.applyConfig(e);
		}
		let t = await this.getLastReconciliationTime(), n = 0;
		if (t && (n = Date.now() - parseInt(t)), !e && n && n < this.reconcilationInterval) {
			$.debug("Skipping reconciliation - data are fresh ");
			return;
		}
		let r = decodeURIComponent(new URL(globalThis.location.href).searchParams.get("caching-config") || "./polyfea-caching.json");
		try {
			let t = await fetch(r);
			if (t.status < 300) {
				let n = await t.json();
				await this.applyConfig(n, e), await this.setStoredConfig(n), await this.setLastReconciliationTime(Date.now().toString()), $.info(`Service worker reconciled: precached ${n.precache?.length || 0} files and added ${n.routes?.length || 0} routes`);
			}
		} catch (e) {
			$.warn({ err: e }, "Failed to reconcile routes");
		}
	}
	async applyConfig(e, t = !1) {
		let n = (e.precache || []).filter((e) => {
			let t = typeof e == "string" ? e : e.url;
			return !this.precacheController.getCacheKeyForURL(t);
		});
		if (this.precacheController.addToCacheList(n), !t && n.length > 0) try {
			let e = await caches.open("polyfea-install-time-v1"), t = n.map(async (t) => {
				let n = typeof t == "string" ? t : t.url, r = this.precacheController.getCacheKeyForURL(n);
				if (r && !await e.match(r)) {
					let t = await fetch(n);
					t.ok && await e.put(r, t);
				}
			});
			await Promise.all(t), $.debug("Dynamically populated missing items into install-time cache");
		} catch (e) {
			$.error({ err: e }, "Failed to dynamically populate install-time cache");
		}
		this.router.routes.clear(), e.routes?.map(Ke.from).forEach((e) => this.router.registerRoute(e)), this.interceptors = [];
		for (let t of e.interceptors || []) try {
			let e = await import(t.module);
			if (e && e.default && e.default.interceptor) {
				let n = t.name, r = t.options;
				this.interceptors.push(Object.assign({}, e.default, {
					name: n,
					intercept: (t, i) => {
						let a = e.default.interceptor(t, i, r);
						return a && $.debug(`Request ${t.url} handled by interceptor: ${n}`), a;
					}
				}));
			} else $.warn(`Interceptor module ${t.module} does not have a default export with an interceptor function`);
		} catch (e) {
			$.warn({ err: e }, `Failed to load interceptor module ${t.module}`);
		}
	}
	async getStoredConfig() {
		return new Promise((e, t) => {
			let n = indexedDB.open("polyfeaDB", 1);
			n.onerror = () => t(n.error), n.onupgradeneeded = () => {
				n.result.createObjectStore("reconciliationTime", { keyPath: "id" });
			}, n.onsuccess = () => {
				let t = n.result.transaction("reconciliationTime", "readonly").objectStore("reconciliationTime").get("cachedConfig");
				t.onerror = () => e(null), t.onsuccess = () => {
					let n = t.result?.value;
					if (n) try {
						e(JSON.parse(n));
					} catch {
						e(null);
					}
					else e(null);
				};
			};
		});
	}
	async setStoredConfig(e) {
		return new Promise((t, n) => {
			let r = indexedDB.open("polyfeaDB", 1);
			r.onerror = () => n(r.error), r.onupgradeneeded = () => {
				r.result.createObjectStore("reconciliationTime", { keyPath: "id" });
			}, r.onsuccess = () => {
				let i = r.result.transaction("reconciliationTime", "readwrite").objectStore("reconciliationTime").put({
					id: "cachedConfig",
					value: JSON.stringify(e)
				});
				i.onerror = () => n(i.error), i.onsuccess = () => t();
			};
		});
	}
	async getLastReconciliationTime() {
		return new Promise((e, t) => {
			let n = indexedDB.open("polyfeaDB", 1);
			n.onerror = () => t(n.error), n.onupgradeneeded = () => {
				n.result.createObjectStore("reconciliationTime", { keyPath: "id" });
			}, n.onsuccess = () => {
				let t = n.result.transaction("reconciliationTime", "readonly").objectStore("reconciliationTime").get("lastReconciliationTime");
				t.onerror = () => e(null), t.onsuccess = () => e(t.result?.value ?? null);
			};
		});
	}
	async setLastReconciliationTime(e) {
		return new Promise((t, n) => {
			let r = indexedDB.open("polyfeaDB", 1);
			r.onerror = () => n(r.error), r.onupgradeneeded = () => {
				r.result.createObjectStore("reconciliationTime");
			}, r.onsuccess = () => {
				let i = r.result.transaction("reconciliationTime", "readwrite").objectStore("reconciliationTime").put({
					id: "lastReconciliationTime",
					value: e
				});
				i.onerror = () => n(i.error), i.onsuccess = () => t();
			};
		});
	}
	install(e) {
		e.waitUntil((async () => {
			$.debug("Installing"), await this.reconcileRoutes(!0), await this.precacheController.install(e), this.scope.skipWaiting();
		})());
	}
	activate(e) {
		e.waitUntil((async () => {
			$.debug("Activating"), this.precacheController.activate(e), await this.reconcileRoutes();
			for (let e of this.interceptors) if (e.activate) try {
				await e.activate(), $.debug(`Interceptor ${e.name} activated successfully`);
			} catch (t) {
				$.error({ error: t }, `Interceptor ${e.name} failed to activate`);
			}
			await this.scope.clients.claim();
		})());
	}
	handleFetch(e) {
		let { request: t } = e, n = $.child({ request: t });
		setTimeout(() => this.reconcileRoutes().catch((e) => n.warn({ err: e }, "Failed to reconcile routes during fetch")), 0);
		let r = this.precacheController.getCacheKeyForURL(t.url);
		if (r) {
			n.debug(`Responded from precache: ${t.url}`), e.respondWith(caches.match(r).then((e) => e ?? fetch(t)));
			return;
		}
		let i = this.tryInterceptors(e);
		if (i) {
			n.debug(`Responded from interceptor: ${t.url}`), e.respondWith(i);
			return;
		}
		if (i = this.router.handleRequest({
			event: e,
			request: t
		}), i) {
			n.debug(`Responded from router: ${t.url}`), e.respondWith(i);
			return;
		}
		n.debug(`Route not found in SW, letting network handle it: ${t.url}`);
	}
	tryInterceptors(e) {
		for (let t of this.interceptors) try {
			let n = t.intercept(e.request, e);
			if (n) return n;
		} catch (n) {
			$.error({ error: n }, `Interceptor ${t.name} failed for ${e.request.url}`);
		}
	}
}().start();
//#endregion

//# sourceMappingURL=sw.mjs.map