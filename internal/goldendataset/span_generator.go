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

package goldendataset

import (
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/translator/conventions"
)

var statusCodeMap = map[PICTInputStatus]pdata.StatusCode{
	SpanStatusUnset: pdata.StatusCodeUnset,
	SpanStatusOk:    pdata.StatusCodeOk,
	SpanStatusError: pdata.StatusCodeError,
}

var statusMsgMap = map[PICTInputStatus]string{
	SpanStatusUnset: "Unset",
	SpanStatusOk:    "Ok",
	SpanStatusError: "Error",
}

// GenerateSpans generates a slice of pdata.Span objects with the number of spans specified by the count input
// parameter. The startPos parameter specifies the line in the PICT tool-generated, test parameter
// combination records file specified by the pictFile parameter to start reading from. When the end record
// is reached it loops back to the first record. The random parameter injects the random number generator
// to use in generating IDs and other random values. Using a random number generator with the same seed value
// enables reproducible tests.
//
// The return values are the slice with the generated spans, the starting position for the next generation
// run and the error which caused the spans generation to fail. If err is not nil, the spans slice will
// have nil values.
func GenerateSpans(count int, startPos int, pictFile string, random io.Reader) ([]pdata.Span, int, error) {
	pairsData, err := loadPictOutputFile(pictFile)
	if err != nil {
		return nil, 0, err
	}
	pairsTotal := len(pairsData)
	spanList := make([]pdata.Span, count)
	index := startPos + 1
	var inputs []string
	var spanInputs *PICTSpanInputs
	var traceID pdata.TraceID
	var parentID pdata.SpanID
	for i := 0; i < count; i++ {
		if index >= pairsTotal {
			index = 1
		}
		inputs = pairsData[index]
		spanInputs = &PICTSpanInputs{
			Parent:     PICTInputParent(inputs[SpansColumnParent]),
			Tracestate: PICTInputTracestate(inputs[SpansColumnTracestate]),
			Kind:       PICTInputKind(inputs[SpansColumnKind]),
			Attributes: PICTInputAttributes(inputs[SpansColumnAttributes]),
			Events:     PICTInputSpanChild(inputs[SpansColumnEvents]),
			Links:      PICTInputSpanChild(inputs[SpansColumnLinks]),
			Status:     PICTInputStatus(inputs[SpansColumnStatus]),
		}
		switch spanInputs.Parent {
		case SpanParentRoot:
			traceID = generatePDataTraceID(random)
			parentID = pdata.NewSpanID([8]byte{})
		case SpanParentChild:
			// use existing if available
			if traceID.IsEmpty() {
				traceID = generatePDataTraceID(random)
			}
			if parentID.IsEmpty() {
				parentID = generatePDataSpanID(random)
			}
		}
		spanName := generateSpanName(spanInputs)
		spanList[i] = GenerateSpan(traceID, parentID, spanName, spanInputs, random)
		parentID = spanList[i].SpanID()
		index++
	}
	return spanList, index, nil
}

func generateSpanName(spanInputs *PICTSpanInputs) string {
	return fmt.Sprintf("/%s/%s/%s/%s/%s/%s/%s", spanInputs.Parent, spanInputs.Tracestate, spanInputs.Kind,
		spanInputs.Attributes, spanInputs.Events, spanInputs.Links, spanInputs.Status)
}

