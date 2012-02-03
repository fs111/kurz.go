var MONTHNAMES = new Array("Jan", "Feb", "Mar", "Apr", "May", "Jun", 
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec");

function formatDate(d){
    return d.getDate() + " " + MONTHNAMES[d.getMonth()] + " "
        + d.getFullYear();
}

function formatURL(url){
    var clean = url.replace("http://", "");
    clean = clean.replace("https://", "");
    clean = clean.substr(0, 52);
    return "<a href=\"" + url +"\">" + clean + "</a>";
}

function loadKurls(howmany){

    $('#data tr:not(:first)').remove();
    $.getJSON("/latest/" + howmany, function( obj ) { 
        var allUrls = obj["urls"];
        for (var i = 0; i < allUrls.length; i++) {
            var kurl = allUrls[i];

            var d = new Date(kurl["CreationDate"] / 1000000);
            $("#data").addRow({ 
                newRow: "<tr>" 
                  + "<td class=\"short\">" + formatURL(kurl["ShortUrl"]) + "</td>"
                  + "<td class=\"long\">" + formatURL(kurl["LongUrl"]) + "</td>"
                  + "<td class=\"date\">" + formatDate(d) + "</td>"
                  + "<td class=\"clicks\">" + kurl["Clicks"] + "</td>"
                  + "</tr>",
                rowSpeed: 700
            
            });
        } 
    });
}

