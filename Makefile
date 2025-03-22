install:
	go install .

testdata:
	TF_ACC=1 go test -count=1 -v

testresource:
    TF_ACC=1 go test -count=1 -run='TestAccWorkspaceResource' -v

plan: 
	terraform plan

apply: 
	terraform apply 

destroy: 
	terraform destroy 

generate:
	cd tools
	go generate ./...