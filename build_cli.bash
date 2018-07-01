pushd ./cmd/attache/internal/cmd_gen >> /dev/null
go generate
popd >> /dev/null
pushd ./cmd/attache/internal/cmd_new >> /dev/null
go generate
popd >> /dev/null

go install ./cmd/attache
