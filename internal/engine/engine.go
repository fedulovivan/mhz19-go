package engine

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/fedulovivan/mhz19-go/internal/engine/actions"
	"github.com/fedulovivan/mhz19-go/internal/engine/conditions"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type GetProviderFn func(ch types.ChannelType) types.ChannelProvider

var _ types.ServiceSupplier = (*engine)(nil)

var BaseTag = utils.NewTag(logger.ENGINE)

type engine struct {
	providers       []types.ChannelProvider
	rules           []types.Rule
	messageService  types.MessagesService
	devicesService  types.DevicesService
	ldmService      types.LdmService
	rulesMu         sync.RWMutex
	queuesContainer *message_queue.Container
}

func NewEngine() *engine {
	return &engine{}
}

func (e *engine) SetQueuesContainer(queuesContainer *message_queue.Container) {
	e.queuesContainer = queuesContainer
}

func (e *engine) SetProviders(providers ...types.ChannelProvider) {
	e.providers = providers
	providerNames := make([]string, len(providers))
	for i, p := range providers {
		providerNames[i] = fmt.Sprintf("%T", p)
	}
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

func (e *engine) GetMessagesService() types.MessagesService {
	return e.messageService
}
func (e *engine) GetDevicesService() types.DevicesService {
	return e.devicesService
}

func (e *engine) Start() {
	for _, p := range e.providers {
		go func(provider types.ChannelProvider) {
			provider.Init()
			for message := range provider.Messages() {
				e.rulesMu.RLock()
				e.handleMessage(message, e.rules)
				e.rulesMu.RUnlock()
			}
		}(p)
	}
	slog.Debug(BaseTag.F("Started"))
}

func (e *engine) GetProvider(ct types.ChannelType) types.ChannelProvider {
	for _, provider := range e.providers {
		if provider.Channel() == ct {
			return provider
		}
	}
	panic(fmt.Sprintf("%v provider not found", ct))
}

func (e *engine) Stop() {
	for _, s := range e.providers {
		s.Stop()
	}
}

func (e *engine) invokeConditionFunc(
	mt types.MessageCompound,
	c types.Condition,
	baseTag utils.Tag,
) bool {
	start := time.Now()
	impl := conditions.Get(c.Fn)
	fnString := c.Fn.String()
	if c.Not {
		fnString = "Not" + fnString
	}
	tag := baseTag.With("Condition=%d %s", c.Id, fnString)
	slog.Debug(tag.F("Started"), "args", c.Args)
	res, err := impl(mt, c.Args, tag)
	elapsed := time.Since(start)
	if err == nil {
		if c.Not {
			res = !res
		}
		slog.Debug(tag.F("Completed in %s", elapsed), "res", res)
		return res
	} else {
		slog.Error(tag.F("Failed in %s", elapsed), "err", err)
		counters.Inc(counters.ERRORS_ALL)
		counters.Errors.WithLabelValues(logger.MOD_CONDS).Inc()
		return false
	}
}

func (e *engine) invokeActionFunc(compound types.MessageCompound, a types.Action, tag utils.Tag) {
	start := time.Now()
	impl := actions.Get(a.Fn)
	atag := tag.With("Action=%d %s", a.Id, a.Fn.String())
	slog.Debug(atag.F("Started"), "args", a.Args)
	err := impl(compound, a.Args, a.Mapping, e, atag)
	elapsed := time.Since(start)
	if err != nil {
		slog.Error(atag.F("Failed in %s", elapsed), "err", err)
		counters.Inc(counters.ERRORS_ALL)
		counters.Errors.WithLabelValues(logger.MOD_ACTIONS).Inc()
	} else {
		slog.Debug(atag.F("Completed in %s", elapsed))
	}
}

func (e *engine) matchesListSome(mtcb types.GetCompoundForOtherDeviceId, cc []types.Condition, tag utils.Tag) bool {
	for _, c := range cc {
		if e.matchesCondition(mtcb, c, tag) {
			return true
		}
	}
	return false
}

func (e *engine) matchesListEvery(mtcb types.GetCompoundForOtherDeviceId, cc []types.Condition, tag utils.Tag) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !e.matchesCondition(mtcb, c, tag) {
			return false
		}
	}
	return true
}

