// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"fyne.io/fyne/v2/widget"
)

type I18n struct {
	currentLang string
	data        map[string]map[string]string
	observers   []func()
	mu          sync.RWMutex
}

// สร้าง instance
func New(defaultLang string) *I18n {
	return &I18n{
		currentLang: defaultLang,
		data:        make(map[string]map[string]string),
	}
}

// โหลดไฟล์ภาษา (lang/en.json)
func (i *I18n) Load(lang string, path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var m map[string]string
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	i.mu.Lock()
	i.data[lang] = m
	i.mu.Unlock()

	return nil
}

// เปลี่ยนภาษา + notify UI
func (i *I18n) SetLang(lang string) {
	i.mu.Lock()
	i.currentLang = lang
	i.mu.Unlock()

	i.notify()
}

// แปลข้อความ
func (i *I18n) T(key string) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// current lang
	if val, ok := i.data[i.currentLang][key]; ok {
		return val
	}

	// fallback -> en
	if val, ok := i.data["en"][key]; ok {
		return val
	}

	return key
}

// รองรับ format เช่น %s %d
func (i *I18n) Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(i.T(key), args...)
}

// subscribe UI
func (i *I18n) Subscribe(fn func()) {
	i.mu.Lock()
	i.observers = append(i.observers, fn)
	i.mu.Unlock()
}

// notify ทุกตัว
func (i *I18n) notify() {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for _, fn := range i.observers {
		fn()
	}
}

// Label ที่เปลี่ยนภาษาตามอัตโนมัติ
func NewLabel(i *I18n, key string) *widget.Label {
	lbl := widget.NewLabel(i.T(key))

	i.Subscribe(func() {
		lbl.SetText(i.T(key))
	})

	return lbl
}

// Button
func NewButton(i *I18n, key string, tapped func()) *widget.Button {
	btn := widget.NewButton(i.T(key), tapped)

	i.Subscribe(func() {
		btn.SetText(i.T(key))
	})

	return btn
}
