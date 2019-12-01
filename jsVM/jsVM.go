package jsVM

import (
	"github.com/robertkrimen/otto"
)

const jsUtils = `
function json(d) {
	return JSON.stringify(d);
}

function xEncode(str, key) {
	if (str == "") {
		return "";
	}
	var v = s(str, true),
		k = s(key, false);
	if (k.length < 4) {
		k.length = 4;
	}
	var n = v.length - 1,
		z = v[n],
		y = v[0],
		c = 0x86014019 | 0x183639A0,
		m,
		e,
		p,
		q = Math.floor(6 + 52 / (n + 1)),
		d = 0;
	while (0 < q--) {
		d = d + c & (0x8CE0D9BF | 0x731F2640);
		e = d >>> 2 & 3;
		for (p = 0; p < n; p++) {
			y = v[p + 1];
			m = z >>> 5 ^ y << 2;
			m += (y >>> 3 ^ z << 4) ^ (d ^ y);
			m += k[(p & 3) ^ e] ^ z;
			z = v[p] = v[p] + m & (0xEFB8D130 | 0x10472ECF);
		}
		y = v[0];
		m = z >>> 5 ^ y << 2;
		m += (y >>> 3 ^ z << 4) ^ (d ^ y);
		m += k[(p & 3) ^ e] ^ z;
		z = v[n] = v[n] + m & (0xBB390742 | 0x44C6F8BD);
	}
	return l(v, false);
}

function l(a, b) {
	var d = a.length, c = (d - 1) << 2;
	if (b) {
		var m = a[d - 1];
		if ((m < c - 3) || (m > c))
			return null;
		c = m;
	}
	for (var i = 0; i < d; i++) {
		a[i] = String.fromCharCode(a[i] & 0xff, a[i] >>> 8 & 0xff, a[i] >>> 16 & 0xff, a[i] >>> 24 & 0xff);
	}
	if (b) {
		return a.join('').substring(0, c);
	} else {
		return a.join('');
	}
}

function s(a, b) {
	var c = a.length, v = [];
	for (var i = 0; i < c; i += 4) {
		v[i >> 2] = a.charCodeAt(i) | a.charCodeAt(i + 1) << 8 | a.charCodeAt(i + 2) << 16 | a.charCodeAt(i + 3) << 24;
	}
	if (b) {
		v[v.length] = c;
	}
	return v;
}
`

