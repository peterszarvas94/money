let accessToken = "";

/**
 * Refresh the access token.
 * First, it checks if the access token is already set.
 * If not, it tries to get an access token.
 * If successful, it reloads the page with the new access token.
 * If something goes wrong, it calls the callback.
 * @param {function} callback - Callback if something goes wrong in client side.
 * Backend handle happy path and errors.
 */
async function refresh(callback) {
  // to prevent infinite loop
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

// "protected" extension to redirect to login page if not logged in.
// 1. Backend send spinner.
// 2. Frontend tries to get access token, reloads the page with new access token.
// 3. If unsuccessful, frontend redirect to login page.
// 4. If successful, backend sends down protected content.
// e.g.
// /dashboard
// <body hx-ext="protected">
htmx.defineExtension("protected", {
  onEvent: async function (name, evt) {
    const isLoadEvent = name === "htmx:load" || name === "htmx:afterOnLoad";
    const isBody = evt.target === document.body;
    if (isLoadEvent && isBody) {
      await refresh(() => {
        redirect("/signin");
      });
    }
  },
});

// "gated" extension to load more/different content if logged in.
// 1. Backend send partial page.
// 2. Frontend tries to get access token, reloads the page with new access token.
// 3. If unsuccessful, frontend does nothing.
// 4. If successful, backend sends down gated content.
// e.g.
// <body hx-ext="gated">
htmx.defineExtension("gated", {
  onEvent: async function (name, evt) {
    const isLoadEvent = name === "htmx:load" || name === "htmx:afterOnLoad";
    const isBody = evt.target === document.body;
    if (isLoadEvent && isBody) {
      await refresh(() => {});
    }
  },
});

// "auth" extension to render error messages on forms
htmx.defineExtension("auth", {
  onEvent: function (name, evt) {
    if (name == "htmx:beforeOnLoad") {
      const status = evt.detail.xhr.status;
      const statusUnauthorized = status === 401;
      const statusForbidden = status === 403;
      const statusConflict = status === 409;
      const displayErrorMessage =
        statusUnauthorized || statusForbidden || statusConflict;
      if (displayErrorMessage) {
        evt.detail.shouldSwap = true;
      }
    }
  },
});

// "delete" extension to delete elements
htmx.defineExtension("delete", {
  onEvent: function (_name, evt) {
    const status = evt?.detail?.xhr?.status;
    if (status === 204) {
      evt.detail.shouldSwap = true;
    }
  },
});

// "error" extension to prevent swap on 500 error
htmx.defineExtension("error", {
  onEvent: function (_name, evt) {
    const erroStatus = evt?.detail?.xhr?.status;
    if (erroStatus === 500) {
      evt.detail.shouldSwap = false;
    }
  },
});

// "access" extension to include authorization header
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

// "signout" extension to clear access token
htmx.defineExtension("signout", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterRequest") {
      if (evt.detail.xhr.status == 200) {
        accessToken = "";
      }
    }
  },
});

// "signin" extension to set access token from server to variable
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

// "description" extension to reload meta tags
htmx.defineExtension("description", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterSwap") {
      const meta = document.querySelector("meta[name=description]");
      const res = evt?.detail?.xhr?.responseText;
      const parser = new DOMParser();
      const doc = parser.parseFromString(res, "text/html");
      const newMeta = doc.querySelector("meta[name=description]");
      const newContent = newMeta.getAttribute("content");
      meta.setAttribute("content", newContent);
    }
  },
});

