<!-- markdownlint-disable MD002 MD013 MD033 MD041 -->
<h1 align="center">
  <a name="logo" href="https://github.com/vallieres/crawl-n-indexnow">
    <img src="https://github.com/vallieres/crawl-n-indexnow/assets/182217/01bf723d-b7c9-476e-8271-817e5ee3a4d2" alt="Crawl n' IndexNow" width="200"></a>
  <br>
  Crawl n' IndexNow
</h1>
<h4 align="center">Get the goods and ship 'em to the indexes!</h4>
<div align="center"></div>

<font size="3">
    <strong>Crawl n' IndexNow</strong> is a simple CLI that pulls your Shopify or regular website sitemap, 
        and submits the URLs to IndexNow to speed up the indexing process.
</font>


## ⚡️ Quickstart

Install the CLI:
```bash
go install github.com/vallieres/crawl-n-indexnow@latest
```

## 🎯 Commands

### Submit All Shopify URLs to Index

```bash
crawl-n-indexnow shopify --domain legacygoods.co --key a1b3c34d
```

This will crawl the [Legacy Goods Co.](https://legacygoods.co) domain, list the found URLs from its sitemaps and submit 
them all to IndexNow.

### Submit All Sitemap URLs to Index

```bash
crawl-n-indexnow sitemap --domain hipsterpixel.co --key a1b3c34d
```

This will crawl the [Hipster Pixel](https://hipsterpixel.co) domain, list the found URLs from its sitemap and submit
them all to IndexNow.