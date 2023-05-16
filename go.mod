module github.com/limero/offlinerss

go 1.18

require (
	github.com/limero/go-newsblur v0.0.0-20210107204044-9310509d25a0
	github.com/limero/go-sqldiff v0.0.0-20230514115909-1d2b5e345671
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/stretchr/testify v1.8.2
	miniflux.app v0.0.0-20220724044632-45a9fd5af60e
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/limero/go-sqldiff => ../go-sqldiff
