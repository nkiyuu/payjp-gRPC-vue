package main

import (
	"context"
	"log"
	"net"
	"os"
	gpay "vue-golang-payment-app/payment-service/proto"

	payjp "github.com/payjp/payjp-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Charge(ctx context.Context, req *gpay.PayRequest) (*gpay.PayResponse, error) {
	// pay の初期化
	pay := payjp.New(os.Getenv("PAYJP_TEST_SECRET_KEY"), nil)

	// 支払いを行う。第一引数に支払い金額、第二引数に支払い方法や設定を入れる。
	charge, err := pay.Charge.Create(int(req.Amount), payjp.Charge{
		// 現在はjpyのみサポート
		Currency: "jpy",
		// カード情報、顧客ID、カードトークンのいずれかを指定。今回はToken。
		CardToken: req.Token,
		Capture:   true,
		// 概要のテキストを設定できる
		Description: req.Name + ":" + req.Discription,
	})
	if err != nil {
		return nil, err
	}

	// 支払った結果から、Response生成
	res := &gpay.PayResponse{
		Paid:     charge.Paid,
		Captured: charge.Captured,
		Amount:   int64(charge.Amount),
	}

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen* %s", err)
	}

	s := grpc.NewServer()
	gpay.RegisterPayManagerServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("gRPC Server started: localhost:%s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failde to serve: %s", err)
	}
}
