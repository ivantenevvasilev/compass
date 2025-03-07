/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package webhookclient_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	webhookclient "github.com/kyma-incubator/compass/components/director/pkg/webhook_client"

	accessstrategy2 "github.com/kyma-incubator/compass/components/director/pkg/accessstrategy"
	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/pkg/auth"
	"github.com/kyma-incubator/compass/components/director/pkg/correlation"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/webhook"

	"github.com/stretchr/testify/require"
)

var (
	invalidTemplate   = "invalidTemplate"
	emptyTemplate     = "{}"
	mockedError       = "mocked error"
	mockedLocationURL = "https://test-domain.com/operation"
	webhookAsyncMode  = graphql.WebhookModeAsync
)

func TestClient_Do_WhenUrlTemplateIsInvalid_ShouldReturnError(t *testing.T) {
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &invalidTemplate,
			OutputTemplate: &emptyTemplate,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to parse webhook URL")
	require.Nil(t, resp)
}

func TestClient_Do_WhenUrlTemplateIsNil_ShouldReturnError(t *testing.T) {
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    nil,
			OutputTemplate: &emptyTemplate,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing webhook url")
	require.Nil(t, resp)
}

func TestClient_Do_WhenOutputTemplateIsNil_ShouldReturnError(t *testing.T) {
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			OutputTemplate: nil,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing output template")
	require.Nil(t, resp)
}

func TestClient_Do_WhenParseInputTemplateIsInvalid_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	invalidInputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"group\": \"{{.Application.Group}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &invalidInputTemplate,
			OutputTemplate: &emptyTemplate,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to parse webhook input body")
	require.Nil(t, resp)
}

func TestClient_Do_WhenHeadersTemplateIsInvalid_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &invalidTemplate,
			OutputTemplate: &emptyTemplate,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to parse webhook headers")
	require.Nil(t, resp)
}

func TestClient_Do_WhenAuthFlowCannotBeDetermined_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &emptyTemplate,
			Auth:           &graphql.Auth{AccessStrategy: str.Ptr("wrong")},
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "could not determine auth flow")
	require.Nil(t, resp)
}

func TestClient_Do_WhenExecutingRequestFails_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &emptyTemplate,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{err: errors.New(mockedError)},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), mockedError)
	require.Nil(t, resp)
}

func TestClient_Do_WhenWebhookResponseDoesNotContainLocationURL_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusAccepted,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing location url after executing async webhook")
	require.Nil(t, resp)
}

func TestClient_Do_WhenWebhookResponseBodyContainsError_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"error\": \"%s\"}", mockedError)))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), mockedError)
	require.Contains(t, err.Error(), "received error while calling external system")
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenWebhookResponseBodyContainsErrorWithJSONObjects_ShouldParseErrorSuccessfully(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	mockedJSONObjectError := json.RawMessage(`{"code":"401","message":"Unauthorized","correlationId":"12345678-e89b-12d3-a456-556642440000"}`)

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader(json.RawMessage(fmt.Sprintf(`{"error": %s}`, mockedJSONObjectError)))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Unauthorized")
	require.Contains(t, err.Error(), "received error while calling external system")
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenWebhookResponseStatusCodeIsGoneAndGoneStatusISDefined_ShouldReturnWebhookStatusGoneError(t *testing.T) {
	goneCodeString := "404"
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := fmt.Sprintf("{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"gone_status_code\": %s,\"error\": \"{{.Body.error}}\"}", goneCodeString)
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusNotFound,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.IsType(t, webhookclient.WebhookStatusGoneErr{}, err)
	require.Contains(t, err.Error(), goneCodeString)
	require.Equal(t, http.StatusNotFound, *resp.ActualStatusCode)
}

func TestClient_Do_WhenWebhookResponseStatusCodeIsNotSuccess_ShouldReturnError(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusInternalServerError,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("response success status code was not met - expected success status code '202' or incomplete status code '204', got '%d'", http.StatusInternalServerError))
	require.Equal(t, http.StatusInternalServerError, *resp.ActualStatusCode)
}

