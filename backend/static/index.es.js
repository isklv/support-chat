(function(){"use strict";try{if(typeof document<"u"){var e=document.createElement("style");e.appendChild(document.createTextNode(":root{--sc-primary: #2563eb;--sc-primary-hover: #1d4ed8;--sc-bg: #ffffff;--sc-bg-subtle: #f8fafc;--sc-card: #f1f5f9;--sc-border: #e2e8f0;--sc-text: #0f172a;--sc-text-muted: #64748b;--sc-user-bubble-bg: #eff6ff;--sc-user-bubble-border: #bfdbfe;--sc-shadow: 0 20px 25px -5px rgba(0, 0, 0, .1), 0 8px 10px -6px rgba(0, 0, 0, .1)}@media (prefers-color-scheme: dark){:root{--sc-bg: #0f172a;--sc-bg-subtle: #1e293b;--sc-card: #1e293b;--sc-border: #334155;--sc-text: #f8fafc;--sc-text-muted: #94a3b8;--sc-user-bubble-bg: #1e3a8a;--sc-user-bubble-border: #2563eb;--sc-shadow: 0 25px 35px -5px rgba(0, 0, 0, .5)}}.sc-launcher{position:fixed;bottom:24px;right:24px;width:56px;height:56px;border-radius:50%;background:var(--sc-primary);color:#fff;border:none;cursor:pointer;box-shadow:0 10px 25px -5px #2563eb66;display:flex;align-items:center;justify-content:center;z-index:999990;transition:transform .2s cubic-bezier(.16,1,.3,1),background-color .15s;outline:none;padding:0}.sc-launcher.sc-pos-left{right:auto;left:24px}.sc-launcher:hover{transform:scale(1.06);background:var(--sc-primary-hover)}.sc-launcher:active{transform:scale(.96)}.sc-badge{position:absolute;top:-2px;right:-2px;background:#ef4444;color:#fff;font-size:11px;font-weight:700;min-width:20px;height:20px;line-height:20px;border-radius:10px;text-align:center;padding:0 4px;border:2px solid var(--sc-bg);box-sizing:border-box;animation:scPulse 2s infinite}@keyframes scPulse{0%,to{transform:scale(1)}50%{transform:scale(1.1)}}.sc-panel{position:fixed;bottom:92px;right:24px;width:380px;height:560px;max-width:calc(100vw - 32px);max-height:calc(100vh - 120px);background:var(--sc-bg);border:1px solid var(--sc-border);border-radius:16px;box-shadow:var(--sc-shadow);display:flex;flex-direction:column;overflow:hidden;z-index:999995;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif;animation:scSlideUp .22s cubic-bezier(.16,1,.3,1);box-sizing:border-box}.sc-panel.sc-pos-left{right:auto;left:24px}@keyframes scSlideUp{0%{opacity:0;transform:translateY(16px) scale(.96)}to{opacity:1;transform:translateY(0) scale(1)}}.sc-header{padding:14px 16px;background:var(--sc-bg);border-bottom:1px solid var(--sc-border);display:flex;align-items:center;gap:12px}.sc-header-icon{width:36px;height:36px;border-radius:10px;background:var(--sc-primary);color:#fff;display:flex;align-items:center;justify-content:center;flex-shrink:0}.sc-header-text{flex:1;min-width:0}.sc-title{margin:0;font-size:15px;font-weight:600;color:var(--sc-text);line-height:1.2}.sc-subtitle{margin:3px 0 0;font-size:12px;color:var(--sc-text-muted);line-height:1.2;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}.sc-close-btn{background:transparent;border:none;cursor:pointer;padding:6px;color:var(--sc-text-muted);border-radius:8px;display:flex;align-items:center;justify-content:center;transition:background-color .15s,color .15s}.sc-close-btn:hover{background:var(--sc-card);color:var(--sc-text)}.sc-messages{flex:1;overflow-y:auto;padding:16px;display:flex;flex-direction:column;gap:10px;list-style:none;margin:0;box-sizing:border-box}.sc-empty{margin:auto 0;text-align:center;color:var(--sc-text-muted);font-size:13px;padding:20px;line-height:1.5}.sc-bubble{max-width:82%;padding:9px 13px;border-radius:14px;font-size:13.5px;line-height:1.42;display:flex;flex-direction:column;gap:4px;box-sizing:border-box}.sc-bubble.sc-mine{align-self:flex-end;background:var(--sc-user-bubble-bg);border:1px solid var(--sc-user-bubble-border);border-bottom-right-radius:4px;color:var(--sc-text)}.sc-bubble.sc-theirs{align-self:flex-start;background:var(--sc-card);border:1px solid var(--sc-border);border-bottom-left-radius:4px;color:var(--sc-text)}.sc-who{font-size:11px;font-weight:700;color:var(--sc-primary);margin-bottom:2px}.sc-bubble p{margin:0;white-space:pre-wrap;word-break:break-word}.sc-time{font-size:10px;color:var(--sc-text-muted);align-self:flex-end;margin-top:2px}.sc-image{max-width:100%;max-height:180px;border-radius:8px;object-fit:cover;cursor:pointer;display:block}.sc-note{align-self:center;font-size:11px;color:var(--sc-text-muted);text-align:center;padding:4px 10px;background:var(--sc-bg-subtle);border-radius:12px;border:1px solid var(--sc-border)}.sc-form{padding:12px 14px;border-top:1px solid var(--sc-border);background:var(--sc-bg);display:flex;flex-direction:column;gap:8px}.sc-contacts{display:flex;gap:8px}.sc-input{flex:1;min-width:0;padding:8px 10px;border:1px solid var(--sc-border);border-radius:8px;font-size:12.5px;color:var(--sc-text);background:var(--sc-bg-subtle);outline:none;box-sizing:border-box}.sc-input:focus{border-color:var(--sc-primary);background:var(--sc-bg)}.sc-hp{position:absolute;left:-9999px;width:1px;height:1px;opacity:0}.sc-compose{display:flex;align-items:flex-end;gap:8px}.sc-textarea{flex:1;min-width:0;padding:8px 10px;border:1px solid var(--sc-border);border-radius:8px;font-size:13px;color:var(--sc-text);background:var(--sc-bg-subtle);outline:none;resize:none;max-height:100px;font-family:inherit;box-sizing:border-box;line-height:1.35}.sc-textarea:focus{border-color:var(--sc-primary);background:var(--sc-bg)}.sc-action-btn{background:transparent;border:none;cursor:pointer;padding:8px;border-radius:8px;color:var(--sc-text-muted);display:flex;align-items:center;justify-content:center;transition:background-color .15s,color .15s}.sc-action-btn:hover{background:var(--sc-card);color:var(--sc-text)}.sc-send-btn{background:var(--sc-primary);color:#fff;border:none;border-radius:8px;cursor:pointer;padding:8px 12px;display:flex;align-items:center;justify-content:center;transition:background-color .15s}.sc-send-btn:hover:not(:disabled){background:var(--sc-primary-hover)}.sc-send-btn:disabled{opacity:.5;cursor:not-allowed}.sc-attached-chip{display:flex;align-items:center;gap:6px;font-size:11.5px;color:var(--sc-primary);background:var(--sc-card);padding:4px 8px;border-radius:6px;width:fit-content;max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.sc-error{font-size:11px;color:#ef4444;margin:0}.sc-foot{font-size:10.5px;color:var(--sc-text-muted);margin:0;text-align:center}@media (max-width: 600px){.sc-panel{bottom:0;right:0;left:0;width:100%;height:100%;max-width:100%;max-height:100%;border-radius:0;border:none}}")),document.head.appendChild(e)}}catch(r){console.error("vite-plugin-css-injected-by-js",r)}})();
var vt = Object.defineProperty;
var yt = (t, e, n) => e in t ? vt(t, e, { enumerable: !0, configurable: !0, writable: !0, value: n }) : t[e] = n;
var Pe = (t, e, n) => yt(t, typeof e != "symbol" ? e + "" : e, n);
var oe, y, Ve, j, Le, Ye, Je, _e, J, O, Xe, fe, de, pe, ee = {}, te = [], bt = /acit|ex(?:s|g|n|p|$)|rph|grid|ows|mnc|ntw|ine[ch]|zoo|^ord|itera/i, ae = Array.isArray;
function P(t, e) {
  for (var n in e) t[n] = e[n];
  return t;
}
function me(t) {
  t && t.parentNode && t.parentNode.removeChild(t);
}
function kt(t, e, n) {
  var o, i, a, _ = {};
  for (a in e) a == "key" ? o = e[a] : a == "ref" ? i = e[a] : _[a] = e[a];
  if (arguments.length > 2 && (_.children = arguments.length > 3 ? oe.call(arguments, 2) : n), typeof t == "function" && t.defaultProps != null) for (a in t.defaultProps) _[a] === void 0 && (_[a] = t.defaultProps[a]);
  return X(t, _, o, i, null);
}
function X(t, e, n, o, i) {
  var a = { type: t, props: e, key: n, ref: o, __k: null, __: null, __b: 0, __e: null, __c: null, constructor: void 0, __v: i ?? ++Ve, __i: -1, __u: 0 };
  return i == null && y.vnode != null && y.vnode(a), a;
}
function ie(t) {
  return t.children;
}
function Q(t, e) {
  this.props = t, this.context = e;
}
function W(t, e) {
  if (e == null) return t.__ ? W(t.__, t.__i + 1) : null;
  for (var n; e < t.__k.length; e++) if ((n = t.__k[e]) != null && n.__e != null) return n.__e;
  return typeof t.type == "function" ? W(t) : null;
}
function wt(t) {
  if (t.__P && t.__d) {
    var e = t.__v, n = e.__e, o = [], i = [], a = P({}, e);
    a.__v = e.__v + 1, y.vnode && y.vnode(a), ge(t.__P, a, e, t.__n, t.__P.namespaceURI, 32 & e.__u ? [n] : null, o, n ?? W(e), !!(32 & e.__u), i), a.__v = e.__v, a.__.__k[a.__i] = a, nt(o, a, i), e.__e = e.__ = null, a.__e != n && Qe(a);
  }
}
function Qe(t) {
  if ((t = t.__) != null && t.__c != null) return t.__e = t.__c.base = null, t.__k.some(function(e) {
    if (e != null && e.__e != null) return t.__e = t.__c.base = e.__e;
  }), Qe(t);
}
function Me(t) {
  (!t.__d && (t.__d = !0) && j.push(t) && !ne.__r++ || Le != y.debounceRendering) && ((Le = y.debounceRendering) || Ye)(ne);
}
function ne() {
  try {
    for (var t, e = 1; j.length; ) j.length > e && j.sort(Je), t = j.shift(), e = j.length, wt(t);
  } finally {
    j.length = ne.__r = 0;
  }
}
function Ze(t, e, n, o, i, a, _, l, p, s, h) {
  var b, r, u, m, S, f, k = o && o.__k || te, g = e.length;
  for (p = xt(n, e, k, p, g), b = 0; b < g; b++) (u = n.__k[b]) != null && (r = u.__i != -1 && k[u.__i] || ee, u.__i = b, f = ge(t, u, r, i, a, _, l, p, s, h), m = u.__e, u.ref && r.ref != u.ref && (r.ref && ve(r.ref, null, u), h.push(u.ref, u.__c || m, u)), S == null && m != null && (S = m), 4 & u.__u ? (p = et(u, p, t), r.__e && (r.__e = null)) : typeof u.type == "function" && f !== void 0 ? p = f : m && (p = m.nextSibling), u.__u &= -7);
  return n.__e = S, p;
}
function xt(t, e, n, o, i) {
  var a, _, l, p, s, h = n.length, b = h, r = 0;
  for (t.__k = new Array(i), a = 0; a < i; a++) (_ = e[a]) != null && typeof _ != "boolean" && typeof _ != "function" ? (typeof _ == "string" || typeof _ == "number" || typeof _ == "bigint" || _.constructor == String ? _ = t.__k[a] = X(null, _, null, null, null) : ae(_) ? _ = t.__k[a] = X(ie, { children: _ }, null, null, null) : _.constructor === void 0 && _.__b > 0 ? _ = t.__k[a] = X(_.type, _.props, _.key, _.ref ? _.ref : null, _.__v) : t.__k[a] = _, p = a + r, _.__ = t, _.__b = t.__b + 1, l = null, (s = _.__i = Nt(_, n, p, b)) != -1 && (b--, (l = n[s]) && (l.__u |= 2)), l == null || l.__v == null ? (s == -1 && (i > h ? r-- : i < h && r++), typeof _.type != "function" && (_.__u |= 4)) : s != p && (s == p - 1 ? r-- : s == p + 1 ? r++ : (s > p ? r-- : r++, _.__u |= 4))) : t.__k[a] = null;
  if (b) for (a = 0; a < h; a++) (l = n[a]) != null && !(2 & l.__u) && (l.__e == o && (o = W(l)), ot(l, l));
  return o;
}
function et(t, e, n) {
  var o, i;
  if (typeof t.type == "function") {
    for (o = t.__k, i = 0; o && i < o.length; i++) o[i] && (o[i].__ = t, e = et(o[i], e, n));
    return e;
  }
  t.__e != e && (e && t.type && !e.parentNode && (e = W(t)), e = n.insertBefore(t.__e, e || null));
  do
    e = e && e.nextSibling;
  while (e != null && e.nodeType == 8);
  return e;
}
function Nt(t, e, n, o) {
  var i, a, _, l = t.key, p = t.type, s = e[n], h = s != null && (2 & s.__u) == 0;
  if (s === null && l == null || h && l == s.key && p == s.type) return n;
  if (o > (h ? 1 : 0)) {
    for (i = n - 1, a = n + 1; i >= 0 || a < e.length; ) if ((s = e[_ = i >= 0 ? i-- : a++]) != null && !(2 & s.__u) && l == s.key && p == s.type) return _;
  }
  return -1;
}
function je(t, e, n) {
  e[0] == "-" ? t.setProperty(e, n ?? "") : t[e] = n == null ? "" : typeof n != "number" || bt.test(e) ? n : n + "px";
}
function V(t, e, n, o, i) {
  var a, _;
  e: if (e == "style") if (typeof n == "string") t.style.cssText = n;
  else {
    if (typeof o == "string" && (t.style.cssText = o = ""), o) for (e in o) n && e in n || je(t.style, e, "");
    if (n) for (e in n) o && n[e] == o[e] || je(t.style, e, n[e]);
  }
  else if (e[0] == "o" && e[1] == "n") a = e != (e = e.replace(Xe, "$1")), _ = e.toLowerCase(), e = _ in t || e == "onFocusOut" || e == "onFocusIn" ? _.slice(2) : e.slice(2), t.l || (t.l = {}), t.l[e + a] = n, n ? o ? n[O] = o[O] : (n[O] = fe, t.addEventListener(e, a ? pe : de, a)) : t.removeEventListener(e, a ? pe : de, a);
  else {
    if (i == "http://www.w3.org/2000/svg") e = e.replace(/xlink(H|:h)/, "h").replace(/sName$/, "s");
    else if (e != "width" && e != "height" && e != "href" && e != "list" && e != "form" && e != "tabIndex" && e != "download" && e != "rowSpan" && e != "colSpan" && e != "role" && e != "popover" && e in t) try {
      t[e] = n ?? "";
      break e;
    } catch {
    }
    typeof n == "function" || (n == null || n === !1 && e[4] != "-" ? t.removeAttribute(e) : t.setAttribute(e, e == "popover" && n == 1 ? "" : n));
  }
}
function De(t) {
  return function(e) {
    if (this.l) {
      var n = this.l[e.type + t];
      if (e[J] == null) e[J] = fe++;
      else if (e[J] < n[O]) return;
      return n(y.event ? y.event(e) : e);
    }
  };
}
function ge(t, e, n, o, i, a, _, l, p, s) {
  var h, b, r, u, m, S, f, k, g, I, D, T, L, E, C, A, $ = e.type;
  if (e.constructor !== void 0) return null;
  128 & n.__u && (p = !!(32 & n.__u), a = [l = e.__e = n.__e]), (h = y.__b) && h(e);
  e: if (typeof $ == "function") {
    b = _.length;
    try {
      if (g = e.props, I = $.prototype && $.prototype.render, D = (h = $.contextType) && o[h.__c], T = h ? D ? D.props.value : h.__ : o, n.__c ? k = (r = e.__c = n.__c).__ = r.__E : (I ? e.__c = r = new $(g, T) : (e.__c = r = new Q(g, T), r.constructor = $, r.render = Ct), D && D.sub(r), r.state || (r.state = {}), r.__n = o, u = r.__d = !0, r.__h = [], r._sb = []), I && r.__s == null && (r.__s = r.state), I && $.getDerivedStateFromProps != null && (r.__s == r.state && (r.__s = P({}, r.__s)), P(r.__s, $.getDerivedStateFromProps(g, r.__s))), m = r.props, S = r.state, r.__v = e, u) I && $.getDerivedStateFromProps == null && r.componentWillMount != null && r.componentWillMount(), I && r.componentDidMount != null && r.__h.push(r.componentDidMount);
      else {
        if (I && $.getDerivedStateFromProps == null && g !== m && r.componentWillReceiveProps != null && r.componentWillReceiveProps(g, T), e.__v == n.__v || !r.__e && r.shouldComponentUpdate != null && r.shouldComponentUpdate(g, r.__s, T) === !1) {
          e.__v != n.__v && (r.props = g, r.state = r.__s, r.__d = !1), e.__e = n.__e, e.__k = n.__k, e.__k.some(function(H) {
            H && (H.__ = e);
          }), te.push.apply(r.__h, r._sb), r._sb = [], r.__h.length && _.push(r), l = W(n);
          break e;
        }
        r.componentWillUpdate != null && r.componentWillUpdate(g, r.__s, T), I && r.componentDidUpdate != null && r.__h.push(function() {
          r.componentDidUpdate(m, S, f);
        });
      }
      if (r.context = T, r.props = g, r.__P = t, r.__e = !1, L = y.__r, E = 0, I) r.state = r.__s, r.__d = !1, L && L(e), h = r.render(r.props, r.state, r.context), te.push.apply(r.__h, r._sb), r._sb = [];
      else do
        r.__d = !1, L && L(e), h = r.render(r.props, r.state, r.context), r.state = r.__s;
      while (r.__d && ++E < 25);
      r.state = r.__s, r.getChildContext != null && (o = P(P({}, o), r.getChildContext())), I && !u && r.getSnapshotBeforeUpdate != null && (f = r.getSnapshotBeforeUpdate(m, S)), C = h != null && h.type === ie && h.key == null ? rt(h.props.children) : h, l = Ze(t, ae(C) ? C : [C], e, n, o, i, a, _, l, p, s), r.base = e.__e, e.__u &= -161, r.__h.length && _.push(r), k && (r.__E = r.__ = null);
    } catch (H) {
      if (_.length = b, e.__v = null, p || a != null) {
        if (H.then) {
          for (e.__u |= p ? 160 : 128; l && l.nodeType == 8 && l.nextSibling; ) l = l.nextSibling;
          a != null && (a[a.indexOf(l)] = null), e.__e = l;
        } else if (a != null) for (A = a.length; A--; ) me(a[A]);
      } else e.__e = n.__e;
      e.__k == null && (e.__k = n.__k || []), H.then || tt(e), y.__e(H, e, n);
    }
  } else a == null && e.__v == n.__v ? (e.__k = n.__k, e.__e = n.__e) : l = e.__e = St(n.__e, e, n, o, i, a, _, p, s);
  return (h = y.diffed) && h(e), 128 & e.__u ? void 0 : l;
}
function tt(t) {
  t && (t.__c && (t.__c.__e = !0), t.__k && t.__k.some(tt));
}
function nt(t, e, n) {
  for (var o = 0; o < n.length; o++) ve(n[o], n[++o], n[++o]);
  y.__c && y.__c(e, t), t.some(function(i) {
    try {
      t = i.__h, i.__h = [], t.some(function(a) {
        a.call(i);
      });
    } catch (a) {
      y.__e(a, i.__v);
    }
  });
}
function rt(t) {
  return typeof t != "object" || t == null || t.__b > 0 ? t : ae(t) ? t.map(rt) : t.constructor !== void 0 ? null : P({}, t);
}
function St(t, e, n, o, i, a, _, l, p) {
  var s, h, b, r, u, m, S, f = n.props || ee, k = e.props, g = e.type;
  if (g == "svg" ? i = "http://www.w3.org/2000/svg" : g == "math" ? i = "http://www.w3.org/1998/Math/MathML" : i || (i = "http://www.w3.org/1999/xhtml"), a != null) {
    for (s = 0; s < a.length; s++) if ((u = a[s]) && "setAttribute" in u == !!g && (g ? u.localName == g : u.nodeType == 3)) {
      t = u, a[s] = null;
      break;
    }
  }
  if (t == null) {
    if (g == null) return document.createTextNode(k);
    t = document.createElementNS(i, g, k.is && k), l && (y.__m && y.__m(e, a), l = !1), a = null;
  }
  if (g == null) f === k || l && t.data == k || (t.data = k);
  else {
    if (a = g == "textarea" && k.defaultValue != null ? null : a && oe.call(t.childNodes), !l && a != null) for (f = {}, s = 0; s < t.attributes.length; s++) f[(u = t.attributes[s]).name] = u.value;
    for (s in f) u = f[s], s == "dangerouslySetInnerHTML" ? b = u : s == "children" || s in k || s == "value" && "defaultValue" in k || s == "checked" && "defaultChecked" in k || V(t, s, null, u, i);
    for (s in k) u = k[s], s == "children" ? r = u : s == "dangerouslySetInnerHTML" ? h = u : s == "value" ? m = u : s == "checked" ? S = u : l && typeof u != "function" || f[s] === u || V(t, s, u, f[s], i);
    if (h) l || b && (h.__html == b.__html || h.__html == t.innerHTML) || (t.innerHTML = h.__html), e.__k = [];
    else if (b && (t.innerHTML = ""), Ze(e.type == "template" ? t.content : t, ae(r) ? r : [r], e, n, o, g == "foreignObject" ? "http://www.w3.org/1999/xhtml" : i, a, _, a ? a[0] : n.__k && W(n, 0), l, p), a != null) for (s = a.length; s--; ) me(a[s]);
    l && g != "textarea" || (s = "value", g == "progress" && m == null ? t.removeAttribute("value") : m != null && (m !== t[s] || g == "progress" && !m || g == "option" && m != f[s]) && V(t, s, m, f[s], i), s = "checked", S != null && S != t[s] && V(t, s, S, f[s], i));
  }
  return t;
}
function ve(t, e, n) {
  try {
    if (typeof t == "function") {
      var o = typeof t.__u == "function";
      o && t.__u(), o && e == null || (t.__u = t(e));
    } else t.current = e;
  } catch (i) {
    y.__e(i, n);
  }
}
function ot(t, e, n) {
  var o, i;
  if (y.unmount && y.unmount(t), (o = t.ref) && (o.current && o.current != t.__e || ve(o, null, e)), (o = t.__c) != null) {
    if (o.componentWillUnmount) try {
      o.componentWillUnmount();
    } catch (a) {
      y.__e(a, e);
    }
    o.base = o.__P = o.__n = null;
  }
  if (o = t.__k) for (i = 0; i < o.length; i++) o[i] && ot(o[i], e, n || typeof t.type != "function");
  n || me(t.__e), t.__c = t.__ = t.__e = void 0;
}
function Ct(t, e, n) {
  return this.constructor(t, n);
}
function We(t, e, n) {
  var o, i, a, _;
  e == document && (e = document.documentElement), y.__ && y.__(t, e), i = (o = !1) ? null : e.__k, a = [], _ = [], ge(e, t = e.__k = kt(ie, null, [t]), i || ee, ee, e.namespaceURI, i ? null : e.firstChild ? oe.call(e.childNodes) : null, a, i ? i.__e : e.firstChild, o, _), nt(a, t, _), t.props.children = null;
}
oe = te.slice, y = { __e: function(t, e, n, o) {
  for (var i, a, _; e = e.__; ) if ((i = e.__c) && !i.__) try {
    if ((a = i.constructor) && a.getDerivedStateFromError != null && (i.setState(a.getDerivedStateFromError(t)), _ = i.__d), i.componentDidCatch != null && (i.componentDidCatch(t, o || {}), _ = i.__d), _) return i.__E = i;
  } catch (l) {
    t = l;
  }
  throw t;
} }, Ve = 0, Q.prototype.setState = function(t, e) {
  var n;
  n = this.__s != null && this.__s != this.state ? this.__s : this.__s = P({}, this.state), typeof t == "function" && (t = t(P({}, n), this.props)), t && P(n, t), t != null && this.__v && (e && this._sb.push(e), Me(this));
}, Q.prototype.forceUpdate = function(t) {
  this.__v && (this.__e = !0, t && this.__h.push(t), Me(this));
}, Q.prototype.render = ie, j = [], Ye = typeof Promise == "function" ? Promise.prototype.then.bind(Promise.resolve()) : setTimeout, Je = function(t, e) {
  return t.__v.__b - e.__v.__b;
}, ne.__r = 0, _e = Math.random().toString(8), J = "__d" + _e, O = "__a" + _e, Xe = /(PointerCapture)$|Capture$/i, fe = 0, de = De(!1), pe = De(!0);
var $t = 0;
function c(t, e, n, o, i, a) {
  e || (e = {});
  var _, l, p = e;
  if ("ref" in p) for (l in p = {}, e) l == "ref" ? _ = e[l] : p[l] = e[l];
  var s = { type: t, props: p, key: n, ref: _, __k: null, __: null, __b: 0, __e: null, __c: null, constructor: void 0, __v: --$t, __i: -1, __u: 0, __source: i, __self: a };
  if (typeof t == "function" && (_ = t.defaultProps)) for (l in _) p[l] === void 0 && (p[l] = _[l]);
  return y.vnode && y.vnode(s), s;
}
var G, w, le, Ae, re = 0, at = [], x = y, Fe = x.__b, Be = x.__r, qe = x.diffed, Re = x.__c, Oe = x.unmount, Ge = x.__;
function ye(t, e) {
  x.__h && x.__h(w, t, re || e), re = 0;
  var n = w.__H || (w.__H = { __: [], __h: [] });
  return t >= n.__.length && n.__.push({}), n.__[t];
}
function U(t) {
  return re = 1, It(st, t);
}
function It(t, e, n) {
  var o = ye(G++, 2);
  if (o.t = t, !o.__c && (o.__ = [st(void 0, e), function(l) {
    var p = o.__N ? o.__N[0] : o.__[0], s = o.t(p, l);
    p !== s && (o.__N = [s, o.__[1]], o.__c.setState({}));
  }], o.__c = w, !w.__f)) {
    var i = function(l, p, s) {
      if (!o.__c.__H) return !0;
      var h = !1, b = o.__c.props !== l;
      if (o.__c.__H.__.some(function(u) {
        if (u.__N) {
          h = !0;
          var m = u.__[0];
          u.__ = u.__N, u.__N = void 0, m !== u.__[0] && (b = !0);
        }
      }), a) {
        var r = a.call(this, l, p, s);
        return h ? r || b : r;
      }
      return !h || b;
    };
    w.__f = !0;
    var a = w.shouldComponentUpdate, _ = w.componentWillUpdate;
    w.componentWillUpdate = function(l, p, s) {
      if (this.__e) {
        var h = a;
        a = void 0, i(l, p, s), a = h;
      }
      _ && _.call(this, l, p, s);
    }, w.shouldComponentUpdate = i;
  }
  return o.__N || o.__;
}
function Y(t, e) {
  var n = ye(G++, 3);
  !x.__s && it(n.__H, e) && (n.__ = t, n.u = e, w.__H.__h.push(n));
}
function ce(t) {
  return re = 5, Tt(function() {
    return { current: t };
  }, []);
}
function Tt(t, e) {
  var n = ye(G++, 7);
  return it(n.__H, e) && (n.__ = t(), n.__H = e, n.__h = t), n.__;
}
function Ut() {
  for (var t; t = at.shift(); ) {
    var e = t.__H;
    if (t.__P && e) try {
      e.__h.some(Z), e.__h.some(he), e.__h = [];
    } catch (n) {
      e.__h = [], x.__e(n, t.__v);
    }
  }
}
x.__b = function(t) {
  w = null, Fe && Fe(t);
}, x.__ = function(t, e) {
  t && e.__k && e.__k.__m && (t.__m = e.__k.__m), Ge && Ge(t, e);
}, x.__r = function(t) {
  Be && Be(t), G = 0;
  var e = (w = t.__c).__H;
  e && (le === w ? (e.__h = [], w.__h = [], e.__.some(function(n) {
    n.__N && (n.__ = n.__N), n.u = n.__N = void 0;
  })) : (e.__h.some(Z), e.__h.some(he), e.__h = [], G = 0)), le = w;
}, x.diffed = function(t) {
  qe && qe(t);
  var e = t.__c;
  e && e.__H && (e.__H.__h.length && (at.push(e) !== 1 && Ae === x.requestAnimationFrame || ((Ae = x.requestAnimationFrame) || Et)(Ut)), e.__H.__.some(function(n) {
    n.u && (n.__H = n.u, n.u = void 0);
  })), le = w = null;
}, x.__c = function(t, e) {
  e.some(function(n) {
    try {
      n.__h.some(Z), n.__h = n.__h.filter(function(o) {
        return !o.__ || he(o);
      });
    } catch (o) {
      e.some(function(i) {
        i.__h && (i.__h = []);
      }), e = [], x.__e(o, n.__v);
    }
  }), Re && Re(t, e);
}, x.unmount = function(t) {
  Oe && Oe(t);
  var e, n = t.__c;
  n && n.__H && (n.__H.__.some(function(o) {
    try {
      Z(o);
    } catch (i) {
      e = i;
    }
  }), n.__H = void 0, e && x.__e(e, n.__v));
};
var ze = typeof requestAnimationFrame == "function";
function Et(t) {
  var e, n = function() {
    clearTimeout(o), ze && cancelAnimationFrame(e), setTimeout(t);
  }, o = setTimeout(n, 35);
  ze && (e = requestAnimationFrame(n));
}
function Z(t) {
  var e = w, n = t.__c;
  typeof n == "function" && (t.__c = void 0, n()), w = e;
}
function he(t) {
  var e = w;
  t.__c = t.__(), w = e;
}
function it(t, e) {
  return !t || t.length !== e.length || e.some(function(n, o) {
    return n !== t[o];
  });
}
function st(t, e) {
  return typeof e == "function" ? e(t) : e;
}
const R = {
  ru: {
    title: "Поддержка",
    subtitle: "Оператор ответит вам здесь",
    subtitleGuest: "Оставьте вопрос или пожелание",
    empty: "Задайте любой вопрос или отправьте снимок экрана — мы на связи!",
    email: "Ваш email (для ответа)",
    name: "Как к вам обращаться",
    placeholder: "Напишите сообщение... (Enter — отправить)",
    send: "Отправить",
    attach: "Прикрепить скриншот",
    removeImage: "Удалить картинку",
    closedNote: "Тема закрыта оператором",
    reopenedNote: "Диалог возобновлен",
    closedHint: "Вопрос решен. Напишите снова, если потребуется помощь.",
    foot: "Ответы оператора появятся здесь и придут на email",
    errors: {
      email: "Укажите корректный email для связи",
      image: "Недопустимый формат или размер файла (макс. 4 МБ)",
      long: "Слишком длинное сообщение (макс. 4000 знаков)",
      failed: "Не удалось отправить сообщение. Попробуйте еще раз"
    }
  },
  en: {
    title: "Support Chat",
    subtitle: "Our team will answer here",
    subtitleGuest: "Ask a question or leave feedback",
    empty: "Ask anything or attach a screenshot — we are here to help!",
    email: "Your email address",
    name: "Your name",
    placeholder: "Type a message... (Enter to send)",
    send: "Send",
    attach: "Attach screenshot",
    removeImage: "Remove image",
    closedNote: "Ticket closed by operator",
    reopenedNote: "Ticket reopened",
    closedHint: "This conversation is closed. Send a message to reopen.",
    foot: "Replies will appear here and in your inbox",
    errors: {
      email: "Please provide a valid email address",
      image: "Unsupported image format or size (max 4 MB)",
      long: "Message is too long (max 4,000 characters)",
      failed: "Could not send message. Please try again"
    }
  },
  es: {
    title: "Soporte",
    subtitle: "Nuestro equipo responderá aquí",
    subtitleGuest: "Haz una pregunta o deja comentarios",
    empty: "Escribe cualquier duda o adjunta una captura de pantalla",
    email: "Tu correo electrónico",
    name: "Tu nombre",
    placeholder: "Escribe un mensaje... (Enter para enviar)",
    send: "Enviar",
    attach: "Adjuntar imagen",
    removeImage: "Eliminar imagen",
    closedNote: "Conversación cerrada por el operador",
    reopenedNote: "Conversación reabierta",
    closedHint: "Esta conversación ha sido cerrada. Escribe para reabrirla.",
    foot: "Las respuestas llegarán aquí y a tu correo",
    errors: {
      email: "Introduce un correo electrónico válido",
      image: "Formato o tamaño de imagen no admitido (máx. 4 MB)",
      long: "Mensaje demasiado largo (máx. 4000 caracteres)",
      failed: "No se pudo enviar el mensaje. Inténtalo de nuevo"
    }
  }
};
function Ht(t) {
  if (t && t !== "auto" && R[t])
    return R[t];
  const e = (typeof navigator < "u" && navigator.language || "").toLowerCase();
  return e.startsWith("ru") || e.startsWith("be") || e.startsWith("uk") || e.startsWith("kk") ? R.ru : e.startsWith("es") ? R.es : R.en;
}
const _t = "support_chat_token", lt = "support_chat_email", ct = "support_chat_name";
function Pt() {
  try {
    return localStorage.getItem(_t);
  } catch {
    return null;
  }
}
function ue(t) {
  try {
    localStorage.setItem(_t, t);
  } catch {
  }
}
function Lt() {
  try {
    return {
      email: localStorage.getItem(lt) || "",
      name: localStorage.getItem(ct) || ""
    };
  } catch {
    return { email: "", name: "" };
  }
}
function Mt(t, e) {
  try {
    t && localStorage.setItem(lt, t), e && localStorage.setItem(ct, e);
  } catch {
  }
}
class jt {
  constructor(e) {
    Pe(this, "baseUrl");
    this.baseUrl = e.replace(/\/+$/, "");
  }
  headers(e, n) {
    const o = {};
    return e && (o["X-Support-Token"] = e), n && (o["X-User-Id"] = n), o;
  }
  async fetchView(e, n) {
    const o = await fetch(`${this.baseUrl}/api/v1/support`, {
      method: "GET",
      headers: this.headers(e, n)
    });
    if (!o.ok)
      throw new Error(`fetch view failed: ${o.status}`);
    return o.json();
  }
  async sendMessage(e, n, o) {
    const i = await fetch(`${this.baseUrl}/api/v1/support/messages`, {
      method: "POST",
      body: e,
      headers: this.headers(n, o)
    });
    if (!i.ok) {
      const a = await i.json().catch(() => ({})), _ = new Error(a.error || `send message failed: ${i.status}`);
      throw _.code = a.error || "", _;
    }
    return i.json();
  }
  async markSeen(e, n) {
    await fetch(`${this.baseUrl}/api/v1/support/seen`, {
      method: "POST",
      headers: this.headers(e, n)
    }).catch(() => {
    });
  }
  async checkUnread(e, n) {
    const o = await fetch(`${this.baseUrl}/api/v1/support/unread`, {
      method: "GET",
      headers: this.headers(e, n)
    });
    if (!o.ok) return 0;
    const i = await o.json();
    return typeof i.unread == "number" ? i.unread : 0;
  }
  attachmentUrl(e, n) {
    const o = `${this.baseUrl}/api/v1/support/attachments/${e}`;
    return n ? `${o}?t=${encodeURIComponent(n)}` : o;
  }
}
function Dt(t) {
  return /* @__PURE__ */ c("svg", { width: "24", height: "24", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round", strokeLinejoin: "round", ...t, children: /* @__PURE__ */ c("path", { d: "M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" }) });
}
function Ke(t) {
  return /* @__PURE__ */ c("svg", { width: "18", height: "18", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round", strokeLinejoin: "round", ...t, children: [
    /* @__PURE__ */ c("line", { x1: "18", y1: "6", x2: "6", y2: "18" }),
    /* @__PURE__ */ c("line", { x1: "6", y1: "6", x2: "18", y2: "18" })
  ] });
}
function Wt(t) {
  return /* @__PURE__ */ c("svg", { width: "16", height: "16", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round", strokeLinejoin: "round", ...t, children: [
    /* @__PURE__ */ c("line", { x1: "22", y1: "2", x2: "11", y2: "13" }),
    /* @__PURE__ */ c("polygon", { points: "22 2 15 22 11 13 2 9 22 2" })
  ] });
}
function At(t) {
  return /* @__PURE__ */ c("svg", { width: "18", height: "18", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round", strokeLinejoin: "round", ...t, children: /* @__PURE__ */ c("path", { d: "M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48" }) });
}
function Ft(t) {
  return /* @__PURE__ */ c("svg", { width: "20", height: "20", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round", strokeLinejoin: "round", ...t, children: [
    /* @__PURE__ */ c("circle", { cx: "12", cy: "12", r: "10" }),
    /* @__PURE__ */ c("circle", { cx: "12", cy: "12", r: "4" }),
    /* @__PURE__ */ c("line", { x1: "4.93", y1: "4.93", x2: "9.17", y2: "9.17" }),
    /* @__PURE__ */ c("line", { x1: "14.83", y1: "14.83", x2: "19.07", y2: "19.07" }),
    /* @__PURE__ */ c("line", { x1: "14.83", y1: "9.17", x2: "19.07", y2: "4.93" }),
    /* @__PURE__ */ c("line", { x1: "14.83", y1: "9.17", x2: "18.36", y2: "5.64" }),
    /* @__PURE__ */ c("line", { x1: "4.93", y1: "19.07", x2: "9.17", y2: "14.83" })
  ] });
}
function Bt(t) {
  var Ee;
  const {
    apiUrl: e,
    title: n,
    subtitle: o,
    primaryColor: i = "#2563eb",
    position: a = "bottom-right",
    lang: _ = "auto",
    userEmail: l = "",
    userName: p = "",
    userId: s,
    metadata: h,
    autoOpen: b = !1
  } = t, r = Ht(_), u = ce(new jt(e)).current, [m, S] = U(b), [f, k] = U(null), [g, I] = U(0), [D, T] = U(!1), [L, E] = U(null), [C, A] = U(Pt()), [$, H] = U(""), [B, z] = U(null), be = Lt(), [F, ut] = U(l || be.email), [q, dt] = U(p || be.name), [ke, we] = U(null), K = ce(null), xe = ce(null);
  Y(() => {
    if (typeof window > "u") return;
    const d = new URLSearchParams(window.location.search), v = d.get("support") || d.get("feedback");
    if (v) {
      v !== "1" && (ue(v), A(v)), S(!0), d.delete("support"), d.delete("feedback");
      const N = d.toString();
      window.history.replaceState(null, "", window.location.pathname + (N ? `?${N}` : "") + window.location.hash);
    }
  }, []), Y(() => {
    C && ue(C);
  }, [C]);
  const se = async () => {
    try {
      const d = await u.fetchView(C, s);
      k(d), m && d.thread && (await u.markSeen(C, s), I(0));
    } catch {
    }
  }, Ne = async () => {
    try {
      const d = await u.checkUnread(C, s);
      I(m ? 0 : d);
    } catch {
    }
  };
  Y(() => {
    m ? se() : Ne();
    const v = setInterval(() => {
      m ? se() : Ne();
    }, m ? 6e3 : 45e3);
    return () => clearInterval(v);
  }, [m, C, s]), Y(() => {
    m && setTimeout(() => {
      K.current && (K.current.scrollTop = K.current.scrollHeight);
    }, 50);
  }, [m, f == null ? void 0 : f.messages.length]);
  const pt = () => {
    S(!0), E(null), se();
  }, ht = () => {
    S(!1);
  }, Se = !s, ft = !(f != null && f.thread), Ce = Se && ft, $e = !D && ($.trim() !== "" || B !== null) && (!Ce || /\S+@\S+\.\S+/.test(F)), Ie = async (d) => {
    if (d.preventDefault(), !$e) return;
    T(!0), E(null);
    const v = new FormData();
    v.append("text", $.trim()), B && v.append("image", B), v.append("page", typeof window < "u" ? window.location.pathname : ""), v.append("lang", _ === "auto" ? typeof navigator < "u" ? navigator.language : "ru" : _), F && v.append("email", F.trim()), q && v.append("name", q.trim()), h && v.append("metadata", JSON.stringify(h));
    try {
      const N = await u.sendMessage(v, C, s);
      N.token && (ue(N.token), A(N.token)), (F || q) && Mt(F, q), k(N), H(""), z(null), T(!1);
    } catch (N) {
      T(!1);
      const M = (N == null ? void 0 : N.code) || "";
      E(M === "email_required" ? r.errors.email : M === "invalid_image" ? r.errors.image : M === "message_too_long" ? r.errors.long : r.errors.failed);
    }
  }, mt = (d) => {
    var N;
    const v = (N = d.clipboardData) == null ? void 0 : N.items;
    if (v) {
      for (let M = 0; M < v.length; M++)
        if (v[M].type.startsWith("image/")) {
          const He = v[M].getAsFile();
          if (He) {
            z(He), d.preventDefault();
            break;
          }
        }
    }
  }, Te = (d) => {
    try {
      return new Intl.DateTimeFormat(void 0, {
        hour: "2-digit",
        minute: "2-digit"
      }).format(new Date(d));
    } catch {
      return "";
    }
  }, gt = (d) => {
    if (d.author === "system")
      return /* @__PURE__ */ c("li", { className: "sc-note", children: [
        d.text === "closed" ? r.closedNote : r.reopenedNote,
        " · ",
        Te(d.created_at)
      ] }, d.id);
    const v = d.author === "user", N = d.has_attachment ? u.attachmentUrl(d.id, C) : null;
    return /* @__PURE__ */ c("li", { className: `sc-bubble ${v ? "sc-mine" : "sc-theirs"}`, children: [
      !v && /* @__PURE__ */ c("span", { className: "sc-who", children: d.operator_name || "Поддержка" }),
      N && /* @__PURE__ */ c(
        "img",
        {
          src: N,
          alt: "",
          className: "sc-image",
          onClick: () => we(N)
        }
      ),
      d.text && /* @__PURE__ */ c("p", { children: d.text }),
      /* @__PURE__ */ c("span", { className: "sc-time", children: Te(d.created_at) })
    ] }, d.id);
  }, Ue = a === "bottom-left" ? "sc-pos-left" : "";
  return /* @__PURE__ */ c("div", { style: { "--sc-primary": i }, children: [
    !m && /* @__PURE__ */ c(
      "button",
      {
        className: `sc-launcher ${Ue}`,
        onClick: pt,
        "aria-label": r.title,
        title: r.title,
        children: [
          /* @__PURE__ */ c(Dt, {}),
          g > 0 && /* @__PURE__ */ c("span", { className: "sc-badge", children: g })
        ]
      }
    ),
    m && /* @__PURE__ */ c("section", { className: `sc-panel ${Ue}`, role: "dialog", "aria-label": r.title, children: [
      /* @__PURE__ */ c("header", { className: "sc-header", children: [
        /* @__PURE__ */ c("div", { className: "sc-header-icon", children: /* @__PURE__ */ c(Ft, {}) }),
        /* @__PURE__ */ c("div", { className: "sc-header-text", children: [
          /* @__PURE__ */ c("h3", { className: "sc-title", children: n || r.title }),
          /* @__PURE__ */ c("p", { className: "sc-subtitle", children: o || (Se ? r.subtitleGuest : r.subtitle) })
        ] }),
        /* @__PURE__ */ c("button", { className: "sc-close-btn", onClick: ht, title: "Закрыть", children: /* @__PURE__ */ c(Ke, {}) })
      ] }),
      /* @__PURE__ */ c("ul", { className: "sc-messages", ref: K, children: [
        (!(f != null && f.messages) || f.messages.length === 0) && /* @__PURE__ */ c("li", { className: "sc-empty", children: r.empty }),
        f == null ? void 0 : f.messages.map(gt),
        ((Ee = f == null ? void 0 : f.thread) == null ? void 0 : Ee.status) === "closed" && /* @__PURE__ */ c("li", { className: "sc-note", children: r.closedHint })
      ] }),
      /* @__PURE__ */ c("form", { className: "sc-form", onSubmit: Ie, children: [
        Ce && /* @__PURE__ */ c("div", { className: "sc-contacts", children: [
          /* @__PURE__ */ c(
            "input",
            {
              className: "sc-input",
              type: "email",
              required: !0,
              placeholder: r.email,
              value: F,
              onInput: (d) => ut(d.target.value)
            }
          ),
          /* @__PURE__ */ c(
            "input",
            {
              className: "sc-input",
              type: "text",
              placeholder: r.name,
              value: q,
              onInput: (d) => dt(d.target.value)
            }
          ),
          /* @__PURE__ */ c("input", { className: "sc-hp", name: "website", tabIndex: -1, autoComplete: "off" })
        ] }),
        B && /* @__PURE__ */ c("div", { className: "sc-attached-chip", children: [
          /* @__PURE__ */ c("span", { children: [
            "📷 ",
            B.name
          ] }),
          /* @__PURE__ */ c(
            "button",
            {
              type: "button",
              className: "sc-close-btn",
              onClick: () => z(null),
              title: r.removeImage,
              children: /* @__PURE__ */ c(Ke, {})
            }
          )
        ] }),
        L && /* @__PURE__ */ c("p", { className: "sc-error", children: L }),
        /* @__PURE__ */ c("div", { className: "sc-compose", children: [
          /* @__PURE__ */ c(
            "button",
            {
              type: "button",
              className: "sc-action-btn",
              onClick: () => {
                var d;
                return (d = xe.current) == null ? void 0 : d.click();
              },
              title: r.attach,
              children: /* @__PURE__ */ c(At, {})
            }
          ),
          /* @__PURE__ */ c(
            "input",
            {
              ref: xe,
              type: "file",
              accept: "image/png,image/jpeg,image/webp",
              style: { display: "none" },
              onChange: (d) => {
                const v = d.target.files;
                v && v[0] && z(v[0]), d.target.value = "";
              }
            }
          ),
          /* @__PURE__ */ c(
            "textarea",
            {
              className: "sc-textarea",
              rows: 2,
              maxLength: 4e3,
              placeholder: r.placeholder,
              value: $,
              onInput: (d) => H(d.target.value),
              onPaste: mt,
              onKeyDown: (d) => {
                d.key === "Enter" && !d.shiftKey && (d.preventDefault(), Ie(d));
              }
            }
          ),
          /* @__PURE__ */ c(
            "button",
            {
              type: "submit",
              className: "sc-send-btn",
              disabled: !$e,
              title: r.send,
              children: /* @__PURE__ */ c(Wt, {})
            }
          )
        ] }),
        /* @__PURE__ */ c("p", { className: "sc-foot", children: r.foot })
      ] })
    ] }),
    ke && /* @__PURE__ */ c(
      "div",
      {
        style: {
          position: "fixed",
          inset: 0,
          background: "rgba(0,0,0,0.85)",
          zIndex: 999999,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          padding: "24px"
        },
        onClick: () => we(null),
        children: /* @__PURE__ */ c(
          "img",
          {
            src: ke,
            alt: "Preview",
            style: {
              maxWidth: "90vw",
              maxHeight: "90vh",
              borderRadius: "8px",
              objectFit: "contain"
            }
          }
        )
      }
    )
  ] });
}
function Ot(t, e) {
  const n = e || (() => {
    let o = document.getElementById("support-chat-root");
    return o || (o = document.createElement("div"), o.id = "support-chat-root", document.body.appendChild(o)), o;
  })();
  return We(/* @__PURE__ */ c(Bt, { ...t }), n), () => {
    We(null, n);
  };
}
export {
  jt as SupportClient,
  Bt as SupportWidget,
  Ot as mountSupportChat,
  Ht as resolveLanguage,
  R as translations
};
