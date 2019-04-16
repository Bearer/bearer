!(function(e) {
  var r = {}
  function n(o) {
    if (r[o]) return r[o].exports
    var t = (r[o] = { i: o, l: !1, exports: {} })
    return e[o].call(t.exports, t, t.exports, n), (t.l = !0), t.exports
  }
  ;(n.m = e),
    (n.c = r),
    (n.d = function(e, r, o) {
      n.o(e, r) || Object.defineProperty(e, r, { enumerable: !0, get: o })
    }),
    (n.r = function(e) {
      'undefined' != typeof Symbol &&
        Symbol.toStringTag &&
        Object.defineProperty(e, Symbol.toStringTag, { value: 'Module' }),
        Object.defineProperty(e, '__esModule', { value: !0 })
    }),
    (n.t = function(e, r) {
      if ((1 & r && (e = n(e)), 8 & r)) return e
      if (4 & r && 'object' == typeof e && e && e.__esModule) return e
      var o = Object.create(null)
      if ((n.r(o), Object.defineProperty(o, 'default', { enumerable: !0, value: e }), 2 & r && 'string' != typeof e))
        for (var t in e)
          n.d(
            o,
            t,
            function(r) {
              return e[r]
            }.bind(null, t)
          )
      return o
    }),
    (n.n = function(e) {
      var r =
        e && e.__esModule
          ? function() {
              return e.default
            }
          : function() {
              return e
            }
      return n.d(r, 'a', r), r
    }),
    (n.o = function(e, r) {
      return Object.prototype.hasOwnProperty.call(e, r)
    }),
    (n.p = ''),
    n((n.s = 0))
})([
  function(e, r, n) {
    'use strict'
    Object.defineProperty(r, '__esModule', { value: !0 })
    var o = n(1),
      t = n(3),
      i = n(4)
    console.debug('[BEARER]', 'Session uuid', i.Storage.cookieUserId),
      console.debug('[BEARER]', 'Storage uuid', i.Storage.storageUserId),
      i.Storage.ensureCurrentUser(),
      o.send(window.parent, t.Events.COOKIE_SETUP, {
        cookie: '' + document.cookie,
        syncCookies: function(e) {
          document.cookie = e
        }
      }),
      o.send(window.parent, t.Events.SESSION_INITIALIZED),
      o.on(t.Events.HAS_AUTHORIZED, function(e) {
        return (
          console.debug('[BEARER]', 'hasAuthorized?', e.data),
          {
            authorized: i.Storage.hasAuthorized(e.data.integrationId, e.data.integrationId)
          }
        )
      }),
      o.on(t.Events.REVOKE, function(e) {
        i.Storage.revoke(e.data.integrationId, e.data.integrationId), o.send(window.parent, t.Events.REVOKED, e.data)
      })
  },
  function(e, r, n) {
    ;(e.exports = n(2)), (e.exports.default = e.exports)
  },
  function(e, r, n) {
    'undefined' != typeof self && self,
      (e.exports = (function(e) {
        var r = {}
        function n(o) {
          if (r[o]) return r[o].exports
          var t = (r[o] = { i: o, l: !1, exports: {} })
          return e[o].call(t.exports, t, t.exports, n), (t.l = !0), t.exports
        }
        return (
          (n.m = e),
          (n.c = r),
          (n.d = function(e, r, o) {
            n.o(e, r) ||
              Object.defineProperty(e, r, {
                configurable: !1,
                enumerable: !0,
                get: o
              })
          }),
          (n.n = function(e) {
            var r =
              e && e.__esModule
                ? function() {
                    return e.default
                  }
                : function() {
                    return e
                  }
            return n.d(r, 'a', r), r
          }),
          (n.o = function(e, r) {
            return Object.prototype.hasOwnProperty.call(e, r)
          }),
          (n.p = ''),
          n((n.s = './src/index.js'))
        )
      })({
        './node_modules/cross-domain-safe-weakmap/src/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./node_modules/cross-domain-safe-weakmap/src/interface.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = (function(e) {
            if (e && e.__esModule) return e
            var r = {}
            if (null != e) for (var n in e) Object.prototype.hasOwnProperty.call(e, n) && (r[n] = e[n])
            return (r.default = e), r
          })(o)
          r.default = t
        },
        './node_modules/cross-domain-safe-weakmap/src/interface.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./node_modules/cross-domain-safe-weakmap/src/weakmap.js')
          Object.defineProperty(r, 'WeakMap', {
            enumerable: !0,
            get: function() {
              return o.CrossDomainSafeWeakMap
            }
          })
        },
        './node_modules/cross-domain-safe-weakmap/src/native.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.hasNativeWeakMap = function() {
              if (!window.WeakMap) return !1
              if (!window.Object.freeze) return !1
              try {
                var e = new window.WeakMap(),
                  r = {}
                return window.Object.freeze(r), e.set(r, '__testvalue__'), '__testvalue__' === e.get(r)
              } catch (e) {
                return !1
              }
            })
        },
        './node_modules/cross-domain-safe-weakmap/src/util.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.safeIndexOf = function(e, r) {
              for (var n = 0; n < e.length; n++)
                try {
                  if (e[n] === r) return n
                } catch (e) {}
              return -1
            }),
            (r.noop = function() {})
        },
        './node_modules/cross-domain-safe-weakmap/src/weakmap.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.CrossDomainSafeWeakMap = void 0)
          var o = n('./node_modules/cross-domain-utils/src/index.js'),
            t = n('./node_modules/cross-domain-safe-weakmap/src/native.js'),
            i = n('./node_modules/cross-domain-safe-weakmap/src/util.js'),
            s = Object.defineProperty,
            a = Date.now() % 1e9
          r.CrossDomainSafeWeakMap = (function() {
            function e() {
              if (
                ((function(e, r) {
                  if (!(e instanceof r)) throw new TypeError('Cannot call a class as a function')
                })(this, e),
                (a += 1),
                (this.name = '__weakmap_' + ((1e9 * Math.random()) >>> 0) + '__' + a),
                (0, t.hasNativeWeakMap)())
              )
                try {
                  this.weakmap = new window.WeakMap()
                } catch (e) {}
              ;(this.keys = []), (this.values = [])
            }
            return (
              (e.prototype._cleanupClosedWindows = function() {
                for (var e = this.weakmap, r = this.keys, n = 0; n < r.length; n++) {
                  var t = r[n]
                  if ((0, o.isWindow)(t) && (0, o.isWindowClosed)(t)) {
                    if (e)
                      try {
                        e.delete(t)
                      } catch (e) {}
                    r.splice(n, 1), this.values.splice(n, 1), (n -= 1)
                  }
                }
              }),
              (e.prototype.isSafeToReadWrite = function(e) {
                if ((0, o.isWindow)(e)) return !1
                try {
                  ;(0, i.noop)(e && e.self), (0, i.noop)(e && e[this.name])
                } catch (e) {
                  return !1
                }
                return !0
              }),
              (e.prototype.set = function(e, r) {
                if (!e) throw new Error('WeakMap expected key')
                var n = this.weakmap
                if (n)
                  try {
                    n.set(e, r)
                  } catch (e) {
                    delete this.weakmap
                  }
                if (this.isSafeToReadWrite(e)) {
                  var o = this.name,
                    t = e[o]
                  t && t[0] === e ? (t[1] = r) : s(e, o, { value: [e, r], writable: !0 })
                } else {
                  this._cleanupClosedWindows()
                  var a = this.keys,
                    c = this.values,
                    u = (0, i.safeIndexOf)(a, e)
                  ;-1 === u ? (a.push(e), c.push(r)) : (c[u] = r)
                }
              }),
              (e.prototype.get = function(e) {
                if (!e) throw new Error('WeakMap expected key')
                var r = this.weakmap
                if (r)
                  try {
                    if (r.has(e)) return r.get(e)
                  } catch (e) {
                    delete this.weakmap
                  }
                if (!this.isSafeToReadWrite(e)) {
                  this._cleanupClosedWindows()
                  var n = this.keys,
                    o = (0, i.safeIndexOf)(n, e)
                  if (-1 === o) return
                  return this.values[o]
                }
                var t = e[this.name]
                if (t && t[0] === e) return t[1]
              }),
              (e.prototype.delete = function(e) {
                if (!e) throw new Error('WeakMap expected key')
                var r = this.weakmap
                if (r)
                  try {
                    r.delete(e)
                  } catch (e) {
                    delete this.weakmap
                  }
                if (this.isSafeToReadWrite(e)) {
                  var n = e[this.name]
                  n && n[0] === e && (n[0] = n[1] = void 0)
                } else {
                  this._cleanupClosedWindows()
                  var o = this.keys,
                    t = (0, i.safeIndexOf)(o, e)
                  ;-1 !== t && (o.splice(t, 1), this.values.splice(t, 1))
                }
              }),
              (e.prototype.has = function(e) {
                if (!e) throw new Error('WeakMap expected key')
                var r = this.weakmap
                if (r)
                  try {
                    return r.has(e)
                  } catch (e) {
                    delete this.weakmap
                  }
                if (this.isSafeToReadWrite(e)) {
                  var n = e[this.name]
                  return !(!n || n[0] !== e)
                }
                return this._cleanupClosedWindows(), -1 !== (0, i.safeIndexOf)(this.keys, e)
              }),
              e
            )
          })()
        },
        './node_modules/cross-domain-utils/src/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./node_modules/cross-domain-utils/src/utils.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./node_modules/cross-domain-utils/src/types.js')
          Object.keys(t).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return t[e]
                }
              })
          })
        },
        './node_modules/cross-domain-utils/src/types.js': function(e, r, n) {
          'use strict'
        },
        './node_modules/cross-domain-utils/src/util.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.isRegex = function(e) {
              return '[object RegExp]' === Object.prototype.toString.call(e)
            }),
            (r.noop = function() {})
        },
        './node_modules/cross-domain-utils/src/utils.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.isFileProtocol = function() {
              return (
                (arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : window).location.protocol ===
                t.FILE_PROTOCOL
              )
            }),
            (r.isAboutProtocol = s),
            (r.getParent = a),
            (r.getOpener = c),
            (r.canReadFromWindow = u),
            (r.getActualDomain = l),
            (r.getDomain = d),
            (r.isBlankDomain = function(e) {
              try {
                if (!e.location.href) return !0
                if ('about:blank' === e.location.href) return !0
              } catch (e) {}
              return !1
            }),
            (r.isActuallySameDomain = f),
            (r.isSameDomain = m),
            (r.getParents = p),
            (r.isAncestorParent = g),
            (r.getFrames = _),
            (r.getAllChildFrames = h),
            (r.getTop = w),
            (r.getAllFramesInWindow = S),
            (r.isTop = function(e) {
              return e === w(e)
            }),
            (r.isFrameWindowClosed = v),
            (r.isWindowClosed = O),
            (r.linkFrameWindow = function(e) {
              if (
                ((function() {
                  for (var e = 0; e < E.length; e++) v(E[e]) && (E.splice(e, 1), y.splice(e, 1))
                  for (var r = 0; r < y.length; r++) O(y[r]) && (E.splice(r, 1), y.splice(r, 1))
                })(),
                e && e.contentWindow)
              )
                try {
                  y.push(e.contentWindow), E.push(e)
                } catch (e) {}
            }),
            (r.getUserAgent = function(e) {
              return (e = e || window).navigator.mockUserAgent || e.navigator.userAgent
            }),
            (r.getFrameByName = b),
            (r.findChildFrameByName = T),
            (r.findFrameByName = function(e, r) {
              var n = void 0
              return (n = b(e, r)) ? n : T(w(e) || e, r)
            }),
            (r.isParent = function(e, r) {
              var n = a(r)
              if (n) return n === e
              for (var o = _(e), t = Array.isArray(o), i = 0, o = t ? o : o[Symbol.iterator](); ; ) {
                var s
                if (t) {
                  if (i >= o.length) break
                  s = o[i++]
                } else {
                  if ((i = o.next()).done) break
                  s = i.value
                }
                var c = s
                if (c === r) return !0
              }
              return !1
            }),
            (r.isOpener = function(e, r) {
              return e === c(r)
            }),
            (r.getAncestor = A),
            (r.getAncestors = function(e) {
              for (var r = [], n = e; n; ) (n = A(n)) && r.push(n)
              return r
            }),
            (r.isAncestor = function(e, r) {
              var n = A(r)
              if (n) return n === e
              if (r === e) return !1
              if (w(r) === r) return !1
              for (var o = _(e), t = Array.isArray(o), i = 0, o = t ? o : o[Symbol.iterator](); ; ) {
                var s
                if (t) {
                  if (i >= o.length) break
                  s = o[i++]
                } else {
                  if ((i = o.next()).done) break
                  s = i.value
                }
                var a = s
                if (a === r) return !0
              }
              return !1
            }),
            (r.isPopup = j),
            (r.isIframe = P),
            (r.isFullpage = function() {
              return Boolean(!P() && !j())
            }),
            (r.getDistanceFromTop = R),
            (r.getNthParent = C),
            (r.getNthParentFromTop = function(e) {
              var r = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 1
              return C(e, R(e) - r)
            }),
            (r.isSameTopWindow = function(e, r) {
              var n = w(e) || e,
                o = w(r) || r
              try {
                if (n && o) return n === o
              } catch (e) {}
              var t = S(e),
                i = S(r)
              if (N(t, i)) return !0
              var s = c(n),
                a = c(o)
              return !((s && N(S(s), i)) || (a && N(S(a), t), 1))
            }),
            (r.matchDomain = function e(r, n) {
              if ('string' == typeof r) {
                if ('string' == typeof n) return r === t.WILDCARD || n === r
                if ((0, o.isRegex)(n)) return !1
                if (Array.isArray(n)) return !1
              }
              return (0, o.isRegex)(r)
                ? (0, o.isRegex)(n)
                  ? r.toString() === n.toString()
                  : !Array.isArray(n) && Boolean(n.match(r))
                : !!Array.isArray(r) &&
                    (Array.isArray(n)
                      ? JSON.stringify(r) === JSON.stringify(n)
                      : !(0, o.isRegex)(n) &&
                        r.some(function(r) {
                          return e(r, n)
                        }))
            }),
            (r.stringifyDomainPattern = function(e) {
              return Array.isArray(e)
                ? '(' + e.join(' | ') + ')'
                : (0, o.isRegex)(e)
                ? 'RegExp(' + e.toString()
                : e.toString()
            }),
            (r.getDomainFromUrl = function(e) {
              return e.match(/^(https?|mock|file):\/\//)
                ? e
                    .split('/')
                    .slice(0, 3)
                    .join('/')
                : d()
            }),
            (r.onCloseWindow = function(e, r) {
              var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 1e3,
                o = arguments.length > 3 && void 0 !== arguments[3] ? arguments[3] : 1 / 0,
                t = void 0
              return (
                (function i() {
                  if (O(e)) return t && clearTimeout(t), r()
                  o <= 0 ? clearTimeout(t) : ((o -= n), (t = setTimeout(i, n)))
                })(),
                {
                  cancel: function() {
                    t && clearTimeout(t)
                  }
                }
              )
            }),
            (r.isWindow = function(e) {
              try {
                if (e === window) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                if ('[object Window]' === Object.prototype.toString.call(e)) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                if (window.Window && e instanceof window.Window) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                if (e && e.self === e) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                if (e && e.parent === e) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                if (e && e.top === e) return !0
              } catch (e) {
                if (e && e.message === i) return !0
              }
              try {
                ;(0, o.noop)(e == e)
              } catch (e) {
                return !0
              }
              try {
                ;(0, o.noop)(e && e.__cross_domain_utils_window_check__)
              } catch (e) {
                return !0
              }
              return !1
            })
          var o = n('./node_modules/cross-domain-utils/src/util.js'),
            t = {
              MOCK_PROTOCOL: 'mock:',
              FILE_PROTOCOL: 'file:',
              ABOUT_PROTOCOL: 'about:',
              WILDCARD: '*'
            },
            i = 'Call was rejected by callee.\r\n'
          function s() {
            return (
              (arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : window).location.protocol ===
              t.ABOUT_PROTOCOL
            )
          }
          function a(e) {
            if (e)
              try {
                if (e.parent && e.parent !== e) return e.parent
              } catch (e) {}
          }
          function c(e) {
            if (e && !a(e))
              try {
                return e.opener
              } catch (e) {}
          }
          function u(e) {
            try {
              return (0, o.noop)(e && e.location && e.location.href), !0
            } catch (e) {}
            return !1
          }
          function l(e) {
            var r = e.location
            if (!r) throw new Error('Can not read window location')
            var n = r.protocol
            if (!n) throw new Error('Can not read window protocol')
            if (n === t.FILE_PROTOCOL) return t.FILE_PROTOCOL + '//'
            if (n === t.ABOUT_PROTOCOL) {
              var o = a(e)
              return o && u(e) ? l(o) : t.ABOUT_PROTOCOL + '//'
            }
            var i = r.host
            if (!i) throw new Error('Can not read window host')
            return n + '//' + i
          }
          function d(e) {
            var r = l((e = e || window))
            return r && e.mockDomain && 0 === e.mockDomain.indexOf(t.MOCK_PROTOCOL) ? e.mockDomain : r
          }
          function f(e) {
            try {
              if (e === window) return !0
            } catch (e) {}
            try {
              var r = Object.getOwnPropertyDescriptor(e, 'location')
              if (r && !1 === r.enumerable) return !1
            } catch (e) {}
            try {
              if (s(e) && u(e)) return !0
            } catch (e) {}
            try {
              if (l(e) === l(window)) return !0
            } catch (e) {}
            return !1
          }
          function m(e) {
            if (!f(e)) return !1
            try {
              if (e === window) return !0
              if (s(e) && u(e)) return !0
              if (d(window) === d(e)) return !0
            } catch (e) {}
            return !1
          }
          function p(e) {
            var r = []
            try {
              for (; e.parent !== e; ) r.push(e.parent), (e = e.parent)
            } catch (e) {}
            return r
          }
          function g(e, r) {
            if (!e || !r) return !1
            var n = a(r)
            return n ? n === e : -1 !== p(r).indexOf(e)
          }
          function _(e) {
            var r = [],
              n = void 0
            try {
              n = e.frames
            } catch (r) {
              n = e
            }
            var o = void 0
            try {
              o = n.length
            } catch (e) {}
            if (0 === o) return r
            if (o) {
              for (var t = 0; t < o; t++) {
                var i = void 0
                try {
                  i = n[t]
                } catch (e) {
                  continue
                }
                r.push(i)
              }
              return r
            }
            for (var s = 0; s < 100; s++) {
              var a = void 0
              try {
                a = n[s]
              } catch (e) {
                return r
              }
              if (!a) return r
              r.push(a)
            }
            return r
          }
          function h(e) {
            var r = [],
              n = _(e),
              o = Array.isArray(n),
              t = 0
            for (n = o ? n : n[Symbol.iterator](); ; ) {
              var i
              if (o) {
                if (t >= n.length) break
                i = n[t++]
              } else {
                if ((t = n.next()).done) break
                i = t.value
              }
              var s = i
              r.push(s)
              var a = h(s),
                c = Array.isArray(a),
                u = 0
              for (a = c ? a : a[Symbol.iterator](); ; ) {
                var l
                if (c) {
                  if (u >= a.length) break
                  l = a[u++]
                } else {
                  if ((u = a.next()).done) break
                  l = u.value
                }
                var d = l
                r.push(d)
              }
            }
            return r
          }
          function w(e) {
            if (e) {
              try {
                if (e.top) return e.top
              } catch (e) {}
              if (a(e) === e) return e
              try {
                if (g(window, e) && window.top) return window.top
              } catch (e) {}
              try {
                if (g(e, window) && window.top) return window.top
              } catch (e) {}
              var r = h(e),
                n = Array.isArray(r),
                o = 0
              for (r = n ? r : r[Symbol.iterator](); ; ) {
                var t
                if (n) {
                  if (o >= r.length) break
                  t = r[o++]
                } else {
                  if ((o = r.next()).done) break
                  t = o.value
                }
                var i = t
                try {
                  if (i.top) return i.top
                } catch (e) {}
                if (a(i) === i) return i
              }
            }
          }
          function S(e) {
            var r = w(e)
            return h(r).concat(r)
          }
          function v(e) {
            if (!e.contentWindow) return !0
            if (!e.parentNode) return !0
            var r = e.ownerDocument
            return !(!r || !r.body || r.body.contains(e))
          }
          var y = [],
            E = []
          function O(e) {
            var r = !(arguments.length > 1 && void 0 !== arguments[1]) || arguments[1]
            try {
              if (e === window) return !1
            } catch (e) {
              return !0
            }
            try {
              if (!e) return !0
            } catch (e) {
              return !0
            }
            try {
              if (e.closed) return !0
            } catch (e) {
              return !e || e.message !== i
            }
            if (r && m(e))
              try {
                if (e.mockclosed) return !0
              } catch (e) {}
            try {
              if (!e.parent || !e.top) return !0
            } catch (e) {}
            try {
              ;(0, o.noop)(e == e)
            } catch (e) {
              return !0
            }
            var n = (function(e, r) {
              for (var n = 0; n < e.length; n++)
                try {
                  if (e[n] === r) return n
                } catch (e) {}
              return -1
            })(y, e)
            if (-1 !== n) {
              var t = E[n]
              if (t && v(t)) return !0
            }
            return !1
          }
          function b(e, r) {
            var n = _(e),
              o = n,
              t = Array.isArray(o),
              i = 0
            for (o = t ? o : o[Symbol.iterator](); ; ) {
              var s
              if (t) {
                if (i >= o.length) break
                s = o[i++]
              } else {
                if ((i = o.next()).done) break
                s = i.value
              }
              var a = s
              try {
                if (m(a) && a.name === r && -1 !== n.indexOf(a)) return a
              } catch (e) {}
            }
            try {
              if (-1 !== n.indexOf(e.frames[r])) return e.frames[r]
            } catch (e) {}
            try {
              if (-1 !== n.indexOf(e[r])) return e[r]
            } catch (e) {}
          }
          function T(e, r) {
            var n = b(e, r)
            if (n) return n
            var o = _(e),
              t = Array.isArray(o),
              i = 0
            for (o = t ? o : o[Symbol.iterator](); ; ) {
              var s
              if (t) {
                if (i >= o.length) break
                s = o[i++]
              } else {
                if ((i = o.next()).done) break
                s = i.value
              }
              var a = T(s, r)
              if (a) return a
            }
          }
          function A(e) {
            var r = c((e = e || window))
            if (r) return r
            var n = a(e)
            return n || void 0
          }
          function j() {
            return Boolean(c(window))
          }
          function P() {
            return Boolean(a(window))
          }
          function N(e, r) {
            var n = e,
              o = Array.isArray(n),
              t = 0
            for (n = o ? n : n[Symbol.iterator](); ; ) {
              var i
              if (o) {
                if (t >= n.length) break
                i = n[t++]
              } else {
                if ((t = n.next()).done) break
                i = t.value
              }
              var s = i,
                a = r,
                c = Array.isArray(a),
                u = 0
              for (a = c ? a : a[Symbol.iterator](); ; ) {
                var l
                if (c) {
                  if (u >= a.length) break
                  l = a[u++]
                } else {
                  if ((u = a.next()).done) break
                  l = u.value
                }
                if (s === l) return !0
              }
            }
            return !1
          }
          function R() {
            for (var e = 0, r = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : window; r; )
              (r = a(r)) && (e += 1)
            return e
          }
          function C(e) {
            for (var r = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 1, n = e, o = 0; o < r; o++) {
              if (!n) return
              n = a(n)
            }
            return n
          }
        },
        './node_modules/webpack/buildin/global.js': function(e, r, n) {
          'use strict'
          var o,
            t =
              'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
                ? function(e) {
                    return typeof e
                  }
                : function(e) {
                    return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                      ? 'symbol'
                      : typeof e
                  }
          o = (function() {
            return this
          })()
          try {
            o = o || Function('return this')() || (0, eval)('this')
          } catch (e) {
            'object' === ('undefined' == typeof window ? 'undefined' : t(window)) && (o = window)
          }
          e.exports = o
        },
        './node_modules/zalgo-promise/src/exceptions.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.dispatchPossiblyUnhandledError = function(e) {
              if (-1 === (0, o.getGlobal)().dispatchedErrors.indexOf(e)) {
                ;(0, o.getGlobal)().dispatchedErrors.push(e),
                  setTimeout(function() {
                    throw e
                  }, 1)
                for (var r = 0; r < (0, o.getGlobal)().possiblyUnhandledPromiseHandlers.length; r++)
                  (0, o.getGlobal)().possiblyUnhandledPromiseHandlers[r](e)
              }
            }),
            (r.onPossiblyUnhandledException = function(e) {
              return (
                (0, o.getGlobal)().possiblyUnhandledPromiseHandlers.push(e),
                {
                  cancel: function() {
                    ;(0, o.getGlobal)().possiblyUnhandledPromiseHandlers.splice(
                      (0, o.getGlobal)().possiblyUnhandledPromiseHandlers.indexOf(e),
                      1
                    )
                  }
                }
              )
            })
          var o = n('./node_modules/zalgo-promise/src/global.js')
        },
        './node_modules/zalgo-promise/src/global.js': function(e, r, n) {
          'use strict'
          ;(function(e) {
            ;(r.__esModule = !0),
              (r.getGlobal = function() {
                var r = void 0
                if ('undefined' != typeof window) r = window
                else {
                  if (void 0 === e) throw new TypeError('Can not find global')
                  r = e
                }
                var n = (r.__zalgopromise__ = r.__zalgopromise__ || {})
                return (
                  (n.flushPromises = n.flushPromises || []),
                  (n.activeCount = n.activeCount || 0),
                  (n.possiblyUnhandledPromiseHandlers = n.possiblyUnhandledPromiseHandlers || []),
                  (n.dispatchedErrors = n.dispatchedErrors || []),
                  n
                )
              })
          }.call(r, n('./node_modules/webpack/buildin/global.js')))
        },
        './node_modules/zalgo-promise/src/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./node_modules/zalgo-promise/src/promise.js')
          Object.defineProperty(r, 'ZalgoPromise', {
            enumerable: !0,
            get: function() {
              return o.ZalgoPromise
            }
          })
        },
        './node_modules/zalgo-promise/src/promise.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.ZalgoPromise = void 0)
          var o = n('./node_modules/zalgo-promise/src/utils.js'),
            t = n('./node_modules/zalgo-promise/src/exceptions.js'),
            i = n('./node_modules/zalgo-promise/src/global.js'),
            s = (function() {
              function e(r) {
                var n = this
                if (
                  ((function(e, r) {
                    if (!(e instanceof r)) throw new TypeError('Cannot call a class as a function')
                  })(this, e),
                  (this.resolved = !1),
                  (this.rejected = !1),
                  (this.errorHandled = !1),
                  (this.handlers = []),
                  r)
                ) {
                  var o = void 0,
                    t = void 0,
                    i = !1,
                    s = !1,
                    a = !1
                  try {
                    r(
                      function(e) {
                        a ? n.resolve(e) : ((i = !0), (o = e))
                      },
                      function(e) {
                        a ? n.reject(e) : ((s = !0), (t = e))
                      }
                    )
                  } catch (e) {
                    return void this.reject(e)
                  }
                  ;(a = !0), i ? this.resolve(o) : s && this.reject(t)
                }
              }
              return (
                (e.prototype.resolve = function(e) {
                  if (this.resolved || this.rejected) return this
                  if ((0, o.isPromise)(e)) throw new Error('Can not resolve promise with another promise')
                  return (this.resolved = !0), (this.value = e), this.dispatch(), this
                }),
                (e.prototype.reject = function(e) {
                  var r = this
                  if (this.resolved || this.rejected) return this
                  if ((0, o.isPromise)(e)) throw new Error('Can not reject promise with another promise')
                  if (!e) {
                    var n = e && 'function' == typeof e.toString ? e.toString() : Object.prototype.toString.call(e)
                    e = new Error('Expected reject to be called with Error, got ' + n)
                  }
                  return (
                    (this.rejected = !0),
                    (this.error = e),
                    this.errorHandled ||
                      setTimeout(function() {
                        r.errorHandled || (0, t.dispatchPossiblyUnhandledError)(e)
                      }, 1),
                    this.dispatch(),
                    this
                  )
                }),
                (e.prototype.asyncReject = function(e) {
                  ;(this.errorHandled = !0), this.reject(e)
                }),
                (e.prototype.dispatch = function() {
                  var r = this,
                    n = this.dispatching,
                    t = this.resolved,
                    s = this.rejected,
                    a = this.handlers
                  if (!n && (t || s)) {
                    ;(this.dispatching = !0), ((0, i.getGlobal)().activeCount += 1)
                    for (
                      var c = function(n) {
                          var i = a[n],
                            c = i.onSuccess,
                            u = i.onError,
                            l = i.promise,
                            d = void 0
                          if (t)
                            try {
                              d = c ? c(r.value) : r.value
                            } catch (e) {
                              return l.reject(e), 'continue'
                            }
                          else if (s) {
                            if (!u) return l.reject(r.error), 'continue'
                            try {
                              d = u(r.error)
                            } catch (e) {
                              return l.reject(e), 'continue'
                            }
                          }
                          d instanceof e && (d.resolved || d.rejected)
                            ? (d.resolved ? l.resolve(d.value) : l.reject(d.error), (d.errorHandled = !0))
                            : (0, o.isPromise)(d)
                            ? d instanceof e && (d.resolved || d.rejected)
                              ? d.resolved
                                ? l.resolve(d.value)
                                : l.reject(d.error)
                              : d.then(
                                  function(e) {
                                    l.resolve(e)
                                  },
                                  function(e) {
                                    l.reject(e)
                                  }
                                )
                            : l.resolve(d)
                        },
                        u = 0;
                      u < a.length;
                      u++
                    )
                      c(u)
                    ;(a.length = 0),
                      (this.dispatching = !1),
                      ((0, i.getGlobal)().activeCount -= 1),
                      0 === (0, i.getGlobal)().activeCount && e.flushQueue()
                  }
                }),
                (e.prototype.then = function(r, n) {
                  if (r && 'function' != typeof r && !r.call)
                    throw new Error('Promise.then expected a function for success handler')
                  if (n && 'function' != typeof n && !n.call)
                    throw new Error('Promise.then expected a function for error handler')
                  var o = new e()
                  return (
                    this.handlers.push({
                      promise: o,
                      onSuccess: r,
                      onError: n
                    }),
                    (this.errorHandled = !0),
                    this.dispatch(),
                    o
                  )
                }),
                (e.prototype.catch = function(e) {
                  return this.then(void 0, e)
                }),
                (e.prototype.finally = function(r) {
                  return this.then(
                    function(n) {
                      return e.try(r).then(function() {
                        return n
                      })
                    },
                    function(n) {
                      return e.try(r).then(function() {
                        throw n
                      })
                    }
                  )
                }),
                (e.prototype.timeout = function(e, r) {
                  var n = this
                  if (this.resolved || this.rejected) return this
                  var o = setTimeout(function() {
                    n.resolved || n.rejected || n.reject(r || new Error('Promise timed out after ' + e + 'ms'))
                  }, e)
                  return this.then(function(e) {
                    return clearTimeout(o), e
                  })
                }),
                (e.prototype.toPromise = function() {
                  if ('undefined' == typeof Promise) throw new TypeError('Could not find Promise')
                  return Promise.resolve(this)
                }),
                (e.resolve = function(r) {
                  return r instanceof e
                    ? r
                    : (0, o.isPromise)(r)
                    ? new e(function(e, n) {
                        return r.then(e, n)
                      })
                    : new e().resolve(r)
                }),
                (e.reject = function(r) {
                  return new e().reject(r)
                }),
                (e.all = function(r) {
                  var n = new e(),
                    t = r.length,
                    i = []
                  if (!t) return n.resolve(i), n
                  for (
                    var s = function(s) {
                        var a = r[s]
                        if (a instanceof e) {
                          if (a.resolved) return (i[s] = a.value), (t -= 1), 'continue'
                        } else if (!(0, o.isPromise)(a)) return (i[s] = a), (t -= 1), 'continue'
                        e.resolve(a).then(
                          function(e) {
                            ;(i[s] = e), 0 == (t -= 1) && n.resolve(i)
                          },
                          function(e) {
                            n.reject(e)
                          }
                        )
                      },
                      a = 0;
                    a < r.length;
                    a++
                  )
                    s(a)
                  return 0 === t && n.resolve(i), n
                }),
                (e.hash = function(r) {
                  var n = {}
                  return e
                    .all(
                      Object.keys(r).map(function(o) {
                        return e.resolve(r[o]).then(function(e) {
                          n[o] = e
                        })
                      })
                    )
                    .then(function() {
                      return n
                    })
                }),
                (e.map = function(r, n) {
                  return e.all(r.map(n))
                }),
                (e.onPossiblyUnhandledException = function(e) {
                  return (0, t.onPossiblyUnhandledException)(e)
                }),
                (e.try = function(r, n, o) {
                  var t = void 0
                  try {
                    t = r.apply(n, o || [])
                  } catch (r) {
                    return e.reject(r)
                  }
                  return e.resolve(t)
                }),
                (e.delay = function(r) {
                  return new e(function(e) {
                    setTimeout(e, r)
                  })
                }),
                (e.isPromise = function(r) {
                  return !!(r && r instanceof e) || (0, o.isPromise)(r)
                }),
                (e.flush = function() {
                  var r = new e()
                  return (
                    (0, i.getGlobal)().flushPromises.push(r), 0 === (0, i.getGlobal)().activeCount && e.flushQueue(), r
                  )
                }),
                (e.flushQueue = function() {
                  var e = (0, i.getGlobal)().flushPromises
                  ;(0, i.getGlobal)().flushPromises = []
                  var r = e,
                    n = Array.isArray(r),
                    o = 0
                  for (r = n ? r : r[Symbol.iterator](); ; ) {
                    var t
                    if (n) {
                      if (o >= r.length) break
                      t = r[o++]
                    } else {
                      if ((o = r.next()).done) break
                      t = o.value
                    }
                    t.resolve()
                  }
                }),
                e
              )
            })()
          r.ZalgoPromise = s
        },
        './node_modules/zalgo-promise/src/utils.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.isPromise = function(e) {
              try {
                if (!e) return !1
                if ('undefined' != typeof Promise && e instanceof Promise) return !0
                if ('undefined' != typeof window && window.Window && e instanceof window.Window) return !1
                if ('undefined' != typeof window && window.constructor && e instanceof window.constructor) return !1
                var r = {}.toString
                if (r) {
                  var n = r.call(e)
                  if ('[object Window]' === n || '[object global]' === n || '[object DOMWindow]' === n) return !1
                }
                if ('function' == typeof e.then) return !0
              } catch (e) {
                return !1
              }
              return !1
            })
        },
        './src/clean.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.cleanUpWindow = function(e) {
              var r = o.global.requestPromises.get(e)
              if (r)
                for (var n = r, t = Array.isArray(n), i = 0, n = t ? n : n[Symbol.iterator](); ; ) {
                  var s
                  if (t) {
                    if (i >= n.length) break
                    s = n[i++]
                  } else {
                    if ((i = n.next()).done) break
                    s = i.value
                  }
                  var a = s
                  a.reject(new Error('No response from window - cleaned up'))
                }
              o.global.popupWindowsByWin && o.global.popupWindowsByWin.delete(e),
                o.global.remoteWindows && o.global.remoteWindows.delete(e),
                o.global.requestPromises.delete(e),
                o.global.methods.delete(e),
                o.global.readyPromises.delete(e)
            }),
            n('./node_modules/cross-domain-utils/src/index.js')
          var o = n('./src/global.js')
        },
        './src/conf/config.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.CONFIG = void 0)
          var o,
            t = n('./src/conf/constants.js'),
            i = (r.CONFIG = {
              ALLOW_POSTMESSAGE_POPUP: !('__ALLOW_POSTMESSAGE_POPUP__' in window) || window.__ALLOW_POSTMESSAGE_POPUP__,
              LOG_LEVEL: 'info',
              BRIDGE_TIMEOUT: 5e3,
              CHILD_WINDOW_TIMEOUT: 5e3,
              ACK_TIMEOUT: -1 !== window.navigator.userAgent.match(/MSIE/i) ? 2e3 : 1e3,
              RES_TIMEOUT: -1,
              LOG_TO_PAGE: !1,
              ALLOWED_POST_MESSAGE_METHODS: ((o = {}),
              (o[t.CONSTANTS.SEND_STRATEGIES.POST_MESSAGE] = !0),
              (o[t.CONSTANTS.SEND_STRATEGIES.BRIDGE] = !0),
              (o[t.CONSTANTS.SEND_STRATEGIES.GLOBAL] = !0),
              o),
              ALLOW_SAME_ORIGIN: !1
            })
          0 === window.location.href.indexOf(t.CONSTANTS.FILE_PROTOCOL) && (i.ALLOW_POSTMESSAGE_POPUP = !0)
        },
        './src/conf/constants.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.CONSTANTS = {
              POST_MESSAGE_TYPE: {
                REQUEST: 'postrobot_message_request',
                RESPONSE: 'postrobot_message_response',
                ACK: 'postrobot_message_ack'
              },
              POST_MESSAGE_ACK: { SUCCESS: 'success', ERROR: 'error' },
              POST_MESSAGE_NAMES: {
                METHOD: 'postrobot_method',
                HELLO: 'postrobot_ready',
                OPEN_TUNNEL: 'postrobot_open_tunnel'
              },
              WINDOW_TYPES: {
                FULLPAGE: 'fullpage',
                POPUP: 'popup',
                IFRAME: 'iframe'
              },
              WINDOW_PROPS: { POSTROBOT: '__postRobot__' },
              SERIALIZATION_TYPES: {
                METHOD: 'postrobot_method',
                ERROR: 'postrobot_error',
                PROMISE: 'postrobot_promise',
                ZALGO_PROMISE: 'postrobot_zalgo_promise',
                REGEX: 'regex'
              },
              SEND_STRATEGIES: {
                POST_MESSAGE: 'postrobot_post_message',
                BRIDGE: 'postrobot_bridge',
                GLOBAL: 'postrobot_global'
              },
              MOCK_PROTOCOL: 'mock:',
              FILE_PROTOCOL: 'file:',
              BRIDGE_NAME_PREFIX: '__postrobot_bridge__',
              POSTROBOT_PROXY: '__postrobot_proxy__',
              WILDCARD: '*'
            })
          var o = (r.POST_MESSAGE_NAMES = {
            METHOD: 'postrobot_method',
            HELLO: 'postrobot_hello',
            OPEN_TUNNEL: 'postrobot_open_tunnel'
          })
          r.POST_MESSAGE_NAMES_LIST = Object.keys(o).map(function(e) {
            return o[e]
          })
        },
        './src/conf/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./src/conf/config.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./src/conf/constants.js')
          Object.keys(t).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return t[e]
                }
              })
          })
        },
        './src/drivers/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./src/drivers/receive/index.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./src/drivers/send/index.js')
          Object.keys(t).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return t[e]
                }
              })
          })
          var i = n('./src/drivers/listeners.js')
          Object.keys(i).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return i[e]
                }
              })
          })
        },
        './src/drivers/listeners.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.resetListeners = function() {
              ;(i.global.responseListeners = {}), (i.global.requestListeners = {})
            }),
            (r.addResponseListener = function(e, r) {
              i.global.responseListeners[e] = r
            }),
            (r.getResponseListener = function(e) {
              return i.global.responseListeners[e]
            }),
            (r.deleteResponseListener = function(e) {
              delete i.global.responseListeners[e]
            }),
            (r.markResponseListenerErrored = function(e) {
              i.global.erroredResponseListeners[e] = !0
            }),
            (r.isResponseListenerErrored = function(e) {
              return Boolean(i.global.erroredResponseListeners[e])
            }),
            (r.getRequestListener = u),
            (r.addRequestListener = function e(r, n) {
              var t = r.name,
                l = r.win,
                d = r.domain
              if (!t || 'string' != typeof t) throw new Error('Name required to add request listener')
              if (Array.isArray(l)) {
                for (var f = [], m = l, p = Array.isArray(m), g = 0, m = p ? m : m[Symbol.iterator](); ; ) {
                  var _
                  if (p) {
                    if (g >= m.length) break
                    _ = m[g++]
                  } else {
                    if ((g = m.next()).done) break
                    _ = g.value
                  }
                  var h = _
                  f.push(e({ name: t, domain: d, win: h }, n))
                }
                return {
                  cancel: function() {
                    for (var e = f, r = Array.isArray(e), n = 0, e = r ? e : e[Symbol.iterator](); ; ) {
                      var o
                      if (r) {
                        if (n >= e.length) break
                        o = e[n++]
                      } else {
                        if ((n = e.next()).done) break
                        o = n.value
                      }
                      var t = o
                      t.cancel()
                    }
                  }
                }
              }
              if (Array.isArray(d)) {
                for (var w = [], S = d, v = Array.isArray(S), y = 0, S = v ? S : S[Symbol.iterator](); ; ) {
                  var E
                  if (v) {
                    if (y >= S.length) break
                    E = S[y++]
                  } else {
                    if ((y = S.next()).done) break
                    E = y.value
                  }
                  var O = E
                  w.push(e({ name: t, win: l, domain: O }, n))
                }
                return {
                  cancel: function() {
                    for (var e = w, r = Array.isArray(e), n = 0, e = r ? e : e[Symbol.iterator](); ; ) {
                      var o
                      if (r) {
                        if (n >= e.length) break
                        o = e[n++]
                      } else {
                        if ((n = e.next()).done) break
                        o = n.value
                      }
                      var t = o
                      t.cancel()
                    }
                  }
                }
              }
              var b = u({ name: t, win: l, domain: d })
              if (
                ((l && l !== a.CONSTANTS.WILDCARD) || (l = i.global.WINDOW_WILDCARD),
                (d = d || a.CONSTANTS.WILDCARD),
                b)
              )
                throw l && d
                  ? new Error(
                      'Request listener already exists for ' +
                        t +
                        ' on domain ' +
                        d.toString() +
                        ' for ' +
                        (l === i.global.WINDOW_WILDCARD ? 'wildcard' : 'specified') +
                        ' window'
                    )
                  : l
                  ? new Error(
                      'Request listener already exists for ' +
                        t +
                        ' for ' +
                        (l === i.global.WINDOW_WILDCARD ? 'wildcard' : 'specified') +
                        ' window'
                    )
                  : d
                  ? new Error('Request listener already exists for ' + t + ' on domain ' + d.toString())
                  : new Error('Request listener already exists for ' + t)
              var T = i.global.requestListeners,
                A = T[t]
              A || ((A = new o.WeakMap()), (T[t] = A))
              var j = A.get(l)
              j || ((j = {}), A.set(l, j))
              var P = d.toString(),
                N = j[c],
                R = void 0
              return (
                (0, s.isRegex)(d)
                  ? (N || ((N = []), (j[c] = N)), (R = { regex: d, listener: n }), N.push(R))
                  : (j[P] = n),
                {
                  cancel: function() {
                    j && (delete j[P], l && 0 === Object.keys(j).length && A.delete(l), R && N.splice(N.indexOf(R, 1)))
                  }
                }
              )
            }),
            n('./node_modules/zalgo-promise/src/index.js')
          var o = n('./node_modules/cross-domain-safe-weakmap/src/index.js'),
            t = n('./node_modules/cross-domain-utils/src/index.js'),
            i = n('./src/global.js'),
            s = n('./src/lib/index.js'),
            a = n('./src/conf/index.js')
          ;(i.global.responseListeners = i.global.responseListeners || {}),
            (i.global.requestListeners = i.global.requestListeners || {}),
            (i.global.WINDOW_WILDCARD = i.global.WINDOW_WILDCARD || new function() {}()),
            (i.global.erroredResponseListeners = i.global.erroredResponseListeners || {})
          var c = '__domain_regex__'
          function u(e) {
            var r = e.name,
              n = e.win,
              o = e.domain
            if ((n === a.CONSTANTS.WILDCARD && (n = null), o === a.CONSTANTS.WILDCARD && (o = null), !r))
              throw new Error('Name required to get request listener')
            var s = i.global.requestListeners[r]
            if (s)
              for (var u = [n, i.global.WINDOW_WILDCARD], l = 0; l < u.length; l++) {
                var d = u[l],
                  f = d && s.get(d)
                if (f) {
                  if (o && 'string' == typeof o) {
                    if (f[o]) return f[o]
                    if (f[c]) {
                      var m = f[c],
                        p = Array.isArray(m),
                        g = 0
                      for (m = p ? m : m[Symbol.iterator](); ; ) {
                        var _
                        if (p) {
                          if (g >= m.length) break
                          _ = m[g++]
                        } else {
                          if ((g = m.next()).done) break
                          _ = g.value
                        }
                        var h = _,
                          w = h.regex,
                          S = h.listener
                        if ((0, t.matchDomain)(w, o)) return S
                      }
                    }
                  }
                  if (f[a.CONSTANTS.WILDCARD]) return f[a.CONSTANTS.WILDCARD]
                }
              }
          }
        },
        './src/drivers/receive/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o =
            'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
              ? function(e) {
                  return typeof e
                }
              : function(e) {
                  return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                    ? 'symbol'
                    : typeof e
                }
          ;(r.receiveMessage = u),
            (r.messageListener = l),
            (r.listenForMessages = function() {
              ;(0, s.addEventListener)(window, 'message', l)
            })
          var t = n('./node_modules/cross-domain-utils/src/index.js'),
            i = n('./src/conf/index.js'),
            s = n('./src/lib/index.js'),
            a = n('./src/global.js'),
            c = n('./src/drivers/receive/types.js')
          function u(e) {
            if (!window || window.closed) throw new Error('Message recieved in closed window')
            try {
              if (!e.source) return
            } catch (e) {
              return
            }
            var r = e.source,
              n = e.origin,
              u = (function(e) {
                var r = void 0
                try {
                  r = (0, s.jsonParse)(e)
                } catch (e) {
                  return
                }
                if (
                  r &&
                  'object' === (void 0 === r ? 'undefined' : o(r)) &&
                  null !== r &&
                  (r = r[i.CONSTANTS.WINDOW_PROPS.POSTROBOT]) &&
                  'object' === (void 0 === r ? 'undefined' : o(r)) &&
                  null !== r &&
                  r.type &&
                  'string' == typeof r.type &&
                  c.RECEIVE_MESSAGE_TYPES[r.type]
                )
                  return r
              })(e.data)
            if (u) {
              if (!u.sourceDomain || 'string' != typeof u.sourceDomain)
                throw new Error('Expected message to have sourceDomain')
              if (
                ((0 !== u.sourceDomain.indexOf(i.CONSTANTS.MOCK_PROTOCOL) &&
                  0 !== u.sourceDomain.indexOf(i.CONSTANTS.FILE_PROTOCOL)) ||
                  (n = u.sourceDomain),
                -1 === a.global.receivedMessages.indexOf(u.id))
              ) {
                a.global.receivedMessages.push(u.id)
                var l = void 0
                ;(l =
                  -1 !== i.POST_MESSAGE_NAMES_LIST.indexOf(u.name) || u.type === i.CONSTANTS.POST_MESSAGE_TYPE.ACK
                    ? 'debug'
                    : 'error' === u.ack
                    ? 'error'
                    : 'info'),
                  s.log.logLevel(l, [
                    '\n\n\t',
                    '#receive',
                    u.type.replace(/^postrobot_message_/, ''),
                    '::',
                    u.name,
                    '::',
                    n,
                    '\n\n',
                    u
                  ]),
                  !(0, t.isWindowClosed)(r) || u.fireAndForget
                    ? (u.data && (u.data = (0, s.deserializeMethods)(r, n, u.data)),
                      c.RECEIVE_MESSAGE_TYPES[u.type](r, n, u))
                    : s.log.debug('Source window is closed - can not send ' + u.type + ' ' + u.name)
              }
            }
          }
          function l(e) {
            try {
              ;(0, s.noop)(e.source)
            } catch (e) {
              return
            }
            var r = {
              source: e.source || e.sourceElement,
              origin: e.origin || (e.originalEvent && e.originalEvent.origin),
              data: e.data
            }
            u(r)
          }
          ;(a.global.receivedMessages = a.global.receivedMessages || []), (a.global.receiveMessage = u)
        },
        './src/drivers/receive/types.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.RECEIVE_MESSAGE_TYPES = void 0)
          var o,
            t =
              Object.assign ||
              function(e) {
                for (var r = 1; r < arguments.length; r++) {
                  var n = arguments[r]
                  for (var o in n) Object.prototype.hasOwnProperty.call(n, o) && (e[o] = n[o])
                }
                return e
              },
            i = n('./node_modules/zalgo-promise/src/index.js'),
            s = n('./node_modules/cross-domain-utils/src/index.js'),
            a = n('./src/conf/index.js'),
            c = n('./src/lib/index.js'),
            u = n('./src/drivers/send/index.js'),
            l = n('./src/drivers/listeners.js')
          r.RECEIVE_MESSAGE_TYPES = (((o = {})[a.CONSTANTS.POST_MESSAGE_TYPE.ACK] = function(e, r, n) {
            if (!(0, l.isResponseListenerErrored)(n.hash)) {
              var o = (0, l.getResponseListener)(n.hash)
              if (!o)
                throw new Error(
                  'No handler found for post message ack for message: ' +
                    n.name +
                    ' from ' +
                    r +
                    ' in ' +
                    window.location.protocol +
                    '//' +
                    window.location.host +
                    window.location.pathname
                )
              if (!(0, s.matchDomain)(o.domain, r))
                throw new Error('Ack origin ' + r + ' does not match domain ' + o.domain.toString())
              o.ack = !0
            }
          }),
          (o[a.CONSTANTS.POST_MESSAGE_TYPE.REQUEST] = function(e, r, n) {
            var o = (0, l.getRequestListener)({
              name: n.name,
              win: e,
              domain: r
            })
            function d(o) {
              return n.fireAndForget || (0, s.isWindowClosed)(e)
                ? i.ZalgoPromise.resolve()
                : (0, u.sendMessage)(e, t({ target: n.originalSource, hash: n.hash, name: n.name }, o), r)
            }
            return i.ZalgoPromise.all([
              d({ type: a.CONSTANTS.POST_MESSAGE_TYPE.ACK }),
              i.ZalgoPromise.try(function() {
                if (!o)
                  throw new Error(
                    'No handler found for post message: ' +
                      n.name +
                      ' from ' +
                      r +
                      ' in ' +
                      window.location.protocol +
                      '//' +
                      window.location.host +
                      window.location.pathname
                  )
                if (!(0, s.matchDomain)(o.domain, r))
                  throw new Error('Request origin ' + r + ' does not match domain ' + o.domain.toString())
                var t = n.data
                return o.handler({ source: e, origin: r, data: t })
              }).then(
                function(e) {
                  return d({
                    type: a.CONSTANTS.POST_MESSAGE_TYPE.RESPONSE,
                    ack: a.CONSTANTS.POST_MESSAGE_ACK.SUCCESS,
                    data: e
                  })
                },
                function(e) {
                  var r = (0, c.stringifyError)(e).replace(/^Error: /, ''),
                    n = e.code
                  return d({
                    type: a.CONSTANTS.POST_MESSAGE_TYPE.RESPONSE,
                    ack: a.CONSTANTS.POST_MESSAGE_ACK.ERROR,
                    error: r,
                    code: n
                  })
                }
              )
            ])
              .then(c.noop)
              .catch(function(e) {
                if (o && o.handleError) return o.handleError(e)
                c.log.error((0, c.stringifyError)(e))
              })
          }),
          (o[a.CONSTANTS.POST_MESSAGE_TYPE.RESPONSE] = function(e, r, n) {
            if (!(0, l.isResponseListenerErrored)(n.hash)) {
              var o = (0, l.getResponseListener)(n.hash)
              if (!o)
                throw new Error(
                  'No handler found for post message response for message: ' +
                    n.name +
                    ' from ' +
                    r +
                    ' in ' +
                    window.location.protocol +
                    '//' +
                    window.location.host +
                    window.location.pathname
                )
              if (!(0, s.matchDomain)(o.domain, r))
                throw new Error(
                  'Response origin ' + r + ' does not match domain ' + (0, s.stringifyDomainPattern)(o.domain)
                )
              if (((0, l.deleteResponseListener)(n.hash), n.ack === a.CONSTANTS.POST_MESSAGE_ACK.ERROR)) {
                var t = new Error(n.error)
                return n.code && (t.code = n.code), o.respond(t, null)
              }
              if (n.ack === a.CONSTANTS.POST_MESSAGE_ACK.SUCCESS) {
                var i = n.data || n.response
                return o.respond(null, { source: e, origin: r, data: i })
              }
            }
          }),
          o)
        },
        './src/drivers/send/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o =
            Object.assign ||
            function(e) {
              for (var r = 1; r < arguments.length; r++) {
                var n = arguments[r]
                for (var o in n) Object.prototype.hasOwnProperty.call(n, o) && (e[o] = n[o])
              }
              return e
            }
          r.sendMessage = function(e, r, n) {
            return i.ZalgoPromise.try(function() {
              var u
              r = (function(e, r) {
                var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : {},
                  i = (0, a.uniqueID)(),
                  s = (0, a.getWindowType)(),
                  c = (0, t.getDomain)(window)
                return o({}, r, n, {
                  sourceDomain: c,
                  id: r.id || i,
                  windowType: s
                })
              })(e, r, {
                data: (0, a.serializeMethods)(e, n, r.data),
                domain: n
              })
              var l = void 0
              if (
                ((l =
                  -1 !== s.POST_MESSAGE_NAMES_LIST.indexOf(r.name) || r.type === s.CONSTANTS.POST_MESSAGE_TYPE.ACK
                    ? 'debug'
                    : 'error' === r.ack
                    ? 'error'
                    : 'info'),
                a.log.logLevel(l, [
                  '\n\n\t',
                  '#send',
                  r.type.replace(/^postrobot_message_/, ''),
                  '::',
                  r.name,
                  '::',
                  n || s.CONSTANTS.WILDCARD,
                  '\n\n',
                  r
                ]),
                e === window && !s.CONFIG.ALLOW_SAME_ORIGIN)
              )
                throw new Error('Attemping to send message to self')
              if ((0, t.isWindowClosed)(e)) throw new Error('Window is closed')
              a.log.debug('Running send message strategies', r)
              var d = [],
                f = (0, a.jsonStringify)((((u = {})[s.CONSTANTS.WINDOW_PROPS.POSTROBOT] = r), u), null, 2)
              return i.ZalgoPromise.map(Object.keys(c.SEND_MESSAGE_STRATEGIES), function(r) {
                return i.ZalgoPromise.try(function() {
                  if (!s.CONFIG.ALLOWED_POST_MESSAGE_METHODS[r]) throw new Error('Strategy disallowed: ' + r)
                  return c.SEND_MESSAGE_STRATEGIES[r](e, f, n)
                }).then(
                  function() {
                    return d.push(r + ': success'), !0
                  },
                  function(e) {
                    return d.push(r + ': ' + (0, a.stringifyError)(e) + '\n'), !1
                  }
                )
              }).then(function(e) {
                var n = e.some(Boolean),
                  o = r.type + ' ' + r.name + ' ' + (n ? 'success' : 'error') + ':\n  - ' + d.join('\n  - ') + '\n'
                if ((a.log.debug(o), !n)) throw new Error(o)
              })
            })
          }
          var t = n('./node_modules/cross-domain-utils/src/index.js'),
            i = n('./node_modules/zalgo-promise/src/index.js'),
            s = n('./src/conf/index.js'),
            a = n('./src/lib/index.js'),
            c = n('./src/drivers/send/strategies.js')
        },
        './src/drivers/send/strategies.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.SEND_MESSAGE_STRATEGIES = void 0)
          var o = n('./node_modules/cross-domain-utils/src/index.js'),
            t = n('./src/conf/index.js'),
            i = (n('./src/lib/index.js'), (r.SEND_MESSAGE_STRATEGIES = {}))
          i[t.CONSTANTS.SEND_STRATEGIES.POST_MESSAGE] = function(e, r, n) {
            ;(Array.isArray(n) ? n : 'string' == typeof n ? [n] : [t.CONSTANTS.WILDCARD])
              .map(function(r) {
                if (0 === r.indexOf(t.CONSTANTS.MOCK_PROTOCOL)) {
                  if (window.location.protocol === t.CONSTANTS.FILE_PROTOCOL) return t.CONSTANTS.WILDCARD
                  if (!(0, o.isActuallySameDomain)(e))
                    throw new Error(
                      'Attempting to send messsage to mock domain ' + r + ', but window is actually cross-domain'
                    )
                  return (0, o.getActualDomain)(e)
                }
                return 0 === r.indexOf(t.CONSTANTS.FILE_PROTOCOL) ? t.CONSTANTS.WILDCARD : r
              })
              .forEach(function(n) {
                return e.postMessage(r, n)
              })
          }
        },
        './src/global.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.global = void 0)
          var o = n('./src/conf/index.js')
          ;(r.global = window[o.CONSTANTS.WINDOW_PROPS.POSTROBOT] =
            window[o.CONSTANTS.WINDOW_PROPS.POSTROBOT] || {}).registerSelf = function() {}
        },
        './src/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./src/interface.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = (function(e) {
            if (e && e.__esModule) return e
            var r = {}
            if (null != e) for (var n in e) Object.prototype.hasOwnProperty.call(e, n) && (r[n] = e[n])
            return (r.default = e), r
          })(o)
          r.default = t
        },
        './src/interface.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.bridge = r.Promise = r.cleanUpWindow = void 0)
          var o = n('./src/public/index.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./src/clean.js')
          Object.defineProperty(r, 'cleanUpWindow', {
            enumerable: !0,
            get: function() {
              return t.cleanUpWindow
            }
          })
          var i = n('./node_modules/zalgo-promise/src/index.js')
          Object.defineProperty(r, 'Promise', {
            enumerable: !0,
            get: function() {
              return i.ZalgoPromise
            }
          }),
            (r.init = u)
          var s = n('./src/lib/index.js'),
            a = n('./src/drivers/index.js'),
            c = n('./src/global.js')
          function u() {
            c.global.initialized ||
              ((0, a.listenForMessages)(), (0, s.initOnReady)(), (0, s.listenForMethods)({ on: o.on, send: o.send })),
              (c.global.initialized = !0)
          }
          ;(r.bridge = null), u()
        },
        './src/lib/index.js': function(e, r, n) {
          'use strict'
          r.__esModule = !0
          var o = n('./src/lib/util.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./src/lib/log.js')
          Object.keys(t).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return t[e]
                }
              })
          })
          var i = n('./src/lib/serialize.js')
          Object.keys(i).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return i[e]
                }
              })
          })
          var s = n('./src/lib/ready.js')
          Object.keys(s).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return s[e]
                }
              })
          })
        },
        './src/lib/log.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.log = void 0)
          var o =
              'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
                ? function(e) {
                    return typeof e
                  }
                : function(e) {
                    return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                      ? 'symbol'
                      : typeof e
                  },
            t = n('./src/conf/index.js'),
            i = n('./src/lib/util.js'),
            s = ['debug', 'info', 'warn', 'error']
          Function.prototype.bind &&
            window.console &&
            'object' === o(console.log) &&
            ['log', 'info', 'warn', 'error'].forEach(function(e) {
              console[e] = this.bind(console[e], console)
            }, Function.prototype.call)
          var a = (r.log = {
            clearLogs: function() {
              if ((window.console && window.console.clear && window.console.clear(), t.CONFIG.LOG_TO_PAGE)) {
                var e = document.getElementById('postRobotLogs')
                e && e.parentNode && e.parentNode.removeChild(e)
              }
            },
            writeToPage: function(e, r) {
              setTimeout(function() {
                var n = document.getElementById('postRobotLogs')
                n ||
                  (((n = document.createElement('div')).id = 'postRobotLogs'),
                  (n.style.cssText = 'width: 800px; font-family: monospace; white-space: pre-wrap;'),
                  document.body && document.body.appendChild(n))
                var o = document.createElement('div'),
                  t = new Date().toString().split(' ')[4],
                  s = Array.prototype.slice
                    .call(r)
                    .map(function(e) {
                      if ('string' == typeof e) return e
                      if (!e) return Object.prototype.toString.call(e)
                      var r = void 0
                      try {
                        r = (0, i.jsonStringify)(e, null, 2)
                      } catch (e) {
                        r = '[object]'
                      }
                      return '\n\n' + r + '\n\n'
                    })
                    .join(' '),
                  a = t + ' ' + e + ' ' + s
                o.innerHTML = a
                var c = {
                  log: '#ddd',
                  warn: 'orange',
                  error: 'red',
                  info: 'blue',
                  debug: '#aaa'
                }[e]
                ;(o.style.cssText = 'margin-top: 10px; color: ' + c + ';'),
                  n.childNodes.length ? n.insertBefore(o, n.childNodes[0]) : n.appendChild(o)
              })
            },
            logLevel: function(e, r) {
              setTimeout(function() {
                try {
                  var n = window.LOG_LEVEL || t.CONFIG.LOG_LEVEL
                  if ('disabled' === n || s.indexOf(e) < s.indexOf(n)) return
                  if (
                    ((r = Array.prototype.slice.call(r)).unshift('' + window.location.host + window.location.pathname),
                    r.unshift('::'),
                    r.unshift('' + (0, i.getWindowType)().toLowerCase()),
                    r.unshift('[post-robot]'),
                    t.CONFIG.LOG_TO_PAGE && a.writeToPage(e, r),
                    !window.console)
                  )
                    return
                  if ((window.console[e] || (e = 'log'), !window.console[e])) return
                  window.console[e].apply(window.console, r)
                } catch (e) {}
              }, 1)
            },
            debug: function() {
              for (var e = arguments.length, r = Array(e), n = 0; n < e; n++) r[n] = arguments[n]
              a.logLevel('debug', r)
            },
            info: function() {
              for (var e = arguments.length, r = Array(e), n = 0; n < e; n++) r[n] = arguments[n]
              a.logLevel('info', r)
            },
            warn: function() {
              for (var e = arguments.length, r = Array(e), n = 0; n < e; n++) r[n] = arguments[n]
              a.logLevel('warn', r)
            },
            error: function() {
              for (var e = arguments.length, r = Array(e), n = 0; n < e; n++) r[n] = arguments[n]
              a.logLevel('error', r)
            }
          })
        },
        './src/lib/ready.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.onHello = l),
            (r.sayHello = d),
            (r.initOnReady = function() {
              l(function(e) {
                var r = e.source,
                  n = e.origin,
                  o = a.global.readyPromises.get(r) || new i.ZalgoPromise()
                o.resolve({ origin: n }), a.global.readyPromises.set(r, o)
              })
              var e = (0, t.getAncestor)()
              e &&
                d(e).catch(function(e) {
                  c.log.debug((0, u.stringifyError)(e))
                })
            }),
            (r.onChildWindowReady = function(e) {
              var r = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 5e3,
                n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 'Window',
                o = a.global.readyPromises.get(e)
              return (
                o ||
                ((o = new i.ZalgoPromise()),
                a.global.readyPromises.set(e, o),
                -1 !== r &&
                  setTimeout(function() {
                    return o.reject(new Error(n + ' did not load after ' + r + 'ms'))
                  }, r),
                o)
              )
            })
          var o = n('./node_modules/cross-domain-safe-weakmap/src/index.js'),
            t = n('./node_modules/cross-domain-utils/src/index.js'),
            i = n('./node_modules/zalgo-promise/src/index.js'),
            s = n('./src/conf/index.js'),
            a = n('./src/global.js'),
            c = n('./src/lib/log.js'),
            u = n('./src/lib/util.js')
          function l(e) {
            a.global.on(s.CONSTANTS.POST_MESSAGE_NAMES.HELLO, { domain: s.CONSTANTS.WILDCARD }, function(r) {
              var n = r.source,
                o = r.origin
              return e({ source: n, origin: o })
            })
          }
          function d(e) {
            return a.global
              .send(e, s.CONSTANTS.POST_MESSAGE_NAMES.HELLO, {}, { domain: s.CONSTANTS.WILDCARD, timeout: -1 })
              .then(function(e) {
                return { origin: e.origin }
              })
          }
          a.global.readyPromises = a.global.readyPromises || new o.WeakMap()
        },
        './src/lib/serialize.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.listenForMethods = void 0)
          var o =
            'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
              ? function(e) {
                  return typeof e
                }
              : function(e) {
                  return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                    ? 'symbol'
                    : typeof e
                }
          ;(r.serializeMethod = f),
            (r.serializeMethods = function(e, r, n) {
              return (0, u.replaceObject)({ obj: n }, function(n, o) {
                return 'function' == typeof n
                  ? f(e, r, n, o.toString())
                  : n instanceof Error
                  ? ((t = n),
                    {
                      __type__: a.CONSTANTS.SERIALIZATION_TYPES.ERROR,
                      __message__: (0, u.stringifyError)(t),
                      __code__: t.code
                    })
                  : window.Promise && n instanceof window.Promise
                  ? (function(e, r, n, o) {
                      return {
                        __type__: a.CONSTANTS.SERIALIZATION_TYPES.PROMISE,
                        __then__: f(
                          e,
                          r,
                          function(e, r) {
                            return n.then(e, r)
                          },
                          o + '.then'
                        )
                      }
                    })(e, r, n, o.toString())
                  : s.ZalgoPromise.isPromise(n)
                  ? (function(e, r, n, o) {
                      return {
                        __type__: a.CONSTANTS.SERIALIZATION_TYPES.ZALGO_PROMISE,
                        __then__: f(
                          e,
                          r,
                          function(e, r) {
                            return n.then(e, r)
                          },
                          o + '.then'
                        )
                      }
                    })(e, r, n, o.toString())
                  : (0, u.isRegex)(n)
                  ? ((i = n),
                    {
                      __type__: a.CONSTANTS.SERIALIZATION_TYPES.REGEX,
                      __source__: i.source
                    })
                  : void 0
                var t, i
              }).obj
            }),
            (r.deserializeMethod = m),
            (r.deserializeError = p),
            (r.deserializeZalgoPromise = g),
            (r.deserializePromise = _),
            (r.deserializeRegex = h),
            (r.deserializeMethods = function(e, r, n) {
              return (0, u.replaceObject)({ obj: n }, function(n) {
                if ('object' === (void 0 === n ? 'undefined' : o(n)) && null !== n)
                  return d(n, a.CONSTANTS.SERIALIZATION_TYPES.METHOD)
                    ? m(e, r, n)
                    : d(n, a.CONSTANTS.SERIALIZATION_TYPES.ERROR)
                    ? p(0, 0, n)
                    : d(n, a.CONSTANTS.SERIALIZATION_TYPES.PROMISE)
                    ? _(e, r, n)
                    : d(n, a.CONSTANTS.SERIALIZATION_TYPES.ZALGO_PROMISE)
                    ? g(e, r, n)
                    : d(n, a.CONSTANTS.SERIALIZATION_TYPES.REGEX)
                    ? h(0, 0, n)
                    : void 0
              }).obj
            })
          var t = n('./node_modules/cross-domain-safe-weakmap/src/index.js'),
            i = n('./node_modules/cross-domain-utils/src/index.js'),
            s = n('./node_modules/zalgo-promise/src/index.js'),
            a = n('./src/conf/index.js'),
            c = n('./src/global.js'),
            u = n('./src/lib/util.js'),
            l = n('./src/lib/log.js')
          function d(e, r) {
            return 'object' === (void 0 === e ? 'undefined' : o(e)) && null !== e && e.__type__ === r
          }
          function f(e, r, n, o) {
            var t = (0, u.uniqueID)(),
              i = c.global.methods.get(e)
            return (
              i || ((i = {}), c.global.methods.set(e, i)),
              (i[t] = { domain: r, method: n }),
              {
                __type__: a.CONSTANTS.SERIALIZATION_TYPES.METHOD,
                __id__: t,
                __name__: o
              }
            )
          }
          function m(e, r, n) {
            function o() {
              var o = Array.prototype.slice.call(arguments)
              return (
                l.log.debug('Call foreign method', n.__name__, o),
                c.global
                  .send(
                    e,
                    a.CONSTANTS.POST_MESSAGE_NAMES.METHOD,
                    { id: n.__id__, name: n.__name__, args: o },
                    { domain: r, timeout: -1 }
                  )
                  .then(
                    function(e) {
                      var r = e.data
                      return l.log.debug('Got foreign method result', n.__name__, r.result), r.result
                    },
                    function(e) {
                      throw (l.log.debug('Got foreign method error', (0, u.stringifyError)(e)), e)
                    }
                  )
              )
            }
            return (o.__name__ = n.__name__), (o.__xdomain__ = !0), (o.source = e), (o.origin = r), o
          }
          function p(e, r, n) {
            var o = new Error(n.__message__)
            return n.__code__ && (o.code = n.__code__), o
          }
          function g(e, r, n) {
            return new s.ZalgoPromise(function(o, t) {
              return m(e, r, n.__then__)(o, t)
            })
          }
          function _(e, r, n) {
            return window.Promise
              ? new window.Promise(function(o, t) {
                  return m(e, r, n.__then__)(o, t)
                })
              : g(e, r, n)
          }
          function h(e, r, n) {
            return new RegExp(n.__source__)
          }
          ;(c.global.methods = c.global.methods || new t.WeakMap()),
            (r.listenForMethods = (0, u.once)(function() {
              c.global.on(a.CONSTANTS.POST_MESSAGE_NAMES.METHOD, { origin: a.CONSTANTS.WILDCARD }, function(e) {
                var r = e.source,
                  n = e.origin,
                  o = e.data,
                  t = c.global.methods.get(r)
                if (!t) throw new Error('Could not find any methods this window has privileges to call')
                var a = t[o.id]
                if (!a) throw new Error('Could not find method with id: ' + o.id)
                if (!(0, i.matchDomain)(a.domain, n))
                  throw new Error('Method domain ' + a.domain + ' does not match origin ' + n)
                return (
                  l.log.debug('Call local method', o.name, o.args),
                  s.ZalgoPromise.try(function() {
                    return a.method.apply({ source: r, origin: n, data: o }, o.args)
                  }).then(function(e) {
                    return { result: e, id: o.id, name: o.name }
                  })
                )
              })
            }))
        },
        './src/lib/util.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.weakMapMemoize = r.once = void 0)
          var o =
            'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
              ? function(e) {
                  return typeof e
                }
              : function(e) {
                  return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                    ? 'symbol'
                    : typeof e
                }
          ;(r.stringifyError = function e(r) {
            var n = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 1
            if (n >= 3) return 'stringifyError stack overflow'
            try {
              if (!r) return '<unknown error: ' + Object.prototype.toString.call(r) + '>'
              if ('string' == typeof r) return r
              if (r instanceof Error) {
                var o = r && r.stack,
                  t = r && r.message
                if (o && t) return -1 !== o.indexOf(t) ? o : t + '\n' + o
                if (o) return o
                if (t) return t
              }
              return 'function' == typeof r.toString ? r.toString() : Object.prototype.toString.call(r)
            } catch (r) {
              return 'Error while stringifying error: ' + e(r, n + 1)
            }
          }),
            (r.noop = function() {}),
            (r.addEventListener = function(e, r, n) {
              return (
                e.addEventListener ? e.addEventListener(r, n) : e.attachEvent('on' + r, n),
                {
                  cancel: function() {
                    e.removeEventListener ? e.removeEventListener(r, n) : e.detachEvent('on' + r, n)
                  }
                }
              )
            }),
            (r.uniqueID = function() {
              var e = '0123456789abcdef'
              return 'xxxxxxxxxx'.replace(/./g, function() {
                return e.charAt(Math.floor(Math.random() * e.length))
              })
            }),
            (r.eachArray = a),
            (r.eachObject = c),
            (r.each = u),
            (r.replaceObject = function e(r, n) {
              var t = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 1
              if (t >= 100) throw new Error('Self-referential object passed, or object contained too many layers')
              var i = void 0
              if ('object' !== (void 0 === r ? 'undefined' : o(r)) || null === r || Array.isArray(r)) {
                if (!Array.isArray(r)) throw new TypeError('Invalid type: ' + (void 0 === r ? 'undefined' : o(r)))
                i = []
              } else i = {}
              return (
                u(r, function(r, s) {
                  var a = n(r, s)
                  void 0 !== a
                    ? (i[s] = a)
                    : 'object' === (void 0 === r ? 'undefined' : o(r)) && null !== r
                    ? (i[s] = e(r, n, t + 1))
                    : (i[s] = r)
                }),
                i
              )
            }),
            (r.safeInterval = function(e, r) {
              var n = void 0
              return (
                (n = setTimeout(function o() {
                  ;(n = setTimeout(o, r)), e.call()
                }, r)),
                {
                  cancel: function() {
                    clearTimeout(n)
                  }
                }
              )
            }),
            (r.isRegex = function(e) {
              return '[object RegExp]' === Object.prototype.toString.call(e)
            }),
            (r.getWindowType = function() {
              return (0, i.isPopup)()
                ? s.CONSTANTS.WINDOW_TYPES.POPUP
                : (0, i.isIframe)()
                ? s.CONSTANTS.WINDOW_TYPES.IFRAME
                : s.CONSTANTS.WINDOW_TYPES.FULLPAGE
            }),
            (r.jsonStringify = function(e, r, n) {
              var o = void 0,
                t = void 0
              try {
                if (
                  ('{}' !== JSON.stringify({}) && ((o = Object.prototype.toJSON), delete Object.prototype.toJSON),
                  '{}' !== JSON.stringify({}))
                )
                  throw new Error('Can not correctly serialize JSON objects')
                if (
                  ('[]' !== JSON.stringify([]) && ((t = Array.prototype.toJSON), delete Array.prototype.toJSON),
                  '[]' !== JSON.stringify([]))
                )
                  throw new Error('Can not correctly serialize JSON objects')
              } catch (e) {
                throw new Error('Can not repair JSON.stringify: ' + e.message)
              }
              var i = JSON.stringify.call(this, e, r, n)
              try {
                o && (Object.prototype.toJSON = o), t && (Array.prototype.toJSON = t)
              } catch (e) {
                throw new Error('Can not repair JSON.stringify: ' + e.message)
              }
              return i
            }),
            (r.jsonParse = function(e) {
              return JSON.parse(e)
            }),
            (r.needsGlobalMessagingForBrowser = function() {
              return (
                !!(0, i.getUserAgent)(window).match(/MSIE|trident|edge\/12|edge\/13/i) ||
                !s.CONFIG.ALLOW_POSTMESSAGE_POPUP
              )
            })
          var t = n('./node_modules/cross-domain-safe-weakmap/src/index.js'),
            i = n('./node_modules/cross-domain-utils/src/index.js'),
            s = n('./src/conf/index.js')
          function a(e, r) {
            for (var n = 0; n < e.length; n++) r(e[n], n)
          }
          function c(e, r) {
            for (var n in e) e.hasOwnProperty(n) && r(e[n], n)
          }
          function u(e, r) {
            Array.isArray(e) ? a(e, r) : 'object' === (void 0 === e ? 'undefined' : o(e)) && null !== e && c(e, r)
          }
          ;(r.once = function(e) {
            if (!e) return e
            var r = !1
            return function() {
              if (!r) return (r = !0), e.apply(this, arguments)
            }
          }),
            (r.weakMapMemoize = function(e) {
              var r = new t.WeakMap()
              return function(n) {
                var o = r.get(n)
                return void 0 !== o ? o : (void 0 !== (o = e.call(this, n)) && r.set(n, o), o)
              }
            })
        },
        './src/public/client.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0),
            (r.send = void 0),
            (r.request = l),
            (r.sendToParent = function(e, r, n) {
              var o = (0, i.getAncestor)()
              return o
                ? d(o, e, r, n)
                : new t.ZalgoPromise(function(e, r) {
                    return r(new Error('Window does not have a parent'))
                  })
            }),
            (r.client = function() {
              var e = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {}
              if (!e.window) throw new Error('Expected options.window')
              var r = e.window
              return {
                send: function(n, o) {
                  return d(r, n, o, e)
                }
              }
            })
          var o = n('./node_modules/cross-domain-safe-weakmap/src/index.js'),
            t = n('./node_modules/zalgo-promise/src/index.js'),
            i = n('./node_modules/cross-domain-utils/src/index.js'),
            s = n('./src/conf/index.js'),
            a = n('./src/drivers/index.js'),
            c = n('./src/lib/index.js'),
            u = n('./src/global.js')
          function l(e) {
            return t.ZalgoPromise.try(function() {
              if (!e.name) throw new Error('Expected options.name')
              var r = e.name,
                n = void 0,
                o = void 0
              if ('string' == typeof e.window) {
                var l = document.getElementById(e.window)
                if (!l)
                  throw new Error(
                    'Expected options.window ' + Object.prototype.toString.call(e.window) + ' to be a valid element id'
                  )
                if ('iframe' !== l.tagName.toLowerCase())
                  throw new Error(
                    'Expected options.window ' + Object.prototype.toString.call(e.window) + ' to be an iframe'
                  )
                if (!l.contentWindow)
                  throw new Error(
                    'Iframe must have contentWindow.  Make sure it has a src attribute and is in the DOM.'
                  )
                n = l.contentWindow
              } else if (e.window instanceof HTMLIFrameElement) {
                if ('iframe' !== e.window.tagName.toLowerCase())
                  throw new Error(
                    'Expected options.window ' + Object.prototype.toString.call(e.window) + ' to be an iframe'
                  )
                if (e.window && !e.window.contentWindow)
                  throw new Error(
                    'Iframe must have contentWindow.  Make sure it has a src attribute and is in the DOM.'
                  )
                e.window && e.window.contentWindow && (n = e.window.contentWindow)
              } else n = e.window
              if (!n) throw new Error('Expected options.window to be a window object, iframe, or iframe element id.')
              var d = n
              o = e.domain || s.CONSTANTS.WILDCARD
              var f = e.name + '_' + (0, c.uniqueID)()
              if ((0, i.isWindowClosed)(d)) throw new Error('Target window is closed')
              var m = !1,
                p = u.global.requestPromises.get(d)
              p || ((p = []), u.global.requestPromises.set(d, p))
              var g = t.ZalgoPromise.try(function() {
                if ((0, i.isAncestor)(window, d))
                  return (0, c.onChildWindowReady)(d, e.timeout || s.CONFIG.CHILD_WINDOW_TIMEOUT)
              })
                .then(function() {
                  var e = (arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {}).origin
                  if ((0, c.isRegex)(o) && !e) return (0, c.sayHello)(d)
                })
                .then(function() {
                  var n = (arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {}).origin
                  if ((0, c.isRegex)(o)) {
                    if (!(0, i.matchDomain)(o, n))
                      throw new Error('Remote window domain ' + n + ' does not match regex: ' + o.toString())
                    o = n
                  }
                  if ('string' != typeof o && !Array.isArray(o))
                    throw new TypeError('Expected domain to be a string or array')
                  var u = o
                  return new t.ZalgoPromise(function(n, o) {
                    var t = void 0
                    if (
                      (e.fireAndForget ||
                        ((t = {
                          name: r,
                          window: d,
                          domain: u,
                          respond: function(e, r) {
                            e || ((m = !0), p.splice(p.indexOf(g, 1))), e ? o(e) : n(r)
                          }
                        }),
                        (0, a.addResponseListener)(f, t)),
                      (0, a.sendMessage)(
                        d,
                        {
                          type: s.CONSTANTS.POST_MESSAGE_TYPE.REQUEST,
                          hash: f,
                          name: r,
                          data: e.data,
                          fireAndForget: e.fireAndForget
                        },
                        u
                      ).catch(o),
                      e.fireAndForget)
                    )
                      return n()
                    var c = s.CONFIG.ACK_TIMEOUT,
                      l = e.timeout || s.CONFIG.RES_TIMEOUT,
                      _ = 100
                    setTimeout(function n() {
                      if (!m) {
                        if ((0, i.isWindowClosed)(d))
                          return t.ack
                            ? o(new Error('Window closed for ' + r + ' before response'))
                            : o(new Error('Window closed for ' + r + ' before ack'))
                        if (((c = Math.max(c - _, 0)), -1 !== l && (l = Math.max(l - _, 0)), t.ack)) {
                          if (-1 === l) return
                          _ = Math.min(l, 2e3)
                        } else {
                          if (0 === c)
                            return o(
                              new Error(
                                'No ack for postMessage ' +
                                  r +
                                  ' in ' +
                                  (0, i.getDomain)() +
                                  ' in ' +
                                  s.CONFIG.ACK_TIMEOUT +
                                  'ms'
                              )
                            )
                          if (0 === l)
                            return o(
                              new Error(
                                'No response for postMessage ' +
                                  r +
                                  ' in ' +
                                  (0, i.getDomain)() +
                                  ' in ' +
                                  (e.timeout || s.CONFIG.RES_TIMEOUT) +
                                  'ms'
                              )
                            )
                        }
                        setTimeout(n, _)
                      }
                    }, _)
                  })
                })
              return (
                g.catch(function() {
                  ;(0, a.markResponseListenerErrored)(f), (0, a.deleteResponseListener)(f)
                }),
                p.push(g),
                g
              )
            })
          }
          function d(e, r, n, o) {
            return ((o = o || {}).window = e), (o.name = r), (o.data = n), l(o)
          }
          ;(u.global.requestPromises = u.global.requestPromises || new o.WeakMap()), (r.send = d), (u.global.send = d)
        },
        './src/public/config.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.CONSTANTS = r.CONFIG = void 0)
          var o = n('./src/conf/index.js')
          Object.defineProperty(r, 'CONFIG', {
            enumerable: !0,
            get: function() {
              return o.CONFIG
            }
          }),
            Object.defineProperty(r, 'CONSTANTS', {
              enumerable: !0,
              get: function() {
                return o.CONSTANTS
              }
            }),
            (r.disable = function() {
              delete window[o.CONSTANTS.WINDOW_PROPS.POSTROBOT],
                window.removeEventListener('message', t.messageListener)
            })
          var t = n('./src/drivers/index.js')
        },
        './src/public/index.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.parent = void 0)
          var o = n('./src/public/client.js')
          Object.keys(o).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return o[e]
                }
              })
          })
          var t = n('./src/public/server.js')
          Object.keys(t).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return t[e]
                }
              })
          })
          var i = n('./src/public/config.js')
          Object.keys(i).forEach(function(e) {
            'default' !== e &&
              '__esModule' !== e &&
              Object.defineProperty(r, e, {
                enumerable: !0,
                get: function() {
                  return i[e]
                }
              })
          })
          var s = n('./node_modules/cross-domain-utils/src/index.js')
          r.parent = (0, s.getAncestor)()
        },
        './src/public/server.js': function(e, r, n) {
          'use strict'
          ;(r.__esModule = !0), (r.on = void 0)
          var o =
            'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator
              ? function(e) {
                  return typeof e
                }
              : function(e) {
                  return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype
                    ? 'symbol'
                    : typeof e
                }
          ;(r.listen = l),
            (r.once = function(e) {
              var r = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : {},
                n = arguments[2]
              'function' == typeof r && ((n = r), (r = {})), (r = r || {}), (n = n || r.handler)
              var o = r.errorHandler,
                t = new i.ZalgoPromise(function(t, i) {
                  ;((r = r || {}).name = e),
                    (r.once = !0),
                    (r.handler = function(e) {
                      if ((t(e), n)) return n(e)
                    }),
                    (r.errorHandler = function(e) {
                      if ((i(e), o)) return o(e)
                    })
                }),
                s = l(r)
              return (t.cancel = s.cancel), t
            }),
            (r.listener = function() {
              var e = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {}
              return {
                on: function(r, n) {
                  return d(r, e, n)
                }
              }
            })
          var t = n('./node_modules/cross-domain-utils/src/index.js'),
            i = n('./node_modules/zalgo-promise/src/index.js'),
            s = n('./src/lib/index.js'),
            a = n('./src/drivers/index.js'),
            c = n('./src/conf/index.js'),
            u = n('./src/global.js')
          function l(e) {
            if (!e.name) throw new Error('Expected options.name')
            if (!e.handler) throw new Error('Expected options.handler')
            var r = e.name,
              n = e.window,
              i = e.domain,
              u = {
                handler: e.handler,
                handleError:
                  e.errorHandler ||
                  function(e) {
                    throw e
                  },
                window: n,
                domain: i || c.CONSTANTS.WILDCARD,
                name: r
              },
              l = (0, a.addRequestListener)({ name: r, win: n, domain: i }, u)
            if (e.once) {
              var d = u.handler
              u.handler = (0, s.once)(function() {
                return l.cancel(), d.apply(this, arguments)
              })
            }
            if (u.window && e.errorOnClose)
              var f = (0, s.safeInterval)(function() {
                n &&
                  'object' === (void 0 === n ? 'undefined' : o(n)) &&
                  (0, t.isWindowClosed)(n) &&
                  (f.cancel(), u.handleError(new Error('Post message target window is closed')))
              }, 50)
            return {
              cancel: function() {
                l.cancel()
              }
            }
          }
          function d(e, r, n) {
            return (
              'function' == typeof r && ((n = r), (r = {})),
              ((r = r || {}).name = e),
              (r.handler = n || r.handler),
              l(r)
            )
          }
          ;(r.on = d), (u.global.on = d)
        }
      }))
  },
  function(e, r, n) {
    'use strict'
    Object.defineProperty(r, '__esModule', { value: !0 }),
      (function(e) {
        ;(e.SESSION_INITIALIZED = 'BEARER_SESSION_INITIALIZED'),
          (e.COOKIE_SETUP = 'BEARER_COOKIE_SETUP'),
          (e.HAS_AUTHORIZED = 'BEARER_HAS_AUTHORIZED'),
          (e.AUTHORIZED = 'BEARER_AUTHORIZED'),
          (e.REVOKE = 'BEARER_REVOKE'),
          (e.REVOKED = 'BEARER_REVOKED')
      })(r.Events || (r.Events = {}))
  },
  function(e, r, n) {
    'use strict'
    Object.defineProperty(r, '__esModule', { value: !0 })
    var o = (function() {
      function e() {}
      return (
        Object.defineProperty(e, 'cookieUserId', {
          get: function() {
            return (
              (function(e) {
                var r = ('; ' + document.cookie).split('; ' + e + '=')
                if (2 == r.length)
                  return r
                    .pop()
                    .split(';')
                    .shift()
              })('uuid') || ''
            )
          },
          enumerable: !0,
          configurable: !0
        }),
        Object.defineProperty(e, 'storageUserId', {
          get: function() {
            return localStorage.getItem('uuid')
          },
          enumerable: !0,
          configurable: !0
        }),
        Object.defineProperty(e, 'isNewUser', {
          get: function() {
            return this.cookieUserId !== this.storageUserId
          },
          enumerable: !0,
          configurable: !0
        }),
        (e.clearStorage = function() {
          console.debug('[BEARER]', 'clearing Storage'),
            localStorage.clear(),
            localStorage.setItem('uuid', this.cookieUserId)
        }),
        (e.authorize = function(e, r) {
          localStorage.setItem(t(e, r), 'true')
        }),
        (e.revoke = function(e, r) {
          localStorage.removeItem(t(e, r))
        }),
        (e.hasAuthorized = function(e, r) {
          return !this.isNewUser && 'true' === localStorage.getItem(t(e, r))
        }),
        (e.ensureCurrentUser = function() {
          this.isNewUser && this.clearStorage()
        }),
        e
      )
    })()
    function t(e, r) {
      return [e, r].join('|')
    }
    r.Storage = o
  }
])
