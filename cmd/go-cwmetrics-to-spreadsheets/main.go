package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
    "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"
    "log"
    "os"
)

func main() {
    // AWS Config の情報を取得
    ctx := context.Background()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("AWS セッションエラー: %v", err)
    }

    // CloudWatch Metrics API でメトリクスの一覧を取得
    cw := cloudwatch.NewFromConfig(cfg)
    input := &cloudwatch.ListMetricsInput{}
    result, err := cw.ListMetrics(ctx, input)
    if err != nil {
        log.Fatalf("CloudWatch メトリクスリストエラー: %v", err)
    }

    // Google Sheets API を使用して、取得したメトリクスを Google スプレッドシートに書き込み
    credential := option.WithCredentialsFile("./../../credentials/gcp-credential.json")
    sheetsService, err := sheets.NewService(ctx, credential)
    if err != nil {
        log.Fatalf("Google Sheets API エラー: %v", err)
    }

    // スプレッドシートIDとシート名を指定
    spreadsheetId := os.Args[1]
    sheetName := os.Args[2]

    // スプレッドシートに書き込むデータを格納するためのリストを作成
    var values [][]interface{}
    for _, metric := range result.Metrics {
        values = append(values, []interface{}{metric.Namespace, metric.MetricName})
    }

    // スプレッドシートにデータを書き込み
    writeRange := fmt.Sprintf("%s!A2:B%d", sheetName, len(result.Metrics)+1)
    valueRange := &sheets.ValueRange{Values: values}
    _, err = sheetsService.Spreadsheets.Values.Update(spreadsheetId, writeRange, valueRange).ValueInputOption("USER_ENTERED").Do()
    if err != nil {
        log.Fatalf("Google Sheets API データ書き込みエラー: %v", err)
    }
}
