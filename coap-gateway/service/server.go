package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pbAS "github.com/go-ocf/cloud/authorization/pb"
	"github.com/go-ocf/cloud/coap-gateway/uri"
	notificationRA "github.com/go-ocf/cloud/resource-aggregate/cqrs/notification"
	projectionRA "github.com/go-ocf/cloud/resource-aggregate/cqrs/projection"
	pbRA "github.com/go-ocf/cloud/resource-aggregate/pb"
	pbRD "github.com/go-ocf/cloud/resource-directory/pb/resource-directory"
	pbRS "github.com/go-ocf/cloud/resource-directory/pb/resource-shadow"
	"github.com/go-ocf/cqrs/eventbus"
	"github.com/go-ocf/cqrs/eventstore"
	"github.com/go-ocf/go-coap/v2/blockwise"
	"github.com/go-ocf/go-coap/v2/keepalive"
	"github.com/go-ocf/go-coap/v2/message"
	coapCodes "github.com/go-ocf/go-coap/v2/message/codes"
	"github.com/go-ocf/go-coap/v2/mux"
	"github.com/go-ocf/go-coap/v2/net"
	"github.com/go-ocf/go-coap/v2/tcp"
	kitNetCoap "github.com/go-ocf/kit/net/coap"

	gocoap "github.com/go-ocf/go-coap"
	"github.com/go-ocf/kit/log"
	cache "github.com/patrickmn/go-cache"
)

//Server a configuration of coapgateway
type Server struct {
	FQDN                            string // fully qualified domain name of GW
	ExternalPort                    uint16 // used to construct oic/res response
	Addr                            string // Address to listen on, ":COAP" if empty.
	Net                             string // if "tcp" or "tcp-tls" (COAP over TLS) it will invoke a TCP listener, otherwise an UDP one
	Keepalive                       keepalive.Config
	DisableTCPSignalMessageCSM      bool
	DisablePeerTCPSignalMessageCSMs bool
	SendErrorTextInResponse         bool
	RequestTimeout                  time.Duration
	ConnectionsHeartBeat            time.Duration
	BlockWiseTransfer               bool
	BlockWiseTransferSZX            blockwise.SZX

	raClient pbRA.ResourceAggregateClient
	asClient pbAS.AuthorizationServiceClient
	rsClient pbRS.ResourceShadowClient
	rdClient pbRD.ResourceDirectoryClient

	clientContainer               *ClientContainer
	clientContainerByDeviceID     *clientContainerByDeviceID
	updateNotificationContainer   *notificationRA.UpdateNotificationContainer
	retrieveNotificationContainer *notificationRA.RetrieveNotificationContainer
	observeResourceContainer      *observeResourceContainer
	goroutinesPool                *ants.Pool
	oicPingCache                  *cache.Cache

	projection      *projectionRA.Projection
	coapServer      *tcp.Server
	listener        tcp.Listener
	authInterceptor kitNetCoap.Interceptor
}

type DialCertManager = interface {
	GetClientTLSConfig() *tls.Config
}

type ListenCertManager = interface {
	GetServerTLSConfig() *tls.Config
}

