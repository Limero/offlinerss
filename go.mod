module github.com/limero/offlinerss

go 1.18

require (
	github.com/limero/go-newsblur v0.0.0-20230708133720-098cbaea1cca
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/stretchr/testify v1.8.2
	miniflux.app v0.0.0-20220724044632-45a9fd5af60e
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/limero/go-newsblur => ../go-newsblur
