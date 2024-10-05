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
				e.rulesMu.RLock()
				e.HandleMessage(message, e.rules)
				e.rulesMu.RUnlock()
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

func (e *engine) InvokeConditionFunc(
	mt types.MessageCompound,
	c types.Condition,
	baseTag logger.Tag,
) bool {
	// start := time.Now()
	impl := conditions.Get(c.Fn)
	fnString := c.Fn.String()
	if c.Not {
		fnString = "Not" + fnString
	}
	tag := baseTag.With("Condition=%d %s", c.Id, fnString)
	slog.Debug(tag.F("Started"), "args", c.Args)
	res, err := impl(mt, c.Args, tag)
	// elapsed := time.Since(start)
	if err == nil {
		if c.Not {
			res = !res
		}
		slog.Debug(tag.F("Completed" /* , elapsed */), "res", res)
		return res
	} else {
		slog.Error(tag.F("Failed" /* , elapsed */), "err", err)
		counters.Inc(counters.ERRORS_ALL)
		return false
	}
}

func (e *engine) InvokeActionFunc(compound types.MessageCompound, a types.Action, tag logger.Tag) {
	start := time.Now()
	impl := actions.Get(a.Fn)
	atag := tag.With("Action=%d %s", a.Id, a.Fn.String())
	slog.Debug(atag.F("Started"), "args", a.Args)
	err := impl(compound, a.Args, a.Mapping, e, atag)
	elapsed := time.Since(start)
	if err != nil {
		slog.Error(atag.F("Failed in %s", elapsed), "err", err)
		counters.Inc(counters.ERRORS_ALL)
	} else {
		slog.Debug(atag.F("Completed in %s", elapsed))
	}
}

func (e *engine) MatchesListSome(mtcb types.GetCompoundForOtherDeviceId, cc []types.Condition, tag logger.Tag) bool {
	for _, c := range cc {
		if e.MatchesCondition(mtcb, c, tag) {
			return true
		}
	}
	return false
}

func (e *engine) MatchesListEvery(mtcb types.GetCompoundForOtherDeviceId, cc []types.Condition, tag logger.Tag) bool {
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

func (e *engine) MatchesCondition(mtcb types.GetCompoundForOtherDeviceId, c types.Condition, tag logger.Tag) bool {
	withFn := c.Fn != 0
	withList := len(c.Nested) > 0
	if !withFn && !withList {
		return false
	} else if withFn && !withList {
		return e.InvokeConditionFunc(
			mtcb(c.OtherDeviceId), c, tag,
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

func (e *engine) ExecuteActions(compound types.MessageCompound, r types.Rule, tag logger.Tag) {
	slog.Debug(tag.F("going to execute %v actions", len(r.Actions)))
	for _, a := range r.Actions {
		go e.InvokeActionFunc(compound, a, tag)
	}
}

// called simultaneusly upon receiving messages from all providers
// should have thread-safe implementation
func (e *engine) HandleMessage(m types.Message, rules []types.Rule) {

	if m.Id == 0 || m.Timestamp.IsZero() {
		panic("message must have Id and Timestamp initialised")
	}

	e.rulesMu.RLock()
	defer e.rulesMu.RUnlock()

	mtag := BaseTag.With("Msg=%d", m.Id)

	defer func(start time.Time) {
		elapsed := utils.TimeTrack(mtag.F, start, "HandleMessage")
		counters.Time(elapsed, counters.MESSAGES_HANDLED)
	}(time.Now())
	defer counters.Inc(counters.MESSAGES_HANDLED)

	isSystem := m.DeviceClass == types.DEVICE_CLASS_SYSTEM
	isBridge := m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE
	sonoffAnnounce := m.DeviceClass == types.DEVICE_CLASS_SONOFF_ANNOUNCE
	ldmKey := e.ldmService.NewKey(m.DeviceClass, m.DeviceId)

	p := m.Payload
	if isBridge {
		p = "<too big to render>"
	}

	slog.Debug(mtag.F("New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
		"FromEndDevice", m.FromEndDevice,
	)

	rulesCnt := len(rules)
	if rulesCnt == 0 {
		slog.Warn(mtag.F("No rules"))
	} else {
		slog.Debug(mtag.F("Matching against %v rules", rulesCnt))
	}

	matches := 0
	for _, r := range rules {
		rtag := mtag.With("Rule=%d", r.Id)
		if r.Disabled {
			slog.Debug(rtag.F("is disabled, skipping"))
			continue
		}
		var mtcb = func(otherDeviceId types.DeviceId) types.MessageCompound {
			compound := types.MessageCompound{}
			takeOtherDeviceMessage := len(otherDeviceId) > 0
			if takeOtherDeviceMessage {
				slog.Debug(rtag.F("requesting message for otherDeviceId=%v", otherDeviceId))
				tmp, err := e.ldmService.GetByDeviceId(otherDeviceId)
				if err == nil {
					slog.Debug(rtag.F("fetched", "message", tmp))
					compound.Curr = &tmp
				} else {
					slog.Warn(rtag.F(err.Error()))
				}
				// otherLdmKey := e.ldmService.NewKey(m.DeviceClass, otherDeviceId)
				// if e.ldmService.Has(otherLdmKey) {
				// 	tmp := e.ldmService.Get(otherLdmKey)
				// 	compound.Curr = &tmp
				// }
			} else {
				compound.Curr = &m
				if e.ldmService.Has(ldmKey) {
					tmp := e.ldmService.Get(ldmKey)
					compound.Prev = &tmp
				}
			}
			return compound
		}
		if e.MatchesCondition(mtcb, r.Condition, rtag) {
			counters.IncRule(r.Id)
			slog.Debug(rtag.F("matches👌"), "name", r.Name)
			matches++
			if r.Throttle.Duration == 0 {
				compound := mtcb("")
				e.ExecuteActions(compound, r, rtag)
			} else {
				key := message_queue.NewKey(m.DeviceClass, m.DeviceId, r.Id)
				qtag := BaseTag.With("Rule=%d", r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle.Duration, func(mm []types.Message) {
						slog.Debug(qtag.F("message queue is flushed now"), "key", key, "mm", len(mm))
						e.ExecuteActions(types.MessageCompound{Queued: mm}, r, qtag)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(rtag.F("message was queued for %s", r.Throttle.Duration))
			}
		}
	}

	if rulesCnt > 0 {
		if matches == 0 {
			slog.Warn(mtag.F("No one matching rule found"))
		} else {
			slog.Debug(mtag.F("%v out of %v rules were matched", matches, len(rules)))
		}
	}

	if !isBridge && !sonoffAnnounce && !isSystem {
		e.ldmService.Set(ldmKey, m)
	}
}