//NewServer setup coap gateway
func New(config Config, dialCertManager DialCertManager, listenCertManager ListenCertManager, authInterceptor kitNetCoap.Interceptor, store eventstore.EventStore, subscriber eventbus.Subscriber, pool *ants.Pool) *Server {
	oicPingCache := cache.New(cache.NoExpiration, time.Minute)
	oicPingCache.OnEvicted(pingOnEvicted)

	dialTLSConfig := dialCertManager.GetClientTLSConfig()

	raConn, err := grpc.Dial(config.ResourceAggregateAddr, grpc.WithTransportCredentials(credentials.NewTLS(dialTLSConfig)))
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	raClient := pbRA.NewResourceAggregateClient(raConn)

	asConn, err := grpc.Dial(config.AuthServerAddr, grpc.WithTransportCredentials(credentials.NewTLS(dialTLSConfig)))
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	asClient := pbAS.NewAuthorizationServiceClient(asConn)

	rdConn, err := grpc.Dial(config.ResourceDirectoryAddr, grpc.WithTransportCredentials(credentials.NewTLS(dialTLSConfig)))
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	var listener tcp.Listener

	if listenCertManager == nil || reflect.ValueOf(listenCertManager).IsNil() {
		l, err := net.NewTCPListener("tcp", config.Addr)
		if err != nil {
			log.Fatalf("cannot setup tcp for server: %v", err)
		}
		listener = l
	} else {
		tlsConfig := listenCertManager.GetServerTLSConfig()
		l, err := net.NewTLSListener("tcp", config.Addr, tlsConfig)
		if err != nil {
			log.Fatalf("cannot setup tcp-tls for server: %v", err)
		}
		listener = l
	}
	rdClient := pbRD.NewResourceDirectoryClient(rdConn)
	rsClient := pbRS.NewResourceShadowClient(rdConn)

	var keepalive *keepalive.KeepAlive
	if config.KeepaliveEnable {
		keepalive, err = keepalive.MakeKeepAlive(config.KeepaliveTimeoutConnection)
		if err != nil {
			log.Fatalf("cannot setup keepalive for server: %v", err)
		}
	}

	var blockWiseTransferSZX gocoap.BlockWiseSzx
	switch strings.ToLower(config.BlockWiseTransferSZX) {
	case "16":
		blockWiseTransferSZX = gocoap.BlockWiseSzx16
	case "32":
		blockWiseTransferSZX = gocoap.BlockWiseSzx32
	case "64":
		blockWiseTransferSZX = gocoap.BlockWiseSzx64
	case "128":
		blockWiseTransferSZX = gocoap.BlockWiseSzx128
	case "256":
		blockWiseTransferSZX = gocoap.BlockWiseSzx256
	case "512":
		blockWiseTransferSZX = gocoap.BlockWiseSzx512
	case "1024":
		blockWiseTransferSZX = gocoap.BlockWiseSzx1024
	case "bert":
		blockWiseTransferSZX = gocoap.BlockWiseSzxBERT
	default:
		log.Fatalf("invalid value BlockWiseTransferSZX %v", config.BlockWiseTransferSZX)
	}

	s := Server{
		Keepalive:                       keepalive,
		Net:                             config.Net,
		FQDN:                            config.FQDN,
		ExternalPort:                    config.ExternalPort,
		Addr:                            config.Addr,
		RequestTimeout:                  config.RequestTimeout,
		DisableTCPSignalMessageCSM:      config.DisableTCPSignalMessageCSM,
		DisablePeerTCPSignalMessageCSMs: config.DisablePeerTCPSignalMessageCSMs,
		SendErrorTextInResponse:         config.SendErrorTextInResponse,
		ConnectionsHeartBeat:            config.ConnectionsHeartBeat,
		BlockWiseTransfer:               !config.DisableBlockWiseTransfer,
		BlockWiseTransferSZX:            blockWiseTransferSZX,

		raClient: raClient,
		asClient: asClient,
		rsClient: rsClient,
		rdClient: rdClient,

		clientContainer:               &ClientContainer{sessions: make(map[string]*Client)},
		clientContainerByDeviceID:     NewClientContainerByDeviceId(),
		updateNotificationContainer:   notificationRA.NewUpdateNotificationContainer(),
		retrieveNotificationContainer: notificationRA.NewRetrieveNotificationContainer(),
		observeResourceContainer:      NewObserveResourceContainer(),
		goroutinesPool:                pool,
		oicPingCache:                  oicPingCache,
		listener:                      listener,
		authInterceptor:               authInterceptor,
	}

	projection, err := projectionRA.NewProjection(context.Background(), fmt.Sprintf("%v:%v", config.FQDN, config.ExternalPort), store, subscriber, newResourceCtx(&s))
	if err != nil {
		log.Fatalf("cannot create projection for server: %v", err)
	}
	s.projection = projection
	return &s
}

func getDeviceID(client *Client) string {
	deviceID := "unknown"
	if client != nil {
		deviceID = client.loadAuthorizationContext().DeviceId
		if deviceID == "" {
			deviceID = fmt.Sprintf("unknown(%v)", client.remoteAddrString())
		}
	}
	return deviceID
}

func validateCommand(s mux.ResponseWriter, req *message.Message, server *Server, fnc func(s mux.ResponseWriter, req *message.Message, client *Client)) {
	client := server.clientContainer.Find(req.Client.RemoteAddr().String())

	switch req.Msg.Code() {
	case coapCodes.POST, coapCodes.DELETE, coapCodes.PUT, coapCodes.GET:
		if client == nil {
			logAndWriteErrorResponse(fmt.Errorf("cannot handle command: client not found"), s, client, coapCodes.InternalServerError)
			return
		}
		fnc(s, req, client)
	case coapCodes.Empty:
		if client == nil {
			logAndWriteErrorResponse(fmt.Errorf("cannot handle command: client not found"), s, client, coapCodes.InternalServerError)
			return
		}
		clientResetHandler(s, req, client)
	case coapCodes.Content:
		// Unregistered observer at a peer send us a notification - inform the peer to remove it
		sendResponse(s, client, coapCodes.Empty, gocoap.TextPlain, nil)
	default:
		deviceID := getDeviceID(client)
		log.Errorf("DeviceId: %v: received invalid code: CoapCode(%v)", deviceID, req.Msg.Code())
	}
}

