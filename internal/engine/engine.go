package engine

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/counters"

	"github.com/fedulovivan/mhz19-go/internal/engine/actions"
	"github.com/fedulovivan/mhz19-go/internal/engine/conditions"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/samber/lo"
)

var queuesContainer = message_queue.NewContainer()

type GetProviderFn func(ch types.ChannelType) types.ChannelProvider

var _ types.Engine = (*engine)(nil)
var _ types.EngineAsSupplier = (*engine)(nil)

var BaseTag = logger.NewTag(logger.ENGINE)

type engine struct {
	providers      []types.ChannelProvider
	rules          []types.Rule
	messageService types.MessagesService
	devicesService types.DevicesService
	ldmService     types.LdmService
	rulesMu        sync.RWMutex
}

func NewEngine() types.Engine {
	return &engine{}
}

func (e *engine) SetProviders(providers ...types.ChannelProvider) {
	e.providers = providers
	providerNames := lo.Map(providers, func(p types.ChannelProvider, index int) string {
		return fmt.Sprintf("%T", p)
	})
	slog.Debug(BaseTag.F(
		"%d provider(s) were set: %s",
		len(providers),
		strings.Join(providerNames, ", "),
	))
}
func (e *engine) SetMessagesService(s types.MessagesService) {
	e.messageService = s
}
func (e *engine) SetDevicesService(s types.DevicesService) {
	e.devicesService = s
}
func (e *engine) SetLdmService(r types.LdmService) {
	e.ldmService = r
}
func (e *engine) AppendRules(rules ...types.Rule) {
	e.rulesMu.Lock()
	defer e.rulesMu.Unlock()
	e.rules = append(e.rules, rules...)
	slog.Debug(BaseTag.F("AppendRules"), "appended", len(rules), "total", len(e.rules))
}

func (e *engine) DeleteRule(ruleId int) {
	e.rulesMu.Lock()
	defer e.rulesMu.Unlock()
	before := len(e.rules)
	e.rules = slices.DeleteFunc(e.rules, func(r types.Rule) bool {
		return r.Id == ruleId
	})
	after := len(e.rules)
	slog.Debug(BaseTag.F("DeleteRule"), "deleted", before-after, "total", after)
}

func (e *engine) MessagesService() types.MessagesService {
	return e.messageService
}
func (e *engine) DevicesService() types.DevicesService {
	return e.devicesService
}

func (e *engine) Start() {
	for _, p := range e.providers {
		go func(provider types.ChannelProvider) {
			provider.Init()
			for message := range provider.Messages() {
				e.HandleMessage(message, e.rules)
			}
		}(p)
	}
	slog.Debug(BaseTag.F("Started"))
}

func (e *engine) FindProvider(ct types.ChannelType) types.ChannelProvider {
	for _, provider := range e.providers {
		if provider.Channel() == ct {
			return provider
		}
	}
	panic(fmt.Sprintf("%v provider is not found", ct))
}

func (e *engine) Stop() {
	for _, s := range e.providers {
		s.Stop()
	}
}

func (e *engine) InvokeConditionFunc(mt types.MessageTuple, fn types.CondFn, not bool, args types.Args, tag logger.Tag) bool {
	impl := conditions.Get(fn)
	fnString := fn.String()
	if not {
		fnString = "Not" + fnString
	}
	tag = tag.With("condition=%s", fnString)
	slog.Debug(tag.F("Start"), "args", args)
	res, err := impl(mt, args)
	if err == nil {
		if not {
			res = !res
		}
		slog.Debug(tag.F("End"), "res", res)
		return res
	} else {
		slog.Error(tag.F("Fail"), "err", err)
		counters.Inc(counters.ERRORS)
		return false
	}
}

func (e *engine) InvokeActionFunc(mm []types.Message, a types.Action, tag logger.Tag) {
	impl := actions.Get(a.Fn)
	tag = tag.With("action=%s", a.Fn.String())
	slog.Debug(tag.F("Start"), "args", a.Args)
	err := impl(mm, a.Args, a.Mapping, e, tag)
	if err != nil {
		slog.Error(tag.F("Fail"), "err", err)
		counters.Inc(counters.ERRORS)
	} else {
		slog.Debug(tag.F("End"))
	}
}

func (e *engine) MatchesListSome(mtcb types.MessageTupleFn, cc []types.Condition, tag logger.Tag) bool {
	for _, c := range cc {
		if e.MatchesCondition(mtcb, c, tag) {
			return true
		}
	}
	return false
}

