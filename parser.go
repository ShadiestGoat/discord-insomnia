package main

import (
	"regexp"
	"strings"
)

var (
	RegNewLine = regexp.MustCompile(`\n`)
)

func ParseDoc(inp []byte) []Request {
	out := []Request{}

	for RegEntryStart.Match(inp) {
		locs := RegEntryStart.FindAllIndex(inp, 2)
		curOut := ParseTitle(inp[:locs[0][1]-1])
		if len(locs) > 1 {
			curOut.DocContent = string(inp[locs[0][1]:locs[1][0]])
			inp = inp[locs[1][0]:]
		} else {
			curOut.DocContent = string(inp[locs[0][1]+1:len(inp)-1])
			inp = []byte{}
		}
		curOut.DocContent = ParseDocumentation(curOut.DocContent)
		out = append(out, curOut)
	}

	return out
}

var (
	RegURLParam = regexp.MustCompile(`\{.+#.+?\}`)
)

var URLParamMap = map[string]string{
	"channel.id": "channel",
	"guild.id": "guild",
	"user.id": "user",
	"webhook.id": "webhook",
	"template.code": "guildTemplate",
	"invite.code": "invite",
	"sticker.id": "sticker",
}

func ParseTitle(headingLine []byte) Request {
	parsed := Request{}

	ret := strings.SplitN(string(headingLine), " % ", 2)

	parsed.Name = ret[0][3:]

	urlRawRaw := strings.Split(ret[1], " ")
	
	parsed.Method = urlRawRaw[0]

	urlRaw := urlRawRaw[1]

	for RegURLParam.MatchString(urlRaw) {
		loc := RegURLParam.FindStringIndex(urlRaw)
		
		paramRaw := urlRaw[loc[0]:loc[1]]
		paramFound := false

		
		for inpP, outP := range URLParamMap {
			if len(paramRaw) < len(inpP) + 1 {
				continue
			}

			if paramRaw[1:len(inpP)+1] == inpP {
				paramFound = true
				urlRaw = urlRaw[:loc[0]] + "{{" + outP + "}}" + urlRaw[loc[1]:]
			}
		}
		if !paramFound {
			panic("Unknown param! " + paramRaw)
		}
	}

	parsed.URI = urlRaw
	
	return parsed
}

var (
	RegDoc = regexp.MustCompile(`\]\(#.+?\)`)
	RegSlash = regexp.MustCompile(`/`)
)

func ParseDocumentation(inp string) string {
	for RegDoc.MatchString(inp) {
		loc := RegDoc.FindStringIndex(inp)
		uriRaw := inp[loc[0]+3:loc[1]-1]
		if RegSlash.MatchString(uriRaw) {
			uriRawLoc := RegSlash.FindStringIndex(uriRaw)
			uriRaw = DocURL(uriRaw[:uriRawLoc[0]]) + "#" + uriRaw[uriRawLoc[1]:]
		} else {
			uriRaw = DocURL(uriRaw)
		}
		inp = inp[:loc[0]+2] + uriRaw + inp[loc[1]-1:]
	}
	return inp
}

func DocURL(inp string) string {
	spl := strings.SplitN(inp, "_", 3)
	if spl[1] == "RESOURCES" {
		return DOC_BASE + "resources/" + URLify(spl[2])
	} else if spl[1] == "INTERACTIONS" {
		return DOC_BASE + "interactions/" + URLify(spl[2])
	} else if spl[1] == "TOPICS" {
		return DOC_BASE + "topics/" + URLify(spl[2])
	} else {
		switch inp {
		case "DOCS_REFERENCE":
			return DOC_BASE + "reference"
		default:
			panic(inp)
		}
	}
}