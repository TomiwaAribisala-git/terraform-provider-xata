# Unofficial Terraform Provider for Xata
- Website: https://xata.io/
- Documentation:

## Requirements
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.23

## Contributing
- Fork and clone this repo
- Run `make install`
- Make and test your changes `cd examples` `&&` `make plan` `&&` `make apply` `&&` `make destroy`
- Generate Documentation `make generate`
- Submit a PR!

## Todo/Extras
- Database: List/Create/Read/Update/Delete
- Branch: List/Create/Read/Update/Delete
- Table: Create/Delete/Update
- [Terraform Provider Tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider)
- [Xata API](https://xata.io/docs/api-reference)
- [Xata Go Client](https://github.com/xataio/xata-go/blob/main/xata)