var jsBase64 = `
//jquery base64

var _PADCHAR = "=",
  _ALPHA = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA",
  _VERSION = "1.0";
function _getbyte64(s, i) {
  var idx = _ALPHA.indexOf(s.charAt(i));
  if (idx === -1) {
    throw "Cannot decode base64";
  }
  return idx;
}
function _setAlpha(s) {
  _ALPHA = s;
}
function _decode(s) {
  var pads = 0,
    i,
    b10,
    imax = s.length,
    x = [];
  s = String(s);
  if (imax === 0) {
    return s;
  }
  if (imax % 4 !== 0) {
    throw "Cannot decode base64";
  }
  if (s.charAt(imax - 1) === _PADCHAR) {
    pads = 1;
    if (s.charAt(imax - 2) === _PADCHAR) {
      pads = 2;
    }
    imax -= 4;
  }
  for (i = 0; i < imax; i += 4) {
    b10 =
      (_getbyte64(s, i) << 18) |
      (_getbyte64(s, i + 1) << 12) |
      (_getbyte64(s, i + 2) << 6) |
      _getbyte64(s, i + 3);
    x.push(String.fromCharCode(b10 >> 16, (b10 >> 8) & 255, b10 & 255));
  }
  switch (pads) {
    case 1:
      b10 =
        (_getbyte64(s, i) << 18) |
        (_getbyte64(s, i + 1) << 12) |
        (_getbyte64(s, i + 2) << 6);
      x.push(String.fromCharCode(b10 >> 16, (b10 >> 8) & 255));
      break;
    case 2:
      b10 = (_getbyte64(s, i) << 18) | (_getbyte64(s, i + 1) << 12);
      x.push(String.fromCharCode(b10 >> 16));
      break;
  }
  return x.join("");
}
function _getbyte(s, i) {
  var x = s.charCodeAt(i);
  if (x > 255) {
    throw "INVALID_CHARACTER_ERR: DOM Exception 5";
  }
  return x;
}
function _encode(s) {
  if (arguments.length !== 1) {
    throw "SyntaxError: exactly one argument required";
  }
  s = String(s);
  var i,
    b10,
    x = [],
    imax = s.length - (s.length % 3);
  if (s.length === 0) {
    return s;
  }
  for (i = 0; i < imax; i += 3) {
    b10 =
      (_getbyte(s, i) << 16) | (_getbyte(s, i + 1) << 8) | _getbyte(s, i + 2);
    x.push(_ALPHA.charAt(b10 >> 18));
    x.push(_ALPHA.charAt((b10 >> 12) & 63));
    x.push(_ALPHA.charAt((b10 >> 6) & 63));
    x.push(_ALPHA.charAt(b10 & 63));
  }
  switch (s.length - imax) {
    case 1:
      b10 = _getbyte(s, i) << 16;
      x.push(
        _ALPHA.charAt(b10 >> 18) +
          _ALPHA.charAt((b10 >> 12) & 63) +
          _PADCHAR +
          _PADCHAR
      );
      break;
    case 2:
      b10 = (_getbyte(s, i) << 16) | (_getbyte(s, i + 1) << 8);
      x.push(
        _ALPHA.charAt(b10 >> 18) +
          _ALPHA.charAt((b10 >> 12) & 63) +
          _ALPHA.charAt((b10 >> 6) & 63) +
          _PADCHAR
      );
      break;
  }
  return x.join("");
}
`