// GenerateSpan generates a single pdata.Span based on the input values provided. They are:
//   traceID - the trace ID to use, should not be nil
//   parentID - the parent span ID or nil if it is a root span
//   spanName - the span name, should not be blank
//   spanInputs - the pairwise combination of field value variations for this span
//   random - the random number generator to use in generating ID values
//
// The generated span is returned.
func GenerateSpan(traceID pdata.TraceID, parentID pdata.SpanID, spanName string, spanInputs *PICTSpanInputs,
	random io.Reader) pdata.Span {
	endTime := time.Now().Add(-50 * time.Microsecond)
	span := pdata.NewSpan()
	span.SetTraceID(traceID)
	span.SetSpanID(generatePDataSpanID(random))
	span.SetTraceState(generateTraceState(spanInputs.Tracestate))
	span.SetParentSpanID(parentID)
	span.SetName(spanName)
	span.SetKind(lookupSpanKind(spanInputs.Kind))
	span.SetStartTimestamp(pdata.Timestamp(endTime.Add(-215 * time.Millisecond).UnixNano()))
	span.SetEndTimestamp(pdata.Timestamp(endTime.UnixNano()))
	generateSpanAttributes(spanInputs.Attributes, spanInputs.Status).CopyTo(span.Attributes())
	span.SetDroppedAttributesCount(0)
	generateSpanEvents(spanInputs.Events).CopyTo(span.Events())
	span.SetDroppedEventsCount(0)
	generateSpanLinks(spanInputs.Links, random).CopyTo(span.Links())
	span.SetDroppedLinksCount(0)
	generateStatus(spanInputs.Status).CopyTo(span.Status())
	return span
}

func generateTraceState(tracestate PICTInputTracestate) pdata.TraceState {
	switch tracestate {
	case TraceStateOne:
		return "lasterror=f39cd56cc44274fd5abd07ef1164246d10ce2955"
	case TraceStateFour:
		return "err@ck=80ee5638,rate@ck=1.62,rojo=00f067aa0ba902b7,congo=t61rcWkgMzE"
	case TraceStateEmpty:
		fallthrough
	default:
		return ""
	}
}

func lookupSpanKind(kind PICTInputKind) pdata.SpanKind {
	switch kind {
	case SpanKindClient:
		return pdata.SpanKindClient
	case SpanKindServer:
		return pdata.SpanKindServer
	case SpanKindProducer:
		return pdata.SpanKindProducer
	case SpanKindConsumer:
		return pdata.SpanKindConsumer
	case SpanKindInternal:
		return pdata.SpanKindInternal
	case SpanKindUnspecified:
		fallthrough
	default:
		return pdata.SpanKindUnspecified
	}
}

func generateSpanAttributes(spanTypeID PICTInputAttributes, statusStr PICTInputStatus) pdata.AttributeMap {
	includeStatus := SpanStatusUnset != statusStr
	var attrs map[string]interface{}
	switch spanTypeID {
	case SpanAttrNil:
		attrs = nil
	case SpanAttrEmpty:
		attrs = make(map[string]interface{})
	case SpanAttrDatabaseSQL:
		attrs = generateDatabaseSQLAttributes()
	case SpanAttrDatabaseNoSQL:
		attrs = generateDatabaseNoSQLAttributes()
	case SpanAttrFaaSDatasource:
		attrs = generateFaaSDatasourceAttributes()
	case SpanAttrFaaSHTTP:
		attrs = generateFaaSHTTPAttributes(includeStatus)
	case SpanAttrFaaSPubSub:
		attrs = generateFaaSPubSubAttributes()
	case SpanAttrFaaSTimer:
		attrs = generateFaaSTimerAttributes()
	case SpanAttrFaaSOther:
		attrs = generateFaaSOtherAttributes()
	case SpanAttrHTTPClient:
		attrs = generateHTTPClientAttributes(includeStatus)
	case SpanAttrHTTPServer:
		attrs = generateHTTPServerAttributes(includeStatus)
	case SpanAttrMessagingProducer:
		attrs = generateMessagingProducerAttributes()
	case SpanAttrMessagingConsumer:
		attrs = generateMessagingConsumerAttributes()
	case SpanAttrGRPCClient:
		attrs = generateGRPCClientAttributes()
	case SpanAttrGRPCServer:
		attrs = generateGRPCServerAttributes()
	case SpanAttrInternal:
		attrs = generateInternalAttributes()
	case SpanAttrMaxCount:
		attrs = generateMaxCountAttributes(includeStatus)
	default:
		attrs = generateGRPCClientAttributes()
	}
	return convertMapToAttributeMap(attrs)
}

