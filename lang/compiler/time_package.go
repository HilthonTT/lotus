package compiler

import (
	"fmt"
	"time"

	"github.com/hilthontt/lotus/object"
)

func timePackage() *object.Package {
	return &object.Package{
		Name: "Time",
		Functions: map[string]object.PackageFunction{

			// Time.now() -> int
			// Returns current time as Unix milliseconds.
			"now": func(args ...object.Object) object.Object {
				return &object.Integer{Value: time.Now().UnixMilli()}
			},

			// Time.sleep(ms: int)
			// Pauses execution for ms milliseconds.
			"sleep": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				time.Sleep(time.Duration(ms.Value) * time.Millisecond)
				return &object.Nil{}
			},

			// Time.format(ms: int, layout: string) -> string
			// Formats a Unix ms timestamp using a Go time layout.
			// Example layouts:
			//   "2006-01-02"           → 2024-04-15
			//   "2006-01-02 15:04:05"  → 2024-04-15 13:45:00
			//   "15:04:05"             → 13:45:00
			//   "Mon, 02 Jan 2006"     → Tue, 15 Apr 2024
			"format": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				ms, ok1 := args[0].(*object.Integer)
				layout, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value)
				return &object.String{Value: t.Format(layout.Value)}
			},

			// Time.parse(str: string, layout: string) -> int
			// Parses a time string into Unix milliseconds.
			// Returns nil on parse error.
			"parse": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				str, ok1 := args[0].(*object.String)
				layout, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t, err := time.Parse(layout.Value, str.Value)
				if err != nil {
					return &object.Nil{}
				}
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.since(ms: int) -> int
			// Returns milliseconds elapsed since the given timestamp.
			"since": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				elapsed := time.Since(time.UnixMilli(ms.Value))
				return &object.Integer{Value: elapsed.Milliseconds()}
			},

			// Time.until(ms: int) -> int
			// Returns milliseconds until the given timestamp.
			"until": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				until := time.Until(time.UnixMilli(ms.Value))
				return &object.Integer{Value: until.Milliseconds()}
			},

			// Time.add(ms: int, duration: int) -> int
			// Adds duration milliseconds to a timestamp.
			"add": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				ms, ok1 := args[0].(*object.Integer)
				dur, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value).Add(time.Duration(dur.Value) * time.Millisecond)
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.diff(a: int, b: int) -> int
			// Returns (a - b) in milliseconds.
			"diff": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				a, ok1 := args[0].(*object.Integer)
				b, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				return &object.Integer{Value: a.Value - b.Value}
			},

			// Time.year(ms: int) -> int
			"year": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Year())}
			},

			// Time.month(ms: int) -> int  (1-12)
			"month": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Month())}
			},

			// Time.day(ms: int) -> int  (1-31)
			"day": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Day())}
			},

			// Time.hour(ms: int) -> int  (0-23)
			"hour": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Hour())}
			},

			// Time.minute(ms: int) -> int  (0-59)
			"minute": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Minute())}
			},

			// Time.second(ms: int) -> int  (0-59)
			"second": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(time.UnixMilli(ms.Value).Second())}
			},

			// Time.weekday(ms: int) -> string
			// Returns the weekday name: "Monday", "Tuesday", etc.
			"weekday": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.String{Value: time.UnixMilli(ms.Value).Weekday().String()}
			},

			// Time.unix(ms: int) -> int
			// Converts Unix milliseconds to Unix seconds.
			"unix": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: ms.Value / 1000}
			},

			// Time.fromUnix(sec: int) -> int
			// Converts Unix seconds to Unix milliseconds.
			"fromUnix": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				sec, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: sec.Value * 1000}
			},

			// Time.isBefore(a: int, b: int) -> bool
			// Returns true if timestamp a is before b.
			"isBefore": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				a, ok1 := args[0].(*object.Integer)
				b, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				return &object.Boolean{Value: a.Value < b.Value}
			},

			// Time.isAfter(a: int, b: int) -> bool
			// Returns true if timestamp a is after b.
			"isAfter": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				a, ok1 := args[0].(*object.Integer)
				b, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				return &object.Boolean{Value: a.Value > b.Value}
			},

			// Time.startOfDay(ms: int) -> int
			// Returns timestamp of midnight (00:00:00) for the given day.
			"startOfDay": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value)
				start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				return &object.Integer{Value: start.UnixMilli()}
			},

			// Time.endOfDay(ms: int) -> int
			// Returns timestamp of 23:59:59.999 for the given day.
			"endOfDay": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value)
				end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999_000_000, t.Location())
				return &object.Integer{Value: end.UnixMilli()}
			},

			// Time.addDays(ms: int, days: int) -> int
			// Adds n calendar days to a timestamp.
			"addDays": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				ms, ok1 := args[0].(*object.Integer)
				days, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value).AddDate(0, 0, int(days.Value))
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.addMonths(ms: int, months: int) -> int
			// Adds n calendar months to a timestamp.
			"addMonths": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				ms, ok1 := args[0].(*object.Integer)
				months, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value).AddDate(0, int(months.Value), 0)
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.addYears(ms: int, years: int) -> int
			// Adds n calendar years to a timestamp.
			"addYears": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				ms, ok1 := args[0].(*object.Integer)
				years, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value).AddDate(int(years.Value), 0, 0)
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.duration(ms: int) -> string
			// Returns a human-readable duration string for elapsed ms.
			// e.g. "2h 35m 10s" or "45ms"
			"duration": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				total := ms.Value
				if total < 0 {
					total = -total
				}
				hours := total / 3_600_000
				mins := (total % 3_600_000) / 60_000
				secs := (total % 60_000) / 1_000
				millis := total % 1_000

				var result string
				if hours > 0 {
					result = fmt.Sprintf("%dh %dm %ds", hours, mins, secs)
				} else if mins > 0 {
					result = fmt.Sprintf("%dm %ds", mins, secs)
				} else if secs > 0 {
					result = fmt.Sprintf("%ds %dms", secs, millis)
				} else {
					result = fmt.Sprintf("%dms", millis)
				}
				return &object.String{Value: result}
			},

			// Time.ms(n: int) -> int  — n milliseconds
			"ms": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				n, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: n.Value}
			},

			// Time.seconds(n: int) -> int  — n seconds in ms
			"seconds": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				n, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: n.Value * 1_000}
			},

			// Time.minutes(n: int) -> int  — n minutes in ms
			"minutes": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				n, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: n.Value * 60_000}
			},

			// Time.hours(n: int) -> int  — n hours in ms
			"hours": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				n, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: n.Value * 3_600_000}
			},

			// Time.days(n: int) -> int  — n days in ms
			"days": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				n, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				return &object.Integer{Value: n.Value * 86_400_000}
			},

			// Time.utc(ms: int) -> int
			// Converts a local timestamp to UTC.
			"utc": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Nil{}
				}
				t := time.UnixMilli(ms.Value).UTC()
				return &object.Integer{Value: t.UnixMilli()}
			},

			// Time.timezone() -> string
			// Returns the local timezone name, e.g. "UTC", "Europe/Berlin".
			"timezone": func(args ...object.Object) object.Object {
				name, _ := time.Now().Zone()
				return &object.String{Value: name}
			},
		},
	}
}