var jsMd5 = `
//md5 v2.10.0
function t(n, t) {
  var r = (65535 & n) + (65535 & t);
  return (((n >> 16) + (t >> 16) + (r >> 16)) << 16) | (65535 & r);
}
function r(n, t) {
  return (n << t) | (n >>> (32 - t));
}
function e(n, e, o, u, c, f) {
  return t(r(t(t(e, n), t(u, f)), c), o);
}
function o(n, t, r, o, u, c, f) {
  return e((t & r) | (~t & o), n, t, u, c, f);
}
function u(n, t, r, o, u, c, f) {
  return e((t & o) | (r & ~o), n, t, u, c, f);
}
function c(n, t, r, o, u, c, f) {
  return e(t ^ r ^ o, n, t, u, c, f);
}
function f(n, t, r, o, u, c, f) {
  return e(r ^ (t | ~o), n, t, u, c, f);
}
function i(n, r) {
  (n[r >> 5] |= 128 << r % 32), (n[14 + (((r + 64) >>> 9) << 4)] = r);
  var e,
    i,
    a,
    d,
    h,
    l = 1732584193,
    g = -271733879,
    v = -1732584194,
    m = 271733878;
  for (e = 0; e < n.length; e += 16)
    (i = l),
      (a = g),
      (d = v),
      (h = m),
      (g = f(
        (g = f(
          (g = f(
            (g = f(
              (g = c(
                (g = c(
                  (g = c(
                    (g = c(
                      (g = u(
                        (g = u(
                          (g = u(
                            (g = u(
                              (g = o(
                                (g = o(
                                  (g = o(
                                    (g = o(
                                      g,
                                      (v = o(
                                        v,
                                        (m = o(
                                          m,
                                          (l = o(
                                            l,
                                            g,
                                            v,
                                            m,
                                            n[e],
                                            7,
                                            -680876936
                                          )),
                                          g,
                                          v,
                                          n[e + 1],
                                          12,
                                          -389564586
                                        )),
                                        l,
                                        g,
                                        n[e + 2],
                                        17,
                                        606105819
                                      )),
                                      m,
                                      l,
                                      n[e + 3],
                                      22,
                                      -1044525330
                                    )),
                                    (v = o(
                                      v,
                                      (m = o(
                                        m,
                                        (l = o(
                                          l,
                                          g,
                                          v,
                                          m,
                                          n[e + 4],
                                          7,
                                          -176418897
                                        )),
                                        g,
                                        v,
                                        n[e + 5],
                                        12,
                                        1200080426
                                      )),
                                      l,
                                      g,
                                      n[e + 6],
                                      17,
                                      -1473231341
                                    )),
                                    m,
                                    l,
                                    n[e + 7],
                                    22,
                                    -45705983
                                  )),
                                  (v = o(
                                    v,
                                    (m = o(
                                      m,
                                      (l = o(
                                        l,
                                        g,
                                        v,
                                        m,
                                        n[e + 8],
                                        7,
                                        1770035416
                                      )),
                                      g,
                                      v,
                                      n[e + 9],
                                      12,
                                      -1958414417
                                    )),
                                    l,
                                    g,
                                    n[e + 10],
                                    17,
                                    -42063
                                  )),
                                  m,
                                  l,
                                  n[e + 11],
                                  22,
                                  -1990404162
                                )),
                                (v = o(
                                  v,
                                  (m = o(
                                    m,
                                    (l = o(
                                      l,
                                      g,
                                      v,
                                      m,
                                      n[e + 12],
                                      7,
                                      1804603682
                                    )),
                                    g,
                                    v,
                                    n[e + 13],
                                    12,
                                    -40341101
                                  )),
                                  l,
                                  g,
                                  n[e + 14],
                                  17,
                                  -1502002290
                                )),
                                m,
                                l,
                                n[e + 15],
                                22,
                                1236535329
                              )),
                              (v = u(
                                v,
                                (m = u(
                                  m,
                                  (l = u(l, g, v, m, n[e + 1], 5, -165796510)),
                                  g,
                                  v,
                                  n[e + 6],
                                  9,
                                  -1069501632
                                )),
                                l,
                                g,
                                n[e + 11],
                                14,
                                643717713
                              )),
                              m,
                              l,
                              n[e],
                              20,
                              -373897302
                            )),
                            (v = u(
                              v,
                              (m = u(
                                m,
                                (l = u(l, g, v, m, n[e + 5], 5, -701558691)),
                                g,
                                v,
                                n[e + 10],
                                9,
                                38016083
                              )),
                              l,
                              g,
                              n[e + 15],
                              14,
                              -660478335
                            )),
                            m,
                            l,
                            n[e + 4],
                            20,
                            -405537848
                          )),
                          (v = u(
                            v,
                            (m = u(
                              m,
                              (l = u(l, g, v, m, n[e + 9], 5, 568446438)),
                              g,
                              v,
                              n[e + 14],
                              9,
                              -1019803690
                            )),
                            l,
                            g,
                            n[e + 3],
                            14,
                            -187363961
                          )),
                          m,
                          l,
                          n[e + 8],
                          20,
                          1163531501
                        )),
                        (v = u(
                          v,
                          (m = u(
                            m,
                            (l = u(l, g, v, m, n[e + 13], 5, -1444681467)),
                            g,
                            v,
                            n[e + 2],
                            9,
                            -51403784
                          )),
                          l,
                          g,
                          n[e + 7],
                          14,
                          1735328473
                        )),
                        m,
                        l,
                        n[e + 12],
                        20,
                        -1926607734
                      )),
                      (v = c(
                        v,
                        (m = c(
                          m,
                          (l = c(l, g, v, m, n[e + 5], 4, -378558)),
                          g,
                          v,
                          n[e + 8],
                          11,
                          -2022574463
                        )),
                        l,
                        g,
                        n[e + 11],
                        16,
                        1839030562
                      )),
                      m,
                      l,
                      n[e + 14],
                      23,
                      -35309556
                    )),
                    (v = c(
                      v,
                      (m = c(
                        m,
                        (l = c(l, g, v, m, n[e + 1], 4, -1530992060)),
                        g,
                        v,
                        n[e + 4],
                        11,
                        1272893353
                      )),
                      l,
                      g,
                      n[e + 7],
                      16,
                      -155497632
                    )),
                    m,
                    l,
                    n[e + 10],
                    23,
                    -1094730640
                  )),
                  (v = c(
                    v,
                    (m = c(
                      m,
                      (l = c(l, g, v, m, n[e + 13], 4, 681279174)),
                      g,
                      v,
                      n[e],
                      11,
                      -358537222
                    )),
                    l,
                    g,
                    n[e + 3],
                    16,
                    -722521979
                  )),
                  m,
                  l,
                  n[e + 6],
                  23,
                  76029189
                )),
                (v = c(
                  v,
                  (m = c(
                    m,
                    (l = c(l, g, v, m, n[e + 9], 4, -640364487)),
                    g,
                    v,
                    n[e + 12],
                    11,
                    -421815835
                  )),
                  l,
                  g,
                  n[e + 15],
                  16,
                  530742520
                )),
                m,
                l,
                n[e + 2],
                23,
                -995338651
              )),
              (v = f(
                v,
                (m = f(
                  m,
                  (l = f(l, g, v, m, n[e], 6, -198630844)),
                  g,
                  v,
                  n[e + 7],
                  10,
                  1126891415
                )),
                l,
                g,
                n[e + 14],
                15,
                -1416354905
              )),
              m,
              l,
              n[e + 5],
              21,
              -57434055
            )),
            (v = f(
              v,
              (m = f(
                m,
                (l = f(l, g, v, m, n[e + 12], 6, 1700485571)),
                g,
                v,
                n[e + 3],
                10,
                -1894986606
              )),
              l,
              g,
              n[e + 10],
              15,
              -1051523
            )),
            m,
            l,
            n[e + 1],
            21,
            -2054922799
          )),
          (v = f(
            v,
            (m = f(
              m,
              (l = f(l, g, v, m, n[e + 8], 6, 1873313359)),
              g,
              v,
              n[e + 15],
              10,
              -30611744
            )),
            l,
            g,
            n[e + 6],
            15,
            -1560198380
          )),
          m,
          l,
          n[e + 13],
          21,
          1309151649
        )),
        (v = f(
          v,
          (m = f(
            m,
            (l = f(l, g, v, m, n[e + 4], 6, -145523070)),
            g,
            v,
            n[e + 11],
            10,
            -1120210379
          )),
          l,
          g,
          n[e + 2],
          15,
          718787259
        )),
        m,
        l,
        n[e + 9],
        21,
        -343485551
      )),
      (l = t(l, i)),
      (g = t(g, a)),
      (v = t(v, d)),
      (m = t(m, h));
  return [l, g, v, m];
}
function a(n) {
  var t,
    r = "",
    e = 32 * n.length;
  for (t = 0; t < e; t += 8)
    r += String.fromCharCode((n[t >> 5] >>> t % 32) & 255);
  return r;
}
function d(n) {
  var t,
    r = [];
  for (r[(n.length >> 2) - 1] = void 0, t = 0; t < r.length; t += 1) r[t] = 0;
  var e = 8 * n.length;
  for (t = 0; t < e; t += 8) r[t >> 5] |= (255 & n.charCodeAt(t / 8)) << t % 32;
  return r;
}
function h(n) {
  return a(i(d(n), 8 * n.length));
}
function l(n, t) {
  var r,
    e,
    o = d(n),
    u = [],
    c = [];
  for (
    u[15] = c[15] = void 0, o.length > 16 && (o = i(o, 8 * n.length)), r = 0;
    r < 16;
    r += 1
  )
    (u[r] = 909522486 ^ o[r]), (c[r] = 1549556828 ^ o[r]);
  return (e = i(u.concat(d(t)), 512 + 8 * t.length)), a(i(c.concat(e), 640));
}
function g(n) {
  var t,
    r,
    e = "";
  for (r = 0; r < n.length; r += 1)
    (t = n.charCodeAt(r)),
      (e +=
        "0123456789abcdef".charAt((t >>> 4) & 15) +
        "0123456789abcdef".charAt(15 & t));
  return e;
}
function v(n) {
  return unescape(encodeURIComponent(n));
}
function m(n) {
  return h(v(n));
}
function p(n) {
  return g(m(n));
}
function s(n, t) {
  return l(v(n), v(t));
}
function C(n, t) {
  return g(s(n, t));
}
function md5(n, t, r) {
  return t ? (r ? s(t, n) : C(t, n)) : r ? m(n) : p(n);
}
`

