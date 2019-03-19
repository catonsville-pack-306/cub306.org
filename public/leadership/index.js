
/** create a table of contents */
$( document ).ready(function()
{
    toc = "Table of Contents:<ol>"
    $("h2").each(function(index)
    {
        let title = $(this).text();
        let hash = title.replace(/[^\w]/ig, "-")
        $("<a name='" + hash + "'></a>").insertBefore(this);
        toc += "<li><a href='#" + hash +"'>" + title + "</a></li>"
    });
    toc += "</ol>"
    $("#toc").append(toc)
});

$( document ).ready(function()
{
    toc = "Table of Contents:<ol>"
    $("h2").each(function(index)
    {
        let title = $(this).text();
        let hash = title.replace(/[^\w]/ig, "-")
        $("<a name='" + hash + "'></a>").insertBefore(this);
        toc += "<li><a href='#" + hash +"'>" + title + "</a></li>"
    });
    toc += "</ol>"
    
    $("a[href='toc']").replaceWith(toc)
});