package webserver

type indexData struct {
	LayoutData layoutData
}

type aramData struct {
	LayoutData layoutData
}

type playersData struct {
	LayoutData layoutData
}

type layoutData struct {
	NavData    navData
	SubnavData subnavData
}

type navData struct {
}

type subnavData struct {
	Title string
}
