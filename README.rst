Nuts consent store
==================

Consent storage for querying and checking active consent. This is the service space copy of the distributed consent. 
Using the distributed records from the Corda vault is not possible due to the encrypted nature of those records.

Binary format migrations
------------------------

go get -u github.com/go-bindata/go-bindata/... (outside module)
cd migrations && go-bindata -pkg migrations .
