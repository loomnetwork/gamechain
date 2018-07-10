package battleground

func OwnerKey(owner string) []byte {
	return []byte("owner:" + owner)
}
