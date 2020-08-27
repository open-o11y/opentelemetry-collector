// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusremotewriteexporter

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr       = "prometheusremotewrite"
	regionStr     = "region"
	serviceStr    = "service"
	origClientStr = "origClient"
)

// SigningRoundTripper is a Custom RoundTripper that performs AWS Sig V4
type SigningRoundTripper struct {
	transport http.RoundTripper
	signer    *v4.Signer
	service   string
	cfg       *aws.Config
}

// RoundTrip signs each outgoing request
func (si *SigningRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBody, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	// Get the body
	content, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(content)

	// Sign the request
	headers, err := si.signer.Sign(req, body, si.service, *si.cfg.Region, time.Now())
	if err != nil {
		// might need a response here
		return nil, err
	}
	for k, vs := range headers {
		req.Header.Del(k)
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	log.Println(req)

	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	//	log.Println(err)
	// }
	// f, err := os.Create("./dat")
	// defer f.Close()
	// f.Write(requestDump)
	// f.Sync()

	// Send the request to Cortex
	response, err := si.transport.RoundTrip(req)
	// log.Println(response)
	 bodyBytes, err := ioutil.ReadAll(response.Body)
	 if err != nil {
		log.Fatal(err)
	 }
	 bodyString := string(bodyBytes)
	 log.Println("response: ", bodyString)
	return response, err
}

// NewAuth takes a map of strings as parameters and return a http.RoundTripper
func NewAuth(params map[string]interface{}) (http.RoundTripper, error) {

	reg, found := params[regionStr]
	if !found {
		return nil, errors.New("plugin error: region not specified")
	}
	region, isString := reg.(string)
	if !isString {
		return nil, errors.New("plugin error: region is not string")
	}
	serv, found := params[serviceStr]
	if !found {
		return nil, errors.New("plugin error: service not specified")
	}

	service, isString := serv.(string)
	if !isString {
		return nil, errors.New("plugin error: region is not string")
	}

	client, found := params[origClientStr]
	if !found {
		return nil, errors.New("plugin error: default client not specified")
	}
	origClient, isClient := client.(*http.Client)
	if !isClient {
		return nil, errors.New("plugin error: default client not specified")
	}

	// Initialize session with default credential chain
	// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
		aws.NewConfig().WithLogLevel(aws.LogDebugWithSigning),
	)
	if err != nil {
		log.Println("AWS session initialization failed")
	}

	if _, err = sess.Config.Credentials.Get(); err != nil {
		log.Println("AWS session initialized, but credentials are not loaded correctly")
	}

	// Get Credentials, either from ./aws or from environmental variables
	creds := sess.Config.Credentials
	signer := v4.NewSigner(creds)

	signer.Debug = aws.LogDebugWithSigning
	signer.Logger = aws.NewDefaultLogger()
	rtp := SigningRoundTripper{
		transport: origClient.Transport,
		signer:    signer,
		cfg:       sess.Config,
		service:   service,
	}
	// return a RoundTripper
	return &rtp, nil
}

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithMetrics(createMetricsExporter))
}

func createMetricsExporter(_ context.Context, _ component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.MetricsExporter, error) {

	prwCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid configuration")
	}
	client, cerr := prwCfg.HTTPClientSettings.ToClient()
	if cerr != nil {
		return nil, cerr
	}
	if prwCfg.AuthCfg != nil {
		authConfig := make(map[string]interface{})
		authConfig[serviceStr] = prwCfg.AuthCfg[serviceStr]
		authConfig[regionStr] = prwCfg.AuthCfg[regionStr]
		authConfig[origClientStr] = client

		roundTripper, err := NewAuth(authConfig)
		if err != nil {
			return nil, err
		}

		client.Transport = roundTripper
	}

	prwe, err := newPrwExporter(prwCfg.Namespace, prwCfg.HTTPClientSettings.Endpoint, client)
	if err != nil {
		return nil, err
	}

	prwexp, err := exporterhelper.NewMetricsExporter(
		cfg,
		prwe.pushMetrics,
		exporterhelper.WithTimeout(prwCfg.TimeoutSettings),
		exporterhelper.WithQueue(prwCfg.QueueSettings),
		exporterhelper.WithRetry(prwCfg.RetrySettings),
		exporterhelper.WithShutdown(prwe.shutdown),
	)

	return prwexp, err
}

func createDefaultConfig() configmodels.Exporter {
	qs := exporterhelper.CreateDefaultQueueSettings()
	qs.Enabled = false

	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		Namespace: "",

		TimeoutSettings: exporterhelper.CreateDefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.CreateDefaultRetrySettings(),
		QueueSettings:   qs,
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: "http://some.url:9411/api/prom/push",
			// We almost read 0 bytes, so no need to tune ReadBufferSize.
			ReadBufferSize:  0,
			WriteBufferSize: 512 * 1024,
			Timeout:         exporterhelper.CreateDefaultTimeoutSettings().Timeout,
			Headers:         map[string]string{},
		},
	}
}