func (e *engine) matchesCondition(mtcb types.GetCompoundForOtherDeviceId, c types.Condition, tag utils.Tag) bool {
	withFn := c.Fn != 0
	withList := len(c.Nested) > 0
	if !withFn && !withList {
		return false
	} else if withFn && !withList {
		return e.invokeConditionFunc(
			mtcb(c.OtherDeviceId), c, tag,
		)
	} else if withList && !withFn {
		if c.Or {
			return e.matchesListSome(mtcb, c.Nested, tag)
		} else {
			return e.matchesListEvery(mtcb, c.Nested, tag)
		}
	} else {
		panic("unexpected conditions")
	}
}

func (e *engine) executeActions(compound types.MessageCompound, r types.Rule, tag utils.Tag) {
	slog.Debug(tag.F("going to execute %d actions", len(r.Actions)))
	for _, a := range r.Actions {
		go e.invokeActionFunc(compound, a, tag)
	}
}

// called simultaneusly upon receiving messages from all providers
// should have thread-safe implementation
func (e *engine) handleMessage(m types.Message, rules []types.Rule) {

	if m.Id == 0 || m.Timestamp.IsZero() {
		panic("message must have Id and Timestamp initialised")
	}

	counters.MessagesByChannel.WithLabelValues(m.ChannelType.String()).Inc()
	defer counters.TimeSince(time.Now(), counters.MESSAGES_HANDLED)
	defer func(t *prometheus.Timer) {
		t.ObserveDuration()
	}(prometheus.NewTimer(counters.MessagesHandled))

	e.rulesMu.RLock()
	defer e.rulesMu.RUnlock()

	mtag := BaseTag.With("Msg=%d", m.Id)

	defer utils.TimeTrack(mtag.F, time.Now(), "handleMessage")

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

	matches := 0
	rulesCnt := len(rules)
	withRules := rulesCnt > 0
	if withRules {
		slog.Debug(mtag.F("Matching against %v rules", rulesCnt))
	} else {
		slog.Warn(mtag.F("No rules"))
	}

	for _, r := range rules {
		rtag := mtag.With("Rule=%d", r.Id)
		qtag := BaseTag.With("Rule=%d", r.Id)
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
			} else {
				compound.Curr = &m
				if e.ldmService.Has(ldmKey) {
					tmp := e.ldmService.Get(ldmKey)
					compound.Prev = &tmp
				}
			}
			return compound
		}
		if e.matchesCondition(mtcb, r.Condition, rtag) {
			counters.IncRule(r.Id)
			counters.Rules.WithLabelValues(r.Name).Inc()
			slog.Debug(rtag.F("matches👌"), "name", r.Name)
			matches++
			if r.Throttle.Duration == 0 {
				compound := mtcb("")
				e.executeActions(compound, r, rtag)
			} else {
				key := message_queue.NewKey(m.DeviceClass, m.DeviceId, r.Id)
				if !e.queuesContainer.HasQueue(key) {
					e.queuesContainer.CreateQueue(key, r.Throttle.Duration, func(mm []types.Message) {
						slog.Debug(qtag.F("message queue is flushed now"), "key", key, "mm", len(mm))
						e.executeActions(types.MessageCompound{Queued: mm}, r, qtag)
					})
				}
				e.queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(rtag.F(
					"message and execution of %d actions were queued for %s",
					len(r.Actions),
					r.Throttle.Duration,
				), "key", key)
			}
		}
	}

	if withRules {
		if matches == 0 {
			slog.Warn(mtag.F("No one matching rule found"))
		} else {
			slog.Debug(mtag.F("%v out of %v rules were matched", matches, rulesCnt))
		}
	}

	if !isBridge && !sonoffAnnounce && !isSystem {
		e.ldmService.Set(ldmKey, m)
	}
}
