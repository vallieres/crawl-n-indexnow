<!-- markdownlint-disable MD002 MD013 MD033 MD041 -->
<h1 align="center">
  <a name="logo" href="https://github.com/vallieres/crawl-n-index">
    <img src="https://github.com/vallieres/crawl-n-index/assets/182217/01bf723d-b7c9-476e-8271-817e5ee3a4d2" alt="Crawl n' Index" width="200"></a>
  <br>
  Crawl n' Index
</h1>
<h4 align="center">Get the goods and ship 'em to the indexes!</h4>
<div align="center"></div>

<font size="3">
    <strong>Crawl n' Index</strong> is a simple CLI that pulls your Shopify site's URL, and submits them to various indexes to speed up the indexing process.
</font>


## ⚡️ Quickstart

Install the CLI:
```bash
go install github.com/vallieres/crawl-n-index@latest
```

## 🎯 Commands

### Submit All URLs to Index

```bash
crawl-n-index all --domain legacygoods.co --key a1b3c34d
```

This will crawl the domain, list the URLs and submit them all to IndexNow.
