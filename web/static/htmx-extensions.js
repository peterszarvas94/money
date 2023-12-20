let accessToken = "";

/**
 * Refresh the access token.
 * First, it checks if the access token is already set.
 * If not, it tries to get an access token.
 * If successful, it reloads the page with the new access token.
 * If something goes wrong, it calls the callback.
 * @param {function | undefined} cb - Callback if something goes wrong in client side.
 * Backend handle happy path and errors.
 */
async function refresh(cb) {
  function callback() {
    if (cb !== undefined) {
      return cb();
    }
  }

  if (accessToken !== "") {
    return;
  }

  let authHeader = "";
  try {
    // attaches refreshtoken cookie
    const res = await fetch("/refresh", { method: "POST" });
    if (!res.ok) {
      // 500
      return callback();
    }
    // "Bearer <accesstoken>"
    authHeader = res.headers.get("Authorization");
  } catch {
    // cant refresh
    return callback();
  }

  if (authHeader === "") {
    // no auth header
    return callback();
  }

  accessToken = authHeader.split(" ")[1];

  try {
    // reload page with new access token
    await htmx.ajax("GET", window.location.pathname, {
      headers: {
        Authorization: authHeader,
      },
      target: "body",
    });
  } catch {
    // cant reload
    return callback();
  }
}

/**
 * Add current route to history and redirect to route.
 * @param {string} route - Route to redirect to
 */
function redirect(route) {
  window.history.pushState({}, "", route);
  window.location.reload();
}

/**
 * "protected-page" extension to redirect to login page if not logged in.
 * 1. Backend send spinner.
 * 2. Frontend tries to get access token, reloads the page with new access token.
 * 3. If unsuccessful, frontend redirect to login page.
 * 4. If successful, backend sends down protected content.
 * e.g.
 * /dashboard
 * <div id="page" hx-ext="protected">
 */
htmx.defineExtension("protected-page", {
  onEvent: async function (name, evt) {
    const processed = name === "htmx:afterProcessNode";
    const pageElement = evt?.target?.id === "page";
    if (processed && pageElement) {
      const path = window.location.pathname;
      const encodedPath = encodeURIComponent(path);
      await refresh(() => {
        redirect("/signin?redirect=" + encodedPath);
      });
    }
  },
});

/**
 * "gated-page" extension to load more/different content if logged in.
 * 1. Backend send partial page.
 * 2. Frontend tries to get access token, reloads the page with new access token.
 * 3. If unsuccessful, frontend does nothing.
 * 4. If successful, backend sends down gated content.
 * e.g.
 * <body hx-ext="gated">
 */
htmx.defineExtension("gated-page", {
  onEvent: async function (name, evt) {
    const processed = name === "htmx:afterProcessNode";
    const pageElement = evt?.target?.id === "page";
    if (processed && pageElement) {
      await refresh();
    }
  },
});

/**
 * "show-client-error" extension to render client error
 * only used on form submission like signin, signup
 */
htmx.defineExtension("show-client-error", {
  onEvent: function (name, evt) {
    if (name == "htmx:beforeOnLoad") {
      const status = evt?.detail?.xhr?.status;
      const statusClientError =
        status === 401 || status === 403 || status === 409;
      if (statusClientError) {
        evt.detail.shouldSwap = true;
      }
    }
  },
});

/**
 * "show-notfound" extension to render not found page
 * only used on page which can throw 404
 */
htmx.defineExtension("show-notfound", {
  transformResponse: function (text, xhr, _elt) {
    console.log(xhr.status);
    if (xhr.status === 404) {
      return text;
    }
  },
  // onEvent: function (name, evt) {
  //   console.log(1);
  //   if (name == "htmx:beforeOnLoad") {
  //     const status = evt?.detail?.xhr?.status;
  //     const statusNotFound = status === 404;
  //     if (statusNotFound) {
  //       evt.detail.shouldSwap = true;
  //     }
  //   }
  // },
});

/**
 * "delete" extension to delete elements
 * only used on delete button
 */
htmx.defineExtension("delete", {
  onEvent: function (_name, evt) {
    const status = evt?.detail?.xhr?.status;
    const statusNoContent = status === 204;
    if (statusNoContent) {
      evt.detail.shouldSwap = true;
    }
  },
});

/**
 * "no-server-error" extension to prevent server error
 * probably should make it global
 */
htmx.defineExtension("no-server-error", {
  onEvent: function (_name, evt) {
    const status = evt?.detail?.xhr?.status;
    const statusServerError = status === 500;
    if (statusServerError) {
      evt.detail.shouldSwap = false;
    }
  },
});

/**
 * "access" extension to set access token to header
 * use it in any request that requires access token
 * e.g. GET /dashboard
 */
htmx.defineExtension("access", {
  onEvent: function (name, evt) {
    if (name === "htmx:beforeRequest") {
      if (accessToken) {
        evt.detail.xhr.setRequestHeader(
          "Authorization",
          "Bearer " + accessToken,
        );
      }
    }
  },
});

/**
 * "signout" extension to clear access token
 * use it in signout button
 */
htmx.defineExtension("signout", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterRequest") {
      if (evt.detail.xhr.status == 200) {
        accessToken = "";
      }
    }
  },
});

/**
 * "signin" extension to set access token from server to variable
 * use it in signin form
 * probably can avoid, body extensions can do the same (gated-page, protected-page)
 */
htmx.defineExtension("signin", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterRequest" && evt.detail.xhr.status == 200) {
      auth = evt.detail.xhr.getResponseHeader("Authorization");
      if (auth) {
        accessToken = auth.split(" ")[1];
      }
    }
  },
});

/**
 * "description" extension to reload meta tag
 * use this on body
 */
htmx.defineExtension("description", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterSwap") {
      const meta = document.querySelector("meta[name=description]");
      const res = evt?.detail?.xhr?.responseText;
      const parser = new DOMParser();
      const doc = parser.parseFromString(res, "text/html");
      const newMeta = doc.querySelector("meta[name=description]");
      const newContent = newMeta?.getAttribute("content");
      meta.setAttribute("content", newContent);
    }
  },
});

/**
 * "receiver" extension to receive event from body
 * usage:
 * <button hx-on="click:emit('event')">Click me</button>
 * <div
 *   hx-ext="receiver"
 *   on-event:open="this.classList.remove('hidden')"
 *   on-event:close="this.classList.add('hidden')"
 * >
 *   I was clicked
 * </div>
 * */
htmx.defineExtension("receiver", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterProcessNode") {
      const e = evt?.target;
      const attributes = e.attributes;
      for (a of attributes) {
        if (a.name.startsWith("on-event:")) {
          const event = a.name.split(":")[1];
          const action = a.value;
          document.body.addEventListener(event, () => {
            if (document.contains(e)) {
              // evals action in the context of e
              new Function(action).call(e);
            }
          });
        }
      }
    }
  },
});

/**
 * Emits an event on body.
 * @param {string} name - Event name
 */
function emit(name) {
  document.body.dispatchEvent(new CustomEvent(name));
}

/**
 * Toggles tabindex of element.
 * @param {HTMLElement} element - Element to toggle tabindex
 */
function toggleTabIndex(element) {
  if (element.tabIndex === -1) {
    element.tabIndex = 0;
  } else {
    element.tabIndex = -1;
  }
}

/**
 * Toggles focus and blur on element.
 * @param {HTMLElement} element - Element to toggle focus
 */
function toggleFocus(element) {
  if (element === document.activeElement) {
    element.blur();
  } else {
    element.focus();
  }
}
