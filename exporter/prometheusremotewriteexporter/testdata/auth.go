package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

// TODO: rather than defining, import constants from the Prometheus Remote Write Package
const (
	region        = "region"
	origClientStr = "origClient"
)

// SigningRoundTripper is a Custom RoundTripper that performs AWS Sig V4
type SigningRoundTripper struct {
	transport http.RoundTripper
	signer    *v4.Signer
	cfg       *aws.Config
}

// RoundTrip signs each outgoing request
func (si *SigningRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Get Credentail

	// Sign the request
	headers, err := si.signer.Sign(req, nil, "serviceStr", *si.cfg.Region, time.Now())
	if err != nil {
		// might need a response here
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	// Send the request to Cortex
	response, err := si.transport.RoundTrip(req)

	return response, err
}

// NewAuth takes a map of strings as parameters and return a http.RoundTripper
func NewAuth(params map[string]interface{}) (http.RoundTripper, error) {

	region, found := params[region]
	if !found {
		return nil, errors.New("plugin error: region not specified")
	}
	regionStr, isString := region.(string)
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
		Region: aws.String(regionStr)},
	)
	if err != nil {
		return nil, err
	}

	// Get Credentials, either from ./aws or from environmental variables
	creds := sess.Config.Credentials
	signer := v4.NewSigner(creds)

	rtp := SigningRoundTripper{
		transport: origClient.Transport,
		signer:    signer,
		cfg:       sess.Config,
	}
	// return a RoundTripper
	return &rtp, nil
}

