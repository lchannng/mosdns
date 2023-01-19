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

package coremain

import (
	"github.com/IrineSistiana/mosdns/v5/mlog"
)

type Config struct {
	Log     mlog.LogConfig `yaml:"log"`
	Include []string       `yaml:"include"`
	Plugins []PluginConfig `yaml:"plugins"`
	API     APIConfig      `yaml:"api"`
}

// PluginConfig represents a plugin config
type PluginConfig struct {
	// Tag for this plugin. Optional. If omitted, this plugin will
	// be registered with a random tag.
	Tag string `yaml:"tag"`

	// Type, required.
	Type string `yaml:"type"`

	// Args, might be required by some plugins.
	// The type of Args is depended on RegNewPluginFunc.
	// If it's a map[string]any, it will be converted by mapstruct.
	Args any `yaml:"args"`
}

type APIConfig struct {
	HTTP string `yaml:"http"`
}
