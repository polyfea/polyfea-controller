var Qe = function(t, e) {
  return Qe = Object.setPrototypeOf || { __proto__: [] } instanceof Array && function(n, r) {
    n.__proto__ = r;
  } || function(n, r) {
    for (var i in r)
      Object.prototype.hasOwnProperty.call(r, i) && (n[i] = r[i]);
  }, Qe(t, e);
};
function M(t, e) {
  if (typeof e != "function" && e !== null)
    throw new TypeError("Class extends value " + String(e) + " is not a constructor or null");
  Qe(t, e);
  function n() {
    this.constructor = t;
  }
  t.prototype = e === null ? Object.create(e) : (n.prototype = e.prototype, new n());
}
function Cr(t, e, n, r) {
  function i(s) {
    return s instanceof n ? s : new n(function(o) {
      o(s);
    });
  }
  return new (n || (n = Promise))(function(s, o) {
    function a(f) {
      try {
        u(r.next(f));
      } catch (h) {
        o(h);
      }
    }
    function l(f) {
      try {
        u(r.throw(f));
      } catch (h) {
        o(h);
      }
    }
    function u(f) {
      f.done ? s(f.value) : i(f.value).then(a, l);
    }
    u((r = r.apply(t, e || [])).next());
  });
}
function dn(t, e) {
  var n = { label: 0, sent: function() {
    if (s[0] & 1)
      throw s[1];
    return s[1];
  }, trys: [], ops: [] }, r, i, s, o;
  return o = { next: a(0), throw: a(1), return: a(2) }, typeof Symbol == "function" && (o[Symbol.iterator] = function() {
    return this;
  }), o;
  function a(u) {
    return function(f) {
      return l([u, f]);
    };
  }
  function l(u) {
    if (r)
      throw new TypeError("Generator is already executing.");
    for (; o && (o = 0, u[0] && (n = 0)), n; )
      try {
        if (r = 1, i && (s = u[0] & 2 ? i.return : u[0] ? i.throw || ((s = i.return) && s.call(i), 0) : i.next) && !(s = s.call(i, u[1])).done)
          return s;
        switch (i = 0, s && (u = [u[0] & 2, s.value]), u[0]) {
          case 0:
          case 1:
            s = u;
            break;
          case 4:
            return n.label++, { value: u[1], done: !1 };
          case 5:
            n.label++, i = u[1], u = [0];
            continue;
          case 7:
            u = n.ops.pop(), n.trys.pop();
            continue;
          default:
            if (s = n.trys, !(s = s.length > 0 && s[s.length - 1]) && (u[0] === 6 || u[0] === 2)) {
              n = 0;
              continue;
            }
            if (u[0] === 3 && (!s || u[1] > s[0] && u[1] < s[3])) {
              n.label = u[1];
              break;
            }
            if (u[0] === 6 && n.label < s[1]) {
              n.label = s[1], s = u;
              break;
            }
            if (s && n.label < s[2]) {
              n.label = s[2], n.ops.push(u);
              break;
            }
            s[2] && n.ops.pop(), n.trys.pop();
            continue;
        }
        u = e.call(t, n);
      } catch (f) {
        u = [6, f], i = 0;
      } finally {
        r = s = 0;
      }
    if (u[0] & 5)
      throw u[1];
    return { value: u[0] ? u[1] : void 0, done: !0 };
  }
}
function ee(t) {
  var e = typeof Symbol == "function" && Symbol.iterator, n = e && t[e], r = 0;
  if (n)
    return n.call(t);
  if (t && typeof t.length == "number")
    return {
      next: function() {
        return t && r >= t.length && (t = void 0), { value: t && t[r++], done: !t };
      }
    };
  throw new TypeError(e ? "Object is not iterable." : "Symbol.iterator is not defined.");
}
function H(t, e) {
  var n = typeof Symbol == "function" && t[Symbol.iterator];
  if (!n)
    return t;
  var r = n.call(t), i, s = [], o;
  try {
    for (; (e === void 0 || e-- > 0) && !(i = r.next()).done; )
      s.push(i.value);
  } catch (a) {
    o = { error: a };
  } finally {
    try {
      i && !i.done && (n = r.return) && n.call(r);
    } finally {
      if (o)
        throw o.error;
    }
  }
  return s;
}
function j(t, e, n) {
  if (n || arguments.length === 2)
    for (var r = 0, i = e.length, s; r < i; r++)
      (s || !(r in e)) && (s || (s = Array.prototype.slice.call(e, 0, r)), s[r] = e[r]);
  return t.concat(s || Array.prototype.slice.call(e));
}
function X(t) {
  return this instanceof X ? (this.v = t, this) : new X(t);
}
function Or(t, e, n) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var r = n.apply(t, e || []), i, s = [];
  return i = {}, o("next"), o("throw"), o("return"), i[Symbol.asyncIterator] = function() {
    return this;
  }, i;
  function o(c) {
    r[c] && (i[c] = function(d) {
      return new Promise(function(p, y) {
        s.push([c, d, p, y]) > 1 || a(c, d);
      });
    });
  }
  function a(c, d) {
    try {
      l(r[c](d));
    } catch (p) {
      h(s[0][3], p);
    }
  }
  function l(c) {
    c.value instanceof X ? Promise.resolve(c.value.v).then(u, f) : h(s[0][2], c);
  }
  function u(c) {
    a("next", c);
  }
  function f(c) {
    a("throw", c);
  }
  function h(c, d) {
    c(d), s.shift(), s.length && a(s[0][0], s[0][1]);
  }
}
function Lr(t) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var e = t[Symbol.asyncIterator], n;
  return e ? e.call(t) : (t = typeof ee == "function" ? ee(t) : t[Symbol.iterator](), n = {}, r("next"), r("throw"), r("return"), n[Symbol.asyncIterator] = function() {
    return this;
  }, n);
  function r(s) {
    n[s] = t[s] && function(o) {
      return new Promise(function(a, l) {
        o = t[s](o), i(a, l, o.done, o.value);
      });
    };
  }
  function i(s, o, a, l) {
    Promise.resolve(l).then(function(u) {
      s({ value: u, done: a });
    }, o);
  }
}
function b(t) {
  return typeof t == "function";
}
function ct(t) {
  var e = function(r) {
    Error.call(r), r.stack = new Error().stack;
  }, n = t(e);
  return n.prototype = Object.create(Error.prototype), n.prototype.constructor = n, n;
}
var je = ct(function(t) {
  return function(n) {
    t(this), this.message = n ? n.length + ` errors occurred during unsubscription:
` + n.map(function(r, i) {
      return i + 1 + ") " + r.toString();
    }).join(`
  `) : "", this.name = "UnsubscriptionError", this.errors = n;
  };
});
function ve(t, e) {
  if (t) {
    var n = t.indexOf(e);
    0 <= n && t.splice(n, 1);
  }
}
var oe = function() {
  function t(e) {
    this.initialTeardown = e, this.closed = !1, this._parentage = null, this._finalizers = null;
  }
  return t.prototype.unsubscribe = function() {
    var e, n, r, i, s;
    if (!this.closed) {
      this.closed = !0;
      var o = this._parentage;
      if (o)
        if (this._parentage = null, Array.isArray(o))
          try {
            for (var a = ee(o), l = a.next(); !l.done; l = a.next()) {
              var u = l.value;
              u.remove(this);
            }
          } catch (y) {
            e = { error: y };
          } finally {
            try {
              l && !l.done && (n = a.return) && n.call(a);
            } finally {
              if (e)
                throw e.error;
            }
          }
        else
          o.remove(this);
      var f = this.initialTeardown;
      if (b(f))
        try {
          f();
        } catch (y) {
          s = y instanceof je ? y.errors : [y];
        }
      var h = this._finalizers;
      if (h) {
        this._finalizers = null;
        try {
          for (var c = ee(h), d = c.next(); !d.done; d = c.next()) {
            var p = d.value;
            try {
              kt(p);
            } catch (y) {
              s = s ?? [], y instanceof je ? s = j(j([], H(s)), H(y.errors)) : s.push(y);
            }
          }
        } catch (y) {
          r = { error: y };
        } finally {
          try {
            d && !d.done && (i = c.return) && i.call(c);
          } finally {
            if (r)
              throw r.error;
          }
        }
      }
      if (s)
        throw new je(s);
    }
  }, t.prototype.add = function(e) {
    var n;
    if (e && e !== this)
      if (this.closed)
        kt(e);
      else {
        if (e instanceof t) {
          if (e.closed || e._hasParent(this))
            return;
          e._addParent(this);
        }
        (this._finalizers = (n = this._finalizers) !== null && n !== void 0 ? n : []).push(e);
      }
  }, t.prototype._hasParent = function(e) {
    var n = this._parentage;
    return n === e || Array.isArray(n) && n.includes(e);
  }, t.prototype._addParent = function(e) {
    var n = this._parentage;
    this._parentage = Array.isArray(n) ? (n.push(e), n) : n ? [n, e] : e;
  }, t.prototype._removeParent = function(e) {
    var n = this._parentage;
    n === e ? this._parentage = null : Array.isArray(n) && ve(n, e);
  }, t.prototype.remove = function(e) {
    var n = this._finalizers;
    n && ve(n, e), e instanceof t && e._removeParent(this);
  }, t.EMPTY = function() {
    var e = new t();
    return e.closed = !0, e;
  }(), t;
}(), pn = oe.EMPTY;
function yn(t) {
  return t instanceof oe || t && "closed" in t && b(t.remove) && b(t.add) && b(t.unsubscribe);
}
function kt(t) {
  b(t) ? t() : t.unsubscribe();
}
var mn = {
  onUnhandledError: null,
  onStoppedNotification: null,
  Promise: void 0,
  useDeprecatedSynchronousErrorHandling: !1,
  useDeprecatedNextContext: !1
}, gn = {
  setTimeout: function(t, e) {
    for (var n = [], r = 2; r < arguments.length; r++)
      n[r - 2] = arguments[r];
    return setTimeout.apply(void 0, j([t, e], H(n)));
  },
  clearTimeout: function(t) {
    var e = gn.delegate;
    return ((e == null ? void 0 : e.clearTimeout) || clearTimeout)(t);
  },
  delegate: void 0
};
function bn(t) {
  gn.setTimeout(function() {
    throw t;
  });
}
function Xe() {
}
function ge(t) {
  t();
}
var ut = function(t) {
  M(e, t);
  function e(n) {
    var r = t.call(this) || this;
    return r.isStopped = !1, n ? (r.destination = n, yn(n) && n.add(r)) : r.destination = Mr, r;
  }
  return e.create = function(n, r, i) {
    return new $e(n, r, i);
  }, e.prototype.next = function(n) {
    this.isStopped || this._next(n);
  }, e.prototype.error = function(n) {
    this.isStopped || (this.isStopped = !0, this._error(n));
  }, e.prototype.complete = function() {
    this.isStopped || (this.isStopped = !0, this._complete());
  }, e.prototype.unsubscribe = function() {
    this.closed || (this.isStopped = !0, t.prototype.unsubscribe.call(this), this.destination = null);
  }, e.prototype._next = function(n) {
    this.destination.next(n);
  }, e.prototype._error = function(n) {
    try {
      this.destination.error(n);
    } finally {
      this.unsubscribe();
    }
  }, e.prototype._complete = function() {
    try {
      this.destination.complete();
    } finally {
      this.unsubscribe();
    }
  }, e;
}(oe), Pr = Function.prototype.bind;
function ze(t, e) {
  return Pr.call(t, e);
}
var Ur = function() {
  function t(e) {
    this.partialObserver = e;
  }
  return t.prototype.next = function(e) {
    var n = this.partialObserver;
    if (n.next)
      try {
        n.next(e);
      } catch (r) {
        he(r);
      }
  }, t.prototype.error = function(e) {
    var n = this.partialObserver;
    if (n.error)
      try {
        n.error(e);
      } catch (r) {
        he(r);
      }
    else
      he(e);
  }, t.prototype.complete = function() {
    var e = this.partialObserver;
    if (e.complete)
      try {
        e.complete();
      } catch (n) {
        he(n);
      }
  }, t;
}(), $e = function(t) {
  M(e, t);
  function e(n, r, i) {
    var s = t.call(this) || this, o;
    if (b(n) || !n)
      o = {
        next: n ?? void 0,
        error: r ?? void 0,
        complete: i ?? void 0
      };
    else {
      var a;
      s && mn.useDeprecatedNextContext ? (a = Object.create(n), a.unsubscribe = function() {
        return s.unsubscribe();
      }, o = {
        next: n.next && ze(n.next, a),
        error: n.error && ze(n.error, a),
        complete: n.complete && ze(n.complete, a)
      }) : o = n;
    }
    return s.destination = new Ur(o), s;
  }
  return e;
}(ut);
function he(t) {
  bn(t);
}
function Fr(t) {
  throw t;
}
var Mr = {
  closed: !0,
  next: Xe,
  error: Fr,
  complete: Xe
}, ft = function() {
  return typeof Symbol == "function" && Symbol.observable || "@@observable";
}();
function ae(t) {
  return t;
}
function Nr(t) {
  return t.length === 0 ? ae : t.length === 1 ? t[0] : function(n) {
    return t.reduce(function(r, i) {
      return i(r);
    }, n);
  };
}
var w = function() {
  function t(e) {
    e && (this._subscribe = e);
  }
  return t.prototype.lift = function(e) {
    var n = new t();
    return n.source = this, n.operator = e, n;
  }, t.prototype.subscribe = function(e, n, r) {
    var i = this, s = Jr(e) ? e : new $e(e, n, r);
    return ge(function() {
      var o = i, a = o.operator, l = o.source;
      s.add(a ? a.call(s, l) : l ? i._subscribe(s) : i._trySubscribe(s));
    }), s;
  }, t.prototype._trySubscribe = function(e) {
    try {
      return this._subscribe(e);
    } catch (n) {
      e.error(n);
    }
  }, t.prototype.forEach = function(e, n) {
    var r = this;
    return n = Rt(n), new n(function(i, s) {
      var o = new $e({
        next: function(a) {
          try {
            e(a);
          } catch (l) {
            s(l), o.unsubscribe();
          }
        },
        error: s,
        complete: i
      });
      r.subscribe(o);
    });
  }, t.prototype._subscribe = function(e) {
    var n;
    return (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(e);
  }, t.prototype[ft] = function() {
    return this;
  }, t.prototype.pipe = function() {
    for (var e = [], n = 0; n < arguments.length; n++)
      e[n] = arguments[n];
    return Nr(e)(this);
  }, t.prototype.toPromise = function(e) {
    var n = this;
    return e = Rt(e), new e(function(r, i) {
      var s;
      n.subscribe(function(o) {
        return s = o;
      }, function(o) {
        return i(o);
      }, function() {
        return r(s);
      });
    });
  }, t.create = function(e) {
    return new t(e);
  }, t;
}();
function Rt(t) {
  var e;
  return (e = t ?? mn.Promise) !== null && e !== void 0 ? e : Promise;
}
function Dr(t) {
  return t && b(t.next) && b(t.error) && b(t.complete);
}
function Jr(t) {
  return t && t instanceof ut || Dr(t) && yn(t);
}
function Hr(t) {
  return b(t == null ? void 0 : t.lift);
}
function C(t) {
  return function(e) {
    if (Hr(e))
      return e.lift(function(n) {
        try {
          return t(n, this);
        } catch (r) {
          this.error(r);
        }
      });
    throw new TypeError("Unable to lift unknown Observable type");
  };
}
function k(t, e, n, r, i) {
  return new Br(t, e, n, r, i);
}
var Br = function(t) {
  M(e, t);
  function e(n, r, i, s, o, a) {
    var l = t.call(this, n) || this;
    return l.onFinalize = o, l.shouldUnsubscribe = a, l._next = r ? function(u) {
      try {
        r(u);
      } catch (f) {
        n.error(f);
      }
    } : t.prototype._next, l._error = s ? function(u) {
      try {
        s(u);
      } catch (f) {
        n.error(f);
      } finally {
        this.unsubscribe();
      }
    } : t.prototype._error, l._complete = i ? function() {
      try {
        i();
      } catch (u) {
        n.error(u);
      } finally {
        this.unsubscribe();
      }
    } : t.prototype._complete, l;
  }
  return e.prototype.unsubscribe = function() {
    var n;
    if (!this.shouldUnsubscribe || this.shouldUnsubscribe()) {
      var r = this.closed;
      t.prototype.unsubscribe.call(this), !r && ((n = this.onFinalize) === null || n === void 0 || n.call(this));
    }
  }, e;
}(ut), jr = ct(function(t) {
  return function() {
    t(this), this.name = "ObjectUnsubscribedError", this.message = "object unsubscribed";
  };
}), ht = function(t) {
  M(e, t);
  function e() {
    var n = t.call(this) || this;
    return n.closed = !1, n.currentObservers = null, n.observers = [], n.isStopped = !1, n.hasError = !1, n.thrownError = null, n;
  }
  return e.prototype.lift = function(n) {
    var r = new Ct(this, this);
    return r.operator = n, r;
  }, e.prototype._throwIfClosed = function() {
    if (this.closed)
      throw new jr();
  }, e.prototype.next = function(n) {
    var r = this;
    ge(function() {
      var i, s;
      if (r._throwIfClosed(), !r.isStopped) {
        r.currentObservers || (r.currentObservers = Array.from(r.observers));
        try {
          for (var o = ee(r.currentObservers), a = o.next(); !a.done; a = o.next()) {
            var l = a.value;
            l.next(n);
          }
        } catch (u) {
          i = { error: u };
        } finally {
          try {
            a && !a.done && (s = o.return) && s.call(o);
          } finally {
            if (i)
              throw i.error;
          }
        }
      }
    });
  }, e.prototype.error = function(n) {
    var r = this;
    ge(function() {
      if (r._throwIfClosed(), !r.isStopped) {
        r.hasError = r.isStopped = !0, r.thrownError = n;
        for (var i = r.observers; i.length; )
          i.shift().error(n);
      }
    });
  }, e.prototype.complete = function() {
    var n = this;
    ge(function() {
      if (n._throwIfClosed(), !n.isStopped) {
        n.isStopped = !0;
        for (var r = n.observers; r.length; )
          r.shift().complete();
      }
    });
  }, e.prototype.unsubscribe = function() {
    this.isStopped = this.closed = !0, this.observers = this.currentObservers = null;
  }, Object.defineProperty(e.prototype, "observed", {
    get: function() {
      var n;
      return ((n = this.observers) === null || n === void 0 ? void 0 : n.length) > 0;
    },
    enumerable: !1,
    configurable: !0
  }), e.prototype._trySubscribe = function(n) {
    return this._throwIfClosed(), t.prototype._trySubscribe.call(this, n);
  }, e.prototype._subscribe = function(n) {
    return this._throwIfClosed(), this._checkFinalizedStatuses(n), this._innerSubscribe(n);
  }, e.prototype._innerSubscribe = function(n) {
    var r = this, i = this, s = i.hasError, o = i.isStopped, a = i.observers;
    return s || o ? pn : (this.currentObservers = null, a.push(n), new oe(function() {
      r.currentObservers = null, ve(a, n);
    }));
  }, e.prototype._checkFinalizedStatuses = function(n) {
    var r = this, i = r.hasError, s = r.thrownError, o = r.isStopped;
    i ? n.error(s) : o && n.complete();
  }, e.prototype.asObservable = function() {
    var n = new w();
    return n.source = this, n;
  }, e.create = function(n, r) {
    return new Ct(n, r);
  }, e;
}(w), Ct = function(t) {
  M(e, t);
  function e(n, r) {
    var i = t.call(this) || this;
    return i.destination = n, i.source = r, i;
  }
  return e.prototype.next = function(n) {
    var r, i;
    (i = (r = this.destination) === null || r === void 0 ? void 0 : r.next) === null || i === void 0 || i.call(r, n);
  }, e.prototype.error = function(n) {
    var r, i;
    (i = (r = this.destination) === null || r === void 0 ? void 0 : r.error) === null || i === void 0 || i.call(r, n);
  }, e.prototype.complete = function() {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.complete) === null || r === void 0 || r.call(n);
  }, e.prototype._subscribe = function(n) {
    var r, i;
    return (i = (r = this.source) === null || r === void 0 ? void 0 : r.subscribe(n)) !== null && i !== void 0 ? i : pn;
  }, e;
}(ht), zr = {
  now: function() {
    return Date.now();
  },
  delegate: void 0
}, de = function(t) {
  M(e, t);
  function e() {
    var n = t !== null && t.apply(this, arguments) || this;
    return n._value = null, n._hasValue = !1, n._isComplete = !1, n;
  }
  return e.prototype._checkFinalizedStatuses = function(n) {
    var r = this, i = r.hasError, s = r._hasValue, o = r._value, a = r.thrownError, l = r.isStopped, u = r._isComplete;
    i ? n.error(a) : (l || u) && (s && n.next(o), n.complete());
  }, e.prototype.next = function(n) {
    this.isStopped || (this._value = n, this._hasValue = !0);
  }, e.prototype.complete = function() {
    var n = this, r = n._hasValue, i = n._value, s = n._isComplete;
    s || (this._isComplete = !0, r && t.prototype.next.call(this, i), t.prototype.complete.call(this));
  }, e;
}(ht), Wr = function(t) {
  M(e, t);
  function e(n, r) {
    return t.call(this) || this;
  }
  return e.prototype.schedule = function(n, r) {
    return this;
  }, e;
}(oe), Ze = {
  setInterval: function(t, e) {
    for (var n = [], r = 2; r < arguments.length; r++)
      n[r - 2] = arguments[r];
    return setInterval.apply(void 0, j([t, e], H(n)));
  },
  clearInterval: function(t) {
    var e = Ze.delegate;
    return ((e == null ? void 0 : e.clearInterval) || clearInterval)(t);
  },
  delegate: void 0
}, Kr = function(t) {
  M(e, t);
  function e(n, r) {
    var i = t.call(this, n, r) || this;
    return i.scheduler = n, i.work = r, i.pending = !1, i;
  }
  return e.prototype.schedule = function(n, r) {
    var i;
    if (r === void 0 && (r = 0), this.closed)
      return this;
    this.state = n;
    var s = this.id, o = this.scheduler;
    return s != null && (this.id = this.recycleAsyncId(o, s, r)), this.pending = !0, this.delay = r, this.id = (i = this.id) !== null && i !== void 0 ? i : this.requestAsyncId(o, this.id, r), this;
  }, e.prototype.requestAsyncId = function(n, r, i) {
    return i === void 0 && (i = 0), Ze.setInterval(n.flush.bind(n, this), i);
  }, e.prototype.recycleAsyncId = function(n, r, i) {
    if (i === void 0 && (i = 0), i != null && this.delay === i && this.pending === !1)
      return r;
    r != null && Ze.clearInterval(r);
  }, e.prototype.execute = function(n, r) {
    if (this.closed)
      return new Error("executing a cancelled action");
    this.pending = !1;
    var i = this._execute(n, r);
    if (i)
      return i;
    this.pending === !1 && this.id != null && (this.id = this.recycleAsyncId(this.scheduler, this.id, null));
  }, e.prototype._execute = function(n, r) {
    var i = !1, s;
    try {
      this.work(n);
    } catch (o) {
      i = !0, s = o || new Error("Scheduled action threw falsy error");
    }
    if (i)
      return this.unsubscribe(), s;
  }, e.prototype.unsubscribe = function() {
    if (!this.closed) {
      var n = this, r = n.id, i = n.scheduler, s = i.actions;
      this.work = this.state = this.scheduler = null, this.pending = !1, ve(s, this), r != null && (this.id = this.recycleAsyncId(i, r, null)), this.delay = null, t.prototype.unsubscribe.call(this);
    }
  }, e;
}(Wr), Ot = function() {
  function t(e, n) {
    n === void 0 && (n = t.now), this.schedulerActionCtor = e, this.now = n;
  }
  return t.prototype.schedule = function(e, n, r) {
    return n === void 0 && (n = 0), new this.schedulerActionCtor(this, e).schedule(r, n);
  }, t.now = zr.now, t;
}(), qr = function(t) {
  M(e, t);
  function e(n, r) {
    r === void 0 && (r = Ot.now);
    var i = t.call(this, n, r) || this;
    return i.actions = [], i._active = !1, i;
  }
  return e.prototype.flush = function(n) {
    var r = this.actions;
    if (this._active) {
      r.push(n);
      return;
    }
    var i;
    this._active = !0;
    do
      if (i = n.execute(n.state, n.delay))
        break;
    while (n = r.shift());
    if (this._active = !1, i) {
      for (; n = r.shift(); )
        n.unsubscribe();
      throw i;
    }
  }, e;
}(Ot), Gr = new qr(Kr), Yr = Gr;
function vn(t) {
  return t && b(t.schedule);
}
function Qr(t) {
  return t[t.length - 1];
}
function Le(t) {
  return vn(Qr(t)) ? t.pop() : void 0;
}
var dt = function(t) {
  return t && typeof t.length == "number" && typeof t != "function";
};
function $n(t) {
  return b(t == null ? void 0 : t.then);
}
function wn(t) {
  return b(t[ft]);
}
function Sn(t) {
  return Symbol.asyncIterator && b(t == null ? void 0 : t[Symbol.asyncIterator]);
}
function Tn(t) {
  return new TypeError("You provided " + (t !== null && typeof t == "object" ? "an invalid object" : "'" + t + "'") + " where a stream was expected. You can provide an Observable, Promise, ReadableStream, Array, AsyncIterable, or Iterable.");
}
function Xr() {
  return typeof Symbol != "function" || !Symbol.iterator ? "@@iterator" : Symbol.iterator;
}
var En = Xr();
function xn(t) {
  return b(t == null ? void 0 : t[En]);
}
function _n(t) {
  return Or(this, arguments, function() {
    var n, r, i, s;
    return dn(this, function(o) {
      switch (o.label) {
        case 0:
          n = t.getReader(), o.label = 1;
        case 1:
          o.trys.push([1, , 9, 10]), o.label = 2;
        case 2:
          return [4, X(n.read())];
        case 3:
          return r = o.sent(), i = r.value, s = r.done, s ? [4, X(void 0)] : [3, 5];
        case 4:
          return [2, o.sent()];
        case 5:
          return [4, X(i)];
        case 6:
          return [4, o.sent()];
        case 7:
          return o.sent(), [3, 2];
        case 8:
          return [3, 10];
        case 9:
          return n.releaseLock(), [7];
        case 10:
          return [2];
      }
    });
  });
}
function An(t) {
  return b(t == null ? void 0 : t.getReader);
}
function P(t) {
  if (t instanceof w)
    return t;
  if (t != null) {
    if (wn(t))
      return Zr(t);
    if (dt(t))
      return Vr(t);
    if ($n(t))
      return ei(t);
    if (Sn(t))
      return In(t);
    if (xn(t))
      return ti(t);
    if (An(t))
      return ni(t);
  }
  throw Tn(t);
}
function Zr(t) {
  return new w(function(e) {
    var n = t[ft]();
    if (b(n.subscribe))
      return n.subscribe(e);
    throw new TypeError("Provided object does not correctly implement Symbol.observable");
  });
}
function Vr(t) {
  return new w(function(e) {
    for (var n = 0; n < t.length && !e.closed; n++)
      e.next(t[n]);
    e.complete();
  });
}
function ei(t) {
  return new w(function(e) {
    t.then(function(n) {
      e.closed || (e.next(n), e.complete());
    }, function(n) {
      return e.error(n);
    }).then(null, bn);
  });
}
function ti(t) {
  return new w(function(e) {
    var n, r;
    try {
      for (var i = ee(t), s = i.next(); !s.done; s = i.next()) {
        var o = s.value;
        if (e.next(o), e.closed)
          return;
      }
    } catch (a) {
      n = { error: a };
    } finally {
      try {
        s && !s.done && (r = i.return) && r.call(i);
      } finally {
        if (n)
          throw n.error;
      }
    }
    e.complete();
  });
}
function In(t) {
  return new w(function(e) {
    ri(t, e).catch(function(n) {
      return e.error(n);
    });
  });
}
function ni(t) {
  return In(_n(t));
}
function ri(t, e) {
  var n, r, i, s;
  return Cr(this, void 0, void 0, function() {
    var o, a;
    return dn(this, function(l) {
      switch (l.label) {
        case 0:
          l.trys.push([0, 5, 6, 11]), n = Lr(t), l.label = 1;
        case 1:
          return [4, n.next()];
        case 2:
          if (r = l.sent(), !!r.done)
            return [3, 4];
          if (o = r.value, e.next(o), e.closed)
            return [2];
          l.label = 3;
        case 3:
          return [3, 1];
        case 4:
          return [3, 11];
        case 5:
          return a = l.sent(), i = { error: a }, [3, 11];
        case 6:
          return l.trys.push([6, , 9, 10]), r && !r.done && (s = n.return) ? [4, s.call(n)] : [3, 8];
        case 7:
          l.sent(), l.label = 8;
        case 8:
          return [3, 10];
        case 9:
          if (i)
            throw i.error;
          return [7];
        case 10:
          return [7];
        case 11:
          return e.complete(), [2];
      }
    });
  });
}
function D(t, e, n, r, i) {
  r === void 0 && (r = 0), i === void 0 && (i = !1);
  var s = e.schedule(function() {
    n(), i ? t.add(this.schedule(null, r)) : this.unsubscribe();
  }, r);
  if (t.add(s), !i)
    return s;
}
function kn(t, e) {
  return e === void 0 && (e = 0), C(function(n, r) {
    n.subscribe(k(r, function(i) {
      return D(r, t, function() {
        return r.next(i);
      }, e);
    }, function() {
      return D(r, t, function() {
        return r.complete();
      }, e);
    }, function(i) {
      return D(r, t, function() {
        return r.error(i);
      }, e);
    }));
  });
}
function Rn(t, e) {
  return e === void 0 && (e = 0), C(function(n, r) {
    r.add(t.schedule(function() {
      return n.subscribe(r);
    }, e));
  });
}
function ii(t, e) {
  return P(t).pipe(Rn(e), kn(e));
}
function si(t, e) {
  return P(t).pipe(Rn(e), kn(e));
}
function oi(t, e) {
  return new w(function(n) {
    var r = 0;
    return e.schedule(function() {
      r === t.length ? n.complete() : (n.next(t[r++]), n.closed || this.schedule());
    });
  });
}
function ai(t, e) {
  return new w(function(n) {
    var r;
    return D(n, e, function() {
      r = t[En](), D(n, e, function() {
        var i, s, o;
        try {
          i = r.next(), s = i.value, o = i.done;
        } catch (a) {
          n.error(a);
          return;
        }
        o ? n.complete() : n.next(s);
      }, 0, !0);
    }), function() {
      return b(r == null ? void 0 : r.return) && r.return();
    };
  });
}
function Cn(t, e) {
  if (!t)
    throw new Error("Iterable cannot be null");
  return new w(function(n) {
    D(n, e, function() {
      var r = t[Symbol.asyncIterator]();
      D(n, e, function() {
        r.next().then(function(i) {
          i.done ? n.complete() : n.next(i.value);
        });
      }, 0, !0);
    });
  });
}
function li(t, e) {
  return Cn(_n(t), e);
}
function ci(t, e) {
  if (t != null) {
    if (wn(t))
      return ii(t, e);
    if (dt(t))
      return oi(t, e);
    if ($n(t))
      return si(t, e);
    if (Sn(t))
      return Cn(t, e);
    if (xn(t))
      return ai(t, e);
    if (An(t))
      return li(t, e);
  }
  throw Tn(t);
}
function Pe(t, e) {
  return e ? ci(t, e) : P(t);
}
function Ve() {
  for (var t = [], e = 0; e < arguments.length; e++)
    t[e] = arguments[e];
  var n = Le(t);
  return Pe(t, n);
}
function ui(t, e) {
  var n = b(t) ? t : function() {
    return t;
  }, r = function(i) {
    return i.error(n());
  };
  return new w(e ? function(i) {
    return e.schedule(r, 0, i);
  } : r);
}
var fi = ct(function(t) {
  return function() {
    t(this), this.name = "EmptyError", this.message = "no elements in sequence";
  };
});
function et(t, e) {
  var n = typeof e == "object";
  return new Promise(function(r, i) {
    var s = new $e({
      next: function(o) {
        r(o), s.unsubscribe();
      },
      error: i,
      complete: function() {
        n ? r(e.defaultValue) : i(new fi());
      }
    });
    t.subscribe(s);
  });
}
function hi(t) {
  return t instanceof Date && !isNaN(t);
}
function pt(t, e) {
  return C(function(n, r) {
    var i = 0;
    n.subscribe(k(r, function(s) {
      r.next(t.call(e, s, i++));
    }));
  });
}
var di = Array.isArray;
function pi(t, e) {
  return di(e) ? t.apply(void 0, j([], H(e))) : t(e);
}
function yi(t) {
  return pt(function(e) {
    return pi(t, e);
  });
}
function mi(t, e, n, r, i, s, o, a) {
  var l = [], u = 0, f = 0, h = !1, c = function() {
    h && !l.length && !u && e.complete();
  }, d = function(y) {
    return u < r ? p(y) : l.push(y);
  }, p = function(y) {
    s && e.next(y), u++;
    var g = !1;
    P(n(y, f++)).subscribe(k(e, function(m) {
      i == null || i(m), s ? d(m) : e.next(m);
    }, function() {
      g = !0;
    }, void 0, function() {
      if (g)
        try {
          u--;
          for (var m = function() {
            var A = l.shift();
            o ? D(e, o, function() {
              return p(A);
            }) : p(A);
          }; l.length && u < r; )
            m();
          c();
        } catch (A) {
          e.error(A);
        }
    }));
  };
  return t.subscribe(k(e, d, function() {
    h = !0, c();
  })), function() {
    a == null || a();
  };
}
function yt(t, e, n) {
  return n === void 0 && (n = 1 / 0), b(e) ? yt(function(r, i) {
    return pt(function(s, o) {
      return e(r, s, i, o);
    })(P(t(r, i)));
  }, n) : (typeof e == "number" && (n = e), C(function(r, i) {
    return mi(r, i, t, n);
  }));
}
function gi(t) {
  return t === void 0 && (t = 1 / 0), yt(ae, t);
}
function On() {
  return gi(1);
}
function Lt() {
  for (var t = [], e = 0; e < arguments.length; e++)
    t[e] = arguments[e];
  return On()(Pe(t, Le(t)));
}
function bi(t) {
  return new w(function(e) {
    P(t()).subscribe(e);
  });
}
var vi = ["addListener", "removeListener"], $i = ["addEventListener", "removeEventListener"], wi = ["on", "off"];
function tt(t, e, n, r) {
  if (b(n) && (r = n, n = void 0), r)
    return tt(t, e, n).pipe(yi(r));
  var i = H(Ei(t) ? $i.map(function(a) {
    return function(l) {
      return t[a](e, l, n);
    };
  }) : Si(t) ? vi.map(Pt(t, e)) : Ti(t) ? wi.map(Pt(t, e)) : [], 2), s = i[0], o = i[1];
  if (!s && dt(t))
    return yt(function(a) {
      return tt(a, e, n);
    })(P(t));
  if (!s)
    throw new TypeError("Invalid event target");
  return new w(function(a) {
    var l = function() {
      for (var u = [], f = 0; f < arguments.length; f++)
        u[f] = arguments[f];
      return a.next(1 < u.length ? u : u[0]);
    };
    return s(l), function() {
      return o(l);
    };
  });
}
function Pt(t, e) {
  return function(n) {
    return function(r) {
      return t[n](e, r);
    };
  };
}
function Si(t) {
  return b(t.addListener) && b(t.removeListener);
}
function Ti(t) {
  return b(t.on) && b(t.off);
}
function Ei(t) {
  return b(t.addEventListener) && b(t.removeEventListener);
}
function Ln(t, e, n) {
  t === void 0 && (t = 0), n === void 0 && (n = Yr);
  var r = -1;
  return e != null && (vn(e) ? n = e : r = e), new w(function(i) {
    var s = hi(t) ? +t - n.now() : t;
    s < 0 && (s = 0);
    var o = 0;
    return n.schedule(function() {
      i.closed || (i.next(o++), 0 <= r ? this.schedule(void 0, r) : i.complete());
    }, s);
  });
}
var xi = new w(Xe);
function mt(t) {
  return C(function(e, n) {
    var r = null, i = !1, s;
    r = e.subscribe(k(n, void 0, void 0, function(o) {
      s = P(t(o, mt(t)(e))), r ? (r.unsubscribe(), r = null, s.subscribe(n)) : i = !0;
    })), i && (r.unsubscribe(), r = null, s.subscribe(n));
  });
}
function _i() {
  for (var t = [], e = 0; e < arguments.length; e++)
    t[e] = arguments[e];
  var n = Le(t);
  return C(function(r, i) {
    On()(Pe(j([r], H(t)), n)).subscribe(i);
  });
}
function Pn() {
  for (var t = [], e = 0; e < arguments.length; e++)
    t[e] = arguments[e];
  return _i.apply(void 0, j([], H(t)));
}
function Ut(t, e) {
  return e === void 0 && (e = ae), t = t ?? Ai, C(function(n, r) {
    var i, s = !0;
    n.subscribe(k(r, function(o) {
      var a = e(o);
      (s || !t(i, a)) && (s = !1, i = a, r.next(o));
    }));
  });
}
function Ai(t, e) {
  return t === e;
}
function Ii(t) {
  t === void 0 && (t = 1 / 0);
  var e;
  t && typeof t == "object" ? e = t : e = {
    count: t
  };
  var n = e.count, r = n === void 0 ? 1 / 0 : n, i = e.delay, s = e.resetOnSuccess, o = s === void 0 ? !1 : s;
  return r <= 0 ? ae : C(function(a, l) {
    var u = 0, f, h = function() {
      var c = !1;
      f = a.subscribe(k(l, function(d) {
        o && (u = 0), l.next(d);
      }, void 0, function(d) {
        if (u++ < r) {
          var p = function() {
            f ? (f.unsubscribe(), f = null, h()) : c = !0;
          };
          if (i != null) {
            var y = typeof i == "number" ? Ln(i) : P(i(d, u)), g = k(l, function() {
              g.unsubscribe(), p();
            }, function() {
              l.complete();
            });
            y.subscribe(g);
          } else
            p();
        } else
          l.error(d);
      })), c && (f.unsubscribe(), f = null, h());
    };
    h();
  });
}
function ki() {
  for (var t = [], e = 0; e < arguments.length; e++)
    t[e] = arguments[e];
  var n = Le(t);
  return C(function(r, i) {
    (n ? Lt(t, r, n) : Lt(t, r)).subscribe(i);
  });
}
function Un(t, e) {
  return C(function(n, r) {
    var i = null, s = 0, o = !1, a = function() {
      return o && !i && r.complete();
    };
    n.subscribe(k(r, function(l) {
      i == null || i.unsubscribe();
      var u = 0, f = s++;
      P(t(l, f)).subscribe(i = k(r, function(h) {
        return r.next(e ? e(l, h, f, u++) : h);
      }, function() {
        i = null, a();
      }));
    }, function() {
      o = !0, a();
    }));
  });
}
function Ri(t, e, n) {
  var r = b(t) || e || n ? { next: t, error: e, complete: n } : t;
  return r ? C(function(i, s) {
    var o;
    (o = r.subscribe) === null || o === void 0 || o.call(r);
    var a = !0;
    i.subscribe(k(s, function(l) {
      var u;
      (u = r.next) === null || u === void 0 || u.call(r, l), s.next(l);
    }, function() {
      var l;
      a = !1, (l = r.complete) === null || l === void 0 || l.call(r), s.complete();
    }, function(l) {
      var u;
      a = !1, (u = r.error) === null || u === void 0 || u.call(r, l), s.error(l);
    }, function() {
      var l, u;
      a && ((l = r.unsubscribe) === null || l === void 0 || l.call(r)), (u = r.finalize) === null || u === void 0 || u.call(r);
    }));
  }) : ae;
}
const Ci = "http://./polyfea".replace(/\/+$/, "");
let ne = class {
  constructor(e = {}) {
    this.configuration = e;
  }
  set config(e) {
    this.configuration = e;
  }
  get basePath() {
    return this.configuration.basePath != null ? this.configuration.basePath : Ci;
  }
  get fetchApi() {
    return this.configuration.fetchApi;
  }
  get middleware() {
    return this.configuration.middleware || [];
  }
  get queryParamsStringify() {
    return this.configuration.queryParamsStringify || Nn;
  }
  get username() {
    return this.configuration.username;
  }
  get password() {
    return this.configuration.password;
  }
  get apiKey() {
    const e = this.configuration.apiKey;
    if (e)
      return typeof e == "function" ? e : () => e;
  }
  get accessToken() {
    const e = this.configuration.accessToken;
    if (e)
      return typeof e == "function" ? e : async () => e;
  }
  get headers() {
    return this.configuration.headers;
  }
  get credentials() {
    return this.configuration.credentials;
  }
};
const Oi = new ne();
let Fn = class Mn {
  constructor(e = Oi) {
    this.configuration = e, this.fetchApi = async (n, r) => {
      let i = { url: n, init: r };
      for (const o of this.middleware)
        o.pre && (i = await o.pre(Object.assign({ fetch: this.fetchApi }, i)) || i);
      let s;
      try {
        s = await (this.configuration.fetchApi || fetch)(i.url, i.init);
      } catch (o) {
        for (const a of this.middleware)
          a.onError && (s = await a.onError({
            fetch: this.fetchApi,
            url: i.url,
            init: i.init,
            error: o,
            response: s ? s.clone() : void 0
          }) || s);
        if (s === void 0)
          throw o instanceof Error ? new Fi(o, "The request failed and the interceptors did not return an alternative response") : o;
      }
      for (const o of this.middleware)
        o.post && (s = await o.post({
          fetch: this.fetchApi,
          url: i.url,
          init: i.init,
          response: s.clone()
        }) || s);
      return s;
    }, this.middleware = e.middleware;
  }
  withMiddleware(...e) {
    const n = this.clone();
    return n.middleware = n.middleware.concat(...e), n;
  }
  withPreMiddleware(...e) {
    const n = e.map((r) => ({ pre: r }));
    return this.withMiddleware(...n);
  }
  withPostMiddleware(...e) {
    const n = e.map((r) => ({ post: r }));
    return this.withMiddleware(...n);
  }
  /**
   * Check if the given MIME is a JSON MIME.
   * JSON MIME examples:
   *   application/json
   *   application/json; charset=UTF8
   *   APPLICATION/JSON
   *   application/vnd.company+json
   * @param mime - MIME (Multipurpose Internet Mail Extensions)
   * @return True if the given MIME is JSON, false otherwise.
   */
  isJsonMime(e) {
    return e ? Mn.jsonRegex.test(e) : !1;
  }
  async request(e, n) {
    const { url: r, init: i } = await this.createFetchParams(e, n), s = await this.fetchApi(r, i);
    if (s && s.status >= 200 && s.status < 300)
      return s;
    throw new Ui(s, "Response returned an error code");
  }
  async createFetchParams(e, n) {
    let r = this.configuration.basePath + e.path;
    e.query !== void 0 && Object.keys(e.query).length !== 0 && (r += "?" + this.configuration.queryParamsStringify(e.query));
    const i = Object.assign({}, this.configuration.headers, e.headers);
    Object.keys(i).forEach((f) => i[f] === void 0 ? delete i[f] : {});
    const s = typeof n == "function" ? n : async () => n, o = {
      method: e.method,
      headers: i,
      body: e.body,
      credentials: this.configuration.credentials
    }, a = Object.assign(Object.assign({}, o), await s({
      init: o,
      context: e
    }));
    let l;
    Pi(a.body) || a.body instanceof URLSearchParams || Li(a.body) ? l = a.body : this.isJsonMime(i["Content-Type"]) ? l = JSON.stringify(a.body) : l = a.body;
    const u = Object.assign(Object.assign({}, a), { body: l });
    return { url: r, init: u };
  }
  /**
   * Create a shallow clone of `this` by constructing a new instance
   * and then shallow cloning data members.
   */
  clone() {
    const e = this.constructor, n = new e(this.configuration);
    return n.middleware = this.middleware.slice(), n;
  }
};
Fn.jsonRegex = new RegExp("^(:?application/json|[^;/ 	]+/[^;/ 	]+[+]json)[ 	]*(:?;.*)?$", "i");
function Li(t) {
  return typeof Blob < "u" && t instanceof Blob;
}
function Pi(t) {
  return typeof FormData < "u" && t instanceof FormData;
}
let Ui = class extends Error {
  constructor(e, n) {
    super(n), this.response = e, this.name = "ResponseError";
  }
}, Fi = class extends Error {
  constructor(e, n) {
    super(n), this.cause = e, this.name = "FetchError";
  }
}, Ft = class extends Error {
  constructor(e, n) {
    super(n), this.field = e, this.name = "RequiredError";
  }
};
function T(t, e) {
  const n = t[e];
  return n != null;
}
function Nn(t, e = "") {
  return Object.keys(t).map((n) => Dn(n, t[n], e)).filter((n) => n.length > 0).join("&");
}
function Dn(t, e, n = "") {
  const r = n + (n.length ? `[${t}]` : t);
  if (e instanceof Array) {
    const i = e.map((s) => encodeURIComponent(String(s))).join(`&${encodeURIComponent(r)}=`);
    return `${encodeURIComponent(r)}=${i}`;
  }
  if (e instanceof Set) {
    const i = Array.from(e);
    return Dn(t, i, n);
  }
  return e instanceof Date ? `${encodeURIComponent(r)}=${encodeURIComponent(e.toISOString())}` : e instanceof Object ? Nn(e, r) : `${encodeURIComponent(r)}=${encodeURIComponent(String(e))}`;
}
function Jn(t, e) {
  return Object.keys(t).reduce((n, r) => Object.assign(Object.assign({}, n), { [r]: e(t[r]) }), {});
}
let Mt = class {
  constructor(e, n = (r) => r) {
    this.raw = e, this.transformer = n;
  }
  async value() {
    return this.transformer(await this.raw.json());
  }
};
function Mi(t) {
  return Ni(t);
}
function Ni(t, e) {
  return t == null ? t : {
    microfrontend: T(t, "microfrontend") ? t.microfrontend : void 0,
    tagName: t.tagName,
    attributes: T(t, "attributes") ? t.attributes : void 0,
    style: T(t, "style") ? t.style : void 0
  };
}
function Di(t) {
  return Ji(t);
}
function Ji(t, e) {
  return t == null ? t : {
    kind: T(t, "kind") ? t.kind : void 0,
    href: T(t, "href") ? t.href : void 0,
    attributes: T(t, "attributes") ? t.attributes : void 0,
    waitOnLoad: T(t, "waitOnLoad") ? t.waitOnLoad : void 0
  };
}
function Hn(t) {
  return Hi(t);
}
function Hi(t, e) {
  return t == null ? t : {
    dependsOn: T(t, "dependsOn") ? t.dependsOn : void 0,
    module: T(t, "module") ? t.module : void 0,
    resources: T(t, "resources") ? t.resources.map(Di) : void 0
  };
}
function Bn(t) {
  return Bi(t);
}
function Bi(t, e) {
  return t == null ? t : {
    elements: t.elements.map(Mi),
    microfrontends: T(t, "microfrontends") ? Jn(t.microfrontends, Hn) : void 0
  };
}
function ji(t) {
  return zi(t);
}
function zi(t, e) {
  return t == null ? t : {
    name: t.name,
    path: T(t, "path") ? t.path : void 0,
    contextArea: T(t, "contextArea") ? Bn(t.contextArea) : void 0
  };
}
function Wi(t) {
  return Ki(t);
}
function Ki(t, e) {
  return t == null ? t : {
    contextAreas: T(t, "contextAreas") ? t.contextAreas.map(ji) : void 0,
    microfrontends: Jn(t.microfrontends, Hn)
  };
}
let we = class extends Fn {
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextAreaRaw(e, n) {
    if (e.name === null || e.name === void 0)
      throw new Ft("name", "Required parameter requestParameters.name was null or undefined when calling getContextArea.");
    if (e.path === null || e.path === void 0)
      throw new Ft("path", "Required parameter requestParameters.path was null or undefined when calling getContextArea.");
    const r = {};
    e.path !== void 0 && (r.path = e.path), e.take !== void 0 && (r.take = e.take);
    const i = {}, s = await this.request({
      path: "/context-area/{name}".replace("{name}", encodeURIComponent(String(e.name))),
      method: "GET",
      headers: i,
      query: r
    }, n);
    return new Mt(s, (o) => Bn(o));
  }
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextArea(e, n) {
    return await (await this.getContextAreaRaw(e, n)).value();
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfigRaw(e) {
    const n = {}, r = {}, i = await this.request({
      path: "/static-config",
      method: "GET",
      headers: r,
      query: n
    }, e);
    return new Mt(i, (s) => Wi(s));
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfig(e) {
    return await (await this.getStaticConfigRaw(e)).value();
  }
}, qi = class {
  constructor(e = "./polyfea") {
    this.spec$ = new w();
    let n;
    typeof e == "string" ? (e.length === 0 && (e = "./polyfea"), n = new we(new ne({ basePath: e }))) : e instanceof ne ? n = new we(e) : n = e, this.spec$ = Pe(n.getStaticConfig()).pipe(Pn(xi));
  }
  getContextArea(e) {
    let n = globalThis.location.pathname;
    if (globalThis.document.baseURI) {
      const i = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(i) && (n = "./" + n.substring(i.length));
    }
    return this.spec$.pipe(
      pt((r) => {
        for (let i of r.contextAreas)
          if (i.name === e && new RegExp(i.path).test(n))
            return { ...i.contextArea, microfrontends: { ...r.microfrontends, ...i.contextArea.microfrontends } };
        return null;
      })
    );
  }
}, Nt = class {
  constructor(e = "./polyfea") {
    var n, r;
    typeof e == "string" ? this.api = new we(new ne({
      basePath: new URL(
        e,
        new URL(
          ((n = globalThis.document) == null ? void 0 : n.baseURI) || "/",
          ((r = globalThis.location) == null ? void 0 : r.href) || "http://localhost"
        )
      ).href
    })) : e instanceof ne ? this.api = new we(e) : this.api = e;
  }
  getContextArea(e) {
    var s, o;
    let n = ((s = globalThis.location) == null ? void 0 : s.pathname) || "/";
    if ((o = globalThis.document) != null && o.baseURI) {
      const l = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(l) && (n = "./" + n.substring(l.length));
    }
    const r = localStorage.getItem(`polyfea-context[${e},${n}]`), i = bi(() => this.api.getContextAreaRaw({ name: e, path: n })).pipe(
      Un((a) => a.raw.ok ? a.value() : ui(() => new Error(a.raw.statusText))),
      Ri((a) => {
        a && localStorage.setItem(`polyfea-context[${e},${n}]`, JSON.stringify(a));
      })
    );
    if (r) {
      const a = JSON.parse(r);
      return Ve(a).pipe(
        Pn(i),
        mt((l) => (console.warn(`Failed to fetch context area ${e} from ${n}, using cached version as the last known value`, l), Ve(a)))
      );
    } else
      return i.pipe(
        Ii({ count: 3, delay: (a) => Ln((a + 1) * 2e3) })
      );
  }
}, Gi = class {
  /**
   * Constructs a new NavigationDestination instance.
   * @param url - The URL of the navigation destination.
   */
  constructor(e) {
    this.url = e;
  }
}, Yi = class extends Event {
  constructor(e, n) {
    super("navigate", { bubbles: !0, cancelable: !0 }), this.transition = e, this.interceptPromises = [], this.downloadRequest = null, this.formData = null, this.hashChange = !1, this.userInitiated = !1;
    let r = new URL(e.href, new URL(globalThis.document.baseURI, globalThis.location.href));
    const i = new URL(n.url, new URL(globalThis.document.baseURI, globalThis.location.href));
    this.canIntercept = i.protocol === r.protocol && i.host === r.host && i.port === r.port, this.destination = new Gi(r.href);
  }
  /** (@see https://developer.mozilla.org/en-US/docs/Web/API/NavigateEvent )
   *  this polyfill signals abort only on programatic navigation
  **/
  get signal() {
    return this.transition.abortController.signal;
  }
  /**
   * Prevents the browser from following the navigation request.
   * @see {@link https://developer.mozilla.org/en-US/docs/Web/API/NavigateEvent/preventDefault}
   */
  intercept(e) {
    e != null && e.handler && this.interceptPromises.push(e.handler(this));
  }
}, Dt = class extends Event {
  constructor(e, n) {
    super("currententrychange", { bubbles: !0, cancelable: !0 }), this.navigationType = e, this.from = n;
  }
}, Jt = class {
  constructor(e) {
    this.request = e;
  }
  get finished() {
    return et(this.request.finished);
  }
  get from() {
    return this.request.entry;
  }
  get type() {
    return this.request.mode;
  }
};
function Qi(t = !1) {
  return Zi.tryRegister(t);
}
let Xi = class {
  constructor(e) {
    this.commited$ = e.committed, this.finished$ = e.finished;
  }
  get commited() {
    return et(this.commited$);
  }
  get finished() {
    return et(this.finished$);
  }
}, Zi = class nt extends EventTarget {
  constructor() {
    super(), this.entriesList = [], this.idCounter = 0, this.transitionRequests = new ht(), this.currentTransition = null, this.currentEntryIndex = -1, this.pushstateDelay = 35, this.rawHistoryMethods = {
      pushState: globalThis.history.pushState,
      replaceState: globalThis.history.replaceState,
      go: globalThis.history.go,
      back: globalThis.history.back,
      forward: globalThis.history.forward
    }, this.transitionRequests.subscribe((e) => this.executeRequest(e));
  }
  entries() {
    return this.entriesList;
  }
  get currentEntry() {
    if (this.currentEntryIndex >= 0)
      return this.entriesList[this.currentEntryIndex];
  }
  get canGoBack() {
    return this.currentEntryIndex > 0;
  }
  get canGoForward() {
    return this.currentEntryIndex + 1 < this.entriesList.length;
  }
  get transition() {
    var e;
    return ((e = this.currentTransition) == null ? void 0 : e.transition) || null;
  }
  navigate(e, n) {
    let r = "push";
    return ((n == null ? void 0 : n.history) === "replace" || n != null && n.replace) && (r = "replace"), this.nextTransitionRequest(r, e, n);
  }
  back() {
    if (this.currentEntryIndex < 1)
      throw { name: "InvaliStateError", message: "Cannot go back from initial state" };
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { traverseTo: this.entriesList[this.currentEntryIndex - 1].key });
  }
  forward(e) {
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { info: e, traverseTo: this.entriesList[this.currentEntryIndex + 1].key });
  }
  reload(e) {
    return this.nextTransitionRequest("reload", this.currentEntry.url, { info: e });
  }
  traverseTo(e, n) {
    const r = this.entriesList.find((i) => i.key === e);
    if (!r)
      throw { name: "InvaliStateError", message: "Cannot traverse to unknown state" };
    return this.nextTransitionRequest("traverse", r.url, { info: n == null ? void 0 : n.info, traverseTo: e });
  }
  updateCurrentEntry(e) {
    this.entriesList[this.currentEntryIndex].setState(JSON.parse(JSON.stringify(e == null ? void 0 : e.state))), this.currentEntry.dispatchEvent(new Dt("replace", this.currentEntry));
  }
  nextTransitionRequest(e, n, r) {
    const i = `@${++this.idCounter}-navigation-polyfill-transition`, s = {
      mode: e,
      href: new URL(n, new URL(globalThis.document.baseURI, globalThis.location.href)).href,
      info: r == null ? void 0 : r.info,
      state: r == null ? void 0 : r.state,
      committed: new de(),
      finished: new de(),
      entry: new We(this, i, i, n.toString(), r == null ? void 0 : r.state),
      abortController: new AbortController(),
      traverseToKey: r == null ? void 0 : r.traverseTo,
      transition: null
    };
    return this.transitionRequests.next(s), new Xi(s);
  }
  async executeRequest(e) {
    this.currentTransition && (this.currentTransition.abortController.abort(), this.currentTransition.finished.error("aborted - new navigation started"), this.currentTransition.committed.closed || this.currentTransition.committed.error("aborted - new navigation started"), globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: Error("aborted - new navigation started") })), this.clearTransition(e)), e.transition = new Jt(e), this.currentTransition = e;
    try {
      await this.commit(e), this.currentEntry.dispatchEvent(new Dt(e.mode, e.transition.from)), await this.dispatchNavigation(e);
    } catch {
      e.finished.error("aborted"), e.committed.error("aborted");
    } finally {
      this.clearTransition(e);
    }
  }
  dispatchNavigation(e) {
    var r;
    const n = new Yi(e, this.currentEntry);
    if ((r = globalThis.navigation) != null && r.dispatchEvent(n))
      return n.interceptPromises.length > 0 ? Promise.all(n.interceptPromises.filter((i) => !!(i != null && i.then))).then(() => {
        globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(e), e.finished.next(), e.finished.complete();
      }).catch((i) => {
        globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: i })), this.clearTransition(e), e.finished.error(i);
      }) : (globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(e), e.finished.next(), e.finished.complete(), Promise.resolve());
  }
  clearTransition(e) {
    var n;
    ((n = this.currentTransition) == null ? void 0 : n.entry.id) === e.entry.id && (this.currentTransition = null);
  }
  commit(e) {
    switch (e.mode) {
      case "push":
        return this.commitPushTransition(e);
      case "replace":
        return this.commitReplaceTransition(e);
      case "reload":
        return this.commitReloadTransition(e);
      case "traverse":
        return this.commitTraverseTransition(e);
    }
  }
  async pushstateAsync(e, n = () => {
  }) {
    return new Promise((r, i) => {
      setTimeout(
        () => {
          n(e), e.committed.next(), e.committed.complete(), r();
        },
        this.pushstateDelay
      );
    });
  }
  commitPushTransition(e) {
    return this.rawHistoryMethods.pushState.apply(globalThis.history, [e.entry.cloneable, "", e.href]), this.pushstateAsync(e, (n) => {
      this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), n.entry];
    });
  }
  commitReplaceTransition(e) {
    return e.entry.key = this.currentEntry.key, this.entriesList[this.currentEntryIndex] = e.entry, this.rawHistoryMethods.replaceState.apply(globalThis.history, [e.entry.cloneable, "", e.href]), this.pushstateAsync(e);
  }
  commitTraverseTransition(e) {
    return new Promise(async (n, r) => {
      const i = this.entriesList.findIndex((o) => o.key === e.traverseToKey);
      i < 0 && r("target entry not found");
      const s = i - this.currentEntryIndex;
      this.rawHistoryMethods.go.apply(globalThis.history, [s]), await this.pushstateAsync(e, (o) => {
        const a = this.entriesList.findIndex((l) => l.key === o.traverseToKey);
        a < 0 && o.committed.error(new Error("target entry not found")), this.currentEntryIndex = a, o.committed.next(), o.committed.complete(), n();
      });
    });
  }
  commitReloadTransition(e) {
    return e.committed.next(), e.committed.complete(), e.finished.subscribe({
      next: () => globalThis.location.reload(),
      error: () => globalThis.location.reload()
    }), Promise.resolve();
  }
  static tryRegister(e = !1) {
    if (!globalThis.navigation) {
      const n = new nt();
      return n.doRegister(e), n;
    }
    return globalThis.navigation;
  }
  static unregister() {
    globalThis.navigation && globalThis.navigation instanceof nt && (globalThis.navigation.doUnregister && globalThis.navigation.doUnregister(), globalThis.navigation = void 0);
  }
  doRegister(e) {
    var n, r, i, s, o;
    if (!globalThis.navigation && !globalThis.navigation) {
      globalThis.navigation = this, this.entriesList = [new We(this, "initial", "initial", globalThis.location.href, void 0)], this.currentEntryIndex = 0;
      const a = ((n = globalThis.history) == null ? void 0 : n.pushState) || ((c, d, p) => {
      }), l = ((r = globalThis.history) == null ? void 0 : r.replaceState) || ((c, d, p) => {
      }), u = ((i = globalThis.history) == null ? void 0 : i.go) || ((c) => {
      }), f = ((s = globalThis.history) == null ? void 0 : s.back) || (() => {
      }), h = ((o = globalThis.history) == null ? void 0 : o.forward) || (() => {
      });
      this.doUnregister = () => {
        globalThis.history.pushState = a, globalThis.history.replaceState = l, globalThis.history.go = u, globalThis.history.back = f, globalThis.history.forward = h;
      }, e ? (this.rawHistoryMethods.pushState = () => {
        a.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.replaceState = () => {
        l.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.go = () => {
        u.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.back = () => {
        f.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.forward = () => {
        h.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }) : this.rawHistoryMethods = {
        pushState: a || ((c, d, p) => {
        }),
        replaceState: l || ((c, d, p) => {
        }),
        go: u || ((c) => {
        }),
        back: f || (() => {
        }),
        forward: h || (() => {
        })
      }, globalThis.history && (globalThis.history.pushState = (c, d, p) => this.navigate(p, { state: c, history: "push" }), globalThis.history.replaceState = (c, d, p) => this.navigate(p, { state: c, history: "replace" }), globalThis.history.go = (c) => this.traverseTo(this.entriesList[this.currentEntryIndex + c].key), globalThis.history.back = () => this.back(), globalThis.history.forward = () => this.forward()), globalThis.addEventListener(
        "popstate",
        (c) => {
          var m, A, W, te, N, fe, K;
          if (this.currentTransition && ((m = c.state) == null ? void 0 : m.id) === ((W = (A = this.currentTransition) == null ? void 0 : A.entry) == null ? void 0 : W.id))
            return;
          (te = this.currentTransition) == null || te.abortController.abort();
          const d = new de();
          d.complete();
          let p;
          if ((N = c.state) != null && N.key) {
            const x = this.entriesList.findIndex((Be) => {
              var q;
              return Be.key === ((q = c.state) == null ? void 0 : q.key);
            });
            x >= 0 && (this.currentEntryIndex = x, p = this.entriesList[x]);
          }
          if (!p) {
            let x = `@${++this.idCounter}-navigation-polyfill-popstate`;
            p = new We(this, x, x, globalThis.location.href, c.state), this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), p];
          }
          const y = new de(), g = {
            mode: "traverse",
            href: globalThis.location.href,
            info: void 0,
            state: ((fe = c.state) == null ? void 0 : fe.state) || c.state,
            committed: d,
            finished: y,
            entry: p,
            abortController: new AbortController(),
            traverseToKey: (K = c.state) == null ? void 0 : K.key,
            transition: null
          };
          g.transition = new Jt(g), this.currentTransition = g, this.dispatchNavigation(this.currentTransition);
        }
      );
    }
  }
}, We = class extends EventTarget {
  constructor(e, n, r, i, s) {
    super(), this.owner = e, this.id = n, this.key = r, this.url = i, this.state = s, this.url = new URL(i, new URL(globalThis.document.baseURI, globalThis.location.href)).href;
  }
  get index() {
    return this.owner.entriesList.findIndex((e) => e.id === this.id);
  }
  get sameDocument() {
    return !0;
  }
  // polyfill is lost between documents
  getState() {
    return this.state;
  }
  setState(e) {
    this.state = e;
  }
  get cloneable() {
    return {
      id: this.id,
      key: this.key,
      url: this.url,
      index: this.index,
      state: this.state
    };
  }
}, Vi = class {
  /** @internal @private */
  constructor() {
  }
  /** @static
   * 
   * Get or create a polyfea driver instance. If the instance is provided on the global context, it is returned. 
   * Otherwise, a new instance is created with the given configuration.
   * 
   * @param config - Configuration for the [`PolyfeaApi`](https://github.com/polyfea/browser-api/blob/main/docs/classes/PolyfeaApi.md).
   **/
  static getOrCreate(e) {
    return globalThis.polyfea ? globalThis.polyfea : new Ht(e);
  }
  /** @static
   * Initialize the polyfea driver in the global context. 
   * This method is typically invoked by the polyfea controller script `boot.ts`.
   * 
   * @remarks 
   * This method also initializes the Navigation polyfill if it's not already present.
   * It augments `window.customElements.define` to allow for duplicate registration of custom elements.
   * This is particularly useful when different microfrontends need to register the same dependencies.
   */
  static initialize() {
    globalThis.polyfea || Ht.install();
  }
}, Ht = class jn {
  constructor(e) {
    this.config = e, this.loadedResources = /* @__PURE__ */ new Set(), globalThis.navigation && globalThis.navigation.addEventListener("navigate", (n) => {
      n.canIntercept && n.destination.url.startsWith(document.baseURI) && n.intercept();
    });
  }
  getBackend() {
    var e;
    if (!this.backend) {
      let n = (e = document.querySelector('meta[name="polyfea.backend"]')) == null ? void 0 : e.getAttribute("content");
      if (n)
        if (n.startsWith("static://")) {
          const r = n.slice(9);
          this.backend = new qi(this.config || r);
        } else
          this.backend = new Nt(this.config || n);
      else
        this.backend = new Nt(this.config || "./polyfea");
    }
    return this.backend;
  }
  getContextArea(e) {
    return globalThis.navigation ? tt(globalThis.navigation, "navigatesuccess").pipe(
      ki(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })),
      Un(
        (n) => this.getBackend().getContextArea(e).pipe(
          mt((r) => (console.error(r), Ve({ elements: [], microfrontends: {} })))
        )
      ),
      Ut((n, r) => JSON.stringify(n) === JSON.stringify(r))
    ) : this.getBackend().getContextArea(e).pipe(
      Ut((n, r) => JSON.stringify(n) === JSON.stringify(r))
    );
  }
  loadMicrofrontend(e, n) {
    if (!n)
      return Promise.resolve();
    const r = [];
    return this.loadMicrofrontendRecursive(e, n, r);
  }
  async loadMicrofrontendRecursive(e, n, r) {
    if (r.includes(n))
      throw new Error("Circular dependency detected: " + r.join(" -> "));
    const i = e.microfrontends[n];
    if (!i)
      throw new Error("Microfrontend specification not found: " + n);
    r.push(n), i.dependsOn && await Promise.all(i.dependsOn.map((o) => this.loadMicrofrontendRecursive(e, o, r).catch(
      (a) => {
        console.error(`Failed to load microfrontend's ${n} dependency ${o}`, a);
      }
    )));
    let s = i.resources || [];
    i.module && (s = [...s, {
      kind: "script",
      href: i.module,
      attributes: {
        type: "module"
      },
      waitOnLoad: !0
    }]), await Promise.all(s.map((o) => {
      if (this.loadedResources.has(o.href))
        return Promise.resolve();
      switch (o.kind) {
        case "script":
          return this.loadScript(o);
        case "stylesheet":
          return this.loadStylesheet(o);
        case "link":
          return this.loadLink(o);
      }
    }));
  }
  loadScript(e) {
    return new Promise((n, r) => {
      var o;
      const i = document.createElement("script");
      i.src = e.href, i.setAttribute("async", ""), e.attributes && Object.entries(e.attributes).forEach(([a, l]) => {
        i.setAttribute(a, l);
      });
      const s = (o = document.querySelector('meta[name="csp-nonce"]')) == null ? void 0 : o.getAttribute("content");
      s && i.setAttribute("nonce", s), this.loadedResources.add(e.href), e.waitOnLoad ? (i.onload = () => {
        n();
      }, i.onerror = () => {
        this.loadedResources.delete(e.href), console.error(`Failed to load script ${e.href} while loading microfrontend resources, check the network tab for details`), n();
      }) : n(), document.head.appendChild(i);
    });
  }
  loadStylesheet(e) {
    return this.loadLink({
      ...e,
      attributes: { ...e.attributes, rel: "stylesheet" }
    });
  }
  loadLink(e) {
    var i;
    const n = document.createElement("link");
    n.href = e.href, n.setAttribute("async", "");
    const r = (i = document.querySelector('meta[name="csp-nonce"]')) == null ? void 0 : i.getAttribute("content");
    return r && n.setAttribute("nonce", r), e.attributes && Object.entries(e.attributes).forEach(([s, o]) => {
      n.setAttribute(s, o);
    }), new Promise((s, o) => {
      this.loadedResources.add(e.href), e.waitOnLoad ? (n.onload = () => {
        s();
      }, n.onerror = () => {
        this.loadedResources.delete(e.href), console.error(`Failed to load ${e.href} while loading microfrontend resources, check the network tab for details`), s();
      }) : s(), document.head.appendChild(n);
    });
  }
  static install() {
    var i;
    globalThis.polyfea || (globalThis.polyfea = new jn()), Qi();
    const e = (i = document.querySelector('meta[name="polyfea.duplicit-custom-elements"]')) == null ? void 0 : i.getAttribute("content");
    let n = "warn";
    e === "silent" ? n = "silent" : e === "error" ? n = "error" : e === "verbose" && (n = "verbose");
    function r(s) {
      if (s.overrider === "polyfea")
        return s;
      const o = function(...a) {
        if (this.get(a[0])) {
          if (n === "error")
            throw new Error(`Custom element '${a[0]}' is duplicately registered`);
          if (n === "warn")
            return console.warn(`Custom element '${a[0]}' is duplicately registered - ignoring the current attempt for registration`), !1;
        } else
          return n === "verbose" && console.log(`Custom element '${a[0]}' is registered`), s.apply(this, a);
      };
      return o.overrider = "polyfea", o;
    }
    customElements.define = r(customElements.define);
  }
};
const B = {
  allRenderFn: !1,
  cmpDidLoad: !0,
  cmpDidUnload: !1,
  cmpDidUpdate: !0,
  cmpDidRender: !0,
  cmpWillLoad: !0,
  cmpWillUpdate: !0,
  cmpWillRender: !0,
  connectedCallback: !0,
  disconnectedCallback: !0,
  element: !0,
  event: !0,
  hasRenderFn: !0,
  lifecycle: !0,
  hostListener: !0,
  hostListenerTargetWindow: !0,
  hostListenerTargetDocument: !0,
  hostListenerTargetBody: !0,
  hostListenerTargetParent: !1,
  hostListenerTarget: !0,
  member: !0,
  method: !0,
  mode: !0,
  observeAttribute: !0,
  prop: !0,
  propMutable: !0,
  reflect: !0,
  scoped: !0,
  shadowDom: !0,
  slot: !0,
  cssAnnotations: !0,
  state: !0,
  style: !0,
  formAssociated: !1,
  svg: !0,
  updatable: !0,
  vdomAttribute: !0,
  vdomXlink: !0,
  vdomClass: !0,
  vdomFunctional: !0,
  vdomKey: !0,
  vdomListener: !0,
  vdomRef: !0,
  vdomPropOrAttr: !0,
  vdomRender: !0,
  vdomStyle: !0,
  vdomText: !0,
  watchCallback: !0,
  taskQueue: !0,
  hotModuleReplacement: !1,
  isDebug: !1,
  isDev: !1,
  isTesting: !1,
  hydrateServerSide: !1,
  hydrateClientSide: !1,
  lifecycleDOMEvents: !1,
  lazyLoad: !1,
  profile: !1,
  slotRelocation: !0,
  // TODO(STENCIL-914): remove this option when `experimentalSlotFixes` is the default behavior
  appendChildSlotFix: !1,
  // TODO(STENCIL-914): remove this option when `experimentalSlotFixes` is the default behavior
  cloneNodeFix: !1,
  hydratedAttribute: !1,
  hydratedClass: !0,
  scriptDataOpts: !1,
  // TODO(STENCIL-914): remove this option when `experimentalSlotFixes` is the default behavior
  scopedSlotTextContentFix: !1,
  // TODO(STENCIL-854): Remove code related to legacy shadowDomShim field
  shadowDomShim: !1,
  // TODO(STENCIL-914): remove this option when `experimentalSlotFixes` is the default behavior
  slotChildNodesFix: !1,
  invisiblePrehydration: !0,
  propBoolean: !0,
  propNumber: !0,
  propString: !0,
  constructableCSS: !0,
  cmpShouldUpdate: !0,
  devTools: !1,
  shadowDelegatesFocus: !0,
  initializeNextTick: !1,
  asyncLoading: !1,
  asyncQueue: !1,
  transformTagName: !1,
  attachStyles: !0,
  // TODO(STENCIL-914): remove this option when `experimentalSlotFixes` is the default behavior
  experimentalSlotFixes: !1
};
let G, zn, Ue, Wn = !1, Se = !1, gt = !1, _ = !1, Bt = null, rt = !1;
const z = (t, e = "") => () => {
}, es = "slot-fb{display:contents}slot-fb[hidden]{display:none}", jt = "http://www.w3.org/1999/xlink", zt = {}, ts = "http://www.w3.org/2000/svg", ns = "http://www.w3.org/1999/xhtml", rs = (t) => t != null, bt = (t) => (t = typeof t, t === "object" || t === "function");
function is(t) {
  var e, n, r;
  return (r = (n = (e = t.head) === null || e === void 0 ? void 0 : e.querySelector('meta[name="csp-nonce"]')) === null || n === void 0 ? void 0 : n.getAttribute("content")) !== null && r !== void 0 ? r : void 0;
}
const Q = (t, e, ...n) => {
  let r = null, i = null, s = null, o = !1, a = !1;
  const l = [], u = (h) => {
    for (let c = 0; c < h.length; c++)
      r = h[c], Array.isArray(r) ? u(r) : r != null && typeof r != "boolean" && ((o = typeof t != "function" && !bt(r)) && (r = String(r)), o && a ? l[l.length - 1].$text$ += r : l.push(o ? Te(null, r) : r), a = o);
  };
  if (u(n), e) {
    e.key && (i = e.key), e.name && (s = e.name);
    {
      const h = e.className || e.class;
      h && (e.class = typeof h != "object" ? h : Object.keys(h).filter((c) => h[c]).join(" "));
    }
  }
  if (typeof t == "function")
    return t(e === null ? {} : e, l, os);
  const f = Te(t, null);
  return f.$attrs$ = e, l.length > 0 && (f.$children$ = l), f.$key$ = i, f.$name$ = s, f;
}, Te = (t, e) => {
  const n = {
    $flags$: 0,
    $tag$: t,
    $text$: e,
    $elm$: null,
    $children$: null
  };
  return n.$attrs$ = null, n.$key$ = null, n.$name$ = null, n;
}, Kn = {}, ss = (t) => t && t.$tag$ === Kn, os = {
  forEach: (t, e) => t.map(Wt).forEach(e),
  map: (t, e) => t.map(Wt).map(e).map(as)
}, Wt = (t) => ({
  vattrs: t.$attrs$,
  vchildren: t.$children$,
  vkey: t.$key$,
  vname: t.$name$,
  vtag: t.$tag$,
  vtext: t.$text$
}), as = (t) => {
  if (typeof t.vtag == "function") {
    const n = Object.assign({}, t.vattrs);
    return t.vkey && (n.key = t.vkey), t.vname && (n.name = t.vname), Q(t.vtag, n, ...t.vchildren || []);
  }
  const e = Te(t.vtag, t.vtext);
  return e.$attrs$ = t.vattrs, e.$children$ = t.vchildren, e.$key$ = t.vkey, e.$name$ = t.vname, e;
}, ls = (t) => Us.map((e) => e(t)).find((e) => !!e), cs = (t, e) => t != null && !bt(t) ? e & 4 ? t === "false" ? !1 : t === "" || !!t : e & 2 ? parseFloat(t) : e & 1 ? String(t) : t : t, Kt = /* @__PURE__ */ new WeakMap(), us = (t, e, n) => {
  let r = xe.get(t);
  Ds && n ? (r = r || new CSSStyleSheet(), typeof r == "string" ? r = e : r.replaceSync(e)) : r = e, xe.set(t, r);
}, fs = (t, e, n) => {
  var r;
  const i = qn(e, n), s = xe.get(i);
  if (t = t.nodeType === 11 ? t : I, s)
    if (typeof s == "string") {
      t = t.head || t;
      let o = Kt.get(t), a;
      if (o || Kt.set(t, o = /* @__PURE__ */ new Set()), !o.has(i)) {
        {
          a = I.createElement("style"), a.innerHTML = s;
          const l = (r = S.$nonce$) !== null && r !== void 0 ? r : is(I);
          l != null && a.setAttribute("nonce", l), t.insertBefore(a, t.querySelector("link"));
        }
        e.$flags$ & 4 && (a.innerHTML += es), o && o.add(i);
      }
    } else
      t.adoptedStyleSheets.includes(s) || (t.adoptedStyleSheets = [...t.adoptedStyleSheets, s]);
  return i;
}, hs = (t) => {
  const e = t.$cmpMeta$, n = t.$hostElement$, r = e.$flags$, i = z("attachStyles", e.$tagName$), s = fs(n.shadowRoot ? n.shadowRoot : n.getRootNode(), e, t.$modeName$);
  r & 10 && (n["s-sc"] = s, n.classList.add(s + "-h"), r & 2 && n.classList.add(s + "-s")), i();
}, qn = (t, e) => "sc-" + (e && t.$flags$ & 32 ? t.$tagName$ + "-" + e : t.$tagName$), qt = (t, e, n, r, i, s) => {
  if (n !== r) {
    let o = Zt(t, e), a = e.toLowerCase();
    if (e === "class") {
      const l = t.classList, u = Gt(n), f = Gt(r);
      l.remove(...u.filter((h) => h && !f.includes(h))), l.add(...f.filter((h) => h && !u.includes(h)));
    } else if (e === "style") {
      for (const l in n)
        (!r || r[l] == null) && (l.includes("-") ? t.style.removeProperty(l) : t.style[l] = "");
      for (const l in r)
        (!n || r[l] !== n[l]) && (l.includes("-") ? t.style.setProperty(l, r[l]) : t.style[l] = r[l]);
    } else if (e !== "key")
      if (e === "ref")
        r && r(t);
      else if (!t.__lookupSetter__(e) && e[0] === "o" && e[1] === "n") {
        if (e[2] === "-" ? e = e.slice(3) : Zt(Fe, a) ? e = a.slice(2) : e = a[2] + e.slice(3), n || r) {
          const l = e.endsWith(Gn);
          e = e.replace(ps, ""), n && S.rel(t, e, n, l), r && S.ael(t, e, r, l);
        }
      } else {
        const l = bt(r);
        if ((o || l && r !== null) && !i)
          try {
            if (t.tagName.includes("-"))
              t[e] = r;
            else {
              const f = r ?? "";
              e === "list" ? o = !1 : (n == null || t[e] != f) && (t[e] = f);
            }
          } catch {
          }
        let u = !1;
        a !== (a = a.replace(/^xlink\:?/, "")) && (e = a, u = !0), r == null || r === !1 ? (r !== !1 || t.getAttribute(e) === "") && (u ? t.removeAttributeNS(jt, e) : t.removeAttribute(e)) : (!o || s & 4 || i) && !l && (r = r === !0 ? "" : r, u ? t.setAttributeNS(jt, e, r) : t.setAttribute(e, r));
      }
  }
}, ds = /\s/, Gt = (t) => t ? t.split(ds) : [], Gn = "Capture", ps = new RegExp(Gn + "$"), Yn = (t, e, n, r) => {
  const i = e.$elm$.nodeType === 11 && e.$elm$.host ? e.$elm$.host : e.$elm$, s = t && t.$attrs$ || zt, o = e.$attrs$ || zt;
  for (r in s)
    r in o || qt(i, r, s[r], void 0, n, e.$flags$);
  for (r in o)
    qt(i, r, s[r], o[r], n, e.$flags$);
}, Ee = (t, e, n, r) => {
  var i;
  const s = e.$children$[n];
  let o = 0, a, l, u;
  if (Wn || (gt = !0, s.$tag$ === "slot" && (G && r.classList.add(G + "-s"), s.$flags$ |= s.$children$ ? (
    // slot element has fallback content
    2
  ) : (
    // slot element does not have fallback content
    1
  ))), s.$text$ !== null)
    a = s.$elm$ = I.createTextNode(s.$text$);
  else if (s.$flags$ & 1)
    a = s.$elm$ = I.createTextNode("");
  else {
    if (_ || (_ = s.$tag$ === "svg"), a = s.$elm$ = I.createElementNS(_ ? ts : ns, s.$flags$ & 2 ? "slot-fb" : s.$tag$), _ && s.$tag$ === "foreignObject" && (_ = !1), Yn(null, s, _), rs(G) && a["s-si"] !== G && a.classList.add(a["s-si"] = G), s.$children$)
      for (o = 0; o < s.$children$.length; ++o)
        l = Ee(t, s, o, a), l && a.appendChild(l);
    s.$tag$ === "svg" ? _ = !1 : a.tagName === "foreignObject" && (_ = !0);
  }
  return a["s-hn"] = Ue, s.$flags$ & 3 && (a["s-sr"] = !0, a["s-fs"] = (i = s.$attrs$) === null || i === void 0 ? void 0 : i.slot, a["s-cr"] = zn, a["s-sn"] = s.$name$ || "", u = t && t.$children$ && t.$children$[n], u && u.$tag$ === s.$tag$ && t.$elm$ && re(t.$elm$, !1)), a;
}, re = (t, e) => {
  var n;
  S.$flags$ |= 1;
  const r = t.childNodes;
  for (let i = r.length - 1; i >= 0; i--) {
    const s = r[i];
    s["s-hn"] !== Ue && s["s-ol"] && (Zn(s).insertBefore(s, vt(s)), s["s-ol"].remove(), s["s-ol"] = void 0, s["s-sh"] = void 0, s.nodeType === 1 && s.setAttribute("slot", (n = s["s-sn"]) !== null && n !== void 0 ? n : ""), gt = !0), e && re(s, e);
  }
  S.$flags$ &= -2;
}, Qn = (t, e, n, r, i, s) => {
  let o = t["s-cr"] && t["s-cr"].parentNode || t, a;
  for (o.shadowRoot && o.tagName === Ue && (o = o.shadowRoot); i <= s; ++i)
    r[i] && (a = Ee(null, n, i, t), a && (r[i].$elm$ = a, o.insertBefore(a, vt(e))));
}, Xn = (t, e, n) => {
  for (let r = e; r <= n; ++r) {
    const i = t[r];
    if (i) {
      const s = i.$elm$;
      tr(i), s && (Se = !0, s["s-ol"] ? s["s-ol"].remove() : re(s, !0), s.remove());
    }
  }
}, ys = (t, e, n, r, i = !1) => {
  let s = 0, o = 0, a = 0, l = 0, u = e.length - 1, f = e[0], h = e[u], c = r.length - 1, d = r[0], p = r[c], y, g;
  for (; s <= u && o <= c; )
    if (f == null)
      f = e[++s];
    else if (h == null)
      h = e[--u];
    else if (d == null)
      d = r[++o];
    else if (p == null)
      p = r[--c];
    else if (pe(f, d, i))
      Y(f, d, i), f = e[++s], d = r[++o];
    else if (pe(h, p, i))
      Y(h, p, i), h = e[--u], p = r[--c];
    else if (pe(f, p, i))
      (f.$tag$ === "slot" || p.$tag$ === "slot") && re(f.$elm$.parentNode, !1), Y(f, p, i), t.insertBefore(f.$elm$, h.$elm$.nextSibling), f = e[++s], p = r[--c];
    else if (pe(h, d, i))
      (f.$tag$ === "slot" || p.$tag$ === "slot") && re(h.$elm$.parentNode, !1), Y(h, d, i), t.insertBefore(h.$elm$, f.$elm$), h = e[--u], d = r[++o];
    else {
      for (a = -1, l = s; l <= u; ++l)
        if (e[l] && e[l].$key$ !== null && e[l].$key$ === d.$key$) {
          a = l;
          break;
        }
      a >= 0 ? (g = e[a], g.$tag$ !== d.$tag$ ? y = Ee(e && e[o], n, a, t) : (Y(g, d, i), e[a] = void 0, y = g.$elm$), d = r[++o]) : (y = Ee(e && e[o], n, o, t), d = r[++o]), y && Zn(f.$elm$).insertBefore(y, vt(f.$elm$));
    }
  s > u ? Qn(t, r[c + 1] == null ? null : r[c + 1].$elm$, n, r, o, c) : o > c && Xn(e, s, u);
}, pe = (t, e, n = !1) => t.$tag$ === e.$tag$ ? t.$tag$ === "slot" ? t.$name$ === e.$name$ : n ? !0 : t.$key$ === e.$key$ : !1, vt = (t) => t && t["s-ol"] || t, Zn = (t) => (t["s-ol"] ? t["s-ol"] : t).parentNode, Y = (t, e, n = !1) => {
  const r = e.$elm$ = t.$elm$, i = t.$children$, s = e.$children$, o = e.$tag$, a = e.$text$;
  let l;
  a === null ? (_ = o === "svg" ? !0 : o === "foreignObject" ? !1 : _, o === "slot" || Yn(t, e, _), i !== null && s !== null ? ys(r, i, e, s, n) : s !== null ? (t.$text$ !== null && (r.textContent = ""), Qn(r, null, e, s, 0, s.length - 1)) : i !== null && Xn(i, 0, i.length - 1), _ && o === "svg" && (_ = !1)) : (l = r["s-cr"]) ? l.parentNode.textContent = a : t.$text$ !== a && (r.data = a);
}, Vn = (t) => {
  const e = t.childNodes;
  for (const n of e)
    if (n.nodeType === 1) {
      if (n["s-sr"]) {
        const r = n["s-sn"];
        n.hidden = !1;
        for (const i of e)
          if (i !== n) {
            if (i["s-hn"] !== n["s-hn"] || r !== "") {
              if (i.nodeType === 1 && (r === i.getAttribute("slot") || r === i["s-sn"])) {
                n.hidden = !0;
                break;
              }
            } else if (i.nodeType === 1 || i.nodeType === 3 && i.textContent.trim() !== "") {
              n.hidden = !0;
              break;
            }
          }
      }
      Vn(n);
    }
}, F = [], er = (t) => {
  let e, n, r;
  for (const i of t.childNodes) {
    if (i["s-sr"] && (e = i["s-cr"]) && e.parentNode) {
      n = e.parentNode.childNodes;
      const s = i["s-sn"];
      for (r = n.length - 1; r >= 0; r--)
        if (e = n[r], !e["s-cn"] && !e["s-nr"] && e["s-hn"] !== i["s-hn"] && !B.experimentalSlotFixes)
          if (Yt(e, s)) {
            let o = F.find((a) => a.$nodeToRelocate$ === e);
            Se = !0, e["s-sn"] = e["s-sn"] || s, o ? (o.$nodeToRelocate$["s-sh"] = i["s-hn"], o.$slotRefNode$ = i) : (e["s-sh"] = i["s-hn"], F.push({
              $slotRefNode$: i,
              $nodeToRelocate$: e
            })), e["s-sr"] && F.map((a) => {
              Yt(a.$nodeToRelocate$, e["s-sn"]) && (o = F.find((l) => l.$nodeToRelocate$ === e), o && !a.$slotRefNode$ && (a.$slotRefNode$ = o.$slotRefNode$));
            });
          } else
            F.some((o) => o.$nodeToRelocate$ === e) || F.push({
              $nodeToRelocate$: e
            });
    }
    i.nodeType === 1 && er(i);
  }
}, Yt = (t, e) => t.nodeType === 1 ? t.getAttribute("slot") === null && e === "" || t.getAttribute("slot") === e : t["s-sn"] === e ? !0 : e === "", tr = (t) => {
  t.$attrs$ && t.$attrs$.ref && t.$attrs$.ref(null), t.$children$ && t.$children$.map(tr);
}, ms = (t, e, n = !1) => {
  var r, i, s, o;
  const a = t.$hostElement$, l = t.$cmpMeta$, u = t.$vnode$ || Te(null, null), f = ss(e) ? e : Q(null, null, e);
  if (Ue = a.tagName, l.$attrsToReflect$ && (f.$attrs$ = f.$attrs$ || {}, l.$attrsToReflect$.map(([h, c]) => f.$attrs$[c] = a[h])), n && f.$attrs$)
    for (const h of Object.keys(f.$attrs$))
      a.hasAttribute(h) && !["key", "ref", "style", "class"].includes(h) && (f.$attrs$[h] = a[h]);
  f.$tag$ = null, f.$flags$ |= 4, t.$vnode$ = f, f.$elm$ = u.$elm$ = a.shadowRoot || a, G = a["s-sc"], zn = a["s-cr"], Wn = (l.$flags$ & 1) !== 0, Se = !1, Y(u, f, n);
  {
    if (S.$flags$ |= 1, gt) {
      er(f.$elm$);
      for (const h of F) {
        const c = h.$nodeToRelocate$;
        if (!c["s-ol"]) {
          const d = I.createTextNode("");
          d["s-nr"] = c, c.parentNode.insertBefore(c["s-ol"] = d, c);
        }
      }
      for (const h of F) {
        const c = h.$nodeToRelocate$, d = h.$slotRefNode$;
        if (d) {
          const p = d.parentNode;
          let y = d.nextSibling;
          {
            let g = (r = c["s-ol"]) === null || r === void 0 ? void 0 : r.previousSibling;
            for (; g; ) {
              let m = (i = g["s-nr"]) !== null && i !== void 0 ? i : null;
              if (m && m["s-sn"] === c["s-sn"] && p === m.parentNode && (m = m.nextSibling, !m || !m["s-nr"])) {
                y = m;
                break;
              }
              g = g.previousSibling;
            }
          }
          (!y && p !== c.parentNode || c.nextSibling !== y) && c !== y && (!c["s-hn"] && c["s-ol"] && (c["s-hn"] = c["s-ol"].parentNode.nodeName), p.insertBefore(c, y), c.nodeType === 1 && (c.hidden = (s = c["s-ih"]) !== null && s !== void 0 ? s : !1));
        } else
          c.nodeType === 1 && (n && (c["s-ih"] = (o = c.hidden) !== null && o !== void 0 ? o : !1), c.hidden = !0);
      }
    }
    Se && Vn(f.$elm$), S.$flags$ &= -2, F.length = 0;
  }
}, gs = (t, e) => {
}, nr = (t, e) => (t.$flags$ |= 16, gs(t, t.$ancestorComponent$), Bs(() => bs(t, e))), bs = (t, e) => {
  const n = t.$hostElement$, r = z("scheduleUpdate", t.$cmpMeta$.$tagName$), i = n;
  let s;
  return e ? s = Z(i, "componentWillLoad") : s = Z(i, "componentWillUpdate"), s = Qt(s, () => Z(i, "componentWillRender")), r(), Qt(s, () => $s(t, i, e));
}, Qt = (t, e) => vs(t) ? t.then(e) : e(), vs = (t) => t instanceof Promise || t && t.then && typeof t.then == "function", $s = async (t, e, n) => {
  const r = t.$hostElement$, i = z("update", t.$cmpMeta$.$tagName$);
  r["s-rc"], n && hs(t);
  const s = z("render", t.$cmpMeta$.$tagName$);
  ws(t, e, r, n), s(), i(), Ss(t);
}, ws = (t, e, n, r) => {
  try {
    Bt = e, e = e.render && e.render(), t.$flags$ &= -17, t.$flags$ |= 2, (B.hasRenderFn || B.reflect) && (B.vdomRender || B.reflect) && (B.hydrateServerSide || ms(t, e, r));
  } catch (l) {
    ce(l, t.$hostElement$);
  }
  return Bt = null, null;
}, Ss = (t) => {
  const e = t.$cmpMeta$.$tagName$, n = t.$hostElement$, r = z("postUpdate", e), i = n;
  t.$ancestorComponent$, Z(i, "componentDidRender"), t.$flags$ & 64 ? (Z(i, "componentDidUpdate"), r()) : (t.$flags$ |= 64, Z(i, "componentDidLoad"), r());
}, Z = (t, e, n) => {
  if (t && t[e])
    try {
      return t[e](n);
    } catch (r) {
      ce(r);
    }
}, Ts = (t, e) => le(t).$instanceValues$.get(e), Es = (t, e, n, r) => {
  const i = le(t), s = t, o = i.$instanceValues$.get(e), a = i.$flags$, l = s;
  n = cs(n, r.$members$[e][0]);
  const u = Number.isNaN(o) && Number.isNaN(n);
  if (n !== o && !u) {
    i.$instanceValues$.set(e, n);
    {
      if (r.$watchers$ && a & 128) {
        const h = r.$watchers$[e];
        h && h.map((c) => {
          try {
            l[c](n, o, e);
          } catch (d) {
            ce(d, s);
          }
        });
      }
      if ((a & 18) === 2) {
        if (l.componentShouldUpdate && l.componentShouldUpdate(n, o, e) === !1)
          return;
        nr(i, !1);
      }
    }
  }
}, xs = (t, e, n) => {
  var r;
  const i = t.prototype;
  if (e.$members$) {
    t.watchers && (e.$watchers$ = t.watchers);
    const s = Object.entries(e.$members$);
    s.map(([o, [a]]) => {
      (a & 31 || a & 32) && Object.defineProperty(i, o, {
        get() {
          return Ts(this, o);
        },
        set(l) {
          Es(this, o, l, e);
        },
        configurable: !0,
        enumerable: !0
      });
    });
    {
      const o = /* @__PURE__ */ new Map();
      i.attributeChangedCallback = function(a, l, u) {
        S.jmp(() => {
          var f;
          const h = o.get(a);
          if (this.hasOwnProperty(h))
            u = this[h], delete this[h];
          else {
            if (i.hasOwnProperty(h) && typeof this[h] == "number" && this[h] == u)
              return;
            if (h == null) {
              const c = le(this), d = c == null ? void 0 : c.$flags$;
              if (d && !(d & 8) && d & 128 && u !== l) {
                const y = this, g = (f = e.$watchers$) === null || f === void 0 ? void 0 : f[a];
                g == null || g.forEach((m) => {
                  y[m] != null && y[m].call(y, u, l, a);
                });
              }
              return;
            }
          }
          this[h] = u === null && typeof this[h] == "boolean" ? !1 : u;
        });
      }, t.observedAttributes = Array.from(/* @__PURE__ */ new Set([
        ...Object.keys((r = e.$watchers$) !== null && r !== void 0 ? r : {}),
        ...s.filter(
          ([a, l]) => l[0] & 15
          /* MEMBER_FLAGS.HasAttribute */
        ).map(([a, l]) => {
          var u;
          const f = l[1] || a;
          return o.set(f, a), l[0] & 512 && ((u = e.$attrsToReflect$) === null || u === void 0 || u.push([a, f])), f;
        })
      ]));
    }
  }
  return t;
}, _s = async (t, e, n, r) => {
  let i;
  if (!(e.$flags$ & 32) && (e.$flags$ |= 32, i = t.constructor, customElements.whenDefined(n.$tagName$).then(() => e.$flags$ |= 128), i.style)) {
    let o = i.style;
    typeof o != "string" && (o = o[e.$modeName$ = ls(t)]);
    const a = qn(n, e.$modeName$);
    if (!xe.has(a)) {
      const l = z("registerStyles", n.$tagName$);
      us(a, o, !!(n.$flags$ & 1)), l();
    }
  }
  e.$ancestorComponent$, nr(e, !0);
}, Xt = (t) => {
}, As = (t) => {
  if (!(S.$flags$ & 1)) {
    const e = le(t), n = e.$cmpMeta$, r = z("connectedCallback", n.$tagName$);
    e.$flags$ & 1 ? (rr(t, e, n.$listeners$), e != null && e.$lazyInstance$ ? Xt(e.$lazyInstance$) : e != null && e.$onReadyPromise$ && e.$onReadyPromise$.then(() => Xt(e.$lazyInstance$))) : (e.$flags$ |= 1, // TODO(STENCIL-854): Remove code related to legacy shadowDomShim field
    n.$flags$ & 12 && Is(t), n.$members$ && Object.entries(n.$members$).map(([i, [s]]) => {
      if (s & 31 && t.hasOwnProperty(i)) {
        const o = t[i];
        delete t[i], t[i] = o;
      }
    }), _s(t, e, n)), r();
  }
}, Is = (t) => {
  const e = t["s-cr"] = I.createComment("");
  e["s-cn"] = !0, t.insertBefore(e, t.firstChild);
}, ks = async (t) => {
  if (!(S.$flags$ & 1)) {
    const e = le(t);
    e.$rmListeners$ && (e.$rmListeners$.map((n) => n()), e.$rmListeners$ = void 0);
  }
}, Rs = (t, e) => {
  const n = {
    $flags$: e[0],
    $tagName$: e[1]
  };
  n.$members$ = e[2], n.$listeners$ = e[3], n.$watchers$ = t.$watchers$, n.$attrsToReflect$ = [];
  const r = t.prototype.connectedCallback, i = t.prototype.disconnectedCallback;
  return Object.assign(t.prototype, {
    __registerHost() {
      Ps(this, n);
    },
    connectedCallback() {
      As(this), r && r.call(this);
    },
    disconnectedCallback() {
      ks(this), i && i.call(this);
    },
    __attachShadow() {
      this.attachShadow({
        mode: "open",
        delegatesFocus: !!(n.$flags$ & 16)
      });
    }
  }), t.is = n.$tagName$, xs(t, n);
}, rr = (t, e, n, r) => {
  n && n.map(([i, s, o]) => {
    const a = Os(t, i), l = Cs(e, o), u = Ls(i);
    S.ael(a, s, l, u), (e.$rmListeners$ = e.$rmListeners$ || []).push(() => S.rel(a, s, l, u));
  });
}, Cs = (t, e) => (n) => {
  try {
    B.lazyLoad || t.$hostElement$[e](n);
  } catch (r) {
    ce(r);
  }
}, Os = (t, e) => e & 4 ? I : e & 8 ? Fe : e & 16 ? I.body : t, Ls = (t) => Ms ? {
  passive: (t & 1) !== 0,
  capture: (t & 2) !== 0
} : (t & 2) !== 0, ir = /* @__PURE__ */ new WeakMap(), le = (t) => ir.get(t), Ps = (t, e) => {
  const n = {
    $flags$: 0,
    $hostElement$: t,
    $cmpMeta$: e,
    $instanceValues$: /* @__PURE__ */ new Map()
  };
  return rr(t, n, e.$listeners$), ir.set(t, n);
}, Zt = (t, e) => e in t, ce = (t, e) => (0, console.error)(t, e), xe = /* @__PURE__ */ new Map(), Us = [], Fe = typeof window < "u" ? window : {}, I = Fe.document || { head: {} }, Fs = Fe.HTMLElement || class {
}, S = {
  $flags$: 0,
  $resourcesUrl$: "",
  jmp: (t) => t(),
  raf: (t) => requestAnimationFrame(t),
  ael: (t, e, n, r) => t.addEventListener(e, n, r),
  rel: (t, e, n, r) => t.removeEventListener(e, n, r),
  ce: (t, e) => new CustomEvent(t, e)
}, Ms = /* @__PURE__ */ (() => {
  let t = !1;
  try {
    I.addEventListener("e", null, Object.defineProperty({}, "passive", {
      get() {
        t = !0;
      }
    }));
  } catch {
  }
  return t;
})(), Ns = (t) => Promise.resolve(t), Ds = /* @__PURE__ */ (() => {
  try {
    return new CSSStyleSheet(), typeof new CSSStyleSheet().replaceSync == "function";
  } catch {
  }
  return !1;
})(), Vt = [], sr = [], Js = (t, e) => (n) => {
  t.push(n), rt || (rt = !0, e && S.$flags$ & 4 ? Hs(it) : S.raf(it));
}, en = (t) => {
  for (let e = 0; e < t.length; e++)
    try {
      t[e](performance.now());
    } catch (n) {
      ce(n);
    }
  t.length = 0;
}, it = () => {
  en(Vt), en(sr), (rt = Vt.length > 0) && S.raf(it);
}, Hs = (t) => Ns().then(t), Bs = /* @__PURE__ */ Js(sr, !0);
/*!
 * Part of Polyfea microfrontends suite - https://github.com/polyfea
 */
