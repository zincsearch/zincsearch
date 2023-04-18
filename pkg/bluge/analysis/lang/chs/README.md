# zinc-analysis-gse

it's a plugin of zinc to support Chinese analyzer.

Analyzer: `gse_standard` , `gse_search`

Tokenizer: `gse_standard` , `gse_search`

TokenFilter: `gse_stop`

> build has embed dictionary of `zh/s_1.txt`, `zh/stop_tokens.txt`.

you can find it: https://github.com/go-ego/gse/tree/master/data/dict

> also you can custom dictionary follow [custom user dictionary](#custom-user-dictionary)

after custom, you need restart zinc.

## gse

https://github.com/go-ego/gse

Go efficient multilingual NLP and text segmentation; support english, chinese, japanese and other.

## Environment

you need pass environment to enable gse support:

`ZINC_PLUGIN_GSE_ENABLE` true of false, default is `false`

`ZINC_PLUGIN_GSE_DICT_EMBED` small or big, default is `small`, which size dictionary will load when `gse` enabled.

`ZINC_PLUGIN_GSE_DICT_PATH` custom dictionary path, default is `./plugins/gse/dict`


## API example

POST http://localhost:4080/es/_analyze

```
{
  "analyzer": "gse_standard",
  "text": "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片."
}
```

POST http://localhost:4080/es/_analyze

```
{
  "analyzer": "gse_search",
  "text": "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片."
}
```

PUT http://localhost:4080/api/index

```
{
	"name": "my-index-chs",
		"mappings": {
			"properties": {
				"title": {
					"type": "text",
					"index": true,
					"highlightable": true,
					"analyzer": "gse_search",
					"search_analyzer": "gse_standard"
				},
				"author": {
					"type": "keyword",
					"index": true,
					"store": false
				},
				"create_time": {
					"type":"date"
				}
			}
		}
}
```

POST http://localhost:4080/api/my-index-chs/document

```
{
	"title": "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片",
	"author": "灭霸",
	"create_time": "2022-03-05T18:18:18+08:00"
}
```

POST http://localhost:4080/es/my-index-chs/_search

```
{
	"query": {
		"match": {
			"title": "复仇者联盟"
		}
	}
}
```

## custom user dictionary

add your words append to the file `${ZINC_PLUGIN_GSE_DICT_PATH}/user.txt`

format:

```
分词文本  频率        词性
word    frequency   property
```

like:

```
复仇者联盟 100 n
```

## custom stop tokens

add your words append to the file `${ZINC_PLUGIN_GSE_DICT_PATH}/stop.txt`

format:

```
停止词
word
```

like:

```
哈哈
```

## Credit

* https://github.com/zincsearch/zincsearch
* https://github.com/blugelabs/bluge
* https://github.com/go-ego/gse
