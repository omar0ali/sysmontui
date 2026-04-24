package parts

type SearchController interface {
	ClosingSearchWith(search string)
	SetLoading()
}

type KillController interface {
	CloseingKillUI()
	SetLoading()
}
