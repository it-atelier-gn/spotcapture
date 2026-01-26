package upload

import (
	"bytes"
	"context"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func Upload(ctx context.Context, data []byte, key string) error {
	proxyURL, _ := url.Parse(viper.GetString("proxy"))

	customHTTPClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(viper.GetString("region")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(viper.GetString("key"), viper.GetString("secret"), ""),
		),
		config.WithHTTPClient(customHTTPClient),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	var bucket = viper.GetString("bucket")

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}
