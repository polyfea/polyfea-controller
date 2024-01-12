var jt = function(e, t) {
  return jt = Object.setPrototypeOf || { __proto__: [] } instanceof Array && function(n, r) {
    n.__proto__ = r;
  } || function(n, r) {
    for (var i in r)
      Object.prototype.hasOwnProperty.call(r, i) && (n[i] = r[i]);
  }, jt(e, t);
};
function V(e, t) {
  if (typeof t != "function" && t !== null)
    throw new TypeError("Class extends value " + String(t) + " is not a constructor or null");
  jt(e, t);
  function n() {
    this.constructor = e;
  }
  e.prototype = t === null ? Object.create(t) : (n.prototype = t.prototype, new n());
}
function vr(e, t, n, r) {
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
    u((r = r.apply(e, t || [])).next());
  });
}
function sn(e, t) {
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
        u = t.call(e, n);
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
function X(e) {
  var t = typeof Symbol == "function" && Symbol.iterator, n = t && e[t], r = 0;
  if (n)
    return n.call(e);
  if (e && typeof e.length == "number")
    return {
      next: function() {
        return e && r >= e.length && (e = void 0), { value: e && e[r++], done: !e };
      }
    };
  throw new TypeError(t ? "Object is not iterable." : "Symbol.iterator is not defined.");
}
function H(e, t) {
  var n = typeof Symbol == "function" && e[Symbol.iterator];
  if (!n)
    return e;
  var r = n.call(e), i, s = [], o;
  try {
    for (; (t === void 0 || t-- > 0) && !(i = r.next()).done; )
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
function Z(e, t, n) {
  if (n || arguments.length === 2)
    for (var r = 0, i = t.length, s; r < i; r++)
      (s || !(r in t)) && (s || (s = Array.prototype.slice.call(t, 0, r)), s[r] = t[r]);
  return e.concat(s || Array.prototype.slice.call(t));
}
function Y(e) {
  return this instanceof Y ? (this.v = e, this) : new Y(e);
}
function $r(e, t, n) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var r = n.apply(e, t || []), i, s = [];
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
    c.value instanceof Y ? Promise.resolve(c.value.v).then(u, f) : h(s[0][2], c);
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
function wr(e) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var t = e[Symbol.asyncIterator], n;
  return t ? t.call(e) : (e = typeof X == "function" ? X(e) : e[Symbol.iterator](), n = {}, r("next"), r("throw"), r("return"), n[Symbol.asyncIterator] = function() {
    return this;
  }, n);
  function r(s) {
    n[s] = e[s] && function(o) {
      return new Promise(function(a, l) {
        o = e[s](o), i(a, l, o.done, o.value);
      });
    };
  }
  function i(s, o, a, l) {
    Promise.resolve(l).then(function(u) {
      s({ value: u, done: a });
    }, o);
  }
}
function g(e) {
  return typeof e == "function";
}
function ie(e) {
  var t = function(r) {
    Error.call(r), r.stack = new Error().stack;
  }, n = e(t);
  return n.prototype = Object.create(Error.prototype), n.prototype.constructor = n, n;
}
var Dt = ie(function(e) {
  return function(n) {
    e(this), this.message = n ? n.length + ` errors occurred during unsubscription:
` + n.map(function(r, i) {
      return i + 1 + ") " + r.toString();
    }).join(`
  `) : "", this.name = "UnsubscriptionError", this.errors = n;
  };
});
function Gt(e, t) {
  if (e) {
    var n = e.indexOf(t);
    0 <= n && e.splice(n, 1);
  }
}
var _t = function() {
  function e(t) {
    this.initialTeardown = t, this.closed = !1, this._parentage = null, this._finalizers = null;
  }
  return e.prototype.unsubscribe = function() {
    var t, n, r, i, s;
    if (!this.closed) {
      this.closed = !0;
      var o = this._parentage;
      if (o)
        if (this._parentage = null, Array.isArray(o))
          try {
            for (var a = X(o), l = a.next(); !l.done; l = a.next()) {
              var u = l.value;
              u.remove(this);
            }
          } catch (y) {
            t = { error: y };
          } finally {
            try {
              l && !l.done && (n = a.return) && n.call(a);
            } finally {
              if (t)
                throw t.error;
            }
          }
        else
          o.remove(this);
      var f = this.initialTeardown;
      if (g(f))
        try {
          f();
        } catch (y) {
          s = y instanceof Dt ? y.errors : [y];
        }
      var h = this._finalizers;
      if (h) {
        this._finalizers = null;
        try {
          for (var c = X(h), d = c.next(); !d.done; d = c.next()) {
            var p = d.value;
            try {
              Te(p);
            } catch (y) {
              s = s ?? [], y instanceof Dt ? s = Z(Z([], H(s)), H(y.errors)) : s.push(y);
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
        throw new Dt(s);
    }
  }, e.prototype.add = function(t) {
    var n;
    if (t && t !== this)
      if (this.closed)
        Te(t);
      else {
        if (t instanceof e) {
          if (t.closed || t._hasParent(this))
            return;
          t._addParent(this);
        }
        (this._finalizers = (n = this._finalizers) !== null && n !== void 0 ? n : []).push(t);
      }
  }, e.prototype._hasParent = function(t) {
    var n = this._parentage;
    return n === t || Array.isArray(n) && n.includes(t);
  }, e.prototype._addParent = function(t) {
    var n = this._parentage;
    this._parentage = Array.isArray(n) ? (n.push(t), n) : n ? [n, t] : t;
  }, e.prototype._removeParent = function(t) {
    var n = this._parentage;
    n === t ? this._parentage = null : Array.isArray(n) && Gt(n, t);
  }, e.prototype.remove = function(t) {
    var n = this._finalizers;
    n && Gt(n, t), t instanceof e && t._removeParent(this);
  }, e.EMPTY = function() {
    var t = new e();
    return t.closed = !0, t;
  }(), e;
}(), on = _t.EMPTY;
function an(e) {
  return e instanceof _t || e && "closed" in e && g(e.remove) && g(e.add) && g(e.unsubscribe);
}
function Te(e) {
  g(e) ? e() : e.unsubscribe();
}
var ln = {
  onUnhandledError: null,
  onStoppedNotification: null,
  Promise: void 0,
  useDeprecatedSynchronousErrorHandling: !1,
  useDeprecatedNextContext: !1
}, cn = {
  setTimeout: function(e, t) {
    for (var n = [], r = 2; r < arguments.length; r++)
      n[r - 2] = arguments[r];
    return setTimeout.apply(void 0, Z([e, t], H(n)));
  },
  clearTimeout: function(e) {
    var t = cn.delegate;
    return ((t == null ? void 0 : t.clearTimeout) || clearTimeout)(e);
  },
  delegate: void 0
};
function un(e) {
  cn.setTimeout(function() {
    throw e;
  });
}
function Yt() {
}
function dt(e) {
  e();
}
var se = function(e) {
  V(t, e);
  function t(n) {
    var r = e.call(this) || this;
    return r.isStopped = !1, n ? (r.destination = n, an(n) && n.add(r)) : r.destination = xr, r;
  }
  return t.create = function(n, r, i) {
    return new yt(n, r, i);
  }, t.prototype.next = function(n) {
    this.isStopped || this._next(n);
  }, t.prototype.error = function(n) {
    this.isStopped || (this.isStopped = !0, this._error(n));
  }, t.prototype.complete = function() {
    this.isStopped || (this.isStopped = !0, this._complete());
  }, t.prototype.unsubscribe = function() {
    this.closed || (this.isStopped = !0, e.prototype.unsubscribe.call(this), this.destination = null);
  }, t.prototype._next = function(n) {
    this.destination.next(n);
  }, t.prototype._error = function(n) {
    try {
      this.destination.error(n);
    } finally {
      this.unsubscribe();
    }
  }, t.prototype._complete = function() {
    try {
      this.destination.complete();
    } finally {
      this.unsubscribe();
    }
  }, t;
}(_t), Tr = Function.prototype.bind;
function Jt(e, t) {
  return Tr.call(e, t);
}
var Sr = function() {
  function e(t) {
    this.partialObserver = t;
  }
  return e.prototype.next = function(t) {
    var n = this.partialObserver;
    if (n.next)
      try {
        n.next(t);
      } catch (r) {
        lt(r);
      }
  }, e.prototype.error = function(t) {
    var n = this.partialObserver;
    if (n.error)
      try {
        n.error(t);
      } catch (r) {
        lt(r);
      }
    else
      lt(t);
  }, e.prototype.complete = function() {
    var t = this.partialObserver;
    if (t.complete)
      try {
        t.complete();
      } catch (n) {
        lt(n);
      }
  }, e;
}(), yt = function(e) {
  V(t, e);
  function t(n, r, i) {
    var s = e.call(this) || this, o;
    if (g(n) || !n)
      o = {
        next: n ?? void 0,
        error: r ?? void 0,
        complete: i ?? void 0
      };
    else {
      var a;
      s && ln.useDeprecatedNextContext ? (a = Object.create(n), a.unsubscribe = function() {
        return s.unsubscribe();
      }, o = {
        next: n.next && Jt(n.next, a),
        error: n.error && Jt(n.error, a),
        complete: n.complete && Jt(n.complete, a)
      }) : o = n;
    }
    return s.destination = new Sr(o), s;
  }
  return t;
}(se);
function lt(e) {
  un(e);
}
function Er(e) {
  throw e;
}
var xr = {
  closed: !0,
  next: Yt,
  error: Er,
  complete: Yt
}, oe = function() {
  return typeof Symbol == "function" && Symbol.observable || "@@observable";
}();
function Rt(e) {
  return e;
}
function _r(e) {
  return e.length === 0 ? Rt : e.length === 1 ? e[0] : function(n) {
    return e.reduce(function(r, i) {
      return i(r);
    }, n);
  };
}
var w = function() {
  function e(t) {
    t && (this._subscribe = t);
  }
  return e.prototype.lift = function(t) {
    var n = new e();
    return n.source = this, n.operator = t, n;
  }, e.prototype.subscribe = function(t, n, r) {
    var i = this, s = Ar(t) ? t : new yt(t, n, r);
    return dt(function() {
      var o = i, a = o.operator, l = o.source;
      s.add(a ? a.call(s, l) : l ? i._subscribe(s) : i._trySubscribe(s));
    }), s;
  }, e.prototype._trySubscribe = function(t) {
    try {
      return this._subscribe(t);
    } catch (n) {
      t.error(n);
    }
  }, e.prototype.forEach = function(t, n) {
    var r = this;
    return n = Se(n), new n(function(i, s) {
      var o = new yt({
        next: function(a) {
          try {
            t(a);
          } catch (l) {
            s(l), o.unsubscribe();
          }
        },
        error: s,
        complete: i
      });
      r.subscribe(o);
    });
  }, e.prototype._subscribe = function(t) {
    var n;
    return (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(t);
  }, e.prototype[oe] = function() {
    return this;
  }, e.prototype.pipe = function() {
    for (var t = [], n = 0; n < arguments.length; n++)
      t[n] = arguments[n];
    return _r(t)(this);
  }, e.prototype.toPromise = function(t) {
    var n = this;
    return t = Se(t), new t(function(r, i) {
      var s;
      n.subscribe(function(o) {
        return s = o;
      }, function(o) {
        return i(o);
      }, function() {
        return r(s);
      });
    });
  }, e.create = function(t) {
    return new e(t);
  }, e;
}();
function Se(e) {
  var t;
  return (t = e ?? ln.Promise) !== null && t !== void 0 ? t : Promise;
}
function Rr(e) {
  return e && g(e.next) && g(e.error) && g(e.complete);
}
function Ar(e) {
  return e && e instanceof se || Rr(e) && an(e);
}
function kr(e) {
  return g(e == null ? void 0 : e.lift);
}
function O(e) {
  return function(t) {
    if (kr(t))
      return t.lift(function(n) {
        try {
          return e(n, this);
        } catch (r) {
          this.error(r);
        }
      });
    throw new TypeError("Unable to lift unknown Observable type");
  };
}
function M(e, t, n, r, i) {
  return new Ir(e, t, n, r, i);
}
var Ir = function(e) {
  V(t, e);
  function t(n, r, i, s, o, a) {
    var l = e.call(this, n) || this;
    return l.onFinalize = o, l.shouldUnsubscribe = a, l._next = r ? function(u) {
      try {
        r(u);
      } catch (f) {
        n.error(f);
      }
    } : e.prototype._next, l._error = s ? function(u) {
      try {
        s(u);
      } catch (f) {
        n.error(f);
      } finally {
        this.unsubscribe();
      }
    } : e.prototype._error, l._complete = i ? function() {
      try {
        i();
      } catch (u) {
        n.error(u);
      } finally {
        this.unsubscribe();
      }
    } : e.prototype._complete, l;
  }
  return t.prototype.unsubscribe = function() {
    var n;
    if (!this.shouldUnsubscribe || this.shouldUnsubscribe()) {
      var r = this.closed;
      e.prototype.unsubscribe.call(this), !r && ((n = this.onFinalize) === null || n === void 0 || n.call(this));
    }
  }, t;
}(se), Or = ie(function(e) {
  return function() {
    e(this), this.name = "ObjectUnsubscribedError", this.message = "object unsubscribed";
  };
}), ae = function(e) {
  V(t, e);
  function t() {
    var n = e.call(this) || this;
    return n.closed = !1, n.currentObservers = null, n.observers = [], n.isStopped = !1, n.hasError = !1, n.thrownError = null, n;
  }
  return t.prototype.lift = function(n) {
    var r = new Ee(this, this);
    return r.operator = n, r;
  }, t.prototype._throwIfClosed = function() {
    if (this.closed)
      throw new Or();
  }, t.prototype.next = function(n) {
    var r = this;
    dt(function() {
      var i, s;
      if (r._throwIfClosed(), !r.isStopped) {
        r.currentObservers || (r.currentObservers = Array.from(r.observers));
        try {
          for (var o = X(r.currentObservers), a = o.next(); !a.done; a = o.next()) {
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
  }, t.prototype.error = function(n) {
    var r = this;
    dt(function() {
      if (r._throwIfClosed(), !r.isStopped) {
        r.hasError = r.isStopped = !0, r.thrownError = n;
        for (var i = r.observers; i.length; )
          i.shift().error(n);
      }
    });
  }, t.prototype.complete = function() {
    var n = this;
    dt(function() {
      if (n._throwIfClosed(), !n.isStopped) {
        n.isStopped = !0;
        for (var r = n.observers; r.length; )
          r.shift().complete();
      }
    });
  }, t.prototype.unsubscribe = function() {
    this.isStopped = this.closed = !0, this.observers = this.currentObservers = null;
  }, Object.defineProperty(t.prototype, "observed", {
    get: function() {
      var n;
      return ((n = this.observers) === null || n === void 0 ? void 0 : n.length) > 0;
    },
    enumerable: !1,
    configurable: !0
  }), t.prototype._trySubscribe = function(n) {
    return this._throwIfClosed(), e.prototype._trySubscribe.call(this, n);
  }, t.prototype._subscribe = function(n) {
    return this._throwIfClosed(), this._checkFinalizedStatuses(n), this._innerSubscribe(n);
  }, t.prototype._innerSubscribe = function(n) {
    var r = this, i = this, s = i.hasError, o = i.isStopped, a = i.observers;
    return s || o ? on : (this.currentObservers = null, a.push(n), new _t(function() {
      r.currentObservers = null, Gt(a, n);
    }));
  }, t.prototype._checkFinalizedStatuses = function(n) {
    var r = this, i = r.hasError, s = r.thrownError, o = r.isStopped;
    i ? n.error(s) : o && n.complete();
  }, t.prototype.asObservable = function() {
    var n = new w();
    return n.source = this, n;
  }, t.create = function(n, r) {
    return new Ee(n, r);
  }, t;
}(w), Ee = function(e) {
  V(t, e);
  function t(n, r) {
    var i = e.call(this) || this;
    return i.destination = n, i.source = r, i;
  }
  return t.prototype.next = function(n) {
    var r, i;
    (i = (r = this.destination) === null || r === void 0 ? void 0 : r.next) === null || i === void 0 || i.call(r, n);
  }, t.prototype.error = function(n) {
    var r, i;
    (i = (r = this.destination) === null || r === void 0 ? void 0 : r.error) === null || i === void 0 || i.call(r, n);
  }, t.prototype.complete = function() {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.complete) === null || r === void 0 || r.call(n);
  }, t.prototype._subscribe = function(n) {
    var r, i;
    return (i = (r = this.source) === null || r === void 0 ? void 0 : r.subscribe(n)) !== null && i !== void 0 ? i : on;
  }, t;
}(ae), ct = function(e) {
  V(t, e);
  function t() {
    var n = e !== null && e.apply(this, arguments) || this;
    return n._value = null, n._hasValue = !1, n._isComplete = !1, n;
  }
  return t.prototype._checkFinalizedStatuses = function(n) {
    var r = this, i = r.hasError, s = r._hasValue, o = r._value, a = r.thrownError, l = r.isStopped, u = r._isComplete;
    i ? n.error(a) : (l || u) && (s && n.next(o), n.complete());
  }, t.prototype.next = function(n) {
    this.isStopped || (this._value = n, this._hasValue = !0);
  }, t.prototype.complete = function() {
    var n = this, r = n._hasValue, i = n._value, s = n._isComplete;
    s || (this._isComplete = !0, r && e.prototype.next.call(this, i), e.prototype.complete.call(this));
  }, t;
}(ae);
function Cr(e) {
  return e && g(e.schedule);
}
function Lr(e) {
  return e[e.length - 1];
}
function At(e) {
  return Cr(Lr(e)) ? e.pop() : void 0;
}
var le = function(e) {
  return e && typeof e.length == "number" && typeof e != "function";
};
function fn(e) {
  return g(e == null ? void 0 : e.then);
}
function hn(e) {
  return g(e[oe]);
}
function dn(e) {
  return Symbol.asyncIterator && g(e == null ? void 0 : e[Symbol.asyncIterator]);
}
function pn(e) {
  return new TypeError("You provided " + (e !== null && typeof e == "object" ? "an invalid object" : "'" + e + "'") + " where a stream was expected. You can provide an Observable, Promise, ReadableStream, Array, AsyncIterable, or Iterable.");
}
function Pr() {
  return typeof Symbol != "function" || !Symbol.iterator ? "@@iterator" : Symbol.iterator;
}
var yn = Pr();
function mn(e) {
  return g(e == null ? void 0 : e[yn]);
}
function gn(e) {
  return $r(this, arguments, function() {
    var n, r, i, s;
    return sn(this, function(o) {
      switch (o.label) {
        case 0:
          n = e.getReader(), o.label = 1;
        case 1:
          o.trys.push([1, , 9, 10]), o.label = 2;
        case 2:
          return [4, Y(n.read())];
        case 3:
          return r = o.sent(), i = r.value, s = r.done, s ? [4, Y(void 0)] : [3, 5];
        case 4:
          return [2, o.sent()];
        case 5:
          return [4, Y(i)];
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
function bn(e) {
  return g(e == null ? void 0 : e.getReader);
}
function N(e) {
  if (e instanceof w)
    return e;
  if (e != null) {
    if (hn(e))
      return Ur(e);
    if (le(e))
      return Mr(e);
    if (fn(e))
      return Fr(e);
    if (dn(e))
      return vn(e);
    if (mn(e))
      return Nr(e);
    if (bn(e))
      return Dr(e);
  }
  throw pn(e);
}
function Ur(e) {
  return new w(function(t) {
    var n = e[oe]();
    if (g(n.subscribe))
      return n.subscribe(t);
    throw new TypeError("Provided object does not correctly implement Symbol.observable");
  });
}
function Mr(e) {
  return new w(function(t) {
    for (var n = 0; n < e.length && !t.closed; n++)
      t.next(e[n]);
    t.complete();
  });
}
function Fr(e) {
  return new w(function(t) {
    e.then(function(n) {
      t.closed || (t.next(n), t.complete());
    }, function(n) {
      return t.error(n);
    }).then(null, un);
  });
}
function Nr(e) {
  return new w(function(t) {
    var n, r;
    try {
      for (var i = X(e), s = i.next(); !s.done; s = i.next()) {
        var o = s.value;
        if (t.next(o), t.closed)
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
    t.complete();
  });
}
function vn(e) {
  return new w(function(t) {
    Jr(e, t).catch(function(n) {
      return t.error(n);
    });
  });
}
function Dr(e) {
  return vn(gn(e));
}
function Jr(e, t) {
  var n, r, i, s;
  return vr(this, void 0, void 0, function() {
    var o, a;
    return sn(this, function(l) {
      switch (l.label) {
        case 0:
          l.trys.push([0, 5, 6, 11]), n = wr(e), l.label = 1;
        case 1:
          return [4, n.next()];
        case 2:
          if (r = l.sent(), !!r.done)
            return [3, 4];
          if (o = r.value, t.next(o), t.closed)
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
          return t.complete(), [2];
      }
    });
  });
}
function P(e, t, n, r, i) {
  r === void 0 && (r = 0), i === void 0 && (i = !1);
  var s = t.schedule(function() {
    n(), i ? e.add(this.schedule(null, r)) : this.unsubscribe();
  }, r);
  if (e.add(s), !i)
    return s;
}
function $n(e, t) {
  return t === void 0 && (t = 0), O(function(n, r) {
    n.subscribe(M(r, function(i) {
      return P(r, e, function() {
        return r.next(i);
      }, t);
    }, function() {
      return P(r, e, function() {
        return r.complete();
      }, t);
    }, function(i) {
      return P(r, e, function() {
        return r.error(i);
      }, t);
    }));
  });
}
function wn(e, t) {
  return t === void 0 && (t = 0), O(function(n, r) {
    r.add(e.schedule(function() {
      return n.subscribe(r);
    }, t));
  });
}
function Hr(e, t) {
  return N(e).pipe(wn(t), $n(t));
}
function Br(e, t) {
  return N(e).pipe(wn(t), $n(t));
}
function zr(e, t) {
  return new w(function(n) {
    var r = 0;
    return t.schedule(function() {
      r === e.length ? n.complete() : (n.next(e[r++]), n.closed || this.schedule());
    });
  });
}
function Kr(e, t) {
  return new w(function(n) {
    var r;
    return P(n, t, function() {
      r = e[yn](), P(n, t, function() {
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
      return g(r == null ? void 0 : r.return) && r.return();
    };
  });
}
function Tn(e, t) {
  if (!e)
    throw new Error("Iterable cannot be null");
  return new w(function(n) {
    P(n, t, function() {
      var r = e[Symbol.asyncIterator]();
      P(n, t, function() {
        r.next().then(function(i) {
          i.done ? n.complete() : n.next(i.value);
        });
      }, 0, !0);
    });
  });
}
function Wr(e, t) {
  return Tn(gn(e), t);
}
function jr(e, t) {
  if (e != null) {
    if (hn(e))
      return Hr(e, t);
    if (le(e))
      return zr(e, t);
    if (fn(e))
      return Br(e, t);
    if (dn(e))
      return Tn(e, t);
    if (mn(e))
      return Kr(e, t);
    if (bn(e))
      return Wr(e, t);
  }
  throw pn(e);
}
function kt(e, t) {
  return t ? jr(e, t) : N(e);
}
function Gr() {
  for (var e = [], t = 0; t < arguments.length; t++)
    e[t] = arguments[t];
  var n = At(e);
  return kt(e, n);
}
function Yr(e, t) {
  var n = g(e) ? e : function() {
    return e;
  }, r = function(i) {
    return i.error(n());
  };
  return new w(t ? function(i) {
    return t.schedule(r, 0, i);
  } : r);
}
var qr = ie(function(e) {
  return function() {
    e(this), this.name = "EmptyError", this.message = "no elements in sequence";
  };
});
function qt(e, t) {
  var n = typeof t == "object";
  return new Promise(function(r, i) {
    var s = new yt({
      next: function(o) {
        r(o), s.unsubscribe();
      },
      error: i,
      complete: function() {
        n ? r(t.defaultValue) : i(new qr());
      }
    });
    e.subscribe(s);
  });
}
function ce(e, t) {
  return O(function(n, r) {
    var i = 0;
    n.subscribe(M(r, function(s) {
      r.next(e.call(t, s, i++));
    }));
  });
}
var Qr = Array.isArray;
function Xr(e, t) {
  return Qr(t) ? e.apply(void 0, Z([], H(t))) : e(t);
}
function Zr(e) {
  return ce(function(t) {
    return Xr(e, t);
  });
}
function Vr(e, t, n, r, i, s, o, a) {
  var l = [], u = 0, f = 0, h = !1, c = function() {
    h && !l.length && !u && t.complete();
  }, d = function(y) {
    return u < r ? p(y) : l.push(y);
  }, p = function(y) {
    s && t.next(y), u++;
    var v = !1;
    N(n(y, f++)).subscribe(M(t, function(m) {
      i == null || i(m), s ? d(m) : t.next(m);
    }, function() {
      v = !0;
    }, void 0, function() {
      if (v)
        try {
          u--;
          for (var m = function() {
            var R = l.shift();
            o ? P(t, o, function() {
              return p(R);
            }) : p(R);
          }; l.length && u < r; )
            m();
          c();
        } catch (R) {
          t.error(R);
        }
    }));
  };
  return e.subscribe(M(t, d, function() {
    h = !0, c();
  })), function() {
    a == null || a();
  };
}
function ue(e, t, n) {
  return n === void 0 && (n = 1 / 0), g(t) ? ue(function(r, i) {
    return ce(function(s, o) {
      return t(r, s, i, o);
    })(N(e(r, i)));
  }, n) : (typeof t == "number" && (n = t), O(function(r, i) {
    return Vr(r, i, e, n);
  }));
}
function ti(e) {
  return e === void 0 && (e = 1 / 0), ue(Rt, e);
}
function Sn() {
  return ti(1);
}
function xe() {
  for (var e = [], t = 0; t < arguments.length; t++)
    e[t] = arguments[t];
  return Sn()(kt(e, At(e)));
}
function ei(e) {
  return new w(function(t) {
    N(e()).subscribe(t);
  });
}
var ni = ["addListener", "removeListener"], ri = ["addEventListener", "removeEventListener"], ii = ["on", "off"];
function Qt(e, t, n, r) {
  if (g(n) && (r = n, n = void 0), r)
    return Qt(e, t, n).pipe(Zr(r));
  var i = H(ai(e) ? ri.map(function(a) {
    return function(l) {
      return e[a](t, l, n);
    };
  }) : si(e) ? ni.map(_e(e, t)) : oi(e) ? ii.map(_e(e, t)) : [], 2), s = i[0], o = i[1];
  if (!s && le(e))
    return ue(function(a) {
      return Qt(a, t, n);
    })(N(e));
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
function _e(e, t) {
  return function(n) {
    return function(r) {
      return e[n](t, r);
    };
  };
}
function si(e) {
  return g(e.addListener) && g(e.removeListener);
}
function oi(e) {
  return g(e.on) && g(e.off);
}
function ai(e) {
  return g(e.addEventListener) && g(e.removeEventListener);
}
var li = new w(Yt);
function ci() {
  for (var e = [], t = 0; t < arguments.length; t++)
    e[t] = arguments[t];
  var n = At(e);
  return O(function(r, i) {
    Sn()(kt(Z([r], H(e)), n)).subscribe(i);
  });
}
function En() {
  for (var e = [], t = 0; t < arguments.length; t++)
    e[t] = arguments[t];
  return ci.apply(void 0, Z([], H(e)));
}
function Re(e, t) {
  return t === void 0 && (t = Rt), e = e ?? ui, O(function(n, r) {
    var i, s = !0;
    n.subscribe(M(r, function(o) {
      var a = t(o);
      (s || !e(i, a)) && (s = !1, i = a, r.next(o));
    }));
  });
}
function ui(e, t) {
  return e === t;
}
function fi() {
  for (var e = [], t = 0; t < arguments.length; t++)
    e[t] = arguments[t];
  var n = At(e);
  return O(function(r, i) {
    (n ? xe(e, r, n) : xe(e, r)).subscribe(i);
  });
}
function xn(e, t) {
  return O(function(n, r) {
    var i = null, s = 0, o = !1, a = function() {
      return o && !i && r.complete();
    };
    n.subscribe(M(r, function(l) {
      i == null || i.unsubscribe();
      var u = 0, f = s++;
      N(e(l, f)).subscribe(i = M(r, function(h) {
        return r.next(t ? t(l, h, f, u++) : h);
      }, function() {
        i = null, a();
      }));
    }, function() {
      o = !0, a();
    }));
  });
}
function hi(e, t, n) {
  var r = g(e) || t || n ? { next: e, error: t, complete: n } : e;
  return r ? O(function(i, s) {
    var o;
    (o = r.subscribe) === null || o === void 0 || o.call(r);
    var a = !0;
    i.subscribe(M(s, function(l) {
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
  }) : Rt;
}
const di = "http://./polyfea".replace(/\/+$/, "");
let nt = class {
  constructor(t = {}) {
    this.configuration = t;
  }
  set config(t) {
    this.configuration = t;
  }
  get basePath() {
    return this.configuration.basePath != null ? this.configuration.basePath : di;
  }
  get fetchApi() {
    return this.configuration.fetchApi;
  }
  get middleware() {
    return this.configuration.middleware || [];
  }
  get queryParamsStringify() {
    return this.configuration.queryParamsStringify || An;
  }
  get username() {
    return this.configuration.username;
  }
  get password() {
    return this.configuration.password;
  }
  get apiKey() {
    const t = this.configuration.apiKey;
    if (t)
      return typeof t == "function" ? t : () => t;
  }
  get accessToken() {
    const t = this.configuration.accessToken;
    if (t)
      return typeof t == "function" ? t : async () => t;
  }
  get headers() {
    return this.configuration.headers;
  }
  get credentials() {
    return this.configuration.credentials;
  }
};
const pi = new nt();
let _n = class Rn {
  constructor(t = pi) {
    this.configuration = t, this.fetchApi = async (n, r) => {
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
          throw o instanceof Error ? new bi(o, "The request failed and the interceptors did not return an alternative response") : o;
      }
      for (const o of this.middleware)
        o.post && (s = await o.post({
          fetch: this.fetchApi,
          url: i.url,
          init: i.init,
          response: s.clone()
        }) || s);
      return s;
    }, this.middleware = t.middleware;
  }
  withMiddleware(...t) {
    const n = this.clone();
    return n.middleware = n.middleware.concat(...t), n;
  }
  withPreMiddleware(...t) {
    const n = t.map((r) => ({ pre: r }));
    return this.withMiddleware(...n);
  }
  withPostMiddleware(...t) {
    const n = t.map((r) => ({ post: r }));
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
  isJsonMime(t) {
    return t ? Rn.jsonRegex.test(t) : !1;
  }
  async request(t, n) {
    const { url: r, init: i } = await this.createFetchParams(t, n), s = await this.fetchApi(r, i);
    if (s && s.status >= 200 && s.status < 300)
      return s;
    throw new gi(s, "Response returned an error code");
  }
  async createFetchParams(t, n) {
    let r = this.configuration.basePath + t.path;
    t.query !== void 0 && Object.keys(t.query).length !== 0 && (r += "?" + this.configuration.queryParamsStringify(t.query));
    const i = Object.assign({}, this.configuration.headers, t.headers);
    Object.keys(i).forEach((f) => i[f] === void 0 ? delete i[f] : {});
    const s = typeof n == "function" ? n : async () => n, o = {
      method: t.method,
      headers: i,
      body: t.body,
      credentials: this.configuration.credentials
    }, a = Object.assign(Object.assign({}, o), await s({
      init: o,
      context: t
    }));
    let l;
    mi(a.body) || a.body instanceof URLSearchParams || yi(a.body) ? l = a.body : this.isJsonMime(i["Content-Type"]) ? l = JSON.stringify(a.body) : l = a.body;
    const u = Object.assign(Object.assign({}, a), { body: l });
    return { url: r, init: u };
  }
  /**
   * Create a shallow clone of `this` by constructing a new instance
   * and then shallow cloning data members.
   */
  clone() {
    const t = this.constructor, n = new t(this.configuration);
    return n.middleware = this.middleware.slice(), n;
  }
};
_n.jsonRegex = new RegExp("^(:?application/json|[^;/ 	]+/[^;/ 	]+[+]json)[ 	]*(:?;.*)?$", "i");
function yi(e) {
  return typeof Blob < "u" && e instanceof Blob;
}
function mi(e) {
  return typeof FormData < "u" && e instanceof FormData;
}
let gi = class extends Error {
  constructor(t, n) {
    super(n), this.response = t, this.name = "ResponseError";
  }
}, bi = class extends Error {
  constructor(t, n) {
    super(n), this.cause = t, this.name = "FetchError";
  }
}, Ae = class extends Error {
  constructor(t, n) {
    super(n), this.field = t, this.name = "RequiredError";
  }
};
function S(e, t) {
  const n = e[t];
  return n != null;
}
function An(e, t = "") {
  return Object.keys(e).map((n) => kn(n, e[n], t)).filter((n) => n.length > 0).join("&");
}
function kn(e, t, n = "") {
  const r = n + (n.length ? `[${e}]` : e);
  if (t instanceof Array) {
    const i = t.map((s) => encodeURIComponent(String(s))).join(`&${encodeURIComponent(r)}=`);
    return `${encodeURIComponent(r)}=${i}`;
  }
  if (t instanceof Set) {
    const i = Array.from(t);
    return kn(e, i, n);
  }
  return t instanceof Date ? `${encodeURIComponent(r)}=${encodeURIComponent(t.toISOString())}` : t instanceof Object ? An(t, r) : `${encodeURIComponent(r)}=${encodeURIComponent(String(t))}`;
}
function In(e, t) {
  return Object.keys(e).reduce((n, r) => Object.assign(Object.assign({}, n), { [r]: t(e[r]) }), {});
}
let ke = class {
  constructor(t, n = (r) => r) {
    this.raw = t, this.transformer = n;
  }
  async value() {
    return this.transformer(await this.raw.json());
  }
};
function vi(e) {
  return $i(e);
}
function $i(e, t) {
  return e == null ? e : {
    microfrontend: S(e, "microfrontend") ? e.microfrontend : void 0,
    tagName: e.tagName,
    attributes: S(e, "attributes") ? e.attributes : void 0,
    style: S(e, "style") ? e.style : void 0
  };
}
function wi(e) {
  return Ti(e);
}
function Ti(e, t) {
  return e == null ? e : {
    kind: S(e, "kind") ? e.kind : void 0,
    href: S(e, "href") ? e.href : void 0,
    attributes: S(e, "attributes") ? e.attributes : void 0,
    waitOnLoad: S(e, "waitOnLoad") ? e.waitOnLoad : void 0
  };
}
function On(e) {
  return Si(e);
}
function Si(e, t) {
  return e == null ? e : {
    dependsOn: S(e, "dependsOn") ? e.dependsOn : void 0,
    module: S(e, "module") ? e.module : void 0,
    resources: S(e, "resources") ? e.resources.map(wi) : void 0
  };
}
function Cn(e) {
  return Ei(e);
}
function Ei(e, t) {
  return e == null ? e : {
    elements: e.elements.map(vi),
    microfrontends: S(e, "microfrontends") ? In(e.microfrontends, On) : void 0
  };
}
function xi(e) {
  return _i(e);
}
function _i(e, t) {
  return e == null ? e : {
    name: e.name,
    path: S(e, "path") ? e.path : void 0,
    contextArea: S(e, "contextArea") ? Cn(e.contextArea) : void 0
  };
}
function Ri(e) {
  return Ai(e);
}
function Ai(e, t) {
  return e == null ? e : {
    contextAreas: S(e, "contextAreas") ? e.contextAreas.map(xi) : void 0,
    microfrontends: In(e.microfrontends, On)
  };
}
let mt = class extends _n {
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextAreaRaw(t, n) {
    if (t.name === null || t.name === void 0)
      throw new Ae("name", "Required parameter requestParameters.name was null or undefined when calling getContextArea.");
    if (t.path === null || t.path === void 0)
      throw new Ae("path", "Required parameter requestParameters.path was null or undefined when calling getContextArea.");
    const r = {};
    t.path !== void 0 && (r.path = t.path), t.take !== void 0 && (r.take = t.take);
    const i = {}, s = await this.request({
      path: "/context-area/{name}".replace("{name}", encodeURIComponent(String(t.name))),
      method: "GET",
      headers: i,
      query: r
    }, n);
    return new ke(s, (o) => Cn(o));
  }
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextArea(t, n) {
    return await (await this.getContextAreaRaw(t, n)).value();
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfigRaw(t) {
    const n = {}, r = {}, i = await this.request({
      path: "/static-config",
      method: "GET",
      headers: r,
      query: n
    }, t);
    return new ke(i, (s) => Ri(s));
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfig(t) {
    return await (await this.getStaticConfigRaw(t)).value();
  }
}, ki = class {
  constructor(t = "./polyfea") {
    this.spec$ = new w();
    let n;
    typeof t == "string" ? (t.length === 0 && (t = "./polyfea"), n = new mt(new nt({ basePath: t }))) : t instanceof nt ? n = new mt(t) : n = t, this.spec$ = kt(n.getStaticConfig()).pipe(En(li));
  }
  getContextArea(t) {
    let n = globalThis.location.pathname;
    if (globalThis.document.baseURI) {
      const i = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(i) && (n = "./" + n.substring(i.length));
    }
    return this.spec$.pipe(
      ce((r) => {
        for (let i of r.contextAreas)
          if (i.name === t && new RegExp(i.path).test(n))
            return { ...i.contextArea, microfrontends: { ...r.microfrontends, ...i.contextArea.microfrontends } };
        return null;
      })
    );
  }
}, Ie = class {
  constructor(t = "./polyfea") {
    var n, r;
    typeof t == "string" ? this.api = new mt(new nt({
      basePath: new URL(
        t,
        new URL(
          ((n = globalThis.document) == null ? void 0 : n.baseURI) || "/",
          ((r = globalThis.location) == null ? void 0 : r.href) || "http://localhost"
        )
      ).href
    })) : t instanceof nt ? this.api = new mt(t) : this.api = t;
  }
  getContextArea(t) {
    var s, o;
    let n = ((s = globalThis.location) == null ? void 0 : s.pathname) || "/";
    if ((o = globalThis.document) != null && o.baseURI) {
      const l = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(l) && (n = "./" + n.substring(l.length));
    }
    const r = localStorage.getItem(`polyfea-context[${t},${n}]`), i = ei(() => this.api.getContextAreaRaw({ name: t, path: n })).pipe(
      xn((a) => a.raw.ok ? a.value() : Yr(() => new Error(a.raw.statusText))),
      hi((a) => {
        a && localStorage.setItem(`polyfea-context[${t},${n}]`, JSON.stringify(a));
      })
    );
    if (r) {
      const a = JSON.parse(r);
      return Gr(a).pipe(
        En(i)
      );
    } else
      return i;
  }
}, Ii = class {
  /**
   * Constructs a new NavigationDestination instance.
   * @param url - The URL of the navigation destination.
   */
  constructor(t) {
    this.url = t;
  }
}, Oi = class extends Event {
  constructor(t, n) {
    super("navigate", { bubbles: !0, cancelable: !0 }), this.transition = t, this.interceptPromises = [], this.downloadRequest = null, this.formData = null, this.hashChange = !1, this.userInitiated = !1;
    let r = new URL(t.href, new URL(globalThis.document.baseURI, globalThis.location.href));
    const i = new URL(n.url, new URL(globalThis.document.baseURI, globalThis.location.href));
    this.canIntercept = i.protocol === r.protocol && i.host === r.host && i.port === r.port, this.destination = new Ii(r.href);
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
  intercept(t) {
    t != null && t.handler && this.interceptPromises.push(t.handler(this));
  }
}, Oe = class extends Event {
  constructor(t, n) {
    super("currententrychange", { bubbles: !0, cancelable: !0 }), this.navigationType = t, this.from = n;
  }
}, Ce = class {
  constructor(t) {
    this.request = t;
  }
  get finished() {
    return qt(this.request.finished);
  }
  get from() {
    return this.request.entry;
  }
  get type() {
    return this.request.mode;
  }
};
function Ci(e = !1) {
  return Pi.tryRegister(e);
}
let Li = class {
  constructor(t) {
    this.commited$ = t.committed, this.finished$ = t.finished;
  }
  get commited() {
    return qt(this.commited$);
  }
  get finished() {
    return qt(this.finished$);
  }
}, Pi = class Xt extends EventTarget {
  constructor() {
    super(), this.entriesList = [], this.idCounter = 0, this.transitionRequests = new ae(), this.currentTransition = null, this.currentEntryIndex = -1, this.pushstateDelay = 35, this.rawHistoryMethods = {
      pushState: globalThis.history.pushState,
      replaceState: globalThis.history.replaceState,
      go: globalThis.history.go,
      back: globalThis.history.back,
      forward: globalThis.history.forward
    }, this.transitionRequests.subscribe((t) => this.executeRequest(t));
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
    var t;
    return ((t = this.currentTransition) == null ? void 0 : t.transition) || null;
  }
  navigate(t, n) {
    let r = "push";
    return ((n == null ? void 0 : n.history) === "replace" || n != null && n.replace) && (r = "replace"), this.nextTransitionRequest(r, t, n);
  }
  back() {
    if (this.currentEntryIndex < 1)
      throw { name: "InvaliStateError", message: "Cannot go back from initial state" };
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { traverseTo: this.entriesList[this.currentEntryIndex - 1].key });
  }
  forward(t) {
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { info: t, traverseTo: this.entriesList[this.currentEntryIndex + 1].key });
  }
  reload(t) {
    return this.nextTransitionRequest("reload", this.currentEntry.url, { info: t });
  }
  traverseTo(t, n) {
    const r = this.entriesList.find((i) => i.key === t);
    if (!r)
      throw { name: "InvaliStateError", message: "Cannot traverse to unknown state" };
    return this.nextTransitionRequest("traverse", r.url, { info: n == null ? void 0 : n.info, traverseTo: t });
  }
  updateCurrentEntry(t) {
    this.entriesList[this.currentEntryIndex].setState(JSON.parse(JSON.stringify(t == null ? void 0 : t.state))), this.currentEntry.dispatchEvent(new Oe("replace", this.currentEntry));
  }
  nextTransitionRequest(t, n, r) {
    const i = `@${++this.idCounter}-navigation-polyfill-transition`, s = {
      mode: t,
      href: new URL(n, new URL(globalThis.document.baseURI, globalThis.location.href)).href,
      info: r == null ? void 0 : r.info,
      state: r == null ? void 0 : r.state,
      committed: new ct(),
      finished: new ct(),
      entry: new Ht(this, i, i, n.toString(), r == null ? void 0 : r.state),
      abortController: new AbortController(),
      traverseToKey: r == null ? void 0 : r.traverseTo,
      transition: null
    };
    return this.transitionRequests.next(s), new Li(s);
  }
  async executeRequest(t) {
    this.currentTransition && (this.currentTransition.abortController.abort(), this.currentTransition.finished.error("aborted - new navigation started"), this.currentTransition.committed.closed || this.currentTransition.committed.error("aborted - new navigation started"), globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: Error("aborted - new navigation started") })), this.clearTransition(t)), t.transition = new Ce(t), this.currentTransition = t;
    try {
      await this.commit(t), this.currentEntry.dispatchEvent(new Oe(t.mode, t.transition.from)), await this.dispatchNavigation(t);
    } catch {
      t.finished.error("aborted"), t.committed.error("aborted");
    } finally {
      this.clearTransition(t);
    }
  }
  dispatchNavigation(t) {
    var r;
    const n = new Oi(t, this.currentEntry);
    if ((r = globalThis.navigation) != null && r.dispatchEvent(n))
      return n.interceptPromises.length > 0 ? Promise.all(n.interceptPromises.filter((i) => !!(i != null && i.then))).then(() => {
        globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(t), t.finished.next(), t.finished.complete();
      }).catch((i) => {
        globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: i })), this.clearTransition(t), t.finished.error(i);
      }) : (globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(t), t.finished.next(), t.finished.complete(), Promise.resolve());
  }
  clearTransition(t) {
    var n;
    ((n = this.currentTransition) == null ? void 0 : n.entry.id) === t.entry.id && (this.currentTransition = null);
  }
  commit(t) {
    switch (t.mode) {
      case "push":
        return this.commitPushTransition(t);
      case "replace":
        return this.commitReplaceTransition(t);
      case "reload":
        return this.commitReloadTransition(t);
      case "traverse":
        return this.commitTraverseTransition(t);
    }
  }
  async pushstateAsync(t, n = () => {
  }) {
    return new Promise((r, i) => {
      setTimeout(
        () => {
          n(t), t.committed.next(), t.committed.complete(), r();
        },
        this.pushstateDelay
      );
    });
  }
  commitPushTransition(t) {
    return this.rawHistoryMethods.pushState.apply(globalThis.history, [t.entry.cloneable, "", t.href]), this.pushstateAsync(t, (n) => {
      this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), n.entry];
    });
  }
  commitReplaceTransition(t) {
    return t.entry.key = this.currentEntry.key, this.entriesList[this.currentEntryIndex] = t.entry, this.rawHistoryMethods.replaceState.apply(globalThis.history, [t.entry.cloneable, "", t.href]), this.pushstateAsync(t);
  }
  commitTraverseTransition(t) {
    return new Promise(async (n, r) => {
      const i = this.entriesList.findIndex((o) => o.key === t.traverseToKey);
      i < 0 && r("target entry not found");
      const s = i - this.currentEntryIndex;
      this.rawHistoryMethods.go.apply(globalThis.history, [s]), await this.pushstateAsync(t, (o) => {
        const a = this.entriesList.findIndex((l) => l.key === o.traverseToKey);
        a < 0 && o.committed.error(new Error("target entry not found")), this.currentEntryIndex = a, o.committed.next(), o.committed.complete(), n();
      });
    });
  }
  commitReloadTransition(t) {
    return t.committed.next(), t.committed.complete(), t.finished.subscribe({
      next: () => globalThis.location.reload(),
      error: () => globalThis.location.reload()
    }), Promise.resolve();
  }
  static tryRegister(t = !1) {
    if (!globalThis.navigation) {
      const n = new Xt();
      return n.doRegister(t), n;
    }
    return globalThis.navigation;
  }
  static unregister() {
    globalThis.navigation && globalThis.navigation instanceof Xt && (globalThis.navigation.doUnregister && globalThis.navigation.doUnregister(), globalThis.navigation = void 0);
  }
  doRegister(t) {
    var n, r, i, s, o;
    if (!globalThis.navigation && !globalThis.navigation) {
      globalThis.navigation = this, this.entriesList = [new Ht(this, "initial", "initial", globalThis.location.href, void 0)], this.currentEntryIndex = 0;
      const a = ((n = globalThis.history) == null ? void 0 : n.pushState) || ((c, d, p) => {
      }), l = ((r = globalThis.history) == null ? void 0 : r.replaceState) || ((c, d, p) => {
      }), u = ((i = globalThis.history) == null ? void 0 : i.go) || ((c) => {
      }), f = ((s = globalThis.history) == null ? void 0 : s.back) || (() => {
      }), h = ((o = globalThis.history) == null ? void 0 : o.forward) || (() => {
      });
      this.doUnregister = () => {
        globalThis.history.pushState = a, globalThis.history.replaceState = l, globalThis.history.go = u, globalThis.history.back = f, globalThis.history.forward = h;
      }, t ? (this.rawHistoryMethods.pushState = () => {
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
          var m, R, z, tt, L, at, K;
          if (this.currentTransition && ((m = c.state) == null ? void 0 : m.id) === ((z = (R = this.currentTransition) == null ? void 0 : R.entry) == null ? void 0 : z.id))
            return;
          (tt = this.currentTransition) == null || tt.abortController.abort();
          const d = new ct();
          d.complete();
          let p;
          if ((L = c.state) != null && L.key) {
            const x = this.entriesList.findIndex((Nt) => {
              var W;
              return Nt.key === ((W = c.state) == null ? void 0 : W.key);
            });
            x >= 0 && (this.currentEntryIndex = x, p = this.entriesList[x]);
          }
          if (!p) {
            let x = `@${++this.idCounter}-navigation-polyfill-popstate`;
            p = new Ht(this, x, x, globalThis.location.href, c.state), this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), p];
          }
          const y = new ct(), v = {
            mode: "traverse",
            href: globalThis.location.href,
            info: void 0,
            state: ((at = c.state) == null ? void 0 : at.state) || c.state,
            committed: d,
            finished: y,
            entry: p,
            abortController: new AbortController(),
            traverseToKey: (K = c.state) == null ? void 0 : K.key,
            transition: null
          };
          v.transition = new Ce(v), this.currentTransition = v, this.dispatchNavigation(this.currentTransition);
        }
      );
    }
  }
}, Ht = class extends EventTarget {
  constructor(t, n, r, i, s) {
    super(), this.owner = t, this.id = n, this.key = r, this.url = i, this.state = s, this.url = new URL(i, new URL(globalThis.document.baseURI, globalThis.location.href)).href;
  }
  get index() {
    return this.owner.entriesList.findIndex((t) => t.id === this.id);
  }
  get sameDocument() {
    return !0;
  }
  // polyfill is lost between documents
  getState() {
    return this.state;
  }
  setState(t) {
    this.state = t;
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
}, Ui = class {
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
  static getOrCreate(t) {
    return globalThis.polyfea ? globalThis.polyfea : new Le(t);
  }
  /** @static
   * Initialize the polyfea driver in the global context. 
   * This method is typically invoked by the polyfea controller script `boot.ts`.
   * 
   * @remarks 
   * This method also initializes the Navigation polyfill if it's not already present.
   * It augments `window.customElements.define` to allow for duplicate registration of custom elements.
   * This is particularly useful when different microfrontends need to register the same dependencies.
   * 
   * @param config - Configuration for the `PolyfeaApi`. 
   * For more details, refer to the [`PolyfeaApi`](https://github.com/polyfea/browser-api/blob/main/docs/classes/PolyfeaApi.md) documentation.
   */
  static initialize() {
    globalThis.polyfea || Le.install();
  }
}, Le = class Ln {
  constructor(t) {
    this.config = t, this.loadedResources = /* @__PURE__ */ new Set(), globalThis.navigation && globalThis.navigation.addEventListener("navigate", (n) => {
      n.canIntercept && n.destination.url.startsWith(document.baseURI) && n.intercept();
    });
  }
  getBackend() {
    var t;
    if (!this.backend) {
      let n = (t = document.querySelector('meta[name="polyfea-backend"]')) == null ? void 0 : t.getAttribute("content");
      if (n)
        if (n.startsWith("static://")) {
          const r = n.slice(9);
          this.backend = new ki(this.config || r);
        } else
          this.backend = new Ie(this.config || n);
      else
        this.backend = new Ie(this.config || "./polyfea");
    }
    return this.backend;
  }
  getContextArea(t) {
    return globalThis.navigation ? Qt(globalThis.navigation, "navigatesuccess").pipe(
      fi(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })),
      xn((n) => this.getBackend().getContextArea(t)),
      Re((n, r) => JSON.stringify(n) === JSON.stringify(r))
    ) : this.getBackend().getContextArea(t).pipe(
      Re((n, r) => JSON.stringify(n) === JSON.stringify(r))
    );
  }
  loadMicrofrontend(t, n) {
    if (!n)
      return Promise.resolve();
    const r = [];
    return this.loadMicrofrontendRecursive(t, n, r);
  }
  async loadMicrofrontendRecursive(t, n, r) {
    if (r.includes(n))
      throw new Error("Circular dependency detected: " + r.join(" -> "));
    const i = t.microfrontends[n];
    if (!i)
      throw new Error("Microfrontend specification not found: " + n);
    r.push(n), i.dependsOn && await Promise.all(i.dependsOn.map((o) => this.loadMicrofrontendRecursive(t, o, r)));
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
  loadScript(t) {
    return new Promise((n, r) => {
      var o;
      const i = document.createElement("script");
      i.src = t.href, i.setAttribute("async", ""), t.attributes && Object.entries(t.attributes).forEach(([a, l]) => {
        i.setAttribute(a, l);
      });
      const s = (o = document.querySelector('meta[name="csp-nonce"]')) == null ? void 0 : o.getAttribute("content");
      s && i.setAttribute("nonce", s), this.loadedResources.add(t.href), t.waitOnLoad ? i.onload = () => {
        n();
      } : n(), i.onerror = () => {
        this.loadedResources.delete(t.href), r();
      }, document.head.appendChild(i);
    });
  }
  loadStylesheet(t) {
    return this.loadLink({
      ...t,
      attributes: { ...t.attributes, rel: "stylesheet" }
    });
  }
  loadLink(t) {
    var i;
    const n = document.createElement("link");
    n.href = t.href, n.setAttribute("async", "");
    const r = (i = document.querySelector('meta[name="csp-nonce"]')) == null ? void 0 : i.getAttribute("content");
    return r && n.setAttribute("nonce", r), t.attributes && Object.entries(t.attributes).forEach(([s, o]) => {
      n.setAttribute(s, o);
    }), new Promise((s, o) => {
      this.loadedResources.add(t.href), t.waitOnLoad ? n.onload = () => {
        s();
      } : s(), n.onerror = () => {
        this.loadedResources.delete(t.href), o();
      }, document.head.appendChild(n);
    });
  }
  static install() {
    var i;
    globalThis.polyfea || (globalThis.polyfea = new Ln()), Ci();
    const t = (i = document.querySelector('meta[name="polyfea-duplicit-custom-elements"]')) == null ? void 0 : i.getAttribute("content");
    let n = "warn";
    t === "silent" ? n = "silent" : t === "error" ? n = "error" : t === "verbose" && (n = "verbose");
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
const J = {
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
let j, Pn, It, Un = !1, gt = !1, fe = !1, _ = !1, Pe = null, Zt = !1;
const B = (e, t = "") => () => {
}, Mi = "slot-fb{display:contents}slot-fb[hidden]{display:none}", Ue = "http://www.w3.org/1999/xlink", Me = {}, Fi = "http://www.w3.org/2000/svg", Ni = "http://www.w3.org/1999/xhtml", Di = (e) => e != null, he = (e) => (e = typeof e, e === "object" || e === "function");
function Ji(e) {
  var t, n, r;
  return (r = (n = (t = e.head) === null || t === void 0 ? void 0 : t.querySelector('meta[name="csp-nonce"]')) === null || n === void 0 ? void 0 : n.getAttribute("content")) !== null && r !== void 0 ? r : void 0;
}
const et = (e, t, ...n) => {
  let r = null, i = null, s = null, o = !1, a = !1;
  const l = [], u = (h) => {
    for (let c = 0; c < h.length; c++)
      r = h[c], Array.isArray(r) ? u(r) : r != null && typeof r != "boolean" && ((o = typeof e != "function" && !he(r)) && (r = String(r)), o && a ? l[l.length - 1].$text$ += r : l.push(o ? bt(null, r) : r), a = o);
  };
  if (u(n), t) {
    t.key && (i = t.key), t.name && (s = t.name);
    {
      const h = t.className || t.class;
      h && (t.class = typeof h != "object" ? h : Object.keys(h).filter((c) => h[c]).join(" "));
    }
  }
  if (typeof e == "function")
    return e(t === null ? {} : t, l, Bi);
  const f = bt(e, null);
  return f.$attrs$ = t, l.length > 0 && (f.$children$ = l), f.$key$ = i, f.$name$ = s, f;
}, bt = (e, t) => {
  const n = {
    $flags$: 0,
    $tag$: e,
    $text$: t,
    $elm$: null,
    $children$: null
  };
  return n.$attrs$ = null, n.$key$ = null, n.$name$ = null, n;
}, Mn = {}, Hi = (e) => e && e.$tag$ === Mn, Bi = {
  forEach: (e, t) => e.map(Fe).forEach(t),
  map: (e, t) => e.map(Fe).map(t).map(zi)
}, Fe = (e) => ({
  vattrs: e.$attrs$,
  vchildren: e.$children$,
  vkey: e.$key$,
  vname: e.$name$,
  vtag: e.$tag$,
  vtext: e.$text$
}), zi = (e) => {
  if (typeof e.vtag == "function") {
    const n = Object.assign({}, e.vattrs);
    return e.vkey && (n.key = e.vkey), e.vname && (n.name = e.vname), et(e.vtag, n, ...e.vchildren || []);
  }
  const t = bt(e.vtag, e.vtext);
  return t.$attrs$ = e.vattrs, t.$children$ = e.vchildren, t.$key$ = e.vkey, t.$name$ = e.vname, t;
}, Ki = (e) => gs.map((t) => t(e)).find((t) => !!t), Wi = (e, t) => e != null && !he(e) ? t & 4 ? e === "false" ? !1 : e === "" || !!e : t & 2 ? parseFloat(e) : t & 1 ? String(e) : e : e, Ne = /* @__PURE__ */ new WeakMap(), ji = (e, t, n) => {
  let r = $t.get(e);
  ws && n ? (r = r || new CSSStyleSheet(), typeof r == "string" ? r = t : r.replaceSync(t)) : r = t, $t.set(e, r);
}, Gi = (e, t, n) => {
  var r;
  const i = Fn(t, n), s = $t.get(i);
  if (e = e.nodeType === 11 ? e : A, s)
    if (typeof s == "string") {
      e = e.head || e;
      let o = Ne.get(e), a;
      if (o || Ne.set(e, o = /* @__PURE__ */ new Set()), !o.has(i)) {
        {
          a = A.createElement("style"), a.innerHTML = s;
          const l = (r = T.$nonce$) !== null && r !== void 0 ? r : Ji(A);
          l != null && a.setAttribute("nonce", l), e.insertBefore(a, e.querySelector("link"));
        }
        t.$flags$ & 4 && (a.innerHTML += Mi), o && o.add(i);
      }
    } else
      e.adoptedStyleSheets.includes(s) || (e.adoptedStyleSheets = [...e.adoptedStyleSheets, s]);
  return i;
}, Yi = (e) => {
  const t = e.$cmpMeta$, n = e.$hostElement$, r = t.$flags$, i = B("attachStyles", t.$tagName$), s = Gi(n.shadowRoot ? n.shadowRoot : n.getRootNode(), t, e.$modeName$);
  r & 10 && (n["s-sc"] = s, n.classList.add(s + "-h"), r & 2 && n.classList.add(s + "-s")), i();
}, Fn = (e, t) => "sc-" + (t && e.$flags$ & 32 ? e.$tagName$ + "-" + t : e.$tagName$), De = (e, t, n, r, i, s) => {
  if (n !== r) {
    let o = Ke(e, t), a = t.toLowerCase();
    if (t === "class") {
      const l = e.classList, u = Je(n), f = Je(r);
      l.remove(...u.filter((h) => h && !f.includes(h))), l.add(...f.filter((h) => h && !u.includes(h)));
    } else if (t === "style") {
      for (const l in n)
        (!r || r[l] == null) && (l.includes("-") ? e.style.removeProperty(l) : e.style[l] = "");
      for (const l in r)
        (!n || r[l] !== n[l]) && (l.includes("-") ? e.style.setProperty(l, r[l]) : e.style[l] = r[l]);
    } else if (t !== "key")
      if (t === "ref")
        r && r(e);
      else if (!e.__lookupSetter__(t) && t[0] === "o" && t[1] === "n") {
        if (t[2] === "-" ? t = t.slice(3) : Ke(Ot, a) ? t = a.slice(2) : t = a[2] + t.slice(3), n || r) {
          const l = t.endsWith(Nn);
          t = t.replace(Qi, ""), n && T.rel(e, t, n, l), r && T.ael(e, t, r, l);
        }
      } else {
        const l = he(r);
        if ((o || l && r !== null) && !i)
          try {
            if (e.tagName.includes("-"))
              e[t] = r;
            else {
              const f = r ?? "";
              t === "list" ? o = !1 : (n == null || e[t] != f) && (e[t] = f);
            }
          } catch {
          }
        let u = !1;
        a !== (a = a.replace(/^xlink\:?/, "")) && (t = a, u = !0), r == null || r === !1 ? (r !== !1 || e.getAttribute(t) === "") && (u ? e.removeAttributeNS(Ue, t) : e.removeAttribute(t)) : (!o || s & 4 || i) && !l && (r = r === !0 ? "" : r, u ? e.setAttributeNS(Ue, t, r) : e.setAttribute(t, r));
      }
  }
}, qi = /\s/, Je = (e) => e ? e.split(qi) : [], Nn = "Capture", Qi = new RegExp(Nn + "$"), Dn = (e, t, n, r) => {
  const i = t.$elm$.nodeType === 11 && t.$elm$.host ? t.$elm$.host : t.$elm$, s = e && e.$attrs$ || Me, o = t.$attrs$ || Me;
  for (r in s)
    r in o || De(i, r, s[r], void 0, n, t.$flags$);
  for (r in o)
    De(i, r, s[r], o[r], n, t.$flags$);
}, vt = (e, t, n, r) => {
  var i;
  const s = t.$children$[n];
  let o = 0, a, l, u;
  if (Un || (fe = !0, s.$tag$ === "slot" && (j && r.classList.add(j + "-s"), s.$flags$ |= s.$children$ ? (
    // slot element has fallback content
    2
  ) : (
    // slot element does not have fallback content
    1
  ))), s.$text$ !== null)
    a = s.$elm$ = A.createTextNode(s.$text$);
  else if (s.$flags$ & 1)
    a = s.$elm$ = A.createTextNode("");
  else {
    if (_ || (_ = s.$tag$ === "svg"), a = s.$elm$ = A.createElementNS(_ ? Fi : Ni, s.$flags$ & 2 ? "slot-fb" : s.$tag$), _ && s.$tag$ === "foreignObject" && (_ = !1), Dn(null, s, _), Di(j) && a["s-si"] !== j && a.classList.add(a["s-si"] = j), s.$children$)
      for (o = 0; o < s.$children$.length; ++o)
        l = vt(e, s, o, a), l && a.appendChild(l);
    s.$tag$ === "svg" ? _ = !1 : a.tagName === "foreignObject" && (_ = !0);
  }
  return a["s-hn"] = It, s.$flags$ & 3 && (a["s-sr"] = !0, a["s-fs"] = (i = s.$attrs$) === null || i === void 0 ? void 0 : i.slot, a["s-cr"] = Pn, a["s-sn"] = s.$name$ || "", u = e && e.$children$ && e.$children$[n], u && u.$tag$ === s.$tag$ && e.$elm$ && rt(e.$elm$, !1)), a;
}, rt = (e, t) => {
  var n;
  T.$flags$ |= 1;
  const r = e.childNodes;
  for (let i = r.length - 1; i >= 0; i--) {
    const s = r[i];
    s["s-hn"] !== It && s["s-ol"] && (Bn(s).insertBefore(s, de(s)), s["s-ol"].remove(), s["s-ol"] = void 0, s["s-sh"] = void 0, s.nodeType === 1 && s.setAttribute("slot", (n = s["s-sn"]) !== null && n !== void 0 ? n : ""), fe = !0), t && rt(s, t);
  }
  T.$flags$ &= -2;
}, Jn = (e, t, n, r, i, s) => {
  let o = e["s-cr"] && e["s-cr"].parentNode || e, a;
  for (o.shadowRoot && o.tagName === It && (o = o.shadowRoot); i <= s; ++i)
    r[i] && (a = vt(null, n, i, e), a && (r[i].$elm$ = a, o.insertBefore(a, de(t))));
}, Hn = (e, t, n) => {
  for (let r = t; r <= n; ++r) {
    const i = e[r];
    if (i) {
      const s = i.$elm$;
      Wn(i), s && (gt = !0, s["s-ol"] ? s["s-ol"].remove() : rt(s, !0), s.remove());
    }
  }
}, Xi = (e, t, n, r, i = !1) => {
  let s = 0, o = 0, a = 0, l = 0, u = t.length - 1, f = t[0], h = t[u], c = r.length - 1, d = r[0], p = r[c], y, v;
  for (; s <= u && o <= c; )
    if (f == null)
      f = t[++s];
    else if (h == null)
      h = t[--u];
    else if (d == null)
      d = r[++o];
    else if (p == null)
      p = r[--c];
    else if (ut(f, d, i))
      G(f, d, i), f = t[++s], d = r[++o];
    else if (ut(h, p, i))
      G(h, p, i), h = t[--u], p = r[--c];
    else if (ut(f, p, i))
      (f.$tag$ === "slot" || p.$tag$ === "slot") && rt(f.$elm$.parentNode, !1), G(f, p, i), e.insertBefore(f.$elm$, h.$elm$.nextSibling), f = t[++s], p = r[--c];
    else if (ut(h, d, i))
      (f.$tag$ === "slot" || p.$tag$ === "slot") && rt(h.$elm$.parentNode, !1), G(h, d, i), e.insertBefore(h.$elm$, f.$elm$), h = t[--u], d = r[++o];
    else {
      for (a = -1, l = s; l <= u; ++l)
        if (t[l] && t[l].$key$ !== null && t[l].$key$ === d.$key$) {
          a = l;
          break;
        }
      a >= 0 ? (v = t[a], v.$tag$ !== d.$tag$ ? y = vt(t && t[o], n, a, e) : (G(v, d, i), t[a] = void 0, y = v.$elm$), d = r[++o]) : (y = vt(t && t[o], n, o, e), d = r[++o]), y && Bn(f.$elm$).insertBefore(y, de(f.$elm$));
    }
  s > u ? Jn(e, r[c + 1] == null ? null : r[c + 1].$elm$, n, r, o, c) : o > c && Hn(t, s, u);
}, ut = (e, t, n = !1) => e.$tag$ === t.$tag$ ? e.$tag$ === "slot" ? e.$name$ === t.$name$ : n ? !0 : e.$key$ === t.$key$ : !1, de = (e) => e && e["s-ol"] || e, Bn = (e) => (e["s-ol"] ? e["s-ol"] : e).parentNode, G = (e, t, n = !1) => {
  const r = t.$elm$ = e.$elm$, i = e.$children$, s = t.$children$, o = t.$tag$, a = t.$text$;
  let l;
  a === null ? (_ = o === "svg" ? !0 : o === "foreignObject" ? !1 : _, o === "slot" || Dn(e, t, _), i !== null && s !== null ? Xi(r, i, t, s, n) : s !== null ? (e.$text$ !== null && (r.textContent = ""), Jn(r, null, t, s, 0, s.length - 1)) : i !== null && Hn(i, 0, i.length - 1), _ && o === "svg" && (_ = !1)) : (l = r["s-cr"]) ? l.parentNode.textContent = a : e.$text$ !== a && (r.data = a);
}, zn = (e) => {
  const t = e.childNodes;
  for (const n of t)
    if (n.nodeType === 1) {
      if (n["s-sr"]) {
        const r = n["s-sn"];
        n.hidden = !1;
        for (const i of t)
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
      zn(n);
    }
}, k = [], Kn = (e) => {
  let t, n, r;
  for (const i of e.childNodes) {
    if (i["s-sr"] && (t = i["s-cr"]) && t.parentNode) {
      n = t.parentNode.childNodes;
      const s = i["s-sn"];
      for (r = n.length - 1; r >= 0; r--)
        if (t = n[r], !t["s-cn"] && !t["s-nr"] && t["s-hn"] !== i["s-hn"] && !J.experimentalSlotFixes)
          if (He(t, s)) {
            let o = k.find((a) => a.$nodeToRelocate$ === t);
            gt = !0, t["s-sn"] = t["s-sn"] || s, o ? (o.$nodeToRelocate$["s-sh"] = i["s-hn"], o.$slotRefNode$ = i) : (t["s-sh"] = i["s-hn"], k.push({
              $slotRefNode$: i,
              $nodeToRelocate$: t
            })), t["s-sr"] && k.map((a) => {
              He(a.$nodeToRelocate$, t["s-sn"]) && (o = k.find((l) => l.$nodeToRelocate$ === t), o && !a.$slotRefNode$ && (a.$slotRefNode$ = o.$slotRefNode$));
            });
          } else
            k.some((o) => o.$nodeToRelocate$ === t) || k.push({
              $nodeToRelocate$: t
            });
    }
    i.nodeType === 1 && Kn(i);
  }
}, He = (e, t) => e.nodeType === 1 ? e.getAttribute("slot") === null && t === "" || e.getAttribute("slot") === t : e["s-sn"] === t ? !0 : t === "", Wn = (e) => {
  e.$attrs$ && e.$attrs$.ref && e.$attrs$.ref(null), e.$children$ && e.$children$.map(Wn);
}, Zi = (e, t, n = !1) => {
  var r, i, s, o;
  const a = e.$hostElement$, l = e.$cmpMeta$, u = e.$vnode$ || bt(null, null), f = Hi(t) ? t : et(null, null, t);
  if (It = a.tagName, l.$attrsToReflect$ && (f.$attrs$ = f.$attrs$ || {}, l.$attrsToReflect$.map(([h, c]) => f.$attrs$[c] = a[h])), n && f.$attrs$)
    for (const h of Object.keys(f.$attrs$))
      a.hasAttribute(h) && !["key", "ref", "style", "class"].includes(h) && (f.$attrs$[h] = a[h]);
  f.$tag$ = null, f.$flags$ |= 4, e.$vnode$ = f, f.$elm$ = u.$elm$ = a.shadowRoot || a, j = a["s-sc"], Pn = a["s-cr"], Un = (l.$flags$ & 1) !== 0, gt = !1, G(u, f, n);
  {
    if (T.$flags$ |= 1, fe) {
      Kn(f.$elm$);
      for (const h of k) {
        const c = h.$nodeToRelocate$;
        if (!c["s-ol"]) {
          const d = A.createTextNode("");
          d["s-nr"] = c, c.parentNode.insertBefore(c["s-ol"] = d, c);
        }
      }
      for (const h of k) {
        const c = h.$nodeToRelocate$, d = h.$slotRefNode$;
        if (d) {
          const p = d.parentNode;
          let y = d.nextSibling;
          {
            let v = (r = c["s-ol"]) === null || r === void 0 ? void 0 : r.previousSibling;
            for (; v; ) {
              let m = (i = v["s-nr"]) !== null && i !== void 0 ? i : null;
              if (m && m["s-sn"] === c["s-sn"] && p === m.parentNode && (m = m.nextSibling, !m || !m["s-nr"])) {
                y = m;
                break;
              }
              v = v.previousSibling;
            }
          }
          (!y && p !== c.parentNode || c.nextSibling !== y) && c !== y && (!c["s-hn"] && c["s-ol"] && (c["s-hn"] = c["s-ol"].parentNode.nodeName), p.insertBefore(c, y), c.nodeType === 1 && (c.hidden = (s = c["s-ih"]) !== null && s !== void 0 ? s : !1));
        } else
          c.nodeType === 1 && (n && (c["s-ih"] = (o = c.hidden) !== null && o !== void 0 ? o : !1), c.hidden = !0);
      }
    }
    gt && zn(f.$elm$), T.$flags$ &= -2, k.length = 0;
  }
}, Vi = (e, t) => {
}, jn = (e, t) => (e.$flags$ |= 16, Vi(e, e.$ancestorComponent$), Es(() => ts(e, t))), ts = (e, t) => {
  const n = e.$hostElement$, r = B("scheduleUpdate", e.$cmpMeta$.$tagName$), i = n;
  let s;
  return t ? s = q(i, "componentWillLoad") : s = q(i, "componentWillUpdate"), s = Be(s, () => q(i, "componentWillRender")), r(), Be(s, () => ns(e, i, t));
}, Be = (e, t) => es(e) ? e.then(t) : t(), es = (e) => e instanceof Promise || e && e.then && typeof e.then == "function", ns = async (e, t, n) => {
  const r = e.$hostElement$, i = B("update", e.$cmpMeta$.$tagName$);
  r["s-rc"], n && Yi(e);
  const s = B("render", e.$cmpMeta$.$tagName$);
  rs(e, t, r, n), s(), i(), is(e);
}, rs = (e, t, n, r) => {
  try {
    Pe = t, t = t.render && t.render(), e.$flags$ &= -17, e.$flags$ |= 2, (J.hasRenderFn || J.reflect) && (J.vdomRender || J.reflect) && (J.hydrateServerSide || Zi(e, t, r));
  } catch (l) {
    ot(l, e.$hostElement$);
  }
  return Pe = null, null;
}, is = (e) => {
  const t = e.$cmpMeta$.$tagName$, n = e.$hostElement$, r = B("postUpdate", t), i = n;
  e.$ancestorComponent$, q(i, "componentDidRender"), e.$flags$ & 64 ? (q(i, "componentDidUpdate"), r()) : (e.$flags$ |= 64, q(i, "componentDidLoad"), r());
}, q = (e, t, n) => {
  if (e && e[t])
    try {
      return e[t](n);
    } catch (r) {
      ot(r);
    }
}, ss = (e, t) => st(e).$instanceValues$.get(t), os = (e, t, n, r) => {
  const i = st(e), s = e, o = i.$instanceValues$.get(t), a = i.$flags$, l = s;
  n = Wi(n, r.$members$[t][0]);
  const u = Number.isNaN(o) && Number.isNaN(n);
  if (n !== o && !u) {
    i.$instanceValues$.set(t, n);
    {
      if (r.$watchers$ && a & 128) {
        const h = r.$watchers$[t];
        h && h.map((c) => {
          try {
            l[c](n, o, t);
          } catch (d) {
            ot(d, s);
          }
        });
      }
      if ((a & 18) === 2) {
        if (l.componentShouldUpdate && l.componentShouldUpdate(n, o, t) === !1)
          return;
        jn(i, !1);
      }
    }
  }
}, as = (e, t, n) => {
  var r;
  const i = e.prototype;
  if (t.$members$) {
    e.watchers && (t.$watchers$ = e.watchers);
    const s = Object.entries(t.$members$);
    s.map(([o, [a]]) => {
      (a & 31 || a & 32) && Object.defineProperty(i, o, {
        get() {
          return ss(this, o);
        },
        set(l) {
          os(this, o, l, t);
        },
        configurable: !0,
        enumerable: !0
      });
    });
    {
      const o = /* @__PURE__ */ new Map();
      i.attributeChangedCallback = function(a, l, u) {
        T.jmp(() => {
          var f;
          const h = o.get(a);
          if (this.hasOwnProperty(h))
            u = this[h], delete this[h];
          else {
            if (i.hasOwnProperty(h) && typeof this[h] == "number" && this[h] == u)
              return;
            if (h == null) {
              const c = st(this), d = c == null ? void 0 : c.$flags$;
              if (d && !(d & 8) && d & 128 && u !== l) {
                const y = this, v = (f = t.$watchers$) === null || f === void 0 ? void 0 : f[a];
                v == null || v.forEach((m) => {
                  y[m] != null && y[m].call(y, u, l, a);
                });
              }
              return;
            }
          }
          this[h] = u === null && typeof this[h] == "boolean" ? !1 : u;
        });
      }, e.observedAttributes = Array.from(/* @__PURE__ */ new Set([
        ...Object.keys((r = t.$watchers$) !== null && r !== void 0 ? r : {}),
        ...s.filter(
          ([a, l]) => l[0] & 15
          /* MEMBER_FLAGS.HasAttribute */
        ).map(([a, l]) => {
          var u;
          const f = l[1] || a;
          return o.set(f, a), l[0] & 512 && ((u = t.$attrsToReflect$) === null || u === void 0 || u.push([a, f])), f;
        })
      ]));
    }
  }
  return e;
}, ls = async (e, t, n, r) => {
  let i;
  if (!(t.$flags$ & 32) && (t.$flags$ |= 32, i = e.constructor, customElements.whenDefined(n.$tagName$).then(() => t.$flags$ |= 128), i.style)) {
    let o = i.style;
    typeof o != "string" && (o = o[t.$modeName$ = Ki(e)]);
    const a = Fn(n, t.$modeName$);
    if (!$t.has(a)) {
      const l = B("registerStyles", n.$tagName$);
      ji(a, o, !!(n.$flags$ & 1)), l();
    }
  }
  t.$ancestorComponent$, jn(t, !0);
}, ze = (e) => {
}, cs = (e) => {
  if (!(T.$flags$ & 1)) {
    const t = st(e), n = t.$cmpMeta$, r = B("connectedCallback", n.$tagName$);
    t.$flags$ & 1 ? (Gn(e, t, n.$listeners$), t != null && t.$lazyInstance$ ? ze(t.$lazyInstance$) : t != null && t.$onReadyPromise$ && t.$onReadyPromise$.then(() => ze(t.$lazyInstance$))) : (t.$flags$ |= 1, // TODO(STENCIL-854): Remove code related to legacy shadowDomShim field
    n.$flags$ & 12 && us(e), n.$members$ && Object.entries(n.$members$).map(([i, [s]]) => {
      if (s & 31 && e.hasOwnProperty(i)) {
        const o = e[i];
        delete e[i], e[i] = o;
      }
    }), ls(e, t, n)), r();
  }
}, us = (e) => {
  const t = e["s-cr"] = A.createComment("");
  t["s-cn"] = !0, e.insertBefore(t, e.firstChild);
}, fs = async (e) => {
  if (!(T.$flags$ & 1)) {
    const t = st(e);
    t.$rmListeners$ && (t.$rmListeners$.map((n) => n()), t.$rmListeners$ = void 0);
  }
}, hs = (e, t) => {
  const n = {
    $flags$: t[0],
    $tagName$: t[1]
  };
  n.$members$ = t[2], n.$listeners$ = t[3], n.$watchers$ = e.$watchers$, n.$attrsToReflect$ = [];
  const r = e.prototype.connectedCallback, i = e.prototype.disconnectedCallback;
  return Object.assign(e.prototype, {
    __registerHost() {
      ms(this, n);
    },
    connectedCallback() {
      cs(this), r && r.call(this);
    },
    disconnectedCallback() {
      fs(this), i && i.call(this);
    },
    __attachShadow() {
      this.attachShadow({
        mode: "open",
        delegatesFocus: !!(n.$flags$ & 16)
      });
    }
  }), e.is = n.$tagName$, as(e, n);
}, Gn = (e, t, n, r) => {
  n && n.map(([i, s, o]) => {
    const a = ps(e, i), l = ds(t, o), u = ys(i);
    T.ael(a, s, l, u), (t.$rmListeners$ = t.$rmListeners$ || []).push(() => T.rel(a, s, l, u));
  });
}, ds = (e, t) => (n) => {
  try {
    J.lazyLoad || e.$hostElement$[t](n);
  } catch (r) {
    ot(r);
  }
}, ps = (e, t) => t & 4 ? A : t & 8 ? Ot : t & 16 ? A.body : e, ys = (e) => vs ? {
  passive: (e & 1) !== 0,
  capture: (e & 2) !== 0
} : (e & 2) !== 0, Yn = /* @__PURE__ */ new WeakMap(), st = (e) => Yn.get(e), ms = (e, t) => {
  const n = {
    $flags$: 0,
    $hostElement$: e,
    $cmpMeta$: t,
    $instanceValues$: /* @__PURE__ */ new Map()
  };
  return Gn(e, n, t.$listeners$), Yn.set(e, n);
}, Ke = (e, t) => t in e, ot = (e, t) => (0, console.error)(e, t), $t = /* @__PURE__ */ new Map(), gs = [], Ot = typeof window < "u" ? window : {}, A = Ot.document || { head: {} }, bs = Ot.HTMLElement || class {
}, T = {
  $flags$: 0,
  $resourcesUrl$: "",
  jmp: (e) => e(),
  raf: (e) => requestAnimationFrame(e),
  ael: (e, t, n, r) => e.addEventListener(t, n, r),
  rel: (e, t, n, r) => e.removeEventListener(t, n, r),
  ce: (e, t) => new CustomEvent(e, t)
}, vs = /* @__PURE__ */ (() => {
  let e = !1;
  try {
    A.addEventListener("e", null, Object.defineProperty({}, "passive", {
      get() {
        e = !0;
      }
    }));
  } catch {
  }
  return e;
})(), $s = (e) => Promise.resolve(e), ws = /* @__PURE__ */ (() => {
  try {
    return new CSSStyleSheet(), typeof new CSSStyleSheet().replaceSync == "function";
  } catch {
  }
  return !1;
})(), We = [], qn = [], Ts = (e, t) => (n) => {
  e.push(n), Zt || (Zt = !0, t && T.$flags$ & 4 ? Ss(Vt) : T.raf(Vt));
}, je = (e) => {
  for (let t = 0; t < e.length; t++)
    try {
      e[t](performance.now());
    } catch (n) {
      ot(n);
    }
  e.length = 0;
}, Vt = () => {
  je(We), je(qn), (Zt = We.length > 0) && T.raf(Vt);
}, Ss = (e) => $s().then(e), Es = /* @__PURE__ */ Ts(qn, !0);
/*!
 * Part of Polyfea microfrontends suite - https://github.com/polyfea
 */
function b(e) {
  return typeof e == "function";
}
function pe(e) {
  const n = e((r) => {
    Error.call(r), r.stack = new Error().stack;
  });
  return n.prototype = Object.create(Error.prototype), n.prototype.constructor = n, n;
}
const Bt = pe((e) => function(n) {
  e(this), this.message = n ? `${n.length} errors occurred during unsubscription:
${n.map((r, i) => `${i + 1}) ${r.toString()}`).join(`
  `)}` : "", this.name = "UnsubscriptionError", this.errors = n;
});
function te(e, t) {
  if (e) {
    const n = e.indexOf(t);
    0 <= n && e.splice(n, 1);
  }
}
class I {
  constructor(t) {
    this.initialTeardown = t, this.closed = !1, this._parentage = null, this._finalizers = null;
  }
  unsubscribe() {
    let t;
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
      if (b(r))
        try {
          r();
        } catch (s) {
          t = s instanceof Bt ? s.errors : [s];
        }
      const { _finalizers: i } = this;
      if (i) {
        this._finalizers = null;
        for (const s of i)
          try {
            Ge(s);
          } catch (o) {
            t = t ?? [], o instanceof Bt ? t = [...t, ...o.errors] : t.push(o);
          }
      }
      if (t)
        throw new Bt(t);
    }
  }
  add(t) {
    var n;
    if (t && t !== this)
      if (this.closed)
        Ge(t);
      else {
        if (t instanceof I) {
          if (t.closed || t._hasParent(this))
            return;
          t._addParent(this);
        }
        (this._finalizers = (n = this._finalizers) !== null && n !== void 0 ? n : []).push(t);
      }
  }
  _hasParent(t) {
    const { _parentage: n } = this;
    return n === t || Array.isArray(n) && n.includes(t);
  }
  _addParent(t) {
    const { _parentage: n } = this;
    this._parentage = Array.isArray(n) ? (n.push(t), n) : n ? [n, t] : t;
  }
  _removeParent(t) {
    const { _parentage: n } = this;
    n === t ? this._parentage = null : Array.isArray(n) && te(n, t);
  }
  remove(t) {
    const { _finalizers: n } = this;
    n && te(n, t), t instanceof I && t._removeParent(this);
  }
}
I.EMPTY = (() => {
  const e = new I();
  return e.closed = !0, e;
})();
const Qn = I.EMPTY;
function Xn(e) {
  return e instanceof I || e && "closed" in e && b(e.remove) && b(e.add) && b(e.unsubscribe);
}
function Ge(e) {
  b(e) ? e() : e.unsubscribe();
}
const Ct = {
  onUnhandledError: null,
  onStoppedNotification: null,
  Promise: void 0,
  useDeprecatedSynchronousErrorHandling: !1,
  useDeprecatedNextContext: !1
}, wt = {
  setTimeout(e, t, ...n) {
    const { delegate: r } = wt;
    return r != null && r.setTimeout ? r.setTimeout(e, t, ...n) : setTimeout(e, t, ...n);
  },
  clearTimeout(e) {
    const { delegate: t } = wt;
    return ((t == null ? void 0 : t.clearTimeout) || clearTimeout)(e);
  },
  delegate: void 0
};
function Zn(e) {
  wt.setTimeout(() => {
    const { onUnhandledError: t } = Ct;
    if (t)
      t(e);
    else
      throw e;
  });
}
function ee() {
}
const xs = ye("C", void 0, void 0);
function _s(e) {
  return ye("E", void 0, e);
}
function Rs(e) {
  return ye("N", e, void 0);
}
function ye(e, t, n) {
  return {
    kind: e,
    value: t,
    error: n
  };
}
function pt(e) {
  e();
}
class me extends I {
  constructor(t) {
    super(), this.isStopped = !1, t ? (this.destination = t, Xn(t) && t.add(this)) : this.destination = Os;
  }
  static create(t, n, r) {
    return new Tt(t, n, r);
  }
  next(t) {
    this.isStopped ? Kt(Rs(t), this) : this._next(t);
  }
  error(t) {
    this.isStopped ? Kt(_s(t), this) : (this.isStopped = !0, this._error(t));
  }
  complete() {
    this.isStopped ? Kt(xs, this) : (this.isStopped = !0, this._complete());
  }
  unsubscribe() {
    this.closed || (this.isStopped = !0, super.unsubscribe(), this.destination = null);
  }
  _next(t) {
    this.destination.next(t);
  }
  _error(t) {
    try {
      this.destination.error(t);
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
const As = Function.prototype.bind;
function zt(e, t) {
  return As.call(e, t);
}
class ks {
  constructor(t) {
    this.partialObserver = t;
  }
  next(t) {
    const { partialObserver: n } = this;
    if (n.next)
      try {
        n.next(t);
      } catch (r) {
        ft(r);
      }
  }
  error(t) {
    const { partialObserver: n } = this;
    if (n.error)
      try {
        n.error(t);
      } catch (r) {
        ft(r);
      }
    else
      ft(t);
  }
  complete() {
    const { partialObserver: t } = this;
    if (t.complete)
      try {
        t.complete();
      } catch (n) {
        ft(n);
      }
  }
}
class Tt extends me {
  constructor(t, n, r) {
    super();
    let i;
    if (b(t) || !t)
      i = {
        next: t ?? void 0,
        error: n ?? void 0,
        complete: r ?? void 0
      };
    else {
      let s;
      this && Ct.useDeprecatedNextContext ? (s = Object.create(t), s.unsubscribe = () => this.unsubscribe(), i = {
        next: t.next && zt(t.next, s),
        error: t.error && zt(t.error, s),
        complete: t.complete && zt(t.complete, s)
      }) : i = t;
    }
    this.destination = new ks(i);
  }
}
function ft(e) {
  Zn(e);
}
function Is(e) {
  throw e;
}
function Kt(e, t) {
  const { onStoppedNotification: n } = Ct;
  n && wt.setTimeout(() => n(e, t));
}
const Os = {
  closed: !0,
  next: ee,
  error: Is,
  complete: ee
}, ge = typeof Symbol == "function" && Symbol.observable || "@@observable";
function Lt(e) {
  return e;
}
function Cs(e) {
  return e.length === 0 ? Lt : e.length === 1 ? e[0] : function(n) {
    return e.reduce((r, i) => i(r), n);
  };
}
class $ {
  constructor(t) {
    t && (this._subscribe = t);
  }
  lift(t) {
    const n = new $();
    return n.source = this, n.operator = t, n;
  }
  subscribe(t, n, r) {
    const i = Ps(t) ? t : new Tt(t, n, r);
    return pt(() => {
      const { operator: s, source: o } = this;
      i.add(s ? s.call(i, o) : o ? this._subscribe(i) : this._trySubscribe(i));
    }), i;
  }
  _trySubscribe(t) {
    try {
      return this._subscribe(t);
    } catch (n) {
      t.error(n);
    }
  }
  forEach(t, n) {
    return n = Ye(n), new n((r, i) => {
      const s = new Tt({
        next: (o) => {
          try {
            t(o);
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
  _subscribe(t) {
    var n;
    return (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(t);
  }
  [ge]() {
    return this;
  }
  pipe(...t) {
    return Cs(t)(this);
  }
  toPromise(t) {
    return t = Ye(t), new t((n, r) => {
      let i;
      this.subscribe((s) => i = s, (s) => r(s), () => n(i));
    });
  }
}
$.create = (e) => new $(e);
function Ye(e) {
  var t;
  return (t = e ?? Ct.Promise) !== null && t !== void 0 ? t : Promise;
}
function Ls(e) {
  return e && b(e.next) && b(e.error) && b(e.complete);
}
function Ps(e) {
  return e && e instanceof me || Ls(e) && Xn(e);
}
function Us(e) {
  return b(e == null ? void 0 : e.lift);
}
function C(e) {
  return (t) => {
    if (Us(t))
      return t.lift(function(n) {
        try {
          return e(n, this);
        } catch (r) {
          this.error(r);
        }
      });
    throw new TypeError("Unable to lift unknown Observable type");
  };
}
function F(e, t, n, r, i) {
  return new Ms(e, t, n, r, i);
}
class Ms extends me {
  constructor(t, n, r, i, s, o) {
    super(t), this.onFinalize = s, this.shouldUnsubscribe = o, this._next = n ? function(a) {
      try {
        n(a);
      } catch (l) {
        t.error(l);
      }
    } : super._next, this._error = i ? function(a) {
      try {
        i(a);
      } catch (l) {
        t.error(l);
      } finally {
        this.unsubscribe();
      }
    } : super._error, this._complete = r ? function() {
      try {
        r();
      } catch (a) {
        t.error(a);
      } finally {
        this.unsubscribe();
      }
    } : super._complete;
  }
  unsubscribe() {
    var t;
    if (!this.shouldUnsubscribe || this.shouldUnsubscribe()) {
      const { closed: n } = this;
      super.unsubscribe(), !n && ((t = this.onFinalize) === null || t === void 0 || t.call(this));
    }
  }
}
const Fs = pe((e) => function() {
  e(this), this.name = "ObjectUnsubscribedError", this.message = "object unsubscribed";
});
class Pt extends $ {
  constructor() {
    super(), this.closed = !1, this.currentObservers = null, this.observers = [], this.isStopped = !1, this.hasError = !1, this.thrownError = null;
  }
  lift(t) {
    const n = new Vn(this, this);
    return n.operator = t, n;
  }
  _throwIfClosed() {
    if (this.closed)
      throw new Fs();
  }
  next(t) {
    pt(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.currentObservers || (this.currentObservers = Array.from(this.observers));
        for (const n of this.currentObservers)
          n.next(t);
      }
    });
  }
  error(t) {
    pt(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.hasError = this.isStopped = !0, this.thrownError = t;
        const { observers: n } = this;
        for (; n.length; )
          n.shift().error(t);
      }
    });
  }
  complete() {
    pt(() => {
      if (this._throwIfClosed(), !this.isStopped) {
        this.isStopped = !0;
        const { observers: t } = this;
        for (; t.length; )
          t.shift().complete();
      }
    });
  }
  unsubscribe() {
    this.isStopped = this.closed = !0, this.observers = this.currentObservers = null;
  }
  get observed() {
    var t;
    return ((t = this.observers) === null || t === void 0 ? void 0 : t.length) > 0;
  }
  _trySubscribe(t) {
    return this._throwIfClosed(), super._trySubscribe(t);
  }
  _subscribe(t) {
    return this._throwIfClosed(), this._checkFinalizedStatuses(t), this._innerSubscribe(t);
  }
  _innerSubscribe(t) {
    const { hasError: n, isStopped: r, observers: i } = this;
    return n || r ? Qn : (this.currentObservers = null, i.push(t), new I(() => {
      this.currentObservers = null, te(i, t);
    }));
  }
  _checkFinalizedStatuses(t) {
    const { hasError: n, thrownError: r, isStopped: i } = this;
    n ? t.error(r) : i && t.complete();
  }
  asObservable() {
    const t = new $();
    return t.source = this, t;
  }
}
Pt.create = (e, t) => new Vn(e, t);
class Vn extends Pt {
  constructor(t, n) {
    super(), this.destination = t, this.source = n;
  }
  next(t) {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.next) === null || r === void 0 || r.call(n, t);
  }
  error(t) {
    var n, r;
    (r = (n = this.destination) === null || n === void 0 ? void 0 : n.error) === null || r === void 0 || r.call(n, t);
  }
  complete() {
    var t, n;
    (n = (t = this.destination) === null || t === void 0 ? void 0 : t.complete) === null || n === void 0 || n.call(t);
  }
  _subscribe(t) {
    var n, r;
    return (r = (n = this.source) === null || n === void 0 ? void 0 : n.subscribe(t)) !== null && r !== void 0 ? r : Qn;
  }
}
class ht extends Pt {
  constructor() {
    super(...arguments), this._value = null, this._hasValue = !1, this._isComplete = !1;
  }
  _checkFinalizedStatuses(t) {
    const { hasError: n, _hasValue: r, _value: i, thrownError: s, isStopped: o, _isComplete: a } = this;
    n ? t.error(s) : (o || a) && (r && t.next(i), t.complete());
  }
  next(t) {
    this.isStopped || (this._value = t, this._hasValue = !0);
  }
  complete() {
    const { _hasValue: t, _value: n, _isComplete: r } = this;
    r || (this._isComplete = !0, t && super.next(n), super.complete());
  }
}
function Ns(e) {
  return e && b(e.schedule);
}
function Ds(e) {
  return e[e.length - 1];
}
function Ut(e) {
  return Ns(Ds(e)) ? e.pop() : void 0;
}
function Js(e, t, n, r) {
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
    u((r = r.apply(e, t || [])).next());
  });
}
function qe(e) {
  var t = typeof Symbol == "function" && Symbol.iterator, n = t && e[t], r = 0;
  if (n)
    return n.call(e);
  if (e && typeof e.length == "number")
    return {
      next: function() {
        return e && r >= e.length && (e = void 0), { value: e && e[r++], done: !e };
      }
    };
  throw new TypeError(t ? "Object is not iterable." : "Symbol.iterator is not defined.");
}
function Q(e) {
  return this instanceof Q ? (this.v = e, this) : new Q(e);
}
function Hs(e, t, n) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var r = n.apply(e, t || []), i, s = [];
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
    c.value instanceof Q ? Promise.resolve(c.value.v).then(u, f) : h(s[0][2], c);
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
function Bs(e) {
  if (!Symbol.asyncIterator)
    throw new TypeError("Symbol.asyncIterator is not defined.");
  var t = e[Symbol.asyncIterator], n;
  return t ? t.call(e) : (e = typeof qe == "function" ? qe(e) : e[Symbol.iterator](), n = {}, r("next"), r("throw"), r("return"), n[Symbol.asyncIterator] = function() {
    return this;
  }, n);
  function r(s) {
    n[s] = e[s] && function(o) {
      return new Promise(function(a, l) {
        o = e[s](o), i(a, l, o.done, o.value);
      });
    };
  }
  function i(s, o, a, l) {
    Promise.resolve(l).then(function(u) {
      s({ value: u, done: a });
    }, o);
  }
}
const be = (e) => e && typeof e.length == "number" && typeof e != "function";
function tr(e) {
  return b(e == null ? void 0 : e.then);
}
function er(e) {
  return b(e[ge]);
}
function nr(e) {
  return Symbol.asyncIterator && b(e == null ? void 0 : e[Symbol.asyncIterator]);
}
function rr(e) {
  return new TypeError(`You provided ${e !== null && typeof e == "object" ? "an invalid object" : `'${e}'`} where a stream was expected. You can provide an Observable, Promise, ReadableStream, Array, AsyncIterable, or Iterable.`);
}
function zs() {
  return typeof Symbol != "function" || !Symbol.iterator ? "@@iterator" : Symbol.iterator;
}
const ir = zs();
function sr(e) {
  return b(e == null ? void 0 : e[ir]);
}
function or(e) {
  return Hs(this, arguments, function* () {
    const n = e.getReader();
    try {
      for (; ; ) {
        const { value: r, done: i } = yield Q(n.read());
        if (i)
          return yield Q(void 0);
        yield yield Q(r);
      }
    } finally {
      n.releaseLock();
    }
  });
}
function ar(e) {
  return b(e == null ? void 0 : e.getReader);
}
function D(e) {
  if (e instanceof $)
    return e;
  if (e != null) {
    if (er(e))
      return Ks(e);
    if (be(e))
      return Ws(e);
    if (tr(e))
      return js(e);
    if (nr(e))
      return lr(e);
    if (sr(e))
      return Gs(e);
    if (ar(e))
      return Ys(e);
  }
  throw rr(e);
}
function Ks(e) {
  return new $((t) => {
    const n = e[ge]();
    if (b(n.subscribe))
      return n.subscribe(t);
    throw new TypeError("Provided object does not correctly implement Symbol.observable");
  });
}
function Ws(e) {
  return new $((t) => {
    for (let n = 0; n < e.length && !t.closed; n++)
      t.next(e[n]);
    t.complete();
  });
}
function js(e) {
  return new $((t) => {
    e.then((n) => {
      t.closed || (t.next(n), t.complete());
    }, (n) => t.error(n)).then(null, Zn);
  });
}
function Gs(e) {
  return new $((t) => {
    for (const n of e)
      if (t.next(n), t.closed)
        return;
    t.complete();
  });
}
function lr(e) {
  return new $((t) => {
    qs(e, t).catch((n) => t.error(n));
  });
}
function Ys(e) {
  return lr(or(e));
}
function qs(e, t) {
  var n, r, i, s;
  return Js(this, void 0, void 0, function* () {
    try {
      for (n = Bs(e); r = yield n.next(), !r.done; ) {
        const o = r.value;
        if (t.next(o), t.closed)
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
    t.complete();
  });
}
function U(e, t, n, r = 0, i = !1) {
  const s = t.schedule(function() {
    n(), i ? e.add(this.schedule(null, r)) : this.unsubscribe();
  }, r);
  if (e.add(s), !i)
    return s;
}
function cr(e, t = 0) {
  return C((n, r) => {
    n.subscribe(F(r, (i) => U(r, e, () => r.next(i), t), () => U(r, e, () => r.complete(), t), (i) => U(r, e, () => r.error(i), t)));
  });
}
function ur(e, t = 0) {
  return C((n, r) => {
    r.add(e.schedule(() => n.subscribe(r), t));
  });
}
function Qs(e, t) {
  return D(e).pipe(ur(t), cr(t));
}
function Xs(e, t) {
  return D(e).pipe(ur(t), cr(t));
}
function Zs(e, t) {
  return new $((n) => {
    let r = 0;
    return t.schedule(function() {
      r === e.length ? n.complete() : (n.next(e[r++]), n.closed || this.schedule());
    });
  });
}
function Vs(e, t) {
  return new $((n) => {
    let r;
    return U(n, t, () => {
      r = e[ir](), U(n, t, () => {
        let i, s;
        try {
          ({ value: i, done: s } = r.next());
        } catch (o) {
          n.error(o);
          return;
        }
        s ? n.complete() : n.next(i);
      }, 0, !0);
    }), () => b(r == null ? void 0 : r.return) && r.return();
  });
}
function fr(e, t) {
  if (!e)
    throw new Error("Iterable cannot be null");
  return new $((n) => {
    U(n, t, () => {
      const r = e[Symbol.asyncIterator]();
      U(n, t, () => {
        r.next().then((i) => {
          i.done ? n.complete() : n.next(i.value);
        });
      }, 0, !0);
    });
  });
}
function to(e, t) {
  return fr(or(e), t);
}
function eo(e, t) {
  if (e != null) {
    if (er(e))
      return Qs(e, t);
    if (be(e))
      return Zs(e, t);
    if (tr(e))
      return Xs(e, t);
    if (nr(e))
      return fr(e, t);
    if (sr(e))
      return Vs(e, t);
    if (ar(e))
      return to(e, t);
  }
  throw rr(e);
}
function Mt(e, t) {
  return t ? eo(e, t) : D(e);
}
function no(...e) {
  const t = Ut(e);
  return Mt(e, t);
}
function ro(e, t) {
  const n = b(e) ? e : () => e, r = (i) => i.error(n());
  return new $(t ? (i) => t.schedule(r, 0, i) : r);
}
const io = pe((e) => function() {
  e(this), this.name = "EmptyError", this.message = "no elements in sequence";
});
function ne(e, t) {
  const n = typeof t == "object";
  return new Promise((r, i) => {
    const s = new Tt({
      next: (o) => {
        r(o), s.unsubscribe();
      },
      error: i,
      complete: () => {
        n ? r(t.defaultValue) : i(new io());
      }
    });
    e.subscribe(s);
  });
}
function ve(e, t) {
  return C((n, r) => {
    let i = 0;
    n.subscribe(F(r, (s) => {
      r.next(e.call(t, s, i++));
    }));
  });
}
const { isArray: so } = Array;
function oo(e, t) {
  return so(t) ? e(...t) : e(t);
}
function ao(e) {
  return ve((t) => oo(e, t));
}
function lo(e, t, n, r, i, s, o, a) {
  const l = [];
  let u = 0, f = 0, h = !1;
  const c = () => {
    h && !l.length && !u && t.complete();
  }, d = (y) => u < r ? p(y) : l.push(y), p = (y) => {
    s && t.next(y), u++;
    let v = !1;
    D(n(y, f++)).subscribe(F(t, (m) => {
      i == null || i(m), s ? d(m) : t.next(m);
    }, () => {
      v = !0;
    }, void 0, () => {
      if (v)
        try {
          for (u--; l.length && u < r; ) {
            const m = l.shift();
            o ? U(t, o, () => p(m)) : p(m);
          }
          c();
        } catch (m) {
          t.error(m);
        }
    }));
  };
  return e.subscribe(F(t, d, () => {
    h = !0, c();
  })), () => {
    a == null || a();
  };
}
function $e(e, t, n = 1 / 0) {
  return b(t) ? $e((r, i) => ve((s, o) => t(r, s, i, o))(D(e(r, i))), n) : (typeof t == "number" && (n = t), C((r, i) => lo(r, i, e, n)));
}
function co(e = 1 / 0) {
  return $e(Lt, e);
}
function hr() {
  return co(1);
}
function Qe(...e) {
  return hr()(Mt(e, Ut(e)));
}
function uo(e) {
  return new $((t) => {
    D(e()).subscribe(t);
  });
}
const fo = ["addListener", "removeListener"], ho = ["addEventListener", "removeEventListener"], po = ["on", "off"];
function re(e, t, n, r) {
  if (b(n) && (r = n, n = void 0), r)
    return re(e, t, n).pipe(ao(r));
  const [i, s] = go(e) ? ho.map((o) => (a) => e[o](t, a, n)) : yo(e) ? fo.map(Xe(e, t)) : mo(e) ? po.map(Xe(e, t)) : [];
  if (!i && be(e))
    return $e((o) => re(o, t, n))(D(e));
  if (!i)
    throw new TypeError("Invalid event target");
  return new $((o) => {
    const a = (...l) => o.next(1 < l.length ? l : l[0]);
    return i(a), () => s(a);
  });
}
function Xe(e, t) {
  return (n) => (r) => e[n](t, r);
}
function yo(e) {
  return b(e.addListener) && b(e.removeListener);
}
function mo(e) {
  return b(e.on) && b(e.off);
}
function go(e) {
  return b(e.addEventListener) && b(e.removeEventListener);
}
const bo = new $(ee);
function vo(...e) {
  const t = Ut(e);
  return C((n, r) => {
    hr()(Mt([n, ...e], t)).subscribe(r);
  });
}
function dr(...e) {
  return vo(...e);
}
function Ze(e, t = Lt) {
  return e = e ?? $o, C((n, r) => {
    let i, s = !0;
    n.subscribe(F(r, (o) => {
      const a = t(o);
      (s || !e(i, a)) && (s = !1, i = a, r.next(o));
    }));
  });
}
function $o(e, t) {
  return e === t;
}
function wo(...e) {
  const t = Ut(e);
  return C((n, r) => {
    (t ? Qe(e, n, t) : Qe(e, n)).subscribe(r);
  });
}
function we(e, t) {
  return C((n, r) => {
    let i = null, s = 0, o = !1;
    const a = () => o && !i && r.complete();
    n.subscribe(F(r, (l) => {
      i == null || i.unsubscribe();
      let u = 0;
      const f = s++;
      D(e(l, f)).subscribe(i = F(r, (h) => r.next(t ? t(l, h, f, u++) : h), () => {
        i = null, a();
      }));
    }, () => {
      o = !0, a();
    }));
  });
}
function To(e, t, n) {
  const r = b(e) || t || n ? { next: e, error: t, complete: n } : e;
  return r ? C((i, s) => {
    var o;
    (o = r.subscribe) === null || o === void 0 || o.call(r);
    let a = !0;
    i.subscribe(F(s, (l) => {
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
  }) : Lt;
}
const So = "http://./polyfea".replace(/\/+$/, "");
class it {
  constructor(t = {}) {
    this.configuration = t;
  }
  set config(t) {
    this.configuration = t;
  }
  get basePath() {
    return this.configuration.basePath != null ? this.configuration.basePath : So;
  }
  get fetchApi() {
    return this.configuration.fetchApi;
  }
  get middleware() {
    return this.configuration.middleware || [];
  }
  get queryParamsStringify() {
    return this.configuration.queryParamsStringify || pr;
  }
  get username() {
    return this.configuration.username;
  }
  get password() {
    return this.configuration.password;
  }
  get apiKey() {
    const t = this.configuration.apiKey;
    if (t)
      return typeof t == "function" ? t : () => t;
  }
  get accessToken() {
    const t = this.configuration.accessToken;
    if (t)
      return typeof t == "function" ? t : async () => t;
  }
  get headers() {
    return this.configuration.headers;
  }
  get credentials() {
    return this.configuration.credentials;
  }
}
const Eo = new it();
class Ft {
  constructor(t = Eo) {
    this.configuration = t, this.fetchApi = async (n, r) => {
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
          throw o instanceof Error ? new Ao(o, "The request failed and the interceptors did not return an alternative response") : o;
      }
      for (const o of this.middleware)
        o.post && (s = await o.post({
          fetch: this.fetchApi,
          url: i.url,
          init: i.init,
          response: s.clone()
        }) || s);
      return s;
    }, this.middleware = t.middleware;
  }
  withMiddleware(...t) {
    const n = this.clone();
    return n.middleware = n.middleware.concat(...t), n;
  }
  withPreMiddleware(...t) {
    const n = t.map((r) => ({ pre: r }));
    return this.withMiddleware(...n);
  }
  withPostMiddleware(...t) {
    const n = t.map((r) => ({ post: r }));
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
  isJsonMime(t) {
    return t ? Ft.jsonRegex.test(t) : !1;
  }
  async request(t, n) {
    const { url: r, init: i } = await this.createFetchParams(t, n), s = await this.fetchApi(r, i);
    if (s && s.status >= 200 && s.status < 300)
      return s;
    throw new Ro(s, "Response returned an error code");
  }
  async createFetchParams(t, n) {
    let r = this.configuration.basePath + t.path;
    t.query !== void 0 && Object.keys(t.query).length !== 0 && (r += "?" + this.configuration.queryParamsStringify(t.query));
    const i = Object.assign({}, this.configuration.headers, t.headers);
    Object.keys(i).forEach((f) => i[f] === void 0 ? delete i[f] : {});
    const s = typeof n == "function" ? n : async () => n, o = {
      method: t.method,
      headers: i,
      body: t.body,
      credentials: this.configuration.credentials
    }, a = Object.assign(Object.assign({}, o), await s({
      init: o,
      context: t
    }));
    let l;
    _o(a.body) || a.body instanceof URLSearchParams || xo(a.body) ? l = a.body : this.isJsonMime(i["Content-Type"]) ? l = JSON.stringify(a.body) : l = a.body;
    const u = Object.assign(Object.assign({}, a), { body: l });
    return { url: r, init: u };
  }
  /**
   * Create a shallow clone of `this` by constructing a new instance
   * and then shallow cloning data members.
   */
  clone() {
    const t = this.constructor, n = new t(this.configuration);
    return n.middleware = this.middleware.slice(), n;
  }
}
Ft.jsonRegex = new RegExp("^(:?application/json|[^;/ 	]+/[^;/ 	]+[+]json)[ 	]*(:?;.*)?$", "i");
function xo(e) {
  return typeof Blob < "u" && e instanceof Blob;
}
function _o(e) {
  return typeof FormData < "u" && e instanceof FormData;
}
class Ro extends Error {
  constructor(t, n) {
    super(n), this.response = t, this.name = "ResponseError";
  }
}
class Ao extends Error {
  constructor(t, n) {
    super(n), this.cause = t, this.name = "FetchError";
  }
}
class Ve extends Error {
  constructor(t, n) {
    super(n), this.field = t, this.name = "RequiredError";
  }
}
function E(e, t) {
  const n = e[t];
  return n != null;
}
function pr(e, t = "") {
  return Object.keys(e).map((n) => yr(n, e[n], t)).filter((n) => n.length > 0).join("&");
}
function yr(e, t, n = "") {
  const r = n + (n.length ? `[${e}]` : e);
  if (t instanceof Array) {
    const i = t.map((s) => encodeURIComponent(String(s))).join(`&${encodeURIComponent(r)}=`);
    return `${encodeURIComponent(r)}=${i}`;
  }
  if (t instanceof Set) {
    const i = Array.from(t);
    return yr(e, i, n);
  }
  return t instanceof Date ? `${encodeURIComponent(r)}=${encodeURIComponent(t.toISOString())}` : t instanceof Object ? pr(t, r) : `${encodeURIComponent(r)}=${encodeURIComponent(String(t))}`;
}
function mr(e, t) {
  return Object.keys(e).reduce((n, r) => Object.assign(Object.assign({}, n), { [r]: t(e[r]) }), {});
}
class tn {
  constructor(t, n = (r) => r) {
    this.raw = t, this.transformer = n;
  }
  async value() {
    return this.transformer(await this.raw.json());
  }
}
function ko(e) {
  return Io(e);
}
function Io(e, t) {
  return e == null ? e : {
    microfrontend: E(e, "microfrontend") ? e.microfrontend : void 0,
    tagName: e.tagName,
    attributes: E(e, "attributes") ? e.attributes : void 0,
    style: E(e, "style") ? e.style : void 0
  };
}
function Oo(e) {
  return Co(e);
}
function Co(e, t) {
  return e == null ? e : {
    kind: E(e, "kind") ? e.kind : void 0,
    href: E(e, "href") ? e.href : void 0,
    attributes: E(e, "attributes") ? e.attributes : void 0,
    waitOnLoad: E(e, "waitOnLoad") ? e.waitOnLoad : void 0
  };
}
function gr(e) {
  return Lo(e);
}
function Lo(e, t) {
  return e == null ? e : {
    dependsOn: E(e, "dependsOn") ? e.dependsOn : void 0,
    module: E(e, "module") ? e.module : void 0,
    resources: E(e, "resources") ? e.resources.map(Oo) : void 0
  };
}
function br(e) {
  return Po(e);
}
function Po(e, t) {
  return e == null ? e : {
    elements: e.elements.map(ko),
    microfrontends: E(e, "microfrontends") ? mr(e.microfrontends, gr) : void 0
  };
}
function Uo(e) {
  return Mo(e);
}
function Mo(e, t) {
  return e == null ? e : {
    name: e.name,
    path: E(e, "path") ? e.path : void 0,
    contextArea: E(e, "contextArea") ? br(e.contextArea) : void 0
  };
}
function Fo(e) {
  return No(e);
}
function No(e, t) {
  return e == null ? e : {
    contextAreas: E(e, "contextAreas") ? e.contextAreas.map(Uo) : void 0,
    microfrontends: mr(e.microfrontends, gr)
  };
}
class St extends Ft {
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextAreaRaw(t, n) {
    if (t.name === null || t.name === void 0)
      throw new Ve("name", "Required parameter requestParameters.name was null or undefined when calling getContextArea.");
    if (t.path === null || t.path === void 0)
      throw new Ve("path", "Required parameter requestParameters.path was null or undefined when calling getContextArea.");
    const r = {};
    t.path !== void 0 && (r.path = t.path), t.take !== void 0 && (r.take = t.take);
    const i = {}, s = await this.request({
      path: "/context-area/{name}".replace("{name}", encodeURIComponent(String(t.name))),
      method: "GET",
      headers: i,
      query: r
    }, n);
    return new tn(s, (o) => br(o));
  }
  /**
   * Retrieve the context area information. This information includes the elements and  microfrontends required for these elements. The actual content depends on the input path and  the user role, which is determined server-side.
   * Get the context area information.
   */
  async getContextArea(t, n) {
    return await (await this.getContextAreaRaw(t, n)).value();
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfigRaw(t) {
    const n = {}, r = {}, i = await this.request({
      path: "/static-config",
      method: "GET",
      headers: r,
      query: n
    }, t);
    return new tn(i, (s) => Fo(s));
  }
  /**
   * Retrieve the static configuration of the application\'s context areas.  This includes a combination of all microfrontends and web components.  This approach is advantageous when the frontend logic is simple and static,  particularly during development or testing phases.
   * Get the static information about all resources and context areas.
   */
  async getStaticConfig(t) {
    return await (await this.getStaticConfigRaw(t)).value();
  }
}
class Do {
  constructor(t = "./polyfea") {
    this.spec$ = new $();
    let n;
    typeof t == "string" ? (t.length === 0 && (t = "./polyfea"), n = new St(new it({ basePath: t }))) : t instanceof it ? n = new St(t) : n = t, this.spec$ = Mt(n.getStaticConfig()).pipe(dr(bo));
  }
  getContextArea(t) {
    let n = globalThis.location.pathname;
    if (globalThis.document.baseURI) {
      const i = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      n.startsWith(i) && (n = "./" + n.substring(i.length));
    }
    return this.spec$.pipe(ve((r) => {
      for (let i of r.contextAreas)
        if (i.name === t && new RegExp(i.path).test(n))
          return Object.assign(Object.assign({}, i.contextArea), { microfrontends: Object.assign(Object.assign({}, r.microfrontends), i.contextArea.microfrontends) });
      return null;
    }));
  }
}
class en {
  constructor(t = "./polyfea") {
    var n, r;
    typeof t == "string" ? this.api = new St(new it({
      basePath: new URL(t, new URL(((n = globalThis.document) === null || n === void 0 ? void 0 : n.baseURI) || "/", ((r = globalThis.location) === null || r === void 0 ? void 0 : r.href) || "http://localhost")).href
    })) : t instanceof it ? this.api = new St(t) : this.api = t;
  }
  getContextArea(t) {
    var n, r;
    let i = ((n = globalThis.location) === null || n === void 0 ? void 0 : n.pathname) || "/";
    if (!((r = globalThis.document) === null || r === void 0) && r.baseURI) {
      const l = new URL(globalThis.document.baseURI, globalThis.location.href).pathname;
      i.startsWith(l) && (i = "./" + i.substring(l.length));
    }
    const s = localStorage.getItem(`polyfea-context[${t},${i}]`), o = uo(() => this.api.getContextAreaRaw({ name: t, path: i })).pipe(we((a) => a.raw.ok ? a.value() : ro(() => new Error(a.raw.statusText))), To((a) => {
      a && localStorage.setItem(`polyfea-context[${t},${i}]`, JSON.stringify(a));
    }));
    if (s) {
      const a = JSON.parse(s);
      return no(a).pipe(dr(o));
    } else
      return o;
  }
}
class Jo {
  /**
   * Constructs a new NavigationDestination instance.
   * @param url - The URL of the navigation destination.
   */
  constructor(t) {
    this.url = t;
  }
}
class Ho extends Event {
  constructor(t, n) {
    super("navigate", { bubbles: !0, cancelable: !0 }), this.transition = t, this.interceptPromises = [], this.downloadRequest = null, this.formData = null, this.hashChange = !1, this.userInitiated = !1;
    let r = new URL(t.href, new URL(globalThis.document.baseURI, globalThis.location.href));
    const i = new URL(n.url, new URL(globalThis.document.baseURI, globalThis.location.href));
    this.canIntercept = i.protocol === r.protocol && i.host === r.host && i.port === r.port, this.destination = new Jo(r.href);
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
  intercept(t) {
    t != null && t.handler && this.interceptPromises.push(t.handler(this));
  }
}
class nn extends Event {
  constructor(t, n) {
    super("currententrychange", { bubbles: !0, cancelable: !0 }), this.navigationType = t, this.from = n;
  }
}
class rn {
  constructor(t) {
    this.request = t;
  }
  get finished() {
    return ne(this.request.finished);
  }
  get from() {
    return this.request.entry;
  }
  get type() {
    return this.request.mode;
  }
}
function Bo(e = !1) {
  return Et.tryRegister(e);
}
class zo {
  constructor(t) {
    this.commited$ = t.committed, this.finished$ = t.finished;
  }
  get commited() {
    return ne(this.commited$);
  }
  get finished() {
    return ne(this.finished$);
  }
}
class Et extends EventTarget {
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
    var t;
    return ((t = this.currentTransition) === null || t === void 0 ? void 0 : t.transition) || null;
  }
  navigate(t, n) {
    let r = "push";
    return ((n == null ? void 0 : n.history) === "replace" || n != null && n.replace) && (r = "replace"), this.nextTransitionRequest(r, t, n);
  }
  back() {
    if (this.currentEntryIndex < 1)
      throw { name: "InvaliStateError", message: "Cannot go back from initial state" };
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { traverseTo: this.entriesList[this.currentEntryIndex - 1].key });
  }
  forward(t) {
    return this.nextTransitionRequest("traverse", this.currentEntry.url, { info: t, traverseTo: this.entriesList[this.currentEntryIndex + 1].key });
  }
  reload(t) {
    return this.nextTransitionRequest("reload", this.currentEntry.url, { info: t });
  }
  traverseTo(t, n) {
    const r = this.entriesList.find((i) => i.key === t);
    if (!r)
      throw { name: "InvaliStateError", message: "Cannot traverse to unknown state" };
    return this.nextTransitionRequest("traverse", r.url, { info: n == null ? void 0 : n.info, traverseTo: t });
  }
  updateCurrentEntry(t) {
    this.entriesList[this.currentEntryIndex].setState(JSON.parse(JSON.stringify(t == null ? void 0 : t.state))), this.currentEntry.dispatchEvent(new nn("replace", this.currentEntry));
  }
  constructor() {
    super(), this.entriesList = [], this.idCounter = 0, this.transitionRequests = new Pt(), this.currentTransition = null, this.currentEntryIndex = -1, this.pushstateDelay = 35, this.rawHistoryMethods = {
      pushState: globalThis.history.pushState,
      replaceState: globalThis.history.replaceState,
      go: globalThis.history.go,
      back: globalThis.history.back,
      forward: globalThis.history.forward
    }, this.transitionRequests.subscribe((t) => this.executeRequest(t));
  }
  nextTransitionRequest(t, n, r) {
    const i = `@${++this.idCounter}-navigation-polyfill-transition`, s = {
      mode: t,
      href: new URL(n, new URL(globalThis.document.baseURI, globalThis.location.href)).href,
      info: r == null ? void 0 : r.info,
      state: r == null ? void 0 : r.state,
      committed: new ht(),
      finished: new ht(),
      entry: new Wt(this, i, i, n.toString(), r == null ? void 0 : r.state),
      abortController: new AbortController(),
      traverseToKey: r == null ? void 0 : r.traverseTo,
      transition: null
    };
    return this.transitionRequests.next(s), new zo(s);
  }
  async executeRequest(t) {
    this.currentTransition && (this.currentTransition.abortController.abort(), this.currentTransition.finished.error("aborted - new navigation started"), this.currentTransition.committed.closed || this.currentTransition.committed.error("aborted - new navigation started"), globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: Error("aborted - new navigation started") })), this.clearTransition(t)), t.transition = new rn(t), this.currentTransition = t;
    try {
      await this.commit(t), this.currentEntry.dispatchEvent(new nn(t.mode, t.transition.from)), await this.dispatchNavigation(t);
    } catch {
      t.finished.error("aborted"), t.committed.error("aborted");
    } finally {
      this.clearTransition(t);
    }
  }
  dispatchNavigation(t) {
    var n;
    const r = new Ho(t, this.currentEntry);
    if (!((n = globalThis.navigation) === null || n === void 0) && n.dispatchEvent(r))
      return r.interceptPromises.length > 0 ? Promise.all(r.interceptPromises.filter((i) => !!(i != null && i.then))).then(() => {
        globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(t), t.finished.next(), t.finished.complete();
      }).catch((i) => {
        globalThis.navigation.dispatchEvent(new ErrorEvent("navigateerror", { bubbles: !0, cancelable: !0, error: i })), this.clearTransition(t), t.finished.error(i);
      }) : (globalThis.navigation.dispatchEvent(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), this.clearTransition(t), t.finished.next(), t.finished.complete(), Promise.resolve());
  }
  clearTransition(t) {
    var n;
    ((n = this.currentTransition) === null || n === void 0 ? void 0 : n.entry.id) === t.entry.id && (this.currentTransition = null);
  }
  commit(t) {
    switch (t.mode) {
      case "push":
        return this.commitPushTransition(t);
      case "replace":
        return this.commitReplaceTransition(t);
      case "reload":
        return this.commitReloadTransition(t);
      case "traverse":
        return this.commitTraverseTransition(t);
    }
  }
  async pushstateAsync(t, n = () => {
  }) {
    return new Promise((r, i) => {
      setTimeout(() => {
        n(t), t.committed.next(), t.committed.complete(), r();
      }, this.pushstateDelay);
    });
  }
  commitPushTransition(t) {
    return this.rawHistoryMethods.pushState.apply(globalThis.history, [t.entry.cloneable, "", t.href]), this.pushstateAsync(t, (n) => {
      this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), n.entry];
    });
  }
  commitReplaceTransition(t) {
    return t.entry.key = this.currentEntry.key, this.entriesList[this.currentEntryIndex] = t.entry, this.rawHistoryMethods.replaceState.apply(globalThis.history, [t.entry.cloneable, "", t.href]), this.pushstateAsync(t);
  }
  commitTraverseTransition(t) {
    return new Promise(async (n, r) => {
      const i = this.entriesList.findIndex((o) => o.key === t.traverseToKey);
      i < 0 && r("target entry not found");
      const s = i - this.currentEntryIndex;
      this.rawHistoryMethods.go.apply(globalThis.history, [s]), await this.pushstateAsync(t, (o) => {
        const a = this.entriesList.findIndex((l) => l.key === o.traverseToKey);
        a < 0 && o.committed.error(new Error("target entry not found")), this.currentEntryIndex = a, o.committed.next(), o.committed.complete(), n();
      });
    });
  }
  commitReloadTransition(t) {
    return t.committed.next(), t.committed.complete(), t.finished.subscribe({
      next: () => globalThis.location.reload(),
      error: () => globalThis.location.reload()
    }), Promise.resolve();
  }
  static tryRegister(t = !1) {
    if (!globalThis.navigation) {
      const n = new Et();
      return n.doRegister(t), n;
    }
    return globalThis.navigation;
  }
  static unregister() {
    globalThis.navigation && globalThis.navigation instanceof Et && (globalThis.navigation.doUnregister && globalThis.navigation.doUnregister(), globalThis.navigation = void 0);
  }
  doRegister(t) {
    var n, r, i, s, o;
    if (!globalThis.navigation && !globalThis.navigation) {
      globalThis.navigation = this, this.entriesList = [new Wt(this, "initial", "initial", globalThis.location.href, void 0)], this.currentEntryIndex = 0;
      const a = ((n = globalThis.history) === null || n === void 0 ? void 0 : n.pushState) || ((c, d, p) => {
      }), l = ((r = globalThis.history) === null || r === void 0 ? void 0 : r.replaceState) || ((c, d, p) => {
      }), u = ((i = globalThis.history) === null || i === void 0 ? void 0 : i.go) || ((c) => {
      }), f = ((s = globalThis.history) === null || s === void 0 ? void 0 : s.back) || (() => {
      }), h = ((o = globalThis.history) === null || o === void 0 ? void 0 : o.forward) || (() => {
      });
      this.doUnregister = () => {
        globalThis.history.pushState = a, globalThis.history.replaceState = l, globalThis.history.go = u, globalThis.history.back = f, globalThis.history.forward = h;
      }, t ? (this.rawHistoryMethods.pushState = () => {
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
        var d, p, y, v, m, R, z;
        if (this.currentTransition && ((d = c.state) === null || d === void 0 ? void 0 : d.id) === ((y = (p = this.currentTransition) === null || p === void 0 ? void 0 : p.entry) === null || y === void 0 ? void 0 : y.id))
          return;
        (v = this.currentTransition) === null || v === void 0 || v.abortController.abort();
        const tt = new ht();
        tt.complete();
        let L;
        if (!((m = c.state) === null || m === void 0) && m.key) {
          const x = this.entriesList.findIndex((Nt) => {
            var W;
            return Nt.key === ((W = c.state) === null || W === void 0 ? void 0 : W.key);
          });
          x >= 0 && (this.currentEntryIndex = x, L = this.entriesList[x]);
        }
        if (!L) {
          let x = `@${++this.idCounter}-navigation-polyfill-popstate`;
          L = new Wt(this, x, x, globalThis.location.href, c.state), this.entriesList = [...this.entriesList.slice(0, ++this.currentEntryIndex), L];
        }
        const at = new ht(), K = {
          mode: "traverse",
          href: globalThis.location.href,
          info: void 0,
          state: ((R = c.state) === null || R === void 0 ? void 0 : R.state) || c.state,
          committed: tt,
          finished: at,
          entry: L,
          abortController: new AbortController(),
          traverseToKey: (z = c.state) === null || z === void 0 ? void 0 : z.key,
          transition: null
        };
        K.transition = new rn(K), this.currentTransition = K, this.dispatchNavigation(this.currentTransition);
      });
    }
  }
}
class Wt extends EventTarget {
  constructor(t, n, r, i, s) {
    super(), this.owner = t, this.id = n, this.key = r, this.url = i, this.state = s, this.url = new URL(i, new URL(globalThis.document.baseURI, globalThis.location.href)).href;
  }
  get index() {
    return this.owner.entriesList.findIndex((t) => t.id === this.id);
  }
  get sameDocument() {
    return !0;
  }
  // polyfill is lost between documents
  getState() {
    return this.state;
  }
  setState(t) {
    this.state = t;
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
class Ko {
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
  static getOrCreate(t) {
    return globalThis.polyfea ? globalThis.polyfea : new xt(t);
  }
  /** @static
   * Initialize the polyfea driver in the global context.
   * This method is typically invoked by the polyfea controller script `boot.ts`.
   *
   * @remarks
   * This method also initializes the Navigation polyfill if it's not already present.
   * It augments `window.customElements.define` to allow for duplicate registration of custom elements.
   * This is particularly useful when different microfrontends need to register the same dependencies.
   *
   * @param config - Configuration for the `PolyfeaApi`.
   * For more details, refer to the [`PolyfeaApi`](https://github.com/polyfea/browser-api/blob/main/docs/classes/PolyfeaApi.md) documentation.
   */
  static initialize() {
    globalThis.polyfea || xt.install();
  }
}
class xt {
  constructor(t) {
    this.config = t, this.loadedResources = /* @__PURE__ */ new Set(), globalThis.navigation && globalThis.navigation.addEventListener("navigate", (n) => {
      n.canIntercept && n.destination.url.startsWith(document.baseURI) && n.intercept();
    });
  }
  getBackend() {
    var t;
    if (!this.backend) {
      let n = (t = document.querySelector('meta[name="polyfea-backend"]')) === null || t === void 0 ? void 0 : t.getAttribute("content");
      if (n)
        if (n.startsWith("static://")) {
          const r = n.slice(9);
          this.backend = new Do(this.config || r);
        } else
          this.backend = new en(this.config || n);
      else
        this.backend = new en(this.config || "./polyfea");
    }
    return this.backend;
  }
  getContextArea(t) {
    return globalThis.navigation ? re(globalThis.navigation, "navigatesuccess").pipe(wo(new Event("navigatesuccess", { bubbles: !0, cancelable: !0 })), we((n) => this.getBackend().getContextArea(t)), Ze((n, r) => JSON.stringify(n) === JSON.stringify(r))) : this.getBackend().getContextArea(t).pipe(Ze((n, r) => JSON.stringify(n) === JSON.stringify(r)));
  }
  loadMicrofrontend(t, n) {
    if (!n)
      return Promise.resolve();
    const r = [];
    return this.loadMicrofrontendRecursive(t, n, r);
  }
  async loadMicrofrontendRecursive(t, n, r) {
    if (r.includes(n))
      throw new Error("Circular dependency detected: " + r.join(" -> "));
    const i = t.microfrontends[n];
    if (!i)
      throw new Error("Microfrontend specification not found: " + n);
    r.push(n), i.dependsOn && await Promise.all(i.dependsOn.map((o) => this.loadMicrofrontendRecursive(t, o, r)));
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
  loadScript(t) {
    return new Promise((n, r) => {
      var i;
      const s = document.createElement("script");
      s.src = t.href, s.setAttribute("async", ""), t.attributes && Object.entries(t.attributes).forEach(([a, l]) => {
        s.setAttribute(a, l);
      });
      const o = (i = document.querySelector('meta[name="csp-nonce"]')) === null || i === void 0 ? void 0 : i.getAttribute("content");
      o && s.setAttribute("nonce", o), this.loadedResources.add(t.href), t.waitOnLoad ? s.onload = () => {
        n();
      } : n(), s.onerror = () => {
        this.loadedResources.delete(t.href), r();
      }, document.head.appendChild(s);
    });
  }
  loadStylesheet(t) {
    return this.loadLink(Object.assign(Object.assign({}, t), { attributes: Object.assign(Object.assign({}, t.attributes), { rel: "stylesheet" }) }));
  }
  loadLink(t) {
    var n;
    const r = document.createElement("link");
    r.href = t.href, r.setAttribute("async", "");
    const i = (n = document.querySelector('meta[name="csp-nonce"]')) === null || n === void 0 ? void 0 : n.getAttribute("content");
    return i && r.setAttribute("nonce", i), t.attributes && Object.entries(t.attributes).forEach(([s, o]) => {
      r.setAttribute(s, o);
    }), new Promise((s, o) => {
      this.loadedResources.add(t.href), t.waitOnLoad ? r.onload = () => {
        s();
      } : s(), r.onerror = () => {
        this.loadedResources.delete(t.href), o();
      }, document.head.appendChild(r);
    });
  }
  static install() {
    var t;
    globalThis.polyfea || (globalThis.polyfea = new xt()), Bo();
    const n = (t = document.querySelector('meta[name="polyfea-duplicit-custom-elements"]')) === null || t === void 0 ? void 0 : t.getAttribute("content");
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
const Wo = ":host{display:contents;width:100%;height:100%;box-sizing:border-box}", jo = /* @__PURE__ */ hs(class extends bs {
  constructor() {
    super(), this.__registerHost(), this.__attachShadow(), this.name = void 0, this.take = void 0, this.extraAttributes = {}, this.extraStyle = {}, this.contextObj = void 0;
  }
  async componentWillLoad() {
    this.name && (this.polyfea = Ko.getOrCreate(), this.polyfea.getContextArea(this.name).pipe(
      // load microfrontends
      we((t) => {
        if (!t)
          return Promise.resolve(t);
        let n = (t.elements || []).slice(0, this.take).map((r) => r.microfrontend ? this.polyfea.loadMicrofrontend(t, r.microfrontend) : Promise.resolve());
        return Promise.all(n).then((r) => t);
      })
    ).subscribe({
      next: (t) => this.contextObj = t,
      error: (t) => {
        console.warn(`<polyfe-context name="${this.name}">: Using slotted content because of error: ${t}`);
      }
    }));
  }
  render() {
    var t;
    let n = ((t = this.contextObj) === null || t === void 0 ? void 0 : t.elements) || [];
    return this.take > 0 && (n = n.slice(0, this.take)), et(Mn, null, n.map((r) => this.renderElement(r)), n.length ? "" : et("slot", null));
  }
  renderElement(t) {
    const n = t.tagName, r = Object.assign({
      context: this.name
    }, t.attributes, this.extraAttributes, {
      class: this.name + "-context"
    }), i = Object.assign({}, t.style, this.extraStyle);
    return et(n, { style: i, ref: (s) => {
      if (s)
        for (let o in r)
          s.setAttribute(o, r[o]);
    } });
  }
  static get style() {
    return Wo;
  }
}, [1, "polyfea-context", {
  name: [1],
  take: [2],
  extraAttributes: [16],
  extraStyle: [16],
  contextObj: [32]
}]), Go = jo;
/*!
 * Part of Polyfea microfrontends suite - https://github.com/polyfea
 */
const Yo = (e) => {
  typeof customElements < "u" && [
    Go
  ].forEach((t) => {
    customElements.get(t.is) || customElements.define(t.is, t, e);
  });
};
globalThis.addEventListener("load", () => {
  if (Ui.initialize(), Yo(), !document.body.hasAttribute("polyfea")) {
    document.body.setAttribute("polyfea", "initialized");
    const e = document.createElement("polyfea-context");
    e.setAttribute("name", "shell"), e.setAttribute("take", "1"), document.body.appendChild(e);
  }
});
//# sourceMappingURL=boot.mjs.map