func generateStatus(statusStr PICTInputStatus) pdata.SpanStatus {
	if SpanStatusUnset == statusStr {
		return pdata.NewSpanStatus()
	}
	spanStatus := pdata.NewSpanStatus()
	spanStatus.SetCode(statusCodeMap[statusStr])
	spanStatus.SetMessage(statusMsgMap[statusStr])
	return spanStatus
}

func generateDatabaseSQLAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeDBSystem] = "mysql"
	attrMap[conventions.AttributeDBConnectionString] = "Server=shopdb.example.com;Database=ShopDb;Uid=billing_user;TableCache=true;UseCompression=True;MinimumPoolSize=10;MaximumPoolSize=50;"
	attrMap[conventions.AttributeDBUser] = "billing_user"
	attrMap[conventions.AttributeNetHostIP] = "192.0.3.122"
	attrMap[conventions.AttributeNetHostPort] = int64(51306)
	attrMap[conventions.AttributeNetPeerName] = "shopdb.example.com"
	attrMap[conventions.AttributeNetPeerIP] = "192.0.2.12"
	attrMap[conventions.AttributeNetPeerPort] = int64(3306)
	attrMap[conventions.AttributeNetTransport] = "IP.TCP"
	attrMap[conventions.AttributeDBName] = "shopdb"
	attrMap[conventions.AttributeDBStatement] = "SELECT * FROM orders WHERE order_id = 'o4711'"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateDatabaseNoSQLAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeDBSystem] = "mongodb"
	attrMap[conventions.AttributeDBUser] = "the_user"
	attrMap[conventions.AttributeNetPeerName] = "mongodb0.example.com"
	attrMap[conventions.AttributeNetPeerIP] = "192.0.2.14"
	attrMap[conventions.AttributeNetPeerPort] = int64(27017)
	attrMap[conventions.AttributeNetTransport] = "IP.TCP"
	attrMap[conventions.AttributeDBName] = "shopDb"
	attrMap[conventions.AttributeDBOperation] = "findAndModify"
	attrMap[conventions.AttributeDBMongoDBCollection] = "products"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSDatasourceAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerDataSource
	attrMap[conventions.AttributeFaaSExecution] = "DB85AF51-5E13-473D-8454-1E2D59415EAB"
	attrMap[conventions.AttributeFaaSDocumentCollection] = "faa-flight-delay-information-incoming"
	attrMap[conventions.AttributeFaaSDocumentOperation] = "insert"
	attrMap[conventions.AttributeFaaSDocumentTime] = "2020-05-09T19:50:06Z"
	attrMap[conventions.AttributeFaaSDocumentName] = "delays-20200509-13.csv"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSHTTPAttributes(includeStatus bool) map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerHTTP
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPHost] = "api.opentelemetry.io"
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	if includeStatus {
		attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
	}
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Safari/605.1.15"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSPubSubAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerPubSub
	attrMap[conventions.AttributeMessagingSystem] = "sqs"
	attrMap[conventions.AttributeMessagingDestination] = "video-views-au"
	attrMap[conventions.AttributeMessagingOperation] = "process"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSTimerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerTimer
	attrMap[conventions.AttributeFaaSExecution] = "73103A4C-E22F-4493-BDE8-EAE5CAB37B50"
	attrMap[conventions.AttributeFaaSTime] = "2020-05-09T20:00:08Z"
	attrMap[conventions.AttributeFaaSCron] = "0/15 * * * *"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSOtherAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerOther
	attrMap["processed.count"] = int64(256)
	attrMap["processed.data"] = 14.46
	attrMap["processed.errors"] = false
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPClientAttributes(includeStatus bool) map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "GET"
	attrMap[conventions.AttributeHTTPURL] = "https://opentelemetry.io/registry/"
	if includeStatus {
		attrMap[conventions.AttributeHTTPStatusCode] = int64(200)
		attrMap[conventions.AttributeHTTPStatusText] = "More Than OK"
	}
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPServerAttributes(includeStatus bool) map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPServerName] = "api22.opentelemetry.io"
	attrMap[conventions.AttributeNetHostPort] = int64(443)
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	if includeStatus {
		attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
	}
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"
	attrMap[conventions.AttributeHTTPRoute] = "/blog/posts"
	attrMap[conventions.AttributeHTTPClientIP] = "2001:506:71f0:16e::1"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMessagingProducerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeMessagingSystem] = "nats"
	attrMap[conventions.AttributeMessagingDestination] = "time.us.east.atlanta"
	attrMap[conventions.AttributeMessagingDestinationKind] = "topic"
	attrMap[conventions.AttributeMessagingMessageID] = "AA7C5438-D93A-43C8-9961-55613204648F"
	attrMap["messaging.sequence"] = int64(1)
	attrMap[conventions.AttributeNetPeerIP] = "10.10.212.33"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMessagingConsumerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeMessagingSystem] = "kafka"
	attrMap[conventions.AttributeMessagingDestination] = "infrastructure-events-zone1"
	attrMap[conventions.AttributeMessagingOperation] = "receive"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:4de0:c223:a800:4e87"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateGRPCClientAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeRPCService] = "PullRequestsService"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:4de0:c223:a800:4e87"
	attrMap[conventions.AttributeNetHostPort] = int64(8443)
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateGRPCServerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeRPCService] = "PullRequestsService"
	attrMap[conventions.AttributeNetPeerIP] = "192.168.1.70"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateInternalAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap["parameters"] = "account=7310,amount=1817.10"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMaxCountAttributes(includeStatus bool) map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPHost] = "api.opentelemetry.io"
	attrMap[conventions.AttributeNetHostName] = "api22.opentelemetry.io"
	attrMap[conventions.AttributeNetHostIP] = "2600:1700:1f00:11c0:1ced:afa5:fd88:9d48"
	attrMap[conventions.AttributeNetHostPort] = int64(443)
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	if includeStatus {
		attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
		attrMap[conventions.AttributeHTTPStatusText] = "Created"
	}
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"
	attrMap[conventions.AttributeHTTPRoute] = "/blog/posts"
	attrMap[conventions.AttributeHTTPClientIP] = "2600:1700:1f00:11c0:1ced:afa5:fd77:9d01"
	attrMap[conventions.AttributePeerService] = "IdentifyImageService"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:1ced:afa5:fd77:9ddc"
	attrMap[conventions.AttributeNetPeerPort] = int64(39111)
	attrMap["ai-sampler.weight"] = 0.07
	attrMap["ai-sampler.absolute"] = false
	attrMap["ai-sampler.maxhops"] = int64(6)
	attrMap["application.create.location"] = "https://api.opentelemetry.io/blog/posts/806673B9-4F4D-4284-9635-3A3E3E3805BE"
	stages := pdata.NewAttributeValueArray()
	stages.ArrayVal().Append(pdata.NewAttributeValueString("Launch"))
	stages.ArrayVal().Append(pdata.NewAttributeValueString("Injestion"))
	stages.ArrayVal().Append(pdata.NewAttributeValueString("Validation"))
	attrMap["application.stages"] = stages
	subMap := pdata.NewAttributeMap()
	subMap.InsertBool("UIx", false)
	subMap.InsertBool("UI4", true)
	subMap.InsertBool("flow-alt3", false)
	attrMap["application.abflags"] = subMap
	attrMap["application.thread"] = "proc-pool-14"
	attrMap["application.session"] = ""
	attrMap["application.persist.size"] = int64(1172184)
	attrMap["application.queue.size"] = int64(0)
	attrMap["application.job.id"] = "0E38800B-9C4C-484E-8F2B-C7864D854321"
	attrMap["application.service.sla"] = 0.34
	attrMap["application.service.slo"] = 0.55
	attrMap[conventions.AttributeEnduserID] = "unittest"
	attrMap[conventions.AttributeEnduserRole] = "poweruser"
	attrMap[conventions.AttributeEnduserScope] = "email profile administrator"
	return attrMap
}

