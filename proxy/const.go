package proxy

const cpListen = "127.0.0.1:9999"
const socketAddress = "/tmp/sidecar.sock"

const grpcAuthorizationKey = "grpc-authorization"
const grpcAuthorizationMethodKey = "grpc-signing-method"

const signingMethodECDSA = "SigningMethodECDSA"
const signingMethodEd25519 = "SigningMethodEd25519"
const signingMethodHMAC = "SigningMethodHMAC"
const signingMethodRSA = "SigningMethodRSA"
const signingMethodRSAPSS = "SigningMethodRSAPSS"
