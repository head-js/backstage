(function () {
  'use strict';

  // require('core-js/modules/es.array.iterator.js');
  // require('core-js/modules/web.dom-collections.iterator.js');
  // require('core-js/modules/es.regexp.exec.js');
  // require('core-js/modules/es.string.replace.js');
  // require('core-js/modules/es.promise.js');
  // require('core-js/modules/web.dom-collections.for-each.js');
  // require('core-js/modules/es.object.get-own-property-descriptors.js');
  // require('core-js/modules/es.symbol.description.js');
  // require('core-js/modules/es.string.split.js');
  // require('core-js/modules/es.regexp.constructor.js');
  // require('core-js/modules/es.array.includes.js');
  // require('core-js/modules/es.array.reduce.js');

  // require('core-js/modules/es.array.iterator.js');
  // require('core-js/modules/web.dom-collections.iterator.js');
  // require('core-js/modules/es.promise.js');

  /**
   * Expose compositor.
   */

  var src$1 = compose$3;

  /**
   * Compose `middleware` returning
   * a fully valid middleware comprised
   * of all those which are passed.
   *
   * @param {Array} middleware
   * @return {Function}
   * @api public
   */

  function compose$3(middleware) {
    if (!Array.isArray(middleware)) throw new TypeError('Middleware stack must be an array!');
    for (const fn of middleware) {
      if (typeof fn !== 'function') throw new TypeError('Middleware must be composed of functions!');
    }

    /**
     * @param {Object} context
     * @return {Promise}
     * @api public
     */

    return function (context, next) {
      // last called middleware #
      let index = -1;
      return dispatch(0);
      function dispatch(i) {
        if (i <= index) return Promise.reject(new Error('next() called multiple times'));
        index = i;
        let fn = middleware[i];
        if (i === middleware.length) fn = next;
        if (!fn) return Promise.resolve();
        try {
          return Promise.resolve(fn(context, dispatch.bind(null, i + 1)));
        } catch (err) {
          return Promise.reject(err);
        }
      }
    };
  }
  var compose_1$1 = src$1;

  var contextExports = {};
  var context$1 = {
    get exports(){ return contextExports; },
    set exports(v){ contextExports = v; },
  };

  /**
   * Expose `Delegator`.
   */

  var delegates = Delegator;

  /**
   * Initialize a delegator.
   *
   * @param {Object} proto
   * @param {String} target
   * @api public
   */

  function Delegator(proto, target) {
    if (!(this instanceof Delegator)) return new Delegator(proto, target);
    this.proto = proto;
    this.target = target;
    this.methods = [];
    this.getters = [];
    this.setters = [];
    this.fluents = [];
  }

  /**
   * Delegate method `name`.
   *
   * @param {String} name
   * @return {Delegator} self
   * @api public
   */

  Delegator.prototype.method = function(name){
    var proto = this.proto;
    var target = this.target;
    this.methods.push(name);

    proto[name] = function(){
      return this[target][name].apply(this[target], arguments);
    };

    return this;
  };

  /**
   * Delegator accessor `name`.
   *
   * @param {String} name
   * @return {Delegator} self
   * @api public
   */

  Delegator.prototype.access = function(name){
    return this.getter(name).setter(name);
  };

  /**
   * Delegator getter `name`.
   *
   * @param {String} name
   * @return {Delegator} self
   * @api public
   */

  Delegator.prototype.getter = function(name){
    var proto = this.proto;
    var target = this.target;
    this.getters.push(name);

    proto.__defineGetter__(name, function(){
      return this[target][name];
    });

    return this;
  };

  /**
   * Delegator setter `name`.
   *
   * @param {String} name
   * @return {Delegator} self
   * @api public
   */

  Delegator.prototype.setter = function(name){
    var proto = this.proto;
    var target = this.target;
    this.setters.push(name);

    proto.__defineSetter__(name, function(val){
      return this[target][name] = val;
    });

    return this;
  };

  /**
   * Delegator fluent accessor
   *
   * @param {String} name
   * @return {Delegator} self
   * @api public
   */

  Delegator.prototype.fluent = function (name) {
    var proto = this.proto;
    var target = this.target;
    this.fluents.push(name);

    proto[name] = function(val){
      if ('undefined' != typeof val) {
        this[target][name] = val;
        return this;
      } else {
        return this[target][name];
      }
    };

    return this;
  };

  /**
   * Module dependencies.
   */

  // const util = require('util');
  // const createError = require('http-errors');
  // const httpAssert = require('http-assert');
  const delegate = delegates;
  // const statuses = require('statuses');
  // const Cookies = require('cookies');

  // const COOKIES = Symbol('context#cookies');

  /**
   * Context prototype.
   */

  const proto = context$1.exports = {
    /**
     * util.inspect() implementation, which
     * just returns the JSON output.
     *
     * @return {Object}
     * @api public
     */

    // inspect() {
    //   if (this === proto) return this;
    //   return this.toJSON();
    // },

    /**
     * Return JSON representation.
     *
     * Here we explicitly invoke .toJSON() on each
     * object, as iteration will otherwise fail due
     * to the getters and cause utilities such as
     * clone() to fail.
     *
     * @return {Object}
     * @api public
     */

    // toJSON() {
    //   return {
    //     request: this.request.toJSON(),
    //     response: this.response.toJSON(),
    //     app: this.app.toJSON(),
    //     originalUrl: this.originalUrl,
    //     req: '<original node req>',
    //     res: '<original node res>',
    //     socket: '<original node socket>'
    //   };
    // },

    /**
     * Similar to .throw(), adds assertion.
     *
     *    this.assert(this.user, 401, 'Please login!');
     *
     * See: https://github.com/jshttp/http-assert
     *
     * @param {Mixed} test
     * @param {Number} status
     * @param {String} message
     * @api public
     */

    // assert: httpAssert,

    /**
     * Throw an error with `status` (default 500) and
     * `msg`. Note that these are user-level
     * errors, and the message may be exposed to the client.
     *
     *    this.throw(403)
     *    this.throw(400, 'name required')
     *    this.throw('something exploded')
     *    this.throw(new Error('invalid'))
     *    this.throw(400, new Error('invalid'))
     *
     * See: https://github.com/jshttp/http-errors
     *
     * Note: `status` should only be passed as the first parameter.
     *
     * @param {String|Number|Error} err, msg or status
     * @param {String|Number|Error} [err, msg or status]
     * @param {Object} [props]
     * @api public
     */

    throw(status, message, extra) {
      // throw createError(...args);
      const err = new Error(message || status); // TODO: || statuses[status].message
      err.status = status;
      // TODO: extends(err, extra);
      throw err;
    },
    /**
     * Default error handling.
     *
     * @param {Error} err
     * @api private
     */

    onerror(err) {
      // don't do anything if there is no error.
      // this allows you to pass `this.onerror`
      // to node-style callbacks.
      if (null == err) return;

      // When dealing with cross-globals a normal `instanceof` check doesn't work properly.
      // See https://github.com/koajs/koa/issues/1466
      // We can probably remove it once jest fixes https://github.com/facebook/jest/issues/2549.
      const isNativeError = Object.prototype.toString.call(err) === '[object Error]' || err instanceof Error;
      // if (!isNativeError) err = new Error(util.format('non-error thrown: %j', err));
      if (!isNativeError) {
        if (err.status || err.code) ; else {
          err = new Error('non-error thrown: ' + err.toString());
        }
      }
      throw err;

      // let headerSent = false;
      // if (this.headerSent || !this.writable) {
      //   headerSent = err.headerSent = true;
      // }

      // // delegate
      // this.app.emit('error', err, this);

      // // nothing we can do here other
      // // than delegate to the app-level
      // // handler and log.
      // if (headerSent) {
      //   return;
      // }

      // const { res } = this;

      // // first unset all headers
      // /* istanbul ignore else */
      // if (typeof res.getHeaderNames === 'function') {
      //   res.getHeaderNames().forEach(name => res.removeHeader(name));
      // } else {
      //   res._headers = {}; // Node < 7.7
      // }

      // // then set those specified
      // this.set(err.headers);

      // // force text/plain
      // this.type = 'text';

      // let statusCode = err.status || err.statusCode;

      // // ENOENT support
      // if ('ENOENT' === err.code) statusCode = 404;

      // // default to 500
      // if ('number' !== typeof statusCode || !statuses[statusCode]) statusCode = 500;

      // // respond
      // const code = statuses[statusCode];
      // const msg = err.expose ? err.message : code;
      // this.status = err.status = statusCode;
      // this.length = Buffer.byteLength(msg);
      // res.end(msg);
    }

    // get cookies() {
    //   if (!this[COOKIES]) {
    //     this[COOKIES] = new Cookies(this.req, this.res, {
    //       keys: this.app.keys,
    //       secure: this.request.secure
    //     });
    //   }
    //   return this[COOKIES];
    // },

    // set cookies(_cookies) {
    //   this[COOKIES] = _cookies;
    // }
  };

  /**
   * Custom inspection implementation for newer Node.js versions.
   *
   * @return {Object}
   * @api public
   */

  /* istanbul ignore else */
  // if (util.inspect.custom) {
  //   module.exports[util.inspect.custom] = module.exports.inspect;
  // }

  /**
   * Response delegation.
   */

  delegate(proto, 'response')
  // .method('attachment')
  // .method('redirect')
  // .method('remove')
  // .method('vary')
  // .method('has')
  // .method('set')
  // .method('append')
  // .method('flushHeaders')
  // .access('status')
  // .access('message')
  // .access('body')
  // .access('length')
  // .access('type')
  // .access('lastModified')
  // .access('etag')
  // .getter('headerSent')
  .getter('writable');

  /**
   * Request delegation.
   */

  delegate(proto, 'request')
  // .method('acceptsLanguages')
  // .method('acceptsEncodings')
  // .method('acceptsCharsets')
  // .method('accepts')
  // .method('get')
  // .method('is')
  // .access('querystring')
  // .access('idempotent')
  // .access('socket')
  // .access('search')
  .access('method')
  // .access('query')
  .access('path');

  /**
   * Module dependencies.
   */

  // const URL = require('url').URL;
  // const net = require('net');
  // const accepts = require('accepts');
  // const contentType = require('content-type');
  // const stringify = require('url').format;
  // const parse = require('parseurl');
  // const qs = require('querystring');
  // const typeis = require('type-is');
  // const fresh = require('fresh');
  // const only = require('only');
  // const util = require('util');

  // const IP = Symbol('context#ip');

  /**
   * Prototype.
   */

  var request$1 = {
    /**
     * Return request header.
     *
     * @return {Object}
     * @api public
     */

    // get header() {
    //   return this.req.headers;
    // },

    /**
     * Set request header.
     *
     * @api public
     */

    // set header(val) {
    //   this.req.headers = val;
    // },

    /**
     * Return request header, alias as request.header
     *
     * @return {Object}
     * @api public
     */

    // get headers() {
    //   return this.req.headers;
    // },

    /**
     * Set request header, alias as request.header
     *
     * @api public
     */

    // set headers(val) {
    //   this.req.headers = val;
    // },

    /**
     * Get request URL.
     *
     * @return {String}
     * @api public
     */

    // get url() {
    //   return this.req.url;
    // },

    /**
     * Set request URL.
     *
     * @api public
     */

    // set url(val) {
    //   this.req.url = val;
    // },

    /**
     * Get origin of URL.
     *
     * @return {String}
     * @api public
     */

    // get origin() {
    //   return `${this.protocol}://${this.host}`;
    // },

    /**
     * Get full request URL.
     *
     * @return {String}
     * @api public
     */

    // get href() {
    //   // support: `GET http://example.com/foo`
    //   if (/^https?:\/\//i.test(this.originalUrl)) return this.originalUrl;
    //   return this.origin + this.originalUrl;
    // },

    /**
     * Get request method.
     *
     * @return {String}
     * @api public
     */

    get method() {
      return this.req.method;
    },
    /**
     * Set request method.
     *
     * @param {String} val
     * @api public
     */

    // set method(val) {
    //   this.req.method = val;
    // },

    /**
     * Get request pathname.
     *
     * @return {String}
     * @api public
     */

    get path() {
      // return parse(this.req).pathname;
      return this.req.pathname;
    },
    /**
     * Set pathname, retaining the query string when present.
     *
     * @param {String} path
     * @api public
     */

    // set path(path) {
    //   const url = parse(this.req);
    //   if (url.pathname === path) return;

    //   url.pathname = path;
    //   url.path = null;

    //   this.url = stringify(url);
    // },

    /**
     * Get parsed query string.
     *
     * @return {Object}
     * @api public
     */

    get query() {
      // const str = this.querystring;
      // const c = this._querycache = this._querycache || {};
      // return c[str] || (c[str] = qs.parse(str));
      return this.req.query;
    },
    /**
     * Set query string as an object.
     *
     * @param {Object} obj
     * @api public
     */

    // set query(obj) {
    //   this.querystring = qs.stringify(obj);
    // },

    /**
     * Get query string.
     *
     * @return {String}
     * @api public
     */

    // get querystring() {
    //   if (!this.req) return '';
    //   return parse(this.req).query || '';
    // },

    /**
     * Set query string.
     *
     * @param {String} str
     * @api public
     */

    // set querystring(str) {
    //   const url = parse(this.req);
    //   if (url.search === `?${str}`) return;

    //   url.search = str;
    //   url.path = null;

    //   this.url = stringify(url);
    // },

    /**
     * Get the search string. Same as the query string
     * except it includes the leading ?.
     *
     * @return {String}
     * @api public
     */

    // get search() {
    //   if (!this.querystring) return '';
    //   return `?${this.querystring}`;
    // },

    /**
     * Set the search string. Same as
     * request.querystring= but included for ubiquity.
     *
     * @param {String} str
     * @api public
     */

    // set search(str) {
    //   this.querystring = str;
    // },

    /**
     * The original body from verb is already pretty
     *
     * @return {Object}
     * @api public
     */

    get body() {
      return this.req.body;
    }

    /**
     * Parse the "Host" header field host
     * and support X-Forwarded-Host when a
     * proxy is enabled.
     *
     * @return {String} hostname:port
     * @api public
     */

    // get host() {
    //   const proxy = this.app.proxy;
    //   let host = proxy && this.get('X-Forwarded-Host');
    //   if (!host) {
    //     if (this.req.httpVersionMajor >= 2) host = this.get(':authority');
    //     if (!host) host = this.get('Host');
    //   }
    //   if (!host) return '';
    //   return splitCommaSeparatedValues(host, 1)[0];
    // },

    /**
     * Parse the "Host" header field hostname
     * and support X-Forwarded-Host when a
     * proxy is enabled.
     *
     * @return {String} hostname
     * @api public
     */

    // get hostname() {
    //   const host = this.host;
    //   if (!host) return '';
    //   if ('[' === host[0]) return this.URL.hostname || ''; // IPv6
    //   return host.split(':', 1)[0];
    // },

    /**
     * Get WHATWG parsed URL.
     * Lazily memoized.
     *
     * @return {URL|Object}
     * @api public
     */

    // get URL() {
    //   /* istanbul ignore else */
    //   if (!this.memoizedURL) {
    //     const originalUrl = this.originalUrl || ''; // avoid undefined in template string
    //     try {
    //       this.memoizedURL = new URL(`${this.origin}${originalUrl}`);
    //     } catch (err) {
    //       this.memoizedURL = Object.create(null);
    //     }
    //   }
    //   return this.memoizedURL;
    // },

    /**
     * Check if the request is fresh, aka
     * Last-Modified and/or the ETag
     * still match.
     *
     * @return {Boolean}
     * @api public
     */

    // get fresh() {
    //   const method = this.method;
    //   const s = this.ctx.status;

    //   // GET or HEAD for weak freshness validation only
    //   if ('GET' !== method && 'HEAD' !== method) return false;

    //   // 2xx or 304 as per rfc2616 14.26
    //   if ((s >= 200 && s < 300) || 304 === s) {
    //     return fresh(this.header, this.response.header);
    //   }

    //   return false;
    // },

    /**
     * Check if the request is stale, aka
     * "Last-Modified" and / or the "ETag" for the
     * resource has changed.
     *
     * @return {Boolean}
     * @api public
     */

    // get stale() {
    //   return !this.fresh;
    // },

    /**
     * Check if the request is idempotent.
     *
     * @return {Boolean}
     * @api public
     */

    // get idempotent() {
    //   const methods = ['GET', 'HEAD', 'PUT', 'DELETE', 'OPTIONS', 'TRACE'];
    //   return !!~methods.indexOf(this.method);
    // },

    /**
     * Return the request socket.
     *
     * @return {Connection}
     * @api public
     */

    // get socket() {
    //   return this.req.socket;
    // },

    /**
     * Get the charset when present or undefined.
     *
     * @return {String}
     * @api public
     */

    // get charset() {
    //   try {
    //     const { parameters } = contentType.parse(this.req);
    //     return parameters.charset || '';
    //   } catch (e) {
    //     return '';
    //   }
    // },

    /**
     * Return parsed Content-Length when present.
     *
     * @return {Number}
     * @api public
     */

    // get length() {
    //   const len = this.get('Content-Length');
    //   if (len === '') return;
    //   return ~~len;
    // },

    /**
     * Return the protocol string "http" or "https"
     * when requested with TLS. When the proxy setting
     * is enabled the "X-Forwarded-Proto" header
     * field will be trusted. If you're running behind
     * a reverse proxy that supplies https for you this
     * may be enabled.
     *
     * @return {String}
     * @api public
     */

    // get protocol() {
    //   if (this.socket.encrypted) return 'https';
    //   if (!this.app.proxy) return 'http';
    //   const proto = this.get('X-Forwarded-Proto');
    //   return proto ? splitCommaSeparatedValues(proto, 1)[0] : 'http';
    // },

    /**
     * Shorthand for:
     *
     *    this.protocol == 'https'
     *
     * @return {Boolean}
     * @api public
     */

    // get secure() {
    //   return 'https' === this.protocol;
    // },

    /**
     * When `app.proxy` is `true`, parse
     * the "X-Forwarded-For" ip address list.
     *
     * For example if the value was "client, proxy1, proxy2"
     * you would receive the array `["client", "proxy1", "proxy2"]`
     * where "proxy2" is the furthest down-stream.
     *
     * @return {Array}
     * @api public
     */

    // get ips() {
    //   const proxy = this.app.proxy;
    //   const val = this.get(this.app.proxyIpHeader);
    //   let ips = proxy && val
    //     ? splitCommaSeparatedValues(val)
    //     : [];
    //   if (this.app.maxIpsCount > 0) {
    //     ips = ips.slice(-this.app.maxIpsCount);
    //   }
    //   return ips;
    // },

    /**
     * Return request's remote address
     * When `app.proxy` is `true`, parse
     * the "X-Forwarded-For" ip address list and return the first one
     *
     * @return {String}
     * @api public
     */

    // get ip() {
    //   if (!this[IP]) {
    //     this[IP] = this.ips[0] || this.socket.remoteAddress || '';
    //   }
    //   return this[IP];
    // },

    // set ip(_ip) {
    //   this[IP] = _ip;
    // },

    /**
     * Return subdomains as an array.
     *
     * Subdomains are the dot-separated parts of the host before the main domain
     * of the app. By default, the domain of the app is assumed to be the last two
     * parts of the host. This can be changed by setting `app.subdomainOffset`.
     *
     * For example, if the domain is "tobi.ferrets.example.com":
     * If `app.subdomainOffset` is not set, this.subdomains is
     * `["ferrets", "tobi"]`.
     * If `app.subdomainOffset` is 3, this.subdomains is `["tobi"]`.
     *
     * @return {Array}
     * @api public
     */

    // get subdomains() {
    //   const offset = this.app.subdomainOffset;
    //   const hostname = this.hostname;
    //   if (net.isIP(hostname)) return [];
    //   return hostname
    //     .split('.')
    //     .reverse()
    //     .slice(offset);
    // },

    /**
     * Get accept object.
     * Lazily memoized.
     *
     * @return {Object}
     * @api private
     */

    // get accept() {
    //   return this._accept || (this._accept = accepts(this.req));
    // },

    /**
     * Set accept object.
     *
     * @param {Object}
     * @api private
     */

    // set accept(obj) {
    //   this._accept = obj;
    // },

    /**
     * Check if the given `type(s)` is acceptable, returning
     * the best match when true, otherwise `false`, in which
     * case you should respond with 406 "Not Acceptable".
     *
     * The `type` value may be a single mime type string
     * such as "application/json", the extension name
     * such as "json" or an array `["json", "html", "text/plain"]`. When a list
     * or array is given the _best_ match, if any is returned.
     *
     * Examples:
     *
     *     // Accept: text/html
     *     this.accepts('html');
     *     // => "html"
     *
     *     // Accept: text/*, application/json
     *     this.accepts('html');
     *     // => "html"
     *     this.accepts('text/html');
     *     // => "text/html"
     *     this.accepts('json', 'text');
     *     // => "json"
     *     this.accepts('application/json');
     *     // => "application/json"
     *
     *     // Accept: text/*, application/json
     *     this.accepts('image/png');
     *     this.accepts('png');
     *     // => false
     *
     *     // Accept: text/*;q=.5, application/json
     *     this.accepts(['html', 'json']);
     *     this.accepts('html', 'json');
     *     // => "json"
     *
     * @param {String|Array} type(s)...
     * @return {String|Array|false}
     * @api public
     */

    // accepts(...args) {
    //   return this.accept.types(...args);
    // },

    /**
     * Return accepted encodings or best fit based on `encodings`.
     *
     * Given `Accept-Encoding: gzip, deflate`
     * an array sorted by quality is returned:
     *
     *     ['gzip', 'deflate']
     *
     * @param {String|Array} encoding(s)...
     * @return {String|Array}
     * @api public
     */

    // acceptsEncodings(...args) {
    //   return this.accept.encodings(...args);
    // },

    /**
     * Return accepted charsets or best fit based on `charsets`.
     *
     * Given `Accept-Charset: utf-8, iso-8859-1;q=0.2, utf-7;q=0.5`
     * an array sorted by quality is returned:
     *
     *     ['utf-8', 'utf-7', 'iso-8859-1']
     *
     * @param {String|Array} charset(s)...
     * @return {String|Array}
     * @api public
     */

    // acceptsCharsets(...args) {
    //   return this.accept.charsets(...args);
    // },

    /**
     * Return accepted languages or best fit based on `langs`.
     *
     * Given `Accept-Language: en;q=0.8, es, pt`
     * an array sorted by quality is returned:
     *
     *     ['es', 'pt', 'en']
     *
     * @param {String|Array} lang(s)...
     * @return {Array|String}
     * @api public
     */

    // acceptsLanguages(...args) {
    //   return this.accept.languages(...args);
    // },

    /**
     * Check if the incoming request contains the "Content-Type"
     * header field and if it contains any of the given mime `type`s.
     * If there is no request body, `null` is returned.
     * If there is no content type, `false` is returned.
     * Otherwise, it returns the first `type` that matches.
     *
     * Examples:
     *
     *     // With Content-Type: text/html; charset=utf-8
     *     this.is('html'); // => 'html'
     *     this.is('text/html'); // => 'text/html'
     *     this.is('text/*', 'application/json'); // => 'text/html'
     *
     *     // When Content-Type is application/json
     *     this.is('json', 'urlencoded'); // => 'json'
     *     this.is('application/json'); // => 'application/json'
     *     this.is('html', 'application/*'); // => 'application/json'
     *
     *     this.is('html'); // => false
     *
     * @param {String|String[]} [type]
     * @param {String[]} [types]
     * @return {String|false|null}
     * @api public
     */

    // is(type, ...types) {
    //   return typeis(this.req, type, ...types);
    // },

    /**
     * Return the request mime type void of
     * parameters such as "charset".
     *
     * @return {String}
     * @api public
     */

    // get type() {
    //   const type = this.get('Content-Type');
    //   if (!type) return '';
    //   return type.split(';')[0];
    // },

    /**
     * Return request header.
     *
     * The `Referrer` header field is special-cased,
     * both `Referrer` and `Referer` are interchangeable.
     *
     * Examples:
     *
     *     this.get('Content-Type');
     *     // => "text/plain"
     *
     *     this.get('content-type');
     *     // => "text/plain"
     *
     *     this.get('Something');
     *     // => ''
     *
     * @param {String} field
     * @return {String}
     * @api public
     */

    // get(field) {
    //   const req = this.req;
    //   switch (field = field.toLowerCase()) {
    //     case 'referer':
    //     case 'referrer':
    //       return req.headers.referrer || req.headers.referer || '';
    //     default:
    //       return req.headers[field] || '';
    //   }
    // },

    /**
     * Inspect implementation.
     *
     * @return {Object}
     * @api public
     */

    // inspect() {
    //   if (!this.req) return;
    //   return this.toJSON();
    // },

    /**
     * Return JSON representation.
     *
     * @return {Object}
     * @api public
     */

    // toJSON() {
    //   return only(this, [
    //     'method',
    //     'url',
    //     'header'
    //   ]);
    // }
  };

  /**
   * Module dependencies.
   */

  // const contentDisposition = require('content-disposition');
  // const getType = require('cache-content-type');
  // const onFinish = require('on-finished');
  // const escape = require('escape-html');
  // const typeis = require('type-is').is;
  // const statuses = require('statuses');
  // const destroy = require('destroy');
  // const assert = require('assert');
  // const extname = require('path').extname;
  // const vary = require('vary');
  // const only = require('only');
  // const util = require('util');
  // const encodeUrl = require('encodeurl');
  // const Stream = require('stream');
  // const URL = require('url').URL;

  /**
   * Prototype.
   */

  var response$1 = {
    /**
     * Return the request socket.
     *
     * @return {Connection}
     * @api public
     */

    // get socket() {
    //   return this.res.socket;
    // },

    /**
     * Return response header.
     *
     * @return {Object}
     * @api public
     */

    // get header() {
    //   const { res } = this;
    //   return typeof res.getHeaders === 'function'
    //     ? res.getHeaders()
    //     : res._headers || {}; // Node < 7.7
    // },

    /**
     * Return response header, alias as response.header
     *
     * @return {Object}
     * @api public
     */

    // get headers() {
    //   return this.header;
    // },

    /**
     * Get response status code.
     *
     * @return {Number}
     * @api public
     */

    // get status() {
    //   return this.res.statusCode;
    // },

    /**
     * Set response status code.
     *
     * @param {Number} code
     * @api public
     */

    // set status(code) {
    //   if (this.headerSent) return;

    //   assert(Number.isInteger(code), 'status code must be a number');
    //   assert(code >= 100 && code <= 999, `invalid status code: ${code}`);
    //   this._explicitStatus = true;
    //   this.res.statusCode = code;
    //   if (this.req.httpVersionMajor < 2) this.res.statusMessage = statuses[code];
    //   if (this.body && statuses.empty[code]) this.body = null;
    // },

    /**
     * Get response status message
     *
     * @return {String}
     * @api public
     */

    // get message() {
    //   return this.res.statusMessage || statuses[this.status];
    // },

    /**
     * Set response status message
     *
     * @param {String} msg
     * @api public
     */

    // set message(msg) {
    //   this.res.statusMessage = msg;
    // },

    /**
     * Get response body.
     *
     * @return {Mixed}
     * @api public
     */

    // get body() {
    //   return this._body;
    // },

    /**
     * Set response body.
     *
     * @param {String|Buffer|Object|Stream} val
     * @api public
     */

    // set body(val) {
    //   const original = this._body;
    //   this._body = val;

    //   // no content
    //   if (null == val) {
    //     if (!statuses.empty[this.status]) this.status = 204;
    //     if (val === null) this._explicitNullBody = true;
    //     this.remove('Content-Type');
    //     this.remove('Content-Length');
    //     this.remove('Transfer-Encoding');
    //     return;
    //   }

    //   // set the status
    //   if (!this._explicitStatus) this.status = 200;

    //   // set the content-type only if not yet set
    //   const setType = !this.has('Content-Type');

    //   // string
    //   if ('string' === typeof val) {
    //     if (setType) this.type = /^\s*</.test(val) ? 'html' : 'text';
    //     this.length = Buffer.byteLength(val);
    //     return;
    //   }

    //   // buffer
    //   if (Buffer.isBuffer(val)) {
    //     if (setType) this.type = 'bin';
    //     this.length = val.length;
    //     return;
    //   }

    //   // stream
    //   if (val instanceof Stream) {
    //     onFinish(this.res, destroy.bind(null, val));
    //     if (original != val) {
    //       val.once('error', err => this.ctx.onerror(err));
    //       // overwriting
    //       if (null != original) this.remove('Content-Length');
    //     }

    //     if (setType) this.type = 'bin';
    //     return;
    //   }

    //   // json
    //   this.remove('Content-Length');
    //   this.type = 'json';
    // },

    /**
     * Set Content-Length field to `n`.
     *
     * @param {Number} n
     * @api public
     */

    // set length(n) {
    //   if (!this.has('Transfer-Encoding')) {
    //     this.set('Content-Length', n);
    //   }
    // },

    /**
     * Return parsed response Content-Length when present.
     *
     * @return {Number}
     * @api public
     */

    // get length() {
    //   if (this.has('Content-Length')) {
    //     return parseInt(this.get('Content-Length'), 10) || 0;
    //   }

    //   const { body } = this;
    //   if (!body || body instanceof Stream) return undefined;
    //   if ('string' === typeof body) return Buffer.byteLength(body);
    //   if (Buffer.isBuffer(body)) return body.length;
    //   return Buffer.byteLength(JSON.stringify(body));
    // },

    /**
     * Check if a header has been written to the socket.
     *
     * @return {Boolean}
     * @api public
     */

    // get headerSent() {
    //   return this.res.headersSent;
    // },

    /**
     * Vary on `field`.
     *
     * @param {String} field
     * @api public
     */

    // vary(field) {
    //   if (this.headerSent) return;

    //   vary(this.res, field);
    // },

    // _getBackReferrer() {
    //   const referrer = this.ctx.get('Referrer');
    //   if (referrer) {
    //     // referrer is an absolute URL, check if it's the same origin
    //     const url = new URL(referrer, this.ctx.href);
    //     if (url.host === this.ctx.host) {
    //       return referrer;
    //     }
    //   }
    // },

    /**
     * Perform a 302 redirect to `url`.
     *
     * The string "back" is special-cased
     * to provide Referrer support, when Referrer
     * is not present `alt` or "/" is used.
     *
     * Examples:
     *
     *    this.redirect('back');
     *    this.redirect('back', '/index.html');
     *    this.redirect('/login');
     *    this.redirect('http://google.com');
     *
     * @param {String} url
     * @param {String} [alt]
     * @api public
     */

    // redirect(url, alt) {
    //   // location
    //   if ('back' === url) {
    //     url = this._getBackReferrer() || alt || '/';
    //   }

    //   if (/^https?:\/\//i.test(url)) {
    //     // formatting url again avoid security escapes
    //     url = new URL(url).toString();
    //   }
    //   this.set('Location', encodeUrl(url));

    //   // status
    //   if (!statuses.redirect[this.status]) this.status = 302;

    //   // html
    //   if (this.ctx.accepts('html')) {
    //     url = escape(url);
    //     this.type = 'text/html; charset=utf-8';
    //     this.body = `Redirecting to ${url}.`;
    //     return;
    //   }

    //   // text
    //   this.type = 'text/plain; charset=utf-8';
    //   this.body = `Redirecting to ${url}.`;
    // },

    /**
     * Set Content-Disposition header to "attachment" with optional `filename`.
     *
     * @param {String} filename
     * @api public
     */

    // attachment(filename, options) {
    //   if (filename) this.type = extname(filename);
    //   this.set('Content-Disposition', contentDisposition(filename, options));
    // },

    /**
     * Set Content-Type response header with `type` through `mime.lookup()`
     * when it does not contain a charset.
     *
     * Examples:
     *
     *     this.type = '.html';
     *     this.type = 'html';
     *     this.type = 'json';
     *     this.type = 'application/json';
     *     this.type = 'png';
     *
     * @param {String} type
     * @api public
     */

    // set type(type) {
    //   type = getType(type);
    //   if (type) {
    //     this.set('Content-Type', type);
    //   } else {
    //     this.remove('Content-Type');
    //   }
    // },

    /**
     * Set the Last-Modified date using a string or a Date.
     *
     *     this.response.lastModified = new Date();
     *     this.response.lastModified = '2013-09-13';
     *
     * @param {String|Date} type
     * @api public
     */

    // set lastModified(val) {
    //   if ('string' === typeof val) val = new Date(val);
    //   this.set('Last-Modified', val.toUTCString());
    // },

    /**
     * Get the Last-Modified date in Date form, if it exists.
     *
     * @return {Date}
     * @api public
     */

    // get lastModified() {
    //   const date = this.get('last-modified');
    //   if (date) return new Date(date);
    // },

    /**
     * Set the ETag of a response.
     * This will normalize the quotes if necessary.
     *
     *     this.response.etag = 'md5hashsum';
     *     this.response.etag = '"md5hashsum"';
     *     this.response.etag = 'W/"123456789"';
     *
     * @param {String} etag
     * @api public
     */

    // set etag(val) {
    //   if (!/^(W\/)?"/.test(val)) val = `"${val}"`;
    //   this.set('ETag', val);
    // },

    /**
     * Get the ETag of a response.
     *
     * @return {String}
     * @api public
     */

    // get etag() {
    //   return this.get('ETag');
    // },

    /**
     * Return the response mime type void of
     * parameters such as "charset".
     *
     * @return {String}
     * @api public
     */

    // get type() {
    //   const type = this.get('Content-Type');
    //   if (!type) return '';
    //   return type.split(';', 1)[0];
    // },

    /**
     * Check whether the response is one of the listed types.
     * Pretty much the same as `this.request.is()`.
     *
     * @param {String|String[]} [type]
     * @param {String[]} [types]
     * @return {String|false}
     * @api public
     */

    // is(type, ...types) {
    //   return typeis(this.type, type, ...types);
    // },

    /**
     * Return response header.
     *
     * Examples:
     *
     *     this.get('Content-Type');
     *     // => "text/plain"
     *
     *     this.get('content-type');
     *     // => "text/plain"
     *
     * @param {String} field
     * @return {String}
     * @api public
     */

    // get(field) {
    //   return this.header[field.toLowerCase()] || '';
    // },

    /**
     * Returns true if the header identified by name is currently set in the outgoing headers.
     * The header name matching is case-insensitive.
     *
     * Examples:
     *
     *     this.has('Content-Type');
     *     // => true
     *
     *     this.get('content-type');
     *     // => true
     *
     * @param {String} field
     * @return {boolean}
     * @api public
     */

    // has(field) {
    //   return typeof this.res.hasHeader === 'function'
    //     ? this.res.hasHeader(field)
    //     // Node < 7.7
    //     : field.toLowerCase() in this.headers;
    // },

    /**
     * Set header `field` to `val` or pass
     * an object of header fields.
     *
     * Examples:
     *
     *    this.set('Foo', ['bar', 'baz']);
     *    this.set('Accept', 'application/json');
     *    this.set({ Accept: 'text/plain', 'X-API-Key': 'tobi' });
     *
     * @param {String|Object|Array} field
     * @param {String} val
     * @api public
     */

    // set(field, val) {
    //   if (this.headerSent) return;

    //   if (2 === arguments.length) {
    //     if (Array.isArray(val)) val = val.map(v => typeof v === 'string' ? v : String(v));
    //     else if (typeof val !== 'string') val = String(val);
    //     this.res.setHeader(field, val);
    //   } else {
    //     for (const key in field) {
    //       this.set(key, field[key]);
    //     }
    //   }
    // },

    /**
     * Append additional header `field` with value `val`.
     *
     * Examples:
     *
     * ```
     * this.append('Link', ['<http://localhost/>', '<http://localhost:3000/>']);
     * this.append('Set-Cookie', 'foo=bar; Path=/; HttpOnly');
     * this.append('Warning', '199 Miscellaneous warning');
     * ```
     *
     * @param {String} field
     * @param {String|Array} val
     * @api public
     */

    // append(field, val) {
    //   const prev = this.get(field);

    //   if (prev) {
    //     val = Array.isArray(prev)
    //       ? prev.concat(val)
    //       : [prev].concat(val);
    //   }

    //   return this.set(field, val);
    // },

    /**
     * Remove header `field`.
     *
     * @param {String} name
     * @api public
     */

    // remove(field) {
    //   if (this.headerSent) return;

    //   this.res.removeHeader(field);
    // },

    /**
     * Checks if the request is writable.
     * Tests for the existence of the socket
     * as node sometimes does not set it.
     *
     * @return {Boolean}
     * @api private
     */

    get writable() {
      // // can't write any more after response finished
      // // response.writableEnded is available since Node > 12.9
      // // https://nodejs.org/api/http.html#http_response_writableended
      // // response.finished is undocumented feature of previous Node versions
      // // https://stackoverflow.com/questions/16254385/undocumented-response-finished-in-node-js
      // if (this.res.writableEnded || this.res.finished) return false;

      // const socket = this.res.socket;
      // // There are already pending outgoing res, but still writable
      // // https://github.com/nodejs/node/blob/v4.4.7/lib/_http_server.js#L486
      // if (!socket) return true;
      // return socket.writable;
      return true;
    }

    /**
     * Inspect implementation.
     *
     * @return {Object}
     * @api public
     */

    // inspect() {
    //   if (!this.res) return;
    //   const o = this.toJSON();
    //   o.body = this.body;
    //   return o;
    // },

    /**
     * Return JSON representation.
     *
     * @return {Object}
     * @api public
     */

    // toJSON() {
    //   return only(this, [
    //     'status',
    //     'message',
    //     'header'
    //   ]);
    // },

    /**
     * Flush any set headers and begin the body
     */

    // flushHeaders() {
    //   this.res.flushHeaders();
    // }
  };

  function Client$1(opts) {
    this.handleRequest = opts.handleRequest;
  }
  Client$1.prototype.verb = function (method, pathname, search, form, headers) {
    // console.log(method, pathname, params, form, headers);
    const req = {
      method,
      pathname,
      query: search || {},
      body: form || {},
      headers: headers || {}
    };
    const res = {
      end: resp => Promise.resolve(resp)
    };
    return this.handleRequest(req, res);
  };
  Client$1.prototype.get = function (pathname, search, form, headers) {
    return this.verb('GET', pathname, search, form, headers);
  };
  Client$1.prototype.post = function (pathname, search, form, headers) {
    return this.verb('POST', pathname, search, form, headers);
  };
  Client$1.prototype.put = function (pathname, search, form, headers) {
    return this.verb('PUT', pathname, search, form, headers);
  };
  Client$1.prototype.del = function (pathname, search, form, headers) {
    return this.verb('DELETE', pathname, search, form, headers);
  };
  var client$1 = Client$1;

  // require('core-js/modules/es.regexp.exec.js');
  // require('core-js/modules/es.array.reduce.js');
  // require('core-js/modules/es.array.iterator.js');
  // require('core-js/modules/web.dom-collections.iterator.js');
  // require('core-js/modules/es.array.includes.js');
  // require('core-js/modules/es.regexp.constructor.js');
  // require('core-js/modules/es.promise.js');
  // require('core-js/modules/es.string.replace.js');

  function _defineProperty(e, r, t) {
    return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, {
      value: t,
      enumerable: !0,
      configurable: !0,
      writable: !0
    }) : e[r] = t, e;
  }
  function ownKeys(e, r) {
    var t = Object.keys(e);
    if (Object.getOwnPropertySymbols) {
      var o = Object.getOwnPropertySymbols(e);
      r && (o = o.filter(function (r) {
        return Object.getOwnPropertyDescriptor(e, r).enumerable;
      })), t.push.apply(t, o);
    }
    return t;
  }
  function _objectSpread2(e) {
    for (var r = 1; r < arguments.length; r++) {
      var t = null != arguments[r] ? arguments[r] : {};
      r % 2 ? ownKeys(Object(t), !0).forEach(function (r) {
        _defineProperty(e, r, t[r]);
      }) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(t)) : ownKeys(Object(t)).forEach(function (r) {
        Object.defineProperty(e, r, Object.getOwnPropertyDescriptor(t, r));
      });
    }
    return e;
  }
  function _toPrimitive(t, r) {
    if ("object" != typeof t || !t) return t;
    var e = t[Symbol.toPrimitive];
    if (void 0 !== e) {
      var i = e.call(t, r || "default");
      if ("object" != typeof i) return i;
      throw new TypeError("@@toPrimitive must return a primitive value.");
    }
    return ("string" === r ? String : Number)(t);
  }
  function _toPropertyKey(t) {
    var i = _toPrimitive(t, "string");
    return "symbol" == typeof i ? i : i + "";
  }

  // require('core-js/modules/es.array.iterator.js');
  // require('core-js/modules/web.dom-collections.iterator.js');
  // require('core-js/modules/es.promise.js');

  /**
   * Expose compositor.
   */

  var src = compose$1;

  /**
   * Compose `middleware` returning
   * a fully valid middleware comprised
   * of all those which are passed.
   *
   * @param {Array} middleware
   * @return {Function}
   * @api public
   */

  function compose$1(middleware) {
    if (!Array.isArray(middleware)) throw new TypeError('Middleware stack must be an array!');
    for (const fn of middleware) {
      if (typeof fn !== 'function') throw new TypeError('Middleware must be composed of functions!');
    }

    /**
     * @param {Object} context
     * @return {Promise}
     * @api public
     */

    return function (context, next) {
      // last called middleware #
      let index = -1;
      return dispatch(0);
      function dispatch(i) {
        if (i <= index) return Promise.reject(new Error('next() called multiple times'));
        index = i;
        let fn = middleware[i];
        if (i === middleware.length) fn = next;
        if (!fn) return Promise.resolve();
        try {
          return Promise.resolve(fn(context, dispatch.bind(null, i + 1)));
        } catch (err) {
          return Promise.reject(err);
        }
      }
    };
  }
  var compose_1 = src;

  // https://github.com/jshttp/methods

  var methods$1 = ['ACL', 'BIND', 'CHECKOUT', 'CONNECT', 'COPY', 'DELETE', 'GET', 'HEAD', 'LINK', 'LOCK', 'M-SEARCH', 'MERGE', 'MKACTIVITY', 'MKCALENDAR', 'MKCOL', 'MOVE', 'NOTIFY', 'OPTIONS', 'PATCH', 'POST', 'PRI', 'PROPFIND', 'PROPPATCH', 'PURGE', 'PUT', 'REBIND', 'REPORT', 'SEARCH', 'SOURCE', 'SUBSCRIBE', 'TRACE', 'UNBIND', 'UNLINK', 'UNLOCK', 'UNSUBSCRIBE'];
  var dist = {};
  Object.defineProperty(dist, "__esModule", {
    value: true
  });
  dist.TokenData = void 0;
  dist.parse = parse$1;
  dist.compile = compile$1;
  dist.match = match;
  dist.pathToRegexp = pathToRegexp$1;
  dist.stringify = stringify;
  const DEFAULT_DELIMITER = "/";
  const NOOP_VALUE = value => value;
  const ID_START = /^(?:[\$A-Z_a-z\xAA\xB5\xBA\xC0-\xD6\xD8-\xF6\xF8-\u02C1\u02C6-\u02D1\u02E0-\u02E4\u02EC\u02EE\u0370-\u0374\u0376\u0377\u037A-\u037D\u037F\u0386\u0388-\u038A\u038C\u038E-\u03A1\u03A3-\u03F5\u03F7-\u0481\u048A-\u052F\u0531-\u0556\u0559\u0560-\u0588\u05D0-\u05EA\u05EF-\u05F2\u0620-\u064A\u066E\u066F\u0671-\u06D3\u06D5\u06E5\u06E6\u06EE\u06EF\u06FA-\u06FC\u06FF\u0710\u0712-\u072F\u074D-\u07A5\u07B1\u07CA-\u07EA\u07F4\u07F5\u07FA\u0800-\u0815\u081A\u0824\u0828\u0840-\u0858\u0860-\u086A\u0870-\u0887\u0889-\u088F\u08A0-\u08C9\u0904-\u0939\u093D\u0950\u0958-\u0961\u0971-\u0980\u0985-\u098C\u098F\u0990\u0993-\u09A8\u09AA-\u09B0\u09B2\u09B6-\u09B9\u09BD\u09CE\u09DC\u09DD\u09DF-\u09E1\u09F0\u09F1\u09FC\u0A05-\u0A0A\u0A0F\u0A10\u0A13-\u0A28\u0A2A-\u0A30\u0A32\u0A33\u0A35\u0A36\u0A38\u0A39\u0A59-\u0A5C\u0A5E\u0A72-\u0A74\u0A85-\u0A8D\u0A8F-\u0A91\u0A93-\u0AA8\u0AAA-\u0AB0\u0AB2\u0AB3\u0AB5-\u0AB9\u0ABD\u0AD0\u0AE0\u0AE1\u0AF9\u0B05-\u0B0C\u0B0F\u0B10\u0B13-\u0B28\u0B2A-\u0B30\u0B32\u0B33\u0B35-\u0B39\u0B3D\u0B5C\u0B5D\u0B5F-\u0B61\u0B71\u0B83\u0B85-\u0B8A\u0B8E-\u0B90\u0B92-\u0B95\u0B99\u0B9A\u0B9C\u0B9E\u0B9F\u0BA3\u0BA4\u0BA8-\u0BAA\u0BAE-\u0BB9\u0BD0\u0C05-\u0C0C\u0C0E-\u0C10\u0C12-\u0C28\u0C2A-\u0C39\u0C3D\u0C58-\u0C5A\u0C5C\u0C5D\u0C60\u0C61\u0C80\u0C85-\u0C8C\u0C8E-\u0C90\u0C92-\u0CA8\u0CAA-\u0CB3\u0CB5-\u0CB9\u0CBD\u0CDC-\u0CDE\u0CE0\u0CE1\u0CF1\u0CF2\u0D04-\u0D0C\u0D0E-\u0D10\u0D12-\u0D3A\u0D3D\u0D4E\u0D54-\u0D56\u0D5F-\u0D61\u0D7A-\u0D7F\u0D85-\u0D96\u0D9A-\u0DB1\u0DB3-\u0DBB\u0DBD\u0DC0-\u0DC6\u0E01-\u0E30\u0E32\u0E33\u0E40-\u0E46\u0E81\u0E82\u0E84\u0E86-\u0E8A\u0E8C-\u0EA3\u0EA5\u0EA7-\u0EB0\u0EB2\u0EB3\u0EBD\u0EC0-\u0EC4\u0EC6\u0EDC-\u0EDF\u0F00\u0F40-\u0F47\u0F49-\u0F6C\u0F88-\u0F8C\u1000-\u102A\u103F\u1050-\u1055\u105A-\u105D\u1061\u1065\u1066\u106E-\u1070\u1075-\u1081\u108E\u10A0-\u10C5\u10C7\u10CD\u10D0-\u10FA\u10FC-\u1248\u124A-\u124D\u1250-\u1256\u1258\u125A-\u125D\u1260-\u1288\u128A-\u128D\u1290-\u12B0\u12B2-\u12B5\u12B8-\u12BE\u12C0\u12C2-\u12C5\u12C8-\u12D6\u12D8-\u1310\u1312-\u1315\u1318-\u135A\u1380-\u138F\u13A0-\u13F5\u13F8-\u13FD\u1401-\u166C\u166F-\u167F\u1681-\u169A\u16A0-\u16EA\u16EE-\u16F8\u1700-\u1711\u171F-\u1731\u1740-\u1751\u1760-\u176C\u176E-\u1770\u1780-\u17B3\u17D7\u17DC\u1820-\u1878\u1880-\u18A8\u18AA\u18B0-\u18F5\u1900-\u191E\u1950-\u196D\u1970-\u1974\u1980-\u19AB\u19B0-\u19C9\u1A00-\u1A16\u1A20-\u1A54\u1AA7\u1B05-\u1B33\u1B45-\u1B4C\u1B83-\u1BA0\u1BAE\u1BAF\u1BBA-\u1BE5\u1C00-\u1C23\u1C4D-\u1C4F\u1C5A-\u1C7D\u1C80-\u1C8A\u1C90-\u1CBA\u1CBD-\u1CBF\u1CE9-\u1CEC\u1CEE-\u1CF3\u1CF5\u1CF6\u1CFA\u1D00-\u1DBF\u1E00-\u1F15\u1F18-\u1F1D\u1F20-\u1F45\u1F48-\u1F4D\u1F50-\u1F57\u1F59\u1F5B\u1F5D\u1F5F-\u1F7D\u1F80-\u1FB4\u1FB6-\u1FBC\u1FBE\u1FC2-\u1FC4\u1FC6-\u1FCC\u1FD0-\u1FD3\u1FD6-\u1FDB\u1FE0-\u1FEC\u1FF2-\u1FF4\u1FF6-\u1FFC\u2071\u207F\u2090-\u209C\u2102\u2107\u210A-\u2113\u2115\u2118-\u211D\u2124\u2126\u2128\u212A-\u2139\u213C-\u213F\u2145-\u2149\u214E\u2160-\u2188\u2C00-\u2CE4\u2CEB-\u2CEE\u2CF2\u2CF3\u2D00-\u2D25\u2D27\u2D2D\u2D30-\u2D67\u2D6F\u2D80-\u2D96\u2DA0-\u2DA6\u2DA8-\u2DAE\u2DB0-\u2DB6\u2DB8-\u2DBE\u2DC0-\u2DC6\u2DC8-\u2DCE\u2DD0-\u2DD6\u2DD8-\u2DDE\u3005-\u3007\u3021-\u3029\u3031-\u3035\u3038-\u303C\u3041-\u3096\u309B-\u309F\u30A1-\u30FA\u30FC-\u30FF\u3105-\u312F\u3131-\u318E\u31A0-\u31BF\u31F0-\u31FF\u3400-\u4DBF\u4E00-\uA48C\uA4D0-\uA4FD\uA500-\uA60C\uA610-\uA61F\uA62A\uA62B\uA640-\uA66E\uA67F-\uA69D\uA6A0-\uA6EF\uA717-\uA71F\uA722-\uA788\uA78B-\uA7DC\uA7F1-\uA801\uA803-\uA805\uA807-\uA80A\uA80C-\uA822\uA840-\uA873\uA882-\uA8B3\uA8F2-\uA8F7\uA8FB\uA8FD\uA8FE\uA90A-\uA925\uA930-\uA946\uA960-\uA97C\uA984-\uA9B2\uA9CF\uA9E0-\uA9E4\uA9E6-\uA9EF\uA9FA-\uA9FE\uAA00-\uAA28\uAA40-\uAA42\uAA44-\uAA4B\uAA60-\uAA76\uAA7A\uAA7E-\uAAAF\uAAB1\uAAB5\uAAB6\uAAB9-\uAABD\uAAC0\uAAC2\uAADB-\uAADD\uAAE0-\uAAEA\uAAF2-\uAAF4\uAB01-\uAB06\uAB09-\uAB0E\uAB11-\uAB16\uAB20-\uAB26\uAB28-\uAB2E\uAB30-\uAB5A\uAB5C-\uAB69\uAB70-\uABE2\uAC00-\uD7A3\uD7B0-\uD7C6\uD7CB-\uD7FB\uF900-\uFA6D\uFA70-\uFAD9\uFB00-\uFB06\uFB13-\uFB17\uFB1D\uFB1F-\uFB28\uFB2A-\uFB36\uFB38-\uFB3C\uFB3E\uFB40\uFB41\uFB43\uFB44\uFB46-\uFBB1\uFBD3-\uFD3D\uFD50-\uFD8F\uFD92-\uFDC7\uFDF0-\uFDFB\uFE70-\uFE74\uFE76-\uFEFC\uFF21-\uFF3A\uFF41-\uFF5A\uFF66-\uFFBE\uFFC2-\uFFC7\uFFCA-\uFFCF\uFFD2-\uFFD7\uFFDA-\uFFDC]|\uD800[\uDC00-\uDC0B\uDC0D-\uDC26\uDC28-\uDC3A\uDC3C\uDC3D\uDC3F-\uDC4D\uDC50-\uDC5D\uDC80-\uDCFA\uDD40-\uDD74\uDE80-\uDE9C\uDEA0-\uDED0\uDF00-\uDF1F\uDF2D-\uDF4A\uDF50-\uDF75\uDF80-\uDF9D\uDFA0-\uDFC3\uDFC8-\uDFCF\uDFD1-\uDFD5]|\uD801[\uDC00-\uDC9D\uDCB0-\uDCD3\uDCD8-\uDCFB\uDD00-\uDD27\uDD30-\uDD63\uDD70-\uDD7A\uDD7C-\uDD8A\uDD8C-\uDD92\uDD94\uDD95\uDD97-\uDDA1\uDDA3-\uDDB1\uDDB3-\uDDB9\uDDBB\uDDBC\uDDC0-\uDDF3\uDE00-\uDF36\uDF40-\uDF55\uDF60-\uDF67\uDF80-\uDF85\uDF87-\uDFB0\uDFB2-\uDFBA]|\uD802[\uDC00-\uDC05\uDC08\uDC0A-\uDC35\uDC37\uDC38\uDC3C\uDC3F-\uDC55\uDC60-\uDC76\uDC80-\uDC9E\uDCE0-\uDCF2\uDCF4\uDCF5\uDD00-\uDD15\uDD20-\uDD39\uDD40-\uDD59\uDD80-\uDDB7\uDDBE\uDDBF\uDE00\uDE10-\uDE13\uDE15-\uDE17\uDE19-\uDE35\uDE60-\uDE7C\uDE80-\uDE9C\uDEC0-\uDEC7\uDEC9-\uDEE4\uDF00-\uDF35\uDF40-\uDF55\uDF60-\uDF72\uDF80-\uDF91]|\uD803[\uDC00-\uDC48\uDC80-\uDCB2\uDCC0-\uDCF2\uDD00-\uDD23\uDD4A-\uDD65\uDD6F-\uDD85\uDE80-\uDEA9\uDEB0\uDEB1\uDEC2-\uDEC7\uDF00-\uDF1C\uDF27\uDF30-\uDF45\uDF70-\uDF81\uDFB0-\uDFC4\uDFE0-\uDFF6]|\uD804[\uDC03-\uDC37\uDC71\uDC72\uDC75\uDC83-\uDCAF\uDCD0-\uDCE8\uDD03-\uDD26\uDD44\uDD47\uDD50-\uDD72\uDD76\uDD83-\uDDB2\uDDC1-\uDDC4\uDDDA\uDDDC\uDE00-\uDE11\uDE13-\uDE2B\uDE3F\uDE40\uDE80-\uDE86\uDE88\uDE8A-\uDE8D\uDE8F-\uDE9D\uDE9F-\uDEA8\uDEB0-\uDEDE\uDF05-\uDF0C\uDF0F\uDF10\uDF13-\uDF28\uDF2A-\uDF30\uDF32\uDF33\uDF35-\uDF39\uDF3D\uDF50\uDF5D-\uDF61\uDF80-\uDF89\uDF8B\uDF8E\uDF90-\uDFB5\uDFB7\uDFD1\uDFD3]|\uD805[\uDC00-\uDC34\uDC47-\uDC4A\uDC5F-\uDC61\uDC80-\uDCAF\uDCC4\uDCC5\uDCC7\uDD80-\uDDAE\uDDD8-\uDDDB\uDE00-\uDE2F\uDE44\uDE80-\uDEAA\uDEB8\uDF00-\uDF1A\uDF40-\uDF46]|\uD806[\uDC00-\uDC2B\uDCA0-\uDCDF\uDCFF-\uDD06\uDD09\uDD0C-\uDD13\uDD15\uDD16\uDD18-\uDD2F\uDD3F\uDD41\uDDA0-\uDDA7\uDDAA-\uDDD0\uDDE1\uDDE3\uDE00\uDE0B-\uDE32\uDE3A\uDE50\uDE5C-\uDE89\uDE9D\uDEB0-\uDEF8\uDFC0-\uDFE0]|\uD807[\uDC00-\uDC08\uDC0A-\uDC2E\uDC40\uDC72-\uDC8F\uDD00-\uDD06\uDD08\uDD09\uDD0B-\uDD30\uDD46\uDD60-\uDD65\uDD67\uDD68\uDD6A-\uDD89\uDD98\uDDB0-\uDDDB\uDEE0-\uDEF2\uDF02\uDF04-\uDF10\uDF12-\uDF33\uDFB0]|\uD808[\uDC00-\uDF99]|\uD809[\uDC00-\uDC6E\uDC80-\uDD43]|\uD80B[\uDF90-\uDFF0]|[\uD80C\uD80E\uD80F\uD81C-\uD822\uD840-\uD868\uD86A-\uD86D\uD86F-\uD872\uD874-\uD879\uD880-\uD883\uD885-\uD88C][\uDC00-\uDFFF]|\uD80D[\uDC00-\uDC2F\uDC41-\uDC46\uDC60-\uDFFF]|\uD810[\uDC00-\uDFFA]|\uD811[\uDC00-\uDE46]|\uD818[\uDD00-\uDD1D]|\uD81A[\uDC00-\uDE38\uDE40-\uDE5E\uDE70-\uDEBE\uDED0-\uDEED\uDF00-\uDF2F\uDF40-\uDF43\uDF63-\uDF77\uDF7D-\uDF8F]|\uD81B[\uDD40-\uDD6C\uDE40-\uDE7F\uDEA0-\uDEB8\uDEBB-\uDED3\uDF00-\uDF4A\uDF50\uDF93-\uDF9F\uDFE0\uDFE1\uDFE3\uDFF2-\uDFF6]|\uD823[\uDC00-\uDCD5\uDCFF-\uDD1E\uDD80-\uDDF2]|\uD82B[\uDFF0-\uDFF3\uDFF5-\uDFFB\uDFFD\uDFFE]|\uD82C[\uDC00-\uDD22\uDD32\uDD50-\uDD52\uDD55\uDD64-\uDD67\uDD70-\uDEFB]|\uD82F[\uDC00-\uDC6A\uDC70-\uDC7C\uDC80-\uDC88\uDC90-\uDC99]|\uD835[\uDC00-\uDC54\uDC56-\uDC9C\uDC9E\uDC9F\uDCA2\uDCA5\uDCA6\uDCA9-\uDCAC\uDCAE-\uDCB9\uDCBB\uDCBD-\uDCC3\uDCC5-\uDD05\uDD07-\uDD0A\uDD0D-\uDD14\uDD16-\uDD1C\uDD1E-\uDD39\uDD3B-\uDD3E\uDD40-\uDD44\uDD46\uDD4A-\uDD50\uDD52-\uDEA5\uDEA8-\uDEC0\uDEC2-\uDEDA\uDEDC-\uDEFA\uDEFC-\uDF14\uDF16-\uDF34\uDF36-\uDF4E\uDF50-\uDF6E\uDF70-\uDF88\uDF8A-\uDFA8\uDFAA-\uDFC2\uDFC4-\uDFCB]|\uD837[\uDF00-\uDF1E\uDF25-\uDF2A]|\uD838[\uDC30-\uDC6D\uDD00-\uDD2C\uDD37-\uDD3D\uDD4E\uDE90-\uDEAD\uDEC0-\uDEEB]|\uD839[\uDCD0-\uDCEB\uDDD0-\uDDED\uDDF0\uDEC0-\uDEDE\uDEE0-\uDEE2\uDEE4\uDEE5\uDEE7-\uDEED\uDEF0-\uDEF4\uDEFE\uDEFF\uDFE0-\uDFE6\uDFE8-\uDFEB\uDFED\uDFEE\uDFF0-\uDFFE]|\uD83A[\uDC00-\uDCC4\uDD00-\uDD43\uDD4B]|\uD83B[\uDE00-\uDE03\uDE05-\uDE1F\uDE21\uDE22\uDE24\uDE27\uDE29-\uDE32\uDE34-\uDE37\uDE39\uDE3B\uDE42\uDE47\uDE49\uDE4B\uDE4D-\uDE4F\uDE51\uDE52\uDE54\uDE57\uDE59\uDE5B\uDE5D\uDE5F\uDE61\uDE62\uDE64\uDE67-\uDE6A\uDE6C-\uDE72\uDE74-\uDE77\uDE79-\uDE7C\uDE7E\uDE80-\uDE89\uDE8B-\uDE9B\uDEA1-\uDEA3\uDEA5-\uDEA9\uDEAB-\uDEBB]|\uD869[\uDC00-\uDEDF\uDF00-\uDFFF]|\uD86E[\uDC00-\uDC1D\uDC20-\uDFFF]|\uD873[\uDC00-\uDEAD\uDEB0-\uDFFF]|\uD87A[\uDC00-\uDFE0\uDFF0-\uDFFF]|\uD87B[\uDC00-\uDE5D]|\uD87E[\uDC00-\uDE1D]|\uD884[\uDC00-\uDF4A\uDF50-\uDFFF]|\uD88D[\uDC00-\uDC79])$/;
  const ID_CONTINUE = /^(?:[\$0-9A-Z_a-z\xAA\xB5\xB7\xBA\xC0-\xD6\xD8-\xF6\xF8-\u02C1\u02C6-\u02D1\u02E0-\u02E4\u02EC\u02EE\u0300-\u0374\u0376\u0377\u037A-\u037D\u037F\u0386-\u038A\u038C\u038E-\u03A1\u03A3-\u03F5\u03F7-\u0481\u0483-\u0487\u048A-\u052F\u0531-\u0556\u0559\u0560-\u0588\u0591-\u05BD\u05BF\u05C1\u05C2\u05C4\u05C5\u05C7\u05D0-\u05EA\u05EF-\u05F2\u0610-\u061A\u0620-\u0669\u066E-\u06D3\u06D5-\u06DC\u06DF-\u06E8\u06EA-\u06FC\u06FF\u0710-\u074A\u074D-\u07B1\u07C0-\u07F5\u07FA\u07FD\u0800-\u082D\u0840-\u085B\u0860-\u086A\u0870-\u0887\u0889-\u088F\u0897-\u08E1\u08E3-\u0963\u0966-\u096F\u0971-\u0983\u0985-\u098C\u098F\u0990\u0993-\u09A8\u09AA-\u09B0\u09B2\u09B6-\u09B9\u09BC-\u09C4\u09C7\u09C8\u09CB-\u09CE\u09D7\u09DC\u09DD\u09DF-\u09E3\u09E6-\u09F1\u09FC\u09FE\u0A01-\u0A03\u0A05-\u0A0A\u0A0F\u0A10\u0A13-\u0A28\u0A2A-\u0A30\u0A32\u0A33\u0A35\u0A36\u0A38\u0A39\u0A3C\u0A3E-\u0A42\u0A47\u0A48\u0A4B-\u0A4D\u0A51\u0A59-\u0A5C\u0A5E\u0A66-\u0A75\u0A81-\u0A83\u0A85-\u0A8D\u0A8F-\u0A91\u0A93-\u0AA8\u0AAA-\u0AB0\u0AB2\u0AB3\u0AB5-\u0AB9\u0ABC-\u0AC5\u0AC7-\u0AC9\u0ACB-\u0ACD\u0AD0\u0AE0-\u0AE3\u0AE6-\u0AEF\u0AF9-\u0AFF\u0B01-\u0B03\u0B05-\u0B0C\u0B0F\u0B10\u0B13-\u0B28\u0B2A-\u0B30\u0B32\u0B33\u0B35-\u0B39\u0B3C-\u0B44\u0B47\u0B48\u0B4B-\u0B4D\u0B55-\u0B57\u0B5C\u0B5D\u0B5F-\u0B63\u0B66-\u0B6F\u0B71\u0B82\u0B83\u0B85-\u0B8A\u0B8E-\u0B90\u0B92-\u0B95\u0B99\u0B9A\u0B9C\u0B9E\u0B9F\u0BA3\u0BA4\u0BA8-\u0BAA\u0BAE-\u0BB9\u0BBE-\u0BC2\u0BC6-\u0BC8\u0BCA-\u0BCD\u0BD0\u0BD7\u0BE6-\u0BEF\u0C00-\u0C0C\u0C0E-\u0C10\u0C12-\u0C28\u0C2A-\u0C39\u0C3C-\u0C44\u0C46-\u0C48\u0C4A-\u0C4D\u0C55\u0C56\u0C58-\u0C5A\u0C5C\u0C5D\u0C60-\u0C63\u0C66-\u0C6F\u0C80-\u0C83\u0C85-\u0C8C\u0C8E-\u0C90\u0C92-\u0CA8\u0CAA-\u0CB3\u0CB5-\u0CB9\u0CBC-\u0CC4\u0CC6-\u0CC8\u0CCA-\u0CCD\u0CD5\u0CD6\u0CDC-\u0CDE\u0CE0-\u0CE3\u0CE6-\u0CEF\u0CF1-\u0CF3\u0D00-\u0D0C\u0D0E-\u0D10\u0D12-\u0D44\u0D46-\u0D48\u0D4A-\u0D4E\u0D54-\u0D57\u0D5F-\u0D63\u0D66-\u0D6F\u0D7A-\u0D7F\u0D81-\u0D83\u0D85-\u0D96\u0D9A-\u0DB1\u0DB3-\u0DBB\u0DBD\u0DC0-\u0DC6\u0DCA\u0DCF-\u0DD4\u0DD6\u0DD8-\u0DDF\u0DE6-\u0DEF\u0DF2\u0DF3\u0E01-\u0E3A\u0E40-\u0E4E\u0E50-\u0E59\u0E81\u0E82\u0E84\u0E86-\u0E8A\u0E8C-\u0EA3\u0EA5\u0EA7-\u0EBD\u0EC0-\u0EC4\u0EC6\u0EC8-\u0ECE\u0ED0-\u0ED9\u0EDC-\u0EDF\u0F00\u0F18\u0F19\u0F20-\u0F29\u0F35\u0F37\u0F39\u0F3E-\u0F47\u0F49-\u0F6C\u0F71-\u0F84\u0F86-\u0F97\u0F99-\u0FBC\u0FC6\u1000-\u1049\u1050-\u109D\u10A0-\u10C5\u10C7\u10CD\u10D0-\u10FA\u10FC-\u1248\u124A-\u124D\u1250-\u1256\u1258\u125A-\u125D\u1260-\u1288\u128A-\u128D\u1290-\u12B0\u12B2-\u12B5\u12B8-\u12BE\u12C0\u12C2-\u12C5\u12C8-\u12D6\u12D8-\u1310\u1312-\u1315\u1318-\u135A\u135D-\u135F\u1369-\u1371\u1380-\u138F\u13A0-\u13F5\u13F8-\u13FD\u1401-\u166C\u166F-\u167F\u1681-\u169A\u16A0-\u16EA\u16EE-\u16F8\u1700-\u1715\u171F-\u1734\u1740-\u1753\u1760-\u176C\u176E-\u1770\u1772\u1773\u1780-\u17D3\u17D7\u17DC\u17DD\u17E0-\u17E9\u180B-\u180D\u180F-\u1819\u1820-\u1878\u1880-\u18AA\u18B0-\u18F5\u1900-\u191E\u1920-\u192B\u1930-\u193B\u1946-\u196D\u1970-\u1974\u1980-\u19AB\u19B0-\u19C9\u19D0-\u19DA\u1A00-\u1A1B\u1A20-\u1A5E\u1A60-\u1A7C\u1A7F-\u1A89\u1A90-\u1A99\u1AA7\u1AB0-\u1ABD\u1ABF-\u1ADD\u1AE0-\u1AEB\u1B00-\u1B4C\u1B50-\u1B59\u1B6B-\u1B73\u1B80-\u1BF3\u1C00-\u1C37\u1C40-\u1C49\u1C4D-\u1C7D\u1C80-\u1C8A\u1C90-\u1CBA\u1CBD-\u1CBF\u1CD0-\u1CD2\u1CD4-\u1CFA\u1D00-\u1F15\u1F18-\u1F1D\u1F20-\u1F45\u1F48-\u1F4D\u1F50-\u1F57\u1F59\u1F5B\u1F5D\u1F5F-\u1F7D\u1F80-\u1FB4\u1FB6-\u1FBC\u1FBE\u1FC2-\u1FC4\u1FC6-\u1FCC\u1FD0-\u1FD3\u1FD6-\u1FDB\u1FE0-\u1FEC\u1FF2-\u1FF4\u1FF6-\u1FFC\u200C\u200D\u203F\u2040\u2054\u2071\u207F\u2090-\u209C\u20D0-\u20DC\u20E1\u20E5-\u20F0\u2102\u2107\u210A-\u2113\u2115\u2118-\u211D\u2124\u2126\u2128\u212A-\u2139\u213C-\u213F\u2145-\u2149\u214E\u2160-\u2188\u2C00-\u2CE4\u2CEB-\u2CF3\u2D00-\u2D25\u2D27\u2D2D\u2D30-\u2D67\u2D6F\u2D7F-\u2D96\u2DA0-\u2DA6\u2DA8-\u2DAE\u2DB0-\u2DB6\u2DB8-\u2DBE\u2DC0-\u2DC6\u2DC8-\u2DCE\u2DD0-\u2DD6\u2DD8-\u2DDE\u2DE0-\u2DFF\u3005-\u3007\u3021-\u302F\u3031-\u3035\u3038-\u303C\u3041-\u3096\u3099-\u309F\u30A1-\u30FF\u3105-\u312F\u3131-\u318E\u31A0-\u31BF\u31F0-\u31FF\u3400-\u4DBF\u4E00-\uA48C\uA4D0-\uA4FD\uA500-\uA60C\uA610-\uA62B\uA640-\uA66F\uA674-\uA67D\uA67F-\uA6F1\uA717-\uA71F\uA722-\uA788\uA78B-\uA7DC\uA7F1-\uA827\uA82C\uA840-\uA873\uA880-\uA8C5\uA8D0-\uA8D9\uA8E0-\uA8F7\uA8FB\uA8FD-\uA92D\uA930-\uA953\uA960-\uA97C\uA980-\uA9C0\uA9CF-\uA9D9\uA9E0-\uA9FE\uAA00-\uAA36\uAA40-\uAA4D\uAA50-\uAA59\uAA60-\uAA76\uAA7A-\uAAC2\uAADB-\uAADD\uAAE0-\uAAEF\uAAF2-\uAAF6\uAB01-\uAB06\uAB09-\uAB0E\uAB11-\uAB16\uAB20-\uAB26\uAB28-\uAB2E\uAB30-\uAB5A\uAB5C-\uAB69\uAB70-\uABEA\uABEC\uABED\uABF0-\uABF9\uAC00-\uD7A3\uD7B0-\uD7C6\uD7CB-\uD7FB\uF900-\uFA6D\uFA70-\uFAD9\uFB00-\uFB06\uFB13-\uFB17\uFB1D-\uFB28\uFB2A-\uFB36\uFB38-\uFB3C\uFB3E\uFB40\uFB41\uFB43\uFB44\uFB46-\uFBB1\uFBD3-\uFD3D\uFD50-\uFD8F\uFD92-\uFDC7\uFDF0-\uFDFB\uFE00-\uFE0F\uFE20-\uFE2F\uFE33\uFE34\uFE4D-\uFE4F\uFE70-\uFE74\uFE76-\uFEFC\uFF10-\uFF19\uFF21-\uFF3A\uFF3F\uFF41-\uFF5A\uFF65-\uFFBE\uFFC2-\uFFC7\uFFCA-\uFFCF\uFFD2-\uFFD7\uFFDA-\uFFDC]|\uD800[\uDC00-\uDC0B\uDC0D-\uDC26\uDC28-\uDC3A\uDC3C\uDC3D\uDC3F-\uDC4D\uDC50-\uDC5D\uDC80-\uDCFA\uDD40-\uDD74\uDDFD\uDE80-\uDE9C\uDEA0-\uDED0\uDEE0\uDF00-\uDF1F\uDF2D-\uDF4A\uDF50-\uDF7A\uDF80-\uDF9D\uDFA0-\uDFC3\uDFC8-\uDFCF\uDFD1-\uDFD5]|\uD801[\uDC00-\uDC9D\uDCA0-\uDCA9\uDCB0-\uDCD3\uDCD8-\uDCFB\uDD00-\uDD27\uDD30-\uDD63\uDD70-\uDD7A\uDD7C-\uDD8A\uDD8C-\uDD92\uDD94\uDD95\uDD97-\uDDA1\uDDA3-\uDDB1\uDDB3-\uDDB9\uDDBB\uDDBC\uDDC0-\uDDF3\uDE00-\uDF36\uDF40-\uDF55\uDF60-\uDF67\uDF80-\uDF85\uDF87-\uDFB0\uDFB2-\uDFBA]|\uD802[\uDC00-\uDC05\uDC08\uDC0A-\uDC35\uDC37\uDC38\uDC3C\uDC3F-\uDC55\uDC60-\uDC76\uDC80-\uDC9E\uDCE0-\uDCF2\uDCF4\uDCF5\uDD00-\uDD15\uDD20-\uDD39\uDD40-\uDD59\uDD80-\uDDB7\uDDBE\uDDBF\uDE00-\uDE03\uDE05\uDE06\uDE0C-\uDE13\uDE15-\uDE17\uDE19-\uDE35\uDE38-\uDE3A\uDE3F\uDE60-\uDE7C\uDE80-\uDE9C\uDEC0-\uDEC7\uDEC9-\uDEE6\uDF00-\uDF35\uDF40-\uDF55\uDF60-\uDF72\uDF80-\uDF91]|\uD803[\uDC00-\uDC48\uDC80-\uDCB2\uDCC0-\uDCF2\uDD00-\uDD27\uDD30-\uDD39\uDD40-\uDD65\uDD69-\uDD6D\uDD6F-\uDD85\uDE80-\uDEA9\uDEAB\uDEAC\uDEB0\uDEB1\uDEC2-\uDEC7\uDEFA-\uDF1C\uDF27\uDF30-\uDF50\uDF70-\uDF85\uDFB0-\uDFC4\uDFE0-\uDFF6]|\uD804[\uDC00-\uDC46\uDC66-\uDC75\uDC7F-\uDCBA\uDCC2\uDCD0-\uDCE8\uDCF0-\uDCF9\uDD00-\uDD34\uDD36-\uDD3F\uDD44-\uDD47\uDD50-\uDD73\uDD76\uDD80-\uDDC4\uDDC9-\uDDCC\uDDCE-\uDDDA\uDDDC\uDE00-\uDE11\uDE13-\uDE37\uDE3E-\uDE41\uDE80-\uDE86\uDE88\uDE8A-\uDE8D\uDE8F-\uDE9D\uDE9F-\uDEA8\uDEB0-\uDEEA\uDEF0-\uDEF9\uDF00-\uDF03\uDF05-\uDF0C\uDF0F\uDF10\uDF13-\uDF28\uDF2A-\uDF30\uDF32\uDF33\uDF35-\uDF39\uDF3B-\uDF44\uDF47\uDF48\uDF4B-\uDF4D\uDF50\uDF57\uDF5D-\uDF63\uDF66-\uDF6C\uDF70-\uDF74\uDF80-\uDF89\uDF8B\uDF8E\uDF90-\uDFB5\uDFB7-\uDFC0\uDFC2\uDFC5\uDFC7-\uDFCA\uDFCC-\uDFD3\uDFE1\uDFE2]|\uD805[\uDC00-\uDC4A\uDC50-\uDC59\uDC5E-\uDC61\uDC80-\uDCC5\uDCC7\uDCD0-\uDCD9\uDD80-\uDDB5\uDDB8-\uDDC0\uDDD8-\uDDDD\uDE00-\uDE40\uDE44\uDE50-\uDE59\uDE80-\uDEB8\uDEC0-\uDEC9\uDED0-\uDEE3\uDF00-\uDF1A\uDF1D-\uDF2B\uDF30-\uDF39\uDF40-\uDF46]|\uD806[\uDC00-\uDC3A\uDCA0-\uDCE9\uDCFF-\uDD06\uDD09\uDD0C-\uDD13\uDD15\uDD16\uDD18-\uDD35\uDD37\uDD38\uDD3B-\uDD43\uDD50-\uDD59\uDDA0-\uDDA7\uDDAA-\uDDD7\uDDDA-\uDDE1\uDDE3\uDDE4\uDE00-\uDE3E\uDE47\uDE50-\uDE99\uDE9D\uDEB0-\uDEF8\uDF60-\uDF67\uDFC0-\uDFE0\uDFF0-\uDFF9]|\uD807[\uDC00-\uDC08\uDC0A-\uDC36\uDC38-\uDC40\uDC50-\uDC59\uDC72-\uDC8F\uDC92-\uDCA7\uDCA9-\uDCB6\uDD00-\uDD06\uDD08\uDD09\uDD0B-\uDD36\uDD3A\uDD3C\uDD3D\uDD3F-\uDD47\uDD50-\uDD59\uDD60-\uDD65\uDD67\uDD68\uDD6A-\uDD8E\uDD90\uDD91\uDD93-\uDD98\uDDA0-\uDDA9\uDDB0-\uDDDB\uDDE0-\uDDE9\uDEE0-\uDEF6\uDF00-\uDF10\uDF12-\uDF3A\uDF3E-\uDF42\uDF50-\uDF5A\uDFB0]|\uD808[\uDC00-\uDF99]|\uD809[\uDC00-\uDC6E\uDC80-\uDD43]|\uD80B[\uDF90-\uDFF0]|[\uD80C\uD80E\uD80F\uD81C-\uD822\uD840-\uD868\uD86A-\uD86D\uD86F-\uD872\uD874-\uD879\uD880-\uD883\uD885-\uD88C][\uDC00-\uDFFF]|\uD80D[\uDC00-\uDC2F\uDC40-\uDC55\uDC60-\uDFFF]|\uD810[\uDC00-\uDFFA]|\uD811[\uDC00-\uDE46]|\uD818[\uDD00-\uDD39]|\uD81A[\uDC00-\uDE38\uDE40-\uDE5E\uDE60-\uDE69\uDE70-\uDEBE\uDEC0-\uDEC9\uDED0-\uDEED\uDEF0-\uDEF4\uDF00-\uDF36\uDF40-\uDF43\uDF50-\uDF59\uDF63-\uDF77\uDF7D-\uDF8F]|\uD81B[\uDD40-\uDD6C\uDD70-\uDD79\uDE40-\uDE7F\uDEA0-\uDEB8\uDEBB-\uDED3\uDF00-\uDF4A\uDF4F-\uDF87\uDF8F-\uDF9F\uDFE0\uDFE1\uDFE3\uDFE4\uDFF0-\uDFF6]|\uD823[\uDC00-\uDCD5\uDCFF-\uDD1E\uDD80-\uDDF2]|\uD82B[\uDFF0-\uDFF3\uDFF5-\uDFFB\uDFFD\uDFFE]|\uD82C[\uDC00-\uDD22\uDD32\uDD50-\uDD52\uDD55\uDD64-\uDD67\uDD70-\uDEFB]|\uD82F[\uDC00-\uDC6A\uDC70-\uDC7C\uDC80-\uDC88\uDC90-\uDC99\uDC9D\uDC9E]|\uD833[\uDCF0-\uDCF9\uDF00-\uDF2D\uDF30-\uDF46]|\uD834[\uDD65-\uDD69\uDD6D-\uDD72\uDD7B-\uDD82\uDD85-\uDD8B\uDDAA-\uDDAD\uDE42-\uDE44]|\uD835[\uDC00-\uDC54\uDC56-\uDC9C\uDC9E\uDC9F\uDCA2\uDCA5\uDCA6\uDCA9-\uDCAC\uDCAE-\uDCB9\uDCBB\uDCBD-\uDCC3\uDCC5-\uDD05\uDD07-\uDD0A\uDD0D-\uDD14\uDD16-\uDD1C\uDD1E-\uDD39\uDD3B-\uDD3E\uDD40-\uDD44\uDD46\uDD4A-\uDD50\uDD52-\uDEA5\uDEA8-\uDEC0\uDEC2-\uDEDA\uDEDC-\uDEFA\uDEFC-\uDF14\uDF16-\uDF34\uDF36-\uDF4E\uDF50-\uDF6E\uDF70-\uDF88\uDF8A-\uDFA8\uDFAA-\uDFC2\uDFC4-\uDFCB\uDFCE-\uDFFF]|\uD836[\uDE00-\uDE36\uDE3B-\uDE6C\uDE75\uDE84\uDE9B-\uDE9F\uDEA1-\uDEAF]|\uD837[\uDF00-\uDF1E\uDF25-\uDF2A]|\uD838[\uDC00-\uDC06\uDC08-\uDC18\uDC1B-\uDC21\uDC23\uDC24\uDC26-\uDC2A\uDC30-\uDC6D\uDC8F\uDD00-\uDD2C\uDD30-\uDD3D\uDD40-\uDD49\uDD4E\uDE90-\uDEAE\uDEC0-\uDEF9]|\uD839[\uDCD0-\uDCF9\uDDD0-\uDDFA\uDEC0-\uDEDE\uDEE0-\uDEF5\uDEFE\uDEFF\uDFE0-\uDFE6\uDFE8-\uDFEB\uDFED\uDFEE\uDFF0-\uDFFE]|\uD83A[\uDC00-\uDCC4\uDCD0-\uDCD6\uDD00-\uDD4B\uDD50-\uDD59]|\uD83B[\uDE00-\uDE03\uDE05-\uDE1F\uDE21\uDE22\uDE24\uDE27\uDE29-\uDE32\uDE34-\uDE37\uDE39\uDE3B\uDE42\uDE47\uDE49\uDE4B\uDE4D-\uDE4F\uDE51\uDE52\uDE54\uDE57\uDE59\uDE5B\uDE5D\uDE5F\uDE61\uDE62\uDE64\uDE67-\uDE6A\uDE6C-\uDE72\uDE74-\uDE77\uDE79-\uDE7C\uDE7E\uDE80-\uDE89\uDE8B-\uDE9B\uDEA1-\uDEA3\uDEA5-\uDEA9\uDEAB-\uDEBB]|\uD83E[\uDFF0-\uDFF9]|\uD869[\uDC00-\uDEDF\uDF00-\uDFFF]|\uD86E[\uDC00-\uDC1D\uDC20-\uDFFF]|\uD873[\uDC00-\uDEAD\uDEB0-\uDFFF]|\uD87A[\uDC00-\uDFE0\uDFF0-\uDFFF]|\uD87B[\uDC00-\uDE5D]|\uD87E[\uDC00-\uDE1D]|\uD884[\uDC00-\uDF4A\uDF50-\uDFFF]|\uD88D[\uDC00-\uDC79]|\uDB40[\uDD00-\uDDEF])$/;
  const DEBUG_URL = "https://git.new/pathToRegexpError";
  const SIMPLE_TOKENS = {
    // Groups.
    "{": "{",
    "}": "}",
    // Reserved.
    "(": "(",
    ")": ")",
    "[": "[",
    "]": "]",
    "+": "+",
    "?": "?",
    "!": "!"
  };
  /**
   * Escape text for stringify to path.
   */
  function escapeText(str) {
    return str.replace(/[{}()\[\]+?!:*]/g, "\\$&");
  }
  /**
   * Escape a regular expression string.
   */
  function escape(str) {
    return str.replace(/[.+*?^${}()[\]|/\\]/g, "\\$&");
  }
  /**
   * Tokenize input string.
   */
  function* lexer(str) {
    const chars = [...str];
    let i = 0;
    function name() {
      let value = "";
      if (ID_START.test(chars[++i])) {
        value += chars[i];
        while (ID_CONTINUE.test(chars[++i])) {
          value += chars[i];
        }
      } else if (chars[i] === '"') {
        let pos = i;
        while (i < chars.length) {
          if (chars[++i] === '"') {
            i++;
            pos = 0;
            break;
          }
          if (chars[i] === "\\") {
            value += chars[++i];
          } else {
            value += chars[i];
          }
        }
        if (pos) {
          throw new TypeError("Unterminated quote at ".concat(pos, ": ").concat(DEBUG_URL));
        }
      }
      if (!value) {
        throw new TypeError("Missing parameter name at ".concat(i, ": ").concat(DEBUG_URL));
      }
      return value;
    }
    while (i < chars.length) {
      const value = chars[i];
      const type = SIMPLE_TOKENS[value];
      if (type) {
        yield {
          type,
          index: i++,
          value
        };
      } else if (value === "\\") {
        yield {
          type: "ESCAPED",
          index: i++,
          value: chars[i++]
        };
      } else if (value === ":") {
        const value = name();
        yield {
          type: "PARAM",
          index: i,
          value
        };
      } else if (value === "*") {
        const value = name();
        yield {
          type: "WILDCARD",
          index: i,
          value
        };
      } else {
        yield {
          type: "CHAR",
          index: i,
          value: chars[i++]
        };
      }
    }
    return {
      type: "END",
      index: i,
      value: ""
    };
  }
  class Iter {
    constructor(tokens) {
      this.tokens = tokens;
    }
    peek() {
      if (!this._peek) {
        const next = this.tokens.next();
        this._peek = next.value;
      }
      return this._peek;
    }
    tryConsume(type) {
      const token = this.peek();
      if (token.type !== type) return;
      this._peek = undefined; // Reset after consumed.
      return token.value;
    }
    consume(type) {
      const value = this.tryConsume(type);
      if (value !== undefined) return value;
      const {
        type: nextType,
        index
      } = this.peek();
      throw new TypeError("Unexpected ".concat(nextType, " at ").concat(index, ", expected ").concat(type, ": ").concat(DEBUG_URL));
    }
    text() {
      let result = "";
      let value;
      while (value = this.tryConsume("CHAR") || this.tryConsume("ESCAPED")) {
        result += value;
      }
      return result;
    }
  }
  /**
   * Tokenized path instance.
   */
  class TokenData {
    constructor(tokens) {
      this.tokens = tokens;
    }
  }
  dist.TokenData = TokenData;
  /**
   * Parse a string for the raw tokens.
   */
  function parse$1(str) {
    let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
    const {
      encodePath = NOOP_VALUE
    } = options;
    const it = new Iter(lexer(str));
    function consume(endType) {
      const tokens = [];
      while (true) {
        const path = it.text();
        if (path) tokens.push({
          type: "text",
          value: encodePath(path)
        });
        const param = it.tryConsume("PARAM");
        if (param) {
          tokens.push({
            type: "param",
            name: param
          });
          continue;
        }
        const wildcard = it.tryConsume("WILDCARD");
        if (wildcard) {
          tokens.push({
            type: "wildcard",
            name: wildcard
          });
          continue;
        }
        const open = it.tryConsume("{");
        if (open) {
          tokens.push({
            type: "group",
            tokens: consume("}")
          });
          continue;
        }
        it.consume(endType);
        return tokens;
      }
    }
    const tokens = consume("END");
    return new TokenData(tokens);
  }
  /**
   * Compile a string to a template function for the path.
   */
  function compile$1(path) {
    let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
    const {
      encode = encodeURIComponent,
      delimiter = DEFAULT_DELIMITER
    } = options;
    const data = path instanceof TokenData ? path : parse$1(path, options);
    const fn = tokensToFunction(data.tokens, delimiter, encode);
    return function path() {
      let data = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
      const [path, ...missing] = fn(data);
      if (missing.length) {
        throw new TypeError("Missing parameters: ".concat(missing.join(", ")));
      }
      return path;
    };
  }
  function tokensToFunction(tokens, delimiter, encode) {
    const encoders = tokens.map(token => tokenToFunction(token, delimiter, encode));
    return data => {
      const result = [""];
      for (const encoder of encoders) {
        const [value, ...extras] = encoder(data);
        result[0] += value;
        result.push(...extras);
      }
      return result;
    };
  }
  /**
   * Convert a single token into a path building function.
   */
  function tokenToFunction(token, delimiter, encode) {
    if (token.type === "text") return () => [token.value];
    if (token.type === "group") {
      const fn = tokensToFunction(token.tokens, delimiter, encode);
      return data => {
        const [value, ...missing] = fn(data);
        if (!missing.length) return [value];
        return [""];
      };
    }
    const encodeValue = encode || NOOP_VALUE;
    if (token.type === "wildcard" && encode !== false) {
      return data => {
        const value = data[token.name];
        if (value == null) return ["", token.name];
        if (!Array.isArray(value) || value.length === 0) {
          throw new TypeError("Expected \"".concat(token.name, "\" to be a non-empty array"));
        }
        return [value.map((value, index) => {
          if (typeof value !== "string") {
            throw new TypeError("Expected \"".concat(token.name, "/").concat(index, "\" to be a string"));
          }
          return encodeValue(value);
        }).join(delimiter)];
      };
    }
    return data => {
      const value = data[token.name];
      if (value == null) return ["", token.name];
      if (typeof value !== "string") {
        throw new TypeError("Expected \"".concat(token.name, "\" to be a string"));
      }
      return [encodeValue(value)];
    };
  }
  /**
   * Transform a path into a match function.
   */
  function match(path) {
    let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
    const {
      decode = decodeURIComponent,
      delimiter = DEFAULT_DELIMITER
    } = options;
    const {
      regexp,
      keys
    } = pathToRegexp$1(path, options);
    const decoders = keys.map(key => {
      if (decode === false) return NOOP_VALUE;
      if (key.type === "param") return decode;
      return value => value.split(delimiter).map(decode);
    });
    return function match(input) {
      const m = regexp.exec(input);
      if (!m) return false;
      const path = m[0];
      const params = Object.create(null);
      for (let i = 1; i < m.length; i++) {
        if (m[i] === undefined) continue;
        const key = keys[i - 1];
        const decoder = decoders[i - 1];
        params[key.name] = decoder(m[i]);
      }
      return {
        path,
        params
      };
    };
  }
  function pathToRegexp$1(path) {
    let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
    const {
      delimiter = DEFAULT_DELIMITER,
      end = true,
      sensitive = false,
      trailing = true
    } = options;
    const keys = [];
    const sources = [];
    const flags = sensitive ? "" : "i";
    const paths = Array.isArray(path) ? path : [path];
    const items = paths.map(path => path instanceof TokenData ? path : parse$1(path, options));
    for (const {
      tokens
    } of items) {
      for (const seq of flatten(tokens, 0, [])) {
        const regexp = sequenceToRegExp(seq, delimiter, keys);
        sources.push(regexp);
      }
    }
    let pattern = "^(?:".concat(sources.join("|"), ")");
    if (trailing) pattern += "(?:".concat(escape(delimiter), "$)?");
    pattern += end ? "$" : "(?=".concat(escape(delimiter), "|$)");
    const regexp = new RegExp(pattern, flags);
    return {
      regexp,
      keys
    };
  }
  /**
   * Generate a flat list of sequence tokens from the given tokens.
   */
  function* flatten(tokens, index, init) {
    if (index === tokens.length) {
      return yield init;
    }
    const token = tokens[index];
    if (token.type === "group") {
      const fork = init.slice();
      for (const seq of flatten(token.tokens, 0, fork)) {
        yield* flatten(tokens, index + 1, seq);
      }
    } else {
      init.push(token);
    }
    yield* flatten(tokens, index + 1, init);
  }
  /**
   * Transform a flat sequence of tokens into a regular expression.
   */
  function sequenceToRegExp(tokens, delimiter, keys) {
    let result = "";
    let backtrack = "";
    let isSafeSegmentParam = true;
    for (let i = 0; i < tokens.length; i++) {
      const token = tokens[i];
      if (token.type === "text") {
        result += escape(token.value);
        backtrack += token.value;
        isSafeSegmentParam || (isSafeSegmentParam = token.value.includes(delimiter));
        continue;
      }
      if (token.type === "param" || token.type === "wildcard") {
        if (!isSafeSegmentParam && !backtrack) {
          throw new TypeError("Missing text after \"".concat(token.name, "\": ").concat(DEBUG_URL));
        }
        if (token.type === "param") {
          result += "(".concat(negate(delimiter, isSafeSegmentParam ? "" : backtrack), "+)");
        } else {
          result += "([\\s\\S]+)";
        }
        keys.push(token);
        backtrack = "";
        isSafeSegmentParam = false;
        continue;
      }
    }
    return result;
  }
  function negate(delimiter, backtrack) {
    if (backtrack.length < 2) {
      if (delimiter.length < 2) return "[^".concat(escape(delimiter + backtrack), "]");
      return "(?:(?!".concat(escape(delimiter), ")[^").concat(escape(backtrack), "])");
    }
    if (delimiter.length < 2) {
      return "(?:(?!".concat(escape(backtrack), ")[^").concat(escape(delimiter), "])");
    }
    return "(?:(?!".concat(escape(backtrack), "|").concat(escape(delimiter), ")[\\s\\S])");
  }
  /**
   * Stringify token data into a path string.
   */
  function stringify(data) {
    return data.tokens.map(function stringifyToken(token, index, tokens) {
      if (token.type === "text") return escapeText(token.value);
      if (token.type === "group") {
        return "{".concat(token.tokens.map(stringifyToken).join(""), "}");
      }
      const isSafe = isNameSafe(token.name) && isNextNameSafe(tokens[index + 1]);
      const key = isSafe ? token.name : JSON.stringify(token.name);
      if (token.type === "param") return ":".concat(key);
      if (token.type === "wildcard") return "*".concat(key);
      throw new TypeError("Unexpected token: ".concat(token));
    }).join("");
  }
  function isNameSafe(name) {
    const [first, ...rest] = name;
    if (!ID_START.test(first)) return false;
    return rest.every(char => ID_CONTINUE.test(char));
  }
  function isNextNameSafe(token) {
    if ((token === null || token === void 0 ? void 0 : token.type) !== "text") return true;
    return !ID_CONTINUE.test(token.value[0]);
  }

  // const { parse: parseUrl, format: formatUrl } = require('node:url');

  const {
    pathToRegexp,
    compile,
    parse
  } = dist;
  var layer = class Layer {
    /**
     * Initialize a new routing Layer with given `method`, `path`, and `middleware`.
     *
     * @param {String|RegExp} path Path string or regular expression.
     * @param {Array} methods Array of HTTP verbs.
     * @param {Array} middleware Layer callback/middleware or series of.
     * @param {Object=} opts
     * @param {String=} opts.name route name
     * @param {String=} opts.sensitive case sensitive (default: false)
     * @param {String=} opts.strict require the trailing slash (default: false)
     * @param {Boolean=} opts.pathAsRegExp if true, treat `path` as a regular expression
     * @returns {Layer}
     * @private
     */
    constructor(path, methods, middleware) {
      let opts = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
      this.opts = opts;
      this.name = this.opts.name || null;
      this.methods = [];
      for (const method of methods) {
        const l = this.methods.push(method.toUpperCase());
        if (this.methods[l - 1] === 'GET') this.methods.unshift('HEAD');
      }
      this.stack = Array.isArray(middleware) ? middleware : [middleware];
      // ensure middleware is a function
      for (let i = 0; i < this.stack.length; i++) {
        const fn = this.stack[i];
        const type = typeof fn;
        if (type !== 'function') throw new Error("".concat(methods.toString(), " `").concat(this.opts.name || path, "`: `middleware` must be a function, not `").concat(type, "`"));
      }
      this.path = path;
      this.paramNames = [];
      if (this.opts.pathAsRegExp === true) {
        this.regexp = new RegExp(path);
      } else if (this.path) {
        if ('strict' in this.opts) {
          // path-to-regexp renamed strict to trailing in v8.1.0
          this.opts.trailing = this.opts.strict !== true;
        }
        const {
          regexp,
          keys
        } = pathToRegexp(this.path, this.opts);
        this.regexp = regexp;
        this.paramNames = keys;
      }
    }

    /**
     * Returns whether request `path` matches route.
     *
     * @param {String} path
     * @returns {Boolean}
     * @private
     */
    match(path) {
      return this.regexp.test(path);
    }

    /**
     * Returns map of URL parameters for given `path` and `paramNames`.
     *
     * @param {String} path
     * @param {Array.<String>} captures
     * @param {Object=} params
     * @returns {Object}
     * @private
     */
    params(path, captures) {
      let params = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
      for (let len = captures.length, i = 0; i < len; i++) {
        if (this.paramNames[i]) {
          const c = captures[i];
          if (c && c.length > 0) params[this.paramNames[i].name] = c ? safeDecodeURIComponent(c) : c;
        }
      }
      return params;
    }

    /**
     * Returns array of regexp url path captures.
     *
     * @param {String} path
     * @returns {Array.<String>}
     * @private
     */
    captures(path) {
      return this.opts.ignoreCaptures ? [] : path.match(this.regexp).slice(1);
    }

    /**
     * Generate URL for route using given `params`.
     *
     * @example
     *
     * ```javascript
     * const route = new Layer('/users/:id', ['GET'], fn);
     *
     * route.url({ id: 123 }); // => "/users/123"
     * ```
     *
     * @param {Object} params url parameters
     * @returns {String}
     * @private
     */
    // url(params, options) {
    //   let args = params;
    //   const url = this.path.replace(/\(\.\*\)/g, '');

    //   if (typeof params !== 'object') {
    //     args = Array.prototype.slice.call(arguments);
    //     if (typeof args[args.length - 1] === 'object') {
    //       options = args[args.length - 1];
    //       args = args.slice(0, -1);
    //     }
    //   }

    //   const toPath = compile(url, { encode: encodeURIComponent, ...options });
    //   let replaced;
    //   const { tokens } = parse(url);
    //   let replace = {};

    //   if (Array.isArray(args)) {
    //     for (let len = tokens.length, i = 0, j = 0; i < len; i++) {
    //       if (tokens[i].name) {
    //         replace[tokens[i].name] = args[j++];
    //       }
    //     }
    //   } else if (tokens.some((token) => token.name)) {
    //     replace = params;
    //   } else if (!options) {
    //     options = params;
    //   }

    //   for (const [key, value] of Object.entries(replace)) {
    //     replace[key] = String(value);
    //   }

    //   replaced = toPath(replace);

    //   if (options && options.query) {
    //     replaced = parseUrl(replaced);
    //     if (typeof options.query === 'string') {
    //       replaced.search = options.query;
    //     } else {
    //       replaced.search = undefined;
    //       replaced.query = options.query;
    //     }

    //     return formatUrl(replaced);
    //   }

    //   return replaced;
    // }

    /**
     * Run validations on route named parameters.
     *
     * @example
     *
     * ```javascript
     * router
     *   .param('user', function (id, ctx, next) {
     *     ctx.user = users[id];
     *     if (!ctx.user) return ctx.status = 404;
     *     next();
     *   })
     *   .get('/users/:user', function (ctx, next) {
     *     ctx.body = ctx.user;
     *   });
     * ```
     *
     * @param {String} param
     * @param {Function} middleware
     * @returns {Layer}
     * @private
     */
    param(param, fn) {
      const {
        stack
      } = this;
      const params = this.paramNames;
      const middleware = function (ctx, next) {
        return fn.call(this, ctx.params[param], ctx, next);
      };
      middleware.param = param;
      const names = params.map(function (p) {
        return p.name;
      });
      const x = names.indexOf(param);
      if (x > -1) {
        // iterate through the stack, to figure out where to place the handler fn
        stack.some((fn, i) => {
          // param handlers are always first, so when we find an fn w/o a param property, stop here
          // if the param handler at this part of the stack comes after the one we are adding, stop here
          if (!fn.param || names.indexOf(fn.param) > x) {
            // inject this param handler right before the current item
            stack.splice(i, 0, middleware);
            return true; // then break the loop
          }
        });
      }
      return this;
    }

    /**
     * Prefix route path.
     *
     * @param {String} prefix
     * @returns {Layer}
     * @private
     */
    // setPrefix(prefix) {
    //   if (this.path) {
    //     this.path =
    //       this.path !== '/' || this.opts.strict === true
    //         ? `${prefix}${this.path}`
    //         : prefix;
    //     if (this.opts.pathAsRegExp === true || prefix instanceof RegExp) {
    //       this.regexp = new RegExp(this.path);
    //     } else if (this.path) {
    //       const { regexp, keys } = pathToRegexp(this.path, this.opts);
    //       this.regexp = regexp;
    //       this.paramNames = keys;
    //     }
    //   }

    //   return this;
    // }
  };

  /**
   * Safe decodeURIComponent, won't throw any error.
   * If `decodeURIComponent` error happen, just return the original value.
   *
   * @param {String} text
   * @returns {String} URL decode original string.
   * @private
   */

  function safeDecodeURIComponent(text) {
    try {
      // TODO: take a look on `safeDecodeURIComponent` if we use it only with route params let's remove the `replace` method otherwise make it flexible.
      // @link https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/decodeURIComponent#decoding_query_parameters_from_a_url
      return decodeURIComponent(text.replace(/\+/g, ' '));
    } catch (_unused) {
      return text;
    }
  }

  // const http = require('node:http');

  // const debug = require('debug')('koa-router');

  const compose$2 = compose_1;
  // const HttpError = require('http-errors');
  const methods = methods$1;
  // const { pathToRegexp } = require('path-to-regexp');

  const Layer = layer;

  // const methods = http.METHODS.map((method) => method.toLowerCase());

  /**
   * @module koa-router
   */
  let Router$1 = class Router {
    /**
     * Create a new router.
     *
     * @example
     *
     * Basic usage:
     *
     * ```javascript
     * const Koa = require('koa');
     * const Router = require('@koa/router');
     *
     * const app = new Koa();
     * const router = new Router();
     *
     * router.get('/', (ctx, next) => {
     *   // ctx.router available
     * });
     *
     * app
     *   .use(router.routes())
     *   .use(router.allowedMethods());
     * ```
     *
     * @alias module:koa-router
     * @param {Object=} opts
     * @param {Boolean=false} opts.exclusive only run last matched route's controller when there are multiple matches
     * @param {String=} opts.prefix prefix router paths
     * @param {String|RegExp=} opts.host host for router match
     * @constructor
     */
    constructor() {
      let opts = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
      if (!(this instanceof Router)) return new Router(opts);
      this.opts = opts;
      this.methods = this.opts.methods || ['HEAD', 'OPTIONS', 'GET', 'PUT', 'PATCH', 'POST', 'DELETE'];
      this.exclusive = Boolean(this.opts.exclusive);
      this.params = {};
      this.stack = [];
      this.host = this.opts.host;
    }

    /**
     * Generate URL from url pattern and given `params`.
     *
     * @example
     *
     * ```javascript
     * const url = Router.url('/users/:id', {id: 1});
     * // => "/users/1"
     * ```
     *
     * @param {String} path url pattern
     * @param {Object} params url parameters
     * @returns {String}
     */
    // static url(path, ...args) {
    //   return Layer.prototype.url.apply({ path }, args);
    // }

    /**
     * Use given middleware.
     *
     * Middleware run in the order they are defined by `.use()`. They are invoked
     * sequentially, requests start at the first middleware and work their way
     * "down" the middleware stack.
     *
     * @example
     *
     * ```javascript
     * // session middleware will run before authorize
     * router
     *   .use(session())
     *   .use(authorize());
     *
     * // use middleware only with given path
     * router.use('/users', userAuth());
     *
     * // or with an array of paths
     * router.use(['/users', '/admin'], userAuth());
     *
     * app.use(router.routes());
     * ```
     *
     * @param {String=} path
     * @param {Function} middleware
     * @param {Function=} ...
     * @returns {Router}
     */
    // use(...middleware) {
    //   const router = this;
    //   let path;

    //   // support array of paths
    //   if (Array.isArray(middleware[0]) && typeof middleware[0][0] === 'string') {
    //     const arrPaths = middleware[0];
    //     for (const p of arrPaths) {
    //       router.use.apply(router, [p, ...middleware.slice(1)]);
    //     }

    //     return this;
    //   }

    //   const hasPath = typeof middleware[0] === 'string';
    //   if (hasPath) path = middleware.shift();

    //   for (const m of middleware) {
    //     if (m.router) {
    //       const cloneRouter = Object.assign(
    //         Object.create(Router.prototype),
    //         m.router,
    //         {
    //           stack: [...m.router.stack]
    //         }
    //       );

    //       for (let j = 0; j < cloneRouter.stack.length; j++) {
    //         const nestedLayer = cloneRouter.stack[j];
    //         const cloneLayer = Object.assign(
    //           Object.create(Layer.prototype),
    //           nestedLayer
    //         );

    //         if (path) cloneLayer.setPrefix(path);
    //         if (router.opts.prefix) cloneLayer.setPrefix(router.opts.prefix);
    //         router.stack.push(cloneLayer);
    //         cloneRouter.stack[j] = cloneLayer;
    //       }

    //       if (router.params) {
    //         const routerParams = Object.keys(router.params);
    //         for (const key of routerParams) {
    //           cloneRouter.param(key, router.params[key]);
    //         }
    //       }
    //     } else {
    //       const { keys } = pathToRegexp(router.opts.prefix || '', router.opts);
    //       const routerPrefixHasParam = Boolean(
    //         router.opts.prefix && keys.length > 0
    //       );
    //       router.register(path || '([^/]*)', [], m, {
    //         end: false,
    //         ignoreCaptures: !hasPath && !routerPrefixHasParam,
    //         pathAsRegExp: true
    //       });
    //     }
    //   }

    //   return this;
    // }

    /**
     * Set the path prefix for a Router instance that was already initialized.
     *
     * @example
     *
     * ```javascript
     * router.prefix('/things/:thing_id')
     * ```
     *
     * @param {String} prefix
     * @returns {Router}
     */
    // prefix(prefix) {
    //   prefix = prefix.replace(/\/$/, '');

    //   this.opts.prefix = prefix;

    //   for (let i = 0; i < this.stack.length; i++) {
    //     const route = this.stack[i];
    //     route.setPrefix(prefix);
    //   }

    //   return this;
    // }

    /**
     * Returns router middleware which dispatches a route matching the request.
     *
     * @returns {Function}
     */
    middleware() {
      const router = this;
      const dispatch = (ctx, next) => {
        // debug('%s %s', ctx.method, ctx.path);

        // const hostMatched = router.matchHost(ctx.host);

        // if (!hostMatched) {
        //   return next();
        // }

        const path = router.opts.routerPath || ctx.newRouterPath || ctx.path || ctx.routerPath;
        const matched = router.match(path, ctx.method);
        if (ctx.matched) {
          ctx.matched.push.apply(ctx.matched, matched.path);
        } else {
          ctx.matched = matched.path;
        }
        ctx.router = router;
        if (!matched.route) return next();
        const matchedLayers = matched.pathAndMethod;
        const mostSpecificLayer = matchedLayers[matchedLayers.length - 1];
        ctx._matchedRoute = mostSpecificLayer.path;
        if (mostSpecificLayer.name) {
          ctx._matchedRouteName = mostSpecificLayer.name;
        }
        const layerChain = (router.exclusive ? [mostSpecificLayer] : matchedLayers).reduce((memo, layer) => {
          memo.push((ctx, next) => {
            ctx.captures = layer.captures(path, ctx.captures);
            ctx.request.params = layer.params(path, ctx.captures, ctx.params);
            ctx.params = ctx.request.params;
            ctx.routerPath = layer.path;
            ctx.routerName = layer.name;
            ctx._matchedRoute = layer.path;
            if (layer.name) {
              ctx._matchedRouteName = layer.name;
            }
            return next();
          });
          return [...memo, ...layer.stack];
        }, []);
        return compose$2(layerChain)(ctx, next);
      };
      dispatch.router = this;
      return dispatch;
    }
    routes() {
      return this.middleware();
    }

    /**
     * Returns separate middleware for responding to `OPTIONS` requests with
     * an `Allow` header containing the allowed methods, as well as responding
     * with `405 Method Not Allowed` and `501 Not Implemented` as appropriate.
     *
     * @example
     *
     * ```javascript
     * const Koa = require('koa');
     * const Router = require('@koa/router');
     *
     * const app = new Koa();
     * const router = new Router();
     *
     * app.use(router.routes());
     * app.use(router.allowedMethods());
     * ```
     *
     * **Example with [Boom](https://github.com/hapijs/boom)**
     *
     * ```javascript
     * const Koa = require('koa');
     * const Router = require('@koa/router');
     * const Boom = require('boom');
     *
     * const app = new Koa();
     * const router = new Router();
     *
     * app.use(router.routes());
     * app.use(router.allowedMethods({
     *   throw: true,
     *   notImplemented: () => new Boom.notImplemented(),
     *   methodNotAllowed: () => new Boom.methodNotAllowed()
     * }));
     * ```
     *
     * @param {Object=} options
     * @param {Boolean=} options.throw throw error instead of setting status and header
     * @param {Function=} options.notImplemented throw the returned value in place of the default NotImplemented error
     * @param {Function=} options.methodNotAllowed throw the returned value in place of the default MethodNotAllowed error
     * @returns {Function}
     */
    // allowedMethods(options = {}) {
    //   const implemented = this.methods;

    //   return (ctx, next) => {
    //     return next().then(() => {
    //       const allowed = {};

    //       if (ctx.matched && (!ctx.status || ctx.status === 404)) {
    //         for (let i = 0; i < ctx.matched.length; i++) {
    //           const route = ctx.matched[i];
    //           for (let j = 0; j < route.methods.length; j++) {
    //             const method = route.methods[j];
    //             allowed[method] = method;
    //           }
    //         }

    //         const allowedArr = Object.keys(allowed);
    //         if (!implemented.includes(ctx.method)) {
    //           if (options.throw) {
    //             const notImplementedThrowable =
    //               typeof options.notImplemented === 'function'
    //                 ? options.notImplemented() // set whatever the user returns from their function
    //                 : new HttpError.NotImplemented();

    //             throw notImplementedThrowable;
    //           } else {
    //             ctx.status = 501;
    //             ctx.set('Allow', allowedArr.join(', '));
    //           }
    //         } else if (allowedArr.length > 0) {
    //           if (ctx.method === 'OPTIONS') {
    //             ctx.status = 200;
    //             ctx.body = '';
    //             ctx.set('Allow', allowedArr.join(', '));
    //           } else if (!allowed[ctx.method]) {
    //             if (options.throw) {
    //               const notAllowedThrowable =
    //                 typeof options.methodNotAllowed === 'function'
    //                   ? options.methodNotAllowed() // set whatever the user returns from their function
    //                   : new HttpError.MethodNotAllowed();

    //               throw notAllowedThrowable;
    //             } else {
    //               ctx.status = 405;
    //               ctx.set('Allow', allowedArr.join(', '));
    //             }
    //           }
    //         }
    //       }
    //     });
    //   };
    // }

    /**
     * Register route with all methods.
     *
     * @param {String} name Optional.
     * @param {String} path
     * @param {Function=} middleware You may also pass multiple middleware.
     * @param {Function} callback
     * @returns {Router}
     */
    // all(name, path, middleware) {
    //   if (typeof path === 'string' || path instanceof RegExp) {
    //     middleware = Array.prototype.slice.call(arguments, 2);
    //   } else {
    //     middleware = Array.prototype.slice.call(arguments, 1);
    //     path = name;
    //     name = null;
    //   }

    //   // Sanity check to ensure we have a viable path candidate (eg: string|regex|non-empty array)
    //   if (
    //     typeof path !== 'string' &&
    //     !(path instanceof RegExp) &&
    //     (!Array.isArray(path) || path.length === 0)
    //   )
    //     throw new Error('You have to provide a path when adding an all handler');

    //   const opts = {
    //     name,
    //     pathAsRegExp: path instanceof RegExp
    //   };

    //   this.register(path, methods, middleware, { ...this.opts, ...opts });

    //   return this;
    // }

    /**
     * Redirect `source` to `destination` URL with optional 30x status `code`.
     *
     * Both `source` and `destination` can be route names.
     *
     * ```javascript
     * router.redirect('/login', 'sign-in');
     * ```
     *
     * This is equivalent to:
     *
     * ```javascript
     * router.all('/login', ctx => {
     *   ctx.redirect('/sign-in');
     *   ctx.status = 301;
     * });
     * ```
     *
     * @param {String} source URL or route name.
     * @param {String} destination URL or route name.
     * @param {Number=} code HTTP status code (default: 301).
     * @returns {Router}
     */
    // redirect(source, destination, code) {
    //   // lookup source route by name
    //   if (typeof source === 'symbol' || source[0] !== '/') {
    //     source = this.url(source);
    //     if (source instanceof Error) throw source;
    //   }

    //   // lookup destination route by name
    //   if (
    //     typeof destination === 'symbol' ||
    //     (destination[0] !== '/' && !destination.includes('://'))
    //   ) {
    //     destination = this.url(destination);
    //     if (destination instanceof Error) throw destination;
    //   }

    //   return this.all(source, (ctx) => {
    //     ctx.redirect(destination);
    //     ctx.status = code || 301;
    //   });
    // }

    /**
     * Create and register a route.
     *
     * @param {String} path Path string.
     * @param {Array.<String>} methods Array of HTTP verbs.
     * @param {Function} middleware Multiple middleware also accepted.
     * @returns {Layer}
     * @private
     */
    register(path, methods, middleware) {
      let newOpts = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
      const router = this;
      const {
        stack
      } = this;
      const opts = _objectSpread2(_objectSpread2({}, this.opts), newOpts);
      // support array of paths
      if (Array.isArray(path)) {
        for (const curPath of path) {
          router.register.call(router, curPath, methods, middleware, opts);
        }
        return this;
      }

      // create route
      const route = new Layer(path, methods, middleware, {
        end: opts.end === false ? opts.end : true,
        name: opts.name,
        sensitive: opts.sensitive || false,
        strict: opts.strict || false,
        prefix: opts.prefix || '',
        ignoreCaptures: opts.ignoreCaptures,
        pathAsRegExp: opts.pathAsRegExp
      });

      // if parent prefix exists, add prefix to new route
      if (this.opts.prefix) {
        route.setPrefix(this.opts.prefix);
      }

      // add parameter middleware
      for (let i = 0; i < Object.keys(this.params).length; i++) {
        const param = Object.keys(this.params)[i];
        route.param(param, this.params[param]);
      }
      stack.push(route);

      // debug('defined route %s %s', route.methods, route.path);

      return route;
    }

    /**
     * Lookup route with given `name`.
     *
     * @param {String} name
     * @returns {Layer|false}
     */
    // route(name) {
    //   const routes = this.stack;

    //   for (let len = routes.length, i = 0; i < len; i++) {
    //     if (routes[i].name && routes[i].name === name) return routes[i];
    //   }

    //   return false;
    // }

    /**
     * Generate URL for route. Takes a route name and map of named `params`.
     *
     * @example
     *
     * ```javascript
     * router.get('user', '/users/:id', (ctx, next) => {
     *   // ...
     * });
     *
     * router.url('user', 3);
     * // => "/users/3"
     *
     * router.url('user', { id: 3 });
     * // => "/users/3"
     *
     * router.use((ctx, next) => {
     *   // redirect to named route
     *   ctx.redirect(ctx.router.url('sign-in'));
     * })
     *
     * router.url('user', { id: 3 }, { query: { limit: 1 } });
     * // => "/users/3?limit=1"
     *
     * router.url('user', { id: 3 }, { query: "limit=1" });
     * // => "/users/3?limit=1"
     * ```
     *
     * @param {String} name route name
     * @param {Object} params url parameters
     * @param {Object} [options] options parameter
     * @param {Object|String} [options.query] query options
     * @returns {String|Error}
     */
    // url(name, ...args) {
    //   const route = this.route(name);
    //   if (route) return route.url.apply(route, args);

    //   return new Error(`No route found for name: ${String(name)}`);
    // }

    /**
     * Match given `path` and return corresponding routes.
     *
     * @param {String} path
     * @param {String} method
     * @returns {Object.<path, pathAndMethod>} returns layers that matched path and
     * path and method.
     * @private
     */
    match(path, method) {
      const layers = this.stack;
      let layer;
      const matched = {
        path: [],
        pathAndMethod: [],
        route: false
      };
      for (let len = layers.length, i = 0; i < len; i++) {
        layer = layers[i];

        // debug('test %s %s', layer.path, layer.regexp);

        if (layer.match(path)) {
          matched.path.push(layer);
          if (layer.methods.length === 0 || layer.methods.includes(method)) {
            matched.pathAndMethod.push(layer);
            if (layer.methods.length > 0) matched.route = true;
          }
        }
      }
      return matched;
    }

    /**
     * Match given `input` to allowed host
     * @param {String} input
     * @returns {boolean}
     */
    // matchHost(input) {
    //   const { host } = this;

    //   if (!host) {
    //     return true;
    //   }

    //   if (!input) {
    //     return false;
    //   }

    //   if (typeof host === 'string') {
    //     return input === host;
    //   }

    //   if (typeof host === 'object' && host instanceof RegExp) {
    //     return host.test(input);
    //   }
    // }

    /**
     * Run middleware for named route parameters. Useful for auto-loading or
     * validation.
     *
     * @example
     *
     * ```javascript
     * router
     *   .param('user', (id, ctx, next) => {
     *     ctx.user = users[id];
     *     if (!ctx.user) return ctx.status = 404;
     *     return next();
     *   })
     *   .get('/users/:user', ctx => {
     *     ctx.body = ctx.user;
     *   })
     *   .get('/users/:user/friends', ctx => {
     *     return ctx.user.getFriends().then(function(friends) {
     *       ctx.body = friends;
     *     });
     *   })
     *   // /users/3 => {"id": 3, "name": "Alex"}
     *   // /users/3/friends => [{"id": 4, "name": "TJ"}]
     * ```
     *
     * @param {String} param
     * @param {Function} middleware
     * @returns {Router}
     */
    // param(param, middleware) {
    //   this.params[param] = middleware;
    //   for (let i = 0; i < this.stack.length; i++) {
    //     const route = this.stack[i];
    //     route.param(param, middleware);
    //   }

    //   return this;
    // }
  };

  /**
   * Create `router.verb()` methods, where *verb* is one of the HTTP verbs such
   * as `router.get()` or `router.post()`.
   *
   * Match URL patterns to callback functions or controller actions using `router.verb()`,
   * where **verb** is one of the HTTP verbs such as `router.get()` or `router.post()`.
   *
   * Additionally, `router.all()` can be used to match against all methods.
   *
   * ```javascript
   * router
   *   .get('/', (ctx, next) => {
   *     ctx.body = 'Hello World!';
   *   })
   *   .post('/users', (ctx, next) => {
   *     // ...
   *   })
   *   .put('/users/:id', (ctx, next) => {
   *     // ...
   *   })
   *   .del('/users/:id', (ctx, next) => {
   *     // ...
   *   })
   *   .all('/users/:id', (ctx, next) => {
   *     // ...
   *   });
   * ```
   *
   * When a route is matched, its path is available at `ctx._matchedRoute` and if named,
   * the name is available at `ctx._matchedRouteName`
   *
   * Route paths will be translated to regular expressions using
   * [path-to-regexp](https://github.com/pillarjs/path-to-regexp).
   *
   * Query strings will not be considered when matching requests.
   *
   * #### Named routes
   *
   * Routes can optionally have names. This allows generation of URLs and easy
   * renaming of URLs during development.
   *
   * ```javascript
   * router.get('user', '/users/:id', (ctx, next) => {
   *  // ...
   * });
   *
   * router.url('user', 3);
   * // => "/users/3"
   * ```
   *
   * #### Multiple middleware
   *
   * Multiple middleware may be given:
   *
   * ```javascript
   * router.get(
   *   '/users/:id',
   *   (ctx, next) => {
   *     return User.findOne(ctx.params.id).then(function(user) {
   *       ctx.user = user;
   *       next();
   *     });
   *   },
   *   ctx => {
   *     console.log(ctx.user);
   *     // => { id: 17, name: "Alex" }
   *   }
   * );
   * ```
   *
   * ### Nested routers
   *
   * Nesting routers is supported:
   *
   * ```javascript
   * const forums = new Router();
   * const posts = new Router();
   *
   * posts.get('/', (ctx, next) => {...});
   * posts.get('/:pid', (ctx, next) => {...});
   * forums.use('/forums/:fid/posts', posts.routes(), posts.allowedMethods());
   *
   * // responds to "/forums/123/posts" and "/forums/123/posts/123"
   * app.use(forums.routes());
   * ```
   *
   * #### Router prefixes
   *
   * Route paths can be prefixed at the router level:
   *
   * ```javascript
   * const router = new Router({
   *   prefix: '/users'
   * });
   *
   * router.get('/', ...); // responds to "/users"
   * router.get('/:id', ...); // responds to "/users/:id"
   * ```
   *
   * #### URL parameters
   *
   * Named route parameters are captured and added to `ctx.params`.
   *
   * ```javascript
   * router.get('/:category/:title', (ctx, next) => {
   *   console.log(ctx.params);
   *   // => { category: 'programming', title: 'how-to-node' }
   * });
   * ```
   *
   * The [path-to-regexp](https://github.com/pillarjs/path-to-regexp) module is
   * used to convert paths to regular expressions.
   *
   *
   * ### Match host for each router instance
   *
   * ```javascript
   * const router = new Router({
   *    host: 'example.domain' // only match if request host exactly equal `example.domain`
   * });
   *
   * ```
   *
   * OR host cloud be a regexp
   *
   * ```javascript
   * const router = new Router({
   *     host: /.*\.?example\.domain$/ // all host end with .example.domain would be matched
   * });
   * ```
   *
   * @name get|put|post|patch|delete|del
   * @memberof module:koa-router.prototype
   * @param {String} path
   * @param {Function=} middleware route middleware(s)
   * @param {Function} callback route callback
   * @returns {Router}
   */
  for (const method_ of methods) {
    function setMethodVerb(method) {
      Router$1.prototype[method] = function (name, path, middleware) {
        if (typeof path === 'string' || path instanceof RegExp) {
          middleware = Array.prototype.slice.call(arguments, 2);
        } else {
          middleware = Array.prototype.slice.call(arguments, 1);
          path = name;
          name = null;
        }

        // Sanity check to ensure we have a viable path candidate (eg: string|regex|non-empty array)
        if (typeof path !== 'string' && !(path instanceof RegExp) && (!Array.isArray(path) || path.length === 0)) throw new Error("You have to provide a path when adding a ".concat(method, " handler"));
        const opts = {
          name,
          pathAsRegExp: path instanceof RegExp
        };

        // pass opts to register call on verb methods
        this.register(path, [method], middleware, _objectSpread2(_objectSpread2({}, this.opts), opts));
        return this;
      };
    }
    setMethodVerb(method_);
  }

  // Alias for `router.delete()` because delete is a reserved word

  Router$1.prototype.del = Router$1.prototype['delete'];
  Router$1.prototype.verb = function (method, pathToMath, action) {
    const verb = method.toUpperCase();
    this[verb](pathToMath, action);
  };
  var router$1 = Router$1;
  var router_1 = router$1;

  /**
   * Module dependencies.
   */

  // const isGeneratorFunction = require('is-generator-function');
  // const debug = require('debug')('koa:application');
  // const onFinished = require('on-finished');
  // const assert = require('assert');
  const compose = compose_1$1;
  const context = contextExports;
  const request = request$1;
  const response = response$1;
  // const statuses = require('statuses');
  // const Emitter = require('events');
  // const util = require('util');
  // const Stream = require('stream');
  // const http = require('http');
  // const only = require('only');
  // const convert = require('koa-convert');
  // const deprecate = require('depd')('koa');
  // const { HttpError } = require('http-errors');
  const Client = client$1;
  const Router = router_1;

  /**
   * Expose `Application` class.
   * Inherits from `Emitter.prototype`.
   */

  var application = class Application /* extends Emitter */ {
    /**
     * Initialize a new `Application`.
     *
     * @api public
     */

    /**
      *
      * @param {object} [options] Application options
      * @param {string} [options.env='development'] Environment
      * @param {string[]} [options.keys] Signed cookie keys
      * @param {boolean} [options.proxy] Trust proxy headers
      * @param {number} [options.subdomainOffset] Subdomain offset
      * @param {string} [options.proxyIpHeader] Proxy IP header, defaults to X-Forwarded-For
      * @param {number} [options.maxIpsCount] Max IPs read from proxy IP header, default to 0 (means infinity)
      *
      */

    constructor(options) {
      // super();
      options = options || {};
      // this.proxy = options.proxy || false;
      // this.subdomainOffset = options.subdomainOffset || 2;
      // this.proxyIpHeader = options.proxyIpHeader || 'X-Forwarded-For';
      // this.maxIpsCount = options.maxIpsCount || 0;
      this.env = options.env || "production" || 'development';
      if (options.keys) this.keys = options.keys;
      this.middleware = [];
      this.context = Object.create(context);
      this.request = Object.create(request);
      this.response = Object.create(response);
      // util.inspect.custom support for node 6+
      // /* istanbul ignore else */
      // if (util.inspect.custom) {
      //   this[util.inspect.custom] = this.inspect;
      // }
      // if (options.asyncLocalStorage) {
      //   const { AsyncLocalStorage } = require('async_hooks');
      //   assert(AsyncLocalStorage, 'Requires node 12.17.0 or higher to enable asyncLocalStorage');
      //   this.ctxStorage = new AsyncLocalStorage();
      // }

      this.router = new Router({
        exclusive: true
      });
      this.use(this.router.routes());
      this.listen();
    }

    /**
     * Shorthand for:
     *
     *    http.createServer(app.callback()).listen(...)
     *
     * @param {Mixed} ...
     * @return {Server}
     * @api public
     */

    listen() {
      // debug('listen');
      // const server = http.createServer(this.callback());
      // return server.listen(...args);
      const handleRequest = this.callback();
      this.client = new Client({
        handleRequest
      });
    }

    /**
     * Return JSON representation.
     * We only bother showing settings.
     *
     * @return {Object}
     * @api public
     */

    // toJSON() {
    //   return only(this, [
    //     'subdomainOffset',
    //     'proxy',
    //     'env'
    //   ]);
    // }

    /**
     * Inspect implementation.
     *
     * @return {Object}
     * @api public
     */

    // inspect() {
    //   return this.toJSON();
    // }

    /**
     * Use the given middleware `fn`.
     *
     * Old-style middleware will be converted.
     *
     * @param {Function} fn
     * @return {Application} self
     * @api public
     */

    use(fn) {
      if (typeof fn !== 'function') throw new TypeError('middleware must be a function!');
      // if (isGeneratorFunction(fn)) {
      //   deprecate('Support for generators will be removed in v3. ' +
      //             'See the documentation for examples of how to convert old middleware ' +
      //             'https://github.com/koajs/koa/blob/master/docs/migration.md');
      //   fn = convert(fn);
      // }
      // debug('use %s', fn._name || fn.name || '-');
      this.middleware.push(fn);
      return this;
    }

    /**
     * Return a request handler callback
     * for node's native http server.
     *
     * @return {Function}
     * @api public
     */

    callback() {
      const fn = compose(this.middleware);

      // if (!this.listenerCount('error')) this.on('error', this.onerror);

      const handleRequest = (req, res) => {
        const ctx = this.createContext(req, res);
        // if (!this.ctxStorage) {
        return this.handleRequest(ctx, fn);
        // }
        // return this.ctxStorage.run(ctx, async() => {
        //   return await this.handleRequest(ctx, fn);
        // });
      };
      return handleRequest;
    }

    /**
     * return currnect contenxt from async local storage
     */
    // get currentContext() {
    //   if (this.ctxStorage) return this.ctxStorage.getStore();
    // }

    /**
     * Handle request in callback.
     *
     * @api private
     */

    handleRequest(ctx, fnMiddleware) {
      const res = ctx.res;
      res.statusCode = 404;
      const onerror = err => ctx.onerror(err);
      const handleResponse = () => respond(ctx);
      // onFinished(res, onerror);
      return fnMiddleware(ctx).then(handleResponse).catch(onerror);
    }

    /**
     * Initialize a new context.
     *
     * @api private
     */

    createContext(req, res) {
      const context = Object.create(this.context);
      const request = context.request = Object.create(this.request);
      const response = context.response = Object.create(this.response);
      context.app = request.app = response.app = this;
      context.req = request.req = response.req = req;
      context.res = request.res = response.res = res;
      // request.ctx = response.ctx = context;
      // request.response = response;
      // response.request = request;
      // context.originalUrl = request.originalUrl = req.url;
      context.state = {};
      return context;
    }

    /**
     * Default error handler.
     *
     * @param {Error} err
     * @api private
     */

    onerror(err) {
      // When dealing with cross-globals a normal `instanceof` check doesn't work properly.
      // See https://github.com/koajs/koa/issues/1466
      // We can probably remove it once jest fixes https://github.com/facebook/jest/issues/2549.
      const isNativeError = Object.prototype.toString.call(err) === '[object Error]' || err instanceof Error;
      if (!isNativeError) throw new TypeError(util.format('non-error thrown: %j', err));
      if (404 === err.status || err.expose) return;
      if (this.silent) return;
      const msg = err.stack || err.toString();
      console.error("\n".concat(msg.replace(/^/gm, '  '), "\n"));
    }

    /**
     * Help TS users comply to CommonJS, ESM, bundler mismatch.
     * @see https://github.com/koajs/koa/issues/1513
     */

    static get default() {
      return Application;
    }

    // createAsyncCtxStorageMiddleware() {
    //   const app = this;
    //   return async function asyncCtxStorage(ctx, next) {
    //     await app.ctxStorage.run(ctx, async() => {
    //       return await next();
    //     });
    //   };
    // }
  };

  /**
   * Response helper.
   */

  function respond(ctx) {
    // allow bypassing koa
    if (false === ctx.respond) return;
    if (!ctx.writable) return;
    const res = ctx.res;
    let body = ctx.body;
    ctx.status;

    // ignore body
    // if (statuses.empty[code]) {
    //   // strip headers
    //   ctx.body = null;
    //   return res.end();
    // }

    // if ('HEAD' === ctx.method) {
    //   if (!res.headersSent && !ctx.response.has('Content-Length')) {
    //     const { length } = ctx.response;
    //     if (Number.isInteger(length)) ctx.length = length;
    //   }
    //   return res.end();
    // }

    // status body
    if (null == body) {
      // if (ctx.response._explicitNullBody) {
      //   ctx.response.remove('Content-Type');
      //   ctx.response.remove('Transfer-Encoding');
      //   return res.end();
      // }
      // if (ctx.req.httpVersionMajor >= 2) {
      //   body = String(code);
      // } else {
      //   body = ctx.message || String(code);
      // }
      // if (!res.headersSent) {
      //   ctx.type = 'text';
      //   ctx.length = Buffer.byteLength(body);
      // }
      // return res.end(body);
      ctx.throw(res.statusCode);
    }

    // responses
    // if (Buffer.isBuffer(body)) return res.end(body);
    if ('string' === typeof body) return res.end(body);
    // if (body instanceof Stream) return body.pipe(res);

    // body: json
    // body = JSON.stringify(body);
    // if (!res.headersSent) {
    //   ctx.length = Buffer.byteLength(body);
    // }
    // res.end(body);
    return res.end(body);
  }

  var application_1 = application;

  const SOURCE = '@head/backstage'; // eslint-disable-line import/prefer-default-export

  function seq() {
    return Date.now() + '-' + Math.random().toString(36).substring(2);
  }
  function Callback() {
    const promised = {};
    return function callback(fn, resolved, rejected) {
      if (typeof fn === 'function') {
        const id = seq();
        return new Promise((resolve, reject) => {
          promised[id] = {
            resolve,
            reject
          };
          fn(id);
        });
      } else {
        // typeof fn === 'string'
        const {
          resolve,
          reject
        } = promised[fn];
        if (rejected) {
          reject(rejected);
        } else {
          resolve(resolved);
        }
        delete promised[fn];
      }
    };
  }

  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __decorateClass = (decorators, target, key, kind) => {
    var result = kind > 1 ? void 0 : kind ? __getOwnPropDesc(target, key) : target;
    for (var i = decorators.length - 1, decorator; i >= 0; i--)
      if (decorator = decorators[i])
        result = (kind ? decorator(target, key, result) : decorator(result)) || result;
    if (kind && result)
      __defProp(target, key, result);
    return result;
  };

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  const e$6=e=>n=>"function"==typeof n?((e,n)=>(customElements.define(e,n),n))(e,n):((e,n)=>{const{kind:t,elements:s}=n;return {kind:t,elements:s,finisher(n){customElements.define(e,n);}}})(e,n);

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  const i$5=(i,e)=>"method"===e.kind&&e.descriptor&&!("value"in e.descriptor)?{...e,finisher(n){n.createProperty(e.key,i);}}:{kind:"field",key:Symbol(),placement:"own",descriptor:{},originalKey:e.key,initializer(){"function"==typeof e.initializer&&(this[e.key]=e.initializer.call(this));},finisher(n){n.createProperty(e.key,i);}},e$5=(i,e,n)=>{e.constructor.createProperty(n,i);};function n$6(n){return (t,o)=>void 0!==o?e$5(n,t,o):i$5(n,t)}

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */function t$4(t){return n$6({...t,state:!0})}

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  const o$5=({finisher:e,descriptor:t})=>(o,n)=>{var r;if(void 0===n){const n=null!==(r=o.originalKey)&&void 0!==r?r:o.key,i=null!=t?{kind:"method",placement:"prototype",key:n,descriptor:t(o.key)}:{...o,key:n};return null!=e&&(i.finisher=function(t){e(t,n);}),i}{const r=o.constructor;void 0!==t&&Object.defineProperty(o,n,t(n)),null==e||e(r,n);}};

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */function i$4(i,n){return o$5({descriptor:o=>{const t={get(){var o,n;return null!==(n=null===(o=this.renderRoot)||void 0===o?void 0:o.querySelector(i))&&void 0!==n?n:null},enumerable:!0,configurable:!0};if(n){const n="symbol"==typeof o?Symbol():"__"+o;t.get=function(){var o,t;return void 0===this[n]&&(this[n]=null!==(t=null===(o=this.renderRoot)||void 0===o?void 0:o.querySelector(i))&&void 0!==t?t:null),this[n]};}return t}})}

  /**
   * @license
   * Copyright 2021 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */var n$5;null!=(null===(n$5=window.HTMLSlotElement)||void 0===n$5?void 0:n$5.prototype.assignedElements)?(o,n)=>o.assignedElements(n):(o,n)=>o.assignedNodes(n).filter((o=>o.nodeType===Node.ELEMENT_NODE));

  let basePath = "";
  function setBasePath(path) {
    basePath = path;
  }
  function getBasePath(subpath = "") {
    if (!basePath) {
      const scripts = [...document.getElementsByTagName("script")];
      const configScript = scripts.find((script) => script.hasAttribute("data-shoelace"));
      if (configScript) {
        setBasePath(configScript.getAttribute("data-shoelace"));
      } else {
        const fallbackScript = scripts.find((s) => {
          return /shoelace(\.min)?\.js($|\?)/.test(s.src) || /shoelace-autoloader(\.min)?\.js($|\?)/.test(s.src);
        });
        let path = "";
        if (fallbackScript) {
          path = fallbackScript.getAttribute("src");
        }
        setBasePath(path.split("/").slice(0, -1).join("/"));
      }
    }
    return basePath.replace(/\/$/, "") + (subpath ? `/${subpath.replace(/^\//, "")}` : ``);
  }

  const library = {
    name: "default",
    resolver: (name) => getBasePath(`assets/icons/${name}.svg`)
  };
  var library_default_default = library;

  const icons = {
    caret: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <polyline points="6 9 12 15 18 9"></polyline>
    </svg>
  `,
    check: `
    <svg part="checked-icon" class="checkbox__icon" viewBox="0 0 16 16">
      <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" stroke-linecap="round">
        <g stroke="currentColor" stroke-width="2">
          <g transform="translate(3.428571, 3.428571)">
            <path d="M0,5.71428571 L3.42857143,9.14285714"></path>
            <path d="M9.14285714,0 L3.42857143,9.14285714"></path>
          </g>
        </g>
      </g>
    </svg>
  `,
    "chevron-down": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-down" viewBox="0 0 16 16">
      <path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
    </svg>
  `,
    "chevron-left": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-left" viewBox="0 0 16 16">
      <path fill-rule="evenodd" d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"/>
    </svg>
  `,
    "chevron-right": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-right" viewBox="0 0 16 16">
      <path fill-rule="evenodd" d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/>
    </svg>
  `,
    eye: `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-eye" viewBox="0 0 16 16">
      <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.133 13.133 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.133 13.133 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5c-2.12 0-3.879-1.168-5.168-2.457A13.134 13.134 0 0 1 1.172 8z"/>
      <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0z"/>
    </svg>
  `,
    "eye-slash": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-eye-slash" viewBox="0 0 16 16">
      <path d="M13.359 11.238C15.06 9.72 16 8 16 8s-3-5.5-8-5.5a7.028 7.028 0 0 0-2.79.588l.77.771A5.944 5.944 0 0 1 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.134 13.134 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755-.165.165-.337.328-.517.486l.708.709z"/>
      <path d="M11.297 9.176a3.5 3.5 0 0 0-4.474-4.474l.823.823a2.5 2.5 0 0 1 2.829 2.829l.822.822zm-2.943 1.299.822.822a3.5 3.5 0 0 1-4.474-4.474l.823.823a2.5 2.5 0 0 0 2.829 2.829z"/>
      <path d="M3.35 5.47c-.18.16-.353.322-.518.487A13.134 13.134 0 0 0 1.172 8l.195.288c.335.48.83 1.12 1.465 1.755C4.121 11.332 5.881 12.5 8 12.5c.716 0 1.39-.133 2.02-.36l.77.772A7.029 7.029 0 0 1 8 13.5C3 13.5 0 8 0 8s.939-1.721 2.641-3.238l.708.709zm10.296 8.884-12-12 .708-.708 12 12-.708.708z"/>
    </svg>
  `,
    eyedropper: `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-eyedropper" viewBox="0 0 16 16">
      <path d="M13.354.646a1.207 1.207 0 0 0-1.708 0L8.5 3.793l-.646-.647a.5.5 0 1 0-.708.708L8.293 5l-7.147 7.146A.5.5 0 0 0 1 12.5v1.793l-.854.853a.5.5 0 1 0 .708.707L1.707 15H3.5a.5.5 0 0 0 .354-.146L11 7.707l1.146 1.147a.5.5 0 0 0 .708-.708l-.647-.646 3.147-3.146a1.207 1.207 0 0 0 0-1.708l-2-2zM2 12.707l7-7L10.293 7l-7 7H2v-1.293z"></path>
    </svg>
  `,
    "grip-vertical": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-grip-vertical" viewBox="0 0 16 16">
      <path d="M7 2a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm3 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM7 5a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm3 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM7 8a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm3 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm-3 3a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm3 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm-3 3a1 1 0 1 1-2 0 1 1 0 0 1 2 0zm3 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0z"></path>
    </svg>
  `,
    indeterminate: `
    <svg part="indeterminate-icon" class="checkbox__icon" viewBox="0 0 16 16">
      <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" stroke-linecap="round">
        <g stroke="currentColor" stroke-width="2">
          <g transform="translate(2.285714, 6.857143)">
            <path d="M10.2857143,1.14285714 L1.14285714,1.14285714"></path>
          </g>
        </g>
      </g>
    </svg>
  `,
    "person-fill": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-person-fill" viewBox="0 0 16 16">
      <path d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z"/>
    </svg>
  `,
    "play-fill": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-play-fill" viewBox="0 0 16 16">
      <path d="m11.596 8.697-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393z"></path>
    </svg>
  `,
    "pause-fill": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-pause-fill" viewBox="0 0 16 16">
      <path d="M5.5 3.5A1.5 1.5 0 0 1 7 5v6a1.5 1.5 0 0 1-3 0V5a1.5 1.5 0 0 1 1.5-1.5zm5 0A1.5 1.5 0 0 1 12 5v6a1.5 1.5 0 0 1-3 0V5a1.5 1.5 0 0 1 1.5-1.5z"></path>
    </svg>
  `,
    radio: `
    <svg part="checked-icon" class="radio__icon" viewBox="0 0 16 16">
      <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g fill="currentColor">
          <circle cx="8" cy="8" r="3.42857143"></circle>
        </g>
      </g>
    </svg>
  `,
    "star-fill": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-star-fill" viewBox="0 0 16 16">
      <path d="M3.612 15.443c-.386.198-.824-.149-.746-.592l.83-4.73L.173 6.765c-.329-.314-.158-.888.283-.95l4.898-.696L7.538.792c.197-.39.73-.39.927 0l2.184 4.327 4.898.696c.441.062.612.636.282.95l-3.522 3.356.83 4.73c.078.443-.36.79-.746.592L8 13.187l-4.389 2.256z"/>
    </svg>
  `,
    "x-lg": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
      <path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z"/>
    </svg>
  `,
    "x-circle-fill": `
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x-circle-fill" viewBox="0 0 16 16">
      <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zM5.354 4.646a.5.5 0 1 0-.708.708L7.293 8l-2.647 2.646a.5.5 0 0 0 .708.708L8 8.707l2.646 2.647a.5.5 0 0 0 .708-.708L8.707 8l2.647-2.646a.5.5 0 0 0-.708-.708L8 7.293 5.354 4.646z"></path>
    </svg>
  `
  };
  const systemLibrary = {
    name: "system",
    resolver: (name) => {
      if (name in icons) {
        return `data:image/svg+xml,${encodeURIComponent(icons[name])}`;
      }
      return "";
    }
  };
  var library_system_default = systemLibrary;

  let registry = [library_default_default, library_system_default];
  let watchedIcons = [];
  function watchIcon(icon) {
    watchedIcons.push(icon);
  }
  function unwatchIcon(icon) {
    watchedIcons = watchedIcons.filter((el) => el !== icon);
  }
  function getIconLibrary(name) {
    return registry.find((lib) => lib.name === name);
  }

  /**
   * @license
   * Copyright 2019 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  const t$3=window,e$4=t$3.ShadowRoot&&(void 0===t$3.ShadyCSS||t$3.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,s$4=Symbol(),n$4=new WeakMap;let o$4 = class o{constructor(t,e,n){if(this._$cssResult$=!0,n!==s$4)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=t,this.t=e;}get styleSheet(){let t=this.o;const s=this.t;if(e$4&&void 0===t){const e=void 0!==s&&1===s.length;e&&(t=n$4.get(s)),void 0===t&&((this.o=t=new CSSStyleSheet).replaceSync(this.cssText),e&&n$4.set(s,t));}return t}toString(){return this.cssText}};const r$2=t=>new o$4("string"==typeof t?t:t+"",void 0,s$4),i$3=(t,...e)=>{const n=1===t.length?t[0]:e.reduce(((e,s,n)=>e+(t=>{if(!0===t._$cssResult$)return t.cssText;if("number"==typeof t)return t;throw Error("Value passed to 'css' function must be a 'css' function result: "+t+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(s)+t[n+1]),t[0]);return new o$4(n,t,s$4)},S$1=(s,n)=>{e$4?s.adoptedStyleSheets=n.map((t=>t instanceof CSSStyleSheet?t:t.styleSheet)):n.forEach((e=>{const n=document.createElement("style"),o=t$3.litNonce;void 0!==o&&n.setAttribute("nonce",o),n.textContent=e.cssText,s.appendChild(n);}));},c$1=e$4?t=>t:t=>t instanceof CSSStyleSheet?(t=>{let e="";for(const s of t.cssRules)e+=s.cssText;return r$2(e)})(t):t;

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */var s$3;const e$3=window,r$1=e$3.trustedTypes,h$1=r$1?r$1.emptyScript:"",o$3=e$3.reactiveElementPolyfillSupport,n$3={toAttribute(t,i){switch(i){case Boolean:t=t?h$1:null;break;case Object:case Array:t=null==t?t:JSON.stringify(t);}return t},fromAttribute(t,i){let s=t;switch(i){case Boolean:s=null!==t;break;case Number:s=null===t?null:Number(t);break;case Object:case Array:try{s=JSON.parse(t);}catch(t){s=null;}}return s}},a$2=(t,i)=>i!==t&&(i==i||t==t),l$4={attribute:!0,type:String,converter:n$3,reflect:!1,hasChanged:a$2},d$1="finalized";let u$1 = class u extends HTMLElement{constructor(){super(),this._$Ei=new Map,this.isUpdatePending=!1,this.hasUpdated=!1,this._$El=null,this._$Eu();}static addInitializer(t){var i;this.finalize(),(null!==(i=this.h)&&void 0!==i?i:this.h=[]).push(t);}static get observedAttributes(){this.finalize();const t=[];return this.elementProperties.forEach(((i,s)=>{const e=this._$Ep(s,i);void 0!==e&&(this._$Ev.set(e,s),t.push(e));})),t}static createProperty(t,i=l$4){if(i.state&&(i.attribute=!1),this.finalize(),this.elementProperties.set(t,i),!i.noAccessor&&!this.prototype.hasOwnProperty(t)){const s="symbol"==typeof t?Symbol():"__"+t,e=this.getPropertyDescriptor(t,s,i);void 0!==e&&Object.defineProperty(this.prototype,t,e);}}static getPropertyDescriptor(t,i,s){return {get(){return this[i]},set(e){const r=this[t];this[i]=e,this.requestUpdate(t,r,s);},configurable:!0,enumerable:!0}}static getPropertyOptions(t){return this.elementProperties.get(t)||l$4}static finalize(){if(this.hasOwnProperty(d$1))return !1;this[d$1]=!0;const t=Object.getPrototypeOf(this);if(t.finalize(),void 0!==t.h&&(this.h=[...t.h]),this.elementProperties=new Map(t.elementProperties),this._$Ev=new Map,this.hasOwnProperty("properties")){const t=this.properties,i=[...Object.getOwnPropertyNames(t),...Object.getOwnPropertySymbols(t)];for(const s of i)this.createProperty(s,t[s]);}return this.elementStyles=this.finalizeStyles(this.styles),!0}static finalizeStyles(i){const s=[];if(Array.isArray(i)){const e=new Set(i.flat(1/0).reverse());for(const i of e)s.unshift(c$1(i));}else void 0!==i&&s.push(c$1(i));return s}static _$Ep(t,i){const s=i.attribute;return !1===s?void 0:"string"==typeof s?s:"string"==typeof t?t.toLowerCase():void 0}_$Eu(){var t;this._$E_=new Promise((t=>this.enableUpdating=t)),this._$AL=new Map,this._$Eg(),this.requestUpdate(),null===(t=this.constructor.h)||void 0===t||t.forEach((t=>t(this)));}addController(t){var i,s;(null!==(i=this._$ES)&&void 0!==i?i:this._$ES=[]).push(t),void 0!==this.renderRoot&&this.isConnected&&(null===(s=t.hostConnected)||void 0===s||s.call(t));}removeController(t){var i;null===(i=this._$ES)||void 0===i||i.splice(this._$ES.indexOf(t)>>>0,1);}_$Eg(){this.constructor.elementProperties.forEach(((t,i)=>{this.hasOwnProperty(i)&&(this._$Ei.set(i,this[i]),delete this[i]);}));}createRenderRoot(){var t;const s=null!==(t=this.shadowRoot)&&void 0!==t?t:this.attachShadow(this.constructor.shadowRootOptions);return S$1(s,this.constructor.elementStyles),s}connectedCallback(){var t;void 0===this.renderRoot&&(this.renderRoot=this.createRenderRoot()),this.enableUpdating(!0),null===(t=this._$ES)||void 0===t||t.forEach((t=>{var i;return null===(i=t.hostConnected)||void 0===i?void 0:i.call(t)}));}enableUpdating(t){}disconnectedCallback(){var t;null===(t=this._$ES)||void 0===t||t.forEach((t=>{var i;return null===(i=t.hostDisconnected)||void 0===i?void 0:i.call(t)}));}attributeChangedCallback(t,i,s){this._$AK(t,s);}_$EO(t,i,s=l$4){var e;const r=this.constructor._$Ep(t,s);if(void 0!==r&&!0===s.reflect){const h=(void 0!==(null===(e=s.converter)||void 0===e?void 0:e.toAttribute)?s.converter:n$3).toAttribute(i,s.type);this._$El=t,null==h?this.removeAttribute(r):this.setAttribute(r,h),this._$El=null;}}_$AK(t,i){var s;const e=this.constructor,r=e._$Ev.get(t);if(void 0!==r&&this._$El!==r){const t=e.getPropertyOptions(r),h="function"==typeof t.converter?{fromAttribute:t.converter}:void 0!==(null===(s=t.converter)||void 0===s?void 0:s.fromAttribute)?t.converter:n$3;this._$El=r,this[r]=h.fromAttribute(i,t.type),this._$El=null;}}requestUpdate(t,i,s){let e=!0;void 0!==t&&(((s=s||this.constructor.getPropertyOptions(t)).hasChanged||a$2)(this[t],i)?(this._$AL.has(t)||this._$AL.set(t,i),!0===s.reflect&&this._$El!==t&&(void 0===this._$EC&&(this._$EC=new Map),this._$EC.set(t,s))):e=!1),!this.isUpdatePending&&e&&(this._$E_=this._$Ej());}async _$Ej(){this.isUpdatePending=!0;try{await this._$E_;}catch(t){Promise.reject(t);}const t=this.scheduleUpdate();return null!=t&&await t,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){var t;if(!this.isUpdatePending)return;this.hasUpdated,this._$Ei&&(this._$Ei.forEach(((t,i)=>this[i]=t)),this._$Ei=void 0);let i=!1;const s=this._$AL;try{i=this.shouldUpdate(s),i?(this.willUpdate(s),null===(t=this._$ES)||void 0===t||t.forEach((t=>{var i;return null===(i=t.hostUpdate)||void 0===i?void 0:i.call(t)})),this.update(s)):this._$Ek();}catch(t){throw i=!1,this._$Ek(),t}i&&this._$AE(s);}willUpdate(t){}_$AE(t){var i;null===(i=this._$ES)||void 0===i||i.forEach((t=>{var i;return null===(i=t.hostUpdated)||void 0===i?void 0:i.call(t)})),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(t)),this.updated(t);}_$Ek(){this._$AL=new Map,this.isUpdatePending=!1;}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$E_}shouldUpdate(t){return !0}update(t){void 0!==this._$EC&&(this._$EC.forEach(((t,i)=>this._$EO(i,this[i],t))),this._$EC=void 0),this._$Ek();}updated(t){}firstUpdated(t){}};u$1[d$1]=!0,u$1.elementProperties=new Map,u$1.elementStyles=[],u$1.shadowRootOptions={mode:"open"},null==o$3||o$3({ReactiveElement:u$1}),(null!==(s$3=e$3.reactiveElementVersions)&&void 0!==s$3?s$3:e$3.reactiveElementVersions=[]).push("1.6.3");

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  var t$2;const i$2=window,s$2=i$2.trustedTypes,e$2=s$2?s$2.createPolicy("lit-html",{createHTML:t=>t}):void 0,o$2="$lit$",n$2=`lit$${(Math.random()+"").slice(9)}$`,l$3="?"+n$2,h=`<${l$3}>`,r=document,u=()=>r.createComment(""),d=t=>null===t||"object"!=typeof t&&"function"!=typeof t,c=Array.isArray,v=t=>c(t)||"function"==typeof(null==t?void 0:t[Symbol.iterator]),a$1="[ \t\n\f\r]",f=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,_=/-->/g,m=/>/g,p=RegExp(`>|${a$1}(?:([^\\s"'>=/]+)(${a$1}*=${a$1}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,"g"),g=/'/g,$=/"/g,y=/^(?:script|style|textarea|title)$/i,w=t=>(i,...s)=>({_$litType$:t,strings:i,values:s}),x=w(1),T=Symbol.for("lit-noChange"),A=Symbol.for("lit-nothing"),E=new WeakMap,C=r.createTreeWalker(r,129,null,!1);function P(t,i){if(!Array.isArray(t)||!t.hasOwnProperty("raw"))throw Error("invalid template strings array");return void 0!==e$2?e$2.createHTML(i):i}const V=(t,i)=>{const s=t.length-1,e=[];let l,r=2===i?"<svg>":"",u=f;for(let i=0;i<s;i++){const s=t[i];let d,c,v=-1,a=0;for(;a<s.length&&(u.lastIndex=a,c=u.exec(s),null!==c);)a=u.lastIndex,u===f?"!--"===c[1]?u=_:void 0!==c[1]?u=m:void 0!==c[2]?(y.test(c[2])&&(l=RegExp("</"+c[2],"g")),u=p):void 0!==c[3]&&(u=p):u===p?">"===c[0]?(u=null!=l?l:f,v=-1):void 0===c[1]?v=-2:(v=u.lastIndex-c[2].length,d=c[1],u=void 0===c[3]?p:'"'===c[3]?$:g):u===$||u===g?u=p:u===_||u===m?u=f:(u=p,l=void 0);const w=u===p&&t[i+1].startsWith("/>")?" ":"";r+=u===f?s+h:v>=0?(e.push(d),s.slice(0,v)+o$2+s.slice(v)+n$2+w):s+n$2+(-2===v?(e.push(void 0),i):w);}return [P(t,r+(t[s]||"<?>")+(2===i?"</svg>":"")),e]};class N{constructor({strings:t,_$litType$:i},e){let h;this.parts=[];let r=0,d=0;const c=t.length-1,v=this.parts,[a,f]=V(t,i);if(this.el=N.createElement(a,e),C.currentNode=this.el.content,2===i){const t=this.el.content,i=t.firstChild;i.remove(),t.append(...i.childNodes);}for(;null!==(h=C.nextNode())&&v.length<c;){if(1===h.nodeType){if(h.hasAttributes()){const t=[];for(const i of h.getAttributeNames())if(i.endsWith(o$2)||i.startsWith(n$2)){const s=f[d++];if(t.push(i),void 0!==s){const t=h.getAttribute(s.toLowerCase()+o$2).split(n$2),i=/([.?@])?(.*)/.exec(s);v.push({type:1,index:r,name:i[2],strings:t,ctor:"."===i[1]?H:"?"===i[1]?L:"@"===i[1]?z:k});}else v.push({type:6,index:r});}for(const i of t)h.removeAttribute(i);}if(y.test(h.tagName)){const t=h.textContent.split(n$2),i=t.length-1;if(i>0){h.textContent=s$2?s$2.emptyScript:"";for(let s=0;s<i;s++)h.append(t[s],u()),C.nextNode(),v.push({type:2,index:++r});h.append(t[i],u());}}}else if(8===h.nodeType)if(h.data===l$3)v.push({type:2,index:r});else {let t=-1;for(;-1!==(t=h.data.indexOf(n$2,t+1));)v.push({type:7,index:r}),t+=n$2.length-1;}r++;}}static createElement(t,i){const s=r.createElement("template");return s.innerHTML=t,s}}function S(t,i,s=t,e){var o,n,l,h;if(i===T)return i;let r=void 0!==e?null===(o=s._$Co)||void 0===o?void 0:o[e]:s._$Cl;const u=d(i)?void 0:i._$litDirective$;return (null==r?void 0:r.constructor)!==u&&(null===(n=null==r?void 0:r._$AO)||void 0===n||n.call(r,!1),void 0===u?r=void 0:(r=new u(t),r._$AT(t,s,e)),void 0!==e?(null!==(l=(h=s)._$Co)&&void 0!==l?l:h._$Co=[])[e]=r:s._$Cl=r),void 0!==r&&(i=S(t,r._$AS(t,i.values),r,e)),i}class M{constructor(t,i){this._$AV=[],this._$AN=void 0,this._$AD=t,this._$AM=i;}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(t){var i;const{el:{content:s},parts:e}=this._$AD,o=(null!==(i=null==t?void 0:t.creationScope)&&void 0!==i?i:r).importNode(s,!0);C.currentNode=o;let n=C.nextNode(),l=0,h=0,u=e[0];for(;void 0!==u;){if(l===u.index){let i;2===u.type?i=new R(n,n.nextSibling,this,t):1===u.type?i=new u.ctor(n,u.name,u.strings,this,t):6===u.type&&(i=new Z(n,this,t)),this._$AV.push(i),u=e[++h];}l!==(null==u?void 0:u.index)&&(n=C.nextNode(),l++);}return C.currentNode=r,o}v(t){let i=0;for(const s of this._$AV)void 0!==s&&(void 0!==s.strings?(s._$AI(t,s,i),i+=s.strings.length-2):s._$AI(t[i])),i++;}}class R{constructor(t,i,s,e){var o;this.type=2,this._$AH=A,this._$AN=void 0,this._$AA=t,this._$AB=i,this._$AM=s,this.options=e,this._$Cp=null===(o=null==e?void 0:e.isConnected)||void 0===o||o;}get _$AU(){var t,i;return null!==(i=null===(t=this._$AM)||void 0===t?void 0:t._$AU)&&void 0!==i?i:this._$Cp}get parentNode(){let t=this._$AA.parentNode;const i=this._$AM;return void 0!==i&&11===(null==t?void 0:t.nodeType)&&(t=i.parentNode),t}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(t,i=this){t=S(this,t,i),d(t)?t===A||null==t||""===t?(this._$AH!==A&&this._$AR(),this._$AH=A):t!==this._$AH&&t!==T&&this._(t):void 0!==t._$litType$?this.g(t):void 0!==t.nodeType?this.$(t):v(t)?this.T(t):this._(t);}k(t){return this._$AA.parentNode.insertBefore(t,this._$AB)}$(t){this._$AH!==t&&(this._$AR(),this._$AH=this.k(t));}_(t){this._$AH!==A&&d(this._$AH)?this._$AA.nextSibling.data=t:this.$(r.createTextNode(t)),this._$AH=t;}g(t){var i;const{values:s,_$litType$:e}=t,o="number"==typeof e?this._$AC(t):(void 0===e.el&&(e.el=N.createElement(P(e.h,e.h[0]),this.options)),e);if((null===(i=this._$AH)||void 0===i?void 0:i._$AD)===o)this._$AH.v(s);else {const t=new M(o,this),i=t.u(this.options);t.v(s),this.$(i),this._$AH=t;}}_$AC(t){let i=E.get(t.strings);return void 0===i&&E.set(t.strings,i=new N(t)),i}T(t){c(this._$AH)||(this._$AH=[],this._$AR());const i=this._$AH;let s,e=0;for(const o of t)e===i.length?i.push(s=new R(this.k(u()),this.k(u()),this,this.options)):s=i[e],s._$AI(o),e++;e<i.length&&(this._$AR(s&&s._$AB.nextSibling,e),i.length=e);}_$AR(t=this._$AA.nextSibling,i){var s;for(null===(s=this._$AP)||void 0===s||s.call(this,!1,!0,i);t&&t!==this._$AB;){const i=t.nextSibling;t.remove(),t=i;}}setConnected(t){var i;void 0===this._$AM&&(this._$Cp=t,null===(i=this._$AP)||void 0===i||i.call(this,t));}}class k{constructor(t,i,s,e,o){this.type=1,this._$AH=A,this._$AN=void 0,this.element=t,this.name=i,this._$AM=e,this.options=o,s.length>2||""!==s[0]||""!==s[1]?(this._$AH=Array(s.length-1).fill(new String),this.strings=s):this._$AH=A;}get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}_$AI(t,i=this,s,e){const o=this.strings;let n=!1;if(void 0===o)t=S(this,t,i,0),n=!d(t)||t!==this._$AH&&t!==T,n&&(this._$AH=t);else {const e=t;let l,h;for(t=o[0],l=0;l<o.length-1;l++)h=S(this,e[s+l],i,l),h===T&&(h=this._$AH[l]),n||(n=!d(h)||h!==this._$AH[l]),h===A?t=A:t!==A&&(t+=(null!=h?h:"")+o[l+1]),this._$AH[l]=h;}n&&!e&&this.j(t);}j(t){t===A?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,null!=t?t:"");}}class H extends k{constructor(){super(...arguments),this.type=3;}j(t){this.element[this.name]=t===A?void 0:t;}}const I=s$2?s$2.emptyScript:"";class L extends k{constructor(){super(...arguments),this.type=4;}j(t){t&&t!==A?this.element.setAttribute(this.name,I):this.element.removeAttribute(this.name);}}class z extends k{constructor(t,i,s,e,o){super(t,i,s,e,o),this.type=5;}_$AI(t,i=this){var s;if((t=null!==(s=S(this,t,i,0))&&void 0!==s?s:A)===T)return;const e=this._$AH,o=t===A&&e!==A||t.capture!==e.capture||t.once!==e.once||t.passive!==e.passive,n=t!==A&&(e===A||o);o&&this.element.removeEventListener(this.name,this,e),n&&this.element.addEventListener(this.name,this,t),this._$AH=t;}handleEvent(t){var i,s;"function"==typeof this._$AH?this._$AH.call(null!==(s=null===(i=this.options)||void 0===i?void 0:i.host)&&void 0!==s?s:this.element,t):this._$AH.handleEvent(t);}}class Z{constructor(t,i,s){this.element=t,this.type=6,this._$AN=void 0,this._$AM=i,this.options=s;}get _$AU(){return this._$AM._$AU}_$AI(t){S(this,t);}}const B=i$2.litHtmlPolyfillSupport;null==B||B(N,R),(null!==(t$2=i$2.litHtmlVersions)&&void 0!==t$2?t$2:i$2.litHtmlVersions=[]).push("2.8.0");const D=(t,i,s)=>{var e,o;const n=null!==(e=null==s?void 0:s.renderBefore)&&void 0!==e?e:i;let l=n._$litPart$;if(void 0===l){const t=null!==(o=null==s?void 0:s.renderBefore)&&void 0!==o?o:null;n._$litPart$=l=new R(i.insertBefore(u(),t),t,void 0,null!=s?s:{});}return l._$AI(t),l};

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */var l$2,o$1;let s$1 = class s extends u$1{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0;}createRenderRoot(){var t,e;const i=super.createRenderRoot();return null!==(t=(e=this.renderOptions).renderBefore)&&void 0!==t||(e.renderBefore=i.firstChild),i}update(t){const i=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(t),this._$Do=D(i,this.renderRoot,this.renderOptions);}connectedCallback(){var t;super.connectedCallback(),null===(t=this._$Do)||void 0===t||t.setConnected(!0);}disconnectedCallback(){var t;super.disconnectedCallback(),null===(t=this._$Do)||void 0===t||t.setConnected(!1);}render(){return T}};s$1.finalized=!0,s$1._$litElement$=!0,null===(l$2=globalThis.litElementHydrateSupport)||void 0===l$2||l$2.call(globalThis,{LitElement:s$1});const n$1=globalThis.litElementPolyfillSupport;null==n$1||n$1({LitElement:s$1});(null!==(o$1=globalThis.litElementVersions)&&void 0!==o$1?o$1:globalThis.litElementVersions=[]).push("3.3.3");

  /**
   * @license
   * Copyright 2020 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */const t$1=(o,l)=>void 0===l?void 0!==(null==o?void 0:o._$litType$):(null==o?void 0:o._$litType$)===l;

  function watch(propertyName, options) {
      const resolvedOptions = Object.assign({ waitUntilFirstUpdate: false }, options);
      return (proto, decoratedFnName) => {
          const { update } = proto;
          const watchedProperties = Array.isArray(propertyName) ? propertyName : [propertyName];
          proto.update = function (changedProps) {
              watchedProperties.forEach(property => {
                  const key = property;
                  if (changedProps.has(key)) {
                      const oldValue = changedProps.get(key);
                      const newValue = this[key];
                      if (oldValue !== newValue) {
                          if (!resolvedOptions.waitUntilFirstUpdate || this.hasUpdated) {
                              this[decoratedFnName](oldValue, newValue);
                          }
                      }
                  }
              });
              update.call(this, changedProps);
          };
      };
  }

  var __decorate = (window && window.__decorate) || function (decorators, target, key, desc) {
      var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
      if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
      else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
      return c > 3 && r && Object.defineProperty(target, key, r), r;
  };
  class ShoelaceElement extends s$1 {
      emit(name, options) {
          const event = new CustomEvent(name, Object.assign({ bubbles: true, cancelable: false, composed: true, detail: {} }, options));
          this.dispatchEvent(event);
          return event;
      }
  }
  __decorate([
      n$6()
  ], ShoelaceElement.prototype, "dir", void 0);
  __decorate([
      n$6()
  ], ShoelaceElement.prototype, "lang", void 0);

  var componentStyles = i$3 `
  :host {
    box-sizing: border-box;
  }

  :host *,
  :host *::before,
  :host *::after {
    box-sizing: inherit;
  }

  [hidden] {
    display: none !important;
  }
`;

  var icon_styles_default = i$3`
  ${componentStyles}

  :host {
    display: inline-block;
    width: 1em;
    height: 1em;
    box-sizing: content-box !important;
  }

  svg {
    display: block;
    height: 100%;
    width: 100%;
  }
`;

  const CACHEABLE_ERROR = Symbol();
  const RETRYABLE_ERROR = Symbol();
  let parser;
  const iconCache = /* @__PURE__ */ new Map();
  let SlIcon = class extends ShoelaceElement {
    constructor() {
      super(...arguments);
      this.initialRender = false;
      this.svg = null;
      this.label = "";
      this.library = "default";
    }
    /** Given a URL, this function returns the resulting SVG element or an appropriate error symbol. */
    async resolveIcon(url, library) {
      var _a;
      let fileData;
      if (library == null ? void 0 : library.spriteSheet) {
        return x`<svg part="svg">
        <use part="use" href="${url}"></use>
      </svg>`;
      }
      try {
        fileData = await fetch(url, { mode: "cors" });
        if (!fileData.ok)
          return fileData.status === 410 ? CACHEABLE_ERROR : RETRYABLE_ERROR;
      } catch (e) {
        return RETRYABLE_ERROR;
      }
      try {
        const div = document.createElement("div");
        div.innerHTML = await fileData.text();
        const svg = div.firstElementChild;
        if (((_a = svg == null ? void 0 : svg.tagName) == null ? void 0 : _a.toLowerCase()) !== "svg")
          return CACHEABLE_ERROR;
        if (!parser)
          parser = new DOMParser();
        const doc = parser.parseFromString(svg.outerHTML, "text/html");
        const svgEl = doc.body.querySelector("svg");
        if (!svgEl)
          return CACHEABLE_ERROR;
        svgEl.part.add("svg");
        return document.adoptNode(svgEl);
      } catch (e) {
        return CACHEABLE_ERROR;
      }
    }
    connectedCallback() {
      super.connectedCallback();
      watchIcon(this);
    }
    firstUpdated() {
      this.initialRender = true;
      this.setIcon();
    }
    disconnectedCallback() {
      super.disconnectedCallback();
      unwatchIcon(this);
    }
    getUrl() {
      const library = getIconLibrary(this.library);
      if (this.name && library) {
        return library.resolver(this.name);
      }
      return this.src;
    }
    handleLabelChange() {
      const hasLabel = typeof this.label === "string" && this.label.length > 0;
      if (hasLabel) {
        this.setAttribute("role", "img");
        this.setAttribute("aria-label", this.label);
        this.removeAttribute("aria-hidden");
      } else {
        this.removeAttribute("role");
        this.removeAttribute("aria-label");
        this.setAttribute("aria-hidden", "true");
      }
    }
    async setIcon() {
      var _a;
      const library = getIconLibrary(this.library);
      const url = this.getUrl();
      if (!url) {
        this.svg = null;
        return;
      }
      let iconResolver = iconCache.get(url);
      if (!iconResolver) {
        iconResolver = this.resolveIcon(url, library);
        iconCache.set(url, iconResolver);
      }
      if (!this.initialRender) {
        return;
      }
      const svg = await iconResolver;
      if (svg === RETRYABLE_ERROR) {
        iconCache.delete(url);
      }
      if (url !== this.getUrl()) {
        return;
      }
      if (t$1(svg)) {
        this.svg = svg;
        return;
      }
      switch (svg) {
        case RETRYABLE_ERROR:
        case CACHEABLE_ERROR:
          this.svg = null;
          this.emit("sl-error");
          break;
        default:
          this.svg = svg.cloneNode(true);
          (_a = library == null ? void 0 : library.mutator) == null ? void 0 : _a.call(library, this.svg);
          this.emit("sl-load");
      }
    }
    render() {
      return this.svg;
    }
  };
  SlIcon.styles = icon_styles_default;
  __decorateClass([
    t$4()
  ], SlIcon.prototype, "svg", 2);
  __decorateClass([
    n$6({ reflect: true })
  ], SlIcon.prototype, "name", 2);
  __decorateClass([
    n$6()
  ], SlIcon.prototype, "src", 2);
  __decorateClass([
    n$6()
  ], SlIcon.prototype, "label", 2);
  __decorateClass([
    n$6({ reflect: true })
  ], SlIcon.prototype, "library", 2);
  __decorateClass([
    watch("label")
  ], SlIcon.prototype, "handleLabelChange", 1);
  __decorateClass([
    watch(["name", "src", "library"])
  ], SlIcon.prototype, "setIcon", 1);
  SlIcon = __decorateClass([
    e$6("sl-icon")
  ], SlIcon);

  /**
   * @license
   * Copyright 2017 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */
  const t={ATTRIBUTE:1,CHILD:2,PROPERTY:3,BOOLEAN_ATTRIBUTE:4,EVENT:5,ELEMENT:6},e$1=t=>(...e)=>({_$litDirective$:t,values:e});let i$1 = class i{constructor(t){}get _$AU(){return this._$AM._$AU}_$AT(t,e,i){this._$Ct=t,this._$AM=e,this._$Ci=i;}_$AS(t,e){return this.update(t,e)}update(t,e){return this.render(...e)}};

  /**
   * @license
   * Copyright 2018 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */const o=e$1(class extends i$1{constructor(t$1){var i;if(super(t$1),t$1.type!==t.ATTRIBUTE||"class"!==t$1.name||(null===(i=t$1.strings)||void 0===i?void 0:i.length)>2)throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.")}render(t){return " "+Object.keys(t).filter((i=>t[i])).join(" ")+" "}update(i,[s]){var r,o;if(void 0===this.it){this.it=new Set,void 0!==i.strings&&(this.nt=new Set(i.strings.join(" ").split(/\s/).filter((t=>""!==t))));for(const t in s)s[t]&&!(null===(r=this.nt)||void 0===r?void 0:r.has(t))&&this.it.add(t);return this.render(s)}const e=i.element.classList;this.it.forEach((t=>{t in s||(e.remove(t),this.it.delete(t));}));for(const t in s){const i=!!s[t];i===this.it.has(t)||(null===(o=this.nt)||void 0===o?void 0:o.has(t))||(i?(e.add(t),this.it.add(t)):(e.remove(t),this.it.delete(t)));}return T}});

  /**
   * @license
   * Copyright 2020 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */const e=Symbol.for(""),l$1=t=>{if((null==t?void 0:t.r)===e)return null==t?void 0:t._$litStatic$},i=(t,...r)=>({_$litStatic$:r.reduce(((r,e,l)=>r+(t=>{if(void 0!==t._$litStatic$)return t._$litStatic$;throw Error(`Value passed to 'literal' function must be a 'literal' result: ${t}. Use 'unsafeStatic' to pass non-literal values, but\n            take care to ensure page security.`)})(e)+t[l+1]),t[0]),r:e}),s=new Map,a=t=>(r,...e)=>{const o=e.length;let i,a;const n=[],u=[];let c,$=0,f=!1;for(;$<o;){for(c=r[$];$<o&&void 0!==(a=e[$],i=l$1(a));)c+=i+r[++$],f=!0;$!==o&&u.push(a),n.push(c),$++;}if($===o&&n.push(r[o]),f){const t=n.join("$$lit$$");void 0===(r=s.get(t))&&(n.raw=n,s.set(t,r=n)),e=u;}return t(r,...e)},n=a(x);

  /**
   * @license
   * Copyright 2018 Google LLC
   * SPDX-License-Identifier: BSD-3-Clause
   */const l=l=>null!=l?l:A;

  var icon_button_styles_default = i$3`
  ${componentStyles}

  :host {
    display: inline-block;
    color: var(--sl-color-neutral-600);
  }

  .icon-button {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    background: none;
    border: none;
    border-radius: var(--sl-border-radius-medium);
    font-size: inherit;
    color: inherit;
    padding: var(--sl-spacing-x-small);
    cursor: pointer;
    transition: var(--sl-transition-x-fast) color;
    -webkit-appearance: none;
  }

  .icon-button:hover:not(.icon-button--disabled),
  .icon-button:focus-visible:not(.icon-button--disabled) {
    color: var(--sl-color-primary-600);
  }

  .icon-button:active:not(.icon-button--disabled) {
    color: var(--sl-color-primary-700);
  }

  .icon-button:focus {
    outline: none;
  }

  .icon-button--disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .icon-button:focus-visible {
    outline: var(--sl-focus-ring);
    outline-offset: var(--sl-focus-ring-offset);
  }

  .icon-button__icon {
    pointer-events: none;
  }
`;

  let SlIconButton = class extends ShoelaceElement {
    constructor() {
      super(...arguments);
      this.hasFocus = false;
      this.label = "";
      this.disabled = false;
    }
    handleBlur() {
      this.hasFocus = false;
      this.emit("sl-blur");
    }
    handleFocus() {
      this.hasFocus = true;
      this.emit("sl-focus");
    }
    handleClick(event) {
      if (this.disabled) {
        event.preventDefault();
        event.stopPropagation();
      }
    }
    /** Simulates a click on the icon button. */
    click() {
      this.button.click();
    }
    /** Sets focus on the icon button. */
    focus(options) {
      this.button.focus(options);
    }
    /** Removes focus from the icon button. */
    blur() {
      this.button.blur();
    }
    render() {
      const isLink = this.href ? true : false;
      const tag = isLink ? i`a` : i`button`;
      return n`
      <${tag}
        part="base"
        class=${o({
      "icon-button": true,
      "icon-button--disabled": !isLink && this.disabled,
      "icon-button--focused": this.hasFocus
    })}
        ?disabled=${l(isLink ? void 0 : this.disabled)}
        type=${l(isLink ? void 0 : "button")}
        href=${l(isLink ? this.href : void 0)}
        target=${l(isLink ? this.target : void 0)}
        download=${l(isLink ? this.download : void 0)}
        rel=${l(isLink && this.target ? "noreferrer noopener" : void 0)}
        role=${l(isLink ? void 0 : "button")}
        aria-disabled=${this.disabled ? "true" : "false"}
        aria-label="${this.label}"
        tabindex=${this.disabled ? "-1" : "0"}
        @blur=${this.handleBlur}
        @focus=${this.handleFocus}
        @click=${this.handleClick}
      >
        <sl-icon
          class="icon-button__icon"
          name=${l(this.name)}
          library=${l(this.library)}
          src=${l(this.src)}
          aria-hidden="true"
        ></sl-icon>
      </${tag}>
    `;
    }
  };
  SlIconButton.styles = icon_button_styles_default;
  __decorateClass([
    i$4(".icon-button")
  ], SlIconButton.prototype, "button", 2);
  __decorateClass([
    t$4()
  ], SlIconButton.prototype, "hasFocus", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "name", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "library", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "src", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "href", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "target", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "download", 2);
  __decorateClass([
    n$6()
  ], SlIconButton.prototype, "label", 2);
  __decorateClass([
    n$6({ type: Boolean, reflect: true })
  ], SlIconButton.prototype, "disabled", 2);
  SlIconButton = __decorateClass([
    e$6("sl-icon-button")
  ], SlIconButton);

  function animateTo(el, keyframes, options) {
      return new Promise(resolve => {
          if ((options === null || options === void 0 ? void 0 : options.duration) === Infinity) {
              throw new Error('Promise-based animations must be finite.');
          }
          const animation = el.animate(keyframes, Object.assign(Object.assign({}, options), { duration: prefersReducedMotion() ? 0 : options.duration }));
          animation.addEventListener('cancel', resolve, { once: true });
          animation.addEventListener('finish', resolve, { once: true });
      });
  }
  function prefersReducedMotion() {
      const query = window.matchMedia('(prefers-reduced-motion: reduce)');
      return query.matches;
  }
  function stopAnimations(el) {
      return Promise.all(el.getAnimations().map(animation => {
          return new Promise(resolve => {
              const handleAnimationEvent = requestAnimationFrame(resolve);
              animation.addEventListener('cancel', () => handleAnimationEvent, { once: true });
              animation.addEventListener('finish', () => handleAnimationEvent, { once: true });
              animation.cancel();
          });
      }));
  }

  const defaultAnimationRegistry = /* @__PURE__ */ new Map();
  const customAnimationRegistry = /* @__PURE__ */ new WeakMap();
  function ensureAnimation(animation) {
    return animation != null ? animation : { keyframes: [], options: { duration: 0 } };
  }
  function getLogicalAnimation(animation, dir) {
    if (dir.toLowerCase() === "rtl") {
      return {
        keyframes: animation.rtlKeyframes || animation.keyframes,
        options: animation.options
      };
    }
    return animation;
  }
  function setDefaultAnimation(animationName, animation) {
    defaultAnimationRegistry.set(animationName, ensureAnimation(animation));
  }
  function getAnimation(el, animationName, options) {
    const customAnimation = customAnimationRegistry.get(el);
    if (customAnimation == null ? void 0 : customAnimation[animationName]) {
      return getLogicalAnimation(customAnimation[animationName], options.dir);
    }
    const defaultAnimation = defaultAnimationRegistry.get(animationName);
    if (defaultAnimation) {
      return getLogicalAnimation(defaultAnimation, options.dir);
    }
    return {
      keyframes: [],
      options: { duration: 0 }
    };
  }

  class HasSlotController {
      constructor(host, ...slotNames) {
          this.slotNames = [];
          (this.host = host).addController(this);
          this.slotNames = slotNames;
          this.handleSlotChange = this.handleSlotChange.bind(this);
      }
      hasDefaultSlot() {
          return [...this.host.childNodes].some(node => {
              if (node.nodeType === node.TEXT_NODE && node.textContent.trim() !== '') {
                  return true;
              }
              if (node.nodeType === node.ELEMENT_NODE) {
                  const el = node;
                  const tagName = el.tagName.toLowerCase();
                  if (tagName === 'sl-visually-hidden') {
                      return false;
                  }
                  if (!el.hasAttribute('slot')) {
                      return true;
                  }
              }
              return false;
          });
      }
      hasNamedSlot(name) {
          return this.host.querySelector(`:scope > [slot="${name}"]`) !== null;
      }
      test(slotName) {
          return slotName === '[default]' ? this.hasDefaultSlot() : this.hasNamedSlot(slotName);
      }
      hostConnected() {
          this.host.shadowRoot.addEventListener('slotchange', this.handleSlotChange);
      }
      hostDisconnected() {
          this.host.shadowRoot.removeEventListener('slotchange', this.handleSlotChange);
      }
      handleSlotChange(event) {
          const slot = event.target;
          if ((this.slotNames.includes('[default]') && !slot.name) || (slot.name && this.slotNames.includes(slot.name))) {
              this.host.requestUpdate();
          }
      }
  }

  const connectedElements = new Set();
  const translations = new Map();
  let fallback;
  let documentDirection = 'ltr';
  let documentLanguage = 'en';
  const isClient = (typeof MutationObserver !== "undefined" && typeof document !== "undefined" && typeof document.documentElement !== "undefined");
  if (isClient) {
      const documentElementObserver = new MutationObserver(update);
      documentDirection = document.documentElement.dir || 'ltr';
      documentLanguage = document.documentElement.lang || navigator.language;
      documentElementObserver.observe(document.documentElement, {
          attributes: true,
          attributeFilter: ['dir', 'lang']
      });
  }
  function registerTranslation(...translation) {
      translation.map(t => {
          const code = t.$code.toLowerCase();
          if (translations.has(code)) {
              translations.set(code, Object.assign(Object.assign({}, translations.get(code)), t));
          }
          else {
              translations.set(code, t);
          }
          if (!fallback) {
              fallback = t;
          }
      });
      update();
  }
  function update() {
      if (isClient) {
          documentDirection = document.documentElement.dir || 'ltr';
          documentLanguage = document.documentElement.lang || navigator.language;
      }
      [...connectedElements.keys()].map((el) => {
          if (typeof el.requestUpdate === 'function') {
              el.requestUpdate();
          }
      });
  }
  let LocalizeController$1 = class LocalizeController {
      constructor(host) {
          this.host = host;
          this.host.addController(this);
      }
      hostConnected() {
          connectedElements.add(this.host);
      }
      hostDisconnected() {
          connectedElements.delete(this.host);
      }
      dir() {
          return `${this.host.dir || documentDirection}`.toLowerCase();
      }
      lang() {
          return `${this.host.lang || documentLanguage}`.toLowerCase();
      }
      getTranslationData(lang) {
          var _a, _b;
          const locale = new Intl.Locale(lang.replace(/_/g, '-'));
          const language = locale === null || locale === void 0 ? void 0 : locale.language.toLowerCase();
          const region = (_b = (_a = locale === null || locale === void 0 ? void 0 : locale.region) === null || _a === void 0 ? void 0 : _a.toLowerCase()) !== null && _b !== void 0 ? _b : '';
          const primary = translations.get(`${language}-${region}`);
          const secondary = translations.get(language);
          return { locale, language, region, primary, secondary };
      }
      exists(key, options) {
          var _a;
          const { primary, secondary } = this.getTranslationData((_a = options.lang) !== null && _a !== void 0 ? _a : this.lang());
          options = Object.assign({ includeFallback: false }, options);
          if ((primary && primary[key]) ||
              (secondary && secondary[key]) ||
              (options.includeFallback && fallback && fallback[key])) {
              return true;
          }
          return false;
      }
      term(key, ...args) {
          const { primary, secondary } = this.getTranslationData(this.lang());
          let term;
          if (primary && primary[key]) {
              term = primary[key];
          }
          else if (secondary && secondary[key]) {
              term = secondary[key];
          }
          else if (fallback && fallback[key]) {
              term = fallback[key];
          }
          else {
              console.error(`No translation found for: ${String(key)}`);
              return String(key);
          }
          if (typeof term === 'function') {
              return term(...args);
          }
          return term;
      }
      date(dateToFormat, options) {
          dateToFormat = new Date(dateToFormat);
          return new Intl.DateTimeFormat(this.lang(), options).format(dateToFormat);
      }
      number(numberToFormat, options) {
          numberToFormat = Number(numberToFormat);
          return isNaN(numberToFormat) ? '' : new Intl.NumberFormat(this.lang(), options).format(numberToFormat);
      }
      relativeTime(value, unit, options) {
          return new Intl.RelativeTimeFormat(this.lang(), options).format(value, unit);
      }
  };

  const translation = {
    $code: "en",
    $name: "English",
    $dir: "ltr",
    carousel: "Carousel",
    clearEntry: "Clear entry",
    close: "Close",
    copy: "Copy",
    currentValue: "Current value",
    goToSlide: (slide, count) => `Go to slide ${slide} of ${count}`,
    hidePassword: "Hide password",
    loading: "Loading",
    nextSlide: "Next slide",
    numOptionsSelected: (num) => {
      if (num === 0)
        return "No options selected";
      if (num === 1)
        return "1 option selected";
      return `${num} options selected`;
    },
    previousSlide: "Previous slide",
    progress: "Progress",
    remove: "Remove",
    resize: "Resize",
    scrollToEnd: "Scroll to end",
    scrollToStart: "Scroll to start",
    selectAColorFromTheScreen: "Select a color from the screen",
    showPassword: "Show password",
    slideNum: (slide) => `Slide ${slide}`,
    toggleColorFormat: "Toggle color format"
  };
  registerTranslation(translation);

  class LocalizeController extends LocalizeController$1 {
  }

  function waitForEvent(el, eventName) {
      return new Promise(resolve => {
          function done(event) {
              if (event.target === el) {
                  el.removeEventListener(eventName, done);
                  resolve();
              }
          }
          el.addEventListener(eventName, done);
      });
  }

  var alert_styles_default = i$3`
  ${componentStyles}

  :host {
    display: contents;

    /* For better DX, we'll reset the margin here so the base part can inherit it */
    margin: 0;
  }

  .alert {
    position: relative;
    display: flex;
    align-items: stretch;
    background-color: var(--sl-panel-background-color);
    border: solid var(--sl-panel-border-width) var(--sl-panel-border-color);
    border-top-width: calc(var(--sl-panel-border-width) * 3);
    border-radius: var(--sl-border-radius-medium);
    font-family: var(--sl-font-sans);
    font-size: var(--sl-font-size-small);
    font-weight: var(--sl-font-weight-normal);
    line-height: 1.6;
    color: var(--sl-color-neutral-700);
    margin: inherit;
  }

  .alert:not(.alert--has-icon) .alert__icon,
  .alert:not(.alert--closable) .alert__close-button {
    display: none;
  }

  .alert__icon {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    font-size: var(--sl-font-size-large);
    padding-inline-start: var(--sl-spacing-large);
  }

  .alert--primary {
    border-top-color: var(--sl-color-primary-600);
  }

  .alert--primary .alert__icon {
    color: var(--sl-color-primary-600);
  }

  .alert--success {
    border-top-color: var(--sl-color-success-600);
  }

  .alert--success .alert__icon {
    color: var(--sl-color-success-600);
  }

  .alert--neutral {
    border-top-color: var(--sl-color-neutral-600);
  }

  .alert--neutral .alert__icon {
    color: var(--sl-color-neutral-600);
  }

  .alert--warning {
    border-top-color: var(--sl-color-warning-600);
  }

  .alert--warning .alert__icon {
    color: var(--sl-color-warning-600);
  }

  .alert--danger {
    border-top-color: var(--sl-color-danger-600);
  }

  .alert--danger .alert__icon {
    color: var(--sl-color-danger-600);
  }

  .alert__message {
    flex: 1 1 auto;
    display: block;
    padding: var(--sl-spacing-large);
    overflow: hidden;
  }

  .alert__close-button {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    font-size: var(--sl-font-size-medium);
    padding-inline-end: var(--sl-spacing-medium);
  }
`;

  const toastStack = Object.assign(document.createElement("div"), { className: "sl-toast-stack" });
  let SlAlert = class extends ShoelaceElement {
    constructor() {
      super(...arguments);
      this.hasSlotController = new HasSlotController(this, "icon", "suffix");
      this.localize = new LocalizeController(this);
      this.open = false;
      this.closable = false;
      this.variant = "primary";
      this.duration = Infinity;
    }
    firstUpdated() {
      this.base.hidden = !this.open;
    }
    restartAutoHide() {
      clearTimeout(this.autoHideTimeout);
      if (this.open && this.duration < Infinity) {
        this.autoHideTimeout = window.setTimeout(() => this.hide(), this.duration);
      }
    }
    handleCloseClick() {
      this.hide();
    }
    handleMouseMove() {
      this.restartAutoHide();
    }
    async handleOpenChange() {
      if (this.open) {
        this.emit("sl-show");
        if (this.duration < Infinity) {
          this.restartAutoHide();
        }
        await stopAnimations(this.base);
        this.base.hidden = false;
        const { keyframes, options } = getAnimation(this, "alert.show", { dir: this.localize.dir() });
        await animateTo(this.base, keyframes, options);
        this.emit("sl-after-show");
      } else {
        this.emit("sl-hide");
        clearTimeout(this.autoHideTimeout);
        await stopAnimations(this.base);
        const { keyframes, options } = getAnimation(this, "alert.hide", { dir: this.localize.dir() });
        await animateTo(this.base, keyframes, options);
        this.base.hidden = true;
        this.emit("sl-after-hide");
      }
    }
    handleDurationChange() {
      this.restartAutoHide();
    }
    /** Shows the alert. */
    async show() {
      if (this.open) {
        return void 0;
      }
      this.open = true;
      return waitForEvent(this, "sl-after-show");
    }
    /** Hides the alert */
    async hide() {
      if (!this.open) {
        return void 0;
      }
      this.open = false;
      return waitForEvent(this, "sl-after-hide");
    }
    /**
     * Displays the alert as a toast notification. This will move the alert out of its position in the DOM and, when
     * dismissed, it will be removed from the DOM completely. By storing a reference to the alert, you can reuse it by
     * calling this method again. The returned promise will resolve after the alert is hidden.
     */
    async toast() {
      return new Promise((resolve) => {
        if (toastStack.parentElement === null) {
          document.body.append(toastStack);
        }
        toastStack.appendChild(this);
        requestAnimationFrame(() => {
          this.clientWidth;
          this.show();
        });
        this.addEventListener(
          "sl-after-hide",
          () => {
            toastStack.removeChild(this);
            resolve();
            if (toastStack.querySelector("sl-alert") === null) {
              toastStack.remove();
            }
          },
          { once: true }
        );
      });
    }
    render() {
      return x`
      <div
        part="base"
        class=${o({
      alert: true,
      "alert--open": this.open,
      "alert--closable": this.closable,
      "alert--has-icon": this.hasSlotController.test("icon"),
      "alert--primary": this.variant === "primary",
      "alert--success": this.variant === "success",
      "alert--neutral": this.variant === "neutral",
      "alert--warning": this.variant === "warning",
      "alert--danger": this.variant === "danger"
    })}
        role="alert"
        aria-hidden=${this.open ? "false" : "true"}
        @mousemove=${this.handleMouseMove}
      >
        <slot name="icon" part="icon" class="alert__icon"></slot>

        <slot part="message" class="alert__message" aria-live="polite"></slot>

        ${this.closable ? x`
              <sl-icon-button
                part="close-button"
                exportparts="base:close-button__base"
                class="alert__close-button"
                name="x-lg"
                library="system"
                label=${this.localize.term("close")}
                @click=${this.handleCloseClick}
              ></sl-icon-button>
            ` : ""}
      </div>
    `;
    }
  };
  SlAlert.styles = alert_styles_default;
  __decorateClass([
    i$4('[part~="base"]')
  ], SlAlert.prototype, "base", 2);
  __decorateClass([
    n$6({ type: Boolean, reflect: true })
  ], SlAlert.prototype, "open", 2);
  __decorateClass([
    n$6({ type: Boolean, reflect: true })
  ], SlAlert.prototype, "closable", 2);
  __decorateClass([
    n$6({ reflect: true })
  ], SlAlert.prototype, "variant", 2);
  __decorateClass([
    n$6({ type: Number })
  ], SlAlert.prototype, "duration", 2);
  __decorateClass([
    watch("open", { waitUntilFirstUpdate: true })
  ], SlAlert.prototype, "handleOpenChange", 1);
  __decorateClass([
    watch("duration")
  ], SlAlert.prototype, "handleDurationChange", 1);
  SlAlert = __decorateClass([
    e$6("sl-alert")
  ], SlAlert);
  setDefaultAnimation("alert.show", {
    keyframes: [
      { opacity: 0, scale: 0.8 },
      { opacity: 1, scale: 1 }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("alert.hide", {
    keyframes: [
      { opacity: 1, scale: 1 },
      { opacity: 0, scale: 0.8 }
    ],
    options: { duration: 250, easing: "ease" }
  });

  const locks = new Set();
  function getScrollbarWidth() {
      const documentWidth = document.documentElement.clientWidth;
      return Math.abs(window.innerWidth - documentWidth);
  }
  function lockBodyScrolling(lockingEl) {
      locks.add(lockingEl);
      if (!document.body.classList.contains('sl-scroll-lock')) {
          const scrollbarWidth = getScrollbarWidth();
          document.body.classList.add('sl-scroll-lock');
          document.body.style.setProperty('--sl-scroll-lock-size', `${scrollbarWidth}px`);
      }
  }
  function unlockBodyScrolling(lockingEl) {
      locks.delete(lockingEl);
      if (locks.size === 0) {
          document.body.classList.remove('sl-scroll-lock');
          document.body.style.removeProperty('--sl-scroll-lock-size');
      }
  }

  function uppercaseFirstLetter(string) {
      return string.charAt(0).toUpperCase() + string.slice(1);
  }

  function isTabbable(el) {
      const tag = el.tagName.toLowerCase();
      if (el.getAttribute('tabindex') === '-1') {
          return false;
      }
      if (el.hasAttribute('disabled')) {
          return false;
      }
      if (el.hasAttribute('aria-disabled') && el.getAttribute('aria-disabled') !== 'false') {
          return false;
      }
      if (tag === 'input' && el.getAttribute('type') === 'radio' && !el.hasAttribute('checked')) {
          return false;
      }
      if (el.offsetParent === null) {
          return false;
      }
      if (window.getComputedStyle(el).visibility === 'hidden') {
          return false;
      }
      if ((tag === 'audio' || tag === 'video') && el.hasAttribute('controls')) {
          return true;
      }
      if (el.hasAttribute('tabindex')) {
          return true;
      }
      if (el.hasAttribute('contenteditable') && el.getAttribute('contenteditable') !== 'false') {
          return true;
      }
      return ['button', 'input', 'select', 'textarea', 'a', 'audio', 'video', 'summary'].includes(tag);
  }
  function getTabbableBoundary(root) {
      var _a, _b;
      const allElements = [];
      function walk(el) {
          if (el instanceof HTMLElement) {
              allElements.push(el);
              if (el.shadowRoot !== null && el.shadowRoot.mode === 'open') {
                  walk(el.shadowRoot);
              }
          }
          [...el.children].forEach((e) => walk(e));
      }
      walk(root);
      const start = (_a = allElements.find(el => isTabbable(el))) !== null && _a !== void 0 ? _a : null;
      const end = (_b = allElements.reverse().find(el => isTabbable(el))) !== null && _b !== void 0 ? _b : null;
      return { start, end };
  }

  let activeModals = [];
  class Modal {
      constructor(element) {
          this.tabDirection = 'forward';
          this.element = element;
          this.handleFocusIn = this.handleFocusIn.bind(this);
          this.handleKeyDown = this.handleKeyDown.bind(this);
          this.handleKeyUp = this.handleKeyUp.bind(this);
      }
      activate() {
          activeModals.push(this.element);
          document.addEventListener('focusin', this.handleFocusIn);
          document.addEventListener('keydown', this.handleKeyDown);
          document.addEventListener('keyup', this.handleKeyUp);
      }
      deactivate() {
          activeModals = activeModals.filter(modal => modal !== this.element);
          document.removeEventListener('focusin', this.handleFocusIn);
          document.removeEventListener('keydown', this.handleKeyDown);
          document.removeEventListener('keyup', this.handleKeyUp);
      }
      isActive() {
          return activeModals[activeModals.length - 1] === this.element;
      }
      checkFocus() {
          if (this.isActive()) {
              if (!this.element.matches(':focus-within')) {
                  const { start, end } = getTabbableBoundary(this.element);
                  const target = this.tabDirection === 'forward' ? start : end;
                  if (typeof (target === null || target === void 0 ? void 0 : target.focus) === 'function') {
                      target.focus({ preventScroll: true });
                  }
              }
          }
      }
      handleFocusIn() {
          this.checkFocus();
      }
      handleKeyDown(event) {
          if (event.key === 'Tab' && event.shiftKey) {
              this.tabDirection = 'backward';
              requestAnimationFrame(() => this.checkFocus());
          }
      }
      handleKeyUp() {
          this.tabDirection = 'forward';
      }
  }

  var drawer_styles_default = i$3`
  ${componentStyles}

  :host {
    --size: 25rem;
    --header-spacing: var(--sl-spacing-large);
    --body-spacing: var(--sl-spacing-large);
    --footer-spacing: var(--sl-spacing-large);

    display: contents;
  }

  .drawer {
    top: 0;
    inset-inline-start: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    overflow: hidden;
  }

  .drawer--contained {
    position: absolute;
    z-index: initial;
  }

  .drawer--fixed {
    position: fixed;
    z-index: var(--sl-z-index-drawer);
  }

  .drawer__panel {
    position: absolute;
    display: flex;
    flex-direction: column;
    z-index: 2;
    max-width: 100%;
    max-height: 100%;
    background-color: var(--sl-panel-background-color);
    box-shadow: var(--sl-shadow-x-large);
    overflow: auto;
    pointer-events: all;
  }

  .drawer__panel:focus {
    outline: none;
  }

  .drawer--top .drawer__panel {
    top: 0;
    inset-inline-end: auto;
    bottom: auto;
    inset-inline-start: 0;
    width: 100%;
    height: var(--size);
  }

  .drawer--end .drawer__panel {
    top: 0;
    inset-inline-end: 0;
    bottom: auto;
    inset-inline-start: auto;
    width: var(--size);
    height: 100%;
  }

  .drawer--bottom .drawer__panel {
    top: auto;
    inset-inline-end: auto;
    bottom: 0;
    inset-inline-start: 0;
    width: 100%;
    height: var(--size);
  }

  .drawer--start .drawer__panel {
    top: 0;
    inset-inline-end: auto;
    bottom: auto;
    inset-inline-start: 0;
    width: var(--size);
    height: 100%;
  }

  .drawer__header {
    display: flex;
  }

  .drawer__title {
    flex: 1 1 auto;
    font: inherit;
    font-size: var(--sl-font-size-large);
    line-height: var(--sl-line-height-dense);
    padding: var(--header-spacing);
    margin: 0;
  }

  .drawer__header-actions {
    flex-shrink: 0;
    display: flex;
    flex-wrap: wrap;
    justify-content: end;
    gap: var(--sl-spacing-2x-small);
    padding: 0 var(--header-spacing);
  }

  .drawer__header-actions sl-icon-button,
  .drawer__header-actions ::slotted(sl-icon-button) {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    font-size: var(--sl-font-size-medium);
  }

  .drawer__body {
    flex: 1 1 auto;
    display: block;
    padding: var(--body-spacing);
    overflow: auto;
    -webkit-overflow-scrolling: touch;
  }

  .drawer__footer {
    text-align: right;
    padding: var(--footer-spacing);
  }

  .drawer__footer ::slotted(sl-button:not(:last-of-type)) {
    margin-inline-end: var(--sl-spacing-x-small);
  }

  .drawer:not(.drawer--has-footer) .drawer__footer {
    display: none;
  }

  .drawer__overlay {
    display: block;
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background-color: var(--sl-overlay-background-color);
    pointer-events: all;
  }

  .drawer--contained .drawer__overlay {
    display: none;
  }

  @media (forced-colors: active) {
    .drawer__panel {
      border: solid 1px var(--sl-color-neutral-0);
    }
  }
`;

  let SlDrawer = class extends ShoelaceElement {
    constructor() {
      super(...arguments);
      this.hasSlotController = new HasSlotController(this, "footer");
      this.localize = new LocalizeController(this);
      this.modal = new Modal(this);
      this.open = false;
      this.label = "";
      this.placement = "end";
      this.contained = false;
      this.noHeader = false;
      this.handleDocumentKeyDown = (event) => {
        if (this.open && !this.contained && event.key === "Escape") {
          event.stopPropagation();
          this.requestClose("keyboard");
        }
      };
    }
    firstUpdated() {
      this.drawer.hidden = !this.open;
      if (this.open) {
        this.addOpenListeners();
        if (!this.contained) {
          this.modal.activate();
          lockBodyScrolling(this);
        }
      }
    }
    disconnectedCallback() {
      super.disconnectedCallback();
      unlockBodyScrolling(this);
    }
    requestClose(source) {
      const slRequestClose = this.emit("sl-request-close", {
        cancelable: true,
        detail: { source }
      });
      if (slRequestClose.defaultPrevented) {
        const animation = getAnimation(this, "drawer.denyClose", { dir: this.localize.dir() });
        animateTo(this.panel, animation.keyframes, animation.options);
        return;
      }
      this.hide();
    }
    addOpenListeners() {
      document.addEventListener("keydown", this.handleDocumentKeyDown);
    }
    removeOpenListeners() {
      document.removeEventListener("keydown", this.handleDocumentKeyDown);
    }
    async handleOpenChange() {
      if (this.open) {
        this.emit("sl-show");
        this.addOpenListeners();
        this.originalTrigger = document.activeElement;
        if (!this.contained) {
          this.modal.activate();
          lockBodyScrolling(this);
        }
        const autoFocusTarget = this.querySelector("[autofocus]");
        if (autoFocusTarget) {
          autoFocusTarget.removeAttribute("autofocus");
        }
        await Promise.all([stopAnimations(this.drawer), stopAnimations(this.overlay)]);
        this.drawer.hidden = false;
        requestAnimationFrame(() => {
          const slInitialFocus = this.emit("sl-initial-focus", { cancelable: true });
          if (!slInitialFocus.defaultPrevented) {
            if (autoFocusTarget) {
              autoFocusTarget.focus({ preventScroll: true });
            } else {
              this.panel.focus({ preventScroll: true });
            }
          }
          if (autoFocusTarget) {
            autoFocusTarget.setAttribute("autofocus", "");
          }
        });
        const panelAnimation = getAnimation(this, `drawer.show${uppercaseFirstLetter(this.placement)}`, {
          dir: this.localize.dir()
        });
        const overlayAnimation = getAnimation(this, "drawer.overlay.show", { dir: this.localize.dir() });
        await Promise.all([
          animateTo(this.panel, panelAnimation.keyframes, panelAnimation.options),
          animateTo(this.overlay, overlayAnimation.keyframes, overlayAnimation.options)
        ]);
        this.emit("sl-after-show");
      } else {
        this.emit("sl-hide");
        this.removeOpenListeners();
        if (!this.contained) {
          this.modal.deactivate();
          unlockBodyScrolling(this);
        }
        await Promise.all([stopAnimations(this.drawer), stopAnimations(this.overlay)]);
        const panelAnimation = getAnimation(this, `drawer.hide${uppercaseFirstLetter(this.placement)}`, {
          dir: this.localize.dir()
        });
        const overlayAnimation = getAnimation(this, "drawer.overlay.hide", { dir: this.localize.dir() });
        await Promise.all([
          animateTo(this.overlay, overlayAnimation.keyframes, overlayAnimation.options).then(() => {
            this.overlay.hidden = true;
          }),
          animateTo(this.panel, panelAnimation.keyframes, panelAnimation.options).then(() => {
            this.panel.hidden = true;
          })
        ]);
        this.drawer.hidden = true;
        this.overlay.hidden = false;
        this.panel.hidden = false;
        const trigger = this.originalTrigger;
        if (typeof (trigger == null ? void 0 : trigger.focus) === "function") {
          setTimeout(() => trigger.focus());
        }
        this.emit("sl-after-hide");
      }
    }
    handleNoModalChange() {
      if (this.open && !this.contained) {
        this.modal.activate();
        lockBodyScrolling(this);
      }
      if (this.open && this.contained) {
        this.modal.deactivate();
        unlockBodyScrolling(this);
      }
    }
    /** Shows the drawer. */
    async show() {
      if (this.open) {
        return void 0;
      }
      this.open = true;
      return waitForEvent(this, "sl-after-show");
    }
    /** Hides the drawer */
    async hide() {
      if (!this.open) {
        return void 0;
      }
      this.open = false;
      return waitForEvent(this, "sl-after-hide");
    }
    render() {
      return x`
      <div
        part="base"
        class=${o({
      drawer: true,
      "drawer--open": this.open,
      "drawer--top": this.placement === "top",
      "drawer--end": this.placement === "end",
      "drawer--bottom": this.placement === "bottom",
      "drawer--start": this.placement === "start",
      "drawer--contained": this.contained,
      "drawer--fixed": !this.contained,
      "drawer--rtl": this.localize.dir() === "rtl",
      "drawer--has-footer": this.hasSlotController.test("footer")
    })}
      >
        <div part="overlay" class="drawer__overlay" @click=${() => this.requestClose("overlay")} tabindex="-1"></div>

        <div
          part="panel"
          class="drawer__panel"
          role="dialog"
          aria-modal="true"
          aria-hidden=${this.open ? "false" : "true"}
          aria-label=${l(this.noHeader ? this.label : void 0)}
          aria-labelledby=${l(!this.noHeader ? "title" : void 0)}
          tabindex="0"
        >
          ${!this.noHeader ? x`
                <header part="header" class="drawer__header">
                  <h2 part="title" class="drawer__title" id="title">
                    <!-- If there's no label, use an invisible character to prevent the header from collapsing -->
                    <slot name="label"> ${this.label.length > 0 ? this.label : String.fromCharCode(65279)} </slot>
                  </h2>
                  <div part="header-actions" class="drawer__header-actions">
                    <slot name="header-actions"></slot>
                    <sl-icon-button
                      part="close-button"
                      exportparts="base:close-button__base"
                      class="drawer__close"
                      name="x-lg"
                      label=${this.localize.term("close")}
                      library="system"
                      @click=${() => this.requestClose("close-button")}
                    ></sl-icon-button>
                  </div>
                </header>
              ` : ""}

          <slot part="body" class="drawer__body"></slot>

          <footer part="footer" class="drawer__footer">
            <slot name="footer"></slot>
          </footer>
        </div>
      </div>
    `;
    }
  };
  SlDrawer.styles = drawer_styles_default;
  __decorateClass([
    i$4(".drawer")
  ], SlDrawer.prototype, "drawer", 2);
  __decorateClass([
    i$4(".drawer__panel")
  ], SlDrawer.prototype, "panel", 2);
  __decorateClass([
    i$4(".drawer__overlay")
  ], SlDrawer.prototype, "overlay", 2);
  __decorateClass([
    n$6({ type: Boolean, reflect: true })
  ], SlDrawer.prototype, "open", 2);
  __decorateClass([
    n$6({ reflect: true })
  ], SlDrawer.prototype, "label", 2);
  __decorateClass([
    n$6({ reflect: true })
  ], SlDrawer.prototype, "placement", 2);
  __decorateClass([
    n$6({ type: Boolean, reflect: true })
  ], SlDrawer.prototype, "contained", 2);
  __decorateClass([
    n$6({ attribute: "no-header", type: Boolean, reflect: true })
  ], SlDrawer.prototype, "noHeader", 2);
  __decorateClass([
    watch("open", { waitUntilFirstUpdate: true })
  ], SlDrawer.prototype, "handleOpenChange", 1);
  __decorateClass([
    watch("contained", { waitUntilFirstUpdate: true })
  ], SlDrawer.prototype, "handleNoModalChange", 1);
  SlDrawer = __decorateClass([
    e$6("sl-drawer")
  ], SlDrawer);
  setDefaultAnimation("drawer.showTop", {
    keyframes: [
      { opacity: 0, translate: "0 -100%" },
      { opacity: 1, translate: "0 0" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.hideTop", {
    keyframes: [
      { opacity: 1, translate: "0 0" },
      { opacity: 0, translate: "0 -100%" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.showEnd", {
    keyframes: [
      { opacity: 0, translate: "100%" },
      { opacity: 1, translate: "0" }
    ],
    rtlKeyframes: [
      { opacity: 0, translate: "-100%" },
      { opacity: 1, translate: "0" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.hideEnd", {
    keyframes: [
      { opacity: 1, translate: "0" },
      { opacity: 0, translate: "100%" }
    ],
    rtlKeyframes: [
      { opacity: 1, translate: "0" },
      { opacity: 0, translate: "-100%" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.showBottom", {
    keyframes: [
      { opacity: 0, translate: "0 100%" },
      { opacity: 1, translate: "0 0" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.hideBottom", {
    keyframes: [
      { opacity: 1, translate: "0 0" },
      { opacity: 0, translate: "0 100%" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.showStart", {
    keyframes: [
      { opacity: 0, translate: "-100%" },
      { opacity: 1, translate: "0" }
    ],
    rtlKeyframes: [
      { opacity: 0, translate: "100%" },
      { opacity: 1, translate: "0" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.hideStart", {
    keyframes: [
      { opacity: 1, translate: "0" },
      { opacity: 0, translate: "-100%" }
    ],
    rtlKeyframes: [
      { opacity: 1, translate: "0" },
      { opacity: 0, translate: "100%" }
    ],
    options: { duration: 250, easing: "ease" }
  });
  setDefaultAnimation("drawer.denyClose", {
    keyframes: [{ scale: 1 }, { scale: 1.01 }, { scale: 1 }],
    options: { duration: 250 }
  });
  setDefaultAnimation("drawer.overlay.show", {
    keyframes: [{ opacity: 0 }, { opacity: 1 }],
    options: { duration: 250 }
  });
  setDefaultAnimation("drawer.overlay.hide", {
    keyframes: [{ opacity: 1 }, { opacity: 0 }],
    options: { duration: 250 }
  });

  // import '@shoelace-style/shoelace/dist/themes/light.css';

  const $callback = new Callback();
  const {
    router,
    client
  } = new application_1();
  router.verb('GET', '/readiness', async (ctx, next) => {
    ctx.body = {
      code: 0,
      message: 'ok'
    };
    next();
  });
  function callBackstage(call) {
    return $callback(id => {
      const message = {
        source: SOURCE,
        type: 'CALL',
        from: 'FRONTSTAGE',
        to: 'BACKSTAGE',
        call,
        callback: id
      };
      window.postMessage(message);
    });
  }
  function callbackBackstage(callback, resolved, rejected) {
    const message = {
      source: SOURCE,
      type: 'CALLBACK',
      from: 'FRONTSTAGE',
      to: 'BACKSTAGE',
      call: {
        resolved,
        rejected
      },
      callback
    };
    window.postMessage(message);
  }
  window.addEventListener('message', async _ref => {
    let {
      /* type, source, origin, */data
    } = _ref;
    if (data.source !== SOURCE) {
      return false;
    }
    const {
      from,
      to,
      type,
      call,
      callback
    } = data;
    if (type === 'CALLBACK' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      console.debug("%c%s", "background-color: #f0f9ff", '[FRONTSTAGE] window.addEventListener.message', from, to, type);
      $callback(callback, call.resolved, call.rejected);
    } else if (type === 'CALL' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      console.debug("%c%s", "background-color: #f0f9ff", '[FRONTSTAGE] window.addEventListener.message', from, to, type);
      // const resolved = { code: 0, message: 'call : backstage -> frontstage : ok' };

      try {
        const resolved = await client.verb(call.method, call.ep, call.search, call.form);
        callbackBackstage(callback, resolved, null);
      } catch (rejected) {
        callbackBackstage(callback, null, rejected);
      }
    } else {
      console.debug("%c%s", "color: #727272", '[FRONTSTAGE] window.addEventListener.message', data);
    }
  });
  const backstage = window.backstage || {};
  window.backstage = backstage;
  backstage.invoke = function invoke(svc, method, ep) {
    let search = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
    let form = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : {};
    const call = {
      svc,
      method,
      ep,
      search,
      form
    };
    return callBackstage(call);
  };
  backstage.route = router.verb.bind(router);
  backstage.verb = client.verb.bind(client);
  backstage.toast = function toast(message) {
    let variant = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 'primary';
    let duration = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 2000;
    // console.log('@head/backstage toast');

    function escapeHtml(html) {
      const div = document.createElement('div');
      div.textContent = html;
      return div.innerHTML;
    }
    const alert = Object.assign(document.createElement('sl-alert'), {
      variant,
      closable: true,
      duration,
      innerHTML: "\n      ".concat(escapeHtml(message), "\n    ")
    });
    document.body.append(alert);
    return alert.toast();
  };
  backstage.drawer = function drawer_(html) {
    let width = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : '50';
    // console.log('@head/backstage drawer');

    const drawer = document.createElement('sl-drawer');
    drawer.label = 'Backstage';
    drawer.style.setProperty('--size', "".concat(width, "vw"));
    if (typeof html === 'string') {
      drawer.innerHTML = "\n      ".concat(html, "\n    ");
      // } else if (content instanceof HTMLElement) {
      //   drawer.innerHTML = ``;
      //   drawer.appendChild(content);
    }
    document.body.appendChild(drawer);
    drawer.show();
    drawer.addEventListener('sl-after-hide', () => {
      drawer.remove();
    });
  };
  window.postMessage({
    source: SOURCE,
    type: 'READY',
    from: 'FRONTSTAGE'
  });

})();