func generateSpanEvents(eventCnt PICTInputSpanChild) pdata.SpanEventSlice {
	spanEvents := pdata.NewSpanEventSlice()
	if SpanChildCountNil == eventCnt {
		return spanEvents
	}
	listSize := calculateListSize(eventCnt)
	for i := 0; i < listSize; i++ {
		spanEvents.Append(generateSpanEvent(i))
	}
	return spanEvents
}

func generateSpanLinks(linkCnt PICTInputSpanChild, random io.Reader) pdata.SpanLinkSlice {
	spanLinks := pdata.NewSpanLinkSlice()
	if SpanChildCountNil == linkCnt {
		return spanLinks
	}
	listSize := calculateListSize(linkCnt)
	for i := 0; i < listSize; i++ {
		spanLinks.Append(generateSpanLink(random, i))
	}
	return spanLinks
}

func calculateListSize(listCnt PICTInputSpanChild) int {
	switch listCnt {
	case SpanChildCountOne:
		return 1
	case SpanChildCountTwo:
		return 2
	case SpanChildCountEight:
		return 8
	case SpanChildCountEmpty:
		fallthrough
	default:
		return 0
	}
}

func generateSpanEvent(index int) pdata.SpanEvent {
	spanEvent := pdata.NewSpanEvent()
	t := time.Now().Add(-75 * time.Microsecond)
	name, attributes := generateEventNameAndAttributes(index)
	spanEvent.SetTimestamp(pdata.Timestamp(t.UnixNano()))
	spanEvent.SetName(name)
	attributes.CopyTo(spanEvent.Attributes())
	spanEvent.SetDroppedAttributesCount(0)
	return spanEvent
}

