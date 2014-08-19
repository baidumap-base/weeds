package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Weed struct {
	Time   time.Time // 时间
	Sv     string    // 版本
	Os     string    // 系统
	Ov     string    // API 等级
	Cuid   string    // 用户唯一 id
	Net    int       // 网络类型
	Sw     int       // 屏幕宽度
	Sh     int       // 屏幕高度
	Act    string    // 类别
	Loc    Location  // 定位点
	Page   string    // 页面栈
	Detail Stack     // 详情
	Pd     string
	Ch     string
	Mb     string
	Ver    string
	Lt     string
	Tm     string
}

type Location struct {
	X int
	Y int
}

type Stack struct {
	Reason Reason
	Traces []Trace
}

type Reason struct {
	JavaClass string
	Message   string
}

type Trace struct {
	File     string // 文件名
	Line     int    // 行号
	Function string // 具体出错位置
}

func parseWeed(l string) (*Weed, error) {
	ll := len(l)
	if ll == 0 {
		return nil, errors.New("parseWeed(): l is empty")
	}
	line := l[1 : ll-1]
	entities := strings.Split(line, "][")
	if len(entities) == 0 {
		return nil, errors.New("parseWeed(): parse error")
	}
	w := Weed{}
	for _, ele := range entities {
		idx := strings.IndexRune(ele, '=')
		if idx == -1 {
			continue
		}

		var value string
		if idx == len(ele)-1 {
			value = ""
		} else {
			value = ele[idx+1:]
		}

		switch key := ele[:idx]; key {
		case "time":
			location, _ := time.LoadLocation("Local")
			w.Time, _ = time.ParseInLocation("20060102150405", value, location)
		case "sv":
			w.Sv = value
		case "os":
			w.Os = value
		case "sw":
			w.Sw, _ = strconv.Atoi(value)
		case "sh":
			w.Sh, _ = strconv.Atoi(value)
		case "pd":
			w.Pd = value
		case "ch":
			w.Ch = value
		case "mb":
			w.Mb = value
		case "ov":
			w.Ov = value
		case "ver":
			w.Ver = value
		case "cuid":
			w.Cuid = value
		case "net":
			w.Net, _ = strconv.Atoi(value)
		case "lt":
			w.Lt = value
		case "tm":
			w.Tm = value
		case "act":
			w.Act = value
		case "ActParam":
			parseActParam(value, &w)
		}
	}
	return &w, nil
}

func parseActParam(c string, w *Weed) {
	content := c[1 : len(c)-1]
	entities := strings.Split(content, "}{")
	for _, ele := range entities {
		idx := strings.IndexRune(ele, '=')
		if idx == -1 {
			continue
		}

		key := ele[:idx]
		var value string
		if len(ele) == idx+1 {
			value = ""
		} else {
			value = ele[idx+1:]
		}

		switch key {
		case "detail":
			parseDetail(value, w)
		case "reason":
			if len(w.Detail.Reason.JavaClass) == 0 {
				items := strings.Split(value, ": ")
				if len(items) >= 2 {
					w.Detail.Reason.JavaClass = items[len(items)-2]
					w.Detail.Reason.Message = items[len(items)-1]
				} else {
					w.Detail.Reason.JavaClass = items[0]
				}
			}
		case "locx":
			w.Loc.X, _ = strconv.Atoi(value)
		case "locy":
			w.Loc.Y, _ = strconv.Atoi(value)
		case "pages":
			w.Page = value
		}
	}
}

func parseDetail(v string, w *Weed) {
	ss := strings.Split(v, "<br>")
	if w.Detail.Traces == nil {
		w.Detail.Traces = make([]Trace, 0, 16)
	}
	for _, s := range ss {
		if len(s) == 0 {
			continue
		}

		if strings.HasPrefix(s, "at ") {
			idx := strings.Index(s, "(")
			if idx == -1 {
				w.Detail.Traces = append(w.Detail.Traces, Trace{Function: s})
			} else {
				cs := s[idx+1 : len(s)-1]
				cidx := strings.Index(cs, ":")
				if cidx != -1 {
					line, _ := strconv.Atoi(cs[cidx+1:])
					w.Detail.Traces = append(w.Detail.Traces,
						Trace{Function: s[:idx], File: cs[:cidx], Line: line})
				} else {
					w.Detail.Traces = append(w.Detail.Traces, Trace{Function: s[:idx]})
				}
			}
		} else {
			idx := strings.Index(s, ": ")
			if idx == -1 {
				w.Detail.Reason.JavaClass = s
			} else {
				w.Detail.Reason.JavaClass = s[:idx]
				w.Detail.Reason.Message = s[idx+1:]
			}
		}
	}
}

func stackEquals(a, b *Stack) bool {
	if a == nil {
		if b == nil {
			return true
		} else {
			return false
		}
	}

	if b == nil {
		return false
	}

	if !reasonEquals(&a.Reason, &b.Reason) {
		return false
	}

	if a.Traces != nil && b.Traces != nil {
		return tracesEquals(a.Traces, b.Traces)
	}

	return false
}

func tracesEquals(a, b []Trace) bool {
	count := 0
	lenA := len(a)
	lenB := len(b)

	total := 0
	if lenA < lenB {
		total = lenA
	} else {
		total = lenB
	}

	for ; count < total; count++ {
		if !traceEquals(&a[count], &b[count]) {
			break
		}
	}

	return count > total>>1
}

func traceEquals(a, b *Trace) bool {
	if a == nil {
		if b == nil {
			return true
		} else {
			return false
		}
	}

	if b == nil {
		return false
	}

	return a.Function == b.Function && a.File == b.File
}

func reasonEquals(a, b *Reason) bool {
	if a == nil {
		if b == nil {
			return true
		} else {
			return false
		}
	}

	if b == nil {
		return false
	}

	return a.JavaClass == b.JavaClass
}
