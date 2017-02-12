package file

/*
// SessionBackend sb
type SessionBackend struct {
	dir string
}

func init() {
	gob.Register(Session{})
}

// NewSessionBackend nsb
func NewSessionBackend(dsn string) backend.SessionBackender {
	cfg, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}
	return SessionBackend{
		dir: cfg.Path,
	}
}

// Init for context
func (sb SessionBackend) Init(c echo.Context) backend.Sessioner {
	cookie, err := c.Cookie("session")
	if err != nil {
		//log.Println(err)
	}
	var s Session

	if cookie == nil {
		s = Session{
			Data: make(map[string]interface{}),
			id:   strconv.FormatInt(rand.Int63(), 16),
			Dir:  sb.dir,
		}
	} else {
		f, err := os.Open(sb.dir + "/" + cookie.Value + ".sess")
		if err != nil {
			//log.Println(1, err)
			s = Session{
				Data: make(map[string]interface{}),
				id:   strconv.FormatInt(rand.Int63(), 16),
				Dir:  sb.dir,
			}
		} else {
			defer f.Close()

			gf := gob.NewDecoder(f)
			err = gf.Decode(&s)
			if err != nil {
				//log.Println(2, err)
				s = Session{
					Data: make(map[string]interface{}),
					id:   strconv.FormatInt(rand.Int63(), 16),
					Dir:  sb.dir,
				}
			} else {
				f.Close()
			}
		}
	}

	cookie = &http.Cookie{
		Value: s.id,
		Name:  "session",
		Path:  "/",
	}
	c.SetCookie(cookie)

	return &s
}

// Session s
type Session struct {
	id   string
	Data map[string]interface{}
	Dir  string
}

// Get g
func (s Session) Get(name string) interface{} {
	return s.Data[name]
}

// Set s
func (s *Session) Set(name string, data interface{}) {
	s.Data[name] = data
}

// Persist p
func (s Session) Persist() bool {
	f, err := os.OpenFile(s.Dir+"/"+s.id+".sess", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gf := gob.NewEncoder(f)
	err = gf.Encode(s)
	if err != nil {
		panic(err)
	}

	return true
}

// ID id
func (s Session) ID() string {
	return s.id
}
*/
