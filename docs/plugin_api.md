# Plugin API Documentation

Welcome to the `mov-cli` V5 plugin engine! Plugins are written in pure JavaScript (ES5.1) and run natively inside our embedded Go VM. 

You **do not** need Python, Node.js, or any external libraries. We provide blazing-fast native Go bindings directly into the JavaScript environment via the `__core` object.

## The `search(query)` Function
Every plugin **must** define a `search(query)` function. This function receives the user's search string and must return a JSON-encoded string representing an array of results.

```javascript
function search(query) {
    var results = [
        { title: "Result 1", url: "https://example.com/watch?v=1" },
        { title: "Result 2", url: "https://example.com/watch?v=2" }
    ];
    // Must return a JSON string!
    return JSON.stringify(results);
}
```

## The `__core` API

### `__core.fetch(url)`
Performs an incredibly fast native HTTP GET request.
- **Parameters**: `url` (string)
- **Returns**: `string` (The HTML/JSON body of the response)

### `__core.parseHTML(html, selector, attribute)`
Parses HTML like jQuery, completely natively in Go.
- **Parameters**:
  - `html` (string): The raw HTML string.
  - `selector` (string): The CSS selector (e.g., `.movie-title`, `a.link`).
  - `attribute` (string): The attribute to extract (e.g., `href`, `src`, or `"text"` to get the inner text).
- **Returns**: `Array` of strings.

### `__core.runCmd(command, ...args)`
Executes a shell command on the host system (useful for `yt-dlp`).
- **Returns**: `string` (The command's stdout).
