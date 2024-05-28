try {
  self["workbox:window:7.0.0"] && _();
} catch {
}
function j(e, r) {
  return new Promise(function(t) {
    var i = new MessageChannel();
    i.port1.onmessage = function(s) {
      t(s.data);
    }, e.postMessage(r, [i.port2]);
  });
}
function L(e) {
  var r = function(t, i) {
    if (typeof t != "object" || !t)
      return t;
    var s = t[Symbol.toPrimitive];
    if (s !== void 0) {
      var h = s.call(t, i);
      if (typeof h != "object")
        return h;
      throw new TypeError("@@toPrimitive must return a primitive value.");
    }
    return String(t);
  }(e, "string");
  return typeof r == "symbol" ? r : r + "";
}
function W(e, r) {
  for (var t = 0; t < r.length; t++) {
    var i = r[t];
    i.enumerable = i.enumerable || !1, i.configurable = !0, "value" in i && (i.writable = !0), Object.defineProperty(e, L(i.key), i);
  }
}
function P(e, r) {
  return P = Object.setPrototypeOf ? Object.setPrototypeOf.bind() : function(t, i) {
    return t.__proto__ = i, t;
  }, P(e, r);
}
function E(e, r) {
  (r == null || r > e.length) && (r = e.length);
  for (var t = 0, i = new Array(r); t < r; t++)
    i[t] = e[t];
  return i;
}
function k(e, r) {
  var t = typeof Symbol < "u" && e[Symbol.iterator] || e["@@iterator"];
  if (t)
    return (t = t.call(e)).next.bind(t);
  if (Array.isArray(e) || (t = function(s, h) {
    if (s) {
      if (typeof s == "string")
        return E(s, h);
      var l = Object.prototype.toString.call(s).slice(8, -1);
      return l === "Object" && s.constructor && (l = s.constructor.name), l === "Map" || l === "Set" ? Array.from(s) : l === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(l) ? E(s, h) : void 0;
    }
  }(e)) || r) {
    t && (e = t);
    var i = 0;
    return function() {
      return i >= e.length ? { done: !0 } : { done: !1, value: e[i++] };
    };
  }
  throw new TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`);
}
try {
  self["workbox:core:7.0.0"] && _();
} catch {
}
var y = function() {
  var e = this;
  this.promise = new Promise(function(r, t) {
    e.resolve = r, e.reject = t;
  });
};
function b(e, r) {
  var t = location.href;
  return new URL(e, t).href === new URL(r, t).href;
}
var m = function(e, r) {
  this.type = e, Object.assign(this, r);
};
function p(e, r, t) {
  return t ? r ? r(e) : e : (e && e.then || (e = Promise.resolve(e)), r ? e.then(r) : e);
}
function U() {
}
var O = { type: "SKIP_WAITING" };
function S(e, r) {
  return e && e.then ? e.then(U) : Promise.resolve();
}
var A = function(e) {
  function r(v, u) {
    var n, o;
    return u === void 0 && (u = {}), (n = e.call(this) || this).nn = {}, n.tn = 0, n.rn = new y(), n.en = new y(), n.on = new y(), n.un = 0, n.an = /* @__PURE__ */ new Set(), n.cn = function() {
      var c = n.fn, a = c.installing;
      n.tn > 0 || !b(a.scriptURL, n.sn.toString()) || performance.now() > n.un + 6e4 ? (n.vn = a, c.removeEventListener("updatefound", n.cn)) : (n.hn = a, n.an.add(a), n.rn.resolve(a)), ++n.tn, a.addEventListener("statechange", n.ln);
    }, n.ln = function(c) {
      var a = n.fn, f = c.target, d = f.state, g = f === n.vn, w = { sw: f, isExternal: g, originalEvent: c };
      !g && n.mn && (w.isUpdate = !0), n.dispatchEvent(new m(d, w)), d === "installed" ? n.wn = self.setTimeout(function() {
        d === "installed" && a.waiting === f && n.dispatchEvent(new m("waiting", w));
      }, 200) : d === "activating" && (clearTimeout(n.wn), g || n.en.resolve(f));
    }, n.yn = function(c) {
      var a = n.hn, f = a !== navigator.serviceWorker.controller;
      n.dispatchEvent(new m("controlling", { isExternal: f, originalEvent: c, sw: a, isUpdate: n.mn })), f || n.on.resolve(a);
    }, n.gn = (o = function(c) {
      var a = c.data, f = c.ports, d = c.source;
      return p(n.getSW(), function() {
        n.an.has(d) && n.dispatchEvent(new m("message", { data: a, originalEvent: c, ports: f, sw: d }));
      });
    }, function() {
      for (var c = [], a = 0; a < arguments.length; a++)
        c[a] = arguments[a];
      try {
        return Promise.resolve(o.apply(this, c));
      } catch (f) {
        return Promise.reject(f);
      }
    }), n.sn = v, n.nn = u, navigator.serviceWorker.addEventListener("message", n.gn), n;
  }
  var t, i;
  i = e, (t = r).prototype = Object.create(i.prototype), t.prototype.constructor = t, P(t, i);
  var s, h, l = r.prototype;
  return l.register = function(v) {
    var u = (v === void 0 ? {} : v).immediate, n = u !== void 0 && u;
    try {
      var o = this;
      return p(function(c, a) {
        var f = c();
        return f && f.then ? f.then(a) : a(f);
      }(function() {
        if (!n && document.readyState !== "complete")
          return S(new Promise(function(c) {
            return window.addEventListener("load", c);
          }));
      }, function() {
        return o.mn = !!navigator.serviceWorker.controller, o.dn = o.pn(), p(o.bn(), function(c) {
          o.fn = c, o.dn && (o.hn = o.dn, o.en.resolve(o.dn), o.on.resolve(o.dn), o.dn.addEventListener("statechange", o.ln, { once: !0 }));
          var a = o.fn.waiting;
          return a && b(a.scriptURL, o.sn.toString()) && (o.hn = a, Promise.resolve().then(function() {
            o.dispatchEvent(new m("waiting", { sw: a, wasWaitingBeforeRegister: !0 }));
          }).then(function() {
          })), o.hn && (o.rn.resolve(o.hn), o.an.add(o.hn)), o.fn.addEventListener("updatefound", o.cn), navigator.serviceWorker.addEventListener("controllerchange", o.yn), o.fn;
        });
      }));
    } catch (c) {
      return Promise.reject(c);
    }
  }, l.update = function() {
    try {
      return this.fn ? p(S(this.fn.update())) : p();
    } catch (v) {
      return Promise.reject(v);
    }
  }, l.getSW = function() {
    return this.hn !== void 0 ? Promise.resolve(this.hn) : this.rn.promise;
  }, l.messageSW = function(v) {
    try {
      return p(this.getSW(), function(u) {
        return j(u, v);
      });
    } catch (u) {
      return Promise.reject(u);
    }
  }, l.messageSkipWaiting = function() {
    this.fn && this.fn.waiting && j(this.fn.waiting, O);
  }, l.pn = function() {
    var v = navigator.serviceWorker.controller;
    return v && b(v.scriptURL, this.sn.toString()) ? v : void 0;
  }, l.bn = function() {
    try {
      var v = this;
      return p(function(u, n) {
        try {
          var o = u();
        } catch (c) {
          return n(c);
        }
        return o && o.then ? o.then(void 0, n) : o;
      }(function() {
        return p(navigator.serviceWorker.register(v.sn, v.nn), function(u) {
          return v.un = performance.now(), u;
        });
      }, function(u) {
        throw u;
      }));
    } catch (u) {
      return Promise.reject(u);
    }
  }, s = r, (h = [{ key: "active", get: function() {
    return this.en.promise;
  } }, { key: "controlling", get: function() {
    return this.on.promise;
  } }]) && W(s.prototype, h), Object.defineProperty(s, "prototype", { writable: !1 }), s;
}(function() {
  function e() {
    this.Pn = /* @__PURE__ */ new Map();
  }
  var r = e.prototype;
  return r.addEventListener = function(t, i) {
    this.jn(t).add(i);
  }, r.removeEventListener = function(t, i) {
    this.jn(t).delete(i);
  }, r.dispatchEvent = function(t) {
    t.target = this;
    for (var i, s = k(this.jn(t.type)); !(i = s()).done; )
      (0, i.value)(t);
  }, r.jn = function(t) {
    return this.Pn.has(t) || this.Pn.set(t, /* @__PURE__ */ new Set()), this.Pn.get(t);
  }, e;
}());
function R(e = "", r = 0) {
  var t, i;
  if ("serviceWorker" in navigator) {
    const s = new URL("./sw.mjs", document.baseURI);
    e || (e = ((t = document.querySelector('meta[name="polyfea-sw-caching-config"]')) == null ? void 0 : t.getAttribute("content")) || ""), e && s.searchParams.set("caching-config", encodeURIComponent(e)), r || (r = parseInt(((i = document.querySelector('meta[name="polyfea-sw-reconcile-interval"]')) == null ? void 0 : i.getAttribute("content")) || "0")), r && s.searchParams.set("reconcile-interval", r.toString()), s.searchParams.set("base-path", new URL(document.baseURI).pathname), new A(s.pathname + s.search).register();
  }
}
R();
export {
  R as registerServiceWorker
};
//# sourceMappingURL=register.mjs.map
