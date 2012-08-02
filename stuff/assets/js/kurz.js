var MONTHNAMES = new Array("Jan", "Feb", "Mar", "Apr", "May", "Jun", 
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec");

function formatDate(d){
    return d.getDate() + " " + MONTHNAMES[d.getMonth()] + " "
        + d.getFullYear();
}

function formatURL(url){
    var clean = url.replace("http://", "");
    clean = clean.replace("https://", "");
    clean = clean.substr(0, 50);
    return "<a href=\"" + url +"\">" + clean + "</a>";
}

function createTweet(url){
    return "<a target=\"blank\" href=\"https://twitter.com/intent/tweet?text="
              +  encodeURIComponent(url) + "\">"
              + "<img src=\"img/tw.png\" alt=\"tw\"></a>";
}

function loadKurls(howmany){

    $('#data tr:not(:first)').remove();
    $.getJSON("/latest/" + howmany, function( allUrls ) { 
        for (var i = 0; i < allUrls.length; i++) {
            var kurl = allUrls[i];
            var d = new Date(kurl["CreationDate"] / 1000000);

            $("#data").addRow({ 
                newRow: "<tr>" 
                  + "<td class=\"short\">" + formatURL(kurl["ShortUrl"]) + "</td>"
                  + "<td class=\"long\">" + formatURL(kurl["LongUrl"]) + "</td>"
                  + "<td class=\"date\">" + formatDate(d) + "</td>"
                  + "<td class=\"clicks\">" + kurl["Clicks"] + "</td>"
                  + "<td class=\"tweet\">" + createTweet(kurl["ShortUrl"]) + "</td>"
                  + "</tr>",
                rowSpeed: 700
            
            });
        } 
    });
}
