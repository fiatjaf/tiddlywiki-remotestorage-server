**DISABLED:** Now https://tiddly.alhur.es/ is just a static TiddlyWiki that talks to remoteStorage and you can use privately, not a way to publish tiddlywikis to the world anymore.

---

Visit https://tiddly.alhur.es/ to create and edit your remoteStorage-based TiddlyWiki.
Visit `https://tiddly.alhur.es/<remoteStorage address>/<bucket>/` to browse your TiddlyWiki in readonly-public mode (`<bucket>` defaults to `"main"`).

Serving your TiddlyWiki on readonly-public mode in a custom domain:

  1. CNAME `<yourdomain>` to `tiddly.alhur.es`
  2. TXT `_tiddlywiki.<yourdomain>` to `<remoteStorage address>/<bucket>`
  3. Visit `<yourdomain>` to get your readonly-public TiddlyWiki.

This only works if your tiddlers are under `/[public/]tiddlers/<bucket>/` and there's an `__index__` file with all the skinny tiddlers (must be on `/public` for the readonly-public mode) -- this is automatic if you're saving your tiddlers with https://github.com/fiatjaf/tiddlywiki-remotestorage.
