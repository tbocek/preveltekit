package encode

import p "github.com/tbocek/preveltekit/v2"

type Req struct {
	IDs []int `js:"ids"`
}

type Resp struct {
	OK bool `js:"ok"`
}

func Send() (Resp, error) {
	return p.Fetch[Resp]("/send", &p.FetchOptions{Method: "POST", Body: Req{}})
}
