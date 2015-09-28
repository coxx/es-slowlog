## ElasticSearch slowlog parser

This tool allows to read [ElasticSearch slowlog](https://www.elastic.co/guide/en/elasticsearch/reference/2.0/index-modules-slowlog.html) records and produce output using  [text/template](https://golang.org/pkg/text/template/) format.

It is mainly intended as helper tool for slowlog analysis and to produce input data for load testing tools. There is predefined output formats for:
 - [Vegeta](https://github.com/coxx/vegeta) (use patched version until https://github.com/tsenart/vegeta/pull/148 is closed)
 - [Yandex Tank](https://github.com/yandex/yandex-tank)


### Install

```
go get github.com/coxx/es-slowlog/cmd/es-slowlog-parse
```


### Examples

Write all queries to standard output
```
cat *_search_slowlog.log | es-slowlog-parse
```

Convert slowlog to tab-separated file
```
cat *_search_slowlog.log | es-slowlog-parse -f '{{.Index}}\t{{.Types}}\t{{.SearchType}}\t{{.Source}}'
```

Produce [vegeta](https://github.com/coxx/vegeta) target file 
```
cat *_search_slowlog.log | es-slowlog-parse -f vegeta -a localhost:9200
```


### Notes

* Parser and scanner inspired by Rob Pike's talk [Lexical Scanning in Go](http://www.youtube.com/watch?v=HxaD_trXwRE). I'm pretty sure all this code can be replaced by couple of regular expressions but then it was not so cool.
* Directory structure inspired by Mat Ryer's post [5 simple tips and tricks for writing unit tests in #golang](https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742).
