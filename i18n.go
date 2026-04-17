// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	_ "embed"

	"fyne.io/fyne/v2/widget"
)

//////////////////////////////////////////////////
// 🔥 I18N CORE
//////////////////////////////////////////////////

type TranslatableItem interface {
	Update()
}

type I18n struct {
	lang  string
	data  map[string]map[string]string
	items []TranslatableItem
}

func NewI18n(en, th map[string]string) *I18n {
	return &I18n{
		lang: "en",
		data: map[string]map[string]string{
			"en": en,
			"th": th,
		},
	}
}

func (i *I18n) T(key string) string {
	if val, ok := i.data[i.lang][key]; ok {
		return val
	}
	if val, ok := i.data["en"][key]; ok {
		return val
	}
	return key
}

func (i *I18n) Register(item TranslatableItem) {
	i.items = append(i.items, item)
}

func (i *I18n) SetLang(lang string) {
	i.lang = lang
	for _, item := range i.items {
		item.Update()
	}
}

//////////////////////////////////////////////////
// 🔥 LABEL
//////////////////////////////////////////////////

type I18nLabel struct {
	key string
	lbl *widget.Label
	i   *I18n
}

func NewLabel(i *I18n, key string) *widget.Label {
	l := &I18nLabel{
		key: key,
		lbl: widget.NewLabel(i.T(key)),
		i:   i,
	}
	i.Register(l)
	return l.lbl
}

func (l *I18nLabel) Update() {
	l.lbl.SetText(l.i.T(l.key))
}

//////////////////////////////////////////////////
// 🔥 BUTTON
//////////////////////////////////////////////////

type I18nButton struct {
	key string
	btn *widget.Button
	i   *I18n
}

func NewButton(i *I18n, key string, tapped func()) *widget.Button {
	b := &I18nButton{
		key: key,
		btn: widget.NewButton(i.T(key), tapped),
		i:   i,
	}
	i.Register(b)
	return b.btn
}

func (b *I18nButton) Update() {
	b.btn.SetText(b.i.T(b.key))
}

//////////////////////////////////////////////////
// 🔥 SELECT (สำคัญ)
//////////////////////////////////////////////////

type I18nSelect struct {
	keys    []string
	selectW *widget.Select
	i       *I18n
}

func NewSelect(i *I18n, keys []string, changed func(string)) *widget.Select {
	s := &I18nSelect{
		keys: keys,
		i:    i,
	}

	options := make([]string, len(keys))
	for idx, k := range keys {
		options[idx] = i.T(k)
	}

	s.selectW = widget.NewSelect(options, func(val string) {
		for idx, k := range keys {
			if options[idx] == val {
				changed(k)
				break
			}
		}
	})

	i.Register(s)
	return s.selectW
}

func (s *I18nSelect) Update() {
	options := make([]string, len(s.keys))
	for idx, k := range s.keys {
		options[idx] = s.i.T(k)
	}
	s.selectW.Options = options
	s.selectW.Refresh()
}
