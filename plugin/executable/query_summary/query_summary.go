/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package query_summary

import (
	"context"
	"github.com/IrineSistiana/mosdns/v5/pkg/query_context"
	"github.com/IrineSistiana/mosdns/v5/plugin/executable/sequence"
	"go.uber.org/zap"
)

const (
	PluginType = "query_summary"
)

func init() {
	sequence.MustRegExecQuickSetup(PluginType, QuickSetup)
}

var _ sequence.RecursiveExecutable = (*SummaryLogger)(nil)

type SummaryLogger struct {
	l   *zap.Logger
	msg string
}

// QuickSetup format: [msg_title]
func QuickSetup(bq sequence.BQ, s string) (any, error) {
	return NewSummaryLogger(bq.L(), s), nil
}

// NewSummaryLogger returns a SummaryLogger that logs query info into l.
// l cannot be nil.
// If msg is empty, "query summary" will be used.
func NewSummaryLogger(l *zap.Logger, msg string) *SummaryLogger {
	if len(msg) == 0 {
		msg = "query summary"
	}
	return &SummaryLogger{
		l:   l,
		msg: msg,
	}
}

func (l *SummaryLogger) Exec(ctx context.Context, qCtx *query_context.Context, next sequence.ChainWalker) error {
	err := next.ExecNext(ctx, qCtx)
	l.l.Info(
		l.msg,
		zap.Inline(qCtx),
		zap.Error(err),
	)
	return err
}
