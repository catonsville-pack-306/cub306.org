# Subscribe to Pack Announcement list #

Use this form to subscribe to the Cub Scout Announcement list
(to read what an announcement list is see [wikipedia](https://en.wikipedia.org/wiki/Electronic_mailing_list)).

<form method="post" action="http://scripts.dreamhost.com/add_list.cgi">
    <input type="hidden" name="list" value="news" />
    <input type="hidden" name="domain" value="cub306.org" />

    <div>
        <label>Name</label>
        <input type="text" name="name" placeholder="your name"/>
    </div>
    <div>
        <label>Your E-mail</label>
        <input type="text" name="email" size="32" placeholder="your@email.com"/>
    </div>
    
    <div>
        <label>Confirm E-mail</label>
        <input type="text" name="address2" size="32" placeholder="your@email.com"/>
    </div>
    
    <div>
        <input type="submit" name="submit" value="Join us"/>
    </div>
</form>

You will be sent an email to confirm your intention to join this mailing list.
You will not be added to the list until you click the confirmation link.

If you have any problem, email us directly at [cubmaster@cub306.org](mailto:cubmaster@cub306.org)

## Email Discussion List ##

You can also subscribe to our [email list](http://lists.cub306.org/listinfo.cgi/talk-cub306.org).
This is a traditional email list server.

<script>
    let searchParams = new URLSearchParams(window.location.search);
    let raw_code = searchParams.get("code");
    let code = Number(raw_code.trim().replace(/[^0-9-]/g, ""));
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
