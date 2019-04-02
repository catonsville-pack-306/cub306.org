$( document ).ready(function()
{
    $("#header_title").click(function()
    {
        document.location.href="/index.md"
    });
    
    $('#header_title').hover(function()
    {
        $(this).css('cursor','pointer');
    });
});

$( document ).ready(function()
{
    $("a").each(function(index)
    {
        if ($(this).has("img").length>0) {return;}
        let href = $(this).attr("href");
        if ( href.startsWith("http://") || href.startsWith("https://") )
        {
            //$(this).addClass("external_link")
            $('<i class="added fas fa-external-link-alt"></i>').appendTo(this);
        }
        else if ( href.startsWith("mailto:") )
        {
            $(this).addClass("mailto");
            $('<i class="added far fa-envelope"></i>').prependTo(this);
        }

        if ( href.endsWith(".pdf") )
        {
            $(this).addClass("pdf");
            $('<i class="added far fa-file-pdf"></i>').prependTo(this);
        }
    });
    
});

/** create a table of contents */
$( document ).ready(function()
{
    toc = "Table of Contents:<ol>"
    $("h2").each(function(index)
    {
        let title = $(this).text();
        if ($.trim(title)!="")
        {
            let hash = title.replace(/[^\w]/ig, "-")
            $("<a name='" + hash + "'></a>").insertBefore(this);
            toc += "<li><a href='#" + hash +"'>" + title + "</a></li>"
        }
    });
    toc += "</ol>"
    $("#toc").append(toc)
    $("a[href='toc']").replaceWith(toc)
});

/*$( document ).ready(function()
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
});*/