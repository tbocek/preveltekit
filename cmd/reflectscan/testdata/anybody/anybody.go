package anybody

import p "github.com/tbocek/preveltekit/v2"

type Resp struct {
	OK bool `js:"ok"`
}

func Send(body any) (Resp, error) { return p.Post[Resp]("/send", body) }
