// plugins/youtube.js

function search(query) {
    // Search for the top 10 results
    var rawOutput = __core.runCmd("yt-dlp", "ytsearch10:" + query, "--dump-json", "--default-search", "ytsearch", "--no-playlist", "--js-runtimes", "node", "--remote-components", "ejs:github");
    
    if (!rawOutput || rawOutput.trim() === "") {
        throw new Error("No output from yt-dlp or error occurred.");
    }

    var lines = rawOutput.trim().split("\n");
    var results = [];

    for (var i = 0; i < lines.length; i++) {
        var line = lines[i].trim();
        if (line.length === 0) continue;
        
        try {
            var data = JSON.parse(line);
            if (data.id && data.title) {
                results.push({
                    title: data.title + " ~ " + (data.channel || "Unknown"),
                    url: data.webpage_url || ("https://youtube.com/watch?v=" + data.id)
                });
            }
        } catch (e) {
            // ignore non-json lines like Warnings
        }
    }
    
    return JSON.stringify(results);
}
