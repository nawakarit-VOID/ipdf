// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

func (m MyTheme) Color(name fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if v == theme.VariantDark {
		//theme black
		switch name {
		case theme.ColorNameButton: //สีพื้นหลังปุ่ม
			return color.NRGBA{R: 0, G: 0, B: 0, A: 50}

		case theme.ColorNameBackground: //สีพื้นหลังสุด *ถ้าไม่มีภาพขั้นกลาง
			return color.NRGBA{R: 0, G: 0, B: 0, A: 255}

		case theme.ColorNameShadow:
			return color.NRGBA{0, 0, 0, 80}
			//select
		case theme.ColorNameInputBackground:
			return color.NRGBA{50, 50, 50, 50} // พื้นหลัง select

		case theme.ColorNameForeground:
			return color.NRGBA{50, 50, 50, 50} // ตัวอักษร

		case theme.ColorNameHover:
			return color.NRGBA{0, 0, 0, 80} // hover

		case theme.ColorNameFocus:
			return color.NRGBA{0, 0, 0, 100} // ตอนคลิก
			//prog
		case theme.ColorNamePrimary:
			return color.NRGBA{0, 0, 0, 50} // สีแท่ง progress

		}
	} else {
		//theme white
		switch name {
		case theme.ColorNameButton: //สีปุ่มปกติ
			return color.NRGBA{R: 50, G: 50, B: 50, A: 50}

		case theme.ColorNameHover: //ตอนเอาเมาส์ชี้
			return color.NRGBA{255, 255, 255, 40}

		case theme.ColorNameFocus: //สีเด่น (progress bar / highlight / ปุ่มสำคัญ)
			return color.NRGBA{50, 50, 50, 40}

		case theme.ColorNameBackground: //สีพื้นหลังสุด *ถ้าไม่มีภาพขั้นกลาง
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}

		case theme.ColorNameShadow: //เงาของทุกสิ่ง
			return color.NRGBA{255, 255, 255, 40}
			//select
		case theme.ColorNameInputBackground: //พื้นหลังช่อง input / select
			return color.NRGBA{50, 50, 50, 50}

		case theme.ColorNameForeground: // สีตัวอักษร/ไอคอน
			return color.Black

		case theme.ColorNamePrimary: // สีแท่ง progress
			return color.NRGBA{0, 0, 0, 50}

		}
	}
	return theme.DefaultTheme().Color(name, v)
}

// ต้องมีครบ
func (m MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return myFont
	//return theme.DefaultTheme().Font(s)
}
func (m MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (m MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
