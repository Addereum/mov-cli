// examples/dummy_scraper.js
// 
// This is an example of a modern mov-cli V5 scraper.
// It uses our native Go bindings to fetch and parse HTML blazingly fast!

function search(query) {
    // 1. Fetch the HTML using Go's native HTTP client
    var searchUrl = "https://html.duckduckgo.com/html/?q=" + encodeURIComponent(query + " movie");
    var htmlContent = __core.fetch(searchUrl);
    
    // 2. Parse the HTML using our GoQuery binding (like jQuery)
    // Let's get all result titles and their URLs
    var titles = __core.parseHTML(htmlContent, ".result__a", "text");
    var urls = __core.parseHTML(htmlContent, ".result__a", "href");
    
    var results = [];
    
    // 3. Combine them into objects
    var limit = titles.length > 10 ? 10 : titles.length;
    for (var i = 0; i < limit; i++) {
        results.push({
            title: titles[i].trim(),
            url: urls[i]
        });
    }
    
    // 4. Return as a JSON string
    return JSON.stringify(results);
}
