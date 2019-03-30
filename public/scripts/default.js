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