func defaultHandler(s mux.ResponseWriter, req *message.Message, client *Client) {
	path := req.Msg.PathString()

	switch {
	case strings.HasPrefix(path, resourceRoute):
		resourceRouteHandler(s, req, client)
	default:
		deviceID := getDeviceID(client)
		logAndWriteErrorResponse(fmt.Errorf("DeviceId: %v: unknown path %v", deviceID, path), s, client, coapCodes.NotFound)
	}
}

func (server *Server) coapConnOnNew(coapConn *gocoap.ClientConn) {
	remoteAddr := coapConn.RemoteAddr().String()
	server.clientContainer.Add(remoteAddr, newClient(server, coapConn))
}

func (server *Server) coapConnOnClose(coapConn *gocoap.ClientConn, err error) {
	if err != nil {
		log.Errorf("coap connection closed with error: %v", err)
	}
	if client, ok := server.clientContainer.Pop(coapConn.RemoteAddr().String()); ok {
		client.OnClose()
	}

}

func (server *Server) logginMiddleware(next func(gocoap.ResponseWriter, *gocoap.Request)) func(gocoap.ResponseWriter, *gocoap.Request) {
	return func(w mux.ResponseWriter, r *message.Message) {
		client := server.clientContainer.Find(req.Client.RemoteAddr().String())
		decodeMsgToDebug(client, req.Msg, "RECEIVED-COMMAND")
		next(w, req)
	}
}

func (server *Server) authMiddleware(next func(gocoap.ResponseWriter, *gocoap.Request)) func(gocoap.ResponseWriter, *gocoap.Request) {
	return func(w mux.ResponseWriter, r *message.Message) {
		client := server.clientContainer.Find(req.Client.RemoteAddr().String())
		if client == nil {
			logAndWriteErrorResponse(fmt.Errorf("cannot handle request: client not found"), w, client, coapCodes.InternalServerError)
			return
		}

		ctx := kitNetCoap.CtxWithToken(req.Ctx, client.loadAuthorizationContext().AccessToken)
		_, err := server.authInterceptor(ctx, req.Msg.Code(), "/"+req.Msg.PathString())
		if err != nil {
			logAndWriteErrorResponse(fmt.Errorf("cannot handle request to path '%v': %v", req.Msg.PathString(), err), w, client, coapCodes.Unauthorized)
			client.Close()
			return
		}
		next(w, req)
	}
}

//setupCoapServer setup coap server
func (server *Server) setupCoapServer() {
	m := mux.NewServeMux()
	m.DefaultHandle(mux.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, defaultHandler)
	}))
	m.Handle(uri.ResourceDirectory, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, resourceDirectoryHandler)
	}))
	m.Handle(uri.SignUp, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, signUpHandler)
	}))
	m.Handle(uri.SecureSignUp, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, signUpHandler)
	}))
	m.Handle(uri.SignIn, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, signInHandler)
	}))
	m.Handle(uri.SecureSignIn, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, signInHandler)
	}))
	m.Handle(uri.ResourceDiscovery, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, resourceDiscoveryHandler)
	}))
	m.Handle(uri.ResourcePing, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, resourcePingHandler)
	}))
	m.Handle(uri.RefreshToken, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, refreshTokenHandler)
	}))
	m.Handle(uri.SecureRefreshToken, gocoap.HandlerFunc(func(w mux.ResponseWriter, r *message.Message) {
		validateCommand(w, r, server, refreshTokenHandler)
	}))

	server.coapServer = &gocoap.Server{
		Net:                             server.Net,
		Addr:                            server.Addr,
		DisableTCPSignalMessageCSM:      server.DisableTCPSignalMessageCSM,
		DisablePeerTCPSignalMessageCSMs: server.DisablePeerTCPSignalMessageCSMs,
		KeepAlive:                       server.Keepalive,
		Handler:                         gocoap.HandlerFunc(server.logginMiddleware(server.authMiddleware(m.ServeCOAP))),
		NotifySessionNewFunc:            server.coapConnOnNew,
		NotifySessionEndFunc:            server.coapConnOnClose,
		HeartBeat:                       server.ConnectionsHeartBeat,
		BlockWiseTransfer:               &server.BlockWiseTransfer,
		BlockWiseTransferSzx:            &server.BlockWiseTransferSZX,
	}
}

func (server *Server) tlsEnabled() bool {
	return strings.HasSuffix(server.Net, "-tls")
}

// Serve starts a coapgateway on the configured address in *Server.
func (server *Server) Serve() error {
	server.setupCoapServer()
	server.coapServer.Listener = server.listener
	return server.coapServer.ActivateAndServe()
}

// Shutdown turn off server.
func (server *Server) Shutdown() error {
	err := server.coapServer.Shutdown()
	server.listener.Close()
	return err
}
