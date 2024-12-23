package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	distributedcachestu "distributed-cache-stu"
	"distributed-cache-stu/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	grpcServer := grpc.NewServer()

	addr, ok := os.LookupEnv("PORT")
	apiPort, _ := os.LookupEnv("API_PORT")

	if !ok || addr == "" {
		addr = "50051"
	}

	group := distributedcachestu.NewGroup("scores", 2<<10, distributedcachestu.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[Slow DB] search key", key)

		if v, ok := db[key]; ok {
			return []byte(v), nil
		}

		return nil, fmt.Errorf("%s not exist", key)
	}))

	peers := distributedcachestu.NewGrpcPool(addr)
	peers.Set("50051", "50052", "50053")
	group.RegisterPeers(peers)

	cacheService := &distributedcachestu.Service{}

	proto.RegisterCacheServiceServer(grpcServer, cacheService)

	log.Printf("grpc server starting on address: %s", addr)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", ":"+addr)

		if err != nil {
			log.Fatal(err)
			return err
		}

		return grpcServer.Serve(lis)
	})

	if apiPort != "" {
		mu := http.NewServeMux()
		mu.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write(view.ByteSlice())
		})

		srv := &http.Server{Addr: ":" + apiPort, Handler: mu}

		g.Go(func() error {
			return srv.ListenAndServe()
		})

		g.Go(func() error {
			<-ctx.Done()
			return srv.Shutdown(context.Background())
		})
	}

	g.Go(func() error {
		<-ctx.Done()
		grpcServer.GracefulStop()
		return nil
	})

	if err := g.Wait(); err != nil {
		os.Exit(1)
	}
}
