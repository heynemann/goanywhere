package web

import (
    "time"
    "net/http"
    "html/template"
    "labix.org/v2/mgo"
    ua "github.com/mssola/user_agent"
)

type URL struct {
    Date time.Time
    Uri string
    Mobile bool
    Bot bool
    Mozilla string

    Platform string
    Os string

    EngineName string
    EngineVersion string

    BrowserName string
    BrowserVersion string
}

func Index(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, req *http.Request) {
        c := session.DB("goanywhere").C("urls")

        results := []URL{}
        err := c.Find(nil).Limit(100).Iter().All(&results)
        if err != nil {
            panic(err)
        }

        tmpl, err := template.New("urls").Parse(`
        <html>
            <head>
            </head>
            <body>
                <ul>
                    {{ range . }}
                    <li>
                        <ul>
                            <li>Date: {{ .Date }}</li>
                            <li>Uri: {{ .Uri }}</li>
                            <li>Mobile: {{ .Mobile }}</li>
                            <li>Bot: {{ .Bot }}</li>
                            <li>Mozilla: {{ .Mozilla }}</li>

                            <li>Platform: {{ .Platform }}</li>
                            <li>Os: {{ .Os }}</li>

                            <li>Engine Name: {{ .EngineName }}</li>
                            <li>Engine Version: {{ .EngineVersion }}</li>

                            <li>Browser Name: {{ .BrowserName }}</li>
                            <li>Browser Version: {{ .BrowserVersion }}</li>
                        </ul>
                    </li>
                    {{ end }}
                </ul>
            </body>
        </html>
        `)

        err = tmpl.Execute(w, results)
    }
}

func Router(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, req *http.Request) {
        user_agent := new(ua.UserAgent);
        user_agent.Parse(req.Header.Get("User-Agent"));

        url := req.FormValue("url")

        engine_name, engine_version := user_agent.Engine();
        browser_name, browser_version := user_agent.Browser();

        url_document := &URL{
            Date: time.Now(),
            Uri: url,
            Mobile: user_agent.Mobile(),
            Bot: user_agent.Bot(),
            Mozilla: user_agent.Mozilla(),
            Platform: user_agent.Platform(),
            Os: user_agent.OS(),
            EngineName: engine_name,
            EngineVersion: engine_version,
            BrowserName: browser_name,
            BrowserVersion: browser_version,
        }

        c := session.DB("goanywhere").C("urls")

        err := c.Insert(url_document)

        if err != nil {
            panic(err)
        }

        http.Redirect(w, req, url, http.StatusFound)
    }
}