var jsSha1 = `
//js-sha1 v0.6.0
function t(t) {
  t
    ? ((f[0] = f[16] = f[1] = f[2] = f[3] = f[4] = f[5] = f[6] = f[7] = f[8] = f[9] = f[10] = f[11] = f[12] = f[13] = f[14] = f[15] = 0),
      (this.blocks = f))
    : (this.blocks = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]),
    (this.h0 = 1732584193),
    (this.h1 = 4023233417),
    (this.h2 = 2562383102),
    (this.h3 = 271733878),
    (this.h4 = 3285377520),
    (this.block = this.start = this.bytes = this.hBytes = 0),
    (this.finalized = this.hashed = !1),
    (this.first = !0);
}
var h = "object" == typeof window ? window : {},
  s =
    !h.JS_SHA1_NO_NODE_JS &&
    "object" == typeof process &&
    process.versions &&
    process.versions.node;
s && (h = global);
var i = !h.JS_SHA1_NO_COMMON_JS && "object" == typeof module && module.exports,
  e = "function" == typeof define && define.amd,
  r = "0123456789abcdef".split(""),
  o = [-2147483648, 8388608, 32768, 128],
  n = [24, 16, 8, 0],
  a = ["hex", "array", "digest", "arrayBuffer"],
  f = [],
  u = function(h) {
    return function(s) {
      return new t(!0).update(s)[h]();
    };
  },
  c = function() {
    var h = u("hex");
    s && (h = p(h)),
      (h.create = function() {
        return new t();
      }),
      (h.update = function(t) {
        return h.create().update(t);
      });
    for (var i = 0; i < a.length; ++i) {
      var e = a[i];
      h[e] = u(e);
    }
    return h;
  },
  p = function(t) {
    var h = eval("require('crypto')"),
      s = eval("require('buffer').Buffer"),
      i = function(i) {
        if ("string" == typeof i)
          return h
            .createHash("sha1")
            .update(i, "utf8")
            .digest("hex");
        if (i.constructor === ArrayBuffer) i = new Uint8Array(i);
        else if (void 0 === i.length) return t(i);
        return h
          .createHash("sha1")
          .update(new s(i))
          .digest("hex");
      };
    return i;
  };
(t.prototype.update = function(t) {
  if (!this.finalized) {
    var s = "string" != typeof t;
    s && t.constructor === h.ArrayBuffer && (t = new Uint8Array(t));
    for (var i, e, r = 0, o = t.length || 0, a = this.blocks; r < o; ) {
      if (
        (this.hashed &&
          ((this.hashed = !1),
          (a[0] = this.block),
          (a[16] = a[1] = a[2] = a[3] = a[4] = a[5] = a[6] = a[7] = a[8] = a[9] = a[10] = a[11] = a[12] = a[13] = a[14] = a[15] = 0)),
        s)
      )
        for (e = this.start; r < o && e < 64; ++r)
          a[e >> 2] |= t[r] << n[3 & e++];
      else
        for (e = this.start; r < o && e < 64; ++r)
          (i = t.charCodeAt(r)) < 128
            ? (a[e >> 2] |= i << n[3 & e++])
            : i < 2048
            ? ((a[e >> 2] |= (192 | (i >> 6)) << n[3 & e++]),
              (a[e >> 2] |= (128 | (63 & i)) << n[3 & e++]))
            : i < 55296 || i >= 57344
            ? ((a[e >> 2] |= (224 | (i >> 12)) << n[3 & e++]),
              (a[e >> 2] |= (128 | ((i >> 6) & 63)) << n[3 & e++]),
              (a[e >> 2] |= (128 | (63 & i)) << n[3 & e++]))
            : ((i = 65536 + (((1023 & i) << 10) | (1023 & t.charCodeAt(++r)))),
              (a[e >> 2] |= (240 | (i >> 18)) << n[3 & e++]),
              (a[e >> 2] |= (128 | ((i >> 12) & 63)) << n[3 & e++]),
              (a[e >> 2] |= (128 | ((i >> 6) & 63)) << n[3 & e++]),
              (a[e >> 2] |= (128 | (63 & i)) << n[3 & e++]));
      (this.lastByteIndex = e),
        (this.bytes += e - this.start),
        e >= 64
          ? ((this.block = a[16]),
            (this.start = e - 64),
            this.hash(),
            (this.hashed = !0))
          : (this.start = e);
    }
    return (
      this.bytes > 4294967295 &&
        ((this.hBytes += (this.bytes / 4294967296) << 0),
        (this.bytes = this.bytes % 4294967296)),
      this
    );
  }
}),
  (t.prototype.finalize = function() {
    if (!this.finalized) {
      this.finalized = !0;
      var t = this.blocks,
        h = this.lastByteIndex;
      (t[16] = this.block),
        (t[h >> 2] |= o[3 & h]),
        (this.block = t[16]),
        h >= 56 &&
          (this.hashed || this.hash(),
          (t[0] = this.block),
          (t[16] = t[1] = t[2] = t[3] = t[4] = t[5] = t[6] = t[7] = t[8] = t[9] = t[10] = t[11] = t[12] = t[13] = t[14] = t[15] = 0)),
        (t[14] = (this.hBytes << 3) | (this.bytes >>> 29)),
        (t[15] = this.bytes << 3),
        this.hash();
    }
  }),
  (t.prototype.hash = function() {
    var t,
      h,
      s = this.h0,
      i = this.h1,
      e = this.h2,
      r = this.h3,
      o = this.h4,
      n = this.blocks;
    for (t = 16; t < 80; ++t)
      (h = n[t - 3] ^ n[t - 8] ^ n[t - 14] ^ n[t - 16]),
        (n[t] = (h << 1) | (h >>> 31));
    for (t = 0; t < 20; t += 5)
      (s =
        ((h =
          ((i =
            ((h =
              ((e =
                ((h =
                  ((r =
                    ((h =
                      ((o =
                        ((h = (s << 5) | (s >>> 27)) +
                          ((i & e) | (~i & r)) +
                          o +
                          1518500249 +
                          n[t]) <<
                        0) <<
                        5) |
                      (o >>> 27)) +
                      ((s & (i = (i << 30) | (i >>> 2))) | (~s & e)) +
                      r +
                      1518500249 +
                      n[t + 1]) <<
                    0) <<
                    5) |
                  (r >>> 27)) +
                  ((o & (s = (s << 30) | (s >>> 2))) | (~o & i)) +
                  e +
                  1518500249 +
                  n[t + 2]) <<
                0) <<
                5) |
              (e >>> 27)) +
              ((r & (o = (o << 30) | (o >>> 2))) | (~r & s)) +
              i +
              1518500249 +
              n[t + 3]) <<
            0) <<
            5) |
          (i >>> 27)) +
          ((e & (r = (r << 30) | (r >>> 2))) | (~e & o)) +
          s +
          1518500249 +
          n[t + 4]) <<
        0),
        (e = (e << 30) | (e >>> 2));
    for (; t < 40; t += 5)
      (s =
        ((h =
          ((i =
            ((h =
              ((e =
                ((h =
                  ((r =
                    ((h =
                      ((o =
                        ((h = (s << 5) | (s >>> 27)) +
                          (i ^ e ^ r) +
                          o +
                          1859775393 +
                          n[t]) <<
                        0) <<
                        5) |
                      (o >>> 27)) +
                      (s ^ (i = (i << 30) | (i >>> 2)) ^ e) +
                      r +
                      1859775393 +
                      n[t + 1]) <<
                    0) <<
                    5) |
                  (r >>> 27)) +
                  (o ^ (s = (s << 30) | (s >>> 2)) ^ i) +
                  e +
                  1859775393 +
                  n[t + 2]) <<
                0) <<
                5) |
              (e >>> 27)) +
              (r ^ (o = (o << 30) | (o >>> 2)) ^ s) +
              i +
              1859775393 +
              n[t + 3]) <<
            0) <<
            5) |
          (i >>> 27)) +
          (e ^ (r = (r << 30) | (r >>> 2)) ^ o) +
          s +
          1859775393 +
          n[t + 4]) <<
        0),
        (e = (e << 30) | (e >>> 2));
    for (; t < 60; t += 5)
      (s =
        ((h =
          ((i =
            ((h =
              ((e =
                ((h =
                  ((r =
                    ((h =
                      ((o =
                        ((h = (s << 5) | (s >>> 27)) +
                          ((i & e) | (i & r) | (e & r)) +
                          o -
                          1894007588 +
                          n[t]) <<
                        0) <<
                        5) |
                      (o >>> 27)) +
                      ((s & (i = (i << 30) | (i >>> 2))) | (s & e) | (i & e)) +
                      r -
                      1894007588 +
                      n[t + 1]) <<
                    0) <<
                    5) |
                  (r >>> 27)) +
                  ((o & (s = (s << 30) | (s >>> 2))) | (o & i) | (s & i)) +
                  e -
                  1894007588 +
                  n[t + 2]) <<
                0) <<
                5) |
              (e >>> 27)) +
              ((r & (o = (o << 30) | (o >>> 2))) | (r & s) | (o & s)) +
              i -
              1894007588 +
              n[t + 3]) <<
            0) <<
            5) |
          (i >>> 27)) +
          ((e & (r = (r << 30) | (r >>> 2))) | (e & o) | (r & o)) +
          s -
          1894007588 +
          n[t + 4]) <<
        0),
        (e = (e << 30) | (e >>> 2));
    for (; t < 80; t += 5)
      (s =
        ((h =
          ((i =
            ((h =
              ((e =
                ((h =
                  ((r =
                    ((h =
                      ((o =
                        ((h = (s << 5) | (s >>> 27)) +
                          (i ^ e ^ r) +
                          o -
                          899497514 +
                          n[t]) <<
                        0) <<
                        5) |
                      (o >>> 27)) +
                      (s ^ (i = (i << 30) | (i >>> 2)) ^ e) +
                      r -
                      899497514 +
                      n[t + 1]) <<
                    0) <<
                    5) |
                  (r >>> 27)) +
                  (o ^ (s = (s << 30) | (s >>> 2)) ^ i) +
                  e -
                  899497514 +
                  n[t + 2]) <<
                0) <<
                5) |
              (e >>> 27)) +
              (r ^ (o = (o << 30) | (o >>> 2)) ^ s) +
              i -
              899497514 +
              n[t + 3]) <<
            0) <<
            5) |
          (i >>> 27)) +
          (e ^ (r = (r << 30) | (r >>> 2)) ^ o) +
          s -
          899497514 +
          n[t + 4]) <<
        0),
        (e = (e << 30) | (e >>> 2));
    (this.h0 = (this.h0 + s) << 0),
      (this.h1 = (this.h1 + i) << 0),
      (this.h2 = (this.h2 + e) << 0),
      (this.h3 = (this.h3 + r) << 0),
      (this.h4 = (this.h4 + o) << 0);
  }),
  (t.prototype.hex = function() {
    this.finalize();
    var t = this.h0,
      h = this.h1,
      s = this.h2,
      i = this.h3,
      e = this.h4;
    return (
      r[(t >> 28) & 15] +
      r[(t >> 24) & 15] +
      r[(t >> 20) & 15] +
      r[(t >> 16) & 15] +
      r[(t >> 12) & 15] +
      r[(t >> 8) & 15] +
      r[(t >> 4) & 15] +
      r[15 & t] +
      r[(h >> 28) & 15] +
      r[(h >> 24) & 15] +
      r[(h >> 20) & 15] +
      r[(h >> 16) & 15] +
      r[(h >> 12) & 15] +
      r[(h >> 8) & 15] +
      r[(h >> 4) & 15] +
      r[15 & h] +
      r[(s >> 28) & 15] +
      r[(s >> 24) & 15] +
      r[(s >> 20) & 15] +
      r[(s >> 16) & 15] +
      r[(s >> 12) & 15] +
      r[(s >> 8) & 15] +
      r[(s >> 4) & 15] +
      r[15 & s] +
      r[(i >> 28) & 15] +
      r[(i >> 24) & 15] +
      r[(i >> 20) & 15] +
      r[(i >> 16) & 15] +
      r[(i >> 12) & 15] +
      r[(i >> 8) & 15] +
      r[(i >> 4) & 15] +
      r[15 & i] +
      r[(e >> 28) & 15] +
      r[(e >> 24) & 15] +
      r[(e >> 20) & 15] +
      r[(e >> 16) & 15] +
      r[(e >> 12) & 15] +
      r[(e >> 8) & 15] +
      r[(e >> 4) & 15] +
      r[15 & e]
    );
  }),
  (t.prototype.toString = t.prototype.hex),
  (t.prototype.digest = function() {
    this.finalize();
    var t = this.h0,
      h = this.h1,
      s = this.h2,
      i = this.h3,
      e = this.h4;
    return [
      (t >> 24) & 255,
      (t >> 16) & 255,
      (t >> 8) & 255,
      255 & t,
      (h >> 24) & 255,
      (h >> 16) & 255,
      (h >> 8) & 255,
      255 & h,
      (s >> 24) & 255,
      (s >> 16) & 255,
      (s >> 8) & 255,
      255 & s,
      (i >> 24) & 255,
      (i >> 16) & 255,
      (i >> 8) & 255,
      255 & i,
      (e >> 24) & 255,
      (e >> 16) & 255,
      (e >> 8) & 255,
      255 & e
    ];
  }),
  (t.prototype.array = t.prototype.digest),
  (t.prototype.arrayBuffer = function() {
    this.finalize();
    var t = new ArrayBuffer(20),
      h = new DataView(t);
    return (
      h.setUint32(0, this.h0),
      h.setUint32(4, this.h1),
      h.setUint32(8, this.h2),
      h.setUint32(12, this.h3),
      h.setUint32(16, this.h4),
      t
    );
  });
var sha1 = c();
`

func NewUtils() (js *otto.Otto) {
	js = otto.New()
	_, _ = js.Run(jsUtils)
	return
}

func NewBase64() (js *otto.Otto) {
	js = otto.New()
	_, _ = js.Run(jsBase64)
	return
}

func NewMd5() (js *otto.Otto) {
	js = otto.New()
	_, _ = js.Run(jsMd5)
	return
}
func NewSha1() (js *otto.Otto) {
	js = otto.New()
	_, _ = js.Run(jsSha1)
	return
}
