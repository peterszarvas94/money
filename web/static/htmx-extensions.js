/**
 * "show-client-error" extension to render client error
 * only used on form submission like signin, signup
 */
htmx.defineExtension("show-client-error", {
  onEvent: function(name, evt) {
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
 * "delete" extension to delete elements
 * only used on delete button
 */
htmx.defineExtension("delete", {
  onEvent: function(_name, evt) {
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
  onEvent: function(_name, evt) {
    const status = evt?.detail?.xhr?.status;
    const statusServerError = status === 500;
    if (statusServerError) {
      evt.detail.shouldSwap = false;
    }
  },
});

/**
 * "description" extension to reload meta tag
 * use this on body
 */
htmx.defineExtension("description", {
  onEvent: function(name, evt) {
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
  onEvent: function(name, evt) {
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
