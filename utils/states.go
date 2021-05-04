package utils

type Messages struct {
	Type string
	Body string
}

type AccountPageVars struct {
	Loggedin bool
	Userinfo struct {
		Firstname  string
		Secondname string
		Storage    float64
	}
	Messages []Messages
}

type AlbumPageVars struct {
	Albums []struct {
		Name        string
		LatestPhoto string
	}
	Messages []Messages
}

type PhotoPageVars struct {
	Photos []struct {
		Photo string
	}
	Messages []Messages
}

type RegPageVars struct {
	Messages []Messages
}

type EqualAlbumVars struct {
	Photos []struct {
		Photo string
	}
	Messages []Messages
}
