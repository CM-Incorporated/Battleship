module cmincorporated.com/server

go 1.26.5

replace (
	cmincorporated.com/gamecode => ../gamecode
	cmincorporated.com/protocol => ../protocol
)

require (
	cmincorporated.com/gamecode v0.0.0-00010101000000-000000000000
	cmincorporated.com/protocol v0.0.0
	github.com/coder/websocket v1.8.15
)
