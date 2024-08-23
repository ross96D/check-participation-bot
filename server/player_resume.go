package server

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ross96D/battle-log-parser/parser"
)

func sort(m map[parser.User]PlayerResume) []parser.User {
	type Z struct {
		k parser.User
		v PlayerResume
	}
	result := make([]Z, 0)
	for k, v := range m {
		result = append(result, Z{k: k, v: v})
	}
	slices.SortFunc(result, func(a Z, b Z) int {
		return b.v.Damage - a.v.Damage
	})
	ss := make([]parser.User, 0, len(result))
	for _, v := range result {
		ss = append(ss, v.k)
	}
	return ss
}

type AllPlayerResume []PlayerResume

func (alp AllPlayerResume) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("   %s\t%s\t%s\t%s\tcrits\n",
		FixedLenStr("Name", 13), FixedLenStr("dmg", 5), FixedLenStr("rcvd", 5), FixedLenStr("hits/total", 10)))
	for _, v := range alp {
		b.WriteString(v.StringSimple())
		b.WriteByte('\n')
	}
	return b.String()
}

type PlayerResume struct {
	Team    parser.Team
	Damage  int
	Tanqued int
	Miss    int
	Hits    int
	Crits   int
	Name    string
}

func (pr PlayerResume) String() string {
	b := strings.Builder{}
	b.WriteString(pr.Team.String())
	b.WriteString(fmt.Sprintf(
		" %s\tdmg: %s\trecieved: %d\tHits/Total: %d/%d %.1f%%\tcrits: %d",
		pr.NameWithFixedWidth(13), FixedLenStr(strconv.FormatInt(int64(pr.Damage), 10), 5), pr.Tanqued, pr.Hits, pr.Hits+pr.Miss, 100*float64(pr.Hits)/float64(pr.Hits+pr.Miss), pr.Crits),
	)
	return b.String()
}

func (pr PlayerResume) StringSimple() string {
	b := strings.Builder{}
	b.WriteString(pr.Team.String())
	b.WriteString(fmt.Sprintf(
		" %s\t%s\t%s\t%d/%d %.1f%%\t%d",
		pr.NameWithFixedWidth(13), FixedLenStr(strconv.Itoa(pr.Damage), 5),
		FixedLenStr(strconv.Itoa(pr.Tanqued), 5), pr.Hits, pr.Hits+pr.Miss,
		100*float64(pr.Hits)/float64(pr.Hits+pr.Miss), pr.Crits),
	)
	return b.String()
}

func (pr PlayerResume) NameWithFixedWidth(width uint) string {
	return FixedLenStr(pr.Name, width)
}

func (pr PlayerResume) Add(other PlayerResume) PlayerResume {
	return PlayerResume{
		Team:    pr.Team,
		Name:    pr.Name,
		Damage:  pr.Damage + other.Damage,
		Tanqued: pr.Tanqued + other.Tanqued,
		Hits:    pr.Hits + other.Hits,
		Miss:    pr.Miss + other.Miss,
		Crits:   pr.Crits + other.Crits,
	}
}

func PlayerResumen(b parser.Battle) map[parser.User]PlayerResume {
	empty := PlayerResume{}

	result := make(map[parser.User]PlayerResume, 0)
	for _, turn := range b.Turns {
		r := result[turn.Attacker]
		new := PlayerResume{
			Damage: turn.Damage(),
			Miss:   turn.Misses(),
			Hits:   turn.Hits(),
			Crits:  turn.Crits(),
		}
		if r == empty {
			r.Name = turn.Attacker.Name
			r.Team = turn.Attacker.Team
		}

		result[turn.Attacker] = r.Add(new)

		if !turn.Target.IsMiss() {
			r = result[turn.Target]
			if r == empty {
				r.Name = turn.Target.Name
				r.Team = turn.Target.Team
			}
			result[turn.Target] = r.Add(PlayerResume{Tanqued: turn.Damage(), Team: turn.Target.Team})
		}
	}
	return result
}

func FixedLenStr(str string, width uint) string {
	strB := []byte(str)
	nameLen := utf8.RuneCount(strB)
	if nameLen > int(width) {
		result := make([]byte, 0, width)
		count := uint(0)
		for _, r := range str {
			if count == width {
				break
			}
			result = utf8.AppendRune(result, r)
			count++
		}
		return string(result)
	} else {
		for i := uint(0); i < (width - uint(nameLen)); i++ {
			strB = append(strB, ' ')
		}
		return string(strB)
	}
}
