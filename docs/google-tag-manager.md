# How to integrate Google Tag Manager

## Create Account and setup the container
- Go to https://tagmanager.google.com and Create a new Account
- Create a new container and select web
- Write down the ID of the new container (GTM-XXXXXX)

## Setup Flamingo to use GTM-Tag
- Open the Flamingo context file ({project}/config/context.yml)
- Ensure the following section in the `config` area

```yaml
  gtm.active: true
  gtm.id: "{your GTM-ID here}"
```

## Add Searchperience-Tracker
- Create a new tag and name it "SP-Tracker"
- Choose custom html tag
- Paste and adjust the following snipped
```html
  <script src="{path to your CDN hosted version of searchperience/t.js}"></script>
```
- Set the trigger to all pages

## Track custom events
- Create a new html tag and name it "SP-Result-Tracking"
- Open advanced settings open tag sequencing
- Enable "Fire a tag before SP-Result-Tracking fires" and set it to "SP-Tracker"
- Set the trigger to all pages
```html
<script type="text/javascript">
(function () {
  var spc = window.searchperienceConfig || {};
  var sp = window.searchperience || false;
  if (!spc.isSearchPage || !sp) return false;

  var channels = sp.getChannels()

  function getQuery () {
    return $('input[name="q"]').val();
  }

  function getProductIds () {
    var ids = $('.search-results .product-tile').map(function (i, el) {
      return $(el).data('product-id');
    }).get();

    return ids;
  }

  function getCurrentPage () {
    return parseInt($('.pagination .btn.current').text(), 10) || 0;
  }

  function trackSearchResults () {
    var q = getQuery();
    if (!q) return;

    ____tq.push({
      track: 'searchForProduct',
      query: q,
      results: getProductIds().join()
    });
  }

  // Add tracking for results
  trackSearchResults();
  sp.subscribe(channels.SP_RequestDone, trackSearchResults);

  // Add tracking for clicks on results
  $('.search-results').on('click', 'a', function (ev) {
    ev.preventDefault();
    var $link = $(ev.target);
    var href = $link.attr('href');
    var $p = $link.closest('[data-product-id]');
    var pID = $p.data('product-id');

    var position = getProductIds().indexOf(pID) + 1;
    ____tq.push({
      track: 'clickOnResult',
      query: getQuery(),
      item: pID,
      position: position,
      pagination: getCurrentPage()
    });

    // we need to wait for the tracker to finish the request
    setTimeout(function () {
      window.location.href = href;
    }, 100);

  })

})();

</script>
```
