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
        let href = $(this).attr("href");
        if ( href.startsWith("http://") || href.startsWith("https://") )
        {
            $('<span>&nbsp;</span>').appendTo(this);
            $('<i style="font-size:0.75em;" class="fas fa-external-link-alt"></i>').appendTo(this);
        }
    });
    
});