func generateEventNameAndAttributes(index int) (string, pdata.AttributeMap) {
	switch index % 4 {
	case 0, 3:
		attrMap := make(map[string]interface{})
		if index%2 == 0 {
			attrMap[conventions.AttributeMessageType] = "SENT"
		} else {
			attrMap[conventions.AttributeMessageType] = "RECEIVED"
		}
		attrMap[conventions.AttributeMessageID] = int64(index / 4)
		attrMap[conventions.AttributeMessageCompressedSize] = int64(17 * index)
		attrMap[conventions.AttributeMessageUncompressedSize] = int64(24 * index)
		return "message", convertMapToAttributeMap(attrMap)
	case 1:
		return "custom", convertMapToAttributeMap(map[string]interface{}{
			"app.inretry":  true,
			"app.progress": 0.6,
			"app.statemap": "14|5|202"})
	default:
		return "annotation", pdata.NewAttributeMap()
	}
}

func generateSpanLink(random io.Reader, index int) pdata.SpanLink {
	spanLink := pdata.NewSpanLink()
	spanLink.SetTraceID(generatePDataTraceID(random))
	spanLink.SetSpanID(generatePDataSpanID(random))
	spanLink.SetTraceState("")
	generateLinkAttributes(index).CopyTo(spanLink.Attributes())
	spanLink.SetDroppedAttributesCount(0)
	return spanLink
}

func generateLinkAttributes(index int) pdata.AttributeMap {
	if index%4 == 2 {
		return pdata.NewAttributeMap()
	}
	attrMap := generateMessagingConsumerAttributes()
	if index%4 == 1 {
		attrMap["app.inretry"] = true
		attrMap["app.progress"] = 0.6
		attrMap["app.statemap"] = "14|5|202"
	}
	return convertMapToAttributeMap(attrMap)
}
