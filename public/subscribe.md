
# Subscribe to Pack News #

<script>
    let searchParams = new URLSearchParams(window.location.search);
    let raw_code = searchParams.get("code");
    code = Number(raw_code.trim().replace(/[^0-9-]/g, ""));
    if (!isNaN(code) && -5<=code && code <=3)
    {
        switch (code)
        {
            case 1: alert("Subscribed") ; break;
            case 2: alert("Unsubscribed") ; break;
            case 3: alert("Mailed Confirmation Link"); break;
        }
    }
</script>

<form method="post" action="http://scripts.dreamhost.com/add_list.cgi">
    <input type="hidden" name="list" value="news" />
    <input type="hidden" name="domain" value="cub306.org" />
    <!--input type="hidden" name="unsuburl" value="http://cub306.local/unsubscribe.md" /-->

    <div>
        <label>Your E-mail</label>
        <input type="email" name="email" size="32" placeholder="your@email.com"/>
    </div>
    
    <div>
        <label>Confirm E-mail</label>
        <input type="email" name="address2" size="32" placeholder="your@email.com"/>
    </div>

    <div><input type="submit" name="submit" value="Join" /></div>
</form>

You will be sent an email to confirm your intention to join this mailing list.
You will not be added to the list until you click the confirmation link.