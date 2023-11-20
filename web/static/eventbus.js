// event bus 
// usage:
//
// emit function emits a custom event, here its called 'clicked'
// with the help of htmx, we listen to click event, and call emit function
// <button hx-on="click:emit('clicked')">Click me</button>
// <div hx-on="clicked:this.classList.add('clicked')">I was clicked</div>

/**
 * Emit a custom event
 * @param {string} name - event name
 */
function emit(name) {
  document.body.dispatchEvent(new CustomEvent(name));
}
