package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
    DOC_BASE = "https://discord.com/developers/docs/"
)

// use var because you can't adress a constant

var (
    EXP_ENV_ID = "env_2d186dc23870ea858f45ef9c385bd421520ccd53"
    EXP_COLLECTION_ID = "wrk_ef4b07e4c2704be891276e79da00e870"
    EXP_HEADER_AUTH_ID = "pair_c81cb238e88e4e66bd14349071020eb7"
)

type MainExport struct {
	Type string `json:"_type"`
    Format int `json:"__export_format"`
	// "2022-04-10T13:32:40.365Z",
	Date string `json:"__export_date"`
	// "__export_source": "insomnia.desktop.app:v2022.2.1",
	Source string `json:"__export_source"`
    Resources []ResourceExport `json:"resources"`
}

type HeaderExport struct {
    Name string `json:"name"`
    Value string `json:"value"`
    Desc string `json:"description"`
    ID string `json:"id"`
}

type ResourceExport struct {
    ID string `json:"_id"`
    ParentID *string `json:"parentId"`
    URL string `json:"url,omitempty"`
    Name string `json:"name"`
    Desc string `json:"description"`
    Scope string `json:"scope,omitempty"`
    Type string `json:"_type"`
    Method string `json:"method,omitempty"`
    Headers []HeaderExport `json:"headers,omitempty"`
    EnvData map[string]string `json:"data,omitempty"`
    EnvOrder map[string][]string `json:"dataPropertyOrder,omitempty"`
}


var IDs = map[string]bool{
    EXP_ENV_ID[4:]: true,
    EXP_COLLECTION_ID[4:]: true,
    EXP_HEADER_AUTH_ID[5:]: true,
}

func MakeID() string {
    v := ""
    for {
        v = uuid.New().String()
        v = strings.Join(strings.Split(v, "-"), "")
        if _, ok := IDs[v]; !ok {
            break
        }
    }
    return v
}

var RequestHeaders = []HeaderExport{
    {
        Name: "Authorization",
        Value: "{{auth}}",
        Desc: "",
        ID: EXP_HEADER_AUTH_ID,
    },
}

func GenerateExport(allRequests []RequestGroup) MainExport {
    exp := MainExport{}

    exp.Type = "export"
    exp.Format = 4
    exp.Date = time.Now().Format(time.RFC3339)
    exp.Source = "insomnia.desktop.app:v2022.2.1"

    exp.Resources = []ResourceExport{
        {
        	ID:       EXP_COLLECTION_ID,
        	ParentID: nil,
        	Name:     "Discord REST API",
        	Desc:     "The Discord REST API Implemented into insomnia, by [Shady Goat](https://shadygoat.eu)",
        	Scope:    "collection",
        	Type:     "workspace",
        },
        {
            ID: EXP_ENV_ID,
            ParentID: &EXP_COLLECTION_ID,
            Name: "Base Environment",
            Type: "environment",
            EnvData: map[string]string{
                "base": "https://canary.discord.com/api/v9",
                "auth": "Bot EG_TOKEN_HERE",
                "channel": "EG_CHANNEL_ID",
                "guild": "EG_GUILD_ID",
                "user": "EG_USER_ID",
                "webhook": "EG_WEBHOOK_ID",
                "guildTemplate": "EG_GUILD_TEMPLATE_CODE",
                "invite": "EG_INVITE_CODE",
                "sticker": "EG_STICKER_ID",
            },
        },
    }

    for _, group := range allRequests {
        groupID := "fld_" + MakeID()

        exp.Resources = append(exp.Resources, ResourceExport{
        	ID:       groupID,
        	ParentID: &EXP_COLLECTION_ID,
        	Name:     group.Name,
        	Type:     "request_group",
        })

        for _, req := range group.Requests {
            exp.Resources = append(exp.Resources, ResourceExport{
            	ID:       "req_" + MakeID(),
            	ParentID: &groupID,
            	URL:      "{{base}}" + req.URI,
            	Name:     req.Name,
            	Desc:     fmt.Sprintf("[Discord Documentation here](%v)\n\n", req.DocURL(group.Name)) + req.DocContent,
            	Type:     "request",
            	Method:   req.Method,
            	Headers:  RequestHeaders,
            })
        }
    }

    return exp
}

type Request struct {
	URI string
	Name string
	Method string
	DocContent string
}

var (
    RegURLifySpace = regexp.MustCompile(` |_`)
)

func URLify(inp string) string {
    return strings.ToLower(RegURLifySpace.ReplaceAllString(inp, "-"))
}

func (req Request) DocURL(parent string) string {
    return DOC_BASE + "resources/" + URLify(parent) + "#" + URLify(req.Name)
}

type RequestGroup struct {
	Requests []Request
	Name string
}
