package filesystem

func (s Service) invalidateSearch() {
	s.searchIndex.Clear()
}
