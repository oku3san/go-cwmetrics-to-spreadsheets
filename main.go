package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
    "log"
)

func main() {
    // AWS Config の情報を取得
    ctx := context.Background()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("AWS セッションエラー: %v", err)
    }

    // CloudWatch Metrics API でメトリクスの一覧を取得します。
    cw := cloudwatch.NewFromConfig(cfg)
    input := &cloudwatch.ListMetricsInput{}
    result, err := cw.ListMetrics(ctx, input)
    if err != nil {
        log.Fatalf("CloudWatch メトリクスリストエラー: %v", err)
    }

    resultJson, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Marshal エラー: %v", err)
    }

    fmt.Println(string(resultJson))
}