func TestClient_Do_WhenWebhookResponseStatusCodeIsIncomplete_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusNoContent,
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, *resp.ActualStatusCode)
}

func TestClient_Do_WhenSuccessfulBasicAuthWebhook_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	username, password := "user", "pass"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
			Auth: &graphql.Auth{
				Credential: graphql.BasicCredentialData{
					Username: username,
					Password: password,
				},
			},
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
			roundTripExpectations: func(r *http.Request) {
				credentials, err := auth.LoadFromContext(r.Context())
				require.NoError(t, err)
				basicCreds, ok := credentials.(*auth.BasicCredentials)
				require.True(t, ok)
				require.Equal(t, username, basicCreds.Username)
				require.Equal(t, password, basicCreds.Password)
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenSuccessfulOAuthWebhook_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	clientID, clientSecret, tokenURL := "client-id", "client-secret", "https://test-domain.com/oauth/token"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			InputTemplate:  &inputTemplate,
			HeaderTemplate: &headersTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
			Auth: &graphql.Auth{
				Credential: graphql.OAuthCredentialData{
					ClientID:     clientID,
					ClientSecret: clientSecret,
					URL:          tokenURL,
				},
			},
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
			roundTripExpectations: func(r *http.Request) {
				credentials, err := auth.LoadFromContext(r.Context())
				require.NoError(t, err)
				oAuthCredentials, ok := credentials.(*auth.OAuthCredentials)
				require.True(t, ok)
				require.Equal(t, clientID, oAuthCredentials.ClientID)
				require.Equal(t, clientSecret, oAuthCredentials.ClientSecret)
				require.Equal(t, tokenURL, oAuthCredentials.TokenURL)
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenSuccessfulMTLSWebhook_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	mtlsCalled := false
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
			Auth: &graphql.Auth{
				AccessStrategy: str.Ptr(string(accessstrategy2.CMPmTLSAccessStrategy)),
			},
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	mtlsClient := &http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
			roundTripExpectations: func(r *http.Request) {
				mtlsCalled = true
			},
		},
	}

	client := webhookclient.NewClient(nil, mtlsClient, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.True(t, mtlsCalled)
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenSuccessfulOpenStrategyWebhook_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	openCalled := false
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			URLTemplate:    &URLTemplate,
			OutputTemplate: &outputTemplate,
			Mode:           &webhookAsyncMode,
			Auth: &graphql.Auth{
				AccessStrategy: str.Ptr(string(accessstrategy2.OpenAccessStrategy)),
			},
		},
		Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
	}

	openClient := &http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
			roundTripExpectations: func(r *http.Request) {
				openCalled = true
			},
		},
	}

	client := webhookclient.NewClient(openClient, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.True(t, openCalled)
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Do_WhenMissingCorrelationID_ShouldBeSuccessful(t *testing.T) {
	URLTemplate := "{\"method\": \"DELETE\",\"path\":\"https://test-domain.com/api/v1/applications/{{.Application.ID}}\"}"
	inputTemplate := "{\"application_id\": \"{{.Application.ID}}\",\"name\": \"{{.Application.Name}}\"}"
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	correlationIDKey := "X-Correlation-Id"
	correlationID := "abc"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.Request{
		Webhook: graphql.Webhook{
			CorrelationIDKey: &correlationIDKey,
			URLTemplate:      &URLTemplate,
			InputTemplate:    &inputTemplate,
			HeaderTemplate:   &headersTemplate,
			OutputTemplate:   &outputTemplate,
			Mode:             &webhookAsyncMode,
		},
		Object:        &webhook.ApplicationLifecycleWebhookRequestObject{Application: app, Headers: map[string]string{}},
		CorrelationID: correlationID,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusAccepted,
			},
			roundTripExpectations: func(r *http.Request) {
				headers := correlation.HeadersFromContext(r.Context())
				correlationIDAttached := false
				for headerKey, headerValue := range headers {
					if headerKey == correlationIDKey && headerValue == correlationID {
						correlationIDAttached = true
						break
					}
				}
				require.True(t, correlationIDAttached)
			},
		},
	}, nil, nil)

	resp, err := client.Do(context.Background(), webhookReq)

	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, *resp.ActualStatusCode)
}

