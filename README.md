# xtream
A project meant to support multi-destination RTMP streaming, in tandem with [nginx-rtmp-module](https://github.com/arut/nginx-rtmp-module)

Very barebones.

---
Current state:
- Has one inbound destination to optionally use (also supports directly pointing to proper app)
- Check that publisher exists in a config file
- Check that the publisher stream key matches the encrypted config data
- If attempting to speak with main app, but has a valid stream key, redirect them to their app
  - Separate apps are used in the module since dynamic multi-push doesn't seem to be supported

```
application xtream {
    live on;
    record off;
    
    # redirects on valid app/stream key - errors otherwise (ex: rtmp://host/other/sk)
    on_publish http://localhost:8000/authz; 
}

application other {
    live on;
    record off;

    on_publish http://localhost:8000/authz; # checks stream key for this app

    push rtmp://<twitch-server-closest>/live/<sk>;
    push rtmp://<yt-server-closest>/live/<sk>;
}
```