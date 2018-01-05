package dgo2poc

func EndpointUser(uid string) string { return "/users/" + uid }

func EndpointChannelMessages(cid string) string { return "/channels/" + cid + "/messages" }

const EndpointGateway = "/gateway"

const EndpointGatewayBot = "/gateway/bot"