func TestClient_Poll_WhenHeadersTemplateIsInvalid_ShouldReturnError(t *testing.T) {
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &invalidTemplate,
				StatusTemplate: &emptyTemplate,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to parse webhook headers")
}

func TestClient_Poll_WhenCreatingRequestFails_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &emptyTemplate,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)
	var ctx context.Context

	_, err := client.Poll(ctx, webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "nil Context")
}

func TestClient_Poll_WhenAuthFlowCannotBeDetermined_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{

		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &emptyTemplate,
				Auth:           &graphql.Auth{AccessStrategy: str.Ptr("wrong")},
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(http.DefaultClient, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "could not determine auth flow")
}

func TestClient_Poll_WhenExecutingRequestFails_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &emptyTemplate,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{err: errors.New(mockedError)},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), mockedError)
}

func TestClient_Poll_WhenParseStatusTemplateFails_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("{}")))},
		},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing Status Template success status code field")
}

func TestClient_Poll_WhenWebhookResponseBodyContainsError_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"error\": \"%s\"}", mockedError)))),
				StatusCode: http.StatusOK,
			},
		},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), mockedError)
	require.Contains(t, err.Error(), "received error while calling external system")
}

func TestClient_Poll_WhenWebhookResponseStatusCodeIsNotSuccess_ShouldReturnError(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusInternalServerError,
			},
		},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("response success status code was not met - expected success status code '200', got '%d'", http.StatusInternalServerError))
}

func TestClient_Poll_WhenSuccessfulBasicAuthWebhook_ShouldBeSuccessful(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	username, password := "user", "pass"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
				Auth: &graphql.Auth{
					Credential: graphql.BasicCredentialData{
						Username: username,
						Password: password,
					},
				},
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusOK,
			},
			roundTripExpectations: func(r *http.Request) {
				credentials, err := auth.LoadFromContext(r.Context())
				require.NoError(t, err)
				basicCreds, ok := credentials.(*auth.BasicCredentials)
				require.True(t, ok)
				require.Equal(t, username, basicCreds.Username)
				require.Equal(t, password, basicCreds.Password)
			},
		},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.NoError(t, err)
}

func TestClient_Poll_WhenSuccessfulOAuthWebhook_ShouldBeSuccessful(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	clientID, clientSecret, tokenURL := "client-id", "client-secret", "https://test-domain.com/oauth/token"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
				Auth: &graphql.Auth{
					Credential: graphql.OAuthCredentialData{
						ClientID:     "client-id",
						ClientSecret: "client-secret",
						URL:          "https://test-domain.com/oauth/token",
					},
				},
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusOK,
			},
			roundTripExpectations: func(r *http.Request) {
				credentials, err := auth.LoadFromContext(r.Context())
				require.NoError(t, err)
				oAuthCredentials, ok := credentials.(*auth.OAuthCredentials)
				require.True(t, ok)
				require.Equal(t, clientID, oAuthCredentials.ClientID)
				require.Equal(t, clientSecret, oAuthCredentials.ClientSecret)
				require.Equal(t, tokenURL, oAuthCredentials.TokenURL)
			},
		},
	}, nil, nil)
	_, err := client.Poll(context.Background(), webhookReq)

	require.NoError(t, err)
}

