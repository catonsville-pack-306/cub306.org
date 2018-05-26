$( document ).ready(function()
{
    $("#header_title").click(function()
    {
        document.location.href="index.md"
    });
    
    $('#header_title').hover(function()
    {
        $(this).css('cursor','pointer');
    });
});