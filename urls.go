package dgo2poc

func EndpointUser(uid string) string { return "/users/" + uid }

const EndpointGateway = "/gateway"

const EndpointGatewayBot = "/gateway/bot"