func TestClient_Poll_WhenSuccessfulMTLSWebhook_ShouldBeSuccessful(t *testing.T) {
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	outputTemplate := "{\"location\":\"{{.Headers.Location}}\",\"success_status_code\": 202,\"incomplete_status_code\": 204,\"error\": \"{{.Body.error}}\"}"
	mtlsCalled := false

	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	pollRequest := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				OutputTemplate: &outputTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
				Auth: &graphql.Auth{
					AccessStrategy: str.Ptr(string(accessstrategy2.CMPmTLSAccessStrategy)),
				},
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: "https://test-domain.com/poll/",
	}

	mtlsClient := &http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				Header:     http.Header{"Location": []string{mockedLocationURL}},
				StatusCode: http.StatusOK,
			},
			roundTripExpectations: func(r *http.Request) {
				mtlsCalled = true
			},
		},
	}

	client := webhookclient.NewClient(nil, mtlsClient, nil)

	_, err := client.Poll(context.Background(), pollRequest)

	require.NoError(t, err)
	require.True(t, mtlsCalled)
}

func TestClient_Poll_WhenSuccessfulWebhookPollResponseContainsNullErrorField_ShouldBeSuccessful(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{\"error\":null}"))),
				StatusCode: http.StatusOK,
			},
		},
	}, nil, nil)
	_, err := client.Poll(context.Background(), webhookReq)

	require.NoError(t, err)
}

func TestClient_Poll_WhenSuccessfulWebhookPollResponseContainsEmptyErrorField_ShouldBeSuccessful(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				HeaderTemplate: &headersTemplate,
				StatusTemplate: &statusTemplate,
				Mode:           &webhookAsyncMode,
			},
			Object: &webhook.ApplicationLifecycleWebhookRequestObject{Application: app},
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{\"error\":\"\"}"))),
				StatusCode: http.StatusOK,
			},
		},
	}, nil, nil)
	_, err := client.Poll(context.Background(), webhookReq)

	require.NoError(t, err)
}

func TestClient_Poll_WhenMissingCorrelationID_ShouldBeSuccessful(t *testing.T) {
	headersTemplate := "{\"user-identity\":[\"{{.Headers.Client_user}}\"]}"
	statusTemplate := "{\"status\":\"{{.Body.status}}\",\"success_status_code\": 200,\"success_status_identifier\":\"SUCCEEDED\",\"in_progress_status_identifier\":\"IN_PROGRESS\",\"failed_status_identifier\":\"FAILED\",\"error\": \"{{.Body.error}}\"}"
	correlationIDKey := "X-Correlation-Id"
	correlationID := "abc"
	app := &graphql.Application{BaseEntity: &graphql.BaseEntity{ID: "appID"}}
	webhookReq := &webhookclient.PollRequest{
		Request: &webhookclient.Request{
			Webhook: graphql.Webhook{
				CorrelationIDKey: &correlationIDKey,
				HeaderTemplate:   &headersTemplate,
				StatusTemplate:   &statusTemplate,
				Mode:             &webhookAsyncMode,
			},
			Object:        &webhook.ApplicationLifecycleWebhookRequestObject{Application: app, Headers: map[string]string{}},
			CorrelationID: correlationID,
		},
		PollURL: mockedLocationURL,
	}

	client := webhookclient.NewClient(&http.Client{
		Transport: mockedTransport{
			resp: &http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusOK,
			},
			roundTripExpectations: func(r *http.Request) {
				headers := correlation.HeadersFromContext(r.Context())
				correlationIDAttached := false
				for headerKey, headerValue := range headers {
					if headerKey == correlationIDKey && headerValue == correlationID {
						correlationIDAttached = true
						break
					}
				}
				require.True(t, correlationIDAttached)
			},
		},
	}, nil, nil)

	_, err := client.Poll(context.Background(), webhookReq)

	require.NoError(t, err)
}

type mockedTransport struct {
	resp                  *http.Response
	err                   error
	roundTripExpectations func(r *http.Request)
}

func (m mockedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if m.roundTripExpectations != nil {
		m.roundTripExpectations(r)
	}
	return m.resp, m.err
}
