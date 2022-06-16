package msg_matcher

import (
	"context"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/handler"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/pkg/dnsutils"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/pkg/matcher/domain"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/pkg/matcher/elem"
	"github.com/IrineSistiana/mosdns/v3/dispatcher/pkg/matcher/netlist"
	"github.com/miekg/dns"
	"net"
	"testing"
)

func TestClientIPMatcher_Match(t *testing.T) {
	type fields struct {
		ipMatcher netlist.Matcher
	}
	type args struct {
		qCtx *handler.Context
	}

	nl := netlist.NewList()
	if err := netlist.LoadFromText(nl, "127.0.0.0/24"); err != nil {
		t.Fatal(err)
	}
	nl.Sort()

	msg := new(dns.Msg)
	meta1271 := &handler.RequestMeta{ClientIP: net.ParseIP("127.0.0.1")}
	meta1281 := &handler.RequestMeta{ClientIP: net.ParseIP("128.0.0.1")}
	metaNilAddr := &handler.RequestMeta{}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMatched bool
		wantErr     bool
	}{
		{"matched", fields{ipMatcher: nl}, args{handler.NewContext(msg, meta1271)}, true, false},
		{"not matched", fields{ipMatcher: nl}, args{handler.NewContext(msg, meta1281)}, false, false},
		{"no meta", fields{ipMatcher: nl}, args{handler.NewContext(msg, nil)}, false, false},
		{"no addr", fields{ipMatcher: nl}, args{handler.NewContext(msg, metaNilAddr)}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewClientIPMatcher(tt.fields.ipMatcher)
			gotMatched, err := m.Match(context.Background(), tt.args.qCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatched != tt.wantMatched {
				t.Errorf("Match() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}

func TestClientECSMatcher_Match(t *testing.T) {
	nl := netlist.NewList()
	if err := netlist.LoadFromText(nl, "127.0.0.0/24"); err != nil {
		t.Fatal(err)
	}
	nl.Sort()

	msg := new(dns.Msg)
	msgWithoutOPT := msg
	msg = new(dns.Msg)
	msg.SetEdns0(512, false)
	msgWithOPT := msg
	msg = new(dns.Msg)
	msg.SetEdns0(512, false)
	opt := msg.IsEdns0()
	dnsutils.AddECS(opt, &dns.EDNS0_SUBNET{Address: net.ParseIP("127.0.0.1")}, false)
	msg1271 := msg
	msg1281 := msg.Copy()
	opt = msg1281.IsEdns0()
	dnsutils.AddECS(opt, &dns.EDNS0_SUBNET{Address: net.ParseIP("128.0.0.1")}, true)

	tests := []struct {
		name        string
		matcher     netlist.Matcher
		qCtx        *handler.Context
		wantMatched bool
		wantErr     bool
	}{
		{"matched", nl, handler.NewContext(msg1271, nil), true, false},
		{"not matched", nl, handler.NewContext(msg1281, nil), false, false},
		{"no ecs", nl, handler.NewContext(msgWithOPT, nil), false, false},
		{"no opt", nl, handler.NewContext(msgWithoutOPT, nil), false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewClientECSMatcher(tt.matcher)
			gotMatched, err := m.Match(context.Background(), tt.qCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatched != tt.wantMatched {
				t.Errorf("Match() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}
func TestQNameMatcher_Match(t *testing.T) {
	dm := domain.NewDomainMatcher[struct{}]()
	dm.Add("com.", struct{}{})

	qm := NewQNameMatcher(dm)
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)
	if !qm.MatchMsg(m) {
		t.Fatal()
	}

	m.SetQuestion("example.xxx.", dns.TypeA)
	if qm.MatchMsg(m) {
		t.Fatal()
	}
}

func TestQTypeMatcher_Match(t *testing.T) {
	em := elem.NewIntMatcher([]int{int(dns.TypeA)})
	qm := NewQTypeMatcher(em)
	m := new(dns.Msg)
	m.SetQuestion(".", dns.TypeA)
	if !qm.MatchMsg(m) {
		t.Fatal()
	}

	m.SetQuestion(".", dns.TypeAAAA)
	if qm.MatchMsg(m) {
		t.Fatal()
	}
}

func TestQClassMatcher_Match(t *testing.T) {
	em := elem.NewIntMatcher([]int{dns.ClassINET})
	qm := NewQClassMatcher(em)
	m := new(dns.Msg)
	m.Question = []dns.Question{{Name: ".", Qtype: dns.TypeA, Qclass: dns.ClassINET}}
	if !qm.MatchMsg(m) {
		t.Fatal()
	}

	m.Question = []dns.Question{{Name: ".", Qtype: dns.TypeA, Qclass: dns.ClassANY}}
	if qm.MatchMsg(m) {
		t.Fatal()
	}
}