function v(t) {
  return typeof t == "function";
}
function $t(t) {
  const n = t((r) => {
    Error.call(r), r.stack = new Error().stack;
  });
  return n.prototype = Object.create(Error.prototype), n.prototype.constructor = n, n;
}
const Ke = $t((t) => function(n) {
  t(this), this.message = n ? `${n.length} errors occurred during unsubscription:
${n.map((r, i) => `${i + 1}) ${r.toString()}`).join(`
  `)}` : "", this.name = "UnsubscriptionError", this.errors = n;
});
function _e(t, e) {
  if (t) {
    const n = t.indexOf(e);
    0 <= n && t.splice(n, 1);
  }
}
class L {
  constructor(e) {
    this.initialTeardown = e, this.closed = !1, this._parentage = null, this._finalizers = null;
  }
  unsubscribe() {
    let e;
    if (!this.closed) {
      this.closed = !0;
      const { _parentage: n } = this;
      if (n)
        if (this._parentage = null, Array.isArray(n))
          for (const s of n)
            s.remove(this);
        else
          n.remove(this);
      const { initialTeardown: r } = this;
      if (v(r))
        try {
          r();
        } catch (s) {
          e = s instanceof Ke ? s.errors : [s];
        }
      const { _finalizers: i } = this;
      if (i) {
        this._finalizers = null;
        for (const s of i)
          try {
            tn(s);
          } catch (o) {
            e = e ?? [], o instanceof Ke ? e = [...e, ...o.errors] : e.push(o);
          }
      }
      if (e)
        throw new Ke(e);
    }
  }
  add(e) {
    var n;
    if (e && e !== this)
      if (this.closed)
        tn(e);
      else {
        if (e instanceof L) {
          if (e.closed || e._hasParent(this))
            return;
          e._addParent(this);
        }
        (this._finalizers = (n = this._finalizers) !== null && n !== void 0 ? n : []).push(e);
      }
  }
  _hasParent(e) {
    const { _parentage: n } = this;
    return n === e || Array.isArray(n) && n.includes(e);
  }
  _addParent(e) {
    const { _parentage: n } = this;
    this._parentage = Array.isArray(n) ? (n.push(e), n) : n ? [n, e] : e;
  }
  _removeParent(e) {
    const { _parentage: n } = this;
    n === e ? this._parentage = null : Array.isArray(n) && _e(n, e);
  }
  remove(e) {
    const { _finalizers: n } = this;
    n && _e(n, e), e instanceof L && e._removeParent(this);
  }
}
L.EMPTY = (() => {
  const t = new L();
  return t.closed = !0, t;
})();
const or = L.EMPTY;
function ar(t) {
  return t instanceof L || t && "closed" in t && v(t.remove) && v(t.add) && v(t.unsubscribe);
}
function tn(t) {
  v(t) ? t() : t.unsubscribe();
}
const Me = {
  onUnhandledError: null,
  onStoppedNotification: null,
  Promise: void 0,
  useDeprecatedSynchronousErrorHandling: !1,
  useDeprecatedNextContext: !1
}, Ae = {
  setTimeout(t, e, ...n) {
    const { delegate: r } = Ae;
    return r != null && r.setTimeout ? r.setTimeout(t, e, ...n) : setTimeout(t, e, ...n);
  },
  clearTimeout(t) {
    const { delegate: e } = Ae;
    return ((e == null ? void 0 : e.clearTimeout) || clearTimeout)(t);
  },
  delegate: void 0
};
function lr(t) {
  Ae.setTimeout(() => {
    const { onUnhandledError: e } = Me;
    if (e)
      e(t);
    else
      throw t;
  });
}
function st() {
}
const js = wt("C", void 0, void 0);
function zs(t) {
  return wt("E", void 0, t);
}
function Ws(t) {
  return wt("N", t, void 0);
}
function wt(t, e, n) {
  return {
    kind: t,
    value: e,
    error: n
  };
}
function be(t) {
  t();
}
class St extends L {
  constructor(e) {
    super(), this.isStopped = !1, e ? (this.destination = e, ar(e) && e.add(this)) : this.destination = Ys;
  }
  static create(e, n, r) {
    return new Ie(e, n, r);
  }
  next(e) {
    this.isStopped ? Ge(Ws(e), this) : this._next(e);
  }
  error(e) {
    this.isStopped ? Ge(zs(e), this) : (this.isStopped = !0, this._error(e));
  }
  complete() {
    this.isStopped ? Ge(js, this) : (this.isStopped = !0, this._complete());
  }
  unsubscribe() {
    this.closed || (this.isStopped = !0, super.unsubscribe(), this.destination = null);
  }
  _next(e) {
    this.destination.next(e);
  }
  _error(e) {
    try {
      this.destination.error(e);
    } finally {
      this.unsubscribe();
    }
  }
  _complete() {
    try {
      this.destination.complete();
    } finally {
      this.unsubscribe();
    }
  }
}
const Ks = Function.prototype.bind;
function qe(t, e) {
  return Ks.call(t, e);
}
class qs {
  constructor(e) {
    this.partialObserver = e;
  }
  next(e) {
    const { partialObserver: n } = this;
    if (n.next)
      try {
        n.next(e);
      } catch (r) {
        ye(r);
      }
  }
  error(e) {
    const { partialObserver: n } = this;
    if (n.error)
      try {
        n.error(e);
      } catch (r) {
        ye(r);
      }
    else
      ye(e);
  }
  complete() {
    const { partialObserver: e } = this;
    if (e.complete)
      try {
        e.complete();
      } catch (n) {
        ye(n);
      }
  }
}
class Ie extends St {
  constructor(e, n, r) {
    super();
    let i;
    if (v(e) || !e)
      i = {
        next: e ?? void 0,
        error: n ?? void 0,
        complete: r ?? void 0
      };
    else {
      let s;
      this && Me.useDeprecatedNextContext ? (s = Object.create(e), s.unsubscribe = () => this.unsubscribe(), i = {
        next: e.next && qe(e.next, s),
        error: e.error && qe(e.error, s),
        complete: e.complete && qe(e.complete, s)
      }) : i = e;
    }
    this.destination = new qs(i);
  }
}
function ye(t) {
  lr(t);
}
function Gs(t) {
  throw t;
}
function Ge(t, e) {
  const { onStoppedNotification: n } = Me;
  n && Ae.setTimeout(() => n(t, e));
}
const Ys = {
  closed: !0,
  next: st,
  error: Gs,
  complete: st
}, Tt = typeof Symbol == "function" && Symbol.observable || "@@observable";
function ue(t) {
  return t;
}
function Qs(t) {
  return t.length === 0 ? ue : t.length === 1 ? t[0] : function(n) {
    return t.reduce((r, i) => i(r), n);
  };
}
class $ {
  constructor(e) {
    e && (this._subscribe = e);
  }
  lift(e) {
    const n = new $();
    return n.source = this, n.operator = e, n;
  }
  subscribe(e, n, r) {
    const i = Zs(e) ? e : new Ie(e, n, r);
    return be(() => {
      const { operator: s, source: o } = this;
      i.add(s ? s.call(i, o) : o ? this._subscribe(i) : this._trySubscribe(i));
    }), i;
  }
  _trySubscribe(e) {
    try {
      return this._subscribe(e);
    } catch (n) {
      e.error(n);
    }
  }
  forEach(e, n) {
    return n = nn(n), new n((r, i) => {
      const s = new Ie({
        next: (o) => {
          try {
            e(o);
          } catch (a) {
            i(a), s.unsubscribe();
          }
        },
        error: i,
        complete: r
      });
      this.subscribe(s);
    });
  }
  _subscribe(e) {
    var n;
    return (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(e);
  }
  [Tt]() {
    return this;
  }
  pipe(...e) {
    return Qs(e)(this);
  }
  toPromise(e) {
    return e = nn(e), new e((n, r) => {
      let i;
      this.subscribe((s) => i = s, (s) => r(s), () => n(i));
    });
  }
}
$.create = (t) => new $(t);
function nn(t) {
  var e;
  return (e = t ?? Me.Promise) !== null && e !== void 0 ? e : Promise;
}
function Xs(t) {
  return t && v(t.next) && v(t.error) && v(t.complete);
}
function Zs(t) {
  return t && t instanceof St || Xs(t) && ar(t);
}
function Vs(t) {
  return v(t == null ? void 0 : t.lift);
}
function O(t) {
  return (e) => {
    if (Vs(e))
      return e.lift(function(n) {
        try {
          return t(n, this);
        } catch (r) {
          this.error(r);
        }
      });
    throw new TypeError("Unable to lift unknown Observable type");
  };
}
function R(t, e, n, r, i) {
  return new eo(t, e, n, r, i);
}
class eo extends St {
  constructor(e, n, r, i, s, o) {
    super(e), this.onFinalize = s, this.shouldUnsubscribe = o, this._next = n ? function(a) {
      try {
        n(a);
      } catch (l) {
        e.error(l);
      }
    } : super._next, this._error = i ? function(a) {
      try {
        i(a);
      } catch (l) {
        e.error(l);
      } finally {
        this.unsubscribe();
      }
    } : super._error, this._complete = r ? function() {
      try {
        r();
      } catch (a) {
        e.error(a);
      } finally {
        this.unsubscribe();
      }
    } : super._complete;
  }
  unsubscribe() {
    var e;
    if (!this.shouldUnsubscribe || this.shouldUnsubscribe()) {
      const { closed: n } = this;
      super.unsubscribe(), !n && ((e = this.onFinalize) === null || e === void 0 || e.call(this));
    }
  }
}
const to = $t((t) => function() {
  t(this), this.name = "ObjectUnsubscribedError", this.message = "object unsubscribed";
});
class Ne extends $ {
  constructor() {
    super(), this.closed = !1, this.currentObservers = null, this.observers = [], this.isStopped = !1, this.hasError = !1, this.thrownError = null;
  }
  lift(e) {
    const n = new cr(this, this);
    return n.operator = e, n;
  }
  _throwIfClosed() {
    if (this.closed)
      throw new to();
  }
  next(e) {
    be(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.currentObservers || (this.currentObservers = Array.from(this.observers));
        for (const n of this.currentObservers)
          n.next(e);
      }
    });
  }
  error(e) {
    be(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.hasError = this.isStopped = !0, this.thrownError = e;
        const { observers: n } = this;
        for (; n.length; )
          n.shift().error(e);
      }
    });
  }
  complete() {
    be(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.isStopped = !0;
        const { observers: e } = this;
        for (; e.length; )
          e.shift().complete();
      }
    });
  }
  unsubscribe() {
    this.isStopped = this.closed = !0, this.observers = this.currentObservers = null;
  }
  get observed() {
    var e;
    return ((e = this.observers) === null || e === void 0 ? void 0 : e.length) > 0;
  }
  _trySubscribe(e) {
    return this._throwIfClosed(), super._trySubscribe(e);
  }
  _subscribe(e) {
    return this._throwIfClosed(), this._checkFinalizedStatuses(e), this._innerSubscribe(e);
  }
  _innerSubscribe(e) {
    const { hasError: n, isStopped: r, observers: i } = this;
    return n || r ? or : (this.currentObservers = null, i.push(e), new L(() => {
      this.currentObservers = null, _e(i, e);
    }));
  }
  _checkFinalizedStatuses(e) {
    const { hasError: n, thrownError: r, isStopped: i } = this;
    n ? e.error(r) : i && e.complete();
  }
  asObservable() {
    const e = new $();
    return e.source = this, e;
  }
}
Ne.create = (t, e) => new cr(t, e);
class cr extends Ne {
  constructor(e, n) {
    super(), this.destination = e, this.source = n;
  }
  next(e) {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.next) === null || r === void 0 || r.call(n, e);
  }
  error(e) {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.error) === null || r === void 0 || r.call(n, e);
  }
  complete() {
    var e, n;
    (n = (e = this.destination) === null || e === void 0 ? void 0 : e.complete) === null || n === void 0 || n.call(e);
  }
  _subscribe(e) {
    var n, r;
    return (r = (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(e)) !== null && r !== void 0 ? r : or;
  }
}
const no = {
  now() {
    return Date.now();
  },
  delegate: void 0
};
class me extends Ne {
  constructor() {
    super(...arguments), this._value = null, this._hasValue = !1, this._isComplete = !1;
  }
  _checkFinalizedStatuses(e) {
    const { hasError: n, _hasValue: r, _value: i, thrownError: s, isStopped: o, _isComplete: a } = this;
    n ? e.error(s) : (o || a) && (r && e.next(i), e.complete());
  }
  next(e) {
    this.isStopped || (this._value = e, this._hasValue = !0);
  }
  complete() {
    const { _hasValue: e, _value: n, _isComplete: r } = this;
    r || (this._isComplete = !0, e && super.next(n), super.complete());
  }
}
class ro extends L {
  constructor(e, n) {
    super();
  }
  schedule(e, n = 0) {
    return this;
  }
}
const ke = {
  setInterval(t, e, ...n) {
    const { delegate: r } = ke;
    return r != null && r.setInterval ? r.setInterval(t, e, ...n) : setInterval(t, e, ...n);
  },
  clearInterval(t) {
    const { delegate: e } = ke;
    return ((e == null ? void 0 : e.clearInterval) || clearInterval)(t);
  },
  delegate: void 0
};
class io extends ro {
  constructor(e, n) {
    super(e, n), this.scheduler = e, this.work = n, this.pending = !1;
  }
  schedule(e, n = 0) {
    var r;
    if (this.closed)
      return this;
    this.state = e;
    const i = this.id, s = this.scheduler;
    return i != null && (this.id = this.recycleAsyncId(s, i, n)), this.pending = !0, this.delay = n, this.id = (r = this.id) !== null && r !== void 0 ? r : this.requestAsyncId(s, this.id, n), this;
  }
  requestAsyncId(e, n, r = 0) {
    return ke.setInterval(e.flush.bind(e, this), r);
  }
  recycleAsyncId(e, n, r = 0) {
    if (r != null && this.delay === r && this.pending === !1)
      return n;
    n != null && ke.clearInterval(n);
  }
  execute(e, n) {
    if (this.closed)
      return new Error("executing a cancelled action");
    this.pending = !1;
    const r = this._execute(e, n);
    if (r)
      return r;
    this.pending === !1 && this.id != null && (this.id = this.recycleAsyncId(this.scheduler, this.id, null));
  }
  _execute(e, n) {
    let r = !1, i;
    try {
      this.work(e);
    } catch (s) {
      r = !0, i = s || new Error("Scheduled action threw falsy error");
    }
    if (r)
      return this.unsubscribe(), i;
  }
  unsubscribe() {
    if (!this.closed) {
      const { id: e, scheduler: n } = this, { actions: r } = n;
      this.work = this.state = this.scheduler = null, this.pending = !1, _e(r, this), e != null && (this.id = this.recycleAsyncId(n, e, null)), this.delay = null, super.unsubscribe();
    }
  }
}
class ie {
  constructor(e, n = ie.now) {
    this.schedulerActionCtor = e, this.now = n;
  }
  schedule(e, n = 0, r) {
    return new this.schedulerActionCtor(this, e).schedule(r, n);
  }
}
ie.now = no.now;
class so extends ie {
  constructor(e, n = ie.now) {
    super(e, n), this.actions = [], this._active = !1;
  }
  flush(e) {
    const { actions: n } = this;
    if (this._active) {
      n.push(e);
      return;
    }
    let r;
    this._active = !0;
    do
      if (r = e.execute(e.state, e.delay))
        break;
    while (e = n.shift());
    if (this._active = !1, r) {
      for (; e = n.shift(); )
        e.unsubscribe();
      throw r;
    }
  }
}
const oo = new so(io), ao = oo;
function ur(t) {
  return t && v(t.schedule);
}
function lo(t) {
  return t[t.length - 1];
}
function De(t) {
  return ur(lo(t)) ? t.pop() : void 0;
}
function co(t, e, n, r) {
  function i(s) {
    return s instanceof n ? s : new n(function(o) {
      o(s);
    });
  }
  return new (n || (n = Promise))(function(s, o) {
    function a(f) {
      try {
        u(r.next(f));
      } catch (h) {
        o(h);
      }
    }
    function l(f) {
      try {
        u(r.throw(f));
      } catch (h) {
        o(h);
      }
    }
    function u(f) {
      f.done ? s(f.value) : i(f.value).then(a, l);
    }
    u((r = r.apply(t, e || [])).next());
  });
}
function rn(t) {
  var e = typeof Symbol == "function" && Symbol.iterator, n = e && t[e], r = 0;
  if (n)
    return n.call(t);
  if (t && typeof t.length == "number")
    return {
      next: function() {
        return t && r >= t.length && (t = void 0), { value: t && t[r++], done: !t };
      }
    };
  throw new TypeError(e ? "Object is not iterable." : "Symbol.iterator is not defined.");
}
function V(t) {
  return this instanceof V ? (this.v = t, this) : new V(t);
}
function uo(t, e, n) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var r = n.apply(t, e || []), i, s = [];
  return i = {}, o("next"), o("throw"), o("return"), i[Symbol.asyncIterator] = function() {
    return this;
  }, i;
  function o(c) {
    r[c] && (i[c] = function(d) {
      return new Promise(function(p, y) {
        s.push([c, d, p, y]) > 1 || a(c, d);
      });
    });
  }
  function a(c, d) {
    try {
      l(r[c](d));
    } catch (p) {
      h(s[0][3], p);
    }
  }
  function l(c) {
    c.value instanceof V ? Promise.resolve(c.value.v).then(u, f) : h(s[0][2], c);
  }
  function u(c) {
    a("next", c);
  }
  function f(c) {
    a("throw", c);
  }
  function h(c, d) {
    c(d), s.shift(), s.length && a(s[0][0], s[0][1]);
  }
}
function fo(t) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var e = t[Symbol.asyncIterator], n;
  return e ? e.call(t) : (t = typeof rn == "function" ? rn(t) : t[Symbol.iterator](), n = {}, r("next"), r("throw"), r("return"), n[Symbol.asyncIterator] = function() {
    return this;
  }, n);
  function r(s) {
    n[s] = t[s] && function(o) {
      return new Promise(function(a, l) {
        o = t[s](o), i(a, l, o.done, o.value);
      });
    };
  }
  function i(s, o, a, l) {
    Promise.resolve(l).then(function(u) {
      s({ value: u, done: a });
    }, o);
  }
}
const Et = (t) => t && typeof t.length == "number" && typeof t != "function";
function fr(t) {
  return v(t == null ? void 0 : t.then);
}
function hr(t) {
  return v(t[Tt]);
}
function dr(t) {
  return Symbol.asyncIterator && v(t == null ? void 0 : t[Symbol.asyncIterator]);
}
function pr(t) {
  return new TypeError(`You provided ${t !== null && typeof t == "object" ? "an invalid object" : `'${t}'`} where a stream was expected. You can provide an Observable, Promise, ReadableStream, Array, AsyncIterable, or Iterable.`);
}
function ho() {
  return typeof Symbol != "function" || !Symbol.iterator ? "@@iterator" : Symbol.iterator;
}
const yr = ho();
function mr(t) {
  return v(t == null ? void 0 : t[yr]);
}
function gr(t) {
  return uo(this, arguments, function* () {
    const n = t.getReader();
    try {
      for (; ; ) {
        const { value: r, done: i } = yield V(n.read());
        if (i)
          return yield V(void 0);
        yield yield V(r);
      }
    } finally {
      n.releaseLock();
    }
  });
}
function br(t) {
  return v(t == null ? void 0 : t.getReader);
}
function U(t) {
  if (t instanceof $)
    return t;
  if (t != null) {
    if (hr(t))
      return po(t);
    if (Et(t))
      return yo(t);
    if (fr(t))
      return mo(t);
    if (dr(t))
      return vr(t);
    if (mr(t))
      return go(t);
    if (br(t))
      return bo(t);
  }
  throw pr(t);
}
function po(t) {
  return new $((e) => {
    const n = t[Tt]();
    if (v(n.subscribe))
      return n.subscribe(e);
    throw new TypeError("Provided object does not correctly implement Symbol.observable");
  });
}
function yo(t) {
  return new $((e) => {
    for (let n = 0; n < t.length && !e.closed; n++)
      e.next(t[n]);
    e.complete();
  });
}
function mo(t) {
  return new $((e) => {
    t.then((n) => {
      e.closed || (e.next(n), e.complete());
    }, (n) => e.error(n)).then(null, lr);
  });
}
function go(t) {
  return new $((e) => {
    for (const n of t)
      if (e.next(n), e.closed)
        return;
    e.complete();
  });
}
function vr(t) {
  return new $((e) => {
    vo(t, e).catch((n) => e.error(n));
  });
}
function bo(t) {
  return vr(gr(t));
}
function vo(t, e) {
  var n, r, i, s;
  return co(this, void 0, void 0, function* () {
    try {
      for (n = fo(t); r = yield n.next(), !r.done; ) {
        const o = r.value;
        if (e.next(o), e.closed)
          return;
      }
    } catch (o) {
      i = { error: o };
    } finally {
      try {
        r && !r.done && (s = n.return) && (yield s.call(n));
      } finally {
        if (i)
          throw i.error;
      }
    }
    e.complete();
  });
}
function J(t, e, n, r = 0, i = !1) {
  const s = e.schedule(function() {
    n(), i ? t.add(this.schedule(null, r)) : this.unsubscribe();
  }, r);
  if (t.add(s), !i)
    return s;
}
function $r(t, e = 0) {
  return O((n, r) => {
    n.subscribe(R(r, (i) => J(r, t, () => r.next(i), e), () => J(r, t, () => r.complete(), e), (i) => J(r, t, () => r.error(i), e)));
  });
}
function wr(t, e = 0) {
  return O((n, r) => {
    r.add(t.schedule(() => n.subscribe(r), e));
  });
}
function $o(t, e) {
  return U(t).pipe(wr(e), $r(e));
}
function wo(t, e) {
  return U(t).pipe(wr(e), $r(e));
}
function So(t, e) {
  return new $((n) => {
    let r = 0;
    return e.schedule(function() {
      r === t.length ? n.complete() : (n.next(t[r++]), n.closed || this.schedule());
    });
  });
}
function To(t, e) {
  return new $((n) => {
    let r;
    return J(n, e, () => {
      r = t[yr](), J(n, e, () => {
        let i, s;
        try {
          ({ value: i, done: s } = r.next());
        } catch (o) {
          n.error(o);
          return;
        }
        s ? n.complete() : n.next(i);
      }, 0, !0);
    }), () => v(r == null ? void 0 : r.return) && r.return();
  });
}
function Sr(t, e) {
  if (!t)
    throw new Error("Iterable cannot be null");
  return new $((n) => {
    J(n, e, () => {
      const r = t[Symbol.asyncIterator]();
      J(n, e, () => {
        r.next().then((i) => {
          i.done ? n.complete() : n.next(i.value);
        });
      }, 0, !0);
    });
  });
}
function Eo(t, e) {
  return Sr(gr(t), e);
}
function xo(t, e) {
  if (t != null) {
    if (hr(t))
      return $o(t, e);
    if (Et(t))
      return So(t, e);
    if (fr(t))
      return wo(t, e);
    if (dr(t))
      return Sr(t, e);
    if (mr(t))
      return To(t, e);
    if (br(t))
      return Eo(t, e);
  }
  throw pr(t);
}
function Je(t, e) {
  return e ? xo(t, e) : U(t);
}
function ot(...t) {
  const e = De(t);
  return Je(t, e);
}
function _o(t, e) {
  const n = v(t) ? t : () => t, r = (i) => i.error(n());
  return new $(e ? (i) => e.schedule(r, 0, i) : r);
}
const Ao = $t((t) => function() {
  t(this), this.name = "EmptyError", this.message = "no elements in sequence";
});
function at(t, e) {
  const n = typeof e == "object";
  return new Promise((r, i) => {
    const s = new Ie({
      next: (o) => {
        r(o), s.unsubscribe();
      },
      error: i,
      complete: () => {
        n ? r(e.defaultValue) : i(new Ao());
      }
    });
    t.subscribe(s);
  });
}
function Io(t) {
  return t instanceof Date && !isNaN(t);
}
function xt(t, e) {
  return O((n, r) => {
    let i = 0;
    n.subscribe(R(r, (s) => {
      r.next(t.call(e, s, i++));
    }));
  });
}
const { isArray: ko } = Array;
function Ro(t, e) {
  return ko(e) ? t(...e) : t(e);
}
function Co(t) {
  return xt((e) => Ro(t, e));
}
function Oo(t, e, n, r, i, s, o, a) {
  const l = [];
  let u = 0, f = 0, h = !1;
  const c = () => {
    h && !l.length && !u && e.complete();
  }, d = (y) => u < r ? p(y) : l.push(y), p = (y) => {
    s && e.next(y), u++;
    let g = !1;
    U(n(y, f++)).subscribe(R(e, (m) => {
      i == null || i(m), s ? d(m) : e.next(m);
    }, () => {
      g = !0;
    }, void 0, () => {
      if (g)
        try {
          for (u--; l.length && u < r; ) {
            const m = l.shift();
            o ? J(e, o, () => p(m)) : p(m);
          }
          c();
        } catch (m) {
          e.error(m);
        }
    }));
  };
  return t.subscribe(R(e, d, () => {
    h = !0, c();
  })), () => {
    a == null || a();
  };
}
function _t(t, e, n = 1 / 0) {
  return v(e) ? _t((r, i) => xt((s, o) => e(r, s, i, o))(U(t(r, i))), n) : (typeof e == "number" && (n = e), O((r, i) => Oo(r, i, t, n)));
}
function Lo(t = 1 / 0) {
  return _t(ue, t);
}
function Tr() {
  return Lo(1);
}
function sn(...t) {
  return Tr()(Je(t, De(t)));
}
function Po(t) {
  return new $((e) => {
    U(t()).subscribe(e);
  });
}
const Uo = ["addListener", "removeListener"], Fo = ["addEventListener", "removeEventListener"], Mo = ["on", "off"];
function lt(t, e, n, r) {
  if (v(n) && (r = n, n = void 0), r)
    return lt(t, e, n).pipe(Co(r));
  const [i, s] = Jo(t) ? Fo.map((o) => (a) => t[o](e, a, n)) : No(t) ? Uo.map(on(t, e)) : Do(t) ? Mo.map(on(t, e)) : [];
  if (!i && Et(t))
    return _t((o) => lt(o, e, n))(U(t));
  if (!i)
    throw new TypeError("Invalid event target");
  return new $((o) => {
    const a = (...l) => o.next(1 < l.length ? l : l[0]);
    return i(a), () => s(a);
  });
}
function on(t, e) {
  return (n) => (r) => t[n](e, r);
}
function No(t) {
  return v(t.addListener) && v(t.removeListener);
}
function Do(t) {
  return v(t.on) && v(t.off);
}
function Jo(t) {
  return v(t.addEventListener) && v(t.removeEventListener);
}
function Er(t = 0, e, n = ao) {
  let r = -1;
  return e != null && (ur(e) ? n = e : r = e), new $((i) => {
    let s = Io(t) ? +t - n.now() : t;
    s < 0 && (s = 0);
    let o = 0;
    return n.schedule(function() {
      i.closed || (i.next(o++), 0 <= r ? this.schedule(void 0, r) : i.complete());
    }, s);
  });
}
const Ho = new $(st);
function At(t) {
  return O((e, n) => {
    let r = null, i = !1, s;
    r = e.subscribe(R(n, void 0, void 0, (o) => {
      s = U(t(o, At(t)(e))), r ? (r.unsubscribe(), r = null, s.subscribe(n)) : i = !0;
    })), i && (r.unsubscribe(), r = null, s.subscribe(n));
  });
}
function Bo(...t) {
  const e = De(t);
  return O((n, r) => {
    Tr()(Je([n, ...t], e)).subscribe(r);
  });
}
function xr(...t) {
  return Bo(...t);
}
function an(t, e = ue) {
  return t = t ?? jo, O((n, r) => {
    let i, s = !0;
    n.subscribe(R(r, (o) => {
      const a = e(o);
      (s || !t(i, a)) && (s = !1, i = a, r.next(o));
    }));
  });
}
function jo(t, e) {
  return t === e;
}
function zo(t = 1 / 0) {
  let e;
  t && typeof t == "object" ? e = t : e = {
    count: t
  };
  const { count: n = 1 / 0, delay: r, resetOnSuccess: i = !1 } = e;
  return n <= 0 ? ue : O((s, o) => {
    let a = 0, l;
    const u = () => {
      let f = !1;
      l = s.subscribe(R(o, (h) => {
        i && (a = 0), o.next(h);
      }, void 0, (h) => {
        if (a++ < n) {
          const c = () => {
            l ? (l.unsubscribe(), l = null, u()) : f = !0;
          };
          if (r != null) {
            const d = typeof r == "number" ? Er(r) : U(r(h, a)), p = R(o, () => {
              p.unsubscribe(), c();
            }, () => {
              o.complete();
            });
            d.subscribe(p);
          } else
            c();
        } else
          o.error(h);
      })), f && (l.unsubscribe(), l = null, u());
    };
    u();
  });
}
function Wo(...t) {
  const e = De(t);
  return O((n, r) => {
    (e ? sn(t, n, e) : sn(t, n)).subscribe(r);
  });
}
function It(t, e) {
  return O((n, r) => {
    let i = null, s = 0, o = !1;
    const a = () => o && !i && r.complete();
    n.subscribe(R(r, (l) => {
      i == null || i.unsubscribe();
      let u = 0;
      const f = s++;
      U(t(l, f)).subscribe(i = R(r, (h) => r.next(e ? e(l, h, f, u++) : h), () => {
        i = null, a();
      }));
    }, () => {
      o = !0, a();
    }));
  });
}
function Ko(t, e, n) {
  const r = v(t) || e || n ? { next: t, error: e, complete: n } : t;
  return r ? O((i, s) => {
    var o;
    (o = r.subscribe) === null || o === void 0 || o.call(r);
    let a = !0;
    i.subscribe(R(s, (l) => {
      var u;
      (u = r.next) === null || u === void 0 || u.call(r, l), s.next(l);
    }, () => {
      var l;
      a = !1, (l = r.complete) === null || l === void 0 || l.call(r), s.complete();
    }, (l) => {
      var u;
      a = !1, (u = r.error) === null || u === void 0 || u.call(r, l), s.error(l);
    }, () => {
      var l, u;
      a && ((l = r.unsubscribe) === null || l === void 0 || l.call(r)), (u = r.finalize) === null || u === void 0 || u.call(r);
    }));
  }) : ue;
}
const qo = "http://./polyfea".replace(/\/+$/, "");
class se {
  constructor(e = {}) {
    this.configuration = e;
  }
  set config(e) {
    this.configuration = e;
  }
  get basePath() {
    return this.configuration.basePath != null ? this.configuration.basePath : qo;
  }
  get fetchApi() {
    return this.configuration.fetchApi;
  }
  get middleware() {
    return this.configuration.middleware || [];
  }
  get queryParamsStringify() {
    return this.configuration.queryParamsStringify || _r;
  }
  get username() {
    return this.configuration.username;
  }
  get password() {
    return this.configuration.password;
  }
  get apiKey() {
    const e = this.configuration.apiKey;
    if (e)
      return typeof e == "function" ? e : () => e;
  }
  get accessToken() {
    const e = this.configuration.accessToken;
    if (e)
      return typeof e == "function" ? e : async () => e;
  }
  get headers() {
    return this.configuration.headers;
  }
  get credentials() {
    return this.configuration.credentials;
  }
}
const Go = new se();
class He {
  constructor(e = Go) {
    this.configuration = e, this.fetchApi = async (n, r) => {
      let i = { url: n, init: r };
      for (const o of this.middleware)
        o.pre && (i = await o.pre(Object.assign({ fetch: this.fetchApi }, i)) || i);
      let s;
      try {
        s = await (this.configuration.fetchApi || fetch)(i.url, i.init);
      } catch (o) {
        for (const a of this.middleware)
          a.onError && (s = await a.onError({
            fetch: this.fetchApi,
            url: i.url,
            init: i.init,
            error: o,
            response: s ? s.clone() : void 0
          }) || s);
        if (s === void 0)
          throw o instanceof Error ? new Zo(o, "The request failed and the interceptors did not return an alternative response") : o;
      }
      for (const o of this.middleware)
        o.post && (s = await o.post({
          fetch: this.fetchApi,
          url: i.url,
          init: i.init,
          response: s.clone()
        }) || s);
      return s;
    }, this.middleware = e.middleware;
  }
  withMiddleware(...e) {
    const n = this.clone();
    return n.middleware = n.middleware.concat(...e), n;
  }
  withPreMiddleware(...e) {
    const n = e.map((r) => ({ pre: r }));
    return this.withMiddleware(...n);
  }
  withPostMiddleware(...e) {
    const n = e.map((r) => ({ post: r }));
    return this.withMiddleware(...n);
  }
  /**
   * Check if the given MIME is a JSON MIME.
   * JSON MIME examples:
   *   application/json
   *   application/json; charset=UTF8
   *   APPLICATION/JSON
   *   application/vnd.company+json
   * @param mime - MIME (Multipurpose Internet Mail Extensions)
   * @return True if the given MIME is JSON, false otherwise.
   */
  isJsonMime(e) {
    return e ? He.jsonRegex.test(e) : !1;
  }
  async request(e, n) {
    const { url: r, init: i } = await this.createFetchParams(e, n), s = await this.fetchApi(r, i);
    if (s && s.status >= 200 && s.status < 300)
      return s;
    throw new Xo(s, "Response returned an error code");
  }
  async createFetchParams(e, n) {
    let r = this.configuration.basePath + e.path;
    e.query !== void 0 && Object.keys(e.query).length !== 0 && (r += "?" + this.configuration.queryParamsStringify(e.query));
    const i = Object.assign({}, this.configuration.headers, e.headers);
    Object.keys(i).forEach((f) => i[f] === void 0 ? delete i[f] : {});
    const s = typeof n == "function" ? n : async () => n, o = {
      method: e.method,
      headers: i,
      body: e.body,
      credentials: this.configuration.credentials
    }, a = Object.assign(Object.assign({}, o), await s({
      init: o,
      context: e
    }));
    let l;
    Qo(a.body) || a.body instanceof URLSearchParams || Yo(a.body) ? l = a.body : this.isJsonMime(i["Content-Type"]) ? l = JSON.stringify(a.body) : l = a.body;
    const u = Object.assign(Object.assign({}, a), { body: l });
    return { url: r, init: u };
  }
  /**
   * Create a shallow clone of `this` by constructing a new instance
   * and then shallow cloning data members.
   */
  clone() {
    const e = this.constructor, n = new e(this.configuration);
    return n.middleware = this.middleware.slice(), n;
  }
}
He.jsonRegex = new RegExp("^(:?application/json|[^;/ 	]+/[^;/ 	]+[+]json)[ 	]*(:?;.*)?$", "i");
function Yo(t) {
  return typeof Blob < "u" && t instanceof Blob;
}
function Qo(t) {
  return typeof FormData < "u" && t instanceof FormData;
}
class Xo extends Error {
  constructor(e, n) {
    super(n), this.response = e, this.name = "ResponseError";
  }
}
class Zo extends Error {
  constructor(e, n) {
    super(n), this.cause = e, this.name = "FetchError";
  }
}
class ln extends Error {
  constructor(e, n) {
    super(n), this.field = e, this.name = "RequiredError";
  }
}
function E(t, e) {
  const n = t[e];
  return n != null;
}
function _r(t, e = "") {
  return Object.keys(t).map((n) => Ar(n, t[n], e)).filter((n) => n.length > 0).join("&");
}
function Ar(t, e, n = "") {
  const r = n + (n.length ? `[${t}]` : t);
  if (e instanceof Array) {
    const i = e.map((s) => encodeURIComponent(String(s))).join(`&${encodeURIComponent(r)}=`);
    return `${encodeURIComponent(r)}=${i}`;
  }
  if (e instanceof Set) {
    const i = Array.from(e);
    return Ar(t, i, n);
  }
  return e instanceof Date ? `${encodeURIComponent(r)}=${encodeURIComponent(e.toISOString())}` : e instanceof Object ? _r(e, r) : `${encodeURIComponent(r)}=${encodeURIComponent(String(e))}`;
}
function Ir(t, e) {
  return Object.keys(t).reduce((n, r) => Object.assign(Object.assign({}, n), { [r]: e(t[r]) }), {});
}
class cn {
  constructor(e, n = (r) => r) {
    this.raw = e, this.transformer = n;
  }
  async value() {
    return this.transformer(await this.raw.json());
  }
}
function Vo(t) {
  return ea(t);
}
function ea(t, e) {
  return t == null ? t : {
    microfrontend: E(t, "microfrontend") ? t.microfrontend : void 0,
    tagName: t.tagName,
    attributes: E(t, "attributes") ? t.attributes : void 0,
    style: E(t, "style") ? t.style : void 0
  };
}
function ta(t) {
  return na(t);
}
function na(t, e) {
  return t == null ? t : {
    kind: E(t, "kind") ? t.kind : void 0,
    href: E(t, "href") ? t.href : void 0,
    attributes: E(t, "attributes") ? t.attributes : void 0,
    waitOnLoad: E(t, "waitOnLoad") ? t.waitOnLoad : void 0
  };
}
function kr(t) {
  return ra(t);
}
function ra(t, e) {
  return t == null ? t : {
    dependsOn: E(t, "dependsOn") ? t.dependsOn : void 0,
    module: E(t, "module") ? t.module : void 0,
    resources: E(t, "resources") ? t.resources.map(ta) : void 0
  };
}
function Rr(t) {
  return ia(t);
}
function ia(t, e) {
  return t == null ? t : {
    elements: t.elements.map(Vo),
    microfrontends: E(t, "microfrontends") ? Ir(t.microfrontends, kr) : void 0
  };
}
function sa(t) {
  return oa(t);
}
function oa(t, e) {
  return t == null ? t : {
    name: t.name,
    path: E(t, "path") ? t.path : void 0,
    contextArea: E(t, "contextArea") ? Rr(t.contextArea) : void 0
  };
}
function aa(t) {
  return la(t);
}
function la(t, e) {
  return t == null ? t : {
    contextAreas: E(t, "contextAreas") ? t.contextAreas.map(sa) : void 0,
    microfrontends: Ir(t.microfrontends, kr)
  };
}
class Re extends He {
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextAreaRaw(e, n) {
    if (e.name === null || e.name === void 0)
      throw new ln("name", "Required parameter requestParameters.name was null or undefined when calling getContextArea.");
    if (e.path === null || e.path === void 0)
      throw new ln("path", "Required parameter requestParameters.path was null or undefined when calling getContextArea.");
    const r = {};
    e.path !== void 0 && (r.path = e.path), e.take !== void 0 && (r.take = e.take);
    const i = {}, s = await this.request({
      path: "/context-area/{name}".replace("{name}", encodeURIComponent(String(e.name))),
      method: "GET",
      headers: i,
      query: r
    }, n);
    return new cn(s, (o) => Rr(o));
  }
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextArea(e, n) {
    return await (await this.getContextAreaRaw(e, n)).value();
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfigRaw(e) {
    const n = {}, r = {}, i = await this.request({
      path: "/static-config",
      method: "GET",
      headers: r,
      query: n
    }, e);
    return new cn(i, (s) => aa(s));
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfig(e) {
    return await (await this.getStaticConfigRaw(e)).value();
  }
}
class ca {
  constructor(e = "./polyfea") {
    this.spec$ = new $();
    let n;
    typeof e == "string" ? (e.length === 0 && (e = "./polyfea"), n = new Re(new se({ basePath: e }))) : e instanceof se ? n = new Re(e) : n = e, this.spec$ = Je(n.getStaticConfig()).pipe(xr(Ho));
  }
  getContextArea(e) {
    let n = globalThis.location.pathname;
    if (globalThis.document.baseURI) {
      const i = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(i) && (n = "./" + n.substring(i.length));
    }
    return this.spec$.pipe(xt((r) => {
      for (let i of r.contextAreas)
        if (i.name === e && new RegExp(i.path).test(n))
          return Object.assign(Object.assign({}, i.contextArea), { microfrontends: Object.assign(Object.assign({}, r.microfrontends), i.contextArea.microfrontends) });
      return null;
    }));
  }
}
class un {
  constructor(e = "./polyfea") {
    var n, r;
    typeof e == "string" ? this.api = new Re(new se({
      basePath: new URL(e, new URL(((n = globalThis.document) === null || n === void 0 ? void 0 : n.baseURI) || "/", ((r = globalThis.location) === null || r === void 0 ? void 0 : r.href) || "http://localhost")).href
    })) : e instanceof se ? this.api = new Re(e) : this.api = e;
  }
  getContextArea(e) {
    var n, r;
    let i = ((n = globalThis.location) === null || n === void 0 ? void 0 : n.pathname) || "/";
    if (!((r = globalThis.document) === null || r === void 0) && r.baseURI) {
      const l = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      i.startsWith(l) && (i = "./" + i.substring(l.length));
    }
    const s = localStorage.getItem(`polyfea-context[${e},${i}]`), o = Po(() => this.api.getContextAreaRaw({ name: e, path: i })).pipe(It((a) => a.raw.ok ? a.value() : _o(() => new Error(a.raw.statusText))), Ko((a) => {
      a && localStorage.setItem(`polyfea-context[${e},${i}]`, JSON.stringify(a));
    }));
    if (s) {
      const a = JSON.parse(s);
      return ot(a).pipe(xr(o), At((l) => (console.warn(`Failed to fetch context area ${e} from ${i}, using cached version as the last known value`, l), ot(a))));
    } else
      return o.pipe(zo({ count: 3, delay: (a) => Er((a + 1) * 2e3) }));
  }
}
class ua {
  /**
   * Constructs a new NavigationDestination instance.
   * @param url - The URL of the navigation destination.
   */
  constructor(e) {
    this.url = e;
  }
}
class fa extends Event {
  constructor(e, n) {
    super("navigate", { bubbles: !0, cancelable: !0 }), this.transition = e, this.interceptPromises = [], this.downloadRequest = null, this.formData = null, this.hashChange = !1, this.userInitiated = !1;
    let r = new URL(e.href, new URL(globalThis.document.baseURI, globalThis.location.href));
    const i = new URL(n.url, new URL(globalThis.document.baseURI, globalThis.location.href));
    this.canIntercept = i.protocol === r.protocol && i.host === r.host && i.port === r.port, this.destination = new ua(r.href);
  }
  /** (@see https://developer.mozilla.org/en-US/docs/Web/API/NavigateEvent )
   *  this polyfill signals abort only on programatic navigation
  **/
  get signal() {
    return this.transition.abortController.signal;
  }
  /**
   * Prevents the browser from following the navigation request.
   * @see {@link https://developer.mozilla.org/en-US/docs/Web/API/NavigateEvent/preventDefault}
   */
  intercept(e) {
    e != null && e.handler && this.interceptPromises.push(e.handler(this));
  }
}
class fn extends Event {
  constructor(e, n) {
    super("currententrychange", { bubbles: !0, cancelable: !0 }), this.navigationType = e, this.from = n;
  }
}
class hn {
  constructor(e) {
    this.request = e;
  }
  get finished() {
    return at(this.request.finished);
  }
  get from() {
    return this.request.entry;
  }
  get type() {
    return this.request.mode;
  }
}
function ha(t = !1) {
  return Ce.tryRegister(t);
}
class da {
  constructor(e) {
    this.commited$ = e.committed, this.finished$ = e.finished;
  }
  get commited() {
    return at(this.commited$);
  }
  get finished() {
    return at(this.finished$);
  }
}
class Ce extends EventTarget {
  entries() {
    return this.entriesList;
  }
  get currentEntry() {
    if (this.currentEntryIndex >= 0)
      return this.entriesList[this.currentEntryIndex];
  }
  get canGoBack() {
    return this.currentEntryIndex > 0;
  }
  get canGoForward() {
    return this.currentEntryIndex + 1 < this.entriesList.length;
  }
  get transition() {
    var e;
    return ((e = this.currentTransition) === null || e === void 0 ? void 0 : e.transition) || null;
  }
  navigate(e, n) {
    let r = "push";
    return ((n == null ? void 0 : n.history) === "replace" || n != null && n.replace) && (r = "replace"), this.nextTransitionRequest(r, e, n);
  }
  back() {
    if (this.currentEntryIndex < 1)
      throw { name: "InvaliStateError", message: "Cannot go back from initial state" };
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { traverseTo: this.entriesList[this.currentEntryIndex - 1].key });
  }
  forward(e) {
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { info: e, traverseTo: this.entriesList[this.currentEntryIndex + 1].key });
  }
  reload(e) {
    return this.nextTransitionRequest("reload", this.currentEntry.url, { info: e });
  }
  traverseTo(e, n) {
    const r = this.entriesList.find((i) => i.key === e);
    if (!r)
      throw { name: "InvaliStateError", message: "Cannot traverse to unknown state" };
    return this.nextTransitionRequest("traverse", r.url, { info: n == null ? void 0 : n.info, traverseTo: e });
  }
  updateCurrentEntry(e) {
    this.entriesList[this.currentEntryIndex].setState(JSON.parse(JSON.stringify(e == null ? void 0 : e.state))), this.currentEntry.dispatchEvent(new fn("replace", this.currentEntry));
  }
  constructor() {
    super(), this.entriesList = [], this.idCounter = 0, this.transitionRequests = new Ne(), this.currentTransition = null, this.currentEntryIndex = -1, this.pushstateDelay = 35, this.rawHistoryMethods = {
      pushState: globalThis.history.pushState,
      replaceState: globalThis.history.replaceState,
      go: globalThis.history.go,
      back: globalThis.history.back,
      forward: globalThis.history.forward
    }, this.transitionRequests.subscribe((e) => this.executeRequest(e));
  }
  nextTransitionRequest(e, n, r) {
    const i = `@${++this.idCounter}-navigation-polyfill-transition`, s = {
      mode: e,
      href: new URL(n, new URL(globalThis.document.baseURI, globalThis.location.href)).href,
      info: r == null ? void 0 : r.info,
      state: r == null ? void 0 : r.state,
      committed: new me(),
      finished: new me(),
      entry: new Ye(this, i, i, n.toString(), r == null ? void 0 : r.state),
      abortController: new AbortController(),
      traverseToKey: r == null ? void 0 : r.traverseTo,
      transition: null
    };
    return this.transitionRequests.next(s), new da(s);
  }
  async executeRequest(e) {
    this.currentTransition && (this.currentTransition.abortController.abort(), this.currentTransition.finished.error("aborted - new navigation started"), this.currentTransition.committed.closed || this.currentTransition.committed.error("aborted - new navigation started"), globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: Error("aborted - new navigation started") })), this.clearTransition(e)), e.transition = new hn(e), this.currentTransition = e;
    try {
      await this.commit(e), this.currentEntry.dispatchEvent(new fn(e.mode, e.transition.from)), await this.dispatchNavigation(e);
    } catch {
      e.finished.error("aborted"), e.committed.error("aborted");
    } finally {
      this.clearTransition(e);
    }
  }
  dispatchNavigation(e) {
    var n;
    const r = new fa(e, this.currentEntry);
    if (!((n = globalThis.navigation) === null || n === void 0) && n.dispatchEvent(r))
      return r.interceptPromises.length > 0 ? Promise.all(r.interceptPromises.filter((i) => !!(i != null && i.then))).then(() => {
        globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(e), e.finished.next(), e.finished.complete();
      }).catch((i) => {
        globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: i })), this.clearTransition(e), e.finished.error(i);
      }) : (globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(e), e.finished.next(), e.finished.complete(), Promise.resolve());
  }
  clearTransition(e) {
    var n;
    ((n = this.currentTransition) === null || n === void 0 ? void 0 : n.entry.id) === e.entry.id && (this.currentTransition = null);
  }
  commit(e) {
    switch (e.mode) {
      case "push":
        return this.commitPushTransition(e);
      case "replace":
        return this.commitReplaceTransition(e);
      case "reload":
        return this.commitReloadTransition(e);
      case "traverse":
        return this.commitTraverseTransition(e);
    }
  }
  async pushstateAsync(e, n = () => {
  }) {
    return new Promise((r, i) => {
      setTimeout(() => {
        n(e), e.committed.next(), e.committed.complete(), r();
      }, this.pushstateDelay);
    });
  }
  commitPushTransition(e) {
    return this.rawHistoryMethods.pushState.apply(globalThis.history, [e.entry.cloneable, "", e.href]), this.pushstateAsync(e, (n) => {
      this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), n.entry];
    });
  }
  commitReplaceTransition(e) {
    return e.entry.key = this.currentEntry.key, this.entriesList[this.currentEntryIndex] = e.entry, this.rawHistoryMethods.replaceState.apply(globalThis.history, [e.entry.cloneable, "", e.href]), this.pushstateAsync(e);
  }
  commitTraverseTransition(e) {
    return new Promise(async (n, r) => {
      const i = this.entriesList.findIndex((o) => o.key === e.traverseToKey);
      i < 0 && r("target entry not found");
      const s = i - this.currentEntryIndex;
      this.rawHistoryMethods.go.apply(globalThis.history, [s]), await this.pushstateAsync(e, (o) => {
        const a = this.entriesList.findIndex((l) => l.key === o.traverseToKey);
        a < 0 && o.committed.error(new Error("target entry not found")), this.currentEntryIndex = a, o.committed.next(), o.committed.complete(), n();
      });
    });
  }
  commitReloadTransition(e) {
    return e.committed.next(), e.committed.complete(), e.finished.subscribe({
      next: () => globalThis.location.reload(),
      error: () => globalThis.location.reload()
    }), Promise.resolve();
  }
  static tryRegister(e = !1) {
    if (!globalThis.navigation) {
      const n = new Ce();
      return n.doRegister(e), n;
    }
    return globalThis.navigation;
  }
  static unregister() {
    globalThis.navigation && globalThis.navigation instanceof Ce && (globalThis.navigation.doUnregister && globalThis.navigation.doUnregister(), globalThis.navigation = void 0);
  }
  doRegister(e) {
    var n, r, i, s, o;
    if (!globalThis.navigation && !globalThis.navigation) {
      globalThis.navigation = this, this.entriesList = [new Ye(this, "initial", "initial", globalThis.location.href, void 0)], this.currentEntryIndex = 0;
      const a = ((n = globalThis.history) === null || n === void 0 ? void 0 : n.pushState) || ((c, d, p) => {
      }), l = ((r = globalThis.history) === null || r === void 0 ? void 0 : r.replaceState) || ((c, d, p) => {
      }), u = ((i = globalThis.history) === null || i === void 0 ? void 0 : i.go) || ((c) => {
      }), f = ((s = globalThis.history) === null || s === void 0 ? void 0 : s.back) || (() => {
      }), h = ((o = globalThis.history) === null || o === void 0 ? void 0 : o.forward) || (() => {
      });
      this.doUnregister = () => {
        globalThis.history.pushState = a, globalThis.history.replaceState = l, globalThis.history.go = u, globalThis.history.back = f, globalThis.history.forward = h;
      }, e ? (this.rawHistoryMethods.pushState = () => {
        a.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.replaceState = () => {
        l.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.go = () => {
        u.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.back = () => {
        f.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }, this.rawHistoryMethods.forward = () => {
        h.apply(globalThis.history, arguments);
        const c = new PopStateEvent("popstate", { state: this.currentTransition.entry.cloneable });
        c.state = this.currentTransition.entry.cloneable, setTimeout(() => globalThis.dispatchEvent(c), 25);
      }) : this.rawHistoryMethods = {
        pushState: a || ((c, d, p) => {
        }),
        replaceState: l || ((c, d, p) => {
        }),
        go: u || ((c) => {
        }),
        back: f || (() => {
        }),
        forward: h || (() => {
        })
      }, globalThis.history && (globalThis.history.pushState = (c, d, p) => this.navigate(p, { state: c, history: "push" }), globalThis.history.replaceState = (c, d, p) => this.navigate(p, { state: c, history: "replace" }), globalThis.history.go = (c) => this.traverseTo(this.entriesList[this.currentEntryIndex + c].key), globalThis.history.back = () => this.back(), globalThis.history.forward = () => this.forward()), globalThis.addEventListener("popstate", (c) => {
        var d, p, y, g, m, A, W;
        if (this.currentTransition && ((d = c.state) === null || d === void 0 ? void 0 : d.id) === ((y = (p = this.currentTransition) === null || p === void 0 ? void 0 : p.entry) === null || y === void 0 ? void 0 : y.id))
          return;
        (g = this.currentTransition) === null || g === void 0 || g.abortController.abort();
        const te = new me();
        te.complete();
        let N;
        if (!((m = c.state) === null || m === void 0) && m.key) {
          const x = this.entriesList.findIndex((Be) => {
            var q;
            return Be.key === ((q = c.state) === null || q === void 0 ? void 0 : q.key);
          });
          x >= 0 && (this.currentEntryIndex = x, N = this.entriesList[x]);
        }
        if (!N) {
          let x = `@${++this.idCounter}-navigation-polyfill-popstate`;
          N = new Ye(this, x, x, globalThis.location.href, c.state), this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), N];
        }
        const fe = new me(), K = {
          mode: "traverse",
          href: globalThis.location.href,
          info: void 0,
          state: ((A = c.state) === null || A === void 0 ? void 0 : A.state) || c.state,
          committed: te,
          finished: fe,
          entry: N,
          abortController: new AbortController(),
          traverseToKey: (W = c.state) === null || W === void 0 ? void 0 : W.key,
          transition: null
        };
        K.transition = new hn(K), this.currentTransition = K, this.dispatchNavigation(this.currentTransition);
      });
    }
  }
}
class Ye extends EventTarget {
  constructor(e, n, r, i, s) {
    super(), this.owner = e, this.id = n, this.key = r, this.url = i, this.state = s, this.url = new URL(i, new URL(globalThis.document.baseURI, globalThis.location.href)).href;
  }
  get index() {
    return this.owner.entriesList.findIndex((e) => e.id === this.id);
  }
  get sameDocument() {
    return !0;
  }
  // polyfill is lost between documents
  getState() {
    return this.state;
  }
  setState(e) {
    this.state = e;
  }
  get cloneable() {
    return {
      id: this.id,
      key: this.key,
      url: this.url,
      index: this.index,
      state: this.state
    };
  }
}
class pa {
  /** @internal @private */
  constructor() {
  }
  /** @static
   *
   * Get or create a polyfea driver instance. If the instance is provided on the global context, it is returned.
   * Otherwise, a new instance is created with the given configuration.
   *
   * @param config - Configuration for the [`PolyfeaApi`](https://github.com/polyfea/browser-api/blob/main/docs/classes/PolyfeaApi.md).
   **/
  static getOrCreate(e) {
    return globalThis.polyfea ? globalThis.polyfea : new Oe(e);
  }
  /** @static
   * Initialize the polyfea driver in the global context.
   * This method is typically invoked by the polyfea controller script `boot.ts`.
   *
   * @remarks
   * This method also initializes the Navigation polyfill if it's not already present.
   * It augments `window.customElements.define` to allow for duplicate registration of custom elements.
   * This is particularly useful when different microfrontends need to register the same dependencies.
   */
  static initialize() {
    globalThis.polyfea || Oe.install();
  }
}
class Oe {
  constructor(e) {
    this.config = e, this.loadedResources = /* @__PURE__ */ new Set(), globalThis.navigation && globalThis.navigation.addEventListener("navigate", (n) => {
      n.canIntercept && n.destination.url.startsWith(document.baseURI) && n.intercept();
    });
  }
  getBackend() {
    var e;
    if (!this.backend) {
      let n = (e = document.querySelector('meta[name="polyfea.backend"]')) === null || e === void 0 ? void 0 : e.getAttribute("content");
      if (n)
        if (n.startsWith("static://")) {
          const r = n.slice(9);
          this.backend = new ca(this.config || r);
        } else
          this.backend = new un(this.config || n);
      else
        this.backend = new un(this.config || "./polyfea");
    }
    return this.backend;
  }
  getContextArea(e) {
    return globalThis.navigation ? lt(globalThis.navigation, "navigatesuccess").pipe(Wo(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), It((n) => this.getBackend().getContextArea(e).pipe(At((r) => (console.error(r), ot({ elements: [], microfrontends: {} }))))), an((n, r) => JSON.stringify(n) === JSON.stringify(r))) : this.getBackend().getContextArea(e).pipe(an((n, r) => JSON.stringify(n) === JSON.stringify(r)));
  }
  loadMicrofrontend(e, n) {
    if (!n)
      return Promise.resolve();
    const r = [];
    return this.loadMicrofrontendRecursive(e, n, r);
  }
  async loadMicrofrontendRecursive(e, n, r) {
    if (r.includes(n))
      throw new Error("Circular dependency detected: " + r.join(" -> "));
    const i = e.microfrontends[n];
    if (!i)
      throw new Error("Microfrontend specification not found: " + n);
    r.push(n), i.dependsOn && await Promise.all(i.dependsOn.map((o) => this.loadMicrofrontendRecursive(e, o, r).catch((a) => {
      console.error(`Failed to load microfrontend's ${n} dependency ${o}`, a);
    })));
    let s = i.resources || [];
    i.module && (s = [...s, {
      kind: "script",
      href: i.module,
      attributes: {
        type: "module"
      },
      waitOnLoad: !0
    }]), await Promise.all(s.map((o) => {
      if (this.loadedResources.has(o.href))
        return Promise.resolve();
      switch (o.kind) {
        case "script":
          return this.loadScript(o);
        case "stylesheet":
          return this.loadStylesheet(o);
        case "link":
          return this.loadLink(o);
      }
    }));
  }
  loadScript(e) {
    return new Promise((n, r) => {
      var i;
      const s = document.createElement("script");
      s.src = e.href, s.setAttribute("async", ""), e.attributes && Object.entries(e.attributes).forEach(([a, l]) => {
        s.setAttribute(a, l);
      });
      const o = (i = document.querySelector('meta[name="csp-nonce"]')) === null || i === void 0 ? void 0 : i.getAttribute("content");
      o && s.setAttribute("nonce", o), this.loadedResources.add(e.href), e.waitOnLoad ? (s.onload = () => {
        n();
      }, s.onerror = () => {
        this.loadedResources.delete(e.href), console.error(`Failed to load script ${e.href} while loading microfrontend resources, check the network tab for details`), n();
      }) : n(), document.head.appendChild(s);
    });
  }
  loadStylesheet(e) {
    return this.loadLink(Object.assign(Object.assign({}, e), { attributes: Object.assign(Object.assign({}, e.attributes), { rel: "stylesheet" }) }));
  }
  loadLink(e) {
    var n;
    const r = document.createElement("link");
    r.href = e.href, r.setAttribute("async", "");
    const i = (n = document.querySelector('meta[name="csp-nonce"]')) === null || n === void 0 ? void 0 : n.getAttribute("content");
    return i && r.setAttribute("nonce", i), e.attributes && Object.entries(e.attributes).forEach(([s, o]) => {
      r.setAttribute(s, o);
    }), new Promise((s, o) => {
      this.loadedResources.add(e.href), e.waitOnLoad ? (r.onload = () => {
        s();
      }, r.onerror = () => {
        this.loadedResources.delete(e.href), console.error(`Failed to load ${e.href} while loading microfrontend resources, check the network tab for details`), s();
      }) : s(), document.head.appendChild(r);
    });
  }
  static install() {
    var e;
    globalThis.polyfea || (globalThis.polyfea = new Oe()), ha();
    const n = (e = document.querySelector('meta[name="polyfea.duplicit-custom-elements"]')) === null || e === void 0 ? void 0 : e.getAttribute("content");
    let r = "warn";
    n === "silent" ? r = "silent" : n === "error" ? r = "error" : n === "verbose" && (r = "verbose");
    function i(s) {
      if (s.overrider === "polyfea")
        return s;
      const o = function(...a) {
        if (this.get(a[0])) {
          if (r === "error")
            throw new Error(`Custom element '${a[0]}' is duplicately registered`);
          if (r === "warn")
            return console.warn(`Custom element '${a[0]}' is duplicately registered - ignoring the current attempt for registration`), !1;
        } else
          return r === "verbose" && console.log(`Custom element '${a[0]}' is registered`), s.apply(this, a);
      };
      return o.overrider = "polyfea", o;
    }
    customElements.define = i(customElements.define);
  }
}
const ya = ":host{display:contents;width:100%;height:100%;box-sizing:border-box}.cyclic-error{color:red}", ma = /* @__PURE__ */ Rs(class extends Fs {
  constructor() {
    super(), this.__registerHost(), this.__attachShadow(), this.cyclicAreas = "none", this.cyclicErrorMsg = "", this.contextName = void 0, this.name = void 0, this.take = void 0, this.extraAttributes = {}, this.extraStyle = {}, this.polyfeaContextStack = void 0, this.contextObj = void 0;
  }
  get areaName() {
    return this.contextName || this.name;
  }
  async componentWillLoad() {
    this.areaName && (this.polyfea = pa.getOrCreate(), this.polyfea.getContextArea(this.areaName).pipe(
      // load microfrontends
      It((e) => {
        if (!e)
          return Promise.resolve(e);
        let n = (e.elements || []).slice(0, this.take).map((r) => r.microfrontend ? this.polyfea.loadMicrofrontend(e, r.microfrontend) : Promise.resolve());
        return Promise.all(n).then((r) => e);
      })
    ).subscribe({
      next: (e) => this.contextObj = e,
      error: (e) => {
        console.warn(`<polyfe-context name="${this.areaName}">: Using slotted content because of error: ${e}`);
      }
    }));
  }
  componentWillRender() {
    if (this.polyfeaContextStack)
      return;
    let e = this.hostElement, n = null;
    for (; n === null && e && e.tagName !== "BODY"; )
      e = e.parentElement, n = (e == null ? void 0 : e.polyfeaContextStack) || null;
    if (n = n || [], n.indexOf(this.areaName) >= 0) {
      this.cyclicAreas = "error", this.cyclicErrorMsg = "";
      let r = document.head.querySelector("meta[name='polyfea.cyclic-context-areas']");
      if (r) {
        const i = r.getAttribute("content");
        (i == "allow" || i == "silent") && (this.cyclicAreas = i);
      }
      r = document.head.querySelector("meta[name='polyfea.cyclic-context-message']"), r && (this.cyclicErrorMsg = r.getAttribute("content")), this.cyclicErrorMsg || (this.cyclicErrorMsg = "Cyclic rendering of context areas detected: <br/>{stack}</br> Area ignored to avoid infinite recursion."), this.cyclicErrorMsg = this.cyclicErrorMsg.replace("{stack}", n.map((i) => i === this.areaName ? `<b>${i}</b>` : i).join(" -> ") + " ==> <b>" + this.areaName + "</b>");
    }
    n.push(this.areaName), this.polyfeaContextStack = n;
  }
  render() {
    var e;
    if (this.cyclicAreas == "error")
      return Q("div", { class: "cyclic-error", innerHTML: this.cyclicErrorMsg });
    if (this.cyclicAreas == "silent")
      return "";
    let n = ((e = this.contextObj) === null || e === void 0 ? void 0 : e.elements) || [];
    return this.take > 0 && (n = n.slice(0, this.take)), Q(Kn, null, n.map((r) => this.renderElement(r)), n.length ? "" : Q("slot", null));
  }
  renderElement(e) {
    const n = e.tagName, r = Object.assign({
      context: this.areaName
    }, e.attributes, this.extraAttributes, {
      class: this.areaName + "-context"
    }), i = Object.assign({}, e.style, this.extraStyle);
    return Q(n, { style: i, ref: (s) => {
      if (s)
        for (let o in r)
          s.setAttribute(o, r[o]);
    } });
  }
  get hostElement() {
    return this;
  }
  static get style() {
    return ya;
  }
}, [1, "polyfea-context", {
  contextName: [1, "context-name"],
  name: [1],
  take: [2],
  extraAttributes: [16],
  extraStyle: [16],
  polyfeaContextStack: [1040],
  contextObj: [32]
}]), ga = ma;
/*!
 * Part of Polyfea microfrontends suite - https://github.com/polyfea
 */
const ba = (t) => {
  typeof customElements < "u" && [
    ga
  ].forEach((e) => {
    customElements.get(e.is) || customElements.define(e.is, e, t);
  });
};
globalThis.addEventListener("load", () => {
  if (Vi.initialize(), ba(), !document.body.hasAttribute("polyfea")) {
    document.body.setAttribute("polyfea", "initialized");
    const t = document.createElement("polyfea-context");
    t.setAttribute("name", "shell"), t.setAttribute("take", "1"), document.body.appendChild(t);
  }
});
//# sourceMappingURL=boot.mjs.map