func (e *engine) MatchesListEvery(mtcb types.MessageTupleFn, cc []types.Condition, tag logger.Tag) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !e.MatchesCondition(mtcb, c, tag) {
			return false
		}
	}
	return true
}

func (e *engine) MatchesCondition(mtcb types.MessageTupleFn, c types.Condition, tag logger.Tag) bool {
	withFn := c.Fn != 0
	withList := len(c.Nested) > 0
	if !withFn && !withList {
		return false
	} else if withFn && !withList {
		return e.InvokeConditionFunc(
			mtcb(c.OtherDeviceId), c.Fn, c.Not, c.Args, tag,
		)
	} else if withList && !withFn {
		if c.Or {
			return e.MatchesListSome(mtcb, c.Nested, tag)
		} else {
			return e.MatchesListEvery(mtcb, c.Nested, tag)
		}
	} else {
		panic("unexpected conditions")
	}
}

func (e *engine) ExecuteActions(mm []types.Message, r types.Rule, tag logger.Tag) {
	slog.Debug(tag.F("going to execute %v actions", len(r.Actions)))
	for _, a := range r.Actions {
		go e.InvokeActionFunc(mm, a, tag)
	}
}

// called simultaneusly upon receiving messages from all providers
// should have thread-safe implementation
func (e *engine) HandleMessage(m types.Message, rules []types.Rule) {
	defer utils.TimeTrack(BaseTag.F, time.Now(), "HandleMessage")
	e.rulesMu.RLock()
	defer e.rulesMu.RUnlock()
	if m.Id == 0 || m.Timestamp.IsZero() {
		panic("message must have Id and Timestamp initialised")
	}
	counters.Inc(counters.MESSAGES_RECEIVED)
	tag := BaseTag.With("Msg=%d", m.Id)
	p := m.Payload
	isSystem := m.DeviceClass == types.DEVICE_CLASS_SYSTEM
	isBridge := m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE
	sonoffAnnounce := m.DeviceClass == types.DEVICE_CLASS_SONOFF_ANNOUNCE
	if isBridge {
		p = "<too big to render>"
	}
	slog.Debug(tag.F("New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
		"FromEndDevice", m.FromEndDevice,
	)
	rulesCnt := len(rules)
	if rulesCnt == 0 {
		slog.Warn(tag.F("No rules"))
		return
	}
	ldmKey := e.ldmService.NewKey(m.DeviceClass, m.DeviceId)
	slog.Debug(tag.F("Matching against %v rules", rulesCnt))
	matches := 0
	for _, r := range rules {
		tag := tag.With("Rule=%d", r.Id)
		if r.Disabled {
			slog.Debug(tag.F("is disabled, skipping"))
			continue
		}
		var mtcb = func(otherDeviceId types.DeviceId) types.MessageTuple {
			tuple := types.MessageTuple{}
			takeOtherDeviceMessage := len(otherDeviceId) > 0
			if takeOtherDeviceMessage {
				slog.Warn(tag.F("requesting message for otherDeviceId=%v", otherDeviceId))
				otherLdmKey := e.ldmService.NewKey(m.DeviceClass, otherDeviceId)
				if e.ldmService.Has(otherLdmKey) {
					tmp := e.ldmService.Get(otherLdmKey)
					tuple.Curr = &tmp
				}
			} else {
				tuple.Curr = &m
				if e.ldmService.Has(ldmKey) {
					tmp := e.ldmService.Get(ldmKey)
					tuple.Prev = &tmp
				}
			}
			return tuple
		}
		if e.MatchesCondition(mtcb, r.Condition, tag) {
			counters.IncRule(r.Id)
			slog.Debug(tag.F("matches👌"), "name", r.Name)
			matches++
			if r.Throttle.Duration == 0 {
				e.ExecuteActions([]types.Message{m}, r, tag)
			} else {
				key := message_queue.NewKey(m.DeviceClass, m.DeviceId, r.Id)
				queueTag := BaseTag.With("Rule=%d", r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle.Duration, func(mm []types.Message) {
						slog.Debug(queueTag.F("message queue is flushed now"), "key", key, "mm", len(mm))
						e.ExecuteActions(mm, r, queueTag)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(tag.F("message was queued for %s", r.Throttle.Duration))
			}
		}
	}
	if matches == 0 {
		slog.Warn(tag.F("No one matching rule found"))
	} else {
		slog.Debug(tag.F("%v out of %v rules were matched", matches, len(rules)))
	}
	if !isBridge && !sonoffAnnounce && !isSystem {
		e.ldmService.Set(ldmKey, m)
	}
}
